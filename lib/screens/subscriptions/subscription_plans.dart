import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/route_manager.dart';
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

class SubscriptionPlans extends StatefulWidget {
  const SubscriptionPlans({Key? key}) : super(key: key);

  @override
  State<SubscriptionPlans> createState() => _SubscriptionPlansState();
}

class _SubscriptionPlansState extends State<SubscriptionPlans> {
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
                          'Subscription Plans',
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
                      'Subscription types',
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
              planItem(
                'Platinum Patron',
                'Tier 1 - \$10 per month / \$60 per year / \$600 lifetime.',
                'Subscribe to Platinum',
                () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: SubscriptionPlanBenefitsViewPageConfig);
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              planItem(
                'Diamond Patron',
                'Tier 1 - \$10 per month / \$60 per year / \$600 lifetime.',
                'Subscribe to Diamond',
                () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: SubscriptionPlanBenefitsViewPageConfig);
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              planItem(
                'Gold Patron',
                'Tier 1 - \$10 per month / \$60 per year / \$600 lifetime.',
                'Subscribe to Gold',
                () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: SubscriptionPlanBenefitsViewPageConfig);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget planItem(
    String name,
    String description,
    String buttonText,
    void Function() onTap,
  ) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              border: Border.all(color: notifier.getbluewhitecolor, width: 1.5),
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 10.0, vertical: 15.0),
              child: Row(
                children: [
                  Container(
                    width: width / 1.2,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        SizedBox(
                          width: width / 50,
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              name,
                              style: TextStyle(
                                fontSize: 19,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                            Text(
                              'Read more',
                              style: TextStyle(
                                fontStyle: FontStyle.italic,
                                fontSize: 15,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ],
                        ),
                        SizedBox(
                          height: 20,
                        ),
                        Text(
                          description,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(
                          height: 20,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        SizedBox(height: 5),
        Button(
          buttonText,
          notifier.getbluecolor,
          wihitecolor,
          onTap: onTap,
        ),
      ],
    );
  }
}
