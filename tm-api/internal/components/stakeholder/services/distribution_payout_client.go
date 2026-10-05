package services

import (
	"context"

	"admin-panel-dashboard/internal/components/stakeholder/models"
)

// DistributionPayoutRegistration is an authorized distribution to pay out.
type DistributionPayoutRegistration struct {
	DistributionID   string `json:"distribution_id"`
	TokenizedAssetID string `json:"tokenized_asset_id"`
	Amount           string `json:"amount"`
	Currency         string `json:"currency"`
}

// DistributionPayoutClient registers an authorized distribution's payout and
// reads it back. proceedpayouts.DistributionClient implements it on
// app-backend's database, where payout-engine pays it.
type DistributionPayoutClient interface {
	Register(ctx context.Context, registration DistributionPayoutRegistration) (*models.DistributionPayoutResponse, error)
	Get(ctx context.Context, distributionID string) (*models.DistributionPayoutResponse, error)
}
