package main

import (
	"context"
	"log"
	"testing"
	"time"

	fb "trovo-wallet-api/internal/pns"

	"github.com/shopspring/decimal"
)

const devURL = "http://localhost:8080"
const prodURL = "https://api.trovotechnologies.com"

type ErrorResponse struct {
	Error   string `json:"error"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

type RegSuccessInfo struct {
	PublicKey string `json:"message"`
}
type UserRegistrationInfo struct {
	Username              string `json:"username"`
	Email                 string `json:"email"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Mobile                string `json:"mobile"`
	MobileCountryCode     string `json:"mobileCountryCode,omitempty"`
	PublicKey             string `json:"publicKey,omitempty"`
	Referrer              string `json:"referrer,omitempty"`
	PushNotificationToken string `json:"pushNotificationToken,omitempty"`
	Corporate             uint   `json:"corporate"`
	VerificationCode      string `json:"verificationCode,omitempty"`
}

type UserInfo struct {
	UserData               UserJSON                 `json:"userData"`
	AssetBalances          map[string]AssetBalances `json:"assetBalances"`          //map of wallet public key and the asset balances
	NFTBalances            map[string]NFTBalances   `json:"nftBalances"`            //map of nft wallet and nftBalances
	ThirdPartyWalletAccess []ThirdPartyWalletAccess `json:"thirdPartyWalletAccess"` //shows all the third party access granted to this user
	DefaultAssets          []DefaultAsset           `json:"defaultAssets"`
}
type UserJSON struct {
	ID                    string           `json:"-"`
	Username              string           `json:"username"`
	Email                 string           `json:"email"`
	ImageThumbnailURL     string           `json:"imageThumbnailURL"`
	FirstName             string           `json:"firstName"`
	LastName              string           `json:"lastName"`
	Mobile                string           `json:"mobile"`
	PublicKey             string           `json:"publicKey"`
	Referrer              string           `json:"referrer"`
	ReferralLink          string           `json:"referralLink"`
	ReferralQrCode        string           `json:"referralQrCode"`
	PushNotificationToken string           `json:"pushNotificationToken"`
	Corporate             uint             `json:"corporate"`
	MobileVerified        uint             `json:"mobileVerified"`
	MembershipType        uint             `json:"membershipType"`
	MembershipExpiry      time.Time        `json:"membershipExpiry"`
	KYCVerified           uint             `json:"kycVerified"`
	WalletRecoveryEnabled uint             `json:"walletRecoveryEnabled"`
	UserWallets           []UserWalletJSON `json:"userWallets"`
	Verified              int              `json:"verified"`
	Suspended             int              `json:"suspended"`
}

type UserWalletJSON struct {
	CreatedAt               time.Time                   `json:"createdAt"`
	ID                      string                      `json:"publicKey"`
	TempPublicKey           string                      `json:"-"`
	Tag                     string                      `json:"tag"`
	Description             string                      `json:"description"`
	Alias                   string                      `json:"alias"`  //primaryUsername_tag for sub wallets
	Signer                  string                      `json:"signer"` //if ID is same as signer, then it is a primary wallet
	UserID                  string                      `json:"userId"`
	ManagedAccessEnabled    uint                        `json:"managedAccessEnabled"`
	PrimaryWallet           uint                        `json:"primaryWallet"`
	UserWalletManagedAccess UserWalletManagedAccessJSON `json:"userWalletManagedAccess"`
}

type UserWalletManagedAccessJSON struct {
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	ID                  string             `json:"accessId"`
	UserWalletID        string             `json:"publicKey"`
	NumberOfAuthorizers uint               `json:"numberOfAuthorizers"`
	AccessList          []WalletAccessJSON `json:"accessList"`
}

type WalletAccessJSON struct {
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	Username                  string    `json:"username"`
	AccessLevel               string    `json:"accessLevel"`
	UserWalletManagedAccessID string    `json:"userWalletManagedAccessId"`
}
type ThirdPartyWalletAccess struct {
	Owner             string `json:"owner"`
	PublicKey         string `json:"publicKey"`
	AccessLevel       string `json:"accessLevel"`
	WalletAlias       string `json:"walletAlias"`
	WalletDescription string `json:"walletDescription"`
}

//Balance model for user
type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	QRCode      string          `json:"qrCode"`
}

//Signer model for user
type Signer struct {
	Weight  int    `json:"weight"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Sponsor string `json:"sponsor"`
}

//Signer model for user
type Thresholds struct {
	LowThreshold    string `json:"low_threshold"`
	MediumThreshold string `json:"medium_threshold"`
	HighThreshold   string `json:"high_threshold"`
}

//AssetBalances holds user balances
type AssetBalances struct {
	Claimed   []Balance `json:"claimed"`
	Unclaimed []Balance `json:"unclaimed"`
}

//NFTBalances holds user NFT balances
type NFTBalances struct {
	NFTs []NFT `json:"nfts"`
}
type NFT struct {
	AssetIssuer string `json:"assetIssuer"`
	AssetCode   string `json:"assetCode"`
}
type DefaultAsset struct {
	AssetCode   string `gorm:"size:12" json:"assetCode"`
	AssetIssuer string `gorm:"size:56" json:"assetIssuer"`
}

// func TestAccountRegistration(t *testing.T) {
// 	/*
// 		{"username":"username","email":"richardsric7@gmail.com","firstName":"Kenny","lastName":"Maduka","mobile":"+2347062685682","mobileCountryCode":"NG","referrer":"","pushNotificationToken":"","corporate":0,"verificationCode":""}
// 	*/
// 	userRegInfo := UserRegistrationInfo{
// 		Username:          "ric",
// 		Email:             "richardsric7@gmail.com",
// 		FirstName:         "Ric",
// 		LastName:          "Richards",
// 		Mobile:            "+2348180067955",
// 		MobileCountryCode: "NG",
// 		Referrer:          "",
// 		Corporate:         0,
// 		VerificationCode:  "417936",
// 	}
// 	pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
// 	secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
// 	kp := keypair.MustParseFull(secretKey)
// 	// log.Println(kp.Address())
// 	baseURL := prodURL
// 	fullPath := "/v1/users"
// 	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
// 	ts := time.Now().Unix() / 1000

// 	tsString := fmt.Sprintf("%v", ts)
// 	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
// 	if err != nil {
// 		t.Errorf(err.Error())
// 		return

// 	}

// 	errorResponse := new(ErrorResponse)
// 	regResponse := new(RegSuccessInfo)

// 	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
// 		Set("X-TW-PUBLIC-KEY", kp.Address()).
// 		Set("X-TW-SIGNER", kp.Address()).
// 		Set("X-TW-SIGNATURE", signedHttpHeader).
// 		Set("X-TW-TIMESTAMP", tsString).
// 		Base(baseURL).
// 		Post(fullPath).BodyJSON(userRegInfo).Receive(regResponse, errorResponse)
// 	//get payload string
// 	if len(errorResponse.Error) > 0 {
// 		log.Println("[TestAccountRegistration] server response error:", *errorResponse)
// 		return

// 	}
// 	if err != nil {
// 		log.Println("[TestAccountRegistration]request error:", err)
// 		t.Errorf(err.Error())

// 		return
// 	}

// 	log.Printf("Result:[%+v]\n", regResponse)

// }

// func TestGetUserInfo(t *testing.T) {

// 	pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
// 	secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
// 	kp := keypair.MustParseFull(secretKey)
// 	// log.Println(kp.Address())
// 	// baseURL := "http://localhost:8080"
// 	baseURL := prodURL
// 	fullPath := fmt.Sprintf("/v1/users/%s", "ric")
// 	ts := time.Now().Unix() / 1000

// 	tsString := fmt.Sprintf("%v", ts)
// 	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
// 	if err != nil {
// 		t.Errorf(err.Error())
// 		return

// 	}

// 	errorResponse := new(ErrorResponse)
// 	resultResponse := new(UserInfo)

// 	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
// 		Set("X-TW-PUBLIC-KEY", kp.Address()).
// 		Set("X-TW-SIGNER", kp.Address()).
// 		Set("X-TW-SIGNATURE", signedHttpHeader).
// 		Set("X-TW-TIMESTAMP", tsString).
// 		Base(baseURL).
// 		Get(fullPath).Receive(resultResponse, errorResponse)
// 	//get payload string
// 	if len(errorResponse.Error) > 0 {
// 		log.Println("[TestGetUserInfo] server response error:", *errorResponse)
// 		return

// 	}
// 	if err != nil {
// 		log.Println("[TestGetUserInfo]request error:", err)
// 		t.Errorf(err.Error())

// 		return
// 	}

// 	log.Printf("Result:[%+v]\n", resultResponse)

// }
func TestSendPushNotificationMessage(t *testing.T) {
	ric := "dqmlIz2eT-CTEX1cesc3za:APA91bGDuSBPzGSZAjYfFMU12eV_gpMw7JdF7gxhhsaPc0XGmgYu9RokVyLfEpH6_Yfg9mYrmuhYlcLRScvolUbfRVftR3QK4d6psgVOY0382fNW3Q1BND8kTbPAThBOXe60DHSXO9hD"
	// kennis := "eNCa_XRaTr2NnXX4pnzhN3:APA91bF9OfBO9IEFJcPOtO-83Qu41_7zZ3ef7qC3i5ySPvT8arcQ1gwnRXYnSZ5uJ9mT4uOW7rgPp5F0hTsqvoqQL9oR02fQtiCyco2DVsNBT6JIgqgVHO1ZTPod7ypm-MpSzA95MRRZ"
	title := "TROVO: Testing Push Notification Service"
	body := `This is a test message to ascertain how the push notification appears`
	imageURL := "https://trovotech.io/img/Trovotech-colored.png"
	ctx := context.TODO()
	client, _, err := fb.GetFirebaseMessagingClient(ctx)
	if err != nil {
		log.Println("[TestGetUserInfo]request error:", err)
		t.Errorf(err.Error())
		return
	}
	response, _ := fb.SendFirebaseMessage(ric, title, body, imageURL, client, ctx)

	log.Printf("Result:[%+v]\n", response)

}

// func TestSendPushNotificationBroadcast(t *testing.T) {

// ric := "dqmlIz2eT-CTEX1cesc3za:APA91bGDuSBPzGSZAjYfFMU12eV_gpMw7JdF7gxhhsaPc0XGmgYu9RokVyLfEpH6_Yfg9mYrmuhYlcLRScvolUbfRVftR3QK4d6psgVOY0382fNW3Q1BND8kTbPAThBOXe60DHSXO9hD"
// kennis := "eNCa_XRaTr2NnXX4pnzhN3:APA91bF9OfBO9IEFJcPOtO-83Qu41_7zZ3ef7qC3i5ySPvT8arcQ1gwnRXYnSZ5uJ9mT4uOW7rgPp5F0hTsqvoqQL9oR02fQtiCyco2DVsNBT6JIgqgVHO1ZTPod7ypm-MpSzA95MRRZ"

// receipients := []string{godswill, tunde, onoja, cryptoking, ric, mavol}
// 	title := "TESTING THE PUSH NOTIFICATION BROADCAST"
// 	body := `This is a test message to ascertain how the push notification appears`
// 	imageURL := "https://trovotech.io/img/Trovotech-colored.png"
// 	ctx := context.TODO()
// 	client, _, err := fb.GetFirebaseMessagingClient(ctx)
// 	if err != nil {
// 		log.Println("[TestGetUserInfo]request error:", err)
// 		t.Errorf(err.Error())
// 		return
// 	}
// 	response, _ := fb.SendFirebaseBroadcast(receipients, title, body, imageURL, client, ctx)

// 	log.Printf("Result:[%+v]\n", response)

// }
