package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	p2pErrors "admin-panel-dashboard/internal/errors"
	"strings"

	"github.com/stellar/go/keypair"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentCast string

func (p PaymentCast) GeneratePaymentMethodID(owner, destinationAccount, paymentChannelId string) string {
	return strings.ReplaceAll(owner+destinationAccount+paymentChannelId, " ", "")
}

// CreatePaymentMethod creates a new payment method
func (p PaymentCast) CreatePaymentMethod(owner string, paymentMethodInput PaymentMethodJSON, db *gorm.DB) (paymentMethod PaymentMethod, err error) {

	if strings.Contains(strings.ToLower(paymentMethodInput.PaymentChannelID), "bantu") {
		_, e := keypair.ParseAddress(paymentMethodInput.DestinationAccount)
		if e != nil {
			return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
				Param:      "destinationAccount",
				ErrMessage: "Failed to create payment method due to invalid bantu address " + paymentMethodInput.DestinationAccount,
				Code:       http.StatusBadRequest,
			}
		}
	}
	var memoAddresses []string
	if len(os.Getenv("WALLETS_REQUIRE_28_BYTE_MEMO")) >= 56 {
		memoAddresses = strings.Split(os.Getenv("WALLETS_REQUIRE_28_BYTE_MEMO"), ",")

		for _, m := range memoAddresses {
			if strings.EqualFold(paymentMethodInput.DestinationAccount, m) && len(paymentMethodInput.Memo) < 28 {

				return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
					Param:      "destinationAccount",
					ErrMessage: "This address requires a valid unique memo from the exchange. Please put the correct memo fro your exchange",
					Code:       http.StatusBadRequest,
				}

			}
		}

	}

	if strings.Contains(paymentMethodInput.PaymentChannelID, "TRC-20") && (!strings.HasPrefix(paymentMethodInput.DestinationAccount, "T") || len(paymentMethodInput.DestinationAccount) != 34) {

		return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
			Param:      "destinationAccount",
			ErrMessage: fmt.Sprintf("Wallet address [%v] is not compatible with TRC-20 format.", paymentMethodInput.DestinationAccount),
			Code:       http.StatusBadRequest,
		}

	}

	if (strings.Contains(paymentMethodInput.PaymentChannelID, "ERC-20") || strings.EqualFold(paymentMethodInput.PaymentChannelID, "ETH")) && (!strings.HasPrefix(paymentMethodInput.DestinationAccount, "0x") || len(paymentMethodInput.DestinationAccount) != 42) {

		return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
			Param:      "destinationAccount",
			ErrMessage: fmt.Sprintf("Wallet address [%v] is not compatible with ERC-20 format.", paymentMethodInput.DestinationAccount),
			Code:       http.StatusBadRequest,
		}

	}

	if strings.Contains(paymentMethodInput.PaymentChannelID, "BEP-20") && (!strings.HasPrefix(paymentMethodInput.DestinationAccount, "0x") || len(paymentMethodInput.DestinationAccount) != 42) {

		return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
			Param:      "destinationAccount",
			ErrMessage: fmt.Sprintf("Wallet address [%v] is not compatible with BEP-20 format.", paymentMethodInput.DestinationAccount),
			Code:       http.StatusBadRequest,
		}

	}
	// check for details required by bank transfer and mobile money
	if (strings.EqualFold(paymentMethodInput.PaymentChannelID, "Bank Transfer") || strings.EqualFold(paymentMethodInput.PaymentChannelID, "Mobile Money")) && (len(strings.ReplaceAll(strings.ReplaceAll(paymentMethodInput.BankName, "	", ""), " ", "")) < 2 || len(strings.ReplaceAll(strings.ReplaceAll(paymentMethodInput.Name, "	", ""), " ", "")) < 2) {

		return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
			Param:      "bankName",
			ErrMessage: "Bank Transfer/Mobile Money require both the Beneficiary & Bank/Mobile Network names.",
			Code:       http.StatusBadRequest,
		}

	}

	// check for details required by bank transfer and mobile money and the currency and country code to be linked to it
	if (strings.EqualFold(paymentMethodInput.PaymentChannelID, "Bank Transfer") || strings.EqualFold(paymentMethodInput.PaymentChannelID, "Mobile Money")) && (len(paymentMethodInput.CurrencyID) == 0 || len(paymentMethodInput.CountryCode) == 0) {

		return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
			Param:      "countryCode",
			ErrMessage: "Bank Transfer/Mobile Money require both the Country and Currency to be linked to them. Please select the appropriate ones before proceeding.",
			Code:       http.StatusBadRequest,
		}

	}
	if len(paymentMethodInput.CountryCode) > 2 {

		return paymentMethod, &p2pErrors.ErrorActionNotAllowed{
			Param:      "countryCode",
			ErrMessage: fmt.Sprintf("Only 2 Character Country Code is required. Value submitted [%v] is invalid.", paymentMethodInput.CountryCode),
			Code:       http.StatusBadRequest,
		}

	}
	e := db.Where("username = ?", owner).Where("payment_channel_id = ?", paymentMethodInput.PaymentChannelID).Where("destination_account = ?", paymentMethodInput.DestinationAccount).First(&PaymentMethod{}).Error
	if e == nil {

		return paymentMethod, &p2pErrors.ErrorPaymentMethodAlreadyExists{Detail: fmt.Sprintf("duplicate payment destination account:[%v] %v", paymentMethodInput.PaymentChannelID, paymentMethodInput.DestinationAccount)}
	}

	// build the payment method
	paymentMethod = PaymentMethod{
		ID:                 PaymentCast("d").GeneratePaymentMethodID(owner, paymentMethodInput.DestinationAccount, paymentMethodInput.PaymentChannelID),
		Username:           owner,
		PaymentChannelID:   paymentMethodInput.PaymentChannelID,
		DestinationAccount: paymentMethodInput.DestinationAccount,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if len(paymentMethodInput.Name) > 0 {
		paymentMethod.Name = &paymentMethodInput.Name
	}
	if len(paymentMethodInput.AccountOpeningBranch) > 0 {
		paymentMethod.AccountOpeningBranch = &paymentMethodInput.AccountOpeningBranch
	}
	if len(paymentMethodInput.Memo) > 0 {
		paymentMethod.Memo = &paymentMethodInput.Memo
	}
	if len(paymentMethodInput.BankName) > 0 {
		paymentMethod.BankName = &paymentMethodInput.BankName
	}
	if len(paymentMethodInput.CountryCode) > 0 {
		paymentMethod.CountryCode = &paymentMethodInput.CountryCode
	}
	if len(paymentMethodInput.CurrencyID) > 0 {
		paymentMethod.CurrencyID = &paymentMethodInput.CurrencyID
	}

	e = db.Create(&paymentMethod).Error
	if e != nil {
		log.Println("[CreatePaymentMethod]could not create paymentMethod,err:", e)
		// LogDiscordError(fmt.Sprintf("[CreatePaymentMethod]could not create paymentMethod,err:[%v]", e))
		return PaymentMethod{}, &p2pErrors.ErrorTemporaryServerError{}

	}
	return paymentMethod, nil
}

func (p PaymentCast) DeletePaymentMethod(owner, paymentMethodID string, db *gorm.DB) (paymentMethod PaymentMethod, err error) {
	e := db.Where("username = ?", owner).Where("id = ?", paymentMethodID).First(&paymentMethod).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[DeletePaymentMethod]:", e)
			return PaymentMethod{}, &p2pErrors.ErrorTemporaryServerError{}
		}
		return paymentMethod, &p2pErrors.ErrorPaymentMethodDoesNotExist{ID: paymentMethodID}
	}

	{
		// check if payment method has existing offer using it.
		type Offer struct {
			ID                       string  `gorm:"primaryKey;check:,length(id) > 2" json:"id"`
			OfferType                string  `gorm:"size:10;not null;index:idx_unique_offer,unique" json:"offerType"`
			OrderType                string  `gorm:"size:10;not null;index:idx_offer_order_type" json:"orderType"`
			Maker                    string  `gorm:"not null;index:idx_unique_offer,unique;check:,length(maker) < 50" json:"maker"`
			MakerPhone               string  `gorm:"not null;" json:"makerPhone"`
			PaymentMethodID          *string `gorm:"null;" json:"paymentMethodId"`
			PaymentChannelID         *string `gorm:"null;index:idx_offer_payment_channel_id" json:"paymentChannelId"`
			CurrencyPaymentChannelID string  `gorm:"not null;index:idx_unique_offer,unique;index:idx_offer_currency_payment_channel_id" json:"currencyPaymentChannelId"`
			CurrencyPaymentMethodID  string  `gorm:"not null;index:idx_offer_currency_payment_method_id;" json:"currencyPaymentMethodlId"`
			CurrencyID               string  `gorm:"size:100;not null;index:idx_unique_offer,unique" json:"currencyId"`
			AssetAmount              float64 `gorm:"not null;check:,asset_amount >= 0" json:"assetAmount"`
			AssetID                  string  `gorm:"size:12;not null;index:idx_unique_offer,unique" json:"assetId"`
			AssetPrice               float64 `gorm:"not null;check:,asset_price > 0" json:"assetPrice"`
			MinTradeAmount           float64 `gorm:"not null;check:,min_trade_amount > 0" json:"minTradeAmount"`
			MaxTradeAmount           float64 `gorm:"not null;check:,max_trade_amount > 0" json:"maxTradeAmount"`
			Remark                   *string `gorm:"null" json:"remark"`
			Inactive                 uint    `gorm:"not null;default:0" json:"-"`
			Offline                  uint    `gorm:"not null;default:0" json:"-"`
		}
		var offer Offer
		e = db.Where("payment_method_id = ?", paymentMethodID).First(&offer).Error
		if e != nil {
			if !errors.Is(e, gorm.ErrRecordNotFound) {
				log.Println("[DeletePaymentMethod]:", e)
				return PaymentMethod{}, &p2pErrors.ErrorTemporaryServerError{}
			}
		} else {
			return paymentMethod, &p2pErrors.ErrorOperationFailed{
				Param:      "paymentMethodId",
				ErrMessage: fmt.Sprintf("Cannot delete payment method in use by your offer of %v for %v on %v", offer.CurrencyID, offer.AssetID, offer.PaymentChannelID),
			}
		}
		e = db.Where("currency_payment_method_id = ?", paymentMethodID).First(&offer).Error
		if e != nil {
			if !errors.Is(e, gorm.ErrRecordNotFound) {
				log.Println("[DeletePaymentMethod]:", e)
				return PaymentMethod{}, &p2pErrors.ErrorTemporaryServerError{}
			}
		} else {
			return paymentMethod, &p2pErrors.ErrorOperationFailed{
				Param:      "paymentMethodId",
				ErrMessage: fmt.Sprintf("Cannot delete payment method in use by your offer of %v for %v on %v", offer.CurrencyID, offer.AssetID, offer.CurrencyPaymentChannelID),
			}
		}
	}

	e = db.Delete(&paymentMethod).Error
	if e != nil {
		log.Println("[DeletePaymentMethod]could not delete paymentMethod,err:", e)
		return PaymentMethod{}, &p2pErrors.ErrorTemporaryServerError{}

	}
	return paymentMethod, nil
}
func (p PaymentCast) GetPaymentChannels(activeOnly bool, db *gorm.DB) (paymentChannels []PaymentChannel, err error) {
	paymentChannels = make([]PaymentChannel, 0)
	var e error
	if activeOnly {
		e = db.Preload(clause.Associations).Order("id").Where("inactive = ?", 0).Find(&paymentChannels).Error
	} else {
		e = db.Preload(clause.Associations).Order("id").Find(&paymentChannels).Error
	}

	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetPaymentChannels]:", e)
			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
	}
	return
}
func (p PaymentCast) GetPaymentChannelByID(chanID string, db *gorm.DB) (paymentChannel PaymentChannel, err error) {

	e := db.Preload(clause.Associations).Where("id = ?", chanID).First(&paymentChannel).Error

	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetPaymentChannelByID]:", e)
			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
		return PaymentChannel{}, &p2pErrors.CustomError{
			Param:      "paymentChannelId",
			Err:        "payment channel does not exist",
			ErrMessage: fmt.Sprintf("Payment Channel [%v] Does Not Exist", chanID),
			Code:       404,
		}
	}
	return paymentChannel, nil
}

func (p PaymentCast) GetPaymentMethods(activeOnly bool, db *gorm.DB) (paymentMethods []PaymentMethod, err error) {
	paymentMethods = make([]PaymentMethod, 0)
	var e error
	if activeOnly {
		e = db.Preload(clause.Associations).Order("id").Where("inactive = ?", 0).Find(&paymentMethods).Error
	} else {
		e = db.Preload(clause.Associations).Order("id").Find(&paymentMethods).Error
	}

	if e != nil {

		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetPaymentMethods]:", e)
			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
	}
	return
}

func (p PaymentCast) GetPaymentMethodsByUser(username string, db *gorm.DB) (paymentMethods []PaymentMethod, err error) {
	paymentMethods = make([]PaymentMethod, 0)

	e := db.Preload(clause.Associations).Order("id").Where("username = ?", username).Find(&paymentMethods).Error

	if e != nil {
		// log.Println("[GetPaymentMethodsByUser]err:", e)
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetPaymentMethodsByUser]could not fetch payment methods by user, err:", e)
			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
	}
	return
}
func (p PaymentCast) GetPaymentMethodByUser(id, username string, db *gorm.DB) (paymentMethod PaymentMethod, err error) {

	e := db.Preload(clause.Associations).Where("id = ?", id).Where("username = ?", username).First(&paymentMethod).Error

	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetPaymentMethodByUser]could not get payment method by user, err:", e)
			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
		err = &p2pErrors.ErrorPaymentMethodDoesNotExist{ID: id}
		return
	}
	return
}
