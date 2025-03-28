import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:local_auth/error_codes.dart' as auth_error;
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ConfirmBuy extends StatefulWidget {
  const ConfirmBuy({Key? key}) : super(key: key);

  @override
  State<ConfirmBuy> createState() => _ConfirmBuy();
}

class _ConfirmBuy extends State<ConfirmBuy> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  double amount = 0;
  double quantity = 0;
  String password = '';
  final formKey = GlobalKey<FormState>();
  late TokenizedAsset tokenizedAsset;
  final Authenticator _authenticator = Authenticator();

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
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    tokenizedAsset = appState.tokenizedAsset!;
    amount = appState.viewData!['amount'];
    quantity = appState.viewData!['quantity'];
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              'Buy ${tokenizedAsset.assetName}',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(
              height: height / 30,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Card(
                shadowColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: notifier.isDark
                    ? notifier.getbluecolor90
                    : notifier.getaddsubwalletgrey,
                child: Center(
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        "You’re buying",
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        '${formatNumber(quantity)} Tokens',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        'of [${tokenizedAsset.assetName}] Asset',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 30,
                      ),
                      Text(
                        'Amount',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        '${formatNumberShort(amount)} ${tokenizedAsset.assetQuoteCurrency}',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      // SizedBox(
                      //   height: height / 70,
                      // ),
                      // Text(
                      //   '\$0.20',
                      //   textAlign: TextAlign.center,
                      //   style: TextStyle(
                      //     fontSize: 14,
                      //     height: 1.4,
                      //     fontFamily: fontbody,
                      //     color: notifier.getbluewhitecolor,
                      //   ),
                      // ),
                      SizedBox(
                        height: height / 30,
                      ),
                      Text(
                        'Pay with',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            '${appState.activeWallet!.alias}',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 14,
                              height: 1.4,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(
                        height: height / 30,
                      ),
                    ],
                  ),
                ),
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            SizedBox(height: 10),
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
                validator: (String? value) {
                  if (value!.isEmpty) return 'Enter your password';

                  if (value.length < 6)
                    return 'Use 6 characters or more for your password';

                  return null;
                },
                onChanged: (value) {
                  setState(() {
                    password = value!.trim().replaceAll(' ', '');
                  });
                },
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
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
                onTap: () {
                  if (!formKey.currentState!.validate()) {
                    return;
                  }

                  if (password == appState.password!) {
                    buyTokenizedAsset();
                  } else {
                    popup(context,
                        title: "oops".tr(), message: "invalidpassword".tr());
                  }
                },
              ),
            ],
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        buyTokenizedAsset();
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

  buyTokenizedAsset() async {
    try {
      showLoader(context);

      String requestBody = jsonEncode({
        'amount': amount,
      });

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: appState.activeWallet!.isSharedWallet
            ? '/v1/shared-access/tokenization/subscriptions/${tokenizedAsset.id}'
            : '/v1/tokenization/subscriptions/${tokenizedAsset.id}',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.activeWallet!.publicKey!,
      );

      print('==============>response: $responseData');
      inspect(responseData);
      hideLoader(context);

      if (responseData['statusCode'] == 200 ||
          responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        await postProcessData(
            messageShown, messageLength, responseData['data']);
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'].toString().isEmpty
              ? responseData['data']['error']
              : responseData['data']['message'],
        );
      }
    } catch (e) {
      // print(e);
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  postProcessData(messageShown, messageLength, data) async {
    // we would like to display all messages returned from the initial
    // request to server using a popup. In order to achieve that we
    // employ the use of a little recursion here. Please recursive
    // functions can turn into a nightmare fast so be carefull here.
    if (messageShown <= messageLength - 1) {
      showResponseMessage(
          context,
          data['messages'][messageShown],
          () => {
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    sendFullDataToServer(transactionData: data);
  }

  Future<void> sendFullDataToServer({required Map transactionData}) async {
    try {
      showLoader(context);
      if (transactionData['signatureRequired'] == 0) {
        transactionData['commit'] = 1;
      } else {
        // sign transaction
        var signature = TrovoWalletSDK().signBase64Txn(
          appState.secretKeys[0], // the primary wallet secret key,
          transactionData['transaction'],
          transactionData['networkPassPhrase'],
        );
        transactionData['transactionSignature'] = signature;
      }
      var requestBody = jsonEncode(transactionData);
      print('requestBody  =======> $requestBody');

      Map responseData = await makePostRequest(
        uri: appState.activeWallet!.isSharedWallet
            ? '/v1/shared-access/tokenization/subscriptions/${tokenizedAsset.id}'
            : '/v1/tokenization/subscriptions/${tokenizedAsset.id}',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.activeWallet!.publicKey!,
      );

      hideLoader(context);

      print('responseData token information  ${responseData}');
      inspect(responseData);

      if (responseData['statusCode'] == 200 ||
          responseData['statusCode'] == 202) {
        if (transactionData['signatureRequired'] == 0) {
          appState.viewData = {
            SuccessViewPageConfig.key: {
              'title': 'Purchase request submitted',
              'message':
                  'You have successfully requested to purchase $amount ${tokenizedAsset.assetCode.toString()} with wallet ${appState.activeWallet!.alias}. This transaction will be completed when it gets the required number of approvals from the authorized approvers.',
              'useOnDone': true,
              'onDone': () {
                appState.currentAction = PageAction(
                  state: PageState.replaceAll,
                  page: BottomHomePageConfig,
                );
              },
            }
          };
          appState.currentAction =
              PageAction(state: PageState.replace, page: SuccessViewPageConfig);
        } else {
          appState.viewData = responseData['data'];
          appState.currentAction = PageAction(
              state: PageState.replace, page: BuyTokensSuccessViewPageConfig);
        }
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
