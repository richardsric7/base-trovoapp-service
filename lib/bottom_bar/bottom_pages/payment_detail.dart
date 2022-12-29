import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/get_utils.dart';
import 'package:intl/intl.dart';
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
import 'package:trovo_wallet/router/PageActions.dart';
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
  String? name;
  String? publicKey;
  double? amount;
  String? assetCode;
  String? date;
  String memo = '';

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
    print('viewData: $viewData');
    name = '${extractUsername(viewData.from!)}';
    publicKey = viewData.fromPublicKey;
    memo = viewData.memo!;

    // if record.from is same as the current active wallet public key
    // then it was a send transaction
    if (viewData.transactionDirection == TransactionDirection.Send) {
      name = '${extractUsername(viewData.to!)}';
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
              SizedBox(height: height / 30),
              Text(
                '${viewData.transactionType!.capitalizeFirst!} ${LanguageEn.details}',
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 22.sp),
              ),
              SizedBox(height: height / 30),
              Text(
                formatAmount(viewData.transactionDirection!, viewData.amount,
                    viewData.assetCode),
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: viewData.transactionDirection! ==
                            TransactionDirection.Send
                        ? Colors.red
                        : notifier.getgreencolor,
                    fontFamily: fontsemibold,
                    fontSize: 20.sp),
              ),
              SizedBox(
                height: height / 50,
              ),
              Stack(
                alignment: AlignmentDirectional.center,
                children: [
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
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          if (TransactionDirection.Swap !=
                              viewData.transactionDirection!) ...[
                            Padding(
                              padding:
                                  const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                              child: Text(
                                viewData.transactionDirection! ==
                                        TransactionDirection.Send
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
                          if (TransactionDirection.Swap !=
                              viewData.transactionDirection!) ...[
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
                              width: width / 1.7,
                              child: Row(
                                children: [
                                  Expanded(
                                    flex: 3,
                                    child: Padding(
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0),
                                      child: Text(
                                        truncate(viewData.fromPublicKey!,
                                                length: 5) +
                                            viewData.fromPublicKey!.substring(
                                                viewData.fromPublicKey!.length -
                                                    5),
                                        style: TextStyle(
                                          fontWeight: FontWeight.w500,
                                          color: notifier.getbluewhitecolor,
                                          fontSize: 15.sp,
                                          fontFamily: fontbody,
                                        ),
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
                                            text: viewData.fromPublicKey!,
                                          ),
                                        ),
                                        showSnackBar(
                                            'From public key', context),
                                      },
                                      icon: Icon(
                                        Icons.copy,
                                        size: 20,
                                      ),
                                      color: notifier.getbluewhitecolor,
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
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
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
                              width: width / 1.7,
                              child: Row(
                                children: [
                                  Expanded(
                                    flex: 3,
                                    child: Padding(
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0),
                                      child: Text(
                                        truncate(viewData.toPublicKey!,
                                                length: 5) +
                                            viewData.toPublicKey!.substring(
                                                viewData.toPublicKey!.length -
                                                    5),
                                        style: TextStyle(
                                          fontWeight: FontWeight.w500,
                                          color: notifier.getbluewhitecolor,
                                          fontSize: 15.sp,
                                          fontFamily: fontbody,
                                        ),
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
                                            text: viewData.toPublicKey!,
                                          ),
                                        ),
                                        showSnackBar('To public key', context),
                                      },
                                      icon: Icon(
                                        Icons.copy,
                                        size: 20,
                                      ),
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            Divider(
                              height: 5,
                            ),
                          ],
                          if (viewData.memo!.isNotEmpty) ...[
                            SizedBox(
                              height: height / 90,
                            ),
                            Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
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
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
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
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
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
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
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
                                        color: notifier.getbluewhitecolor,
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
                                    icon: Icon(
                                      Icons.copy,
                                      size: 20,
                                    ),
                                    color: notifier.getbluewhitecolor,
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
                  Image.asset(
                    'assets/images/trovo_white.png',
                    height: TransactionDirection.Swap !=
                            viewData.transactionDirection!
                        ? height / 4.5
                        : height / 6.5,
                    color: notifier.isDark
                        ? notifier.getdarkgrey
                        : notifier.getsplashgrey,
                  ),
                ],
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                'Generate receipt',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: ShareReceiptViewPageConfig,
                  );

                  appState.viewData![ShareReceiptViewPageConfig.key] = viewData;
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
                  Expanded(
                    flex: 1,
                    child: IconButton(
                      padding: EdgeInsets.zero,
                      onPressed: () => {
                        Clipboard.setData(
                          ClipboardData(
                            text: name.toString().isEmpty ? publicKey : name!,
                          ),
                        ),
                        showSnackBar(
                            name.toString().isEmpty ? 'Address' : 'Username',
                            context),
                      },
                      icon: Icon(
                        Icons.copy,
                        size: 20,
                      ),
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
              Text(
                '$date',
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
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

  String formatAmount(TransactionDirection transactionType, amount, assetCode) {
    var am = formatNumber(double.parse(amount.toString()));
    return transactionType == TransactionDirection.Send
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
}
