package users

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// A tokenized asset on Base has three addresses:
//
//   - its token contract (TokenizedAsset.ContractAddress): the deployed
//     TokenizedAsset token, owned by the issuing Safe - the only account
//     that can mint it;
//   - its issuing wallet (IssuingWalletAddress): a Safe sub-wallet of the
//     tokenization issuing profile (TOKENIZATION_ISSUING_PROFILE), owned
//     by the profile's key and the asset's minting approvers, with the
//     minting approval threshold;
//   - its distribution wallet (MarketMakingWallet,
//     WalletToHoldAssetsNotForSale): the issuing wallet's linked Safe, with
//     the same owners and the issuing Safe as a module. It receives the
//     minted supply and sells it through TrovoOfferBook.
//
// Both Safes are counterfactual: their addresses are fixed from their
// owners when the issuing wallet is assigned, and the mint operation - the
// issuing Safe's first, signed by the minting approvers - deploys them. No
// platform key owns or signs for either.

// issuingProfile loads the tokenization issuing profile. When
// TOKENIZATION_ISSUING_PROFILE_WALLET is set it must be that profile's
// primary wallet address (a sanity check of the configuration).
func issuingProfile(gc *sharedconfig.GlobalConfig) (userModels.User, error) {
	name := strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))
	if name == "" {
		return userModels.User{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-default-issuing-profile-not-set", ErrMessage: "Issuing profile not set."}
	}
	userModels.Username(name).InvalidateUserCache(gc)
	profile, err := userModels.Username(name).GetFullUser(gc.DB, gc)
	if err != nil || !common.IsHexAddress(profile.PrimarySigner) {
		return userModels.User{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-invallid-issuing-profile", ErrMessage: "Issuing profile not valid."}
	}
	if w := strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")); w != "" && !strings.EqualFold(w, profile.Address) {
		log.Printf("[issuingProfile] TOKENIZATION_ISSUING_PROFILE_WALLET %v is not %v's wallet %v", w, name, profile.Address)
		return userModels.User{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-invallid-issuing-profile", ErrMessage: "Issuing profile wallet is misconfigured."}
	}
	return profile, nil
}

// mintingUsers resolves a CSV of usernames.
func mintingUsers(csv string, gc *sharedconfig.GlobalConfig) ([]*userModels.User, error) {
	var out []*userModels.User
	seen := map[string]bool{}
	for _, name := range strings.Split(strings.NewReplacer("\r", "", "\n", "", " ", "").Replace(csv), ",") {
		if name == "" || seen[strings.ToLower(name)] {
			continue
		}
		seen[strings.ToLower(name)] = true
		u, err := userModels.Username(name).GetSimpleUser(gc.DB, gc)
		if err != nil {
			return nil, &tErrors.CustomError{Param: "mintingApprovers", Err: "error-invalid-minting-approver", ErrMessage: fmt.Sprintf("%v is not a valid username.", name)}
		}
		out = append(out, &u)
	}
	return out, nil
}

// mintingApprovalsNeeded is how many minting approvers must approve a mint.
func mintingApprovalsNeeded(approvers int) int {
	if n := approvers - 2; n > 1 {
		return n
	}
	return 1
}

// issuingSafeDeployments are the asset's issuing and distribution Safes for
// this owner set.
func issuingSafeDeployments(profile *userModels.User, approvers []*userModels.User, tokenizationID string) (issuing, distribution userModels.SafeDeployment) {
	owners, threshold := sharedAccessTarget(profile, approvers, mintingApprovalsNeeded(len(approvers)))
	hexOwners := make([]string, len(owners))
	for i, o := range owners {
		hexOwners[i] = o.Hex()
	}
	salt := func(role string) *big.Int {
		return new(big.Int).SetBytes(crypto.Keccak256([]byte("trovo-tokenization-" + role + ":" + tokenizationID)))
	}
	issuing = userModels.NewSafeDeployment(hexOwners, int(threshold), salt("issuing"))
	distribution = userModels.NewSafeDeployment(hexOwners, int(threshold), salt("distribution"), issuing.Address)
	return issuing, distribution
}

// AssignIssuingWallet gives the tokenization its issuing and distribution
// Safes (recorded as sub-wallets of the issuing profile, with the minting
// approvers and initiators as their shared access), or returns the ones it
// has. While the token contract is not registered yet, a change of minting
// approvers moves the asset to new Safes; afterwards the owners are fixed.
func AssignIssuingWallet(tokenizationID string, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, issuingWallet userModels.UserWallet, err error) {
	ato, _, err = GetTokenizedAssetByID(tokenizationID, gc.DB)
	if err != nil {
		return
	}
	if ato.AssetCode == nil || strings.TrimSpace(*ato.AssetCode) == "" {
		err = &tErrors.CustomError{Param: "assetCode", Err: "error-asset-code-not-set", ErrMessage: "Asset code not set."}
		return
	}
	code := strings.ToUpper(strings.NewReplacer("\r", "", "\n", "", " ", "").Replace(*ato.AssetCode))
	ato.AssetCode = &code
	profile, err := issuingProfile(gc)
	if err != nil {
		return
	}
	if ato.MintingApprovers == nil || strings.TrimSpace(*ato.MintingApprovers) == "" {
		csv := ato.GetMintingApproversInCSV(gc)
		ato.MintingApprovers = &csv
	}
	if ato.MintingInitators == nil || strings.TrimSpace(*ato.MintingInitators) == "" {
		csv := ato.GetMintingInitiatorsInCSV(gc)
		ato.MintingInitators = &csv
	}
	approvers, err := mintingUsers(*ato.MintingApprovers, gc)
	if err != nil {
		return
	}
	initiators, err := mintingUsers(*ato.MintingInitators, gc)
	if err != nil {
		return
	}
	if len(approvers) == 0 {
		err = &tErrors.CustomError{Param: "mintingApprovers", Err: "error-no-approver-or-initiator-specified", ErrMessage: "No minting approvers specified."}
		return
	}
	issuing, distribution := issuingSafeDeployments(&profile, approvers, ato.ID)

	if ato.IssuingWalletAddress != nil && *ato.IssuingWalletAddress != "" {
		if strings.EqualFold(*ato.IssuingWalletAddress, issuing.Address) {
			issuingWallet, err = userModels.UserWalletID(*ato.IssuingWalletAddress).GetWallet(gc.DB, gc)
			return
		}
		// the minting approvers changed
		if ato.ContractAddress != nil {
			err = &tErrors.CustomError{Param: "mintingApprovers", Err: "error-minting-approvers-fixed", ErrMessage: "The token contract is already registered to this asset's issuing wallet, so its minting approvers can no longer change.", Code: http.StatusConflict}
			return
		}
		if err = retireIssuingWallets(*ato.IssuingWalletAddress, gc); err != nil {
			return
		}
	}

	tag := strings.ToLower(code) + "issuer"
	iw, err := profile.BuildNewSubWallet(issuing, tag, fmt.Sprintf("Issuing wallet for %v", code), 1, distribution.Address, gc)
	if err != nil {
		return
	}
	dw, err := iw.BuildNewLinkedSubWallet(distribution, &profile, gc)
	if err != nil {
		return
	}
	needed := mintingApprovalsNeeded(len(approvers))
	var perms []userModels.WalletPermission
	for _, w := range []*userModels.UserWallet{&iw, &dw} {
		w.SharedAccessEnabled = 1
		w.NumberOfApprovalsNeeded = needed
		for _, set := range []struct {
			users      []*userModels.User
			permission string
		}{{approvers, "APPROVER"}, {initiators, "INITIATOR"}} {
			for _, u := range set.users {
				perms = append(perms, userModels.WalletPermission{ID: uuid.NewString(), WalletAddress: w.ID, TargetUsername: u.Username, Permission: set.permission})
			}
		}
	}

	ato.IssuingWalletAddress = &iw.ID
	ato.IssuingWalletAlias = &iw.Alias
	ato.MarketMakingWallet = &dw.ID
	ato.WalletToHoldAssetsNotForSale = &dw.ID
	err = gc.DB.Transaction(func(tx *gorm.DB) error {
		for _, w := range []*userModels.UserWallet{&iw, &dw} {
			if e := tx.Omit(clause.Associations).Create(w).Error; e != nil {
				return e
			}
		}
		if e := tx.Omit(clause.Associations).Create(&perms).Error; e != nil {
			return e
		}
		return tx.Omit(clause.Associations).Save(&ato).Error
	})
	if err != nil {
		log.Printf("[AssignIssuingWallet] saving issuing wallets of %v: %v", tokenizationID, err)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	for _, w := range []string{iw.ID, dw.ID} {
		if e := gc.RoachDB.Create(&userModels.TrackedAddress{Address: w}).Error; e != nil {
			// payment-history-engine also discovers untracked wallets on its own
			log.Printf("[AssignIssuingWallet] tracking %v for payment history: %v", w, e)
		}
	}
	profile.InvalidateUserCache(gc)
	for _, u := range append(approvers, initiators...) {
		u.InvalidateUserWalletCache(gc)
	}
	log.Printf("[AssignIssuingWallet] %v: issuing %v, distribution %v, %d-of-%d", code, iw.ID, dw.ID, issuing.Threshold, len(issuing.Owners))
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	issuingWallet, err = userModels.UserWalletID(iw.ID).GetWallet(gc.DB, gc)
	return
}

// retireIssuingWallets removes the records of an issuing wallet (and its
// distribution wallet) that was never deployed, when the asset moves to
// new Safes.
func retireIssuingWallets(issuingAddress string, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	w, err := userModels.UserWalletID(issuingAddress).GetWallet(gc.DB, gc)
	if err != nil {
		return nil
	}
	ids := []string{w.ID}
	if w.LinkedWalletAddress != nil {
		ids = append(ids, *w.LinkedWalletAddress)
	}
	for _, id := range ids {
		if deployed, e := aa.Deployed(ctx, gc.BantuExpansionClient, common.HexToAddress(id)); e != nil || deployed {
			return &tErrors.CustomError{Param: "mintingApprovers", Err: "error-minting-approvers-fixed", ErrMessage: "This asset's issuing wallet is already active, so its minting approvers can no longer change.", Code: http.StatusConflict}
		}
	}
	return gc.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Where("wallet_address IN ?", ids).Delete(&userModels.WalletPermission{}).Error; e != nil {
			return e
		}
		return tx.Where("id IN ?", ids).Delete(&userModels.UserWallet{}).Error
	})
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
