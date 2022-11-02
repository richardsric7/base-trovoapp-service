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
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	Commit               int      `json:"commit"`
	Multiparty           int      `json:"-"`
	Memo                 string   `json:"memo"`
	ReturnedDescription  string   `json:"-"`
}
