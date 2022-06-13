import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:otp_text_field/otp_text_field.dart';
import 'package:otp_text_field/style.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/screens/Auth/complateverification.dart';
import 'package:trovo_wallet/screens/Auth/fingerprint.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../Models/User.dart';
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
  // Timer? countdownTimer;
  // Duration myDuration = Duration(seconds: 10);
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
    // startTimer();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    state = Provider.of<DataProvider>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    // String strDigits(int n) => n.toString().padLeft(2, '0');
    // final minutes = strDigits(myDuration.inMinutes.remainder(60));
    // final seconds = strDigits(myDuration.inSeconds.remainder(60));
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
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
                  otpFieldStyle: OtpFieldStyle(
                    borderColor: Colors.black38,
                  ),
                  onChanged: (pin) {
                    otp = pin;
                  },
                  onCompleted: (pin) {
                    otp = pin;
                    completeRegistration();
                  },
                ),
              ),
              SizedBox(height: height / 10),
              Button(
                LanguageEn.verify,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  if (otp.length == 6) {
                    completeRegistration();
                  } else {
                    popup(context,
                        title: LanguageEn.alert,
                        message: LanguageEn.enterverification);
                  }
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future postUserInfo() async {
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
        'mobile': state.userInfo!.mobile,
        'mobileCountryCode': state.userInfo!.countryCode,
        'referrer': state.userInfo!.referrer,
        'pushNotificationToken': state.userInfo!.pushNotificationToken,
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

      return responseData;
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(context,
          title: LanguageEn.error,
          // TODO: Show somethingwentwrongerror for internal error
          // message: LanguageEn.somethingwentwrong);
          message: LanguageEn.errormessage + e.toString());
    }
  }

  void completeRegistration() async {
    Map responseData = await postUserInfo();

    if (responseData['statusCode'] == 200) {
      getUserInfo();

      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => const Complateerification(),
        ),
      );
    } else {
      popup(context,
          title: LanguageEn.error,
          message: LanguageEn.errormessage + responseData['data']['message']);
    }
  }

  getUserInfo() async {
    var publicKey = await StoreData().storeGetData('publicKey') ?? '';
    var secretKey = await StoreData().storeGetData('secretKey') ?? '';

    Map responseData = await makeGetRequest(
        uri:
            '/v1/users/${state.userInfo!.username!.trim().replaceAll(' ', '')}',
        signer: publicKey,
        publicKey: publicKey,
        secretKey: secretKey);

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200) {
      storeUserInfo(responseData['data']);
    } else if (responseData['statusCode'] == 404) {
      popup(context,
          title: LanguageEn.error,
          message: LanguageEn.errormessage + responseData['data']['message']);
    } else {
      // must be some sort of server error
      // let's throw it
      popup(context,
          title: LanguageEn.error,
          message: LanguageEn.errormessage + responseData['data']['message']);
    }
  }

  storeUserInfo(userInfoMap) async {
    print('this is userinfo map: ${userInfoMap}');
    var userInfo = userInfoMap['userData'] ?? {};
    var assetBalances = userInfoMap['assetBalances'] ?? {};
    var nftBalances = userInfoMap['nftBalances'] ?? {};
    var thirdPartyWalletAccess = userInfoMap['thirdPartyWalletAccess'] ?? [];
    var defaultAssets = userInfoMap['defaultAssets'] ?? [];

    await StoreData().storeInsertData('userInfo', userInfo);
    await StoreData().storeInsertData('assetBalances', assetBalances);
    await StoreData().storeInsertData('nftBalances', nftBalances);
    await StoreData()
        .storeInsertData('thirdPartyWalletAccess', thirdPartyWalletAccess);
    await StoreData().storeInsertData('defaultAssets', defaultAssets);
    await StoreData().storeInsertData('isFirstTime', false);
  }

  // void resendOTP() async {
  //   Map responseData = await postUserInfo();

  //   if (responseData['statusCode'] == 200) {
  //     popup(context,
  //         title: LanguageEn.success,
  //         message: LanguageEn.errormessage + responseData['data']['message']);
  //   } else {
  //     popup(context,
  //         title: LanguageEn.error,
  //         message: LanguageEn.errormessage + responseData['data']['message']);
  //   }

  //   resetTimer();
  //   startTimer();
  // }

//   void startTimer() {
//     countdownTimer =
//         Timer.periodic(Duration(seconds: 1), (_) => setCountDown());
//   }

// // Step 4
//   void stopTimer() {
//     setState(() {
//       countdownTimer!.cancel();
//       countdownTimer = null;
//     });
//   }

// // Step 5
//   void resetTimer() {
//     stopTimer();
//     setState(() => myDuration = Duration(seconds: 10));
//   }

// // Step 6
//   void setCountDown() {
//     final reduceSecondsBy = 1;
//     setState(() {
//       final seconds = myDuration.inSeconds - reduceSecondsBy;
//       if (seconds < 0) {
//         print('stopping...');
//         countdownTimer!.cancel();
//         print('stopped!');
//       } else {
//         myDuration = Duration(seconds: seconds);
//       }
//     });
//   }

  @override
  void dispose() {
    print('disposing...');
    // countdownTimer!.cancel();
    // countdownTimer = null;
    print('disposed');
    // TODO: implement dispose
    super.dispose();
  }
}
