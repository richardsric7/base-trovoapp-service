import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/Icons/my_flutter_app_icons.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/popups.dart';

import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class GetStarted extends StatefulWidget {
  const GetStarted({Key? key}) : super(key: key);

  @override
  State<GetStarted> createState() => _GetStartedState();
}

class _GetStartedState extends State<GetStarted> {
  late ColorNotifier notifier;
  late DataProvider appState;
  final _dropDownKey = GlobalKey<FormFieldState>();
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
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20.5),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    Container(
                      width: width / 2.8,
                      child: Row(
                        children: [
                          Expanded(
                            child: DropdownButtonFormField(
                              isExpanded: true,
                              key: _dropDownKey,
                              dropdownColor: notifier.isDark
                                  ? darktilewhitecolor
                                  : notifier.getaddsubwalletgrey,
                              value: appState.walletMode,
                              icon: Icon(Icons.keyboard_arrow_down_rounded),
                              decoration: InputDecoration(
                                contentPadding: EdgeInsets.symmetric(
                                  vertical: 8.0,
                                  horizontal: 10,
                                ),
                                enabledBorder: OutlineInputBorder(
                                  borderSide: BorderSide.lerp(
                                    BorderSide(color: notifier.getgrey),
                                    BorderSide(color: notifier.getgrey),
                                    1.0,
                                  ),
                                  borderRadius: const BorderRadius.all(
                                    Radius.circular(20.0),
                                  ),
                                ),
                                border: OutlineInputBorder(
                                  borderSide: BorderSide.lerp(
                                    BorderSide(color: notifier.getgrey),
                                    BorderSide(color: notifier.getgrey),
                                    1.0,
                                  ),
                                  borderRadius: const BorderRadius.all(
                                    Radius.circular(20.0),
                                  ),
                                ),
                              ),
                              elevation: 0,
                              style: TextStyle(
                                color: notifier.getdarkgrey,
                                fontSize: 13.5.sp,
                                fontFamily: fontbody,
                              ),
                              onChanged: handleEnvironmentSwitch,
                              items: <DropdownMenuItem<String>>[
                                DropdownMenuItem(
                                  child: Row(
                                    children: [
                                      Icon(
                                        CustomIcon.globeOutlined,
                                        size: 18,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                      SizedBox(width: 6),
                                      Text(
                                        "testnet".tr(),
                                        overflow: TextOverflow.ellipsis,
                                        style: TextStyle(
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                    ],
                                  ),
                                  value: 'Testnet',
                                ),
                                DropdownMenuItem(
                                  child: Row(
                                    children: [
                                      Icon(
                                        CustomIcon.globeOutlined,
                                        size: 18,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                      SizedBox(width: 6),
                                      Text(
                                        "mainnet".tr(),
                                        overflow: TextOverflow.ellipsis,
                                        style: TextStyle(
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                    ],
                                  ),
                                  value: 'Mainnet',
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              Image.asset(
                "assets/images/palm-recognition.png",
                height: height / 2.3,
              ),
              SizedBox(height: height / 30),
              Text(
                "whatwouldyou".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: notifier.getblck,
                  fontSize: 25.sp,
                  fontFamily: fontsemibold,
                ),
              ),
              Text(
                "liketodo".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: notifier.getblck,
                  fontSize: 25.sp,
                  fontFamily: fontsemibold,
                ),
              ),
              SizedBox(height: height / 30.5),
              Button(
                "getstarted".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: CreatePasswordPageConfig,
                  );
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                "importwallet".tr(),
                notifier.getwihitecolor,
                notifier.isDark ? wihitecolor : notifier.getbluecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: ImportWalletPageConfig,
                  );
                },
              ),
              SizedBox(height: height / 50),
              ButtonOutlined(
                "recoveraccount".tr(),
                notifier.getbluecolor80,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: RecoverAccountViewPageConfig,
                  );
                },
              ),
              SizedBox(height: height / 50),
            ],
          ),
        ),
      ),
    );
  }

  void handleEnvironmentSwitch(String? newValue) async {
    if (newValue != appState.walletMode) {
      showSwitchEnvironmentPopup(
        context,
        onProceed: () async {
          await appState.changeWalletMode(newValue.toString());
        },
        onCancel: () {
          _dropDownKey.currentState!.reset();
        },
        toEnvironment: newValue!,
      );
    }
  }
}
