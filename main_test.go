package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"trovo-wallet-api/internal/middleware"
	pns "trovo-wallet-api/internal/pns"

	"github.com/dghubble/sling"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
)

const devURL = "http://localhost:8080"
const prodURL = "https://api.trovotechnologies.com"
const stagingURL = "https://apidev.trovotechnologies.com"

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
	Commit                  int               `json:"commit"`
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
	WalletType              int      `json:"walletType"` //0=normal, 1= assetIssuing, 2= marketMaking, 3 = bulkPayments
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
	ID                           string                       `json:"-"`
	Username                     string                       `json:"username"`
	Email                        string                       `json:"email"`
	ImageThumbnailURL            string                       `json:"imageThumbnailURL"`
	FirstName                    string                       `json:"firstName"`
	LastName                     string                       `json:"lastName"`
	Mobile                       string                       `json:"mobile"`
	PublicKey                    string                       `json:"publicKey"`
	PrimarySigner                string                       `json:"primarySigner"`
	Referrer                     string                       `json:"referrer"`
	ReferralLink                 string                       `json:"referralLink"`
	ReferralQrCode               string                       `json:"referralQrCode"`
	PushNotificationToken        string                       `json:"pushNotificationToken"`
	Corporate                    uint                         `json:"corporate"`
	MobileVerified               uint                         `json:"mobileVerified"`
	MembershipType               uint                         `json:"membershipType"`
	MembershipExpiry             time.Time                    `json:"membershipExpiry"`
	KYCVerified                  uint                         `json:"kycVerified"`
	WalletRecoveryEnabled        uint                         `json:"walletRecoveryEnabled"`
	UserWallets                  []UserWalletJSON             `json:"userWallets"`
	Verified                     int                          `json:"verified"`
	Suspended                    int                          `json:"suspended"`
	HasSecurityQuestions         int                          `json:"hasSecurityQuestions"`
	CuratedSwapList              []CuratedSwapAsset           `json:"curatedSwapList"`
	CryptoWalletDepositAddresses []CryptoWalletDepositAddress `json:"cryptoWalletDepositAddresses"`
}

type CryptoWalletDepositAddress struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `gorm:"default:now()" json:"createdAt"`
	UserID               string    `gorm:"not null;size:100;" json:"-"`
	TrovoWalletPublicKey string    `gorm:"not null;size:100;index:idx_unique_address,unique" json:"TrovoWalletPublicKey"`
	Currency             string    `gorm:"not null;size:12;index:idx_unique_address,unique" json:"currency"`
	DepositAddress       string    `gorm:"not null;size:100" json:"depositAddress"`
	Network              string    `gorm:"not null;size:100;index:idx_unique_address,unique" json:"network"`
	QRCode               *string   `gorm:"null;" json:"qrCode"`
}

type CryptoSubwalletResponse struct {
	Message string          `json:"message"`
	Data    CryptoSubWallet `json:"data"`
}
type CryptoSubwalletsResponse struct {
	Message string            `json:"message"`
	Data    []CryptoSubWallet `json:"data"`
}
type CryptoSubWallet struct {
	WalletID     string          `json:"walletId"`
	Addresses    []CryptoAddress `json:"addresses"`
	Currency     string          `json:"currency"`
	IntegratorPk string          `json:"integratorPk,omitempty"`
	CreatedAt    string          `json:"createdAt,omitempty"`
	UpdatedAt    string          `json:"updatedAt,omitempty"`
}
type CryptoAddress struct {
	Address string `json:"address"`
	Network string `json:"network"`
}

type CryptoDepositResponse struct {
	Message string `json:"message"`
	Data    struct {
		DepositID   string `json:"depositId"`
		TxID        string `json:"txId"`
		Amount      string `json:"amount"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		Currency    string `json:"currency"`
		Decimal     int    `json:"decimal"`
		Fees        string `json:"fees"`
		FromAddress string `json:"fromAddress"`
		IsCompleted bool   `json:"isCompleted"`
		IsValid     bool   `json:"isValid"`
		IsVerified  bool   `json:"isVerified"`
		ToAddress   string `json:"toAddress"`
	} `json:"data"`
}

type CryptoDeposit struct {
	DepositID   string `json:"depositId"`
	TxID        string `json:"txId"`
	Amount      string `json:"amount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Currency    string `json:"currency"`
	Decimal     int    `json:"decimal"`
	Fees        string `json:"fees"`
	FromAddress string `json:"fromAddress"`
	IsCompleted bool   `json:"isCompleted"`
	IsValid     bool   `json:"isValid"`
	IsVerified  bool   `json:"isVerified"`
	ToAddress   string `json:"toAddress"`
}

type CuratedSwapAsset struct {
	ID                          uint64     `gorm:"primaryKey" json:"-"`
	CreatedAt                   time.Time  `json:"-"`
	UpdatedAt                   time.Time  `json:"-"`
	AssetCode                   string     `gorm:"size:12;unique;not null" json:"assetCode"`
	AssetName                   string     `gorm:"size:50;null" json:"assetName"`
	AssetIssuer                 string     `gorm:"size:56;not null;" json:"assetIssuer"`
	Description                 string     `gorm:"size:200; not null" json:"description"`
	ImageURL                    string     `gorm:"null" json:"imageUrl"`
	Website                     string     `gorm:"null;size:100" json:"website"`
	AssetConditions             string     `gorm:"null;size:100" json:"assetConditions"`
	AssetLimit                  uint64     `gorm:"type:integer;not null;default:0" json:"assetLimit"` //0 = unlimited
	AssetRedemptionInstructions string     `gorm:"null;size:100" json:"assetRedemptionInstructions"`
	ContactEmail                string     `gorm:"null;size:100" json:"contactEmail"`
	Priority                    uint64     `gorm:"null;unique" json:"-"`
	AssetClassID                uint64     `gorm:"not null; default:1" json:"assetClassId"`
	AssetClass                  AssetClass `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetClass"`
	Organization                string     `gorm:"null;size:100" json:"organization"`
	Withdrawable                uint64     `gorm:"type:integer;not null;default:0" json:"withdrawable"`
	DecimalPlaces               uint64     `gorm:"type:integer;not null;default:7" json:"decimalPlaces"`
	Inactive                    uint64     `gorm:"type:integer;not null;default:1" json:"-"`
}

type AssetClass struct {
	ID         uint64 `gorm:"primaryKey" json:"-"`
	AssetClass string `gorm:"size:45;unique;not null" json:"assetClass"`
}

type UserWalletJSON struct {
	CreatedAt              time.Time                  `json:"createdAt"`
	ID                     string                     `json:"publicKey"`
	TempPublicKey          string                     `json:"-"`
	Tag                    string                     `json:"tag"`
	Description            string                     `json:"description"`
	Alias                  string                     `json:"alias"`  //primaryUsername_tag for sub wallets
	Signer                 string                     `json:"signer"` //if ID is same as signer, then it is a primary wallet
	UserID                 string                     `json:"userId"`
	SharedAccessEnabled    uint                       `json:"sharedAccessEnabled"`
	PrimaryWallet          uint                       `json:"primaryWallet"`
	UserWalletSharedAccess UserWalletSharedAccessJSON `json:"userWalletSharedAccess"`
}

type UserWalletSharedAccessJSON struct {
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	ID                string                 `json:"accessId"`
	UserWalletID      string                 `json:"publicKey"`
	NumberOfApprovers uint                   `json:"numberOfApprovers"`
	Permissions       []WalletPermissionJSON `json:"permissions"`
}

type WalletPermissionJSON struct {
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	Username                  string    `json:"username"`
	Permission                string    `json:"permission"`
	UserWalletManagedAccessID string    `json:"userWalletManagedAccessId"`
}
type ThirdPartyWalletAccess struct {
	Owner             string `json:"owner"`
	PublicKey         string `json:"publicKey"`
	Permission        string `json:"permission"`
	WalletAlias       string `json:"walletAlias"`
	WalletDescription string `json:"walletDescription"`
}

// Balance model for user
type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	QRCode      string          `json:"qrCode"`
	UsdPrice    string          `json:"usdPrice"`
	NativePrice string          `json:"nativePrice"`
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
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Commit               int      `json:"commit"`
	Messages             []string `json:"messages"`
}

type SecurityQuestion struct {
	ID       uint64
	Question string
}

type UserSecurityAnswer struct {
	ID uint64 `json:"id"`
	Q1 uint64 `gorm:"not null;" json:"q1"`
	A1 string `gorm:"size:50;not null;" json:"a1"`
	Q2 uint64 `gorm:"not null;" json:"q2"`
	A2 string `gorm:"size:50;not null;" json:"a2"`
	Q3 uint64 `gorm:"not null;" json:"q3"`
	A3 string `gorm:"size:50;not null;" json:"a3"`
}

type AnswerResp struct {
	Questions   []SecurityQuestion  `json:"securityQuestions"`
	UserAnswers *UserSecurityAnswer `json:"userSecurityAnswers"`
}

type UserAccountRecoveryPayload struct {
	Transaction          string             `json:"transaction"`
	TransactionSignature string             `json:"transactionSignature"`
	TransactionID        string             `json:"transactionId"`
	NetworkPassPhrase    string             `json:"networkPassPhrase"`
	Messages             []string           `json:"messages"`
	SecurityAnswers      UserSecurityAnswer `json:"securityAnswers"`
}

type AccountRecoveryRequest struct {
	NewSignerPublicKey                string             `json:"newSignerPublicKey"`
	DisableOldSignerFromPrimaryWallet uint64             `json:"disableOldSignerFromPrimaryWallet"`
	Commit                            uint64             `json:"commit"`
	Messages                          []string           `json:"messages"`
	SecurityAnswers                   UserSecurityAnswer `json:"securityAnswers"`
	EmailOTP                          string             `json:"emailOtp"`
	Username                          string             `json:"username"`
	TransactionID                     string             `json:"transactionId"`
}
type UserWalletSharedAccessInfo struct {
	WalletPublicKey         string                 `json:"walletPublicKey"`
	NumberOfApprovalsNeeded int                    `json:"numberOfApprovalsNeeded"`
	Permissions             []WalletPermissionInfo `json:"permissions"`
	Transaction             string                 `json:"transaction"`
	TransactionSignature    string                 `json:"transactionSignature"`
	TransactionID           string                 `json:"transactionId"`
	NetworkPassPhrase       string                 `json:"networkPassPhrase"`
	Messages                []string               `json:"messages"`
	SignatureRequired       int                    `json:"signatureRequired"`
}
type ModifySharedAccessInfo struct {
	WalletPublicKey         string                 `json:"walletPublicKey"`
	NumberOfApprovalsNeeded int                    `json:"numberOfApprovalsNeeded"`
	Transaction             string                 `json:"transaction"`
	TransactionSignature    string                 `json:"transactionSignature"`
	TransactionID           string                 `json:"transactionId"`
	NetworkPassPhrase       string                 `json:"networkPassPhrase"`
	Messages                []string               `json:"messages"`
	SignatureRequired       int                    `json:"signatureRequired"`
	MultiParty              int                    `json:"multiParty"`
	Commit                  int                    `json:"commit"`
	SHash                   string                 `json:"sHash"`
	ModifiedPermissions     []WalletPermissionInfo `json:"modifiedPermissions"`
	AddedPermissions        []WalletPermissionInfo `json:"addedPermissions"`
	RevokedPermissions      []WalletPermissionInfo `json:"revokedPermissions"`
}
type WalletPermissionInfo struct {
	ID                    string  `json:"Id"`
	WalletPublicKey       string  `json:"-"`
	WalletAlias           string  `json:"-"`
	TargetUsername        string  `json:"targetUsername"`
	Name                  string  `json:"name"`
	Permission            string  `json:"permission"`
	PushNotificationToken *string `json:"-"`
}
type DisableSharedAccessInfo struct {
	WalletPublicKey      string   `json:"walletPublicKey"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	SignatureRequired    int      `json:"signatureRequired"`
	MultiParty           int      `json:"multiParty"`
	Commit               int      `json:"commit"`
}

type ApprovalPayload struct {
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	RemainingApprovals   int    `json:"remainingApproval"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}

type RejectPayload struct {
	RejectionReason string `json:"rejectionReason"`
}
type MarketOfferRequest struct {
	OfferType            string   `json:"offerType"`
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	CurrencyCode         string   `json:"currencyCode"`
	CurrencyIssuer       string   `json:"currencyIssuer"`
	PricePerUnit         string   `json:"pricePerUnit"`
	Quantity             string   `json:"quantity"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	Commit               int      `json:"commit"`
	SignatureRequired    int      `json:"signatureRequired"`
	Memo                 string   `json:"memo"`
}

type WithdrawalRequestInput struct {
	Currency             string  `json:"currency"`
	AmountSubmitted      float64 `json:"amountSubmitted"`
	AmountToWithdraw     float64 `json:"amountToWithdraw"` //submitted amount less serviceFee
	WithdrawalAddress    string  `json:"withdrawalAddress"`
	WithdrawalNetwork    string  `json:"withdrawalNetwork"`
	WithdrawalMemo       string  `json:"withdrawalMemo"`
	WithdrawalServiceFee float64 `json:"withdrawalServiceFee"`
	WithdrawalNetworkFee float64 `json:"withdrawalNetworkFee"`
	Transaction          string  `json:"transaction"`
	TransactionSignature string  `json:"transactionSignature"`
	TransactionID        string  `json:"transactionId"`
	NetworkPassPhrase    string  `json:"networkPassPhrase"`
	Multiparty           int     `json:"-"`
	SignatureRequired    int     `json:"signatureRequired"`
	Commit               int     `json:"commit"`
	ReturnedDescription  string  `json:"-"`
}

func TestCreateAccount(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := stagingURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := UserRegistrationInfo{
		Username:          "ric1",
		Email:             "chukwunenyeo@gmail.com",
		FirstName:         "Second",
		LastName:          "Account",
		Mobile:            "+234-8050564391",
		MobileCountryCode: "NG",
		PublicKey:         pk,
		Referrer:          "ric",
		VerificationCode:  "157490",
	}
	// payload := UserRegistrationInfo{
	// 	Username:          "ric",
	// 	Email:             "richardsric7@gmail.com",
	// 	FirstName:         "Ric",
	// 	LastName:          "Rcichards",
	// 	Mobile:            "+234-8180067955",
	// 	MobileCountryCode: "NG",
	// 	PublicKey:         pk,
	// 	Referrer:          "",
	// 	VerificationCode:  "120462",
	// }
	errorResponse := new(ErrorResponse)
	rResponse := new(map[string]string)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(rResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateAccount] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}

	if err != nil {
		log.Println("[TestCreateAccount]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", rResponse)

	log.Println("[TestCreateAccount] completed")
	// time.Sleep(time.Second * 10)

}
func TestRequestAccountRecoveryEmailOTP(t *testing.T) {
	// use the new signer key pairs for the request
	pk := "GCFM3WOKRTURPZMQDINQRJA6LOLK3SKDN6W4JDAOIOAOBBYL7YVYGT6Z"
	secretKey := "SC6XAWBZQY6QEEGRRXFN5BX66HE6HHR4TYFZ5JA7PFKBEMG74TMFSN5M"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	ownerUsername := "ric1"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	// fullPath := "/v1/account/recovery/request-email-otp/:targetUser"
	fullPath := fmt.Sprintf("/v1/account/recovery/request-email-otp/%s", ownerUsername)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := UserAccountRecoveryPayload{}
	errorResponse := new(ErrorResponse)
	payResponse := new(map[string]string)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestRequestEmailOTP] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestRequestEmailOTP]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	log.Println("[TestRequestEmailOTP] completed")
	time.Sleep(time.Second * 1)

}
func TestEnableAccountRecovery(t *testing.T) {

	pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
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
	fullPath := "/v1/users/account/recovery"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := UserAccountRecoveryPayload{}
	errorResponse := new(ErrorResponse)
	payResponse := new(UserAccountRecoveryPayload)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestAccountEnableRecovery] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestAccountEnableRecovery]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		//sign transaction

		signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestAccountEnableRecovery] makePayment error:", err)
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

		log.Printf("Enable Recovery Response:[%+v]\n", payResponse)
	}
	log.Println("[TestAccountEnableRecovery] completed")
	time.Sleep(time.Second * 1)

}
func TestDisableAccountRecovery(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
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
	fullPath := "/v1/users/account/recovery"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := UserAccountRecoveryPayload{
		SecurityAnswers: UserSecurityAnswer{
			A1: "test",
			A2: "test",
			A3: "test",
		},
	}
	errorResponse := new(ErrorResponse)
	payResponse := new(UserAccountRecoveryPayload)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Delete(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestAccountDisableRecovery] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestAccountDisableRecovery]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		//sign transaction
		// p.SecretAnswers = UserSecretAnswer{
		// 	A1: "test",
		// 	A2: "test",
		// 	A3: "test",
		// }

		signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestAccountDisableRecovery] confirm transaction error:", err)
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
			Delete(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestAccountDisableRecovery] server response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Disable Recovery Response:[%+v]\n", payResponse)
	}
	log.Println("[TestAccountDisableRecovery] completed")

}

func TestDoAccountRecovery(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// newSigner := "SA4JYDZJSOVWHOWWLGEZV3NSSE3YRS2BVRNDQM2O6UREN53RJTTHUS4P"
	newSigner := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	ownerUsername := "ric1"
	kp := keypair.MustParseFull(newSigner)
	// log.Println(kp.Address())
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/account/recover"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := AccountRecoveryRequest{
		Username:                          ownerUsername,
		EmailOTP:                          "342474",
		NewSignerPublicKey:                kp.Address(),
		DisableOldSignerFromPrimaryWallet: 1,
		Commit:                            0,
		SecurityAnswers: UserSecurityAnswer{
			A1: "test",
			A2: "test",
			A3: "test",
		},
	}
	errorResponse := new(ErrorResponse)
	rResponse := new(AccountRecoveryRequest)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(rResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestDoAccountRecovery] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestDoAccountRecovery]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", rResponse)

	{
		//run the payment signing and submission
		p := *rResponse
		p.Commit = 1
		//sign transaction
		log.Println("set commit = 1")
		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
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
			Post(fullPath).BodyJSON(p).Receive(rResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestDoAccountRecovery] server response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Disable Recovery Response:[%+v]\n", rResponse)
	}
	log.Println("[TestAccountDisableRecovery] completed")

}
func TestAccountSetSecurityAnswer(t *testing.T) {

	answers := UserSecurityAnswer{
		Q1: 1,
		Q2: 2,
		Q3: 7,
		A1: "test",
		A2: "test",
		A3: "test",
	}
	// pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	// secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	primaryPK := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	primarySecretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// primaryPK := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// primarySecretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// primaryPK := os.Getenv("RICPK")
	// primarySecretKey := os.Getenv("RICSC")
	if len(primarySecretKey) == 0 || len(primaryPK) == 0 {
		log.Println("No primary secret or public key specified")
		time.Sleep(20 * time.Second)
		t.Errorf("No primary secret or public key specified")
		return
	}
	kp := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	baseURL := prodURL
	fullPath := "/v1/security-questions"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, primaryPK+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}
	type Success struct {
		Message string `json:"message"`
	}
	errorResponse := new(ErrorResponse)
	regResponse := new(Success)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(answers).Receive(regResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestAccountSetSecretAnswer] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestAccountSetSecretAnswer]request error:", err)
		t.Errorf(err.Error())
		return
	}

	log.Printf("[TestAccountSetSecretAnswer]Result:[%+v]\n", regResponse)

}

func TestAccountProfileUpdate(t *testing.T) {
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
	fullPath := "/v1/users/upload-picture"
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
		Put(fullPath).BodyJSON(userRegInfo).Receive(regResponse, errorResponse)
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

func TestGetSecurityAnswers(t *testing.T) {

	// pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	// secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	primaryPK := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	primarySecretKey := "SA4JYDZJSOVWHOWWLGEZV3NSSE3YRS2BVRNDQM2O6UREN53RJTTHUS4P"
	// primaryPK := os.Getenv("RICPK")
	// primarySecretKey := os.Getenv("RICSC")
	kp := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	// baseURL := "http://localhost:8080"
	ownerUsername := "ric1"
	baseURL := devURL
	fullPath := "/v1/security-questions/" + ownerUsername
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	resultResponse := new(AnswerResp)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", primaryPK).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Get(fullPath).Receive(resultResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestGetSecurityAnswers] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestGetSecurityAnswers]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestGetSecurityAnswers] Result:[Questions: %+v\n Answers: %+v]\n", resultResponse.Questions, resultResponse.UserAnswers)
	log.Println("[TestGetSecurityAnswers] completed.")
}
func TestGetUserInfo(t *testing.T) {

	// pk := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	// secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	// primaryPK := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// primaryPK := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// primarySecretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	primaryPK := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	primarySecretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// primaryPK := os.Getenv("RICPK")
	// primarySecretKey := os.Getenv("RICSC")
	ownerUsername := "ric1"
	kp := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	// baseURL := "http://localhost:8080"
	baseURL := stagingURL
	fullPath := fmt.Sprintf("/v1/users/%s", ownerUsername)
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	resultResponse := new(UserInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", primaryPK).
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

	// p k := "GCATEXQ3TNQU7IYBOCXMAKTWJ4FXXZ5POUZ4VS4VMVU2H43XLNFAJJUF"
	// secretKey := "SDZZHRY6BJ5MHMOZCVZC5TT3XKOXGPDVJRE7CK7NZHR35ORDGGZ2VJGP"
	// publickey := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	primaryPK := os.Getenv("RICPK")
	addressToViewHistory := os.Getenv("RICPK")
	// primarySecretKey := os.Getenv("RICSC")
	// ownerUsername := "ric1"
	// pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	// baseURL := "http://localhost:8080"
	// baseURL := devURL
	baseURL := prodURL
	fullPath := fmt.Sprintf("/v1/users/payments/%v", addressToViewHistory)
	ts := time.Now().Unix() / 1000

	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	resultResponse := new(PaginatedPaymentHistory)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", primaryPK).
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

func TestSendPushNotificationMessage(t *testing.T) {
	ric := "dWLRIQWuSm-sqAKA-ABwhS:APA91bGJt8PE4KBS0OPIJOVp4JsWbpiJMK1DrIJIhZM7hlOVeLo6OUGlN5PbbsttT3Oq0YNXZZ8P0zDEcVD6wQkdoFOfrTSFq0q9A1XZU605ZhTpaLLfqwqONniRQKoj4I-YdMgxXoJc"
	// kennis := "eNCa_XRaTr2NnXX4pnzhN3:APA91bF9OfBO9IEFJcPOtO-83Qu41_7zZ3ef7qC3i5ySPvT8arcQ1gwnRXYnSZ5uJ9mT4uOW7rgPp5F0hTsqvoqQL9oR02fQtiCyco2DVsNBT6JIgqgVHO1ZTPod7ypm-MpSzA95MRRZ"
	title := "TROVO: Testing Push Notification Service"
	body := `This is a test message to ascertain how the push notification appears`
	imageURL := "https://trovotech.io/img/Trovotech-colored.png"
	dataPayload := make(map[string]string)
	dataPayload["route"] = "announcements"
	dataPayload["openLink"] = "https://trovowallet.page.link"
	ctx := context.Background()
	client, _, err := pns.GetFirebaseMessagingClient(ctx)
	if err != nil {
		log.Println("[TestGetUserInfo]request error:", err)
		t.Errorf(err.Error())
		return
	}
	response, _ := pns.SendFirebaseMessage(ric, title, body, imageURL, dataPayload, client, ctx)
	// response, _ := fb.SendFirebaseMessage(ric, title, body, imageURL, nil, nil)

	log.Printf("Result:[%+v]\n", response)

}

func TestSendPushNotificationBroadcast(t *testing.T) {

	// ric := "dqmlIz2eT-CTEX1cesc3za:APA91bGDuSBPzGSZAjYfFMU12eV_gpMw7JdF7gxhhsaPc0XGmgYu9RokVyLfEpH6_Yfg9mYrmuhYlcLRScvolUbfRVftR3QK4d6psgVOY0382fNW3Q1BND8kTbPAThBOXe60DHSXO9hD"
	// kennis := "eNCa_XRaTr2NnXX4pnzhN3:APA91bF9OfBO9IEFJcPOtO-83Qu41_7zZ3ef7qC3i5ySPvT8arcQ1gwnRXYnSZ5uJ9mT4uOW7rgPp5F0hTsqvoqQL9oR02fQtiCyco2DVsNBT6JIgqgVHO1ZTPod7ypm-MpSzA95MRRZ"

	receipients := []string{"godswill", "tunde", "onoja", "cryptoking", "ric", "mavol"}
	title := "TESTING THE PUSH NOTIFICATION BROADCAST"
	body := `This is a test message to ascertain how the push notification appears`
	imageURL := "https://trovotech.io/img/Trovotech-colored.png"
	ctx := context.TODO()
	client, _, err := pns.GetFirebaseMessagingClient(ctx)
	if err != nil {
		log.Println("[TestGetUserInfo]request error:", err)
		t.Errorf(err.Error())
		return
	}
	response, _ := pns.SendFirebaseBroadcast(receipients, title, body, imageURL, nil, client, ctx)

	log.Printf("Result:[%+v]\n", response)

}

/**
2022/06/30 18:37:43 Full Path With Query:[/v1/users/payment] KeyParam:[GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC1656614] Signature: [WnEUUYVAjhVRSvqJ/ndOf8TMlIil45Gua2Cn2PGmZULrb5RGVcZuOxYskp0BiHWqGKpRpRi+6Kfhf/DA4jrpCA==]
**/

// TestSendPaymentMultiAccessDisabled sends payment from primary wallet
func TestClaimAssetMultiAccessDisabled(t *testing.T) {

	pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	fromWallet := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// fromWallet := "GDW6UKK6RI2LBTGHTDKKXYZKCGPDFBRFDTYSZKGGGE6SC5TCSG3MMJST"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")

	kp := keypair.MustParseFull(secretKey)
	baseURL := devURL

	fullPath := "/v1/users/actions/claim-asset"
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	// claimPayload := PendingAssetToClaim{
	// 	AssetCode:   "ABC",
	// 	AssetIssuer: "GAD3DZNQY4SXJEUJOPLJZEK3OWTASEUK2LZYT3V7C52UN5QYOFP3PM5P",
	// }

	// claimPayload := PendingAssetToClaim{
	// 	AssetCode:   "YAM",
	// 	AssetIssuer: "GAJ65QHSOIXOA6FZMKDIBNGMHXQ7U46TNRBKDL3MTHERF2VRVMWU2F57",
	// }

	claimPayload := PendingAssetToClaim{
		AssetCode:   "LUMI",
		AssetIssuer: "GBGUHXVAK32BZTWRBML7RNIQ3532QDR5RRXOJ2P2MGEPHC3YJREMRTLE",
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

// TestClaimAssetMultiAccessEnabled sends payment from primary wallet
func TestClaimAssetMultiAccessEnabled(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// fromWallet := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	fromWallet := "GD6IO3P4J2C63Z3VEIH5TVZVDITHKGJMOAKHX6J6TEDA6JEEQCD5GJFN"
	// fromWallet := "GDW6UKK6RI2LBTGHTDKKXYZKCGPDFBRFDTYSZKGGGE6SC5TCSG3MMJST"
	// pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")

	kp := keypair.MustParseFull(secretKey)
	// baseURL := devURL
	baseURL := prodURL

	fullPath := "/v1/shared-access/users/actions/claim-asset"
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	// claimPayload := PendingAssetToClaim{
	// 	AssetCode:   "ABC",
	// 	AssetIssuer: "GAD3DZNQY4SXJEUJOPLJZEK3OWTASEUK2LZYT3V7C52UN5QYOFP3PM5P",
	// }

	claimPayload := PendingAssetToClaim{
		AssetCode:   "YAM",
		AssetIssuer: "GAJ65QHSOIXOA6FZMKDIBNGMHXQ7U46TNRBKDL3MTHERF2VRVMWU2F57",
	}

	// claimPayload := PendingAssetToClaim{
	// 	AssetCode:   "LUMI",
	// 	AssetIssuer: "GBGUHXVAK32BZTWRBML7RNIQ3532QDR5RRXOJ2P2MGEPHC3YJREMRTLE",
	// }
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
		log.Println("[TestClaimAssetMultiAccessEnabled] server response error:", *errorResponse)
		return

	}
	if err != nil {
		log.Println("[TestClaimAssetMultiAccessEnabled]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Claim Response:[%+v]\n", claimResponse)

	{
		//run the payment signing and submission
		p := *claimResponse
		p.Commit = 1
		//sign transaction

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
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
		if len(errorResponse.Error) > 0 {
			log.Println("[TestClaimAssetMultiAccessEnabled] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Make Claim Response:[%+v]\n", claimResponse)
	}
	log.Println("[TestClaimAssetMultiAccessEnabled] completed")

}
func TestSendPaymentMultiAccessDisabled(t *testing.T) {

	pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
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

	paymentPayload := PaymentInfo{
		Destination: "ric1_shared",
		Memo:        "Test Payment",
		Amount:      "20000",
		AssetCode:   "",
		AssetIssuer: "",
	}

	// paymentPayload := PaymentInfo{
	// 	Destination: "ric1",
	// 	Memo:        "Test Payment",
	// 	Amount:      "10",
	// 	AssetCode:   "YAM",
	// 	AssetIssuer: "GAJ65QHSOIXOA6FZMKDIBNGMHXQ7U46TNRBKDL3MTHERF2VRVMWU2F57",
	// }

	// paymentPayload := PaymentInfo{
	// 	Destination: "ric_joint",
	// 	Memo:        "Test Payment",
	// 	Amount:      "2000",
	// 	AssetCode:   "ABC",
	// 	AssetIssuer: "GAD3DZNQY4SXJEUJOPLJZEK3OWTASEUK2LZYT3V7C52UN5QYOFP3PM5P",
	// }

	// paymentPayload := PaymentInfo{
	// 	Destination: "ric1",
	// 	Memo:        "Test Payment",
	// 	Amount:      "100",
	// 	AssetCode:   "LUMI",
	// 	AssetIssuer: "GBGUHXVAK32BZTWRBML7RNIQ3532QDR5RRXOJ2P2MGEPHC3YJREMRTLE",
	// }
	errorResponse := new(ErrorResponse)
	payResponse := new(PaymentInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", pk).
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
		// t.Errorf(err.Error())

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
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", pk).
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
	// time.Sleep(time.Second * 10)

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
		t.Errorf(errorResponse.Error)
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

// TestSendPaymentWithSharedAccessEnabled sends payment from primary account
func TestSendPaymentWithSharedAccessEnabled(t *testing.T) {

	// pk := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// secretKey := "SDBLGMM6HVLYSUUR2TIKC6E7GZHQA5VJUUGBVOGDC5KQHTJVC2KK3EXK"
	// pk := os.Getenv("RICPK")
	// fromWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	//ric_joint
	// fromWallet := "GD6IO3P4J2C63Z3VEIH5TVZVDITHKGJMOAKHX6J6TEDA6JEEQCD5GJFN"
	//ric_joint1
	fromWallet := "GA7ZU2CZPXVCBCDABXOQ7BNLCF24PASOZZTLFGPN4PFODLC7UVO3XYGU"
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
	fullPath := "/v1/shared-access/payment"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, signerKP.Address()+tsString, signerKP.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	paymentPayload := PaymentInfo{
		Destination: "nechey",
		Memo:        "Test XBN shared Payment",
		Amount:      "55.2",
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
		log.Println("[TestSendPaymentWithSharedAccessEnabled] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestSendPaymentWithSharedAccessEnabled]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestSendPaymentWithSharedAccessEnabled]Confirmation Payment Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		// commit transaction
		p.Commit = 1

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

		log.Printf("Make shared Payment Response:[%+v]\n", payResponse)
	}
	log.Println("[TestSendPaymentWithSharedAccessEnabled] completed")
	// time.Sleep(time.Second * 10)

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
	// subPK := "GDQF3BDD5JQN5N7HRMFBEKU42DYDJI4WMZHXMPCGYVZKCE7QRN6Y3ERH"
	// subSecretKey := "SB2LEXZ6UBRXTGBGIDXXXE6VAU333FBTHSXV6KEKES3WO3NNTKTMMF3C"
	// subPK := "GBAI3QHD73YQO3S5L55OCT62DBTNGVR4JHEI4Q4DYPURT72WHK6U6NWS"
	// subSecretKey := "SD47WSETFWODYVZXYOBSNL3E5TFMBV7SZZF3YESJRKHYYWVEXPEQC2GT"
	// subPK := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	// subSecretKey := "SBSMH2IHU4HHK4DBDP6PUNT2ZZIKJKWYXJAYSEGXL6D6GWXMPQMCDGPM"

	// subPK := "GD6IO3P4J2C63Z3VEIH5TVZVDITHKGJMOAKHX6J6TEDA6JEEQCD5GJFN"
	// subSecretKey := "SDCSU5C6F4HWD5T7QBIG4U2Z3VITF5HY6KC5XPEPIDDGD7GK2UNEYARG"

	// subPK := "GAYKJR7KECN57NPKF4ABYQPFLUCELKXMSPD3D7ACEATI77TYFXKJSKRO"
	// subSecretKey := "SAHTUGJVWK7WCERUDZM5VVDJO2JYRUTUQQ7O7CJVLLVSYZ7CTV2SO5IP"
	// subPK := "GB3ZYN2EUPKQHXPLPJ3FATGQKOOIY7CHJ6FJW6RB54ANIT7N5VUZOHVN"
	// subSecretKey := "SBJINW3YZ7IM7DEA7GVBTPYOTTEHRVBWYAF6G3E2WQIXRM7EGYN66E7Q"
	subPK := "GAKWJJTLNYA74XIZZIVFWHTG6BLJWDAVJG4KT5I3H5JJW6INAOKWLHFV"
	subSecretKey := "SBQGZN2EAYMKS37BQYMRUHAZY5YDJIFDTZKBIRSWRURTVDH6GPDJYR4Q"
	primaryPK := os.Getenv("RICPK")
	primarySecretKey := os.Getenv("RICSC")

	// primaryPK := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// primarySecretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// primaryPK := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	// primarySecretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// primaryPK := "GD36GHMT65T2O5YOSFE57TLF4VTSI67IAQXSUT4L5SNBKPNMV5R5R6VV"
	// primarySecretKey := "SAWWK6BIPRALRRHVHELHI2Q3U66KBZLTLOPE7DVGJYRKZYFZGBCZZALY"
	channelAccountSK := ""
	// ownerUsername := "ric"
	// subKP := keypair.MustParseFull(subSecretKey)
	// primaryKP := keypair.MustParseFull(primarySecretKey)
	// log.Println(kp.Address())
	// baseURL := prodURL
	baseURL := prodURL
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
		WalletTag:         "bulkpay",
		WalletDescription: "market making wallet",
		WalletType:        3,
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
	if err != nil {
		log.Println("[TestCreateSubWalletMultiAccessDisabled]request error:", err)
		t.Errorf(err.Error())

		return
	}
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateSubWalletMultiAccessDisabled] server response error:", *errorResponse)
		return

	}

	if len(subWalletResponse.Transaction) == 0 {
		log.Println("[TestCreateSubWalletMultiAccessDisabled]no transaction generated")
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestCreateSubWalletMultiAccessDisabled] Request Subwallet Response:[%+v]\n", subWalletResponse)
	log.Println("==========waiting for 15seconds to before confirmation==============")
	time.Sleep(time.Second * 2)
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
	// time.Sleep(time.Second * 1)

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

func TestCreateSharedAccess(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk :=  os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"
	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// accessToWallet := "GD6IO3P4J2C63Z3VEIH5TVZVDITHKGJMOAKHX6J6TEDA6JEEQCD5GJFN"
	// accessToWallet := "GAYKJR7KECN57NPKF4ABYQPFLUCELKXMSPD3D7ACEATI77TYFXKJSKRO"
	//ric_join1
	// accessToWallet := "GA7ZU2CZPXVCBCDABXOQ7BNLCF24PASOZZTLFGPN4PFODLC7UVO3XYGU"
	accessToWallet := "GBIYYWYIDTZAMKABTMUXCPNFNWDNB72473MZL6ZSQWBYU2TUNCQIOD4K"
	// channelAccountSK := ""
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
	fullPath := "/v1/shared-access/users/account"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}
	var accessList []WalletPermissionInfo
	payload := UserWalletSharedAccessInfo{
		NumberOfApprovalsNeeded: 1,
		// Commit:            1,
	}
	accessList = append(accessList,
		WalletPermissionInfo{
			TargetUsername: "kenmaddy",
			Permission:     "APPROVER"},
		WalletPermissionInfo{
			TargetUsername: "ric",
			Permission:     "INITIATOR"},
		WalletPermissionInfo{
			TargetUsername: "ric",
			Permission:     "APPROVER"},
		WalletPermissionInfo{
			TargetUsername: "ric1",
			Permission:     "INITIATOR"},
		WalletPermissionInfo{
			TargetUsername: "ric1",
			Permission:     "APPROVER"})
	payload.Permissions = accessList

	log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(UserWalletSharedAccessInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", accessToWallet).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateSharedAccess] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestCreateSharedAccess]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		// p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestCreateSharedAccess] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", accessToWallet).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestCreateSharedAccess] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Shared Access Response:[%+v]\n", payResponse)
	}
	log.Println("[TestCreateSharedAccess] completed")

}

func TestCreateSharedAccessWithApprover(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk :=  os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"
	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	accessToWallet := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	// channelAccountSK := ""
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
	fullPath := "/v1/shared-access/users/account"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}
	var accessList []WalletPermissionInfo
	payload := UserWalletSharedAccessInfo{
		NumberOfApprovalsNeeded: 2,
		// Commit:            1,
	}
	accessList = append(accessList,
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "ric",
			Permission:      "INITIATOR"},
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "ric1",
			Permission:      "INITIATOR"},
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "ric",
			Permission:      "APPROVER"},
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "ric1",
			Permission:      "APPROVER"},
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "onoja",
			Permission:      "APPROVER"})
	payload.Permissions = accessList

	log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(UserWalletSharedAccessInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", accessToWallet).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateSharedAccess] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestCreateSharedAccess]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		// p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestCreateSharedAccess] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", accessToWallet).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestCreateSharedAccess] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Shared Access Response:[%+v]\n", payResponse)
	}
	log.Println("[TestCreateSharedAccess] completed")

}
func TestRemoveSharedAccessOnReadOnly(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk :=  os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"
	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	accessToWallet := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	// channelAccountSK := ""
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
	fullPath := "/v1/shared-access/users/account"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := DisableSharedAccessInfo{}

	log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(DisableSharedAccessInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", accessToWallet).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Delete(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestRemoveSharedAccessOnReadOnly] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestRemoveSharedAccessOnReadOnly]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		// p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestRemoveSharedAccessOnReadOnly] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", accessToWallet).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Delete(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestRemoveSharedAccessOnReadOnly] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Remove Shared Access Response:[%+v]\n", payResponse)
	}
	log.Println("[TestRemoveSharedAccessOnReadOnly] completed")

}

func TestModifySharedAccess(t *testing.T) {
	//
	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk :=  os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"
	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// accessToWallet := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	accessToWallet := "GBIYYWYIDTZAMKABTMUXCPNFNWDNB72473MZL6ZSQWBYU2TUNCQIOD4K"
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := stagingURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/shared-access/users/account"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}
	var addedList, revokeList, modifyList []WalletPermissionInfo
	payload := ModifySharedAccessInfo{
		NumberOfApprovalsNeeded: 1,
	}
	addedList = make([]WalletPermissionInfo, 0)
	modifyList = make([]WalletPermissionInfo, 0)
	revokeList = make([]WalletPermissionInfo, 0)
	addedList = append(addedList,
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "efizee",
			Permission:      "INITIATOR"},
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "muche",
			Permission:      "INITIATOR"},
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "muche",
			Permission:      "APPROVER"})

	// modifyList = append(modifyList,
	// 	WalletPermissionInfo{
	// 		WalletPublicKey: accessToWallet,
	// 		TargetUsername:  "ric",
	// 		Permission:      "INITIATOR"},
	// 	WalletPermissionInfo{
	// 		WalletPublicKey: accessToWallet,
	// 		TargetUsername:  "ric1",
	// 		Permission:      "INITIATOR"},
	// 	WalletPermissionInfo{
	// 		WalletPublicKey: accessToWallet,
	// 		TargetUsername:  "ric",
	// 		Permission:      "APPROVER"},
	// 	WalletPermissionInfo{
	// 		WalletPublicKey: accessToWallet,
	// 		TargetUsername:  "ric1",
	// 		Permission:      "APPROVER"},
	// 	WalletPermissionInfo{
	// 		WalletPublicKey: accessToWallet,
	// 		TargetUsername:  "kenmaddy",
	// 		Permission:      "APPROVER"})

	revokeList = append(revokeList,
		WalletPermissionInfo{
			WalletPublicKey: accessToWallet,
			TargetUsername:  "onoja",
			Permission:      "APPROVER"})

	payload.AddedPermissions = addedList
	payload.RevokedPermissions = revokeList
	payload.ModifiedPermissions = modifyList

	// log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(ModifySharedAccessInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", accessToWallet).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Put(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestModifySharedAccessWithApprover] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestModifySharedAccessWithApprover]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			p.Commit = 0
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestModifySharedAccessWithApprover] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", accessToWallet).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Put(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestModifySharedAccessWithApprover] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Shared Access Response:[%+v]\n", payResponse)
	}
	log.Println("[TestModifySharedAccessWithApprover] completed")

}
func TestRemoveSharedAccessWithApprover(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk :=  os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"
	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	accessToWallet := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	// channelAccountSK := ""
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
	fullPath := "/v1/shared-access/users/account"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := DisableSharedAccessInfo{}

	log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(DisableSharedAccessInfo)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", accessToWallet).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Delete(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestRemoveSharedAccessWithApprover] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestRemoveSharedAccessWithApprover]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse
		p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestRemoveSharedAccessWithApprover] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", accessToWallet).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Delete(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestRemoveSharedAccessWithApprover] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Remove Shared Access Response:[%+v]\n", payResponse)
	}
	log.Println("[TestRemoveSharedAccessWithApprovers] completed")

}

func TestApproveTransaction(t *testing.T) {
	approvalID := "d21601d0-7bd9-41b3-ad8c-f3d6349bf8b7"
	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"

	//ric1
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"

	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// accessToWallet := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	pk := kp.Address()
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	// approvalID := "1f8a4d47-cd71-44f5-8d7b-9eb9326be91f"
	// approvalID := "6c8d4dd6-d0dc-4dcb-a67f-6e67f643d9d0"

	fullPath := "/v1/shared-access/approval/" + approvalID
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := ApprovalPayload{}

	log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(ApprovalPayload)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", pk).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestApproveTransaction] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestApproveTransaction]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	{
		//run the payment signing and submission
		p := *payResponse

		//sign transaction

		signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
		if err != nil {
			log.Println("[TestApproveTransaction] confirm transaction error:", err)
			t.Errorf(err.Error())

			return
		}

		p.TransactionSignature = signedBase64

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", pk).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(payResponse, errorResponse)
		if len(errorResponse.Error) > 0 {
			log.Println("[TestApproveTransaction] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}
		if err != nil {
			t.Errorf(err.Error())
			return

		}

		log.Printf("Approval Response:[%+v]\n", payResponse)
	}
	log.Println("[TestApproveTransaction] completed")

}
func TestRejectTransaction(t *testing.T) {
	approvalID := "d21601d0-7bd9-41b3-ad8c-f3d6349bf8b7"
	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk := os.Getenv("RICPK")
	secretKey := os.Getenv("RICSC")
	// pk := "GCC3HG535RVZ3MPTDBANZH7V2HRDEQH3LZXDPBEKPJKZBI2UYJR3OJGF"
	// pk := "GDBWYZWLYASCZ6KP4AIRNRY5WQ5OX6H2T6WASG7WFAEEYO6R6AC4GXRM"
	// secretKey := "SBKXWM6TWUVY6NEVRO3CXTKALILMFG2R4WQAAXYKII665U2RDHQ5EB3B"

	//ric1
	// secretKey := "SB2KSQNONOLO2RRS44TTHSCQRDO4WDUFSRT64LPA4TNWI4C6A34GDIKS"

	// accessToWallet := "GDIJRIJ7OFKK4IYUCYGP6GQIMNLCIO4U7EDH7JX3626JS4ACY6WZNIH2"
	// accessToWallet := "GCN2Z2ZV7GKZMJQMUJUFSAKV5BGK5ECZMWLGEBDHC5QOHM66J4FCQXUZ"
	// accessToWallet := "GBU5IARLMK3DG6E5VJNFWLKYF6FP53CPX6X6XIV7YPMA6XYAC27M55SN"
	// accessToWallet := "GBQBJFGWYXCKSKTXFCG5WMPKQC3LYJPSPRNADVPQD7W6K5SXSOW744MQ"
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	pk := kp.Address()
	baseURL := prodURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	// approvalID := "1f8a4d47-cd71-44f5-8d7b-9eb9326be91f"
	// approvalID := "6134726e-1833-43f3-94ef-d4676a5f2662"

	fullPath := "/v1/shared-access/approval/" + approvalID
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := RejectPayload{
		RejectionReason: "This is not what was approved. We approved 52 XBN",
	}

	log.Printf("[DEBUG] Payload: %+v\n", payload)
	errorResponse := new(ErrorResponse)
	payResponse := new(ApprovalPayload)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", pk).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Delete(fullPath).BodyJSON(payload).Receive(payResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestRejectTransaction] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}
	if err != nil {
		log.Println("[TestRejectTransaction]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", payResponse)

	log.Println("[TestRejectTransaction] completed")

}

func TestCreateMarketOffer(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := stagingURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/users/trades"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := MarketOfferRequest{
		OfferType:      "BUY",
		AssetCode:      "NETFLIXSUBSC",
		AssetIssuer:    "GAWLRSFF6Y72AYJ5OWV4HVLP56TJAIWKGPFC7N3K22BUBY245ZTKJNHW",
		CurrencyCode:   "TROV",
		CurrencyIssuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
		PricePerUnit:   "0.00018",
		Quantity:       "2000000",
	}

	errorResponse := new(ErrorResponse)
	rResponse := new(MarketOfferRequest)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(rResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateMarketOffer] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}

	if err != nil {
		log.Println("[TestCreateMarketOffer]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", rResponse)
	{
		//run the payment signing and submission
		p := *rResponse
		p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			p.Commit = 0
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestCreateMarketOffer] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
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
			Post(fullPath).BodyJSON(p).Receive(rResponse, errorResponse)

		if err != nil {
			log.Println("[TestCreateMarketOffer] server 2nd response error:", err.Error())

			t.Errorf("[TestCreateMarketOffer] server second response error: %v", err)
			return

		}
		if len(errorResponse.Error) > 0 {
			log.Println("[TestCreateMarketOffer] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}

		log.Printf("Create market offer Response:[%+v]\n", p)
	}

	log.Println("[TestCreateMarketOffer] completed")
	// time.Sleep(time.Second * 10)

}

func TestGenerateCryptoDepositAddress(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := stagingURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/crypto/generate-addresses/usdt"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, pk+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	errorResponse := new(ErrorResponse)
	rResponse := new([]CryptoWalletDepositAddress)

	sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", kp.Address()).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).Receive(rResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestGenerateCryptoDepositAddress] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}

	log.Printf("Confirmation Response:[%+v]\n", rResponse)

	log.Println("[TestGenerateCryptoDepositAddress] completed")
	// time.Sleep(time.Second * 10)

}

func TestCreateWithdrawalRequest(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := stagingURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/crypto/withdrawals"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := WithdrawalRequestInput{
		Currency:          "USDT",
		AmountSubmitted:   79,
		WithdrawalAddress: "Fphf1sHNtudEWteRNHMdU1SiwuXkdRdJZXLJyKmbu2V8",
		WithdrawalNetwork: "SOL",
	}

	// payload := WithdrawalRequestInput{
	// 	Currency:          "BTC",
	// 	AmountSubmitted:   0.001,
	// 	WithdrawalAddress: "bc1q7adgaawtg8l66zvmsc07r9qfd7lf05mzy5st3h",
	// 	WithdrawalNetwork: "BTC",
	// }

	errorResponse := new(ErrorResponse)
	rResponse := new(WithdrawalRequestInput)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", pk).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(rResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateWithdrawalRequest] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}

	if err != nil {
		log.Println("[TestCreateWithdrawalRequest]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("Confirmation Response:[%+v]\n", rResponse)
	{
		//run the payment signing and submission
		p := *rResponse
		log.Printf("[TestCreateWithdrawalRequest] response: %+v\n", p)
		time.Sleep(5 * time.Second)
		p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			p.Commit = 0
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestCreateWithdrawalRequest] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", pk).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(rResponse, errorResponse)

		if err != nil {
			log.Println("[TestCreateWithdrawalRequest] server 2nd response error:", err.Error())

			t.Errorf("[TestCreateWithdrawalRequest] server second response error: %v", err)
			return

		}
		if len(errorResponse.Error) > 0 {
			log.Println("[TestCreateWithdrawalRequest] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}

		log.Printf("TestCreateWithdrawalRequest Response:[%+v]\n", p)
	}

	log.Println("[TestCreateWithdrawalRequest] completed")
	// time.Sleep(time.Second * 10)

}
func TestCreateWithdrawalRequestShared(t *testing.T) {

	// pk := "GCSTDHLYVVFGNPWASPOVAIRJOQVDDJJON2S3AB3LNXX3PDJCIGDMUQZM"
	// secretKey := "SCIPZFUIWIZEHHAIHDQVOTGODPHMHNAZC2VBC7PN3YYD74PQYFHGCP4F"
	// pk := "GCZ77KBBPINJRHZEYZMCF7SSR5WZVDCUPFG6OSB6FORQVEJV2UOHBG3B"
	pk := "GDLAUQBDFCNO5LJILXMTVSVQHCQOEKKZ7ANYHL2W75WJDAF3QDHVMGGK" //ric1_shared
	secretKey := "SA37LXNUXO62HXXL2SUXVLDCUA6SSQAOUSO2B3LNVMAO3WPE3RDK5OPZ"
	// pk := os.Getenv("RICPK")
	// secretKey := os.Getenv("RICSC")
	// channelAccountSK := ""
	// ownerUsername := "ric"
	kp := keypair.MustParseFull(secretKey)
	// log.Println(kp.Address())
	baseURL := stagingURL
	// var sEnc string
	// if strings.Contains(ownerUsername, "/") {
	// 	sEnc = base64.URLEncoding.EncodeToString([]byte(ownerUsername))

	// } else {
	// 	sEnc = ownerUsername
	// }
	fullPath := "/v1/shared-access/crypto/withdrawals"
	// fullPath := fmt.Sprintf("/v1/users", targetUser, loginID)
	ts := time.Now().Unix() / 1000
	tsString := fmt.Sprintf("%v", ts)
	signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
	if err != nil {
		t.Errorf(err.Error())
		return

	}

	payload := WithdrawalRequestInput{
		Currency:          "USDT",
		AmountSubmitted:   78,
		WithdrawalAddress: "Fphf1sHNtudEWteRNHMdU1SiwuXkdRdJZXLJyKmbu2V8",
		WithdrawalNetwork: "SOL",
	}

	// payload := WithdrawalRequestInput{
	// 	Currency:          "BTC",
	// 	AmountSubmitted:   0.001,
	// 	WithdrawalAddress: "bc1q7adgaawtg8l66zvmsc07r9qfd7lf05mzy5st3h",
	// 	WithdrawalNetwork: "BTC",
	// }

	errorResponse := new(ErrorResponse)
	rResponse := new(WithdrawalRequestInput)

	_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
		Set("X-TW-PUBLIC-KEY", pk).
		Set("X-TW-SIGNER", kp.Address()).
		Set("X-TW-SIGNATURE", signedHttpHeader).
		Set("X-TW-TIMESTAMP", tsString).
		Base(baseURL).
		Post(fullPath).BodyJSON(payload).Receive(rResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		log.Println("[TestCreateWithdrawalRequestShared] server response error:", *errorResponse)
		t.Errorf(errorResponse.Error)
		return

	}

	if err != nil {
		log.Println("[TestCreatTestCreateWithdrawalRequestSharedeWithdrawalRequest]request error:", err)
		t.Errorf(err.Error())

		return
	}

	log.Printf("[TestCreateWithdrawalRequestShared] Confirmation Response:[%+v]\n", rResponse)
	{
		//run the payment signing and submission
		p := *rResponse
		log.Printf("[TestCreateWithdrawalRequestShared] response: %+v\n", p)
		time.Sleep(5 * time.Second)
		p.Commit = 1
		//sign transaction
		if p.SignatureRequired == 1 {
			p.Commit = 0
			signedBase64, err := middleware.SignBase64Txn(kp.Seed(), p.Transaction, p.NetworkPassPhrase)
			if err != nil {
				log.Println("[TestCreateWithdrawalRequestShared] confirm transaction error:", err)
				t.Errorf(err.Error())

				return
			}

			p.TransactionSignature = signedBase64
		}

		ts := time.Now().Unix() / 1000
		tsString := fmt.Sprintf("%v", ts)
		signedHttpHeader, err := middleware.SignHttp(fullPath, kp.Address()+tsString, kp.Seed())
		if err != nil {
			t.Errorf(err.Error())
			return

		}
		_, err = sling.New().Set("User-Agent", "TROVO Go TEST").
			Set("X-TW-PUBLIC-KEY", pk).
			Set("X-TW-SIGNER", kp.Address()).
			Set("X-TW-SIGNATURE", signedHttpHeader).
			Set("X-TW-TIMESTAMP", tsString).
			Base(baseURL).
			Post(fullPath).BodyJSON(p).Receive(rResponse, errorResponse)

		if err != nil {
			log.Println("[TestCreateWithdrawalRequestShared] server 2nd response error:", err.Error())

			t.Errorf("[TestCreateWithdrawalRequestShared] server second response error: %v", err)
			return

		}
		if len(errorResponse.Error) > 0 {
			log.Println("[TestCreateWithdrawalRequestShared] server 2nd response error:", *errorResponse)
			t.Errorf(errorResponse.Error)
			return

		}

		log.Printf("TestCreateWithdrawalRequestShared Response:[%+v]\n", p)
	}

	log.Println("[TestCreateWithdrawalRequestShared] completed")
	// time.Sleep(time.Second * 10)

}
