import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/Auth/face_id.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:gocrypto/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

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
      builder: () => Scaffold(
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
                LanguageEn.youreverified,
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: 'Gilroy_Bold',
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Text(
                LanguageEn.youhavebeensucces,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
              SizedBox(height: height / 4.3),
              GestureDetector(
                  onTap: () {
                    Get.to(const Faceid());
                  },
                  child: Button(LanguageEn.done, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }
}
