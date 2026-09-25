package usermetrics

import (
	userServices "admin-panel-dashboard/internal/components/users/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/middleware"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TradeLiabilities struct {
	SellingLiabilities string `json:"sellingLiabilities"`
	BuyingLiabilities  string `json:"buyingLiabilities"`
}

type CryptoWalletDepositAddress struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	TrovoWalletPublicKey string    `json:"TrovoWalletPublicKey"`
	Currency             string    `json:"currency"`
	DepositAddress       string    `json:"depositAddress"`
	Network              string    `json:"network"`
	QRCode               *string   `json:"qrCode"`
}

type Balance struct {
	AssetIssuer                  string                       `json:"assetIssuer"`
	AssetCode                    string                       `json:"assetCode"`
	Amount                       string                       `json:"amount"`
	InTrade                      TradeLiabilities             `json:"inTrade"`
	QRCode                       string                       `json:"qrCode"`
	ImageURL                     string                       `json:"imageUrl"`
	USDPrice                     string                       `json:"usdPrice"`
	NativePrice                  string                       `json:"nativePrice"`
	CryptoWalletDepositAddresses []CryptoWalletDepositAddress `json:"cryptoWalletDepositAddresses"`
	ClosedGroup                  string                       `json:"closedGroup"`
}

type AssetBalances struct {
	Claimed   []Balance `json:"claimed"`
	Unclaimed []Balance `json:"unclaimed"`
}

type ClosedGroup struct {
	ID                      string `json:"id"`
	GroupName               string `json:"groupName"`
	GroupOwner              string `json:"groupOwner"`
	GroupDescription        string `json:"groupDescription"`
	RegisteredEntity        int    `json:"registeredEntity"`
	RegistrationName        string `json:"registrationName"`
	RegistrationNumber      string `json:"registrationNumber"`
	RegistrationDocumentUrl string `json:"registrationDocumentUrl"`
}

type ApprovedAssetCustodian struct {
	ID                    uint64 `json:"id"`
	AssetCustodianName    string `json:"assetCustodianName"`
	AssetCustodianAddress string `json:"assetCustodianAddress"`
	AssetCustodianCountry string `json:"assetCustodianCountry"` // two character country code
	RequirementDocument   string `json:"requirementDocument"`
}

type AssetManager struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type AssetTokenizationDocument struct {
	ID               uint64    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	TokenizedAssetID string    `json:"tokenizedAssetId"`
	DocumentType     string    `json:"documentType"`
	DocumentTitle    string    `json:"documentTitle"`
	DocumentUrl      string    `json:"documentUrl"`
}

type TokenizationFee struct {
	ID                 uint64  `json:"id"`
	FeeFiatPercentage  float64 `json:"feeFiatPercentage"`
	FeeFiatCap         float64 `json:"feeFiatCap"`
	FeeAssetPercentage float64 `json:"feeAssetPercentage"`
	FeeDescription     string  `json:"feeDescription"`
}

type TokenizationFeeProofOfPayment struct {
	ID          uint64    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	DocumentUrl string    `json:"documentUrl"`
}

type TokenizedAssetSector struct {
	ID                  string `json:"sector"`
	RequirementDocument string `json:"requirementDocument"`
}

type TokenizationStatuses struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
}

type TokenizedAsset struct {
	ID                                          string                          `json:"id"`
	CreatedAt                                   time.Time                       `json:"createdAt"`
	UpdatedAt                                   time.Time                       `json:"updatedAt"`
	InitiatorUsername                           string                          `json:"initiatorUsername"`
	AssetSector                                 string                          `json:"assetSector"`
	AssetSubSector                              string                          `json:"assetSubSector"`
	AssetType                                   string                          `json:"assetType"`
	AssetName                                   string                          `json:"assetName"`
	AssetWebsite                                string                          `json:"assetWebsite"`
	ApprovedAssetCustodianID                    uint64                          `json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian                      ApprovedAssetCustodian          `json:"approvedAssetCustodianInfo"`
	OfferingType                                string                          `json:"offeringType"`
	ClosedGroupID                               string                          `json:"closedGroupId"`
	ClosedGroup                                 ClosedGroup                     `json:"closedGroupInfo"`
	SecApproval                                 int                             `json:"secApproval"`
	SecApprovalIdNumber                         string                          `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                      string                          `json:"issuingWalletPublicKey"`
	IssuingWalletAlias                          string                          `json:"issuingWalletAlias"`
	MarketMakingWallet                          string                          `json:"marketMakingWallet"`
	AssetDescription                            string                          `json:"assetDescription"`
	AssetCountryLocation                        string                          `json:"assetCountryLocation"`
	AssetPhysicalAddress                        string                          `json:"assetPhysicalAddress"`
	AssetLongitude                              string                          `json:"assetLongitude"`
	AssetLatitude                               string                          `json:"assetLatitude"`
	OwnershipType                               string                          `json:"ownershipType"`
	OwnershipKind                               string                          `json:"ownershipKind"`
	InitialOwnerPreferredWalletAddress          string                          `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              string                          `json:"assetOwnerName"`
	AssetOwnerRetainedOrContributedValue        float64                         `json:"assetOwnerRetainedOrContributedValue"`
	AssetOwnerAddress                           string                          `json:"assetOwnerAddress"`
	AssetManagerID                              uint64                          `json:"assetManagerId"`
	AssetManager                                AssetManager                    `json:"assetManagerInfo"`
	AssetQuoteCurrency                          string                          `json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64                         `json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation              float64                         `json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                       float64                         `json:"valueOfTokenizedAsset"`
	ProtectionMethods                           string                          `json:"protectionMethods"`
	InsuranceCompanyName                        string                          `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       string                          `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       string                          `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64                         `json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int                             `json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int                             `json:"assetAlreadyExists"`
	VettingStatus                               int                             `json:"vettingStatus"`
	AssetTokenizationDocuments                  []AssetTokenizationDocument     `json:"AssetTokenizationDocuments"`
	ProofOfPaymentDocuments                     []TokenizationFeeProofOfPayment `json:"ProofOfPaymentDocuments"`
	AssetCode                                   string                          `json:"assetCode"`
	AssetLogo                                   string                          `json:"assetLogo"`
	NumberOfTokenToBeIssued                     float64                         `json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale            float64                         `json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                  float64                         `json:"feeInAsset"`
	FeeInFiat                                   float64                         `json:"feeInFiat"`
	NumberOfTokenToBeSold                       float64                         `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64                         `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                string                          `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                               float64                         `json:"pricePerToken"`
	SalesStart                                  time.Time                       `json:"salesStart"`
	SalesEnd                                    time.Time                       `json:"salesEnd"`
	CapOnPurchase                               float64                         `json:"capOnPurchase"`
	CapQuantity                                 float64                         `json:"capQuantity"`
	CapDurationInDays                           int                             `json:"capDurationInDays"`
	ProceedCycle                                string                          `json:"proceedCycle"`
	TokenizationFeeID                           uint64                          `json:"tokenizationFeeId"`
	TokenizationFee                             TokenizationFee                 `json:"tokenizationFee"`
	SECTokenizationFeePercent                   float64                         `json:"SECTokenizationFeePercent"`
	SECTokenizationFeeValue                     float64                         `json:"SECTokenizationFeeValue"`
	CustodianFeePercent                         float64                         `json:"custodianFeePercent"`
	CustodianFeeValue                           float64                         `json:"custodianFeeValue"`
	AssetManagerFeeValue                        float64                         `json:"assetManagerFeeValue"`
	AssetManagerFeePercent                      float64                         `json:"assetManagerFeePercent"`
	ProceedPayoutCurrency                       string                          `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int                             `json:"proceedPayoutType"`
	ExemptedCountries                           string                          `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int                             `json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   string                          `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int                             `json:"investorAccreditationRequired"`
	AssetTokenizationStatus                     int                             `json:"assetTokenizationStatus"`
	LastUpdatedBy                               string                          `json:"lastUpdatedBy"`
	AgreeTransferTitleToCustodian               int                             `json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int                             `json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int                             `json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int                             `json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int                             `json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int                             `json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int                             `json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   string                          `json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int                             `json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int                             `json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int                             `json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int                             `json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int                             `json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int                             `json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int                             `json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        string                          `json:"otherAssetProtection"`
	LegalAdvisor                                string                          `json:"legalAdvisor"`
	FinancialAdvisor                            string                          `json:"financialAdvisor"`
	UndertakingNoLien                           int                             `json:"undertakingNoLien"`
	UndertakingNotCollateral                    int                             `json:"undertakingNotCollateral"`
	UndertakingNoClaims                         int                             `json:"undertakingNoClaims"`
	UndertakingNoForeclosure                    int                             `json:"undertakingNoForeclosure"`
	ComplianceNoViolation                       int                             `json:"complianceNoViolation"`
	ComplianceAllPermits                        int                             `json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts             int                             `json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities int                             `json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                  int                             `json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                 int                             `json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                      int                             `json:"physicalConditionSound"`
	PhysicalConditionNolease                    int                             `json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements     int                             `json:"physicalConditionNoUndisclosedEasements"`
}

type TokenizationRequest struct {
	AssetSector                                 string    `json:"assetSector"`
	AssetSubSector                              string    `json:"assetSubSector"`
	AssetType                                   string    `json:"assetType"`
	AssetName                                   string    `json:"assetName"`
	AssetWebsite                                string    `json:"assetWebsite"`
	ApprovedAssetCustodianID                    uint64    `json:"approvedAssetCustodianId"`
	OfferingType                                string    `json:"offeringType"`
	ClosedGroupID                               string    `json:"closedGroupId"`
	SecApproval                                 int       `json:"secApproval"`
	SecApprovalIdNumber                         string    `json:"secApprovalIdNumber"`
	MarketMakingWallet                          string    `json:"marketMakingWallet"`
	AssetDescription                            string    `json:"assetDescription"`
	AssetCountryLocation                        string    `json:"assetCountryLocation"`
	AssetPhysicalAddress                        string    `json:"assetPhysicalAddress"`
	AssetLongitude                              string    `json:"assetLongitude"`
	AssetLatitude                               string    `json:"assetLatitude"`
	OwnershipType                               string    `json:"ownershipType"`
	OwnershipKind                               string    `json:"ownershipKind"`
	InitialOwnerPreferredWalletAddress          string    `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              string    `json:"assetOwnerName"`
	AssetOwnerRetainedOrContributedValue        float64   `json:"assetOwnerRetainedOrContributedValue"`
	AssetOwnerAddress                           string    `json:"assetOwnerAddress"`
	AssetManagerID                              uint64    `json:"assetManagerId"`
	AssetQuoteCurrency                          string    `json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64   `json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation              float64   `json:"assetMscCostOutisdeOfValuation"`
	ProtectionMethods                           string    `json:"protectionMethods"`
	InsuranceCompanyName                        string    `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       string    `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       string    `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64   `json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int       `json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int       `json:"assetAlreadyExists"`
	AssetCode                                   string    `json:"assetCode"`
	AssetLogo                                   string    `json:"assetLogo"`
	NumberOfTokenToBeIssued                     float64   `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold                       float64   `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64   `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                string    `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                               float64   `json:"pricePerToken"`
	SalesStart                                  time.Time `json:"salesStart"`
	SalesEnd                                    time.Time `json:"salesEnd"`
	CapOnPurchase                               float64   `json:"capOnPurchase"`
	CapQuantity                                 float64   `json:"capQuantity"`
	CapDurationInDays                           int       `json:"capDurationInDays"`
	ProceedCycle                                string    `json:"proceedCycle"`
	TokenizationFeeID                           uint64    `json:"tokenizationFeeId"`
	ProceedPayoutCurrency                       string    `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int       `json:"proceedPayoutType"`
	ExemptedCountries                           string    `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int       `json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   string    `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int       `json:"investorAccreditationRequired"`
	Messages                                    []string  `json:"messages"`
	AgreeTransferTitleToCustodian               int       `json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int       `json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int       `json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int       `json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int       `json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int       `json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int       `json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   string    `json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int       `json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int       `json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int       `json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int       `json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int       `json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int       `json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int       `json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        string    `json:"otherAssetProtection"`
	LegalAdvisor                                string    `json:"legalAdvisor"`
	FinancialAdvisor                            string    `json:"financialAdvisor"`
	UndertakingNoLien                           int       `json:"undertakingNoLien"`
	UndertakingNotCollateral                    int       `json:"undertakingNotCollateral"`
	UndertakingNoClaims                         int       `json:"undertakingNoClaims"`
	UndertakingNoForeclosure                    int       `json:"undertakingNoForeclosure"`
	ComplianceNoViolation                       int       `json:"complianceNoViolation"`
	ComplianceAllPermits                        int       `json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts             int       `json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities int       `json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                  int       `json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                 int       `json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                      int       `json:"physicalConditionSound"`
	PhysicalConditionNolease                    int       `json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements     int       `json:"physicalConditionNoUndisclosedEasements"`
	MintingInitators                            string    `json:"mintingInitators"`
	MintingApprovers                            string    `json:"mintingApprovers"`
	BankID                                      uint64    `json:"bankId"`
	AccountNumber                               string    `json:"accountNumber"`
	BeneficiaryName                             string    `json:"beneficiaryName"`
}

type TokenizationVetRequest struct {
	ApprovedAssetCustodianID uint64   `json:"approvedAssetCustodianId"`
	AssetManagerID           uint64   `json:"assetManagerId"`
	CountryCode              string   `json:"CountryCode"`
	ProceedPayoutCurrency    string   `json:"proceedPayoutCurrency"`
	AssetQuoteCurrency       string   `json:"assetQuoteCurrency"`
	Messages                 []string `json:"messages"`
}

type TokenizationListRequest struct {
	AssetTokenizationStatus *int   `form:"assetTokenizationStatus"`
	OnlyWithUserPermission  *int   `form:"onlyWithUserPermission"`
	AssetDescription        string `form:"assetDescription"`
	SalesList               *int   `form:"salesList"`
	AssetType               string `form:"assetType"`
	AssetName               string `form:"assetName"`
	AssetCode               string `form:"assetCode"`
	AssetSector             string `form:"assetSector"`
	AssetSubSector          string `form:"assetSubSector"`
	InitiatorUsername       string `form:"initiatorUsername"`
	OfferingType            string `form:"offeringType"`
	HasSecApproval          *int   `form:"hasSecApproval"`
	Limit                   int    `form:"limit,default=25"`
	Page                    int    `form:"page,default=1"`
	CreatedBetween          string `form:"createdBetween"`
}

type TokenizationListResponse struct {
	Pages        int              `json:"pages"`
	CurrentPage  int              `json:"currentPage"`
	TotalRecords int              `json:"totalRecords"`
	Limit        int              `json:"limit"`
	Records      []TokenizedAsset `json:"records"`
}

type UpdateSalesDatesRequest struct {
	SalesStart string `json:"salesStart"`
	SalesEnd   string `json:"salesEnd"`
}

type DueDiligenceFailRequest struct {
	Reason string `json:"reason" binding:"required"`
}

const (
	MaxFileSize = 1024 * 1024 // 900KB in bytes
)

var allowedFileTypes = map[string]bool{
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
	"image/gif":       true,
	"application/pdf": true,
}

var allowedLogoTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
}

type DocumentUploadRequest struct {
	TokenizedAssetID string `form:"tokenizedAssetId"`
	DocumentType     string `form:"documentType"`
	DocumentTitle    string `form:"documentTitle"`
}

func validateFile(file *multipart.FileHeader) error {
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", file.Size, MaxFileSize)
	}

	// Open the file to check its content type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer src.Close()

	// Read the first 512 bytes to determine the content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading file: %v", err)
	}

	contentType := http.DetectContentType(buffer)
	if !allowedFileTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s. Only jpg, jpeg, png, gif and pdf are supported", contentType)
	}

	return nil
}

func validateLogo(file *multipart.FileHeader) error {
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", file.Size, MaxFileSize)
	}

	// Open the file to check its content type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer src.Close()

	// Read the first 512 bytes to determine the content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading file: %v", err)
	}

	contentType := http.DetectContentType(buffer)
	if !allowedLogoTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s. Only jpg, jpeg, png, and gif are supported", contentType)
	}

	return nil
}

// makeRequest makes a request to the Trovo Wallet API
// method: HTTP method (GET, POST, etc.)
// endpoint: API endpoint (e.g., "/wallet-balances/{walletPublicKey}")
// token: JWT token for authorization
// body: request body for POST/PUT requests (can be nil for GET requests)
// result: pointer to the struct where the response will be unmarshaled
func makeRequest(method, endpoint, token string, body interface{}, result interface{}) error {
	trovoWalletBaseURL := os.Getenv("TROVO_WALLET_BASE_URL")
	if trovoWalletBaseURL == "" {
		trovoWalletBaseURL = "https://apidev.trovotechnologies.com" // fallback to dev URL
	}

	url := fmt.Sprintf("%s/v1/trovo-manager%s", trovoWalletBaseURL, endpoint)

	// Prepare request body if provided
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %v", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	log.Println("[METRICS] request to:", url)
	defer resp.Body.Close()
	log.Println("[METRICS] response status:", resp.Status)
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(respBody, &errorResponse); err != nil {
			return fmt.Errorf("api error: %s", string(respBody))
		}
		return fmt.Errorf("api error: %v", errorResponse)
	}

	// Parse the response into the provided result struct
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	return nil
}

// makeRequest makes a request to the Trovo Wallet API
// method: HTTP method (GET, POST, etc.)
// endpoint: API endpoint (e.g., "/wallet-balances/{walletPublicKey}")
// token: JWT token for authorization
// body: request body for POST/PUT requests (can be nil for GET requests)
// result: pointer to the struct where the response will be unmarshaled
func makeRequestWithRaw(method, endpoint, token string, body []byte) (result []byte, err error) {
	trovoWalletBaseURL := os.Getenv("TROVO_WALLET_BASE_URL")
	if trovoWalletBaseURL == "" {
		trovoWalletBaseURL = "https://apidev.trovotechnologies.com" // fallback to dev URL
	}

	url := fmt.Sprintf("%s/v1/trovo-manager%s", trovoWalletBaseURL, endpoint)

	// Prepare request body if provided
	// var bodyReader io.Reader

	bodyReader := bytes.NewReader(body)
	// Create request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return result, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error making request: %v", err)
	}
	log.Println("[METRICS] request to:", url)
	defer resp.Body.Close()
	log.Println("[METRICS] response status:", resp.Status)
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response: %v", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(respBody, &errorResponse); err != nil {
			return result, fmt.Errorf("api error: %s", string(respBody))
		}
		return result, fmt.Errorf("api error: %v", errorResponse)
	}

	// // Parse the response into the provided result struct
	// if err := json.Unmarshal(respBody, result); err != nil {
	// 	return fmt.Errorf("error parsing response: %v", err)
	// }

	return respBody, nil
}

// @Summary Get wallet balances
// @Description Get the wallet balances of the walletPublicKey submitted in the URI request
// @ID GetWalletBalances
// @Tags Wallets
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param walletPublicKey path string true "Wallet Public Key"
// @Success 200 {object} AssetBalances
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /wallet-balances/{walletPublicKey} [get]
func GetWalletBalances(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}

			c.JSON(statusCode, response)
			return
		}

		walletPublicKey := c.Param("walletPublicKey")
		if walletPublicKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wallet public key is required"})
			return
		}

		// Use the makeRequestWithRaw function to fetch wallet balances
		endpoint := fmt.Sprintf("/wallet-balances/%s", walletPublicKey)
		result, err := makeRequestWithRaw(http.MethodGet, endpoint, c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[METRICS] error fetching wallet balances: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Get tokenization parameters (Public)
// @Description Get all parameters needed for tokenization and display (requires authentication)
// @ID GetPublicTokenizationWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /public/tokenization [get]
func GetPublicTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		result, err := makeRequestWithRaw(http.MethodGet, "/public/tokenization", c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION] error fetching public tokenization data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Get tokenization parameters (Protected)
// @Description Get all parameters needed for tokenization and display (requires authentication)
// @ID GetTokenizationWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization [get]
func GetTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		result, err := makeRequestWithRaw(http.MethodGet, "/tokenization", c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION] error fetching tokenization data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Get tokenization parameters (Trovo Manager)
// @Description Get all parameters needed for tokenization and display (requires Trovo Manager access)
// @ID GetTrovoManagerTokenizationWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization [get]
func GetTrovoManagerTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		result, err := makeRequestWithRaw(http.MethodGet, "/tokenization", c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION] error fetching trovo manager tokenization data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Get tokenization details (Trovo Manager)
// @Description Get detailed information about a specific tokenization request by ID
// @ID GetTrovoManagerTokenizationDetail
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/detail/{tokenizedAssetID} [get]
func GetTrovoManagerTokenizationDetail(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var result TokenizedAsset
		endpoint := fmt.Sprintf("/tokenization/detail/%s", tokenizedAssetID)
		err = makeRequest(http.MethodGet, endpoint, c.GetHeader("Authorization"), nil, &result)
		if err != nil {
			log.Printf("[TOKENIZATION] error fetching tokenization detail: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Tokenization detail retrieved successfully", result, nil)
	}
}

// @Summary Get tokenization details (Trovo Manager)
// @Description Get detailed information about a specific tokenization request by ID
// @ID GetTrovoManagerTokenizationDetailWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/detail/{tokenizedAssetID} [get]
func GetTrovoManagerTokenizationDetailWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		endpoint := fmt.Sprintf("/tokenization/detail/%s", tokenizedAssetID)
		result, err := makeRequestWithRaw(http.MethodGet, endpoint, c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION] error fetching tokenization detail: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Submit tokenization request
// @Description Submit parameters needed for tokenization
// @ID SubmitTokenization
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param request body TokenizationRequest true "Tokenization request parameters"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization [post]
func SubmitTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	// not needed
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		// if err := c.ShouldBindJSON(&request); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		// 	return
		// }

		var result []byte
		result, err = makeRequestWithRaw(http.MethodPost, "/tokenization", c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error submitting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// serverResponse.JSON(c, http.StatusOK, "Tokenization request submitted successfully", result, nil)
		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Submit tokenization request
// @Description Submit parameters needed for tokenization
// @ID SubmitTokenization
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param request body TokenizationRequest true "Tokenization request parameters"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization [post]
func SubmitTokenization(walletDB *gorm.DB) gin.HandlerFunc { // not needed
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		var request TokenizationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result TokenizedAsset
		err = makeRequest(http.MethodPost, "/tokenization", c.GetHeader("Authorization"), request, &result)
		if err != nil {
			log.Printf("[TOKENIZATION] error submitting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Tokenization request submitted successfully", result, nil)
	}
}

// @Summary Update tokenization request (Trovo Manager)
// @Description Update tokenization information for paid requests
// @ID UpdateTrovoManagerTokenizationRequest
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param request body TokenizationRequest true "Updated tokenization parameters"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/{tokenizedAssetID} [put]
func UpdateTrovoManagerTokenization(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var request TokenizationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result TokenizedAsset
		endpoint := fmt.Sprintf("/tokenization/%s", tokenizedAssetID)
		err = makeRequest(http.MethodPut, endpoint, c.GetHeader("Authorization"), request, &result)
		if err != nil {
			log.Printf("[TOKENIZATION] error updating tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Tokenization request updated successfully", result, nil)
	}
}

// @Summary Update tokenization request (Trovo Manager)
// @Description Update tokenization information for paid requests
// @ID UpdateTrovoManagerTokenization
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param request body TokenizationRequest true "Updated tokenization parameters"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/update/{tokenizedAssetID} [put]
func UpdateTrovoManagerTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result []byte
		endpoint := fmt.Sprintf("/tokenization/update/%s", tokenizedAssetID)
		result, err = makeRequestWithRaw(http.MethodPut, endpoint, c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error updating tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Vet tokenization request (Trovo Manager)
// @Description Mark a tokenization request as vetted with additional parameters
// @ID VetTrovoManagerTokenization
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param request body TokenizationVetRequest true "Vetting parameters"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/vet/{tokenizedAssetID} [put]
func VetTrovoManagerTokenization(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var request TokenizationVetRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Validate country code format (2 characters)
		if len(request.CountryCode) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CountryCode must be a 2-character code"})
			return
		}

		var result TokenizedAsset
		endpoint := fmt.Sprintf("/tokenization/vet/%s", tokenizedAssetID)
		err = makeRequest(http.MethodPut, endpoint, c.GetHeader("Authorization"), request, &result)
		if err != nil {
			log.Printf("[TOKENIZATION] error vetting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Tokenization request vetted successfully", result, nil)
	}
}

// @Summary Vet tokenization request (Trovo Manager)
// @Description Mark a tokenization request as vetted with additional parameters
// @ID VetTrovoManagerTokenization
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param request body TokenizationVetRequest true "Vetting parameters"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/vet/{tokenizedAssetID} [put]
func VetTrovoManagerTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result []byte
		endpoint := fmt.Sprintf("/tokenization/vet/%s", tokenizedAssetID)
		result, err = makeRequestWithRaw(http.MethodPut, endpoint, c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error vetting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Upload tokenization document
// @Description Upload and update document needed for tokenization
// @ID UploadTokenizationDocumentWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetId formData string true "Tokenized Asset ID"
// @Param documentType formData string true "Document Type"
// @Param documentTitle formData string true "Document Title"
// @Param documentFile formData file true "Document File (jpg, jpeg, png, gif, pdf, max 900KB)"
// @Success 200 {object} AssetTokenizationDocument
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/document [put]
func UploadTokenizationDocumentWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		// // Parse multipart form
		// if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing form: %v", err)})
		// 	return
		// }

		// // Get form values
		// var request DocumentUploadRequest
		// if err := c.ShouldBind(&request); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		// 	return
		// }

		// // Get the file
		// file, err := c.FormFile("documentFile")
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "documentFile is required"})
		// 	return
		// }

		// // Validate file
		// if err := validateFile(file); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// // Create multipart form data
		// body := &bytes.Buffer{}
		// writer := multipart.NewWriter(body)

		// // Add form fields
		// writer.WriteField("tokenizedAssetId", request.TokenizedAssetID)
		// writer.WriteField("documentType", request.DocumentType)
		// writer.WriteField("documentTitle", request.DocumentTitle)

		// // Add file
		// part, err := writer.CreateFormFile("documentFile", file.Filename)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating form file"})
		// 	return
		// }

		// // Open uploaded file
		// src, err := file.Open()
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error opening uploaded file"})
		// 	return
		// }
		// defer src.Close()

		// // Copy file content
		// if _, err = io.Copy(part, src); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error copying file content"})
		// 	return
		// }

		// writer.Close()
		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result []byte
		// Make request to API
		result, err = makeRequestWithRaw(http.MethodPut, "/tokenization/document", c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error uploading document: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Upload tokenization document (Trovo Manager)
// @Description Upload and update document needed for tokenization (Trovo Manager access)
// @ID UploadTrovoManagerTokenizationDocumentWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetId formData string true "Tokenized Asset ID"
// @Param documentType formData string true "Document Type"
// @Param documentTitle formData string true "Document Title"
// @Param documentFile formData file true "Document File (jpg, jpeg, png, gif, pdf, max 900KB)"
// @Success 200 {object} AssetTokenizationDocument
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/document [put]
func UploadTrovoManagerTokenizationDocumentWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing form: %v", err)})
			return
		}

		// Get form values
		var request DocumentUploadRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
			return
		}

		// Get the file
		file, err := c.FormFile("documentFile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "documentFile is required"})
			return
		}

		// Validate file
		if err := validateFile(file); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add form fields
		err = writer.WriteField("tokenizedAssetId", request.TokenizedAssetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing tokenizedAssetId"})
			return
		}
		err = writer.WriteField("documentType", request.DocumentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing documentType"})
			return
		}
		err = writer.WriteField("documentTitle", request.DocumentTitle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing documentTitle"})
			return
		}

		// Add file
		part, err := writer.CreateFormFile("documentFile", file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating form file"})
			return
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error opening uploaded file"})
			return
		}
		defer src.Close()

		// Copy file content
		if _, err = io.Copy(part, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error copying file content"})
			return
		}

		writer.Close()

		// Make request to API
		result, err := makeRequestWithRaw(http.MethodPut, "/tokenization/document", c.GetHeader("Authorization"), body.Bytes())
		if err != nil {
			log.Printf("[TOKENIZATION] error uploading document: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Upload tokenized asset logo (Trovo Manager)
// @Description Upload and update logo for tokenized asset
// @ID UploadTrovoManagerTokenizationLogoWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param documentFile formData file true "Logo File (jpg, jpeg, png, gif, max 900KB)"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/logo/{tokenizedAssetID} [put]
func UploadTrovoManagerTokenizationLogo(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing form: %v", err)})
			return
		}

		// Get the file
		file, err := c.FormFile("documentFile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "documentFile is required"})
			return
		}

		// Validate logo
		if err := validateLogo(file); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add file
		part, err := writer.CreateFormFile("documentFile", file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating form file"})
			return
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error opening uploaded file"})
			return
		}
		defer src.Close()

		// Copy file content
		if _, err = io.Copy(part, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error copying file content"})
			return
		}

		writer.Close()

		// Make request to API
		var result TokenizedAsset
		endpoint := fmt.Sprintf("/tokenization/logo/%s", tokenizedAssetID)
		err = makeRequest(http.MethodPut, endpoint, c.GetHeader("Authorization"), body, &result)
		if err != nil {
			log.Printf("[TOKENIZATION] error uploading logo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Logo uploaded successfully", result, nil)
	}
}

// @Summary Delete tokenization document (Trovo Manager)
// @Description Delete an uploaded document for tokenization
// @ID DeleteTrovoManagerTokenizationDocumentWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param documentID path string true "Document ID"
// @Success 200 {object} string "Document deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/document/{documentID} [delete]
func DeleteTrovoManagerTokenizationDocumentWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		documentID := c.Param("documentID")
		if documentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "documentID is required"})
			return
		}

		// Make request to API
		endpoint := fmt.Sprintf("/tokenization/document/%s", documentID)
		result, err := makeRequestWithRaw(http.MethodDelete, endpoint, c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION] error deleting document: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// GetTrovoManagerTokenizationListWithRaw @Summary Get tokenization list (Trovo Manager)
// @Description Get all tokenizations with search and filter options
// @ID GetTrovoManagerTokenizationListWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param assetTokenizationStatus query int false "Filter by tokenization status"
// @Param onlyWithUserPermission query int false "Show only user's permissions (1) or market ready (0)"
// @Param assetDescription query string false "Filter by asset description"
// @Param salesList query int false "Show primary (0) or secondary (1) sales"
// @Param assetType query string false "Filter by asset type"
// @Param assetName query string false "Filter by asset name"
// @Param assetCode query string false "Filter by asset code"
// @Param assetSector query string false "Filter by asset sector"
// @Param assetSubSector query string false "Filter by asset subsector"
// @Param initiatorUsername query string false "Filter by initiator username"
// @Param offeringType query string false "Filter by offering type (PRIVATE/PUBLIC)"
// @Param hasSecApproval query int false "Filter by SEC approval status (0/1)"
// @Param limit query int false "Page size limit" default(25)
// @Param page query int false "Page number" default(1)
// @Param createdBetween query string false "Date range filter (format: 2020-01-01|2020-02-31)"
// @Success 200 {object} TokenizationListResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/list [get]
func GetTrovoManagerTokenizationList(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		// Parse query parameters
		var request TokenizationListRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
			return
		}

		// Set default values
		if request.Limit == 0 {
			request.Limit = 25
		}
		if request.Page == 0 {
			request.Page = 1
		}
		if request.OnlyWithUserPermission == nil {
			defaultValue := 1
			request.OnlyWithUserPermission = &defaultValue
		}

		// Validate date format if provided
		if request.CreatedBetween != "" {
			dates := strings.Split(request.CreatedBetween, "|")
			if len(dates) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range format. Expected format: 2020-01-01|2020-02-31"})
				return
			}
			for _, date := range dates {
				if _, err := time.Parse("2006-01-02", date); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. Expected format: YYYY-MM-DD"})
					return
				}
			}
		}

		// Build query parameters
		queryParams := url.Values{}
		if request.AssetTokenizationStatus != nil {
			queryParams.Set("assetTokenizationStatus", fmt.Sprintf("%d", *request.AssetTokenizationStatus))
		}
		if request.OnlyWithUserPermission != nil {
			queryParams.Set("onlyWithUserPermission", fmt.Sprintf("%d", *request.OnlyWithUserPermission))
		}
		if request.AssetDescription != "" {
			queryParams.Set("assetDescription", request.AssetDescription)
		}
		if request.SalesList != nil {
			queryParams.Set("salesList", fmt.Sprintf("%d", *request.SalesList))
		}
		if request.AssetType != "" {
			queryParams.Set("assetType", request.AssetType)
		}
		if request.AssetName != "" {
			queryParams.Set("assetName", request.AssetName)
		}
		if request.AssetCode != "" {
			queryParams.Set("assetCode", request.AssetCode)
		}
		if request.AssetSector != "" {
			queryParams.Set("assetSector", request.AssetSector)
		}
		if request.AssetSubSector != "" {
			queryParams.Set("assetSubSector", request.AssetSubSector)
		}
		if request.InitiatorUsername != "" {
			queryParams.Set("initiatorUsername", request.InitiatorUsername)
		}
		if request.OfferingType != "" {
			queryParams.Set("offeringType", request.OfferingType)
		}
		if request.HasSecApproval != nil {
			queryParams.Set("hasSecApproval", fmt.Sprintf("%d", *request.HasSecApproval))
		}
		queryParams.Set("limit", fmt.Sprintf("%d", request.Limit))
		queryParams.Set("page", fmt.Sprintf("%d", request.Page))
		if request.CreatedBetween != "" {
			queryParams.Set("createdBetween", request.CreatedBetween)
		}
		// queryParamsProxy := c.Request.URL.Query() // Returns url.Values (map[string][]string)
		// Make request to API
		endpoint := fmt.Sprintf("/tokenization/list?%s", queryParams.Encode())
		// endpoint := fmt.Sprintf("/tokenization/list?%s", queryParamsProxy.Encode())
		var result TokenizationListRequest
		err = makeRequest(http.MethodGet, endpoint, c.GetHeader("Authorization"), nil, &result)
		if err != nil {
			log.Printf("[TOKENIZATION] error getting tokenization list: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Tokenization list retrieved successfully", result, nil)
	}
}

// GetTrovoManagerTokenizationList @Summary Get tokenization list (Trovo Manager)
// @Description Get all tokenizations with search and filter options
// @ID GetTrovoManagerTokenizationList
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param assetTokenizationStatus query int false "Filter by tokenization status"
// @Param onlyWithUserPermission query int false "Show only user's permissions (1) or market ready (0)"
// @Param assetDescription query string false "Filter by asset description"
// @Param salesList query int false "Show primary (0) or secondary (1) sales"
// @Param assetType query string false "Filter by asset type"
// @Param assetName query string false "Filter by asset name"
// @Param assetCode query string false "Filter by asset code"
// @Param assetSector query string false "Filter by asset sector"
// @Param assetSubSector query string false "Filter by asset subsector"
// @Param initiatorUsername query string false "Filter by initiator username"
// @Param offeringType query string false "Filter by offering type (PRIVATE/PUBLIC)"
// @Param hasSecApproval query int false "Filter by SEC approval status (0/1)"
// @Param limit query int false "Page size limit" default(25)
// @Param page query int false "Page number" default(1)
// @Param createdBetween query string false "Date range filter (format: 2020-01-01|2020-02-31)"
// @Success 200 {object} TokenizationListResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/list [get]
func GetTrovoManagerTokenizationListWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		// Parse query parameters
		var request TokenizationListRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
			return
		}

		// Set default values
		if request.Limit == 0 {
			request.Limit = 25
		}
		if request.Page == 0 {
			request.Page = 1
		}
		if request.OnlyWithUserPermission == nil {
			defaultValue := 1
			request.OnlyWithUserPermission = &defaultValue
		}

		// Validate date format if provided
		if request.CreatedBetween != "" {
			dates := strings.Split(request.CreatedBetween, "|")
			if len(dates) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date range format. Expected format: 2020-01-01|2020-02-31"})
				return
			}
			for _, date := range dates {
				if _, err := time.Parse("2006-01-02", date); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. Expected format: YYYY-MM-DD"})
					return
				}
			}
		}

		// Build query parameters
		queryParams := url.Values{}
		if request.AssetTokenizationStatus != nil {
			queryParams.Set("assetTokenizationStatus", fmt.Sprintf("%d", *request.AssetTokenizationStatus))
		}
		if request.OnlyWithUserPermission != nil {
			queryParams.Set("onlyWithUserPermission", fmt.Sprintf("%d", *request.OnlyWithUserPermission))
		}
		if request.AssetDescription != "" {
			queryParams.Set("assetDescription", request.AssetDescription)
		}
		if request.SalesList != nil {
			queryParams.Set("salesList", fmt.Sprintf("%d", *request.SalesList))
		}
		if request.AssetType != "" {
			queryParams.Set("assetType", request.AssetType)
		}
		if request.AssetName != "" {
			queryParams.Set("assetName", request.AssetName)
		}
		if request.AssetCode != "" {
			queryParams.Set("assetCode", request.AssetCode)
		}
		if request.AssetSector != "" {
			queryParams.Set("assetSector", request.AssetSector)
		}
		if request.AssetSubSector != "" {
			queryParams.Set("assetSubSector", request.AssetSubSector)
		}
		if request.InitiatorUsername != "" {
			queryParams.Set("initiatorUsername", request.InitiatorUsername)
		}
		if request.OfferingType != "" {
			queryParams.Set("offeringType", request.OfferingType)
		}
		if request.HasSecApproval != nil {
			queryParams.Set("hasSecApproval", fmt.Sprintf("%d", *request.HasSecApproval))
		}
		queryParams.Set("limit", fmt.Sprintf("%d", request.Limit))
		queryParams.Set("page", fmt.Sprintf("%d", request.Page))
		if request.CreatedBetween != "" {
			queryParams.Set("createdBetween", request.CreatedBetween)
		}
		queryParamsProxy := c.Request.URL.Query() // Returns url.Values (map[string][]string)
		// Make request to API
		// endpoint := fmt.Sprintf("/tokenization/list?%s", queryParams.Encode())
		endpoint := fmt.Sprintf("/tokenization/list?%s", queryParamsProxy.Encode())
		var result []byte
		result, err = makeRequestWithRaw(http.MethodGet, endpoint, c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION] error getting tokenization list: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}

// @Summary Acknowledge tokenization fee payment (Trovo Manager)
// @Description Acknowledges the fee payment and advances the request to due diligence state
// @ID AcknowledgeTokenizationFeePaymentWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/fee/{tokenizedAssetID} [post]
func AcknowledgeTokenizationFeePaymentWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		endpoint := fmt.Sprintf("/tokenization/fee/%s", tokenizedAssetID)
		var result []byte
		result, err = makeRequestWithRaw(http.MethodPost, endpoint, c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error submitting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Fee payment acknowledged successfully", result, nil)
	}
}

// @Summary Fail due diligence for tokenization request (Trovo Manager)
// @Description Sets the DueDiligenceFail flag to 1 and records the reason. Resets tokenization and vetting status to 0.
// @ID FailDueDiligenceWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param request body DueDiligenceFailRequest true "Due diligence fail reason"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/failed/{tokenizedAssetID} [post]
func FailDueDiligenceWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		endpoint := fmt.Sprintf("/tokenization/faildd/%s", tokenizedAssetID)

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result []byte
		result, err = makeRequestWithRaw(http.MethodPost, endpoint, c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error submitting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Due diligence failed successfully", result, nil)
	}
}

// @Summary Confirm tokenization and mint asset (Trovo Manager)
// @Description Confirms the tokenization and initiates minting process to make the asset market ready
// @ID ConfirmAndMintTokenizationWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/mint/{tokenizedAssetID} [post]
func ConfirmAndMintTokenizationWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		endpoint := fmt.Sprintf("/tokenization/mint/%s", tokenizedAssetID)

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		var result []byte
		result, err = makeRequestWithRaw(http.MethodPost, endpoint, c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error submitting tokenization request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Tokenization confirmed and minting initiated successfully", result, nil)
	}
}

// @Summary Update tokenization sales dates (Trovo Manager)
// @Description Updates the sales start and end dates for a tokenization
// @ID UpdateTrovoManagerTokenizationSalesDatesWithRaw
// @Tags Tokenization
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param tokenizedAssetID path string true "Tokenized Asset ID"
// @Param request body UpdateSalesDatesRequest true "Sales dates in YYYY-MM-DD format"
// @Success 200 {object} TokenizedAsset
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tokenization/salesdate/{tokenizedAssetID} [put]
func UpdateTrovoManagerTokenizationSalesDatesWithRaw(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token metadata
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION] error for user:", userInfo.Username, "error: ", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		tokenizedAssetID := c.Param("tokenizedAssetID")
		if tokenizedAssetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizedAssetID is required"})
			return
		}

		var request []byte
		request, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		endpoint := fmt.Sprintf("/tokenization/salesdate/%s", tokenizedAssetID)
		result, err := makeRequestWithRaw(http.MethodPut, endpoint, c.GetHeader("Authorization"), request)
		if err != nil {
			log.Printf("[TOKENIZATION] error updating sales dates: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}
