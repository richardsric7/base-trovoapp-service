package models

// TokenizedAssetStatistics represents statistics for a specific status
type TokenizedAssetStatistics struct {
	Count                 int64   `json:"count" example:"1724"`
	TotalCurrentValue     float64 `json:"total_current_value" example:"580000000000000.00"`
	TotalTokenizedValue   float64 `json:"total_tokenized_value" example:"580000000000000.00"`
	TotalTokensToBeIssued float64 `json:"total_tokens_to_be_issued" example:"1000000.00"`
	TotalTokensToBeSold   float64 `json:"total_tokens_to_be_sold" example:"800000.00"`
	TotalPricePerToken    float64 `json:"total_price_per_token" example:"5000000.00"`
	AveragePricePerToken  float64 `json:"average_price_per_token" example:"50.00"`
}

// TokenizedAssetStatisticsResponse represents the complete statistics response
type TokenizedAssetStatisticsResponse struct {
	Total    TokenizedAssetStatistics `json:"total"`
	Approved TokenizedAssetStatistics `json:"approved"`
	Pending  TokenizedAssetStatistics `json:"pending"`
	Rejected TokenizedAssetStatistics `json:"rejected"`
}
