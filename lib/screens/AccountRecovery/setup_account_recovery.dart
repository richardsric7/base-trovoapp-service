import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../router/PageActions.dart';
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
            // crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                LanguageEn.setup,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluecolor,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold),
              ),
              Text(
                LanguageEn.accountrecovery,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluecolor80,
                    fontSize: 30.sp,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 45),
              Center(
                child: Text(
                  LanguageEn.unlockfinger,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontSize: 16.sp,
                      color: notifier.getgrey,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              Center(
                child: Image.asset("assets/images/palm-recognition.png",
                    height: height / 2.8),
              ),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.taptoenableaccountrecovery,
                notifier.getbluecolor,
                wihitecolor,
                onTap: _handleSubmit,
              )
            ],
          ),
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

      if (responseData['statusCode'] == 200) {
        sendFullDataToServer(responseData['data']);
        // print('sending full data to server.........');
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
        showSuccessAlert(context, onTap: () {
          appState.currentAction = PageAction(
              state: PageState.replaceAll, page: BottomHomePageConfig);
        });
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
