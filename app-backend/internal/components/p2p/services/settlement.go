package p2p

import (
	"math/big"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/components/p2p/safesigner"
	usersDB "trovo-wallet-api/internal/components/users/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// weiShift matches internal/network's decimalToWei convention exactly (a
// fixed 18-decimal shift applied uniformly to every asset on-chain,
// independent of CuratedAsset.DecimalPlaces, which is a display/rounding
// precision, not the on-chain token precision).
const weiShift = 18

func toWei(amount string) *big.Int {
	d, err := decimal.NewFromString(amount)
	if err != nil || d.IsZero() {
		return big.NewInt(0)
	}
	return d.Shift(weiShift).BigInt()
}

// ReleaseEscrowSettlement implements Plan Sections 59-61: assembles the
// Universal Safe release - buyerNetAssetAmount to assetRecipient,
// combinedPlatformFee/combinedRegulatoryFee/combinedVat to their respective
// fee wallets - and submits it via the safesigner package (app-backend's
// own Safe-transaction-assembly code, per your decision). Returns the
// canonical release transaction hash.
func ReleaseEscrowSettlement(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) (string, error) {
	recipientAddress, err := resolveAssetRecipientAddress(gc, order)
	if err != nil {
		return "", err
	}

	transfers := []safesigner.Transfer{
		{Token: order.AssetContractAddress, Recipient: recipientAddress, Amount: toWei(order.BuyerNetAssetAmount)},
	}
	if order.PlatformFeeWalletAddress != "" {
		transfers = append(transfers, safesigner.Transfer{Token: order.AssetContractAddress, Recipient: order.PlatformFeeWalletAddress, Amount: toWei(order.CombinedPlatformFee)})
	}
	if order.RegulatoryFeeWalletAddress != "" {
		transfers = append(transfers, safesigner.Transfer{Token: order.AssetContractAddress, Recipient: order.RegulatoryFeeWalletAddress, Amount: toWei(order.CombinedRegulatoryFee)})
	}
	if order.VatWalletAddress != "" {
		transfers = append(transfers, safesigner.Transfer{Token: order.AssetContractAddress, Recipient: order.VatWalletAddress, Amount: toWei(order.CombinedVat)})
	}

	txHash, err := safesigner.ExecuteTransfers(transfers)
	if err != nil {
		return "", &tErrors.CustomError{Param: "settlement", Err: "error-settlement-release-failed", ErrMessage: "Escrow settlement release failed: " + err.Error()}
	}
	return txHash, nil
}

// resolveAssetRecipientAddress prefers an explicit BuyerPayoutAddress
// (Plan Section 20) and otherwise resolves the recipient's own primary
// wallet address from their username.
func resolveAssetRecipientAddress(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) (string, error) {
	if order.BuyerPayoutAddress != "" {
		return order.BuyerPayoutAddress, nil
	}
	username := order.MerchantUsername
	if order.AssetRecipient == order.CustomerUserID {
		username = order.CustomerUsername
	}
	wallet, _, err := usersDB.GetWallet(username, gc.DB)
	if err != nil || wallet.ID == "" {
		return "", &tErrors.CustomError{Param: "assetRecipient", Err: "error-recipient-wallet-not-found", ErrMessage: "Could not resolve the asset recipient's wallet"}
	}
	return wallet.ID, nil
}
