import 'dart:convert';
import 'dart:developer';
import 'package:country_picker/country_picker.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  // bool hasCustodianAgreement = true;
  bool hasSecApproval = false;
  bool hasSecApprovalId = false;
  bool hasAllRequiredCustodianDocuments = false;
  bool hasAllRequiredManagerDocuments = false;
  int offeringType = 0;
  String selectedAssetSectorId = '';
  String selectedAssetSubSectorId = '';
  String selectedAssetTypeId = '';
  String selectedAssetCustodian = '';
  String selectedAssetManager = '';
  String secApprovalId = '';
  bool formHasError = false;
  bool isCountryPickerOpen = false;
  late dynamic data = {};
  final _formKey = GlobalKey<FormState>();

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
    if (data != null && data.isNotEmpty) {
      selectedAssetSectorId = data!["assetSector"];
      selectedAssetSubSectorId = data!["assetSubSector"];
      selectedAssetTypeId = data!["assetType"];
      offeringType = data!["offeringType"].toString() == 'private' ? 1 : 0;
      secApprovalId = data!["secApprovalIdNumber"];
      hasSecApproval = data!["secApproval"] == 1;
      selectedCountry = data!["assetCountryLocation"];
      selectedAssetCustodian = data!["approvedAssetCustodianId"].toString();
      selectedAssetManager = data!["assetManagerId"].toString();
      hasAllRequiredCustodianDocuments = data!["approvedAssetCustodianId"] != 0;
      hasAllRequiredManagerDocuments = data!["assetManagerId"] != 0;
      // hasCustodianAgreement = data!["approvedAssetCustodianInfo"].length != 0;
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
          backgroundColor:
              notifier.isDark ? notifier.getbluecolor50 : Colors.white,
          textStyle: TextStyle(color: notifier.getblck),
          searchTextStyle: TextStyle(color: notifier.getblck),
          inputDecoration: InputDecoration(
            errorStyle: TextStyle(
              fontFamily: fontbody,
            ),
            labelText: 'Search',
            helperStyle: TextStyle(
              fontSize: 12,
              fontFamily: fontbody,
            ),
            // label: Text(labletext),
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15),
            ),
            prefixIcon: Icon(Icons.search, color: notifier.getblck),
            labelStyle: TextStyle(color: notifier.getblck),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15),
            ),
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
            selectedCountry = country.name;
          });
        },
      );
    });
    return SizedBox();
  }

  Widget setupAndCompliance(dynamic tokenizationData) {
    List<DropdownMenuItem<String>> assetSectors = [];
    for (var i = 0; i < tokenizationData!['assetSectors'].length; i++) {
      assetSectors.add(DropdownMenuItem(
          child: Text(
            tokenizationData!['assetSectors'][i]['sector'].toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: tokenizationData!['assetSectors'][i]['sector'].toString()));
    }

    List<DropdownMenuItem<String>> assetSubsectors =
        getAssetSubsectorList(tokenizationData, selectedAssetSectorId);

    List<DropdownMenuItem<String>> assetTypes =
        getAssetTypes(tokenizationData, selectedAssetSubSectorId);

    // List<DropdownMenuItem<int>> assetCustodians =
    //     getAssetCustodians(tokenizationData);

    List<DropdownMenuItem<int>> assetManagers =
        getAssetManagers(tokenizationData);

    return SingleChildScrollView(
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(
              height: height / 50,
            ),
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
            SizedBox(
              height: height / 50,
            ),
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
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    selectedAssetSectorId = value.toString();
                    assetSectors = getAssetSubsectorList(
                        tokenizationData, selectedAssetSectorId);
                  });
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
            SizedBox(
              height: height / 50,
            ),
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
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    selectedAssetSubSectorId = value.toString();
                    selectedAssetTypeId = '';
                    assetSectors = getAssetTypes(
                        tokenizationData, selectedAssetSubSectorId);
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
            SizedBox(
              height: height / 50,
            ),
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
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    selectedAssetTypeId = value.toString();
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
            SizedBox(
              height: height / 15,
            ),
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
            SizedBox(
              height: height / 70,
            ),
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
                  )
                ],
              ),
            ),
            SizedBox(
              height: height / 15,
            ),
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
            SizedBox(
              height: height / 50,
            ),
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
            SizedBox(
              height: height / 70,
            ),
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
                        elevation: MaterialStateProperty.all<double>(0)),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          selectedCountry.isEmpty
                              ? "selectcountrylocation".tr()
                              : selectedCountry,
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
            SizedBox(
              height: height / 30,
            ),
            // Padding(
            //   padding: const EdgeInsets.symmetric(horizontal: 15.0),
            //   child: Row(
            //     children: [
            //       Text(
            //         "assetcustodian".tr(),
            //         style: TextStyle(
            //           fontSize: 18,
            //           fontFamily: fontsemibold,
            //           color: notifier.getbluewhitecolor,
            //         ),
            //       ),
            //     ],
            //   ),
            // ),
            // Padding(
            //   padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 5),
            //   child: Row(
            //     mainAxisAlignment: MainAxisAlignment.spaceBetween,
            //     children: [
            //       Text(
            //         "selectapprovedcustodian".tr(),
            //         style: TextStyle(
            //           fontSize: 13,
            //           fontFamily: fontsemibold,
            //           color: notifier.getbluewhitecolor,
            //         ),
            //       ),
            //     ],
            //   ),
            // ),
            // SizedBox(
            //   height: height / 70,
            // ),
            // Padding(
            //   padding: const EdgeInsets.symmetric(horizontal: 10.0),
            //   child: dropdown(
            //     (value) {
            //       setState(() {
            //         selectedAssetCustodian = value.toString();
            //         assetCustodians = getAssetCustodians(tokenizationData);
            //       });
            //     },
            //     assetCustodians,
            //     selectedAssetCustodian.isNotEmpty
            //         ? assetCustodians
            //             .where((element) {
            //               return element.value.toString() ==
            //                   selectedAssetCustodian;
            //             })
            //             .first
            //             .value
            //         : null,
            //     selectedAssetCustodian.isNotEmpty
            //         ? getSelectedAssetCustodianLabel(
            //             selectedAssetCustodian, tokenizationData)
            //         : "selectassetcustodian".tr(),
            //     context,
            //     null,
            //     validator: (value) {
            //       if (selectedAssetCustodian.isEmpty) {
            //         return "pleaseselectassetcustodian".tr();
            //       }
            //       return null;
            //     },
            //   ),
            // ),
            // Padding(
            //   padding: const EdgeInsets.symmetric(horizontal: 10.0),
            //   child: TextButton(
            //     onPressed: () =>
            //         appState.goToWebView(tokenizationRequirementsUrl),
            //     child: Text(
            //       "pleasecheckrequirements".tr(),
            //       style: TextStyle(
            //         decoration: TextDecoration.underline,
            //         fontSize: 13,
            //         fontFamily: fontsemibold,
            //         color: notifier.getbluewhitecolor,
            //       ),
            //     ),
            //   ),
            // ),
            // Row(
            //   children: [
            //     Row(
            //       children: [
            //         Transform.scale(
            //           scale: 1,
            //           child: Radio<bool>(
            //             value: true,
            //             groupValue: hasAllRequiredCustodianDocuments,
            //             activeColor: notifier.getbluewhitecolor,
            //             fillColor: MaterialStateColor.resolveWith(
            //                 (states) => notifier.getbluewhitecolor),
            //             onChanged: (value) => {
            //               setState(
            //                 () {
            //                   hasAllRequiredCustodianDocuments = value!;
            //                 },
            //               )
            //             },
            //           ),
            //         ),
            //         Text(
            //           "yes".tr(),
            //           style: TextStyle(
            //             fontSize: 14,
            //             fontFamily: fontsemibold,
            //             color: notifier.getbluewhitecolor,
            //           ),
            //         ),
            //       ],
            //     ),
            //     Row(
            //       children: [
            //         Transform.scale(
            //           scale: 1,
            //           child: Radio<bool>(
            //             value: false,
            //             activeColor: notifier.getbluewhitecolor,
            //             fillColor: MaterialStateColor.resolveWith(
            //                 (states) => notifier.getbluewhitecolor),
            //             groupValue: hasAllRequiredCustodianDocuments,
            //             onChanged: (value) => {
            //               setState(
            //                 () {
            //                   hasAllRequiredCustodianDocuments = value!;
            //                 },
            //               )
            //             },
            //           ),
            //         ),
            //         Text(
            //           "no".tr(),
            //           style: TextStyle(
            //             fontSize: 14,
            //             fontFamily: fontsemibold,
            //             color: notifier.getbluewhitecolor,
            //           ),
            //         ),
            //       ],
            //     ),
            //   ],
            // ),
            // SizedBox(
            //   height: height / 50,
            // ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
              child: Row(
                children: [
                  Text(
                    "assetissuer".tr(),
                    // "assetmanager".tr(),
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
              padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 5),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    // "selectapprovedmanager".tr(),
                    "selectapprovedissuer".tr(),
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    selectedAssetManager = value.toString();
                    assetManagers = getAssetManagers(tokenizationData);
                  });
                },
                assetManagers,
                selectedAssetManager.isNotEmpty
                    ? assetManagers
                        .where((element) {
                          return element.value.toString() ==
                              selectedAssetManager;
                        })
                        .first
                        .value
                    : null,
                selectedAssetManager.isNotEmpty
                    ? getSelectedAssetManagerLabel(
                        selectedAssetManager, tokenizationData)
                    : "selectapprovedmanager".tr(),
                context,
                null,
                validator: (value) {
                  if (selectedAssetManager.isEmpty) {
                    return "pleaseselectassetmanager".tr();
                  }
                  return null;
                },
              ),
            ),
            // Padding(
            //   padding: const EdgeInsets.symmetric(horizontal: 10.0),
            //   child: TextButton(
            //     onPressed: () =>
            //         appState.goToWebView(tokenizationRequirementsUrl),
            //     child: Text(
            //       "doyouhaverequiredmanagerdocs".tr(),
            //       style: TextStyle(
            //         decoration: TextDecoration.underline,
            //         fontSize: 13,
            //         fontFamily: fontsemibold,
            //         color: notifier.getbluewhitecolor,
            //       ),
            //     ),
            //   ),
            // ),
            // Row(
            //   children: [
            //     Row(
            //       children: [
            //         Transform.scale(
            //           scale: 1,
            //           child: Radio<bool>(
            //             value: true,
            //             groupValue: hasAllRequiredManagerDocuments,
            //             activeColor: notifier.getbluewhitecolor,
            //             fillColor: MaterialStateColor.resolveWith(
            //                 (states) => notifier.getbluewhitecolor),
            //             onChanged: (value) => {
            //               setState(
            //                 () {
            //                   hasAllRequiredManagerDocuments = value!;
            //                 },
            //               )
            //             },
            //           ),
            //         ),
            //         Text(
            //           "yes".tr(),
            //           style: TextStyle(
            //             fontSize: 14,
            //             fontFamily: fontsemibold,
            //             color: notifier.getbluewhitecolor,
            //           ),
            //         ),
            //       ],
            //     ),
            //     Row(
            //       children: [
            //         Transform.scale(
            //           scale: 1,
            //           child: Radio<bool>(
            //             value: false,
            //             activeColor: notifier.getbluewhitecolor,
            //             fillColor: MaterialStateColor.resolveWith(
            //                 (states) => notifier.getbluewhitecolor),
            //             groupValue: hasAllRequiredManagerDocuments,
            //             onChanged: (value) => {
            //               setState(
            //                 () {
            //                   hasAllRequiredManagerDocuments = value!;
            //                 },
            //               )
            //             },
            //           ),
            //         ),
            //         Text(
            //           "no".tr(),
            //           style: TextStyle(
            //             fontSize: 14,
            //             fontFamily: fontsemibold,
            //             color: notifier.getbluewhitecolor,
            //           ),
            //         ),
            //       ],
            //     ),
            //   ],
            // ),
            // SizedBox(
            //   height: height / 50,
            // ),
            SizedBox(
              height: height / 20,
            ),
            Button(
              "saveandcontinuee".tr(),
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                submit();
              },
            ),
            SizedBox(height: 10),
            Button(
              "deletetokenization".tr(),
              Colors.red,
              wihitecolor,
              onTap: () {
                deleteTokenization();
              },
            ),
            SizedBox(height: height / 10),
            Padding(
              padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom),
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

    if (!form!.validate() && !formHasError) return;

    // if (!hasAllRequiredCustodianDocuments) {
    //   message +=
    //       "You need to acquire all the documents in the asset custodian's required documents list before you can proceed.\n\n";
    // }

    // if (!hasAllRequiredManagerDocuments) {
    //   message +=
    //       "You need to acquire all the documents in the asset manager's required documents list before you can proceed.\n\n";
    // }

    // if (offeringType == 0 && !hasCustodianAgreement) {
    //   message +=
    //       "An asset custodian agreement is needed in this process. You need to obtain an agreement with an asset custodian to proceed.\n\n";
    // }

    if (message.isNotEmpty) {
      popup(context, title: "info".tr(), message: message);
      return;
    }

    try {
      showLoader(context);
      // make initial request to the server using the
      // following credential
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;
      var marketMakingWallet = appState.activeDistributionWalletPublicKey!;
      var newData = {...data as Map};

      Map map = {
        "assetSector": selectedAssetSectorId,
        "assetSubSector": selectedAssetSubSectorId,
        "assetType": selectedAssetTypeId,
        "offeringType": offeringType == 1 ? 'private' : 'public',
        "approvedAssetCustodianId": selectedAssetCustodian.length > 0
            ? int.parse(selectedAssetCustodian)
            : 1,
        "assetManagerId": selectedAssetManager.length > 0
            ? int.parse(selectedAssetManager)
            : 1,
        "marketMakingWallet": marketMakingWallet,
        "secApproval": hasSecApproval ? 1 : 0,
        "secApprovalIdNumber": secApprovalId,
        "assetCountryLocation": selectedCountry,
        "numberOfTokenToBeSold": newData['numberOfTokenToBeSold'],
        "numberOfTokenToBeIssued": newData['numberOfTokenToBeIssued'],
        "totalTokenHeldByManager": newData['totalTokenHeldByManager'],
        "pricePerToken": newData['pricePerToken'],
        "proceedPayoutType": newData['proceedPayoutType'],
        "assetCode": newData['assetCode'],
        "assetName": newData['assetName'],
        "salesStart": newData['salesStart'],
        "salesEnd": newData['salesEnd'],
        "capOnPurchase": newData['capOnPurchase'],
        "capQuantity": newData['capQuantity'],
        "capDurationInDays": newData['capDurationInDays'],
        "proceedCycle": newData['proceedCycle'],
        "walletToHoldAssetsNotForSale": newData['walletToHoldAssetsNotForSale'],
        "assetLogo": newData['assetLogo'],
        "exemptedCountries": newData['exemptedCountries'],
        "hasAdditionalKYCRequirements": newData['hasAdditionalKYCRequirements'],
        "assetQuoteCurrency": newData['assetQuoteCurrency'],
        "proceedPayoutCurrency": newData['proceedPayoutCurrency'],
        "additionalKYCRequirements": newData['additionalKYCRequirements'],
        "investorAccreditationRequired":
            newData['investorAccreditationRequired'],
        "tokenizationFeeId": newData['tokenizationFeeId'],
        "assetAlreadyExists": newData['assetAlreadyExists'],
        "ownershipType": newData['ownershipType'],
        "ownershipKind": newData['ownershipKind'],
        "assetDescription": newData['assetDescription'],
        "assetPhysicalAddress": newData['assetPhysicalAddress'],
        "assetLatitude": newData['assetLatitude'],
        "assetLongitude": newData['assetLongitude'],
        "assetOwnerName": newData['assetOwnerName'],
        "assetOwnerAddress": newData['assetOwnerAddress'],
        "assetManagerName": newData['assetManagerName'],
        "assetManagerAddress": newData['assetManagerAddress'],
        "assetCurrentValue": newData['assetCurrentValue'],
        "valueOfTokenizedAsset": newData['valueOfTokenizedAsset'],
        "protectionMethods": newData['protectionMethods'],
        "insuranceCompanyName": newData['insuranceCompanyName'],
        "insurance_policy_number": newData['insurance_policy_number'],
        "insurancePolicyHolder": newData['insurancePolicyHolder'],
        "percentageValueOfInsurance": newData['percentageValueOfInsurance'],
        "IsFreeFromLiensAndEncumbrances":
            newData['IsFreeFromLiensAndEncumbrances'],
      };
      String requestBody = jsonEncode(map);
      inspect(map);

      print('requestBody =======> $requestBody');
      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: mintingWalletPublicKey,
      );
      hideLoader(context);

      print('responseData ${responseData['data']}');
      inspect(responseData['data']);

      if (responseData['statusCode'] == 200) {
        appState.viewData = responseData['data'];
        // appState.viewData!["assetSector"] = selectedAssetSectorId;
        // appState.viewData!["assetSubSector"] = selectedAssetSubSectorId;
        // appState.viewData!["assetType"] = selectedAssetTypeId;
        // appState.viewData!["offeringType"] =
        //     offeringType == 1 ? 'private' : 'public';
        // appState.viewData!["approvedAssetCustodianId"] =
        //     selectedAssetCustodian.length > 0
        //         ? int.parse(selectedAssetCustodian)
        //         : 1;
        // appState.viewData!["marketMakingWallet"] = marketMakingWallet;
        // appState.viewData!["secApproval"] = hasSecApproval ? 1 : 0;
        // appState.viewData!["secApprovalIdNumber"] = secApprovalId;
        // appState.viewData!["assetCountryLocation"] = selectedCountry;
        // appState.viewData!["assetCountryLocation"] = selectedCountry;
        // appState.viewData!["assetCountryLocation"] = selectedCountry;

        appState.currentAction = PageAction(
            state: PageState.addPage, page: TokenizeAssetViewPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
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
          overlayColor:
              MaterialStateProperty.all<Color>(notifier.getsplashgrey),
          elevation: MaterialStateProperty.all<double>(0),
          backgroundColor: MaterialStateProperty.all<Color>(backColor),
          side: MaterialStateProperty.all(
            BorderSide(color: borderColor, width: 1, style: BorderStyle.solid),
          ),
          shape: MaterialStateProperty.all<RoundedRectangleBorder>(
            const RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(10),
              ),
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
                  color: foreColor, fontFamily: fontbody, fontSize: fontSize),
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
        publicKey: appState.activeTokenizationWalletPublicKey!,
      );

      hideLoader(context);

      print('responseData token information  ${responseData['data']}');
      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        appState.currentAction = PageAction(
          state: PageState.replace,
          page: BottomHomePageConfig,
        );
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}
