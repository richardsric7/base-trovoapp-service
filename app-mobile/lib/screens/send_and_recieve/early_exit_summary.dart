import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/bottom_tab_page.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

import 'package:local_auth/error_codes.dart' as auth_error;

class EarlyExitSummaryView extends StatefulWidget {
  const EarlyExitSummaryView({Key? key}) : super(key: key);

  @override
  State<EarlyExitSummaryView> createState() => _EarlyExitSummaryView();
}

class _EarlyExitSummaryView extends State<EarlyExitSummaryView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  bool understandPenaltyAndEffects = false;
  late Asset? asset;
  final formKey = GlobalKey<FormState>();
  String password = '';
  final Authenticator _authenticator = Authenticator();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (appState.viewData!['walletAddress'] != null) {
      wallet = appState.userInfo!.getWallet(
        appState.viewData!['walletAddress'],
      );
      asset = wallet.claimedAssets!.firstWhere(
        (asset) =>
            asset.assetCode == appState.viewData!['assetCode'] &&
            asset.contractAddress == appState.viewData!['contractAddress'],
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          '',
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Estimated Payout Summary',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(10.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      SizedBox(height: height / 50),
                      item(
                        "Token Quantity to Exit",
                        '1,000 ${asset?.assetCode}',
                      ),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Current NAV/Token", 'N1,000'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Early Exit Penalty + Fees", '3%'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Payout Price/Token", 'N970'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Estimated Payout Amount", 'N970,000'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Payout Currency", 'NGN'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Wallet ID", 'johnnydoe'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Account Number", '1234567890'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Account Name", 'John Doe Kenechukwu'),
                      Divider(color: notifier.getsplashgrey, thickness: 1),
                      item("Settlement Time Estimate", '48 hours'),
                      SizedBox(height: height / 50),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.08,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    Transform.scale(
                      scale: 1.sp,
                      child: Checkbox(
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(Radius.circular(5.sp)),
                        ),
                        activeColor: notifier.isDark
                            ? notifier.getbluecolor50
                            : notifier.getbluecolor90,
                        side: BorderSide(
                          color: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
                        ),
                        value: understandPenaltyAndEffects,
                        onChanged: (bool? value) {
                          setState(() {
                            understandPenaltyAndEffects = value!;
                          });
                        },
                      ),
                    ),
                    Container(
                      width: width / 1.3,
                      child: Text(
                        "I understand that early exit may involve penalties and affects my future returns",
                        overflow: TextOverflow.visible,
                        style: TextStyle(
                          fontSize: 15,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 20),
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  "password".tr(),
                  notifier.getbluewhitecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  validator: (String? value) {
                    if (value!.isEmpty) return 'Enter your password';

                    if (value.length < 6)
                      return 'Use 6 characters or more for your password';

                    return null;
                  },
                  onChanged: (value) {
                    setState(() {
                      password = value!.trim().replaceAll(' ', '');
                    });
                  },
                ),
              ),
              SizedBox(height: height / 30),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  "authorizewithbiometrics".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  "authorize".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: handleAuthorization,
                ),
              ],
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

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer();
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  void handleAuthorization() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer();
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  sendDataToServer() {
    appState.viewData = {
      SuccessViewPageConfig.key: {
        'title': 'Payment request submitted',
        'message':
            'You have successfully requested early exit of [1000 ${asset!.assetCode.toString().isEmpty ? 'ETH' : asset!.assetCode}] from [${wallet.alias}]. This transaction will be completed when it gets the required number of approvals by those who have approver access on this wallet.',
        'useOnDone': true,
        'onDone': () {
          appState.currentAction = PageAction(
            state: PageState.replaceAll,
            page: BottomHomePageConfig,
          );
          changeTabPage(appState, ButtomTabPage.Dashboard.index);
        },
      },
    };
    appState.currentAction = PageAction(
      state: PageState.replace,
      page: SuccessViewPageConfig,
    );
  }

  Widget item(String key, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 0, horizontal: 15),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            constraints: BoxConstraints(maxWidth: width / 2.36),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
            child: Text(
              key,
              textAlign: TextAlign.start,
              overflow: TextOverflow.visible,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 13.sp,
                fontFamily: fontbody,
              ),
            ),
          ),
          Container(
            width: width / 2.56,
            child: Text(
              value,
              textAlign: TextAlign.end,
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 13.sp,
                fontFamily: fontsemibold,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
