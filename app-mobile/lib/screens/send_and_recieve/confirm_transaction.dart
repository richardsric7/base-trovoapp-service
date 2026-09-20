import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/bottom_tab_page.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/cache.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';

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
  late Wallet wallet;
  late Asset? asset;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  var viewData;
  var transactionData;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet = appState.userInfo!.getWallet(appState.viewData!['walletAddress']);
    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );
    viewData = appState.viewData!;
    transactionData = appState.viewData!['transactionData'];
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    inspect(appState.viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "confirmyourtransaction".tr(),
                    style: TextStyle(
                      fontSize: 20.sp,
                      fontWeight: FontWeight.bold,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 20),
              Text(
                "youareabouttosend".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(height: height / 50),
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
                          SizedBox(height: height / 50),
                          Text(
                            '${formatNumber(double.parse(transactionData['amount']))} ${asset!.assetCode.toString().isEmpty ? 'ETH' : asset!.assetCode}',
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w700,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                            ),
                          ),
                          SizedBox(height: height / 50),
                          Container(
                            width: width / 1.3,
                            child: Text(
                              '- ${calculateFiatValue(transactionData['amount'], asset!.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody,
                              ),
                            ),
                          ),
                          SizedBox(height: height / 50.0),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Text(
                "to".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(height: height / 50),
              showAddressInfo(),
              SizedBox(height: height / 50),
              if (transactionData['memo'].toString().isNotEmpty) ...[
                showMemo(),
              ],
              if (wallet.isSharedWalletAndCanInitiate) ...[
                SizedBox(height: height / 50),
                Text(
                  "servicefee".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: const BorderRadius.all(
                        Radius.circular(15.0),
                      ),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Column(
                        children: [
                          SizedBox(height: height / 50),
                          myKeyValueRow(
                            "${"fee".tr()}: ",
                            transactionData['fee'] + '%',
                          ),
                          myKeyValueRow(
                            "${"amountcalculated".tr()}: ",
                            "${transactionData['feeAmount']} ${asset!.assetCode.toString().isEmpty ? 'ETH' : asset!.assetCode}",
                          ),
                          if (transactionData['vatAmount']
                              .toString()
                              .isNotEmpty) ...[
                            myKeyValueRow(
                              "${"vat".tr()}: ",
                              "${transactionData['vat']}%",
                            ),
                            myKeyValueRow(
                              "${"vatamount".tr()}: ",
                              "${transactionData['vatAmount']} ${asset!.assetCode.toString().isEmpty ? 'ETH' : asset!.assetCode}",
                            ),
                          ],
                          SizedBox(height: height / 50),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
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
              SizedBox(height: height / 20),
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

  Widget myKeyValueRow(String key, String value) {
    return Row(
      children: [
        Expanded(
          child: Text.rich(
            TextSpan(
              text: key,
              style: TextStyle(
                fontSize: 15,
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
              ),
              children: [
                TextSpan(
                  text: value,
                  style: TextStyle(
                    fontSize: 13,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget showAddressInfo() {
    if (transactionData['destination'].toString().length == 42) {
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
                    transactionData['destination'].toString(),
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
                  child:
                      transactionData['destinationThumbnail'].toString().isEmpty
                      ? CircleAvatar(
                          radius: 30,
                          backgroundColor: notifier.getaddsubwalletgrey,
                          foregroundImage: AssetImage(
                            "assets/images/default-user.png",
                          ),
                        )
                      : CircleAvatar(
                          radius: 30,
                          backgroundColor: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                          foregroundImage: NetworkImage(
                            transactionData['destinationThumbnail'].toString(),
                          ),
                        ),
                ),
                SizedBox(width: width / 70),
                Container(
                  width: width / 2,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        transactionData['destination'].toString(),
                        style: TextStyle(
                          fontWeight: FontWeight.w500,
                          color: notifier.getbluewhitecolor,
                          fontSize: 19.sp,
                          fontFamily: fontbody,
                        ),
                      ),
                      SizedBox(height: 5),
                      Text(
                        '${transactionData['destinationFirstName']} ${transactionData['destinationLastName']}',
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: 12.sp,
                          fontWeight: FontWeight.w500,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
                ),
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
          "descriptionmemo".tr(),
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: 15,
            fontWeight: FontWeight.w400,
            color: notifier.getbluewhitecolor,
            fontFamily: fontbody,
          ),
        ),
        SizedBox(height: height / 50),
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
                    vertical: 30.0,
                    horizontal: 15,
                  ),
                  child: Container(
                    width: width / 1.3,
                    child: Text(
                      transactionData['memo'],
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

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  sendDataToServer() async {
    try {
      showLoader(context);

      if (wallet.isSharedWalletAndCanInitiate) {
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

      String requestBody = jsonEncode(transactionData);

      Map responseData = await makePostRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/payment'
            : '/v1/users/payment',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        address: wallet.address!,
      );

      if (responseData['statusCode'] == 200) {
        updateUserInfo(
          appState.primaryWallet.signer!,
          appState.secretKeys[0],
          appState.primaryWallet.address,
          appState.userInfo!.username,
          appState,
          forceRefresh: true,
        );
        if (wallet.isSharedWalletAndCanInitiate) {
          appState.viewData = {
            SuccessViewPageConfig.key: {
              'title': 'Payment request submitted',
              'message':
                  'You have successfully requested payment of [${transactionData['amount']} ${asset!.assetCode.toString().isEmpty ? 'ETH' : asset!.assetCode}] from [${wallet.alias}] to [${transactionData['destination']}]. This transaction will be completed when it gets the required number of approvals by those who have approver access on this wallet.',
              'useOnDone': true,
              'onDone': () {
                appState.currentAction = PageAction(
                  state: PageState.replaceAll,
                  page: BottomHomePageConfig,
                );
                changeTabPage(appState, ButtomTabPage.Dashboard.index);
              },
            },
          };
          appState.currentAction = PageAction(
            state: PageState.replace,
            page: SuccessViewPageConfig,
          );
        } else {
          appState.viewData = {
            'transactionData': responseData['data'],
            'walletAddress': wallet.address,
          };

          appState.currentAction = PageAction(
            state: PageState.replaceAll,
            page: TransactionSuccessViewPageConfig,
          );
        }
        hideLoader(context);
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'].toString().isEmpty
              ? responseData['data']['error']
              : responseData['data']['message'],
        );
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
      hideLoader(context);
    }
  }
}
