package sharedconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"trovo-wallet-api/internal/cache"

	cs "cloud.google.com/go/storage"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const TOKEN_LIMIT float64 = 922337203685.00

type GlobalConfig struct {
	// DynamicLinkServiceURLChan  chan string
	PushNotificationClient     *messaging.Client
	FirebaseStorageUploader    *ClientUploader
	PNSContext                 context.Context
	RedisCache                 *cache.RedisCache
	DB                         *gorm.DB
	RoachDB                    *gorm.DB
	BantuExpansionClient       *horizonclient.Client
	BantuNetworkPassphrase     string
	ChannelAccounts            chan *keypair.Full
	InUseChannelAccounts       map[string]*keypair.Full
	Mutex                      sync.Mutex
	ChannelOfTokenizedAssetIDs chan string
}

type ClientUploader struct {
	Client     *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
}
type KYCConfig struct {
	ID              uint64 `json:"-"`
	ServiceProvider string `json:"serviceProvider"`
	Token           string `json:"token"`
	SecretKey       string `json:"secretKey"`
}

type KycWebhookRequest struct {
	ID              uint64
	ServiceProvider string
	Data            string
}
type Asset struct {
	AssetCode   string `json:"assetCode"`
	AssetIssuer string `json:"assetIssuer"`
}

type OfferVolume struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}
type OrderBook struct {
	Bids     []OfferVolume `json:"bids"`
	Asks     []OfferVolume `json:"asks"`
	Asset    Asset         `json:"asset"`
	Currency Asset         `json:"currency"`
}

// CuratedAsset model struct for CuratedAsset.
type CuratedAsset struct {
	ID                          uint64    `gorm:"primaryKey" json:"-"`
	CreatedAt                   time.Time `json:"-"`
	UpdatedAt                   time.Time `json:"-"`
	AssetCode                   string    `gorm:"size:12;unique;not null; default:''" json:"assetCode"`
	AssetName                   string    `gorm:"size:50;null; default:''" json:"assetName"`
	AssetIssuer                 string    `gorm:"size:56;not null; default:''" json:"assetIssuer"`
	Description                 string    `gorm:"not null" json:"description"`
	ImageURL                    *string   `gorm:"null" json:"imageUrl"`
	Website                     string    `gorm:"null;size:100" json:"website"`
	AssetConditions             string    `gorm:"null;size:100" json:"assetConditions"`
	AssetLimit                  float64   `gorm:"type:integer;not null;default:0" json:"assetLimit"` //0 = unlimited
	AssetRedemptionInstructions string    `gorm:"null;" json:"assetRedemptionInstructions"`
	ContactEmail                string    `gorm:"null;size:100" json:"contactEmail"`
	Priority                    uint64    `gorm:"null;" json:"-"`
	AssetClassID                uint64    `gorm:"not null; default:1" json:"assetClassId"`
	Organization                string    `gorm:"null;size:100" json:"organization"`
	Withdrawable                uint64    `gorm:"type:integer;not null;default:0" json:"withdrawable"`
	GenerateDepositAddress      uint64    `gorm:"type:integer;not null;default:0" json:"generateDepositAddress"`
	DecimalPlaces               uint64    `gorm:"type:integer;not null;default:7" json:"decimalPlaces"`
	RealAssetImageURL           *string   `gorm:"null;" json:"realAssetImageUrl"`
	Inactive                    uint64    `gorm:"type:integer;not null;default:0" json:"-"`
	ClosedGroup                 *string   `gorm:"null;" json:"closedGroup"`
}

func (c *ClientUploader) UploadFile(fileInput multipart.File, fileName, imageThumbnailURL string) (string, error) {

	// create an id
	id := uuid.New()
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[UploadFile] error getting bucket handle %v: %v\n", c.BucketName, err)
		return "", fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	_, err = sh.Attrs(ctx)

	if err != nil {
		//no bucket with that name exists, create it
		rules := make([]cs.ACLRule, 0)
		rules = append(rules, cs.ACLRule{Entity: "allUsers", Role: "READER"})
		err := sh.Create(ctx, c.ProjectID, &cs.BucketAttrs{ACL: rules})
		if err != nil {
			log.Printf("[UploadFile] error creating bucket handle %v: %v\n", c.BucketName, err)
			return "", fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}
	newImageThumbnailName := c.UploadPath + "/" + id.String() + fileName
	object := sh.Object(newImageThumbnailName)

	if len(imageThumbnailURL) > 3 {
		// ImageThumbnailURL is full https url. strip the unnecessary portion
		oldName := strings.ReplaceAll(imageThumbnailURL, fmt.Sprintf("https://storage.googleapis.com/%v/", c.BucketName), "")
		oldObject := sh.Object(oldName)
		//check if object already exists and delete it.
		if _, err := oldObject.Attrs(ctx); err == nil {
			oldObject.Delete(ctx)

		}
	}
	writer := object.NewWriter(ctx)

	//Set the attribute
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		log.Printf("[UploadFile] error uploading file %v: %v\n", newImageThumbnailName, err)
		return "", fmt.Errorf("error uploading file %v: %v", newImageThumbnailName, err)
	}

	return newImageThumbnailName, nil
}

func (c *ClientUploader) DeleteFile(imageThumbnailURL string) error {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[UploadFile] error getting bucket handle %v: %v\n", c.BucketName, err)
		return fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	_, err = sh.Attrs(ctx)

	if err != nil {
		//no bucket with that name exists, create it
		rules := make([]cs.ACLRule, 0)
		rules = append(rules, cs.ACLRule{Entity: "allUsers", Role: "READER"})
		err := sh.Create(ctx, c.ProjectID, &cs.BucketAttrs{ACL: rules})
		if err != nil {
			log.Printf("[UploadFile] error creating bucket handle %v: %v\n", c.BucketName, err)
			return fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}

	if len(imageThumbnailURL) > 3 {
		// ImageThumbnailURL is full https url. strip the unnecessary portion
		oldName := strings.ReplaceAll(imageThumbnailURL, fmt.Sprintf("https://storage.googleapis.com/%v/", c.BucketName), "")
		oldObject := sh.Object(oldName)
		//check if object already exists and delete it.
		if _, err := oldObject.Attrs(ctx); err == nil {
			oldObject.Delete(ctx)

		}
	}

	return nil
}

func (c *ClientUploader) SaveQrCodeAsFileToCloud(fileInput *os.File, fileName, imageThumbnailURL string) (string, error) {
	fileName = strings.ReplaceAll(fileName, "/tmp/", "")
	fileName = strings.ReplaceAll(fileName, "/", "")
	fileName = strings.ReplaceAll(fileName, "tmp", "")
	// create an id
	id := uuid.New()
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sh, err := c.Client.Bucket(c.BucketName)
	if err != nil {
		//no bucket with that name exists
		log.Printf("[SaveQrCodeAsFileToCloud] error getting bucket handle %v: %v\n", c.BucketName, err)
		return "", fmt.Errorf("error getting bucket handle %v: %v", c.BucketName, err)
	}

	_, err = sh.Attrs(ctx)

	if err != nil {
		//no bucket with that name exists, create it
		rules := make([]cs.ACLRule, 0)
		rules = append(rules, cs.ACLRule{Entity: "allUsers", Role: "READER"})
		err := sh.Create(ctx, c.ProjectID, &cs.BucketAttrs{ACL: rules})
		if err != nil {
			log.Printf("[SaveQrCodeAsFileToCloud] error creating bucket handle %v: %v\n", c.BucketName, err)
			return "", fmt.Errorf("error creating bucket handle %v: %v", c.BucketName, err)
		}
	}
	newImageThumbnailName := c.UploadPath + "/" + id.String() + fileName
	object := sh.Object(newImageThumbnailName)

	if len(imageThumbnailURL) > 3 {
		// ImageThumbnailURL is full https url. strip the unnecessary portion
		oldName := strings.ReplaceAll(imageThumbnailURL, fmt.Sprintf("https://storage.googleapis.com/%v/", c.BucketName), "")
		oldObject := sh.Object(oldName)
		//check if object already exists and delete it.
		if _, err := oldObject.Attrs(ctx); err == nil {
			oldObject.Delete(ctx)

		}
	}
	writer := object.NewWriter(ctx)

	//Set the attribute
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		log.Printf("[SaveQrCodeAsFileToCloud] error uploading file %v: %v\n", newImageThumbnailName, err)
		return "", fmt.Errorf("error uploading file %v: %v", newImageThumbnailName, err)
	}

	return newImageThumbnailName, nil
}

func (gc *GlobalConfig) ReleaseInUseChannelAccount(pk string) {
	if len(pk) == 0 {
		return
	}
	gc.Mutex.Lock()
	defer gc.Mutex.Unlock()
	ca, ok := gc.InUseChannelAccounts[pk]
	if ok {
		gc.ChannelAccounts <- ca
	}
	delete(gc.InUseChannelAccounts, pk)
}

func (gc *GlobalConfig) StoreInUseChannelAccount(kp *keypair.Full) {
	if kp == nil {
		return
	}
	gc.Mutex.Lock()
	defer gc.Mutex.Unlock()
	gc.InUseChannelAccounts[kp.Address()] = kp
}

// IsValidTokenizedAsset checks if the tokenized asset has bcome market ready at least....with status > 3
func (gc *GlobalConfig) IsValidTokenizedAsset(assetCode string) bool {
	type Result struct {
		ID string
	}
	var result Result
	gc.DB.Raw("SELECT id FROM Tokenized_Assets WHERE Asset_Tokenization_Status > 3 AND Asset_Code = upper(?)", assetCode).Scan(&result)

	return len(result.ID) > 0
}

// IsTokenizedAssetInPrimarySales checks if the status is 5 (primary sales)
func (gc *GlobalConfig) IsTokenizedAssetInPrimarySales(assetCode string) bool {
	type Result struct {
		ID string
	}
	var result Result
	gc.DB.Raw("SELECT id FROM Tokenized_Assets WHERE Asset_Tokenization_Status = 5 AND Asset_Code = upper(?)", assetCode).Scan(&result)

	return len(result.ID) > 0
}

type TokenizedAsset struct {
	ID                                           string     `json:"id"`
	CreatedAt                                    time.Time  `json:"createdAt"`
	UpdatedAt                                    time.Time  `json:"updatedAt"`
	InitiatorUsername                            string     `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                                  *string    `json:"assetSector"`
	AssetSubSector                               *string    `json:"assetSubSector"`
	AssetType                                    *string    `json:"assetType"`
	AssetName                                    *string    `json:"assetName"`
	ApprovedAssetCustodianID                     uint64     `gorm:"not null;default:0" json:"approvedAssetCustodianId"`
	OfferingType                                 *string    `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                                *string    `gorm:"null" json:"closedGroupId"`
	SecApproval                                  int        `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                          *string    `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                       *string    `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                           *string    `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                           *string    `json:"marketMakingWallet"`
	AssetDescription                             *string    `json:"assetDescription"`
	AssetCountryLocation                         *string    `gorm:"not null;size:2;default'NG'" json:"assetCountryLocation"`
	AssetPhysicalAddress                         *string    `json:"assetPhysicalAddress"`
	AssetLongitude                               *string    `json:"assetLongitude"`
	AssetLatitude                                *string    `json:"assetLatitude"`
	OwnershipType                                *string    `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                                *string    `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress           *string    `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                               *string    `json:"assetOwnerName"`
	AssetOwnerAddress                            *string    `json:"assetOwnerAddress"`
	AssetManagerID                               uint64     `gorm:"default:0" json:"assetManagerId"`
	AssetIssuingHouseID                          uint64     `gorm:"not null;default:0" json:"assetIssuingHouseId"`
	LegalAndProfesionalPartnerID                 uint64     `gorm:"not null;default:0" json:"legalAndProfesionalPartnerId"`
	RatingAgencyID                               uint64     `gorm:"not null;default:0" json:"ratingAgencyId"`
	TrusteeID                                    uint64     `gorm:"not null;default:0" json:"trusteeId"`
	AssetQuoteCurrency                           *string    `gorm:"default:'CNGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                            float64    `gorm:"default:0" json:"assetCurrentValue"`
	AssetOwnerRetainedOrContributedValue         float64    `gorm:"default:0" json:"assetOwnerRetainedOrContributedValue"`
	AssetMscCostOutisdeOfValuation               float64    `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                        float64    `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                            *string    `json:"protectionMethods"` //csv format
	InsuranceCompanyName                         *string    `json:"insuranceCompanyName"`
	InsurancePolicyNumber                        *string    `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                        *string    `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                   float64    `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances               int        `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                           int        `gorm:"default:0" json:"assetAlreadyExists"`
	VettingStatus                                int        `gorm:"default:0" json:"vettingStatus"`
	DueDiligenceFail                             int        `gorm:"default:0" json:"dueDiligenceFail"` //0=False(success/in-progress), 1= true (failed).
	DueDiligenceFailureReason                    *string    `json:"dueDiligenceFailureReason"`
	AssetCode                                    *string    `gorm:"size:12; index:idx_unique_tokenized_asset_code,unique" json:"assetCode"`
	AssetLogo                                    *string    `json:"assetLogo"`
	AssetWebsite                                 *string    `gorm:"default:'trovotech.io'" json:"assetWebsite"`
	InitialValueOfTokenizedAsset                 float64    `gorm:"default:0" json:"initialValueOfTokenizedAsset"`
	NumberOfTokenToBeIssued                      float64    `gorm:"default:0" json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale             float64    `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                   float64    `gorm:"default:0" json:"feeInAsset"`
	FeeInAssetPercent                            float64    `gorm:"default:0" json:"feeInAssetPercent"`
	FeeInFiat                                    float64    `gorm:"default:0" json:"feeInFiat"` //tokenization fee in fiat
	NumberOfTokenToBeSold                        float64    `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                      float64    `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                 *string    `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                                float64    `gorm:"default:0" json:"pricePerToken"`
	SalesStart                                   time.Time  `json:"salesStart"`
	SalesEnd                                     time.Time  `json:"salesEnd"`
	CapOnPurchase                                float64    `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                  float64    `gorm:"default:0" json:"capQuantity"`
	CapAmountInFiat                              float64    `gorm:"default:0" json:"capAmountInFiat"`
	CapDurationInDays                            int        `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                 *string    `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                            *uint64    `gorm:"default:0" json:"tokenizationFeeId"`
	ExcludeSecFee                                int        `gorm:"default:0" json:"excludeSecFee"`
	SECTokenizationFeePercent                    float64    `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeFixed                      float64    `gorm:"default:0" json:"SECTokenizationFeeFixed"`
	SECTokenizationFeeValue                      float64    `gorm:"default:0" json:"SECTokenizationFeeValue"`
	CustodianFeePercent                          float64    `gorm:"default:0" json:"custodianFeePercent"`
	CustodianFeeFixed                            float64    `gorm:"default:0" json:"custodianFeeFixed"`
	CustodianFeeValue                            float64    `gorm:"default:0" json:"custodianFeeValue"`
	AssetManagerFeeValue                         float64    `gorm:"default:0" json:"assetManagerFeeValue"`
	AssetManagerFeePercent                       float64    `gorm:"default:0" json:"assetManagerFeePercent"`
	AssetManagerFeeFixed                         float64    `gorm:"default:0" json:"assetManagerFeeFixed"`
	IssuingHouseFeeValue                         float64    `gorm:"default:0" json:"issuingHouseFeeValue"`
	IssuingHouseFeePercent                       float64    `gorm:"default:0" json:"issuingHouseFeePercent"`
	IssuingHouseFeeFixed                         float64    `gorm:"default:0" json:"issuingHouseFeeFixed"`
	LegalAndProfessionalFeePercent               float64    `gorm:"default:0" json:"legalAndProfessionalFeePercent"`
	LegalAndProfessionalFeeFixed                 float64    `gorm:"default:0" json:"legalAndProfessionalFeeFixed"`
	LegalAndProfessionalFeeValue                 float64    `gorm:"default:0" json:"legalAndProfessionalFeeValue"`
	RatingAgencyFeePercent                       float64    `gorm:"default:0" json:"ratingAgencyFeePercent"`
	RatingAgencyFeeFixed                         float64    `gorm:"default:0" json:"ratingAgencyFeeFixed"`
	RatingAgencyFeeValue                         float64    `gorm:"default:0" json:"ratingAgencyFeeValue"`
	TrusteeFeePercent                            float64    `gorm:"default:0" json:"trusteeFeePercent"`
	TrusteeFeeFixed                              float64    `gorm:"default:0" json:"trusteeFeeFixed"`
	TrusteeFeeValue                              float64    `gorm:"default:0" json:"trusteeFeeValue"`
	VATPercent                                   float64    `gorm:"default:0" json:"vatPercent"`
	VATValue                                     float64    `gorm:"default:0" json:"vatValue"`
	VATInAsset                                   float64    `gorm:"default:0" json:"vatInAsset"`
	ProceedPayoutCurrency                        *string    `json:"proceedPayoutCurrency"`
	ProceedPayoutType                            int        `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                            *string    `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                 int        `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                    *string    `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired                int        `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                      int        `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                                *string    `gorm:"null" json:"lastUpdatedBy"`
	TokenizationTransaction                      *string    `gorm:"null" json:"-"`
	AgreeTransferTitleToCustodian                int        `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees           int        `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond                int        `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                     int        `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                     int        `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments         int        `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees     int        `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                    *string    `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                        int        `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                    int        `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl                int        `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems           int        `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel      int        `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity            int        `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections     int        `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                         *string    `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                 *string    `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                             *string    `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                            int        `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                     int        `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                          int        `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                     int        `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                        int        `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                         int        `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts              int        `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities  int        `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                   int        `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                  int        `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                       int        `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                     int        `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements      int        `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
	MintingInitators                             *string    `gorm:"null" json:"mintingInitators"` //CSV of approvers
	MintingApprovers                             *string    `gorm:"null" json:"mintingApprovers"` //csv of initators
	BankID                                       *uint64    `gorm:"null" json:"bankId"`
	AccountNumber                                *string    `gorm:"null" json:"accountNumber"`
	BeneficiaryName                              *string    `gorm:"null" json:"beneficiaryName"`
	TokenizationApplicationFee                   float64    `gorm:"default:0" json:"tokenizationApplicationFee"`
	TokenizationApplicationFeeAsset              string     `gorm:"not null;default:'TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ'" json:"tokenizationApplicationFeeAsset"`
	ProjectStrategicObjectives                   *string    `json:"projectStrategicObjectives"`
	ProjectDevelopmentTimeline                   *string    `json:"projectDevelopmentTimeline"`
	ProjectKeyMilestoneAndDates                  *string    `json:"projectKeyMilestoneAndDates"`
	ProjectScope                                 *string    `json:"projectScope"`
	ProjectEconomicBenefits                      *string    `json:"projectEconomicBenefits"`
	ProjectExpectedNoOfJobs                      int        `json:"projectExpectedNoOfJobs"`
	ProjectIntendedSocialBenefits                *string    `json:"projectIntendedSocialBenefits"`
	ProjectTechnicalPartners                     *string    `json:"projectTechnicalPartners"`
	ProjectFinancialPartners                     *string    `json:"projectFinancialPartners"`
	EstimatedProjectIRR                          float64    `json:"estimatedProjectIRR"`
	EstimatedProjectROI                          float64    `json:"estimatedProjectROI"`
	EstimatedProjectNPV                          float64    `json:"estimatedProjectNPV"`
	EstimatedProjectPaybackPeriodsInMonths       int        `json:"estimatedProjectPaybackPeriodsInMonths"`
	KeyAssumptionsList                           *string    `json:"keyAssumptionsList"`
	ProjectIdentifiedLegalRisks                  *string    `json:"projectIdentifiedLegalRisks"`
	ProjectIdentifiedRegulatoryRisks             *string    `json:"projectIdentifiedRegulatoryRisks"`
	ProjectIdentifiedOperationalOrExecutionRisks *string    `json:"projectIdentifiedOperationalOrExecutionRisks"`
	ProjectIdentifiedMarketRisks                 *string    `json:"projectIdentifiedMarketRisks"`
	ProjectIdentifiedOtherRelevantRisks          *string    `json:"projectIdentifiedOtherRelevantRisks"`
	DeepLink                                     *string    `gorm:"null" json:"deepLink"`
	IsinOrSerialNumber                           *string    `json:"isinOrSerialNumber"` /// begin the mutual fund fields
	InstrumentName                               *string    `json:"instrumentName"`
	InstrumentType                               *string    `json:"instrumentType"`
	TotalIssueSize                               float64    `gorm:"default:0" json:"totalIssueSize"`
	IssueDate                                    time.Time  `json:"issueDate"`
	MaturityDate                                 time.Time  `json:"maturityDate"`
	FaceValuePerUnit                             float64    `gorm:"default:0" json:"faceValuePerUnit"`
	CouponOrInterestRateType                     *string    `json:"couponOrInterestRateType"`
	CouponOrInterestRate                         float64    `json:"couponOrInterestRate"`
	ReferenceIndex                               *string    `json:"referenceIndex"`
	ExpectedYield                                float64    `gorm:"default:0" json:"expectedYield"`
	EarlyRedemptionOptionInvestor                *string    `json:"earlyRedemptionOptionInvestor"`
	MinimumInvestmentAmount                      float64    `json:"minimumInvestmentAmount"`
	TaxTreatmentTokenHolders                     *string    `json:"taxTreatmentTokenHolders"`
	PaymentStructureToTokenHolders               *string    `json:"paymentStructureToTokenHolders"`
	RedemptionMethod                             *string    `json:"redemptionMethod"`
	PaymentStructure                             *string    `json:"paymentStructure"`
	RepaymentMethod                              *string    `json:"repaymentMethod"`
	PaymentCycle                                 *string    `json:"paymentCycle"`
	WeightedAverageLife                          float64    `json:"weightedAverageLife"`
	UnderlyingAssetPoolSize                      float64    `json:"underlyingAssetPoolSize"`
	PoolComposition                              *string    `json:"poolComposition"`
	CreditEnhancementMethod                      *string    `json:"creditEnhancementMethod"`
	SummaryOfUseOfProceeds                       *string    `json:"summaryOfUseOfProceeds"`
	CreditRatingIfAny                            *string    `json:"creditRatingIfAny"`
	IssuerName                                   *string    `json:"issuerName"`
	IssuerType                                   *string    `json:"issuerType"`
	IssuerContactPerson                          *string    `json:"issuerContactPerson"`
	ContactEmail                                 *string    `json:"contactEmail"`
	ContactPhoneNumber                           *string    `json:"contactPhoneNumber"`
	BriefCompanyOverview                         *string    `json:"briefCompanyOverview"`
	MortgageOriginators                          *string    `json:"mortgageOriginators"`
	Servicer                                     *string    `json:"servicer"`
	InstrumentTrustee                            *string    `json:"instrumentTrustee"`
	InstrumentCustodians                         *string    `json:"instrumentCustodian"`
	InstrumentLegalAdvisor                       *string    `json:"instrumentLegalAdvisor"`
	SpecialPurposeVehicle                        *string    `json:"specialPurposeVehicle"`
	InstrumentAssetManagerOrAdministrator        *string    `json:"instrumentAssetManagerOrAdministrator"`
	UnderwriterIfAny                             *string    `json:"underwriterIfAny"`
	InstrumentCreditRatingAgency                 *string    `json:"instrumentCreditRatingAgency"`
	AuditorOrVerifier                            *string    `json:"auditorOrVerifier"`
	CreditRiskAssessment                         *string    `json:"creditRiskAssessment"`
	CreditRating                                 *string    `json:"creditRating"`
	PrepaymentRisk                               *string    `json:"prepaymentRisk"`
	InterestRateRisk                             *string    `json:"interestRateRisk"`
	StructuralComplexityRisk                     *string    `json:"structuralComplexityRisk"`
	LegalOrRegulatoryRisk                        *string    `json:"legalOrRegulatoryRisk"`
	OperationalRisk                              *string    `json:"operationalRisk"`
	MarketRisk                                   *string    `json:"marketRisk"`
	EsgRisk                                      *string    `json:"esgRisk"`
	MitigationMeasures                           *string    `json:"mitigationMeasures"`
	IssuingAuthority                             *string    `json:"issuingAuthority"`
	RegulatoryApprovalId                         *string    `json:"regulatoryApprovalId"`
	LicenseApprovalReferenceNumber               *string    `json:"licenseApprovalReferenceNumber"`
	ListingStatus                                *string    `json:"listingStatus"`
	CurrencyOfIssuance                           *string    `json:"currencyOfIssuance"`
	CouponRateType                               *string    `json:"couponRateType"`
	CouponRate                                   float64    `gorm:"default:0" json:"couponRate"`
	SpreadOrMargin                               float64    `gorm:"default:0" json:"spreadOrMargin"`
	ResetFrequency                               *string    `json:"resetFrequency"`
	CouponPaymentFrequency                       *string    `json:"couponPaymentFrequency"`
	RedemptionStructure                          *string    `json:"redemptionStructure"`
	EarlyRedemptionOption                        *string    `json:"earlyRedemptionOption"`
	EarlyRedemptionPenalty                       *string    `json:"earlyRedemptionPenalty"`
	TaxTreatment                                 *string    `json:"taxTreatment"`
	NavOrMarketValueUpdates                      *string    `json:"navOrMarketValueUpdates"`
	ImpactMetrics                                *string    `json:"impactMetrics"`
	LegalBacking                                 *string    `json:"legalBacking"`
	DefaultHistory                               *string    `json:"defaultHistory"`
	RiskFactorsSummary                           *string    `json:"riskFactorsSummary"`
	PayingAgent                                  *string    `json:"payingAgent"`
	Auditor                                      *string    `json:"auditor"`
	RegistrarOrCSCSAgent                         *string    `json:"registrarOrCSCSAgent"`
	FundStructure                                *string    `json:"fundStructure"`
	AssetManagementCompanyName                   *string    `json:"assetManagementCompanyName"`
	FundManagers                                 *string    `json:"fundManagers"`
	RegulatoryLicenseNumber                      *string    `json:"regulatoryLicenseNumber"`
	FundLaunchDate                               *time.Time `json:"fundLaunchDate"`
	TotalExpenseRatio                            float64    `gorm:"default:0" json:"totalExpenseRatio"`
	ExitLoadRedemptionFee                        float64    `json:"exitLoadRedemptionFee"`
	Tenure                                       int        `gorm:"default:0" json:"tenure"`
	InitialNetAssetValue                         float64    `json:"initialNetAssetValue"`
	NavUpdateFrequency                           *string    `json:"navUpdateFrequency"`
	NavCalculationMethod                         *string    `json:"navCalculationMethod"`
	RedemptionRules                              *string    `json:"redemptionRules"`
	LockInPeriod                                 int        `gorm:"default:0" json:"lockInPeriod"`
	EntryLoad                                    float64    `gorm:"default:0" json:"entryLoad"`
	PerformanceFee                               float64    `gorm:"default:0" json:"performanceFee"`
	DividendPolicy                               *string    `json:"dividendPolicy"`
	LiquidityProfile                             *string    `json:"liquidityProfile"`
	DistributionFrequency                        *string    `json:"distributionFrequency"`
	DistributionMethod                           *string    `json:"distributionMethod"`
	BenchmarkComparisonMethod                    *string    `json:"benchmarkComparisonMethod"`
	FeeBreakdownSummary                          *string    `json:"feeBreakdownSummary"`
	InvestmentObjective                          *string    `json:"investmentObjective"`
	EquityStrategy                               *string    `json:"equityStrategy"`
	MarketCapitalizationFocus                    *string    `json:"marketCapitalizationFocus"`
	BenchmarkIndex                               *string    `json:"benchmarkIndex"`
	SectorExposureLimits                         *string    `json:"sectorExposureLimits"`
	TopHoldings                                  *string    `json:"topHoldings"`
	GeographicExposure                           *string    `json:"geographicExposure"`
	RiskProfile                                  *string    `json:"riskProfile"`
	VolatilityEstimate                           float64    `gorm:"default:0" json:"volatilityEstimate"`
	DividendYield                                float64    `gorm:"default:0" json:"dividendYield"`
	TrusteeName                                  *string    `json:"trusteeName"`
	FundAdministrator                            *string    `json:"fundAdministrator"`
	InvestmentCommitteeMembers                   *string    `json:"investmentCommitteeMembers"`
	IsinOrSecFundCode                            *string    `json:"isinOrSecFundCode"`
	ExitLoadOrRedemptionFee                      float64    `json:"exitLoadOrRedemptionFee"`
	FundRiskRating                               *string    `json:"fundRiskRating"`
	AssetAllocation                              *string    `json:"assetAllocation"`
	AverageMaturity                              float64    `gorm:"default:0" json:"averageMaturity"`
	YieldToMaturity                              float64    `gorm:"default:0" json:"yieldToMaturity"`
	CreditRatingProfile                          *string    `json:"creditRatingProfile"`
	LockInPeriodPortfolio                        int        `gorm:"default:0" json:"lockInPeriodPortfolio"`
	PerformanceFeeIfAny                          float64    `json:"performanceFeeIfAny"`
	TopEquityHoldings                            *string    `json:"topEquityHoldings"`
	TargetAllocation                             *string    `json:"targetAllocation"`
	AllowedAllocationRange                       *string    `json:"allowedAllocationRange"`
	AssetClassesIncluded                         *string    `json:"assetClassesIncluded"`
	RebalancingFrequency                         *string    `json:"rebalancingFrequency"`
	BenchmarkIndexComposite                      *string    `json:"benchmarkIndexComposite"`
	TopEquityHoldingsList                        *string    `json:"topEquityHoldingsList"`
	TopDebtHoldings                              *string    `json:"topDebtHoldings"`
	CreditRatingDistribution                     *string    `json:"creditRatingDistribution"`
	AverageMaturityHybrid                        float64    `json:"averageMaturityHybrid"`
	YieldToMaturityHybrid                        float64    `json:"yieldToMaturityHybrid"`
	TitleOfIssuance                              *string    `json:"titleOfIssuance"`
	TypeOfCommercialPaper                        *string    `json:"typeOfCommercialPaper"`
	PricingYield                                 float64    `gorm:"default:0" json:"pricingYield"`
	UseOfProceeds                                *string    `json:"useOfProceeds"`
	Ranking                                      *string    `json:"ranking"`
	BackingSecurity                              *string    `json:"backingSecurity"`
	IssuerRegistrationNumber                     *string    `json:"issuerRegistrationNumber"`
	IncorporationDate                            *time.Time `json:"incorporationDate"`
	RcNumber                                     *string    `json:"rcNumber"`
	TaxIdNumber                                  *string    `json:"taxIdNumber"`
	OfficeAddress                                *string    `json:"officeAddress"`
	Rating                                       *string    `json:"rating"`
	PartiesInvolvedIssuer                        *string    `json:"partiesInvolvedIssuer"`
	PartiesInvolvedArranger                      *string    `json:"partiesInvolvedArranger"`
	PartiesInvolvedLegalAdviser                  *string    `json:"partiesInvolvedLegalAdviser"`
	PartiesInvolvedAuditor                       *string    `json:"partiesInvolvedAuditor"`
	PartiesInvolvedRatingAgency                  *string    `json:"partiesInvolvedRatingAgency"`
	PartiesInvolvedCustodian                     *string    `json:"partiesInvolvedCustodian"`
	PartiesInvolvedTrustee                       *string    `json:"partiesInvolvedTrustee"`
	PartiesInvolvedAuditorVerifier               *string    `json:"partiesInvolvedAuditorVerifier"`
	SecurityRiskLegalBacking                     *string    `json:"securityRiskLegalBacking"`
	SecurityRiskCollateral                       *string    `json:"securityRiskCollateral"`
	SecurityRiskDefaultHistory                   *string    `json:"securityRiskDefaultHistory"`
	SecurityRiskCreditRating                     *string    `json:"securityRiskCreditRating"`
	SecurityRiskRiskFactorsSummary               *string    `json:"securityRiskRiskFactorsSummary"`
	SecurityRiskBusinessRisk                     *string    `json:"securityRiskBusinessRisk"`
	SecurityRiskDefaultRisk                      *string    `json:"securityRiskDefaultRisk"`
	SecurityRiskLiquidityRisk                    *string    `json:"securityRiskLiquidityRisk"`
	SecurityRiskRegulatoryRisk                   *string    `json:"securityRiskRegulatoryRisk"`
	SecurityRiskMarketRisk                       *string    `json:"securityRiskMarketRisk"`
	SecurityRiskOperationalRisk                  *string    `json:"securityRiskOperationalRisk"`
	SecurityRiskMitigationMeasures               *string    `json:"securityRiskMitigationMeasures"`
	CommodityType                                *string    `json:"commodityType"`
	CommodityDescription                         *string    `json:"commodityDescription"`
	Quantity                                     int        `gorm:"default:0" json:"quantity"`
	QualityGrade                                 *string    `json:"qualityGrade"`
	IssuerContactInfo                            *string    `json:"issuerContactInfo"`
	WarehouseName                                *string    `json:"warehouseName"`
	WarehouseOperatorName                        *string    `json:"warehouseOperatorName"`
	WarehouseLicenseNumber                       *string    `json:"warehouseLicenseNumber"`
	WarehouseLocation                            *string    `json:"warehouseLocation"`
	WrNumber                                     *string    `json:"wrNumber"`
	WrIssueDate                                  *time.Time `json:"wrIssueDate"`
	WrExpiryDate                                 *time.Time `json:"wrExpiryDate"`
	WrSystemRegistration                         *string    `json:"wrSystemRegistration"`
	WrRegistrationNumber                         *string    `json:"wrRegistrationNumber"`
	WrVerifier                                   *string    `json:"wrVerifier"`
	StorageCondition                             *string    `json:"storageCondition"`
	WarehouseAccreditationBody                   *string    `json:"warehouseAccreditationBody"`
	MinimumPurchaseAmount                        float64    `json:"minimumPurchaseAmount"`
	AutoRollover                                 *string    `json:"autoRollover"`
	CurrentBeneficialOwner                       *string    `json:"currentBeneficialOwner"`
	WrCustodianName                              *string    `json:"wrCustodianName"`
	OwnershipRightsRepresented                   *string    `json:"ownershipRightsRepresented"`
	TrusteeOrThirdPartyOversight                 *string    `json:"trusteeOrThirdPartyOversight"`
	LienOrEncumbrances                           *string    `json:"lienOrEncumbrances"`
	AssetValuation                               float64    `gorm:"default:0" json:"assetValuation"`
	ValuationDate                                *time.Time `json:"valuationDate"`
	ValuationMethodology                         *string    `json:"valuationMethodology"`
	TokenizationObjective                        *string    `json:"tokenizationObjective"`
	HoldingPeriod                                int        `gorm:"default:0" json:"holdingPeriod"`
	RedemptionMechanism                          *string    `json:"redemptionMechanism"`
	PartiesInvolvedUnderwriter                   *string    `json:"partiesInvolvedUnderwriter"`
	PartiesInvolvedAssetManager                  *string    `json:"partiesInvolvedAssetManager"`
	PartiesInvolvedLegalAdvisor                  *string    `json:"partiesInvolvedLegalAdvisor"`
	PartiesInvolvedRegulator                     *string    `json:"partiesInvolvedRegulator"`
	RisksMarketRisk                              *string    `json:"risksMarketRisk"`
	RisksStorageRisk                             *string    `json:"risksStorageRisk"`
	RisksTitleRisk                               *string    `json:"risksTitleRisk"`
	RisksFraudRisk                               *string    `json:"risksFraudRisk"`
	RisksInsuranceRisk                           *string    `json:"risksInsuranceRisk"`
	RisksOperationalRisk                         *string    `json:"risksOperationalRisk"`
	RisksRegulatoryRisk                          *string    `json:"risksRegulatoryRisk"`
	RisksLiquidityRisk                           *string    `json:"risksLiquidityRisk"`
	RisksForceMajeureRisk                        *string    `json:"risksForceMajeureRisk"`
	RisksEarlyRedemptionRisk                     *string    `json:"risksEarlyRedemptionRisk"`
	RisksMitigationMeasures                      *string    `json:"risksMitigationMeasures"`
	RisksInsuranceCoverageSummary                *string    `json:"risksInsuranceCoverageSummary"`
	RisksInsuranceProvider                       *string    `json:"risksInsuranceProvider"`
	RisksCoverageValue                           float64    `json:"risksCoverageValue"`
	QuanlityStandard                             *string    `json:"quanlityStandard"`
	IssuerContactInformation                     *string    `json:"issuerContactInformation"`
	VaultCustodianName                           *string    `json:"vaultCustodianName"`
	VaultOperator                                *string    `json:"vaultOperator"`
	VaultLicenseNumber                           *string    `json:"vaultLicenseNumber"`
	VaultLocation                                *string    `json:"vaultLocation"`
	Number                                       *string    `json:"number"`
	IssuerDate                                   *time.Time `json:"issuerDate"`
	ExpiryDate                                   *time.Time `json:"expiryDate"`
	RegistryRecord                               *string    `json:"registryRecord"`
	Verifier                                     *string    `json:"verifier"`
	StorageConditions                            *string    `json:"storageConditions"`
	VaultAccreditationBody                       *string    `json:"vaultAccreditationBody"`
	OwnershipLegalHolder                         *string    `json:"ownershipLegalHolder"`
	OwnershipCustodianName                       *string    `json:"ownershipCustodianName"`
	OwnershipTrustee                             *string    `json:"ownershipTrustee"`
	OwnershipLienOrEncumbrances                  *string    `json:"ownershipLienOrEncumbrances"`
	ValuationAssetValuation                      *string    `json:"valuationAssetValuation"`
	HoldingLockinPeriod                          int        `gorm:"default:0" json:"holdingLockinPeriod"`
	InsuranceMarketRisk                          *string    `json:"insuranceMarketRisk"`
	InsuranceStorageRisk                         *string    `json:"insuranceStorageRisk"`
	InsuranceTitleRisk                           *string    `json:"insuranceTitleRisk"`
	InsuranceFraudRisk                           *string    `json:"insuranceFraudRisk"`
	InsuranceInsuranceRisk                       *string    `json:"insuranceInsuranceRisk"`
	InsuranceOperationalRisk                     *string    `json:"insuranceOperationalRisk"`
	InsuranceRegulatoryRisk                      *string    `json:"insuranceRegulatoryRisk"`
	InsuranceLiquidityRisk                       *string    `json:"insuranceLiquidityRisk"`
	InsuranceForceMajeureRisk                    *string    `json:"insuranceForceMajeureRisk"`
	InsuranceEarlyRedemptionRisk                 *string    `json:"insuranceEarlyRedemptionRisk"`
	InsuranceMitigationMeasures                  *string    `json:"insuranceMitigationMeasures"`
	InsuranceInsuranceCoverageSummary            *string    `json:"insuranceInsuranceCoverageSummary"`
	InsuranceInsuranceProvider                   *string    `json:"insuranceInsuranceProvider"`
	InsuranceCoverageValue                       *string    `json:"insuranceCoverageValue"`
	IssuerRegistrationNo                         *string    `json:"issuerRegistrationNo"`
	SectorAndIndustry                            *string    `json:"sectorAndIndustry"`
	LicenseOrPermitNumber                        *string    `json:"licenseOrPermitNumber"`
	IssuerAdditionalInfo                         *string    `json:"issuerAdditionalInfo"`
	IsinSerialNumber                             *string    `json:"isinSerialNumber"`
	EsgOrImpactMetrics                           *string    `json:"esgOrImpactMetrics"`
	InstrumentAdditionalInfo                     *string    `json:"instrumentAdditionalInfo"`
	SecurityType                                 *string    `json:"securityType"`
	CollateralDescription                        *string    `json:"collateralDescription"`
	CovenantSummary                              *string    `json:"covenantSummary"`
	CovenantTestingFrequency                     *string    `json:"covenantTestingFrequency"`
	EventOfDefaultClauses                        *string    `json:"eventOfDefaultClauses"`
	LegalEnforcementMechanism                    *string    `json:"legalEnforcementMechanism"`
	Guarantee                                    *string    `json:"guarantee"`
	RecoveryEstimate                             float64    `gorm:"default:0" json:"recoveryEstimate"`
	RiskProfileAdditionalInfo                    *string    `json:"riskProfileAdditionalInfo"`
	LicenseNumber                                *string    `json:"licenseNumber"`
	ExitLoadFee                                  float64    `gorm:"default:0" json:"exitLoadFee"`
	EntryLoadFee                                 float64    `gorm:"default:0" json:"entryLoadFee"`
	PortfolioLockInPeriod                        int        `gorm:"default:0" json:"portfolioLockInPeriod"`
	PortfolioPerformanceFee                      float64    `gorm:"default:0" json:"portfolioPerformanceFee"`
	FundingStructure                             int        `gorm:"default:0" json:"fundingStructure"` //0=equity, 1= debt, 2= hybrid
	EquityPercentage                             float64    `gorm:"default:0" json:"equityPercentage"`
	DebtPercentage                               float64    `gorm:"default:0" json:"debtPercentage"`
	DebtInstrumentType                           *string    `json:"debtInstrumentType"`
	PrincipalPaymentMethod                       *string    `json:"principalPaymentMethod"`
	DebtInstrumentRepaymentSource                *string    `json:"debtInstrumentRepaymentSource"`
	DebtInstrumentGuaranteesOrEnhancements       *string    `json:"debtInstrumentGuaranteesOrEnhancements"`
	DebtInstrumentDefaultAndRecoveryTerms        *string    `json:"debtInstrumentDefaultAndRecoveryTerms"`
	DebtInstrumentRepaymentFrequency             *string    `json:"debtInstrumentRepaymentFrequency"`
	InterestRepaymentFrequency                   *string    `json:"interestRepaymentFrequency"`
	DcsrDetails                                  *string    `json:"dcsrDetails"`
	SinkingFundStructure                         *string    `json:"sinkingFundStructure"`
	CovenantMonitoringAgent                      *string    `json:"covenantMonitoringAgent"`
	RightOfRecourse                              *string    `json:"rightOfRecourse"`
	DebtInstrumentInterestRate                   float64    `gorm:"default:0" json:"debtInstrumentInterestRate"`
	DcsrRatio                                    float64    `gorm:"default:0" json:"dcsrRatio"`
	LtvRatio                                     float64    `gorm:"default:0" json:"ltvRatio"`
	InterestCoverageRatio                        float64    `gorm:"default:0" json:"interestCoverageRatio"`
	MaximumLeverageRatio                         float64    `gorm:"default:0" json:"maximumLeverageRatio"`
	GracePeriod                                  int        `gorm:"default:0" json:"gracePeriod"`
	TrusteeAppointed                             int        `gorm:"default:0" json:"trusteeAppointed"`
	ReserveFundInPlace                           int        `gorm:"default:0" json:"reserveFundInPlace"`
	SecurityOrCollateralOffered                  *string    `json:"securityOrCollateralOffered"`
	FundInstrumentType                           *string    `json:"fundInstrumentType"`
	InstrumentRatingAgency                       *string    `json:"instrumentRatingAgency"`
	PortfolioTopHoldings                         *string    `json:"portfolioTopHoldings"`
	CreditRatingAgency                           *string    `json:"CreditRatingAgency"`
	TrusteeRegNumber                             *string    `json:"trusteeRegNumber"`
	CreditEnhancerOrGuarantor                    *string    `json:"creditEnhancerOrGuarantor"`
	BondStructuringAdvisor                       *string    `json:"bondStructuringAdvisor"`
	EntitiesAdditionalInfo                       *string    `json:"entitiesAdditionalInfo"`
	AuthorizedRepresentativeName                 *string    `json:"authorizedRepresentativeName"`
	AuthorizedRepresentativeTitleOrPosition      *string    `json:"authorizedRepresentativeTitleOrPosition"`
	AuthorizedRepresentativeEmail                *string    `json:"authorizedRepresentativeEmail"`
	AcceptTokenizationTermsAndAgreement          int        `gorm:"default:0" json:"acceptTokenizationTermsAndAgreement"`
	AttestInformationAccurateAndVerifiable       int        `gorm:"default:0" json:"attestInformationAccurateAndVerifiable"`
	AcknowledgedSuitabilityCriteria              int        `gorm:"default:0" json:"acknowledgedSuitabilityCriteria"`
	MinimumKycTier                               *string    `json:"minimumKycTier"`
	InvestorCategory                             *string    `json:"investorCategory"`
	WithholdingTaxDisclosure                     int        `gorm:"default:0" json:"withholdingTaxDisclosure"`
	ExitWithFiat                                 int        `gorm:"default:0" json:"exitWithFiat"`
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
	DeepLink                                    string    `gorm:"null" json:"deepLink"`
}

type CountryConfig struct {
	CountryCode                              string  `gorm:"size:2;primaryKey" json:"countryCode"`
	SECTokenizationFeePercent                float64 `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeFixed                  float64 `gorm:"default:0" json:"SECTokenizationFeeFixed"`
	SECTradeFeePercent                       float64 `gorm:"default:0" json:"SECTradeFeePercent"`
	SECTradeFeeFixed                         float64 `gorm:"default:0" json:"SECTradeFeeFixed"`
	RegulatorName                            string  `gorm:"default:'SECURITY AND EXCHANGE COMMISSION'" json:"regulatorName"`
	QuoteCurrencyCode                        string  `json:"quoteCurrencyCode"`
	FiatLabel                                string  `json:"fiatLabel"`
	FiatGlyph                                string  `json:"fiatGlyph"`
	MinTokenizationFee                       float64 `gorm:"default:0" json:"minTokenizationFee"`
	MinTROVBalanceForTokenizationApplication float64 `gorm:"default:600" json:"minTROVBalanceForTokenizationApplication"`
	TokenizationApplicationFee               float64 `gorm:"default:0" json:"tokenizationApplicationFee"`
	TokenizationApplicationFeeAsset          string  `gorm:"default:'TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ'" json:"tokenizationApplicationFeeAsset"`
	VATPercent                               float64 `gorm:"default:0" json:"vatPercent"`
	FiatActivationAmount                     float64 `gorm:"default:1000" json:"fiatActivationAmount"`
	TrovTokenActivationPercent               float64 `gorm:"default:50" json:"trovTokenActivationPercent"` //the rest is for gas/nativetoken
}

// type CountryCode string

func (gc *GlobalConfig) GetConfig(c string) (cConfig CountryConfig) {
	gc.DB.Where("country_code = ?", c).First(&cConfig)
	return
}

func (gc *GlobalConfig) GetTokenizedAssetByCode(assetCode string) (t TokenizedAsset) {

	gc.DB.Where("asset_code = upper(?)", assetCode).First(&t)

	return t
}

//GetCuratedAssetByClassID
/**
  1	"Token"
  2	"Stablecoin"
  3	"Tokenized Asset"
  4	"Non Fungible Token (NFT)"
  5	"Reward"
  **/
func (gc *GlobalConfig) GetCuratedAssetByClassID(assetClassID uint64, includeInactive bool) (tas []CuratedAsset) {
	if includeInactive {
		gc.DB.Where("Asset_Class_ID = ?", assetClassID).Find(&tas)

	} else {
		gc.DB.Where("Asset_Class_ID = ? AND inactive = 0", assetClassID).Find(&tas)

	}

	return tas
}

type PostTokenizationTrustlineCandidate struct {
	ID          uint64
	PublicKey   string
	Description string
}

func (gc *GlobalConfig) GetPostTokenizationTrustlineCandidates() (tas []PostTokenizationTrustlineCandidate) {
	tas = make([]PostTokenizationTrustlineCandidate, 0)
	gc.DB.Find(&tas)
	return tas
}

func (gc *GlobalConfig) GetKycConfig(provider string) (t KYCConfig) {

	gc.DB.Where("service_provider = ?", provider).First(&t)

	return t
}

func (gc *GlobalConfig) SaveKycWebhookData(provider, data string) error {

	t := KycWebhookRequest{
		ServiceProvider: provider,
		Data:            data,
	}

	return gc.DB.Save(&t).Error
}

// DynamicLink is model for saving dybamic links
type DynamicLink struct {
	ID   string `json:"linkId"`
	Link string `json:"link"`
}

func (gc *GlobalConfig) GetLinkFromShortlinkID(linkID string) (t DynamicLink) {

	gc.DB.Where("id = ?", linkID).First(&t)
	return
}

func (gc *GlobalConfig) GetTokenizedAssetByID(tokenizedAssetID string) (t TokenizedAsset) {

	gc.DB.Where("id = ?", tokenizedAssetID).First(&t)

	return t
}

// GetCuratedAssets returns list of Curated Assets
func (gc *GlobalConfig) GetCuratedAssets(includeInactive bool) (assets map[string]CuratedAsset) {
	var fetchedAssets []CuratedAsset
	tempAssets := make(map[string]CuratedAsset)
	assets = make(map[string]CuratedAsset)

	cacheKeyInfo := "curatedAssets_"
	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			// log.Printf("[GetCuratedAssets] %v, served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &assets)
			return
		}

	}

	var m sync.Mutex
	var dberr error
	if !includeInactive {
		dberr = gc.DB.Preload(clause.Associations).Order("priority").Order("asset_code").Where("inactive = ?", 0).Find(&fetchedAssets).Error

	} else {
		dberr = gc.DB.Preload(clause.Associations).Order("priority").Order("asset_code").Find(&fetchedAssets).Error
	}

	if dberr != nil {
		log.Printf("[GetCuratedAssets]error getting assets: %v\n", dberr)
		return assets
	}
	// usdPrice, _ := blockchain.GetXBNDollarAskPrice(db)
	var wg sync.WaitGroup
	for _, v := range fetchedAssets {

		wg.Add(1)
		go func(v CuratedAsset) {
			defer wg.Done()
			//get native price
			// log.Printf(">>>>>>>>>>>>>>>Fetched Asset: Code: %v, Issuer: %v\n", v.AssetCode, v.AssetIssuer)
			// nativePrice, _ := blockchain.GetNativeAskPrice(v.AssetCode, v.AssetIssuer)
			// v.NativePrice = nativePrice
			// usdPriceFloat := decimal.RequireFromString(usdPrice)
			// nativePriceFloat := decimal.RequireFromString(nativePrice)
			// v.UsdPrice = nativePriceFloat.Mul(usdPriceFloat).Truncate(7).String()
			m.Lock()
			tempAssets[v.AssetCode+":"+v.AssetIssuer] = v
			m.Unlock()
		}(v)

	}
	wg.Wait()
	// tempAssets[":"] = models.CuratedAsset{
	// 	ImageURL:     nativeLogo(),
	// 	AssetName:    "Bantu Network Token",
	// 	Description:  "XBN is the native network utility token issued by the Bantu Blockchain Foundation, it is used as gas to power transactions on the blockchain network.",
	// 	Website:      "www.bantufoundation.org",
	// 	ContactEmail: "ops@bantufoundation.org",
	// 	Priority:     1,
	// }
	assets = tempAssets
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, assets, 0)

	return assets
}

// GetCuratedAssets returns list of Curated Assets
func (gc *GlobalConfig) GetCuratedAssetByCode(assetCode string) (asset CuratedAsset) {

	dberr := gc.DB.Preload(clause.Associations).Where("asset_code = ?", assetCode).First(&asset).Error

	if dberr != nil {
		log.Printf("[GetCuratedAssets]error getting asset: %v\n", dberr)
		return asset
	}

	return asset
}

func (gc *GlobalConfig) TokenLimit() float64 {

	return decimal.NewFromFloat(TOKEN_LIMIT).Truncate(7).InexactFloat64()
}

func (gc *GlobalConfig) TokenLimitAsDecimal() decimal.Decimal {

	return decimal.NewFromFloat(TOKEN_LIMIT).Truncate(7)
}

func (gc *GlobalConfig) TokenLimitAsString() string {

	return decimal.NewFromFloat(TOKEN_LIMIT).Truncate(7).String()
}

func (ti *TokenizedAsset) ToJSON() (t TokenizedAssetJSON) {
	titleCaser := cases.Title(language.English)
	t.ID = ti.ID
	t.CreatedAt = ti.CreatedAt
	t.UpdatedAt = ti.UpdatedAt

	t.AssetOwnerRetainedOrContributedValue = ti.AssetOwnerRetainedOrContributedValue

	if ti.InitialOwnerPreferredWalletAddress != nil {

		t.InitialOwnerPreferredWalletAddress = *ti.InitialOwnerPreferredWalletAddress

	}

	if ti.AdditionalKYCRequirements != nil {

		t.AdditionalKYCRequirements = *ti.AdditionalKYCRequirements

	}
	if ti.AssetSector != nil {

		t.AssetSector = *ti.AssetSector
	}

	if ti.AssetSubSector != nil {

		t.AssetSubSector = *ti.AssetSubSector
	}

	if ti.AssetType != nil {

		t.AssetType = *ti.AssetType
	}
	if ti.AssetName != nil {

		t.AssetName = titleCaser.String(*ti.AssetName)

	}
	if ti.AssetWebsite != nil {
		t.AssetWebsite = strings.ToLower(*ti.AssetWebsite)

	}

	if ti.AssetCode != nil {
		t.AssetCode = strings.ToUpper(*ti.AssetCode)

	}
	if ti.AssetLogo != nil {
		t.AssetLogo = *ti.AssetLogo

	}

	if ti.OfferingType != nil {
		t.OfferingType = *ti.OfferingType
	}

	if ti.SecApproval > 0 {
		t.SecApproval = ti.SecApproval
		t.SecApprovalIdNumber = *ti.SecApprovalIdNumber
	}

	if ti.IssuingWalletPublicKey != nil {
		t.IssuingWalletPublicKey = *ti.IssuingWalletPublicKey
	}
	if ti.IssuingWalletAlias != nil {
		t.IssuingWalletAlias = *ti.IssuingWalletAlias
	}
	if ti.MarketMakingWallet != nil {
		t.MarketMakingWallet = *ti.MarketMakingWallet
	}
	if ti.AssetDescription != nil {
		t.AssetDescription = titleCaser.String(*ti.AssetDescription)
	}
	if ti.AssetCountryLocation != nil {
		t.AssetCountryLocation = *ti.AssetCountryLocation
	}

	if ti.AssetPhysicalAddress != nil {
		t.AssetPhysicalAddress = titleCaser.String(*ti.AssetPhysicalAddress)
	}
	if ti.AssetLongitude != nil {
		t.AssetLongitude = *ti.AssetLongitude
	}

	if ti.AssetLatitude != nil {
		t.AssetLatitude = *ti.AssetLatitude
	}
	if ti.OwnershipType != nil {
		t.OwnershipType = *ti.OwnershipType
	}
	if ti.OwnershipKind != nil {
		t.OwnershipKind = *ti.OwnershipKind
	}
	if ti.AssetOwnerName != nil {
		t.AssetOwnerName = strings.ToUpper(*ti.AssetOwnerName)
	}
	if ti.AssetOwnerAddress != nil {
		t.AssetOwnerAddress = titleCaser.String(*ti.AssetOwnerAddress)
	}

	if ti.AssetQuoteCurrency != nil {
		t.AssetQuoteCurrency = *ti.AssetQuoteCurrency
	}
	t.AssetCurrentValue = ti.AssetCurrentValue

	t.AssetMscCostOutisdeOfValuation = ti.AssetMscCostOutisdeOfValuation
	t.ValueOfTokenizedAsset = ti.ValueOfTokenizedAsset

	if ti.ProtectionMethods != nil {
		t.ProtectionMethods = *ti.ProtectionMethods
	}
	if ti.InsuranceCompanyName != nil {
		t.InsuranceCompanyName = *ti.InsuranceCompanyName
	}
	if ti.InsurancePolicyHolder != nil {
		t.InsurancePolicyHolder = *ti.InsurancePolicyHolder
	}
	if ti.InsurancePolicyNumber != nil {
		t.InsurancePolicyNumber = *ti.InsurancePolicyNumber
	}
	t.PercentageValueOfInsurance = ti.PercentageValueOfInsurance
	t.IsFreeFromLiensAndEncumbrances = ti.IsFreeFromLiensAndEncumbrances
	t.AssetAlreadyExists = ti.AssetAlreadyExists
	t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
	t.NumberOfTokenToBeSold = ti.NumberOfTokenToBeSold
	t.TotalTokenHeldByManager = ti.TotalTokenHeldByManager

	if ti.WalletToHoldAssetsNotForSale != nil {
		t.WalletToHoldAssetsNotForSale = *ti.WalletToHoldAssetsNotForSale
	}
	t.PricePerToken = ti.PricePerToken
	t.SalesStart = ti.SalesStart
	t.SalesEnd = ti.SalesEnd
	t.CapOnPurchase = ti.CapOnPurchase
	t.CapQuantity = ti.CapQuantity
	t.CapDurationInDays = ti.CapDurationInDays

	if ti.ProceedCycle != nil {
		t.ProceedCycle = *ti.ProceedCycle
	}
	if ti.ProceedPayoutCurrency != nil {
		t.ProceedPayoutCurrency = *ti.ProceedPayoutCurrency
	}

	if ti.ExemptedCountries != nil {
		t.ExemptedCountries = *ti.ExemptedCountries
	}

	if ti.AdditionalKYCRequirements != nil {

		t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements
		t.AdditionalKYCRequirements = *ti.AdditionalKYCRequirements
	}

	t.InvestorAccreditationRequired = ti.InvestorAccreditationRequired

	t.AssetTokenizationStatus = ti.AssetTokenizationStatus
	t.AgreeTransferTitleToCustodian = ti.AgreeTransferTitleToCustodian
	t.ContractualProtectionRevGuarantees = ti.ContractualProtectionRevGuarantees
	t.ContractualProtectionPerfBond = ti.ContractualProtectionPerfBond
	t.ContractualProtectionSLA = ti.ContractualProtectionSLA
	t.RiskSharingMechanismPPPs = ti.RiskSharingMechanismPPPs
	t.RiskSharingMechanismHedgeInstruments = ti.RiskSharingMechanismHedgeInstruments
	t.RiskSharingMechanismCompletionGuarantees = ti.RiskSharingMechanismCompletionGuarantees
	if ti.IndependentMonitoringList != nil {

		t.IndependentMonitoringList = *ti.IndependentMonitoringList
	}
	t.ESGSafeguardsSusCerts = ti.ESGSafeguardsSusCerts
	t.ESGSafeguardsCommEngPlans = ti.ESGSafeguardsCommEngPlans
	t.SecurityMeasuresAccessControl = ti.SecurityMeasuresAccessControl
	t.SecurityMeasuresSurveilanceSystems = ti.SecurityMeasuresSurveilanceSystems
	t.SecurityMeasuresOnSiteSecurityPersonnel = ti.SecurityMeasuresOnSiteSecurityPersonnel
	t.SecurityMeasuresPerimeterSecurity = ti.SecurityMeasuresPerimeterSecurity
	t.SecurityMeasuresCriticalInfraProtections = ti.SecurityMeasuresCriticalInfraProtections
	if ti.OtherAssetProtection != nil {

		t.OtherAssetProtection = *ti.OtherAssetProtection
	}
	if ti.LegalAdvisor != nil {

		t.LegalAdvisor = strings.ToUpper(*ti.LegalAdvisor)
	}
	if ti.FinancialAdvisor != nil {

		t.FinancialAdvisor = strings.ToUpper(*ti.FinancialAdvisor)
	}

	t.UndertakingNoLien = ti.UndertakingNoLien
	t.UndertakingNotCollateral = ti.UndertakingNotCollateral
	t.UndertakingNoClaims = ti.UndertakingNoClaims
	t.UndertakingNoForeclosure = ti.UndertakingNoForeclosure
	t.ComplianceNoViolation = ti.ComplianceNoViolation
	t.ComplianceAllPermits = ti.ComplianceAllPermits
	t.OutstandingFinancialRespNoDebts = ti.OutstandingFinancialRespNoDebts
	t.OutstandingFinancialRespNoHiddenLiabilities = ti.OutstandingFinancialRespNoHiddenLiabilities
	t.RiskManagementFullyInsured = ti.RiskManagementFullyInsured
	t.RiskManagementDeclaredValue = ti.RiskManagementDeclaredValue
	t.PhysicalConditionSound = ti.PhysicalConditionSound
	t.PhysicalConditionNolease = ti.PhysicalConditionNolease
	t.PhysicalConditionNoUndisclosedEasements = ti.PhysicalConditionNoUndisclosedEasements
	if ti.DeepLink != nil {
		t.DeepLink = *ti.DeepLink
	}
	return t

}

func (gc *GlobalConfig) LogDiscordFailedRequest(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

// GetOrderBook
func (gc *GlobalConfig) GetOrderBook(assetCode, assetIssuer, currencyCode, currencyIssuer string) (trovoOrderBook OrderBook, err error) {
	trovoOrderBook.Asks = make([]OfferVolume, 0)
	trovoOrderBook.Bids = make([]OfferVolume, 0)
	trovoOrderBook.Currency = Asset{
		AssetCode:   currencyCode,
		AssetIssuer: currencyIssuer,
	}
	trovoOrderBook.Asset = Asset{
		AssetCode:   assetCode,
		AssetIssuer: assetIssuer,
	}
	if strings.EqualFold(assetCode, os.Getenv("NATIVE_ASSET_CODE")) {
		assetCode = ""
		assetIssuer = ""
	}
	if strings.EqualFold(currencyCode, os.Getenv("NATIVE_ASSET_CODE")) {
		currencyCode = ""
		currencyIssuer = ""
	}

	var input OrderBookRequestInput

	input.SellingAssetCode = assetCode
	input.SellingAssetIssuer = assetIssuer
	input.BuyingAssetCode = currencyCode
	input.BuyingAssetIssuer = currencyIssuer

	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetOrderBook]Error getting order book summary: %v\n", err)
		return
	}

	//process bids

	for _, bid := range orderBook.Bids {
		trovoOrderBook.Bids = append(trovoOrderBook.Bids, OfferVolume{
			Price:    bid.Price,
			Quantity: bid.Amount,
		})
	}

	//process asks

	for _, ask := range orderBook.Asks {
		trovoOrderBook.Asks = append(trovoOrderBook.Asks, OfferVolume{
			Price:    ask.Price,
			Quantity: ask.Amount,
		})
	}

	return trovoOrderBook, nil

}

func (gc *GlobalConfig) GetVATValue(serviceFee decimal.Decimal) (vat float64) {
	// var serviceFee ServiceFee
	cc := gc.GetConfig("NG")
	vat = serviceFee.Mul(decimal.NewFromFloat((cc.VATPercent / 100))).Truncate(7).InexactFloat64()
	return
}

func (gc *GlobalConfig) GetVATRate() (vat float64) {
	// var serviceFee ServiceFee
	cc := gc.GetConfig("NG")

	return cc.VATPercent
}

type ServiceFee struct {
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	LastUpdatedBy      string    `json:"lastUpdatedBy"`
	ID                 string    `json:"id"`
	FeeWalletSecretKey string    `json:"feeWalletSecretKey"` //Secret key
	FeePercent         float64   `gorm:"default:0" json:"feePercent"`
	FeeFixed           float64   `gorm:"default:0" json:"feeFixed"`
	FeeAssetCode       string    `gorm:"default:''" json:"feeAssetCode"`
	FeeAssetIssuer     string    `gorm:"default:''" json:"feeAssetIssuer"`
	Inactive           int       `gorm:"default:0" json:"inactive"`
	Remarks            string    `gorm:"default:''" json:"remarks"`
}

func (gc *GlobalConfig) GetVATWallet() string {
	var serviceFee ServiceFee
	gc.DB.Where("id = ? AND inactive = 0", "VAT").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero swap fees and modify the swap fee

	}
	// set the VAT FEE WALLET IN ENV
	serviceFee.FeeWalletSecretKey = os.Getenv("VAT_WALLET")
	return serviceFee.FeeWalletSecretKey
}

// ServiceLinkServiceFee holds fee data model
type ServiceLinkServiceFee struct {
	ServiceLinkID string `gorm:"size:100;primaryKey" json:"serviceLinkId"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PaymentFee    float64 `json:"-" gorm:"not null;default:0.0"`
	SwapFee       float64 `json:"-" gorm:"not null;default:0.0"`
	SubwalletFee  float64 `json:"-" gorm:"not null;default:0.0"`
}

type FeeCollection struct {
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	ID                         string    `gorm:"size:100;primaryKey" json:"id"`
	FromUsername               string    `json:"fromUsername"`
	FromWalletPublicKey        string    `json:"fromWalletPublicKey"`
	FromWalletAlias            string    `json:"fromWalletAlias"`
	BelongsToEnterpriseProfile *string   `gorm:"null" json:"belongsToEnterpriseProfile"` //enterprise profile username if this user belongs to an enterprise profile
	FeeType                    string    `gorm:"not null" json:"feeType"`
	Amount                     float64   `gorm:"default:0.0" json:"amount"`
	AssetCode                  string    `gorm:"not null" json:"assetCode"`
	AssetIssuer                *string   `gorm:"null" json:"assetIssuer"`
	DestinationWallet          string    `gorm:"not null" json:"destinationWallet"`
	SharedAccessOperation      int       `json:"sharedAccessOperation" gorm:"type:integer;not null;default:0"`
	TransactionHash            *string   `gorm:"null" json:"transactionHash"`
	Processed                  int       `json:"processed" gorm:"type:integer;not null;default:0"` // if the fee split has been processed or not
}

func (gc *GlobalConfig) GetServiceLinkFees(serviceLinkID string) (s ServiceLinkServiceFee, exists bool, err error) {
	e := gc.DB.Where("service_link_id = ?", serviceLinkID).First(&s).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			exists = false
			err = nil
			return s, false, nil
		}

		err = errors.New("error retrieving the service fees")
		log.Printf("[GC:GetServiceLinkFees] Error retreiving service fees: %v", e)
		return
	}
	// fees retreived
	return s, true, nil
}

func (gc *GlobalConfig) GenerateUUIDString() (s string) {
	s = uuid.NewString()
	return s
}

// Rates model
type CurrencyRates struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Rates     []byte    `gorm:"null" json:"rates"`
}
type CngnPriceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		UsdToNgn       float64   `json:"usdToNgn"`
		NgnToUsd       float64   `json:"ngnToUsd"`
		FormattedPrice string    `json:"formattedPrice"`
		ReversePrice   string    `json:"reversePrice"`
		Decimals       int       `json:"decimals"`
		Description    string    `json:"description"`
		Timestamp      time.Time `json:"timestamp"`
	} `json:"data"`
}
type fiatConversionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		NgnAmount       int       `json:"ngnAmount"`
		UsdAmount       float64   `json:"usdAmount"`
		Rate            float64   `json:"rate"`
		FormattedResult string    `json:"formattedResult"`
		Timestamp       time.Time `json:"timestamp"`
	} `json:"data"`
}

// GetRates is service to get rate
func (gc *GlobalConfig) GetCngnUsdRate() (rate CngnPriceResponse) {
	cacheKeyGetUSDRate := "cngnusdRate"

	{

		// search cache
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyGetUSDRate)

		if ok {

			json.Unmarshal(rawdata, &rate)
			return
		}

	}
	rates := &CurrencyRates{}
	gc.DB.Where("id = ?", 2).First(&rates)
	json.Unmarshal(rates.Rates, &rate)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyGetUSDRate, rate, 2000)
	return
}

// ConvertCngnToUsd is service to get USD from cngn
func (gc *GlobalConfig) ConvertCngnToUsd(cngnAmount float64) (result decimal.Decimal) {
	result = decimal.Zero
	if len(os.Getenv("CNGN_PRICE_API_URL")) == 0 {
		fmt.Println("[ConvertCngnToUsd] CNGN_PRICE_API_URL ENV not Defined")
		return
	}

	resp, err := http.Get(fmt.Sprintf("%v/api/convert/ngn-to-usd/%v", os.Getenv("CNGN_PRICE_API_URL"), decimal.NewFromFloat(cngnAmount).String()))
	if err != nil {
		log.Println("[ConvertCngnToUsd] http get error", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ConvertCngnToUsd] read body error:", err)
		return
	}

	fetchedConvertion := fiatConversionResponse{}
	err = json.Unmarshal(body, &fetchedConvertion)
	if err != nil {
		log.Println("Error Fetching ConvertCngnToUsd\n\n", err)
		return
	}

	// var a interface{}
	if fetchedConvertion.Success {
		result = decimal.NewFromFloat(fetchedConvertion.Data.UsdAmount)

	}
	return
}

// ConvertUsdToCngn is service to get cngn from USD
func (gc *GlobalConfig) ConvertUsdToCngn(usdAmount float64) (result decimal.Decimal) {
	result = decimal.Zero
	if len(os.Getenv("CNGN_PRICE_API_URL")) == 0 {
		fmt.Println("[ConvertUsdToCngn] CNGN_PRICE_API_URL ENV not Defined")
		return
	}

	resp, err := http.Get(fmt.Sprintf("%v/api/convert/usd-to-ngn/%v", os.Getenv("CNGN_PRICE_API_URL"), decimal.NewFromFloat(usdAmount).String()))
	if err != nil {
		log.Println("[ConvertUsdToCngn] http get error", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ConvertUsdToCngn] read body error:", err)
		return
	}

	fetchedConvertion := fiatConversionResponse{}
	err = json.Unmarshal(body, &fetchedConvertion)
	if err != nil {
		log.Println("Error Fetching ConvertUsdToCngn\n\n", err)
		return
	}

	// var a interface{}
	if fetchedConvertion.Success {
		result = decimal.NewFromFloat(fetchedConvertion.Data.UsdAmount)

	}
	return
}
