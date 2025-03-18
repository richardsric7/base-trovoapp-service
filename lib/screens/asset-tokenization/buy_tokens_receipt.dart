import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:share_plus/share_plus.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class BuyTokensReceipt extends StatefulWidget {
  const BuyTokensReceipt({Key? key}) : super(key: key);

  @override
  State<BuyTokensReceipt> createState() => _BuyTokensReceipt();
}

class _BuyTokensReceipt extends State<BuyTokensReceipt>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  double amount = 0;
  double quantity = 0;
  late DateTime date;
  late TokenizedAsset tokenizedAsset;
  GlobalKey shareArea = GlobalKey();

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
    tokenizedAsset = appState.tokenizedAsset!;
    amount = appState.viewData!['amount'];
    quantity = appState.viewData!['quantity'];
    date = appState.viewData!['date'];

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
              RepaintBoundary(
                key: shareArea,
                child: Container(
                  color: notifier.isDark ? null : wihitecolor,
                  child: Column(
                    children: [
                      SizedBox(height: height / 50),
                      Image.asset(
                        'assets/images/trovo-horizontal-logo.png',
                        height: height / 16.5,
                        color: notifier.isDark ? Colors.white : null,
                      ),
                      SizedBox(height: height / 30),
                      Text(
                        'Purchase Receipt',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 22),
                      ),
                      SizedBox(height: 3),
                      Container(
                        width: width / 1.5,
                        child: Text(
                          "generatedon".tr(args: [
                            DateFormat('MMMM dd, yyyy hh:mm a')
                                .format(DateTime.now())
                          ]),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                              fontSize: 12),
                        ),
                      ),
                      SizedBox(height: height / 50),
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
                          child: Stack(
                            alignment: AlignmentDirectional.center,
                            children: [
                              Image.asset(
                                'assets/images/trovo_white.png',
                                height: height / 4.5,
                                color: notifier.isDark
                                    ? notifier.getdarkgrey
                                    : notifier.getsplashgrey,
                              ),
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      'Receiving Wallet',
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    width: width / 1.2,
                                    child: Column(
                                      children: [
                                        Row(
                                          children: [
                                            Padding(
                                              padding:
                                                  const EdgeInsets.symmetric(
                                                      horizontal: 20.0),
                                              child: Text(
                                                appState.activeWallet!.alias!,
                                                style: TextStyle(
                                                    fontWeight: FontWeight.w500,
                                                    color: notifier
                                                        .getbluewhitecolor,
                                                    fontSize: 16.sp,
                                                    fontFamily: fontbody,
                                                    overflow:
                                                        TextOverflow.visible),
                                              ),
                                            ),
                                          ],
                                        ),
                                        Row(
                                          children: [
                                            Expanded(
                                              flex: 3,
                                              child: Padding(
                                                padding:
                                                    const EdgeInsets.symmetric(
                                                        horizontal: 20.0,
                                                        vertical: 5),
                                                child: Text(
                                                  truncatePublicKey(appState
                                                      .activeWallet!
                                                      .publicKey!),
                                                  style: TextStyle(
                                                    fontWeight: FontWeight.w500,
                                                    color: notifier
                                                        .getbluewhitecolor,
                                                    fontSize: 13,
                                                    fontFamily: fontbody,
                                                  ),
                                                ),
                                              ),
                                            ),
                                          ],
                                        ),
                                      ],
                                    ),
                                  ),
                                  Divider(
                                    height: 5,
                                  ),
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      "amount".tr(),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    height: 5,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      formatAmount(
                                          TransactionDirection.Receive,
                                          amount,
                                          tokenizedAsset.assetQuoteCurrency),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 15,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ),
                                  Divider(
                                    height: 5,
                                  ),
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      "quantity".tr(),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    height: 5,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      formatAmount(TransactionDirection.Receive,
                                          quantity, tokenizedAsset.assetCode),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 15,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ),
                                  Divider(
                                    height: 5,
                                  ),
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      "blockchainproof".tr(),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    height: 5,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Row(
                                      children: [
                                        Expanded(
                                          flex: 5,
                                          child: GestureDetector(
                                            // onTap: () => appState.goToWebView(
                                            //     getExplorerBaseUrl(
                                            //             appState.walletMode) +
                                            //         viewData.transactionId!),
                                            child: Text(
                                              'viewData.transactionId!',
                                              style: TextStyle(
                                                decoration:
                                                    TextDecoration.underline,
                                                color:
                                                    notifier.getbluewhitecolor,
                                                fontSize: 12,
                                                fontWeight: FontWeight.w500,
                                                fontFamily: fontbody,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                  Divider(
                                    height: 5,
                                  ),
                                  SizedBox(
                                    height: height / 50,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Text(
                                      "date".tr(),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0, vertical: 5),
                                    child: Text(
                                      DateFormat('MMMM dd, yyyy hh:mm a')
                                          .format(date),
                                      style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 13,
                                        fontWeight: FontWeight.w500,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    height: height / 50,
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
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Button(
                "shareimage".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  share(
                      'Blockchain proof\n${getExplorerBaseUrl(appState.walletMode)}${'viewData.transactionId!'}',
                      shareArea);
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              Button(
                "sharepdf".tr(),
                notifier.getbluecolor70,
                wihitecolor,
                onTap: () {
                  sharePDF(
                      '${"blockchainproof".tr()}\n${getExplorerBaseUrl(appState.walletMode)}${'viewData.transactionId!'}',
                      shareArea);
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              ButtonOutlined(
                "sharetext".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  shareText();
                },
              ),
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }

  String formatAmount(TransactionDirection transactionType, amount, assetCode) {
    var am = formatNumber(double.parse(amount.toString()));
    return '$am $assetCode';
  }

  void shareText() {
    String? shareString = "sharestringpurchase".tr(args: [
      '${amount.toString()} ${tokenizedAsset.assetCode!}',
      appState.activeWallet!.alias!,
      'viewData.transactionId!.toLowerCase()',
      date.toString(),
      getExplorerBaseUrl(appState.walletMode) + 'viewData.transactionId!'
    ]);

    Share.share(shareString);
  }
}
