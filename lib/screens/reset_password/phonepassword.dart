import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/reset_password/emailpassword.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class PhonePassword extends StatefulWidget {
  const PhonePassword({Key? key}) : super(key: key);

  @override
  State<PhonePassword> createState() => _PhonePasswordState();
}

class _PhonePasswordState extends State<PhonePassword> {
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
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Center(
                  child: Image.asset("assets/images/lock.png",
                      height: height / 3.9)),
              SizedBox(height: height / 25),
              Text(
                LanguageEn.resetpass,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 22.sp,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 100),
              Text(
                LanguageEn.enteranemailadress,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: fontbody),
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
              SizedBox(height: height / 4.7),
              GestureDetector(
                  onTap: () {
                    Get.to(() => const Emailpassword());
                  },
                  child: Button(LanguageEn.continuee, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }
}
