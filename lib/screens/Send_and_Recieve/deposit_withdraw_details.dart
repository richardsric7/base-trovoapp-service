import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/deposit_transaction_model.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/models/withdrawal_transaction_model.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DepositWithdrawDetails extends StatefulWidget {
  const DepositWithdrawDetails({Key? key}) : super(key: key);

  @override
  State<DepositWithdrawDetails> createState() => _DepositWithdrawDetails();
}

class _DepositWithdrawDetails extends State<DepositWithdrawDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late DepositTransactionModel depositInfo;
  late WithdrawalTransactionModel withdrawalInfo;
  late TransactionDirection transactionDirection;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet =
        appState.userInfo!.getWallet(appState.viewData!['walletPublicKey']);
    transactionDirection =
        appState.viewData!['transactionDirection'] as TransactionDirection;

    if (transactionDirection == TransactionDirection.Deposit) {
      depositInfo = appState.viewData!['transaction'];
    } else {
      withdrawalInfo = appState.viewData!['transaction'];
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

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
              SizedBox(height: height / 50),
              Text(
                '${transactionDirection == TransactionDirection.Deposit ? 'Deposit' : 'Withdrawal'} ${"details".tr()}',
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 22.sp),
              ),
              SizedBox(height: height / 30),
              Text(
                transactionDirection == TransactionDirection.Deposit
                    ? formatAmount(
                        transactionDirection,
                        depositInfo.amount,
                        depositInfo.currency,
                      )
                    : formatAmount(
                        transactionDirection,
                        withdrawalInfo.amountToWithdraw,
                        withdrawalInfo.currency,
                      ),
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgreencolor,
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
                      child:
                          transactionDirection == TransactionDirection.Deposit
                              ? showDepositInfo()
                              : showWithdrawalInfo(),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "done".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  Navigator.of(context).pop(context);
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

  Widget showDepositInfo() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          height: height / 90,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            "depositedon".tr(),
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
          child: Column(
            children: [
              Row(
                children: [
                  Expanded(
                    flex: 9,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        truncatePublicKey(depositInfo.toAddress),
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
                            text: depositInfo.toAddress,
                          ),
                        ),
                        showSnackBar("toaddress".tr(), context),
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
        Padding(
          padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
          child: Text(
            "from".tr(),
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
          child: Column(
            children: [
              Row(
                children: [
                  Expanded(
                    flex: 9,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        truncatePublicKey(depositInfo.fromAddress),
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
                            text: depositInfo.fromAddress,
                          ),
                        ),
                        showSnackBar("fromaddress".tr(), context),
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
        Padding(
          padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
          child: Text(
            "trovowalletinfo".tr(),
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
          child: Column(
            children: [
              Row(
                children: [
                  Expanded(
                    flex: 9,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        wallet.alias!,
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
                            text: depositInfo.fromAddress,
                          ),
                        ),
                        showSnackBar("walletalias".tr(), context),
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
                    flex: 15,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        truncatePublicKey(depositInfo.trovoWalletPublicKey),
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
                            text: depositInfo.fromAddress,
                          ),
                        ),
                        showSnackBar("publickey".tr(), context),
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
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
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
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Row(
            children: [
              Expanded(
                flex: 5,
                child: GestureDetector(
                  onTap: () => appState.goToWebView(
                      getExplorerBaseUrl(appState.walletMode) +
                          depositInfo.transactionId),
                  child: Text(
                    depositInfo.transactionId,
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
                        text: depositInfo.transactionId,
                      ),
                    ),
                    showSnackBar("transactionid".tr(), context),
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
        Padding(
          padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
          child: Text(
            "transactionstatus",
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
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            depositInfo.isCompleted ? "completed".tr() : "pending_2".tr(),
            style: TextStyle(
              color: notifier.getbluewhitecolor,
              fontSize: 13.sp,
              fontWeight: FontWeight.w500,
              fontFamily: fontbody,
            ),
          ),
        ),
        Divider(
          height: 5,
        ),
        SizedBox(
          height: height / 50,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
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
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            '${DateFormat('MMMM dd, yyyy hh:mm a').format(depositInfo.createdAt)}',
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
    );
  }

  Widget showWithdrawalInfo() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          height: height / 90,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            "withdrawnfrom".tr(),
            style: TextStyle(
              fontWeight: FontWeight.w500,
              color: notifier.getbluewhitecolor,
              fontSize: 16.sp,
              fontFamily: fontsemibold,
            ),
          ),
        ),
        SizedBox(
          width: width / 1.5,
          child: Column(
            children: [
              Row(
                children: [
                  Expanded(
                    flex: 9,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        withdrawalInfo.walletAlias,
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
                            text: withdrawalInfo.walletAlias,
                          ),
                        ),
                        showSnackBar('', context),
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
                    flex: 9,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        truncatePublicKey(withdrawalInfo.walletPublicKey),
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
                            text: withdrawalInfo.walletPublicKey,
                          ),
                        ),
                        showSnackBar("publickey".tr(), context),
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
        Padding(
          padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
          child: Text(
            "networkinfo".tr(),
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
          child: Column(
            children: [
              keyValuePair(
                '${"network".tr()}:',
                truncatePublicKey(withdrawalInfo.withdrawalNetwork),
              ),
              keyValuePair(
                '${"address".tr()}:',
                truncatePublicKey(withdrawalInfo.withdrawalAddress),
                copy: true,
                copyText: withdrawalInfo.withdrawalAddress,
              ),
            ],
          ),
        ),
        Divider(
          height: 5,
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
          child: Text(
            "fees".tr(),
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
          child: Column(
            children: [
              keyValuePair(
                '${"networkfee".tr()}:',
                '${withdrawalInfo.withdrawalNetworkFee} ${withdrawalInfo.currency}',
              ),
              SizedBox(
                height: height / 50,
              ),
              keyValuePair(
                '${"servicefee".tr()}:',
                '${withdrawalInfo.withdrawalServiceFee} XBN',
              ),
            ],
          ),
        ),
        Divider(
          height: 5,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
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
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Row(
            children: [
              Expanded(
                flex: 5,
                child: GestureDetector(
                  onTap: () => appState.goToWebView(
                      getExplorerBaseUrl(appState.walletMode) +
                          withdrawalInfo.transactionId),
                  child: Text(
                    withdrawalInfo.transactionId,
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
                        text: withdrawalInfo.transactionId,
                      ),
                    ),
                    showSnackBar("transactionid".tr(), context),
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
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
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
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            '${DateFormat('MMMM dd, yyyy hh:mm a').format(withdrawalInfo.createdAt)}',
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
    );
  }

  String formatAmount(TransactionDirection transactionType, amount, assetCode) {
    var am = formatNumber(double.parse(amount.toString()));
    return transactionType == TransactionDirection.Withdraw
        ? '- $am $assetCode'
        : '+ $am $assetCode';
  }

  Widget keyValuePair(String key, String value,
      {bool copy = false, String copyText = ''}) {
    return Row(
      children: [
        Expanded(
          flex: 5,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Text(
              key,
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
          flex: 5,
          child: Row(
            children: [
              Text(
                value,
                style: TextStyle(
                  fontWeight: FontWeight.w500,
                  color: notifier.getbluewhitecolor,
                  fontSize: 15.sp,
                  fontFamily: fontbody,
                ),
              ),
              if (copy) ...[
                Expanded(
                  flex: 1,
                  child: IconButton(
                    onPressed: () => {
                      Clipboard.setData(
                        ClipboardData(
                          text: copyText,
                        ),
                      ),
                      showSnackBar('', context),
                    },
                    icon: Icon(
                      Icons.copy,
                      size: 20,
                    ),
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ],
            ],
          ),
        ),
      ],
    );
  }
}
