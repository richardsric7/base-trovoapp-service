import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/termsofservice.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class Cameraverification extends StatefulWidget {
  const Cameraverification({Key? key}) : super(key: key);

  @override
  State<Cameraverification> createState() => _CameraverificationState();
}

class _CameraverificationState extends State<Cameraverification> {
  late ColorNotifier notifier;
  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        LanguageEn.verifyaccount,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 23.sp,
                            fontFamily: fontsemibold),
                      ),
                      SizedBox(height: height / 40),
                      Text(
                        LanguageEn.takeaphotooffront,
                        style: TextStyle(
                            fontSize: 14.sp,
                            color: notifier.getgrey,
                            fontFamily: fontbody),
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 10),
              Center(
                  child: Image.asset("assets/images/verifyacount.png",
                      height: height / 3)),
              SizedBox(height: height / 4.5),
              GestureDetector(
                  onTap: () {
                    Get.to(() => const TermsofService());
                  },
                  child: Button(LanguageEn.continuee, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }
}
