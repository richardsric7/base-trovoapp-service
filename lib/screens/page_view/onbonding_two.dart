import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/fonts.dart';
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
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Center(
            child: Column(
              children: [
                SizedBox(height: height / 10.5),
                Image.asset("assets/images/wallet.png", height: height / 2.5),
                SizedBox(height: height / 20),
                Padding(
                  padding: EdgeInsets.symmetric(horizontal: width / 10),
                  child: Text(
                    LanguageEn.managetrovowallet,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 25.sp,
                        fontFamily: fontsemibold),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
