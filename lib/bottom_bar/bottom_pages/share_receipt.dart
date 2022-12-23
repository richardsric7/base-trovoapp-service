import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:share/share.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Transaction.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'payment_history.dart';

class ShareReceipt extends StatefulWidget {
  const ShareReceipt({Key? key}) : super(key: key);

  @override
  State<ShareReceipt> createState() => _ShareReceipt();
}

class _ShareReceipt extends State<ShareReceipt> with TickerProviderStateMixin {
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
  late TransactionInfo viewData;
  late TransactionType transactionType;
  String? name;
  String? publicKey;
  double? amount;
  String? assetCode;
  String? date;
  String memo = '';
  GlobalKey shareArea = GlobalKey();

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
    viewData = appState.viewData![ShareReceiptViewPageConfig.key];
    print('viewData: $viewData');
    transactionType = TransactionType.Receive;
    name = viewData.from.toString().contains('[')
        ? '${extractUsername(viewData.from!)}'
        : viewData.from;
    publicKey = viewData.fromPublicKey;
    memo = viewData.memo!;

    // if record.from is same as the current active wallet public key
    // then it was a send transaction
    if (viewData.fromPublicKey == activeWallet!.publicKey) {
      transactionType = TransactionType.Send;
      name = viewData.to.toString().contains('[')
          ? '${extractUsername(viewData.to!)}'
          : viewData.to;
      publicKey = viewData.toPublicKey;
    }

    if (viewData.transactionType!.contains('SWAP') &&
        viewData.memo!.contains('>')) {
      transactionType = TransactionType.Swap;
      var splitResult = viewData.memo!.split('>');
      memo = "Swapped ${splitResult[0]} to ${splitResult[1]}";
    }

    amount = viewData.amount;

    assetCode = viewData.assetCode;
    date =
        DateFormat('MMMM dd, yyyy hh:mm a').format(viewData.transactionDate!);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              RepaintBoundary(
                key: shareArea,
                child: Container(
                  color: Colors.white,
                  child: Column(
                    children: [
                      SizedBox(height: height / 50),
                      Image.asset(
                        'assets/images/trovo-horizontal-logo.png',
                        height: height / 16.5,
                      ),
                      SizedBox(height: height / 30),
                      Text(
                        'Transaction Receipt',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 22.sp),
                      ),
                      SizedBox(height: 3),
                      Container(
                        width: width / 1.5,
                        child: Text(
                          'Generated from Trovo-Wallet on ${DateFormat('MMMM dd, yyyy hh:mm a').format(DateTime.now())}',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                              fontSize: 12.sp),
                        ),
                      ),
                      SizedBox(height: height / 50),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                        child: Container(
                          decoration: BoxDecoration(
                            borderRadius:
                                const BorderRadius.all(Radius.circular(15.0)),
                            color: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                          ),
                          child: Stack(
                            alignment: AlignmentDirectional.center,
                            children: [
                              Image.asset(
                                'assets/images/trovo_white.png',
                                height: height / 4.5,
                                color: notifier.isDark
                                    ? notifier.getdarkgrey
                                    : notifier.getsplashgrey,
                              ),
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  if (TransactionType.Swap !=
                                      transactionType) ...[
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          20.0, 15, 0, 0),
                                      child: Text(
                                        transactionType == TransactionType.Send
                                            ? LanguageEn.sentto
                                            : LanguageEn.receivedfrom,
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
                                    showUserInfo(),
                                    Divider(
                                      height: 5,
                                    ),
                                  ],
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 20.0,
                                    ),
                                    child: Text(
                                      'Amount',
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16.sp,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    width: width / 1.2,
                                    child: Row(
                                      children: [
                                        Expanded(
                                          flex: 7,
                                          child: Padding(
                                            padding: const EdgeInsets.fromLTRB(
                                                20.0, 0, 0, 0),
                                            child: Text(
                                              formatAmount(
                                                  transactionType,
                                                  viewData.amount,
                                                  viewData.assetCode),
                                              style: TextStyle(
                                                  color: transactionType ==
                                                          TransactionType.Send
                                                      ? Colors.red
                                                      : notifier.getgreencolor,
                                                  fontFamily: fontsemibold,
                                                  fontSize: 15.sp),
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                  Divider(
                                    height: 5,
                                  ),
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 20.0,
                                    ),
                                    child: Text(
                                      'Date',
                                      style: TextStyle(
                                        fontWeight: FontWeight.w500,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 16.sp,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    width: width / 1.2,
                                    child: Row(
                                      children: [
                                        Expanded(
                                          flex: 7,
                                          child: Padding(
                                            padding: const EdgeInsets.fromLTRB(
                                                20.0, 0, 0, 0),
                                            child: Text(
                                              '$date',
                                              style: TextStyle(
                                                color:
                                                    notifier.getbluewhitecolor,
                                                fontSize: 15.sp,
                                                fontWeight: FontWeight.w500,
                                                fontFamily: fontbody,
                                              ),
                                            ),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                  Divider(
                                    height: 5,
                                  ),
                                  if (TransactionType.Swap !=
                                      transactionType) ...[
                                    SizedBox(
                                      height: height / 90,
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0,
                                      ),
                                      child: Text(
                                        'From Public Key',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w500,
                                          color: notifier.getbluewhitecolor,
                                          fontSize: 16.sp,
                                          fontFamily: fontsemibold,
                                        ),
                                      ),
                                    ),
                                    SizedBox(
                                      width: width / 1.2,
                                      child: Row(
                                        children: [
                                          Expanded(
                                            flex: 3,
                                            child: Padding(
                                              padding:
                                                  const EdgeInsets.symmetric(
                                                      horizontal: 20.0),
                                              child: Text(
                                                truncate(
                                                        viewData.fromPublicKey!,
                                                        length: 5) +
                                                    viewData.fromPublicKey!
                                                        .substring(viewData
                                                                .fromPublicKey!
                                                                .length -
                                                            5),
                                                style: TextStyle(
                                                  fontWeight: FontWeight.w500,
                                                  color: notifier
                                                      .getbluewhitecolor,
                                                  fontSize: 15.sp,
                                                  fontFamily: fontbody,
                                                ),
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                    Divider(
                                      height: 5,
                                    ),
                                    SizedBox(
                                      height: height / 90,
                                    ),
                                    if (viewData.toPublicKey
                                        .toString()
                                        .isNotEmpty) ...[
                                      Padding(
                                        padding: const EdgeInsets.symmetric(
                                            horizontal: 20.0),
                                        child: Text(
                                          'To Public Key',
                                          style: TextStyle(
                                            fontWeight: FontWeight.w500,
                                            color: notifier.getbluewhitecolor,
                                            fontSize: 16.sp,
                                            fontFamily: fontsemibold,
                                          ),
                                        ),
                                      ),
                                      SizedBox(
                                        width: width / 1.2,
                                        child: Row(
                                          children: [
                                            Expanded(
                                              flex: 3,
                                              child: Padding(
                                                padding:
                                                    const EdgeInsets.symmetric(
                                                        horizontal: 20.0),
                                                child: Text(
                                                  truncate(
                                                          viewData.toPublicKey!,
                                                          length: 5) +
                                                      viewData.toPublicKey!
                                                          .substring(viewData
                                                                  .toPublicKey!
                                                                  .length -
                                                              5),
                                                  style: TextStyle(
                                                    fontWeight: FontWeight.w500,
                                                    color: notifier
                                                        .getbluewhitecolor,
                                                    fontSize: 15.sp,
                                                    fontFamily: fontbody,
                                                  ),
                                                ),
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Divider(
                                        height: 5,
                                      ),
                                    ],
                                  ],
                                  if (viewData.memo!.isNotEmpty) ...[
                                    SizedBox(
                                      height: height / 90,
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0),
                                      child: Text(
                                        LanguageEn.formemo,
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
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0),
                                      child: Text(
                                        memo,
                                        style: TextStyle(
                                          fontWeight: FontWeight.w500,
                                          color: notifier.getbluewhitecolor,
                                          fontSize: 15.sp,
                                          fontFamily: fontbody,
                                        ),
                                      ),
                                    ),
                                    Divider(
                                      height: 5,
                                    ),
                                  ],
                                  SizedBox(
                                    height: height / 90,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
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
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Row(
                                      children: [
                                        Expanded(
                                          flex: 5,
                                          child: GestureDetector(
                                            onTap: () => appState.goToWebView(
                                                bantuBlockchainExplorerBaseUrl +
                                                    viewData.transactionId!),
                                            child: Text(
                                              viewData.transactionId!,
                                              style: TextStyle(
                                                decoration:
                                                    TextDecoration.underline,
                                                color:
                                                    notifier.getbluewhitecolor,
                                                fontSize: 12.sp,
                                                fontWeight: FontWeight.w500,
                                                fontFamily: fontbody,
                                              ),
                                            ),
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
                            ],
                          ),
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
                height: height / 30,
              ),
              Button(
                'Share Image',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  share('', shareArea);
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              Button(
                'Share PDF',
                notifier.getbluecolor70,
                wihitecolor,
                onTap: () {
                  sharePDF('', shareArea);
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              ButtonOutlined(
                'Share Text',
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  shareText();
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
        SizedBox(
          width: width / 20,
        ),
        SizedBox(
          width: width / 1.2,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    flex: 5,
                    child: Text(
                      name.toString().isEmpty
                          ? truncate(publicKey!, length: 5) +
                              publicKey!.substring(publicKey!.length - 5)
                          : name!,
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontSize: 15.sp,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        )
      ],
    );
  }

  String formatAmount(TransactionType transactionType, amount, assetCode) {
    var am = formatNumber(double.parse(amount.toString()));
    return transactionType == TransactionType.Send
        ? '- $am $assetCode'
        : '+ $am $assetCode';
  }

  String extractUsername(String data) {
    print('data $data');
    if (data.isNotEmpty) {
      const start = '[';
      const end = ']';
      final startIndex = data.indexOf(start);
      final endIndex = data.indexOf(end);
      print('data $data');
      return data.substring(startIndex + start.length, endIndex);
    }

    return '';
  }

  void shareText() {
    String? shareString;
    switch (transactionType) {
      // case TransactionType.swap:
      //   shareString =
      //       'Swapped from ${transaction.asset.name} to ${transaction.destinationAsset?.name} \nAmount: ${_getSwapValue(transaction, truncateLength: 7)} \nTransaction Id: ${transaction.transactionId.toLowerCase()} \nTime: ${_getTimestampString(transaction.timestamp)}';
      //   break;
      case TransactionType.Send:
        shareString =
            'Sent $amount $assetCode \n\nTo: ${name.toString().isEmpty ? publicKey! : name} \n\nFor: ${viewData.memo} \n\nTransaction Id: ${viewData.transactionId!.toLowerCase()} \n\nTime: ${date} \n\nBlockchain Proof: ${bantuBlockchainExplorerBaseUrl + viewData.transactionId!}';
        break;
      default:
        shareString =
            'Recieved $amount $assetCode \n\nFrom: ${name.toString().isEmpty ? publicKey! : name} \n\nFor: ${viewData.memo} \n\nTransaction Id: ${viewData.transactionId!.toLowerCase()} \n\nTime: ${date} \n\nBlockchain Proof: ${bantuBlockchainExplorerBaseUrl + viewData.transactionId!}';
    }

    Share.share(shareString);
  }
}
