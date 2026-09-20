import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart' hide Trans;
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/patronInfo.dart';
import 'package:trovo_app/models/patronTier.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/cache.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:local_auth/error_codes.dart' as auth_error;
import 'package:trovo_app/widgets/utilities.dart';

import '../../custom_bloc_observer/fonts.dart';
import '../../models/wallet.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AuthorizeSubscription extends StatefulWidget {
  const AuthorizeSubscription({Key? key}) : super(key: key);

  @override
  State<AuthorizeSubscription> createState() => _AuthorizeSubscriptionState();
}

class _AuthorizeSubscriptionState extends State<AuthorizeSubscription> {
  late ColorNotifier notifier;
  late DataProvider appState;
  final Authenticator _authenticator = Authenticator();
  late PatronInfo patronInfo;
  late PatronTier patronTier;
  late List<Asset> paymentAssets;
  late Wallet wallet;
  late List<Asset> claimedAssets;
  String password = '';
  final formKey = GlobalKey<FormState>();
  String? selectedAsset;

  List<DropdownMenuItem<String>> get getAssetDropdownItems {
    List<DropdownMenuItem<String>> menuItems = [];
    for (var asset in paymentAssets) {
      menuItems.add(
        DropdownMenuItem(
          child: Text(
            getAssetCode(asset.assetCode),
            overflow: TextOverflow.visible,
          ),
          value:
              '${getAssetCode(asset.assetCode)}|${getAssetIssuer(asset.assetIssuer)}',
        ),
      );
    }
    return menuItems;
  }

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
    wallet = appState.primaryWallet;
    claimedAssets = wallet.claimedAssets!;
    patronInfo = appState.viewData!['patronInfo'];
    patronTier = appState.viewData!['selectedTier'];
    paymentAssets = appState.viewData!['paymentAssets'];
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "trovopatron".tr(),
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(height: height / 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "youhavechosen".tr(),
                    style: TextStyle(
                      fontSize: 22,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                ],
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    '${patronInfo.patronPackage.capitalizeFirst} ${"patron".tr()}',
                    style: TextStyle(
                      fontSize: 22,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                  Image.asset(patronInfo.logo, height: 30),
                ],
              ),
              SizedBox(height: 20),
              planItem(),
              if (patronTier.isExpirable) ...[
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.all(30.0),
                  child: Text(
                    "nobenefitsoncesubexpires".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 12,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                ),
              ],
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "paymentasset".tr(),
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.6,
                child: DropdownButtonFormField<String>(
                  isExpanded: true,
                  value: selectedAsset,
                  dropdownColor: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                  decoration: InputDecoration(
                    contentPadding: EdgeInsets.symmetric(
                      vertical: 0,
                      horizontal: 20,
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderSide: BorderSide.none,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    border: OutlineInputBorder(
                      borderSide: BorderSide.none,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    filled: true,
                    fillColor: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                    errorStyle: TextStyle(
                      fontFamily: fontbody,
                      fontSize: 12,
                      overflow: TextOverflow.visible,
                    ),
                  ),
                  hint: Container(
                    width: 150, //and here
                    child: Text(
                      "chooseasset".tr(),
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
                      textAlign: TextAlign.end,
                    ),
                  ),
                  icon: Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                  elevation: 0,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    fontWeight: FontWeight.w500,
                  ),
                  onChanged: (newValue) {
                    selectedAsset = newValue;
                    setState(() {});
                  },
                  items: getAssetDropdownItems,
                ),
              ),
              SizedBox(height: height / 20),
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  "password".tr(),
                  notifier.getbluewhitecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  validator: validatePassword,
                  onChanged: (value) {
                    setState(() {
                      password = value!.trim().replaceAll(' ', '');
                    });
                  },
                ),
              ),
              SizedBox(height: height / 20),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  "authorizewithbiometrics".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  "authorize".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: handleAuthorization,
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
      ),
    );
  }

  Widget planItem() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Column(
              children: [
                Container(
                  height: 10,
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.only(
                      topLeft: Radius.circular(15.0),
                      topRight: Radius.circular(15.0),
                    ),
                    color: getColor(patronInfo.patronPackage),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 10.0,
                    vertical: 15.0,
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: width / 1.2,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            SizedBox(width: width / 50),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  "subscriptionplan".tr(),
                                  style: TextStyle(
                                    fontSize: 20,
                                    fontWeight: FontWeight.w400,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                              ],
                            ),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  patronTier.tier.capitalizeFirst!,
                                  style: TextStyle(
                                    fontSize: 20,
                                    fontWeight: FontWeight.w400,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 50),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  "expirydate".tr(),
                                  style: TextStyle(
                                    fontSize: 20,
                                    fontWeight: FontWeight.w400,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                              ],
                            ),
                            if (patronTier.isExpirable) ...[
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    DateFormat('MMMM dd, yyyy').format(
                                      DateTime.now().add(
                                        patronTier.tier == 'Annual'
                                            ? Duration(days: 31)
                                            : Duration(days: 365),
                                      ),
                                    ),
                                    style: TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  ),
                                ],
                              ),
                            ] else ...[
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    "hasnoexpirydate".tr(),
                                    style: TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  ),
                                ],
                              ),
                            ],
                            SizedBox(height: 20),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Color getColor(String patronPlan) {
    if (patronPlan.toLowerCase() == 'gold')
      return notifier.getgoldcolor;
    else if (patronPlan.toLowerCase() == 'platinum')
      return notifier.getplatinumcolor;
    else
      return notifier.getdiamondcolor;
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return "pleaseenteryourpassword".tr();

    if (value.length < 6) return "use6charsormoreforpassword".tr();

    return null;
  }

  void handleAuthorization() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer();
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer();
        // aparently we need the code below to make the
        // screen updata to show loader
        // after authorizing with biometrics
        setState(() {});
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  void sendDataToServer() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        'patronMembershipGradeId': patronTier.id,
        'paymentAssetCode': selectedAsset?.split('|')[0],
        'paymentAssetIssuer': selectedAsset?.split('|')[0] == 'ETH'
            ? ''
            : selectedAsset?.split('|')[1],
      };
      String requestBody = jsonEncode(map);

      Map responseData = await makePostRequest(
        uri: '/v1/patron',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        address: appState.primaryWallet.address!,
      );

      if (responseData['statusCode'] == 202) {
        completeRequest(responseData['data']);
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
      hideLoader(context);
    }
  }

  void completeRequest(responseBody) async {
    try {
      showLoader(context);

      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0],
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      responseBody['transactionSignature'] = signature;

      String requestBody = jsonEncode(responseBody);

      Map responseData = await makePostRequest(
        uri: '/v1/patron',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        address: appState.primaryWallet.address!,
      );

      if (responseData['statusCode'] == 200) {
        await updateUserInfo(
          appState.primaryWallet.signer!,
          appState.secretKeys[0],
          appState.primaryWallet.address!,
          appState.userInfo!.username,
          appState,
        );
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': "success".tr(),
          'message': 'Subscription successful!',
          'useOnDone': true,
          'onDone': () {
            appState.currentAction =
                appState.returnView ??
                PageAction(
                  state: PageState.addAll,
                  pages: [BottomHomePageConfig],
                );
          },
        };
        appState.currentAction = PageAction(
          state: PageState.replaceAll,
          page: SuccessViewPageConfig,
        );
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }
}
