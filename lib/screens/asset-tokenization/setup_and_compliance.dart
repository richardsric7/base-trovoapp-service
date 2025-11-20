import 'dart:convert';
import 'dart:developer';
import 'package:country_picker/country_picker.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SetupAndCompliance extends StatefulWidget {
  const SetupAndCompliance({Key? key}) : super(key: key);

  @override
  State<SetupAndCompliance> createState() => _SetupAndComplianceState();
}

class _SetupAndComplianceState extends State<SetupAndCompliance>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  String selectedCountry = '';
  final equityPercentageController = TextEditingController();
  final debtPercentageController = TextEditingController();
  bool hasSecApproval = false;
  bool hasSecApprovalId = false;
  bool agreeTransferTitleToCustodian = false;
  bool hasAllRequiredCustodianDocuments = false;
  bool hasAllRequiredManagerDocuments = false;
  int offeringType = 0;
  double equityPercentage = 0;
  double debtPercentage = 0;
  bool assetExisting = false;
  int fundingStructure = 0;
  String proceedPayoutCurrency = '';
  String assetQuoteCurrency = '';
  String selectedAssetSectorId = '';
  String selectedAssetSubSectorId = '';
  String selectedAssetTypeId = '';
  String selectedAssetCustodian = '';
  String selectedAssetManager = '';
  String secApprovalId = '';
  String tokenizationRequirementsUrl =
      'https://tokenization-requirements-app-xu8c6.ondigitalocean.app';
  bool formHasError = false;
  bool isCountryPickerOpen = false;
  late dynamic data = {};
  final _formKey = GlobalKey<FormState>();
  final Map<String, String> allowedCountries = {'NG': 'Nigeria'};
  late bool acceptTokenizationTermsAndAgreement;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;
    inspect(data);
    if (data != null && data.isNotEmpty) {
      selectedAssetSectorId = data!["assetSector"];
      selectedAssetSubSectorId = data!["assetSubSector"];
      selectedAssetTypeId = data!["assetType"];
      assetExisting = data!['assetAlreadyExists'] == 1;
      offeringType = data!["offeringType"].toString() == 'private' ? 1 : 0;
      fundingStructure = data['fundingStructure'] ?? 0;
      equityPercentage =
          double.tryParse(data['equityPercentage'].toString()) ?? 0;
      debtPercentage = double.tryParse(data['debtPercentage'].toString()) ?? 0;
      secApprovalId = data!["secApprovalIdNumber"];
      hasSecApproval = data!["secApproval"] == 1;
      selectedCountry = data!["assetCountryLocation"];
      selectedAssetCustodian = data!["approvedAssetCustodianId"].toString();
      selectedAssetManager = data!["assetManagerId"].toString();
      hasAllRequiredCustodianDocuments = data!["approvedAssetCustodianId"] != 0;
      hasAllRequiredManagerDocuments = data!["assetManagerId"] != 0;
      acceptTokenizationTermsAndAgreement =
          data['acceptTokenizationTermsAndAgreement'] == 1;
      equityPercentageController.text = equityPercentage.toString();
      debtPercentageController.text = debtPercentage.toString();
      agreeTransferTitleToCustodian =
          data!["agreeTransferTitleToCustodian"] != 0;
      for (
        var i = 0;
        i < appState.tokenizationData['countryConfigs'].length;
        i++
      ) {
        if (appState.tokenizationData['countryConfigs'][i]['countryCode']
                .toString()
                .toLowerCase() ==
            selectedCountry.toString().toLowerCase()) {
          proceedPayoutCurrency = appState
              .tokenizationData['countryConfigs'][i]['quoteCurrencyCode'];
          assetQuoteCurrency = appState
              .tokenizationData['countryConfigs'][i]['quoteCurrencyCode'];
        }
      }
    }
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              'setupandcompliance'.tr(),
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            setupAndCompliance(appState.tokenizationData),
          ],
        ),
      ),
    );
  }

  showCountryListPopup() {
    appState.dialogOpen = true;
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) {
      showCountryPicker(
        context: context,
        useSafeArea: true,
        countryFilter: ['NG'],
        countryListTheme: CountryListThemeData(
          backgroundColor: notifier.isDark
              ? notifier.getbluecolor50
              : Colors.white,
          textStyle: TextStyle(color: notifier.getblck),
          searchTextStyle: TextStyle(color: notifier.getblck),
          inputDecoration: InputDecoration(
            errorStyle: TextStyle(fontFamily: fontbody),
            labelText: 'Search',
            helperStyle: TextStyle(fontSize: 12, fontFamily: fontbody),
            // label: Text(labletext),
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15),
            ),
            prefixIcon: Icon(Icons.search, color: notifier.getblck),
            labelStyle: TextStyle(color: notifier.getblck),
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(15)),
            enabledBorder: OutlineInputBorder(
              borderSide: BorderSide(color: notifier.getgrey, width: 1),
              borderRadius: BorderRadius.circular(15),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: BorderSide(color: notifier.getgrey, width: 1),
              borderRadius: BorderRadius.circular(15),
            ),
          ),
        ),
        onSelect: (Country country) {
          setState(() {
            selectedCountry = country.countryCode;
            for (
              var i = 0;
              i < appState.tokenizationData['countryConfigs'].length;
              i++
            ) {
              if (appState.tokenizationData['countryConfigs'][i]['countryCode']
                      .toString()
                      .toLowerCase() ==
                  selectedCountry.toString().toLowerCase()) {
                proceedPayoutCurrency = appState
                    .tokenizationData['countryConfigs'][i]['quoteCurrencyCode'];
                assetQuoteCurrency = appState
                    .tokenizationData['countryConfigs'][i]['quoteCurrencyCode'];
              }
            }
          });
        },
      );
    });
    return SizedBox();
  }

  Widget setupAndCompliance(dynamic tokenizationData) {
    List<DropdownMenuItem<String>> assetSectors = [];
    for (var i = 0; i < tokenizationData!['assetSectors'].length; i++) {
      assetSectors.add(
        DropdownMenuItem(
          child: Text(
            tokenizationData!['assetSectors'][i]['sector'].toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: tokenizationData!['assetSectors'][i]['sector'].toString(),
        ),
      );
    }

    List<DropdownMenuItem<String>> assetSubsectors = getAssetSubsectorList(
      tokenizationData,
      selectedAssetSectorId,
    );

    List<DropdownMenuItem<String>> assetTypes = getAssetTypes(
      tokenizationData,
      selectedAssetSubSectorId,
    );

    return SingleChildScrollView(
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Row(
                children: [
                  Text(
                    "assetclassification".tr(),
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Container(
                child: Text(
                  "pleaseselectclassification".tr(),
                  textAlign: TextAlign.left,
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15),
              child: Container(
                width: width,
                child: Text(
                  textAlign: TextAlign.left,
                  "assetsector".tr(),
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
                    selectedAssetSectorId = value.toString();
                    assetSectors = getAssetSubsectorList(
                      tokenizationData,
                      selectedAssetSectorId,
                    );
                    selectedAssetSubSectorId = '';
                  });
                  if (selectedAssetSectorId.toLowerCase() ==
                      'finance and investment markets') {
                    assetExisting = true;
                  }
                },
                assetSectors,
                selectedAssetSectorId.isEmpty ? null : selectedAssetSectorId,
                'Select asset sector',
                context,
                null,
                validator: (value) {
                  if (selectedAssetSectorId.isEmpty) {
                    return "pleaseselectassetsector".tr();
                  }
                  return null;
                },
              ),
            ),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    "assetsubsector".tr(),
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    selectedAssetSubSectorId = value.toString();
                    selectedAssetTypeId = '';
                    assetSectors = getAssetTypes(
                      tokenizationData,
                      selectedAssetSubSectorId,
                    );
                  });
                },
                assetSubsectors,
                selectedAssetSubSectorId.isEmpty
                    ? null
                    : selectedAssetSubSectorId,
                "selectassetsubsector".tr(),
                context,
                null,
                validator: (value) {
                  if (selectedAssetSubSectorId.isEmpty) {
                    return "pleaseselectassetsubsector".tr();
                  }
                  return null;
                },
              ),
            ),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    "assettype".tr(),
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    selectedAssetTypeId = value.toString();
                    print('==========> $selectedAssetTypeId');
                  });
                },
                assetTypes,
                selectedAssetTypeId.isEmpty ? null : selectedAssetTypeId,
                'Select asset type',
                context,
                null,
                validator: (value) {
                  if (selectedAssetTypeId.isEmpty) {
                    return "pleaseselectassettype".tr();
                  }
                  return null;
                },
              ),
            ),
            if (selectedAssetSectorId.toLowerCase() !=
                'finance and investment markets') ...[
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    Text(
                      "assetstatus".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Container(
                  width: width,
                  child: Text(
                    "selectwhatappliestoasset".tr(),
                    textAlign: TextAlign.left,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Column(
                children: [
                  Row(
                    children: [
                      SizedBox(
                        height: 20,
                        child: Transform.scale(
                          scale: 1,
                          child: Radio<bool>(
                            value: true,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            groupValue: assetExisting,
                            onChanged: (value) => {
                              setState(() {
                                assetExisting = value!;
                              }),
                            },
                          ),
                        ),
                      ),
                      Text(
                        "assetexisting".tr(),
                        style: TextStyle(
                          fontSize: 14,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 70),
                  Row(
                    children: [
                      SizedBox(
                        height: 20,
                        child: Transform.scale(
                          scale: 1,
                          child: Radio<bool>(
                            value: false,
                            groupValue: assetExisting,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) => {
                              setState(() {
                                assetExisting = value!;
                              }),
                            },
                          ),
                        ),
                      ),
                      Text(
                        "assetnotyetexisting".tr(),
                        style: TextStyle(
                          fontSize: 14,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 30),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    Text(
                      "Funding Structure",
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Container(
                  width: width,
                  child: Text(
                    "Select funding structure",
                    textAlign: TextAlign.left,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Column(
                children: [
                  Row(
                    children: [
                      SizedBox(
                        height: 20,
                        child: Transform.scale(
                          scale: 1,
                          child: Radio<int>(
                            value: 0,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            groupValue: fundingStructure,
                            onChanged: (value) => {
                              setState(() {
                                fundingStructure = value!;
                              }),
                            },
                          ),
                        ),
                      ),
                      Text(
                        "Equity",
                        style: TextStyle(
                          fontSize: 14,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 70),
                  Row(
                    children: [
                      SizedBox(
                        height: 20,
                        child: Transform.scale(
                          scale: 1,
                          child: Radio<int>(
                            value: 1,
                            groupValue: fundingStructure,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) => {
                              setState(() {
                                fundingStructure = value!;
                              }),
                            },
                          ),
                        ),
                      ),
                      Text(
                        "Debt",
                        style: TextStyle(
                          fontSize: 14,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 70),
                  Row(
                    children: [
                      SizedBox(
                        height: 20,
                        child: Transform.scale(
                          scale: 1,
                          child: Radio<int>(
                            value: 2,
                            groupValue: fundingStructure,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: WidgetStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) => {
                              setState(() {
                                fundingStructure = value!;
                              }),
                            },
                          ),
                        ),
                      ),
                      Text(
                        "Hybrid",
                        style: TextStyle(
                          fontSize: 14,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              if (fundingStructure == 2) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'What is the equity percentage (%)',
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
                          setState(() {
                            equityPercentage = double.tryParse(value!) ?? 0;
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
                            equityPercentage = double.tryParse(value!) ?? 0;
                          });
                        },
                        autoFormatNumber: true,
                        isFiat: true,
                        controller: equityPercentageController,
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              if (fundingStructure == 2) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'What is the debt percentage (%)',
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
                          setState(() {
                            debtPercentage = double.tryParse(value!) ?? 0;
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
                            debtPercentage = double.tryParse(value!) ?? 0;
                          });
                        },
                        autoFormatNumber: true,
                        isFiat: true,
                        controller: debtPercentageController,
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ],
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Row(
                children: [
                  Text(
                    "assetlocation".tr(),
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Container(
                width: width,
                child: Text(
                  "pleaseselectcountry".tr(),
                  textAlign: TextAlign.left,
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    "country".tr(),
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Container(
                child: Card(
                  shadowColor: Colors.black,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10.0),
                  ),
                  color: notifier.isDark
                      ? notifier.getbluecolor90
                      : notifier.getaddsubwalletgrey,
                  child: TextButton(
                    onPressed: showCountryListPopup,
                    style: ButtonStyle(
                      elevation: WidgetStateProperty.all<double>(0),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          selectedCountry.isEmpty
                              ? "selectcountrylocation".tr()
                              : allowedCountries[selectedCountry]!,
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Icon(
                          Icons.keyboard_arrow_down_rounded,
                          color: notifier.getbluewhitecolor,
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            if (formHasError && selectedCountry.isEmpty) ...[
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "pleaseselectcountrylocation".tr(),
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
            SizedBox(height: height / 30),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Row(
                children: [
                  Text(
                    "offeringtype".tr(),
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Container(
                width: width,
                child: Text(
                  "selectofferingtype".tr(),
                  textAlign: TextAlign.left,
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ),
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Row(
                children: [
                  CheckItem(
                    "public".tr(),
                    () {
                      setState(() {
                        offeringType = 0;
                      });
                    },
                    borderColor: notifier.getbluewhitecolor,
                    foreColor: notifier.getbluewhitecolor,
                    backColor: offeringType == 0
                        ? notifier.getbluecolor60
                        : notifier.getwihitecolor,
                    icon: Icon(
                      Icons.public_rounded,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  CheckItem(
                    "private".tr(),
                    () {
                      setState(() {
                        offeringType = 1;
                      });
                    },
                    borderColor: notifier.getbluewhitecolor,
                    foreColor: notifier.getbluewhitecolor,
                    backColor: offeringType == 1
                        ? notifier.getbluecolor60
                        : notifier.getwihitecolor,
                    icon: Icon(
                      Icons.people_outline_outlined,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: height / 30),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Row(
                children: [
                  Text(
                    "requireddocuments".tr(),
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: TextButton(
                onPressed: () {
                  if (selectedAssetSectorId.isEmpty) {
                    popup(
                      context,
                      title: "info".tr(),
                      message:
                          'You must select an asset sector before you can view the requirements.',
                    );
                    return;
                  }

                  var url =
                      '$tokenizationRequirementsUrl/#/${selectedAssetSectorId.replaceAll(' ', '-').toLowerCase()}/${assetExisting ? '' : 'non-'}existing-assets';
                  appState.goToWebView(url);
                },
                child: Text(
                  "doyouhaverequiredmanagerdocs".tr(),
                  style: TextStyle(
                    decoration: TextDecoration.underline,
                    fontSize: 13,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ),
            Row(
              children: [
                Row(
                  children: [
                    Transform.scale(
                      scale: 1,
                      child: Radio<bool>(
                        value: true,
                        groupValue: hasAllRequiredManagerDocuments,
                        activeColor: notifier.getbluewhitecolor,
                        fillColor: WidgetStateColor.resolveWith(
                          (states) => notifier.getbluewhitecolor,
                        ),
                        onChanged: (value) => {
                          setState(() {
                            hasAllRequiredManagerDocuments = value!;
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
                        groupValue: hasAllRequiredManagerDocuments,
                        onChanged: (value) => {
                          setState(() {
                            hasAllRequiredManagerDocuments = value!;
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
            SizedBox(height: height / 30),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Transform.scale(
                  scale: 1,
                  child: FormField(
                    builder: (state) {
                      return Checkbox(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(Radius.circular(5)),
                        ),
                        activeColor: notifier.isDark
                            ? notifier.getbluecolor50
                            : notifier.getbluecolor90,
                        side: BorderSide(
                          color: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
                        ),
                        value: agreeTransferTitleToCustodian,
                        onChanged: (bool? value) {
                          setState(() {
                            agreeTransferTitleToCustodian = value!;
                          });
                        },
                      );
                    },
                    validator: (value) {
                      if (!agreeTransferTitleToCustodian) {
                        setState(() {
                          formHasError = true;
                        });
                        return '';
                      }

                      return null;
                    },
                  ),
                ),
                Container(
                  width: width / 1.2,
                  child: Text(
                    "Tokenizing your asset requires ownership transfer of the asset to a licensed custodian. Agree?",
                    overflow: TextOverflow.visible,
                    style: TextStyle(
                      fontSize: 15,
                      color: formHasError && !agreeTransferTitleToCustodian
                          ? Colors.red
                          : notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Transform.scale(
                  scale: 1,
                  child: FormField(
                    builder: (state) {
                      return Checkbox(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(Radius.circular(5)),
                        ),
                        activeColor: notifier.isDark
                            ? notifier.getbluecolor50
                            : notifier.getbluecolor90,
                        side: BorderSide(
                          color: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
                        ),
                        value: acceptTokenizationTermsAndAgreement,
                        onChanged: (bool? value) {
                          setState(() {
                            acceptTokenizationTermsAndAgreement = value!;
                          });
                        },
                      );
                    },
                    validator: (value) {
                      if (!acceptTokenizationTermsAndAgreement) {
                        setState(() {
                          formHasError = true;
                        });
                        return '';
                      }

                      return null;
                    },
                  ),
                ),
                Container(
                  width: width / 1.2,
                  child: RichText(
                    text: TextSpan(
                      text: 'I agree to the Trovotech ',
                      children: [
                        TextSpan(
                          text: 'Tokenization Terms and Conditions',
                          style: TextStyle(
                            decoration: TextDecoration.underline,
                            color:
                                formHasError &&
                                    !acceptTokenizationTermsAndAgreement
                                ? Colors.red
                                : notifier.getbluewhitecolor,
                            fontVariations: [FontVariation('wght', 700)],
                          ),
                          recognizer: TapGestureRecognizer()
                            ..onTap = () {
                              // appState.viewData = {
                              //   'url': 'https://olaratech.com/legal?tab=terms',
                              // };
                              // appState.currentAction = PageAction(
                              //   state: PageState.addPage,
                              //   page: appWebViewPageConfig,
                              // );
                            },
                        ),
                      ],
                      style: TextStyle(
                        fontSize: 13.sp,
                        fontFamily: fontbody,
                        color:
                            formHasError && !acceptTokenizationTermsAndAgreement
                            ? Colors.red
                            : notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(height: height / 20),
            Button(
              "saveandcontinuee".tr(),
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                submit();
              },
            ),
            if (appState.viewData?['id'] != null) ...[
              SizedBox(height: 10),
              Button(
                "deletetokenization".tr(),
                Colors.red,
                wihitecolor,
                onTap: () {
                  confirmTokenizationDeletePopup(
                    context,
                    onConfirmationSuccess: deleteTokenization,
                  );
                },
              ),
            ],
            SizedBox(height: height / 10),
            Padding(
              padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom,
              ),
            ),
          ],
        ),
      ),
    );
  }

  submit() async {
    var message = "";
    final form = _formKey.currentState;
    formHasError = false;

    if (selectedCountry.isEmpty) {
      setState(() {
        formHasError = true;
      });
    }

    if (!form!.validate() || formHasError) return;

    if (!hasAllRequiredManagerDocuments) {
      message +=
          "You need to acquire all the documents listed in the tokenization requirements document before you can proceed.\n\n";
    }

    if (message.isNotEmpty) {
      popup(context, title: "info".tr(), message: message);
      return;
    }

    try {
      showLoader(context);
      data = appState.viewData;
      data["assetSector"] = selectedAssetSectorId;
      data["assetSubSector"] = selectedAssetSubSectorId;
      data["assetType"] = selectedAssetTypeId;
      data["offeringType"] = offeringType == 1 ? 'private' : 'public';
      data["approvedAssetCustodianId"] = selectedAssetCustodian.length > 0
          ? int.parse(selectedAssetCustodian)
          : 1;
      data["assetManagerId"] = selectedAssetManager.length > 0
          ? int.parse(selectedAssetManager)
          : 1;
      data["agreeTransferTitleToCustodian"] = agreeTransferTitleToCustodian
          ? 1
          : 0;
      data["assetAlreadyExists"] = assetExisting ? 1 : 0;
      data["secApproval"] = hasSecApproval ? 1 : 0;
      data["secApprovalIdNumber"] = secApprovalId;
      data["assetCountryLocation"] = selectedCountry;
      data['proceedPayoutCurrency'] = proceedPayoutCurrency;
      data['assetQuoteCurrency'] = assetQuoteCurrency;
      data['fundingStructure'] = fundingStructure;
      data['equityPercentage'] = equityPercentage;
      data['debtPercentage'] = debtPercentage;
      data['acceptTokenizationTermsAndAgreement'] =
          acceptTokenizationTermsAndAgreement ? 1 : 0;

      String requestBody = jsonEncode(data);

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      if (responseData['statusCode'] == 200) {
        appState.viewData = responseData['data'];
        await refreshCurrentTokenizationInfo(appState);
        await fetchBanksList();
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: TokenizeAssetViewPageConfig,
        );
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

  Future<void> fetchBanksList() async {
    var uri = '/v1/banks/${selectedCountry}';

    Map responseData = await makeGetRequest(
      uri: Uri.encodeFull(uri),
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.signer!,
    );
    inspect(responseData['data']);
    if (responseData['statusCode'] == 200) {
      appState.tokenizationData['banks'] = responseData['data'];
    }
  }

  List<DropdownMenuItem<String>> getAssetSubsectorList(
    dynamic data,
    String selectedSectorId,
  ) {
    List<DropdownMenuItem<String>> assetSubsectors = [];
    for (var i = 0; i < data!['assetSubSectors'].length; i++) {
      if (data!['assetSubSectors'][i]['assetSectorId'] == selectedSectorId) {
        assetSubsectors.add(
          DropdownMenuItem(
            child: Text(
              data!['assetSubSectors'][i]['subSector'],
              overflow: TextOverflow.ellipsis,
            ),
            value: data!['assetSubSectors'][i]['subSector'],
          ),
        );
      }
    }
    return assetSubsectors;
  }

  List<DropdownMenuItem<String>> getAssetTypes(data, selectedSubsectorId) {
    List<DropdownMenuItem<String>> assetTypes = [];
    for (var i = 0; i < data!['assetTypes'].length; i++) {
      if (data!['assetTypes'][i]['assetSubSectorId'] == selectedSubsectorId) {
        assetTypes.add(
          DropdownMenuItem(
            child: Text(
              data!['assetTypes'][i]['assetType'],
              overflow: TextOverflow.ellipsis,
            ),
            value: data!['assetTypes'][i]['id'].toString(),
          ),
        );
      }
    }
    assetTypes.sort((a, b) {
      if (a.key?.toString().toLowerCase() == 'other') return 1;
      if (b.key?.toString().toLowerCase() == 'other') return -1;
      return a.value?.compareTo(b.value ?? '') ?? 0;
    });

    return assetTypes;
  }

  List<DropdownMenuItem<int>> getAssetCustodians(data) {
    List<DropdownMenuItem<int>> assetCustodians = [];
    for (var i = 0; i < data!['assetCustodians'].length; i++) {
      assetCustodians.add(
        DropdownMenuItem(
          child: Text(
            '${data!['assetCustodians'][i]['assetCustodianName']}, ${data!['assetCustodians'][i]['assetCustodianAddress']} ${data!['assetCustodians'][i]['assetCustodianCountry']}',
            overflow: TextOverflow.ellipsis,
          ),
          value: data!['assetCustodians'][i]['id'],
        ),
      );
    }
    return assetCustodians;
  }

  List<DropdownMenuItem<int>> getAssetManagers(data) {
    List<DropdownMenuItem<int>> assetManagers = [];
    for (var i = 0; i < data!['assetManagers'].length; i++) {
      assetManagers.add(
        DropdownMenuItem(
          child: Text(
            '${data!['assetManagers'][i]['assetManagerName']}, ${data!['assetManagers'][i]['assetManagerAddress']} ${data!['assetManagers'][i]['assetManagerCountry']}',
            overflow: TextOverflow.ellipsis,
          ),
          value: data!['assetManagers'][i]['id'],
        ),
      );
    }
    return assetManagers;
  }

  String getSelectedAssetCustodianLabel(id, data) {
    var label = "Select asset custodian";
    for (var i = 0; i < data!['assetCustodians'].length; i++) {
      if (data!['assetCustodians'][i]['id'] == id) {
        label =
            '${data!['assetCustodians'][i]['assetCustodianName']}, ${data!['assetCustodians'][i]['assetCustodianAddress']} ${data!['assetCustodians'][i]['assetCustodianCountry']}';
        break;
      }
    }
    return label;
  }

  String getSelectedAssetManagerLabel(id, data) {
    var label = "Select asset manager";
    for (var i = 0; i < data!['assetManagers'].length; i++) {
      if (data!['assetManagers'][i]['id'] == id) {
        label =
            '${data!['assetManagers'][i]['assetManagerName']}, ${data!['assetManagers'][i]['assetManagerAddress']} ${data!['assetManagers'][i]['assetManagersCountry']}';
        break;
      }
    }
    return label;
  }

  Widget CheckItem(
    String name,
    void Function()? onClick, {
    required Color backColor,
    required Color foreColor,
    required Color borderColor,
    required Icon icon,
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
            icon,
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

  void deleteTokenization() async {
    try {
      showLoader(context);

      Map responseData = await makeDeleteRequest(
        uri: '/v1/tokenization/${appState.viewData!['id']}',
        body: "",
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      hideLoader(context);

      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        appState.currentAction = PageAction(
          state: PageState.replace,
          page: BottomHomePageConfig,
        );
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}
