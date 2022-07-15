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

class ConfirmTransaction extends StatefulWidget {
  const ConfirmTransaction({Key? key}) : super(key: key);

  @override
  State<ConfirmTransaction> createState() => _ConfirmTransaction();
}

class _ConfirmTransaction extends State<ConfirmTransaction>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  Wallet? activeWallet;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  late Account primaryWalletKeyPair;
  var viewData;

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
    viewData = appState.viewData![ConfirmTransactionViewPageConfig.key];

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
                    LanguageEn.confirmyourtransaction,
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
                LanguageEn.youareabouttosend,
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
                            '${viewData['amount']} ${viewData['assetCode'].toString().isEmpty ? 'XBN' : viewData['assetCode']}',
                            style: TextStyle(
                                fontSize: 15,
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
                height: height / 50,
              ),
              if (viewData['memo'].toString().isNotEmpty) ...[
                showMemo(),
              ],
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
    if (viewData['destination'].toString().length == 56) {
      // destination user is not known so we display only
      // destination public key
      return Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.getaddsubwalletgrey,
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 15),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 1.5,
                  child: Text(
                    viewData['destination'].toString(),
                    style: TextStyle(
                      fontWeight: FontWeight.w500,
                      color: notifier.getbluecolor,
                      fontSize: 15.sp,
                      fontFamily: fontbody,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      );
    }

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
            Row(
              children: [
                Padding(
                  padding: EdgeInsets.fromLTRB(20, 11, 8, 10),
                  child: viewData['destinationThumbnail'].toString().isEmpty
                      ? CircleAvatar(
                          radius: 30,
                          backgroundColor: notifier.getaddsubwalletgrey,
                          foregroundImage:
                              AssetImage("assets/images/default-user.png"),
                        )
                      : CircleAvatar(
                          radius: 30,
                          backgroundColor: notifier.getaddsubwalletgrey,
                          foregroundImage: NetworkImage(
                            viewData['destinationThumbnail'].toString(),
                          ),
                        ),
                ),
                SizedBox(
                  width: width / 70,
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      viewData['destination'].toString(),
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
                      '${viewData['destinationFirstName']} ${viewData['destinationLastName']}',
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
          ],
        ),
      ),
    );
  }

  Widget showMemo() {
    return Column(
      children: [
        Text(
          LanguageEn.descriptionmemo,
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
                Padding(
                  padding: const EdgeInsets.symmetric(vertical: 30.0),
                  child: Text(
                    viewData['memo'],
                    style: TextStyle(
                      fontWeight: FontWeight.w500,
                      color: notifier.getbluecolor,
                      fontSize: 17.sp,
                      fontFamily: fontbody,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
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

  void sendDataToServer() async {
    print('sending to server....');
    showLoader(context);

    try {
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
        uri: '/v1/users/payment',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        appState.viewData![TransactionSuccessViewPageConfig.key] =
            responseData['data'];

        await updateUserInfo();
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

    appState.currentAction = PageAction(
      state: PageState.replaceAll,
      page: TransactionSuccessViewPageConfig,
    );
  }
}
