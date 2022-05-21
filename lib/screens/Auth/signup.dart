import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottombar.dart';
import 'package:gocrypto/screens/Auth/login.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SignUp extends StatefulWidget {
  const SignUp({Key? key}) : super(key: key);

  @override
  State<SignUp> createState() => _SignUpState();
}

class _SignUpState extends State<SignUp> {
  late ColorNotifier notifier;
  bool isChecked = false;

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
                        LanguageEn.signup,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 26.sp,
                            fontFamily: 'Gilroy_Bold'),
                      ),
                      SizedBox(height: height / 35),
                      Text(
                        LanguageEn.ittakesaminute,
                        style:
                            TextStyle(fontSize: 14.sp, color: notifier.getgrey),
                      ),
                      SizedBox(height: height / 30),
                      Row(
                        children: [
                          Customtextfild.textField(
                              LanguageEn.fanme,
                              notifier.getbluecolor,
                              Icons.person,
                              notifier.getgrey,
                              notifier.getprefixicon,
                              notifier.getblck,
                              notifier.getgrey,
                              45.sp,
                              145.sp),
                          SizedBox(width: width / 40),
                          Customtextfild.textField(
                              LanguageEn.lname,
                              notifier.getbluecolor,
                              Icons.person,
                              notifier.getgrey,
                              notifier.getprefixicon,
                              notifier.getblck,
                              notifier.getgrey,
                              45.sp,
                              145.sp),
                        ],
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
                      Row(
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
                              side: BorderSide(color: notifier.getbluecolor),
                              value: isChecked,
                              onChanged: (bool? value) {
                                setState(() {
                                  isChecked = value!;
                                });
                              },
                            ),
                          ),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Text(
                                    LanguageEn.iagreetothe,
                                    style: TextStyle(
                                        fontSize: height / 55,
                                        color: notifier.getblck),
                                  ),
                                  Text(
                                    LanguageEn.termsofservices,
                                    style: TextStyle(
                                        fontFamily: 'Gilroy_Medium',
                                        fontSize: height / 55,
                                        color: notifier.getbluecolor),
                                  ),
                                  Text(
                                    LanguageEn.and,
                                    style: TextStyle(
                                        fontFamily: 'Gilroy_Medium',
                                        fontSize: height / 55,
                                        color: notifier.getblck),
                                  ),
                                ],
                              ),
                              Text(
                                LanguageEn.privacypolicy,
                                style: TextStyle(
                                    fontFamily: 'Gilroy_Medium',
                                    fontSize: height / 55,
                                    color: notifier.getbluecolor),
                              ),
                            ],
                          )
                        ],
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 25),
              GestureDetector(
                  onTap: () {
                    Get.to(const BottomHome());
                  },
                  child: Button(LanguageEn.signup, notifier.getbluecolor,
                      notifier.getwihitecolor)),
              SizedBox(height: height / 50),

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
              SizedBox(height: height / 80),
              googlelogin(),
              SizedBox(height: height / 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.alreadyregistered,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Medium'),
                  ),
                  GestureDetector(
                    onTap: () {
                      Get.to(const Login());
                    },
                    child: Text(
                      LanguageEn.signin,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontSize: 13.sp,
                          fontFamily: 'Gilroy_Medium'),
                    ),
                  )
                ],
              )
              // Row(
              //   mainAxisAlignment: MainAxisAlignment.center,
              //   children: [
              //     // Text(
              //     //   "Don’t have an account?",
              //     //   style: TextStyle(
              //     //       color: notifier.getgrey,
              //     //       fontSize: 15.sp,
              //     //       fontFamily: 'Gilroy_Medium'),
              //     // ),
              //     // GestureDetector(
              //     //   onTap: () {
              //     //     Get.to(const BottomHome());
              //     //   },
              //     //   child: Text(
              //     //     "Login",
              //     //     style: TextStyle(
              //     //         color: notifier.getbluecolor,
              //     //         fontSize: 15.sp,
              //     //         fontFamily: 'Gilroy_Bold'),
              //     //   ),
              //     // ),
              //   ],
              // )
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
