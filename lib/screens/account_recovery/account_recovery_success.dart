import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AccountRecoverySuccess extends StatefulWidget {
  const AccountRecoverySuccess({Key? key}) : super(key: key);

  @override
  State<AccountRecoverySuccess> createState() => _AccountRecoverySuccess();
}

class _AccountRecoverySuccess extends State<AccountRecoverySuccess> {
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
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 20),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.account,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 30.sp,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(width: width / 50),
                  Text(
                    LanguageEn.recovery,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 30.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              Image.asset("assets/images/startup-launch.png",
                  height: height / 3.5),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 15.0),
                        child: Container(
                          width: width / 1.3,
                          child: Column(
                            children: [
                              Text(
                                '${LanguageEn.congratulations} ${appState.tempUsername}',
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                    fontSize: 16,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(height: 2),
                              Text(
                                LanguageEn.otpcongratulationsdetails,
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                    fontSize: 16,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              SizedBox(height: 2),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 30),
              SizedBox(
                height: height / 20,
              ),
              Button(
                LanguageEn.done,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () async {
                  bool isFirstTime =
                      await StoreData().storeGetData('isFirstTime') ?? true;
                  if (isFirstTime) {
                    setState(() {
                      appState.currentAction = PageAction(
                          state: PageState.replaceAll,
                          page: OnboardingPageConfig);
                    });
                  } else {
                    setState(() {
                      appState.currentAction = PageAction(
                          state: PageState.replaceAll, page: LoginPageConfig);
                    });
                  }
                },
              ),
              SizedBox(height: height / 10),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }
}
