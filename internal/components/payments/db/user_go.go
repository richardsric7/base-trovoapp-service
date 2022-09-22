package payments

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	CreatedAt             time.Time    `json:"createdAt"`
	UpdatedAt             time.Time    `json:"updatedAt"`
	LastUpdatedMobileOn   time.Time    `json:"lastUpdatedMobileOn"`
	ID                    string       `json:"id"`
	Username              string       `gorm:"size:16; index:idx_user_unique_username, unique" json:"username"`
	Email                 string       `gorm:"size:45; index:idx_user_unique_email, unique" json:"email"`
	ImageThumbnailURL     *string      `json:"imageThumbnailURL"`
	FirstName             string       `gorm:"size:50" json:"firstName"`
	LastName              *string      `gorm:"size:50" json:"lastName"`
	Mobile                *string      `gorm:"size:16; index:idx_user_unique_phone, unique" json:"mobile"`
	PublicKey             string       `gorm:"size:56; index:idx_user_unique_public_key, unique" json:"publicKey"`
	PrimarySigner         string       `gorm:"size:56; index:idx_user_unique_primary_signer, unique" json:"primarySigner"`
	Referrer              *string      `gorm:"size:16; index:idx_user_referrer" json:"referrer"`
	ReferralLink          *string      `json:"referralLink"`
	ReferralQrCode        *string      `json:"referralQrCode"`
	PushNotificationToken *string      `json:"pushNotificationToken"`
	Corporate             uint         `gorm:"type:integer;not null; default:0" json:"corporate"`
	MobileVerified        uint         `gorm:"type:integer;not null; default:0" json:"mobileVerified"`
	MembershipType        uint         `gorm:"type:integer;not null; default:0" json:"membershipType"`
	MembershipExpiry      *time.Time   `json:"membershipExpiry"`
	KYCVerified           uint         `gorm:"type:integer;not null; default:0" json:"kycVerified"`
	AccountRecoveryEnabled uint         `gorm:"type:integer;not null; default:0" json:"accountRecoveryEnabled"`
	UserWallets           []UserWallet `json:"userWallets"`
	PublicIP              string       `gorm:"size:45" json:"publicIP"`
	CountryCode           *string      `gorm:"size:2;null"`
	Latitude              *float64     `gorm:"null"`
	Longitude             *float64     `gorm:"null"`
	City                  *string      `gorm:"null;size:100"`
	Region                *string      `gorm:"null;size:100"`
	RegionName            *string      `gorm:"null;size:100"`
	TimeZone              *string      `gorm:"null;size:100"`
	ISP                   *string      `gorm:"null;size:150"`
	Verified              int          `gorm:"type:integer;not null;default:0" json:"verified"`
	Suspended             int          `gorm:"type:integer;not null;default:0" json:"suspended"`
	SuspensionReason      *string      `gorm:"null" json:"suspensionReason"`
}

type UserWallet struct {
	CreatedAt               time.Time               `json:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt"`
	ID                      string                  `gorm:"size:56" json:"publicKey"`
	TempPublicKey           *string                 `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag                     *string                 `gorm:"null;size:16" json:"tag"`
	Description             *string                 `gorm:"null;size:100" json:"description"`
	Alias                   string                  `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer                  string                  `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  //if ID is same as signer, then it is a primary wallet
	UserID                  string                  `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	ManagedAccessEnabled    uint                    `gorm:"type:integer;not null; default:0" json:"managedAccessEnabled"`
	UserWalletManagedAccess UserWalletManagedAccess `json:"userWalletManagedAccess"`
}

type UserWalletManagedAccess struct {
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	ID                  string         `gorm:"" json:"accessId"`
	UserWalletID        string         `gorm:"size:56; index:idx_manage_access_user_wallet_id" json:"publicKey"`
	NumberOfAuthorizers uint           `gorm:"type:integer; default:1" json:"numberOfAuthorizers"`
	AccessList          []WalletAccess `json:"accessList"`
}
type WalletAccess struct {
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	Username                  string    `gorm:"size:16; primaryKey" json:"username"`
	AccessLevel               string    `gorm:"size:10" json:"accessLevel"`
	UserWalletManagedAccessID string    `gorm:"index:idx_wallet_access_wallet_access_id" json:"userWalletManagedAccessId"`
}

type AccessLevel struct {
	ID          uint64
	AccessLevel string `gorm:"size:text" json:"accessList"`
}

// ReservedName holds model struct for ReservedName table
type ReservedName struct {
	ID           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ReservedName *string `gorm:"size:50;not null;index:unique_reserved_name, unique;index:idx_reserved_status"`
	Status       *uint64 `gorm:"default:0;index:idx_reserved_status"`
}

// Signer model for user
type Signer struct {
	Weight  int    `json:"weight"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Sponsor string `json:"sponsor"`
}

// Signer model for user
type Thresholds struct {
	LowThreshold    string `json:"low_threshold"`
	MediumThreshold string `json:"medium_threshold"`
	HighThreshold   string `json:"high_threshold"`
}

// GetBlockchainAccountDetail fetches the bantu account information using public key
func (u *UserWallet) GetBlockchainAccountDetail(temp bool) (clientAccount horizon.Account, destinationAccountExists bool, err error) {
	client := network.GetBlockchainClient()
	var accountRequest horizonclient.AccountRequest
	if temp {
		//temp account
		if u.TempPublicKey != nil {
			//temp account has been generated
			accountRequest = horizonclient.AccountRequest{AccountID: *u.TempPublicKey}

		} else {
			//temp account not yet generated
			err = &tErrors.ErrorBlockchainAccountNotActivated{}
			return
		}

	} else {
		//real account
		accountRequest = horizonclient.AccountRequest{AccountID: u.ID}
	}

	clientAccount, err = client.AccountDetail(accountRequest)
	if err != nil {
		// log.Printf("[GetBlockchainAccountDetail]: %v, error: [%v]", accountRequest.AccountID, err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[GetBlockchainAccountDetail Network Failure]: %s\n", "Error Connecting to Expansion Service")
			return clientAccount, destinationAccountExists, &tErrors.ErrorTemporaryServerError{}
		} else {
			horizonException, ok := err.(*horizonclient.Error)

			if ok {

				if horizonException.Problem.Status == http.StatusNotFound {
					return clientAccount, false, &tErrors.ErrorBlockchainAccountNotActivated{}
				}
				log.Println("[BlockchainAccountProperties] error is known", horizonException.Problem.Status)
			}

		}
		return clientAccount, destinationAccountExists, &tErrors.ErrorTemporaryServerError{}
	}
	return clientAccount, true, nil
}

// GetSigners returns user signers
func (u *UserWallet) GetSigners(temp bool) (signers map[string]Signer) {
	signers = make(map[string]Signer)
	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return signers
	}

	for _, v := range account.Signers {
		signers[v.Key] = Signer{
			Weight:  int(v.Weight),
			Key:     v.Key,
			Type:    v.Type,
			Sponsor: v.Sponsor,
		}
	}
	return
}

// GetSignersWA returns user signers
func (u *User) GetSignersWA(account *horizon.Account) (signers map[string]Signer) {
	signers = make(map[string]Signer)
	if account == nil {
		return signers
	}

	for _, v := range account.Signers {
		signers[v.Key] = Signer{
			Weight:  int(v.Weight),
			Key:     v.Key,
			Type:    v.Type,
			Sponsor: v.Sponsor,
		}
	}
	return
}

// GetSignersWA returns user signers
func (u *UserWallet) GetSignersWA(account *horizon.Account) (signers map[string]Signer) {

	signers = make(map[string]Signer)
	if account == nil {
		return signers
	}
	for _, v := range account.Signers {
		signers[v.Key] = Signer{
			Weight:  int(v.Weight),
			Key:     v.Key,
			Type:    v.Type,
			Sponsor: v.Sponsor,
		}
	}
	return
}

// SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *User) SignerIsValidWA(signerKey string, account *horizon.Account) bool {
	if account == nil {
		return false
	}
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValidWA(signerKey string, account *horizon.Account) bool {
	if account == nil {
		return false
	}
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValid checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValid(signerKey string, temp bool) bool {
	signers := u.GetSigners(temp)
	if signers == nil {
		return false
	}
	signer, ok := signers[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// // SignerIsValid checks if the signerKey is valid for this user public key
// func (u *User) SignerIsValid(signerKey string, temp bool) bool {
// 	for _, w := range u.UserWallets {
// 		if w.ID == w.Signer {
// 			signer, ok := w.GetSigners(temp)[signerKey]
// 			if !ok || signer.Weight < 1 {
// 				return false
// 			}

// 			return true
// 		}
// 	}
// 	return false
// }

func GetUser(userInfo string, db *gorm.DB) (user User, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(userInfo) == 56 {
		//56 char public key is supplied

		subQuery := db.Table("user_wallets").Where("id = ?", userInfo).Or("temp_public_key = ?", &userInfo).Or("signer = ?", userInfo).Select("user_id")
		e = db.Preload(clause.Associations).Where("id IN (?)", subQuery).First(&user).Error
	} else if strings.Contains(userInfo, "_") {
		//alias format is supplied
		subQuery := db.Table("user_wallets").Where("alias = ?", strings.ToLower(userInfo)).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error

	} else {
		//search by ID and phone number, username, email

		e = db.Preload(clause.Associations).Where("id = ?", userInfo).Or("username = ?", strings.ToLower(userInfo)).Or("mobile = ?", &userInfo).Or("email = ?", userInfo).First(&user).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorUserDoesNotExist{Username: userInfo}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}

// GetWallet gets user wallet data by alias or public key or temp public key
func GetWallet(identifier string, db *gorm.DB) (userWallet UserWallet, temp bool, err error) {
	conDB.PrintDBStats("[payments]GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(identifier) == 56 {
		//56 char public key is supplied
		e = db.Preload(clause.Associations).Where("id = ?", identifier).Or("temp_public_key = ?", &identifier).First(&userWallet).Error
		if e == nil {
			if identifier == *userWallet.TempPublicKey {
				temp = true
			}
			return
		}
	} else {
		//username is supplied
		e = db.Preload(clause.Associations).Where("alias = ?", strings.ToLower(identifier)).First(&userWallet).Error

	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-does-not-exist", ErrMessage: fmt.Sprintf("%v is not assigned to any wallet", identifier)}
			return
		}
		log.Println("[payments]GetUserInfo error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return userWallet, temp, nil

}

// UsernameIsReserved check is name is reserved. Status = 0 means not available (reserved). Status = 1 means available
func UsernameIsReserved(username string, db *gorm.DB) (reserved bool, err error) {

	username = strings.TrimSpace(username)
	var reservedName ReservedName
	if err := db.Where("reserved_name = ? AND status = 0", strings.ToLower(strings.ReplaceAll(username, " ", ""))).First(&reservedName).Error; err != nil {

		return false, nil
	}

	return true, &tErrors.ErrorUsernameIsReserved{}
}

// PublicKeyIAlreadyExists check if public key already exists
func PublicKeyAlreadyExists(publicKey string, db *gorm.DB) (exists bool, err error) {
	// discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	// if len(os.Getenv("IMPORT_ERROR_WEBHOOK")) > 50 {
	// 	discord.WebhookURL = os.Getenv("IMPORT_ERROR_WEBHOOK")
	// }
	publicKey = strings.TrimSpace(publicKey)
	var userWallet UserWallet
	if err := db.Where("id = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&userWallet).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, &tErrors.ErrorTemporaryServerError{}
		}
		return false, nil
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return true, &tErrors.CustomError{Param: "publicKey", Err: "error-public-key-already-exists", ErrMessage: fmt.Sprintf("Bantu Address [%v] already exists with another active account", userWallet.ID)}

}
func (u *User) SendPushMessage(title, body, imageURI string, dataPayload map[string]string, gc *sharedconfig.GlobalConfig) {
	//Send push notification to user
	// log.Println(title, body)

	if u.PushNotificationToken == nil {
		return
	}

	pns.SendFirebaseMessage(*u.PushNotificationToken, title, body, imageURI, dataPayload, gc.PushNotificationClient, gc.PNSContext)

}

// GetUserFromPrimarySigner fetches the user linked to the primary signer
func GetUserFromPrimarySigner(publicKey string, db *gorm.DB) (user User, err error) {

	publicKey = strings.TrimSpace(publicKey)
	// var user usermodels.User
	if err := db.Where("primary_signer = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return user, &tErrors.ErrorTemporaryServerError{}
		}
		return user, &tErrors.CustomError{Param: "primarySigner",
			Err:        "error primary signer does not exist",
			ErrMessage: "Primary Signer does not exist",
			Code:       http.StatusNotFound,
		}
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return user, nil

}
