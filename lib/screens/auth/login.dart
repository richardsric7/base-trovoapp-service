import 'dart:io';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:trovo_wallet/Icons/my_flutter_app_icons.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/custtom_textfild/custtom_password.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../utils/local_auth.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/popups.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class Login extends StatefulWidget {
  const Login({Key? key}) : super(key: key);

  @override
  State<Login> createState() => _LoginState();
}

class _LoginState extends State<Login> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  String password = '';
  final _formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  final _dropDownKey = GlobalKey<FormFieldState>();

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
    WidgetsBinding.instance.addPostFrameCallback((_) {
      versionControl();
    });
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    userInfo = appState.userInfo!;
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        resizeToAvoidBottomInset: false,
        floatingActionButton: FloatingActionButton(
          onPressed: () {
            appState.currentAction =
                PageAction(state: PageState.addPage, page: QrScannerPageConfig);
          },
          backgroundColor: notifier.getbluecolor,
          child: SvgPicture.asset(
            "assets/images/scan.svg",
            color: wihitecolor,
          ),
        ),
        body: SingleChildScrollView(
          child: Stack(children: [
            Container(
              height: height / 2.65,
              decoration: BoxDecoration(
                color: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
              ),
              child: Stack(
                children: [
                  Column(
                    children: [
                      SizedBox(
                        height: height / 10.5,
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Image.asset(
                            'assets/images/trovo_white.png',
                            height: height / 4.5,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                        ],
                      ),
                    ],
                  ),
                  Column(
                    children: [
                      SizedBox(
                        height: height / 80,
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Image.asset(
                            'assets/images/trovo_white.png',
                            height: height / 15,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                          Image.asset(
                            'assets/images/trovo_white.png',
                            height: height / 15,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                        ],
                      ),
                      Spacer(
                        flex: 50,
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Image.asset(
                            'assets/images/trovo_white.png',
                            height: height / 15,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                          Image.asset(
                            'assets/images/trovo_white.png',
                            height: height / 15,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                        ],
                      ),
                      SizedBox(
                        height: height / 80,
                      ),
                    ],
                  ),
                ],
              ),
            ),
            Column(
              children: [
                SizedBox(height: height / 20),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      Container(
                        width: width / 2.8,
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
                                icon: Icon(
                                  Icons.keyboard_arrow_down_rounded,
                                ),
                                decoration: InputDecoration(
                                  contentPadding: EdgeInsets.symmetric(
                                      vertical: 8.0, horizontal: 10),
                                  enabledBorder: OutlineInputBorder(
                                    borderSide: BorderSide.lerp(
                                        BorderSide(color: notifier.getgrey),
                                        BorderSide(color: notifier.getgrey),
                                        1.0),
                                    borderRadius: const BorderRadius.all(
                                        Radius.circular(20.0)),
                                  ),
                                  border: OutlineInputBorder(
                                    borderSide: BorderSide.lerp(
                                        BorderSide(color: notifier.getgrey),
                                        BorderSide(color: notifier.getgrey),
                                        1.0),
                                    borderRadius: const BorderRadius.all(
                                        Radius.circular(20.0)),
                                  ),
                                ),
                                elevation: 0,
                                style: TextStyle(
                                    color: notifier.getdarkgrey,
                                    fontSize: 13.5.sp,
                                    fontFamily: fontbody),
                                onChanged: handleEnvironmentSwitch,
                                items: <DropdownMenuItem<String>>[
                                  DropdownMenuItem(
                                    child: Row(
                                      children: [
                                        Icon(
                                          CustomIcon.globeOutlined,
                                          size: 18,
                                        ),
                                        SizedBox(
                                          width: 6,
                                        ),
                                        Text(
                                          "testnet".tr(),
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ],
                                    ),
                                    value: 'Testnet',
                                  ),
                                  DropdownMenuItem(
                                    child: Row(
                                      children: [
                                        Icon(
                                          CustomIcon.globeOutlined,
                                          size: 18,
                                        ),
                                        SizedBox(
                                          width: 6,
                                        ),
                                        Text(
                                          "mainnet".tr(),
                                          overflow: TextOverflow.ellipsis,
                                        ),
                                      ],
                                    ),
                                    value: 'Mainnet',
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                SizedBox(height: height / 15),
                Column(
                  children: [
                    Column(
                      children: [
                        Text(
                          "welcome".tr(),
                          style: TextStyle(
                              color: notifier.getblck,
                              fontSize: 26.sp,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(height: height / 95),
                        ConstrainedBox(
                          constraints: BoxConstraints(maxWidth: width / 1.1),
                          child: Text(
                            userInfo.username ?? "",
                            style: TextStyle(
                                color: notifier.getblck,
                                fontSize: 26.sp,
                                fontFamily: fontsemibold),
                          ),
                        ),
                        SizedBox(height: height / 40),
                        Text(
                          "youhavebeenmissed".tr(),
                          style: TextStyle(
                              fontSize: 16.sp,
                              color: notifier.getgrey,
                              fontFamily: fontbody),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 7),
                    Form(
                      key: _formKey,
                      child: CustomPasswordFormField(
                        "password".tr(),
                        notifier.getbluecolor,
                        Icons.lock,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        70.sp,
                        300.sp,
                        validator: validateInput,
                        onChanged: (value) {
                          setState(() {
                            password = value!.trim().replaceAll(' ', '');
                          });
                        },
                      ),
                    ),
                  ],
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 25.0),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      TextButton(
                        onPressed: () {
                          appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: ImportWalletPageConfig);
                        },
                        child: Text(
                          "forgotpassword".tr(),
                          style: TextStyle(
                              color: notifier.getdarkgrey,
                              fontSize: 13.5.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                      // SizedBox(width: width / 10),
                    ],
                  ),
                ),
                SizedBox(height: height / 25),
                if (appState.biometricEnabled && password.isEmpty) ...[
                  Button(
                    "signinwithbiometrics".tr(),
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: toggleSwitch,
                  ),
                ] else ...[
                  Button(
                    "signin".tr(),
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: handleSignin,
                  ),
                ],

                // SizedBox(height: height / 90),
                Row(
                  children: [
                    Expanded(
                      child: Container(
                        margin: const EdgeInsets.only(left: 27.0, right: 27.0),
                        child: Divider(
                          color: notifier.getgrey,
                          height: 50,
                        ),
                      ),
                    ),
                    Text(
                      "oR".tr(),
                      style: TextStyle(color: notifier.getgrey),
                    ),
                    Expanded(
                      child: Container(
                        margin: const EdgeInsets.only(left: 27.0, right: 27.0),
                        child: Divider(
                          color: notifier.getgrey,
                          height: 50,
                        ),
                      ),
                    ),
                  ],
                ),
                ButtonOutlined(
                  "signup".tr(),
                  notifier.getwihitecolor,
                  notifier.getbluewhitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: CreatePasswordPageConfig);
                  },
                ),
                SizedBox(height: height / 50),
                ButtonOutlined(
                  "recoveraccount".tr(),
                  notifier.getbluecolor80,
                  wihitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: RecoverAccountViewPageConfig);
                  },
                ),
                SizedBox(height: height / 50),
                Text(
                  '${"version".tr()} ${appState.appVersion}',
                  style: TextStyle(
                      color: notifier.getdarkgrey,
                      fontSize: 13.5.sp,
                      fontFamily: fontbody),
                ),
                SizedBox(height: height / 20),
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          ]),
        ),
      ),
    );
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        appState.currentAction =
            PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
        appState.isLoggedIn = true;
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  String? validateInput(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  void handleSignin() {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      appState.currentAction =
          PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
      appState.isLoggedIn = true;
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  versionControl() async {
    Map appVersionData = await StoreData().storeGetData('appVersion') ?? {};

    if (appVersionData.isNotEmpty) {
      var phoneVersion = appState.appVersion.replaceAll('.', '');

      String minVersion =
          await appVersionData['minVersion'].replaceAll('.', '');

      String currentVersion =
          await appVersionData['version'].replaceAll('.', '');

      var currentVersionRaw = await appVersionData['version'];

      var forceUpdate = await appVersionData['forceUpdate'];

      // the min version is ahead of current phone version. Force update
      if (num.tryParse(minVersion)! > num.tryParse(phoneVersion)!) {
        updateAppMessagePopup(context,
            'You must upgrade to Trovo Wallet version $currentVersionRaw to continue to use the wallet.',
            () async {
          Uri uri = Uri.parse(
            Platform.isAndroid
                ? appVersionData['androidUrl'].toString()
                : appVersionData['iosUrl'].toString(),
          );
          if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
            throw 'Could not launch $uri';
          }
        });
        //You must upgrade to BantuPay Version $currentVersionRaw to continue to use the wallet.');
        return;
      }

      // new version available but doesn't require force update

      if (num.tryParse(currentVersion)! > num.tryParse(phoneVersion)!) {
        // the current version requires force update.
        if (forceUpdate == 1) {
          updateAppMessagePopup(context,
              'You must upgrade to Trovo Wallet version: $currentVersionRaw to continue using the wallet.',
              () async {
            Uri uri = Uri.parse(
              Platform.isAndroid
                  ? appVersionData['androidUrl'].toString()
                  : appVersionData['iosUrl'].toString(),
            );
            if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
              throw 'Could not launch $uri';
            }
          });
          return;
        }

        showSnackBarForInfo(
            'Trovo Wallet version $currentVersionRaw available for download',
            context);
      }
    }
  }

  void handleEnvironmentSwitch(String? newValue) async {
    if (newValue != appState.walletMode) {
      showSwitchEnvironmentPopup(context, onProceed: () async {
        await appState.changeWalletMode(newValue.toString());
      }, onCancel: () {
        _dropDownKey.currentState!.reset();
      }, toEnvironment: newValue!);
    }
  }
}
