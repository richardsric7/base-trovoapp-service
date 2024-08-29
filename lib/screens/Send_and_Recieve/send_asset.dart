import 'dart:convert';
import 'dart:math';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SendAsset extends StatefulWidget {
  const SendAsset({Key? key}) : super(key: key);

  @override
  State<SendAsset> createState() => _SendAsset();
}

class _SendAsset extends State<SendAsset> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late Asset? asset;
  final formKey = GlobalKey<FormState>();
  String to = ''; // the reciever
  String amount = '';
  bool amountError = false;
  String? memo;
  var deeplinkInfo;
  TextEditingController _utf8TextController = TextEditingController();
  TextEditingController toController = TextEditingController();
  TextEditingController sendingWalletController = TextEditingController();
  final amountController = TextEditingController();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (appState.viewData![SendAssetViewPageConfig.key]?['deepLinkInfo'] ==
            null &&
        appState.viewData!['walletPublicKey'] != null) {
      wallet = appState.userInfo!.getWallet(
        appState.viewData!['walletPublicKey'],
      );
      asset = wallet.claimedAssets!.firstWhere(
        (asset) =>
            asset.assetCode == appState.viewData!['assetCode'] &&
            asset.assetIssuer == appState.viewData!['assetIssuer'],
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    if (appState.viewData![SendAssetViewPageConfig.key] != null &&
        appState.viewData![SendAssetViewPageConfig.key]['deepLinkInfo'] !=
            null) {
      deeplinkInfo =
          appState.viewData![SendAssetViewPageConfig.key]['deepLinkInfo'];
      toController.text = deeplinkInfo['receiver'];
      amountController.text = deeplinkInfo['amount'];
      _utf8TextController.text = deeplinkInfo['memo'];
      wallet = appState.userInfo!.getWallet(deeplinkInfo['sendingWallet']);
      asset = wallet.claimedAssets!.firstWhere(
        (asset) =>
            asset.assetCode == deeplinkInfo['assetCode'] &&
            asset.assetIssuer == deeplinkInfo['assetIssuer'],
      );
      appState.viewData![SendAssetViewPageConfig.key]['deepLinkInfo'] = null;
    }

    sendingWalletController.text = wallet.alias!;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          '',
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '${"send".tr()} ${getAssetCode(asset!.assetCode)}',
                      style: TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    GestureDetector(
                      onTap: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: QrScannerPageConfig);
                      },
                      child: SvgPicture.asset(
                        "assets/images/scan.svg",
                        color: notifier.getbluewhitecolor,
                        height: height / 40,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              formFields(),
              SizedBox(
                height: height / 10,
              ),
              Button(
                "proceed".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  handleSubmit();
                },
              ),
              SizedBox(
                height: height / 20,
              ),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Widget availableBalance() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Flexible(
          child: Text(
            amount.isNotEmpty
                ? "≈ ${formatNumber(double.parse(amount))} ${getAssetCode(asset!.assetCode)}"
                : "≈ 0.0000 ${getAssetCode(asset!.assetCode)}",
            textScaleFactor: 1.0,
            style: TextStyle(
                color: notifier.getdarkgrey,
                fontWeight: FontWeight.w400,
                fontSize: 12.0.sp),
          ),
        ),
        Flexible(
            child: Visibility(
          visible: true,
          replacement: Container(),
          child: Text(
            "${formatNumber(asset!.amount!)} ${getAssetCode(asset!.assetCode)}",
            textScaleFactor: 1.0,
            textAlign: TextAlign.right,
            style: TextStyle(color: notifier.getdarkgrey, fontSize: 12.0.sp),
          ),
        )),
      ],
    );
  }

  Widget formFields() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
              width: 300.sp,
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              ),
              child: Form(
                key: formKey,
                child: Column(
                  children: [
                    if (sendingWalletController.text.isNotEmpty) ...[
                      SizedBox(
                        height: height / 50,
                      ),
                      CustomTextFormField.textField(
                        "sendingwallet".tr(),
                        notifier.getbluecolor,
                        Icons.wallet,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        notifier.getgrey,
                        75.sp,
                        300.sp,
                        controller: sendingWalletController,
                        readOnly: true,
                        onSaved: (value) =>
                            to = value.trim().replaceAll(' ', ''),
                      ),
                    ],
                    SizedBox(
                      height: height / 50,
                    ),
                    GestureDetector(
                      child: CustomTextFormField.textField(
                        "to".tr(),
                        notifier.getbluecolor,
                        Icons.send,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        notifier.getgrey,
                        75.sp,
                        300.sp,
                        controller: toController,
                        readOnly: deeplinkInfo != null,
                        validator: validateTo,
                        onSaved: (value) =>
                            to = value.trim().replaceAll(' ', ''),
                      ),
                    ),
                    SizedBox(height: height / 50),
                    CustomTextFormField.textField(
                      "amount".tr(),
                      notifier.getbluecolor,
                      Icons.currency_exchange,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      75.sp,
                      300.sp,
                      onChanged: (value) {
                        setState(() {
                          amount = trim(value.toString(), '.');
                        });
                      },
                      controller: amountController,
                      readOnly: deeplinkInfo != null &&
                          deeplinkInfo['amount'].toString().isNotEmpty,
                      autoFormatNumber: true,
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
                      validator: validateAmount,
                      onSaved: (value) =>
                          amount = value.trim().replaceAll(' ', ''),
                      inputFormatters: [
                        FilteringTextInputFormatter.allow(RegExp(r'[0-9 \.]'))
                      ],
                    ),
                    if (!appState.hideBalances) ...[availableBalance()],
                    SizedBox(height: height / 50),
                    CustomTextFormField.textField(
                      "memo".tr(),
                      notifier.getbluecolor,
                      Icons.edit_note,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      75.sp,
                      300.sp,
                      onSaved: (value) => memo = value,
                      maxLength: 28,
                      controller: _utf8TextController,
                      readOnly: deeplinkInfo != null &&
                          deeplinkInfo['memo'].toString().isNotEmpty,
                      buildCounter: (context,
                          {currentLength, isFocused, maxLength}) {
                        int utf8Length =
                            utf8.encode(_utf8TextController.text).length;
                        return Container(
                          child: Text(
                            '$utf8Length/$maxLength',
                            style: TextStyle(color: notifier.getdarkgrey),
                          ),
                        );
                      },
                      inputFormatters: [
                        _Utf8LengthLimitingTextInputFormatter(28),
                      ],
                    ),
                    SizedBox(height: height / 20),
                  ],
                ),
              )),
        ),
      ],
    );
  }

  String? validateTo(String? value) {
    // reciever cannot be empty
    if (value!.isEmpty) return "enterreceiverusername".tr();

    if (value.length < 3) return "invalidusernameorpublickey".tr();

    return null;
  }

  String? validateAmount(String? value) {
    if (value!.isEmpty) {
      setState(() {
        amountError = true;
      });
      "pleaseenteramounttosend".tr();
    }

    if (double.tryParse(value) == null) {
      setState(() {
        amountError = true;
      });
      return "pleaseentervalidamount".tr();
    }

    if (double.tryParse(value)! <= 0) {
      setState(() {
        amountError = true;
      });
      return "valuemustbegreaterthan".tr(args: ['0']);
    }

    if (double.tryParse(value)! > asset!.amount!) {
      setState(() {
        amountError = true;
      });
      return "youdonthavesufficientbalance".tr();
    }

    if ((getAssetCode(asset!.assetCode) == 'XBN') &&
        double.tryParse(value)! > (asset!.amount! - 7)) {
      setState(() {
        amountError = true;
      });
      return "youdonthavesufficientbalance".tr();
    }

    setState(() {
      amountError = false;
    });
    return null;
  }

  void handleSubmit() {
    final form = formKey.currentState;
    if (!form!.validate()) {
      return;
    }

    form.save();

    submit();
  }

  submit() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credential
      Map map = {
        "destination": to,
        "memo": memo,
        "amount": amount.toString(),
        "assetCode": asset!.assetCode == 'XBN' ? '' : asset!.assetCode,
        "assetIssuer": asset!.assetIssuer,
      };
      String requestBody = jsonEncode(map);
      print('requestBody =======> $requestBody');
      Map responseData = await makePostRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/payment'
            : '/v1/users/payment',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      hideLoader(context);

      // print('responseData $responseData');

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  postProcessData(messageShown, messageLength, data) {
    // we would like to display all messages returned from the initial
    // request to server using a popup. In order to achieve that we
    // employ the use of a little recursion here. Please recursive
    // functions can turn into a nightmare fast so be carefull here.
    if (messageShown <= messageLength - 1) {
      showResponseMessage(
          context,
          data['messages'][messageShown],
          () => {
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    // go to the definition of appState.viewData
    // to learn more about viewData
    // appState.viewData![ConfirmTransactionViewPageConfig.key] = data;
    // appState.viewData![ConfirmTransactionViewPageConfig.key]["walletInfo"] = {
    //   'alias': activeWallet['alias'],
    //   'publicKey': activeWallet['publicKey'],
    // };
    // appState.viewData![ConfirmTransactionViewPageConfig.key]["usdPrice"] =
    //     viewData['usdPrice'];
    // appState.viewData![ConfirmTransactionViewPageConfig.key]["isSharedWallet"] =
    //     isSharedWallet;
    // if (isSharedWallet) {
    //   appState.viewData![ConfirmTransactionViewPageConfig.key]["rel"] =
    //       'dashboard';
    // }

    appState.viewData = {
      'walletPublicKey': wallet.publicKey,
      'assetCode': asset!.assetCode,
      'assetIssuer': asset!.assetIssuer,
      'rel': 'dashboard',
      'transactionData': data,
    };
    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: ConfirmTransactionViewPageConfig,
    );
  }
}

class _Utf8LengthLimitingTextInputFormatter extends TextInputFormatter {
  _Utf8LengthLimitingTextInputFormatter(this.maxLength)
      : assert(maxLength == -1 || maxLength > 0);

  final int maxLength;

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    if (maxLength > 0 && bytesLength(newValue.text) > maxLength) {
      // If already at the maximum and tried to enter even more, keep the old value.
      if (bytesLength(oldValue.text) == maxLength) {
        return oldValue;
      }
      return truncate(newValue, maxLength);
    }
    return newValue;
  }

  static TextEditingValue truncate(TextEditingValue value, int maxLength) {
    var newValue = '';
    if (bytesLength(value.text) > maxLength) {
      var length = 0;

      value.text.characters.takeWhile((char) {
        var nbBytes = bytesLength(char);
        if (length + nbBytes <= maxLength) {
          newValue += char;
          length += nbBytes;
          return true;
        }
        return false;
      });
    }
    return TextEditingValue(
      text: newValue,
      selection: value.selection.copyWith(
        baseOffset: min(value.selection.start, newValue.length),
        extentOffset: min(value.selection.end, newValue.length),
      ),
      composing: TextRange.empty,
    );
  }

  static int bytesLength(String value) {
    return utf8.encode(value).length;
  }
}
