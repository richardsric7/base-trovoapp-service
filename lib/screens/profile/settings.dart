import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/pages.dart';
import 'package:gocrypto/screens/Auth/login.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class Settings extends StatefulWidget {
  const Settings({Key? key}) : super(key: key);

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> {
  late ColorNotifier notifier;
  bool value = false;

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
        appBar: CustomAppBar(
            notifier.getwihitecolor, "Settings", notifier.getblck,
            height: height / 15),

      ),
    );
  }

  Widget darkmode() {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 13),
          Text(
            "Dark Mode",
            style: TextStyle(
                fontFamily: 'Gilroy_Medium',
                fontSize: 17.sp,
                color: notifier.getblck),
          ),
          const Spacer(),
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
          SizedBox(width: width / 20),
        ],
      ),
    );
  }

  Widget fingerprintandfaceid() {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 13),
          Text(
            "Fingerprint and Face ID",
            style: TextStyle(
                fontFamily: 'Gilroy_Medium',
                fontSize: 17.sp,
                color: notifier.getblck),
          ),
          const Spacer(),
          Transform.scale(
            scale: 0.7,
            child: CupertinoSwitch(
              activeColor: notifier.getbluecolor,
              value: value,
              onChanged: (v) {
                setState(() {
                  value = v;
                });
              },
            ),
          ),
          // FlutterSwitch(toggleSize: 15.sp,
          //     height: height / 38,inactiveColor: notifier.getbluecolor,
          //     width: width / 10,
          //     activeToggleColor: notifier.getwihitecolor,
          //     activeColor: notifier.getbluecolor,
          //     value: value,
          //     onToggle: (v) {
          //       value = v;
          //     }),
          SizedBox(width: width / 20),
        ],
      ),
    );
  }

  Widget settingoption(name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 13),
          Text(
            name,
            style: TextStyle(
                fontFamily: 'Gilroy_Medium',
                fontSize: 17.sp,
                color: notifier.getblck),
          ),
          const Spacer(),
          Icon(Icons.arrow_forward_ios, color: notifier.getgrey, size: 17.sp),
          SizedBox(width: width / 20),
        ],
      ),
    );
  }

  Widget settingpolecy() {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 13),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                "Privacy Policy",
                style: TextStyle(
                    fontFamily: 'Gilroy_Medium',
                    fontSize: 17.sp,
                    color: notifier.getblck),
              ),
              Text(
                "Choose what data you share with us",
                style: TextStyle(
                    fontFamily: 'Gilroy_Medium',
                    fontSize: 13.sp,
                    color: notifier.getgrey),
              ),
            ],
          ),
          const Spacer(),
          Icon(Icons.arrow_forward_ios, color: notifier.getgrey, size: 17.sp),
          SizedBox(width: width / 20),
        ],
      ),
    );
  }
}
