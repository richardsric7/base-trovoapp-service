import 'dart:convert';
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/request_specific_payment_details.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class RequestSpecificPayment extends StatefulWidget {
  const RequestSpecificPayment({Key? key}) : super(key: key);

  @override
  State<RequestSpecificPayment> createState() => _RequestSpecificPayment();
}

class _RequestSpecificPayment extends State<RequestSpecificPayment>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  Wallet? activeWallet;
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
    print(
        'this is appState: ${appState.viewData![RequestSpecificPaymentViewPageConfig.key]}');
    asset = appState.viewData![RequestSpecificPaymentViewPageConfig.key];

    if (asset['deepLinkInfo'] != null) {
      deeplinkInfo = asset['deepLinkInfo'];
      print('deeplink is here....$deeplinkInfo');
      toController.text = deeplinkInfo['receiver'];
      amountController.text = deeplinkInfo['amount'];
      // amountController.text = '34';
      // amountController.selection = TextSelection.fromPosition(
      //     TextPosition(offset: amountController.text.length));
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
                      '${LanguageEn.request} ${getAssetCode(asset['assetCode'])}',
                      style: TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              formFields(),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.proceed,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  handleSubmit();
                },
              ),
              SizedBox(height: height / 7.3),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Widget formFields() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Form(
            key: formKey,
            child: Column(
              children: [
                SizedBox(
                  height: height / 50,
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
                      amount = value;
                    });
                  },
                  controller: amountController,
                  keyboardtype: TextInputType.numberWithOptions(decimal: true),
                  validator: validateAmount,
                  onSaved: (value) => amount = value.trim().replaceAll(' ', ''),
                ),
                SizedBox(height: height / 50),
                CustomTextFormField.textField(
                  LanguageEn.memo,
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
              ],
            ),
          ),
        ),
      ],
    );
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
    print('submitting form...');
    try {
      showLoader(context);
      Map responseData = await makeGetRequest(
        uri:
            '/v1/users/payment/generate/${appState.userInfo!.username}?paymentDestination=${activeWallet!.publicKey}&assetCode=${asset['assetCode']}&assetIssuer=${asset['assetIssuer']}&amount=${amount.toString()}&memo=${memo != null ? Uri.encodeComponent(memo!) : ''}',
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        print('this is responseData ${responseData['data']}');
        appState.viewData![RequestSpecificPaymentDetailsViewPageConfig.key] = {
          'qrCode': responseData['data']['qrCode'],
          'dynamicLink': responseData['data']['dynamicLink'],
          'amount': amount,
          'assetIssuer': asset['assetIssuer'],
          'assetCode': asset['assetCode'],
          'memo': memo,
        };
        appState.currentAction = PageAction(
            state: PageState.addPage,
            page: RequestSpecificPaymentDetailsViewPageConfig);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
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
