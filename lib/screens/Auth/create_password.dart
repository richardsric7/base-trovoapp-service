import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/screens/Auth/signup.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../Custom_BlocObserver/fonts.dart';
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
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 15.5),
                Center(
                  child: Icon(
                    CupertinoIcons.lock,
                    color: notifier.getbluecolor,
                    size: 200.sp,
                  ),
                ),
                SizedBox(height: height / 50),
                Text(
                  LanguageEn.enterpassword,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 25.sp,
                      fontFamily: fontsemibold),
                ),
                SizedBox(height: height / 30),
                Text(
                  LanguageEn.enteryourpassword,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 13.sp,
                      fontFamily: fontbody),
                ),
                SizedBox(height: height / 50),
                CustomPasswordFormField(
                  LanguageEn.password,
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
                  LanguageEn.confirmPassword,
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
                  LanguageEn.continuee,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: saveAndProceed,
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
    print('confirm password: ${value.trim().replaceAll(' ', '')} & $password');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.confirmpasswordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    if (password != value.trim().replaceAll(' ', '')) {
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

  void saveAndProceed() async {
    if (validate()) {
      try {
        showLoader(context);
        // hide keyboard because of the bad effect it has on the
        // next screen. This is just a hack, will work out a better
        // solution later
        // TODO: Find better way to solve the keyboard overlay issue
        FocusScope.of(context).requestFocus(FocusNode());
        var account = TrovoWalletSDK().createAccount();
        await StoreData().storeInsertData('secretKey', account.secretKey);
        await StoreData().storeInsertData('publicKey', account.publicKey);
        await StoreData().storeInsertData('password', password);
        hideLoader(context);
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => const SignUp(),
          ),
        );
      } catch (e) {
        hideLoader(context);
        print('we ran into and error $e');
      }
    }
  }
}
