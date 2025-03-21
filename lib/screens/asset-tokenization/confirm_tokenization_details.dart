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
import 'package:provider/provider.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  double fiatFee = 0;
  double tokenFee = 0;
  String feeInfo = '';
  final Authenticator _authenticator = Authenticator();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = TokenizedAsset().deserializeJson(appState.viewData!);
    getFeeInfo(tokenizedAsset.tokenizationFeeId!);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    var isAlreadySubmitted = tokenizedAsset.tokenizationStatus! >= 1;
    var isVetted = tokenizedAsset.vettingStatus == 1;
    inspect(appState.viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
                context,
                notifier.getwihitecolor,
                isAlreadySubmitted && isVetted
                    ? 'Vetted Summary'
                    : "confirmyourinformation".tr(),
                notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              if (isAlreadySubmitted && isVetted) ...[
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Row(
                    children: [
                      Container(
                        width: width / 1.2,
                        child: Text(
                          'Your application has been vetted, please confirm the information below and proceed to pay for tokenization'
                              .tr(),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 13,
                            fontFamily: fontbody,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              SizedBox(height: height / 50),
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
                            '${truncateToDecimalPlaces(tokenizedAsset.assetCurrentValue!, decimalPlaces: 2)} ${appState.defaultCurrency}'),
                      ] else ...[
                        item("originaltotalprojectcost".tr(),
                            '${truncateToDecimalPlaces(tokenizedAsset.assetCurrentValue!, decimalPlaces: 2)} ${appState.defaultCurrency}'),
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
                      if (isVetted) ...[
                        item("Total tokens".tr(),
                            '${(truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeIssued!))} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                        item("Value of total tokens".tr(),
                            '${(truncateToDecimalPlaces(tokenizedAsset.valueOfTokenizedAsset!, decimalPlaces: 2))} ${appState.defaultCurrency}'),
                        SizedBox(height: height / 90),
                        item("pricepertoken".tr(),
                            '${(truncateToDecimalPlaces(tokenizedAsset.pricePerToken!, decimalPlaces: 2))} ${appState.defaultCurrency}'),
                        SizedBox(height: height / 90),
                        item("Tokens not for sale".tr(),
                            '${(truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeIssued! - tokenizedAsset.numberOfTokenToBeSold! - tokenFee))} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                        item(
                            tokenizedAsset.assetAlreadyExists == 1
                                ? "Amount retained".tr()
                                : "Amount contributed".tr(),
                            '${(truncateToDecimalPlaces(tokenizedAsset.assetOwnerRetainedOrContributedValue!, decimalPlaces: 2))} ${appState.defaultCurrency}'),
                        SizedBox(height: height / 90),
                        item("Tokens for sale".tr(),
                            '${(truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeSold!))} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                        item("totalamounttoberaised".tr(),
                            '${truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!, decimalPlaces: 2)} ${appState.defaultCurrency}'),
                        SizedBox(height: height / 90),
                      ] else ...[
                        item("proposedtotaltokenstobeissued".tr(),
                            '${formatNumber(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                      ],
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
                      if (isVetted) ...[
                        item("startdate".tr(),
                            '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesStart!)}'),
                        SizedBox(
                          height: height / 90,
                        ),
                        item("enddate".tr(),
                            '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesEnd!)}'),
                        SizedBox(height: height / 90),
                      ] else ...[
                        item("proposedstartdate".tr(),
                            '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesStart!)}'),
                        SizedBox(
                          height: height / 90,
                        ),
                        item("proposedenddate".tr(),
                            '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesEnd!)}'),
                        SizedBox(height: height / 90),
                      ],
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
                      if (isVetted) ...[
                        item("Asset tokenization fee".tr(), feeInfo),
                        SizedBox(height: height / 90),
                        item(
                            "SEC Regulatory Fee".tr(),
                            formatNumberShort(
                                tokenizedAsset.SECTokenizationFeeValue!)),
                        SizedBox(height: height / 90),
                        item(
                            "Asset Custody Fee".tr(),
                            formatNumberShort(
                                tokenizedAsset.custodianFeeValue!)),
                        SizedBox(height: height / 90),
                        item(
                            "Asset Management Fee".tr(),
                            formatNumberShort(
                                tokenizedAsset.assetManagerFeeValue!)),
                        SizedBox(height: height / 90),
                      ] else ...[
                        item("applicationfee".tr(),
                            '500 TROV ${tokenizedAsset.tokenizationStatus == 1 ? '(Paid)' : ''}'),
                        SizedBox(height: height / 90),
                        item("tokenizationfee".tr(), feeInfo),
                        SizedBox(height: height / 90),
                        // item("otherstatutoryfees".tr(), ''),
                        // SizedBox(height: height / 90),
                      ]
                    ],
                  ),
                ),
              ),
              if (isAlreadySubmitted) ...[
                if (tokenizedAsset.vettingStatus == 1 &&
                    tokenizedAsset.tokenizationStatus == 1) ...[
                  SizedBox(
                    height: height / 30,
                  ),
                  Button(
                    "Proceed to Pay".tr(),
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
                ] else ...[
                  if (!isVetted) ...[
                    notifyAdditionalInfo(
                        "Please note that other statutory fees will be added after vetting"
                            .tr()),
                  ],
                  SizedBox(height: 20),
                  ButtonOutlined(
                    'back'.tr(),
                    notifier.getwihitecolor,
                    notifier.getbluewhitecolor,
                    borderColor: notifier.getbluewhitecolor,
                    onTap: () {
                      Navigator.of(context).pop();
                    },
                  ),
                ],
              ] else ...[
                SizedBox(height: 10),
                notifyAdditionalInfo(
                    "Please note that other statutory fees will be added after vetting"
                        .tr()),
                SizedBox(height: 5),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    child: Card(
                      shadowColor: Colors.black,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(5.0),
                      ),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 5.0, vertical: 10),
                        child: Row(
                          children: [
                            SizedBox(
                              width: 300,
                              child: Text(
                                "Please authorize the deduction of application fee to submit application."
                                    .tr(),
                                textAlign: TextAlign.start,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                  overflow: TextOverflow.visible,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
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

  Widget notifyAdditionalInfo(String info) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        child: Card(
          shadowColor: Colors.black,
          shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(5.0),
              side: BorderSide(
                color: notifier.getbluewhitecolor,
                width: 1,
              )),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 5.0, vertical: 10),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(20.0),
                        side: BorderSide(
                          color: notifier.getbluewhitecolor,
                          width: 1,
                        )),
                    color: notifier.getaddsubwalletgrey,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 5.0),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "i".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                              overflow: TextOverflow.visible,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                SizedBox(
                  width: 300,
                  child: Text(
                    info,
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                      overflow: TextOverflow.visible,
                    ),
                  ),
                ),
              ],
            ),
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

  getFeeInfo(int index) {
    var fiatPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeFiatPercentage'];
    var assetPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeAssetPercentage'];
    var fiatFeeCap = double.parse(appState.tokenizationData["tokenizationFees"]
            [index]['feeFiatCap']
        .toString());
    tokenFee =
        (tokenizedAsset.numberOfTokenToBeIssued! * assetPercentage) / 100;
    fiatFee = (tokenizedAsset.assetCurrentValue! * fiatPercentage) / 100;
    feeInfo =
        "\$${formatNumber(fiatFeeCap > fiatFee ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${tokenizedAsset.assetCode}";
  }

  Widget item(String key, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 0, horizontal: 15),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            constraints: BoxConstraints(
              maxWidth: width / 2.36,
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
            width: width / 2.56,
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
