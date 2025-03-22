package sharedconfig

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"
	"trovo-wallet-api/internal/cache"

	cs "cloud.google.com/go/storage"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"gorm.io/gorm"
)

type GlobalConfig struct {
	DynamicLinkServiceURLChan chan string
	PushNotificationClient    *messaging.Client
	FirebaseStorageUploader   *ClientUploader
	PNSContext                context.Context
	RedisCache                *cache.RedisCache
	DB                        *gorm.DB
	RoachDB                   *gorm.DB
	BantuExpansionClient      *horizonclient.Client
	BantuNetworkPassphrase    string
	ChannelAccounts           chan *keypair.Full
	InUseChannelAccounts      map[string]*keypair.Full
	Mutex                     sync.Mutex
	ChannelOfTokenizedAssetIDs  chan string
}

type ClientUploader struct {
	Client     *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
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
	ID                                          string    `json:"id"`
	CreatedAt                                   time.Time `json:"createdAt"`
	UpdatedAt                                   time.Time `json:"updatedAt"`
	InitiatorUsername                           string    `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                                 *string   `json:"assetSector"`
	AssetSubSector                              *string   `json:"assetSubSector"`
	AssetType                                   *string   `json:"assetType"`
	AssetName                                   *string   `json:"assetName"`
	ApprovedAssetCustodianID                    uint64    `gorm:"not null" json:"approvedAssetCustodianId"`
	OfferingType                                *string   `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                               *string   `gorm:"null" json:"closedGroupId"`
	SecApproval                                 int       `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                         *string   `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                      *string   `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                          *string   `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                          *string   `json:"marketMakingWallet"`
	AssetDescription                            *string   `json:"assetDescription"`
	AssetCountryLocation                        *string   `json:"assetCountryLocation"`
	AssetPhysicalAddress                        *string   `json:"assetPhysicalAddress"`
	AssetLongitude                              *string   `json:"assetLongitude"`
	AssetLatitude                               *string   `json:"assetLatitude"`
	OwnershipType                               *string   `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                               *string   `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress          *string   `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              *string   `json:"assetOwnerName"`
	AssetOwnerAddress                           *string   `json:"assetOwnerAddress"`
	AssetManagerID                              uint64    `json:"assetManagerId"`
	AssetQuoteCurrency                          *string   `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64   `gorm:"default:0" json:"assetCurrentValue"`
	AssetOwnerRetainedOrContributedValue        float64   `gorm:"default:0" json:"assetOwnerRetainedOrContributedValue"`
	AssetMscCostOutisdeOfValuation              float64   `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                       float64   `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                           *string   `json:"protectionMethods"` //csv format
	InsuranceCompanyName                        *string   `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       *string   `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       *string   `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64   `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int       `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int       `gorm:"default:1" json:"assetAlreadyExists"`
	VettingStatus                               int       `gorm:"default:0" json:"vettingStatus"`
	AssetCode                                   *string   `gorm:"size:12; index:idx_unique_tokenized_asset_code,unique" json:"assetCode"`
	AssetLogo                                   *string   `json:"assetLogo"`
	AssetWebsite                                *string   `gorm:"default:'trovotech.io'" json:"assetWebsite"`
	NumberOfTokenToBeIssued                     float64   `gorm:"default:0" json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale            float64   `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                  float64   `gorm:"default:0" json:"feeInAsset"`
	FeeInFiat                                   float64   `gorm:"default:0" json:"feeInFiat"`
	NumberOfTokenToBeSold                       float64   `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64   `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                *string   `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                               float64   `gorm:"default:0" json:"pricePerToken"`
	SalesStart                                  time.Time `json:"salesStart"`
	SalesEnd                                    time.Time `json:"salesEnd"`
	CapOnPurchase                               float64   `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                 float64   `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays                           int       `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                *string   `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                           *uint64   `gorm:"default:0" json:"tokenizationFeeId"`
	SECTokenizationFeePercent                   float64   `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeValue                     float64   `gorm:"default:0" json:"SECTokenizationFeeValue"`
	CustodianFeePercent                         float64   `gorm:"default:0" json:"custodianFeePercent"`
	CustodianFeeValue                           float64   `gorm:"default:0" json:"custodianFeeValue"`
	AssetManagerFeeValue                        float64   `gorm:"default:0" json:"assetManagerFeeValue"`
	AssetManagerFeePercent                      float64   `gorm:"default:0" json:"assetManagerFeePercent"`
	ProceedPayoutCurrency                       *string   `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int       `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                           *string   `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int       `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   *string   `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int       `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                     int       `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                               *string   `gorm:"null" json:"lastUpdatedBy"`
	TokenizationTransaction                     *string   `gorm:"null" json:"tokenizationTransaction"`
	AgreeTransferTitleToCustodian               int       `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int       `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int       `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int       `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int       `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int       `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int       `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   *string   `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int       `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int       `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int       `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int       `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int       `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int       `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int       `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        *string   `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                *string   `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                            *string   `gorm:"null" json:"financialAdvisor"`
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

func (gc *GlobalConfig) GetTokenizedAssetByCode(assetCode string) (t TokenizedAsset) {

	gc.DB.Where("asset_code = upper(?)", assetCode).First(&t)

	return t
}

func (gc *GlobalConfig) GetTokenizedAssetByID(tokenizedAssetID string) (t TokenizedAsset) {

	gc.DB.Where("id = ?", tokenizedAssetID).First(&t)

	return t
}
