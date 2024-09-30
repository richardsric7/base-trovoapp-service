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

class DeleteAccountPrerequisites extends StatefulWidget {
  const DeleteAccountPrerequisites({Key? key}) : super(key: key);

  @override
  State<DeleteAccountPrerequisites> createState() =>
      _DeleteAccountPrerequisitesState();
}

class _DeleteAccountPrerequisitesState
    extends State<DeleteAccountPrerequisites> {
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
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Text(
                  "To delete your account please complete the following steps:"
                      .tr(),
                  // textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w800,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: Colors.blue[50],
                  ),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 10.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.start,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  "Disable Account Recovery",
                                  style: TextStyle(
                                    fontSize: 14,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                                SizedBox(
                                  height: height / 90,
                                ),
                                Container(
                                  width: width / 1.5,
                                  child: Text(
                                    '${"Please disable account recovery on your account.".tr()}',
                                    style: TextStyle(
                                      fontSize: 14,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            OutlinedButton(
                              style: OutlinedButton.styleFrom(
                                backgroundColor: Colors.transparent,
                                shadowColor: Colors.transparent,
                                side: BorderSide(color: Colors.transparent),
                                padding: EdgeInsets.all(0),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.zero,
                                ),
                              ),
                              child: Image.asset(
                                appState.userInfo!.accountRecoveryEnabled == 1
                                    ? "assets/images/export.png"
                                    : "assets/images/tick-circle.png",
                              ),
                              onPressed: () {
                                if (appState.userInfo!.accountRecoveryEnabled ==
                                    1) {
                                  appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page:
                                        DisableAccountRecoveryInfoViewPageConfig,
                                  );
                                }
                              },
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: Colors.blue[50],
                  ),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 10.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.start,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      "Disable Initiator Access",
                                      style: TextStyle(
                                        fontSize: 14,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                    SizedBox(
                                      height: height / 90,
                                    ),
                                    Container(
                                      width: width / 1.5,
                                      child: Text(
                                        "disableorrevokeinitiatoraccess".tr(),
                                        style: TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w400,
                                          color: notifier.getbluewhitecolor,
                                          fontFamily: fontbody,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                OutlinedButton(
                                  style: OutlinedButton.styleFrom(
                                    backgroundColor: Colors.transparent,
                                    shadowColor: Colors.transparent,
                                    side: BorderSide(color: Colors.transparent),
                                    padding: EdgeInsets.all(0),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.zero,
                                    ),
                                  ),
                                  child: Image.asset(
                                    appState.userInfo!.sharedWallets!.length > 0
                                        ? "assets/images/export.png"
                                        : "assets/images/tick-circle.png",
                                  ),
                                  onPressed: () {
                                    if (appState
                                            .userInfo!.sharedWallets!.length >
                                        0) {
                                      appState.currentAction = PageAction(
                                        state: PageState.addPage,
                                        page: SharedAccessViewPageConfig,
                                      );
                                    }
                                  },
                                ),
                              ],
                            ),
                            SizedBox(
                              height: height / 50,
                            ),
                            Container(
                                width: width / 1.2,
                                child: Wrap(
                                  alignment: WrapAlignment.center,
                                  children: [
                                    for (var wallet in appState
                                        .userInfo!.sharedWallets!) ...[
                                      if (wallet.isInitiator) ...[
                                        Padding(
                                          padding: const EdgeInsets.all(3.0),
                                          child: Container(
                                            decoration: BoxDecoration(
                                              border: Border.all(
                                                  color: notifier
                                                      .getbluewhitecolor),
                                              borderRadius:
                                                  const BorderRadius.all(
                                                      Radius.circular(10.0)),
                                            ),
                                            child: Padding(
                                              padding:
                                                  const EdgeInsets.all(5.0),
                                              child: Wrap(
                                                alignment: WrapAlignment.center,
                                                crossAxisAlignment:
                                                    WrapCrossAlignment.center,
                                                children: [
                                                  Text(
                                                    wallet.alias!,
                                                    textAlign: TextAlign.center,
                                                    softWrap: true,
                                                    style: TextStyle(
                                                        color: notifier
                                                            .getbluewhitecolor,
                                                        fontFamily: fontbody,
                                                        fontSize: 14),
                                                  ),
                                                  SizedBox(
                                                    width: width / 70,
                                                  ),
                                                ],
                                              ),
                                            ),
                                          ),
                                        ),
                                      ],
                                    ],
                                    SizedBox(
                                      height: height / 50,
                                    ),
                                  ],
                                )),
                            SizedBox(height: 2),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(
                    horizontal: 20.0, vertical: 20.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      "Please Note:".tr(),
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Image.asset(
                          "assets/images/jam_alert.png",
                        ),
                        Spacer(),
                        Container(
                          width: width / 1.2,
                          child: Text(
                            "deletioninfo1".tr(),
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Image.asset(
                          "assets/images/jam_alert.png",
                        ),
                        Spacer(),
                        Container(
                          width: width / 1.2,
                          child: Text(
                            "deletioninfo2".tr(),
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 25,
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "deleteaccount".tr(),
                appState.userInfo!.sharedWallets!.length > 0 ||
                        appState.userInfo!.accountRecoveryEnabled == 1
                    ? notifier.getbluecolor70
                    : notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  if (appState.userInfo!.accountRecoveryEnabled! == 1 ||
                      appState.userInfo!.sharedWallets!.length > 0) {
                    popup(context,
                        title: "important".tr(),
                        message: "pleasecompleteprerequisites".tr());
                    return;
                  }
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: DeleteAccountViewPageConfig,
                  );
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
}
