package users

import (
	"fmt"
	"net/http"
	"strings"

	"trovo-wallet-api/internal/basetxn"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"
	"trovo-wallet-api/internal/validators"
)

// ApproveWalletAssetAuthorization grants or revokes walletAddress's
// authorization to hold/send a regulated asset - the compliance-approval
// equivalent of the issuer signing a Stellar SetTrustLineFlags operation.
// approverSigner is the signer address the request was authenticated
// against (see middleware.ExtractSigner): only the signer of the asset's
// own issuing wallet may call this, mirroring how only the Stellar issuer
// account could authorize a trustline.
func ApproveWalletAssetAuthorization(approverSigner string, req *userModels.WalletAssetAuthorizationRequest, gc *sharedconfig.GlobalConfig) (network.WalletAssetAuthorization, error) {
	var result network.WalletAssetAuthorization

	if err := validators.ValidateAddressFormat(req.WalletAddress); err != nil {
		return result, err
	}
	if err := validators.ValidateAddressFormat(req.ContractAddress); err != nil {
		return result, err
	}
	if err := validators.ValidateAssetCodeFormat(req.AssetCode); err != nil {
		return result, err
	}

	if !gc.IsValidTokenizedAsset(req.AssetCode) {
		return result, &tErrors.CustomError{
			Param:      "assetCode",
			Err:        "error-not-a-regulated-asset",
			ErrMessage: fmt.Sprintf("%v is not a market-ready tokenized/regulated asset - only those require wallet-level authorization.", req.AssetCode),
			Code:       http.StatusBadRequest,
		}
	}

	issuingWallet, _, err := usersDB.GetWallet(req.ContractAddress, gc.DB)
	if err != nil {
		return result, err
	}
	if issuingWallet.WalletType != 1 {
		return result, &tErrors.CustomError{
			Param:      "contractAddress",
			Err:        "error-not-an-issuing-wallet",
			ErrMessage: fmt.Sprintf("%v is not an asset-issuing wallet.", req.ContractAddress),
			Code:       http.StatusBadRequest,
		}
	}
	if !strings.EqualFold(issuingWallet.Signer, approverSigner) {
		return result, &tErrors.CustomError{
			Param:      "contractAddress",
			Err:        "error-not-authorized-issuer",
			ErrMessage: "Only the asset's own issuing wallet may approve or revoke wallet authorization for it.",
			Code:       http.StatusForbidden,
		}
	}

	asset := basetxn.CreditAsset{Code: strings.ToUpper(req.AssetCode), Issuer: req.ContractAddress}
	if err := network.SetWalletAssetAuthorization(req.WalletAddress, asset, req.Authorized, approverSigner, req.Reason); err != nil {
		return result, &tErrors.ErrorTemporaryServerError{}
	}

	result = network.WalletAssetAuthorization{
		WalletAddress:   strings.ToLower(req.WalletAddress),
		AssetCode:       strings.ToUpper(req.AssetCode),
		ContractAddress: strings.ToLower(req.ContractAddress),
		Authorized:      req.Authorized,
		ApprovedBy:      strings.ToLower(approverSigner),
		Reason:          req.Reason,
	}
	return result, nil
}

// GetWalletAssetAuthorizations lists every wallet currently authorized (or
// previously revoked) for a regulated asset, for the asset's own issuing
// wallet to review. Same issuer-signature gate as
// ApproveWalletAssetAuthorization.
func GetWalletAssetAuthorizations(approverSigner, assetCode, contractAddress string, gc *sharedconfig.GlobalConfig) ([]network.WalletAssetAuthorization, error) {
	if err := validators.ValidateAddressFormat(contractAddress); err != nil {
		return nil, err
	}
	if err := validators.ValidateAssetCodeFormat(assetCode); err != nil {
		return nil, err
	}

	issuingWallet, _, err := usersDB.GetWallet(contractAddress, gc.DB)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(issuingWallet.Signer, approverSigner) {
		return nil, &tErrors.CustomError{
			Param:      "contractAddress",
			Err:        "error-not-authorized-issuer",
			ErrMessage: "Only the asset's own issuing wallet may view its wallet authorizations.",
			Code:       http.StatusForbidden,
		}
	}

	return network.WalletAssetAuthorizations(strings.ToUpper(assetCode), contractAddress)
}
