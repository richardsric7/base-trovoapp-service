import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DisableAccountRecoveryInfo extends StatefulWidget {
  const DisableAccountRecoveryInfo({Key? key}) : super(key: key);

  @override
  State<DisableAccountRecoveryInfo> createState() =>
      _DisableAccountRecoveryInfoState();
}

class _DisableAccountRecoveryInfoState
    extends State<DisableAccountRecoveryInfo> {
  late ColorNotifier notifier;
  bool isSwitched = false;
  late DataProvider appState;
  late Wallet primaryWallet;

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
    appState = Provider.of<DataProvider>(context, listen: false);
    primaryWallet = appState.userInfo!.wallets!
        .firstWhere((wallet) => wallet.primaryWallet == 1);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        resizeToAvoidBottomInset: false,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            // crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                LanguageEn.disable,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold),
              ),
              Text(
                LanguageEn.accountrecovery,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 45),
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
                        child: Column(
                          children: [
                            Container(
                              width: width / 1.3,
                              child: Text(
                                LanguageEn.disableaccountrecoverydescription,
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                    fontSize: 16,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                            ),
                            SizedBox(height: 2),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.proceed,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () => appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: DisableAccountRecoveryViewPageConfig),
              )
            ],
          ),
        ),
      ),
    );
  }
}
