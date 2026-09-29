package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// OperationSubWallet is the wallet operation kind for creating sub-wallets.
const OperationSubWallet = "SUB WALLET"

// subWalletContext is what a sub-wallet operation saves when submitted.
type subWalletContext struct {
	Wallets []userModels.UserWallet      `json:"wallets"`
	Fees    []sharedconfig.FeeCollection `json:"fees"`
}

// platformCoSigner is the platform key that co-owns market-making (type 2)
// and bulk-payment (type 3) sub-wallets so the platform can operate them.
// A user has at most one of each, so it is derived from the user's
// permanent primary address.
func platformCoSigner(owner *userModels.User, walletType int) (string, error) {
	key := "primary:" + strings.ToUpper(owner.Address)
	switch walletType {
	case 2:
		kp, err := bc.MarketMakingSignerKeypair(owner.Username, key)
		if err != nil {
			return "", err
		}
		return kp.Address(), nil
	case 3:
		kp, err := bc.BulkPaymentSignerKeypair(owner.Username, key)
		if err != nil {
			return "", err
		}
		return kp.Address(), nil
	}
	return "", nil
}

// subWalletSeed is the ETH a new sub-wallet receives from the primary
// wallet so it can pay its own first network fees (the *_ACTIVATION_AMOUNT
// settings, in ETH; 0 = no seed).
func subWalletSeed(owner *userModels.User, walletType int, gc *sharedconfig.GlobalConfig) decimal.Decimal {
	primary := owner.UserWallets[0]
	id := map[int]string{1: "ISSUING_SUB_WALLET_ACTIVATION_AMOUNT", 2: "MM_SUB_WALLET_ACTIVATION_AMOUNT", 3: "BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT"}[walletType]
	if id != "" {
		if a := primary.GetActivationFee(id, gc); a.Amount > 0 {
			return decimal.NewFromFloat(a.Amount)
		}
	}
	if a := primary.GetActivationFee("SUB_WALLET_ACTIVATION_AMOUNT", gc); a.Amount > 0 {
		return decimal.NewFromFloat(a.Amount)
	}
	return decimal.Zero
}

func primaryWalletOf(owner *userModels.User) (*userModels.UserWallet, error) {
	for i := range owner.UserWallets {
		if owner.UserWallets[i].PrimaryWallet == 1 {
			return &owner.UserWallets[i], nil
		}
	}
	return nil, &tErrors.CustomError{Param: "publicKey", Err: "error-primary-wallet-not-found", ErrMessage: "Primary wallet not found.", Code: http.StatusNotFound}
}

// CreateNewSubWallet creates a sub-wallet in two steps:
//
//  1. without primarySignature: derives the new Safe (and, for an issuing
//     wallet, its linked distribution Safe), and builds the primary
//     wallet's operation that deploys them, seeds them with ETH and pays
//     the creation fee; returns it as transaction for the user to sign;
//  2. with transaction + primarySignature: records the wallets and submits
//     the operation.
//
// The primary wallet pays for everything; creation is refused when it
// cannot.
func CreateNewSubWallet(accountOwner *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig) (*userModels.SubWalletInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	subWalletInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	subWalletInfo.SubWalletMustSign = 0
	subWalletInfo.LinkedWalletMustSign = 0

	primary, err := primaryWalletOf(accountOwner)
	if err != nil {
		return subWalletInfo, err
	}
	if len(subWalletInfo.PrimarySignature) > 0 && len(subWalletInfo.Transaction) > 0 {
		return submitSubWallet(ctx, accountOwner, primary, subWalletInfo, gc)
	}
	if subWalletInfo.WalletType < 0 || subWalletInfo.WalletType > 3 {
		return subWalletInfo, &tErrors.CustomError{Param: "walletType", Err: "error-invalid-wallet-type", ErrMessage: "Unknown wallet type.", Code: http.StatusBadRequest}
	}

	// the new Safe(s)
	owners := []string{accountOwner.PrimarySigner}
	if co, err := platformCoSigner(accountOwner, subWalletInfo.WalletType); err != nil {
		log.Printf("[CreateNewSubWallet] platform co-signer for %v: %v", accountOwner.Username, err)
		return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
	} else if co != "" {
		owners = append(owners, co)
	}
	salt, err := userModels.RandomSaltNonce()
	if err != nil {
		return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
	}
	sub := userModels.NewSafeDeployment(owners, 1, salt)
	subWalletInfo.Address = sub.Address

	// an issuing wallet always gets its linked distribution wallet
	withLinked := subWalletInfo.WalletType == 1 || len(subWalletInfo.LinkedWalletAddress) > 0
	var linked userModels.SafeDeployment
	subWalletInfo.LinkedWalletAddress = ""
	if withLinked {
		if linked, err = userModels.NewSubWalletSafeDeployment(accountOwner.PrimarySigner); err != nil {
			return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
		}
		subWalletInfo.LinkedWalletAddress = linked.Address
	}

	subWalletObj, err := accountOwner.BuildNewSubWallet(sub, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, subWalletInfo.WalletType, subWalletInfo.LinkedWalletAddress, gc)
	if err != nil {
		return subWalletInfo, err
	}
	subWalletInfo.Alias = subWalletObj.Alias
	wallets := []userModels.UserWallet{subWalletObj}
	deploys := []userModels.SafeDeployment{sub}
	if withLinked {
		linkedWallet, err := subWalletObj.BuildNewLinkedSubWallet(linked, accountOwner, gc)
		if err != nil {
			return subWalletInfo, err
		}
		wallets = append(wallets, linkedWallet)
		deploys = append(deploys, linked)
	}

	// calls from the primary wallet: deploy, seed, fee
	cfg := network.AAConfig()
	var calls []aa.Call
	for _, d := range deploys {
		salt, _ := new(big.Int).SetString(d.SaltNonce, 10)
		calls = append(calls, cfg.DeploySafeCall(d.OwnerAddresses(), int64(d.Threshold), salt))
	}
	subWalletInfo.Messages = []string{fmt.Sprintf("Creates the sub-wallet %v (%v).", subWalletObj.Alias, sub.Address)}
	if withLinked {
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Creates its linked distribution wallet (%v).", linked.Address))
	}
	seed := subWalletSeed(accountOwner, subWalletInfo.WalletType, gc)
	seedWei := seed.Shift(18).Truncate(0).BigInt()
	ethNeeded := new(big.Int)
	if seedWei.Sign() > 0 {
		for _, d := range deploys {
			calls = append(calls, aa.NativeTransfer(common.HexToAddress(d.Address), seedWei))
			ethNeeded.Add(ethNeeded, seedWei)
		}
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v ETH is moved from your primary wallet to each new wallet for its network fees.", seed))
	}

	var fees []sharedconfig.FeeCollection
	creationFee := primary.GetSubwalletCreationFee(gc)
	// the fee is set in USD (fee_fixed) and paid in the fee asset; the
	// platform's tokenization issuing profile pays none
	if creationFee.Inactive == 0 && creationFee.FeeFixed > 0 && !feeExemptProfile(accountOwner.Username) {
		feeAddr, err := feeWalletAddress(creationFee.FeeWalletSecretKey, "sub-wallet creation fee wallet", gc)
		if err != nil {
			return subWalletInfo, err
		}
		if !common.IsHexAddress(creationFee.FeeContractAddress) {
			gc.LogDiscordFailedRequest("[CreateNewSubWallet] SUBWALLET_CREATION_FEE has no fee_contract_address (the fee must be a stablecoin)")
			return subWalletInfo, &tErrors.CustomError{Param: "feeAmount", Err: "error-fee-asset-not-priceable", ErrMessage: "The sub-wallet creation fee is not configured correctly. Please try again later."}
		}
		feeAsset := basetxn.CreditAsset{Code: creationFee.FeeAssetCode, Issuer: creationFee.FeeContractAddress}
		amount, err := usdAmountIn(creationFee.FeeFixed, creationFee.FeeAssetCode, creationFee.FeeContractAddress, gc)
		if err != nil {
			gc.LogDiscordFailedRequest(fmt.Sprintf("[CreateNewSubWallet] the sub-wallet creation fee asset %v cannot be priced in USD", creationFee.FeeAssetCode))
			return subWalletInfo, err
		}
		decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, feeAsset)
		if err != nil {
			return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
		}
		amount = amount.Truncate(int32(decimals))
		units, err := baseUnits(amount.String(), decimals)
		if err != nil {
			return subWalletInfo, err
		}
		bal, err := network.B20BalanceOf(gc.BantuExpansionClient, feeAsset.Issuer, primary.ID, decimals)
		if err != nil {
			return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
		}
		if bal.LessThan(amount) {
			return subWalletInfo, &tErrors.CustomError{Param: "username", Err: "error-primary-wallet-underfunded", ErrMessage: fmt.Sprintf("%v %v is required on wallet %v to pay the sub-wallet creation fee. Please first fund the wallet with at least %v %v.", amount, feeAsset.Code, accountOwner.Username, amount.Sub(bal), feeAsset.Code), Code: http.StatusBadRequest}
		}
		calls = append(calls, transferCall(feeAsset, feeAddr, units))
		subWalletInfo.FeeAmount, subWalletInfo.FeeCode = amount.String(), feeAsset.Code
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("%v %v ($%v USD) will be deducted from wallet %v as the sub-wallet creation fee.", amount, feeAsset.Code, creationFee.FeeFixed, accountOwner.Username))
		fees = append(fees, feeRecord("SUBWALLET_CREATION", primary, accountOwner, feeAsset, amount.String(), feeAddr, 0))
	}

	op, err := PrepareWalletOperation(ctx, OperationSubWallet, accountOwner, accountOwner, primary, calls, 0, subWalletContext{Wallets: wallets, Fees: fees}, gc)
	if err != nil {
		return subWalletInfo, err
	}
	// the primary wallet must hold the ETH it sends, plus the network fee
	// when that is paid in ETH
	needed := new(big.Int).Set(ethNeeded)
	if op.Prepared.Quote == nil && op.Prepared.MaxCostWei != nil {
		needed.Add(needed, op.Prepared.MaxCostWei.ToInt())
	}
	if needed.Sign() > 0 {
		bal, err := gc.BantuExpansionClient.BalanceAt(ctx, common.HexToAddress(primary.ID), nil)
		if err != nil {
			return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
		}
		if bal.Cmp(needed) < 0 {
			return subWalletInfo, &tErrors.CustomError{Param: "publicKey", Err: "error-primary-account-underfunded",
				ErrMessage: fmt.Sprintf("Your primary wallet needs at least %v ETH to create this sub-wallet.", decimal.NewFromBigInt(needed, -18).Round(8)), Code: http.StatusBadRequest}
		}
	}
	var deployed []string
	for _, d := range deploys {
		deployed = append(deployed, d.Address)
	}
	gc.DB.Model(&userModels.WalletOperation{}).Where("id = ?", op.Record.ID).Update("deploys", strings.Join(deployed, ","))

	subWalletInfo.Transaction = op.Transaction
	subWalletInfo.Messages = append(subWalletInfo.Messages, operationMessages(op.Prepared)...)
	return subWalletInfo, nil
}

// submitSubWallet records the wallets of a signed sub-wallet operation and
// submits it.
func submitSubWallet(ctx context.Context, accountOwner *userModels.User, primary *userModels.UserWallet, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig) (*userModels.SubWalletInfo, error) {
	rec, p, err := LoadWalletOperation(subWalletInfo.Transaction, primary.ID, OperationSubWallet, gc)
	if err != nil {
		return subWalletInfo, err
	}
	var c subWalletContext
	if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil || len(c.Wallets) == 0 {
		return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
	}

	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	for i := range c.Wallets {
		if err := dbTX.Omit(clause.Associations).Create(&c.Wallets[i]).Error; err != nil {
			log.Printf("[CreateNewSubWallet] saving %v for %v: %v", c.Wallets[i].ID, accountOwner.Username, err)
			return subWalletInfo, &tErrors.CustomError{Param: "publicKey", Err: "error-saving-subwallet", ErrMessage: "There is an error saving the sub-wallet. Please try again later."}
		}
	}
	hash, err := SignSingleOwnerOperation(ctx, rec, p, accountOwner.PrimarySigner, subWalletInfo.PrimarySignature, gc)
	if err != nil {
		return subWalletInfo, err
	}
	if err := dbTX.Commit().Error; err != nil {
		log.Printf("[CreateNewSubWallet] operation %v submitted but wallets not saved: %v", hash, err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[CreateNewSubWallet] operation %v for %v submitted but its wallets were not saved: %v", hash, accountOwner.Username, err))
		return subWalletInfo, &tErrors.ErrorTemporaryServerError{}
	}
	recordOperationFees(&userModels.WalletOperation{Context: rec.Context}, hash, gc)

	subWalletInfo.TransactionID = hash
	subWalletInfo.Address = c.Wallets[0].ID
	subWalletInfo.Alias = c.Wallets[0].Alias
	if len(c.Wallets) > 1 {
		subWalletInfo.LinkedWalletAddress = c.Wallets[1].ID
	}
	for _, w := range c.Wallets {
		// payment history tracking; the engine also discovers untracked wallets itself
		if err := gc.RoachDB.Create(&userModels.TrackedAddress{Address: w.ID}).Error; err != nil {
			discord.Say(fmt.Sprintf("[CreateNewSubWallet] tracking %v for %v failed: %v", w.ID, accountOwner.Username, err))
		}
	}
	accountOwner.InvalidateUserCache(gc)
	accountOwner.InvalidateUserWalletCache(gc)
	return subWalletInfo, nil
}
