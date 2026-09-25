package services

import (
	"context"
	"net/http"
	"strings"

	"admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type FundManagementService struct {
	adminDB  *gorm.DB
	walletDB *gorm.DB
	assets   *AssetService
}

func NewFundManagementService(adminDB, walletDB *gorm.DB, assets *AssetService) *FundManagementService {
	return &FundManagementService{adminDB: adminDB, walletDB: walletDB, assets: assets}
}

func (s *FundManagementService) TrusteeSummary(ctx context.Context, auth AuthContext) (*models.TrusteeFundManagementSummary, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	if s.adminDB == nil || s.walletDB == nil || s.assets == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "fund management data sources are not configured")
	}
	assets, err := s.assets.ListAllAssets(ctx, auth)
	if err != nil {
		return nil, err
	}
	assetIDs := make([]string, 0, len(assets))
	fees := make(map[string]decimal.Decimal)
	for _, asset := range assets {
		assetIDs = append(assetIDs, asset.ID)
		currency := strings.ToUpper(strings.TrimSpace(asset.AssetQuoteCurrency))
		if currency != "" && asset.ConfiguredFeeValue != 0 {
			fees[currency] = fees[currency].Add(decimal.NewFromFloat(asset.ConfiguredFeeValue))
		}
	}
	result := &models.TrusteeFundManagementSummary{
		AssetCount: len(assets), TotalAmountProcessed: []models.CurrencyTotal{}, TotalAmountFromPrimarySales: []models.CurrencyTotal{},
		TotalIncomeFromAssets: []models.CurrencyTotal{}, TotalAmountPaidToIssuers: []models.CurrencyTotal{},
		TotalPaidToInvestors: []models.CurrencyTotal{}, MilestonePaymentBalance: []models.CurrencyTotal{},
		TotalFeesGenerated: currencyTotals(fees),
	}
	if len(assetIDs) == 0 {
		return result, nil
	}

	primarySales, err := walletPrimarySalesTotals(ctx, s.walletDB, assetIDs)
	if err != nil {
		return nil, err
	}
	assetIncome, err := adminCurrencyTotals(ctx, s.adminDB.Model(&models.RevenueRecord{}).
		Where("asset_id IN ?", assetIDs), "amount")
	if err != nil {
		return nil, err
	}
	issuerPayments, err := adminCurrencyTotals(ctx, s.adminDB.Model(&models.FundReleaseRequest{}).
		Where("asset_id IN ? AND trustee_org_id = ? AND status = ?", assetIDs, auth.OrganizationID, models.FundReleaseStatusCompleted), "amount")
	if err != nil {
		return nil, err
	}
	investorPayments, err := adminCurrencyTotals(ctx, s.adminDB.Model(&models.Distribution{}).
		Where("asset_id IN ? AND trustee_org_id = ? AND status = ?", assetIDs, auth.OrganizationID, models.DistributionStatusCompleted), "amount")
	if err != nil {
		return nil, err
	}
	milestoneBalance, err := adminCurrencyTotals(ctx, s.adminDB.Model(&models.SegregatedAccount{}).
		Where("asset_id IN ? AND account_type = ? AND status = ?", assetIDs, models.AccountTypeDevelopmentFund, models.AccountStatusActive), "balance")
	if err != nil {
		return nil, err
	}

	result.TotalAmountProcessed = currencyTotals(mergeCurrencyAmounts(issuerPayments, investorPayments))
	result.TotalAmountFromPrimarySales = currencyTotals(primarySales)
	result.TotalIncomeFromAssets = currencyTotals(assetIncome)
	result.TotalAmountPaidToIssuers = currencyTotals(issuerPayments)
	result.TotalPaidToInvestors = currencyTotals(investorPayments)
	result.MilestonePaymentBalance = currencyTotals(milestoneBalance)
	return result, nil
}

type currencyAmountRow struct {
	Currency string
	Amount   decimal.Decimal
}

func adminCurrencyTotals(ctx context.Context, query *gorm.DB, amountColumn string) (map[string]decimal.Decimal, error) {
	var rows []currencyAmountRow
	if err := query.WithContext(ctx).
		Select("UPPER(currency) AS currency, COALESCE(SUM(" + amountColumn + "), 0) AS amount").
		Group("UPPER(currency)").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return currencyRowsToMap(rows), nil
}

func walletPrimarySalesTotals(ctx context.Context, walletDB *gorm.DB, assetIDs []string) (map[string]decimal.Decimal, error) {
	var rows []currencyAmountRow
	if err := walletDB.WithContext(ctx).Table("tokenized_asset_subscriptions AS subscriptions").
		Select("UPPER(assets.asset_quote_currency) AS currency, COALESCE(SUM(subscriptions.amount), 0) AS amount").
		Joins("JOIN tokenized_assets AS assets ON assets.id = subscriptions.tokenized_asset_id").
		Where("subscriptions.tokenized_asset_id IN ?", assetIDs).
		Group("UPPER(assets.asset_quote_currency)").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return currencyRowsToMap(rows), nil
}

func currencyRowsToMap(rows []currencyAmountRow) map[string]decimal.Decimal {
	result := make(map[string]decimal.Decimal, len(rows))
	for _, row := range rows {
		currency := strings.ToUpper(strings.TrimSpace(row.Currency))
		if currency != "" {
			result[currency] = result[currency].Add(row.Amount)
		}
	}
	return result
}

func mergeCurrencyAmounts(groups ...map[string]decimal.Decimal) map[string]decimal.Decimal {
	result := make(map[string]decimal.Decimal)
	for _, group := range groups {
		for currency, amount := range group {
			result[currency] = result[currency].Add(amount)
		}
	}
	return result
}
