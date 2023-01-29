import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';

import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SubscriptionPlanBenefits extends StatefulWidget {
  const SubscriptionPlanBenefits({Key? key}) : super(key: key);

  @override
  State<SubscriptionPlanBenefits> createState() =>
      _SubscriptionPlanBenefitsState();
}

class _SubscriptionPlanBenefitsState extends State<SubscriptionPlanBenefits> {
  late ColorNotifier notifier;
  late DataProvider appState;

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
    appState = Provider.of<DataProvider>(context, listen: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(height / 15),
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            title: Container(
              width: width / 1.5,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Padding(
                    padding: EdgeInsets.all(8),
                    child: Column(
                      children: [
                        Text(
                          'Platinum Patron',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getblck,
                              fontSize: 22.sp,
                              fontFamily: fontsemibold),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Platinum Patron - Tier 1',
                      style: TextStyle(
                          fontSize: 20,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  'This is the highest level of patronage in the Trovotech ecosystem. It unlocks every single perk available to our ecosytem.',
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Benefits',
                      style: TextStyle(
                          fontSize: 20,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  "1. Increased referral bonuses.                       ",
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: 10),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  "2. Reduced fees when using any of our products and services.",
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: 10),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  "3. Access to TROV token as incentive (70% of membership fee value in TROV is distributed to the member's wallet).",
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: 10),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  "4. Increased limits to transactions where applicable.",
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: 10),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  "5. Special access to members-only areas of our site.",
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: 10),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  "6. Other privileges, as will be added as we progress.",
                  style: TextStyle(
                      fontSize: 17,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Button(
                'Subscribe to Platinum',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: SubscriptionPlanOptionsViewPageConfig);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
