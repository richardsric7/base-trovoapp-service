import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/Payment%20Method/paymentmethod.dart';
import 'package:trovo_wallet/screens/profile/faq.dart';
import 'package:trovo_wallet/screens/profile/language.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Settings extends StatefulWidget {
  const Settings({Key? key}) : super(key: key);

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> {
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
              SizedBox(height: height / 20),
              GestureDetector(
                onTap: () {
                  share();
                },
                child: invitefriend(notifier.getbluecolor,
                    LanguageEn.invitefriends, notifier.getwihitecolor),
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
                        fontFamily: 'Gilroy_Bold'),
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
              GestureDetector(
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: ReferralInfoViewPageConfig);
                },
                child: iteamlist("assets/images/referrals-dark.png", "",
                    LanguageEn.myreferrals),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const PaymentMethod());
                },
                child: iteamlist("assets/images/trovo-blue.png", "",
                    LanguageEn.mysubscriptions),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.preferences,
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
                  Get.to(() => const FAQ());
                },
                child: iteamlist(
                    "assets/images/languages.png", "", LanguageEn.languages),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const FAQ());
                },
                child: iteamlist(
                    "assets/images/currency.png", "", LanguageEn.currency),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const FAQ());
                },
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
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/access.png", "", LanguageEn.access),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist("assets/images/lock.png", "",
                    LanguageEn.passwordmanagement),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/eyeoff.png", "", LanguageEn.hidebalance),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/hourglass.png", "", LanguageEn.timeout),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
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
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/coins.png", "", LanguageEn.curatedassets),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/import.png", "", LanguageEn.importwallet),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist("assets/images/backup-wallets.png", "",
                    LanguageEn.backupwallet),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/walletmode.png", "", LanguageEn.walletmode),
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
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/help.png", "", LanguageEn.helpandsupport),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist(
                    "assets/images/terms.png", "", LanguageEn.termsofuse),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Language());
                },
                child: iteamlist("assets/images/copyright.png", "",
                    LanguageEn.abouttrovowallet),
              ),
              SizedBox(height: height / 20),
              GestureDetector(
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: LoginPageConfig);
                },
                child:
                    logout("assets/images/logout.png", "", LanguageEn.logout),
              ),
              SizedBox(height: height / 10),
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
              color: notifier.getbluecolor,
            ),
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
            color: notifier.getbluecolor,
          ),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
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
              color: notifier.getbluecolor,
            ),
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
      ),
    );
  }
}
