import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/Auth/login.dart';
import 'package:gocrypto/screens/page_view/onbonding_two.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class Threeonbonding extends StatefulWidget {
  const Threeonbonding({Key? key}) : super(key: key);

  @override
  State<Threeonbonding> createState() => _ThreeonbondingState();
}

class _ThreeonbondingState extends State<Threeonbonding> {
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
              SizedBox(height: height / 13.5),
              Image.asset("assets/images/Laptop.png", height: height / 2),
              SizedBox(height: height / 20),
              Text(
                LanguageEn.startedDiscover,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 25.sp,
                    fontFamily: 'Gilroy_Bold'),
              ),
              SizedBox(height: height / 50),
              Text(
                LanguageEn.starttradingyourmoney,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 13.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
              SizedBox(height: height / 10),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  GestureDetector(
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => const Login(),
                          ),
                        );
                      },
                      child: button(LanguageEn.skip, notifier.getwihitecolor,
                          notifier.getbluecolor)),
                  SizedBox(width: width / 50),
                  GestureDetector(
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => const Onbondingtwo(),
                          ),
                        );
                      },
                      child: button(LanguageEn.next, notifier.getbluecolor,
                          notifier.getwihitecolor)),
                ],
              ),
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
              width: width / 2.4,
              decoration: BoxDecoration(
                border: Border.all(color: notifier.getbluecolor),
                color: colorbutton,
                borderRadius: BorderRadius.circular(15),
              ),
              child: Center(
                child: Text(
                  buttontext,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontFamily: 'Gilroy_Medium',
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
