import 'dart:developer';

class TokenizedAsset {
  String? id;
  String? shadowId;
  // double? amount;
  double? usdPrice;
  String? assetIssuer;
  String? assetCode;
  String? assetName;
  String? assetSector;
  String? assetSubSector;
  String? assetType;
  String? offeringType; // 0:public, 1: private
  int? secApproval;
  String? secApprovalIdNumber;
  String? assetCountryLocation;
  int? approvedAssetCustodianId;
  int? hasAllRequiredDocuments;
  int? hasCustodianAgreement;
  String? marketMakingWallet;
  int? assetAlreadyExists;
  String? ownershipType; // DIRECT | THIRD-PARTY
  String? ownershipKind; // INDIVIDUAL | CORPORATE
  String? assetDescription;
  String? assetPhysicalAddress;
  String? assetLatitude;
  String? assetLongitude;
  String? assetOwnerName;
  String? assetOwnerAddress;
  ManagerInfo? assetManagerInfo;
  double? assetCurrentValue;
  double? assetPercentageForTokenization;
  double? valueOfTokenizedAsset;
  String? protectionMethods;
  String? insuranceCompanyName;
  String? insurancePolicyNumber;
  String? insurancePolicyHolder;
  double? percentageValueOfInsurance;
  int? isFreeFromLiensAndEncumbrances;
  int? tokenizationFeeId;
  int? vettingStatus;
  double? numberOfTokenToBeSold;
  double? numberOfTokenToBeIssued;
  String? walletToHoldAssetsNotForSale;
  double? totalTokenHeldByManager;
  double? pricePerToken;
  DateTime? salesStart;
  DateTime? salesEnd;
  int? capOnPurchase;
  double? capQuantity;
  int? capDurationInDays;
  String? proceedCycle;
  String? assetLogo;
  bool? expressedInterest;
  double? expressedInterestAmount;
  bool? isSubscribed;
  double? subscriptionAmount;
  String? exemptedCountries;
  int? hasAdditionalKYCRequirements;
  String? proceedPayoutCurrency;
  String? assetQuoteCurrency;
  String? additionalKYCRequirements;
  int? investorAccreditationRequired;
  int?
      tokenizationStatus; // 0 pending, 1 submitted (awaiting fee payment), 2 submitted (processing)
  List<ProofOfPaymentDocument>? proofOfPaymentDocuments;
  String? lastUpdatedBy;
  DateTime? createdAt;
  DateTime? updatedAt;

  CustodianInfo? approvedAssetCustodianInfo;
  String? initiatorUsername;
  String? closedGroupId;
  ClosedGroupInfo? closedGroupInfo;
  List<Document>? assetTokenizationDocuments;
  int? agreeTransferTitleToCustodian;
  int? contractualProtectionRevGuarantees;
  int? contractualProtectionPerfBond;
  int? contractualProtectionSLA;
  int? riskSharingMechanismPPPs;
  int? riskSharingMechanismHedgeInstruments;
  int? riskSharingMechanismCompletionGuarantees;
  String? independentMonitoringList;
  int? eSGSafeguardsSusCerts;
  int? eSGSafeguardsCommEngPlans;
  int? securityMeasuresAccessControl;
  int? securityMeasuresSurveilanceSystems;
  int? securityMeasuresOnSiteSecurityPersonnel;
  int? securityMeasuresPerimeterSecurity;
  int? securityMeasuresCriticalInfraProtections;
  String? otherAssetProtection;
  String? legalAdvisor;
  String? financialAdvisor;
  int? undertakingNoLien;
  int? undertakingNotCollateral;
  int? undertakingNoClaims;
  int? undertakingNoForeclosure;
  int? complianceNoViolation;
  int? complianceAllPermits;
  int? outstandingFinancialRespNoDebts;
  int? outstandingFinancialRespNoHiddenLiabilities;
  int? riskManagementFullyInsured;
  int? riskManagementDeclaredValue;
  int? physicalConditionSound;
  int? physicalConditionNolease;
  int? assetMscCostOutisdeOfValuation;
  double? SECTokenizationFeePercent;
  double? SECTokenizationFeeValue;
  double? custodianFeeValue;
  double? custodianFeePercent;
  double? assetManagerFeeValue;
  double? assetManagerFeePercent;
  double? assetOwnerRetainedOrContributedValue;

  TokenizedAsset({
    this.id,
    this.shadowId,
    // this.amount,
    this.usdPrice,
    this.assetIssuer,
    this.assetCode,
    this.assetName,
    this.assetSector,
    this.assetSubSector,
    this.assetType,
    this.offeringType,
    this.secApproval,
    this.secApprovalIdNumber,
    this.assetCountryLocation,
    this.approvedAssetCustodianId,
    this.hasAllRequiredDocuments,
    this.hasCustodianAgreement,
    this.marketMakingWallet,
    this.assetAlreadyExists,
    this.ownershipType,
    this.ownershipKind,
    this.assetDescription,
    this.assetPhysicalAddress,
    this.assetLatitude,
    this.assetLongitude,
    this.assetOwnerName,
    this.assetOwnerAddress,
    this.assetManagerInfo,
    this.assetCurrentValue,
    this.assetPercentageForTokenization,
    this.valueOfTokenizedAsset,
    this.protectionMethods,
    this.insuranceCompanyName,
    this.insurancePolicyNumber,
    this.insurancePolicyHolder,
    this.percentageValueOfInsurance,
    this.isFreeFromLiensAndEncumbrances,
    this.tokenizationFeeId,
    this.numberOfTokenToBeSold,
    this.numberOfTokenToBeIssued,
    this.walletToHoldAssetsNotForSale,
    this.totalTokenHeldByManager,
    this.pricePerToken,
    this.salesStart,
    this.salesEnd,
    this.capOnPurchase,
    this.capQuantity,
    this.capDurationInDays,
    this.proceedCycle,
    this.vettingStatus,
    this.assetLogo,
    this.exemptedCountries,
    this.hasAdditionalKYCRequirements,
    this.proceedPayoutCurrency,
    this.additionalKYCRequirements,
    this.investorAccreditationRequired,
    this.lastUpdatedBy,
    this.createdAt,
    this.updatedAt,
    this.approvedAssetCustodianInfo,
    this.initiatorUsername,
    this.closedGroupId,
    this.closedGroupInfo,
    this.assetTokenizationDocuments,
    this.tokenizationStatus,
    this.assetQuoteCurrency,
    this.isSubscribed,
    this.subscriptionAmount,
    this.expressedInterest,
    this.expressedInterestAmount,
    this.proofOfPaymentDocuments,
    this.agreeTransferTitleToCustodian,
    this.contractualProtectionRevGuarantees,
    this.contractualProtectionPerfBond,
    this.contractualProtectionSLA,
    this.riskSharingMechanismPPPs,
    this.riskSharingMechanismHedgeInstruments,
    this.riskSharingMechanismCompletionGuarantees,
    this.independentMonitoringList,
    this.eSGSafeguardsSusCerts,
    this.eSGSafeguardsCommEngPlans,
    this.securityMeasuresAccessControl,
    this.securityMeasuresSurveilanceSystems,
    this.securityMeasuresOnSiteSecurityPersonnel,
    this.securityMeasuresPerimeterSecurity,
    this.securityMeasuresCriticalInfraProtections,
    this.otherAssetProtection,
    this.legalAdvisor,
    this.financialAdvisor,
    this.undertakingNoLien,
    this.undertakingNotCollateral,
    this.undertakingNoClaims,
    this.undertakingNoForeclosure,
    this.complianceNoViolation,
    this.complianceAllPermits,
    this.outstandingFinancialRespNoDebts,
    this.outstandingFinancialRespNoHiddenLiabilities,
    this.riskManagementFullyInsured,
    this.riskManagementDeclaredValue,
    this.physicalConditionSound,
    this.physicalConditionNolease,
    this.assetMscCostOutisdeOfValuation,
    this.SECTokenizationFeePercent,
    this.SECTokenizationFeeValue,
    this.custodianFeeValue,
    this.custodianFeePercent,
    this.assetManagerFeeValue,
    this.assetManagerFeePercent,
    this.assetOwnerRetainedOrContributedValue,
  });

  TokenizedAsset deserializeJson(Map<String, dynamic> m) {
    print('deserializing json $m');
    inspect(m);
    return TokenizedAsset(
      id: m["id"],
      shadowId: m["shadowId"],
      // amount: m["amount"],
      usdPrice: m["usdPrice"],
      assetCode: m["assetCode"],
      assetIssuer: m["assetIssuer"],
      assetName: m["assetName"],
      assetSector: m["assetSector"],
      assetSubSector: m["assetSubSector"],
      assetType: m["assetType"],
      offeringType: m["offeringType"],
      secApproval: m["secApproval"],
      secApprovalIdNumber: m["secApprovalIdNumber"],
      assetCountryLocation: m["assetCountryLocation"],
      approvedAssetCustodianId: m["approvedAssetCustodianId"],
      hasAllRequiredDocuments: m["hasAllRequiredDocuments"],
      hasCustodianAgreement: m["hasCustodianAgreement"],
      marketMakingWallet: m["marketMakingWallet"],
      assetAlreadyExists: m["assetAlreadyExists"],
      ownershipType: m["ownershipType"],
      ownershipKind: m["ownershipKind"],
      assetDescription: m["assetDescription"],
      assetPhysicalAddress: m["assetPhysicalAddress"],
      assetLatitude: m["assetLatitude"],
      assetLongitude: m["assetLongitude"],
      assetOwnerName: m["assetOwnerName"],
      assetOwnerAddress: m["assetOwnerAddress"],
      vettingStatus: m["vettingStatus"],
      assetManagerInfo: ManagerInfo().deserializeJson(m),
      assetQuoteCurrency: m["assetQuoteCurrency"],
      assetCurrentValue: double.parse(m["assetCurrentValue"].toString()),
      assetOwnerRetainedOrContributedValue:
          double.parse(m["assetOwnerRetainedOrContributedValue"].toString()),
      valueOfTokenizedAsset:
          double.parse(m["valueOfTokenizedAsset"].toString()),
      protectionMethods: m["protectionMethods"],
      insuranceCompanyName: m["insuranceCompanyName"],
      insurancePolicyNumber: m["insurance_policy_number"],
      insurancePolicyHolder: m["insurancePolicyHolder"],
      percentageValueOfInsurance:
          double.parse(m["percentageValueOfInsurance"].toString()),
      isFreeFromLiensAndEncumbrances: m["isFreeFromLiensAndEncumbrances"],
      tokenizationFeeId: m["tokenizationFeeId"],
      numberOfTokenToBeSold:
          double.parse(m["numberOfTokenToBeSold"].toString()),
      numberOfTokenToBeIssued:
          double.parse(m["numberOfTokenToBeIssued"].toString()),
      walletToHoldAssetsNotForSale: m["walletToHoldAssetsNotForSale"],
      totalTokenHeldByManager:
          double.parse(m["totalTokenHeldByManager"].toString()),
      pricePerToken: double.parse(m["pricePerToken"].toString()),
      salesStart: DateTime.parse(m["salesStart"]),
      salesEnd: DateTime.parse(m["salesEnd"]),
      capOnPurchase: m["capOnPurchase"],
      capQuantity: double.parse(m["capQuantity"].toString()),
      capDurationInDays: m["capDurationInDays"],
      proceedCycle: m["proceedCycle"],
      assetLogo: m["assetLogo"],
      exemptedCountries: m["exemptedCountries"],
      hasAdditionalKYCRequirements: m["hasAdditionalKYCRequirements"],
      proceedPayoutCurrency: m["proceedPayoutCurrency"],
      additionalKYCRequirements: m["additionalKYCRequirements"],
      investorAccreditationRequired: m["investorAccreditationRequired"],
      lastUpdatedBy: m["lastUpdatedBy"],
      createdAt: DateTime.parse(m["createdAt"]),
      updatedAt: DateTime.parse(m["updatedAt"]),
      approvedAssetCustodianInfo:
          CustodianInfo().deserializeJson(m["approvedAssetCustodianInfo"]),
      initiatorUsername: m["initiatorUsername"],
      closedGroupId: m["closedGroupId"],
      closedGroupInfo: ClosedGroupInfo().deserializeJson(m["closedGroupInfo"]),
      assetTokenizationDocuments:
          deserializeDocuments(m["AssetTokenizationDocuments"]),
      tokenizationStatus: m["assetTokenizationStatus"],
      agreeTransferTitleToCustodian: m["agreeTransferTitleToCustodian"],
      contractualProtectionRevGuarantees:
          m["contractualProtectionRevGuarantees"],
      proofOfPaymentDocuments:
          deserializeProofOfPaymentDocuments(m["ProofOfPaymentDocuments"]),
      contractualProtectionPerfBond: m["contractualProtectionPerfBond"],
      contractualProtectionSLA: m["contractualProtectionSLA"],
      riskSharingMechanismPPPs: m["riskSharingMechanismPPPs"],
      riskSharingMechanismHedgeInstruments:
          m["riskSharingMechanismHedgeInstruments"],
      riskSharingMechanismCompletionGuarantees:
          m["riskSharingMechanismCompletionGuarantees"],
      independentMonitoringList: m["independentMonitoringList"],
      eSGSafeguardsSusCerts: m["eSGSafeguardsSusCerts"],
      eSGSafeguardsCommEngPlans: m["eSGSafeguardsCommEngPlans"],
      securityMeasuresAccessControl: m["securityMeasuresAccessControl"],
      securityMeasuresSurveilanceSystems:
          m["securityMeasuresSurveilanceSystems"],
      securityMeasuresOnSiteSecurityPersonnel:
          m["securityMeasuresOnSiteSecurityPersonnel"],
      securityMeasuresPerimeterSecurity: m["securityMeasuresPerimeterSecurity"],
      securityMeasuresCriticalInfraProtections:
          m["securityMeasuresCriticalInfraProtections"],
      otherAssetProtection: m["otherAssetProtection"],
      legalAdvisor: m["legalAdvisor"],
      financialAdvisor: m["financialAdvisor"],
      undertakingNoLien: m["undertakingNoLien"],
      undertakingNotCollateral: m["undertakingNotCollateral"],
      undertakingNoClaims: m["undertakingNoClaims"],
      undertakingNoForeclosure: m["undertakingNoForeclosure"],
      complianceNoViolation: m["complianceNoViolation"],
      complianceAllPermits: m["complianceAllPermits"],
      outstandingFinancialRespNoDebts: m["outstandingFinancialRespNoDebts"],
      outstandingFinancialRespNoHiddenLiabilities:
          m["outstandingFinancialRespNoHiddenLiabilities"],
      riskManagementFullyInsured: m["riskManagementFullyInsured"],
      riskManagementDeclaredValue: m["riskManagementDeclaredValue"],
      physicalConditionSound: m["physicalConditionSound"],
      physicalConditionNolease: m["physicalConditionNolease"],
      assetMscCostOutisdeOfValuation: m["assetMscCostOutisdeOfValuation"],
      SECTokenizationFeePercent:
          double.tryParse(m["SECTokenizationFeePercent"].toString()),
      SECTokenizationFeeValue:
          double.tryParse(m["SECTokenizationFeeValue"].toString()),
      custodianFeeValue: double.tryParse(m["custodianFeeValue"].toString()),
      custodianFeePercent: double.tryParse(m["custodianFeePercent"].toString()),
      assetManagerFeeValue:
          double.tryParse(m["assetManagerFeeValue"].toString()),
      assetManagerFeePercent:
          double.tryParse(m["assetManagerFeePercent"].toString()),
    );
  }

  List<Document> deserializeDocuments(m) {
    List<Document> docs = [];
    for (int i = 0; i < m.length; i++) {
      docs.add(Document().deserializeJson(m[i]));
    }
    return docs;
  }

  List<ProofOfPaymentDocument> deserializeProofOfPaymentDocuments(m) {
    List<ProofOfPaymentDocument> docs = [];
    for (int i = 0; i < m.length; i++) {
      docs.add(ProofOfPaymentDocument().deserializeJson(m[i]));
    }
    return docs;
  }
  // bool get isWithdrawable => withdrawable == 1;
  // bool get canGenerateDepositAddresses => generateDepositAddress == 1;
}

class CustodianInfo {
  int? id;
  String? assetCustodianName;
  String? assetCustodianAddress;
  String? assetCustodianCountry;
  String? requirementDocument;

  CustodianInfo({
    this.id,
    this.assetCustodianName,
    this.assetCustodianAddress,
    this.assetCustodianCountry,
    this.requirementDocument,
  });

  CustodianInfo deserializeJson(Map<String, dynamic> m) {
    return CustodianInfo(
      id: m["id"],
      assetCustodianName: m["assetCustodianName"],
      assetCustodianAddress: m["assetCustodianAddress"],
      assetCustodianCountry: m["assetCustodianCountry"],
      requirementDocument: m["requirementDocument"],
    );
  }
}

class ProofOfPaymentDocument {
  int? id;
  DateTime? createdAt;
  String? tokenizationFeePaymentMethodID;
  String? tokenizedAssetId;
  String? transactionReference;
  String? documentUrl;

  ProofOfPaymentDocument({
    this.id,
    this.createdAt,
    this.tokenizationFeePaymentMethodID,
    this.tokenizedAssetId,
    this.transactionReference,
    this.documentUrl,
  });

  ProofOfPaymentDocument deserializeJson(Map<String, dynamic> m) {
    var info =
        m["ProofOfPaymentDocument"] != null ? m["ProofOfPaymentDocument"] : m;
    return ProofOfPaymentDocument(
      createdAt: DateTime.parse(m["CreatedAt"]),
      id: info["id"] != null ? info["id"] : null,
      tokenizationFeePaymentMethodID: info["tokenizationFeePaymentMethodID"],
      tokenizedAssetId: info["tokenizedAssetId"],
      transactionReference: info["transactionReference"],
      documentUrl: info["documentUrl"],
    );
  }
}

class ManagerInfo {
  int? id;
  String? assetManagerName;
  String? assetManagerAddress;
  String? assetManagerCountry;
  String? requirementDocument;

  ManagerInfo({
    this.id,
    this.assetManagerName,
    this.assetManagerAddress,
    this.assetManagerCountry,
  });

  ManagerInfo deserializeJson(Map<String, dynamic> m) {
    var info = m["assetManagerInfo"] != null ? m["assetManagerInfo"] : m;
    return ManagerInfo(
      id: m["assetManagerInfo"] != null ? info["id"] : null,
      assetManagerName: info["assetManagerName"],
      assetManagerAddress: info["assetManagerAddress"],
      assetManagerCountry: info["assetManagerCountry"],
    );
  }
}

class ClosedGroupInfo {
  String? id;
  String? groupName;
  String? groupOwner;
  String? groupDescription;
  int? registeredEntity;
  String? registrationName;
  String? registrationNumber;
  String? registrationDocumentUrl;

  ClosedGroupInfo({
    this.id,
    this.groupName,
    this.groupOwner,
    this.groupDescription,
    this.registeredEntity,
    this.registrationName,
    this.registrationNumber,
    this.registrationDocumentUrl,
  });

  ClosedGroupInfo deserializeJson(Map<String, dynamic> m) {
    return ClosedGroupInfo(
      id: m["id"],
      groupName: m["groupName"],
      groupOwner: m["groupOwner"],
      groupDescription: m["groupDescription"],
      registeredEntity: m["registeredEntity"],
      registrationName: m["registrationName"],
      registrationNumber: m["registrationNumber"],
      registrationDocumentUrl: m["registrationDocumentUrl"],
    );
  }
}

class Document {
  int? id;
  String? tokenizedAssetId;
  String? documentType;
  String? documentTitle;
  String? documentUrl;
  DateTime? createdAt;

  Document({
    this.id,
    this.tokenizedAssetId,
    this.documentType,
    this.documentTitle,
    this.documentUrl,
    this.createdAt,
  });

  Document deserializeJson(Map<String, dynamic> m) {
    return Document(
      id: m["id"],
      tokenizedAssetId: m["tokenizedAssetId"],
      documentType: m["documentType"].toString(),
      documentTitle: m["documentTitle"],
      documentUrl: m["documentUrl"],
      createdAt: DateTime.parse(m["CreatedAt"]),
    );
  }
}
