import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class BuyTokensSuccess extends StatefulWidget {
  const BuyTokensSuccess({Key? key}) : super(key: key);

  @override
  State<BuyTokensSuccess> createState() => _BuyTokensSuccess();
}

class _BuyTokensSuccess extends State<BuyTokensSuccess>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  double amount = 0;
  double quantity = 0;
  DateTime date = DateTime.now();
  late TokenizedAsset tokenizedAsset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
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

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Center(
                child: Image.asset("assets/images/success.gif",
                    height: height / 10),
              ),
              SizedBox(height: height / 50),
              Text(
                "yourtransactionwassuccessful".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 22.sp),
              ),
              SizedBox(height: height / 30),
              Text(
                '- $quantity ${tokenizedAsset.assetCode}',
                style: TextStyle(
                    color: Colors.red,
                    fontFamily: fontsemibold,
                    fontSize: 20.sp),
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
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                        child: Text(
                          "Receiving Wallet".tr(),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Row(
                          children: [
                            Container(
                              width: width / 1.8,
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    appState.activeWallet!.alias!,
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
                                    truncatePublicKey(
                                        appState.activeWallet!.publicKey),
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
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      Divider(
                        height: 5,
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 10),
                        child: Text(
                          "Amount".tr(),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Text(
                          formatNumber(amount),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 15.sp,
                            fontFamily: fontbody,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      Divider(
                        height: 5,
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 10),
                        child: Text(
                          "blockchainproof".tr(),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Row(
                          children: [
                            Expanded(
                              flex: 5,
                              child: GestureDetector(
                                // onTap: () => appState.goToWebView(
                                //     getExplorerBaseUrl(appState.walletMode) +
                                //         viewData['transactionId']),

                                child: Text(
                                  "viewData['transactionId']",
                                  style: TextStyle(
                                    decoration: TextDecoration.underline,
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 12.sp,
                                    fontWeight: FontWeight.w500,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              ),
                            ),
                            Expanded(
                              flex: 1,
                              child: IconButton(
                                onPressed: () => {
                                  Clipboard.setData(
                                    ClipboardData(
                                      text: "viewData['transactionId']",
                                    ),
                                  ),
                                  showSnackBar("transactionid".tr(), context),
                                },
                                icon: Icon(Icons.copy),
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ],
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      Divider(
                        height: 5,
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 10),
                        child: Text(
                          "Date".tr(),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: GestureDetector(
                          // onTap: () => appState.goToWebView(
                          //     getExplorerBaseUrl(appState.walletMode) +
                          //         viewData['transactionId']),

                          child: Text(
                            DateFormat('MMMM dd, yyyy hh:mm a').format(date),
                            style: TextStyle(
                              decoration: TextDecoration.underline,
                              color: notifier.getbluewhitecolor,
                              fontSize: 12.sp,
                              fontWeight: FontWeight.w500,
                              fontFamily: fontbody,
                            ),
                          ),
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "generatereceipt".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.viewData = {
                    'amount': amount,
                    'quantity': quantity,
                    'date': date,
                  };
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: BuyTokensReceiptViewPageConfig,
                  );
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              ButtonOutlined(
                "dashboard".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.replaceAll,
                    page: BottomHomePageConfig,
                  );
                },
              ),
              SizedBox(
                height: height / 20,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
