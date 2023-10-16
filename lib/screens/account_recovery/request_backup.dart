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

class RequestBackup extends StatefulWidget {
  const RequestBackup({super.key});

  @override
  State<RequestBackup> createState() => _RequestBackup();
}

class _RequestBackup extends State<RequestBackup> {
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
                  '${"congratulations".tr()}',
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                      fontSize: 27),
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
                      fontSize: 27),
                ),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.all(15.0),
                child: Text(
                  "backupnewsecretkeygenerated".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 20,
                      wordSpacing: 3,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Button(
                "taptobackup".tr(),
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
