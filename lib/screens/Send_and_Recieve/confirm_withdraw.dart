import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class ConfirmWithdrawal extends StatefulWidget {
  const ConfirmWithdrawal({Key? key}) : super(key: key);

  @override
  State<ConfirmWithdrawal> createState() => _ConfirmWithdrawal();
}

class _ConfirmWithdrawal extends State<ConfirmWithdrawal>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  late Account primaryWalletKeyPair;
  var viewData;
  bool isSharedWallet = false;
  var transactionInfo = {};

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
    viewData = appState.viewData![ConfirmWithdrawViewPageConfig.key];
    isSharedWallet = viewData['walletInfo']['sharedAccessEnabled'] == 1;
    transactionInfo = viewData['data'];
    print('viewData here ======> $isSharedWallet');

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
                'You are about to withdraw',
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
                            '${transactionInfo['amountSubmitted']} ${transactionInfo['currency']}',
                            style: TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w700,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold),
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
              showToAddress(),
              SizedBox(
                height: height / 50,
              ),
              showNetwork(),
              SizedBox(
                height: height / 50,
              ),
              showFee(),
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

  Widget showToAddress() {
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
                  '${transactionInfo['withdrawalAddress']}',
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

  Widget showFee() {
    return Column(
      children: [
        Text(
          'Fees',
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
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 30),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Service fee:',
                        style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.bold,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold),
                      ),
                      Text(
                        '${transactionInfo['withdrawalServiceFee']} ${transactionInfo['currency']}',
                        style: TextStyle(
                            fontSize: 15,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Network fee:',
                        style: TextStyle(
                            fontSize: 15,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold),
                      ),
                      Text(
                        '${transactionInfo['withdrawalNetworkFee']} ${transactionInfo['currency']}',
                        style: TextStyle(
                            fontSize: 15,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget showNetwork() {
    return Column(
      children: [
        Text(
          'Withdrawal Network',
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
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 20),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        '${transactionInfo['withdrawalNetworkName']} (${transactionInfo['withdrawalNetwork']})',
                        style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.bold,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold),
                      ),
                    ],
                  ),
                ],
              ),
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
        transactionInfo['commit'] = 1;
      } else {
        // sign transaction
        var signature = TrovoWalletSDK().signBase64Txn(
          appState.secretKeys[0], // the primary wallet secret key,
          transactionInfo['transaction'],
          transactionInfo['networkPassPhrase'],
        );
        transactionInfo['transactionSignature'] = signature;
      }

      String requestBody = jsonEncode(transactionInfo);

      print('requestBody: $requestBody');

      Map responseData = await makePostRequest(
        uri: isSharedWallet
            ? '/v1/shared-access/crypto/withdrawals'
            : '/v1/crypto/withdrawals',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: viewData['walletInfo']['publicKey']!,
      );

      print('responseData: $responseData');

      if (responseData['statusCode'] == 200) {
        await updateUserInfo(
          appState.activeWallet!.signer!,
          appState.secretKeys[0], // the primary wallet secret key
          appState.activeWallet!.publicKey!,
          'kenmaddy',
          appState,
        );
        if (isSharedWallet) {
          appState.viewData![SuccessViewPageConfig.key] = {
            'title': 'Withdrawal request submitted',
            'message':
                'You have successfully requested withdrawal of [${transactionInfo['amountSubmitted']} ${transactionInfo['currency']}] on network [${transactionInfo['withdrawalNetworkName']} (${transactionInfo['withdrawalNetwork']})] to address [${transactionInfo['withdrawalAddress']}]. This transaction will be completed when it gets the required number of approvals by those who have approver access on this wallet.',
            'useOnDone': true,
            'onDone': () {
              appState.currentAction = PageAction(
                state: PageState.addAll,
                pages: [BottomHomePageConfig],
              );
            },
          };
          appState.currentAction =
              PageAction(state: PageState.replace, page: SuccessViewPageConfig);
        } else {
          appState.viewData![TransactionStatusViewPageConfig.key] =
              responseData['data'];
          appState.viewData![TransactionStatusViewPageConfig.key]
              ['walletInfo'] = viewData['walletInfo'];
          appState.currentAction = PageAction(
            state: PageState.addPage,
            page: TransactionStatusViewPageConfig,
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
}
