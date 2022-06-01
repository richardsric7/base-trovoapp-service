import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/pin.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Faceid extends StatefulWidget {
  const Faceid({Key? key}) : super(key: key);

  @override
  State<Faceid> createState() => _FaceidState();
}

class _FaceidState extends State<Faceid> {
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
        appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Center(
                child: Text(
                  LanguageEn.setupfaceid,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 26.sp,
                      fontFamily: fontsemibold),
                ),
              ),
              SizedBox(height: height / 45),
              Center(
                child: Text(
                  LanguageEn.unlockfaceid,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontSize: 16.sp,
                      color: notifier.getgrey,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Center(
                child: Icon(
                  Icons.face,
                  color: notifier.getbluecolor,
                  size: 200.sp,
                ),
              ),
              SizedBox(height: height / 20),
              Row(
                children: [
                  SizedBox(width: width / 25),
                  Icon(
                    Icons.face,
                    color: notifier.getbluecolor,
                    size: 20.sp,
                  ),
                  SizedBox(width: width / 40),
                  Text(
                    LanguageEn.enablefaceid,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 15.sp,
                        fontFamily: fontbody),
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
              SizedBox(height: height / 20),
              GestureDetector(
                  onTap: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => const Pin(),
                      ),
                    );
                  },
                  child: Button(LanguageEn.goahead, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }
}
