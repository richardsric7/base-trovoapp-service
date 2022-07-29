import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:share/share.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
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

class PaymentDetails extends StatefulWidget {
  const PaymentDetails({Key? key}) : super(key: key);

  @override
  State<PaymentDetails> createState() => _PaymentDetails();
}

class _PaymentDetails extends State<PaymentDetails>
    with TickerProviderStateMixin {
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
    viewData = appState.viewData![PaymentDetailsViewPageConfig.key];
    transactionType = viewData.fromPublicKey == activeWallet!.publicKey
        ? TransactionType.Send
        : TransactionType.Receive;

    name =
        transactionType == TransactionType.Send ? viewData.to : viewData.from;

    publicKey = transactionType == TransactionType.Send
        ? viewData.toPublicKey
        : viewData.fromPublicKey;

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
              SizedBox(height: height / 15),
              Text(
                LanguageEn.transactionDetails,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluecolor,
                    fontFamily: fontsemibold,
                    fontSize: 22.sp),
              ),
              SizedBox(height: height / 30),
              Text(
                formatAmount(
                    transactionType, viewData.amount, viewData.assetCode),
                style: TextStyle(
                    color: transactionType == TransactionType.Send
                        ? Colors.red
                        : notifier.getgreencolor,
                    fontFamily: fontsemibold,
                    fontSize: 20.sp),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.getaddsubwalletgrey,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                        child: Text(
                          transactionType == TransactionType.Send
                              ? LanguageEn.sentto
                              : LanguageEn.receivedfrom,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      showUserInfo(),
                      SizedBox(
                        height: height / 50,
                      ),
                      Divider(
                        height: 5,
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                      if (viewData.memo!.isNotEmpty) ...[
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20.0, vertical: 10),
                          child: Text(
                            LanguageEn.formemo,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluecolor,
                              fontSize: 16.sp,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        SizedBox(
                          height: 5,
                        ),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Text(
                            viewData.memo!,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluecolor,
                              fontSize: 15.sp,
                              fontFamily: fontbody,
                            ),
                          ),
                        ),
                        SizedBox(
                          height: height / 50,
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
                            horizontal: 20.0, vertical: 10),
                        child: Text(
                          LanguageEn.blockchainproof,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
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
                                    decoration: TextDecoration.underline,
                                    color: notifier.getbluecolor,
                                    fontSize: 12.sp,
                                    fontWeight: FontWeight.w500,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              ),
                            ),
                            Expanded(
                              flex: 1,
                              child: IconButton(
                                onPressed: () => {
                                  Clipboard.setData(
                                    ClipboardData(
                                      text: viewData.transactionId!,
                                    ),
                                  ),
                                  showSnackBar('Transaction ID', context),
                                },
                                icon: Icon(Icons.copy),
                                color: notifier.getbluecolor,
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
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                LanguageEn.share,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  share();
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
          width: width / 1.7,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    flex: 3,
                    child: Text(
                      name.toString().isEmpty
                          ? truncate(publicKey!, length: 5) +
                              publicKey!.substring(publicKey!.length - 5)
                          : extractUsername(name!),
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluecolor,
                        fontSize: 19.sp,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                  Expanded(
                    flex: 1,
                    child: IconButton(
                      padding: EdgeInsets.zero,
                      onPressed: () => {
                        Clipboard.setData(
                          ClipboardData(
                            text: name.toString().isEmpty
                                ? truncate(publicKey!, length: 5) +
                                    publicKey!.substring(publicKey!.length - 5)
                                : extractUsername(name!),
                          ),
                        ),
                        showSnackBar('Address', context),
                      },
                      icon: Icon(Icons.copy),
                      color: notifier.getbluecolor,
                    ),
                  ),
                ],
              ),
              Text(
                '$date',
                style: TextStyle(
                  color: notifier.getbluecolor,
                  fontSize: 12.sp,
                  fontWeight: FontWeight.w500,
                  fontFamily: fontbody,
                ),
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
    const start = '[';
    const end = ']';
    final startIndex = data.indexOf(start);
    final endIndex = data.indexOf(end);
    return data.substring(startIndex + start.length, endIndex);
  }

  void share() {
    String? shareString;
    switch (transactionType) {
      // case TransactionType.swap:
      //   shareString =
      //       'Swapped from ${transaction.asset.name} to ${transaction.destinationAsset?.name} \nAmount: ${_getSwapValue(transaction, truncateLength: 7)} \nTransaction Id: ${transaction.transactionId.toLowerCase()} \nTime: ${_getTimestampString(transaction.timestamp)}';
      //   break;
      case TransactionType.Send:
        shareString =
            'Sent $amount $assetCode \n\nTo: ${name.toString().isEmpty ? publicKey! : extractUsername(name!)} \n\nFor: ${viewData.memo} \n\nTransaction Id: ${viewData.transactionId!.toLowerCase()} \n\nTime: ${date} \n\nBlockchain Proof: ${bantuBlockchainExplorerBaseUrl + viewData.transactionId!}';
        break;
      default:
        shareString =
            'Recieved $amount $assetCode \n\nFrom: ${name.toString().isEmpty ? publicKey! : extractUsername(name!)} \n\nFor: ${viewData.memo} \n\nTransaction Id: ${viewData.transactionId!.toLowerCase()} \n\nTime: ${date} \n\nBlockchain Proof: ${bantuBlockchainExplorerBaseUrl + viewData.transactionId!}';
    }

    Share.share(shareString);
  }
}
