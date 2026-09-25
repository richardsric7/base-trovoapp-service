import {
  TokenizationRecord,
  UpdateTokenizationPayload,
} from "@/redux/api/assettokenization";

/**
 * Sanitizes an asset object for update submission
 */
export const getCleanedUpdatePayload = (
  asset?: TokenizationRecord
): UpdateTokenizationPayload => {
  if (!asset) throw new Error("Asset is required");

  const fallback = <T>(value: T | undefined | null, defaultValue: T): T =>
    value ?? defaultValue;

  return {
    IsFreeFromLiensAndEncumbrances: fallback(
      asset.IsFreeFromLiensAndEncumbrances,
      0
    ),
    accountNumber: fallback(asset.accountNumber, ""),
    additionalKYCRequirements: fallback(asset.additionalKYCRequirements, ""),
    agreeTransferTitleToCustodian: fallback(
      asset.agreeTransferTitleToCustodian,
      0
    ),
    approvedAssetCustodianId: fallback(asset.approvedAssetCustodianId, 0),
    assetAlreadyExists: fallback(asset.assetAlreadyExists, 0),
    assetCode: fallback(asset.assetCode, ""),
    assetCountryLocation: fallback(asset.assetCountryLocation, ""),
    assetCurrentValue: fallback(asset.assetCurrentValue, 0),
    assetDescription: fallback(asset.assetDescription, ""),
    assetLatitude: fallback(asset.assetLatitude, ""),
    assetLogo: fallback(asset.assetLogo, ""),
    assetLongitude: fallback(asset.assetLongitude, ""),
    assetManagerId: fallback(asset.assetManagerId, 0),
    assetMscCostOutisdeOfValuation: fallback(
      asset.assetMscCostOutisdeOfValuation,
      0
    ),
    assetName: fallback(asset.assetName, ""),
    assetOwnerAddress: fallback(asset.assetOwnerAddress, ""),
    assetOwnerName: fallback(asset.assetOwnerName, ""),
    assetOwnerRetainedOrContributedValue: fallback(
      asset.assetOwnerRetainedOrContributedValue,
      0
    ),
    assetPhysicalAddress: fallback(asset.assetPhysicalAddress, ""),
    assetQuoteCurrency: fallback(asset.assetQuoteCurrency, ""),
    assetSector: fallback(asset.assetSector, ""),
    assetSubSector: fallback(asset.assetSubSector, ""),
    assetType: fallback(asset.assetType, ""),
    assetWebsite: fallback(asset.assetWebsite, ""),
    bankId: fallback(asset.bankId, 0),
    beneficiaryName: fallback(asset.beneficiaryName, ""),
    capDurationInDays: fallback(asset.capDurationInDays, 0),
    capOnPurchase: fallback(asset.capOnPurchase, 0),
    capQuantity: fallback(asset.capQuantity, 0),
    capAmountInFiat: fallback(asset.capAmountInFiat, 0),
    closedGroupId: fallback(asset.closedGroupId, ""),
    complianceAllPermits: fallback(asset.complianceAllPermits, 0),
    complianceNoViolation: fallback(asset.complianceNoViolation, 0),
    contractualProtectionPerfBond: fallback(
      asset.contractualProtectionPerfBond,
      0
    ),
    contractualProtectionRevGuarantees: fallback(
      asset.contractualProtectionRevGuarantees,
      0
    ),
    contractualProtectionSLA: fallback(asset.contractualProtectionSLA, 0),
    eSGSafeguardsCommEngPlans: fallback(asset.eSGSafeguardsCommEngPlans, 0),
    eSGSafeguardsSusCerts: fallback(asset.eSGSafeguardsSusCerts, 0),
    exemptedCountries: fallback(asset.exemptedCountries, ""),
    financialAdvisor: fallback(asset.financialAdvisor, ""),
    hasAdditionalKYCRequirements: fallback(
      asset.hasAdditionalKYCRequirements,
      0
    ),
    independentMonitoringList: fallback(asset.independentMonitoringList, ""),
    initialOwnerPreferredWalletAddress: fallback(
      asset.initialOwnerPreferredWalletAddress,
      ""
    ),
    insuranceCompanyName: fallback(asset.insuranceCompanyName, ""),
    insurancePolicyHolder: fallback(asset.insurancePolicyHolder, ""),
    insurancePolicyNumber: fallback(asset.insurancePolicyNumber, ""),
    investorAccreditationRequired: fallback(
      asset.investorAccreditationRequired,
      0
    ),
    legalAdvisor: fallback(asset.legalAdvisor, ""),
    marketMakingWallet: fallback(asset.marketMakingWallet, ""),
    mintingApprovers: fallback(asset.mintingApprovers, ""),
    mintingInitators: fallback(asset.mintingInitators, ""),
    numberOfTokenToBeIssued: fallback(asset.numberOfTokenToBeIssued, 0),
    numberOfTokenToBeSold: fallback(asset.numberOfTokenToBeSold, 0),
    offeringType: fallback(asset.offeringType, ""),
    otherAssetProtection: fallback(asset.otherAssetProtection, ""),
    outstandingFinancialRespNoDebts: fallback(
      asset.outstandingFinancialRespNoDebts,
      0
    ),
    outstandingFinancialRespNoHiddenLiabilities: fallback(
      asset.outstandingFinancialRespNoHiddenLiabilities,
      0
    ),
    ownershipKind: fallback(asset.ownershipKind, ""),
    ownershipType: fallback(asset.ownershipType, ""),
    percentageValueOfInsurance: fallback(asset.percentageValueOfInsurance, 0),
    physicalConditionNoUndisclosedEasements: fallback(
      asset.physicalConditionNoUndisclosedEasements,
      0
    ),
    physicalConditionNolease: fallback(asset.physicalConditionNolease, 0),
    physicalConditionSound: fallback(asset.physicalConditionSound, 0),
    pricePerToken: fallback(asset.pricePerToken, 0),
    proceedCycle: fallback(asset.proceedCycle, ""),
    proceedPayoutCurrency: fallback(asset.proceedPayoutCurrency, ""),
    proceedPayoutType: fallback(asset.proceedPayoutType, 0),
    protectionMethods: fallback(asset.protectionMethods, ""),
    riskManagementDeclaredValue: fallback(asset.riskManagementDeclaredValue, 0),
    riskManagementFullyInsured: fallback(asset.riskManagementFullyInsured, 0),
    riskSharingMechanismCompletionGuarantees: fallback(
      asset.riskSharingMechanismCompletionGuarantees,
      0
    ),
    riskSharingMechanismHedgeInstruments: fallback(
      asset.riskSharingMechanismHedgeInstruments,
      0
    ),
    riskSharingMechanismPPPs: fallback(asset.riskSharingMechanismPPPs, 0),
    salesEnd: fallback(asset.salesEnd, ""),
    salesStart: fallback(asset.salesStart, ""),
    secApproval: fallback(asset.secApproval, 0),
    secApprovalIdNumber: fallback(asset.secApprovalIdNumber, ""),
    securityMeasuresAccessControl: fallback(
      asset.securityMeasuresAccessControl,
      0
    ),
    securityMeasuresCriticalInfraProtections: fallback(
      asset.securityMeasuresCriticalInfraProtections,
      0
    ),
    securityMeasuresOnSiteSecurityPersonnel: fallback(
      asset.securityMeasuresOnSiteSecurityPersonnel,
      0
    ),
    securityMeasuresPerimeterSecurity: fallback(
      asset.securityMeasuresPerimeterSecurity,
      0
    ),
    securityMeasuresSurveilanceSystems: fallback(
      asset.securityMeasuresSurveilanceSystems,
      0
    ),
    tokenizationFeeId: fallback(asset.tokenizationFeeId, 0),
    totalTokenHeldByManager: fallback(asset.totalTokenHeldByManager, 0),
    undertakingNoClaims: fallback(asset.undertakingNoClaims, 0),
    undertakingNoForeclosure: fallback(asset.undertakingNoForeclosure, 0),
    undertakingNoLien: fallback(asset.undertakingNoLien, 0),
    undertakingNotCollateral: fallback(asset.undertakingNotCollateral, 0),
    walletToHoldAssetsNotForSale: fallback(
      asset.walletToHoldAssetsNotForSale,
      ""
    ),
  };
};
