import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:loadmore/loadmore.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/deposit_transaction_model.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/models/withdrawal_transaction_model.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:timeago/timeago.dart' as timeago;
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DepositWithdrawHistory extends StatefulWidget {
  const DepositWithdrawHistory({Key? key}) : super(key: key);

  @override
  State<DepositWithdrawHistory> createState() => _DepositWithdrawHistoryState();
}

class _DepositWithdrawHistoryState extends State<DepositWithdrawHistory>
    with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  late RefreshController _refreshController;
  late DataProvider appState;
  late Wallet wallet;
  late Asset asset;
  var claimedAssets;
  bool showFilter = false;
  String historyMode = 'Deposit history';
  late List<DepositTransactionModel>? depositHistory;
  late List<WithdrawalTransactionModel>? withdrawalHistory;

  ScrollController scrollController = new ScrollController();

  HistoryFilterType filterType = HistoryFilterType.TransactionDirection;

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    resetFilters();
    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    depositHistory = appState.depositHistoryData;
    withdrawalHistory = appState.withdrawalHistoryData;

    return ScreenUtilInit(
      builder: (context, child) => DefaultTabController(
        length: 2,
        child: Scaffold(
          resizeToAvoidBottomInset: false,
          backgroundColor: notifier.getwihitecolor,
          appBar: PreferredSize(
            preferredSize: Size.fromHeight(height / 15),
            child: AppBar(
              centerTitle: true,
              elevation: 0,
              backgroundColor: notifier.getwihitecolor,
              leading: GestureDetector(
                onTap: () {
                  Navigator.of(context).pop();
                },
                child: Image.asset("assets/images/back.png", scale: 5),
              ),
              title: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    width: width / 1.7,
                    child: getContent(filterType),
                  ),
                  TextButton(
                    onPressed: () {
                      setState(() {
                        showFilter = !showFilter;
                      });
                    },
                    child: Container(
                      child: Image.asset(
                        "assets/images/filter-list.png",
                        height: height / 35,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  )
                ],
              ),
            ),
          ),
          body: SmartRefresher(
            enablePullDown: true,
            controller: _refreshController,
            onRefresh: refreshData,
            child: Column(
              children: [
                if (showFilter) ...[
                  // Container(
                  //   width: width,
                  //   child: Row(
                  //     children: [
                  //       SizedBox(
                  //         width: width / 20,
                  //       ),
                  //       Expanded(
                  //         flex: 2,
                  //         child: getContent(filterType),
                  //       ),
                  //       SizedBox(
                  //         width: width / 20,
                  //       ),
                  //     ],
                  //   ),
                  // ),
                ],
                SizedBox(
                  height: height / 50,
                ),
                listHistory(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget listHistory() {
    if ((historyMode == 'Deposit history' && depositHistory!.length > 0) ||
        (historyMode == 'Withdrawal history' &&
            withdrawalHistory!.length > 0)) {
      return Container(
        height: (showFilter ? height / 1.22 : height / 1.14),
        child: LoadMore(
          isFinish: (historyMode == 'Deposit history' &&
                  depositHistory!.length == appState.totalRecords) ||
              (historyMode == 'Withdrawal history' &&
                  withdrawalHistory!.length == appState.totalRecords),
          onLoadMore: () async {
            appState.limit += 20;
            historyMode == 'Deposit history'
                ? await appState.fetchDepositHistory(
                    context,
                    publicKey: wallet.publicKey!,
                    currency: asset.assetCode!,
                  )
                : await appState.fetchWithdrawalHistory(
                    context,
                    publicKey: wallet.publicKey!,
                    currency: asset.assetCode!,
                  );
            return depositHistory!.length <= appState.totalRecords!;
          },
          textBuilder: (LoadMoreStatus status) {
            String text;
            switch (status) {
              case LoadMoreStatus.fail:
                text = "Tap to load more";
                break;
              case LoadMoreStatus.idle:
                text = "Tap to load more";
                break;
              default:
                text = "";
            }
            return text;
          },
          child: ListView.separated(
              separatorBuilder: (context, int) => Container(),
              itemCount: historyMode == 'Deposit history'
                  ? depositHistory!.length
                  : withdrawalHistory!.length,
              controller: scrollController,
              itemBuilder: (context, index) {
                return historyMode == 'Deposit history'
                    ? depositHistoryTile(depositHistory![index])
                    : withdrawalHistoryTile(withdrawalHistory![index]);
              }),
        ),
      );
    }

    return Container(
      height: height / 1.8,
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Sorry no results here',
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(
              height: height / 90,
            ),
            ElevatedButton(
              onPressed: () async {
                await appState.fetchDepositHistory(
                  context,
                  publicKey: wallet.publicKey!,
                  currency: asset.assetCode,
                  // onDone: () => adjustScrollPosition(),
                );
              },
              style: ButtonStyle(
                backgroundColor:
                    MaterialStateProperty.all<Color>(notifier.getbluecolor!),
              ),
              child: Text(
                'Refresh',
                style: TextStyle(
                  fontFamily: fontsemibold,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget depositHistoryTile(DepositTransactionModel transaction) {
    return GestureDetector(
      onTap: () {
        appState.viewData = {
          'transaction': transaction,
          'transactionDirection': TransactionDirection.Deposit,
          'walletPublicKey': wallet.publicKey,
        };

        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: DepositWithdrawDetailsViewPageConfig,
        );
      },
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
          ),
          child: Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 10.0, vertical: 15.0),
            child: Row(
              children: [
                Image.asset(
                  'assets/images/deposit.png',
                  width: width / 12,
                  color: notifier.getbluewhitecolor,
                  height: 25,
                ),
                SizedBox(
                  width: width / 50,
                ),
                Container(
                  width: width / 1.5,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SizedBox(
                        width: width / 50,
                      ),
                      Text(
                        formatAmount(
                            TransactionDirection.Deposit,
                            transaction.amount.toString(),
                            transaction.currency),
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w400,
                          color: notifier.getgreencolor,
                          fontFamily: fontbody,
                        ),
                      ),
                      Row(
                        // mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          Text(
                            timeago.format(transaction.createdAt),
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(
                        height: 2,
                      ),
                      Wrap(
                        children: [
                          Text(
                            'From address: ${truncatePublicKey(transaction.fromAddress)}',
                            overflow: TextOverflow.visible,
                            // textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
                      ),
                      Wrap(
                        children: [
                          Text(
                            'To address: ${truncatePublicKey(transaction.toAddress)}',
                            overflow: TextOverflow.visible,
                            // textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
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

  Widget withdrawalHistoryTile(WithdrawalTransactionModel transaction) {
    return GestureDetector(
      onTap: () {
        appState.viewData = {
          'transaction': transaction,
          'transactionDirection': TransactionDirection.Withdraw,
          'walletPublicKey': wallet.publicKey,
        };

        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: DepositWithdrawDetailsViewPageConfig,
        );
      },
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
          ),
          child: Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 10.0, vertical: 15.0),
            child: Row(
              children: [
                Image.asset(
                  'assets/images/withdraw.png',
                  width: width / 12,
                  color: notifier.getbluewhitecolor,
                  height: 25,
                ),
                SizedBox(
                  width: width / 50,
                ),
                Container(
                  width: width / 1.5,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SizedBox(
                        width: width / 50,
                      ),
                      Text(
                        formatAmount(
                            TransactionDirection.Withdraw,
                            transaction.amountSubmitted.toString(),
                            transaction.currency),
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w400,
                          color: notifier.getgreencolor,
                          fontFamily: fontbody,
                        ),
                      ),
                      Row(
                        // mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          Text(
                            timeago.format(transaction.createdAt),
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(
                        height: 2,
                      ),
                      Wrap(
                        children: [
                          Text(
                            'Network: ${transaction.withdrawalNetwork}',
                            overflow: TextOverflow.visible,
                            // textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(
                        height: 2,
                      ),
                      Wrap(
                        children: [
                          Text(
                            'Status: ${transaction.withdrawalStatus}',
                            overflow: TextOverflow.visible,
                            // textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
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

  String formatAmount(TransactionDirection transactionType, amount, assetCode) {
    var am = formatHistoryNumber(double.parse(amount.toString()), 1000000);
    return transactionType == TransactionDirection.Withdraw
        ? '- $am $assetCode'
        : '+ $am $assetCode';
  }

  refreshData() async {
    try {
      showLoader(context);
      historyMode == 'Deposit history'
          ? await appState.fetchDepositHistory(
              context,
              publicKey: wallet.publicKey!,
              currency: asset.assetCode!,
              // onDone: () => adjustScrollPosition(),
            )
          : await appState.fetchWithdrawalHistory(
              context,
              publicKey: wallet.publicKey!,
              currency: asset.assetCode!,
              // onDone: () => adjustScrollPosition(),
            );
      hideLoader(context);
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  Widget getContent(HistoryFilterType type) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(10.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: TextButton(
          onPressed: () {
            wrappedAssettransactionTypePopup(
              context,
              onWithdrawSelected: () {
                historyMode = 'Withdrawal history';
                appState.limit = 20;
                appState.fetchWithdrawalHistory(
                  context,
                  publicKey: wallet.publicKey!,
                  currency: asset.assetCode,
                  // onDone: () => adjustScrollPosition(),
                );
                setState(() {});
                Navigator.of(context).pop(); // dismiss dialog,
              },
              onDepositSelected: () {
                appState.limit = 20;
                historyMode = 'Deposit history';
                appState.fetchDepositHistory(
                  context,
                  publicKey: wallet.publicKey!,
                  currency: asset.assetCode,
                  // onDone: () => adjustScrollPosition(),
                );
                setState(() {});
                Navigator.of(context).pop(); // dismiss dialog,
              },
            );
          },
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                historyMode,
                textAlign: TextAlign.start,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: appState.filterStartDate != null &&
                            appState.filterEndDate != null
                        ? 13
                        : 15,
                    fontFamily: fontsemibold),
              ),
              Icon(
                Icons.keyboard_arrow_down_rounded,
                color: notifier.getbluewhitecolor,
              ),
            ],
          ),
        ),
      ),
    );
  }

  adjustScrollPosition() {
    if (scrollController.hasClients)
      scrollController.jumpTo(scrollController.position.minScrollExtent);
  }

  void resetFilters() {
    appState.filterAsset = "*|*";
    appState.filterEndDate = null;
    appState.filterStartDate = null;
    appState.filterFromPublicKey = null;
    appState.filterToPublicKey = null;
    appState.filterUsername = null;
    appState.filterQuery = "";
    appState.filterMaxAmount = null;
    appState.filterMinAmount = null;
    appState.filterMemo = null;
  }

  @override
  void dispose() {
    super.dispose();
    appState.viewData = {};
    resetFilters();
  }
}
