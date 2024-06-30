import 'dart:convert';

import 'package:country_picker/country_picker.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
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
  String selectedCountry = 'Nigeria';
  bool hasCustodianAgreement = true;
  bool hasSecApproval = false;
  bool hasSecApprovalId = false;
  bool hasAllRequiredDocuments = false;
  int offeringType = 0;
  String selectedAssetSectorId = 'Real Estate Sector';
  String selectedAssetSubSectorId = 'Land';
  String selectedAssetTypeId = '';
  String selectedAssetCustodian = '';
  String secApprovalId = '';
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
    print(appState.viewData);
    if (data != null) {
      selectedAssetSectorId = data!["assetSector"];
      selectedAssetSubSectorId = data!["assetSubSector"];
      selectedAssetTypeId = data!["assetType"];
      offeringType = data!["offeringType"].toString() == 'private' ? 1 : 0;
      secApprovalId = data!["secApprovalIdNumber"];
      hasSecApproval = data!["secApproval"] == 1;
      selectedCountry = data!["assetCountryLocation"];
      selectedAssetCustodian = data!["approvedAssetCustodianId"].toString();
      hasAllRequiredDocuments = data!["approvedAssetCustodianId"] != 0;
      hasCustodianAgreement = data!["approvedAssetCustodianInfo"].length != 0;
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

    List<DropdownMenuItem<int>> assetCustodians =
        getAssetCustodians(tokenizationData);

    return SingleChildScrollView(
      child: Form(
        key: _formKey,
        child: Column(
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
                null,
                assetSectors.length > 0 ? assetSectors.first.value : '',
                context,
                null,
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
                    assetSectors = getAssetTypes(
                        tokenizationData, selectedAssetSubSectorId);
                  });
                },
                assetSubsectors,
                null,
                assetSubsectors.length > 0 ? assetSubsectors.first.value : '',
                context,
                null,
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
                    onPressed: () {
                      showCountryPicker(
                        context: context,
                        onSelect: (Country country) {
                          setState(() {
                            selectedCountry = country.name;
                          });
                        },
                      );
                    },
                    style: ButtonStyle(
                        elevation: MaterialStateProperty.all<double>(0)),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          selectedCountry,
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Icon(Icons.keyboard_arrow_down_rounded),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: TextButton(
                onPressed: () async {},
                child: Text(
                  "doyouhaveallrequireddocs".tr(),
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
                        groupValue: hasAllRequiredDocuments,
                        activeColor: notifier.getbluewhitecolor,
                        fillColor: MaterialStateColor.resolveWith(
                            (states) => notifier.getbluewhitecolor),
                        onChanged: (value) => {
                          setState(
                            () {
                              hasAllRequiredDocuments = value!;
                            },
                          )
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
                        fillColor: MaterialStateColor.resolveWith(
                            (states) => notifier.getbluewhitecolor),
                        groupValue: hasAllRequiredDocuments,
                        onChanged: (value) => {
                          setState(
                            () {
                              hasAllRequiredDocuments = value!;
                            },
                          )
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
            if (hasAllRequiredDocuments) ...[
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    Text(
                      "custodianagreement".tr(),
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
                    "entercustodianinformation".tr(),
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
                    "doyouhavecustodianagreement".tr(),
                    style: TextStyle(
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
                          groupValue: hasCustodianAgreement,
                          activeColor: notifier.getbluewhitecolor,
                          fillColor: MaterialStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor),
                          onChanged: (value) => {
                            setState(
                              () {
                                hasCustodianAgreement = value!;
                              },
                            )
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
                          fillColor: MaterialStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor),
                          groupValue: hasCustodianAgreement,
                          onChanged: (value) => {
                            setState(
                              () {
                                hasCustodianAgreement = value!;
                              },
                            )
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
              SizedBox(
                height: height / 50,
              ),
              if (hasCustodianAgreement) ...[
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        "selectapprovedcustodian".tr(),
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
                        selectedAssetCustodian = value.toString();
                        assetCustodians = getAssetCustodians(tokenizationData);
                      });
                    },
                    assetCustodians,
                    selectedAssetCustodian.isNotEmpty
                        ? assetCustodians
                            .where((element) {
                              return element.value.toString() ==
                                  selectedAssetCustodian;
                            })
                            .first
                            .value
                        : null,
                    selectedAssetCustodian.isNotEmpty
                        ? getSelectedAssetCustodianLabel(
                            selectedAssetCustodian, tokenizationData)
                        : selectedAssetCustodian,
                    context,
                    null,
                    validator: (value) {
                      if (hasCustodianAgreement &&
                          selectedAssetCustodian.isEmpty) {
                        return "pleaseselectassetcustodian".tr();
                      }
                      return null;
                    },
                  ),
                ),
              ] else ...[
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15),
                  child: Container(
                    width: width,
                    child: Text(
                      "reachouttoanapprovedcustodian".tr(),
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
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10),
                      child: Container(
                        width: width / 1.07,
                        decoration: BoxDecoration(
                          borderRadius:
                              const BorderRadius.all(Radius.circular(15.0)),
                          color: notifier.getaddsubwalletgrey,
                        ),
                        child: Padding(
                          padding: const EdgeInsets.all(10.0),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              for (var item
                                  in getAssetCustodians(tokenizationData)) ...[
                                Column(
                                  children: [
                                    Row(
                                      children: [
                                        Text('- '),
                                        item.child,
                                      ],
                                    ),
                                    SizedBox(
                                      height: 10,
                                    )
                                  ],
                                ),
                              ],
                            ],
                          ),
                        ),
                      ),
                    )
                  ],
                ),
              ],
              // if (offeringType == 0 && hasCustodianAgreement) ...[
              //   SizedBox(
              //     height: height / 15,
              //   ),
              //   Padding(
              //     padding: const EdgeInsets.symmetric(horizontal: 15.0),
              //     child: Row(
              //       children: [
              //         Container(
              //           child: Text(
              //             "secapproval".tr(),
              //             style: TextStyle(
              //               fontSize: 18,
              //               fontFamily: fontsemibold,
              //               color: notifier.getbluewhitecolor,
              //             ),
              //           ),
              //         ),
              //       ],
              //     ),
              //   ),
              //   Padding(
              //     padding: const EdgeInsets.symmetric(horizontal: 15.0),
              //     child: Container(
              //       width: width,
              //       child: Text(
              //         "pleasefillapprovalinfo".tr(),
              //         textAlign: TextAlign.left,
              //         style: TextStyle(
              //           fontSize: 15,
              //           fontFamily: fontbody,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ),
              //   SizedBox(
              //     height: height / 50,
              //   ),
              //   Padding(
              //     padding: const EdgeInsets.symmetric(horizontal: 15),
              //     child: Container(
              //       width: width,
              //       child: Text(
              //         "isyourassetapproved".tr(),
              //         style: TextStyle(
              //           fontSize: 13,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ),
              //   Row(
              //     children: [
              //       Row(
              //         children: [
              //           Transform.scale(
              //             scale: 1,
              //             child: Radio<bool>(
              //               value: true,
              //               groupValue: hasSecApproval,
              //               activeColor: notifier.getbluewhitecolor,
              //               fillColor: MaterialStateColor.resolveWith(
              //                   (states) => notifier.getbluewhitecolor),
              //               onChanged: (value) => {
              //                 setState(
              //                   () {
              //                     hasSecApproval = value!;
              //                   },
              //                 )
              //               },
              //             ),
              //           ),
              //           Text(
              //             "yes".tr(),
              //             style: TextStyle(
              //               fontSize: 14,
              //               fontFamily: fontsemibold,
              //               color: notifier.getbluewhitecolor,
              //             ),
              //           ),
              //         ],
              //       ),
              //       Row(
              //         children: [
              //           Transform.scale(
              //             scale: 1,
              //             child: Radio<bool>(
              //               value: false,
              //               activeColor: notifier.getbluewhitecolor,
              //               fillColor: MaterialStateColor.resolveWith(
              //                   (states) => notifier.getbluewhitecolor),
              //               groupValue: hasSecApproval,
              //               onChanged: (value) => {
              //                 setState(
              //                   () {
              //                     hasSecApproval = value!;
              //                   },
              //                 )
              //               },
              //             ),
              //           ),
              //           Text(
              //             "no".tr(),
              //             style: TextStyle(
              //               fontSize: 14,
              //               fontFamily: fontsemibold,
              //               color: notifier.getbluewhitecolor,
              //             ),
              //           ),
              //         ],
              //       ),
              //     ],
              //   ),
              //   SizedBox(
              //     height: height / 50,
              //   ),
              //   if (hasSecApproval) ...[
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 15),
              //       child: Row(
              //         mainAxisAlignment: MainAxisAlignment.spaceBetween,
              //         children: [
              //           Text(
              //             "secapprovalid".tr(),
              //             style: TextStyle(
              //               fontSize: 13,
              //               fontFamily: fontsemibold,
              //               color: notifier.getbluewhitecolor,
              //             ),
              //           ),
              //         ],
              //       ),
              //     ),
              //     SizedBox(
              //       height: height / 70,
              //     ),
              //     CustomTextFormField.textField(
              //       'enterapprovalidnumber'.tr(),
              //       notifier.getbluecolor,
              //       null,
              //       notifier.getgrey,
              //       notifier.getprefixicon,
              //       notifier.getblck,
              //       notifier.getgrey,
              //       70,
              //       350,
              //       initialValue: secApprovalId,
              //       validator: (value) {
              //         if (hasSecApproval && value.isEmpty) {
              //           return "pleaseentersecapprovalid".tr();
              //         }
              //         return null;
              //       },
              //       onSaved: (value) {
              //         secApprovalId = value.trim().replaceAll(' ', '');
              //       },
              //     ),
              //   ] else ...[
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 10.0),
              //       child: TextButton(
              //         onPressed: () async {},
              //         child: Text(
              //           "pleaseapplytosec".tr(),
              //           textAlign: TextAlign.left,
              //           style: TextStyle(
              //             decoration: TextDecoration.underline,
              //             fontSize: 13,
              //             fontFamily: fontsemibold,
              //             color: notifier.getbluewhitecolor,
              //           ),
              //         ),
              //       ),
              //     ),
              //   ],
              // ],
            ],
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
    if (!hasAllRequiredDocuments) {
      message +=
          "You need to acquire all the documents in the required documents list before you can proceed.\n\n";
    }

    if (offeringType == 0 && !hasCustodianAgreement) {
      message +=
          "An asset custodian agreement is needed in this process. You need to obtain an agreement with an asset custodian to proceed.\n\n";
    }

    if (message.isNotEmpty) {
      popup(context, title: "info".tr(), message: message);
      return;
    }

    final form = _formKey.currentState;
    if (!form!.validate()) return;

    try {
      showLoader(context);
      // make initial request to the server using the
      // following credential
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;
      var marketMakingWallet = appState.activeDistributionWalletPublicKey!;
      // var newData = {...data as Map};

      // newData["assetSector"] = selectedAssetSectorId;
      // newData["assetSubSector"] = selectedAssetSubSectorId;
      // newData["assetType"] = selectedAssetTypeId;
      // newData["offeringType"] = offeringType == 1 ? 'private' : 'public';
      // newData["approvedAssetCustodianId"] = selectedAssetCustodian.length > 0
      //     ? int.parse(selectedAssetCustodian)
      //     : 1;
      // newData["marketMakingWallet"] = marketMakingWallet.publicKey;
      // newData["secApprovalIdNumber"] = secApprovalId;
      // newData["secApproval"] = hasSecApproval;
      // newData["assetCountryLocation"] = selectedCountry;

      // print('map here $newData');
      // String requestBody = jsonEncode(newData);
      print('selected asset custodian $selectedAssetCustodian');
      Map map = {
        "assetSector": selectedAssetSectorId,
        "assetSubSector": selectedAssetSubSectorId,
        "assetType": selectedAssetTypeId,
        "offeringType": offeringType == 1 ? 'private' : 'public',
        "approvedAssetCustodianId": selectedAssetCustodian.length > 0
            ? int.parse(selectedAssetCustodian)
            : 1,
        "marketMakingWallet": marketMakingWallet,
        "secApproval": hasSecApproval ? 1 : 0,
        "secApprovalIdNumber": secApprovalId,
        "assetCountryLocation": selectedCountry,
      };
      print('map here $map');
      String requestBody = jsonEncode(map);
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

      if (responseData['statusCode'] == 200) {
        if (appState.viewData != null) {
          var newData = {...appState.viewData!};
          newData["assetSector"] = selectedAssetSectorId;
          newData["assetSubSector"] = selectedAssetSubSectorId;
          newData["assetType"] = selectedAssetTypeId;
          newData["offeringType"] = offeringType == 1 ? 'private' : 'public';
          newData["approvedAssetCustodianId"] =
              selectedAssetCustodian.length > 0
                  ? int.parse(selectedAssetCustodian)
                  : 1;
          newData["marketMakingWallet"] = marketMakingWallet;
          newData["secApproval"] = hasSecApproval ? 1 : 0;
          newData["secApprovalIdNumber"] = secApprovalId;
          newData["assetCountryLocation"] = selectedCountry;
          appState.viewData = newData;
        } else {
          appState.viewData = responseData['data'];
        }

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

  String getSelectedAssetCustodianLabel(id, data) {
    var label = "";
    for (var i = 0; i < data!['assetCustodians'].length; i++) {
      if (data!['assetCustodians'][i]['id'] == id) {
        label =
            '${data!['assetCustodians'][i]['assetCustodianName']}, ${data!['assetCustodians'][i]['assetCustodianAddress']} ${data!['assetCustodians'][i]['assetCustodianCountry']}';
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
}
