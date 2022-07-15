import 'dart:convert';
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SendAsset extends StatefulWidget {
  const SendAsset({Key? key}) : super(key: key);

  @override
  State<SendAsset> createState() => _SendAsset();
}

class _SendAsset extends State<SendAsset> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  Wallet? activeWallet;
  final formKey = GlobalKey<FormState>();
  String to = ''; // the reciever
  double amount = 0;
  String? memo;
  String activeAsset = '';
  TextEditingController _utf8TextController = TextEditingController();

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
    activeAsset = appState.viewData![SendAssetViewPageConfig.key]['assetCode'];
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
                      '${LanguageEn.send} ${activeAsset}',
                      style: TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                          color: notifier.getbluecolor,
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
                        color: notifier.getbluecolor,
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
                notifier.getwihitecolor,
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
                      70.sp,
                      300.sp,
                      helperText: amount > 0 ? "$amount $activeAsset" : "",
                      keyboardtype: TextInputType.number,
                      validator: validateAmount,
                      onSaved: (value) => amount =
                          double.parse(value.trim().replaceAll(' ', '')),
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
                      70.sp,
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
                            style: Theme.of(context).textTheme.caption,
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
    // reciever cannot be empty
    if (value!.isEmpty) return 'Please enter amount to send';

    if (double.tryParse(value) == null) return 'Please enter a valid amount';

    if (double.tryParse(value)! <= 0) return 'Value must be greater than 0';

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
      // make initial request to the server using the
      // following credentials
      Map map = {
        "destination": to,
        "memo": memo,
        "amount": amount.toString(),
        "assetCode": activeAsset == 'XBN' ? '' : activeAsset,
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/payment',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
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
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  postProcessData(messageShown, messageLength, data) {
    print('messageShown: $messageShown messageLength $messageLength');
    // we would like to display all messages returned from the initial
    // request to server using a popup. In order to achieve that we
    // employ the use of a little recursion here. Please recursive
    // functions can turn into a nightmare fast so be carefull here.
    if (messageShown <= messageLength - 1) {
      showResponseMessage(
          context,
          data['messages'][messageShown],
          () => {
                print('postProcessData: $messageShown'),
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    // go to the definition of appState.viewData
    // to learn more about viewData
    appState.viewData![ConfirmTransactionViewPageConfig.key] = data;
    print(appState.viewData);

    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: ConfirmTransactionViewPageConfig,
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
