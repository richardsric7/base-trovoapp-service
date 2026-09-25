export interface TokenizationResponse {
  pages: number;
  currentPage: number;
  totalRecords: number;
  limit: number;
  records: TokenizationRecord[];
}
export interface TokenizationRecord {
  id: string;
  createdAt: string;
  updatedAt: string;
  initiatorUsername: string;
  assetSector: string;
  assetSubSector: string;
  assetType: string;
  assetName: string;
  assetWebsite: string;
  approvedAssetCustodianId: number;
  approvedAssetCustodianInfo: ApprovedAssetCustodianInfo;
  assetIssuingHouseId: number;
  assetIssuingHouseInfo: AssetIssuingHouseInfo;
  legalAndProfesionalPartnerId: number;
  legalAndProfesionalPartnerInfo: LegalAndProfesionalPartnerInfo;
  ratingAgencyId: number;
  ratingAgencyInfo: RatingAgencyInfo;

  offeringType: string;
  trusteeId: number;
  trusteeInfo: TrusteeInfo;
  closedGroupId: string;
  closedGroupInfo: ClosedGroupInfo;
  secApproval: number;
  secApprovalIdNumber: string;
  issuingWalletPublicKey: string;
  issuingWalletAlias: string;
  marketMakingWallet: string;
  assetDescription: string;
  assetCountryLocation: string;
  assetPhysicalAddress: string;
  assetLongitude: string;
  assetLatitude: string;
  ownershipType: string;
  ownershipKind: string;
  initialOwnerPreferredWalletAddress: string;
  assetOwnerName: string;
  assetOwnerRetainedOrContributedValue: number;
  assetOwnerAddress: string;
  assetManagerId: number;
  assetManagerInfo: AssetManagerInfo;
  assetQuoteCurrency: string;
  assetCurrentValue: number;
  assetMscCostOutisdeOfValuation: number;
  valueOfTokenizedAsset: number;
  protectionMethods: string;
  insuranceCompanyName: string;
  insurancePolicyNumber: string;
  insurancePolicyHolder: string;
  percentageValueOfInsurance: number;
  IsFreeFromLiensAndEncumbrances: number;
  assetAlreadyExists: number;
  vettingStatus: number;
  dueDiligenceFail: number;
  dueDiligenceFailureReason: string;
  AssetTokenizationDocuments: AssetTokenizationDocument[];
  ProofOfPaymentDocuments: ProofOfPaymentDocument[];
  assetCode: string;
  assetLogo: string;
  numberOfTokenToBeIssued: number;
  maxNumberOfTokenAvailableForSale: number;
  feeInAsset: number;
  feeInFiat: number;
  numberOfTokenToBeSold: number;
  totalTokenHeldByManager: number;
  walletToHoldAssetsNotForSale: string;
  pricePerToken: number;
  salesStart: string;
  salesEnd: string;
  capOnPurchase: number;
  capQuantity: number;
  capAmountInFiat: number;
  capDurationInDays: number;
  feeFiatCap: number;
  proceedCycle: string;
  tokenizationFeeId: number;
  tokenizationFee: TokenizationFee;
  countryConfig: CountryConfig;
  SECTokenizationFeePercent: number;
  SECTokenizationFeeValue: number;
  SECTokenizationFeeFixed: number;
  custodianFeePercent: number;
  custodianFeeValue: number;
  custodianFeeFixed: number;
  assetManagerFeeValue: number;
  assetManagerFeePercent: number;
  assetManagerFeeFixed: number;
  issuingHouseFeeValue: number;
  issuingHouseFeePercent: number;
  issuingHouseFeeFixed: number;
  legalAndProfessionalFeePercent: number;
  legalAndProfessionalFeeFixed: number;
  legalAndProfessionalFeeValue: number;
  ratingAgencyFeePercent: number;
  ratingAgencyFeeFixed: number;
  ratingAgencyFeeValue: number;
  trusteeFeePercent: number;
  trusteeFeeFixed: number;
  trusteeFeeValue: number;
  vatInAsset: number;
  vatPercent: number;
  vatValue: number;
  proceedPayoutCurrency: string;
  proceedPayoutType: number;
  exemptedCountries: string;
  hasAdditionalKYCRequirements: number;
  additionalKYCRequirements: string;
  investorAccreditationRequired: number;
  assetTokenizationStatus: number;
  lastUpdatedBy: string;
  agreeTransferTitleToCustodian: number;
  contractualProtectionRevGuarantees: number;
  contractualProtectionPerfBond: number;
  contractualProtectionSLA: number;
  riskSharingMechanismPPPs: number;
  riskSharingMechanismHedgeInstruments: number;
  riskSharingMechanismCompletionGuarantees: number;
  independentMonitoringList: string;
  eSGSafeguardsSusCerts: number;
  eSGSafeguardsCommEngPlans: number;
  securityMeasuresAccessControl: number;
  securityMeasuresSurveilanceSystems: number;
  securityMeasuresOnSiteSecurityPersonnel: number;
  securityMeasuresPerimeterSecurity: number;
  securityMeasuresCriticalInfraProtections: number;
  otherAssetProtection: string;
  legalAdvisor: string;
  financialAdvisor: string;
  ratingAgency: string;
  undertakingNoLien: number;
  undertakingNotCollateral: number;
  undertakingNoClaims: number;
  undertakingNoForeclosure: number;
  complianceNoViolation: number;
  complianceAllPermits: number;
  outstandingFinancialRespNoDebts: number;
  outstandingFinancialRespNoHiddenLiabilities: number;
  riskManagementFullyInsured: number;
  riskManagementDeclaredValue: number;
  physicalConditionSound: number;
  physicalConditionNolease: number;
  physicalConditionNoUndisclosedEasements: number;
  mintingInitators: string;
  mintingApprovers: string;
  bankId: number;
  bankInfo: BankInfo;
  accountNumber: string;
  beneficiaryName: string;
  tokenizationApplicationFee: number;
  tokenizationApplicationFeeAsset: string;
  purchaseCommitments: number;
  quantityOfTokensSold: number;
  quantityOfTokensSoldInFiat: number;
  numberOfSubscribers: number;
  numberOfExpressedInterests: number;
}

export interface ApprovedAssetCustodianInfo {
  id: number;
  assetCustodianName: string;
  assetCustodianAddress: string;
  assetCustodianCountry: string;
  requirementDocument: string;
  FeePercent: number;
}

export interface ClosedGroupInfo {
  id: string;
  groupName: string;
  groupOwner: string;
  groupDescription: string;
  registeredEntity: number;
  registrationName: string;
  registrationNumber: string;
  registrationDocumentUrl: string;
}

export interface AssetManagerInfo {
  id: number;
  assetManagerName: string;
  assetManagerAddress: string;
  assetManagerCountry: string;
  FeePercent: number;
}

export interface AssetTokenizationDocument {
  id: number;
  CreatedAt: string;
  tokenizedAssetId: string;
  documentType: string;
  documentTitle: string;
  documentUrl: string;
}

export interface ProofOfPaymentDocument {
  id: number;
  CreatedAt: string;
  tokenizationFeePaymentMethodID: string;
  tokenizedAssetId: string;
  transactionReference: string;
  documentUrl: string;
}

export interface AssetIssuingHouseInfo {
  id: number;
  assetIssuingHouseName: string;
  assetIssuingHouseAddress: string;
  assetIssuingHouseCountry: string;
  FeePercent: number;
}

export interface TokenizationFee {
  id: number;
  feeFiatPercentage: number;
  feeFiatCap: number;
  feeAssetPercentage: number;
  feeDescription: string;
  countryCode: string;
}

export interface TrusteeInfo {
  id: number;
  trusteeName: string;
  trusteeAddress: string;
  trusteeCountry: string;
  FeePercent: number;
  FeeFixed: number;
}

export interface ClosedGroupInfo {
  id: string;
  groupName: string;
  groupOwner: string;
  groupDescription: string;
  registeredEntity: number;
  registrationName: string;
  registrationNumber: string;
  registrationDocumentUrl: string;
}

export interface CountryConfig {
  countryCode: string;
  SECTokenizationFee: number;
  SECTokenizationFeeType: number;
  SECTradeFeeFixed: number;
  SECTokenizationFeePercent: number;
  SECTokenizationFeeFixed: number;
  SECTradeFeePercent: number;
  SECTradeFeeType: number;
  regulatorName: string;
  regionName: string;
  countryName: string;
  quoteCurrencyCode: string;
  fiatLabel: string;
  fiatGlyph: string;
  minTokenizationFee: number;
  minTROVBalanceForTokenizationApplication: number;
  tokenizationApplicationFee: number;
  tokenizationApplicationFeeAsset: string;
  vatPercent: number;
}

export interface BankInfo {
  id: number;
  bankName: string;
  countryCode: string;
}

export interface TokenizationQueryParams {
  assetTokenizationStatus?: number;
  vettingStatus?: number;
  onlyWithUserPermission?: number;
  assetDescription?: string;
  salesList?: number;
  assetType?: string;
  assetName?: string;
  assetCode?: string;
  assetSector?: string;
  assetSubSector?: string;
  initiatorUsername?: string;
  offeringType?: string;
  hasSecApproval?: number;
  limit?: number;
  page?: number;
  createdBetween?: string;
}
export interface VetTokenizedAssetRequest {
  tokenizedAssetID: string;
  vettingPayload: {
    CountryCode: string;
    approvedAssetCustodianId: number;
    assetManagerId: number;
    assetIssuingHouseId: number;
    legalAndProfesionalPartnerId: number;
    ratingAgencyId: number;
    trusteeId: number;
    legalAdviserId?: number;
    financialAdviserId?: number;
    assetQuoteCurrency: string;
    proceedPayoutCurrency: string;
  };
}

export interface AssetTokenizationDocument {
  createdAt: string;
  documentTitle: string;
  documentType: string;
  documentUrl: string;
  id: number;
  tokenizedAssetId: string;
}

export interface ProofOfPaymentDocument {
  createdAt: string;
  documentUrl: string;
  id: number;
}

export interface ApprovedAssetCustodianInfo {
  assetCustodianAddress: string;
  assetCustodianCountry: string;
  assetCustodianName: string;
  id: number;
  requirementDocument: string;
}

export interface AssetManagerInfo {
  id: number;
  name: string;
}

export interface ClosedGroupInfo {
  groupDescription: string;
  groupName: string;
  groupOwner: string;
  id: string;
  registeredEntity: number;
  registrationDocumentUrl: string;
  registrationName: string;
  registrationNumber: string;
}

export interface TokenizationFee {
  feeAssetPercentage: number;
  feeDescription: string;
  feeFiatCap: number;
  feeFiatPercentage: number;
  id: number;
}

export interface LegalAndProfesionalPartnerInfo {
  id: number;
  partnerName: string;
  partnerAddress: string;
  partnerCountry: string;
  FeePercent: number;
  FeeFixed: number;
}

export interface RatingAgencyInfo {
  id: number;
  agencyName: string;
  agencyAddress: string;
  agencyCountry: string;
  FeePercent: number;
  FeeFixed: number;
}

export interface TokenizationDetailResponse {
  AssetTokenizationDocuments: AssetTokenizationDocument[];
  IsFreeFromLiensAndEncumbrances: number;
  bankId: number;
  ProofOfPaymentDocuments: ProofOfPaymentDocument[];
  SECTokenizationFeePercent: number;
  SECTokenizationFeeValue: number;
  additionalKYCRequirements: string;
  agreeTransferTitleToCustodian: number;
  approvedAssetCustodianId: number;
  approvedAssetCustodianInfo: ApprovedAssetCustodianInfo;
  assetIssuingHouseId: number;
  assetIssuingHouseInfo: AssetIssuingHouseInfo;
  issuingHouseFeeValue: number;
  issuingHouseFeePercent: number;
  ratingAgency: string;
  ratingAgencyFeePercent: number;
  legalAndProfessionalFeePercent: number;
  vatPercent: number;
  vatInAsset: number;
  SECTokenizationFeeFixed: number;
  custodianFeeFixed: number;
  assetManagerFeeFixed: number;
  issuingHouseFeeFixed: number;
  legalAndProfessionalFeeFixed: number;
  ratingAgencyFeeFixed: number;
  vatValue: number;
  legalAndProfessionalFeeValue: number;
  ratingAgencyFeeValue: number;
  trusteeFeePercent: number;
  trusteeFeeFixed: number;
  trusteeFeeValue: number;
  legalAndProfesionalPartnerInfo: LegalAndProfesionalPartnerInfo;
  ratingAgencyInfo: RatingAgencyInfo;

  tokenizationApplicationFee: number;
  tokenizationApplicationFeeAsset: string;

  assetAlreadyExists: number;
  assetCode: string;
  assetCountryLocation: string;
  assetCurrentValue: number;
  assetDescription: string;
  assetLatitude: string;
  assetLogo: string;
  assetLongitude: string;
  assetManagerFeePercent: number;
  assetManagerFeeValue: number;
  assetManagerId: number;
  assetManagerInfo: AssetManagerInfo;
  assetMscCostOutisdeOfValuation: number;
  assetName: string;
  assetOwnerAddress: string;
  assetOwnerName: string;
  assetOwnerRetainedOrContributedValue: number;
  assetPhysicalAddress: string;
  assetQuoteCurrency: string;
  assetSector: string;
  assetSubSector: string;
  assetTokenizationStatus: number;
  assetType: string;
  assetWebsite: string;
  capDurationInDays: number;
  capOnPurchase: number;
  capQuantity: number;
  capAmountInFiat: number;
  closedGroupId: string;
  closedGroupInfo: ClosedGroupInfo;
  complianceAllPermits: number;
  complianceNoViolation: number;
  contractualProtectionPerfBond: number;
  contractualProtectionRevGuarantees: number;
  contractualProtectionSLA: number;
  createdAt: string;
  custodianFeePercent: number;
  feeFiatCap: number;
  custodianFeeValue: number;
  eSGSafeguardsCommEngPlans: number;
  eSGSafeguardsSusCerts: number;
  exemptedCountries: string;
  feeInAsset: number;
  feeInFiat: number;
  financialAdvisor: string;
  hasAdditionalKYCRequirements: number;
  id: string;
  independentMonitoringList: string;
  initialOwnerPreferredWalletAddress: string;
  initiatorUsername: string;
  insuranceCompanyName: string;
  insurancePolicyHolder: string;
  insurancePolicyNumber: string;
  investorAccreditationRequired: number;
  issuingWalletAlias: string;
  issuingWalletPublicKey: string;
  lastUpdatedBy: string;
  legalAdvisor: string;
  marketMakingWallet: string;
  maxNumberOfTokenAvailableForSale: number;
  numberOfTokenToBeIssued: number;
  numberOfTokenToBeSold: number;
  offeringType: string;
  otherAssetProtection: string;
  outstandingFinancialRespNoDebts: number;
  outstandingFinancialRespNoHiddenLiabilities: number;
  ownershipKind: string;
  ownershipType: string;
  percentageValueOfInsurance: number;
  physicalConditionNoUndisclosedEasements: number;
  physicalConditionNolease: number;
  physicalConditionSound: number;
  pricePerToken: number;
  proceedCycle: string;
  proceedPayoutCurrency: string;
  proceedPayoutType: number;
  protectionMethods: string;
  riskManagementDeclaredValue: number;
  riskManagementFullyInsured: number;
  riskSharingMechanismCompletionGuarantees: number;
  riskSharingMechanismHedgeInstruments: number;
  riskSharingMechanismPPPs: number;
  salesEnd: string;
  salesStart: string;
  secApproval: number;
  secApprovalIdNumber: string;
  securityMeasuresAccessControl: number;
  securityMeasuresCriticalInfraProtections: number;
  securityMeasuresOnSiteSecurityPersonnel: number;
  securityMeasuresPerimeterSecurity: number;
  securityMeasuresSurveilanceSystems: number;
  tokenizationFee: TokenizationFee;
  tokenizationFeeId: number;
  totalTokenHeldByManager: number;
  undertakingNoClaims: number;
  undertakingNoForeclosure: number;
  undertakingNoLien: number;
  undertakingNotCollateral: number;
  updatedAt: string;
  valueOfTokenizedAsset: number;
  vettingStatus: number;
  walletToHoldAssetsNotForSale: string;
  dueDiligenceFail: number;
  dueDiligenceFailureReason: string;
  countryConfig: CountryConfig;
  mintingInitators: string;
  mintingApprovers: string;
  bankInfo: BankInfo;
  accountNumber: string;
  beneficiaryName: string;
  legalAndProfesionalPartnerId: number;
  ratingAgencyId: number;
  trusteeId: number;
  trusteeInfo: TrusteeInfo;
  purchaseCommitments: number;
  quantityOfTokensSold: number;
  quantityOfTokensSoldInFiat: number;
  numberOfSubscribers: number;
  numberOfExpressedInterests: number;
}
export interface TokenizationDetailApiResponse {
  message: string;
  data: TokenizationDetailResponse;
  timestamp: string;
  status: string;
}

export interface DueDiligenceFailReason {
  tokenizedAssetID: string;
  reason: string;
}

export interface UpdateTokenizationPayload {
  IsFreeFromLiensAndEncumbrances?: number;
  accountNumber?: string;
  additionalKYCRequirements?: string;
  agreeTransferTitleToCustodian?: number;
  approvedAssetCustodianId?: number;
  assetAlreadyExists?: number;
  assetCode?: string;
  assetCountryLocation?: string;
  assetCurrentValue?: number;
  assetDescription?: string;
  assetLatitude?: string;
  assetLogo?: string;
  assetLongitude?: string;
  assetManagerId?: number;
  assetMscCostOutisdeOfValuation?: number;
  assetName?: string;
  assetOwnerAddress?: string;
  assetOwnerName?: string;
  assetOwnerRetainedOrContributedValue?: number;
  assetPhysicalAddress?: string;
  assetQuoteCurrency?: string;
  assetSector?: string;
  assetSubSector?: string;
  assetType?: string;
  assetWebsite?: string;
  bankId?: number;
  beneficiaryName?: string;
  capDurationInDays?: number;
  capOnPurchase?: number;
  capQuantity?: number;
  capAmountInFiat?: number;
  closedGroupId?: string;
  complianceAllPermits?: number;
  complianceNoViolation?: number;
  contractualProtectionPerfBond?: number;
  contractualProtectionRevGuarantees?: number;
  contractualProtectionSLA?: number;
  eSGSafeguardsCommEngPlans?: number;
  eSGSafeguardsSusCerts?: number;
  exemptedCountries?: string;
  financialAdvisor?: string;
  hasAdditionalKYCRequirements?: number;
  independentMonitoringList?: string;
  initialOwnerPreferredWalletAddress?: string;
  insuranceCompanyName?: string;
  insurancePolicyHolder?: string;
  insurancePolicyNumber?: string;
  investorAccreditationRequired?: number;
  legalAdvisor?: string;
  marketMakingWallet?: string;

  mintingApprovers?: string;
  mintingInitators?: string;
  numberOfTokenToBeIssued?: number;
  numberOfTokenToBeSold?: number;
  offeringType?: string;
  otherAssetProtection?: string;
  outstandingFinancialRespNoDebts?: number;
  outstandingFinancialRespNoHiddenLiabilities?: number;
  ownershipKind?: string;
  ownershipType?: string;
  percentageValueOfInsurance?: number;
  physicalConditionNoUndisclosedEasements?: number;
  physicalConditionNolease?: number;
  physicalConditionSound?: number;
  pricePerToken?: number;
  proceedCycle?: string;
  proceedPayoutCurrency?: string;
  proceedPayoutType?: number;
  protectionMethods?: string;
  riskManagementDeclaredValue?: number;
  riskManagementFullyInsured?: number;
  riskSharingMechanismCompletionGuarantees?: number;
  riskSharingMechanismHedgeInstruments?: number;
  riskSharingMechanismPPPs?: number;
  salesEnd?: string;
  salesStart?: string;
  secApproval?: number;
  secApprovalIdNumber?: string;
  securityMeasuresAccessControl?: number;
  securityMeasuresCriticalInfraProtections?: number;
  securityMeasuresOnSiteSecurityPersonnel?: number;
  securityMeasuresPerimeterSecurity?: number;
  securityMeasuresSurveilanceSystems?: number;
  tokenizationFeeId?: number;
  totalTokenHeldByManager?: number;
  undertakingNoClaims?: number;
  undertakingNoForeclosure?: number;
  undertakingNoLien?: number;
  undertakingNotCollateral?: number;
  walletToHoldAssetsNotForSale?: string;
  AssetTokenizationDocuments?: AssetTokenizationDocument[];
  ProofOfPaymentDocuments?: ProofOfPaymentDocument[];
  SECTokenizationFeePercent?: number;
  SECTokenizationFeeValue?: number;
  approvedAssetCustodianInfo?: ApprovedAssetCustodianInfo;
  assetTokenizationStatus?: number;
  assetManagerFeePercent?: number;
  assetManagerFeeValue?: number;
  assetManagerInfo?: AssetManagerInfo;
  bankInfo?: BankInfo;
  closedGroupInfo?: ClosedGroupInfo;
  countryConfig?: CountryConfig;
  createdAt?: string;
  custodianFeePercent?: number;
  custodianFeeValue?: number;
  dueDiligenceFail?: number;
  dueDiligenceFailureReason?: string;
  feeInAsset?: number;
  feeInFiat?: number;
  id?: string;
  initiatorUsername?: string;
  issuingWalletAlias?: string;
  issuingWalletPublicKey?: string;
  lastUpdatedBy?: string;
  maxNumberOfTokenAvailableForSale?: number;
  tokenizationFee?: TokenizationFee;
  updatedAt?: string;
  valueOfTokenizedAsset?: number;
  vettingStatus?: number;
}

export interface AssetSector {
  id: number;
  sectorName: string;
}

export interface AssetSubSector {
  id: number;
  subSectorName: string;
  assetSectorId: string;
}

export interface AssetType {
  id: number;
  assetSubSectorId: string;
  assetType: string;
}

export interface ITokenizationParameterCustodian {
  id: number;
  assetCustodianName: string;
  assetCustodianAddress: string;
  assetCustodianCountry: string;
  requirementDocument: string;
  FeePercent: number;
  FeeFixed: number;
}

export interface ITokenizationParameterIssuingHouse {
  id: number;
  assetIssuingHouseName: string;
  assetIssuingHouseAddress: string;
  assetIssuingHouseCountry: string;
  FeePercent: number;
  FeeFixed: number;
}

export interface ITokenizationParameterManager {
  id: number;
  assetManagerName: string;
  assetManagerAddress: string;
  assetManagerCountry: string;
  FeePercent: number;
  FeeFixed: number;
}

export interface IAssetProceedCycle {
  id: string;
}

export interface IAssetProtectionOption {
  id: string;
}

export interface IAssetSectorParam {
  sector: string;
  requirementDocument: string;
}

export interface IAssetSubSectorParam {
  subSector: string;
  assetSectorId: string;
}

export interface IAssetTypeParam {
  id: number;
  assetSubSectorId: string;
  assetType: string;
}

export interface AssetDocument {
  documentType: string;
  documentTypeDescription: string;
  documentCategory: string;
  isPublic: number; // 0 or 1 (since your API uses numeric boolean)
}

export interface TokenizationParamsResponse {
  assetTypes: IAssetTypeParam[];
  assetSectors?: IAssetSectorParam[];
  assetSubSectors?: IAssetSubSectorParam[];
  assetCustodians: ITokenizationParameterCustodian[];
  assetIssuingHouses: ITokenizationParameterIssuingHouse[];
  assetManagers: ITokenizationParameterManager[];
  assetProceedCycle: IAssetProceedCycle[];
  assetProtectionOptions: IAssetProtectionOption[];
  tokenizationDocumentTypes: AssetDocument[];
}

// tokenized asset statistics
export interface TokenizedAssetStats {
  count: number;
  totalCurrentValue: number;
  totalTokenizedValue: number;
}

export interface TokenizedAssetsStatistics {
  total: TokenizedAssetStats;
  approved: TokenizedAssetStats;
  pending: TokenizedAssetStats;
  rejected: TokenizedAssetStats;
  refunded: TokenizedAssetStats;
  liquidated: TokenizedAssetStats;
  submitted: TokenizedAssetStats;
  unsubmitted: TokenizedAssetStats;
}

export interface TokenizedAssetStatisticsResponse {
  success: boolean;
  message: string;
  errorCode: string | null;
  data: TokenizedAssetsStatistics;
}
