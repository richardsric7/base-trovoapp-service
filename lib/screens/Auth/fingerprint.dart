import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../storage/store.dart';
import '../../utils/local_auth.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class FingerPrint extends StatefulWidget {
  const FingerPrint({Key? key}) : super(key: key);

  @override
  State<FingerPrint> createState() => _FingerPrintState();
}

class _FingerPrintState extends State<FingerPrint> {
  late ColorNotifier notifier;
  bool isSwitched = false;
  late DataProvider appState;
  final Authenticator _authenticator = Authenticator();

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
        resizeToAvoidBottomInset: false,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Center(
                child: Text(
                  LanguageEn.biometrics,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 26.sp,
                      fontFamily: fontsemibold),
                ),
              ),
              SizedBox(height: height / 45),
              Center(
                child: Text(
                  LanguageEn.unlockwithbiometrics,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontSize: 16.sp,
                      color: notifier.getgrey,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Center(
                child: Icon(
                  Icons.face_rounded,
                  color: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor,
                  size: 200.sp,
                ),
              ),
              SizedBox(height: height / 20),
              Row(
                children: [
                  SizedBox(width: width / 10),
                  Icon(
                    Icons.face_rounded,
                    color: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor,
                    size: 20.sp,
                  ),
                  SizedBox(width: width / 40),
                  Text(
                    LanguageEn.enablebiometrics,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 15.sp,
                        fontFamily: fontbody),
                  ),
                  const Spacer(),
                  SizedBox(width: width / 100),
                  Transform.scale(
                    scale: 0.7,
                    child: CupertinoSwitch(
                        activeColor: notifier.isDark
                            ? notifier.getbluecolor50
                            : notifier.getbluecolor,
                        value: isSwitched,
                        onChanged: _toggleSwitch),
                  ),
                  SizedBox(width: width / 15),
                ],
              ),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.goahead,
                notifier.getbluecolor,
                wihitecolor,
                onTap: _handleSubmit,
              )
            ],
          ),
        ),
      ),
    );
  }

  void _toggleSwitch(bool value) async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        setState(() {
          isSwitched = !isSwitched;
        });
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  void _handleSubmit() {
    if (!isSwitched) {
      showSkipBiometricsDialog(context);
    } else {
      _submit();
    }
  }

  void _submit() {
    appState.biometricEnabled = isSwitched;
    _persistBiometricState();
    appState.isFirstTime = false;
    appState.currentAction =
        PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
    appState.isLoggedIn = true;
  }

  void _persistBiometricState() async =>
      await StoreData().storeInsertData('biometricsEnabled', isSwitched);

  void showSkipBiometricsDialog(context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    showDialog(
        context: context,
        builder: (context) {
          return AlertDialog(
              scrollable: true,
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(20),
              content: Container(
                decoration: BoxDecoration(
                  color: notifier.getwihitecolor,
                  borderRadius: BorderRadius.all(
                    Radius.circular(23),
                  ),
                ),
                child: Column(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Center(
                        child: Text(
                          LanguageEn.important,
                          style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w500,
                              color: notifier.getblck),
                        ),
                      ),
                    ),
                    Container(
                      constraints: BoxConstraints(
                        maxHeight: height / 5,
                      ),
                      // height: height / 5,
                      child: SingleChildScrollView(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Padding(
                              padding: const EdgeInsets.symmetric(
                                  vertical: 10.0, horizontal: 5.0),
                              child: Text(
                                LanguageEn.skipBiometricsMessage,
                                style: TextStyle(
                                  fontSize: 17,
                                  fontWeight: FontWeight.w300,
                                  color: Colors.red,
                                ),
                                textAlign: TextAlign.center,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(10.0),
                      child: Column(
                        children: [
                          ElevatedButton(
                            onPressed: () => {
                              _submit(),
                            },
                            style: ButtonStyle(
                              fixedSize: MaterialStateProperty.all(
                                Size(width / 1.5, height / 20),
                              ),
                              backgroundColor: MaterialStateProperty.all<Color>(
                                  notifier.getbluecolor),
                              shape: MaterialStateProperty.all<
                                  RoundedRectangleBorder>(
                                const RoundedRectangleBorder(
                                  borderRadius: BorderRadius.all(
                                    Radius.circular(10),
                                  ),
                                ),
                              ),
                            ),
                            child: Text(
                              LanguageEn.skipBiometrics,
                              style: TextStyle(
                                  color: wihitecolor, fontFamily: fontbody),
                            ),
                          ),
                          OutlinedButton(
                            onPressed: () =>
                                Navigator.of(context).pop(), // dismiss dialog,
                            child: Text(
                              LanguageEn.cancel,
                              style: TextStyle(
                                  color: notifier.getbluecolor,
                                  fontFamily: fontbody),
                            ),
                            style: ButtonStyle(
                              fixedSize: MaterialStateProperty.all(
                                Size(width / 1.5, height / 20),
                              ),
                              side: MaterialStateProperty.all(
                                BorderSide(
                                    color: notifier.getgrey,
                                    width: 1,
                                    style: BorderStyle.solid),
                              ),
                              shape: MaterialStateProperty.all<
                                  RoundedRectangleBorder>(
                                const RoundedRectangleBorder(
                                  borderRadius: BorderRadius.all(
                                    Radius.circular(10),
                                  ),
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    SizedBox(height: height / 50),
                  ],
                ),
              ));
        });
  }
}
