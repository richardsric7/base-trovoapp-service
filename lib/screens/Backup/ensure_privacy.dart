import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/screens/Backup/backup.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../Auth/fingerprint.dart';
import '../ImportWallet/importwallet.dart';

class EnsurePrivacy extends StatefulWidget {
  const EnsurePrivacy({Key? key}) : super(key: key);

  @override
  State<EnsurePrivacy> createState() => _EnsurePrivacyState();
}

class _EnsurePrivacyState extends State<EnsurePrivacy> {
  bool isNotLooking = false;
  bool isSecurelyStored = false;
  bool isLiable = false;
  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 6),
              Text(
                LanguageEn.backup,
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Transform.scale(
                    scale: 1.sp,
                    child: Checkbox(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.all(
                          Radius.circular(5.sp),
                        ),
                      ),
                      activeColor: notifier.getbluecolor,
                      side: BorderSide(color: notifier.getbluecolor),
                      value: isNotLooking,
                      onChanged: (value) => setState(() {
                        isNotLooking = value!;
                      }),
                    ),
                  ),
                  Column(
                    children: [
                      Container(
                        width: width / 1.2,
                        child: Text(
                          LanguageEn.iensuredprivacy,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 15.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                    ],
                  )
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Transform.scale(
                    scale: 1.sp,
                    child: Checkbox(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.all(
                          Radius.circular(5.sp),
                        ),
                      ),
                      activeColor: notifier.getbluecolor,
                      side: BorderSide(color: notifier.getbluecolor),
                      value: isSecurelyStored,
                      onChanged: (value) => setState(() {
                        isSecurelyStored = value!;
                      }),
                    ),
                  ),
                  Column(
                    children: [
                      Container(
                        width: width / 1.2,
                        child: Text(
                          LanguageEn.iunderstandimportanceofsecretkey,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 15.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                    ],
                  )
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Transform.scale(
                    scale: 1.sp,
                    child: Checkbox(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.all(
                          Radius.circular(5.sp),
                        ),
                      ),
                      activeColor: notifier.getbluecolor,
                      side: BorderSide(color: notifier.getbluecolor),
                      value: isLiable,
                      onChanged: (value) => setState(() {
                        isLiable = value!;
                      }),
                    ),
                  ),
                  Column(
                    children: [
                      Container(
                        width: width / 1.2,
                        child: Text(
                          LanguageEn.iunderstandliability,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 15.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                    ],
                  )
                ],
              ),
              SizedBox(height: height / 4.3),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  if (isNotLooking && isLiable && isSecurelyStored) {
                    Navigator.pushReplacement(
                      context,
                      MaterialPageRoute(
                        builder: (context) => const Backup(),
                      ),
                    );
                  } else {
                    popup(context,
                        title: LanguageEn.important,
                        message: LanguageEn.ensureAll);
                  }
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
