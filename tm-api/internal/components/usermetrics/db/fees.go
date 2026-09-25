package usermetrics

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// FeeCollection represents the fee_collections table
type FeeCollection struct {
	ID                         string    `json:"id" gorm:"primaryKey"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	FromUsername               string    `json:"from_username"`
	FromWalletAddress          string    `json:"from_wallet_address"`
	FromWalletAlias            string    `json:"from_wallet_alias"`
	BelongsToEnterpriseProfile string    `json:"belongs_to_enterprise_profile"`
	FeeType                    string    `json:"fee_type"`
	Amount                     float64   `json:"amount"`
	AssetCode                  string    `json:"asset_code"`
	AssetIssuer                string    `json:"asset_issuer"`
	DestinationWallet          string    `json:"destination_wallet"`
	SharedAccessOperation      int       `json:"shared_access_operation"`
	TransactionHash            string    `json:"transaction_hash"`
	Processed                  int       `json:"processed"`
}

// ServiceLinkServiceFee represents the service_link_service_fees table
type ServiceLinkServiceFee struct {
	ServiceLinkID string    `json:"service_link_id" gorm:"primaryKey"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PaymentFee    float64   `json:"payment_fee"`
	SwapFee       float64   `json:"swap_fee"`
	SubwalletFee  float64   `json:"subwallet_fee"`
}

// GetAllFeeCollections retrieves all fee collections
func GetAllFeeCollections(db *gorm.DB) ([]FeeCollection, error) {
	var fees []FeeCollection
	if err := db.Order("created_at desc").Find(&fees).Error; err != nil {
		return nil, err
	}
	return fees, nil
}

// GetAllServiceLinkServiceFees retrieves all service link service fees
func GetAllServiceLinkServiceFees(db *gorm.DB) ([]ServiceLinkServiceFee, error) {
	var fees []ServiceLinkServiceFee
	if err := db.Find(&fees).Error; err != nil {
		return nil, err
	}
	return fees, nil
}

// GetServiceLinkServiceFeeByID retrieves a service link service fee by ID
func GetServiceLinkServiceFeeByID(db *gorm.DB, id string) (*ServiceLinkServiceFee, error) {
	var fee ServiceLinkServiceFee
	if err := db.Where("service_link_id = ?", id).First(&fee).Error; err != nil {
		return nil, err
	}
	return &fee, nil
}

// SaveServiceLinkServiceFee creates or updates a service link service fee
func SaveServiceLinkServiceFee(db *gorm.DB, fee *ServiceLinkServiceFee) error {
	var count int64
	if err := db.Table("service_links").Where("id = ?", fee.ServiceLinkID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("service link not found")
	}
	return db.Save(fee).Error
}

// DeleteServiceLinkServiceFee deletes a service link service fee by ID
func DeleteServiceLinkServiceFee(db *gorm.DB, id string) error {
	return db.Where("service_link_id = ?", id).Delete(&ServiceLinkServiceFee{}).Error
}
