import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DeleteAccount extends StatefulWidget {
  const DeleteAccount({Key? key}) : super(key: key);

  @override
  State<DeleteAccount> createState() => _DeleteAccountState();
}

class _DeleteAccountState extends State<DeleteAccount> {
  late ColorNotifier notifier;
  late DataProvider appState;

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
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
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
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Text(
                "deleteaccount".tr(),
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                  fontSize: 24,
                ),
              ),
              Image.asset('assets/images/deleted.png'),
              deletionInfo(),
              SizedBox(
                height: height / 50,
              ),
              Button(
                "yesdeleteaccount".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  confirmAccountDeletionPopup(context,
                      onConfirmationSuccess: () {
                    requestAccountDeletion();
                  });
                },
              ),
              SizedBox(
                height: height / 50,
              ),
              ButtonOutlined(
                "nokeepaccount".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  Navigator.of(context).pop();
                },
              ),
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget deletionInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: Colors.red[50],
        ),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0, vertical: 20.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            children: [
              Text(
                "abouttodeleteaccount".tr(),
                style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold),
              ),
              SizedBox(
                height: height / 90,
              ),
              Text(
                "beforeyoudeleteaccount".tr(),
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w800,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Image.asset(
                    "assets/images/jam_alert.png",
                  ),
                  Spacer(),
                  Container(
                    width: width / 1.4,
                    child: Text(
                      "deletioninfo1".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Image.asset(
                    "assets/images/jam_alert.png",
                  ),
                  Spacer(),
                  Container(
                    width: width / 1.4,
                    child: Text(
                      "deletioninfo2".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Text(
                "doyouwanttoproceed".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w800,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 25,
              ),
            ],
          ),
        ),
      ),
    );
  }

  void requestAccountDeletion() async {
    print('sending request to delete user account...');
    try {
      showLoader(context);
      // make initial request to the server using empty body
      Map map = {};
      String requestBody = jsonEncode(map);
      print(appState.primaryWallet.signer);
      print(appState.primaryWallet.publicKey);
      print(requestBody);
      Map responseData = await makeDeleteRequest(
        uri: '/v1/users',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: appState.primaryWallet.publicKey!,
      );

      print('response: $responseData');
      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        // get primary signature
        var signature = TrovoWalletSDK().signBase64Txn(
          appState.secretKeys[0],
          responseData['data']['transaction'],
          responseData['data']['networkPassPhrase'],
        );
        responseData['transactionSignature'] = signature;

        String req = jsonEncode(responseData);

        Map res = await makeDeleteRequest(
          uri: '/v1/users',
          body: req,
          signer: appState.primaryWallet.signer!,
          secretKey: appState.secretKeys[0],
          publicKey: appState.primaryWallet.publicKey!,
        );

        if (responseData['statusCode'] == 200) {
          print('res is here ============> $res');
          StoreData().storeDeleteData();
          appState.viewData![SuccessViewPageConfig.key] = {
            'title': 'Request successfull',
            'message': 'Your request has been successfully submitted.',
            'useOnDone': true,
            'onDone': () {
              appState.currentAction = PageAction(
                state: PageState.replaceAll,
                page: GetStartedViewPageConfig,
              );
            },
          };
          appState.currentAction =
              PageAction(state: PageState.replace, page: SuccessViewPageConfig);
        } else {
          popup(context,
              title: "error".tr(), message: responseData['data']['error']);
        }
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }
    hideLoader(context);
  }
}
