import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AccountRecoverySuccess extends StatefulWidget {
  const AccountRecoverySuccess({Key? key}) : super(key: key);

  @override
  State<AccountRecoverySuccess> createState() => _AccountRecoverySuccess();
}

class _AccountRecoverySuccess extends State<AccountRecoverySuccess> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  String email = '';
  String otp = '';

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
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 20),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 50),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      LanguageEn.account,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontSize: 30.sp,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(width: width / 50),
                    Text(
                      LanguageEn.recovery,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                          color: notifier.getbluecolor80,
                          fontSize: 30.sp,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
                Image.asset("assets/images/startup-launch.png",
                    height: height / 3.5),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
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
                          child: Container(
                            width: width / 1.3,
                            child: Column(
                              children: [
                                Text(
                                  '${LanguageEn.congratulations} ${appState.tempUsername}',
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontsemibold),
                                ),
                                SizedBox(height: 2),
                                Text(
                                  LanguageEn.otpcongratulationsdetails,
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
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
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
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
                          child: Container(
                            width: width / 1.3,
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  '${LanguageEn.secretkey} for ${appState.tempUsername}',
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontsemibold),
                                ),
                                SizedBox(height: 2),
                                Text(
                                  appState.tempSecretKey,
                                  style: TextStyle(
                                      fontSize: 14,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
                                ),
                                SizedBox(height: 2),
                                ElevatedButton(
                                  onPressed: () => {
                                    Clipboard.setData(ClipboardData(
                                        text: appState.tempSecretKey)),
                                    showSnackBar(LanguageEn.secretkey, context),
                                  },
                                  style: ButtonStyle(
                                    backgroundColor:
                                        MaterialStateProperty.all<Color>(
                                            notifier.getbluecolor!),
                                  ),
                                  child: Text(
                                    LanguageEn.copy,
                                    style: TextStyle(
                                      fontFamily: fontsemibold,
                                    ),
                                  ),
                                )
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(
                  height: height / 20,
                ),
                Button(
                  LanguageEn.done,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () async {
                    bool isFirstTime =
                        await StoreData().storeGetData('isFirstTime') ?? true;
                    if (isFirstTime) {
                      setState(() {
                        appState.currentAction = PageAction(
                            state: PageState.replaceAll,
                            page: OnboardingPageConfig);
                      });
                    } else {
                      setState(() {
                        appState.currentAction = PageAction(
                            state: PageState.replaceAll, page: LoginPageConfig);
                      });
                    }
                  },
                ),
                SizedBox(height: height / 10),
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String? validatePassword(value) {
    print('password: $value');
    if (value.isEmpty) {
      //return "Enter a password";
      return LanguageEn.passwordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    return null;
  }

  String? validateConfirmPassword(value) {
    print('confirm password: ${value.trim().replaceAll(' ', '')} & $otp');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.confirmpasswordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    if (otp != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return LanguageEn.passwordmismatcherror;
    }

    return null;
  }

  bool validate() {
    final form = _formKey.currentState;
    if (form!.validate()) {
      form.save();
      return true;
    }
    return false;
  }

  void saveAndProceed() async {}
}
