import 'package:app_settings/app_settings.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/route_manager.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/screens/Auth/fingerprint.dart';
import 'package:trovo_wallet/screens/Auth/signup.dart';
import 'package:trovo_wallet/screens/Backup/ensure_privacy.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import '../Custom_BlocObserver/notifire_clor.dart';
import '../router/PageActions.dart';
import '../router/ui_pages.dart';
import '../storage/state.dart';
import '../utils/enstring.dart';

late ColorNotifier notifier;

popup(context,
    {required String title,
    Color bodyColor = Colors.red,
    required String message}) async {
  notifier = Provider.of<ColorNotifier>(context, listen: false);
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
  notifier = Provider.of<ColorNotifier>(context, listen: false);
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
  notifier = Provider.of<ColorNotifier>(context, listen: false);
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
                              LanguageEn.ensureprivacybackup,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
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

void warnSkipBackupDialog(context, onSkip) {
  notifier = Provider.of<ColorNotifier>(context, listen: false);
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
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        LanguageEn.important,
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

void showResponseMessage(context, message, successAction) {
  notifier = Provider.of<ColorNotifier>(context, listen: false);
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
                              message,
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
                            color: notifier.getwihitecolor,
                            fontFamily: fontbody),
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

void showSuccessAlert(context, {onTap}) {
  notifier = Provider.of<ColorNotifier>(context, listen: false);
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
