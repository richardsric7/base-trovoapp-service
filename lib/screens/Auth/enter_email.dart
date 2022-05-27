import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottombar.dart';
import 'package:gocrypto/screens/Auth/signup.dart';
import 'package:gocrypto/screens/Auth/vericication.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../reset_password/phone_num_reset_password.dart';

class EnterEmail extends StatefulWidget {
  const EnterEmail({Key? key}) : super(key: key);

  @override
  State<EnterEmail> createState() => _EnterEmailState();
}

class _EnterEmailState extends State<EnterEmail> {
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
        appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 10.5),
              Image.asset(
                "assets/images/mailbox.png",
                height: height / 4,
              ),
              SizedBox(height: height / 50),
              Text(
                LanguageEn.enteremail,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 25.sp,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 30),
              Text(
                LanguageEn.enteremailgetstarted,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 13.sp,
                    fontFamily: fontbody),
              ),
              SizedBox(height: height / 40),
              Customtextfild.textField(
                LanguageEn.emailadress,
                notifier.getbluecolor,
                Icons.email,
                notifier.getgrey,
                notifier.getprefixicon,
                notifier.getblck,
                notifier.getgrey,
                45.sp,
                300.sp,
              ),
              SizedBox(height: height / 5),
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const Veryfication(),
                    ),
                  );
                },
                child: Button(LanguageEn.getstarted, notifier.getbluecolor,
                    notifier.getwihitecolor),
              ),
              SizedBox(height: height / 10),
            ],
          ),
        ),
      ),
    );
  }
}
