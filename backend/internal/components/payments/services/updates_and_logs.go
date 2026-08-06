package payments

import (
	payments "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateAndLogUserPaymentGeoInformation logs payment information and updates user location information
func UpdateAndLogUserPaymentGeoInformation(userInfo *userModels.User, paymentInfoReturned *payments.PaymentInfo, db *gorm.DB) {

	// updatedUser:= &userInfo -m''
	userInfo.AppendGeoInfo()
	db.Omit(clause.Associations).Save(userInfo)
	paymentLog := payments.PaymentLog{
		Sender:               userInfo.Username,
		Destination:          paymentInfoReturned.Destination,
		Memo:                 &paymentInfoReturned.Memo,
		AssetIssuer:          &paymentInfoReturned.AssetIssuer,
		AssetCode:            &paymentInfoReturned.AssetCode,
		Amount:               paymentInfoReturned.Amount,
		Transaction:          &paymentInfoReturned.Transaction,
		TransactionSignature: &paymentInfoReturned.TransactionSignature,
		TransactionID:        paymentInfoReturned.TransactionID,
		DestinationFirstName: &paymentInfoReturned.DestinationFirstName,
		DestinationLastName:  &paymentInfoReturned.DestinationLastName,
		PublicIP:             userInfo.PublicIP,
		CountryCode:          userInfo.CountryCode,
		Latitude:             userInfo.Latitude,
		Longitude:            userInfo.Longitude,
		City:                 userInfo.City,
		Region:               userInfo.Region,
		RegionName:           userInfo.RegionName,
		TimeZone:             userInfo.TimeZone,
		ISP:                  userInfo.ISP,
	}
	db.Omit(clause.Associations).Create(&paymentLog)

}
