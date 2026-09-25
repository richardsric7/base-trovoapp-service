package users

// WalletAssetAuthorizationRequest is submitted by a tokenized/regulated
// asset's own issuing wallet to grant or revoke another wallet's
// authorization to hold/send that asset - the compliance-approval
// equivalent of the issuer signing a Stellar SetTrustLineFlags operation.
type WalletAssetAuthorizationRequest struct {
	WalletAddress   string `json:"walletAddress"`
	AssetCode       string `json:"assetCode"`
	ContractAddress string `json:"contractAddress"`
	Authorized      bool   `json:"authorized"`
	Reason          string `json:"reason"`
}
