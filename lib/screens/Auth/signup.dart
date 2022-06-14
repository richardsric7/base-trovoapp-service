import 'dart:convert';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import 'package:toggle_switch/toggle_switch.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/login.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';
import 'package:trovo_wallet/services/push_fcm_service.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Models/User.dart';
import '../../storage/state.dart';
import '../../network/requests.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../../widgets/termsOfService.dart';

class SignUp extends StatefulWidget {
  const SignUp({Key? key}) : super(key: key);

  @override
  State<SignUp> createState() => _SignUpState();
}

class _SignUpState extends State<SignUp> {
  late ColorNotifier notifier;
  late DataProvider state;
  final _formKey = GlobalKey<FormState>();
  String fName = '';
  String lName = '';
  late String email;
  late String phoneNumber;
  String countryCode = "NG";
  late String username;
  late String referrer;
  late String password;
  late int corporate = 0;
  bool showError = false;
  bool hasAgreed = false; // to the terms of services

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
    state = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        resizeToAvoidBottomInset: false,
        appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Create Account
                        Text(
                          LanguageEn.signup,
                          style: TextStyle(
                              color: notifier.getblck,
                              fontSize: 26.sp,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(height: height / 35),
                        // It only takes a minute to create your
                        // account
                        Text(
                          LanguageEn.ittakesaminute,
                          style: TextStyle(
                              fontSize: 14.sp,
                              color: notifier.getgrey,
                              fontFamily: fontbody),
                        ),
                        SizedBox(height: height / 50),
                        // Account Type
                        Text(
                          LanguageEn.accounttype,
                          style: TextStyle(
                              fontSize: height / 55,
                              color: notifier.getgrey,
                              fontFamily: fontbody),
                        ),
                        SizedBox(height: height / 70),
                        ToggleSwitch(
                          minHeight: height / 16,
                          customWidths: [
                            width / 2.4,
                            width / 2.4,
                          ],
                          customTextStyles: [
                            TextStyle(
                                fontSize: height / 55,
                                color: corporate == 0
                                    ? notifier.getwihitecolor
                                    : notifier.getblck,
                                fontFamily: fontbody),
                            TextStyle(
                                fontSize: height / 55,
                                color: corporate == 1
                                    ? notifier.getwihitecolor
                                    : notifier.getblck,
                                fontFamily: fontbody),
                          ],
                          fontSize: 16.0,
                          initialLabelIndex: corporate,
                          activeBgColor: [notifier.getbluecolor],
                          activeFgColor: Colors.red,
                          // activeFgColor: notifier.getwihitecolor,
                          inactiveBgColor: notifier.getsplashgrey,
                          inactiveFgColor: notifier.getblck,
                          totalSwitches: 2,
                          labels: ['Individual', 'Corporate'],
                          onToggle: (index) {
                            print('switched to: $index');

                            setState(() {
                              corporate = index!;
                            });
                            print('switched to: $corporate');
                          },
                        ),
                        SizedBox(height: height / 50),
                        getNameFields(),
                        SizedBox(height: height / 50),
                        // Email address
                        CustomTextFormField.textField(
                          LanguageEn.emailadress,
                          notifier.getbluecolor,
                          Icons.email,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          validator: validateEmail,
                          onSaved: (value) {
                            print('email: $value');
                            email = value.trim().replaceAll(' ', '');
                          },
                          keyboardtype: TextInputType.emailAddress,
                        ),
                        SizedBox(height: height / 50),
                        // Username
                        CustomTextFormField.textField(
                          LanguageEn.username,
                          notifier.getbluecolor,
                          Icons.person,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          maxLength: 16,
                          validator: validateUsername,
                          onSaved: (value) {
                            print('username: $value');
                            username = value.trim().replaceAll(' ', '');
                          },
                        ),
                        SizedBox(height: height / 50),
                        // Phone Number
                        phoneFormField(
                          labletext: LanguageEn.phonenumber,
                          focuscolor: notifier.getbluecolor,
                          preicon: Icons.phone,
                          lablecolor: notifier.getgrey,
                          iconcolor: notifier.getprefixicon,
                          textcolor: notifier.getblck,
                          bordercolor: notifier.getgrey,
                          h: 70.sp,
                          w: 300.sp,
                        ),
                        SizedBox(height: height / 50),
                        // Referrer's Username
                        CustomTextFormField.textField(
                          LanguageEn.referrer,
                          notifier.getbluecolor,
                          Icons.link,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          validator: validateReferrer,
                          onSaved: (value) =>
                              referrer = value.trim().replaceAll(' ', ''),
                          maxLength: 16,
                        ),
                        SizedBox(height: height / 50),
                        // Terms of Service
                        TermsOfService(
                          value: hasAgreed,
                          showError: showError,
                          onChanged: (bool? value) {
                            setState(() {
                              hasAgreed = value!;
                            });
                          },
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              // Create Account
              SizedBox(height: height / 25),
              Button(
                LanguageEn.signup,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () => _validateAndSave(),
              ),
              SizedBox(height: height / 40),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.alreadyregistered,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontbody),
                  ),
                  GestureDetector(
                    onTap: () {
                      Get.to(() => const Login());
                    },
                    child: Text(
                      ' ' + LanguageEn.signin,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontSize: 13.sp,
                          fontFamily: fontbody),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 20),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Widget getNameFields() {
    if (corporate == 0) {
      return Column(
        children: [
          // Firstname
          CustomTextFormField.textField(
            LanguageEn.fanme,
            notifier.getbluecolor,
            Icons.person,
            notifier.getgrey,
            notifier.getprefixicon,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            maxLength: 50,
            validator: validateFName,
            onSaved: (value) => fName = value.trim().replaceAll(' ', ''),
          ),
          SizedBox(height: height / 50),
          // Lastname
          CustomTextFormField.textField(
            LanguageEn.lname,
            notifier.getbluecolor,
            Icons.person,
            notifier.getgrey,
            notifier.getprefixicon,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            maxLength: 50,
            validator: validateLName,
            onSaved: (value) => lName = value.trim().replaceAll(' ', ''),
          ),
        ],
      );
    } else {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // EntityName
          CustomTextFormField.textField(
            LanguageEn.entityname,
            notifier.getbluecolor,
            Icons.person,
            notifier.getgrey,
            notifier.getprefixicon,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            maxLength: 50,
            validator: validateEntityName,
            onChanged: (value) {
              setState(() {
                fName = value;
              });
            },
            onSaved: (value) => fName = value.toString().trimLeft().trimRight(),
          ),
        ],
      );
    }
  }

  String getEntityFullName() => '${fName} ${lName}';

  Widget phoneFormField({
    labletext,
    focuscolor,
    preicon,
    lablecolor,
    iconcolor,
    textcolor,
    bordercolor,
    h,
    w,
  }) {
    return ScreenUtilInit(
      builder: (context, child) => Container(
        color: Colors.transparent,
        height: h,
        width: w,
        child: IntlPhoneField(
          autovalidateMode: AutovalidateMode.disabled,
          style: TextStyle(color: textcolor, fontFamily: fontbody),
          cursorColor: lablecolor,
          initialCountryCode: countryCode,
          dropdownIcon: Icon(
            Icons.arrow_drop_down,
            color: textcolor,
          ),
          dropdownTextStyle:
              TextStyle(color: textcolor, fontSize: 16, fontFamily: fontbody),
          decoration: InputDecoration(
            label: Text(labletext),
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            prefixIcon: Icon(preicon, color: iconcolor),
            labelStyle: TextStyle(color: lablecolor),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            enabledBorder: OutlineInputBorder(
              borderSide: BorderSide(color: bordercolor, width: 1),
              borderRadius: BorderRadius.circular(15.sp),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: BorderSide(color: focuscolor, width: 1),
              borderRadius: BorderRadius.circular(15.sp),
            ),
          ),
          onChanged: (value) {
            setState(() {
              phoneNumber = value.completeNumber;
            });
          },
          onCountryChanged: (value) {
            print('country: ' + value.code);
            setState(() {
              countryCode = value.code;
            });
          },
        ),
      ),
    );
  }

  String? validateEmail(String? value) {
    String pattern =
        r'^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$';
    RegExp regex = new RegExp(pattern);

    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return LanguageEn.emailvalidateempty;
    }

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return LanguageEn.emailvalidateinvalid;
    }

    return null;
  }

  String? validateUsername(String? value) {
    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return LanguageEn.usernamevalidateempty;
    }

    if (value.trim().replaceAll(' ', '').length < 3 ||
        value.trim().replaceAll(' ', '').length > 16) {
      return LanguageEn.usernamevalidatelength;
    }

    if (num.tryParse(value.trim().replaceAll(' ', '')) != null) {
      return LanguageEn.usernamevalidatenumber;
    }

    String pattern = r'^(?!.*\.\.)(?!.*\.$)[^\W][\w]{2,16}$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', '')) ||
        value.contains('_')) {
      return LanguageEn.usernamevalidateinvalid;
    }

    return null;
  }

  String? validateReferrer(String? value) {
    print('referrer: $value');
    String pattern = r'^(?!.*\.\.)(?!.*\.$)[^\W][\w]{2,16}$';
    RegExp regex = new RegExp(pattern);

    if (value!.isNotEmpty && value.trim().replaceAll(' ', '').length < 3 ||
        value.trim().replaceAll(' ', '').length > 16) {
      return LanguageEn.usernamevalidatelength;
    }

    if (value.isNotEmpty && !regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return LanguageEn.usernamevalidateinvalid;
    }

    if (value.contains('_')) {
      return LanguageEn.usernamevalidateinvalid;
    }

    return null;
  }

  // String? validateMobile(PhoneNumber? value) {
  //   print('phone: ${phoneNumber!.completeNumber}');
  //   if (phoneNumber!.number.isEmpty) {
  //     return LanguageEn.entervalidmobilenumber;
  //   }
  //   return null;
  // }

  String? validateFName(String? value) {
    String pattern = r'(?:\d+[a-z]|[a-z]+\d)[a-z\d]*';
    RegExp regex = new RegExp(pattern);

    print('fname: $value');
    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return LanguageEn.firstnamevalidateempty;
    }

    if (value.trim().replaceAll(' ', '').length < 3) {
      return LanguageEn.namevalidatelength;
    }

    if (corporate == 0 && regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return LanguageEn.invalidname;
    }

    return null;
  }

  String? validateEntityName(String? value) {
    String pattern = r'(?:\d+[a-z]|[a-z]+\d)[a-z\d]*';
    RegExp regex = new RegExp(pattern);
    var trimmedValue = value!.trimLeft().trimRight();

    print('Entityname: $value');
    if (trimmedValue.isEmpty) {
      return LanguageEn.entitynamevalidateempty;
    }

    if (trimmedValue.length < 3) {
      return LanguageEn.namevalidatelength;
    }

    if (regex.hasMatch(trimmedValue)) {
      return LanguageEn.invalidname;
    }

    return null;
  }

  String? validateLName(String? value) {
    String pattern = r'(?:\d+[a-z]|[a-z]+\d)[a-z\d]*';
    RegExp regex = new RegExp(pattern);

    print('lname: $value');
    if (value!.isEmpty) {
      return LanguageEn.lastnamevalidateempty;
    }

    if (value.trim().replaceAll(' ', '').length < 3) {
      return LanguageEn.namevalidatelength;
    }

    if (corporate == 0 && regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return LanguageEn.invalidname;
    }

    return null;
  }

  bool checkTerms() {
    if (!hasAgreed) {
      // show error message if the user has not
      // agreed to terms and conditions
      setState(() {
        showError = true;
      });
      return false;
    } else {
      setState(() {
        showError = false;
      });
      return true;
    }
  }

  // Check if form is valid
  _validateAndSave() async {
    try {
      final form = _formKey.currentState;
      if (!form!.validate()) {
        checkTerms();
        return;
      }

      // check that terms and conditions has been accepted
      if (!checkTerms()) return;

      showLoader(context);

      form.save();

      String? token = await StoreData().storeGetData('token');

      if (token == null) {
        token = await FCM().getPushNotificationToken();
      }

      Map map = {
        'username': username,
        'email': email,
        'firstName': fName,
        'lastName': lName,
        'mobile': phoneNumber,
        'mobileCountryCode': countryCode,
        'referrer': referrer,
        'pushNotificationToken': token,
        'corporate': corporate,
        'verificationCode': '',
      };

      String jsonBody = jsonEncode(map);
      print(jsonBody);

      var publicKey = await StoreData().storeGetData('publicKey') ?? '';
      var secretKey = await StoreData().storeGetData('secretKey') ?? '';
      print('public: $publicKey, secret: $secretKey');

      Map responseData = await makePostRequest(
          uri: '/v1/users',
          body: jsonBody,
          signer: publicKey,
          publicKey: publicKey,
          secretKey: secretKey);
      print('$responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        state.setUser = UserInfo(
          username: username,
          firstName: fName,
          lastName: lName,
          email: email,
          mobile: phoneNumber,
          countryCode: countryCode,
          referrer: referrer,
          pushNotificationToken: token,
          corporate: corporate,
        );

        print('userInfo: ${state.userInfo}');

        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => const Veryfication(),
          ),
        );
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }
}
