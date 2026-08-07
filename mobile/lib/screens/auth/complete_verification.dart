import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

import '../backup/congratulation.dart';

class Complateerification extends StatefulWidget {
  const Complateerification({Key? key}) : super(key: key);

  @override
  State<Complateerification> createState() => _ComplateerificationState();
}

class _ComplateerificationState extends State<Complateerification> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 6),
              Center(
                child: Image.asset("assets/images/vrfcomplate.png",
                    height: height / 3.3),
              ),
              SizedBox(height: height / 11),
              Text(
                "youreverified".tr(),
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Text(
                "youhavebeensucces".tr(),
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: fontbody),
              ),
              SizedBox(height: height / 4.3),
              Button(
                "continuee".tr(),
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  Navigator.pushReplacement(
                    context,
                    MaterialPageRoute(
                      builder: (context) => Congratulations(),
                    ),
                  );
                },
              )
            ],
          ),
        ),
      ),
    );
  }
}
