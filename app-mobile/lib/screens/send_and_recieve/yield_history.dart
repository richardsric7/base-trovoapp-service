import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
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

class YieldHistoryView extends StatefulWidget {
  const YieldHistoryView({Key? key}) : super(key: key);

  @override
  State<YieldHistoryView> createState() => _YieldHistoryView();
}

class _YieldHistoryView extends State<YieldHistoryView>
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
          'Interest History',
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              item('+3,000.00 CNGN', '22 Jan, 2024  10:30 AM'),
              item('+3,000.00 CNGN', '22 Jan, 2024  10:30 AM'),
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

  Widget item(String amount, String date) {
    return InkWell(
      onTap: () {
        appState.setPage(page: DividendPaymentDetailsViewPageConfig);
      },
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(10.0)),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 15, horizontal: 15),
            child: Row(
              spacing: 10,
              children: [
                Image.asset(
                  'assets/images/receive.png',
                  color: notifier.getbluewhitecolor,
                  height: 20,
                  width: 20,
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      amount,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontSize: 13.sp,
                        fontFamily: fontsemibold,
                      ),
                    ),
                    Text(
                      date,
                      textAlign: TextAlign.end,
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontSize: 10.sp,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
