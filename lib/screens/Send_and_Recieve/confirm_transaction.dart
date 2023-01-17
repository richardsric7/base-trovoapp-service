import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  late Map sendingWallet;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  late Account primaryWalletKeyPair;
  var viewData;
  bool isSharedWallet = false;

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
    viewData = appState.viewData![ConfirmTransactionViewPageConfig.key];
    sendingWallet = viewData['walletInfo'];
    isSharedWallet = viewData['isSharedWallet'];
    print('=====view: $viewData');
    print('=====sendingWallet: $sendingWallet');

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
                        color: notifier.getbluewhitecolor,
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
                  color: notifier.getbluewhitecolor,
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
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
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
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Container(
                            width: width / 1.3,
                            child: Text(
                              '- ${calculateFiatValue(viewData['amount'], viewData["usdPrice"], appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
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
                  color: notifier.getbluewhitecolor,
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
              if (isSharedWallet) ...[
                SizedBox(
                  height: height / 50,
                ),
                Text(
                  'Service Fee',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
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
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            SizedBox(
                              height: height / 50,
                            ),
                            myKeyValueRow("Fee: ", viewData['fee'] + '%'),
                            myKeyValueRow("Amount (Calculated): ",
                                "${viewData['feeAmount']} ${viewData['assetCode'].toString().isEmpty ? 'XBN' : viewData['assetCode']}"),
                            SizedBox(
                              height: height / 50,
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ],
              SizedBox(
                height: height / 20,
              ),
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  LanguageEn.password,
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
              SizedBox(
                height: height / 20,
              ),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  LanguageEn.authorizewithbiometrics,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  LanguageEn.authorize,
                  notifier.getbluecolor,
                  wihitecolor,
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

  Widget myKeyValueRow(String key, String value) {
    return Row(children: [
      Text(
        key,
        style: TextStyle(
            fontSize: 15,
            color: notifier.getbluewhitecolor,
            fontFamily: fontsemibold),
      ),
      Text(
        value,
        style: TextStyle(
            fontSize: 15,
            color: notifier.getbluewhitecolor,
            fontFamily: fontbody),
      ),
    ]);
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
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
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
                      color: notifier.getbluewhitecolor,
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
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
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
                          backgroundColor: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                          foregroundImage: NetworkImage(
                            viewData['destinationThumbnail'].toString(),
                          ),
                        ),
                ),
                SizedBox(
                  width: width / 70,
                ),
                Container(
                  width: width / 2,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        viewData['destination'].toString(),
                        style: TextStyle(
                          fontWeight: FontWeight.w500,
                          color: notifier.getbluewhitecolor,
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
                          color: notifier.getbluewhitecolor,
                          fontSize: 12.sp,
                          fontWeight: FontWeight.w500,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
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
            color: notifier.getbluewhitecolor,
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
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(
                      vertical: 30.0, horizontal: 15),
                  child: Container(
                    width: width / 1.3,
                    child: Text(
                      viewData['memo'],
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontSize: 17.sp,
                        fontFamily: fontbody,
                      ),
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
    try {
      showLoader(context);

      if (isSharedWallet) {
        viewData['commit'] = 1;
      } else {
        // sign transaction
        var signature = TrovoWalletSDK().signBase64Txn(
          appState.secretKeys[0], // the primary wallet secret key,
          viewData['transaction'],
          viewData['networkPassPhrase'],
        );
        viewData['transactionSignature'] = signature;
      }

      String requestBody = jsonEncode(viewData);

      // print(requestBody);

      Map responseData = await makePostRequest(
        uri: isSharedWallet ? '/v1/shared-access/payment' : '/v1/users/payment',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: sendingWallet['publicKey']!,
      );

      if (responseData['statusCode'] == 200) {
        await updateUserInfo();
        if (isSharedWallet) {
          appState.viewData![SuccessViewPageConfig.key] = {
            'title': 'Payment request submitted',
            'message':
                'You have successfully requested payment of [${viewData['amount']} ${viewData['assetCode'].toString().isEmpty ? 'XBN' : viewData['assetCode']}] from [${sendingWallet['alias']}] to [${viewData['destination']}]. This transaction will be completed when it gets the required number of approvals by those who have approver access on this wallet.',
            'useOnDone': true,
            'onDone': () {
              // if we got here through the wallets tab on dashboard
              if (viewData['rel'] == 'walletsView') {
                appState.currentAction = PageAction(
                  state: PageState.addAll,
                  pages: [
                    BottomHomePageConfig,
                    SharedWalletDetailsViewPageConfig
                  ],
                );
              } else if (viewData['rel'] == 'dashboard') {
                appState.currentAction = PageAction(
                  state: PageState.addAll,
                  pages: [BottomHomePageConfig],
                );
              } else {
                // if we got here through the shared access page
                appState.currentAction =
                    PageAction(state: PageState.addAll, pages: [
                  BottomHomePageConfig,
                  SharedAccessViewPageConfig,
                  SharedWalletInfoViewPageConfig,
                  SharedWalletDetailsViewPageConfig
                ]);
              }
            },
          };
          appState.currentAction =
              PageAction(state: PageState.replace, page: SuccessViewPageConfig);
        } else {
          appState.viewData![TransactionSuccessViewPageConfig.key] =
              responseData['data'];
          appState.viewData![TransactionSuccessViewPageConfig.key]
              ['sendingWallet'] = sendingWallet;
          appState.currentAction = PageAction(
            state: PageState.replaceAll,
            page: TransactionSuccessViewPageConfig,
          );
        }
        hideLoader(context);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }

  Future<void> updateUserInfo() async {
    Map responseData = await makeGetRequest(
      uri:
          '/v1/users/${appState.userInfo!.username!.trim().replaceAll(' ', '')}',
      signer: appState.activeWallet!.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.activeWallet!.publicKey!,
    );

    if (responseData['statusCode'] == 200) {
      await storeUserInfo(responseData['data'], appState);
    }
  }
}
