import 'dart:async';

import 'package:easy_localization/easy_localization.dart';
import 'package:firebase_remote_config/firebase_remote_config.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:share_plus/share_plus.dart';
// import 'package:flutter_share/flutter_share.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/bottom_tab_page.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:local_auth/error_codes.dart' as auth_error;
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Settings extends StatefulWidget {
  const Settings({Key? key}) : super(key: key);

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> {
  late ColorNotifier notifier;
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

  final internationalLanguages = <String, String>{
    "en": "English",
    "ig": "Igbo",
    "fr": "French",
  };

  List<DropdownMenuItem<String>> get getCurrencies {
    List<DropdownMenuItem<String>> currencies = [];
    appState.fiatRate.forEach((key, value) {
      currencies.add(
        DropdownMenuItem(
          child: Text(key, overflow: TextOverflow.ellipsis),
          value: key,
        ),
      );
    });
    return currencies;
  }

  final _dropDownKey = GlobalKey<FormFieldState>();

  List<DropdownMenuItem<String>> get getLanguages {
    List<DropdownMenuItem<String>> languages = [];
    internationalLanguages.forEach((key, value) {
      languages.add(
        DropdownMenuItem(
          child: Text(
            value,
            textAlign: TextAlign.right,
            overflow: TextOverflow.ellipsis,
          ),
          value: key,
        ),
      );
    });
    return languages;
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
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Stack(
          children: [
            Column(
              children: [
                CustomAppBar(
                  context,
                  notifier.getwihitecolor,
                  '',
                  notifier.getbluewhitecolor,
                  height: height / 15,
                ).getBar(),
                Center(
                  child: CircleAvatar(
                    radius: width / 10,
                    backgroundColor: notifier.getbluecolor70,
                    child: GestureDetector(
                      onTap: () {
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: ProfileDetailsViewPageConfig,
                        );
                      },
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(100.0),
                        child: Image.network(
                          appState.userInfo!.imageThumbnailURL!,
                          width: width / 5.3,
                          // height: width / 10,
                          fit: BoxFit.fill,
                          errorBuilder: (context, error, stackTrace) {
                            return Image.asset(
                              'assets/images/trovo.png',
                              width: width / 9,
                            );
                          },
                        ),
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 80),
                Text(
                  '${appState.userInfo!.firstName} ${appState.userInfo!.lastName} ${appState.userInfo!.isCorporate ? '(Corporate)' : ''}',
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 18,
                  ),
                ),
                TextButton(
                  onPressed: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: ProfileDetailsViewPageConfig,
                    );
                  },
                  style: ButtonStyle(
                    padding: WidgetStateProperty.all(EdgeInsets.zero),
                  ),
                  child: Text(
                    'View Profile',
                    style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                      fontSize: 14,
                      decoration: TextDecoration.underline,
                    ),
                  ),
                ),
                SizedBox(height: height / 50),
                GestureDetector(
                  onTap: () {
                    share();
                  },
                  child: invitefriend(
                    notifier.getbluecolor,
                    "invitefriends".tr(),
                    wihitecolor,
                  ),
                ),
                SizedBox(height: height / 25),
                Row(
                  children: [
                    SizedBox(width: width / 20),
                    Text(
                      "preferences".tr(),
                      style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                // GestureDetector(
                //   child: languages(
                //       "assets/images/languages.png", "", "languages".tr()),
                // ),
                GestureDetector(
                  child: currency(
                    "assets/images/currency.png",
                    "",
                    "currency".tr(),
                  ),
                ),
                GestureDetector(
                  child: darkmode("assets/images/theme.png", "", "theme".tr()),
                ),
                SizedBox(height: height / 25),
                Row(
                  children: [
                    SizedBox(width: width / 20),
                    Text(
                      "security".tr(),
                      style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                GestureDetector(
                  onTap: () => appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: PasswordMgtViewPageConfig,
                  ),
                  child: iteamlist(
                    "assets/images/lock.png",
                    "",
                    "passwordmanagement".tr(),
                  ),
                ),
                GestureDetector(
                  child: hideBalance(
                    "assets/images/eyeoff.png",
                    "",
                    "hidebalance".tr(),
                  ),
                ),
                GestureDetector(
                  child: timeout(
                    "assets/images/hourglass.png",
                    "",
                    "timeout".tr(),
                  ),
                ),
                GestureDetector(
                  child: biometrics(
                    "assets/images/biometrics.png",
                    "",
                    "biometrics".tr(),
                  ),
                ),
                SizedBox(height: height / 25),
                Row(
                  children: [
                    SizedBox(width: width / 20),
                    Text(
                      "wallet".tr(),
                      style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                walletMode(
                  "assets/images/walletmode.png",
                  "",
                  "walletmode".tr(),
                ),
                SizedBox(height: height / 25),
                GestureDetector(
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.replaceAll,
                      page: LoginPageConfig,
                    );
                    appState.isLoggedIn = false;
                  },
                  child: logout("assets/images/logout.png", "", "logout".tr()),
                ),
                SizedBox(height: height / 30),
                Text(
                  '${"version".tr()} ${appState.appVersion}',
                  style: TextStyle(
                    color: notifier.getdarkgrey,
                    fontSize: 13.5,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: height / 30),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Future<void> share() async {
    var label = await FirebaseRemoteConfig.instance.getString(
      'share_wallet_referral_label',
    );

    label = label
        .replaceAll('[link]', appState.userInfo!.referralLink!)
        .replaceAll('[username]', appState.userInfo!.username!);

    var splitLabel = label.split('[newline]');
    var buffer = StringBuffer();

    for (var line in splitLabel) {
      buffer.write('${line}\n\n');
    }

    SharePlus.instance.share(
      ShareParams(text: buffer.toString().trim(), title: 'Trovo Wallet'),
    );
  }

  Widget invitefriend(colorbutton, buttontext, buttontextcolor) {
    return Center(
      child: Container(
        decoration: BoxDecoration(borderRadius: BorderRadius.circular(15)),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            LayoutBuilder(
              builder: (context, constraints) {
                return Container(
                  height: height / 10,
                  width: width / 1.1,
                  decoration: BoxDecoration(
                    color: colorbutton!,
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      Image.asset(
                        "assets/images/referrals.png",
                        height: height / 30,
                      ),
                      Container(
                        width: width / 1.7,
                        child: Text(
                          buttontext!,
                          textAlign: TextAlign.start,
                          style: TextStyle(
                            fontFamily: fontbody,
                            fontSize: 13,
                            color: buttontextcolor,
                          ),
                        ),
                      ),
                      Icon(
                        Icons.arrow_forward_ios,
                        size: 12,
                        color: wihitecolor,
                      ),
                    ],
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget iteamlist(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 14,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Text(
              txt,
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontSize: 13,
                fontFamily: fontsemibold,
                fontWeight: FontWeight.w500,
              ),
            ),
            Icon(Icons.arrow_forward_ios, color: notifier.getgrey, size: 17),
            SizedBox(width: width / 15),
          ],
        ),
      ),
    );
  }

  Widget walletMode(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              width: 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Container(
              width: width / 4.5,
              height: 20,
              child: Row(
                children: [
                  Expanded(
                    child: DropdownButtonFormField(
                      isExpanded: true,
                      key: _dropDownKey,
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      value: appState.walletMode,
                      icon: Visibility(
                        visible: false,
                        child: Icon(Icons.arrow_downward),
                      ),
                      decoration: InputDecoration(
                        contentPadding: EdgeInsets.symmetric(
                          vertical: 0,
                          horizontal: 10,
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(borderSide: BorderSide.none),
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: handleEnvironmentSwitch,
                      items: <DropdownMenuItem<String>>[
                        DropdownMenuItem(
                          child: Text(
                            "testnet".tr(),
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: 'Testnet',
                        ),
                        DropdownMenuItem(
                          child: Text(
                            "mainnet".tr(),
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: 'Mainnet',
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  Widget timeout(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              width: 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Container(
              width: width / 3.5,
              height: 20,
              child: Row(
                children: [
                  Expanded(
                    child: DropdownButtonFormField(
                      isExpanded: true,
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      value: appState.timeout,
                      icon: Visibility(
                        visible: false,
                        child: Icon(Icons.arrow_downward),
                      ),
                      decoration: InputDecoration(
                        contentPadding: EdgeInsets.symmetric(
                          vertical: 0,
                          horizontal: 10,
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(borderSide: BorderSide.none),
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: (newValue) async {
                        await StoreData().storeInsertData('timeOut', newValue);
                        appState.timeout = newValue.toString();
                        setState(() {});
                      },
                      items: <DropdownMenuItem<String>>[
                        DropdownMenuItem(
                          child: Text(
                            '1 ${"minute".tr()}',
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: '1',
                        ),
                        DropdownMenuItem(
                          child: Text(
                            '2 ${"minutes".tr()}',
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: '2',
                        ),
                        DropdownMenuItem(
                          child: Text(
                            '5 ${"minutes".tr()}',
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: '5',
                        ),
                        DropdownMenuItem(
                          child: Text(
                            '10 ${"minutes".tr()}',
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: '10',
                        ),
                        DropdownMenuItem(
                          child: Text(
                            '15 ${"minutes".tr()}',
                            overflow: TextOverflow.ellipsis,
                          ),
                          value: '15',
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  void handleEnvironmentSwitch(String? newValue) async {
    if (newValue != appState.walletMode) {
      showSwitchEnvironmentPopup(
        context,
        onProceed: () async {
          await appState.changeWalletMode(newValue.toString());
        },
        onCancel: () {
          _dropDownKey.currentState!.reset();
        },
        toEnvironment: newValue!,
      );
    }
  }

  Widget currency(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              width: 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Container(
              width: width / 5.9,
              height: 20,
              child: Row(
                children: [
                  Expanded(
                    child: DropdownButtonFormField(
                      isExpanded: true,
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      value: appState.defaultCurrency,
                      icon: Visibility(
                        visible: false,
                        child: Icon(Icons.arrow_downward),
                      ),
                      decoration: InputDecoration(
                        contentPadding: EdgeInsets.symmetric(
                          vertical: 0,
                          horizontal: 10,
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(borderSide: BorderSide.none),
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: (newValue) async {
                        await StoreData().storeInsertData(
                          'defaultCurrency',
                          newValue.toString(),
                        );
                        appState.setDefaultCurrency = newValue.toString();
                      },
                      items: getCurrencies,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  Widget languages(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              width: 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Container(
              width: width / 5,
              height: 20,
              child: Row(
                children: [
                  Expanded(
                    child: DropdownButtonFormField(
                      isExpanded: true,
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      value: appState.defaultLanguage,
                      icon: Visibility(
                        visible: false,
                        child: Icon(Icons.arrow_downward),
                      ),
                      decoration: InputDecoration(
                        contentPadding: EdgeInsets.symmetric(
                          vertical: 0,
                          horizontal: 10,
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(borderSide: BorderSide.none),
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: (newValue) async {
                        context.setLocale(Locale(newValue.toString()));
                        await StoreData().storeInsertData(
                          'defaultLanguage',
                          newValue.toString(),
                        );
                        appState.setDefaultLanguage = newValue.toString();
                      },
                      items: getLanguages,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  Widget logout(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 25),
          Image.asset(
            image,
            height: height / 30,
            color: notifier.getbluewhitecolor,
          ),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
              color: notifier.getblck,
              fontSize: 13,
              fontFamily: fontsemibold,
            ),
          ),
        ],
      ),
    );
  }

  Widget darkmode(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Transform.scale(
              scale: 0.7,
              child: CupertinoSwitch(
                activeColor: notifier.getgreencolor,
                value: notifier.getIsDark,
                onChanged: (val) async {
                  final prefs = await SharedPreferences.getInstance();
                  setState(() {
                    notifier.setIsDark = val;
                    prefs.setBool("setIsDark", val);
                    SystemChrome.setSystemUIOverlayStyle(
                      SystemUiOverlayStyle(
                        statusBarColor: notifier.isDark
                            ? Color(0xFF00225A)
                            : Colors.white,
                        statusBarIconBrightness: notifier.isDark
                            ? Brightness.light
                            : Brightness.dark,
                        systemNavigationBarColor: notifier.isDark
                            ? Color(0xFF00225A)
                            : Colors.white,
                        systemNavigationBarIconBrightness: notifier.isDark
                            ? Brightness.light
                            : Brightness.dark,
                        systemNavigationBarContrastEnforced: true,
                      ),
                    );
                  });
                },
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  Widget biometrics(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Transform.scale(
              scale: 0.7,
              child: CupertinoSwitch(
                activeColor: notifier.getgreencolor,
                value: appState.biometricEnabled,
                onChanged: (val) async {
                  toggleBiometrics();
                },
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  Widget hideBalance(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 15.0),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(
              image,
              height: height / 30,
              color: notifier.getbluewhitecolor,
            ),
            SizedBox(width: width / 40),
            Text(
              name,
              style: TextStyle(
                color: notifier.getblck,
                fontSize: 13,
                fontFamily: fontsemibold,
              ),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Transform.scale(
              scale: 0.7,
              child: CupertinoSwitch(
                activeColor: notifier.getgreencolor,
                value: appState.hideBalances,
                onChanged: (val) async {
                  if (val)
                    toggleHideBalances();
                  else
                    authenticateAndUnhideBalances();
                },
              ),
            ),
            SizedBox(width: width / 20),
          ],
        ),
      ),
    );
  }

  void toggleBiometrics() async {
    try {
      bool result = await Authenticator().authenticateMe();
      if (result) {
        setState(() {
          appState.biometricEnabled = !appState.biometricEnabled;
          StoreData().storeInsertData(
            'biometricsEnabled',
            appState.biometricEnabled,
          );
        });
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        if (!appState.biometricEnabled) {
          popup(
            context,
            title: 'Invalid Operation',
            message:
                'You cannot enable bometrics unless you do biometrics enrollment on your device',
          );
          return;
        }

        biometricsErrorAlert(
          context,
          callback: () {
            setState(() {
              appState.biometricEnabled = false;
              StoreData().storeInsertData(
                'biometricsEnabled',
                appState.biometricEnabled,
              );
            });
          },
        );
      }
    }
  }

  void authenticateAndUnhideBalances() async {
    if (appState.biometricEnabled) {
      try {
        bool result = await Authenticator().authenticateMe();
        if (result) {
          toggleHideBalances();
          return;
        }
      } on PlatformException catch (e) {
        if (e.code == auth_error.notEnrolled ||
            e.code == auth_error.notAvailable) {
          biometricsErrorAlert(context);
        }
      }
    } else {
      showPasswordDialog(context, () {
        toggleHideBalances();
      });
    }
  }

  void toggleHideBalances() {
    setState(() {
      appState.sethideBalances = !appState.hideBalances;
      appState.sethideWalletList = List.filled(6, appState.hideBalances);
      StoreData().storeInsertData('hideBalances', appState.hideBalances);
      StoreData().storeInsertData('hideWalletList', appState.hideWalletList);
      changeTabPage(appState, ButtomTabPage.Dashboard.index);
    });
  }
}
