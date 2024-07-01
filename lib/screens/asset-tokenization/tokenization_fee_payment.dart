import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
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

class TokenizationFeePayment extends StatefulWidget {
  const TokenizationFeePayment({Key? key}) : super(key: key);

  @override
  State<TokenizationFeePayment> createState() => _TokenizationFeePayment();
}

class _TokenizationFeePayment extends State<TokenizationFeePayment>
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

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(context, notifier.getwihitecolor, "payment".tr(),
                notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "makepaymentto".tr(),
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "payto".tr(args: ["N1,000,000.00"]),
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "walletaddress".tr(),
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
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
                    children: [
                      Padding(
                        padding: const EdgeInsets.all(10.0),
                        child: Container(
                          width: width / 1.4,
                          child: Text(
                            tokenizedAsset.walletToHoldAssetsNotForSale!
                                .toLowerCase(),
                            style: TextStyle(
                                fontSize: 15,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold),
                          ),
                        ),
                      ),
                      IconButton(
                        onPressed: () {
                          Clipboard.setData(
                            ClipboardData(
                              text:
                                  tokenizedAsset.walletToHoldAssetsNotForSale!,
                            ),
                          );
                          showSnackBar("walletalias".tr(), context);
                        },
                        icon: Icon(Icons.copy,
                            size: 20, color: notifier.getbluewhitecolor),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "or".tr().toUpperCase(),
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "bankaccountdetails".tr(),
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
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
                    padding: const EdgeInsets.all(15.0),
                    child: Column(
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              "accountname".tr(),
                              style: TextStyle(
                                  fontSize: 15,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            Text(
                              "Trovo Tokenizer",
                              style: TextStyle(
                                  fontSize: 15,
                                  fontWeight: FontWeight.w700,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold),
                            ),
                          ],
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              "accountnumber".tr(),
                              style: TextStyle(
                                  fontSize: 15,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  "0123908438",
                                  style: TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.w700,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontsemibold),
                                ),
                                IconButton(
                                  onPressed: () {
                                    Clipboard.setData(
                                      ClipboardData(
                                        text: tokenizedAsset
                                            .walletToHoldAssetsNotForSale!,
                                      ),
                                    );
                                    showSnackBar("walletalias".tr(), context);
                                  },
                                  icon: Icon(Icons.copy,
                                      size: 20,
                                      color: notifier.getbluewhitecolor),
                                ),
                              ],
                            ),
                          ],
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              "bank".tr(),
                              style: TextStyle(
                                  fontSize: 15,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            Text(
                              "First Bank",
                              style: TextStyle(
                                  fontSize: 15,
                                  fontWeight: FontWeight.w700,
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
              SizedBox(
                height: height / 30,
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
                    padding: const EdgeInsets.all(15.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          "${"important".tr()}!",
                          style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w700,
                              color: Colors.red,
                              fontFamily: fontbody),
                        ),
                        SizedBox(
                          height: height / 70,
                        ),
                        Text(
                          "yourtokenizationfeeis".tr(),
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody),
                        ),
                        SizedBox(
                          height: height / 90,
                        ),
                        Text(
                          "N1,000,000.00 + 1,500,000.00 ${tokenizedAsset.assetCode}",
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(
                          height: height / 70,
                        ),
                        Text(
                          "pleasepayasap".tr(),
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "submitapplication".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () async {
                  appState.viewData!['tokenizationStatus'] = 0;
                  await inspect(appState.viewData);
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
                  await Future.delayed(Duration(seconds: 12));
                  hideLoader(context);
                  appState.viewData![SuccessViewPageConfig.key] = {
                    'title': '',
                    'message':
                        'Your Asset Tokenization Request has been submitted and is awaiting approval. You’ll be notified when  it has been approved.',
                  };
                  appState.currentAction = PageAction(
                      state: PageState.replace, page: SuccessViewPageConfig);
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
                      'transactionData',
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

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }
}
