package users

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gorm.io/gorm/clause"
)

// A tokenized asset on Base has two distinct addresses:
//
//   - its token contract (TokenizedAsset.ContractAddress): the deployed B20
//     token itself - the asset's on-chain identity, holding every balance;
//   - its issuing wallet (TokenizedAsset.IssuingWalletAddress): a Safe
//     multisig that owns the token contract, is the only account allowed to
//     mint it, and also holds the asset's unsold supply (its treasury).
//
// The issuing Safe is deployed here, per asset, owned by the configured
// TOKENIZATION_ISSUING_SAFE_SIGNERS. The token contract is deployed by
// operations with the Safe as its owner/minter and then registered with
// RegisterTokenizedAssetContract, which verifies that on-chain.

// Canonical Safe v1.4.1 deployments (identical on Base and Base Sepolia).
// Each is verified to have contract code before a Safe is deployed from it.
const (
	defaultSafeProxyFactoryAddress    = "0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67"
	defaultSafeSingletonAddress       = "0x29fcB43b46531BcA003ddC8FCB67FFE91900C762" // SafeL2
	defaultSafeFallbackHandlerAddress = "0xfd0732Dc9E303f09fCEf3a7388Ad10A83459Ec99" // CompatibilityFallbackHandler
)

// issuingProfileKeyConfigured reports whether TOKENIZATION_ISSUING_PROFILE_WALLET
// holds a valid private key: the issuing profile's primary signer key, which
// signs the minting approvers' shared-access setup on each issuing Safe.
func issuingProfileKeyConfigured() bool {
	_, err := evmkeypair.ParseFull(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
	return err == nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func safeDeployConfig() gnosissafe.DeployConfig {
	return gnosissafe.DeployConfig{
		ProxyFactory:    common.HexToAddress(envOrDefault("SAFE_PROXY_FACTORY_ADDRESS", defaultSafeProxyFactoryAddress)),
		Singleton:       common.HexToAddress(envOrDefault("SAFE_SINGLETON_ADDRESS", defaultSafeSingletonAddress)),
		FallbackHandler: common.HexToAddress(envOrDefault("SAFE_FALLBACK_HANDLER_ADDRESS", defaultSafeFallbackHandlerAddress)),
	}
}

func multiSendCallOnlyAddress() common.Address {
	return common.HexToAddress(envOrDefault("SAFE_MULTISEND_CALL_ONLY_ADDRESS", gnosissafe.DefaultMultiSendCallOnlyAddress))
}

// issuingSafeSigners parses TOKENIZATION_ISSUING_SAFE_SIGNERS: the private
// keys of the issuing Safes' owners, ";"-separated ("," also accepted).
// The first key also broadcasts (and pays gas for) Safe deployments and
// mint transactions, so it must hold native balance.
func issuingSafeSigners() ([]*evmkeypair.Full, error) {
	raw := strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_SAFE_SIGNERS"))
	if raw == "" {
		return nil, fmt.Errorf("TOKENIZATION_ISSUING_SAFE_SIGNERS is not configured")
	}
	sep := ";"
	if !strings.Contains(raw, ";") {
		sep = ","
	}
	signers := make([]*evmkeypair.Full, 0)
	for i, part := range strings.Split(raw, sep) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kp, err := evmkeypair.ParseFull(part)
		if err != nil {
			return nil, fmt.Errorf("TOKENIZATION_ISSUING_SAFE_SIGNERS entry %d is invalid: %w", i, err)
		}
		signers = append(signers, kp)
	}
	if len(signers) == 0 {
		return nil, fmt.Errorf("TOKENIZATION_ISSUING_SAFE_SIGNERS has no keys")
	}
	return signers, nil
}

// issuingSafeThreshold is TOKENIZATION_ISSUING_SAFE_THRESHOLD, defaulting
// to 3 (or every owner, when fewer than 3 are configured).
func issuingSafeThreshold(owners int) (int64, error) {
	def := int64(3)
	if int64(owners) < def {
		def = int64(owners)
	}
	raw := strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_SAFE_THRESHOLD"))
	if raw == "" {
		return def, nil
	}
	t, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || t < 1 || t > int64(owners) {
		return 0, fmt.Errorf("TOKENIZATION_ISSUING_SAFE_THRESHOLD must be between 1 and the %d configured signer(s)", owners)
	}
	return t, nil
}

// issuingSafeSaltNonce makes each asset's Safe address deterministic, so a
// retried assignment reuses the Safe an earlier attempt deployed.
func issuingSafeSaltNonce(tokenizationID string) *big.Int {
	return new(big.Int).SetBytes(crypto.Keccak256([]byte("trovo-tokenization-issuing-safe:" + tokenizationID)))
}

// deployIssuingSafe deploys (or finds, on retry) the tokenized asset's
// issuing Safe and records it as a wallet of the tokenization issuing
// profile, so the existing minting-approver permissions and approval flow
// attach to it like any other wallet.
func deployIssuingSafe(issuerProfile *userModels.User, ato *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (userModels.UserWallet, error) {
	signers, err := issuingSafeSigners()
	if err != nil {
		log.Printf("[deployIssuingSafe] %v\n", err)
		return userModels.UserWallet{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-issuing-safe-signers-not-set", ErrMessage: "Issuing Safe signers are not configured."}
	}
	threshold, err := issuingSafeThreshold(len(signers))
	if err != nil {
		log.Printf("[deployIssuingSafe] %v\n", err)
		return userModels.UserWallet{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-issuing-safe-threshold-invalid", ErrMessage: err.Error()}
	}
	owners := make([]common.Address, 0, len(signers))
	for _, s := range signers {
		owners = append(owners, common.HexToAddress(s.Address()))
	}

	safe, err := gnosissafe.DeploySafe(context.Background(), gc.BantuExpansionClient, network.GetBlockchainChainID(), signers[0], safeDeployConfig(), owners, threshold, issuingSafeSaltNonce(ato.ID))
	if err != nil {
		log.Printf("[deployIssuingSafe] deploying issuing Safe for tokenization %v: %v\n", ato.ID, err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[deployIssuingSafe] deploying issuing Safe for tokenization %v failed: %v", ato.ID, err))
		return userModels.UserWallet{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-issuing-safe-deployment-failed", ErrMessage: "Could not deploy the issuing Safe for this asset. Please try again later."}
	}
	safeAddress := safe.Hex()
	log.Printf("[deployIssuingSafe] issuing Safe for tokenization %v (%v): %v, %d-of-%d\n", ato.ID, *ato.AssetCode, safeAddress, threshold, len(owners))

	// A retry after a failed DB write finds the Safe already recorded.
	if existing, e := userModels.UserWalletID(safeAddress).GetWallet(gc.DB, gc); e == nil {
		return existing, nil
	}

	row, err := issuerProfile.BuildNewSubWallet(userModels.SafeDeployment{Address: safeAddress}, *ato.AssetCode+"issuer", fmt.Sprintf("Issuing Safe for %v", *ato.AssetCode), 1, "", gc)
	if err != nil {
		log.Printf("[deployIssuingSafe] building wallet record for issuing Safe %v: %v\n", safeAddress, err)
		return userModels.UserWallet{}, err
	}
	if e := gc.DB.Omit(clause.Associations).Create(&row).Error; e != nil {
		log.Printf("[deployIssuingSafe] saving wallet record for issuing Safe %v: %v\n", safeAddress, e)
		return userModels.UserWallet{}, &tErrors.ErrorTemporaryServerError{}
	}
	if e := gc.RoachDB.Create(&userModels.TrackedAddress{Address: safeAddress}).Error; e != nil {
		// payment-history-engine also discovers untracked wallets on its own
		log.Printf("[deployIssuingSafe] tracking issuing Safe %v for payment history: %v\n", safeAddress, e)
	}
	issuerProfile.InvalidateUserCache(gc)

	return userModels.UserWalletID(safeAddress).GetWallet(gc.DB, gc)
}

// tokenizedAssetContract returns the asset's registered B20 token contract
// - the address every transfer, balance, authorization and listing of the
// asset must use - or an error when it has not been registered yet.
func tokenizedAssetContract(ta *userModels.TokenizedAsset) (string, error) {
	if ta.ContractAddress == nil || !common.IsHexAddress(*ta.ContractAddress) {
		return "", &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-token-contract-not-registered", ErrMessage: "This tokenized asset's token contract has not been registered yet."}
	}
	return *ta.ContractAddress, nil
}
