import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl_phone_field/intl_phone_field.dart';
import 'package:toggle_switch/toggle_switch.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/services/push_fcm_service.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/models/user.dart';
import '../../functions/trovo-sdk.dart';
import '../../router/page_actions.dart';
import '../../storage/state.dart';
import '../../network/requests.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../../widgets/terms_of_service.dart';

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
  final referrerController = TextEditingController();
  final secretKeyController = TextEditingController();
  late FocusNode passPhraseFocusNode;
  late FocusNode secretKeyFocusNode;
  bool usePassPhrase = false;
  String passPhrase = '';
  String secretKey = '';
  bool importMode = false;
  String errorText = '';

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
    passPhraseFocusNode = FocusNode();
    secretKeyFocusNode = FocusNode();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    state = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    if (state.tempReferrerUsername.isNotEmpty) {
      referrerController.text = state.tempReferrerUsername;
      state.tempReferrerUsername = '';
    }
    if (state.viewData![SignupPageConfig.key] != null &&
        state.viewData![SignupPageConfig.key]['importMode']) {
      secretKeyController.text = state.tempSecretKey;
      importMode = true;
      state.viewData![SignupPageConfig.key] = null;
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        resizeToAvoidBottomInset: false,
        appBar: CustomAppBar(
                context, notifier.getwihitecolor, "", notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Form(
                    key: _formKey,
                    child: Column(
                      children: [
                        SizedBox(height: height / 50),
                        Text(
                          "ittakesaminute1".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 30.sp,
                              fontFamily: fontsemibold),
                        ),
                        Text(
                          "ittakesaminute2".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 30.sp,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(height: height / 20),
                        ToggleSwitch(
                          minHeight: height / 16,
                          customWidths: [
                            width / 2.4,
                            width / 2.4,
                          ],
                          customTextStyles: [
                            TextStyle(
                                fontSize: height / 55,
                                color: corporate == 0 ? wihitecolor : darkblck,
                                fontFamily: fontbody),
                            TextStyle(
                                fontSize: height / 55,
                                color: corporate == 1 ? wihitecolor : darkblck,
                                fontFamily: fontbody),
                          ],
                          fontSize: 16.0,
                          initialLabelIndex: corporate,
                          activeBgColor: [
                            notifier.isDark
                                ? notifier.getbluecolor50
                                : notifier.getbluecolor,
                          ],
                          inactiveBgColor: notifier.getsplashgrey,
                          inactiveFgColor: notifier.getblck,
                          totalSwitches: 2,
                          labels: ['Individual', 'Corporate'],
                          onToggle: (index) {
                            setState(() {
                              corporate = index!;
                            });
                          },
                        ),
                        SizedBox(height: height / 50),
                        getNameFields(),
                        SizedBox(height: height / 50),
                        // Email address
                        CustomTextFormField.textField(
                          "emailadress".tr(),
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
                            email = value.trim().replaceAll(' ', '');
                          },
                          keyboardtype: TextInputType.emailAddress,
                        ),
                        SizedBox(height: height / 50),
                        // Username
                        CustomTextFormField.textField(
                          "username".tr(),
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
                            username = value.trim().replaceAll(' ', '');
                          },
                        ),
                        SizedBox(height: height / 50),
                        // Phone Number
                        phoneFormField(
                          labletext: "phonenumber".tr(),
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
                          "referrer".tr(),
                          notifier.getbluecolor,
                          Icons.link,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          controller: referrerController,
                          validator: validateReferrer,
                          onSaved: (value) =>
                              referrer = value.trim().replaceAll(' ', ''),
                          maxLength: 16,
                        ),
                        Row(
                          children: [
                            Container(
                              width: width / 1.2,
                              child: checkUseImportMode(),
                            ),
                          ],
                        ),
                        if (importMode) ...[
                          if (usePassPhrase) ...[
                            // Pass phrase/Mnemonic
                            passPhraseInput(
                              '${"passphrase".tr()} (optional)',
                              notifier.getbluecolor,
                              notifier.getgrey,
                              notifier.getblck,
                              notifier.getgrey,
                              100.sp,
                              300.sp,
                              onSaved: (value) {
                                passPhrase = value;
                              },
                              minLines: 3,
                              maxLines: null,
                              keyboardtype: TextInputType.multiline,
                              focusNode: passPhraseFocusNode,
                            ),
                          ] else ...[
                            // Secret Key
                            CustomPasswordFormField(
                              '${"secretkey".tr()} (optional)',
                              notifier.getbluecolor,
                              Icons.lock,
                              notifier.getgrey,
                              notifier.getprefixicon,
                              notifier.getblck,
                              70.sp,
                              300.sp,
                              validator: (value) {
                                var trimmedVal =
                                    value!.trim().replaceAll(' ', '');
                                if (trimmedVal.isNotEmpty &&
                                    trimmedVal.length < 56) {
                                  return "secretkeyinvalid".tr();
                                }
                                return null;
                              },
                              onSaved: (value) {
                                secretKey = value!.trim().replaceAll(' ', '');
                              },
                              controller: secretKeyController,
                              maxLength: 56,
                              focusNode: secretKeyFocusNode,
                            )
                          ],
                          Row(
                            children: [
                              Container(
                                width: width / 1.2,
                                child: checkUsePassphrase(),
                              ),
                            ],
                          ),
                        ],
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
                "signup".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () => _validateAndSave(),
              ),
              SizedBox(height: height / 40),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "alreadyregistered".tr(),
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 13.sp,
                        fontFamily: fontbody),
                  ),
                  GestureDetector(
                    onTap: () {
                      state.currentAction = state.userInfo == null
                          ? PageAction(
                              state: PageState.addPage,
                              page: ImportWalletPageConfig,
                            )
                          : PageAction(
                              state: PageState.replaceAll,
                              page: LoginPageConfig,
                            );
                    },
                    child: Text(
                      ' ${state.userInfo == null ? "importwallet".tr() : "signin".tr()}',
                      style: TextStyle(
                          color: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
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

  Widget passPhraseInput(
    labletext,
    focuscolor,
    lablecolor,
    textcolor,
    bordercolor,
    h,
    w, {
    onChanged,
    maxLength,
    minLines,
    maxLines,
    validator,
    onSaved,
    keyboardtype,
    focusNode,
  }) {
    return Container(
      color: Colors.transparent,
      height: h,
      width: w,
      child: TextFormField(
        focusNode: focusNode,
        maxLength: maxLength,
        minLines: minLines,
        maxLines: maxLines,
        style: TextStyle(color: textcolor, fontFamily: fontbody),
        cursorColor: lablecolor,
        onChanged: onChanged,
        decoration: InputDecoration(
          hintText: labletext,
          hintStyle: TextStyle(color: lablecolor),
          disabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(15.sp),
          ),
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
        keyboardType: keyboardtype,
        validator: validator,
        onSaved: onSaved,
      ),
    );
  }

  Widget checkUseImportMode() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          children: [
            Row(
              children: [
                Text(
                  'Import existing wallet',
                  style: TextStyle(
                      fontSize: height / 55,
                      color: notifier.getblck,
                      fontFamily: fontbody),
                ),
              ],
            ),
          ],
        ),
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(5.sp),
              ),
            ),
            activeColor: notifier.isDark
                ? notifier.getbluecolor50
                : notifier.getbluecolor90,
            side: BorderSide(
              color: notifier.isDark
                  ? notifier.getbluecolor50
                  : notifier.getbluecolor90,
            ),
            value: importMode,
            onChanged: (bool? value) {
              setState(() {
                importMode = value!;
                if (usePassPhrase) {
                  passPhraseFocusNode.requestFocus();
                } else {
                  secretKeyFocusNode.requestFocus();
                }
              });
            },
          ),
        ),
      ],
    );
  }

  Widget checkUsePassphrase() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  "enterpassphrase".tr(),
                  style: TextStyle(
                      fontSize: height / 55,
                      color: notifier.getblck,
                      fontFamily: fontbody),
                ),
              ],
            ),
          ],
        ),
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(5.sp),
              ),
            ),
            activeColor: notifier.isDark
                ? notifier.getbluecolor50
                : notifier.getbluecolor90,
            side: BorderSide(
              color: notifier.isDark
                  ? notifier.getbluecolor50
                  : notifier.getbluecolor90,
            ),
            value: usePassPhrase,
            onChanged: (bool? value) {
              setState(() {
                usePassPhrase = value!;
                if (usePassPhrase) {
                  passPhraseFocusNode.requestFocus();
                } else {
                  secretKeyFocusNode.requestFocus();
                }
              });
            },
          ),
        ),
      ],
    );
  }

  Widget getNameFields() {
    if (corporate == 0) {
      return Column(
        children: [
          // Firstname
          CustomTextFormField.textField(
            "fanme".tr(),
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
            "lname".tr(),
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
            "entityname".tr(),
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
    return Container(
      color: Colors.transparent,
      height: h,
      width: w,
      child: IntlPhoneField(
        autovalidateMode: AutovalidateMode.disabled,
        disableLengthCheck: countryCode ==
            'ID', // disable when user selects indonesia and let backend validate
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
          counterStyle: TextStyle(color: textcolor),
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
          setState(() {
            countryCode = value.code;
          });
        },
      ),
    );
  }

  String? validateEmail(String? value) {
    String pattern =
        r'^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$';
    RegExp regex = new RegExp(pattern);

    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return "emailvalidateempty".tr();
    }

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "emailvalidateinvalid".tr();
    }

    return null;
  }

  String? validateUsername(String? value) {
    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return "usernamevalidateempty".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 3 ||
        value.trim().replaceAll(' ', '').length > 16) {
      return "usernamevalidatelength".tr();
    }

    if (num.tryParse(value.trim().replaceAll(' ', '')) != null) {
      return "usernamevalidatenumber".tr();
    }

    String pattern = r'^(?!.*\.\.)(?!.*\.$)[^\W][\w]{2,16}$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', '')) ||
        value.contains('_')) {
      return "usernamevalidateinvalid".tr();
    }

    return null;
  }

  String? validateReferrer(String? value) {
    String pattern = r'^(?!.*\.\.)(?!.*\.$)[^\W][\w]{2,16}$';
    RegExp regex = new RegExp(pattern);

    if (value!.isNotEmpty && value.trim().replaceAll(' ', '').length < 3 ||
        value.trim().replaceAll(' ', '').length > 16) {
      return "usernamevalidatelength".tr();
    }

    if (value.isNotEmpty && !regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "usernamevalidateinvalid".tr();
    }

    if (value.contains('_')) {
      return "usernamevalidateinvalid".tr();
    }

    return null;
  }

  Future<Account?> getCredsFromPassPhrase() async {
    try {
      var trimmedPassprase = passPhrase.trimLeft().trimRight();
      Account account = await TrovoWalletSDK()
          .retrieveCredentialsFromPassPhrase(trimmedPassprase);
      return account;
    } catch (e) {
      print(e);
      // must be some sort of server error
      // let's throw it
      popup(context, title: "error".tr(), message: "invalidcredentials".tr());
      return null;
    }
  }

  String? validateFName(String? value) {
    String pattern = r'(?:\d+[a-z]|[a-z]+\d)[a-z\d]*';
    RegExp regex = new RegExp(pattern);

    if (value!.trim().replaceAll(' ', '').isEmpty) {
      return "firstnamevalidateempty".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 3) {
      return "namevalidatelength".tr();
    }

    if (corporate == 0 && regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "invalidname".tr();
    }

    return null;
  }

  String? validateEntityName(String? value) {
    String pattern = r'(?:\d+[a-z]|[a-z]+\d)[a-z\d]*';
    RegExp regex = new RegExp(pattern);
    var trimmedValue = value!.trimLeft().trimRight();

    if (trimmedValue.isEmpty) {
      return "entitynamevalidateempty".tr();
    }

    if (trimmedValue.length < 3) {
      return "namevalidatelength".tr();
    }

    if (regex.hasMatch(trimmedValue)) {
      return "invalidname".tr();
    }

    return null;
  }

  String? validateLName(String? value) {
    String pattern = r'(?:\d+[a-z]|[a-z]+\d)[a-z\d]*';
    RegExp regex = new RegExp(pattern);

    if (value!.isEmpty) {
      return "lastnamevalidateempty".tr();
    }

    if (value.trim().replaceAll(' ', '').length < 3) {
      return "namevalidatelength".tr();
    }

    if (corporate == 0 && regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "invalidname".tr();
    }

    return null;
  }

  Account? parseKey(String secretKey) {
    try {
      Account account =
          TrovoWalletSDK().parseSecretKey(secretKey.toUpperCase());
      return account;
    } catch (e) {
      print(e);
      // must be some sort of server error
      // let's throw it
      popup(context, title: "error".tr(), message: "invalidcredentials".tr());
      return null;
    }
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

      String result = await FCM().getPushNotificationToken();
      var token = result.split('|').first;

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
      Account? creds = null;

      if (!usePassPhrase && secretKey.isNotEmpty) {
        creds = parseKey(secretKey)!;
      } else if (usePassPhrase && passPhrase.isNotEmpty) {
        creds = await getCredsFromPassPhrase();
      }

      if (creds != null) {
        state.tempPublicKey = creds.publicKey;
        state.tempSecretKey = creds.secretKey;
      }

      Map responseData = await makePostRequest(
          uri: '/v1/users',
          body: jsonBody,
          signer: state.tempPublicKey,
          publicKey: state.tempPublicKey,
          secretKey: state.tempSecretKey);

      // print('$responseData');
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

        state.currentAction =
            PageAction(state: PageState.addPage, page: VerificationPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(context,
          title: "error".tr(),
          message: e.toString().contains('firebase')
              ? 'Network error! Please check your connection and try again.'
              : e.toString());
    }
  }
}
