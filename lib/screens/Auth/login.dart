import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../router/PageActions.dart';
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
                            'assets/images/trovo-logo-bg.png',
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
                            'assets/images/trovo-logo-bg.png',
                            height: height / 15,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                          Image.asset(
                            'assets/images/trovo-logo-bg.png',
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
                            'assets/images/trovo-logo-bg.png',
                            height: height / 15,
                            color: notifier.isDark
                                ? notifier.getdarkgrey
                                : notifier.getsplashgrey,
                          ),
                          Image.asset(
                            'assets/images/trovo-logo-bg.png',
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
                SizedBox(height: height / 7),
                Column(
                  children: [
                    Column(
                      children: [
                        Text(
                          LanguageEn.welcome,
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
                          LanguageEn.youhavebeenmissed,
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
                        LanguageEn.password,
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
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    TextButton(
                      onPressed: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: ImportWalletPageConfig);
                      },
                      child: Text(
                        LanguageEn.forgotpassword,
                        style: TextStyle(
                            color: notifier.getdarkgrey,
                            fontSize: 13.5.sp,
                            fontFamily: fontbody),
                      ),
                    ),
                    SizedBox(width: width / 10),
                  ],
                ),
                SizedBox(height: height / 25),
                if (appState.biometricEnabled && password.isEmpty) ...[
                  Button(
                    LanguageEn.signinwithbiometrics,
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: toggleSwitch,
                  ),
                ] else ...[
                  Button(
                    LanguageEn.signin,
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: handleSignin,
                  ),
                ],

                // SizedBox(height: height / 90),
                Row(
                  children: <Widget>[
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
                      LanguageEn.oR,
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
                  LanguageEn.signup,
                  notifier.getwihitecolor,
                  notifier.getbluewhitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: CreatePasswordPageConfig);
                  },
                ),
                SizedBox(height: height / 50),
                Text(
                  '${LanguageEn.version} $appVersion',
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
    print('handling signin...');
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      appState.currentAction =
          PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }
}
