package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"trovo-wallet-api/internal/middleware"

	"github.com/dghubble/sling"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
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
type PaymentInfo struct {
	Destination             string            `json:"destination"`
	Memo                    string            `json:"memo"`
	AssetIssuer             string            `json:"assetIssuer"`
	AssetCode               string            `json:"assetCode"`
	Amount                  string            `json:"amount"`
	Transaction             string            `json:"transaction"`
	TransactionSignature    string            `json:"transactionSignature"`
	TransactionID           string            `json:"transactionId"`
	NetworkPassPhrase       string            `json:"networkPassPhrase"`
	DestinationFirstName    string            `json:"destinationFirstName"`
	DestinationLastName     string            `json:"destinationLastName"`
	DestinationThumbnail    string            `json:"destinationThumbnail"`
	DestinationVerified     int               `json:"destinationVerified"`
	ChannelAccount          string            `json:"channelAccount"`
	ChannelAccountSignature string            `json:"channelAccountSignature"`
	Messages                []string          `json:"messages"`
	CallbackURLS            map[string]string `json:"-"`
}
type SubWalletInfo struct {
	PublicKey               string   `json:"publicKey"`
	WalletTag               string   `json:"walletTag"`
	WalletDescription       string   `json:"walletDescription"`
	Transaction             string   `json:"transaction"`
	PrimarySignature        string   `json:"primarySignature"`
	SubWalletSignature      string   `json:"subWalletSignature"`
	TransactionID           string   `json:"transactionId"`
	NetworkPassPhrase       string   `json:"networkPassPhrase"`
	ChannelAccount          string   `json:"channelAccount"`
	ChannelAccountSignature string   `json:"channelAccountSignature"`
	SubWalletMustSign       int      `json:"subWalletMustSign"`
	Messages                []string `json:"messages"`
}
type UserInfo struct {
	UserData               UserJSON                 `json:"userData"`
	AssetBalances          map[string]AssetBalances `json:"assetBalances"`          //map of wallet public key and the asset balances
	NFTs                   map[string][]NFT         `json:"nfts"`                   //map of nft wallet and nftBalances
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
	PrimarySigner         string           `json:"primarySigner"`
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

// Balance model for user
type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	QRCode      string          `json:"qrCode"`
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

// AssetBalances holds user balances
type AssetBalances struct {
	Claimed   []Balance `json:"claimed"`
	Unclaimed []Balance `json:"unclaimed"`
}

// NFTBalances holds user NFT balances
type NFTBalances struct {
	NFTs []NFT `json:"nfts"`
}
type NFT struct {
	AssetIssuer    string `json:"assetIssuer"`
	AssetCode      string `json:"assetCode"`
	NFTName        string `json:"nftName"`
	NFTDescription string `json:"nftDescription"`
	NFTImageURI    string `json:"nftImageURI"`
}
type DefaultAsset struct {
	AssetCode   string `gorm:"size:12" json:"assetCode"`
	AssetIssuer string `gorm:"size:56" json:"assetIssuer"`
}
type PaymentHistoryJSON struct {
	TransactionDate time.Time `json:"transactionDate"`
	TransactionType string    `json:"transactionType"`
	From            string    `json:"from"` //trovoWallet alias and name
	FromPublicKey   string    `json:"fromPublicKey"`
	To              string    `json:"to"` //trovoWallet alias and name
	ToPublicKey     string    `json:"toPublicKey"`
	Memo            string    `json:"memo"`
	AssetIssuer     string    `json:"assetIssuer"`
	AssetCode       string    `json:"assetCode"`
	Amount          string    `json:"amount"`
	TransactionID   string    `json:"transactionId"`
}

type PaginatedPaymentHistory struct {
	Pages        int                  `json:"pages"`
	CurrentPage  int                  `json:"currentPage"`
	TotalRecords int                  `json:"totalRecords"`
	Limit        int                  `json:"limit"`
	Records      []PaymentHistoryJSON `json:"records"`
}

type Trustline struct {
	AssetCode            string `json:"assetCode"`
	AssetIssuer          string `json:"assetIssuer"`
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	TransactionID        string `json:"transactionId"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}

type SwapSendInfo struct {
	DestinationAssetCode   string   `json:"destinationAssetCode"`
	DestinationAssetIssuer string   `json:"destinationAssetIssuer"`
	SwappedEstimate        string   `json:"swappedEstimate"`
	SourceAssetCode        string   `json:"sourceAssetCode"`
	SourceAssetIssuer      string   `json:"sourceAssetIssuer"`
	SourceAmount           string   `json:"sourceAmount" `
	Transaction            string   `json:"transaction"`
	TransactionSignature   string   `json:"transactionSignature"`
	TransactionID          string   `json:"transactionId"`
	NetworkPassPhrase      string   `json:"networkPassPhrase"`
	Messages               []string `json:"messages"`
	Memo                   string   `json:"memo"`
}

// PendingAssetToClaim holds pensing assets to be claimed
type PendingAssetToClaim struct {
	AssetCode            string `json:"assetCode"`
	AssetIssuer          string `json:"assetIssuer"`
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	TransactionID        string `json:"transactionId"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}

func TestAccountRegistration(t *testing.T) {
	/*
		{"username":"username","email":"richardsric7@gmail.com","firstName":"Kenny","lastName":"Maduka","mobile":"+2347062685682","mobileCountryCode":"NG","referrer":"","pushNotificationToken":"","corporate":0,"verificationCode":""}
	*/
	userRegInfo := UserRegistrationInfo{
		Username:          "ric",
		Email:             "richardsric7@gmail.com",
		FirstName:         "Ric",
		LastName:          "Richards",
		Mobile:            "+2348180067955",
		MobileCountryCode: "NG",
		Referrer:          "",
		Corporate:         0,
		VerificationCode:  "253988",
	}
	// pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	// secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	primaryPK := os.Getenv("RICPK")
	primarySecretKey := os.Getenv("RICSC")
	if len(primarySecretKey) == 0 || len(primaryPK) == 0 {
		log.Println("No primary secret or public key specified")
		time.Sleep(20 * time.Second)
		t.Errorf("No primary secret or public key specified")
		return
	}
	kp := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	baseURL := devURL
	fullPath := "/v1/users"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	regResponse := new(RegSuccessInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(userRegInfo).Receive(regResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestAccountRegistration] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestAccountRegistration]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Result:[%+v]\n", regResponse)

}

func TestGetUserInfo(t *testing.T) {

	pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	// baseURL := "http://localhost:8080"
	baseURL := devURL
	fullPath := fmt.Sprintf("/v1/users/%s", "ric")
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	resultResponse := new(UserInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Get(fullPath).Receive(resultResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestGetUserInfo] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestGetUserInfo]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Result:[%+v]\n", resultResponse)

}

/*
*
2022/06/30 18:40:23 Full Path With Query:[/v1/users/payments/GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC] KeyParam:[GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC1656614] Signature: [dBVVuSejU7jfc7QPLU4pMuAuDa00S/3EPyasS+ime9gGFFOONwogY1otWVduUu0VnirdTIj1jXCLzsmblH3qBA==]
*
*/
func TestGetPaymentHistory(t *testing.T) {

	// pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	// secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	// baseURL := "http://localhost:8080"
	baseURL := prodURL
	// baseURL := prodURL
	fullPath := fmt.Sprintf("/v1/users/payments/%v", kp.Address())
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	resultResponse := new(PaginatedPaymentHistory)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Get(fullPath).Receive(resultResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestGetUserInfo] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestGetUserInfo]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Result:[%+v]\n", resultResponse)

}

// func TestSendPushNotificationMessage(t *testing.T) {
// 	ric := "dWLRIQWuSm-sqAKA-ABwhS:APA91bGJt8PE4KBS0OPIJOVp4JsWbpiJMK1DrIJIhZM7hlOVeLo6OUGlN5PbbsttT3Oq0YNXZZ8P0zDEcVD6wQkdoFOfrTSFq0q9A1XZU605ZhTpaLLfqwqONniRQKoj4I-YdMgxXoJc"
// 	// kennis := "eNCa_XRaTr2NnXX4pnzhN3:APA91bF9OfBO9IEFJcPOtO-83Qu41_7zZ3ef7qC3i5ySPvT8arcQ1gwnRXYnSZ5uJ9mT4uOW7rgPp5F0hTsqvoqQL9oR02fQtiCyco2DVsNBT6JIgqgVHO1ZTPod7ypm-MpSzA95MRRZ"
// 	title := "TROVO: Testing Push Notification Service"
// 	body := `This is a test message to ascertain how the push notification appears`
// 	imageURL := "https://trovotech.io/img/Trovotech-colored.png"
// 	dataPayload := make(map[string]string)
// 	dataPayload["route"] = "announcements"
// 	dataPayload["openLink"] = "https://trovowallet.page.link"
// 	ctx := context.Background()
// 	client, _, err := fb.GetFirebaseMessagingClient(ctx)
// 	if err != nil {
// 		log.Println("[TestGetUserInfo]request error:", err)
// 		t.Errorf(err.Error())
// 		return
// 	}
// 	response, _ := fb.SendFirebaseMessage(ric, title, body, imageURL, dataPayload, client, ctx)
// 	// response, _ := fb.SendFirebaseMessage(ric, title, body, imageURL, nil, nil)

// 	log.Printf("Result:[%+v]\n", response)

// }

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

/**
2022/06/30 18:37:43 Full Path With Query:[/v1/users/payment] KeyParam:[GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC1656614] Signature: [WnEUUYVAjhVRSvqJ/ndOf8TMlIil45Gua2Cn2PGmZULrb5RGVcZuOxYskp0BiHWqGKpRpRi+6Kfhf/DA4jrpCA==]
**/

// TestSendPaymentMultiAccessDisabled sends payment from primary wallet
func TestAcceptAssetMultiAccessDisabled(t *testing.T) {

	fromWallet := "GDW6UKK6RI2LBTGHTDKKXYZKCGPDFBRFDTYSZKGGGE6SC5TCSG3MMJST"
	pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")

	kp := keypair.MustParseFull(secretKey)
	baseURL := prodURL

	fullPath := "/v1/users/actions/claim-asset"
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	claimPayload := PendingAssetToClaim{
		AssetCode:   "YAM",
		AssetIssuer: "GAJ65QHSOIXOA6FZMKDIBNGMHXQ7U46TNRBKDL3MTHERF2VRVMWU2F57",
	}
	errorResponse := new(ErrorResponse)
	claimResponse := new(PendingAssetToClaim)
	log.Println("making request from wallet", fromWallet)
	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", fromWallet).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Put(fullPath).BodyJSON(claimPayload).Receive(claimResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestAcceptAssetMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestAcceptAssetMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Claim Response:[%+v]\n", claimResponse)

	{
		//run the payment signing and submission
		p := *claimResponse
		//sign transaction

		signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestAcceptAssetMultiAccessDisabled] make claim error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = signedBase64

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		log.Printf("Second Claim Payload:[%+v]\n", p)
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", fromWallet).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Put(fullPath).BodyJSON(p).Receive(claimResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Make Claim Response:[%+v]\n", claimResponse)
	}
	log.Println("[TestAcceptAssetMultiAccessDisabled] completed")
	time.Sleep(time.Second * 10)

}
func TestSendPaymentMultiAccessDisabled(t *testing.T) {

	// pk := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// secretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/payment"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	// paymentPayload := PaymentInfo{
	// 	Destination: "obi",
	// 	Memo:        "Test Payment",
	// 	Amount:      "200",
	// 	AssetCode:   "YAM",
	// 	AssetIssuer: "GAJ65QHSOIXOA6FZMKDIBNGMHXQ7U46TNRBKDL3MTHERF2VRVMWU2F57",
	// }

	// paymentPayload := PaymentInfo{
	// 	Destination: "onoja",
	// 	Memo:        "Test Payment",
	// 	Amount:      "200",
	// 	AssetCode:   "ABC",
	// 	AssetIssuer: "GAD3DZNQY4SXJEUJOPLJZEK3OWTASEUK2LZYT3V7C52UN5QYOFP3PM5P",
	// }

	paymentPayload := PaymentInfo{
		Destination: "obi",
		Memo:        "Test Payment",
		Amount:      "120",
		AssetCode:   "LUMI",
		AssetIssuer: "GBGUHXVAK32BZTWRBML7RNIQ3532QDR5RRXOJ2P2MGEPHC3YJREMRTLE",
	}
	errorResponse := new(ErrorResponse)
	payResponse := new(PaymentInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(paymentPayload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestSendPaymentMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestSendPaymentMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Payment Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		//sign transaction
		if len(p.ChannelAccount) == 56 {
			ckp := keypair.MustParseFull(channelAccountSK)

			dsigned, err := middleware.SignBase64Txn(ckp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestSendPaymentMultiAccessDisabled]request error:", err)
				t.Errorf(err.Error())

				return
			}
			p.ChannelAccountSignature = dsigned

		}
		signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestSendPaymentMultiAccessDisabled] makePayment error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = signedBase64

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", kp.Address()).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Make Payment Response:[%+v]\n", payResponse)
	}
	log.Println("[TestSendPaymentMultiAccessDisabled] completed")
	time.Sleep(time.Second * 10)

}

// TestSendPaymentMultiAccessDisabled sends payment from primary account
func TestSendPaymentFromSubWalletMultiAccessDisabled(t *testing.T) {

	// pk := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// secretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	// pk := os.Getenv("RICPK")
	fromWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	signerSecretKey := os.Getenv("RICSC")
	channelAccountSK := ""
	// ownerUsername := "ric"
	signerKP := keypair.MustParseFull(signerSecretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/payment"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, signerKP.Address()+tsString, signerKP.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	paymentPayload := PaymentInfo{
		Destination: "obi",
		Memo:        "Test XBN Payment",
		Amount:      "51",
	}
	errorResponse := new(ErrorResponse)
	payResponse := new(PaymentInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", fromWallet).
		Set("X-TW-SIGNER", signerKP.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(paymentPayload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestSendPaymentFromSubWalletMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestSendPaymentFromSubWalletMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestSendPaymentFromSubWalletMultiAccessDisabled]Confirmation Payment Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		//sign transaction
		if len(p.ChannelAccount) == 56 {
			ckp := keypair.MustParseFull(channelAccountSK)

			dsigned, err := middleware.SignBase64Txn(ckp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestSendPaymentFromSubWalletMultiAccessDisabled]request error:", err)
				t.Errorf(err.Error())

				return
			}
			p.ChannelAccountSignature = dsigned

		}
		signedBase64, err := middleware.SignBase64Txn(signerKP.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestSendPaymentFromSubWalletMultiAccessDisabled] makePayment error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = signedBase64

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, signerKP.Address()+tsString, signerKP.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", fromWallet).
			Set("X-TW-SIGNER", signerKP.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Make Payment Response:[%+v]\n", payResponse)
	}
	log.Println("[TestSendPaymentFromSubWalletMultiAccessDisabled] completed")
	time.Sleep(time.Second * 10)

}
func TestSwapFromSubWalletMultiAccessDisabled(t *testing.T) {

	// pk := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// secretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	// pk := os.Getenv("RICPK")
	fromWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	signerSecretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	signerKP := keypair.MustParseFull(signerSecretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	// fullPath := fmt.Sprintf("/v1/users/swap", sEnc)
	fullPath := "/v1/users/swap"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, signerKP.Address()+tsString, signerKP.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	swapPayload := SwapSendInfo{
		DestinationAssetCode:   "YAM",
		DestinationAssetIssuer: "GAJ65QHSOIXOA6FZMKDIBNGMHXQ7U46TNRBKDL3MTHERF2VRVMWU2F57",
		SourceAssetCode:        "",
		SourceAssetIssuer:      "",
		SourceAmount:           "1700",
	}
	errorResponse := new(ErrorResponse)
	swapResponse := new(SwapSendInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", fromWallet).
		Set("X-TW-SIGNER", signerKP.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(swapPayload).Receive(swapResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestSwapFromSubWalletMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestSwapFromSubWalletMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestSwapFromSubWalletMultiAccessDisabled]Confirmation SWAP Response:[%+v]\n", swapResponse)

	{
		//run the payment signing and submission
		p := *swapResponse
		//sign transaction

		signedBase64, err := middleware.SignBase64Txn(signerKP.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestSwapFromSubWalletMultiAccessDisabled] makePayment error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = signedBase64

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, signerKP.Address()+tsString, signerKP.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", fromWallet).
			Set("X-TW-SIGNER", signerKP.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(swapResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("DO SWAP Response:[%+v]\n", swapResponse)
	}
	log.Println("[TestSwapFromSubWalletMultiAccessDisabled] completed")
	time.Sleep(time.Second * 10)

}

func TestCreateSubWalletMultiAccessDisabled(t *testing.T) {

	// subPK := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// subSecretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	// subPK := "GDW6UKK6RI2LBTGHTDKKXYZKCGPDFBRFDTYSZKGGGE6SC5TCSG3MMJST"
	// subSecretKey := "SBOMXAYMOQ64KJSYGIMLLDJBC5DOVGCOWUDVC6ENB4Z642MJNSQQP5HY"
	// subPK := "GAFR2PQHE6GBGCTWN7AAC6WBVOICBZFE35DKRYJKQUHDZFEODRAGRAU4"
	// subSecretKey := "SDOSD4PD6RE7PG2TSGIGVRUXCSLQBTSJH3FBBLODPHEQNTBNFR3CIF2F"
	// subPK := "GBDC4XVY2BVLHOCCRC5A6PJO4GFZ3JT4K3QD65HJQM655QTQ6UHU2GSP"
	// subSecretKey := "SAHXHVXR63DS3LDOYXBO3ENXDK7AFL67NDBYDXHXLMZQ5DHSVQ75C5HU"
	subPK := "GBPAFGAOQ65ODZEH7XKLGVFAHUPPADALTK3AS2ISSI65NXDXH6CXAKA7"
	subSecretKey := "SCCNAU4JKDVJP5O7Y4RJDLGARCNTQVUBK5L5OXI5VOIZYHPBJLESZJML"
	primaryPK := os.Getenv("RICPK")
	primarySecretKey := os.Getenv("RICSC")
	// primaryPK := "GD36GHMT65T2O5YOSFE57TLF4VTSI67IAQXSUT4L5SNBKPNMV5R5R6VV"
	// primarySecretKey := "SAWWK6BIPRALRRHVHELHI2Q3U66KBZLTLOPE7DVGJYRKZYFZGBCZZALY"
	channelAccountSK := ""
	// ownerUsername := "ric"
	// subKP := keypair.MustParseFull(subSecretKey)
	// primaryKP := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// baseURL := devURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/subwallet"
	// fullPath := fmt.Sprintf("/v1/users/%v/payments", sEnc)
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, primarySecretKey)
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	subwalletPayload := SubWalletInfo{
		PublicKey:         subPK,
		WalletTag:         "sub4",
		WalletDescription: "Sub wallet Four",
	}
	errorResponse := new(ErrorResponse)
	subWalletResponse := new(SubWalletInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", primaryPK).
		Set("X-TW-SIGNER", primaryPK).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(subwalletPayload).Receive(subWalletResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateSubWalletMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestCreateSubWalletMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}
	if len(subWalletResponse.Transaction) == 0 {
		log.Println("[TestCreateSubWalletMultiAccessDisabled]no transaction generated")
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestCreateSubWalletMultiAccessDisabled] Request Subwallet Response:[%+v]\n", subWalletResponse)
	log.Println("==========waiting for 15seconds to before confirmation==============")
	time.Sleep(time.Second * 15)
	log.Println("==========Confirming Subwallet creation==============")

	{
		//run the subwallet signing and submission
		p := *subWalletResponse
		//sign transaction
		if len(p.ChannelAccount) == 56 {
			ckp := keypair.MustParseFull(channelAccountSK)

			dsigned, err := middleware.SignBase64Txn(ckp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestCreateSubWalletMultiAccessDisabled]request error:", err)
				t.Errorf(err.Error())

				return
			}
			p.ChannelAccountSignature = dsigned

		}
		primarySignature, subwalletSignature, err := middleware.SignSubwalletBase64Txn(primarySecretKey, subSecretKey, p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestCreateSubWalletMultiAccessDisabled] sub transactions error:", err)
			t.Errorf(err.Error())

			return
		}

		p.PrimarySignature = primarySignature
		p.SubWalletSignature = subwalletSignature

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, primarySecretKey)
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", primaryPK).
			Set("X-TW-SIGNER", primaryPK).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(subWalletResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		if len(subWalletResponse.TransactionID) == 0 {
			log.Println("[TestCreateSubWalletMultiAccessDisabled]no transaction ID")
			t.Errorf(err.Error())
			return

		}

		log.Printf("Create Subwallet Response:[%+v]\n", subWalletResponse)
	}
	log.Println("COMPLETED TEST: TestCreateSubWalletMultiAccessDisabled")
	time.Sleep(time.Second * 10)

}

func TestTrustAssetMultiAccessDisabled(t *testing.T) {

	subPK := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// subSecretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	// subPK := "GDW6UKK6RI2LBTGHTDKKXYZKCGPDFBRFDTYSZKGGGE6SC5TCSG3MMJST"
	// subSecretKey := "SBOMXAYMOQ64KJSYGIMLLDJBC5DOVGCOWUDVC6ENB4Z642MJNSQQP5HY"
	primaryPK := os.Getenv("RICPK")
	primarySecretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	// subKP := keypair.MustParseFull(subSecretKey)
	// primaryKP := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// baseURL := devURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/trust-asset"
	// fullPath := fmt.Sprintf("/v1/users/%v/payments", sEnc)
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, primarySecretKey)
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	trustLinePayload := Trustline{
		AssetCode:   "TROV",
		AssetIssuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
	}
	errorResponse := new(ErrorResponse)
	trustLineResponse := new(Trustline)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", subPK).
		Set("X-TW-SIGNER", primaryPK).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(trustLinePayload).Receive(trustLineResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestTrustAssetMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestTrustAssetMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}
	if len(trustLineResponse.Transaction) == 0 {
		log.Println("[TestTrustAssetMultiAccessDisabled]no transaction generated")
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestTrustAssetMultiAccessDisabled] Request Subwallet Response:[%+v]\n", trustLineResponse)
	log.Println("==========waiting for 15seconds before confirmation==============")
	time.Sleep(time.Second * 15)
	log.Println("==========Confirming TRUSTLINE creation==============")

	{
		//run the subwallet signing and submission
		p := *trustLineResponse
		//sign transaction

		primarySignature, err := middleware.SignBase64Txn(primarySecretKey, p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestTrustAssetMultiAccessDisabled] sub transactions error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = primarySignature

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, primarySecretKey)
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", subPK).
			Set("X-TW-SIGNER", primaryPK).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(trustLineResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		if len(trustLineResponse.TransactionID) == 0 {
			log.Println("[TestTrustAssetMultiAccessDisabled]no transaction ID")
			t.Errorf(err.Error())
			return

		}

		log.Printf("Create TRUSTLINE Response:[%+v]\n", trustLineResponse)
	}

}

func TestRemoveTrustAssetMultiAccessDisabled(t *testing.T) {

	// subPK := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// subSecretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	subPK := "GDW6UKK6RI2LBTGHTDKKXYZKCGPDFBRFDTYSZKGGGE6SC5TCSG3MMJST"
	// subSecretKey := "SBOMXAYMOQ64KJSYGIMLLDJBC5DOVGCOWUDVC6ENB4Z642MJNSQQP5HY"
	primaryPK := os.Getenv("RICPK")
	primarySecretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	// subKP := keypair.MustParseFull(subSecretKey)
	// primaryKP := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// baseURL := devURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/remove-asset"
	// fullPath := fmt.Sprintf("/v1/users/%v/payments", sEnc)
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, primarySecretKey)
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	trustLinePayload := Trustline{
		AssetCode:   "TROV",
		AssetIssuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
	}
	errorResponse := new(ErrorResponse)
	trustLineResponse := new(Trustline)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", subPK).
		Set("X-TW-SIGNER", primaryPK).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(trustLinePayload).Receive(trustLineResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestRemoveTrustAssetMultiAccessDisabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestRemoveTrustAssetMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}
	if len(trustLineResponse.Transaction) == 0 {
		log.Println("[TestRemoveTrustAssetMultiAccessDisabled]no transaction generated")
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestRemoveTrustAssetMultiAccessDisabled] Request REMOVE ASSET Response:[%+v]\n", trustLineResponse)
	log.Println("==========waiting for 15seconds before confirmation==============")
	time.Sleep(time.Second * 15)
	log.Println("==========Confirming REMOVE TRUSTLINE ==============")

	{
		//run the subwallet signing and submission
		p := *trustLineResponse
		//sign transaction

		primarySignature, err := middleware.SignBase64Txn(primarySecretKey, p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestRemoveTrustAssetMultiAccessDisabled] sub transactions error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = primarySignature

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, primarySecretKey)
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", subPK).
			Set("X-TW-SIGNER", primaryPK).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(trustLineResponse, errorResponse)
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		if len(trustLineResponse.TransactionID) == 0 {
			log.Println("[TestRemoveTrustAssetMultiAccessDisabled]no transaction ID")
			t.Errorf(err.Error())
			return

		}

		log.Printf("REMOVE TRUSTLINE Response:[%+v]\n", trustLineResponse)
	}

}
