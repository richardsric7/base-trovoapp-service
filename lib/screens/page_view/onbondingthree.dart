import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/Auth/login.dart';
import 'package:gocrypto/screens/page_view/onbonding_two.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/fonts.dart';
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
          child: Center(
            child: Column(
              children: [
                SizedBox(height: height / 10.5),
                Image.asset("assets/images/crypto-p2p.png", height: height / 3),
                SizedBox(height: height / 20),
                Text(
                  LanguageEn.startedDiscover,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 25.sp,
                      fontFamily: fontsemibold),
                ),
                SizedBox(height: height / 50),
                Text(
                  LanguageEn.starttradingyourmoney,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 13.sp,
                      fontFamily: fontbody),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
