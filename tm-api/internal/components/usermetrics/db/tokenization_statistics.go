package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	trovosdk "admin-panel-dashboard/internal/trovosdk"

	"gorm.io/gorm"
)

type statsResult struct {
	Count                 int64   `gorm:"column:count"`
	TotalCurrentValue     float64 `gorm:"column:total_current_value"`
	TotalTokenizedValue   float64 `gorm:"column:total_tokenized_value"`
	TotalTokensToBeIssued float64 `gorm:"column:total_tokens_to_be_issued"`
	TotalTokensToBeSold   float64 `gorm:"column:total_tokens_to_be_sold"`
	TotalPricePerToken    float64 `gorm:"column:total_price_per_token"`
	AveragePricePerToken  float64 `gorm:"column:average_price_per_token"`
}

// GetTokenizedAssetStatistics retrieves aggregated statistics from tokenized_assets table
// Uses TokenizedAssetJSON model which has proper GORM tags and TableName() method
func GetTokenizedAssetStatistics(db *gorm.DB) (*models.TokenizedAssetStatisticsResponse, error) {
	var result models.TokenizedAssetStatisticsResponse

	// Get total statistics (all records, including drafts status 0)
	var totalStats statsResult
	if err := db.Model(&trovosdk.TokenizedAssetJSON{}).
		Select(`COUNT(*) as count,
			COALESCE(SUM(asset_current_value), 0) as total_current_value,
			COALESCE(SUM(value_of_tokenized_asset), 0) as total_tokenized_value,
			COALESCE(SUM(number_of_token_to_be_issued), 0) as total_tokens_to_be_issued,
			COALESCE(SUM(number_of_token_to_be_sold), 0) as total_tokens_to_be_sold,
			COALESCE(SUM(price_per_token), 0) as total_price_per_token,
			COALESCE(AVG(price_per_token), 0) as average_price_per_token`).
		Scan(&totalStats).Error; err != nil {
		return nil, err
	}

	// Get approved statistics (status 4, 5, 6, 7, 8 - Market Ready)
	var approvedStats statsResult
	if err := db.Model(&trovosdk.TokenizedAssetJSON{}).
		Where("asset_tokenization_status IN ?", []int{4, 5, 6, 7, 8}).
		Select(`COUNT(*) as count,
			COALESCE(SUM(asset_current_value), 0) as total_current_value,
			COALESCE(SUM(value_of_tokenized_asset), 0) as total_tokenized_value,
			COALESCE(SUM(number_of_token_to_be_issued), 0) as total_tokens_to_be_issued,
			COALESCE(SUM(number_of_token_to_be_sold), 0) as total_tokens_to_be_sold,
			COALESCE(SUM(price_per_token), 0) as total_price_per_token,
			COALESCE(AVG(price_per_token), 0) as average_price_per_token`).
		Scan(&approvedStats).Error; err != nil {
		return nil, err
	}

	// Get pending statistics (status 0, 1, 2, 3 - Including Drafts)
	// Includes status 0 to match Total count which includes drafts
	var pendingStats statsResult
	if err := db.Model(&trovosdk.TokenizedAssetJSON{}).
		Where("asset_tokenization_status IN ?", []int{0, 1, 2, 3}).
		Select(`COUNT(*) as count,
			COALESCE(SUM(asset_current_value), 0) as total_current_value,
			COALESCE(SUM(value_of_tokenized_asset), 0) as total_tokenized_value,
			COALESCE(SUM(number_of_token_to_be_issued), 0) as total_tokens_to_be_issued,
			COALESCE(SUM(number_of_token_to_be_sold), 0) as total_tokens_to_be_sold,
			COALESCE(SUM(price_per_token), 0) as total_price_per_token,
			COALESCE(AVG(price_per_token), 0) as average_price_per_token`).
		Scan(&pendingStats).Error; err != nil {
		return nil, err
	}

	// Get rejected statistics (status > 8)
	var rejectedStats statsResult
	if err := db.Model(&trovosdk.TokenizedAssetJSON{}).
		Where("asset_tokenization_status > ?", trovosdk.TokenizationStatusRejectedMin).
		Select(`COUNT(*) as count,
			COALESCE(SUM(asset_current_value), 0) as total_current_value,
			COALESCE(SUM(value_of_tokenized_asset), 0) as total_tokenized_value,
			COALESCE(SUM(number_of_token_to_be_issued), 0) as total_tokens_to_be_issued,
			COALESCE(SUM(number_of_token_to_be_sold), 0) as total_tokens_to_be_sold,
			COALESCE(SUM(price_per_token), 0) as total_price_per_token,
			COALESCE(AVG(price_per_token), 0) as average_price_per_token`).
		Scan(&rejectedStats).Error; err != nil {
		return nil, err
	}

	// Build response
	result = models.TokenizedAssetStatisticsResponse{
		Total: models.TokenizedAssetStatistics{
			Count:                 totalStats.Count,
			TotalCurrentValue:     totalStats.TotalCurrentValue,
			TotalTokenizedValue:   totalStats.TotalTokenizedValue,
			TotalTokensToBeIssued: totalStats.TotalTokensToBeIssued,
			TotalTokensToBeSold:   totalStats.TotalTokensToBeSold,
			TotalPricePerToken:    totalStats.TotalPricePerToken,
			AveragePricePerToken:  totalStats.AveragePricePerToken,
		},
		Approved: models.TokenizedAssetStatistics{
			Count:                 approvedStats.Count,
			TotalCurrentValue:     approvedStats.TotalCurrentValue,
			TotalTokenizedValue:   approvedStats.TotalTokenizedValue,
			TotalTokensToBeIssued: approvedStats.TotalTokensToBeIssued,
			TotalTokensToBeSold:   approvedStats.TotalTokensToBeSold,
			TotalPricePerToken:    approvedStats.TotalPricePerToken,
			AveragePricePerToken:  approvedStats.AveragePricePerToken,
		},
		Pending: models.TokenizedAssetStatistics{
			Count:                 pendingStats.Count,
			TotalCurrentValue:     pendingStats.TotalCurrentValue,
			TotalTokenizedValue:   pendingStats.TotalTokenizedValue,
			TotalTokensToBeIssued: pendingStats.TotalTokensToBeIssued,
			TotalTokensToBeSold:   pendingStats.TotalTokensToBeSold,
			TotalPricePerToken:    pendingStats.TotalPricePerToken,
			AveragePricePerToken:  pendingStats.AveragePricePerToken,
		},
		Rejected: models.TokenizedAssetStatistics{
			Count:                 rejectedStats.Count,
			TotalCurrentValue:     rejectedStats.TotalCurrentValue,
			TotalTokenizedValue:   rejectedStats.TotalTokenizedValue,
			TotalTokensToBeIssued: rejectedStats.TotalTokensToBeIssued,
			TotalTokensToBeSold:   rejectedStats.TotalTokensToBeSold,
			TotalPricePerToken:    rejectedStats.TotalPricePerToken,
			AveragePricePerToken:  rejectedStats.AveragePricePerToken,
		},
	}

	return &result, nil
}
