package swaps

import (
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	userModels "trovo-wallet-api/internal/components/users/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateAndLogUserSwapGeoInformation logs payment information and updates user location inforamtion
func UpdateAndLogUserSwapGeoInformation(userInfo *userModels.User, swapInfo *swapModels.SwapSendInfo, db *gorm.DB) {

	// updatedUser:= &userInfo -m''
	userInfo.AppendGeoInfo()
	db.Omit(clause.Associations).Save(userInfo)
	swapLog := swapModels.PaymentLog{
		TrasanctionType:            1,
		Sender:                     userInfo.Username,
		Destination:                userInfo.Username,
		Memo:                       &swapInfo.Memo,
		ContractAddress:            &swapInfo.SourceContractAddress,
		AssetCode:                  &swapInfo.SourceAssetCode,
		Amount:                     swapInfo.SourceAmount,
		DestinationContractAddress: &swapInfo.DestinationContractAddress,
		DestinationAssetCode:       &swapInfo.DestinationAssetCode,
		Transaction:                &swapInfo.Transaction,
		TransactionSignature:       &swapInfo.TransactionSignature,
		TransactionID:              swapInfo.TransactionID,
		DestinationFirstName:       &userInfo.FirstName,
		DestinationLastName:        userInfo.LastName,
		PublicIP:                   userInfo.PublicIP,
		CountryCode:                userInfo.CountryCode,
		Latitude:                   userInfo.Latitude,
		Longitude:                  userInfo.Longitude,
		City:                       userInfo.City,
		Region:                     userInfo.Region,
		RegionName:                 userInfo.RegionName,
		TimeZone:                   userInfo.TimeZone,
		ISP:                        userInfo.ISP,
	}

	db.Omit(clause.Associations).Create(&swapLog)

}
