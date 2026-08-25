import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetInformation extends StatefulWidget {
  const AssetInformation({Key? key}) : super(key: key);

  @override
  State<AssetInformation> createState() => _AssetInformation();
}

class _AssetInformation extends State<AssetInformation>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  late String assetOwnership;
  late String thirdPartyOwnerType;
  late String assetDescription;
  late String assetPhysicalAddress;
  late double latitude;
  late double longitude;
  late String nameOfOwner;
  late String assetName;
  late String addressOfOwner;
  late double currentValueOfAsset;
  late double assetMiscCost;
  late double valueOfTokenizedAsset;
  late double assetOwnerRetainedOrContributedValue;
  late List<String> assetProtectionInPlace;
  late String insuranceCompanyName;
  late String insurancePolicyNumber;
  late String insurancePolicyHolder;
  late double percentageValueOfInsurance;
  int assetDescriptionMaxWords = 100;
  bool assetAlreadyExists = false;
  bool formHasError = false;
  late dynamic data = {};
  double percentageFromPromoters = 0;

  String otherAssetProtection = "";
  String legalAdvisor = "";
  String financialAdvisor = "";

  bool hasLegalAdvisor = false;
  bool hasFinancialAdvisor = false;
  bool hasOtherAssetProtection = false;

  bool contractualProtectionRevGuarantees = false;
  bool contractualProtectionPerfBond = false;
  bool contractualProtectionSLA = false;
  bool riskSharingMechanismPPPs = false;
  bool riskSharingMechanismHedgeInstruments = false;
  bool riskSharingMechanismCompletionGuarantees = false;
  bool eSGSafeguardsSusCerts = false;
  bool eSGSafeguardsCommEngPlans = false;
  bool securityMeasuresAccessControl = false;
  bool securityMeasuresSurveilanceSystems = false;
  bool securityMeasuresOnSiteSecurityPersonnel = false;
  bool securityMeasuresPerimeterSecurity = false;
  bool securityMeasuresCriticalInfraProtections = false;
  bool undertakingNoLien = false;
  bool undertakingNotCollateral = false;
  bool undertakingNoClaims = false;
  bool undertakingNoForeclosure = false;
  bool complianceNoViolation = false;
  bool complianceAllPermits = false;
  bool outstandingFinancialRespNoDebts = false;
  bool outstandingFinancialRespNoHiddenLiabilities = false;
  bool riskManagementFullyInsured = false;
  bool riskManagementDeclaredValue = false;
  bool physicalConditionSound = false;
  bool physicalConditionNoUndisclosedEasements = false;
  bool physicalConditionNolease = false;
  bool hasInsurance = false;

  String trusteeName = "";
  int fundingStructure = 0;
  String? debtInstrumentType;
  String? principalPaymentMethod;
  String? debtInstrumentRepaymentSource;
  String? securityOrCollateralOffered;
  String? debtInstrumentGuaranteesOrEnhancements;
  String debtInstrumentDefaultAndRecoveryTerms = "";
  String? debtInstrumentRepaymentFrequency;
  String? interestRepaymentFrequency;
  String? earlyRedemptionOption;
  String dcsrDetails = "";
  String sinkingFundStructure = "";
  String covenantMonitoringAgent = "";
  String? rightOfRecourse;
  String earlyRedemptionPenalty = "";
  double debtInstrumentInterestRate = 0;
  double dcsrRatio = 0;
  double ltvRatio = 0;
  double interestCoverageRatio = 0;
  double maximumLeverageRatio = 0;
  int tenure = 0;
  int gracePeriod = 0;
  bool trusteeAppointed = false;
  bool reserveFundInPlace = false;

  final valueOfAssetController = TextEditingController();
  final miscCostOfAssetController = TextEditingController();
  final percentageFromPromotersController = TextEditingController();
  final assetOwnerRetainedOrContributedValueController =
      TextEditingController();
  final percentValueOfInsuranceController = TextEditingController();
  final assetDescriptionController = TextEditingController();

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getDebtInstrumentOptions {
    List<DropdownMenuItem<String>> debtInstrumentOptions = [];
    var data = ['Term Loan', 'Bond', 'Promissory Note', 'Debenture', 'Sukuk'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (debtInstrumentType == data[i]) {
        isFound = true;
      }
      debtInstrumentOptions.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    debtInstrumentType = isFound ? debtInstrumentType : null;
    return debtInstrumentOptions;
  }

  List<DropdownMenuItem<String>> get getDebtRepaymentFrequencyOptions {
    List<DropdownMenuItem<String>> debtRepaymentFrequencyOptions = [];
    var data = [
      'Monthly',
      'Quarterly',
      'Semi-Annually',
      'Bullet (at maturity)',
    ];
    bool isFoundInterest = false;
    bool isFoundDebtInstrument = false;
    for (var i = 0; i < data.length; i++) {
      if (interestRepaymentFrequency == data[i]) {
        isFoundInterest = true;
      }

      if (debtInstrumentRepaymentFrequency == data[i]) {
        isFoundDebtInstrument = true;
      }

      debtRepaymentFrequencyOptions.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    interestRepaymentFrequency =
        isFoundInterest ? interestRepaymentFrequency : null;

    debtInstrumentRepaymentFrequency =
        isFoundDebtInstrument ? debtInstrumentRepaymentFrequency : null;
    return debtRepaymentFrequencyOptions;
  }

  List<DropdownMenuItem<String>> get getPrincipalPaymentOptions {
    List<DropdownMenuItem<String>> principalPaymentOptions = [];
    var data = ['Bullet', 'Equal Installments', 'Balloon'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (principalPaymentMethod == data[i]) {
        isFound = true;
      }

      principalPaymentOptions.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    principalPaymentMethod = isFound ? principalPaymentMethod : null;
    return principalPaymentOptions;
  }

  List<DropdownMenuItem<String>> get getRepaymentSourceOptions {
    List<DropdownMenuItem<String>> repaymentSourceOptions = [];
    var data = ['Rental Income', 'Operations Revenue'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (debtInstrumentRepaymentSource == data[i]) {
        isFound = true;
      }
      repaymentSourceOptions.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    debtInstrumentRepaymentSource =
        isFound ? debtInstrumentRepaymentSource : null;

    return repaymentSourceOptions;
  }

  List<DropdownMenuItem<String>> get getSecurityCollateralOptions {
    List<DropdownMenuItem<String>> options = [];
    var data = ['Secured (Real Estate, Equipment)', 'Unsecured'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (securityOrCollateralOffered == data[i]) {
        isFound = true;
      }
      options.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    securityOrCollateralOffered = isFound ? securityOrCollateralOffered : null;
    return options;
  }

  List<DropdownMenuItem<String>> get getGuaranteesOptions {
    List<DropdownMenuItem<String>> guaranteesOptions = [];
    var data = ['None', 'Guarantee', 'DSRA', 'Insurance'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (debtInstrumentGuaranteesOrEnhancements == data[i]) {
        isFound = true;
      }
      guaranteesOptions.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    debtInstrumentGuaranteesOrEnhancements =
        isFound ? debtInstrumentGuaranteesOrEnhancements : null;
    return guaranteesOptions;
  }

  List<DropdownMenuItem<String>> get getYesOrNoOptions {
    List<DropdownMenuItem<String>> options = [];
    var data = ['Yes', 'No'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (earlyRedemptionOption == data[i]) {
        isFound = true;
      }
      options.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    earlyRedemptionOption = isFound ? earlyRedemptionOption : null;
    return options;
  }

  List<DropdownMenuItem<String>> get getRightOfRecourseOptions {
    List<DropdownMenuItem<String>> options = [];
    var data = ['Full Recourse', 'Limited Recourse', 'Non-Recourse'];
    bool isFound = false;
    for (var i = 0; i < data.length; i++) {
      if (rightOfRecourse == data[i]) {
        isFound = true;
      }
      options.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }
    rightOfRecourse = isFound ? rightOfRecourse : null;
    return options;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;
    inspect(data);

    fundingStructure = data['fundingStructure'] ?? 0;
    currentValueOfAsset =
        double.tryParse(data['assetCurrentValue'].toString()) ?? 0;
    assetOwnerRetainedOrContributedValue = double.tryParse(
          data['assetOwnerRetainedOrContributedValue'].toString(),
        ) ??
        0;
    assetMiscCost =
        double.tryParse(data['assetMscCostOutisdeOfValuation'].toString()) ?? 0;
    valueOfTokenizedAsset =
        double.tryParse(data['valueOfTokenizedAsset'].toString()) ?? 0;
    percentageValueOfInsurance =
        double.tryParse(data['percentageValueOfInsurance'].toString()) ?? 0;
    latitude = double.tryParse(data['assetLatitude'].toString()) ?? 0;
    longitude = double.tryParse(data['assetLongitude'].toString()) ?? 0;
    debtInstrumentInterestRate =
        double.tryParse(data['debtInstrumentInterestRate'].toString()) ?? 0;
    dcsrRatio = double.tryParse(data['dcsrRatio'].toString()) ?? 0;
    ltvRatio = double.tryParse(data['ltvRatio'].toString()) ?? 0;
    interestCoverageRatio =
        double.tryParse(data['interestCoverageRatio'].toString()) ?? 0;
    maximumLeverageRatio =
        double.tryParse(data['maximumLeverageRatio'].toString()) ?? 0;
    tenure = int.tryParse(data['tenure'].toString()) ?? 0;
    gracePeriod = int.tryParse(data['gracePeriod'].toString()) ?? 0;
    trusteeAppointed = int.tryParse(data['trusteeAppointed'].toString()) == 1;
    reserveFundInPlace =
        int.tryParse(data['reserveFundInPlace'].toString()) == 1;

    assetOwnership =
        data['ownershipType'] != null && data['ownershipType'].isNotEmpty
            ? data['ownershipType']
            : 'DIRECT';
    thirdPartyOwnerType =
        data['ownershipKind'] != null && data['ownershipKind'].isNotEmpty
            ? data['ownershipKind']
            : 'INDIVIDUAL';
    assetName = data['assetName'] ?? "";
    assetAlreadyExists = data!['assetAlreadyExists'] == 1;
    assetDescription = data['assetDescription'] ?? "";
    assetPhysicalAddress = data['assetPhysicalAddress'] ?? "";
    nameOfOwner = data['assetOwnerName'] ?? "";
    addressOfOwner = data['assetOwnerAddress'] ?? "";
    assetProtectionInPlace = data['protectionMethods'] == null ||
            data['protectionMethods'].toString().isEmpty
        ? []
        : data['protectionMethods'].toString().split(',');
    insuranceCompanyName = data['insuranceCompanyName'] ?? "";
    insurancePolicyNumber = data['insurancePolicyNumber'] ?? "";
    insurancePolicyHolder = data['insurancePolicyHolder'] ?? "";
    trusteeName = data['trusteeName'] ?? "";

    valueOfAssetController.text = currentValueOfAsset == 0
        ? ''
        : formatNumberForInput(currentValueOfAsset);
    miscCostOfAssetController.text =
        assetMiscCost == 0 ? '' : formatNumberForInput(assetMiscCost);
    assetOwnerRetainedOrContributedValueController.text =
        assetOwnerRetainedOrContributedValue == 0
            ? ''
            : formatNumberForInput(assetOwnerRetainedOrContributedValue);
    percentValueOfInsuranceController.text = percentageValueOfInsurance == 0
        ? ''
        : percentageValueOfInsurance.toString();

    otherAssetProtection = data['otherAssetProtection'] ?? "";
    legalAdvisor = data['legalAdvisor'] ?? "";
    financialAdvisor = data['financialAdvisor'] ?? "";

    contractualProtectionRevGuarantees =
        data['contractualProtectionRevGuarantees'] == 1;
    contractualProtectionPerfBond = data['contractualProtectionPerfBond'] == 1;
    contractualProtectionSLA = data['contractualProtectionSLA'] == 1;
    riskSharingMechanismPPPs = data['riskSharingMechanismPPPs'] == 1;
    riskSharingMechanismHedgeInstruments =
        data['riskSharingMechanismHedgeInstruments'] == 1;
    riskSharingMechanismCompletionGuarantees =
        data['riskSharingMechanismCompletionGuarantees'] == 1;
    eSGSafeguardsSusCerts = data['eSGSafeguardsSusCerts'] == 1;
    eSGSafeguardsCommEngPlans = data['eSGSafeguardsCommEngPlans'] == 1;
    securityMeasuresAccessControl = data['securityMeasuresAccessControl'] == 1;
    securityMeasuresSurveilanceSystems =
        data['securityMeasuresSurveilanceSystems'] == 1;
    securityMeasuresOnSiteSecurityPersonnel =
        data['securityMeasuresOnSiteSecurityPersonnel'] == 1;
    securityMeasuresPerimeterSecurity =
        data['securityMeasuresPerimeterSecurity'] == 1;
    securityMeasuresCriticalInfraProtections =
        data['securityMeasuresCriticalInfraProtections'] == 1;
    undertakingNoLien = data['undertakingNoLien'] == 1;
    undertakingNotCollateral = data['undertakingNotCollateral'] == 1;
    undertakingNoClaims = data['undertakingNoClaims'] == 1;
    undertakingNoForeclosure = data['undertakingNoForeclosure'] == 1;
    complianceNoViolation = data['complianceNoViolation'] == 1;
    complianceAllPermits = data['complianceAllPermits'] == 1;
    outstandingFinancialRespNoDebts =
        data['outstandingFinancialRespNoDebts'] == 1;
    outstandingFinancialRespNoHiddenLiabilities =
        data['outstandingFinancialRespNoHiddenLiabilities'] == 1;
    riskManagementFullyInsured = data['riskManagementFullyInsured'] == 1;
    riskManagementDeclaredValue = data['riskManagementDeclaredValue'] == 1;
    physicalConditionSound = data['physicalConditionSound'] == 1;
    physicalConditionNoUndisclosedEasements =
        data['physicalConditionNoUndisclosedEasements'] == 1;
    physicalConditionNolease = data['physicalConditionNolease'] == 1;
    hasInsurance = data['insuranceCompanyName'].toString().isNotEmpty;
    hasLegalAdvisor = data['legalAdvisor'].toString().isNotEmpty;
    hasFinancialAdvisor = data['financialAdvisor'].toString().isNotEmpty;
    hasOtherAssetProtection =
        data['otherAssetProtection'].toString().isNotEmpty;

    debtInstrumentType = data['debtInstrumentType'].toString().nullIfEmpty();
    principalPaymentMethod =
        data['principalPaymentMethod'].toString().nullIfEmpty();
    debtInstrumentRepaymentSource =
        data['debtInstrumentRepaymentSource'].toString().nullIfEmpty();
    securityOrCollateralOffered =
        data['securityOrCollateralOffered'].toString().nullIfEmpty();
    debtInstrumentGuaranteesOrEnhancements =
        data['debtInstrumentGuaranteesOrEnhancements'].toString().nullIfEmpty();
    debtInstrumentDefaultAndRecoveryTerms =
        data['debtInstrumentDefaultAndRecoveryTerms'] ?? "";
    debtInstrumentRepaymentFrequency =
        data['debtInstrumentRepaymentFrequency'].toString().nullIfEmpty();
    interestRepaymentFrequency =
        data['interestRepaymentFrequency'].toString().nullIfEmpty();
    earlyRedemptionOption =
        data['earlyRedemptionOption'].toString().nullIfEmpty();
    dcsrDetails = data['dcsrDetails'] ?? "";
    sinkingFundStructure = data['sinkingFundStructure'] ?? "";
    covenantMonitoringAgent = data['covenantMonitoringAgent'] ?? "";
    rightOfRecourse = data['rightOfRecourse'].toString().nullIfEmpty();
    earlyRedemptionPenalty = data['earlyRedemptionPenalty'] ?? "";

    percentageFromPromoters =
        ((assetOwnerRetainedOrContributedValue / currentValueOfAsset) * 100);
    percentageFromPromotersController.text = percentageFromPromoters.isNaN
        ? '0'
        : formatNumberShort(percentageFromPromoters);

    assetDescriptionController.text = assetDescription;

    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                "assetinformation".tr(),
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    Text(
                      "basicinformation".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetname".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "enterassetname".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: assetName,
                      onChanged: (value) {
                        setState(() {
                          assetName = value;
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          assetName = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetdescription".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: multilineInput(
                      'Asset Description',
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }

                        int wordLength = value.toString().split(' ').length;

                        if (wordLength > assetDescriptionMaxWords) {
                          return "Max words exceeded!";
                        }

                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          assetDescription = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      maxWords: assetDescriptionMaxWords,
                      keyboardtype: TextInputType.multiline,
                      controller: assetDescriptionController,
                      buildCounter: (context,
                          {currentLength, isFocused, maxLength}) {
                        var maxWords = assetDescriptionMaxWords;
                        int length =
                            assetDescriptionController.text.split(' ').length;
                        return Container(
                          child: Text(
                            '$length/$maxWords words',
                            style: TextStyle(
                              color: length > maxWords
                                  ? Colors.red
                                  : notifier.getdarkgrey,
                            ),
                          ),
                        );
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "enterassetphysicaladdress".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "assetphysicaladdress".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      width / 1.12,
                      initialValue: assetPhysicalAddress,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          assetPhysicalAddress = value;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "entergooglemapcords".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    CustomTextFormField.textField(
                      "latitude".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 2.5,
                      initialValue: latitude == 0 ? '' : latitude.toString(),
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          latitude = double.parse(value!.toString());
                        });
                      },
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                    ),
                    CustomTextFormField.textField(
                      "longitude".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 2.5,
                      initialValue: longitude == 0 ? '' : longitude.toString(),
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          longitude = double.parse(value!.toString());
                        });
                      },
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                    ),
                  ],
                ),
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetvalue".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Container(
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "provideassetvalueinfo".tr(),
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetcurrentvalue".tr(args: ['NGN']),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "currentvalueofasset".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      onChanged: (value) {
                        setState(() {
                          if (value.toString().isEmpty) {
                            currentValueOfAsset = 0;
                            valueOfTokenizedAsset = 0;
                            return;
                          }

                          currentValueOfAsset = double.parse(
                            value!.toString().replaceAll(',', ''),
                          );
                          valueOfTokenizedAsset =
                              (assetMiscCost + currentValueOfAsset);
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        currentValueOfAsset = double.parse(value!.toString());
                      },
                      autoFormatNumber: true,
                      isFiat: true,
                      controller: valueOfAssetController,
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "What percentage do you want to retain? (%)",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "How much (%)",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      onChanged: (value) {
                        if (value.toString().isEmpty) {
                          percentageFromPromoters = 0;
                          return;
                        }

                        percentageFromPromoters = double.parse(
                          value!.toString().replaceAll(',', ''),
                        );

                        assetOwnerRetainedOrContributedValue =
                            ((currentValueOfAsset * percentageFromPromoters) /
                                100);
                        assetOwnerRetainedOrContributedValueController.text =
                            truncateToDecimalPlaces(
                          assetOwnerRetainedOrContributedValue,
                          decimalPlaces: 10,
                        );
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        percentageFromPromoters = double.parse(
                          value!.toString(),
                        );
                      },
                      autoFormatNumber: true,
                      isFiat: true,
                      controller: percentageFromPromotersController,
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: width - 60,
                      child: Text(
                        "Value of the Asset retained",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "howmuch".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      onChanged: (value) {
                        setState(() {
                          if (value.toString().isEmpty) {
                            assetOwnerRetainedOrContributedValue = 0;
                            return;
                          }

                          assetOwnerRetainedOrContributedValue = double.parse(
                            value!.toString().replaceAll(',', ''),
                          );
                        });

                        percentageFromPromoters =
                            ((assetOwnerRetainedOrContributedValue /
                                    currentValueOfAsset) *
                                100);
                        percentageFromPromotersController.text =
                            formatNumberShort(percentageFromPromoters);
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        assetOwnerRetainedOrContributedValue = double.parse(
                          value!.toString(),
                        );
                      },
                      autoFormatNumber: true,
                      isFiat: true,
                      controller:
                          assetOwnerRetainedOrContributedValueController,
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                    ),
                  ),
                ],
              ),
              if (assetAlreadyExists) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "assetmisccost".tr(),
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "misccost".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        width / 1.12,
                        onChanged: (value) {
                          setState(() {
                            if (value.toString().isEmpty) {
                              assetMiscCost = 0;
                              return;
                            }

                            assetMiscCost = double.parse(
                              value!.toString().replaceAll(',', ''),
                            );
                            valueOfTokenizedAsset =
                                (assetMiscCost + currentValueOfAsset);
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          assetMiscCost = double.parse(value!.toString());
                        },
                        autoFormatNumber: true,
                        isFiat: true,
                        controller: miscCostOfAssetController,
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetprotectioninplace".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "insurance".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: hasInsurance,
                              label: "comprehensiveinsurance".tr(),
                              onChanged: (value) {
                                setState(() {
                                  hasInsurance = value!;
                                });
                              },
                            ),
                            if (hasInsurance) ...[
                              Row(
                                children: [
                                  SizedBox(width: 15),
                                  Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      SizedBox(height: 10),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "insurancecompanyname".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 70),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child:
                                                CustomTextFormField.textField(
                                              "companyname".tr(),
                                              notifier.getbluecolor,
                                              null,
                                              notifier.getgrey,
                                              null,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              85,
                                              267,
                                              initialValue:
                                                  insuranceCompanyName,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "fieldcannotbeempty"
                                                      .tr();
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  insuranceCompanyName = value!;
                                                });
                                              },
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 70),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "insurancypolicynumber".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child:
                                                CustomTextFormField.textField(
                                              "insurancypolicynumber".tr(),
                                              notifier.getbluecolor,
                                              null,
                                              notifier.getgrey,
                                              null,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              85,
                                              267,
                                              initialValue:
                                                  insurancePolicyNumber,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "fieldcannotbeempty"
                                                      .tr();
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  insurancePolicyNumber =
                                                      value!;
                                                });
                                              },
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "insurancypolicyholder".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child:
                                                CustomTextFormField.textField(
                                              "insurancypolicyholder".tr(),
                                              notifier.getbluecolor,
                                              null,
                                              notifier.getgrey,
                                              null,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              85,
                                              267,
                                              initialValue:
                                                  insurancePolicyHolder,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "fieldcannotbeempty"
                                                      .tr();
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  insurancePolicyHolder =
                                                      value!;
                                                });
                                              },
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "percentagevalueofinsurance".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child:
                                                CustomTextFormField.textField(
                                              "percentagevalueofinsurance".tr(),
                                              notifier.getbluecolor,
                                              null,
                                              notifier.getgrey,
                                              null,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              85,
                                              267,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "Please enter percentage";
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  percentageValueOfInsurance =
                                                      double.parse(value);
                                                });
                                              },
                                              autoFormatNumber: true,
                                              controller:
                                                  percentValueOfInsuranceController,
                                              keyboardtype: TextInputType
                                                  .numberWithOptions(
                                                decimal: true,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "contractualprotections".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: contractualProtectionRevGuarantees,
                              label: "revenueguarantees".tr(),
                              onChanged: (value) {
                                setState(() {
                                  contractualProtectionRevGuarantees = value!;
                                });
                              },
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: contractualProtectionSLA,
                              label: "slas".tr(),
                              onChanged: (value) {
                                setState(() {
                                  contractualProtectionSLA = value!;
                                });
                              },
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "securitymeasures".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: securityMeasuresAccessControl,
                              label: "accesscontrol".tr(),
                              onChanged: (value) {
                                setState(() {
                                  securityMeasuresAccessControl = value!;
                                });
                              },
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: securityMeasuresSurveilanceSystems,
                              label: "surveillancesystems".tr(),
                              onChanged: (value) {
                                setState(() {
                                  securityMeasuresSurveilanceSystems = value!;
                                });
                              },
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: securityMeasuresOnSiteSecurityPersonnel,
                              label: "onsitesecuritypersonnel".tr(),
                              onChanged: (value) {
                                setState(() {
                                  securityMeasuresOnSiteSecurityPersonnel =
                                      value!;
                                });
                              },
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: securityMeasuresPerimeterSecurity,
                              label: "perimetersecurity".tr(),
                              onChanged: (value) {
                                setState(() {
                                  securityMeasuresPerimeterSecurity = value!;
                                });
                              },
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: securityMeasuresCriticalInfraProtections,
                              label: "criticalinfrastructureprotections".tr(),
                              onChanged: (value) {
                                setState(() {
                                  securityMeasuresCriticalInfraProtections =
                                      value!;
                                });
                              },
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "legalfinancialcounsel".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: hasLegalAdvisor,
                              label: "legaladvisor".tr(),
                              onChanged: (value) {
                                setState(() {
                                  hasLegalAdvisor = value!;
                                });
                              },
                            ),
                            if (hasLegalAdvisor) ...[
                              SizedBox(height: 10),
                              Row(
                                children: [
                                  SizedBox(width: 15),
                                  Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      SizedBox(height: 10),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "nameoflegaladvisor".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: multilineInput(
                                              "",
                                              notifier.getbluecolor,
                                              notifier.getgrey,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              100.sp,
                                              250.sp,
                                              initialValue: legalAdvisor,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "fieldcannotbeempty"
                                                      .tr();
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  legalAdvisor = value!;
                                                });
                                              },
                                              minLines: 3,
                                              maxLines: null,
                                              keyboardtype:
                                                  TextInputType.multiline,
                                            ),
                                          ),
                                        ],
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ],
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: hasFinancialAdvisor,
                              label: "financialadvisor".tr(),
                              onChanged: (value) {
                                setState(() {
                                  hasFinancialAdvisor = value!;
                                });
                              },
                            ),
                            if (hasFinancialAdvisor) ...[
                              Row(
                                children: [
                                  SizedBox(width: 15),
                                  Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      SizedBox(height: 10),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "nameoffinancialadvisor".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: multilineInput(
                                              "",
                                              notifier.getbluecolor,
                                              notifier.getgrey,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              100.sp,
                                              250.sp,
                                              initialValue: financialAdvisor,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "fieldcannotbeempty"
                                                      .tr();
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  financialAdvisor = value!;
                                                });
                                              },
                                              minLines: 3,
                                              maxLines: null,
                                              keyboardtype:
                                                  TextInputType.multiline,
                                            ),
                                          ),
                                        ],
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "others".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: hasOtherAssetProtection,
                              label: "other".tr(),
                              onChanged: (value) {
                                setState(() {
                                  hasOtherAssetProtection = value!;
                                });
                              },
                            ),
                            if (hasOtherAssetProtection) ...[
                              Row(
                                children: [
                                  SizedBox(width: 15),
                                  Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      SizedBox(height: 10),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: Text(
                                              "enterotherprotectionsinplace"
                                                  .tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                      SizedBox(height: height / 50),
                                      Row(
                                        children: [
                                          Padding(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0,
                                            ),
                                            child: multilineInput(
                                              "",
                                              notifier.getbluecolor,
                                              notifier.getgrey,
                                              notifier.getblck,
                                              notifier.getgrey,
                                              100.sp,
                                              250.sp,
                                              initialValue:
                                                  otherAssetProtection,
                                              validator: (value) {
                                                if (value.isEmpty) {
                                                  return "fieldcannotbeempty"
                                                      .tr();
                                                }
                                                return null;
                                              },
                                              onSaved: (value) {
                                                setState(() {
                                                  otherAssetProtection = value!;
                                                });
                                              },
                                              minLines: 3,
                                              maxLines: null,
                                              keyboardtype:
                                                  TextInputType.multiline,
                                            ),
                                          ),
                                        ],
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetverificationundertaking".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "ownershipandlegalindependence".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            if (assetAlreadyExists) ...[
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: undertakingNoLien,
                                label: "confirmfreeofloans".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    undertakingNoLien = value!;
                                  });
                                },
                              ),
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: undertakingNotCollateral,
                                label: "confirmfreeofcolateral".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    undertakingNotCollateral = value!;
                                  });
                                },
                                validator: (value) {
                                  if (!undertakingNotCollateral) {
                                    setState(() {
                                      formHasError = true;
                                    });
                                    return '';
                                  }

                                  return null;
                                },
                              ),
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: undertakingNoClaims,
                                label: "confirmfreeofthirdparties".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    undertakingNoClaims = value!;
                                  });
                                },
                                validator: (value) {
                                  if (!undertakingNoClaims) {
                                    setState(() {
                                      formHasError = true;
                                    });
                                    return '';
                                  }

                                  return null;
                                },
                              ),
                            ],
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: undertakingNoForeclosure,
                              label: "confirmfreeoflegaldisputes".tr(),
                              onChanged: (value) {
                                setState(() {
                                  undertakingNoForeclosure = value!;
                                });
                              },
                              validator: (value) {
                                if (!undertakingNoForeclosure) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "regulatorycomplianceandapprovals".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: complianceNoViolation,
                              label: "confirmenvcompliance".tr(),
                              onChanged: (value) {
                                setState(() {
                                  complianceNoViolation = value!;
                                });
                              },
                              validator: (value) {
                                if (!complianceNoViolation) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: complianceAllPermits,
                              label: "confirmhaspermits".tr(),
                              onChanged: (value) {
                                setState(() {
                                  complianceAllPermits = value!;
                                });
                              },
                              validator: (value) {
                                if (!complianceAllPermits) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "outstandingfinancialresponsibilities".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            if (assetAlreadyExists) ...[
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: outstandingFinancialRespNoDebts,
                                label: "confirmnooutstandingpayments".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    outstandingFinancialRespNoDebts = value!;
                                  });
                                },
                              ),
                            ],
                            SizedBox(height: 10),
                            CheckboxItem(
                              value:
                                  outstandingFinancialRespNoHiddenLiabilities,
                              label: "confirmnohiddenliabilities".tr(),
                              onChanged: (value) {
                                setState(() {
                                  outstandingFinancialRespNoHiddenLiabilities =
                                      value!;
                                });
                              },
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "riskmanagementandinsurancecoverage".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: riskManagementFullyInsured,
                              label: "confirmassetinsured".tr(),
                              onChanged: (value) {
                                setState(() {
                                  riskManagementFullyInsured = value!;
                                });
                              },
                            ),
                            if (assetAlreadyExists) ...[
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: riskManagementDeclaredValue,
                                label: "confirmvalueiscurrent".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    riskManagementDeclaredValue = value!;
                                  });
                                },
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "physicalconditionandlegalstatus".tr(),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            CheckboxItem(
                              value: physicalConditionNoUndisclosedEasements,
                              label: "confirmnoundisclosedeasements".tr(),
                              onChanged: (value) {
                                setState(() {
                                  physicalConditionNoUndisclosedEasements =
                                      value!;
                                });
                              },
                              validator: (value) {
                                if (!physicalConditionNoUndisclosedEasements) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            if (assetAlreadyExists) ...[
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: physicalConditionSound,
                                label: "confirmassetstructuralysound".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    physicalConditionSound = value!;
                                  });
                                },
                              ),
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: physicalConditionNolease,
                                label: "confirmnoexistingleaseagreements".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    physicalConditionNolease = value!;
                                  });
                                },
                                validator: (value) {
                                  if (!physicalConditionNolease) {
                                    setState(() {
                                      formHasError = true;
                                    });
                                    return '';
                                  }

                                  return null;
                                },
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              if (fundingStructure != 0) ...[
                SizedBox(height: height / 30),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Debt Instrument Structure",
                        style: TextStyle(
                          fontSize: 18,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Debt Instrument Type",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        debtInstrumentType = value.toString();
                      });
                    },
                    getDebtInstrumentOptions,
                    debtInstrumentType,
                    'Select type',
                    context,
                    null,
                    validator: (value) {
                      if (debtInstrumentType == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Loan Tenure (Total Repayment Period in Months)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue: tenure == 0 ? '' : tenure.toString(),
                        onChanged: (value) {
                          setState(() {
                            tenure = int.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            tenure = int.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Grace Period",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue:
                            gracePeriod == 0 ? '' : gracePeriod.toString(),
                        onChanged: (value) {
                          setState(() {
                            gracePeriod = int.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            gracePeriod = int.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Repayment Frequency",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        debtInstrumentRepaymentFrequency = value.toString();
                      });
                    },
                    getDebtRepaymentFrequencyOptions,
                    debtInstrumentRepaymentFrequency,
                    'Select type',
                    context,
                    null,
                    validator: (value) {
                      if (debtInstrumentRepaymentFrequency == null ||
                          debtInstrumentRepaymentFrequency!.isEmpty) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Interest Rate (Coupon Rate-Fixed or Floating rate)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue: debtInstrumentInterestRate == 0
                            ? ''
                            : debtInstrumentInterestRate.toCleanString(),
                        onChanged: (value) {
                          setState(() {
                            debtInstrumentInterestRate =
                                double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            debtInstrumentInterestRate =
                                double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Interest Repayment Frequency",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        interestRepaymentFrequency = value.toString();
                      });
                    },
                    getDebtRepaymentFrequencyOptions,
                    interestRepaymentFrequency,
                    'Select type',
                    context,
                    null,
                    validator: (value) {
                      if (interestRepaymentFrequency == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Principal Payment Method",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        principalPaymentMethod = value.toString();
                      });
                    },
                    getPrincipalPaymentOptions,
                    principalPaymentMethod,
                    'Select type',
                    context,
                    null,
                    validator: (value) {
                      if (principalPaymentMethod == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Repayment Source",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        debtInstrumentRepaymentSource = value.toString();
                      });
                    },
                    getRepaymentSourceOptions,
                    debtInstrumentRepaymentSource,
                    'Select source',
                    context,
                    null,
                    validator: (value) {
                      if (debtInstrumentRepaymentSource == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Security/Collateral Offered",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        securityOrCollateralOffered = value.toString();
                      });
                    },
                    getSecurityCollateralOptions,
                    securityOrCollateralOffered,
                    'Select option',
                    context,
                    null,
                    validator: (value) {
                      if (securityOrCollateralOffered == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Guarantees/Enhancements",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        debtInstrumentGuaranteesOrEnhancements =
                            value.toString();
                      });
                    },
                    getGuaranteesOptions,
                    debtInstrumentGuaranteesOrEnhancements,
                    'Select type',
                    context,
                    null,
                    validator: (value) {
                      if (debtInstrumentGuaranteesOrEnhancements == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Default & Recovery Terms",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: multilineInput(
                        'What happens in case of missed payment',
                        notifier.getbluecolor,
                        notifier.getgrey,
                        notifier.getblck,
                        notifier.getgrey,
                        100.sp,
                        width / 1.12,
                        initialValue: debtInstrumentDefaultAndRecoveryTerms,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            debtInstrumentDefaultAndRecoveryTerms = value!;
                          });
                        },
                        minLines: 3,
                        maxLines: null,
                        keyboardtype: TextInputType.multiline,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Early Redemption",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        earlyRedemptionOption = value!.toString();
                      });
                    },
                    getYesOrNoOptions,
                    earlyRedemptionOption,
                    'Select type',
                    context,
                    null,
                    validator: (value) {
                      if (earlyRedemptionOption == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
                if (earlyRedemptionOption?.toLowerCase() == 'yes') ...[
                  SizedBox(height: height / 50),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "Early Redemption Penalty",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 50),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          "Enter penalty",
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          85,
                          width / 1.12,
                          initialValue: earlyRedemptionPenalty.toString(),
                          onChanged: (value) {
                            setState(() {
                              earlyRedemptionPenalty = value!;
                            });
                          },
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            earlyRedemptionPenalty = value!;
                          },
                          // autoFormatNumber: true,
                          // isFiat: true,
                          // controller: valueOfAssetController,
                          keyboardtype: TextInputType.numberWithOptions(
                            decimal: true,
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
                SizedBox(height: height / 30),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Container(
                        width: 340,
                        child: Text(
                          "Creditworthiness and Coverage Ratios",
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Container(
                        width: 340,
                        child: Text(
                          "Debt Service Coverage Ratio (DSCR) (Minimum DSCR expected)",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue:
                            dcsrRatio == 0 ? '' : dcsrRatio.toCleanString(),
                        onChanged: (value) {
                          setState(() {
                            dcsrRatio = double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            dcsrRatio = double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Loan-to-Value Ratio (LTV)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue:
                            ltvRatio == 0 ? '' : ltvRatio.toCleanString(),
                        onChanged: (value) {
                          setState(() {
                            ltvRatio = double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            ltvRatio = double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Interest Coverage Ratio (ICR)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue: interestCoverageRatio == 0
                            ? ''
                            : interestCoverageRatio.toCleanString(),
                        onChanged: (value) {
                          setState(() {
                            interestCoverageRatio =
                                double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            interestCoverageRatio =
                                double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Maximum Leverage Ratio if any (optional)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue: maximumLeverageRatio == 0
                            ? ''
                            : maximumLeverageRatio.toCleanString(),
                        onChanged: (value) {
                          setState(() {
                            maximumLeverageRatio =
                                double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            maximumLeverageRatio =
                                double.tryParse(value.toString()) ?? 0;
                          });
                        },
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Container(
                        width: 340,
                        child: Text(
                          "Investor Protections (Specific to Debt)",
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Trustee Appointed?",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: Column(
                    children: [
                      Row(
                        children: [
                          Transform.scale(
                            scale: 1,
                            child: Radio<bool>(
                              value: true,
                              groupValue: trusteeAppointed,
                              activeColor: notifier.getbluewhitecolor,
                              fillColor: WidgetStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor,
                              ),
                              onChanged: (value) => {
                                setState(() {
                                  trusteeAppointed = value!;
                                }),
                              },
                            ),
                          ),
                          Text(
                            "yes".tr(),
                            style: TextStyle(
                              fontSize: 14,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Transform.scale(
                            scale: 1,
                            child: Radio<bool>(
                              value: false,
                              activeColor: notifier.getbluewhitecolor,
                              fillColor: WidgetStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor,
                              ),
                              groupValue: trusteeAppointed,
                              onChanged: (value) => {
                                setState(() {
                                  trusteeAppointed = value!;
                                }),
                              },
                            ),
                          ),
                          Text(
                            "no".tr(),
                            style: TextStyle(
                              fontSize: 14,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Name of bond/debt trustee",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "Enter value",
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        initialValue: trusteeName,
                        onChanged: (value) {
                          setState(() {
                            trusteeName = value;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            trusteeName = value!;
                          });
                        },
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Reserve Fund in Place?",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: Column(
                    children: [
                      Row(
                        children: [
                          Transform.scale(
                            scale: 1,
                            child: Radio<bool>(
                              value: true,
                              groupValue: reserveFundInPlace,
                              activeColor: notifier.getbluewhitecolor,
                              fillColor: WidgetStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor,
                              ),
                              onChanged: (value) => {
                                setState(() {
                                  reserveFundInPlace = value!;
                                }),
                              },
                            ),
                          ),
                          Text(
                            "yes".tr(),
                            style: TextStyle(
                              fontSize: 14,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                      Row(
                        children: [
                          Transform.scale(
                            scale: 1,
                            child: Radio<bool>(
                              value: false,
                              activeColor: notifier.getbluewhitecolor,
                              fillColor: WidgetStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor,
                              ),
                              groupValue: reserveFundInPlace,
                              onChanged: (value) => {
                                setState(() {
                                  reserveFundInPlace = value!;
                                }),
                              },
                            ),
                          ),
                          Text(
                            "no".tr(),
                            style: TextStyle(
                              fontSize: 14,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                if (reserveFundInPlace) ...[
                  SizedBox(height: height / 50),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "Debt service reserve account (DSRA) details",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 70),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: multilineInput(
                          'Give details',
                          notifier.getbluecolor,
                          notifier.getgrey,
                          notifier.getblck,
                          notifier.getgrey,
                          100.sp,
                          width / 1.12,
                          initialValue: dcsrDetails,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {
                              dcsrDetails = value!;
                            });
                          },
                          minLines: 3,
                          maxLines: null,
                          keyboardtype: TextInputType.multiline,
                        ),
                      ),
                    ],
                  ),
                ],
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Container(
                        width: 340,
                        child: Text(
                          "Sinking Fund Structure (Gradual principal repayment details)",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: multilineInput(
                        'Give details',
                        notifier.getbluecolor,
                        notifier.getgrey,
                        notifier.getblck,
                        notifier.getgrey,
                        100.sp,
                        width / 1.12,
                        initialValue: sinkingFundStructure,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            sinkingFundStructure = value!;
                          });
                        },
                        minLines: 3,
                        maxLines: null,
                        keyboardtype: TextInputType.multiline,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Container(
                        width: 330,
                        child: Text(
                          "Covenant Monitoring Agent (Entity overseeing covenant compliance)",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 70),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: multilineInput(
                        'Give details',
                        notifier.getbluecolor,
                        notifier.getgrey,
                        notifier.getblck,
                        notifier.getgrey,
                        100.sp,
                        width / 1.12,
                        initialValue: covenantMonitoringAgent,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            covenantMonitoringAgent = value!;
                          });
                        },
                        minLines: 3,
                        maxLines: null,
                        keyboardtype: TextInputType.multiline,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      textAlign: TextAlign.left,
                      "Right of Recourse",
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 70),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: dropdown(
                    (value) {
                      setState(() {
                        rightOfRecourse = value.toString();
                      });
                    },
                    getRightOfRecourseOptions,
                    rightOfRecourse,
                    'Select recourse',
                    context,
                    null,
                    validator: (value) {
                      if (rightOfRecourse == null) {
                        return "Please choose an option";
                      }
                      return null;
                    },
                  ),
                ),
              ],
              SizedBox(height: height / 30),
              Button(
                "saveandcontinuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  var form = _formKey.currentState;
                  if (form!.validate()) {
                    form.save();
                    submitForm();
                  }
                },
              ),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void submitForm() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credential
      if (!hasInsurance) {
        insuranceCompanyName = "";
        insurancePolicyNumber = "";
        insurancePolicyHolder = "";
        percentageValueOfInsurance = 0;
      }

      if (!hasLegalAdvisor) {
        legalAdvisor = "";
      }

      if (!hasFinancialAdvisor) {
        financialAdvisor = "";
      }

      if (!hasOtherAssetProtection) {
        otherAssetProtection = "";
      }
      var newData = {...data as Map};

      newData['ownershipType'] = assetOwnership;
      newData['ownershipKind'] = thirdPartyOwnerType;
      newData['assetName'] = assetName;
      newData['assetDescription'] = assetDescription;
      newData['assetPhysicalAddress'] = assetPhysicalAddress;
      newData['assetLatitude'] = latitude.toString();
      newData['assetLongitude'] = longitude.toString();
      newData['assetOwnerName'] = nameOfOwner;
      newData['assetOwnerAddress'] = addressOfOwner;
      newData['assetCurrentValue'] = currentValueOfAsset;
      newData['assetMscCostOutisdeOfValuation'] = assetMiscCost;
      newData['assetOwnerRetainedOrContributedValue'] =
          assetOwnerRetainedOrContributedValue;
      newData['valueOfTokenizedAsset'] = valueOfTokenizedAsset;
      newData['protectionMethods'] = assetProtectionInPlace.join(',');
      newData['insuranceCompanyName'] = insuranceCompanyName;
      newData['insurancePolicyNumber'] = insurancePolicyNumber;
      newData['insurancePolicyHolder'] = insurancePolicyHolder;
      newData['percentageValueOfInsurance'] = percentageValueOfInsurance;

      newData['otherAssetProtection'] = otherAssetProtection;
      newData['legalAdvisor'] = legalAdvisor;
      newData['financialAdvisor'] = financialAdvisor;

      newData['contractualProtectionRevGuarantees'] =
          contractualProtectionRevGuarantees ? 1 : 0;
      newData['contractualProtectionPerfBond'] =
          contractualProtectionPerfBond ? 1 : 0;
      newData['contractualProtectionSLA'] = contractualProtectionSLA ? 1 : 0;
      newData['riskSharingMechanismPPPs'] = riskSharingMechanismPPPs ? 1 : 0;
      newData['riskSharingMechanismHedgeInstruments'] =
          riskSharingMechanismHedgeInstruments ? 1 : 0;
      newData['riskSharingMechanismCompletionGuarantees'] =
          riskSharingMechanismCompletionGuarantees ? 1 : 0;
      newData['eSGSafeguardsSusCerts'] = eSGSafeguardsSusCerts ? 1 : 0;
      newData['eSGSafeguardsCommEngPlans'] = eSGSafeguardsCommEngPlans ? 1 : 0;
      newData['securityMeasuresAccessControl'] =
          securityMeasuresAccessControl ? 1 : 0;
      newData['securityMeasuresSurveilanceSystems'] =
          securityMeasuresSurveilanceSystems ? 1 : 0;
      newData['securityMeasuresOnSiteSecurityPersonnel'] =
          securityMeasuresOnSiteSecurityPersonnel ? 1 : 0;
      newData['securityMeasuresPerimeterSecurity'] =
          securityMeasuresPerimeterSecurity ? 1 : 0;
      newData['securityMeasuresCriticalInfraProtections'] =
          securityMeasuresCriticalInfraProtections ? 1 : 0;
      newData['undertakingNoLien'] = undertakingNoLien ? 1 : 0;
      newData['undertakingNotCollateral'] = undertakingNotCollateral ? 1 : 0;
      newData['undertakingNoClaims'] = undertakingNoClaims ? 1 : 0;
      newData['undertakingNoForeclosure'] = undertakingNoForeclosure ? 1 : 0;
      newData['complianceNoViolation'] = complianceNoViolation ? 1 : 0;
      newData['complianceAllPermits'] = complianceAllPermits ? 1 : 0;
      newData['outstandingFinancialRespNoDebts'] =
          outstandingFinancialRespNoDebts ? 1 : 0;
      newData['outstandingFinancialRespNoHiddenLiabilities'] =
          outstandingFinancialRespNoHiddenLiabilities ? 1 : 0;
      newData['riskManagementFullyInsured'] =
          riskManagementFullyInsured ? 1 : 0;
      newData['riskManagementDeclaredValue'] =
          riskManagementDeclaredValue ? 1 : 0;
      newData['physicalConditionSound'] = physicalConditionSound ? 1 : 0;
      newData['physicalConditionNoUndisclosedEasements'] =
          physicalConditionNoUndisclosedEasements ? 1 : 0;
      newData['physicalConditionNolease'] = physicalConditionNolease ? 1 : 0;
      newData['hasInsurance'] = hasInsurance ? 1 : 0;

      newData['trusteeName'] = trusteeName;
      newData['debtInstrumentType'] = debtInstrumentType;
      newData['principalPaymentMethod'] = principalPaymentMethod;
      newData['debtInstrumentRepaymentSource'] = debtInstrumentRepaymentSource;
      newData['debtInstrumentGuaranteesOrEnhancements'] =
          debtInstrumentGuaranteesOrEnhancements;
      newData['debtInstrumentDefaultAndRecoveryTerms'] =
          debtInstrumentDefaultAndRecoveryTerms;
      newData['debtInstrumentRepaymentFrequency'] =
          debtInstrumentRepaymentFrequency;
      newData['interestRepaymentFrequency'] = interestRepaymentFrequency;
      newData['earlyRedemptionOption'] = earlyRedemptionOption;
      newData['earlyRedemptionPenalty'] = earlyRedemptionPenalty;
      newData['securityOrCollateralOffered'] = securityOrCollateralOffered;
      newData['dcsrDetails'] = dcsrDetails;
      newData['sinkingFundStructure'] = sinkingFundStructure;
      newData['covenantMonitoringAgent'] = covenantMonitoringAgent;
      newData['rightOfRecourse'] = rightOfRecourse;
      newData['debtInstrumentInterestRate'] = debtInstrumentInterestRate;
      newData['dcsrRatio'] = dcsrRatio;
      newData['ltvRatio'] = ltvRatio;
      newData['interestCoverageRatio'] = interestCoverageRatio;
      newData['maximumLeverageRatio'] = maximumLeverageRatio;
      newData['tenure'] = tenure;
      newData['gracePeriod'] = gracePeriod;
      newData['trusteeAppointed'] = trusteeAppointed ? 1 : 0;
      newData['reserveFundInPlace'] = reserveFundInPlace ? 1 : 0;

      inspect(newData);

      String requestBody = jsonEncode(newData);
      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
        Navigator.of(context).pop();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
      hideLoader(context);
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Future<void> refreshCurrentTokenizationInfo() async {
    try {
      var uri = '/v1/tokenization/detail/${appState.viewData!['id']}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      if (responseData['statusCode'] == 200) {
        appState.viewData = responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Widget CheckItem(
    String name,
    void Function()? onClick, {
    required Color backColor,
    required Color foreColor,
    required Color borderColor,
    double? fontSize = 15,
  }) {
    return Padding(
      padding: const EdgeInsets.all(3.0),
      child: ElevatedButton(
        onPressed: onClick,
        style: ButtonStyle(
          overlayColor: WidgetStateProperty.all<Color>(notifier.getsplashgrey),
          elevation: WidgetStateProperty.all<double>(0),
          backgroundColor: WidgetStateProperty.all<Color>(backColor),
          foregroundColor: WidgetStateProperty.all<Color>(
            notifier.getwihitecolor,
          ),
          side: WidgetStateProperty.all(
            BorderSide(color: borderColor, width: 1, style: BorderStyle.solid),
          ),
          shape: WidgetStateProperty.all<RoundedRectangleBorder>(
            const RoundedRectangleBorder(
              borderRadius: BorderRadius.all(Radius.circular(10)),
            ),
          ),
        ),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                color: foreColor,
                fontFamily: fontbody,
                fontSize: fontSize,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget CheckboxItem({
    required String label,
    required bool value,
    required void Function(bool?) onChanged,
    String? Function(Object?)? validator,
  }) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Transform.scale(
          scale: 1,
          child: FormField(
            builder: (state) {
              return SizedBox(
                width: 24,
                height: 24,
                child: Checkbox(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(Radius.circular(5.sp)),
                  ),
                  activeColor: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor90,
                  side: BorderSide(
                    color: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                  ),
                  value: value,
                  onChanged: onChanged,
                ),
              );
            },
            validator: validator,
          ),
        ),
        SizedBox(width: 10),
        Container(
          width: width / 1.4,
          child: Text(
            label,
            overflow: TextOverflow.visible,
            style: TextStyle(
              fontSize: 15,
              color: formHasError && !value && validator != null
                  ? Colors.red
                  : notifier.getbluewhitecolor,
              fontFamily: fontbody,
            ),
          ),
        ),
      ],
    );
  }
}
