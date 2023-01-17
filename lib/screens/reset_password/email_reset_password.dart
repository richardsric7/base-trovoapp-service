import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class EmailResetPassword extends StatefulWidget {
  const EmailResetPassword({Key? key}) : super(key: key);

  @override
  State<EmailResetPassword> createState() => _EmailResetPasswordState();
}

class _EmailResetPasswordState extends State<EmailResetPassword> {
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
        appBar: CustomAppBar(
          context,
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
                        "Reset Password",
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 26.sp,
                            fontFamily: 'Gilroy_Bold'),
                      ),
                      SizedBox(height: height / 25),
                      Text(
                        "Enter your email and we will send you a link\nto reset your password.",
                        style:
                            TextStyle(fontSize: 16.sp, color: notifier.getgrey),
                      ),
                      SizedBox(height: height / 17),
                      Customtextfild.textField(
                          "Email address",
                          notifier.getbluecolor,
                          Icons.email,
                          notifier.getgrey,
                          notifier.getblck,
                          notifier.getblck,
                          notifier.getgrey,
                          45.sp,
                          300.sp)
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 2),
              GestureDetector(
                  onTap: () {
                    Get.to(
                      const Veryfication(),
                    );
                  },
                  child: Button("Send OTP Code", notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }
}
