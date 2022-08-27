import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
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
  Wallet? activeWallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int touchedIndex = -1;
  String password = '';
  var viewData;
  var sourceAmount;
  var swappedEstimate;

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
    activeWallet = appState.activeWallet;
    viewData = appState.viewData![SwapSuccessViewPageConfig.key];
    sourceAmount = double.parse(viewData['sourceAmount']).toStringAsFixed(4);
    swappedEstimate =
        double.parse(viewData['swappedEstimate']).toStringAsFixed(4);
    print('viewData $viewData');

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
                LanguageEn.yourtransactionwassuccessful,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 22.sp),
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
                          LanguageEn.swapped,
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
                        padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                        child: showUserInfo(),
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
                          LanguageEn.blockchainproof,
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
                                onTap: () => appState.goToWebView(
                                    bantuBlockchainExplorerBaseUrl +
                                        viewData['transactionId']),
                                child: Text(
                                  viewData['transactionId'],
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
                                      text: viewData['transactionId'],
                                    ),
                                  ),
                                  showSnackBar('Transaction ID', context),
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
                LanguageEn.dashboard,
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
              '${sourceAmount} ${viewData['sourceAssetCode'].toString().isEmpty ? 'XBN' : viewData['sourceAssetCode']}',
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
              '- 34,000',
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontSize: 12.sp,
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
              LanguageEn.to,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 16.sp,
                fontFamily: fontsemibold,
              ),
            ),
            SizedBox(
              height: 20,
            ),
            Text(
              '${swappedEstimate} ${viewData['destinationAssetCode'].toString().isEmpty ? 'XBN' : viewData['destinationAssetCode']}',
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
              '- 34,000',
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontSize: 12.sp,
                fontWeight: FontWeight.w500,
                fontFamily: fontbody,
              ),
            ),
          ],
        )
      ],
    );
  }
}
