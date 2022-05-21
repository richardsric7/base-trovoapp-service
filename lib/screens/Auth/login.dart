import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottombar.dart';
import 'package:gocrypto/screens/Auth/signup.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../reset_password/phone_num_reset_password.dart';

class Login extends StatefulWidget {
  const Login({Key? key}) : super(key: key);

  @override
  State<Login> createState() => _LoginState();
}

class _LoginState extends State<Login> {
  late ColorNotifier notifier;
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
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
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
                            fontFamily: 'Gilroy_Bold'),
                      ),
                      SizedBox(height: height / 40),
                      Text(
                        LanguageEn.youhavebeenmissed,
                        style:
                            TextStyle(fontSize: 16.sp, color: notifier.getgrey),
                      ),
                      SizedBox(height: height / 30),
                      Customtextfild.textField(
                          LanguageEn.emailadress,
                          notifier.getbluecolor,
                          Icons.email,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          45.sp,
                          300.sp),
                      SizedBox(height: height / 30),
                      Custompasswordtextfild.textField(
                          LanguageEn.password,
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck),
                      SizedBox(height: height / 30),
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
                        const PhoneNumResetPassword(),
                      );
                    },
                    child: Text(
                      LanguageEn.forgotpassword,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontSize: 13.5.sp,
                          fontFamily: 'Gilroy_Medium'),
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
                child: Button(LanguageEn.signin, notifier.getbluecolor,
                    notifier.getwihitecolor),
              ),
              SizedBox(height: height / 40),
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
              SizedBox(height: height / 50),
              googlelogin(),
              SizedBox(height: height / 6.5),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.donthaveanaccount,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 15.sp,
                        fontFamily: 'Gilroy_Medium'),
                  ),
                  GestureDetector(
                    onTap: () {
                      Get.to(const SignUp());
                    },
                    child: Text(
                      LanguageEn.signup,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontSize: 15.sp,
                          fontFamily: 'Gilroy_Bold'),
                    ),
                  ),
                ],
              )
            ],
          ),
        ),
      ),
    );
  }

  Widget googlelogin() {
    return Center(
      child: Container(
        height: height / 15,
        width: width / 1.1,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.all(
            Radius.circular(15.sp),
          ),
          border: Border.all(color: notifier.getgrey),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset("assets/images/google.png", height: height / 25),
            SizedBox(width: width / 25),
            Text(
              LanguageEn.continuewithgoogle,
              style: TextStyle(
                  color: notifier.getblck,
                  fontSize: 15.sp,
                  fontFamily: 'Gilroy_Bold'),
            ),
          ],
        ),
      ),
    );
  }
}
