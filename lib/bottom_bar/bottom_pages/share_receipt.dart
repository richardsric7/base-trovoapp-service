import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:intl/intl.dart';
import 'package:share/share.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/transaction.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
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
    name = viewData.from;
    publicKey = viewData.fromPublicKey;
    memo = viewData.memo!;

    // if record.from is same as the current active wallet public key
    // then it was a send transaction
    if (viewData.transactionDirection == TransactionDirection.Send) {
      name = viewData.to;
      publicKey = viewData.toPublicKey;
    }

    if (viewData.transactionType!.contains('SWAP') &&
        viewData.memo!.contains('>')) {
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
                  color: notifier.isDark ? null : wihitecolor,
                  child: Column(
                    children: [
                      SizedBox(height: height / 50),
                      Image.asset(
                        'assets/images/trovo-horizontal-logo.png',
                        height: height / 16.5,
                        color: notifier.isDark ? Colors.white : null,
                      ),
                      SizedBox(height: height / 30),
                      Text(
                        '${viewData.transactionType!.capitalizeFirst!} ${LanguageEn.details}',
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
                                  if (TransactionDirection.Swap !=
                                      viewData.transactionDirection!) ...[
                                    SizedBox(
                                      height: height / 90,
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0),
                                      child: Text(
                                        viewData.transactionDirection! ==
                                                TransactionDirection.Send
                                            ? LanguageEn.sentfrom
                                            : 'Received on',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w500,
                                          color: notifier.getbluewhitecolor,
                                          fontSize: 16.sp,
                                          fontFamily: fontsemibold,
                                        ),
                                      ),
                                    ),
                                    SizedBox(
                                      width: width / 1.7,
                                      child: Column(
                                        children: [
                                          Row(
                                            children: [
                                              Expanded(
                                                flex: 3,
                                                child: Padding(
                                                  padding: const EdgeInsets
                                                          .symmetric(
                                                      horizontal: 20.0),
                                                  child: Text(
                                                    viewData.transactionDirection! ==
                                                            TransactionDirection
                                                                .Send
                                                        ? '${viewData.from!}'
                                                        : '${viewData.to!}',
                                                    style: TextStyle(
                                                      fontWeight:
                                                          FontWeight.w500,
                                                      color: notifier
                                                          .getbluewhitecolor,
                                                      fontSize: 18.sp,
                                                      fontFamily: fontbody,
                                                    ),
                                                  ),
                                                ),
                                              ),
                                            ],
                                          ),
                                          Row(
                                            children: [
                                              Expanded(
                                                flex: 3,
                                                child: Padding(
                                                  padding: const EdgeInsets
                                                          .symmetric(
                                                      horizontal: 20.0,
                                                      vertical: 5),
                                                  child: Text(
                                                    viewData.transactionDirection! ==
                                                            TransactionDirection
                                                                .Send
                                                        ? truncate(
                                                                viewData
                                                                    .fromPublicKey!,
                                                                length: 5) +
                                                            viewData
                                                                .fromPublicKey!
                                                                .substring(viewData
                                                                        .fromPublicKey!
                                                                        .length -
                                                                    5)
                                                        : truncate(
                                                                viewData
                                                                    .toPublicKey!,
                                                                length: 5) +
                                                            viewData
                                                                .toPublicKey!
                                                                .substring(viewData
                                                                        .toPublicKey!
                                                                        .length -
                                                                    5),
                                                    style: TextStyle(
                                                      fontWeight:
                                                          FontWeight.w500,
                                                      color: notifier
                                                          .getbluewhitecolor,
                                                      fontSize: 13.sp,
                                                      fontFamily: fontbody,
                                                    ),
                                                  ),
                                                ),
                                              ),
                                            ],
                                          ),
                                        ],
                                      ),
                                    ),
                                    Divider(
                                      height: 5,
                                    ),
                                  ],
                                  if (TransactionDirection.Swap !=
                                      viewData.transactionDirection!) ...[
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          20.0, 10, 0, 0),
                                      child: Text(
                                        viewData.transactionDirection! ==
                                                TransactionDirection.Send
                                            ? LanguageEn.to
                                            : 'From',
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
                                        horizontal: 20.0),
                                    child: Text(
                                      LanguageEn.amount,
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
                                      formatAmount(
                                          viewData.transactionDirection!,
                                          viewData.amount,
                                          viewData.assetCode),
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
                                  Divider(
                                    height: 5,
                                  ),
                                  SizedBox(
                                    height: height / 50,
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
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
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0, vertical: 5),
                                    child: Text(
                                      '$date',
                                      style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 13.sp,
                                        fontWeight: FontWeight.w500,
                                        fontFamily: fontbody,
                                      ),
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
                  share(
                      'Blockchain proof\n$bantuBlockchainExplorerBaseUrl${viewData.transactionId!}',
                      shareArea);
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
                  sharePDF(
                      'Blockchain proof\n${bantuBlockchainExplorerBaseUrl}${viewData.transactionId!}',
                      shareArea);
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
              Row(
                children: [
                  Expanded(
                    flex: 5,
                    child: Text(
                      name.toString().isNotEmpty
                          ? truncate(publicKey!, length: 5) +
                              publicKey!.substring(publicKey!.length - 5)
                          : '',
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontSize: 13.sp,
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

  String formatAmount(TransactionDirection transactionType, amount, assetCode) {
    var am = formatNumber(double.parse(amount.toString()));
    return '$am $assetCode';
  }

  // String extractUsername(String data) {
  //   print('data $data');
  //   if (data.isNotEmpty) {
  //     const start = '[';
  //     const end = ']';
  //     final startIndex = data.indexOf(start);
  //     final endIndex = data.indexOf(end);
  //     print('data $data');
  //     return data.substring(startIndex + start.length, endIndex);
  //   }

  //   return '';
  // }

  void shareText() {
    String? shareString;
    switch (viewData.transactionDirection) {
      case TransactionDirection.Send:
        shareString =
            'Payment Details\n____________________\n\nSent $amount $assetCode \n\nTo: ${name.toString().isEmpty ? publicKey! : name} \n\nFor: ${viewData.memo} \n\nTransaction Id: ${viewData.transactionId!.toLowerCase()} \n\nTime: ${date} \n\nBlockchain Proof: ${bantuBlockchainExplorerBaseUrl + viewData.transactionId!}';
        break;
      default:
        shareString =
            'Payment Details\n____________________\n\nRecieved $amount $assetCode \n\nFrom: ${name.toString().isEmpty ? publicKey! : name} \n\nFor: ${viewData.memo} \n\nTransaction Id: ${viewData.transactionId!.toLowerCase()} \n\nTime: ${date} \n\nBlockchain Proof: ${bantuBlockchainExplorerBaseUrl + viewData.transactionId!}';
    }

    Share.share(shareString);
  }
}
