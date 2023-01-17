import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SetupAccountRecovery extends StatefulWidget {
  const SetupAccountRecovery({Key? key}) : super(key: key);

  @override
  State<SetupAccountRecovery> createState() => _SetupAccountRecoveryState();
}

class _SetupAccountRecoveryState extends State<SetupAccountRecovery> {
  late ColorNotifier notifier;
  bool isSwitched = false;
  late DataProvider appState;
  late Wallet primaryWallet;

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
    appState = Provider.of<DataProvider>(context, listen: false);
    primaryWallet = appState.userInfo!.wallets!
        .firstWhere((wallet) => wallet.primaryWallet == 1);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        resizeToAvoidBottomInset: false,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Text(
                LanguageEn.enable,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold),
              ),
              Text(
                LanguageEn.accountrecovery,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 45),
              // Center(
              //   child: Image.asset("assets/images/palm-recognition.png",
              //       height: height / 2.8),
              // ),
              description(LanguageEn.enableaccountrecoverydescription1),
              description(LanguageEn.enableaccountrecoverydescription2),
              description(LanguageEn.enableaccountrecoverydescription3),
              description(LanguageEn.enableaccountrecoverydescription4),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.taptoenableaccountrecovery,
                notifier.getbluecolor,
                wihitecolor,
                onTap: _handleSubmit,
              ),
              SizedBox(height: height / 15),
            ],
          ),
        ),
      ),
    );
  }

  Widget description(String desc) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                children: [
                  Container(
                    width: width / 1.3,
                    child: Text(
                      desc,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                          fontSize: 16,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody),
                    ),
                  ),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  _handleSubmit() async {
    showLoader(context);

    try {
      Map map = {
        "transaction": '',
        "transactionSignature": '',
        "transactionId": '',
        "networkPassPhrase": '',
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/account/recovery',
        body: requestBody,
        signer: primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: primaryWallet.publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }

  postProcessData(messageShown, messageLength, data) {
    print('messageShown: $messageShown messageLength $messageLength');
    // we would like to display all messages returned from the initial
    // request to server using a popup. In order to achieve that we
    // employ the use of a little recursion here. Please recursive
    // functions can turn into a nightmare fast so be carefull here.
    if (messageShown <= messageLength - 1) {
      showResponseMessage(
          context,
          data['messages'][messageShown],
          () => {
                print('postProcessData: $messageShown'),
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }
    sendFullDataToServer(data);
  }

  void sendFullDataToServer(responseBody) async {
    try {
      showLoader(context);
      // get primary signature
      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0],
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      print('this is primary sign: $signature');
      responseBody['transactionSignature'] = signature;
      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/account/recovery',
        body: requestBody,
        signer: primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: primaryWallet.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        await updateUserInfo(primaryWallet.signer!, appState.secretKeys[0],
            primaryWallet.publicKey!, appState.userInfo!.username, appState);
        appState.viewData = {
          SuccessViewPageConfig.key: {
            'title': LanguageEn.success,
            'message': LanguageEn.enableaccountrecoverysuccess,
          }
        };
        appState.currentAction =
            PageAction(state: PageState.replace, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }

    hideLoader(context);
  }
}
