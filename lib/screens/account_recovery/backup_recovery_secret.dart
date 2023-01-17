import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/secret.dart';

class BackupRecoverySecret extends StatefulWidget {
  const BackupRecoverySecret({Key? key}) : super(key: key);

  @override
  State<BackupRecoverySecret> createState() => _BackupRecoverySecretState();
}

class _BackupRecoverySecretState extends State<BackupRecoverySecret> {
  late DataProvider state;
  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    state = Provider.of<DataProvider>(context, listen: true);
    var data = state.viewData![EnsurePrivacyPageConfig.key];

    print(data);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Text(
                LanguageEn.backupwallet,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  LanguageEn.writeitdown,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Secret(
                  state.tempUsername, state.tempSecretKey, state.tempPublicKey),
              if (data != null &&
                  data['rel'] != 'restoreUnactivatedAccount') ...[
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
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
                        side: BorderSide(color: notifier.getbluewhitecolor),
                        value: state.tempInvalidateOldSigner,
                        onChanged: (value) {
                          state.setTempInvalidateOldSigner = value;
                        },
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.all(16.0),
                      width: width / 1.2,
                      child: Text(
                        LanguageEn.invalidateoldsigner,
                        style: TextStyle(
                            fontSize: height / 55,
                            color: notifier.getgrey,
                            fontFamily: fontbody),
                      ),
                    )
                  ],
                ),
              ],
              SizedBox(height: height / 20),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  if (data != null &&
                      data['rel'] == 'restoreUnactivatedAccount') {
                    restoreInactiveAccountRequest();
                  } else {
                    requestAccountRecovery();
                  }
                },
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

  restoreInactiveAccountRequest() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "newSignerPublicKey": state.tempPublicKey,
        "emailOtp": state.tempEmailOtp,
        "username": state.tempUsername,
        "securityAnswers": state.tempSecurityQuestionsAndAnswers,
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/inactive-account/recover',
        body: requestBody,
        signer: state.tempPublicKey,
        secretKey: state.tempSecretKey, // the primary wallet secret key
        publicKey: state.tempPublicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        state.currentAction = PageAction(
            state: PageState.addPage,
            page: AccountRecoverySuccessViewPageConfig);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  requestAccountRecovery() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "newSignerPublicKey": state.tempPublicKey,
        "disableOldSignerFromPrimaryWallet":
            state.tempInvalidateOldSigner ? 1 : 0,
        "commit": 0,
        "emailOtp": state.tempEmailOtp,
        "username": state.tempUsername,
        "transactionId": "",
        "securityAnswers": state.tempSecurityQuestionsAndAnswers,
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/account/recover',
        body: requestBody,
        signer: state.tempPublicKey,
        secretKey: state.tempSecretKey, // the primary wallet secret key
        publicKey: state.tempPublicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        state.viewData![CompleteAccountRecoveryViewPageConfig.key] =
            responseData['data'];
        state.currentAction = PageAction(
            state: PageState.addPage,
            page: CompleteAccountRecoveryViewPageConfig);
      } else {
        hideLoader(context);
        popup(context,
            title: LanguageEn.error, message: responseData['data']['error']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }
}
