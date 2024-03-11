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
import 'package:trovo_wallet/screens/asset-tokenization/state.dart';
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
  late AssetTokenizationViewsState tokenizationState;
  late DataProvider appState;
  bool assetExisting = false;
  int assetOwnership = 0;
  int thirdPartyOwnerType = 0; // individual = 0; 1 = organization

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
    var data = tokenizationState.tokenizationData!['assetProtectionOptions'];
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
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    tokenizationState =
        Provider.of<AssetTokenizationViewsState>(context, listen: false);

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
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
              height: height / 30,
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
                    fontFamily: fontsemibold,
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
                        value: false,
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
                        value: true,
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
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
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
                    '',
                    notifier.getbluecolor,
                    notifier.getgrey,
                    notifier.getblck,
                    notifier.getgrey,
                    100.sp,
                    width / 1.12,
                    validator: (value) {
                      if (value.isEmpty) {
                        return "enterpassphraseempty".tr();
                      }
                    },
                    onSaved: (value) {},
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
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
              height: height / 70,
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
                  width / 2.7,
                  // controller: referrerController,
                  // validator: validateReferrer,
                  onSaved: (value) {},
                  keyboardtype: TextInputType.numberWithOptions(
                    decimal: true,
                    signed: true,
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
                  50.sp,
                  width / 2.5,
                  // controller: referrerController,
                  // validator: validateReferrer,
                  onSaved: (value) {},
                  keyboardtype: TextInputType.numberWithOptions(
                    decimal: true,
                    signed: true,
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
                    "ownership".tr(),
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
                        assetOwnership = 0;
                      });
                    },
                    borderColor: notifier.getbluewhitecolor,
                    foreColor: notifier.getbluewhitecolor,
                    backColor: assetOwnership == 0
                        ? notifier.getbluecolor60
                        : notifier.getwihitecolor,
                  ),
                  CheckItem(
                    "thirdparty".tr(),
                    () {
                      setState(() {
                        assetOwnership = 1;
                      });
                    },
                    borderColor: notifier.getbluewhitecolor,
                    foreColor: notifier.getbluewhitecolor,
                    backColor: assetOwnership == 1
                        ? notifier.getbluecolor60
                        : notifier.getwihitecolor,
                  )
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            if (assetOwnership == 1) ...[
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
                          thirdPartyOwnerType = 0;
                        });
                      },
                      borderColor: notifier.getbluewhitecolor,
                      foreColor: notifier.getbluewhitecolor,
                      backColor: thirdPartyOwnerType == 0
                          ? notifier.getbluecolor60
                          : notifier.getwihitecolor,
                    ),
                    CheckItem(
                      "organization".tr(),
                      () {
                        setState(() {
                          thirdPartyOwnerType = 1;
                        });
                      },
                      borderColor: notifier.getbluewhitecolor,
                      foreColor: notifier.getbluewhitecolor,
                      backColor: thirdPartyOwnerType == 1
                          ? notifier.getbluecolor60
                          : notifier.getwihitecolor,
                    )
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              if (thirdPartyOwnerType == 0) ...[
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
                        // controller: referrerController,
                        // validator: validateReferrer,
                        onSaved: (value) {},
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
                        // controller: referrerController,
                        // validator: validateReferrer,
                        onSaved: (value) {},
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
                        // controller: referrerController,
                        // validator: validateReferrer,
                        onSaved: (value) {},
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
                        // controller: referrerController,
                        // validator: validateReferrer,
                        onSaved: (value) {},
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
                    "assetcustodian".tr(),
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
                      "provideinfoaboutcustodian".tr(),
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
                    "whoiscustodian".tr(),
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
                    "assetcustodian".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                    "custodianaddress".tr(),
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
                    "address".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                    "assetmanager".tr(),
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
                      "provideassetmanagerinfo".tr(),
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
                    "whoisassetmanager".tr(),
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
                    "assetmanager".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                    "assetmanageraddress".tr(),
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
                    "address".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                    "tokenizedpercentage".tr(),
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
                    "percentagetobetokenized".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    "valueoftokenizedasset".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    "valuexpercentage".tr(),
                    style: TextStyle(
                      fontSize: 9,
                      fontFamily: fontbody,
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
                  child: Text(
                    "assetprotectioninplace".tr(),
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
                (value) {},
                getAssetProtectionOptions,
                null,
                'Insurance',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: Row(
                children: [
                  userItem(
                    'Insurance',
                    () {},
                    foreColor: notifier.getwihitecolor,
                    backColor: notifier.getbluewhitecolor,
                  ),
                  userItem(
                    'Alarm System',
                    () {},
                    foreColor: notifier.getwihitecolor,
                    backColor: notifier.getbluewhitecolor,
                  )
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
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    "companyname".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    "insurancypolicynumber".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    "insurancypolicyholder".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
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
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    "percentagevalueofinsurance".tr(),
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    width / 1.12,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            confirmLiensAndEncumbrance(),
            SizedBox(
              height: height / 30,
            ),
            Button(
              "save".tr(),
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                Navigator.of(context).pop();
              },
            ),
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
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
                value: true,
                onChanged: (bool? value) {
                  setState(() {});
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
                    color: notifier.getbluewhitecolor,
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
                value: true,
                onChanged: (bool? value) {
                  setState(() {});
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
                    color: notifier.getbluewhitecolor,
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
                value: true,
                onChanged: (bool? value) {
                  setState(() {});
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
                    color: notifier.getbluewhitecolor,
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
                value: true,
                onChanged: (bool? value) {
                  setState(() {});
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
                    color: notifier.getbluewhitecolor,
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
}
