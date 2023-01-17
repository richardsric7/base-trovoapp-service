import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class ReferralCode extends StatefulWidget {
  const ReferralCode({Key? key}) : super(key: key);

  @override
  State<ReferralCode> createState() => _ReferralCodeState();
}

class _ReferralCodeState extends State<ReferralCode> {
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

  // copybord() async {
  //   ClipboardData data = ClipboardData(text: '<Text to copy goes here>');
  //   await Clipboard.setData(data);
  // }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "Referral Code", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 10),
              Center(
                child: Image.asset("assets/images/ReferralCode.png",
                    height: height / 6),
              ),
              SizedBox(height: height / 10),
              Text(
                "Refer & Earn",
                style: TextStyle(fontFamily: 'Gilroy_Bold', fontSize: 27.sp),
              ),
              SizedBox(height: height / 100),
              Text(
                "Share this code with your friend and both of\nyou will get \$10 free stocks.",
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 13.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
              SizedBox(height: height / 15),
              Row(
                children: [
                  SizedBox(width: width / 8),
                  Text(
                    "ZOZO123",
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 16.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                  const Spacer(),
                  GestureDetector(
                    onTap: () {
                      _showSnackBar();
                      Clipboard.setData(const ClipboardData(text: "ZOZO123"));
                    },
                    child: Icon(
                      Icons.copy,
                      color: notifier.getbluecolor,
                      size: 16.sp,
                    ),
                  ),
                  SizedBox(width: width / 70),
                  GestureDetector(
                    onTap: () {
                      _showSnackBar();
                      Clipboard.setData(const ClipboardData(text: "ZOZO123"));
                    },
                    child: Text(
                      "Copy Code",
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontFamily: 'Gilroy_Bold',
                          fontSize: 14.sp),
                    ),
                  ),
                  SizedBox(width: width / 8),
                ],
              ),
              SizedBox(height: height / 20),
              GestureDetector(
                  onTap: () {
                    share();
                  },
                  child: Button("Refer friend", notifier.getbluecolor,
                      notifier.getwihitecolor)),
            ],
          ),
        ),
      ),
    );
  }

  void _showSnackBar() {
    const snack = SnackBar(
        content: Text("Copy Successfully"), duration: Duration(seconds: 2));
    ScaffoldMessenger.of(context).showSnackBar(snack);
  }

  Future<void> share() async {
    await FlutterShare.share(
        title: 'Example share',
        text: 'Example share text',
        linkUrl: 'https://flutter.dev/',
        chooserTitle: 'Example Chooser Title');
  }
  // void _copyText(String text) {
  //   FlutterClipboard.copy(field.text).then((value) {
  //     _showSnackBar();
  //   });
  // }
  //
  // void _pastText() {
  //   FlutterClipboard.paste().then((value) {
  //     setState(() {
  //       field.text = value;
  //       pasteValue = value;
  //     });
  //   });
  // }
  // void _showSnackBar() {
  //   const snack =
  //   SnackBar(content: Text("Text copied"), duration: Duration(seconds: 2));
  //   ScaffoldMessenger.of(context).showSnackBar(snack);
  // }
}
