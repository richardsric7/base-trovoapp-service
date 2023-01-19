import 'dart:convert';
import 'dart:math';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SendAssetSharedWallet extends StatefulWidget {
  const SendAssetSharedWallet({Key? key}) : super(key: key);

  @override
  State<SendAssetSharedWallet> createState() => _SendAssetSharedWallet();
}

class _SendAssetSharedWallet extends State<SendAssetSharedWallet>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  String to = ''; // the reciever
  String amount = '';
  bool amountError = false;
  String? memo;
  var asset;
  var deeplinkInfo;
  TextEditingController _utf8TextController = TextEditingController();
  TextEditingController toController = TextEditingController();
  final amountController = TextEditingController();
  var walletDetails; // details of the current shared wallet
  var acceptedNumbers = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0];

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
    asset = appState.viewData![SendAssetSharedWalletViewPageConfig.key];
    walletDetails = asset['walletInfo'];
    print('viewData: ${walletDetails}');

    if (asset['deepLinkInfo'] != null) {
      deeplinkInfo = asset['deepLinkInfo'];
      toController.text = deeplinkInfo['receiver'];
      amountController.text = deeplinkInfo['amount'];
      amount = deeplinkInfo['amount'];
      _utf8TextController.text = deeplinkInfo['memo'];
      asset['deepLinkInfo'] = null;
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(height / 15),
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
          ),
        ),
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
                      '${LanguageEn.send} ${getAssetCode(asset['assetCode'])}',
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
                height: height / 20,
              ),
              Button(
                LanguageEn.proceed,
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
                ? "≈ ${formatNumber(double.parse(amount))} ${getAssetCode(asset['assetCode'])}"
                : "≈ 0.0000 ${getAssetCode(asset['assetCode'])}",
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
            "${formatNumber(double.parse(asset['amount']))} ${getAssetCode(asset['assetCode'])}",
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
              height: height / 2.5,
              width: 300.sp,
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              ),
              child: Form(
                key: formKey,
                child: Column(
                  children: [
                    SizedBox(
                      height: height / 50,
                    ),
                    CustomTextFormField.textField(
                      LanguageEn.to,
                      notifier.getbluecolor,
                      Icons.send,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      controller: toController,
                      readOnly: deeplinkInfo != null,
                      validator: validateTo,
                      onSaved: (value) => to = value.trim().replaceAll(' ', ''),
                    ),
                    SizedBox(height: height / 50),
                    CustomTextFormField.textField(
                      LanguageEn.amount,
                      notifier.getbluecolor,
                      Icons.currency_exchange,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      // dynamically change the size
                      // of the textbox so it will
                      // consistent when showing an
                      // error message
                      amountError ? 70.sp : 58.sp,
                      300.sp,
                      onChanged: (value) {
                        setState(() {
                          amount = trim(value.toString(), '.');
                        });
                      },
                      controller: amountController,
                      readOnly: deeplinkInfo != null,
                      // keyboardtype:
                      //     TextInputType.numberWithOptions(decimal: true),
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
                      LanguageEn.memo,
                      notifier.getbluecolor,
                      Icons.edit_note,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      onSaved: (value) => memo = value,
                      maxLength: 28,
                      controller: _utf8TextController,
                      readOnly: deeplinkInfo != null,
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
                    SizedBox(height: height / 50),
                  ],
                ),
              )),
        ),
      ],
    );
  }

  String? validateTo(String? value) {
    // reciever cannot be empty
    if (value!.isEmpty) return 'Please enter reciever username or public key';

    if (value.length < 3) return 'Invalid username or public key';

    return null;
  }

  String? validateAmount(String? value) {
    if (value!.isEmpty) {
      setState(() {
        amountError = true;
      });
      'Please enter amount to send';
    }

    if (double.tryParse(value) == null) {
      setState(() {
        amountError = true;
      });
      return 'Please enter a valid amount';
    }

    if (double.tryParse(value)! <= 0) {
      setState(() {
        amountError = true;
      });
      return 'Value must be greater than 0';
    }

    if (double.tryParse(value)! > (double.parse(asset['amount']) - 6)) {
      setState(() {
        amountError = true;
      });
      return 'You don\'t have sufficient balance';
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

    submitForm();
  }

  submitForm() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "destination": to,
        "memo": memo,
        "amount": amount.toString(),
        "assetCode": asset['assetCode'] == 'XBN' ? '' : asset['assetCode'],
        "assetIssuer": asset['assetIssuer'],
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/payment',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: walletDetails['walletPublicKey'],
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      popup(context, title: LanguageEn.error, message: e.toString());
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
    appState.viewData![ConfirmInitiatePaymentViewPageConfig.key] = data;
    appState.viewData![ConfirmInitiatePaymentViewPageConfig.key]['walletInfo'] =
        asset['walletInfo'];
    appState.viewData![ConfirmInitiatePaymentViewPageConfig.key]["usdPrice"] =
        asset['usdPrice'];
    appState.viewData![ConfirmInitiatePaymentViewPageConfig.key]["rel"] =
        asset['rel'];

    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: ConfirmInitiatePaymentViewPageConfig,
    );
  }
}

class _Utf8LengthLimitingTextInputFormatter extends TextInputFormatter {
  _Utf8LengthLimitingTextInputFormatter(this.maxLength)
      : assert(maxLength == null || maxLength == -1 || maxLength > 0);

  final int maxLength;

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    if (maxLength != null &&
        maxLength > 0 &&
        bytesLength(newValue.text) > maxLength) {
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
