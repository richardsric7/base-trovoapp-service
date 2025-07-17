import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/custtom_textfild/custtom_password.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../functions/trovo-sdk.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';

class CreatePassword extends StatefulWidget {
  const CreatePassword({Key? key}) : super(key: key);

  @override
  State<CreatePassword> createState() => _CreatePassword();
}

class _CreatePassword extends State<CreatePassword> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  String password = '';

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
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 18,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                Text(
                  "letsgetyoustarted1".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold,
                  ),
                ),
                Text(
                  "letsgetyoustarted2".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold,
                  ),
                ),
                Center(
                  child: Image.asset(
                    "assets/images/palm-recognition.png",
                    height: height / 2.8,
                  ),
                ),
                Text(
                  "enteryourpassword".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 13.sp,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: height / 50),
                CustomPasswordFormField(
                  "password".tr(),
                  notifier.getbluecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  onChanged: (value) {
                    setState(() {
                      password = value!.trim().replaceAll(' ', '');
                    });
                  },
                  validator: validatePassword,
                ),
                SizedBox(height: height / 50),
                CustomPasswordFormField(
                  "confirmPassword".tr(),
                  notifier.getbluecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  validator: validateConfirmPassword,
                ),
                SizedBox(height: height / 20),
                Button(
                  "continuee".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: saveAndProceed,
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
      ),
    );
  }

  String? validatePassword(value) {
    if (value.isEmpty) {
      //return "Enter a password";
      return "passwordemptyerror".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return "hinterrorpassword".tr();
    }

    return null;
  }

  String? validateConfirmPassword(value) {
    if (value.isEmpty) {
      // return "Confirm your password";
      return "confirmpasswordemptyerror".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return "hinterrorpassword".tr();
    }

    if (password != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return "passwordmismatcherror".tr();
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

  void saveAndProceed() async {
    if (validate()) {
      try {
        showLoader(context);
        var account = TrovoWalletSDK().createAccount();
        appState.setTempPassword = password;
        appState.setTempPublicKey = account.publicKey;
        appState.setTempSecretKey = account.secretKey;
        appState.setTempSigner = account.publicKey;

        hideLoader(context);
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: SignupPageConfig,
        );
      } catch (e) {
        hideLoader(context);
      }
    }
  }
}
