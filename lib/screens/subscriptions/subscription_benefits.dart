import 'package:flutter/material.dart';
import 'package:flutter_html/flutter_html.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/models/patronInfo.dart';
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
  late PatronInfo patronInfo;

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
    patronInfo = appState.viewData!['patronInfo'] as PatronInfo;
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
                          '${patronInfo.patronPackage.capitalizeFirst!} Patron',
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
                child: Html(data: patronInfo.description, style: {
                  "*": Style(
                      color: notifier.getbluewhitecolor,
                      fontSize: FontSize.large,
                      lineHeight: LineHeight.number(1.2),
                      wordSpacing: 1.2,
                      textAlign: TextAlign.justify),
                  "h1, h2, h3, h4": Style(
                    fontFamily: fontsemibold,
                    fontSize: FontSize.large,
                  ),
                }),
              ),
              SizedBox(height: height / 20),
              Button(
                'Subscribe to ${patronInfo.patronPackage.capitalizeFirst!}',
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
