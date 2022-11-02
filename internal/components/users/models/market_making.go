package users

type MakeOfferRequest struct {
	MarketMakingWalletPK string   `json:"marketMakingWallet"`
	OfferType            string   `json:"offerType"`
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	CurrencyCode         string   `json:"currencyCode"`
	CurrencyIssuer       string   `json:"currencyIssuer"`
	PricePerAsset        string   `json:"pricePerAsset"`
	Quantity             string   `json:"quantity"`
	Transaction          string   `json:"transaction"`
	Messages             []string `json:"messages"`
}
