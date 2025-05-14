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
  double trovUsdPrice = 0;
  String feeInfo = '';
  String fiatCurrency = '';
  double tokenizationApplicationFee = 0;
  String tokenizationApplicationFeeAsset = '';
  final Authenticator _authenticator = Authenticator();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = TokenizedAsset().deserializeJson(appState.viewData!);
    for (var asset in appState.primaryWallet.claimedAssets!) {
      if (asset.assetCode!.toUpperCase() == 'TROV') {
        trovUsdPrice = asset.usdPrice!;
      }
    }

    var quoteCurrencyCode = '';

    for (var i = 0;
        i < appState.tokenizationData['countryConfigs'].length;
        i++) {
      if (appState.tokenizationData['countryConfigs'][i]['countryCode']
              .toString()
              .toLowerCase() ==
          tokenizedAsset.assetCountryLocation.toString().toLowerCase()) {
        quoteCurrencyCode =
            appState.tokenizationData['countryConfigs'][i]['quoteCurrencyCode'];
        tokenizationApplicationFee = appState.tokenizationData['countryConfigs']
            [i]['tokenizationApplicationFee'];
        tokenizationApplicationFeeAsset = appState
            .tokenizationData['countryConfigs'][i]
                ['tokenizationApplicationFeeAsset']
            .toString()
            .split(':')[0];
      }
    }

    for (var i = 0;
        i < appState.tokenizationData['tokenizationCurrencies'].length;
        i++) {
      if (appState.tokenizationData['tokenizationCurrencies'][i]['assetCode']
              .toString()
              .toLowerCase() ==
          quoteCurrencyCode.toString().toLowerCase()) {
        fiatCurrency = appState.tokenizationData['tokenizationCurrencies'][i]
                ['label']
            .toString()
            .toUpperCase();
      }
    }

    feeInfo = formatNumber(getFeeInfo(tokenizedAsset.tokenizationFeeId!));
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    var isAlreadySubmitted = tokenizedAsset.tokenizationStatus! >= 1;
    var isVetted = tokenizedAsset.vettingStatus == 1;
    inspect(appState.tokenizationData);

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
                            '${truncateToDecimalPlaces(tokenizedAsset.assetCurrentValue!, decimalPlaces: 2)} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Percentage Retained",
                            '${formatNumber(((tokenizedAsset.assetOwnerRetainedOrContributedValue! / tokenizedAsset.assetCurrentValue!) * 100))}%'),
                        SizedBox(height: height / 90),
                        item("Value Retained",
                            '${(truncateToDecimalPlaces(tokenizedAsset.assetOwnerRetainedOrContributedValue!, decimalPlaces: 2))} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Additional Cost",
                            '${truncateToDecimalPlaces(tokenizedAsset.assetMscCostOutisdeOfValuation!, decimalPlaces: 2)} ${fiatCurrency}'),
                      ] else ...[
                        item("Project Budget",
                            '${truncateToDecimalPlaces(tokenizedAsset.assetCurrentValue!, decimalPlaces: 2)} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Equity Contribution",
                            '${formatNumber(((tokenizedAsset.assetOwnerRetainedOrContributedValue! / tokenizedAsset.assetCurrentValue!) * 100))}%'),
                        SizedBox(height: height / 90),
                        item("Value of Equity",
                            '${(truncateToDecimalPlaces(tokenizedAsset.assetOwnerRetainedOrContributedValue!, decimalPlaces: 2))} ${fiatCurrency}'),
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
                        item("Total Token",
                            '${(truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeIssued!))} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                        item("Final Value",
                            '${(truncateToDecimalPlaces(tokenizedAsset.valueOfTokenizedAsset!, decimalPlaces: 2))} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Price per Token",
                            '${(truncateToDecimalPlaces(tokenizedAsset.pricePerToken!, decimalPlaces: 2))} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Tokens not for Sale",
                            '${(truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeIssued! - tokenizedAsset.numberOfTokenToBeSold! - tokenizedAsset.feeInAsset!))} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                        item("Tokens for Sale",
                            '${(truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeSold!))} ${tokenizedAsset.assetCode}'),
                        SizedBox(height: height / 90),
                        item("Amount to be Raised",
                            '${truncateToDecimalPlaces(tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!, decimalPlaces: 2)} ${fiatCurrency}'),
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
                          isVetted ? 'Fees in Fiat' : 'Application Fee',
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
                        item("Tokenization Fee", '${feeInfo} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("SEC Fee",
                            '${formatNumberShort(tokenizedAsset.SECTokenizationFeeValue!)} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Custody Fee",
                            '${formatNumberShort(tokenizedAsset.custodianFeeValue!)} ${fiatCurrency}',
                            isNotUpfront: true),
                        SizedBox(height: height / 90),
                        item("Management Fee",
                            '${formatNumberShort(tokenizedAsset.assetManagerFeeValue!)} ${fiatCurrency}',
                            isNotUpfront: true),
                        SizedBox(height: height / 90),
                        if (tokenizedAsset.issuingHouseFeeValue! > 0) ...[
                          item("Issuing house Fee",
                              '${formatNumberShort(tokenizedAsset.issuingHouseFeeValue!)} ${fiatCurrency}'),
                          SizedBox(height: height / 90),
                        ],
                        if (tokenizedAsset.legalAndProfessionalFeeValue! >
                            0) ...[
                          item("Legal/Professional Fee",
                              '${formatNumberShort(tokenizedAsset.legalAndProfessionalFeeValue!)} ${fiatCurrency}'),
                          SizedBox(height: height / 90),
                        ],
                        if (tokenizedAsset.ratingAgencyFeeValue! > 0) ...[
                          item("Rating agency Fee",
                              '${formatNumberShort(tokenizedAsset.ratingAgencyFeeValue!)} ${fiatCurrency}'),
                          SizedBox(height: height / 90),
                        ],
                        item("VAT (Fiat)",
                            '${formatNumberShort(tokenizedAsset.vatValue!)} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                        item("Total", '${getTotalFee()} ${fiatCurrency}'),
                        SizedBox(height: height / 90),
                      ] else ...[
                        item("applicationfee".tr(),
                            '${formatNumber(tokenizationApplicationFee)} $tokenizationApplicationFeeAsset ${tokenizedAsset.tokenizationStatus == 1 ? '(Paid)' : ''}'),
                        SizedBox(height: height / 90),
                        // item(
                        //     "tokenizationfee".tr(), "$feeInfo ${fiatCurrency}"),
                        // SizedBox(height: height / 90),
                        // item("otherstatutoryfees".tr(), ''),
                        // SizedBox(height: height / 90),
                      ]
                    ],
                  ),
                ),
              ),
              if (isVetted) ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(10.0)),
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
                            "Fees in Asset",
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        item("Tokenization Fee",
                            '${formatNumberShort(tokenizedAsset.feeInAsset!)} ${tokenizedAsset.assetCode?.toUpperCase()}',
                            isNotUpfront: true),
                        SizedBox(height: height / 90),
                        item("VAT (Asset)",
                            '${formatNumberShort(tokenizedAsset.feeInAsset! * 0.075)} ${tokenizedAsset.assetCode?.toUpperCase()}',
                            isNotUpfront: true),
                        SizedBox(height: height / 90),
                        item("Total",
                            '${formatNumberShort((tokenizedAsset.feeInAsset! * 0.075) + tokenizedAsset.feeInAsset!)} ${tokenizedAsset.assetCode?.toUpperCase()}'),
                        SizedBox(height: height / 90),
                      ],
                    ),
                  ),
                ),
              ],
              if (isAlreadySubmitted) ...[
                if (tokenizedAsset.vettingStatus == 1 &&
                    tokenizedAsset.tokenizationStatus == 1) ...[
                  notifyAdditionalInfo(
                      "Items marked in asterisks (*) are not to be paid upfront."),
                  SizedBox(
                    height: height / 30,
                  ),
                  Button(
                    "Proceed to Pay",
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
                        "Please note that tokenization fee and other statutory fees will be displayed after vetting"),
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
                    "Please note that tokenization fee and other statutory fees will be displayed after vetting"),
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
                                "Please authorize the deduction of application fee to submit application.",
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
              ],
              SizedBox(height: 8),
              TextButton(
                child: Text(
                  'View Asset Details',
                  style: TextStyle(
                    fontSize: 14.0,
                    decoration: TextDecoration.underline,
                    fontFamily: fontbody,
                    fontWeight: FontWeight.bold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                onPressed: () {
                  appState.tokenizedAsset = tokenizedAsset;
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: AssetDashboardViewPageConfig,
                  );
                },
              ),
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

  String getTotalFee() {
    var total = tokenizedAsset.SECTokenizationFeeValue! +
        tokenizedAsset.custodianFeeValue! +
        tokenizedAsset.assetManagerFeeValue! +
        tokenizedAsset.issuingHouseFeeValue! +
        tokenizedAsset.legalAndProfessionalFeeValue! +
        tokenizedAsset.ratingAgencyFeeValue! +
        tokenizedAsset.vatValue! +
        getFeeInfo(tokenizedAsset.tokenizationFeeId!);

    return "${formatNumberShort(total)}";
  }

  double getFeeInfo(int index) {
    var fiatPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeFiatPercentage'];
    var fiatFeeCap = double.parse(appState.tokenizationData["tokenizationFees"]
            [index]['feeFiatCap']
        .toString());
    fiatFee = (tokenizedAsset.assetCurrentValue! * fiatPercentage) / 100;
    return fiatFeeCap > fiatFee ? fiatFeeCap : fiatFee;
  }

  double getFeeInAsset(int index) {
    var assetFeePercentage = appState.tokenizationData["tokenizationFees"]
        [index]['feeAssetPercentage'];

    return (tokenizedAsset.numberOfTokenToBeIssued! * assetFeePercentage) / 100;
  }

  Widget item(String key, String value, {isNotUpfront = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 0, horizontal: 15),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            constraints: BoxConstraints(
              maxWidth: width / 2.36,
            ),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
            child: Row(
              children: [
                Text(
                  key,
                  textAlign: TextAlign.start,
                  style: TextStyle(
                    fontWeight: FontWeight.w500,
                    color: notifier.getbluewhitecolor,
                    fontSize: 13.sp,
                    fontFamily: fontbody,
                  ),
                ),
                if (isNotUpfront) ...[
                  Text(
                    ' *',
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontWeight: FontWeight.w500,
                      color: Colors.red,
                      fontSize: 13.sp,
                      fontFamily: fontbody,
                    ),
                  ),
                ]
              ],
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
