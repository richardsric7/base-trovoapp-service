import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:loadmore/loadmore.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Models/Transaction.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:timeago/timeago.dart' as timeago;
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class PaymentHistory extends StatefulWidget {
  const PaymentHistory({Key? key}) : super(key: key);

  @override
  State<PaymentHistory> createState() => Payment_HistoryState();
}

class Payment_HistoryState extends State<PaymentHistory>
    with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  late RefreshController _refreshController;
  late DataProvider appState;
  List<Wallet>? wallets;
  Wallet? activeWallet;
  dynamic selectedWallet = '';
  var claimedAssets;
  late bool isSharedWallet;
  bool showFilter = false;
  late List<TransactionInfo>? historyData;
  var filterTypesMap = {
    HistoryFilterType.TransactionType: "Transaction type",
    HistoryFilterType.DateRange: "Date range",
    HistoryFilterType.AmountRange: "Amount range",
    HistoryFilterType.Username: "Username",
    HistoryFilterType.FromPublicKey: "From public key",
    HistoryFilterType.ToPublicKey: "To public key",
    HistoryFilterType.Memo: "Memo",
  };

  ScrollController scrollController = new ScrollController();

  HistoryFilterType filterType = HistoryFilterType.TransactionType;

  List<DropdownMenuItem<String>> get walletDropdownItems {
    var dropdownItems = wallets!
        .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  wallet.alias!,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
            value: wallet.publicKey))
        .toList();

    return dropdownItems;
  }

  List<DropdownMenuItem<HistoryFilterType>> get filterTypeDropdownItems {
    List<DropdownMenuItem<HistoryFilterType>> items = [];
    filterTypesMap.forEach((key, value) {
      items.add(
        DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  value,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
            value: key),
      );
    });

    return items;
  }

  List<DropdownMenuItem<String>> get assetsDropdownItems {
    var items = <DropdownMenuItem<String>>[];
    items.add(DropdownMenuItem<String>(
      value: '*|*',
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            "All assets",
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    ));

    items.addAll(claimedAssets.map<DropdownMenuItem<String>>((asset) {
      return DropdownMenuItem<String>(
        value: '${asset['assetIssuer']}|${asset["assetCode"]}',
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              asset["assetCode"].toString().isEmpty
                  ? 'XBN'
                  : asset["assetCode"],
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      );
    }).toList());

    return items;
  }

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    resetFilters();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    wallets = appState.userInfo!.wallets!;
    activeWallet = appState.activeWallet;
    if (activeWallet == null && wallets!.length > 0) {
      activeWallet = wallets![0];
    }
    selectedWallet = activeWallet!.publicKey;
    historyData = appState.historyData;
    var assetBalances = appState.assetBalances;
    isSharedWallet =
        appState.viewData![PaymentHistoryViewPageConfig.key] != null;

    // if this page is viewed from shared wallet then get the claimed assets
    // from viewData
    claimedAssets = isSharedWallet
        ? appState.viewData![PaymentHistoryViewPageConfig.key]['claimed']
        : assetBalances[activeWallet!.publicKey]['claimed'];

    return ScreenUtilInit(
      builder: (context, child) => DefaultTabController(
        length: 2,
        child: Scaffold(
          resizeToAvoidBottomInset: false,
          backgroundColor: notifier.getwihitecolor,
          appBar: AppBar(
            centerTitle: true,
            // this part will only appear when we are viewing payment history
            // from shared wallet in which case appState.viewData![PaymentHistoryViewPageConfig.key]
            // will not be null;
            leading: isSharedWallet
                ? GestureDetector(
                    onTap: () {
                      Navigator.of(context).pop();
                      appState.viewData![PaymentHistoryViewPageConfig.key] =
                          null;
                    },
                    child: Image.asset("assets/images/back.png", scale: 5),
                  )
                : null,
            title: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                SizedBox(width: width / 15),
                Text(
                  LanguageEn.transactionHistory,
                  style: TextStyle(
                      color: notifier.getblck, fontFamily: fontsemibold),
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
                    ))
              ],
            ),
            backgroundColor: notifier.getfavorites,
            elevation: 0,
          ),
          body: SmartRefresher(
            enablePullDown: true,
            controller: _refreshController,
            onRefresh: refreshData,
            child: Column(
              children: [
                SizedBox(
                  height: height / 35,
                ),
                if (showFilter) ...[
                  Container(
                    width: width,
                    child: Row(
                      children: [
                        SizedBox(
                          width: width / 50,
                        ),
                        // hide the dropdown when we view this page from shared
                        // wallet
                        if (!isSharedWallet) ...[
                          Expanded(
                            flex: 2,
                            child: dropdown(
                              (newValue) async {
                                selectedWallet = newValue!;
                                appState.filterAsset = "*|*";
                                appState.activeWallet = wallets!.firstWhere(
                                    (wallet) => wallet.publicKey == newValue);
                                showLoader(context);
                                appState.limit = 20;
                                appState.totalRecords = 0;
                                appState.currentPage = 1;
                                await appState.getHistory(
                                  context,
                                  onDone: () => adjustScrollPosition(),
                                );
                                hideLoader(context);

                                if (mounted) {
                                  setState(() {});
                                }
                              },
                              walletDropdownItems,
                              selectedWallet,
                              null,
                              context,
                              null,
                            ),
                          ),
                        ],
                        Expanded(
                          flex: 2,
                          child: dropdown(
                            (newValue) async {
                              print(newValue);
                              showLoader(context);
                              appState.limit = 20;
                              appState.totalRecords = 0;
                              appState.currentPage = 1;
                              appState.setFilterAsset = newValue.toString();
                              await appState.getHistory(
                                context,
                                onDone: () => adjustScrollPosition(),
                              );
                              hideLoader(context);
                            },
                            assetsDropdownItems,
                            appState.filterAsset,
                            'Assets',
                            context,
                            null,
                          ),
                        ),
                        SizedBox(
                          width: width / 50,
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: height / 50),
                  Container(
                    width: width,
                    child: Row(
                      children: [
                        SizedBox(
                          width: width / 50,
                        ),
                        Expanded(
                          flex: 2,
                          child: dropdown(
                            (newValue) async {
                              setState(() {
                                filterType = newValue as HistoryFilterType;
                                showPopup(newValue);
                              });
                            },
                            filterTypeDropdownItems,
                            null,
                            filterTypesMap[filterType],
                            context,
                            null,
                          ),
                        ),
                        Expanded(
                          flex: 2,
                          child: getContent(filterType),
                        ),
                        SizedBox(
                          width: width / 50,
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: height / 35),
                ],
                listHistory(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget listHistory() {
    if (historyData != null && historyData!.length > 0) {
      return Container(
        height: isSharedWallet
            ? (showFilter ? height / 1.3950 : height / 1.14)
            : (showFilter ? height / 1.5523 : height / 1.24),
        child: LoadMore(
          isFinish: historyData!.length == appState.totalRecords,
          onLoadMore: () async {
            appState.limit += 20;
            await appState.getHistory(context);
            return historyData!.length <= appState.totalRecords!;
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
              // case LoadMoreStatus.loading:
              //   text = "Loading";
              //   break;
              // case LoadMoreStatus.nomore:
              //   text = "";
              //   break;
              default:
                text = "";
            }
            return text;
          },
          child: ListView.separated(
              separatorBuilder: (context, int) => Container(),
              itemCount: historyData!.length,
              controller: scrollController,
              itemBuilder: (context, index) {
                return tile(historyData![index]);
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
                await appState.getHistory(context);
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

  Widget tile(TransactionInfo transaction) {
    // lets start by setting transactionType to receive
    TransactionType transactionType = TransactionType.Receive;
    var amount = transaction.amount;
    var assetCode = transaction.assetCode;
    var date = transaction.transactionDate;
    var name =
        '${LanguageEn.receivedfrom} ${extractUsername(transaction.from!) ?? truncate(transaction.fromPublicKey!)}';

    // if record.from is same as the current active wallet public key
    // then it was a send transaction
    if (transaction.fromPublicKey == activeWallet!.publicKey) {
      transactionType = TransactionType.Send;
      name =
          '${LanguageEn.sentto} ${extractUsername(transaction.to!) ?? truncate(transaction.toPublicKey!)}';
    }

    if (transaction.transactionType!.contains('SWAP')) {
      transactionType = TransactionType.Swap;
      var splitResult =
          transaction.transactionType!.replaceAll('SWAP', '').trim().split('>');
      name = "Swapped ${splitResult[0]} to ${splitResult[1]}";
    }

    return GestureDetector(
      onTap: () {
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: PaymentDetailsViewPageConfig,
        );

        appState.viewData![PaymentDetailsViewPageConfig.key] = transaction;
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
                  getIcon(transactionType),
                  width: width / 12,
                  color: notifier.getbluewhitecolor,
                  height: 25,
                ),
                SizedBox(
                  width: width / 50,
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    SizedBox(
                      width: width / 50,
                    ),
                    Text(
                      formatAmount(transactionType, amount, assetCode),
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: transactionType == TransactionType.Send
                            ? Colors.red
                            : notifier.getgreencolor,
                        fontFamily: fontbody,
                      ),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        Text(
                          timeago.format(date!),
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
                    Row(
                      children: [
                        Text(
                          name,
                          textAlign: TextAlign.center,
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
              ],
            ),
          ),
        ),
      ),
    );
  }

  String getIcon(TransactionType transactionType) {
    switch (transactionType) {
      case TransactionType.Swap:
        return "assets/images/swap.png";
      case TransactionType.Send:
        return 'assets/images/send.png';
      default:
        return 'assets/images/receive.png';
    }
  }

  String formatAmount(TransactionType transactionType, amount, assetCode) {
    var am = formatHistoryNumber(double.parse(amount.toString()));
    return transactionType == TransactionType.Send
        ? '- $am $assetCode'
        : '+ $am $assetCode';
  }

  refreshData() async {
    try {
      showLoader(context);
      await appState.getHistory(
        context,
        onDone: () => adjustScrollPosition(),
      );
      hideLoader(context);
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  String? extractUsername(String data) {
    print('data $data');
    if (data.isNotEmpty) {
      const start = '[';
      const end = ']';
      final startIndex = data.indexOf(start);
      final endIndex = data.indexOf(end);
      return data.substring(startIndex + start.length, endIndex);
    }
    return null;
  }

  Widget getContent(HistoryFilterType type) {
    switch (type) {
      case HistoryFilterType.Username:
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
                textFieldPopup(context, rel: HistoryFilterType.Username,
                    onDone: (value) async {
                  appState.setFilterUsername = value;
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                    await appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      appState.filterUsername == null
                          ? "Enter username"
                          : appState.filterUsername!,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterUsername != null ? 12 : 15,
                          fontFamily: fontsemibold),
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
      case HistoryFilterType.Memo:
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
                textFieldPopup(context, rel: HistoryFilterType.Memo,
                    onDone: (value) async {
                  appState.setFilterMemo = value;
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&memo=${value}";
                    await appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      appState.filterMemo == null
                          ? "Enter memo"
                          : truncate(appState.filterMemo!, length: 30),
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterMemo != null ? 12 : 15,
                          fontFamily: fontsemibold),
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
      case HistoryFilterType.FromPublicKey:
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
                textFieldPopup(context, rel: HistoryFilterType.FromPublicKey,
                    onDone: (value) async {
                  if (value != null && value.toString().isNotEmpty) {
                    appState.setFilterFromPublicKey = value;
                    appState.setFilterQuery = "&fromPublicKey=$value";
                    await appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      getTruncatedPublicKey(appState.filterFromPublicKey),
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize:
                              appState.filterFromPublicKey != null ? 12 : 15,
                          fontFamily: fontsemibold),
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
      case HistoryFilterType.ToPublicKey:
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
                textFieldPopup(context, rel: HistoryFilterType.ToPublicKey,
                    onDone: (value) async {
                  if (value != null && value.toString().isNotEmpty) {
                    appState.setFilterToPublicKey = value;
                    appState.setFilterQuery = "&toPublicKey=$value";
                    await appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      getTruncatedPublicKey(appState.filterToPublicKey),
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize:
                              appState.filterToPublicKey != null ? 12 : 15,
                          fontFamily: fontsemibold),
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
      case HistoryFilterType.AmountRange:
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
                amountRangePopup(context, onDone: () async {
                  if (appState.filterMinAmount != null &&
                      appState.filterMaxAmount != null) {
                    appState.setFilterQuery =
                        "&amount=${appState.filterMinAmount}%7C${appState.filterMaxAmount}";
                    await appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      truncate(getAmountRangeValue(), length: 30),
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterMinAmount != null &&
                                  appState.filterMaxAmount != null
                              ? 12
                              : 15,
                          fontFamily: fontsemibold),
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
      case HistoryFilterType.DateRange:
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
                customDateRangePopup(context, onDone: () async {
                  appState.setFilterQuery =
                      "&dateBetween=${DateFormat('yyyy-MM-dd').format(appState.filterStartDate!)}%7C${DateFormat('yyyy-MM-dd').format(appState.filterEndDate!)}";
                  await appState.getHistory(
                    context,
                    onDone: () => adjustScrollPosition(),
                  );
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    getDateRangeValue(),
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
      // HistoryFilterType.TransactionType
      default:
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
                transactionTypePopup(
                  context,
                  onAllSelected: () {
                    appState.setFilterQuery = "";
                    appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                    Navigator.of(context).pop(); // dismiss dialog,
                  },
                  onPaymentSelected: () {
                    appState.setFilterQuery = "&transactionType=payment";
                    appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                    Navigator.of(context).pop(); // dismiss dialog,
                  },
                  onSwapSelected: () {
                    appState.setFilterQuery = "&transactionType=swap";
                    appState.getHistory(
                      context,
                      onDone: () => adjustScrollPosition(),
                    );
                    Navigator.of(context).pop(); // dismiss dialog,
                  },
                );
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    getTransactionTypeValue(),
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
  }

  adjustScrollPosition() {
    if (scrollController.hasClients)
      scrollController.jumpTo(scrollController.position.minScrollExtent);
  }

  getDateRangeValue() {
    if (appState.filterStartDate != null && appState.filterEndDate != null) {
      return "${DateFormat('dd/MM/yy').format(appState.filterStartDate!)} - ${DateFormat('dd/MM/yy').format(appState.filterEndDate!)} ";
    }

    return 'Enter range';
  }

  getAmountRangeValue() {
    if (appState.filterMinAmount != null && appState.filterMaxAmount != null) {
      return "${appState.filterMinAmount} - ${appState.filterMaxAmount} ";
    }

    return 'Enter range';
  }

  getTruncatedPublicKey(String? publicKey) {
    if (publicKey == null) return "Enter public key";
    if (publicKey.length <= 7) return publicKey;
    return truncate(publicKey, length: 7) +
        publicKey.substring(publicKey.length - 7);
  }

  getTransactionTypeValue() {
    if (filterType == HistoryFilterType.TransactionType) {
      if (appState.filterQuery.contains('swap')) return "Swap";
      if (appState.filterQuery.contains('payment')) return "Payment";

      return "All";
    }
  }

  void showPopup(HistoryFilterType filterType) {
    switch (filterType) {
      case HistoryFilterType.Username:
        textFieldPopup(context, rel: HistoryFilterType.Username,
            onDone: (value) async {
          print('timer fired! $value');
          appState.setFilterUsername = value;
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
            await appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
          }
        });
        break;
      case HistoryFilterType.FromPublicKey:
        textFieldPopup(context, rel: HistoryFilterType.FromPublicKey,
            onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterFromPublicKey = value;
            appState.setFilterQuery = "&fromPublicKey=$value";
            await appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
          }
        });
        break;
      case HistoryFilterType.ToPublicKey:
        textFieldPopup(context, rel: HistoryFilterType.ToPublicKey,
            onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterToPublicKey = value;
            appState.setFilterQuery = "&toPublicKey=$value";
            await appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
          }
        });
        break;
      case HistoryFilterType.AmountRange:
        amountRangePopup(context, onDone: () async {
          if (appState.filterMinAmount != null &&
              appState.filterMaxAmount != null) {
            appState.setFilterQuery =
                "&amount=${appState.filterMinAmount}%7C${appState.filterMaxAmount}";
            await appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
          }
        });
        break;
      case HistoryFilterType.DateRange:
        customDateRangePopup(context, onDone: () async {
          appState.setFilterQuery =
              "&dateBetween=${DateFormat('yyyy-MM-dd').format(appState.filterStartDate!)}%7C${DateFormat('yyyy-MM-dd').format(appState.filterEndDate!)}";
          await appState.getHistory(
            context,
            onDone: () => adjustScrollPosition(),
          );
        });
        break;
      case HistoryFilterType.Memo:
        textFieldPopup(context, rel: HistoryFilterType.Memo,
            onDone: (value) async {
          appState.setFilterMemo = value;
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&memo=${value}";
            await appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
          }
        });
        break;
      default:
        // appState.setFilterQuery = "";
        // appState.setFilterAsset = "*|*";
        // appState.getHistory(context);
        transactionTypePopup(
          context,
          onAllSelected: () {
            appState.setFilterQuery = "";
            appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
            Navigator.of(context).pop(); // dismiss dialog,
          },
          onPaymentSelected: () {
            appState.setFilterQuery = "&transactionType=payment";
            appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
            Navigator.of(context).pop(); // dismiss dialog,
          },
          onSwapSelected: () {
            appState.setFilterQuery = "&transactionType=swap";
            appState.getHistory(
              context,
              onDone: () => adjustScrollPosition(),
            );
            Navigator.of(context).pop(); // dismiss dialog,
          },
        );
        break;
    }
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
    print('disposing...');
    appState.viewData![PaymentHistoryViewPageConfig.key] = null;
    resetFilters();
  }
}

enum TransactionType {
  Send,
  Receive,
  Swap,
}

enum HistoryFilterType {
  TransactionType,
  DateRange,
  AmountRange,
  Username,
  FromPublicKey,
  ToPublicKey,
  Memo,
}
