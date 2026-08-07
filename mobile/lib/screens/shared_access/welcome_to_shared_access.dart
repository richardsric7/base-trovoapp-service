import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/storage/store.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WelcomeToSharedAccess extends StatefulWidget {
  const WelcomeToSharedAccess({super.key});

  @override
  State<WelcomeToSharedAccess> createState() => _WelcomeToSharedAccess();
}

class _WelcomeToSharedAccess extends State<WelcomeToSharedAccess> {
  late DataProvider appState;
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "sharedaccess".tr(),
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(
                    vertical: 15.0, horizontal: 25.0),
                child: RichText(
                  text: TextSpan(
                    text: "welcometosharedaccess".tr(),
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                    children: [
                      TextSpan(
                        text: '${"vieweraccess".tr()} ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: "describevieweraccess".tr(),
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: '${"initiatoraccess".tr()} ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: "describeinitiatoraccess".tr(),
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: '${"approveraccess".tr()} ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: "describeapproveraccess".tr(),
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: '${"note".tr()}: ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: "moresharedaccessdetails".tr(),
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  textAlign: TextAlign.justify,
                ),
              ),
              SizedBox(height: height / 20),
              Button(
                "continuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () async {
                  await StoreData()
                      .storeInsertData('introducedSharedAccess', true);
                  appState.setIntroducedSharedAccess = true;
                  appState.currentAction = PageAction(
                      state: PageState.replace,
                      page: SharedAccessViewPageConfig);
                },
              ),
              SizedBox(height: height / 15),
            ],
          ),
        ),
      ),
    );
  }
}
