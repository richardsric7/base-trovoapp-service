import 'package:firebase_remote_config/firebase_remote_config.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/BottomTabPage.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:local_auth/error_codes.dart' as auth_error;
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

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

  List<DropdownMenuItem<String>> get getCurrencies {
    List<DropdownMenuItem<String>> currencies = [];
    appState.fiatRate.forEach((key, value) {
      currencies.add(DropdownMenuItem(
          child: Text(
            key,
            overflow: TextOverflow.ellipsis,
          ),
          value: key));
    });
    return currencies;
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
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 10,
              ),
              Center(
                child: CircleAvatar(
                    radius: width / 10,
                    backgroundColor: notifier.getbluecolor70,
                    child: GestureDetector(
                      onTap: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: ProfileDetailsViewPageConfig);
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
                    )),
              ),
              SizedBox(
                height: height / 80,
              ),
              Text(
                '${appState.userInfo!.firstName} ${appState.userInfo!.lastName}',
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 18.sp),
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  share();
                },
                child: invitefriend(notifier.getbluecolor,
                    LanguageEn.invitefriends, wihitecolor),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.personal,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: ProfileDetailsViewPageConfig);
                },
                child: iteamlist(
                    "assets/images/profile.png", "", LanguageEn.myprofile),
              ),
              // GestureDetector(
              //   onTap: () {
              //     appState.currentAction = PageAction(
              //         state: PageState.addPage,
              //         page: ReferralInfoViewPageConfig);
              //   },
              //   child: iteamlist("assets/images/referrals-dark.png", "",
              //       LanguageEn.myreferrals),
              // ),
              // GestureDetector(
              //   child: iteamlist(
              //       "assets/images/trovo-blue.png", "", LanguageEn.trovopatron),
              // ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.preferences,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                child: iteamlist(
                    "assets/images/languages.png", "", LanguageEn.languages),
              ),
              GestureDetector(
                child: currency(
                    "assets/images/currency.png", "", LanguageEn.currency),
              ),
              GestureDetector(
                child:
                    darkmode("assets/images/theme.png", "", LanguageEn.theme),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.security,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  // if shared access is enabled on this user's account
                  if (appState.introducedSharedAccess) {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: SharedAccessViewPageConfig);
                  } else {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: WelcomeToSharedAccessViewPageConfig);
                  }
                },
                child: iteamlist(
                    "assets/images/access.png", "", LanguageEn.sharedaccess),
              ),
              GestureDetector(
                onTap: () => appState.currentAction = PageAction(
                    state: PageState.addPage, page: PasswordMgtViewPageConfig),
                child: iteamlist("assets/images/lock.png", "",
                    LanguageEn.passwordmanagement),
              ),
              GestureDetector(
                child: hideBalance(
                    "assets/images/eyeoff.png", "", LanguageEn.hidebalance),
              ),
              GestureDetector(
                child: timeout(
                    "assets/images/hourglass.png", "", LanguageEn.timeout),
              ),
              GestureDetector(
                child: biometrics(
                    "assets/images/biometrics.png", "", LanguageEn.biometrics),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.wallet,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () => appState.currentAction = PageAction(
                    state: PageState.addPage, page: ImportWalletPageConfig),
                child: iteamlist(
                    "assets/images/import.png", "", LanguageEn.importwallet),
              ),
              GestureDetector(
                onTap: () {
                  // go to the definition of appState.viewData
                  // to learn more about viewData
                  appState.viewData![EnsurePrivacyPageConfig.key] = {
                    'rel': 'backupAll',
                  };
                  appState.currentAction = PageAction(
                      state: PageState.addPage, page: EnsurePrivacyPageConfig);
                },
                child: iteamlist("assets/images/backup-wallets.png", "",
                    LanguageEn.backupwallet),
              ),
              walletMode(
                  "assets/images/walletmode.png", "", LanguageEn.walletmode),
              GestureDetector(
                onTap: () {
                  if (appState.userInfo!.hasSecurityQuestions == 0) {
                    var primaryWallet = appState.userInfo!.wallets!
                        .firstWhere((wallet) => wallet.primaryWallet == 1);
                    appState.viewData = {
                      SecurityQuestionsViewPageConfig.key: {
                        'signer': primaryWallet.signer,
                        'publicKey': primaryWallet.publicKey,
                        'secretKey': appState.secretKeys[0],
                        'username': appState.userInfo!.username,
                      }
                    };
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: SecurityQuestionsViewPageConfig);
                  } else if (appState.userInfo!.accountRecoveryEnabled == 0) {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: SetupAccountRecoveryViewPageConfig);
                  } else {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: DisableAccountRecoveryInfoViewPageConfig);
                  }
                },
                child: iteamlist("assets/images/history.png", "",
                    LanguageEn.accountrecovery),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.more,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                child: iteamlist(
                    "assets/images/help.png", "", LanguageEn.helpandsupport),
              ),
              GestureDetector(
                onTap: () => appState.goToWebView(termsOfServiceUrl),
                child: iteamlist(
                    "assets/images/terms.png", "", LanguageEn.termsofuse),
              ),
              GestureDetector(
                onTap: () => appState.goToWebView(trovoLandingPage),
                child: iteamlist("assets/images/copyright.png", "",
                    LanguageEn.abouttrovowallet),
              ),
              SizedBox(height: height / 20),
              GestureDetector(
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: LoginPageConfig);
                  appState.isLoggedIn = false;
                },
                child:
                    logout("assets/images/logout.png", "", LanguageEn.logout),
              ),
              SizedBox(height: height / 30),
              Text(
                '${LanguageEn.version} $appVersion',
                style: TextStyle(
                    color: notifier.getdarkgrey,
                    fontSize: 13.5.sp,
                    fontFamily: fontbody),
              ),
              SizedBox(height: height / 30),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> share() async {
    var label = await FirebaseRemoteConfig.instance
        .getString('share_wallet_referral_label');

    label = label
        .replaceAll('[link]', appState.userInfo!.referralLink!)
        .replaceAll('[username]', appState.userInfo!.username!);

    var splitLabel = label.split('[newline]');
    var buffer = StringBuffer();

    for (var line in splitLabel) {
      buffer.write('${line}\n\n');
    }

    await FlutterShare.share(
      title: 'Trovo Wallet',
      text: buffer.toString().trim(),
    );
  }

  Widget invitefriend(colorbutton, buttontext, buttontextcolor) {
    return Center(
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(15),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            LayoutBuilder(builder: (context, constraints) {
              return ScreenUtilInit(
                builder: (context, child) => Container(
                  height: height / 10,
                  width: width / 1.1,
                  decoration: BoxDecoration(
                    color: colorbutton!,
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      Image.asset("assets/images/referrals.png",
                          height: height / 30),
                      Container(
                        width: width / 1.7,
                        child: Text(
                          buttontext!,
                          textAlign: TextAlign.start,
                          style: TextStyle(
                              fontFamily: fontbody,
                              fontSize: 13.sp,
                              color: buttontextcolor),
                        ),
                      ),
                      Icon(
                        Icons.arrow_forward_ios,
                        size: 12.sp,
                        color: wihitecolor,
                      )
                    ],
                  ),
                ),
              );
            }),
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
                  fontSize: 14.sp,
                  fontFamily: fontsemibold),
            ),
            const Spacer(),
            SizedBox(width: width / 100),
            Icon(Icons.arrow_forward_ios, color: notifier.getgrey, size: 17.sp),
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
                  fontSize: 13.sp,
                  fontFamily: fontsemibold),
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
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      value: 'Testnet',
                      icon: Visibility(
                          visible: false, child: Icon(Icons.arrow_downward)),
                      decoration: InputDecoration(
                        contentPadding:
                            EdgeInsets.symmetric(vertical: 0, horizontal: 10),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      onChanged: (newValue) {
                        setState(() {});
                      },
                      items: <DropdownMenuItem<String>>[
                        DropdownMenuItem(
                            child: Text(
                              'Testnet',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: 'Testnet'),
                        DropdownMenuItem(
                            child: Text(
                              'Mainnet',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: 'Mainet'),
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
                  fontSize: 13.sp,
                  fontFamily: fontsemibold),
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
                          visible: false, child: Icon(Icons.arrow_downward)),
                      decoration: InputDecoration(
                        contentPadding:
                            EdgeInsets.symmetric(vertical: 0, horizontal: 10),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
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
                        setState(() {});
                      },
                      items: <DropdownMenuItem<String>>[
                        DropdownMenuItem(
                            child: Text(
                              '1 ${LanguageEn.minute}',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: '1'),
                        DropdownMenuItem(
                            child: Text(
                              '2 ${LanguageEn.minutes}',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: '2'),
                        DropdownMenuItem(
                            child: Text(
                              '5 ${LanguageEn.minutes}',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: '5'),
                        DropdownMenuItem(
                            child: Text(
                              '10 ${LanguageEn.minutes}',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: '10'),
                        DropdownMenuItem(
                            child: Text(
                              '15 ${LanguageEn.minutes}',
                              overflow: TextOverflow.ellipsis,
                            ),
                            value: '15'),
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
                  fontSize: 13.sp,
                  fontFamily: fontsemibold),
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
                          visible: false, child: Icon(Icons.arrow_downward)),
                      decoration: InputDecoration(
                        contentPadding:
                            EdgeInsets.symmetric(vertical: 0, horizontal: 10),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                        ),
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
                            'defaultCurrency', newValue.toString());
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
                fontSize: 13.sp,
                fontFamily: fontsemibold),
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
                  fontSize: 13.sp,
                  fontFamily: fontsemibold),
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
                  fontSize: 13.sp,
                  fontFamily: fontsemibold),
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
                  fontSize: 13.sp,
                  fontFamily: fontsemibold),
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
          StoreData()
              .storeInsertData('biometricsEnabled', appState.biometricEnabled);
        });
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
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
    }

    showPasswordDialog(context, () {
      toggleHideBalances();
    });
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
