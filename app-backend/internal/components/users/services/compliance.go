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

	if _, err := authorizeIssuerSigner(approverSigner, req.AssetCode, req.ContractAddress, "approve or revoke wallet authorization for it", gc); err != nil {
		return result, err
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

	if _, err := authorizeIssuerSigner(approverSigner, assetCode, contractAddress, "view its wallet authorizations", gc); err != nil {
		return nil, err
	}

	return network.WalletAssetAuthorizations(strings.ToUpper(assetCode), contractAddress)
}

// authorizeIssuerSigner resolves the market-ready tokenized asset whose
// token contract is contractAddress (authorizations are keyed by the token
// contract, never the issuer) and checks approverSigner is the recorded
// signer of that asset's issuing Safe wallet.
func authorizeIssuerSigner(approverSigner, assetCode, contractAddress, action string, gc *sharedconfig.GlobalConfig) (userModels.TokenizedAsset, error) {
	var ta userModels.TokenizedAsset
	e := gc.DB.Where("Asset_Tokenization_Status > 3 AND Asset_Code = upper(?) AND lower(contract_address) = lower(?)", assetCode, contractAddress).First(&ta).Error
	if e != nil || ta.IssuingWalletAddress == nil {
		return ta, &tErrors.CustomError{
			Param:      "contractAddress",
			Err:        "error-not-a-regulated-asset",
			ErrMessage: fmt.Sprintf("%v at %v is not a market-ready tokenized/regulated asset token contract - only those require wallet-level authorization.", assetCode, contractAddress),
			Code:       http.StatusBadRequest,
		}
	}
	issuingWallet, err := usersDB.GetWallet(*ta.IssuingWalletAddress, gc.DB)
	if err != nil {
		return ta, err
	}
	if !strings.EqualFold(issuingWallet.Signer, approverSigner) {
		return ta, &tErrors.CustomError{
			Param:      "contractAddress",
			Err:        "error-not-authorized-issuer",
			ErrMessage: fmt.Sprintf("Only the signer of the asset's issuing wallet may %v.", action),
			Code:       http.StatusForbidden,
		}
	}
	return ta, nil
}
