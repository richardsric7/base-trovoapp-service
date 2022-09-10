import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
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
  late List<TransactionInfo>? historyData;

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
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
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
                        fillColor: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      value: selectedWallet,
                      icon: Icon(
                        Icons.keyboard_arrow_down_rounded,
                        color: notifier.getbluewhitecolor,
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: (newValue) async {
                        selectedWallet = newValue!;
                        appState.activeWallet = wallets!.firstWhere(
                            (wallet) => wallet.publicKey == newValue);
                        showLoader(context);
                        appState.limit = 20;
                        appState.totalRecords = 0;
                        appState.currentPage = 1;
                        await appState.fetchHistory(appState.limit);
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
          isFinish: historyData!.length == appState.totalRecords,
          onLoadMore: () async {
            appState.limit += 20;
            await appState.fetchHistory(appState.limit);
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
              itemBuilder: (context, index) {
                return tile(historyData![index]);
              }),
        ),
      );
    }

    return Container(
      height: height / 1.8,
      child: Center(
        child: CircularProgressIndicator(
          backgroundColor: notifier.getbluecolor,
          valueColor: new AlwaysStoppedAnimation<Color>(
            notifier.getgreencolor,
          ),
          strokeWidth: 3.0,
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
      var splitResult = transaction.memo!.split('>');
      name = "Swapped ${splitResult[0]} to ${splitResult[1]}";
    }

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
      await appState.fetchHistory(appState.limit);
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
}

enum TransactionType {
  Send,
  Receive,
  Swap,
}
