import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/Auth/signup.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class Onbondingtwo extends StatefulWidget {
  const Onbondingtwo({Key? key}) : super(key: key);

  @override
  State<Onbondingtwo> createState() => _OnbondingtwoState();
}

class _OnbondingtwoState extends State<Onbondingtwo> {
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
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 6.5),
              Image.asset("assets/images/getstarted.png"),
              SizedBox(height: height / 20),
              Text(
                LanguageEn.bestappto,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 25.sp,
                    fontFamily: 'Gilroy_Bold'),
              ),
              SizedBox(height: height / 40),
              Text(
                LanguageEn.trackmorethan,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 13.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
              SizedBox(height: height / 5.9),
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const SignUp(),
                    ),
                  );
                },
                child: Button(LanguageEn.getstarted, notifier.getbluecolor,
                    notifier.getwihitecolor),
              ),
              // GestureDetector(
              //     onTap: () {
              //       Get.to(
              //         const Login(),
              //       );
              //     },
              //     child: button(
              //         "Login", notifier.getwihitecolor, notifier.getbluecolor)),
            ],
          ),
        ),
      ),
    );
  }
}
