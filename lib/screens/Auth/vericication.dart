import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:otp_text_field/otp_field.dart';
import 'package:otp_text_field/style.dart';
import 'package:trovo_wallet/Custom_BlocObserver/provider.dart';
import 'package:trovo_wallet/screens/Auth/complateverification.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../network/requests.dart';
import '../../services/push_fcm_service.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../../widgets/popups.dart';

class Veryfication extends StatefulWidget {
  const Veryfication({Key? key}) : super(key: key);

  @override
  State<Veryfication> createState() => _VeryficationState();
}

class _VeryficationState extends State<Veryfication> {
  late ColorNotifier notifier;
  late DataProvider state;
  late String otp;

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
    state = Provider.of<DataProvider>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: height / 30.5),
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        LanguageEn.enterverification,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 23.sp,
                            fontFamily: fontsemibold),
                      ),
                      SizedBox(height: height / 30),
                      Text(
                        LanguageEn.enterfourdigitnumber +
                            state.userInfo!.email!,
                        style: TextStyle(
                            fontSize: 14.sp,
                            color: notifier.getgrey,
                            fontFamily: fontbody),
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 30),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30),
                child: OTPTextField(
                  length: 6,
                  width: MediaQuery.of(context).size.width,
                  fieldWidth: 40,
                  style:
                      TextStyle(color: notifier.getblck, fontFamily: fontbody),
                  textFieldAlignment: MainAxisAlignment.spaceAround,
                  fieldStyle: FieldStyle.box,
                  onChanged: (pin) {
                    print("Changed: " + pin);
                  },
                  onCompleted: (pin) {
                    print("Completed: " + pin);
                    otp = pin;
                    completeRegistration();
                  },
                ),
              ),
              SizedBox(height: height / 40),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.resetcode,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 13.sp,
                        fontFamily: fontbody),
                  ),
                  SizedBox(width: width / 100),
                  Text(
                    "29:58",
                    style: TextStyle(
                        color: notifier.getbluecolor,
                        fontSize: 13.sp,
                        fontFamily: fontbody),
                  ),
                ],
              ),
              SizedBox(height: height / 15),
              GestureDetector(
                  onTap: () {
                    if (otp.length == 6) {
                      completeRegistration();
                    }
                  },
                  child: Button(LanguageEn.verify, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }

  completeRegistration() async {
    try {
      showLoader(context);

      String? token = await StoreData().storeGetData('token');

      if (token == null) {
        token = await FCM().getPushNotificationToken();
      }

      Map map = {
        'username': state.userInfo!.username!,
        'email': state.userInfo!.email,
        'firstName': state.userInfo!.firstName,
        'lastName': state.userInfo!.lastName,
        'mobile': state.userInfo!.phoneNumber,
        'mobileCountryCode': state.userInfo!.countryCode,
        'referrer': state.userInfo!.referrer,
        'pushNotificationToken': state.userInfo!.token,
        'corporate': state.userInfo!.corporate,
        'verificationCode': otp,
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
        secretKey: secretKey,
      );

      print('$responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => const Complateerification(),
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

class UserInfo {
  String? username;
  String? firstName;
  String? lastName;
  String? email;
  String? phoneNumber;
  String? countryCode;
  String? referrer;
  String? token;
  int? corporate;
  String? pushNotificationToken;

  UserInfo({
    this.username,
    this.firstName,
    this.lastName,
    this.email,
    this.phoneNumber,
    this.countryCode,
    this.referrer,
    this.token,
    this.corporate,
    this.pushNotificationToken,
  });
}
