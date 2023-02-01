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

class WelcomeSubscriptions extends StatefulWidget {
  const WelcomeSubscriptions({Key? key}) : super(key: key);

  @override
  State<WelcomeSubscriptions> createState() => _WelcomeSubscriptionsState();
}

class _WelcomeSubscriptionsState extends State<WelcomeSubscriptions> {
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
              width: width / 1.7,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Padding(
                    padding: EdgeInsets.all(8),
                    child: Column(children: [
                      Text(
                        'Trovo Patron',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 22.sp,
                            fontFamily: fontsemibold),
                      ),
                    ]),
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
              SizedBox(height: height / 8),
              Image.asset("assets/images/unlock.png", height: height / 3.7),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.symmetric(horizontal: width / 15),
                child: Column(children: [
                  Text(
                    'Unlock the full potential of the Trovotech ecosystem with Trovo Patron',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 22.sp,
                        fontFamily: fontsemibold),
                  ),
                ]),
              ),
              SizedBox(height: height / 15),
              Padding(
                padding: EdgeInsets.symmetric(horizontal: width / 15),
                child: Column(children: [
                  Text(
                    'You currently are not a Trovo Patron',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 17.sp,
                        fontFamily: fontbody),
                  ),
                ]),
              ),
              SizedBox(height: height / 10),
              Button(
                'View all Trovo Patron plans',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: SubscriptionPlansViewPageConfig);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
