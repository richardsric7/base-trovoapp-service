import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class ConfirmTokenizationDetails extends StatefulWidget {
  const ConfirmTokenizationDetails({Key? key}) : super(key: key);

  @override
  State<ConfirmTokenizationDetails> createState() =>
      _ConfirmTokenizationDetails();
}

class _ConfirmTokenizationDetails extends State<ConfirmTokenizationDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  late TokenizedAsset tokenizedAsset;
  String password = '';
  final Authenticator _authenticator = Authenticator();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = TokenizedAsset().deserializeJson(appState.viewData!);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    var isAlreadySubmitted = tokenizedAsset.tokenizationStatus == 1;
    inspect(appState.viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
                context,
                notifier.getwihitecolor,
                isAlreadySubmitted
                    ? tokenizedAsset.assetName!
                    : "confirmyourinformation".tr(),
                notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 40,
              ),
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
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            vertical: 10.0, horizontal: 15),
                        child: Text(
                          'assetinformation'.tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 15.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      item("assetname".tr(), '${tokenizedAsset.assetName}'),
                      SizedBox(height: height / 90),
                      item("assetcode".tr(), '${tokenizedAsset.assetCode}'),
                      SizedBox(height: height / 90),
                      if (tokenizedAsset.assetAlreadyExists == 1) ...[
                        item("originalassetvalue".tr(),
                            '${formatNumber(tokenizedAsset.assetCurrentValue!)} ${tokenizedAsset.assetQuoteCurrency}'),
                      ] else ...[
                        item("originaltotalprojectcost".tr(),
                            '${formatNumber(tokenizedAsset.assetCurrentValue!)} ${tokenizedAsset.assetQuoteCurrency}'),
                      ],
                      SizedBox(height: height / 90),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 40,
              ),
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
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            vertical: 10.0, horizontal: 15),
                        child: Text(
                          'assettokeninfo'.tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 15.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      item("proposedtotaltokenstobeissued".tr(),
                          '${formatNumber(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode}'),
                      SizedBox(height: height / 90),
                      item("totalamounttoberaised".tr(),
                          '${formatNumber(tokenizedAsset.numberOfTokenToBeSold!)} ${tokenizedAsset.assetCode}'),
                      SizedBox(height: height / 90),
                      // item("pricepertoken".tr(),
                      //     '${(formatNumber(tokenizedAsset.pricePerToken!))} ${tokenizedAsset.assetQuoteCurrency}'),
                      // SizedBox(height: height / 90),
                      // item("totalamounttoberaised".tr(),
                      //     '${formatNumber(tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!)} ${tokenizedAsset.assetQuoteCurrency}'),
                      // SizedBox(height: height / 90),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 40),
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
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            vertical: 10.0, horizontal: 15),
                        child: Text(
                          'primaryoffering'.tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 15.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      item("proposedstartdate".tr(),
                          '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesStart!)}'),
                      SizedBox(
                        height: height / 90,
                      ),
                      item("proposedenddate".tr(),
                          '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesEnd!)}'),
                      SizedBox(height: height / 90),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 40),
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
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            vertical: 10.0, horizontal: 15),
                        child: Text(
                          'fees'.tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 15.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      item("applicationfee".tr(), '500 TROV'),
                      SizedBox(height: height / 90),
                      item("tokenizationfee".tr(),
                          getFeeInfo(tokenizedAsset.tokenizationFeeId!)),
                      SizedBox(height: height / 90),
                      item("otherstatutoryfees".tr(),
                          getFeeInfo(tokenizedAsset.tokenizationFeeId!)),
                      SizedBox(height: height / 90),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              if (isAlreadySubmitted) ...[
                Button(
                  "viewpaymentdetails".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: TokenizationFeePaymentViewPageConfig,
                    );
                  },
                ),
                SizedBox(height: height / 70),
                ButtonOutlined(
                  'back'.tr(),
                  notifier.getwihitecolor,
                  notifier.getbluewhitecolor,
                  borderColor: notifier.getbluewhitecolor,
                  onTap: () {
                    Navigator.of(context).pop();
                  },
                ),
              ] else ...[
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
                        submitForm();
                      } else {
                        popup(context,
                            title: "oops".tr(),
                            message: "invalidpassword".tr());
                      }
                    },
                  ),
                ],
                // Button(
                //   "submitapplication".tr(),
                //   notifier.getbluecolor,
                //   wihitecolor,
                //   onTap: () => submitForm(),
                // ),
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

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        submitForm();
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

  Future<void> submitForm({Map? transactionData}) async {
    try {
      showLoader(context);
      String requestBody = "{}";

      if (transactionData != null && transactionData['transaction'] != null) {
        var signature = TrovoWalletSDK().signBase64Txn(
          appState.secretKeys[0], // the primary wallet secret key,
          transactionData['transaction'],
          transactionData['networkPassPhrase'],
        );
        transactionData['transactionSignature'] = signature;
        requestBody = jsonEncode(transactionData);
        print('requestBody  =======> $requestBody');
      }

      Map responseData = await makePutRequest(
        uri: '/v1/tokenization/confirm/${tokenizedAsset.id}',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      hideLoader(context);

      print('responseData token information  ${responseData}');
      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        await postProcessData(
            messageShown, messageLength, responseData['data']);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
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
      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0], // the primary wallet secret key,
        transactionData['transaction'],
        transactionData['networkPassPhrase'],
      );
      transactionData['transactionSignature'] = signature;
      var requestBody = jsonEncode(transactionData);
      print('requestBody  =======> $requestBody');

      Map responseData = await makePutRequest(
        uri: '/v1/tokenization/confirm/${tokenizedAsset.id}',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      hideLoader(context);

      print('responseData token information  ${responseData}');
      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': '',
          'message':
              'Your Asset Tokenization application has been submitted successfully, please wait for vetting to be done, you will be notified once it is vetted so that you can proceed to pay the asset tokenization fee and other statutory fees',
        };
        appState.currentAction =
            PageAction(state: PageState.addPage, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  String getFeeInfo(int index) {
    var fiatPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeFiatPercentage'];
    var assetPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeAssetPercentage'];
    var fiatFeeCap = double.parse(appState.tokenizationData["tokenizationFees"]
            [index]['feeFiatCap']
        .toString());
    var tokenFee =
        (tokenizedAsset.numberOfTokenToBeIssued! * assetPercentage) / 100;
    var fiatFee = (tokenizedAsset.assetCurrentValue! * fiatPercentage) / 100;
    return "\$${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${tokenizedAsset.assetCode}";
  }

  Widget item(String key, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 0, horizontal: 15),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            constraints: BoxConstraints(
              maxWidth: width / 2.1,
            ),
            child: Text(
              key,
              textAlign: TextAlign.start,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 13.sp,
                fontFamily: fontbody,
              ),
            ),
          ),
          Container(
            constraints: BoxConstraints(
              maxWidth: width / 2.4,
            ),
            child: Text(
              value,
              textAlign: TextAlign.end,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 13.sp,
                fontFamily: fontsemibold,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
