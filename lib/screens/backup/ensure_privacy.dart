import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class EnsurePrivacy extends StatefulWidget {
  const EnsurePrivacy({Key? key}) : super(key: key);

  @override
  State<EnsurePrivacy> createState() => _EnsurePrivacyState();
}

class _EnsurePrivacyState extends State<EnsurePrivacy> {
  bool hasAccepted = false;
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
        appBar: CustomAppBarWithoutBanner(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Text(
                "backup".tr(),
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                  fontSize: 27.sp,
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  "iensuredprivacy".tr(),
                  style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: fontbody,
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  "iunderstandimportanceofsecretkey".tr(),
                  style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: fontbody,
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  "iunderstandliability".tr(),
                  style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: fontbody,
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              acceptAll(
                hasAccepted,
                (value) => {
                  print('hasAccepted $hasAccepted'),
                  setState(() {
                    hasAccepted = !hasAccepted;
                  }),
                },
              ),
              SizedBox(height: height / 4.3),
              Button(
                "continuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  if (hasAccepted) {
                    gotoNext();
                  } else {
                    popup(
                      context,
                      title: "important".tr(),
                      message: "ensureaccepted".tr(),
                    );
                  }
                },
              ),
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }

  Widget acceptAll(value, onChanged) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(Radius.circular(5.sp)),
            ),
            activeColor: notifier.getbluecolor,
            side: BorderSide(color: notifier.getbluewhitecolor),
            value: value,
            onChanged: onChanged,
          ),
        ),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              padding: const EdgeInsets.all(16.0),
              width: width / 1.2,
              child: Text(
                "iunderstandall".tr(),
                style: TextStyle(
                  fontSize: height / 55,
                  color: notifier.getgrey,
                  fontFamily: fontbody,
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  gotoNext() async {
    var data = appState.viewData![EnsurePrivacyPageConfig.key];
    print('gotoNext: $data');
    if (data != null && data['rel'] == 'backupAll') {
      appState.currentAction = PageAction(
        state: PageState.addPage,
        page: BackupAllViewPageConfig,
      );
    } else if (data != null && data['rel'] == 'accountRecovery') {
      appState.currentAction = PageAction(
        state: PageState.addPage,
        page: BackupRecoverySecretViewPageConfig,
      );
    } else if (data != null && data['rel'] == 'restoreUnactivatedAccount') {
      appState.currentAction = PageAction(
        state: PageState.addPage,
        page: BackupRecoverySecretViewPageConfig,
      );
    } else {
      appState.currentAction = PageAction(
        state: PageState.addPage,
        page: BackupPageConfig,
      );
    }
  }
}
