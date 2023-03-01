import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

import '../../custom_bloc_observer/button/custtom_button.dart';

class PrivacyPolicy extends StatefulWidget {
  const PrivacyPolicy({Key? key}) : super(key: key);

  @override
  State<PrivacyPolicy> createState() => _PrivacyPolicyState();
}

class _PrivacyPolicyState extends State<PrivacyPolicy> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBarWithoutBanner(
          context,
          notifier.getwihitecolor,
          LanguageEn.privacypolicy,
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Container(
                height: height / 1.6,
                color: Colors.transparent,
                child: SingleChildScrollView(
                  child: Column(
                    children: [
                      Row(
                        children: [
                          SizedBox(width: width / 20),
                          Text(
                            LanguageEn.lastupdate,
                            style: TextStyle(
                                color: notifier.getgrey,
                                fontSize: 15.sp,
                                fontFamily: fontbody),
                          )
                        ],
                      ),
                      SizedBox(height: height / 20),
                      Row(
                        children: [
                          SizedBox(width: width / 20),
                          Text(
                            LanguageEn.terms,
                            style: TextStyle(
                                color: notifier.getblck,
                                fontSize: 19.sp,
                                fontFamily: fontsemibold),
                          )
                        ],
                      ),
                      SizedBox(height: height / 50),
                      Padding(
                        padding: EdgeInsets.only(
                            left: width / 19, right: width / 22),
                        child: Text(
                          LanguageEn.termslorem,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 15.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                      SizedBox(height: height / 25),
                      Row(
                        children: [
                          SizedBox(width: width / 20),
                          Text(
                            LanguageEn.uselicense,
                            style: TextStyle(
                                color: notifier.getblck,
                                fontSize: 19.sp,
                                fontFamily: fontsemibold),
                          )
                        ],
                      ),
                      SizedBox(height: height / 50),
                      Padding(
                        padding: EdgeInsets.only(
                            left: width / 19, right: width / 22),
                        child: Text(
                          LanguageEn.adipiscingtempus,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 15.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.only(
                            left: width / 19, right: width / 22),
                        child: Text(
                          LanguageEn.adipiscingtempus,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 15.sp,
                              fontFamily: fontbody),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 15),
              GestureDetector(
                  onTap: () {
                    Navigator.pop(context);
                  },
                  child: Button(LanguageEn.done, notifier.getbluecolor,
                      notifier.getwihitecolor)),
            ],
          ),
        ),
      ),
    );
  }
}
