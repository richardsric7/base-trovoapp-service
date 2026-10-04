import 'dart:async';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/store.dart';
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

// AccountRecoverySuccess follows a started recovery for the recovering
// device: PENDING until the recovery period is over (the account owner can
// still cancel), then COMPLETED, when the new key can import the account.
class _AccountRecoverySuccess extends State<AccountRecoverySuccess> {
  late ColorNotifier notifier;
  late DataProvider appState;
  String status = 'PENDING';
  DateTime? executeAfter;
  Timer? timer;

  Future<void> checkStatus() async {
    try {
      Map res = await makeGetRequest(
        uri: '/v1/account/recovery/status/${appState.tempUsername}',
        signer: appState.tempAddress,
        address: appState.tempAddress,
        secretKey: appState.tempSecretKey,
      );
      if (res['statusCode'] != 200 || !mounted) return;
      final s = res['data']['status'] ?? '';
      setState(() {
        if (s != '') status = s;
        if (res['data']['executeAfter'] != null) {
          executeAfter = DateTime.tryParse(res['data']['executeAfter'])
              ?.toLocal();
        }
      });
      if (status != 'PENDING') timer?.cancel();
    } catch (_) {}
  }

  String get details {
    switch (status) {
      case 'COMPLETED':
        return "otpcongratulationsdetails".tr();
      case 'CANCELED':
        return "recoverycanceleddetails".tr();
      case 'FAILED':
        return "recoveryfaileddetails".tr();
      default:
        final when = executeAfter == null
            ? ''
            : DateFormat.yMMMd().add_jm().format(executeAfter!);
        return "recoverypendingdetails".tr(args: [when]);
    }
  }

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
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final data = appState.viewData?[AccountRecoverySuccessViewPageConfig.key];
      if (data != null && data['executeAfter'] != null) {
        setState(() {
          executeAfter = DateTime.tryParse(data['executeAfter'])?.toLocal();
        });
      }
      checkStatus();
      timer = Timer.periodic(Duration(seconds: 30), (_) => checkStatus());
    });
  }

  @override
  void dispose() {
    timer?.cancel();
    super.dispose();
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
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 20,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "account".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 30,
                      fontFamily: fontsemibold,
                    ),
                  ),
                  SizedBox(width: width / 50),
                  Text(
                    "recovery".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 30,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ],
              ),
              Image.asset(
                "assets/images/startup-launch.png",
                height: height / 3.5,
              ),
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
                          horizontal: 20.0,
                          vertical: 15.0,
                        ),
                        child: Container(
                          width: width / 1.3,
                          child: Column(
                            children: [
                              Text(
                                status == 'COMPLETED'
                                    ? '${"congratulations".tr()} ${appState.tempUsername}'
                                    : status == 'PENDING'
                                    ? "recoverypendingtitle".tr()
                                    : "accountrecovery".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                              SizedBox(height: 2),
                              Text(
                                details,
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody,
                                ),
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
              SizedBox(height: height / 20),
              Button(
                "done".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () async {
                  bool isFirstTime =
                      await StoreData().storeGetData('isFirstTime') ?? true;
                  if (isFirstTime) {
                    setState(() {
                      appState.currentAction = PageAction(
                        state: PageState.replaceAll,
                        page: OnboardingPageConfig,
                      );
                    });
                  } else {
                    setState(() {
                      appState.currentAction = PageAction(
                        state: PageState.replaceAll,
                        page: LoginPageConfig,
                      );
                    });
                  }
                },
              ),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
