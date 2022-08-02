import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class Language extends StatefulWidget {
  const Language({Key? key}) : super(key: key);

  @override
  State<Language> createState() => _LanguageState();
}

class _LanguageState extends State<Language> {
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
        appBar: CustomAppBar(context, notifier.getwihitecolor,
            LanguageEn.selectlanguage, notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 30),
              Center(
                child: Customtextfild.textField(
                    "Search ",
                    notifier.getbluecolor,
                    Icons.search_rounded,
                    notifier.getgrey,
                    notifier.getblck,
                    notifier.getblck,
                    notifier.getgrey,
                    45.sp,
                    300.sp),
              ),
              SizedBox(height: height / 25),
              selectLanguage("assets/images/English.png", LanguageEn.english),
              SizedBox(height: height / 20),
              selectLanguage("assets/images/Deutsch.png", LanguageEn.italy),
              SizedBox(height: height / 20),
              selectLanguage("assets/images/Spanish.png", LanguageEn.japanese),
              SizedBox(height: height / 20),
              selectLanguage("assets/images/French.png", LanguageEn.indonesian),
              SizedBox(height: height / 20),
              selectLanguage(
                  "assets/images/Portuguese.png", LanguageEn.russian),
              SizedBox(height: height / 4),
              GestureDetector(
                onTap: () {
                  Get.back();
                },
                child: Button(LanguageEn.savechange, notifier.getbluecolor,
                    notifier.getwihitecolor),
              )
            ],
          ),
        ),
      ),
    );
  }

  Widget selectLanguage(image, languagename) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 15),
          Image.asset(image, height: height / 30),
          SizedBox(width: width / 20),
          Text(
            languagename,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 16.sp,
                fontFamily: 'Gilroy_Medium'),
          )
        ],
      ),
    );
  }
}
