import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class AddSharedAccessDetails extends StatefulWidget {
  const AddSharedAccessDetails({Key? key}) : super(key: key);

  @override
  State<AddSharedAccessDetails> createState() => _AddSharedAccessDetails();
}

class _AddSharedAccessDetails extends State<AddSharedAccessDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  Wallet? activeWallet;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  late Account primaryWalletKeyPair;
  var viewData;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    activeWallet = appState.activeWallet;
    viewData = appState.viewData![AddSharedAccessDetailsViewPageConfig.key];
    print(viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
                context, notifier.getwihitecolor, "", notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "confirmrequest".tr(),
                    style: TextStyle(
                        fontSize: 20.sp,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: height / 20,
              ),
              Text(
                "abouttorequestaccess".tr(),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              for (var i = 0; i < viewData.length; i++) ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Column(
                          children: [
                            SizedBox(
                              height: height / 50,
                            ),
                            Text(
                              "wallet".tr(),
                              style: TextStyle(
                                  fontSize: 15,
                                  fontWeight: FontWeight.w700,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold),
                            ),
                            SizedBox(
                              height: height / 50,
                            ),
                            Padding(
                              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                              child: Container(
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Column(
                                      children: [
                                        Container(
                                            width: width / 1.3,
                                            child: Wrap(
                                              alignment: WrapAlignment.center,
                                              children: [
                                                Text(
                                                  viewData[i]['wallet'].alias ??
                                                      "",
                                                  style: TextStyle(
                                                      fontSize: 15,
                                                      color: notifier
                                                          .getbluewhitecolor,
                                                      fontFamily: fontbody),
                                                )
                                              ],
                                            )),
                                        SizedBox(height: 2),
                                      ],
                                    ),
                                  ],
                                ),
                              ),
                            ),
                            SizedBox(height: height / 50.0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                if (viewData[i]["viewers"].length > 0) ...[
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Column(
                            children: [
                              SizedBox(
                                height: height / 50,
                              ),
                              Text(
                                "vieweraccess".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: height / 50,
                              ),
                              Padding(
                                padding:
                                    const EdgeInsets.fromLTRB(20, 10, 20, 3),
                                child: Container(
                                  child: Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Column(
                                        children: [
                                          Container(
                                              width: width / 1.3,
                                              child: Wrap(
                                                alignment: WrapAlignment.center,
                                                children: [
                                                  for (var k = 0;
                                                      k <
                                                          viewData[i]['viewers']
                                                              .length;
                                                      k++) ...[
                                                    userItem(
                                                        '${viewData[i]['viewers'][k]} [${viewData[i]['userFullnames'][viewData[i]['viewers'][k]]}]',
                                                        null,
                                                        foreColor: wihitecolor,
                                                        backColor: notifier
                                                            .getbluebackcolor)
                                                  ],
                                                ],
                                              )),
                                          SizedBox(height: 2),
                                        ],
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                              SizedBox(
                                height: height / 50.0,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
                SizedBox(
                  height: height / 50,
                ),
                if (viewData[i]['addApprovers'] == true) ...[
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Column(
                            children: [
                              SizedBox(
                                height: height / 50,
                              ),
                              Text(
                                "approveraccess".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: height / 50,
                              ),
                              Padding(
                                padding:
                                    const EdgeInsets.fromLTRB(20, 10, 20, 3),
                                child: Container(
                                  child: Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Column(
                                        children: [
                                          Container(
                                              width: width / 1.3,
                                              child: Wrap(
                                                alignment: WrapAlignment.center,
                                                children: [
                                                  for (var k = 0;
                                                      k <
                                                          viewData[i]
                                                                  ['approvers']
                                                              .length;
                                                      k++) ...[
                                                    userItem(
                                                        '${viewData[i]['approvers'][k]} [${viewData[i]['userFullnames'][viewData[i]['approvers'][k]]}]',
                                                        null,
                                                        foreColor: wihitecolor,
                                                        backColor: notifier
                                                            .getbluebackcolor),
                                                  ],
                                                ],
                                              )),
                                          SizedBox(height: 2),
                                        ],
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                              SizedBox(
                                height: height / 50.0,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Column(
                            children: [
                              SizedBox(
                                height: height / 50,
                              ),
                              Text(
                                "initiatoraccess".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: height / 50,
                              ),
                              Padding(
                                padding:
                                    const EdgeInsets.fromLTRB(20, 10, 20, 3),
                                child: Container(
                                  child: Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Column(
                                        children: [
                                          Container(
                                              width: width / 1.3,
                                              child: Wrap(
                                                alignment: WrapAlignment.center,
                                                children: [
                                                  for (var k = 0;
                                                      k <
                                                          viewData[i]
                                                                  ['initiators']
                                                              .length;
                                                      k++) ...[
                                                    userItem(
                                                      '${viewData[i]['initiators'][k]} [${viewData[i]['userFullnames'][viewData[i]['initiators'][k]]}]',
                                                      null,
                                                      foreColor: wihitecolor,
                                                      backColor: notifier
                                                          .getbluebackcolor,
                                                    )
                                                  ],
                                                ],
                                              )),
                                          SizedBox(height: 2),
                                        ],
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                              SizedBox(
                                height: height / 50.0,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Column(
                            children: [
                              SizedBox(
                                height: height / 50,
                              ),
                              Text(
                                "noofapprovalsrequired".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: height / 50,
                              ),
                              Padding(
                                padding:
                                    const EdgeInsets.fromLTRB(20, 10, 20, 3),
                                child: Container(
                                  child: Text(
                                    '${viewData[i]['noOfApprovalsNeeded']}/${viewData[i]['noOfApprovers']}',
                                    style: TextStyle(
                                        fontSize: 15,
                                        fontWeight: FontWeight.w700,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold),
                                  ),
                                ),
                              ),
                              SizedBox(
                                height: height / 50.0,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
                SizedBox(
                  height: height / 50.0,
                ),
              ],
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  "password".tr(),
                  notifier.getbluewhitecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  validator: validatePassword,
                  onChanged: (value) {
                    setState(() {
                      password = value!.trim().replaceAll(' ', '');
                    });
                  },
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  "authorizewithbiometrics".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  "authorize".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: handleAuthorization,
                ),
              ],
              SizedBox(
                height: height / 20,
              ),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  void handleAuthorization() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer(0);
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer(0);
        // aparently we need the code below to make the
        // screen updata to show loader
        // after authorizing with biometrics
        setState(() {});
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return "pleaseenteryourpassword".tr();

    if (value.length < 6) return "use6charsormoreforpassword".tr();

    return null;
  }

  sendDataToServer(int index) async {
    print('sending to server....');

    try {
      showLoader(context);

      var permissions = [];
      print(viewData);

      for (var i = 0; i < viewData[index]['viewers'].length; i++) {
        print(viewData[index]['viewers'][i]);
        permissions.add(
          {
            "targetUsername": viewData[index]['viewers'][i],
            // "name": "",
            "permission": "VIEW-ONLY",
          },
        );
      }

      if (viewData[index]['addApprovers'] == true) {
        for (var i = 0; i < viewData[index]['approvers'].length; i++) {
          print(viewData[index]['approvers'][i]);
          permissions.add(
            {
              "targetUsername": viewData[index]['approvers'][i],
              // "name": "",
              "permission": "APPROVER",
            },
          );
        }

        for (var i = 0; i < viewData[index]['initiators'].length; i++) {
          print(viewData[index]['initiators'][i]);
          permissions.add(
            {
              "targetUsername": viewData[index]['initiators'][i],
              // "name": "",
              "permission": "INITIATOR",
            },
          );
        }
      }

      var postData = {};

      if (viewData[index]['addApprovers']) {
        postData = {
          "numberOfApprovalsNeeded": viewData[index]['noOfApprovalsNeeded'],
          "permissions": permissions,
        };
      } else {
        postData = {
          "permissions": permissions,
        };
      }

      String requestBody = jsonEncode(postData);

      print(requestBody);

      var wallet = viewData[index]['wallet'] as Wallet;

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/users/account',
        body: requestBody,
        signer: wallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );
      // print(responseData);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;
        hideLoader(context);
        postProcessData(
            context, messageShown, messageLength, responseData['data'],
            callback: () {
          signAndSendToServerAgain(responseData['data'], index);
        });
      } else {
        hideLoader(context);
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  void signAndSendToServerAgain(responseFromServer, int index) async {
    try {
      showLoader(context);

      //sign the transaction and the submit again
      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0], // the primary wallet secret key,
        responseFromServer['transaction'],
        responseFromServer['networkPassPhrase'],
      );
      responseFromServer['transactionId'] = "";
      responseFromServer['transactionSignature'] = signature;

      String requestBody = jsonEncode(responseFromServer);

      // print('second: ${requestBody}');
      var wallet = viewData[index]['wallet'] as Wallet;

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/users/account',
        body: requestBody,
        signer: wallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      if (responseData['statusCode'] == 200) {
        if (viewData.length > 1 && index == 0) {
          sendDataToServer(1);
          return;
        }

        await updateUserInfo(
          appState.primaryWallet.signer!,
          appState.secretKeys[0],
          appState.primaryWallet.publicKey,
          appState.userInfo!.username,
          appState,
          forceRefresh: true,
        );
        var walletAliases = viewData.length > 1
            ? '${viewData[0]['wallet'].alias} and ${wallet.alias}'
            : wallet.alias!;
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': "sharedaccessenabledsuccessfully".tr(),
          'message':
              "sharedaccessenabledsuccessfully2".tr(args: [walletAliases]),
          'useOnDone': true,
          'onDone': () {
            print('ondone fired!');
            if ((appState.returnView != null &&
                    appState.returnView!.pages != null) &&
                appState.returnView!.pages!
                    .contains(WalletPreparationViewPageConfig)) {
              appState.currentAction =
                  PageAction(state: PageState.addAll, pages: [
                BottomHomePageConfig,
                WalletPreparationViewPageConfig,
              ]);
              appState.setActiveTokenizationWalletPublicKey = null;
              appState.viewData = {};
              appState.backupSecrets.clear();
              appState.returnView = null;
            } else if (appState.backupSecrets.length > 1) {
              appState.currentAction =
                  PageAction(state: PageState.addAll, pages: [
                BottomHomePageConfig,
              ]);
              appState.clearAccessList = true;
              appState.backupSecrets.clear();
            } else {
              appState.currentAction =
                  PageAction(state: PageState.addAll, pages: [
                BottomHomePageConfig,
                SharedAccessViewPageConfig,
              ]);
              appState.clearAccessList = true;
              WidgetsBinding.instance.addPostFrameCallback((_) {
                appState.sharedAccesstabController.animateTo(0,
                    duration: Duration(milliseconds: 500),
                    curve: Curves.easeInOut);
              });
            }
          },
        };
        hideLoader(context);
        appState.currentAction =
            PageAction(state: PageState.replace, page: SuccessViewPageConfig);
      } else {
        hideLoader(context);
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}
