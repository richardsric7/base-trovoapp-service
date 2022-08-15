import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/services/push_fcm_service.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../Models/User.dart';
import '../../storage/state.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../../widgets/popups.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class PasswordMgtView extends StatefulWidget {
  const PasswordMgtView({Key? key}) : super(key: key);

  @override
  State<PasswordMgtView> createState() => _PasswordMgtViewState();
}

class _PasswordMgtViewState extends State<PasswordMgtView> {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  String? secretKey;
  String oldPassword = '';
  String newPassword = '';
  late DataProvider appState;

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
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: height / 10),
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          LanguageEn.changepassword,
                          style: TextStyle(
                              color: notifier.getblck,
                              fontSize: 26.sp,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(height: height / 10),
                        // Old Password
                        CustomPasswordFormField(
                          LanguageEn.oldpassword,
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70.sp,
                          300.sp,
                          onChanged: (value) {
                            setState(() {
                              oldPassword = value!.trim().replaceAll(' ', '');
                            });
                          },
                          validator: validateOldPassword,
                        ),
                        SizedBox(height: height / 40),
                        // New Password
                        CustomPasswordFormField(
                          LanguageEn.newpassword,
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70.sp,
                          300.sp,
                          onChanged: (value) {
                            setState(() {
                              newPassword = value!.trim().replaceAll(' ', '');
                            });
                          },
                          validator: validateNewPassword,
                        ),
                        SizedBox(height: height / 80),
                        // New Password
                        CustomPasswordFormField(
                          LanguageEn.confirmPassword,
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70.sp,
                          300.sp,
                          validator: validateConfirmPassword,
                        ),
                      ],
                    ),
                  )
                ],
              ),
              SizedBox(height: height / 20),
              if (appState.biometricEnabled) ...[
                Button(
                  LanguageEn.authorizewithbiometrics,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: () {
                    final form = _formKey.currentState;
                    if (!form!.validate()) {
                      return;
                    }
                    form.save();
                    toggleSwitch();
                  },
                ),
              ] else ...[
                Button(
                  LanguageEn.changepassword,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: () {
                    final form = _formKey.currentState;
                    if (!form!.validate()) {
                      return;
                    }
                    form.save();
                    changePassword();
                  },
                ),
              ],
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void toggleSwitch() async {
    try {
      bool result = await Authenticator().authenticateMe();
      if (result) {
        changePassword();
        // aparently we need the code below to make the
        // screen updata to show loader
        // after authorizing with biometrics
        setState(() {});
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  changePassword() async {
    try {
      await StoreData().storeInsertData('password', newPassword);
      appState.setPassword = await StoreData().storeGetData('password');
      showSuccessAlert(context, onTap: () {});
    } catch (e) {
      print(e);
      popup(context,
          title: LanguageEn.error, message: LanguageEn.somethingwentwrong);
    }
  }

  String? validateOldPassword(value) {
    print('old password: $value');
    if (value.isEmpty) {
      //return "Enter a password";
      return LanguageEn.passwordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }
    print('checking old password');
    if (oldPassword != appState.password!) {
      return LanguageEn.invalidpassword;
    }

    return null;
  }

  String? validateNewPassword(value) {
    print('new password: ${value.trim().replaceAll(' ', '')}');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.newpasswordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    return null;
  }

  String? validateConfirmPassword(value) {
    print(
        'confirm password: ${value.trim().replaceAll(' ', '')} & $newPassword');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.confirmnewpasswordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    if (newPassword != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return LanguageEn.passwordmismatcherror;
    }

    return null;
  }
}
