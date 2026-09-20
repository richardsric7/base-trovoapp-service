class CuratedAsset {
  String? assetIssuer;
  String? assetCode;
  String? assetName;
  String? description;
  String? imageUrl;
  String? website;
  String? assetConditions;
  double? assetLimit;
  String? assetRedemptionInstructions;
  String? contactEmail;
  int? assetClassId;
  Map? assetClass;
  String? organization;
  int? withdrawable;
  int? generateDepositAddress;
  int? decimalPlaces;
  String? realAssetImageUrl;
  bool isRemovable = false;
  bool isCustom = false; // if the asset is added by the user and not curated by the platform

  CuratedAsset({
    this.assetIssuer,
    this.assetCode,
    this.assetName,
    this.description,
    this.imageUrl,
    this.website,
    this.assetConditions,
    this.assetLimit,
    this.assetRedemptionInstructions,
    this.contactEmail,
    this.assetClassId,
    this.assetClass,
    this.organization,
    this.withdrawable,
    this.generateDepositAddress,
    this.decimalPlaces,
    this.realAssetImageUrl,
  });

  CuratedAsset deserializeJson(Map<String, dynamic> m) {
    return CuratedAsset(
      assetCode: m["assetCode"],
      assetIssuer: m["assetIssuer"],
      assetName: m["assetName"],
      description: m["description"],
      website: m["website"],
      imageUrl: m["imageUrl"],
      assetConditions: m["assetConditions"],
      assetLimit: double.parse(m["assetLimit"].toString()),
      assetRedemptionInstructions: m["assetRedemptionInstructions"],
      contactEmail: m["contactEmail"],
      assetClassId: m["assetClassId"],
      assetClass: m["assetClass"],
      organization: m["organization"],
      withdrawable: m["withdrawable"],
      generateDepositAddress: m["generateDepositAddress"],
      decimalPlaces: m["decimalPlaces"],
      realAssetImageUrl: m["realAssetImageUrl"],
    );
  }

  bool get isWithdrawable => withdrawable == 1;
  bool get canGenerateDepositAddresses => generateDepositAddress == 1;
}
