import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
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

class UpcomingAssetInformationView extends StatefulWidget {
  const UpcomingAssetInformationView({Key? key}) : super(key: key);

  @override
  State<UpcomingAssetInformationView> createState() =>
      _UpcomingAssetInformationView();
}

class _UpcomingAssetInformationView extends State<UpcomingAssetInformationView>
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
  bool assetAlreadyExists = false;
  bool formHasError = false;
  late dynamic data = {};

  String projectStrategicObjectives = "";
  String projectDevelopmentTimeline = "";
  String projectKeyMilestoneAndDates = "";
  String projectScope = "";
  String projectEconomicBenefits = "";
  int projectExpectedNoOfJobs = 0;
  String projectIntendedSocialBenefits = "";
  String projectTechnicalPartners = "";
  String projectFinancialPartners = "";

  double percentageFromPromoters = 0;
  double estimatedProjectIRR = 0;
  double estimatedProjectROI = 0;
  double estimatedProjectNPV = 0;
  int estimatedProjectPaybackPeriodsInMonths = 0;
  String keyAssumptionsList = "";
  String projectIdentifiedLegalRisks = "";
  String projectIdentifiedRegulatoryRisks = "";
  String projectIdentifiedOperationalOrExecutionRisks = "";
  String projectIdentifiedMarketRisks = "";
  String projectIdentifiedOtherRelevantRisks = "";

  String independentMonitoringList = "";
  String otherAssetProtection = "";
  String legalAdvisor = "";
  String financialAdvisor = "";

  bool hasIndependentMonitoring = false;
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

  final valueOfAssetController = TextEditingController();
  final miscCostOfAssetController = TextEditingController();
  final percentageFromPromotersController = TextEditingController();
  final assetOwnerRetainedOrContributedValueController =
      TextEditingController();
  final percentValueOfInsuranceController = TextEditingController();

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getAssetProtectionOptions {
    List<DropdownMenuItem<String>> assetProtectionOptions = [];
    var data = appState.tokenizationData['assetProtectionOptions'];
    for (var i = 0; i < data.length; i++) {
      assetProtectionOptions.add(
        DropdownMenuItem(
          child: Text(data![i]['id'], overflow: TextOverflow.ellipsis),
          value: data![i]['id'],
        ),
      );
    }
    return assetProtectionOptions;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;
    inspect(data);
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
    latitude = double.tryParse(data['assetLatitude'].toString()) ?? 0;
    longitude = double.tryParse(data['assetLongitude'].toString()) ?? 0;
    nameOfOwner = data['assetOwnerName'] ?? "";
    addressOfOwner = data['assetOwnerAddress'] ?? "";
    currentValueOfAsset =
        double.tryParse(data['assetCurrentValue'].toString()) ?? 0;
    assetOwnerRetainedOrContributedValue =
        double.tryParse(
          data['assetOwnerRetainedOrContributedValue'].toString(),
        ) ??
        0;
    assetMiscCost =
        double.tryParse(data['assetMscCostOutisdeOfValuation'].toString()) ?? 0;
    valueOfTokenizedAsset =
        double.tryParse(data['valueOfTokenizedAsset'].toString()) ?? 0;
    assetProtectionInPlace =
        data['protectionMethods'] == null ||
            data['protectionMethods'].toString().isEmpty
        ? []
        : data['protectionMethods'].toString().split(',');
    insuranceCompanyName = data['insuranceCompanyName'] ?? "";
    insurancePolicyNumber = data['insurancePolicyNumber'] ?? "";
    insurancePolicyHolder = data['insurancePolicyHolder'] ?? "";
    percentageValueOfInsurance =
        double.tryParse(data['percentageValueOfInsurance'].toString()) ?? 0;

    valueOfAssetController.text = currentValueOfAsset == 0
        ? ''
        : formatNumberForInput(currentValueOfAsset);
    miscCostOfAssetController.text = assetMiscCost == 0
        ? ''
        : formatNumberForInput(assetMiscCost);
    assetOwnerRetainedOrContributedValueController.text =
        assetOwnerRetainedOrContributedValue == 0
        ? ''
        : formatNumberForInput(assetOwnerRetainedOrContributedValue);
    percentValueOfInsuranceController.text = percentageValueOfInsurance == 0
        ? ''
        : percentageValueOfInsurance.toString();

    independentMonitoringList = data['independentMonitoringList'] ?? "";
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
    hasIndependentMonitoring = data['independentMonitoringList']
        .toString()
        .isNotEmpty;
    hasLegalAdvisor = data['legalAdvisor'].toString().isNotEmpty;
    hasFinancialAdvisor = data['financialAdvisor'].toString().isNotEmpty;
    hasOtherAssetProtection = data['otherAssetProtection']
        .toString()
        .isNotEmpty;

    projectStrategicObjectives = data['projectStrategicObjectives'] ?? "";
    projectDevelopmentTimeline = data['projectDevelopmentTimeline'] ?? "";
    projectKeyMilestoneAndDates = data['projectKeyMilestoneAndDates'] ?? "";
    projectScope = data['projectScope'] ?? "";
    projectEconomicBenefits = data['projectEconomicBenefits'] ?? "";
    projectExpectedNoOfJobs = data['projectExpectedNoOfJobs'] ?? 0;
    projectIntendedSocialBenefits = data['projectIntendedSocialBenefits'] ?? "";
    projectTechnicalPartners = data['projectTechnicalPartners'] ?? "";
    projectFinancialPartners = data['projectFinancialPartners'] ?? "";

    estimatedProjectPaybackPeriodsInMonths =
        data['estimatedProjectPaybackPeriodsInMonths'];
    keyAssumptionsList = data['keyAssumptionsList'] ?? "";
    projectIdentifiedLegalRisks = data['projectIdentifiedLegalRisks'] ?? "";
    projectIdentifiedRegulatoryRisks =
        data['projectIdentifiedRegulatoryRisks'] ?? "";
    projectIdentifiedOperationalOrExecutionRisks =
        data['projectIdentifiedOperationalOrExecutionRisks'] ?? "";
    projectIdentifiedMarketRisks = data['projectIdentifiedMarketRisks'] ?? "";
    projectIdentifiedOtherRelevantRisks =
        data['projectIdentifiedOtherRelevantRisks'] ?? "";

    estimatedProjectIRR =
        double.tryParse(data['estimatedProjectIRR'].toString()) ?? 0;
    estimatedProjectROI =
        double.tryParse(data['estimatedProjectROI'].toString()) ?? 0;
    estimatedProjectNPV =
        double.tryParse(data['estimatedProjectNPV'].toString()) ?? 0;

    percentageFromPromoters =
        ((assetOwnerRetainedOrContributedValue / currentValueOfAsset) * 100);
    percentageFromPromotersController.text = percentageFromPromoters.isNaN
        ? '0'
        : formatNumberShort(percentageFromPromoters);

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
                      initialValue: assetDescription,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
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
                    child: Text(
                      "Project Strategic Objectives",
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
                      "Enter objectives",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectStrategicObjectives,
                      onChanged: (value) {
                        setState(() {
                          projectStrategicObjectives = value;
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
                          projectStrategicObjectives = value!;
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
                    child: Text(
                      "Project Scope",
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
                      "Enter project scope",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectScope,
                      onChanged: (value) {
                        setState(() {
                          projectScope = value;
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
                          projectScope = value!;
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
                    child: Text(
                      "Upload Asset Images",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  GestureDetector(
                    onTap: () {
                      getFile();
                    },
                    child: Column(
                      children: [
                        SizedBox(height: height / 50),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                          child: Container(
                            decoration: BoxDecoration(
                              border: Border.all(
                                color: notifier.getbluewhitecolor,
                                width: 1,
                              ),
                              borderRadius: const BorderRadius.all(
                                Radius.circular(15.0),
                              ),
                              color: notifier.isDark
                                  ? darktilewhitecolor
                                  : notifier.getaddsubwalletgrey,
                            ),
                            child: Column(
                              children: [
                                SizedBox(
                                  width: width / 1.2,
                                  height: height / 6,
                                  child: Center(
                                    child: Wrap(
                                      alignment: WrapAlignment.center,
                                      children: [
                                        Text(
                                          "browsefiles".tr(),
                                          textAlign: TextAlign.center,
                                          style: TextStyle(
                                            color: notifier.getbluewhitecolor,
                                            fontFamily: fontsemibold,
                                            fontSize: 12.sp,
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 2),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              if (true) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'Please upload asset images',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: Colors.red,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              if (true) ...[
                GestureDetector(
                  onTap: () {
                    getFile();
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Image.network(
                      'https://picsum.photos/200/300',
                      width: width / 1.3,
                      height: height / 6,
                    ),
                  ),
                ),
              ],
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15),
                child: Container(
                  width: width,
                  child: Text(
                    textAlign: TextAlign.left,
                    "Ownership",
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select asset ownership',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
                    return null;
                  },
                ),
              ),
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Project Timeline",
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
                        "Please provide the following information about the project's relevant dates"
                            .tr(),
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
                      "Project Development Timeline",
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
                      "Time to go live in months",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          projectDevelopmentTimeline = value;
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
                          projectDevelopmentTimeline = value!;
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
                      "Project Estimated Payback Period (in Months)",
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
                      "Enter estimate",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: estimatedProjectPaybackPeriodsInMonths
                          .toString(),
                      onChanged: (value) {
                        setState(() {
                          estimatedProjectPaybackPeriodsInMonths =
                              int.tryParse(value) ?? 0;
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
                          estimatedProjectPaybackPeriodsInMonths =
                              int.tryParse(value) ?? 0;
                        });
                      },
                      keyboardtype: TextInputType.number,
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      "5 Key Milestones & Dates",
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    TextButton(
                      onPressed: addMilestone,
                      style: ButtonStyle(
                        padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                        minimumSize: WidgetStatePropertyAll(Size.zero),
                      ),
                      child: Icon(Icons.add_circle, size: 20),
                    ),
                  ],
                ),
              ),
              if (true) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'Atleast 5 milestones',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: Colors.red,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(10),
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
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  "Land acquisition",
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                                TextButton(
                                  onPressed: () {},
                                  style: ButtonStyle(
                                    padding: WidgetStatePropertyAll(
                                      EdgeInsets.all(7),
                                    ),
                                    tapTargetSize:
                                        MaterialTapTargetSize.shrinkWrap,
                                    minimumSize: WidgetStatePropertyAll(
                                      Size.zero,
                                    ),
                                  ),
                                  child: Icon(
                                    CupertinoIcons.trash,
                                    size: 15,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  "January, 2025",
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
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
                      "Project Benefits",
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
                        "Please provide the following information about the project's benefits"
                            .tr(),
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
                      "Project Economic Benefits",
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
                      "Enter intended economic benefits",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectEconomicBenefits,
                      onChanged: (value) {
                        setState(() {
                          projectEconomicBenefits = value;
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
                          projectEconomicBenefits = value!;
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
                    child: Text(
                      "Expected No. of Job to be Created",
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
                      "Enter number of jobs",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: projectExpectedNoOfJobs.toString(),
                      onChanged: (value) {
                        setState(() {
                          projectExpectedNoOfJobs = int.tryParse(value) ?? 0;
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
                          projectExpectedNoOfJobs = int.tryParse(value) ?? 0;
                        });
                      },
                      keyboardtype: TextInputType.number,
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
                      "Project Intended Social Benefits",
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
                      "Enter intended social benefits",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectIntendedSocialBenefits,
                      onChanged: (value) {
                        setState(() {
                          projectIntendedSocialBenefits = value;
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
                          projectIntendedSocialBenefits = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                      "Partners",
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
                        "Please provide the following information about the project's partners if any"
                            .tr(),
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
                      "Technical Partners",
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
                      "Enter name of technical partners",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectTechnicalPartners,
                      onChanged: (value) {
                        setState(() {
                          projectTechnicalPartners = value;
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
                          projectTechnicalPartners = value!;
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
                    child: Text(
                      "Financial Partners",
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
                      "Enter details of financial partners",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectFinancialPartners,
                      onChanged: (value) {
                        setState(() {
                          projectFinancialPartners = value;
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
                          projectFinancialPartners = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                      "Asset Location",
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
                        "Please provide the following information about the asset's location"
                            .tr(),
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
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Asset Value & Financing",
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
                        "Please provide the following information about the asset value"
                            .tr(),
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
                      "What is the Total Estimated Project Budget",
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
                      'How Much Equity is Contributed by Promoters (%)',
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
                        "Value of the Equity Contributed by the Promoter(s)",
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Asset Financial Performance",
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
                        "Provide projections of your project’s future financial performance to help assess the viability and potential growth trajectory of your project",
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
                      "Estimated Project IRR (in %)",
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
                      "Enter estimate",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      initialValue: estimatedProjectIRR.toString(),
                      onChanged: (value) {
                        setState(() {
                          estimatedProjectIRR = double.parse(value!.toString());
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        estimatedProjectIRR = double.parse(value!.toString());
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Estimated Project ROI in 5 years (in %)",
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
                      "Enter estimate",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      initialValue: estimatedProjectROI.toString(),
                      onChanged: (value) {
                        setState(() {
                          estimatedProjectROI = double.parse(value!.toString());
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        estimatedProjectROI = double.parse(value!.toString());
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Estimated Project NPV at Launch (Day 1)",
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
                      "Enter estimate",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      initialValue: estimatedProjectNPV.toString(),
                      onChanged: (value) {
                        setState(() {
                          estimatedProjectNPV = double.parse(value!.toString());
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        estimatedProjectNPV = double.parse(value!.toString());
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Estimated Project Payback Periods?",
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
                      "Enter estimate",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      initialValue: estimatedProjectROI.toString(),
                      onChanged: (value) {
                        setState(() {
                          estimatedProjectROI = double.parse(value!.toString());
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        estimatedProjectROI = double.parse(value!.toString());
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 255,
                      child: Text(
                        "List all Key Assumptions Including Values Assumed (If any)",
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
                    child: multilineInput(
                      "List all key assumptions",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: keyAssumptionsList,
                      onChanged: (value) {
                        setState(() {
                          keyAssumptionsList = value.toString();
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        keyAssumptionsList = value.toString();
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
                    child: Text(
                      "Project Risk Assessment",
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
                        "Outline the types of risks you anticipate, their likelihood and potential impact, This assessment will help us understand your risk management approach and consider how these factors may influence project outcomes.",
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
                      "Legal Risks Identified (Provide Details)",
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
                    child: multilineInput(
                      "Enter legal risks identified ",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectIdentifiedLegalRisks,
                      onChanged: (value) {
                        setState(() {
                          projectIdentifiedLegalRisks = value!.toString();
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        projectIdentifiedLegalRisks = value!.toString();
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
                    child: Text(
                      "Regulatory Risks Identified (Provide Details)",
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
                    child: multilineInput(
                      "Enter regulatory risks identified ",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectIdentifiedRegulatoryRisks,
                      onChanged: (value) {
                        setState(() {
                          projectIdentifiedRegulatoryRisks = value!.toString();
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        projectIdentifiedRegulatoryRisks = value!.toString();
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
                    child: SizedBox(
                      width: 255,
                      child: Text(
                        "Operational/Execution Risks Identified (Provide Details)",
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
                    child: multilineInput(
                      "Enter operational/execution risks identified ",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue:
                          projectIdentifiedOperationalOrExecutionRisks,
                      onChanged: (value) {
                        setState(() {
                          projectIdentifiedOperationalOrExecutionRisks = value!
                              .toString();
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        projectIdentifiedOperationalOrExecutionRisks = value!
                            .toString();
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
                    child: Text(
                      "Market Risks Identified (Provide Details)",
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
                    child: multilineInput(
                      "Enter market risks identified ",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectIdentifiedMarketRisks,
                      onChanged: (value) {
                        setState(() {
                          projectIdentifiedMarketRisks = value!.toString();
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        projectIdentifiedMarketRisks = value!.toString();
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
                    child: Text(
                      "Other Relevant Risks Identified (Provide Details)",
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
                    child: multilineInput(
                      "Enter other relevant risks identified ",
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: projectIdentifiedOtherRelevantRisks,
                      onChanged: (value) {
                        setState(() {
                          projectIdentifiedOtherRelevantRisks = value!
                              .toString();
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        projectIdentifiedOtherRelevantRisks = value!.toString();
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                      "Asset Protection In Place",
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
                        "Check each asset protection that has been put in place for the asset and enter the information needed for each protection"
                            .tr(),
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
                                            child: CustomTextFormField.textField(
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
                                            child: CustomTextFormField.textField(
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
                                            child: CustomTextFormField.textField(
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
                                            child: CustomTextFormField.textField(
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
                                              keyboardtype:
                                                  TextInputType.numberWithOptions(
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
                            if (!assetAlreadyExists) ...[
                              CheckboxItem(
                                value: contractualProtectionPerfBond,
                                label: "performancebond".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    contractualProtectionPerfBond = value!;
                                  });
                                },
                              ),
                            ],
                            CheckboxItem(
                              value: contractualProtectionRevGuarantees,
                              label: "revenueguarantees".tr(),
                              onChanged: (value) {
                                setState(() {
                                  contractualProtectionRevGuarantees = value!;
                                });
                              },
                            ),
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
              if (!assetAlreadyExists) ...[
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
                                    "risksharingmechanisms".tr(),
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
                                value: riskSharingMechanismPPPs,
                                label: "ppps".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    riskSharingMechanismPPPs = value!;
                                  });
                                },
                              ),
                              CheckboxItem(
                                value: riskSharingMechanismHedgeInstruments,
                                label: "hedginginstruments".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    riskSharingMechanismHedgeInstruments =
                                        value!;
                                  });
                                },
                              ),
                              CheckboxItem(
                                value: riskSharingMechanismCompletionGuarantees,
                                label: "completionguarantees".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    riskSharingMechanismCompletionGuarantees =
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
                                    "governanceandoversight".tr(),
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
                                value: hasIndependentMonitoring,
                                label: "independentmonitoring".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    hasIndependentMonitoring = value!;
                                  });
                                },
                              ),
                              if (hasIndependentMonitoring) ...[
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
                                              padding:
                                                  const EdgeInsets.symmetric(
                                                    horizontal: 20.0,
                                                  ),
                                              child: Container(
                                                width: 260,
                                                child: Text(
                                                  "listmonitoringoperators"
                                                      .tr(),
                                                  style: TextStyle(
                                                    fontSize: 12,
                                                    fontFamily: fontsemibold,
                                                    color: notifier
                                                        .getbluewhitecolor,
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
                                              padding:
                                                  const EdgeInsets.symmetric(
                                                    horizontal: 20.0,
                                                  ),
                                              child: multilineInput(
                                                "List details of independent monitors",
                                                notifier.getbluecolor,
                                                notifier.getgrey,
                                                notifier.getblck,
                                                notifier.getgrey,
                                                100.sp,
                                                250.sp,
                                                initialValue:
                                                    independentMonitoringList,
                                                validator: (value) {
                                                  if (hasIndependentMonitoring &&
                                                      value.isEmpty) {
                                                    return "fieldcannotbeempty"
                                                        .tr();
                                                  }
                                                  return null;
                                                },
                                                onSaved: (value) {
                                                  setState(() {
                                                    independentMonitoringList =
                                                        value!;
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
                                  Container(
                                    width: width / 1.3,
                                    child: Text(
                                      "esgs".tr(),
                                      style: TextStyle(
                                        fontSize: 12,
                                        fontFamily: fontsemibold,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: eSGSafeguardsSusCerts,
                                label: "sustainabilitycertifications".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    eSGSafeguardsSusCerts = value!;
                                  });
                                },
                              ),
                              SizedBox(height: 10),
                              CheckboxItem(
                                value: eSGSafeguardsCommEngPlans,
                                label: "communityengagementplans".tr(),
                                onChanged: (value) {
                                  setState(() {
                                    eSGSafeguardsCommEngPlans = value!;
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
              ],
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
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
                      "Loan Tenor (Total Repayment Period in Months)",
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
                        });
                      },
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
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
                      "Repayment Schedule",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  GestureDetector(
                    onTap: () {
                      getFile();
                    },
                    child: Column(
                      children: [
                        SizedBox(height: height / 50),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                          child: Container(
                            decoration: BoxDecoration(
                              border: Border.all(
                                color: notifier.getbluewhitecolor,
                                width: 1,
                              ),
                              borderRadius: const BorderRadius.all(
                                Radius.circular(15.0),
                              ),
                              color: notifier.isDark
                                  ? darktilewhitecolor
                                  : notifier.getaddsubwalletgrey,
                            ),
                            child: Column(
                              children: [
                                SizedBox(
                                  width: width / 1.2,
                                  height: height / 6,
                                  child: Center(
                                    child: Wrap(
                                      alignment: WrapAlignment.center,
                                      children: [
                                        Text(
                                          "browsefiles".tr(),
                                          textAlign: TextAlign.center,
                                          style: TextStyle(
                                            color: notifier.getbluewhitecolor,
                                            fontFamily: fontsemibold,
                                            fontSize: 12.sp,
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                ),
                                const SizedBox(height: 2),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              if (true) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'Please upload file',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: Colors.red,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              if (true) ...[
                GestureDetector(
                  onTap: () {
                    getFile();
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Image.network(
                      'https://picsum.photos/200/200',
                      width: width / 1.3,
                      height: height / 6,
                    ),
                  ),
                ),
              ],
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
                        });
                      },
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: multilineInput(
                      'What cash flows will be used to repay debt (e.g. rental income, operation revenue)',
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: assetDescription,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
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
                      initialValue: assetDescription,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
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
                      keyboardtype: TextInputType.multiline,
                    ),
                  ),
                ],
              ),
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
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
                            groupValue: true,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) => {
                              // setState(() {
                              //   hasAllRequiredManagerDocuments = value!;
                              // }),
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
                            groupValue: false,
                            onChanged: (value) => {
                              // setState(() {
                              //   hasAllRequiredManagerDocuments = value!;
                              // }),
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
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      // initialValue: projectDevelopmentTimeline,
                      onChanged: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value;
                        });
                      },
                      validator: (value) {
                        // if (value.isEmpty) {
                        //   return "fieldcannotbeempty".tr();
                        // }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          // projectDevelopmentTimeline = value!;
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
                            groupValue: true,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) => {
                              // setState(() {
                              //   hasAllRequiredManagerDocuments = value!;
                              // }),
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
                            groupValue: false,
                            onChanged: (value) => {
                              // setState(() {
                              //   hasAllRequiredManagerDocuments = value!;
                              // }),
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
                      initialValue: assetDescription,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
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
                      initialValue: assetDescription,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
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
                      initialValue: assetDescription,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
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
                    setState(() {});
                  },
                  [],
                  null,
                  'Select asset ownership',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
                    // }
                    return null;
                  },
                ),
              ),
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

      if (!hasIndependentMonitoring) {
        independentMonitoringList = "";
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

      newData['independentMonitoringList'] = independentMonitoringList;
      newData['otherAssetProtection'] = otherAssetProtection;
      newData['legalAdvisor'] = legalAdvisor;
      newData['financialAdvisor'] = financialAdvisor;

      newData['contractualProtectionRevGuarantees'] =
          contractualProtectionRevGuarantees ? 1 : 0;
      newData['contractualProtectionPerfBond'] = contractualProtectionPerfBond
          ? 1
          : 0;
      newData['contractualProtectionSLA'] = contractualProtectionSLA ? 1 : 0;
      newData['riskSharingMechanismPPPs'] = riskSharingMechanismPPPs ? 1 : 0;
      newData['riskSharingMechanismHedgeInstruments'] =
          riskSharingMechanismHedgeInstruments ? 1 : 0;
      newData['riskSharingMechanismCompletionGuarantees'] =
          riskSharingMechanismCompletionGuarantees ? 1 : 0;
      newData['eSGSafeguardsSusCerts'] = eSGSafeguardsSusCerts ? 1 : 0;
      newData['eSGSafeguardsCommEngPlans'] = eSGSafeguardsCommEngPlans ? 1 : 0;
      newData['securityMeasuresAccessControl'] = securityMeasuresAccessControl
          ? 1
          : 0;
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
      newData['riskManagementFullyInsured'] = riskManagementFullyInsured
          ? 1
          : 0;
      newData['riskManagementDeclaredValue'] = riskManagementDeclaredValue
          ? 1
          : 0;
      newData['physicalConditionSound'] = physicalConditionSound ? 1 : 0;
      newData['physicalConditionNoUndisclosedEasements'] =
          physicalConditionNoUndisclosedEasements ? 1 : 0;
      newData['physicalConditionNolease'] = physicalConditionNolease ? 1 : 0;
      newData['hasInsurance'] = hasInsurance ? 1 : 0;

      newData['projectStrategicObjectives'] = projectStrategicObjectives;
      newData['projectDevelopmentTimeline'] = projectDevelopmentTimeline;
      newData['projectKeyMilestoneAndDates'] = projectKeyMilestoneAndDates;
      newData['projectScope'] = projectScope;
      newData['projectEconomicBenefits'] = projectEconomicBenefits;
      newData['projectExpectedNoOfJobs'] = projectExpectedNoOfJobs;
      newData['projectIntendedSocialBenefits'] = projectIntendedSocialBenefits;
      newData['projectTechnicalPartners'] = projectTechnicalPartners;
      newData['projectFinancialPartners'] = projectFinancialPartners;

      newData['estimatedProjectPaybackPeriodsInMonths'] =
          estimatedProjectPaybackPeriodsInMonths;
      newData['keyAssumptionsList'] = keyAssumptionsList;
      newData['projectIdentifiedLegalRisks'] = projectIdentifiedLegalRisks;
      newData['projectIdentifiedRegulatoryRisks'] =
          projectIdentifiedRegulatoryRisks;
      newData['projectIdentifiedOperationalOrExecutionRisks'] =
          projectIdentifiedOperationalOrExecutionRisks;
      newData['projectIdentifiedMarketRisks'] = projectIdentifiedMarketRisks;
      newData['projectIdentifiedOtherRelevantRisks'] =
          projectIdentifiedOtherRelevantRisks;

      newData['estimatedProjectIRR'] = estimatedProjectIRR;
      newData['estimatedProjectROI'] = estimatedProjectROI;
      newData['estimatedProjectNPV'] = estimatedProjectNPV;

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

  void addMilestone() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: notifier.getwihitecolor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
          top: Radius.circular(16),
        ), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.75,
          minChildSize: 0.25,
          expand: false,
          builder: (context, scrollController) {
            return SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  SizedBox(height: 10),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 14.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          "Add Milestone",
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(Icons.cancel_outlined),
                          style: ButtonStyle(
                            padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            minimumSize: WidgetStatePropertyAll(Size.zero),
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: 20),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "Milestone",
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
                          "Enter milestone",
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          // initialValue: assetPhysicalAddress,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {});
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
                          "Select Date",
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
                          "",
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          // initialValue: assetPhysicalAddress,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {});
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 60),
                ],
              ),
            );
          },
        );
      },
    );
  }
}
