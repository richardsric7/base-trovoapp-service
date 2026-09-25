import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:loadmore/loadmore.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/deposit_transaction_model.dart';
import 'package:trovo_app/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/models/withdrawal_transaction_model.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  var filterTypesMap = {
    FilterType.DateRange: "Date range",
    FilterType.AmountRange: "Amount range",
    FilterType.WithdrawalAddress: "Withdrawal address",
    FilterType.WithdrawalStatus: "Withdrawal status",
  };

  ScrollController scrollController = new ScrollController();
  FilterType filterType = FilterType.DateRange;

  List<DropdownMenuItem<FilterType>> get filterTypeDropdownItems {
    List<DropdownMenuItem<FilterType>> items = [];
    filterTypesMap.forEach((key, value) {
      if (historyMode == 'Deposit history') {
        // filter out withdrawal related filters
        if (key != FilterType.WithdrawalAddress &&
            key != FilterType.WithdrawalStatus) {
          items.add(
            DropdownMenuItem(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [Text(value, overflow: TextOverflow.ellipsis)],
              ),
              value: key,
            ),
          );
        }
      } else {
        items.add(
          DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [Text(value, overflow: TextOverflow.ellipsis)],
            ),
            value: key,
          ),
        );
      }
    });

    return items;
  }

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    resetFilters();
    wallet = appState.userInfo!.getWallet(appState.viewData!['walletAddress']);

    if (appState.viewData!['historyMode'] != null) {
      historyMode = appState.viewData!['historyMode'];
    }

    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.contractAddress == appState.viewData!['contractAddress'],
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
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 5.0),
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius: const BorderRadius.all(
                            Radius.circular(10.0),
                          ),
                          color: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                        ),
                        child: TextButton(
                          onPressed: () {
                            wrappedAssetTransactionTypePopup(
                              context,
                              onWithdrawSelected: () {
                                historyMode = 'Withdrawal history';
                                filterType = FilterType.DateRange;
                                resetFilters();
                                appState.limit = 20;
                                appState.fetchWithdrawalHistory(
                                  context,
                                  address: wallet.address!,
                                  currency: asset.assetCode,
                                  // onDone: () => adjustScrollPosition(),
                                );
                                setState(() {});
                                Navigator.of(context).pop(); // dismiss dialog,
                              },
                              onDepositSelected: () {
                                appState.limit = 20;
                                historyMode = 'Deposit history';
                                resetFilters();
                                appState.fetchDepositHistory(
                                  context,
                                  address: wallet.address!,
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
                                  fontSize:
                                      appState.filterStartDate != null &&
                                          appState.filterEndDate != null
                                      ? 13
                                      : 15,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                              Icon(
                                Icons.keyboard_arrow_down_rounded,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
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
                  ),
                ],
              ),
              actions: [
                if (appState.walletMode == "Testnet") ...[
                  Visibility(
                    visible: true,
                    child: Padding(
                      padding: EdgeInsets.only(top: 5),
                      child: Banner(
                        location: BannerLocation.topEnd,
                        message: "testnet".tr(),
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
          body: SmartRefresher(
            enablePullDown: true,
            controller: _refreshController,
            onRefresh: refreshData,
            child: Column(
              children: [
                if (showFilter) ...[
                  SizedBox(height: height / 50),
                  Container(
                    width: width,
                    child: Row(
                      children: [
                        SizedBox(width: width / 50),
                        Expanded(
                          flex: 2,
                          child: dropdown(
                            (newValue) async {
                              setState(() {
                                filterType = newValue as FilterType;
                                showPopup(filterType);
                              });
                            },
                            filterTypeDropdownItems,
                            null,
                            filterTypesMap[filterType],
                            context,
                            null,
                          ),
                        ),
                        Expanded(flex: 2, child: getContent(filterType)),
                        SizedBox(width: width / 50),
                      ],
                    ),
                  ),
                ],
                SizedBox(height: height / 50),
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
        height: (showFilter ? height / 1.24 : height / 1.14),
        child: LoadMore(
          isFinish:
              (historyMode == 'Deposit history' &&
                  depositHistory!.length == appState.totalRecords) ||
              (historyMode == 'Withdrawal history' &&
                  withdrawalHistory!.length == appState.totalRecords),
          onLoadMore: () async {
            appState.limit += 20;
            historyMode == 'Deposit history'
                ? await appState.fetchDepositHistory(
                    context,
                    address: wallet.address!,
                    currency: asset.assetCode!,
                  )
                : await appState.fetchWithdrawalHistory(
                    context,
                    address: wallet.address!,
                    currency: asset.assetCode!,
                  );
            return depositHistory!.length <= appState.totalRecords;
          },
          textBuilder: (LoadMoreStatus status) {
            String text;
            switch (status) {
              case LoadMoreStatus.fail:
                text = "taptoloadmore".tr();
                break;
              case LoadMoreStatus.idle:
                text = "taptoloadmore".tr();
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
            },
          ),
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
              "sorrynoresults".tr(),
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(height: height / 90),
            ElevatedButton(
              onPressed: () async {
                await appState.fetchDepositHistory(
                  context,
                  address: wallet.address!,
                  currency: asset.assetCode,
                  // onDone: () => adjustScrollPosition(),
                );
              },
              style: ButtonStyle(
                backgroundColor: WidgetStateProperty.all<Color>(
                  notifier.getbluecolor!,
                ),
                foregroundColor: WidgetStateProperty.all<Color>(
                  notifier.getwihitecolor,
                ),
              ),
              child: Text(
                "refresh".tr(),
                style: TextStyle(fontFamily: fontsemibold),
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
          'walletAddress': wallet.address,
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
            padding: const EdgeInsets.symmetric(
              horizontal: 10.0,
              vertical: 15.0,
            ),
            child: Row(
              children: [
                Image.asset(
                  'assets/images/deposit.png',
                  width: width / 12,
                  color: notifier.getbluewhitecolor,
                  height: 25,
                ),
                SizedBox(width: width / 50),
                Container(
                  width: width / 1.5,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SizedBox(width: width / 50),
                      Text(
                        formatAmount(
                          TransactionDirection.Deposit,
                          transaction.amount.toString(),
                          transaction.currency,
                        ),
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
                      SizedBox(height: 2),
                      Wrap(
                        children: [
                          Text(
                            '${"fromaddress".tr()}: ${truncateAddress(transaction.fromAddress)}',
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
                            '${"toaddress".tr()}: ${truncateAddress(transaction.toAddress)}',
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
          'walletAddress': wallet.address,
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
            padding: const EdgeInsets.symmetric(
              horizontal: 10.0,
              vertical: 15.0,
            ),
            child: Row(
              children: [
                Image.asset(
                  'assets/images/withdraw.png',
                  width: width / 12,
                  color: notifier.getbluewhitecolor,
                  height: 25,
                ),
                SizedBox(width: width / 50),
                Container(
                  width: width / 1.5,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SizedBox(width: width / 50),
                      Text(
                        formatAmount(
                          TransactionDirection.Withdraw,
                          transaction.amountSubmitted.toString(),
                          transaction.currency,
                        ),
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w400,
                          color: notifier.getgreencolor,
                          fontFamily: fontbody,
                        ),
                      ),
                      Row(
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
                      SizedBox(height: 2),
                      Wrap(
                        children: [
                          Text(
                            '${"network".tr()}: ${transaction.withdrawalNetwork}',
                            overflow: TextOverflow.visible,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: 2),
                      Wrap(
                        children: [
                          Text(
                            '${"status".tr()}: ${transaction.withdrawalStatus}',
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
              address: wallet.address!,
              currency: asset.assetCode!,
              // onDone: () => adjustScrollPosition(),
            )
          : await appState.fetchWithdrawalHistory(
              context,
              address: wallet.address!,
              currency: asset.assetCode!,
              // onDone: () => adjustScrollPosition(),
            );
      hideLoader(context);
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  Widget getContent(FilterType type) {
    var text = '';
    switch (type) {
      case FilterType.AmountRange:
        text = getAmountRangeValue();
        break;
      case FilterType.DateRange:
        text = getDateRangeValue();
        break;
      case FilterType.WithdrawalStatus:
        text = appState.filterQuery.contains('withdrawalStatus')
            ? appState.filterWithdrawalStatus
            : "choosestatus".tr();
        break;
      case FilterType.WithdrawalAddress:
        text = appState.filterQuery.contains('withdrawalAddress')
            ? truncateAddress(appState.filterWithdrawalAddress)
            : "choosestatus".tr();
        break;
      default:
    }

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
            showPopup(type);
          },
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                text,
                textAlign: TextAlign.start,
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontSize:
                      appState.filterStartDate != null &&
                          appState.filterEndDate != null
                      ? 13
                      : 15,
                  fontFamily: fontsemibold,
                ),
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

  void showPopup(FilterType filterType) async {
    switch (filterType) {
      case FilterType.WithdrawalAddress:
        wrappedAssetsTextFieldPopup(
          context,
          rel: FilterType.WithdrawalAddress,
          onDone: (value) async {
            appState.setFilterUsername = value;
            if (value != null && value.isNotEmpty) {
              appState.limit = 20;
              appState.setFilterWithdrawalAddress = value;
              appState.setFilterQuery = "&withdrawalAddress=${value}";
              await appState.fetchWithdrawalHistory(
                context,
                address: wallet.address!,
                currency: asset.assetCode,
                // onDone: () => adjustScrollPosition(),
              );
            }
          },
        );
        break;
      case FilterType.WithdrawalStatus:
        wrappedAssetTransactionStatusPopup(
          context,
          onPendingSelected: () {
            appState.limit = 20;
            appState.setFilterWithdrawalStatus = 'PENDING';
            appState.setFilterQuery = "&withdrawalStatus=PENDING";
            appState.fetchWithdrawalHistory(
              context,
              address: wallet.address!,
              currency: asset.assetCode,
              // onDone: () => adjustScrollPosition(),
            );
            setState(() {});
            Navigator.of(context).pop(); // dismiss dialog,
          },
          onCompletedSelected: () {
            appState.limit = 20;
            appState.setFilterWithdrawalStatus = 'COMPLETED';
            appState.setFilterQuery = "&withdrawalStatus=COMPLETED";
            appState.fetchWithdrawalHistory(
              context,
              address: wallet.address!,
              currency: asset.assetCode,
              // onDone: () => adjustScrollPosition(),
            );
            setState(() {});
            Navigator.of(context).pop(); // dismiss dialog,
          },
        );
        break;
      case FilterType.AmountRange:
        amountRangePopup(
          context,
          onDone: () async {
            if (appState.filterMinAmount != null &&
                appState.filterMaxAmount != null) {
              appState.limit = 20;
              appState.setFilterQuery =
                  "&amount=${appState.filterMinAmount}%7C${appState.filterMaxAmount}";
              historyMode == 'Deposit history'
                  ? await appState.fetchDepositHistory(
                      context,
                      address: wallet.address!,
                      currency: asset.assetCode!,
                    )
                  : await appState.fetchWithdrawalHistory(
                      context,
                      address: wallet.address!,
                      currency: asset.assetCode!,
                    );
            }
          },
        );
        break;
      case FilterType.DateRange:
        customDateRangePopup(
          context,
          onDone: () async {
            appState.limit = 20;
            appState.setFilterQuery =
                "&dateBetween=${DateFormat('yyyy-MM-dd').format(appState.filterStartDate!)}%7C${DateFormat('yyyy-MM-dd').format(appState.filterEndDate!)}";
            historyMode == 'Deposit history'
                ? await appState.fetchDepositHistory(
                    context,
                    address: wallet.address!,
                    currency: asset.assetCode!,
                  )
                : await appState.fetchWithdrawalHistory(
                    context,
                    address: wallet.address!,
                    currency: asset.assetCode!,
                  );
          },
        );
        break;
      default:
    }
  }

  getDateRangeValue() {
    if (appState.filterStartDate != null && appState.filterEndDate != null) {
      return "${DateFormat('dd/MM/yy').format(appState.filterStartDate!)} - ${DateFormat('dd/MM/yy').format(appState.filterEndDate!)} ";
    }

    return "enterrange".tr();
  }

  getAmountRangeValue() {
    if (appState.filterMinAmount != null && appState.filterMaxAmount != null) {
      return "${appState.filterMinAmount} - ${appState.filterMaxAmount} ";
    }

    return "enterrange".tr();
  }

  adjustScrollPosition() {
    if (scrollController.hasClients)
      scrollController.jumpTo(scrollController.position.minScrollExtent);
  }

  void resetFilters() {
    filterType = FilterType.DateRange;
    appState.filterEndDate = null;
    appState.filterStartDate = null;
    appState.filterWithdrawalAddress = '';
    appState.filterWithdrawalStatus = '';
    appState.filterQuery = "";
    appState.filterMaxAmount = null;
    appState.filterMinAmount = null;
  }

  @override
  void dispose() {
    super.dispose();
    appState.viewData = {};
    resetFilters();
  }
}

enum FilterType {
  WithdrawalStatus,
  WithdrawalAddress,
  DateRange,
  AmountRange,
  Reset, // just to reset the filter controls
}
