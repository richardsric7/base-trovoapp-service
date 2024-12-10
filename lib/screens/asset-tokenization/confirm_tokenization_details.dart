import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

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
    var isAlreadySubmitted = tokenizedAsset.tokenizationStatus != null;
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
                      item("currentvalueofasset".tr(),
                          '${formatNumber(tokenizedAsset.assetCurrentValue!)} ${tokenizedAsset.assetQuoteCurrency}'),
                      SizedBox(height: height / 90),
                      item("assetmanager".tr(),
                          '${tokenizedAsset.assetManagerInfo?.assetManagerName}'),
                      SizedBox(height: height / 90),
                      item("assetcustodian".tr(),
                          '${tokenizedAsset.approvedAssetCustodianInfo?.assetCustodianName}'),
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
                      item("totaltokenstobeissued".tr(),
                          '${formatNumber(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode}'),
                      SizedBox(height: height / 90),
                      item("totaltokenstobesold".tr(),
                          '${formatNumber(tokenizedAsset.numberOfTokenToBeSold!)} ${tokenizedAsset.assetCode}'),
                      SizedBox(height: height / 90),
                      item("pricepertoken".tr(),
                          '${(formatNumber(tokenizedAsset.pricePerToken!))} ${tokenizedAsset.assetQuoteCurrency}'),
                      SizedBox(height: height / 90),
                      item("totalamounttoberaised".tr(),
                          '${formatNumber(tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!)} ${tokenizedAsset.assetQuoteCurrency}'),
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
                      item("startdate".tr(),
                          '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesStart!)}'),
                      SizedBox(
                        height: height / 90,
                      ),
                      item("enddate".tr(),
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
                      item("tokenizationfee".tr(),
                          getFeeInfo(tokenizedAsset.tokenizationFeeId!)),
                      SizedBox(height: height / 90),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 20,
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
                Button(
                  "submitapplication".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () async {
                    appState.viewData!['tokenizationStatus'] = 0;
                    var savedAssets =
                        await StoreData().storeGetData('tokenizedAsset');
                    print(savedAssets);
                    if (savedAssets != null) {
                      for (int i = 0; i < savedAssets.length; i++) {
                        print(savedAssets);
                      }
                      savedAssets = [...savedAssets, appState.viewData];
                      await StoreData()
                          .storeInsertData('tokenizedAsset', savedAssets);
                    } else {
                      savedAssets = [appState.viewData];
                      await StoreData()
                          .storeInsertData('tokenizedAsset', savedAssets);
                    }
                    showLoader(context);
                    await Future.delayed(Duration(seconds: 1));
                    hideLoader(context);
                    appState.viewData![SuccessViewPageConfig.key] = {
                      'title': '',
                      'buttonText': 'Proceed to Pay',
                      'useOnDone': true,
                      'onDone': () {
                        appState.currentAction = PageAction(
                          state: PageState.replace,
                          page: TokenizationFeePaymentViewPageConfig,
                        );
                      },
                      'message':
                          'Your asset tokenization request has been submitted successfully. Please complete the payment to proceed. We will begin processing your application once the payment is received.',
                    };
                    appState.currentAction = PageAction(
                        state: PageState.addPage, page: SuccessViewPageConfig);
                  },
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
              maxWidth: width / 2.8,
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
