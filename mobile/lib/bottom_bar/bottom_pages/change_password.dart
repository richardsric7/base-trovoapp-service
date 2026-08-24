import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/utils/local_auth.dart';
import '../../custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/custtom_textfild/custtom_password.dart';
import '../../storage/state.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
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
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
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
                      children: [
                        Text(
                          "changepassword".tr(),
                          style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 26,
                            fontFamily: fontsemibold,
                          ),
                        ),
                        SizedBox(height: height / 10),
                        // Old Password
                        CustomPasswordFormField(
                          "oldpassword".tr(),
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70,
                          300,
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
                          "newpassword".tr(),
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70,
                          300,
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
                          "confirmPassword".tr(),
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70,
                          300,
                          validator: validateConfirmPassword,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 20),
              Button(
                "changepassword".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  final form = _formKey.currentState;
                  if (!form!.validate()) {
                    return;
                  }
                  form.save();
                  changePassword();
                },
              ),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom,
                ),
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
      appState.viewData = {
        SuccessViewPageConfig.key: {
          'title': "success".tr(),
          'message': "passwordchangesuccessful".tr(),
        },
      };
      appState.currentAction = PageAction(
        state: PageState.replace,
        page: SuccessViewPageConfig,
      );
    } catch (e) {
      popup(context, title: "error".tr(), message: "somethingwentwrong".tr());
    }
  }

  String? validateOldPassword(value) {
    if (value.isEmpty) {
      //return "Enter a password";
      return "passwordemptyerror".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return "hinterrorpassword".tr();
    }
    if (oldPassword != appState.password!) {
      return "invalidpassword".tr();
    }

    return null;
  }

  String? validateNewPassword(value) {
    if (value.isEmpty) {
      // return "Confirm your password";
      return "newpasswordemptyerror".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return "hinterrorpassword".tr();
    }

    return null;
  }

  String? validateConfirmPassword(value) {
    if (value.isEmpty) {
      // return "Confirm your password";
      return "confirmnewpasswordemptyerror".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return "hinterrorpassword".tr();
    }

    if (newPassword != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return "passwordmismatcherror".tr();
    }

    return null;
  }
}
