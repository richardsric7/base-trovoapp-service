package users

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gorm.io/gorm/clause"
)

// tokenizedAssetTokenABI is the subset of a tokenized asset's B20 token
// contract this backend relies on: ERC-20 metadata/supply, a
// mint(address,uint256) callable by the issuing Safe, and either
// OpenZeppelin AccessControl (MINTER_ROLE) or Ownable to prove it.
var tokenizedAssetTokenABI abi.ABI

// minterRole is OpenZeppelin's conventional keccak256("MINTER_ROLE").
var minterRole = crypto.Keccak256Hash([]byte("MINTER_ROLE"))

func init() {
	var err error
	tokenizedAssetTokenABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},
		{"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"},
		{"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
		{"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"name":"role","type":"bytes32"},{"name":"account","type":"address"}],"name":"hasRole","outputs":[{"name":"","type":"bool"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"}
	]`))
	if err != nil {
		panic("users: invalid tokenized asset token ABI: " + err.Error())
	}
}

// contractReader is the read-only chain access token verification needs
// (satisfied by *ethclient.Client; faked in tests).
type contractReader interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

func readToken(ctx context.Context, r contractReader, token common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data, err := tokenizedAssetTokenABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	out, err := r.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	res, err := tokenizedAssetTokenABI.Unpack(method, out)
	if err != nil || len(res) == 0 {
		return nil, fmt.Errorf("could not decode %s() result", method)
	}
	return res, nil
}

func contractError(msg string) error {
	return &tErrors.CustomError{Param: "contractAddress", Err: "error-invalid-token-contract", ErrMessage: msg, Code: 400}
}

// verifyTokenizedAssetContract checks, on-chain, that token is a B20 token
// contract the platform can issue assetCode through: it has code, its
// symbol is the asset code, it reports decimals, nothing has been minted
// yet (so the platform's mint is the whole supply), and the issuing Safe
// holds MINTER_ROLE or is the contract owner.
func verifyTokenizedAssetContract(ctx context.Context, r contractReader, token, issuingSafe common.Address, assetCode string) error {
	code, err := r.CodeAt(ctx, token, nil)
	if err != nil {
		return err
	}
	if len(code) == 0 {
		return contractError(fmt.Sprintf("No contract is deployed at %s on this network.", token.Hex()))
	}
	safeCode, err := r.CodeAt(ctx, issuingSafe, nil)
	if err != nil {
		return err
	}
	if len(safeCode) == 0 {
		return contractError(fmt.Sprintf("The issuing wallet %s is not a deployed Safe.", issuingSafe.Hex()))
	}

	res, err := readToken(ctx, r, token, "symbol")
	if err != nil {
		return contractError("The contract does not implement the B20 symbol() function.")
	}
	if symbol, _ := res[0].(string); !strings.EqualFold(strings.TrimSpace(symbol), assetCode) {
		return contractError(fmt.Sprintf("The contract's symbol %q does not match the asset code %q.", symbol, assetCode))
	}
	if _, err := readToken(ctx, r, token, "decimals"); err != nil {
		return contractError("The contract does not implement the B20 decimals() function.")
	}
	res, err = readToken(ctx, r, token, "totalSupply")
	if err != nil {
		return contractError("The contract does not implement the B20 totalSupply() function.")
	}
	if supply, _ := res[0].(*big.Int); supply == nil || supply.Sign() != 0 {
		return contractError("The contract already has tokens in circulation. Its supply must be minted by the issuing Safe through the minting approval flow.")
	}

	if res, err := readToken(ctx, r, token, "hasRole", [32]byte(minterRole), issuingSafe); err == nil {
		if ok, _ := res[0].(bool); ok {
			return nil
		}
	}
	if res, err := readToken(ctx, r, token, "owner"); err == nil {
		if owner, _ := res[0].(common.Address); owner == issuingSafe {
			return nil
		}
	}
	return contractError(fmt.Sprintf("The issuing Safe %s is neither granted MINTER_ROLE nor the owner of this contract, so it cannot mint the asset.", issuingSafe.Hex()))
}

// RegisterTokenizedAssetContract records the deployed B20 token contract of
// a tokenized asset that has not been minted yet, after verifying it
// on-chain against the asset's issuing Safe.
func RegisterTokenizedAssetContract(tokenizationID string, initiator *userModels.User, input *userModels.TokenizedAssetContractInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {
	if sErr := initiator.EnsureNotSuspended(); sErr != nil {
		return ato, sErr
	}
	if !IsTokenizationMintingApprover(initiator.Username, gc.DB) && !IsTokenizationMintingInitiator(initiator.Username, gc.DB) {
		return ato, &tErrors.CustomError{Param: "publicKey", Err: "error-access-denied", ErrMessage: fmt.Sprintf("%v does not have minting permission to perform this action.", initiator.Username), Code: 401}
	}
	raw := strings.TrimSpace(input.ContractAddress)
	if !common.IsHexAddress(raw) || !strings.HasPrefix(strings.ToLower(raw), "0x") {
		return ato, contractError("Contract address must be a 0x-prefixed, 20-byte hex address.")
	}
	token := common.HexToAddress(raw)

	ato, _, err = GetTokenizedAssetByID(tokenizationID, gc.DB)
	if err != nil {
		return ato, err
	}
	if ato.AssetTokenizationStatus > 3 {
		return ato, &tErrors.CustomError{Param: "contractAddress", Err: "error-status-too-high", ErrMessage: "The asset has already been minted; its token contract can no longer be changed.", Code: 400}
	}
	if ato.AssetCode == nil || ato.IssuingWalletAddress == nil {
		return ato, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-no-issuing-wallet", ErrMessage: "The asset has no issuing Safe yet. The token contract must be deployed with the issuing Safe as its owner/minter.", Code: 400}
	}

	var clashes int64
	gc.DB.Model(&userModels.TokenizedAsset{}).Where("lower(contract_address) = ? AND id <> ?", strings.ToLower(token.Hex()), ato.ID).Count(&clashes)
	if clashes > 0 {
		return ato, contractError("This contract is already registered for another tokenized asset.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = verifyTokenizedAssetContract(ctx, gc.BantuExpansionClient, token, common.HexToAddress(*ato.IssuingWalletAddress), *ato.AssetCode); err != nil {
		log.Printf("[RegisterTokenizedAssetContract] %v rejected contract %v for tokenization %v: %v\n", initiator.Username, token.Hex(), ato.ID, err)
		if _, ok := err.(tErrors.GenericError); !ok {
			err = &tErrors.ErrorTemporaryServerError{}
		}
		return ato, err
	}

	contract := token.Hex()
	ato.ContractAddress = &contract
	ato.LastUpdatedBy = &initiator.Username
	if e := gc.DB.Omit(clause.Associations).Save(&ato).Error; e != nil {
		log.Printf("[RegisterTokenizedAssetContract] saving contract %v for tokenization %v: %v\n", contract, ato.ID, e)
		return ato, &tErrors.ErrorTemporaryServerError{}
	}
	log.Printf("[RegisterTokenizedAssetContract] %v registered token contract %v for %v (%v)\n", initiator.Username, contract, *ato.AssetCode, ato.ID)
	ato, _, err = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, err
}
