import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../router/PageActions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../Auth/fingerprint.dart';
import '../ImportWallet/importwallet.dart';
import 'ensure_privacy.dart';

class Congratulations extends StatelessWidget {
  Congratulations({Key? key}) : super(key: key);
  late DataProvider appState;

  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 6),
              Text(
                LanguageEn.congratulations,
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              Text(
                '${appState.userInfo!.username!}!',
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.all(15.0),
                child: Text(
                  LanguageEn.walletcreatesuccess,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 20.sp,
                      wordSpacing: 3.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 7.3),
              Button(
                LanguageEn.backup,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  ensureBackupPrivacyDialog(
                    context,
                    () {
                      Navigator.of(context).pop();
                      appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: EnsurePrivacyPageConfig);
                    },
                  );
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                LanguageEn.skip,
                notifier.getwihitecolor,
                notifier.getbluecolor,
                onTap: () {
                  warnSkipBackupDialog(context, gotoNext);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  gotoNext() async {
    var isFirstTime = await StoreData().storeGetData('isFirstTime') ?? true;
    if (isFirstTime) {
      appState.currentAction =
          PageAction(state: PageState.addPage, page: EnsurePrivacyPageConfig);
    } else {
      appState.currentAction =
          PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
    }
  }
}
