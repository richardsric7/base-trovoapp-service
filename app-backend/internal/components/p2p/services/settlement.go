package p2p

import (
	"context"
	"math/big"
	"trovo-wallet-api/internal/basetxn"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/components/p2p/safesigner"
	usersDB "trovo-wallet-api/internal/components/users/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// toWei converts a human-readable decimal amount to the token's raw
// on-chain uint256 unit, using decimals as resolved by
// network.AssetDecimals for the specific asset being settled - not a
// fixed shift assumed the same for every asset (a B20 token's decimals
// vary - USDC/USDT use 6, WBTC uses 8 - independent of
// CuratedAsset.DecimalPlaces, which is a display/rounding precision, not
// the on-chain token precision).
func toWei(amount string, decimals uint8) *big.Int {
	d, err := decimal.NewFromString(amount)
	if err != nil || d.IsZero() {
		return big.NewInt(0)
	}
	return d.Shift(int32(decimals)).BigInt()
}

// resolveOrderAsset returns the basetxn.Asset order settles in: the
// native asset if it has no ERC-20 contract address, otherwise a B20
// CreditAsset for that contract - just enough of the Asset interface for
// network.AssetDecimals to resolve against.
func resolveOrderAsset(order *p2pModels.Order) basetxn.Asset {
	if order.AssetContractAddress == "" {
		return basetxn.NativeAsset{}
	}
	return basetxn.CreditAsset{Code: order.Asset, Issuer: order.AssetContractAddress}
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

	decimals, err := network.AssetDecimals(context.Background(), network.GetBlockchainClient(), resolveOrderAsset(order))
	if err != nil {
		return "", &tErrors.CustomError{Param: "settlement", Err: "error-settlement-release-failed", ErrMessage: "Escrow settlement release failed: could not resolve asset decimals: " + err.Error()}
	}

	transfers := []safesigner.Transfer{
		{Token: order.AssetContractAddress, Recipient: recipientAddress, Amount: toWei(order.BuyerNetAssetAmount, decimals)},
	}
	if order.PlatformFeeWalletAddress != "" {
		transfers = append(transfers, safesigner.Transfer{Token: order.AssetContractAddress, Recipient: order.PlatformFeeWalletAddress, Amount: toWei(order.CombinedPlatformFee, decimals)})
	}
	if order.RegulatoryFeeWalletAddress != "" {
		transfers = append(transfers, safesigner.Transfer{Token: order.AssetContractAddress, Recipient: order.RegulatoryFeeWalletAddress, Amount: toWei(order.CombinedRegulatoryFee, decimals)})
	}
	if order.VatWalletAddress != "" {
		transfers = append(transfers, safesigner.Transfer{Token: order.AssetContractAddress, Recipient: order.VatWalletAddress, Amount: toWei(order.CombinedVat, decimals)})
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
