import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class GetStarted extends StatefulWidget {
  const GetStarted({Key? key}) : super(key: key);

  @override
  State<GetStarted> createState() => _GetStartedState();
}

class _GetStartedState extends State<GetStarted> {
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
    var appState = Provider.of<DataProvider>(context, listen: true);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 10.5),
              Image.asset("assets/images/palm-recognition.png",
                  height: height / 2.3),
              SizedBox(height: height / 20),
              Text(
                "whatwouldyou".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 25.sp,
                    fontFamily: fontsemibold),
              ),
              Text(
                "liketodo".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 25.sp,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 30.5),
              Button(
                "getstarted".tr(),
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage, page: CreatePasswordPageConfig);
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                "importwallet".tr(),
                notifier.getwihitecolor,
                notifier.getbluecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage, page: ImportWalletPageConfig);
                },
              ),
              SizedBox(height: height / 50),
              ButtonOutlined(
                "recoveraccount".tr(),
                notifier.getbluecolor80,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: RecoverAccountViewPageConfig);
                },
              ),
              SizedBox(height: height / 50)
            ],
          ),
        ),
      ),
    );
  }
}
