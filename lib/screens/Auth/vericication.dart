import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/screens/Auth/complateverification.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Veryfication extends StatefulWidget {
  const Veryfication({Key? key}) : super(key: key);

  @override
  State<Veryfication> createState() => _VeryficationState();
}

class _VeryficationState extends State<Veryfication> {
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
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        LanguageEn.enterverification,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 23.sp,
                            fontFamily: 'Gilroy_Bold'),
                      ),
                      SizedBox(height: height / 30),
                      Text(
                        LanguageEn.enterfourdigitnumber,
                        style:
                            TextStyle(fontSize: 14.sp, color: notifier.getgrey),
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  otpfild(),
                  otpfild(),
                  otpfild(),
                  otpfild(),
                ],
              ),
              SizedBox(height: height / 40),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.resetcode,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Medium'),
                  ),
                  SizedBox(width: width / 100),
                  Text(
                    "29:58",
                    style: TextStyle(
                        color: notifier.getbluecolor,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Medium'),
                  ),
                ],
              ),
              SizedBox(height: height / 15),
              GestureDetector(
                  onTap: () {
                    Get.to(
                      const Complateerification(),
                    );
                  },
                  child: Button(LanguageEn.verify, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }

  Widget otpfild() {
    return Container(
      height: height / 13,
      width: width / 6,
      decoration: BoxDecoration(
        border: Border.all(color: notifier.getbluecolor, width: 1.sp),
        borderRadius: BorderRadius.all(
          Radius.circular(15.sp),
        ),
      ),
      child: TextFormField(
        keyboardType: TextInputType.number,
        onChanged: (value) {
          if (value.length == 1) {
            FocusScope.of(context).nextFocus();
          }
        },
        textAlign: TextAlign.center,
        style: TextStyle(fontSize: 25.sp),
        decoration: const InputDecoration(border: InputBorder.none),
      ),
    );
  }
}
