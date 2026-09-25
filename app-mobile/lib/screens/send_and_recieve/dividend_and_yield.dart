import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class DividendAndYieldView extends StatefulWidget {
  const DividendAndYieldView({Key? key}) : super(key: key);

  @override
  State<DividendAndYieldView> createState() => _DividendAndYieldView();
}

class _DividendAndYieldView extends State<DividendAndYieldView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late Asset? asset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (appState.viewData!['walletAddress'] != null) {
      wallet = appState.userInfo!.getWallet(
        appState.viewData!['walletAddress'],
      );
      asset = wallet.claimedAssets!.firstWhere(
        (asset) =>
            asset.assetCode == appState.viewData!['assetCode'] &&
            asset.contractAddress == appState.viewData!['contractAddress'],
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          'Dividend and Interest',
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                spacing: 10,
                children: [
                  if (asset?.imageUrl != null) ...[
                    ClipRRect(
                      borderRadius: BorderRadius.circular(100.0),
                      child: Image.network(
                        asset!.imageUrl!,
                        width: 40,
                        height: 40,
                        fit: BoxFit.fill,
                      ),
                    ),
                  ] else ...[
                    Image.asset(
                      'assets/images/trovo.png',
                      height: 35,
                      width: 35,
                    ),
                  ],
                  Text(
                    asset!.assetCode!.toUpperCase(),
                    style: TextStyle(
                      fontSize: 16,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              if (asset?.fundingStructure == 0 ||
                  asset?.fundingStructure == 2) ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: const BorderRadius.all(
                        Radius.circular(10.0),
                      ),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                            vertical: 10.0,
                            horizontal: 15,
                          ),
                          child: Text(
                            "Dividend Details",
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        item("Dividend frequency", 'Quarterly'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Total dividends received", 'N45,000'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Last dividend amount", 'N7,500'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Last dividend date", '20th May, 2024'),
                        SizedBox(height: height / 50),
                        Button(
                          "View Dividend History",
                          notifier.getbluecolor,
                          wihitecolor,
                          onTap: () {
                            appState.setPage(
                              page: DividendHistoryViewPageConfig,
                            );
                          },
                          width: 300.w,
                        ),
                        SizedBox(height: height / 50),
                      ],
                    ),
                  ),
                ),
              ],
              if (asset?.fundingStructure == 1 ||
                  asset?.fundingStructure == 2) ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: const BorderRadius.all(
                        Radius.circular(10.0),
                      ),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                            vertical: 10.0,
                            horizontal: 15,
                          ),
                          child: Text(
                            "Interest Details",
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        item("Interest frequency", 'Quarterly'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Total interest earned", 'N45,000'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Last payment date", '20th May, 2024'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Accrued interest (unpaid)", 'N45,000'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Next payment date", '20th May, 2024'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("YTD effective yield", '9.20%'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Expected annual yield", '9.20%'),
                        SizedBox(height: height / 50),
                        Button(
                          "View Yield History",
                          notifier.getbluecolor,
                          wihitecolor,
                          onTap: () {
                            appState.setPage(page: YieldHistoryViewPageConfig);
                          },
                          width: 300.w,
                        ),
                        SizedBox(height: height / 50),
                      ],
                    ),
                  ),
                ),
              ],
              if (asset!.isExitWithFiat) ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: const BorderRadius.all(
                        Radius.circular(10.0),
                      ),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                            vertical: 10.0,
                            horizontal: 15,
                          ),
                          child: Text(
                            "Early Exit",
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        item("Eligible for early exit", 'Yes'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Current buyback NAV", 'N450/token'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Exit charges or penalty", '2% early exit fee'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Estimated payout (Net)", 'N450/token'),
                        Divider(color: notifier.getsplashgrey, thickness: 1),
                        item("Processing time", '3 business days'),
                        SizedBox(height: height / 50),
                        Button(
                          "Request Early Exit",
                          notifier.getbluecolor,
                          wihitecolor,
                          onTap: () {
                            appState.setPage(
                              page: RequestEarlyExitViewPageConfig,
                            );
                          },
                          width: 300.w,
                        ),
                        SizedBox(height: height / 50),
                      ],
                    ),
                  ),
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

  Widget item(String key, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 0, horizontal: 15),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            constraints: BoxConstraints(maxWidth: width / 2.36),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
            child: Text(
              key,
              textAlign: TextAlign.start,
              overflow: TextOverflow.visible,
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
