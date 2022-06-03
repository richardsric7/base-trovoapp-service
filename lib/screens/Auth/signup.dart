import 'dart:async';
import 'dart:convert';

import 'package:cool_dropdown/cool_dropdown.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:intl_phone_field/countries.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import 'package:intl_phone_field/phone_number.dart';
import 'package:toggle_switch/toggle_switch.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/login.dart';
import 'package:trovo_wallet/screens/Auth/privacypolicy.dart';
import 'package:trovo_wallet/screens/Auth/termsofservice.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';
import 'package:trovo_wallet/services/push_fcm_service.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/provider.dart';
import '../../network/requests.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';

class SignUp extends StatefulWidget {
  const SignUp({Key? key}) : super(key: key);

  @override
  State<SignUp> createState() => _SignUpState();
}

class _SignUpState extends State<SignUp> {
  late ColorNotifier notifier;
  late DataProvider state;
  final _formKey = GlobalKey<FormState>();
  late String fName;
  late String lName;
  late String email;
  late String phoneNumber;
  String countryCode = "NG";
  late String username;
  late String referrer;
  late String password;
  late int corporate = 0;
  var accountTypes = [
    {'label': 'Individual Account', 'value': 0},
    {'label': 'Corporate Account', 'value': 1},
  ];
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
                        Text(
                          LanguageEn.signup,
                          style: TextStyle(
                              color: notifier.getblck,
                              fontSize: 26.sp,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(height: height / 35),
                        Text(
                          LanguageEn.ittakesaminute,
                          style: TextStyle(
                              fontSize: 14.sp,
                              color: notifier.getgrey,
                              fontFamily: fontbody),
                        ),
                        SizedBox(height: height / 50),
                        Text(
                          'Account Type',
                          style: TextStyle(
                              fontSize: height / 55,
                              color: notifier.getblck,
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
                                color: notifier.getblck,
                                fontFamily: fontbody),
                            TextStyle(
                                fontSize: height / 55,
                                color: notifier.getblck,
                                fontFamily: fontbody),
                          ],
                          fontSize: 16.0,
                          initialLabelIndex: corporate,
                          activeBgColor: [notifier.getbluecolor],
                          activeFgColor: notifier.getwihitecolor,
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
                          validator: validateFName,
                          onSaved: (value) => fName = value,
                        ),
                        SizedBox(height: height / 50),
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
                          validator: validateLName,
                          onSaved: (value) => lName = value,
                        ),
                        SizedBox(height: height / 50),
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
                            email = value;
                          },
                        ),
                        SizedBox(height: height / 50),
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
                          validator: validateUsername,
                          onSaved: (value) {
                            print('username: $value');
                            username = value;
                          },
                        ),
                        SizedBox(height: height / 50),
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
                          onSaved: (value) => referrer = value,
                        ),
                        SizedBox(height: height / 50),
                        Row(
                          children: [
                            Transform.scale(
                              scale: 1.sp,
                              child: Checkbox(
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.all(
                                    Radius.circular(5.sp),
                                  ),
                                ),
                                activeColor: notifier.getbluecolor,
                                side: BorderSide(color: notifier.getbluecolor),
                                value: hasAgreed,
                                onChanged: (bool? value) {
                                  setState(() {
                                    hasAgreed = value!;
                                  });
                                },
                              ),
                            ),
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Row(
                                  children: [
                                    Text(
                                      LanguageEn.iagreetothe,
                                      style: TextStyle(
                                          fontSize: height / 55,
                                          color: notifier.getblck,
                                          fontFamily: fontbody),
                                    ),
                                    GestureDetector(
                                      onTap: () => {
                                        Get.to(() => const TermsofService())
                                      },
                                      child: Text(
                                        ' ' + LanguageEn.termsofservices,
                                        style: TextStyle(
                                            fontFamily: fontbody,
                                            fontSize: height / 55,
                                            color: notifier.getbluecolor),
                                      ),
                                    ),
                                    Text(
                                      LanguageEn.and,
                                      style: TextStyle(
                                          fontFamily: fontbody,
                                          fontSize: height / 55,
                                          color: notifier.getblck),
                                    ),
                                  ],
                                ),
                                GestureDetector(
                                  onTap: () =>
                                      {Get.to(() => const PrivacyPolicy())},
                                  child: Text(
                                    LanguageEn.privacypolicy,
                                    style: TextStyle(
                                        fontFamily: fontbody,
                                        fontSize: height / 55,
                                        color: notifier.getbluecolor),
                                  ),
                                ),
                              ],
                            )
                          ],
                        ),
                        if (showError) ...[
                          Text(
                            'You need to accept terms',
                            style: TextStyle(
                                color: Colors.red,
                                fontSize: 13,
                                fontWeight: FontWeight.w400),
                          ),
                        ],
                      ],
                    ),
                  ),
                ],
              ),
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
          validator: validatePhone,
        ),
      ),
    );
  }

  FutureOr<String?> validatePhone(PhoneNumber? number) {
    print('validating phone number ...');
  }

  String? validateEmail(String? value) {
    String pattern =
        r'^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$';
    RegExp regex = new RegExp(pattern);

    if (value!.isEmpty) {
      return LanguageEn.emailvalidateempty;
    } else if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return LanguageEn.emailvalidateinvalid;
    }
    return null;
  }

  String? validateUsername(String? value) {
    if (value!.isEmpty) {
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
    } else if (value.isNotEmpty &&
        !regex.hasMatch(value.trim().replaceAll(' ', ''))) {
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
    print('fname: $value');
    if (value!.isEmpty) {
      return LanguageEn.firstnamevalidateempty;
    } else if (value.trim().replaceAll(' ', '').length < 2) {
      return LanguageEn.firstnamevalidatelength;
    }
    return null;
  }

  String? validateLName(String? value) {
    print('lname: $value');
    if (value!.isEmpty) {
      return LanguageEn.lastnamevalidateempty;
    } else if (value.trim().replaceAll(' ', '').length < 2) {
      return LanguageEn.lastnamevalidatelength;
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

      Map responseData = await makePostRequest(
          uri: '/v1/users',
          body: jsonBody,
          signer: publicKey,
          publicKey: publicKey,
          secretKey: secretKey);
      print('$responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        await StoreData().storeInsertData('username', username);
        await StoreData().storeInsertData('firstname', fName);
        await StoreData().storeInsertData('lastname', lName);
        await StoreData().storeInsertData('email', email);
        await StoreData().storeInsertData('mobile', phoneNumber);
        await StoreData().storeInsertData('mobileCountryCode', countryCode);
        await StoreData().storeInsertData('referrer', referrer);
        await StoreData().storeInsertData('pushNotificationToken', token);
        await StoreData().storeInsertData('corporate', corporate);

        state.setUser = UserInfo(
          username: username,
          firstName: fName,
          lastName: lName,
          email: email,
          phoneNumber: phoneNumber,
          countryCode: countryCode,
          referrer: referrer,
          token: token,
          corporate: corporate,
        );

        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => const Veryfication(),
          ),
        );
      } else {
        errorPopup(context,
            title: LanguageEn.error,
            message: LanguageEn.errormessage + responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      errorPopup(context,
          title: LanguageEn.error,
          message: LanguageEn.errormessage + e.toString());
    }
  }
}
