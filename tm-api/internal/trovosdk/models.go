package trovosdk

import (
	"time"

	"github.com/shopspring/decimal"
)

type ServiceLink struct {
	ApiBaseUrl      string `json:"apiBaseurl"`
	ServiceUsername string `json:"serviceUsername"`
	ApiKey          string `json:"apiKey"`
}
type ErrorResponse struct {
	Error   string `json:"error"`
	Data    string `json:"data"`
	Message string `json:"message"`
}
type ServiceLinksUser struct {
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	LastUpdatedMobileOn   time.Time `json:"lastUpdatedMobileOn"`
	ID                    string    `json:"id"`
	Username              string    `json:"username"`
	Email                 string    `json:"email"`
	ImageThumbnailURL     *string   `json:"imageThumbnailURL"`
	FirstName             string    `json:"firstName"`
	LastName              *string   `json:"lastName"`
	Mobile                *string   `json:"mobile"`
	PublicKey             string    `json:"publicKey"`
	PrimarySigner         string    `json:"primarySigner"`
	PushNotificationToken *string   `json:"pushNotificationToken"`
	Corporate             int       `json:"corporate"`
	MobileVerified        int       `json:"mobileVerified"`
	KYCVerified           int       `json:"kycVerified"`
	Suspended             int       `json:"suspended"`
}
type ServiceLinkTokenizedAssetAuthRequest struct {
	UnsignedTransaction string `json:"unsignedTransaction,omitempty"`
	AssetCode           string `json:"assetCode,omitempty"`
	SignedTransaction   string `json:"signedTransaction,omitempty"`
}

// Balance model for user
type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	InTrade     TradeLiabilties `json:"inTrade"`
	QRCode      string          `json:"qrCode"`
	ImageURL    string          `json:"imageUrl"`
	UsdPrice    string          `json:"usdPrice"`
	NativePrice string          `json:"nativePrice"`
}

type TradeLiabilties struct {
	SellingLiabilities string `json:"sellingLiabilities"`
	BuyingLiabilities  string `json:"buyingLiabilities"`
}
type ServiceLinkUserInfo struct {
	UserData ServiceLinksUser `json:"userData,omitempty"`
	Wallet   []Balance        `json:"wallet,omitempty"`
}

type ServiceLinkRequestInput struct {
	AuthDescription   string `json:"authDescription,omitempty"`
	DeviceInfo        string `json:"deviceInfo,omitempty"`
	CallbackURL       string `json:"callbackUrl,omitempty"`
	ValidityInMinutes int    `json:"validityInMinutes,omitempty"`
}
type ServiceLinkEventRequestInput struct {
	EventDescription  string `json:"eventDescription,omitempty"`
	DeviceInfo        string `json:"deviceInfo,omitempty"`
	CallbackURL       string `json:"callbackUrl,omitempty"`
	ValidityInMinutes int    `json:"validityInMinutes,omitempty"`
}

type ServiceLinkPushNotificationInput struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	ImageURI string `json:"imageUri"`
	Action   string `json:"action"` //login, payment, 2fa, event
}
type PayWithTrovoWalletData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}
type ReferralLinkData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
}
type LoginWithTrovoWalletData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	LoginID     string `json:"loginId"`
}

type TrovoWalletAuthorizationData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	AuthID      string `json:"authId"`
}

type TrovoWalletAuthorizationRespomseData struct {
	Message string `json:"message"`
}
type TrovoWalletEventData struct {
	DynamicLink string `json:"dynamicLink"`
	QRCode      string `json:"qrCode"`
	EventID     string `json:"eventId"`
}
type LoginVerifyResponse struct {
	AccessToken  string              `json:"accessToken"`
	RefreshToken string              `json:"refreshToken"`
	UserInfo     ServiceLinkUserInfo `json:"userInfo"`
}

type PNSResponse struct {
	Data string `json:"data"`
}
type JwtVerifyResponse struct {
	Message string `json:"message"`
	Token   *Token `json:"token"`
}
type JwtDeleteResponse struct {
	Message string `json:"message"`
}

type JwtRefreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

type Token struct {
	Raw    string `json:"Raw"`
	Method struct {
		Name string `json:"Name"`
		Hash int    `json:"Hash"`
	} `json:"Method"`
	Header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	} `json:"Header"`
	Claims struct {
		AccessUUID string `json:"access_uuid"`
		Authorized bool   `json:"authorized"`
		Exp        int    `json:"exp"`
		UserID     string `json:"user_id"`
	} `json:"Claims"`
	Signature string `json:"Signature"`
	Valid     bool   `json:"Valid"`
}

type TokenizedAssetJSON struct {
	ID                                          string    `json:"id"`
	CreatedAt                                   time.Time `json:"createdAt"`
	UpdatedAt                                   time.Time `json:"updatedAt"`
	AssetSector                                 string    `json:"assetSector"`
	AssetSubSector                              string    `json:"assetSubSector"`
	AssetType                                   string    `json:"assetType"`
	AssetName                                   string    `json:"assetName"`
	OfferingType                                string    `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	SecApproval                                 int       `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                         string    `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                      string    `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                          string    `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                          string    `json:"marketMakingWallet"`
	AssetDescription                            string    `json:"assetDescription"`
	AssetCountryLocation                        string    `json:"assetCountryLocation"`
	AssetPhysicalAddress                        string    `json:"assetPhysicalAddress"`
	AssetLongitude                              string    `json:"assetLongitude"`
	AssetLatitude                               string    `json:"assetLatitude"`
	OwnershipType                               string    `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                               string    `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress          string    `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              string    `json:"assetOwnerName"`
	AssetOwnerAddress                           string    `json:"assetOwnerAddress"`
	AssetQuoteCurrency                          string    `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64   `gorm:"default:0" json:"assetCurrentValue"`
	AssetOwnerRetainedOrContributedValue        float64   `gorm:"default:0" json:"assetOwnerRetainedOrContributedValue"`
	AssetMscCostOutisdeOfValuation              float64   `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                       float64   `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                           string    `json:"protectionMethods"` //csv format
	InsuranceCompanyName                        string    `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       string    `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       string    `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64   `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int       `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int       `gorm:"default:1" json:"assetAlreadyExists"`
	AssetCode                                   string    `gorm:"size:12; index:idx_unique_tokenized_asset_code,unique" json:"assetCode"`
	AssetLogo                                   string    `json:"assetLogo"`
	AssetWebsite                                string    `gorm:"default:'trovotech.io'" json:"assetWebsite"`
	NumberOfTokenToBeIssued                     float64   `gorm:"default:0" json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale            float64   `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	NumberOfTokenToBeSold                       float64   `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64   `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                string    `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                               float64   `gorm:"default:0" json:"pricePerToken"`
	SalesStart                                  time.Time `json:"salesStart"`
	SalesEnd                                    time.Time `json:"salesEnd"`
	CapOnPurchase                               float64   `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                 float64   `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays                           int       `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                string    `gorm:"size:50" json:"proceedCycle"`
	ProceedPayoutCurrency                       string    `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int       `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                           string    `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int       `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   string    `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int       `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                     int       `gorm:"default:0" json:"assetTokenizationStatus"`
	AgreeTransferTitleToCustodian               int       `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int       `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int       `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int       `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int       `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int       `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int       `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   string    `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int       `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int       `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int       `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int       `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int       `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int       `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int       `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        string    `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                string    `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                            string    `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                           int       `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                    int       `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                         int       `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                    int       `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                       int       `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                        int       `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts             int       `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities int       `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                  int       `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                 int       `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                      int       `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                    int       `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements     int       `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
}

// Tokenization Status Constants
// Configurable status values for easy updates
var (
	// TokenizationStatusPendingValues defines status values for pending (0, 1, or 2)
	TokenizationStatusPendingValues = []int{0, 1, 2}
	// TokenizationStatusApproved defines status value for approved (5)
	TokenizationStatusApproved = 5
	// TokenizationStatusRejectedMin defines minimum status value for rejected (> 8, so we use 8 in WHERE clause)
	TokenizationStatusRejectedMin = 8
)
