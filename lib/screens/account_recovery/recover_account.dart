import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:otp_text_field/otp_text_field.dart';
import 'package:otp_text_field/style.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import '../../network/requests.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../../widgets/popups.dart';

class RecoverAccount extends StatefulWidget {
  const RecoverAccount({Key? key}) : super(key: key);

  @override
  State<RecoverAccount> createState() => _RecoverAccountState();
}

class _RecoverAccountState extends State<RecoverAccount> {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  String? username;
  late DataProvider appState;
  bool otpSent = false;

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
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
                context, notifier.getwihitecolor, "", notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              SizedBox(width: width / 15),
              Form(
                key: _formKey,
                child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          "account".tr(),
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 26,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(
                          width: width / 50,
                        ),
                        Text(
                          "recovery".tr(),
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 26,
                              fontFamily: fontsemibold),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 15),
                    if (otpSent) ...[
                      enterOTP(),
                    ] else ...[
                      requestOTP(),
                    ],
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget enterOTP() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
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
                          "otpsentinfo".tr(),
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
        SizedBox(height: height / 50),
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
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
                        child: Wrap(
                          children: [
                            Text(
                              "enterotp".tr(),
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            SizedBox(
                              width: width / 70,
                            ),
                            Text(
                              appState.tempUsername,
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold),
                            ),
                          ],
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
        SizedBox(height: height / 50),
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 30),
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
                  appState.tempEmailOtp = pin;
                },
                onCompleted: (pin) {
                  appState.tempEmailOtp = pin;
                  verifyOTPAndProceed();
                },
              ),
            ),
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
            "resendotp".tr(),
            style: TextStyle(
                color: notifier.getdarkgrey,
                fontSize: 15,
                fontFamily: fontbody),
          ),
        ),
        SizedBox(width: width / 10),
      ],
    );
  }

  void verifyOTPAndProceed() async {
    print('sending otp request.............');

    try {
      showLoader(context);

      Map responseData = await makePostRequest(
        uri:
            '/v1/account/recovery/verify-email-otp/${appState.tempUsername}/${appState.tempEmailOtp}',
        body: "",
        signer: appState.tempPublicKey,
        secretKey: appState.tempSecretKey, // the primary wallet secret key
        publicKey: appState.tempPublicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        haveYouSetupSecurityQuestionsPopup(context, onYes: () {
          setState(() {
            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: AnswerSecurityQuestionsViewPageConfig);
          });
        }, onNo: () {
          setState(() {
            appState.viewData = {
              SecurityQuestionsForInactiveAccountsViewPageConfig.key: {
                'signer': appState.tempPublicKey,
                'publicKey': appState.tempPublicKey,
                'secretKey': appState.tempSecretKey,
                'username': appState.tempUsername,
              }
            };
            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: SecurityQuestionsForInactiveAccountsViewPageConfig);
            appState.viewData![EnsurePrivacyPageConfig.key] = {
              'rel': 'restoreUnactivatedAccount',
            };
            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: SecurityQuestionsForInactiveAccountsViewPageConfig);
          });
        });
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Widget requestOTP() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
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
                          "provideusernameforaccountrecovery".tr(),
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
        SizedBox(height: height / 50),
        // Email address
        CustomTextFormField.textField(
          "enteryourusername".tr(),
          notifier.getbluecolor,
          Icons.person,
          notifier.getgrey,
          notifier.getprefixicon,
          notifier.getblck,
          notifier.getgrey,
          70,
          300,
          validator: (value) {
            var trimmedVal = value!.trim().replaceAll(' ', '');
            if (trimmedVal.isEmpty) {
              return "usernameoremailempty".tr();
            }

            if (trimmedVal.length < 3) {
              return "usernameoremailinvalid".tr();
            }
          },
          onSaved: storeUsernameOrEmail,
          keyboardtype: TextInputType.emailAddress,
        ),
        SizedBox(height: height / 20),
        Button(
          "continuee".tr(),
          notifier.getbluecolor,
          wihitecolor,
          onTap: () => validateForm(),
        ),
        if (appState.tempUsername.isNotEmpty) ...[
          SizedBox(height: height / 70),
          TextButton(
            onPressed: () {
              setState(() {
                otpSent = true;
              });
            },
            child: Text(
              "alreadyhaveotp".tr(),
              style: TextStyle(
                  color: notifier.getdarkgrey,
                  fontSize: 15,
                  fontFamily: fontbody),
            ),
          ),
        ],
        SizedBox(height: height / 10),
        Padding(
          padding:
              EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
        ),
      ],
    );
  }

  String? storeUsernameOrEmail(String? value) {
    var currValue = value!.trim().replaceAll(' ', '');
    if (currValue.isEmpty) {
      return "emailvalidateempty".tr();
    }

    setState(() {
      username = currValue;
    });
    return null;
  }

  validateForm() async {
    print('saving form...');
    final form = _formKey.currentState;
    if (!form!.validate()) {
      return;
    }

    showLoader(context);
    form.save();

    // store these credentials and the password before making request
    // to get userinfo from the server. This way we can use these stored data
    // to create new user account if the provided user account does not exist
    var account = TrovoWalletSDK().createAccount();
    appState.setTempUsername = username;
    appState.setTempPublicKey = account.publicKey;
    appState.setTempSecretKey = account.secretKey;

    sendOTPRequest();
  }

  void sendOTPRequest() async {
    print('sending otp request.............');

    try {
      showLoader(context);

      Map responseData = await makePostRequest(
        uri: '/v1/account/recovery/request-email-otp/${username}',
        body: "",
        signer: appState.tempPublicKey,
        secretKey: appState.tempSecretKey, // the primary wallet secret key
        publicKey: appState.tempPublicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        setState(() {
          otpSent = true;
        });
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}
