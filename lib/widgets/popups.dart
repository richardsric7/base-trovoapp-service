import 'dart:convert';

import 'package:app_settings/app_settings.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/models/wallets_list_view_data.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/deposit_withdrawal_history.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../custom_bloc_observer/notifire_clor.dart';
import '../router/page_actions.dart';
import '../router/ui_pages.dart';
import '../screens/shared_access/shared_access.dart';
import '../storage/state.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

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
                        "continuee".tr(),
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
          title: Text("biometricsrequired".tr()),
          actions: [
            TextButton(
              child: Text(
                "cancel".tr(),
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
                  "gotosettings".tr(),
                  style: TextStyle(
                      fontSize: 14.0,
                      fontFamily: fontbody,
                      fontWeight: FontWeight.bold,
                      color: Colors.teal),
                ),
                onPressed: () async {
                  await AppSettings.openAppSettings();
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
              "biometricsnotenabled".tr(),
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
                        "oops".tr(),
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
                              "accountnotfound".tr(),
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
                        appState.viewData![SignupPageConfig.key] = {
                          'importMode': true,
                        };
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
                        "continuee".tr(),
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
                        "cancel".tr(),
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

Future<bool?> accountNotFoundAfterSwitchPopup(
  BuildContext context, {
  required void Function() onContinueWithCredentials,
  required void Function() onImportNewCredential,
  required void Function() onGoBackToPrevEnvironment,
  String? message,
}) {
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
                        "oops".tr(),
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
                              message != null
                                  ? message
                                  : "accountnotfoundafterswitch"
                                      .tr(args: [appState.walletMode]),
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
                  if (message == null) ...[
                    Padding(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 10.0, vertical: 5.0),
                      child: ElevatedButton(
                        onPressed: () {
                          // Navigator.of(context).pop();
                          onContinueWithCredentials();
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
                          "continuewithcredentials".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getwihitecolor,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                  ],
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        onImportNewCredential();
                      },
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor90),
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
                        "importanotherwallet".tr(),
                        textAlign: TextAlign.center,
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
                        onGoBackToPrevEnvironment();
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
                        "gobacktopreviousnet".tr(args: [
                          appState.walletMode == 'Testnet'
                              ? 'Mainnet'
                              : 'Testnet'
                        ]),
                        textAlign: TextAlign.center,
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
                        "important".tr(),
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
                              "ensureprivacybackup".tr(),
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
                        "continuee".tr(),
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
                        "cancel".tr(),
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
                        "password".tr(),
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
                                "password".tr(),
                                notifier.getbluewhitecolor,
                                Icons.lock,
                                notifier.getgrey,
                                notifier.getprefixicon,
                                notifier.getblck,
                                70.sp,
                                300.sp,
                                validator: (String? value) {
                                  if (value!.isEmpty)
                                    return "pleaseenteryourpassword".tr();

                                  if (value.length < 6)
                                    return "use6charsormoreforpassword".tr();
                                  if (password != appState.password!) {
                                    return "invalidpassword".tr();
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
                        "continuee".tr(),
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
                        "cancel".tr(),
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
                        "important".tr(),
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
                              "warnskipbackup".tr(),
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
                        "cancel".tr(),
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
                        "skip".tr(),
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

void updateAppMessagePopup(context, message, Function() onTap) {
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
                    padding: const EdgeInsets.fromLTRB(20.0, 20.0, 20.0, 0.0),
                    child: Center(
                      child: Text(
                        "updateapp".tr(),
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
                        onTap();
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
                        "update".tr(),
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
                        "important".tr(),
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
                        "continuee".tr(),
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
                        "cancel".tr(),
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

void mintWalletExplainerPopup(context) {
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
                        "information".tr(),
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
                              "explainmintwallet".tr(),
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
                        "continuee".tr(),
                        style: TextStyle(
                          color: wihitecolor,
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
                        "success".tr(),
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
                        "continuee".tr(),
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
                              "chooseimagesource".tr(),
                              style: TextStyle(
                                fontSize: 17,
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
                        "gallery".tr(),
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
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
                        "camera".tr(),
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
                          "enterdaterange".tr(),
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
                                quickDateRange(context, text: "pastweek".tr(),
                                    onPressed: () {
                                  appState.setFilterStartDate = DateTime.now()
                                      .subtract(Duration(days: 7));
                                  appState.setFilterEndDate = DateTime.now();
                                }),
                                quickDateRange(context, text: "pastmonth".tr(),
                                    onPressed: () {
                                  var date = DateTime.now();
                                  appState.setFilterEndDate = date;
                                  appState.setFilterStartDate = DateTime(
                                      date.year, date.month - 1, date.day);
                                }),
                                quickDateRange(context, text: "past3month".tr(),
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
                                    "startdate".tr(),
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
                                    "enddate".tr(),
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
                          "done".tr(),
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
                          "enteramountrange".tr(),
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
                                "minimumamount".tr(),
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
                                "maximumamount".tr(),
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
                          "done".tr(),
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
    {required HistoryFilterType rel,
    required void Function(String?) onDone}) async {
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
            case HistoryFilterType.Username:
              textController.text = appState.filterUsername ?? "";
              textValue = appState.filterUsername ?? "";
              break;
            case HistoryFilterType.FromPublicKey:
              textController.text = appState.filterFromPublicKey ?? "";
              textValue = appState.filterFromPublicKey ?? "";
              break;
            case HistoryFilterType.Memo:
              textController.text = appState.filterMemo ?? "";
              textValue = appState.filterMemo ?? "";
              break;
            default: // HistoryFilterType.ToPublicKey
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
                          "done".tr(),
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
                          "selecttransactiontype".tr(),
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
                                quickDateRange(context, text: "all".tr(),
                                    onPressed: () {
                                  onAllSelected();
                                }),
                                quickDateRange(context, text: "swap".tr(),
                                    onPressed: () {
                                  onSwapSelected();
                                }),
                                quickDateRange(context, text: "payment".tr(),
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

wrappedAssetTransactionTypePopup(
  context, {
  required void Function() onWithdrawSelected,
  required void Function() onDepositSelected,
}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
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
                          "selecttransactiontype".tr(),
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
                                quickDateRange(context,
                                    text: "deposithistory".tr(), onPressed: () {
                                  onDepositSelected();
                                }),
                                quickDateRange(context,
                                    text: "withdrawalhistory".tr(),
                                    onPressed: () {
                                  onWithdrawSelected();
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

wrappedAssetTransactionStatusPopup(
  context, {
  required void Function() onPendingSelected,
  required void Function() onCompletedSelected,
}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
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
                          "choosestatus".tr(),
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
                                quickDateRange(context, text: "pending_2".tr(),
                                    onPressed: () {
                                  onPendingSelected();
                                }),
                                quickDateRange(context, text: "completed".tr(),
                                    onPressed: () {
                                  onCompletedSelected();
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

String getLabelText(HistoryFilterType rel) {
  switch (rel) {
    case HistoryFilterType.FromPublicKey:
      return "enterfrompublickey".tr();
    case HistoryFilterType.ToPublicKey:
      return "entertopublickey".tr();
    case HistoryFilterType.Memo:
      return "entermemotext".tr();
    default:
      return "enterusernameorfullname".tr();
  }
}

String getPlaceholder(HistoryFilterType rel) {
  switch (rel) {
    case HistoryFilterType.FromPublicKey:
      return "frompublickey".tr();
    case HistoryFilterType.ToPublicKey:
      return "topublickey".tr();
    case HistoryFilterType.Memo:
      return "memo".tr();
    default:
      return "username".tr();
  }
}

void haveYouSetupSecurityQuestionsPopup(context,
    {required void Function() onYes, required void Function() onNo}) {
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
                        "important".tr(),
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
                              "haveyousetupsecurityquestions".tr(),
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
                        onYes();
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
                        "ihavesetupsecurityquestions".tr(),
                        textAlign: TextAlign.center,
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
                        onNo();
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
                        "ihavenotsetupsecurityquestions".tr(),
                        textAlign: TextAlign.center,
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

shareAccessInfoPopup(context) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;

  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
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
                    padding: const EdgeInsets.fromLTRB(20.0, 20.0, 20.0, 10.0),
                    child: Center(
                      child: Text(
                        "sharedaccess".tr(),
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 1.8,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 20.0),
                            child: RichText(
                              text: TextSpan(
                                text: "welcometosharedaccess".tr(),
                                style: TextStyle(
                                  fontSize: 17,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                                children: [
                                  TextSpan(
                                    text: '${"vieweraccess".tr()} ',
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: "describevieweraccess".tr(),
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontbody,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: '${"initiatoraccess".tr()} ',
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: "describeinitiatoraccess".tr(),
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontbody,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: '${"approveraccess".tr()} ',
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: "describeapproveraccess".tr(),
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontbody,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: '${"note".tr()}: ',
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  TextSpan(
                                    text: "moresharedaccessdetails".tr(),
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontbody,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                ],
                              ),
                              textAlign: TextAlign.justify,
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
                        "done".tr(),
                        style:
                            TextStyle(color: wihitecolor, fontFamily: fontbody),
                      ),
                    ),
                  ),
                ],
              ),
            ));
      });
}

void rejectionReasonPopup(context, void Function(String) action) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  final formKey = GlobalKey<FormState>();
  String reason = '';

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
                        "rejecttransaction".tr(),
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
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
                            child: Center(
                              child: Text(
                                "rejectionreason".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 15,
                                    fontFamily: fontbody),
                              ),
                            ),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 25.0, vertical: 5.0),
                            child: Form(
                              key: formKey,
                              child: CustomTextFormField.textFieldWithoutIcon(
                                "enterreason".tr(),
                                notifier.getbluewhitecolor,
                                notifier.getgrey,
                                notifier.getgrey,
                                notifier.getblck,
                                notifier.getgrey,
                                70.sp,
                                300.sp,
                                validator: (String? value) {
                                  if (value!.isEmpty)
                                    return "pleaseenterreason".tr();

                                  if (value.length < 5)
                                    return "reasonmustbe5ormorechars".tr();

                                  return null;
                                },
                                onChanged: (value) {
                                  reason = value!.trim();
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
                        if (!formKey.currentState!.validate()) {
                          return;
                        }
                        action(reason);
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
                        "continuee".tr(),
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
                        "cancel".tr(),
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

approvalListTransactionTypePopup(context, List<String> options, String label,
    void Function(String) onSelected,
    {ApprovalsListFilterType? rel}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  bool isChecked = appState.excludeUserApproved == 0;
  return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
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
                          label,
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
                              alignment: WrapAlignment.center,
                              children: [
                                for (var i = 0; i < options.length; i++) ...[
                                  quickDateRange(context,
                                      text: options[i].capitalizeFirst!,
                                      onPressed: () {
                                    onSelected(options[i]);
                                  }),
                                ]
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Transform.scale(
                          scale: 1.sp,
                          child: Checkbox(
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.all(
                                  Radius.circular(5.sp),
                                ),
                              ),
                              activeColor: notifier.getbluecolor,
                              side:
                                  BorderSide(color: notifier.getbluewhitecolor),
                              value: isChecked,
                              onChanged: (value) {
                                appState.setExcludeUserApproved =
                                    appState.excludeUserApproved == 1 ? 0 : 1;
                                isChecked = appState.excludeUserApproved == 0;
                                setStateForDialog(() {});
                              }),
                        ),
                        Container(
                          width: width / 1.7,
                          child: Text(
                            "includealreadysignedtransactions".tr(),
                            overflow: TextOverflow.visible,
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                  ],
                ),
              ));
        });
      });
}

approvalTextFieldPopup(context,
    {required String label,
    required String value,
    required String placeholder,
    required void Function(String?) onDone}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  String? textValue;
  var textController = TextEditingController();
  textController.text = value;
  textValue = value;
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
                          label,
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
                            placeholder,
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
                          "done".tr(),
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

void showChooseWalletPopup(context, assetCode, assetIssuer,
    {required void Function(String, bool) onDone,
    required void Function() onCancel}) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  String selectedWallet = '';
  var filteredWallets = {};
  appState.transactionableWallets.forEach((key, value) {
    if (value['claimedAssets'] != null) {
      for (var i = 0; i < value['claimedAssets'].length; i++) {
        if (value['claimedAssets'][i]['assetIssuer'] == assetIssuer &&
            value['claimedAssets'][i]['assetCode'] == assetCode) {
          filteredWallets[key] = value;
        }
      }
    }
  });

  List<DropdownMenuItem<String>> walletDropdownItems(bool isSelected) {
    var walletsList = <DropdownMenuItem<String>>[];
    filteredWallets.forEach((key, value) {
      walletsList.add(
        DropdownMenuItem(
          child: Row(
            children: [
              Container(
                constraints:
                    isSelected ? BoxConstraints(maxWidth: width / 3) : null,
                child: Text(
                  value['alias'],
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              if (value['sharedAccessEnabled'] == 1) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.people_outline,
                  size: 17,
                  color: notifier.getbluewhitecolor,
                )
              ],
              if (!isSelected && key == selectedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.check,
                  size: 18,
                  color: notifier.getbluewhitecolor,
                )
              ],
            ],
          ),
          value: key,
        ),
      );
    });

    return walletsList;
  }

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
                        "selectsendingwallet".tr(),
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Row(
                    children: [
                      SizedBox(
                        width: width / 15,
                      ),
                      Expanded(
                          child: dropdown(
                        (newValue) {
                          onDone(
                              newValue.toString(),
                              filteredWallets[newValue]
                                      ['sharedAccessEnabled'] ==
                                  1);
                          Navigator.of(context).pop(); // dismiss dialog,
                        },
                        walletDropdownItems(false),
                        selectedWallet.toString().isEmpty
                            ? null
                            : selectedWallet,
                        "choosewallet".tr(),
                        context,
                        (context) {
                          return walletDropdownItems(true);
                        },
                      )),
                      SizedBox(
                        width: width / 15,
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: OutlinedButton(
                      onPressed: () {
                        onCancel();
                      },
                      // dismiss dialog,
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
                        "cancel".tr(),
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

wrappedAssetsTextFieldPopup(context,
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
          textController.text = appState.filterWithdrawalAddress;
          textValue = appState.filterWithdrawalAddress;

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
                          "enterwithdrawaladdress".tr(),
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
                            "withdrawaladdress".tr(),
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
                          "done".tr(),
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

void addCustomAssetPopup(context, void Function(String, String) action) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  final formKey = GlobalKey<FormState>();
  String assetIssuer = '';
  String assetCode = '';

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
                        "addasset".tr(),
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
                              child: Column(
                                children: [
                                  CustomTextFormField.textFieldWithoutIcon(
                                    "enterassetcode".tr(),
                                    notifier.getbluewhitecolor,
                                    notifier.getgrey,
                                    notifier.getgrey,
                                    notifier.getblck,
                                    notifier.getgrey,
                                    70.sp,
                                    300.sp,
                                    validator: (String? value) {
                                      if (value!.isEmpty)
                                        return "enterassetcodeplease".tr();

                                      return null;
                                    },
                                    onChanged: (value) {
                                      assetCode = value!.trim();
                                    },
                                  ),
                                  SizedBox(height: 5),
                                  CustomTextFormField.textFieldWithoutIcon(
                                    "enterissuerpublickey".tr(),
                                    notifier.getbluewhitecolor,
                                    notifier.getgrey,
                                    notifier.getgrey,
                                    notifier.getblck,
                                    notifier.getgrey,
                                    70.sp,
                                    300.sp,
                                    validator: (String? value) {
                                      if (value!.isEmpty)
                                        return "enterissuerpublickeyplease"
                                            .tr();

                                      return null;
                                    },
                                    onChanged: (value) {
                                      assetIssuer = value!.trim();
                                    },
                                  ),
                                ],
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
                        if (!formKey.currentState!.validate()) {
                          return;
                        }
                        action(assetCode, assetIssuer);
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
                        "continuee".tr(),
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
                        "cancel".tr(),
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

void warnDisableSharedAccessDialog(context, void Function() onDisable) {
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
                        "important".tr(),
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
                              "warnDisableSharedAccess".tr(),
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
                        onDisable();
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
                        "disable".tr(),
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
                        "cancel".tr(),
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

void viewOnlySharedWalletOptions(
    context, void Function() onModify, void Function() onDisable) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  showDialog(
      context: context,
      barrierDismissible: true,
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
                        "sharedaccess".tr(),
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontFamily: fontsemibold),
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 5.0),
                    child: ElevatedButton(
                      onPressed: () {
                        Navigator.of(context).pop();
                        onModify();
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
                        "modifysharedaccess".tr(),
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
                        onDisable();
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
                        "disablesharedaccess".tr(),
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

showDocumentUploadPopup(context, String title,
    {required void Function() onDone,
    required List<DropdownMenuItem<String>> dropdownItems}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
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
                width: width / 1.1,
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
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontsemibold),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: dropdown(
                        (value) {},
                        dropdownItems,
                        null,
                        "purchasereceipt".tr(),
                        context,
                        null,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(3.0),
                      child: Container(
                        width: width / 2.5,
                        decoration: BoxDecoration(
                          borderRadius:
                              const BorderRadius.all(Radius.circular(10.0)),
                          color: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                        ),
                        child: Container(
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              TextButton(
                                onPressed: () => {},
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Icon(
                                      Icons.file_copy_outlined,
                                      size: 20,
                                      // color: notifier.getbluewhitecolor,
                                    ),
                                    SizedBox(
                                      width: width / 50,
                                    ),
                                    Text(
                                      "selectfile".tr(),
                                      style: TextStyle(
                                        fontSize: 15,
                                        fontFamily: fontbody,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(3.0),
                      child: Container(
                        width: width / 2.5,
                        decoration: BoxDecoration(
                          borderRadius:
                              const BorderRadius.all(Radius.circular(10.0)),
                          color: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                        ),
                        child: Container(
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              TextButton(
                                onPressed: () => {},
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Icon(
                                      CupertinoIcons.camera,
                                      size: 20,
                                      // color: notifier.getbluewhitecolor,
                                    ),
                                    SizedBox(
                                      width: width / 50,
                                    ),
                                    Text(
                                      "takephoto".tr(),
                                      style: TextStyle(
                                        fontSize: 15,
                                        fontFamily: fontbody,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: CustomTextFormField.textField(
                            "urltofile".tr(),
                            notifier.getbluecolor,
                            null,
                            notifier.getgrey,
                            null,
                            notifier.getblck,
                            notifier.getgrey,
                            70.sp,
                            200.sp,
                            // controller: referrerController,
                            // validator: validateReferrer,
                            onSaved: (value) {},
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
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
                          "upload".tr(),
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

showSubscribePopup(context,
    {required void Function() onDone,
    required List<DropdownMenuItem<String>> dropdownItems}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
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
                width: width / 1.1,
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
                    Image.asset(
                      "assets/images/thinking_man.png",
                      height: height / 4,
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "confirmsubscribewithwallet".tr(args: ['Atlantis 1']),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: dropdown(
                        (value) {},
                        dropdownItems,
                        null,
                        "choosewallet".tr(),
                        context,
                        null,
                      ),
                    ),
                    // SizedBox(
                    //   height: height / 70,
                    // ),
                    // Padding(
                    //   padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    //   child: Row(
                    //     children: [
                    //       Text(
                    //         'Change wallet',
                    //         textAlign: TextAlign.center,
                    //         style: TextStyle(
                    //             color: notifier.getbluewhitecolor,
                    //             fontSize: 12,
                    //             fontFamily: fontsemibold),
                    //       ),
                    //     ],
                    //   ),
                    // ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
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
                          "iwanttosubscribe".tr(),
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10.0),
                      child: OutlinedButton(
                        onPressed: () {
                          Navigator.of(context).pop(); // dismiss dialog,
                        },
                        // dismiss dialog,
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
                          "cancel".tr(),
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
      });
}

showUnSubscribePopup(
  context, {
  required void Function() onDone,
}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
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
                width: width / 1.1,
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
                    Image.asset(
                      "assets/images/thinking_man.png",
                      height: height / 4,
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "confirmunsubscribewithwallet"
                              .tr(args: ["Atlantis 1"]),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
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
                          "unsubscribe".tr(),
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10.0),
                      child: OutlinedButton(
                        onPressed: () {
                          Navigator.of(context).pop(); // dismiss dialog,
                        },
                        // dismiss dialog,
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
                          "cancel".tr(),
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
      });
}

showBuyTokenPopup(context,
    {required void Function() onDone,
    required List<DropdownMenuItem<String>> dropdownItems}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
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
                width: width / 1.1,
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
                    Image.asset(
                      "assets/images/thinking_man.png",
                      height: height / 4,
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "addmoretokens".tr(args: ["Atlantis"]),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: dropdown(
                        (value) {},
                        dropdownItems,
                        null,
                        "choosewallet".tr(),
                        context,
                        null,
                      ),
                    ),
                    // SizedBox(
                    //   height: height / 70,
                    // ),
                    // Padding(
                    //   padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    //   child: Row(
                    //     children: [
                    //       Text(
                    //         'Change wallet',
                    //         textAlign: TextAlign.center,
                    //         style: TextStyle(
                    //             color: notifier.getbluewhitecolor,
                    //             fontSize: 12,
                    //             fontFamily: fontsemibold),
                    //       ),
                    //     ],
                    //   ),
                    // ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
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
                          "proceedtobuytokens".tr(),
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
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
                          "cancel".tr(),
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
      });
}

showSwitchEnvironmentPopup(context,
    {required void Function() onProceed,
    required void Function() onCancel,
    required String toEnvironment}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
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
                width: width / 1.1,
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
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "doyouwanttoswitch".tr(args: [toEnvironment]),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontSize: 15,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "appwillrestart".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
                          Navigator.of(context).pop(); // dismiss dialog,
                          onProceed();
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
                          "switchto".tr(args: [toEnvironment]),
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10.0),
                      child: OutlinedButton(
                        onPressed: () {
                          Navigator.of(context).pop(); // dismiss dialog,
                          onCancel();
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
                          "cancel".tr(),
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
      });
}

showSwitchModePopup(context,
    {required void Function() onCreateWallet,
    required void Function() onImportWallet,
    required String toEnvironment}) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
          return AlertDialog(
              // scrollable: true,
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(0),
              content: Container(
                width: width / 1.1,
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
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "welcometothe".tr(args: [toEnvironment]),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontSize: 15,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          "whatwouldyouliketodo".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontbody),
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: ElevatedButton(
                        onPressed: () {
                          Navigator.of(context).pop(); // dismiss dialog,
                          onCreateWallet();
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
                          "createwallet".tr(),
                          style: TextStyle(
                              color: wihitecolor, fontFamily: fontbody),
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10.0),
                      child: OutlinedButton(
                        onPressed: () {
                          Navigator.of(context).pop(); // dismiss dialog,
                          onImportWallet();
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
                          "importwallet".tr(),
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
      });
}

late List<String> walletTypes = [
  'Standard',
  'Minting/Asset Tokenization',
  'Market Making/Trade',
  'Bulk Payment'
];

List<DropdownMenuItem<String>> get walletTypeDropdownItems {
  var dropdownItems = walletTypes
      .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                wallet,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
          value: walletTypes.indexOf(wallet).toString()))
      .toList();

  return dropdownItems;
}

addSubWalletPopup(context) async {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  var appState = Provider.of<DataProvider>(context, listen: false);
  final Authenticator _authenticator = Authenticator();
  late Account primaryWalletKeyPair;
  late Account newSubWalletKeyPair;
  int selectedWalletType = 0;
  String? tag;
  String password = '';
  String? description;
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  final _formKey = GlobalKey<FormState>();
  final _formKey2 = GlobalKey<FormState>();
  WalletAction? action = WalletAction.createNew;
  String? secretKey;
  var userInfo = appState.userInfo!;
  var walletView = WalletView.addSubWallet;

  selectedWalletType = 0;
  tag = '';
  password = '';
  description = '';

  void sendFullDataToServer(responseBody, context) async {
    showLoader(context);
    var appState = Provider.of<DataProvider>(context, listen: false);
    var userInfo = appState.userInfo!;

    try {
      // get primary signature
      var primarySignature = TrovoWalletSDK().signBase64Txn(
        primaryWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      // get secondary signature
      var subWalletSignature = TrovoWalletSDK().signBase64Txn(
        newSubWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      responseBody['primarySignature'] = primarySignature;
      responseBody['subWalletSignature'] = subWalletSignature;

      String requestBody = jsonEncode(responseBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      if (responseData['statusCode'] == 200) {
        // add the secret key of this new subwallet to
        // the existing list of secrets
        appState.secretKeys.add(newSubWalletKeyPair.secretKey);
        // store back the list of secret keys but this time it
        // contains the secret key of the newly created subwallet
        await StoreData().storeInsertData('secretKey', appState.secretKeys);
        await updateUserInfo(
          appState.primaryWallet.signer,
          appState.secretKeys[0],
          appState.primaryWallet.publicKey,
          userInfo.username,
          appState,
          forceRefresh: true,
        );
        // add the new subwallet to appState and
        // set the newly created subwallet as the activeWallet
        appState.activeWallet = appState.userInfo!.wallets!.firstWhere(
            (wallet) => wallet.publicKey == newSubWalletKeyPair.publicKey);
        appState.activeWallet!.secretKey = newSubWalletKeyPair.secretKey;
        // move to next page
        appState.currentAction = PageAction(
            state: PageState.addPage, page: CongratulationsPageConfig);
        Navigator.of(
          context,
          rootNavigator: true,
        ).pop(false);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['error']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }

  postProcessSubwalletData(messageShown, messageLength, data, context) {
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
                postProcessSubwalletData(
                    messageShown, messageLength, data, context),
              });

      messageShown++;
      return;
    }

    sendFullDataToServer(data, context);
    return;
  }

  void sendDataToServer(
    context,
  ) async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "publickey": newSubWalletKeyPair.publicKey,
        "walletTag": tag,
        "WalletDescription": description,
        // "assetIssuerWallet": isAssetIssuerWallet,
        "walletType": selectedWalletType,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      print('response: $responseData');

      if (responseData['statusCode'] == 200) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        await postProcessSubwalletData(
            messageShown, messageLength, responseData['data'], context);
        // print('sending full data to server.........');
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }
    hideLoader(context);
  }

  void handleAuthorization(_formKey, appState, context) {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer(context);
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  void toggleSwitch(context) async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer(context);
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return "pleaseenteryourpassword".tr();

    if (value.length < 6) return "use6charsormoreforpassword".tr();

    return null;
  }

  String? validateDescription(String? value) {
    if (value!.isEmpty) return "enterwalletdesc".tr();

    return null;
  }

  String? validateTag(String? value) {
    if (value!.isEmpty) return "enterwallettag".tr();

    String pattern = r'^[a-zA-Z0-9\_]*$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "invalidtagname".tr();
    }

    return null;
  }

  var show = showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return StatefulBuilder(builder: (context, setStateForDialog) {
          return AlertDialog(
              // scrollable: true,
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(0),
              content: Container(
                width: width / 1.1,
                decoration: BoxDecoration(
                  color: notifier.getwihitecolor,
                  borderRadius: BorderRadius.all(
                    Radius.circular(23),
                  ),
                ),
                child: SingleChildScrollView(
                  child: walletView == WalletView.addSubWallet
                      ? Column(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Padding(
                              padding: const EdgeInsets.all(20.0),
                              child: Center(
                                child: Text(
                                  "addsubwallet".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 20,
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                              ),
                            ),
                            SizedBox(
                              height: height / 50,
                            ),
                            Form(
                              key: _formKey2,
                              child: Column(
                                children: [
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 20.0),
                                    child: Card(
                                      shadowColor: Colors.black,
                                      shape: RoundedRectangleBorder(
                                        borderRadius:
                                            BorderRadius.circular(15.0),
                                      ),
                                      color: notifier.isDark
                                          ? notifier.getbluecolor90
                                          : notifier.getaddsubwalletgrey,
                                      child: Center(
                                        child: Column(
                                          children: [
                                            SizedBox(
                                              height: height / 50,
                                            ),
                                            Container(
                                              width: width / 1.4,
                                              child: Text(
                                                "abouttocreatesubwallet".tr(),
                                                textAlign: TextAlign.center,
                                                style: TextStyle(
                                                  fontSize: 15,
                                                  fontFamily: fontsemibold,
                                                  color: notifier
                                                      .getbluewhitecolor,
                                                ),
                                              ),
                                            ),
                                            SizedBox(
                                              height: height / 50,
                                            ),
                                            Text(
                                              "chooseamethod".tr(),
                                              textAlign: TextAlign.center,
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                              ),
                                            ),
                                            SizedBox(
                                              height: height / 50,
                                            ),
                                            Row(
                                              children: [
                                                SizedBox(
                                                  width: width / 10,
                                                ),
                                                Transform.scale(
                                                  scale: 1.5,
                                                  child: Radio<WalletAction>(
                                                    value: WalletAction.import,
                                                    groupValue: action,
                                                    activeColor: notifier
                                                        .getbluewhitecolor,
                                                    fillColor: MaterialStateColor
                                                        .resolveWith((states) =>
                                                            notifier
                                                                .getbluewhitecolor),
                                                    onChanged: (value) => {
                                                      setStateForDialog(
                                                        () {
                                                          action = value;
                                                        },
                                                      )
                                                    },
                                                  ),
                                                ),
                                                Text(
                                                  "importexistingwallet".tr(),
                                                  style: TextStyle(
                                                    fontSize: 15,
                                                    fontFamily: fontsemibold,
                                                    color: notifier
                                                        .getbluewhitecolor,
                                                  ),
                                                ),
                                              ],
                                            ),
                                            Row(
                                              children: [
                                                SizedBox(
                                                  width: width / 10,
                                                ),
                                                Transform.scale(
                                                  scale: 1.5,
                                                  child: Radio<WalletAction>(
                                                    value:
                                                        WalletAction.createNew,
                                                    activeColor: notifier
                                                        .getbluewhitecolor,
                                                    fillColor: MaterialStateColor
                                                        .resolveWith((states) =>
                                                            notifier
                                                                .getbluewhitecolor),
                                                    groupValue: action,
                                                    onChanged: (value) => {
                                                      setStateForDialog(
                                                        () {
                                                          action = value;
                                                        },
                                                      )
                                                    },
                                                  ),
                                                ),
                                                Text(
                                                  "createnewwallet".tr(),
                                                  style: TextStyle(
                                                    fontSize: 15,
                                                    fontFamily: fontsemibold,
                                                    color: notifier
                                                        .getbluewhitecolor,
                                                  ),
                                                ),
                                              ],
                                            ),
                                            SizedBox(
                                              height: height / 50,
                                            ),
                                            Column(
                                              crossAxisAlignment:
                                                  CrossAxisAlignment.start,
                                              children: [
                                                Row(
                                                  mainAxisAlignment:
                                                      MainAxisAlignment.center,
                                                  children: [
                                                    Container(
                                                      width: width / 1.5,
                                                      decoration: BoxDecoration(
                                                        border: Border.all(
                                                          color: notifier
                                                              .getbluecolor,
                                                        ),
                                                        borderRadius:
                                                            const BorderRadius
                                                                .all(
                                                                Radius.circular(
                                                                    15.0)),
                                                      ),
                                                      child: dropdown(
                                                        (newValue) async {
                                                          selectedWalletType =
                                                              int.parse(newValue
                                                                  .toString());
                                                        },
                                                        walletTypeDropdownItems,
                                                        selectedWalletType
                                                            .toString(),
                                                        null,
                                                        context,
                                                        null,
                                                      ),
                                                    ),
                                                  ],
                                                ),
                                                SizedBox(
                                                  height: height / 50,
                                                ),
                                              ],
                                            )
                                          ],
                                        ),
                                      ),
                                    ),
                                  ),
                                  SizedBox(
                                    height: height / 30,
                                  ),
                                  // Tag name
                                  CustomTextFormField.textField(
                                    "tag".tr(),
                                    notifier.getbluecolor,
                                    Icons.tag,
                                    notifier.getgrey,
                                    notifier.getbluewhitecolor,
                                    notifier.getblck,
                                    notifier.getgrey,
                                    70,
                                    300,
                                    initialValue: tag,
                                    onChanged: (value) {
                                      // setState(() {
                                      //   tag = value.trim().replaceAll(' ', '');
                                      // });
                                    },
                                    onSaved: (value) {
                                      print('tag: $value');
                                      tag = value.trim().replaceAll(' ', '');
                                    },
                                    keyboardtype: TextInputType.text,
                                    maxLength: 12,
                                    validator: validateTag,
                                    helperText: tag == null || tag!.isEmpty
                                        ? ''
                                        : "${appState.userInfo!.username}_$tag",
                                  ),
                                  SizedBox(height: height / 50),
                                  CustomTextFormField.textField(
                                    "description".tr(),
                                    notifier.getbluecolor,
                                    Icons.description,
                                    notifier.getgrey,
                                    notifier.getbluewhitecolor,
                                    notifier.getblck,
                                    notifier.getgrey,
                                    70,
                                    300,
                                    initialValue: description,
                                    onSaved: (value) {
                                      print('description: $value');
                                      description = value;
                                    },
                                    keyboardtype: TextInputType.text,
                                    maxLength: 100,
                                    validator: validateDescription,
                                  ),
                                  if (action == WalletAction.import) ...[
                                    SizedBox(height: height / 50),
                                    // Secret Key
                                    CustomPasswordFormField(
                                      "secretkey".tr(),
                                      notifier.getbluecolor,
                                      Icons.lock,
                                      notifier.getgrey,
                                      notifier.getbluewhitecolor,
                                      notifier.getblck,
                                      70,
                                      300,
                                      validator: (value) {
                                        var trimmedVal =
                                            value!.trim().replaceAll(' ', '');
                                        if (trimmedVal.isEmpty) {
                                          return "entersecretkeyempty".tr();
                                        }

                                        if (trimmedVal.length < 56) {
                                          return "secretkeyinvalid".tr();
                                        }

                                        try {
                                          TrovoWalletSDK()
                                              .parseSecretKey(value);
                                        } catch (e) {
                                          return 'Secret Key is invalid';
                                        }

                                        return null;
                                      },
                                      onSaved: (value) {
                                        print('${"email".tr()}: $value');
                                        secretKey =
                                            value!.trim().replaceAll(' ', '');
                                      },
                                      maxLength: 56,
                                    ),
                                  ],
                                  SizedBox(height: height / 30),
                                  Button(
                                    "continuee".tr(),
                                    notifier.getbluecolor,
                                    wihitecolor,
                                    width: width / 1.5,
                                    onTap: () {
                                      var form = _formKey2.currentState;
                                      if (!form!.validate()) {
                                        return;
                                      }

                                      form.save();

                                      primaryWalletKeyPair = TrovoWalletSDK()
                                          .parseSecretKey(
                                              appState.secretKeys[0]);

                                      setStateForDialog(() {
                                        if (action == WalletAction.import) {
                                          try {
                                            // parse supplied secret to get the keypair
                                            newSubWalletKeyPair =
                                                TrovoWalletSDK()
                                                    .parseSecretKey(secretKey);
                                          } catch (e) {
                                            popup(context,
                                                title: "error".tr(),
                                                message:
                                                    "invalidsecretkey".tr());
                                          }
                                        } else {
                                          // generate keypair for the new subwallet
                                          newSubWalletKeyPair =
                                              TrovoWalletSDK().createAccount();
                                        }
                                        walletView =
                                            WalletView.confirmAddSubWallet;
                                        // appState.walletView.actionIcon =
                                        //     Icons.arrow_circle_left_outlined;
                                        // appState.walletView.actionText = "back".tr();
                                        password = '';
                                        secretKey = '';
                                      });
                                    },
                                  ),
                                  SizedBox(height: height / 70),
                                  TextButton(
                                    child: Text(
                                      "cancel".tr(),
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
                                  SizedBox(height: height / 20),
                                ],
                              ),
                            ),
                          ],
                        )
                      : Column(children: [
                          SizedBox(
                            height: height / 50,
                          ),
                          Padding(
                            padding: const EdgeInsets.all(20.0),
                            child: Center(
                              child: Text(
                                "addsubwallet".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 20,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Padding(
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
                            child: Card(
                              shadowColor: Colors.black,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(15.0),
                              ),
                              color: notifier.isDark
                                  ? notifier.getbluecolor90
                                  : notifier.getaddsubwalletgrey,
                              child: Center(
                                child: Form(
                                  key: _formKey,
                                  child: Column(
                                    children: [
                                      SizedBox(
                                        height: height / 50,
                                      ),
                                      Container(
                                        width: width / 1.4,
                                        child: Text(
                                          "requesttocreatesubwallet".tr(),
                                          textAlign: TextAlign.center,
                                          style: TextStyle(
                                            fontSize: 15,
                                            fontFamily: fontsemibold,
                                            color: notifier.getbluewhitecolor,
                                          ),
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 50,
                                      ),
                                      Text(
                                        "tag".tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontsemibold,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      Text(
                                        "${userInfo.username!}_$tag",
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontbody,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 50,
                                      ),
                                      Text(
                                        "description".tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontsemibold,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      Text(
                                        description ?? '',
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontbody,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 50,
                                      ),
                                      Text(
                                        "method".tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontsemibold,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      Text(
                                        action == WalletAction.import
                                            ? "importsubwallet".tr()
                                            : "createnewsubwallet".tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontbody,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 50,
                                      ),
                                      Text(
                                        "wallettype".tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontsemibold,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      Text(
                                        walletTypes[selectedWalletType],
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontbody,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 50,
                                      ),
                                      Text(
                                        "publickey".tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                          fontSize: 15,
                                          fontFamily: fontsemibold,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                      Padding(
                                        padding: const EdgeInsets.symmetric(
                                            horizontal: 15.0),
                                        child: Text(
                                          newSubWalletKeyPair.publicKey,
                                          textAlign: TextAlign.center,
                                          style: TextStyle(
                                            fontSize: 15,
                                            fontFamily: fontbody,
                                            color: notifier.getbluewhitecolor,
                                          ),
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 20,
                                      ),
                                      Padding(
                                        padding: const EdgeInsets.symmetric(
                                            horizontal: 15.0),
                                        child: Text(
                                          "willattractcharges".tr(),
                                          textAlign: TextAlign.center,
                                          style: TextStyle(
                                            fontSize: 13,
                                            fontFamily: fontbody,
                                            color: notifier.getbluewhitecolor,
                                          ),
                                        ),
                                      ),
                                      SizedBox(
                                        height: height / 30,
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                          ),
                          SizedBox(height: height / 50),
                          // Secret Key
                          CustomPasswordFormField(
                            "password".tr(),
                            notifier.getbluecolor,
                            Icons.lock,
                            notifier.getgrey,
                            notifier.getbluewhitecolor,
                            notifier.getblck,
                            70,
                            300,
                            validator: validatePassword,
                            textInputAction: TextInputAction.done,
                            onChanged: (value) {
                              setStateForDialog(() {
                                password = value!.trim().replaceAll(' ', '');
                              });
                            },
                            // onSubmitted: (value) {
                            //   print('email: $value');
                            //   secretKey = value!.trim().replaceAll(' ', '');
                            // },
                            onSaved: (value) {
                              print('email: $value');
                              secretKey = value!.trim().replaceAll(' ', '');
                            },
                          ),
                          SizedBox(height: height / 50),
                          if (appState.biometricEnabled &&
                              password.isEmpty) ...[
                            Button(
                              "authorizewithbiometrics".tr(),
                              notifier.getbluecolor,
                              wihitecolor,
                              onTap: () => {toggleSwitch(context)},
                              width: width / 1.5,
                            ),
                          ] else ...[
                            Button(
                              "authorize".tr(),
                              notifier.getbluecolor,
                              wihitecolor,
                              onTap: () => {
                                handleAuthorization(_formKey, appState, context)
                              },
                              width: width / 1.5,
                            ),
                          ],
                          SizedBox(height: height / 70),
                          TextButton(
                              child: Text(
                                "back".tr(),
                                style: TextStyle(
                                    fontSize: 14.0,
                                    fontFamily: fontbody,
                                    fontWeight: FontWeight.bold,
                                    color: Colors.teal),
                              ),
                              onPressed: () {
                                print('clicked 0o');
                                setStateForDialog(() {
                                  walletView = WalletView.addSubWallet;
                                });
                              }),
                          SizedBox(height: height / 20),
                        ]),
                ),
              ));
        });
      });

  return show;
}
