import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:loadmore/loadmore.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:sembast/sembast.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
//
import 'package:trovo_wallet/Models/Transaction.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:timeago/timeago.dart' as timeago;

import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../storage/store.dart';
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
  int totalRecords = 0;
  int currentPage = 1;
  int limit = 20;
  List<TransactionInfo>? historyData = <TransactionInfo>[];

  List<DropdownMenuItem<String>> get walletDropdownItems {
    return wallets!
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
  }

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    wallets = appState.userInfo!.wallets!;
    activeWallet = appState.activeWallet;
    selectedWallet = activeWallet!.publicKey;
    getHistory();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    wallets = appState.userInfo!.wallets!;
    activeWallet = appState.activeWallet;
    selectedWallet = activeWallet!.publicKey;

    return ScreenUtilInit(
      builder: (context, child) => DefaultTabController(
        length: 2,
        child: Scaffold(
          resizeToAvoidBottomInset: false,
          backgroundColor: notifier.getwihitecolor,
          appBar: AppBar(
            centerTitle: true,
            title: Text(
              LanguageEn.transactionHistory,
              style:
                  TextStyle(color: notifier.getblck, fontFamily: fontsemibold),
            ),
            backgroundColor: notifier.getfavorites,
            elevation: 0,
          ),
          body: SmartRefresher(
            enablePullDown: true,
            controller: _refreshController,
            onRefresh: refreshData,
            child: SingleChildScrollView(
              child: Column(
                children: [
                  SizedBox(
                    height: height / 50,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: DropdownButtonFormField(
                      isDense: true,
                      isExpanded: true,
                      decoration: InputDecoration(
                        contentPadding:
                            EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                          borderRadius: BorderRadius.circular(10),
                        ),
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                          borderRadius: BorderRadius.circular(10),
                        ),
                        filled: true,
                        fillColor: notifier.getaddsubwalletgrey,
                      ),
                      value: selectedWallet,
                      icon: Icon(
                        Icons.keyboard_arrow_down_rounded,
                        color: notifier.getbluecolor,
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluecolor,
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: (newValue) async {
                        selectedWallet = newValue!;
                        appState.activeWallet = wallets!.firstWhere(
                            (wallet) => wallet.publicKey == newValue);
                        showLoader(context);
                        limit = 20;
                        totalRecords = 0;
                        currentPage = 1;
                        await fetchHistory(limit);
                        hideLoader(context);
                        if (mounted) {
                          setState(() {});
                        }
                      },
                      items: walletDropdownItems,
                    ),
                  ),
                  // walletTile(),
                  SizedBox(height: height / 50),
                  listHistory(),
                  SizedBox(
                    height: height / 50,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget listHistory() {
    if (historyData != null && historyData!.length > 0) {
      return Container(
        height: height / 1.343,
        // color: Colors.black,
        child: LoadMore(
          isFinish: historyData!.length == totalRecords,
          onLoadMore: () async {
            // setState(() {
            limit += 20;
            // });
            await fetchHistory(limit);
            return historyData!.length <= totalRecords;
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
              // child: ListView.builder(
              itemCount: historyData!.length,
              itemBuilder: (context, index) {
                print(
                    'historyData.length: ${historyData!.length}, index: $index');
                // return (index + 1 == historyData!.length)
                //     ? Container(
                //         color: Colors.greenAccent,
                //         child: TextButton(
                //           child: Text("Load More"),
                //           onPressed: () {},
                //         ),
                //       )
                //     : tile(historyData![index]);
                return tile(historyData![index]);
              }),
        ),
      );
    }

    return Center(
      child: Text(
        'You do not have any transactions yet.',
        style: TextStyle(
          color: notifier.getbluecolor,
          fontSize: 15,
          fontFamily: fontsemibold,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Padding walletTile() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getbluecolor,
        ),
        child: Stack(children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 15.0, horizontal: 20),
                child: Image.asset('assets/images/trovo_white.png'),
              ),
            ],
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  activeWallet!.alias!.capitalizeFirst!,
                  style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    Text(
                      LanguageEn.totalbalance,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getwihitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 98.0,
                ),
                Text(
                  '2,082,898 NGN',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: notifier.getwihitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
                SizedBox(height: 2),
                Text(
                  '4,014 USD',
                  style: TextStyle(
                    fontWeight: FontWeight.w300,
                    fontSize: 13,
                    color: notifier.getwihitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ],
            ),
          ),
        ]),
      ),
    );
  }

  Widget tile(TransactionInfo transaction) {
    // if record.from is same as the current active wallet public key
    // then it was a send transaction otherwise, its a recieve transaction
    TransactionType transactionType =
        transaction.fromPublicKey == activeWallet!.publicKey
            ? TransactionType.Send
            : TransactionType.Receive;
    var name = transactionType == TransactionType.Send
        ? transaction.to
        : transaction.from;

    var publicKey = transactionType == TransactionType.Send
        ? transaction.toPublicKey
        : transaction.fromPublicKey;

    var amount = transaction.amount;

    var assetCode = transaction.assetCode;

    var date = transaction.transactionDate;

    return GestureDetector(
      onTap: () {
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: PaymentDetailsViewPageConfig,
        );

        appState.viewData = {PaymentDetailsViewPageConfig.key: transaction};
      },
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.getaddsubwalletgrey,
          ),
          child: Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
            child: Row(
              children: [
                SvgPicture.asset(
                  transactionType == TransactionType.Send
                      ? 'assets/images/send.svg'
                      : 'assets/images/recieve.svg',
                  width: width / 8,
                  color: transactionType == TransactionType.Send
                      ? notifier.getbluecolor
                      : notifier.getbluecolor,
                  height: 25,
                ),
                SizedBox(
                  width: 15,
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        Text(
                          timeago.format(date!),
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(
                          width: 10,
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
                      ],
                    ),
                    SizedBox(
                      height: 2,
                    ),
                    Row(
                      children: [
                        Text(
                          transactionType == TransactionType.Send
                              ? LanguageEn.sentto
                              : LanguageEn.receivedfrom,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(
                          width: 3,
                        ),
                        Text(
                          name.toString().isEmpty
                              ? truncate(publicKey!)
                              : extractUsername(name!),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: 2),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String formatAmount(TransactionType transactionType, amount, assetCode) {
    var am = formatNumber(double.parse(amount.toString()));
    return transactionType == TransactionType.Send
        ? '- $am $assetCode'
        : '+ $am $assetCode';
  }

  getHistory() async {
    var data =
        await StoreData().storeGetData('historyData${activeWallet!.alias}');
    for (var i = 0; i < data.length; i++) {
      historyData!.add(TransactionInfo().deserializeJson(data[i]));
    }

    totalRecords =
        await StoreData().storeGetData('totalRecords${activeWallet!.alias}');
    currentPage =
        await StoreData().storeGetData('currentPage${activeWallet!.alias}');
    if (mounted) {
      setState(() => {});
    }
  }

  refreshData() async {
    try {
      await fetchHistory(limit);
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  Future<void> fetchHistory(limit) async {
    Map responseData = await makeGetRequest(
        uri:
            '/v1/users/payments/${appState.activeWallet!.publicKey}?limit=$limit',
        signer: appState.activeWallet!.signer!,
        publicKey: appState.activeWallet!.publicKey!,
        secretKey: appState.secretKeys[0]);

    print('response: ${responseData['data']}');
    if (responseData['statusCode'] == 200) {
      totalRecords = responseData['data']['totalRecords'];
      currentPage = responseData['data']['currentPage'];
      var transactions = <TransactionInfo>[];
      for (var i = 0; i < responseData['data']['records'].length; i++) {
        transactions.add(TransactionInfo()
            .deserializeJson(responseData['data']['records'][i]));
      }

      if (limit <= 20) {
        await StoreData().storeInsertData(
            'historyData${appState.activeWallet!.alias}',
            responseData['data']['records']);
        await StoreData().storeInsertData(
            'totalRecords${appState.activeWallet!.alias}',
            responseData['data']['totalRecords']);
        await StoreData().storeInsertData(
            'currentPage${appState.activeWallet!.alias}',
            responseData['data']['currentPage']);
      }

      historyData = transactions;

      if (mounted) {
        setState(() => {});
      }
    }
  }

  String extractUsername(String data) {
    const start = '[';
    const end = ']';
    final startIndex = data.indexOf(start);
    final endIndex = data.indexOf(end);
    return data.substring(startIndex + start.length, endIndex);
  }
}

enum TransactionType {
  Send,
  Receive,
  Swap,
}
