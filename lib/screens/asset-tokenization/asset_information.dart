import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  late bool assetExisting;
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
  late List<String> assetProtectionInPlace;
  late String insuranceCompanyName;
  late String insurancePolicyNumber;
  late String insurancePolicyHolder;
  late double percentageValueOfInsurance;
  late bool freeOfLiensAndEncumbrances;
  bool freeOfMortgages = false;
  bool freeOfLoans = false;
  bool freeOfDisputes = false;
  bool formHasError = false;
  late dynamic data = {};
  final valueOfAssetController = TextEditingController();
  final miscCostOfAssetController = TextEditingController();
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
    var data = appState.tokenizationData!['assetProtectionOptions'];
    for (var i = 0; i < data.length; i++) {
      assetProtectionOptions.add(
        DropdownMenuItem(
          child: Text(
            data![i]['id'],
            overflow: TextOverflow.ellipsis,
          ),
          value: data![i]['id'],
        ),
      );
    }
    return assetProtectionOptions;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    inspect(appState.viewData);
    data = appState.viewData;
    inspect(data);

    assetExisting = data!['assetAlreadyExists'] == 1;
    assetOwnership =
        data['ownershipType'] != null && data['ownershipType'].isNotEmpty
            ? data['ownershipType']
            : 'DIRECT';
    thirdPartyOwnerType =
        data['ownershipKind'] != null && data['ownershipKind'].isNotEmpty
            ? data['ownershipKind']
            : 'INDIVIDUAL';
    assetName = data['assetName'] ?? "";
    assetDescription = data['assetDescription'] ?? "";
    assetPhysicalAddress = data['assetPhysicalAddress'] ?? "";
    latitude = double.tryParse(data['assetLatitude'].toString()) ?? 0;
    longitude = double.tryParse(data['assetLongitude'].toString()) ?? 0;
    nameOfOwner = data['assetOwnerName'] ?? "";
    addressOfOwner = data['assetOwnerAddress'] ?? "";
    currentValueOfAsset =
        double.tryParse(data['assetCurrentValue'].toString()) ?? 0;
    assetMiscCost =
        double.tryParse(data['assetMscCostOutisdeOfValuation'].toString()) ?? 0;
    valueOfTokenizedAsset =
        double.tryParse(data['valueOfTokenizedAsset'].toString()) ?? 0;
    assetProtectionInPlace = data['protectionMethods'] == null ||
            data['protectionMethods'].toString().isEmpty
        ? []
        : data['protectionMethods'].toString().split(',');
    insuranceCompanyName = data['insuranceCompanyName'] ?? "";
    insurancePolicyNumber = data['insurancePolicyNumber'] ?? "";
    insurancePolicyHolder = data['insurancePolicyHolder'] ?? "";
    percentageValueOfInsurance =
        double.tryParse(data['percentageValueOfInsurance'].toString()) ?? 0;
    freeOfLiensAndEncumbrances = data['IsFreeFromLiensAndEncumbrances'] == 1;

    valueOfAssetController.text = currentValueOfAsset == 0
        ? ''
        : formatNumberForInput(currentValueOfAsset);
    miscCostOfAssetController.text =
        assetMiscCost == 0 ? '' : formatNumberForInput(assetMiscCost);
    percentValueOfInsuranceController.text = percentageValueOfInsurance == 0
        ? ''
        : percentageValueOfInsurance.toString();

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
              SizedBox(
                height: height / 50,
              ),
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
              SizedBox(
                height: height / 70,
              ),
              Container(
                width: width,
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    "selectwhatappliestoasset".tr(),
                    textAlign: TextAlign.left,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ),
              Column(
                children: [
                  Row(
                    children: [
                      Transform.scale(
                        scale: 1,
                        child: Radio<bool>(
                          value: true,
                          activeColor: notifier.getbluewhitecolor,
                          fillColor: MaterialStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor),
                          groupValue: assetExisting,
                          onChanged: (value) => {
                            setState(
                              () {
                                assetExisting = value!;
                              },
                            )
                          },
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
                  Row(
                    children: [
                      Transform.scale(
                        scale: 1,
                        child: Radio<bool>(
                          value: false,
                          groupValue: assetExisting,
                          activeColor: notifier.getbluewhitecolor,
                          fillColor: MaterialStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor),
                          onChanged: (value) => {
                            setState(
                              () {
                                assetExisting = value!;
                              },
                            )
                          },
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
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetdescription".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
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
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "enterassetphysicaladdress".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
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
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "entergooglemapcords".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  CustomTextFormField.textField(
                    "latitude".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    50.sp,
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
                    keyboardtype:
                        TextInputType.numberWithOptions(decimal: true),
                  ),
                  CustomTextFormField.textField(
                    "longitude".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    50.sp,
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
                    keyboardtype:
                        TextInputType.numberWithOptions(decimal: true),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Container(
                width: width,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.start,
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: TextButton(
                        onPressed: () => findCoordinatesPopup(context),
                        style: ButtonStyle(
                          padding: MaterialStateProperty.all(EdgeInsets.zero),
                          minimumSize: MaterialStateProperty.all(Size.zero),
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                        ),
                        child: Text(
                          "howtofindcoordinates".tr(),
                          style: TextStyle(
                            decoration: TextDecoration.underline,
                            fontSize: 13,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetownership".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    CheckItem(
                      "directownership".tr(),
                      () {
                        setState(() {
                          assetOwnership = 'DIRECT';
                        });
                      },
                      borderColor: notifier.getbluewhitecolor,
                      foreColor: notifier.getbluewhitecolor,
                      backColor: assetOwnership == 'DIRECT'
                          ? notifier.getbluecolor60
                          : notifier.getwihitecolor,
                    ),
                    CheckItem(
                      "thirdparty".tr(),
                      () {
                        setState(() {
                          assetOwnership = 'THIRD-PARTY';
                        });
                      },
                      borderColor: notifier.getbluewhitecolor,
                      foreColor: notifier.getbluewhitecolor,
                      backColor: assetOwnership == 'THIRD-PARTY'
                          ? notifier.getbluecolor60
                          : notifier.getwihitecolor,
                    )
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              if (assetOwnership == 'THIRD-PARTY') ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Container(
                        width: width / 1.17,
                        child: Text(
                          "whatbestdescribesthirdparty".tr(),
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 70,
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 15.0),
                  child: Row(
                    children: [
                      CheckItem(
                        "individual".tr(),
                        () {
                          setState(() {
                            thirdPartyOwnerType = 'INDIVIDUAL';
                          });
                        },
                        borderColor: notifier.getbluewhitecolor,
                        foreColor: notifier.getbluewhitecolor,
                        backColor: thirdPartyOwnerType == 'INDIVIDUAL'
                            ? notifier.getbluecolor60
                            : notifier.getwihitecolor,
                      ),
                      CheckItem(
                        "organization".tr(),
                        () {
                          setState(() {
                            thirdPartyOwnerType = 'CORPORATE';
                          });
                        },
                        borderColor: notifier.getbluewhitecolor,
                        foreColor: notifier.getbluewhitecolor,
                        backColor: thirdPartyOwnerType == 'CORPORATE'
                            ? notifier.getbluecolor60
                            : notifier.getwihitecolor,
                      )
                    ],
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                if (thirdPartyOwnerType == 'INDIVIDUAL') ...[
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "nameofowner".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          "nameofowner".tr(),
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          initialValue: nameOfOwner,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {
                              nameOfOwner = value!;
                            });
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "addressofowner".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          "addressofowner".tr(),
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          initialValue: addressOfOwner,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {
                              addressOfOwner = value!;
                            });
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                ] else ...[
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "nameoforg".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          "nameoforg".tr(),
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          initialValue: nameOfOwner,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {
                              nameOfOwner = value;
                            });
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          "addressoforg".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          "addressoforg".tr(),
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          initialValue: addressOfOwner,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {
                              addressOfOwner = value;
                            });
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                ],
              ],
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetvalue".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
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
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetcurrentvalue".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
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
                      70.sp,
                      width / 1.12,
                      onChanged: (value) {
                        setState(() {
                          if (value.toString().isEmpty) {
                            currentValueOfAsset = 0;
                            valueOfTokenizedAsset = 0;
                            return;
                          }

                          currentValueOfAsset = double.parse(
                              value!.toString().replaceAll(',', ''));
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
                      controller: valueOfAssetController,
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
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
              SizedBox(
                height: height / 50,
              ),
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
                      70.sp,
                      width / 1.12,
                      onChanged: (value) {
                        setState(() {
                          if (value.toString().isEmpty) {
                            assetMiscCost = 0;
                            return;
                          }

                          assetMiscCost = double.parse(
                              value!.toString().replaceAll(',', ''));
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
                      controller: miscCostOfAssetController,
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "valueoftokenizedasset".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      height: 55.sp,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 10),
                            child: Text(
                              formatNumberForInput(valueOfTokenizedAsset),
                              style: TextStyle(fontSize: 15),
                            ),
                          ),
                          const SizedBox(height: 2),
                        ],
                      ),
                    ),
                  )
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetprotectioninplace".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetprotectioninplace2".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Row(
                children: [
                  Container(
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "selectprotectionoption".tr(),
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
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      if (!assetProtectionInPlace.contains(value.toString())) {
                        assetProtectionInPlace.add(value.toString());
                      }
                    });
                  },
                  getAssetProtectionOptions,
                  null,
                  'Select asset protection',
                  context,
                  null,
                  validator: (value) {
                    if (assetProtectionInPlace.isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: Container(
                  width: width,
                  child: Wrap(
                    alignment: WrapAlignment.start,
                    children: [
                      for (var item in assetProtectionInPlace) ...[
                        userItem(
                          item,
                          () {
                            setState(() {
                              assetProtectionInPlace
                                  .removeWhere((element) => element == item);
                            });
                          },
                          foreColor: notifier.getwihitecolor,
                          backColor: notifier.getbluewhitecolor,
                        ),
                      ],
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         "insurancecompanyname".tr(),
              //         style: TextStyle(
              //           fontSize: 12,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: CustomTextFormField.textField(
              //         "companyname".tr(),
              //         notifier.getbluecolor,
              //         null,
              //         notifier.getgrey,
              //         null,
              //         notifier.getblck,
              //         notifier.getgrey,
              //         70.sp,
              //         width / 1.12,
              //         initialValue: insuranceCompanyName,
              //         validator: (value) {
              //           if (value.isEmpty) {
              //             return "fieldcannotbeempty".tr();
              //           }
              //           return null;
              //         },
              //         onSaved: (value) {
              //           setState(() {
              //             insuranceCompanyName = value!;
              //           });
              //         },
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         "insurancypolicynumber".tr(),
              //         style: TextStyle(
              //           fontSize: 12,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: CustomTextFormField.textField(
              //         "insurancypolicynumber".tr(),
              //         notifier.getbluecolor,
              //         null,
              //         notifier.getgrey,
              //         null,
              //         notifier.getblck,
              //         notifier.getgrey,
              //         70.sp,
              //         width / 1.12,
              //         initialValue: insurancePolicyNumber,
              //         validator: (value) {
              //           if (value.isEmpty) {
              //             return "fieldcannotbeempty".tr();
              //           }
              //           return null;
              //         },
              //         onSaved: (value) {
              //           setState(() {
              //             insurancePolicyNumber = value!;
              //           });
              //         },
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         "insurancypolicyholder".tr(),
              //         style: TextStyle(
              //           fontSize: 12,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: CustomTextFormField.textField(
              //         "insurancypolicyholder".tr(),
              //         notifier.getbluecolor,
              //         null,
              //         notifier.getgrey,
              //         null,
              //         notifier.getblck,
              //         notifier.getgrey,
              //         70.sp,
              //         width / 1.12,
              //         initialValue: insurancePolicyHolder,
              //         validator: (value) {
              //           if (value.isEmpty) {
              //             return "fieldcannotbeempty".tr();
              //           }
              //           return null;
              //         },
              //         onSaved: (value) {
              //           setState(() {
              //             insurancePolicyHolder = value!;
              //           });
              //         },
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         "percentagevalueofinsurance".tr(),
              //         style: TextStyle(
              //           fontSize: 12,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: CustomTextFormField.textField(
              //         "percentagevalueofinsurance".tr(),
              //         notifier.getbluecolor,
              //         null,
              //         notifier.getgrey,
              //         null,
              //         notifier.getblck,
              //         notifier.getgrey,
              //         70.sp,
              //         width / 1.12,
              //         validator: (value) {
              //           if (value.isEmpty) {
              //             return "enterassetdescription".tr();
              //           }
              //           return null;
              //         },
              //         onSaved: (value) {
              //           setState(() {
              //             percentageValueOfInsurance = double.parse(value);
              //           });
              //         },
              //         autoFormatNumber: true,
              //         controller: percentValueOfInsuranceController,
              //         keyboardtype:
              //             TextInputType.numberWithOptions(decimal: true),
              //       ),
              //     ),
              //   ],
              // ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "comprehensiveinsurance".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            Row(
                              children: [
                                SizedBox(
                                  width: 15,
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    SizedBox(
                                      height: 10,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "insurancecompanyname".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "companyname".tr(),
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insuranceCompanyName,
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
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "insurancypolicynumber".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "insurancypolicynumber".tr(),
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insurancePolicyNumber,
                                            validator: (value) {
                                              if (value.isEmpty) {
                                                return "fieldcannotbeempty"
                                                    .tr();
                                              }
                                              return null;
                                            },
                                            onSaved: (value) {
                                              setState(() {
                                                insurancePolicyNumber = value!;
                                              });
                                            },
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "insurancypolicyholder".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "insurancypolicyholder".tr(),
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insurancePolicyHolder,
                                            validator: (value) {
                                              if (value.isEmpty) {
                                                return "fieldcannotbeempty"
                                                    .tr();
                                              }
                                              return null;
                                            },
                                            onSaved: (value) {
                                              setState(() {
                                                insurancePolicyHolder = value!;
                                              });
                                            },
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "percentagevalueofinsurance".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "percentagevalueofinsurance".tr(),
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            validator: (value) {
                                              if (value.isEmpty) {
                                                return "enterassetdescription"
                                                    .tr();
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
                                                    decimal: true),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ],
                                )
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "revenueguarantees".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            CheckboxItem(
                              label: "performancebond".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            CheckboxItem(
                              label: "slas".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "ppps".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            CheckboxItem(
                              label: "hedginginstruments".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            CheckboxItem(
                              label: "completionguarantees".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "independentmonitoring".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            Row(
                              children: [
                                SizedBox(
                                  width: 15,
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    SizedBox(
                                      height: 10,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Container(
                                            width: 260,
                                            child: Text(
                                              "listmonitoringoperators".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "eg John Doe, Anna Harry",
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insuranceCompanyName,
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
                                  ],
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "sustainabilitycertifications".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "communityengagementplans".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "accesscontrol".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "surveillancesystems".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "onsitesecuritypersonnel".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "perimetersecurity".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "criticalinfrastructureprotections".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "legaladvisor".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            Row(
                              children: [
                                SizedBox(
                                  width: 15,
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    SizedBox(
                                      height: 10,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "nameoflegaladvisor".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "",
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insuranceCompanyName,
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
                                  ],
                                ),
                              ],
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "financialadvisor".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            Row(
                              children: [
                                SizedBox(
                                  width: 15,
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    SizedBox(
                                      height: 10,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "nameoffinancialadvisor".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "",
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insuranceCompanyName,
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
                                  ],
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "other".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            Row(
                              children: [
                                SizedBox(
                                  width: 15,
                                ),
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    SizedBox(
                                      height: 10,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            "enterotherprotectionsinplace".tr(),
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                    Row(
                                      children: [
                                        Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: CustomTextFormField.textField(
                                            "",
                                            notifier.getbluecolor,
                                            null,
                                            notifier.getgrey,
                                            null,
                                            notifier.getblck,
                                            notifier.getgrey,
                                            60,
                                            267,
                                            initialValue: insuranceCompanyName,
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
                                  ],
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  )
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetverificationundertaking".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmfreeofloans".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmfreeofcolateral".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmfreeofthirdparties".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmfreeoflegaldisputes".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmenvcompliance".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmhaspermits".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmnooutstandingpayments".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmnohiddenliabilities".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmassetinsured".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmvalueiscurrent".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 10),
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
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmassetstructuralysound".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
                                  setState(() {
                                    formHasError = true;
                                  });
                                  return '';
                                }

                                return null;
                              },
                            ),
                            SizedBox(
                              height: 10,
                            ),
                            CheckboxItem(
                              label: "confirmnoexistingleaseagreements".tr(),
                              onChanged: (value) {},
                              validator: (value) {
                                if (!freeOfLiensAndEncumbrances) {
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
                  )
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
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
                    bottom: MediaQuery.of(context).viewInsets.bottom),
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
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;
      inspect(data);
      var newData = {...data as Map};

      newData['assetAlreadyExists'] = assetExisting ? 1 : 0;
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
      newData['valueOfTokenizedAsset'] = valueOfTokenizedAsset;
      newData['protectionMethods'] = assetProtectionInPlace.join(',');
      newData['insuranceCompanyName'] = insuranceCompanyName;
      newData['insurancePolicyNumber'] = insurancePolicyNumber;
      newData['insurancePolicyHolder'] = insurancePolicyHolder;
      newData['percentageValueOfInsurance'] = percentageValueOfInsurance;
      newData['IsFreeFromLiensAndEncumbrances'] =
          freeOfLiensAndEncumbrances ? 1 : 0;

      String requestBody = jsonEncode(newData);
      print('requestBody =======> ');
      inspect(newData);
      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: mintingWalletPublicKey,
      );

      print('responseData ${responseData['data']}');
      inspect(responseData['data']);

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
        Navigator.of(context).pop();
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
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
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        print('success');
        appState.viewData = responseData['data'];
        await inspect(appState.viewData);
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Widget confirmLiensAndEncumbrance() {
    return Column(
      children: [
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                "financialencumbrances".tr(),
                style: TextStyle(
                  fontSize: 12,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Transform.scale(
              scale: 1.sp,
              child: FormField(
                builder: (state) {
                  return Checkbox(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(5.sp),
                      ),
                    ),
                    activeColor: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                    side: BorderSide(
                      color: notifier.isDark
                          ? notifier.getbluecolor50
                          : notifier.getbluecolor90,
                    ),
                    value: freeOfLiensAndEncumbrances,
                    onChanged: (bool? value) {
                      setState(() {
                        freeOfLiensAndEncumbrances = value!;
                      });
                    },
                  );
                },
                validator: (value) {
                  if (!freeOfLiensAndEncumbrances) {
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
                "confirmfreeofencumbrances".tr(),
                overflow: TextOverflow.visible,
                style: TextStyle(
                    fontSize: 15,
                    color: formHasError && !freeOfLiensAndEncumbrances
                        ? Colors.red
                        : notifier.getbluewhitecolor,
                    fontFamily: fontbody),
              ),
            ),
          ],
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Transform.scale(
              scale: 1.sp,
              child: FormField(
                builder: (state) {
                  return Checkbox(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(5.sp),
                      ),
                    ),
                    activeColor: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                    side: BorderSide(
                      color: notifier.isDark
                          ? notifier.getbluecolor50
                          : notifier.getbluecolor90,
                    ),
                    value: freeOfMortgages,
                    onChanged: (bool? value) {
                      setState(() {
                        freeOfMortgages = value!;
                      });
                    },
                  );
                },
                validator: (value) {
                  if (!freeOfMortgages) {
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
                "confirmfreeofmortgages".tr(),
                overflow: TextOverflow.visible,
                style: TextStyle(
                    fontSize: 15,
                    color: formHasError && !freeOfMortgages
                        ? Colors.red
                        : notifier.getbluewhitecolor,
                    fontFamily: fontbody),
              ),
            ),
          ],
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Transform.scale(
              scale: 1.sp,
              child: FormField(
                builder: (state) {
                  return Checkbox(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(5.sp),
                      ),
                    ),
                    activeColor: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                    side: BorderSide(
                      color: notifier.isDark
                          ? notifier.getbluecolor50
                          : notifier.getbluecolor90,
                    ),
                    value: freeOfLoans,
                    onChanged: (bool? value) {
                      setState(() {
                        freeOfLoans = value!;
                      });
                    },
                  );
                },
                validator: (value) {
                  if (!freeOfLoans) {
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
                "confirmfreeofloans".tr(),
                overflow: TextOverflow.visible,
                style: TextStyle(
                    fontSize: 15,
                    color: formHasError && !freeOfLoans
                        ? Colors.red
                        : notifier.getbluewhitecolor,
                    fontFamily: fontbody),
              ),
            ),
          ],
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Transform.scale(
              scale: 1.sp,
              child: FormField(
                builder: (state) {
                  return Checkbox(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(5.sp),
                      ),
                    ),
                    activeColor: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                    side: BorderSide(
                      color: notifier.isDark
                          ? notifier.getbluecolor50
                          : notifier.getbluecolor90,
                    ),
                    value: freeOfDisputes,
                    onChanged: (bool? value) {
                      setState(() {
                        freeOfDisputes = value!;
                      });
                    },
                  );
                },
                validator: (value) {
                  if (!freeOfDisputes) {
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
                "confirmfreeofdisputes".tr(),
                overflow: TextOverflow.visible,
                style: TextStyle(
                    fontSize: 15,
                    color: formHasError && !freeOfDisputes
                        ? Colors.red
                        : notifier.getbluewhitecolor,
                    fontFamily: fontbody),
              ),
            ),
          ],
        ),
      ],
    );
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

  Widget CheckboxItem({
    required String label,
    required void Function(bool?) onChanged,
    required String? Function(Object?) validator,
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
                    borderRadius: BorderRadius.all(
                      Radius.circular(5.sp),
                    ),
                  ),
                  activeColor: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor90,
                  side: BorderSide(
                    color: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                  ),
                  value: freeOfLiensAndEncumbrances,
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
                color: formHasError && !freeOfLiensAndEncumbrances
                    ? Colors.red
                    : notifier.getbluewhitecolor,
                fontFamily: fontbody),
          ),
        ),
      ],
    );
  }
}
