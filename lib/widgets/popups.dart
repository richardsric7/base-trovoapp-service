import 'package:app_settings/app_settings.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:sembast/sembast.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import '../Custom_BlocObserver/notifire_clor.dart';
import '../router/PageActions.dart';
import '../router/ui_pages.dart';
import '../storage/state.dart';
import '../utils/enstring.dart';

// late ColorNotifier notifier;

popup(context,
    {required String title,
    Color bodyColor = Colors.red,
    required String message}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        title,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              message,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: bodyColor,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: ElevatedButton(
                      onPressed: () =>
                          Navigator.of(context).pop(), // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

Future<bool?> biometricsErrorAlert(BuildContext context) {
  return showDialog<bool>(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(LanguageEn.fingerprintrequired),
          actions: [
            TextButton(
              child: Text(
                LanguageEn.cancel,
                style: TextStyle(
                    fontSize: 14.0,
                    fontFamily: fontbody,
                    fontWeight: FontWeight.bold,
                    color: Colors.teal),
              ),
              onPressed: () => Navigator.of(
                context,
                rootNavigator: true,
              ).pop(false),
            ),
            TextButton(
                child: Text(
                  LanguageEn.gotosettings,
                  style: TextStyle(
                      fontSize: 14.0,
                      fontFamily: fontbody,
                      fontWeight: FontWeight.bold,
                      color: Colors.teal),
                ),
                onPressed: () async {
                  await AppSettings.openSecuritySettings();
                }),
          ],
          content: Container(
            decoration: const BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.all(
                Radius.circular(30),
              ),
            ),
            child: Text(
              LanguageEn.fingerprintnotenabled,
              style: TextStyle(
                fontFamily: fontbody,
                fontSize: 15.0,
                fontWeight: FontWeight.w400,
              ),
            ),
          ),
        );
      });
}

Future<bool?> accountNotFoundPopup(BuildContext context) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.oops,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              LanguageEn.usernamenotfound,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: Colors.red,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.of(context).pop();
                        appState.currentAction = PageAction(
                            state: PageState.addPage, page: SignupPageConfig);
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style: TextStyle(
                            color: notifier.getwihitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () =>
                          Navigator.of(context).pop(), // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        elevation: MaterialStateProperty.all<double>(0),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getwihitecolor!),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getgrey,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.cancel,
                        style: TextStyle(
                            color: notifier.getbluecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void ensureBackupPrivacyDialog(context, action) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.important,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              LanguageEn.ensureprivacybackup,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: notifier.getbluewhitecolor,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: action,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () =>
                          Navigator.of(context).pop(), // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        elevation: MaterialStateProperty.all<double>(0),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getwihitecolor!),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getgrey,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.cancel,
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void showPasswordDialog(context, action) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  final formKey = GlobalKey<FormState>();
  String password = '';

  showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.password,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 25.0, vertical: 5.0),
                            child: Form(
                              key: formKey,
                              child: CustomPasswordFormField(
                                LanguageEn.password,
                                notifier.getbluewhitecolor,
                                Icons.lock,
                                notifier.getgrey,
                                notifier.getprefixicon,
                                notifier.getblck,
                                70.sp,
                                300.sp,
                                validator: (String? value) {
                                  if (value!.isEmpty)
                                    return 'Enter your password';

                                  if (value.length < 6)
                                    return 'Use 6 characters or more for your password';
                                  if (password != appState.password!) {
                                    return 'Invalid password';
                                  }
                                  return null;
                                },
                                onChanged: (value) {
                                  password = value!.trim().replaceAll(' ', '');
                                },
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        print('elevated button pressed...$password');
                        if (!formKey.currentState!.validate()) {
                          return;
                        }
                        action();
                        Navigator.of(context).pop(); // dismiss dialog,
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () =>
                          Navigator.of(context).pop(), // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        elevation: MaterialStateProperty.all<double>(0),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getwihitecolor!),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getgrey,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.cancel,
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void warnSkipBackupDialog(context, onSkip) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  // var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.important,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              LanguageEn.warnskipbackup,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: Colors.red,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () =>
                          Navigator.of(context).pop(), // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.cancel,
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () {
                        Navigator.of(context).pop();
                        onSkip();
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        elevation: MaterialStateProperty.all<double>(0),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getwihitecolor!),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getgrey,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.skip,
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void showSetSecurityQuestionsPopup(context) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20.0, 20.0, 20.0, 0.0),
                    child: Center(
                      child: Text(
                        "${LanguageEn.setup} ${LanguageEn.securityquestions}",
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              LanguageEn.pleasesetupsecurityquestions,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: notifier.getbluecolor,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: SecurityQuestionsViewPageConfig);
                        Navigator.of(context).pop();
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.proceed,
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void showResponseMessage(context, message, successAction) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.important,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 4.5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              message,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: notifier.getbluewhitecolor,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.of(context).pop();
                        successAction();
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style: TextStyle(
                          color: wihitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () {
                        Navigator.of(context).pop(); // dismiss dialog,
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        elevation: MaterialStateProperty.all<double>(0),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getwihitecolor!),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getgrey,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.cancel,
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void showSuccessAlert(context, {required onTap}) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.success,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Image.asset(
                              "assets/images/success.gif",
                              height: 125.0,
                              width: 125.0,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.of(context).pop(); // dismiss dialog,
                        onTap();
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getgreencolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style: TextStyle(
                            color: notifier.getwihitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

void imageSourceDialog(context, {onCamera, onGallery}) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  appState.dialogOpen = true;
  showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 20.0, horizontal: 5.0),
                            child: Text(
                              LanguageEn.chooseimagesource,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.bold,
                                color: notifier.getbluecolor,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: OutlinedButton(
                      onPressed: () {
                        onGallery();
                        Navigator.of(context).pop();
                        appState.dialogOpen = false;
                      },
                      // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.gallery,
                        style: TextStyle(
                            color: notifier.getwihitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () {
                        onCamera();
                        Navigator.of(context).pop();
                        appState.dialogOpen = false;
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        elevation: MaterialStateProperty.all<double>(0),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getwihitecolor!),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getgrey,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.camera,
                        style: TextStyle(
                            color: notifier.getbluecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

customDateRangePopup(context, {required void Function() onDone}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
          var appState = Provider.of<DataProvider>(context, listen: false);
          var initialDate = DateTime.now();
          var startDate = appState.filterStartDate ??
              DateTime.now().subtract(Duration(days: 1));
          var endDate = appState.filterEndDate ?? DateTime.now();
          return AlertDialog(
              // scrollable: true,
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(0),
              content: Container(
                decoration: BoxDecoration(
                  color: notifier.getwihitecolor,
                  borderRadius: BorderRadius.all(
                    Radius.circular(23),
                  ),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          'Enter the date range below',
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Container(
                      constraints: BoxConstraints(
                          maxHeight: height / 1.7, minWidth: width / 1.1),
                      // height: height / 5,
                      child: SingleChildScrollView(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Wrap(
                              children: [
                                quickDateRange(context, text: 'Past week',
                                    onPressed: () {
                                  appState.setFilterStartDate = DateTime.now()
                                      .subtract(Duration(days: 7));
                                  appState.setFilterEndDate = DateTime.now();
                                }),
                                quickDateRange(context, text: 'Past month',
                                    onPressed: () {
                                  var date = DateTime.now();
                                  appState.setFilterEndDate = date;
                                  appState.setFilterStartDate = DateTime(
                                      date.year, date.month - 1, date.day);
                                }),
                                quickDateRange(context, text: 'Past 3 months',
                                    onPressed: () {
                                  var date = DateTime.now();
                                  appState.setFilterEndDate = date;
                                  appState.setFilterStartDate = DateTime(
                                      date.year, date.month - 3, date.day);
                                }),
                              ],
                            ),
                            SizedBox(
                              height: height / 70,
                            ),
                            Padding(
                              padding: const EdgeInsets.symmetric(
                                  vertical: 10.0, horizontal: 5.0),
                              child: Column(
                                children: [
                                  Text(
                                    'Start Date',
                                    textAlign: TextAlign.start,
                                    style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 15,
                                        fontFamily: fontbody),
                                  ),
                                  SizedBox(
                                    height: height / 70,
                                  ),
                                  Container(
                                    decoration: BoxDecoration(
                                      borderRadius: const BorderRadius.all(
                                          Radius.circular(10.0)),
                                      color: notifier.isDark
                                          ? darktilewhitecolor
                                          : notifier.getaddsubwalletgrey,
                                    ),
                                    child: TextButton(
                                      onPressed: () async {
                                        appState.setFilterStartDate =
                                            await showDatePicker(
                                                  context: context,
                                                  initialDate: appState
                                                          .filterStartDate ??
                                                      initialDate,
                                                  firstDate: DateTime
                                                      .fromMicrosecondsSinceEpoch(
                                                          1000),
                                                  lastDate: DateTime.now(),
                                                ) ??
                                                appState.filterStartDate;
                                      },
                                      child: Wrap(children: [
                                        Text(
                                          DateFormat('MMMM dd, yyyy')
                                              .format(startDate),
                                          textAlign: TextAlign.start,
                                          style: TextStyle(
                                              color: notifier.getbluewhitecolor,
                                              fontSize: 15,
                                              fontFamily: fontsemibold),
                                        ),
                                        SizedBox(
                                          width: width / 50,
                                        ),
                                        Icon(
                                          Icons.edit,
                                          size: 16,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ]),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            Padding(
                              padding: const EdgeInsets.symmetric(
                                  vertical: 10.0, horizontal: 5.0),
                              child: Column(
                                children: [
                                  Text(
                                    'End Date',
                                    textAlign: TextAlign.start,
                                    style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 15,
                                        fontFamily: fontbody),
                                  ),
                                  SizedBox(
                                    height: height / 70,
                                  ),
                                  Container(
                                    decoration: BoxDecoration(
                                      borderRadius: const BorderRadius.all(
                                          Radius.circular(10.0)),
                                      color: notifier.isDark
                                          ? darktilewhitecolor
                                          : notifier.getaddsubwalletgrey,
                                    ),
                                    child: TextButton(
                                      onPressed: () async {
                                        appState.setFilterEndDate =
                                            await showDatePicker(
                                                  context: context,
                                                  initialDate:
                                                      appState.filterEndDate ??
                                                          initialDate,
                                                  firstDate: DateTime
                                                      .fromMicrosecondsSinceEpoch(
                                                          1000),
                                                  lastDate: DateTime.now(),
                                                ) ??
                                                appState.filterEndDate;
                                      },
                                      child: Wrap(
                                        crossAxisAlignment:
                                            WrapCrossAlignment.center,
                                        children: [
                                          Text(
                                            DateFormat('MMMM dd, yyyy')
                                                .format(endDate),
                                            textAlign: TextAlign.start,
                                            style: TextStyle(
                                                color:
                                                    notifier.getbluewhitecolor,
                                                fontSize: 15,
                                                fontFamily: fontsemibold),
                                          ),
                                          SizedBox(
                                            width: width / 50,
                                          ),
                                          Icon(
                                            Icons.edit,
                                            size: 16,
                                            color: notifier.getbluewhitecolor,
                                          ),
                                        ],
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
                          appState.setFilterStartDate =
                              appState.filterStartDate == null
                                  ? startDate
                                  : appState.filterStartDate;
                          appState.setFilterEndDate =
                              appState.filterEndDate == null
                                  ? endDate
                                  : appState.filterEndDate;
                          Navigator.of(context).pop(); // dismiss dialog,
                          onDone();
                        },
                        style: ButtonStyle(
                          fixedSize: MaterialStateProperty.all(
                            Size(width / 1.5, height / 20),
                          ),
                          backgroundColor: MaterialStateProperty.all<Color>(
                              notifier.getbluecolor),
                          shape:
                              MaterialStateProperty.all<RoundedRectangleBorder>(
                            const RoundedRectangleBorder(
                              borderRadius: BorderRadius.all(
                                Radius.circular(10),
                              ),
                            ),
                          ),
                        ),
                        child: Text(
                          LanguageEn.done,
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    SizedBox(height: height / 50),
                  ],
                ),
              ));
        });
      });
}

Widget quickDateRange(BuildContext context,
    {required String text, required void Function() onPressed}) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return Padding(
    padding: const EdgeInsets.all(3.0),
    child: Container(
      decoration: BoxDecoration(
        borderRadius: const BorderRadius.all(Radius.circular(10.0)),
        color:
            notifier.isDark ? darktilewhitecolor : notifier.getaddsubwalletgrey,
      ),
      child: Wrap(
        children: [
          TextButton(
            onPressed: onPressed,
            child: Text(
              text,
              overflow: TextOverflow.ellipsis,
              softWrap: true,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                  fontSize: 10.sp),
            ),
          )
        ],
      ),
    ),
  );
}

amountRangePopup(context, {required void Function() onDone}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  var minAmountTextController = TextEditingController();
  var maxAmountTextController = TextEditingController();
  var appState = Provider.of<DataProvider>(context, listen: false);
  minAmountTextController.text = appState.filterMinAmount ?? "";
  maxAmountTextController.text = appState.filterMaxAmount ?? "";

  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
          return AlertDialog(
              // scrollable: true,
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(0),
              content: Container(
                decoration: BoxDecoration(
                  color: notifier.getwihitecolor,
                  borderRadius: BorderRadius.all(
                    Radius.circular(23),
                  ),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          'Enter the amount range below',
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Container(
                      constraints: BoxConstraints(
                          maxHeight: height / 1.7, minWidth: width / 1.1),
                      // height: height / 5,
                      child: SingleChildScrollView(
                        child: Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 35.0),
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              CustomTextFormField.textFieldWithoutIcon(
                                'minimum amount',
                                notifier.getbluecolor,
                                notifier.getgrey,
                                notifier.getprefixicon,
                                notifier.getblck,
                                notifier.getgrey,
                                55.sp, 300.sp,
                                onChanged: (value) {
                                  if (value != null &&
                                      value.toString().isNotEmpty) {
                                    appState.setFilterMinAmount = value;
                                  }
                                },
                                controller: minAmountTextController,
                                keyboardtype: TextInputType.numberWithOptions(
                                    decimal: true),
                                // onSaved: (value) => amount = value.trim().replaceAll(' ', ''),
                              ),
                              SizedBox(
                                height: height / 30,
                              ),
                              CustomTextFormField.textFieldWithoutIcon(
                                'maximum amount',
                                notifier.getbluecolor,
                                notifier.getgrey,
                                notifier.getprefixicon,
                                notifier.getblck,
                                notifier.getgrey,
                                55.sp, 300.sp,
                                onChanged: (value) {
                                  if (value != null &&
                                      value.toString().isNotEmpty) {
                                    appState.setFilterMaxAmount = value;
                                  }
                                },
                                controller: maxAmountTextController,
                                keyboardtype: TextInputType.numberWithOptions(
                                    decimal: true),
                                // onSaved: (value) => amount = value.trim().replaceAll(' ', ''),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 30,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
                          onDone();
                          Navigator.of(context).pop(); // dismiss dialog,
                        },
                        style: ButtonStyle(
                          fixedSize: MaterialStateProperty.all(
                            Size(width / 1.5, height / 20),
                          ),
                          backgroundColor: MaterialStateProperty.all<Color>(
                              notifier.getbluecolor),
                          shape:
                              MaterialStateProperty.all<RoundedRectangleBorder>(
                            const RoundedRectangleBorder(
                              borderRadius: BorderRadius.all(
                                Radius.circular(10),
                              ),
                            ),
                          ),
                        ),
                        child: Text(
                          LanguageEn.done,
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    SizedBox(height: height / 50),
                  ],
                ),
              ));
        });
      });
}

textFieldPopup(context,
    {required FilterType rel, required void Function(String?) onDone}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  String? textValue;
  var textController = TextEditingController();
  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
          var appState = Provider.of<DataProvider>(context, listen: false);
          switch (rel) {
            case FilterType.Username:
              textController.text = appState.filterUsername ?? "";
              textValue = appState.filterUsername ?? "";
              break;
            case FilterType.FromPublicKey:
              textController.text = appState.filterFromPublicKey ?? "";
              textValue = appState.filterFromPublicKey ?? "";
              break;
            case FilterType.Memo:
              textController.text = appState.filterMemo ?? "";
              textValue = appState.filterMemo ?? "";
              break;
            default: // FilterType.ToPublicKey
              textController.text = appState.filterToPublicKey ?? "";
              textValue = appState.filterToPublicKey ?? "";
              break;
          }
          return AlertDialog(
              // scrollable: true,
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(0),
              content: Container(
                decoration: BoxDecoration(
                  color: notifier.getwihitecolor,
                  borderRadius: BorderRadius.all(
                    Radius.circular(23),
                  ),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          getLabelText(rel),
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Container(
                      constraints: BoxConstraints(
                          maxHeight: height / 1.7, minWidth: width / 1.1),
                      // height: height / 5,
                      child: SingleChildScrollView(
                        child: Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 35.0),
                          child: CustomTextFormField.textFieldWithoutIcon(
                            getPlaceholder(rel),
                            notifier.getbluecolor,
                            notifier.getgrey,
                            notifier.getprefixicon,
                            notifier.getblck,
                            notifier.getgrey,
                            55.sp,
                            300.sp,
                            onChanged: (value) {
                              if (value != null &&
                                  value.toString().isNotEmpty) {
                                textValue = value;
                              }
                            },
                            controller: textController,
                            keyboardtype: TextInputType.text,
                          ),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 30,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
                          onDone(textValue);
                          Navigator.of(context).pop(); // dismiss dialog,
                        },
                        style: ButtonStyle(
                          fixedSize: MaterialStateProperty.all(
                            Size(width / 1.5, height / 20),
                          ),
                          backgroundColor: MaterialStateProperty.all<Color>(
                              notifier.getbluecolor),
                          shape:
                              MaterialStateProperty.all<RoundedRectangleBorder>(
                            const RoundedRectangleBorder(
                              borderRadius: BorderRadius.all(
                                Radius.circular(10),
                              ),
                            ),
                          ),
                        ),
                        child: Text(
                          LanguageEn.done,
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    SizedBox(height: height / 50),
                  ],
                ),
              ));
        });
      });
}

transactionTypePopup(context,
    {required void Function() onAllSelected,
    required void Function() onSwapSelected,
    required void Function() onPaymentSelected}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
          var appState = Provider.of<DataProvider>(context, listen: false);
          return AlertDialog(
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(0),
              content: Container(
                decoration: BoxDecoration(
                  color: notifier.getwihitecolor,
                  borderRadius: BorderRadius.all(
                    Radius.circular(23),
                  ),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          'Select transaction type',
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Container(
                      constraints: BoxConstraints(
                          maxHeight: height / 1.7, minWidth: width / 1.1),
                      // height: height / 5,
                      child: SingleChildScrollView(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Wrap(
                              children: [
                                quickDateRange(context, text: 'All',
                                    onPressed: () {
                                  onAllSelected();
                                }),
                                quickDateRange(context, text: 'Swap',
                                    onPressed: () {
                                  onSwapSelected();
                                }),
                                quickDateRange(context, text: 'Payment',
                                    onPressed: () {
                                  onPaymentSelected();
                                }),
                              ],
                            ),
                            SizedBox(
                              height: height / 70,
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
              ));
        });
      });
}

String getLabelText(FilterType rel) {
  switch (rel) {
    case FilterType.FromPublicKey:
      return 'Enter from public key below';
    case FilterType.ToPublicKey:
      return 'Enter to public key below';
    case FilterType.Memo:
      return 'Enter to memo text below';
    default:
      return 'Enter username or full name below';
  }
}

String getPlaceholder(FilterType rel) {
  switch (rel) {
    case FilterType.FromPublicKey:
      return 'from public key';
    case FilterType.ToPublicKey:
      return 'to public key';
    case FilterType.Memo:
      return 'memo';
    default:
      return 'username';
  }
}
