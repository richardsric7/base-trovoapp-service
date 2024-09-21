import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DeleteAccount extends StatefulWidget {
  const DeleteAccount({Key? key}) : super(key: key);

  @override
  State<DeleteAccount> createState() => _DeleteAccountState();
}

class _DeleteAccountState extends State<DeleteAccount> {
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
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
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
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Text(
                "deleteaccount".tr(),
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                  fontSize: 24,
                ),
              ),
              Image.asset('assets/images/deleted.png'),
              deletionInfo(),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "yesdeleteaccount".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  confirmAccountDeletionPopup(context,
                      onConfirmationSuccess: () {
                    appState.viewData![SuccessViewPageConfig.key] = {
                      'title': 'Request successfull',
                      'message':
                          'Your request has been successfully submitted.',
                      'useOnDone': true,
                      'onDone': () {
                        appState.currentAction = PageAction(
                          state: PageState.replaceAll,
                          page: GetStartedViewPageConfig,
                        );
                      },
                    };
                    appState.currentAction = PageAction(
                        state: PageState.replace, page: SuccessViewPageConfig);
                  });
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              ButtonOutlined(
                "nokeepaccount".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  Navigator.of(context).pop();
                },
              ),
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget deletionInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: Colors.red[50],
        ),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0, vertical: 20.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              Text(
                "abouttodeleteaccount".tr(),
                style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold),
              ),
              SizedBox(
                height: height / 90,
              ),
              Text(
                "beforeyoudeleteaccount".tr(),
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w800,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Text(
                '- ${"deletioninfo1".tr()}',
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Text(
                '- ${"deletioninfo2".tr()}',
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Text(
                "doyouwanttoproceed".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w800,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 25,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
