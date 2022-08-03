import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/Auth/login.dart';
import 'package:trovo_wallet/screens/Payment%20Method/paymentmethod.dart';
import 'package:trovo_wallet/screens/profile/faq.dart';
import 'package:trovo_wallet/screens/profile/language.dart';
import 'package:trovo_wallet/screens/profile/myaccount.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Profile extends StatefulWidget {
  const Profile({Key? key}) : super(key: key);

  @override
  State<Profile> createState() => _ProfileState();
}

class _ProfileState extends State<Profile> {
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
    var appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 15,
              ),
              Center(
                child: Image.asset("assets/images/avatar.png",
                    height: height / 10),
              ),
              SizedBox(height: height / 70),
              Text(
                LanguageEn.hirenjoshi,
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: 'Gilroy_Bold',
                    fontSize: 16.sp),
              ),
              // _profile(),
              SizedBox(height: height / 20),
              GestureDetector(
                onTap: () {
                  share();
                },
                child: invitefriend(
                    notifier.getbluecolor,
                    "Invite Friends\nInvite your friends and get \$20 each.",
                    notifier.getwihitecolor),
              ),
              // GestureDetector(
              //     onTap: () {
              //       Get.to(() => const ReferralCode());
              //     },
              //     child: referalcode()),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.general,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  Get.to(() => const MyAccount());
                },
                child: iteamlist("assets/images/BillingPayment.png", "",
                    LanguageEn.myaccount),
              ),
              SizedBox(height: height / 30),
              GestureDetector(
                onTap: () {
                  Get.to(() => const PaymentMethod());
                },
                child: iteamlist("assets/images/Language.png", "",
                    LanguageEn.billingpayment),
              ),
              SizedBox(height: height / 30),
              GestureDetector(
                onTap: () {
                  Get.to(() => const FAQ());
                },
                child: iteamlist(
                    "assets/images/Settings.png", "", LanguageEn.faqsupport),
              ),

              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.general,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child:
                    iteamlist("assets/images/FAQ.png", "", LanguageEn.language),
              ),
              SizedBox(height: height / 40),
              GestureDetector(
                onTap: () {},
                child: darkmode(
                    "assets/images/darkmode.png", "", LanguageEn.darkmode),
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: LoginPageConfig);
                },
                child:
                    iteamlist("assets/images/FAQ.png", "", LanguageEn.logout),
              ),
              SizedBox(height: height / 50),
              // GestureDetector(
              //   onTap: () {
              //     Get.to(() => const MessageSupport());
              //   },
              //   child: support(
              //       notifier.getbluecolor,
              //       "We’d love to hear your feedback!\nWe are always looking to improve.",
              //       notifier.getwihitecolor),
              // )
            ],
          ),
        ),
      ),
    );
  }

  Future<void> share() async {
    await FlutterShare.share(
        title: 'Example share',
        text: 'Example share text',
        linkUrl: 'https://flutter.dev/',
        chooserTitle: 'Example Chooser Title');
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
                    Image.asset("assets/images/Invitefriends.png",
                        height: height / 20),
                    Text(
                      buttontext!,
                      textAlign: TextAlign.start,
                      style: TextStyle(
                          fontFamily: 'Gilroy_Medium',
                          fontSize: 13.sp,
                          color: buttontextcolor),
                    ),
                    Icon(
                      Icons.arrow_forward_ios,
                      size: 12.sp,
                      color: notifier.getwihitecolor,
                    )
                  ],
                ),
              );
            }),
          ],
        ),
      ),
    );
  }

  // Widget _profile() {
  //   return Row(
  //     children: [
  //       SizedBox(width: width / 20),
  //       Stack(
  //         children: [
  //           Image.asset("assets/images/profile.png", height: height / 14),
  //           Padding(
  //             padding: EdgeInsets.only(top: height / 26, left: width / 12),
  //             child:
  //                 Image.asset("assets/images/Camera.png", height: height / 30),
  //           ),
  //         ],
  //       ),
  //       SizedBox(width: width / 20),
  //       Column(
  //         crossAxisAlignment: CrossAxisAlignment.start,
  //         children: [
  //           Text(
  //             "Sunder Pichai",
  //             style: TextStyle(
  //                 color: notifier.getblck,
  //                 fontSize: 15.sp,
  //                 fontFamily: 'Gilroy_Bold'),
  //           ),
  //           Text(
  //             "SunderPichai@yahoo.com",
  //             style: TextStyle(
  //                 color: notifier.getgrey,
  //                 fontSize: 12.sp,
  //                 fontFamily: 'Gilroy_Medium'),
  //           ),
  //         ],
  //       ),
  //       const Spacer(),
  //       GestureDetector(
  //         onTap: () {
  //           Get.to(() => const MyAccount());
  //         },
  //         child: Text(
  //           "Edit",
  //           style: TextStyle(
  //               color: notifier.getbluecolor,
  //               fontFamily: 'Gilroy_Bold',
  //               fontSize: 15.sp),
  //         ),
  //       ),
  //       SizedBox(width: width / 17),
  //     ],
  //   );
  // }

  // Widget referalcode() {
  //   return Container(
  //     color: Colors.transparent,
  //     child: Row(
  //       children: [
  //         SizedBox(width: width / 20),
  //         Column(
  //           crossAxisAlignment: CrossAxisAlignment.start,
  //           children: [
  //             Text(
  //               "Referral Code",
  //               style: TextStyle(
  //                   color: notifier.getblck,
  //                   fontSize: 15.sp,
  //                   fontFamily: 'Gilroy_Bold'),
  //             ),
  //             Text(
  //               "Share your love and get \$10 of free stocks",
  //               style: TextStyle(
  //                   color: notifier.getgrey,
  //                   fontSize: 12.sp,
  //                   fontFamily: 'Gilroy_Medium'),
  //             ),
  //           ],
  //         ),
  //         const Spacer(),
  //         Image.asset("assets/images/gift.png", height: height / 30),
  //         SizedBox(width: width / 20),
  //       ],
  //     ),
  //   );
  // }

  Widget iteamlist(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 25),
          Image.asset(image, height: height / 22),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
          ),
          const Spacer(),
          SizedBox(width: width / 100),
          Icon(Icons.arrow_forward_ios, color: notifier.getgrey, size: 17.sp),
          SizedBox(width: width / 15),
        ],
      ),
    );
  }

  Widget darkmode(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 25),
          Image.asset(image, height: height / 22),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
          ),
          const Spacer(),
          SizedBox(width: width / 100),
          Transform.scale(
            scale: 0.7,
            child: CupertinoSwitch(
              activeColor: notifier.getbluecolor,
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
          SizedBox(width: width / 15),
        ],
      ),
    );
  }
}
