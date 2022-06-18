import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/screens/ImportWallet/importwallet.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../bottom_bar/bottombar.dart';
import '../../utils/local_auth.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/popups.dart';
import 'create_password.dart';
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
        // appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
        //     height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: height / 7),
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
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
                      SizedBox(height: height / 10),
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
                ],
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  GestureDetector(
                    onTap: () {
                      Get.to(
                        () => ImportWallet(),
                      );
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
              SizedBox(height: height / 20),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  LanguageEn.signinwithbiometrics,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  LanguageEn.signin,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
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
                notifier.getbluecolor,
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const CreatePassword(),
                    ),
                  );
                },
              ),
              SizedBox(height: height / 40),
            ],
          ),
        ),
      ),
    );
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => const BottomHome(),
          ),
        );
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
    print('handling signin $password ${appState.password!}');
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (context) => const BottomHome(),
        ),
      );
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }
}
