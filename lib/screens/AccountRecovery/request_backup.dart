import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../router/PageActions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class RequestBackup extends StatelessWidget {
  RequestBackup({Key? key}) : super(key: key);
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
              SizedBox(height: height / 6),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  '${LanguageEn.congratulations}',
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                      fontSize: 27.sp),
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  '${appState.tempUsername}!',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                      fontSize: 27.sp),
                ),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.all(15.0),
                child: Text(
                  LanguageEn.backupnewsecretkeygenerated,
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
                LanguageEn.taptobackup,
                notifier.getbluecolor,
                wihitecolor,
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
              SizedBox(height: height / 7.3),
            ],
          ),
        ),
      ),
    );
  }
}
