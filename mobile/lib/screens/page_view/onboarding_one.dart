import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Oneonbonding extends StatefulWidget {
  const Oneonbonding({Key? key}) : super(key: key);

  @override
  State<Oneonbonding> createState() => _OneonbondingState();
}

class _OneonbondingState extends State<Oneonbonding> {
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
    notifier = Provider.of<ColorNotifier>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return Scaffold(
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            SizedBox(height: height / 10.5),
            Image.asset("assets/images/crypto-p2p.png", height: height / 2.3),
            SizedBox(height: height / 20),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: width / 15),
              child: Column(
                children: [
                  Text(
                    "welcometotrovowallet".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 29.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
