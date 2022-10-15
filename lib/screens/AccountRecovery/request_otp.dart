import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:otp_text_field/otp_text_field.dart';
import 'package:otp_text_field/style.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class RequestOtp extends StatefulWidget {
  const RequestOtp({Key? key}) : super(key: key);

  @override
  State<RequestOtp> createState() => _RequestOtp();
}

class _RequestOtp extends State<RequestOtp> {
  late ColorNotifier notifier;
  bool otpSent = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  String email = '';
  String otp = '';

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 20),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 50),
                Text(
                  LanguageEn.account,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 30.sp,
                      fontFamily: fontsemibold),
                ),
                Text(
                  LanguageEn.recovery,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 30.sp,
                      fontFamily: fontsemibold),
                ),
                SizedBox(height: height / 20),
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
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20.0, vertical: 15.0),
                          child: Column(
                            children: [
                              Container(
                                width: width / 1.3,
                                child: Text(
                                  LanguageEn.requestotpwarn,
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
                                ),
                              ),
                              SizedBox(height: 2),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                if (otpSent) ...[
                  enterOTP(),
                ] else ...[
                  requestOTP(),
                ],
                SizedBox(height: height / 10),
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget requestOTP() {
    return Column(
      children: [
        SizedBox(height: height / 20),
        CustomTextFormField.textField(
          LanguageEn.enteryouremailaddress,
          notifier.getbluecolor,
          Icons.email,
          notifier.getgrey,
          notifier.getprefixicon,
          notifier.getblck,
          notifier.getgrey,
          70.sp,
          300.sp,
          validator: validateEmail,
          onSaved: (value) {
            print('email: $value');
            email = value.trim().replaceAll(' ', '');
          },
          keyboardtype: TextInputType.emailAddress,
        ),
        SizedBox(height: height / 40),
        Button(
          LanguageEn.requestotp,
          notifier.getbluecolor,
          wihitecolor,
          onTap: () {
            // setState(() {
            //   otpSent = true;
            // });
            validateAndProceed();
          },
        ),
        SizedBox(height: height / 70),
        TextButton(
          onPressed: () {
            setState(() {
              otpSent = true;
            });
          },
          child: Text(
            LanguageEn.alreadyhaveotp,
            style: TextStyle(
                color: notifier.getdarkgrey,
                fontSize: 15.sp,
                fontFamily: fontbody),
          ),
        ),
      ],
    );
  }

  Widget enterOTP() {
    return Column(
      children: [
        SizedBox(height: height / 20),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 30),
          child: OTPTextField(
            length: 6,
            width: MediaQuery.of(context).size.width,
            fieldWidth: 40,
            style: TextStyle(color: notifier.getblck, fontFamily: fontbody),
            textFieldAlignment: MainAxisAlignment.spaceAround,
            fieldStyle: FieldStyle.box,
            otpFieldStyle: OtpFieldStyle(
              borderColor: Colors.black38,
            ),
            onChanged: (pin) {
              otp = pin;
            },
            onCompleted: (pin) {
              otp = pin;
              verifyOTPAndProceed();
            },
          ),
        ),
        SizedBox(height: height / 70),
        TextButton(
          onPressed: () {
            setState(() {
              otpSent = false;
            });
          },
          child: Text(
            LanguageEn.resendotp,
            style: TextStyle(
                color: notifier.getdarkgrey,
                fontSize: 15.sp,
                fontFamily: fontbody),
          ),
        ),
        SizedBox(width: width / 10),
      ],
    );
  }

  String? validateEmail(String? value) {
    String pattern =
        r'^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$';
    RegExp regex = new RegExp(pattern);

    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return LanguageEn.emailvalidateempty;
    }

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return LanguageEn.emailvalidateinvalid;
    }

    return null;
  }

  validateAndProceed() {
    final form = _formKey.currentState;
    if (!form!.validate()) {
      return;
    }

    form.save();
    sendOTPRequest();
  }

  void sendOTPRequest() async {
    print('sending otp request.............');

    try {
      showLoader(context);

      var primaryWallet = appState.userInfo!.wallets!
          .firstWhere((wallet) => wallet.primaryWallet == 1);

      Map responseData = await makePostRequest(
        uri:
            '/v1/account/recovery/request-email-otp/${appState.userInfo!.username}',
        body: "",
        signer: primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: primaryWallet.publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        setState(() {
          otpSent = true;
        });
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  void verifyOTPAndProceed() async {
    print('sending otp request.............');

    try {
      showLoader(context);

      var primaryWallet = appState.userInfo!.wallets!
          .firstWhere((wallet) => wallet.primaryWallet == 1);

      Map responseData = await makePostRequest(
        uri:
            '/v1/account/recovery/verify-email-otp/${appState.userInfo!.username}/$otp',
        body: "",
        signer: primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: primaryWallet.publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        setState(() {});
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
