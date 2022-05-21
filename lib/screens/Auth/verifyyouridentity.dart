import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/reset_password/cameraverification.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:gocrypto/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class VerifyYourIdentity extends StatefulWidget {
  const VerifyYourIdentity({Key? key}) : super(key: key);

  @override
  State<VerifyYourIdentity> createState() => _VerifyYourIdentityState();
}

class _VerifyYourIdentityState extends State<VerifyYourIdentity> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
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
                        LanguageEn.letsverify,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 23.sp,
                            fontFamily: 'Gilroy_Bold'),
                      ),
                      SizedBox(height: height / 30),
                      Text(
                        LanguageEn.chooseyourdocument,
                        style:
                            TextStyle(fontSize: 14.sp, color: notifier.getgrey),
                      ),
                      SizedBox(height: height / 8),
                      Text(
                        LanguageEn.methodofverification,
                        style:
                            TextStyle(fontSize: 14.sp, color: notifier.getgrey),
                      ),
                      SizedBox(height: height / 50),
                      methodverification(
                          LanguageEn.passport, "assets/images/passport.png"),
                      SizedBox(height: height / 50),
                      methodverification(LanguageEn.identitycard,
                          "assets/images/identifitycard.png"),
                      SizedBox(height: height / 50),
                      methodverification(LanguageEn.digitaldocument,
                          "assets/images/digitaldocument.png"),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 5),
              GestureDetector(
                  onTap: () {
                    Get.to(const Cameraverification());
                  },
                  child: Button(LanguageEn.continuee, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }

  Widget methodverification(verificationname, image) {
    return Container(
      height: height / 12,
      width: width / 1.14,
      decoration: BoxDecoration(
        color: notifier.getidentyfiymethod,
        borderRadius: BorderRadius.all(
          Radius.circular(15.sp),
        ),
      ),
      child: Row(
        children: [
          Padding(
            padding: const EdgeInsets.only(top: 7),
            child: Image.asset(image),
          ),
          SizedBox(width: width / 80),
          Text(
            verificationname,
            style: TextStyle(
                color: notifier.getblck,
                fontFamily: 'Gilroy_Bold',
                fontSize: 16.sp),
          ),
          const Spacer(),
          Icon(Icons.arrow_forward_ios_rounded,
              color: notifier.getgrey, size: 20.sp),
          SizedBox(width: width / 20),
        ],
      ),
    );
  }
}
