package proceedpayouts

import (
	"context"

	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"
	stakeholderServices "admin-panel-dashboard/internal/components/stakeholder/services"
)

// DistributionClient is the stakeholder portal's DistributionPayoutClient
// on app-backend's database: a trustee's authorization registers the
// distribution's payout, and the distribution's detail shows it.
type DistributionClient struct{ Service *Service }

var _ stakeholderServices.DistributionPayoutClient = DistributionClient{}

func (d DistributionClient) Register(ctx context.Context, r stakeholderServices.DistributionPayoutRegistration) (*stakeholderModels.DistributionPayoutResponse, error) {
	if _, err := d.Service.Register(ctx, Registration{
		DistributionID: r.DistributionID, TokenizedAssetID: r.TokenizedAssetID, Amount: r.Amount, Currency: r.Currency,
	}); err != nil {
		return nil, err
	}
	return d.Service.DistributionPayout(ctx, r.DistributionID)
}

func (d DistributionClient) Get(ctx context.Context, distributionID string) (*stakeholderModels.DistributionPayoutResponse, error) {
	return d.Service.DistributionPayout(ctx, distributionID)
}
