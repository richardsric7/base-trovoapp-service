import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/get_utils.dart' hide Trans;
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
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
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
    } else if (viewData.transactionDirection == TransactionDirection.Send &&
        viewData.memo!.contains('>')) {
      var splitResult = viewData.memo!.split('>');
      memo = "${splitResult[0]} to ${splitResult[1]} swap fee";
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
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 30),
              Text(
                '${viewData.transactionType!.capitalizeFirst!} ${"details".tr()}',
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
                            SizedBox(
                              height: height / 90,
                            ),
                            Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
                              child: Text(
                                viewData.transactionDirection! ==
                                        TransactionDirection.Send
                                    ? "sentfrom".tr()
                                    : "receivedon".tr(),
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 16.sp,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                            SizedBox(
                              width: width / 1.3,
                              child: Column(
                                children: [
                                  Row(
                                    children: [
                                      Expanded(
                                        flex: 3,
                                        child: Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            viewData.transactionDirection! ==
                                                    TransactionDirection.Send
                                                ? '${extractUsername(viewData.from!)}'
                                                : '${extractUsername(viewData.to!)}',
                                            style: TextStyle(
                                              fontWeight: FontWeight.w500,
                                              color: notifier.getbluewhitecolor,
                                              fontSize: 18.sp,
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
                                            showSnackBar(
                                                "tousername".tr(), context),
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
                                  Row(
                                    children: [
                                      Expanded(
                                        flex: 3,
                                        child: Padding(
                                          padding: const EdgeInsets.symmetric(
                                              horizontal: 20.0),
                                          child: Text(
                                            viewData.transactionDirection! ==
                                                    TransactionDirection.Send
                                                ? truncate(
                                                        viewData.fromPublicKey!,
                                                        length: 5) +
                                                    viewData.fromPublicKey!
                                                        .substring(viewData
                                                                .fromPublicKey!
                                                                .length -
                                                            5)
                                                : truncate(
                                                        viewData.toPublicKey!,
                                                        length: 5) +
                                                    viewData.toPublicKey!
                                                        .substring(viewData
                                                                .toPublicKey!
                                                                .length -
                                                            5),
                                            style: TextStyle(
                                              fontWeight: FontWeight.w500,
                                              color: notifier.getbluewhitecolor,
                                              fontSize: 13.sp,
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
                                            showSnackBar(
                                                "topublickey2".tr(), context),
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
                              padding:
                                  const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                              child: Text(
                                viewData.transactionDirection! ==
                                        TransactionDirection.Send
                                    ? "to".tr()
                                    : "from".tr(),
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
                            SizedBox(
                              width: width / 1.2,
                              child: Row(
                                children: [
                                  SizedBox(
                                    width: width / 20,
                                  ),
                                  Expanded(
                                    flex: 5,
                                    child: Text(
                                      name.toString().isEmpty
                                          ? truncate(publicKey!, length: 5) +
                                              publicKey!.substring(
                                                  publicKey!.length - 5)
                                          : name!,
                                      style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 18.sp,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ),
                                  Expanded(
                                    flex: 1,
                                    child: IconButton(
                                      padding: EdgeInsets.zero,
                                      onPressed: () {
                                        Clipboard.setData(
                                          ClipboardData(
                                            text: name.toString().isEmpty
                                                ? publicKey
                                                : name!,
                                          ),
                                        );

                                        showSnackBar(
                                            name.toString().isEmpty
                                                ? "address".tr()
                                                : "username".tr(),
                                            context);
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
                            if (name.toString().isNotEmpty) ...[
                              SizedBox(
                                width: width / 1.2,
                                child: Row(
                                  children: [
                                    SizedBox(
                                      width: width / 20,
                                    ),
                                    Expanded(
                                      flex: 5,
                                      child: Text(
                                        truncate(publicKey!, length: 5) +
                                            publicKey!.substring(
                                                publicKey!.length - 5),
                                        style: TextStyle(
                                          fontWeight: FontWeight.w500,
                                          color: notifier.getbluewhitecolor,
                                          fontSize: 13.sp,
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
                                              text: publicKey,
                                            ),
                                          ),
                                          showSnackBar("address".tr(), context),
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
                            ],
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
                                "formemo".tr(),
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
                              "blockchainproof".tr(),
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
                                      showSnackBar(
                                          "transactionid".tr(), context),
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
                            height: height / 50,
                          ),
                          Padding(
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
                            child: Text(
                              "date".tr(),
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 16.sp,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Padding(
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
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
                "generatereceipt".tr(),
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
        Container(
          width: width / 1.7,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [],
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
    if (data.isNotEmpty) {
      const start = '[';
      const end = ']';
      final startIndex = data.indexOf(start);
      final endIndex = data.indexOf(end);
      return data.substring(startIndex + start.length, endIndex);
    }

    return '';
  }
}
