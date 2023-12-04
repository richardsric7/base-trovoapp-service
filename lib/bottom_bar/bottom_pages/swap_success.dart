import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SwapSuccess extends StatefulWidget {
  const SwapSuccess({Key? key}) : super(key: key);

  @override
  State<SwapSuccess> createState() => _SwapSuccess();
}

class _SwapSuccess extends State<SwapSuccess> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  late Wallet wallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int touchedIndex = -1;
  String password = '';
  var sourceAmount;
  var swappedEstimate;
  late Map transactionData = {};
  late Map viewData = {};

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );
    viewData = appState.viewData!;
    transactionData = viewData['transactionData'];
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    sourceAmount =
        double.parse(transactionData['sourceAmount']).toStringAsFixed(4);
    swappedEstimate =
        double.parse(transactionData['swappedEstimate']).toStringAsFixed(4);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 18),
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
                    fontSize: 22),
              ),
              SizedBox(height: height / 30),
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
                          "swapped".tr(),
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
                        padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                        child: showUserInfo(),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      Divider(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              "servicefee".tr(),
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 16,
                                fontFamily: fontsemibold,
                              ),
                            ),
                            SizedBox(
                              height: 20,
                            ),
                            Text(
                              '${transactionData['feeAmount']} ${transactionData['sourceAssetCode'].toString().isEmpty ? 'XBN' : transactionData['sourceAssetCode']} (${transactionData['fee']}%)',
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 19,
                                fontFamily: fontbody,
                              ),
                            ),
                          ],
                        ),
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                      Divider(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 10),
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
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Row(
                          children: [
                            Expanded(
                              flex: 5,
                              child: GestureDetector(
                                onTap: () => appState.goToWebView(
                                    getExplorerBaseUrl(appState.walletMode) +
                                        transactionData['transactionId']),
                                child: Text(
                                  transactionData['transactionId'],
                                  style: TextStyle(
                                    decoration: TextDecoration.underline,
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 12,
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
                                      text: transactionData['transactionId'],
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
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "dashboard".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.replaceAll,
                    page: BottomHomePageConfig,
                  );
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

  Widget showUserInfo() {
    return Row(
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '${sourceAmount} ${transactionData['sourceAssetCode'].toString().isEmpty ? 'XBN' : transactionData['sourceAssetCode']}',
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 19,
                fontFamily: fontbody,
              ),
            ),
            SizedBox(
              height: 5,
            ),
            Text(
              '- ${calculateFiatValue(sourceAmount, viewData["sourceUsdPrice"].toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontSize: 12,
                fontWeight: FontWeight.w500,
                fontFamily: fontbody,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Divider(
              height: 5,
            ),
            Text(
              "to".tr(),
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 16,
                fontFamily: fontsemibold,
              ),
            ),
            SizedBox(
              height: 20,
            ),
            Text(
              '${swappedEstimate} ${transactionData['destinationAssetCode'].toString().isEmpty ? 'XBN' : transactionData['destinationAssetCode']}',
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 19,
                fontFamily: fontbody,
              ),
            ),
            SizedBox(
              height: 5,
            ),
            if (viewData["destinationUsdPrice"] != null) ...[
              Text(
                '+ ${calculateFiatValue(swappedEstimate, viewData["destinationUsdPrice"].toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  fontFamily: fontbody,
                ),
              ),
            ]
          ],
        )
      ],
    );
  }
}
