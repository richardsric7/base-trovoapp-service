import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart' show DateFormat, NumberFormat;
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/proceed_payout.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/payout_requests.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

// The payouts (dividends, or interest for YieldHistoryView) of the asset
// in viewData to the user's wallets: paid ones, and scheduled ones waiting
// to be paid.
class DividendHistoryView extends StatefulWidget {
  final String title;
  const DividendHistoryView({Key? key, this.title = 'Dividend History'})
    : super(key: key);

  @override
  State<DividendHistoryView> createState() => _DividendHistoryView();
}

class _DividendHistoryView extends State<DividendHistoryView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  Asset? asset;
  List<ProceedPayoutReceipt>? receipts;
  bool failed = false;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (appState.viewData!['walletAddress'] != null) {
      Wallet wallet = appState.userInfo!.getWallet(
        appState.viewData!['walletAddress'],
      );
      for (final a in wallet.claimedAssets ?? <Asset>[]) {
        if (a.assetCode == appState.viewData!['assetCode'] &&
            a.contractAddress == appState.viewData!['contractAddress']) {
          asset = a;
        }
      }
    }
    _load();
  }

  Future<void> _load() async {
    final list = await PayoutApi(appState).receipts(
      assetCode: appState.viewData!['assetCode'],
      tokenContract: appState.viewData!['contractAddress'],
    );
    if (!mounted) return;
    setState(() {
      receipts = list ?? [];
      failed = list == null;
    });
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    Widget body;
    if (receipts == null) {
      body = Padding(
        padding: EdgeInsets.only(top: height / 4),
        child: const CircularProgressIndicator(),
      );
    } else if (failed || receipts!.isEmpty) {
      body = Padding(
        padding: EdgeInsets.only(top: height / 4, left: 30, right: 30),
        child: Text(
          failed
              ? 'The payouts could not be loaded. Pull down to try again.'
              : 'No payouts yet.',
          textAlign: TextAlign.center,
          style: TextStyle(
            color: notifier.getbluewhitecolor,
            fontSize: 14.sp,
            fontFamily: fontbody,
          ),
        ),
      );
    } else {
      body = Column(children: receipts!.map(item).toList());
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          widget.title,
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: RefreshIndicator(
          onRefresh: _load,
          child: SingleChildScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            child: Column(
              children: [
                SizedBox(height: height / 50),
                body,
                SizedBox(height: height / 20),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String _status(ProceedPayoutReceipt p) {
    switch (p.status) {
      case 'PAID':
        return '';
      case 'FAILED':
        return 'Could not be paid yet';
      default:
        return 'Scheduled';
    }
  }

  Widget item(ProceedPayoutReceipt p) {
    final amount =
        '${p.paid ? '+' : ''}${NumberFormat('#,##0.00######').format(p.amountValue)} ${p.payoutAssetCode}';
    final date = p.date == null
        ? ''
        : DateFormat('d MMM, yyyy  hh:mm a').format(p.date!);
    final status = _status(p);
    return InkWell(
      onTap: () {
        appState.viewData!['payoutReceipt'] = p.toMap();
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
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        amount,
                        style: TextStyle(
                          fontWeight: FontWeight.w500,
                          color: p.paid
                              ? notifier.getgreencolor
                              : notifier.getbluewhitecolor,
                          fontSize: 13.sp,
                          fontFamily: fontsemibold,
                        ),
                      ),
                      Text(
                        [
                          date,
                          if (p.walletAlias.isNotEmpty) p.walletAlias,
                          if (status.isNotEmpty) status,
                        ].join('  ·  '),
                        style: TextStyle(
                          fontWeight: FontWeight.w500,
                          color: notifier.getbluewhitecolor,
                          fontSize: 10.sp,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
