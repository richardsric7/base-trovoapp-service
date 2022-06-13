import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/bottom_bar/bottombar.dart';
import 'package:trovo_wallet/screens/ImportWallet/importwallet.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../network/requests.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../reset_password/emailpassword.dart';
import 'create_password.dart';

class Login extends StatefulWidget {
  const Login({Key? key}) : super(key: key);

  @override
  State<Login> createState() => _LoginState();
}

class _LoginState extends State<Login> {
  late ColorNotifier notifier;
  DataProvider? appState;
  UserInfo? userInfo;
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
    userInfo = appState!.userInfo;
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
                          userInfo?.username ?? "",
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
                      CustomPasswordFormField(
                        LanguageEn.password,
                        notifier.getbluecolor,
                        Icons.lock,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        70.sp,
                        300.sp,
                        onChanged: (value) {
                          setState(() {
                            // password = value!.trim().replaceAll(' ', '');
                          });
                        },
                        // validator: validatePassword,
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
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const BottomHome(),
                    ),
                  );
                },
                child: Button(LanguageEn.signinwithbiometrics,
                    notifier.getbluecolor, notifier.getwihitecolor),
              ),
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
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const CreatePassword(),
                    ),
                  );
                },
                child: button(LanguageEn.signup, notifier.getwihitecolor,
                    notifier.getbluecolor),
              ),
              SizedBox(height: height / 40),
            ],
          ),
        ),
      ),
    );
  }

  Widget button(buttontext, colorbutton, buttontextcolor) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(15),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          LayoutBuilder(builder: (context, constraints) {
            return Container(
              height: height / 15,
              width: width / 1.1,
              decoration: BoxDecoration(
                border: Border.all(color: notifier.getgrey),
                color: colorbutton,
                borderRadius: BorderRadius.circular(15),
              ),
              child: Center(
                child: Text(
                  buttontext,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontFamily: fontbody,
                      fontSize: 15.sp,
                      color: buttontextcolor),
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}
