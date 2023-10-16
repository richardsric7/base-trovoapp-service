import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Congratulations extends StatefulWidget {
  const Congratulations({super.key});

  @override
  State<Congratulations> createState() => _Congratulations();
}

class _Congratulations extends State<Congratulations> {
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
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 10),
              Text(
                '${"congratulations".tr()}',
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              Text(
                '${appState.userInfo!.username!}!',
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.all(15.0),
                child: Text(
                  appState.activeWallet!.primaryWallet == 1
                      ? "walletcreatesuccess".tr()
                      : "subwalletcreatesuccess".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 20.sp,
                      wordSpacing: 3.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Button(
                "backup".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  ensureBackupPrivacyDialog(
                    context,
                    () {
                      Navigator.of(context).pop();
                      appState.viewData![EnsurePrivacyPageConfig.key] = null;
                      appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: EnsurePrivacyPageConfig);
                    },
                  );
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                "skip".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  warnSkipBackupDialog(context, gotoNext);
                },
              ),
              SizedBox(height: height / 15),
            ],
          ),
        ),
      ),
    );
  }

  gotoNext() async {
    if (appState.isFirstTime) {
      appState.currentAction =
          PageAction(state: PageState.addPage, page: FingerprintPageConfig);
    } else if ((appState.returnView != null &&
            appState.returnView!.pages != null) &&
        appState.returnView!.pages!.contains(WalletPreparationViewPageConfig)) {
      appState.currentAction = PageAction(
          state: PageState.addPage, page: SharedAccessViewPageConfig);
    } else {
      appState.currentAction =
          PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
    }
  }
}
