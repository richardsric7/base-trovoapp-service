import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class ConfirmSwap extends StatefulWidget {
  const ConfirmSwap({Key? key}) : super(key: key);

  @override
  State<ConfirmSwap> createState() => _ConfirmSwap();
}

class _ConfirmSwap extends State<ConfirmSwap> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  Wallet? activeWallet;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  late Account primaryWalletKeyPair;
  var viewData;
  var sourceAmount;
  var swappedEstimate;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    activeWallet = appState.activeWallet;
    viewData = appState.viewData![ConfirmSwapViewPageConfig.key];
    sourceAmount = double.parse(viewData['sourceAmount']).toStringAsFixed(4);
    swappedEstimate =
        double.parse(viewData['swappedEstimate']).toStringAsFixed(4);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.confirmswap,
                    style: TextStyle(
                        fontSize: 20.sp,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: height / 20,
              ),
              Text(
                LanguageEn.youareabouttoswap,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.getaddsubwalletgrey,
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Column(
                        children: [
                          SizedBox(
                            height: height / 50,
                          ),
                          Text(
                            '${sourceAmount} ${viewData['sourceAssetCode'].toString().isEmpty ? 'XBN' : viewData['sourceAssetCode']}',
                            style: TextStyle(
                                fontSize: 19,
                                fontWeight: FontWeight.w700,
                                color: notifier.getbluecolor,
                                fontFamily: fontsemibold),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Container(
                            width: width / 1.3,
                            child: Text(
                              '- 3400 NGN',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluecolor,
                                fontFamily: fontbody,
                              ),
                            ),
                          ),
                          SizedBox(
                            height: height / 50.0,
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Text(
                LanguageEn.to,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              showAddressInfo(),
              SizedBox(
                height: height / 20,
              ),
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  LanguageEn.password,
                  notifier.getbluecolor,
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
              SizedBox(
                height: height / 20,
              ),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  LanguageEn.authorizewithbiometrics,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  LanguageEn.authorize,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: handleAuthorization,
                ),
              ],
              SizedBox(
                height: height / 20,
              ),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Widget showAddressInfo() {
    return Padding(
      // destination user is known so we display the user
      // information.
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
        ),
        child: Column(
          children: [
            SizedBox(
              height: height / 50,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(
                  width: width / 70,
                ),
                Column(
                  children: [
                    Text(
                      '${swappedEstimate} ${viewData['destinationAssetCode'].toString().isEmpty ? 'XBN' : viewData['destinationAssetCode']}',
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluecolor,
                        fontSize: 19.sp,
                        fontFamily: fontbody,
                      ),
                    ),
                    SizedBox(
                      height: 5,
                    ),
                    Text(
                      '+ 3400 NGN',
                      style: TextStyle(
                        color: notifier.getbluecolor,
                        fontSize: 12.sp,
                        fontWeight: FontWeight.w500,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                )
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
          ],
        ),
      ),
    );
  }

  void handleAuthorization() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer();
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
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

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  sendDataToServer() async {
    print('sending to server....');

    try {
      showLoader(context);
      // sign transaction
      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0], // the primary wallet secret key,
        viewData['transaction'],
        viewData['networkPassPhrase'],
      );
      viewData['transactionSignature'] = signature;

      String requestBody = jsonEncode(viewData);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/swap',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        await updateUserInfo();
        appState.viewData![SwapSuccessViewPageConfig.key] =
            responseData['data'];
        appState.currentAction = PageAction(
          state: PageState.replaceAll,
          page: SwapSuccessViewPageConfig,
        );
        hideLoader(context);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  Future<void> updateUserInfo() async {
    Map responseData = await makeGetRequest(
      uri:
          '/v1/users/${appState.userInfo!.username!.trim().replaceAll(' ', '')}',
      signer: activeWallet!.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: activeWallet!.publicKey!,
    );

    print('secretkey: ${appState.secretKeys[0]}');

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200) {
      await storeUserInfo(responseData['data']);
    }
  }

  Future<void> storeUserInfo(userInfoMap) async {
    print('userInfoMap: ${userInfoMap['userData']}');
    var userInfo = userInfoMap['userData'] ?? {};
    var assetBalances = userInfoMap['assetBalances'] ?? {};
    var nfts = userInfoMap['nfts'] ?? {};
    var thirdPartyWalletAccess = userInfoMap['thirdPartyWalletAccess'] ?? [];
    var defaultAssets = userInfoMap['defaultAssets'] ?? [];

    await StoreData().storeInsertData('userInfo', userInfo);
    await StoreData().storeInsertData('assetBalances', assetBalances);
    await StoreData().storeInsertData('nftBalances', nfts);
    await StoreData()
        .storeInsertData('thirdPartyWalletAccess', thirdPartyWalletAccess);
    await StoreData().storeInsertData('defaultAssets', defaultAssets);

    // save useInfo to appstate
    appState.setUser = UserInfo().deserializeJson(userInfo);
    appState.setNFTs = nfts;
    appState.setassetBalances = assetBalances;
    print('stored new user data.................');
  }
}
