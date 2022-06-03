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

class CreatePassword extends StatefulWidget {
  const CreatePassword({Key? key}) : super(key: key);

  @override
  State<CreatePassword> createState() => _CreatePassword();
}

class _CreatePassword extends State<CreatePassword> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late String password;

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
                SizedBox(height: height / 10.5),
                Image.asset(
                  "assets/images/mailbox.png",
                  height: height / 4,
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
                SizedBox(height: height / 40),
                SizedBox(height: height / 50),
                CustomPasswordFormField.textField(
                  LanguageEn.password,
                  notifier.getbluecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  validator: validatePassword,
                  onSaved: (value) => password = value,
                ),
                SizedBox(height: height / 50),
                CustomPasswordFormField.textField(
                  LanguageEn.confirmPassword,
                  notifier.getbluecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  onChanged: (value) {
                    setState(() {
                      password = value.trim().replaceAll(' ', '');
                    });
                  },
                  validator: validateConfirmPassword,
                  onSaved: (value) => password = value,
                ),
                SizedBox(height: height / 10),
                Button(
                  LanguageEn.continuee,
                  notifier.getbluecolor,
                  notifier.getwihitecolor,
                  onTap: saveAndProceed,
                ),
                // ButtonCustom(LanguageEn.continuee, notifier.getbluecolor,
                //     notifier.getwihitecolor, saveAndProceed),
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
    } else if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }
    return null;
  }

  String? validateConfirmPassword(value) {
    print('confirm password: $value');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.confirmpasswordemptyerror;
    } else if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    } else if (password != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return LanguageEn.passwordmismatcherror;
    }
    return null;
  }

  bool validateAndSave() {
    final form = _formKey.currentState;
    if (form!.validate()) {
      form.save();
      return true;
    }
    return false;
  }

  void saveAndProceed() async {
    if (validateAndSave()) {
      // hide keyboard because of the bad effect it has on the
      // next screen. This is just a hack, will work out a better
      // solution later
      // TODO: Find better way to solve the keyboard overlay issue
      FocusScope.of(context).requestFocus(FocusNode());
      var account = TrovoWalletSDK().createAccount();
      await StoreData().storeInsertData('secretKey', account.secretKey);
      await StoreData().storeInsertData('publicKey', account.publicKey);
      await StoreData().storeInsertData('password', password);
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (context) => const SignUp(),
        ),
      );
    }
  }
}
