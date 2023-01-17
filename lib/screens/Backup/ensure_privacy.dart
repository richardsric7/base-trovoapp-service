import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
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
        appBar: CustomAppBar(
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
                LanguageEn.backup,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  LanguageEn.iensuredprivacy,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  LanguageEn.iunderstandimportanceofsecretkey,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  LanguageEn.iunderstandliability,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
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
                      }),
              SizedBox(height: height / 4.3),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  if (hasAccepted) {
                    gotoNext();
                  } else {
                    popup(context,
                        title: LanguageEn.important,
                        message: LanguageEn.ensureaccepted);
                  }
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget acceptAll(
    value,
    onChanged,
  ) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(5.sp),
              ),
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
                LanguageEn.iunderstandall,
                style: TextStyle(
                    fontSize: height / 55,
                    color: notifier.getgrey,
                    fontFamily: fontbody),
              ),
            ),
          ],
        )
      ],
    );
  }

  gotoNext() async {
    var data = appState.viewData![EnsurePrivacyPageConfig.key];
    print('gotoNext: $data');
    if (data != null && data['rel'] == 'backupAll') {
      appState.currentAction =
          PageAction(state: PageState.addPage, page: BackupAllViewPageConfig);
    } else if (data != null && data['rel'] == 'accountRecovery') {
      appState.currentAction = PageAction(
          state: PageState.addPage, page: BackupRecoverySecretViewPageConfig);
    } else if (data != null && data['rel'] == 'restoreUnactivatedAccount') {
      appState.currentAction = PageAction(
          state: PageState.addPage, page: BackupRecoverySecretViewPageConfig);
    } else {
      appState.currentAction =
          PageAction(state: PageState.addPage, page: BackupPageConfig);
    }
  }
}
