import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
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
    print('this is viewData');
    print(viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
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
                    'Confirm request',
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
                'You are about to grant access to the following users',
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
              Padding(
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
                      Column(
                        children: [
                          SizedBox(
                            height: height / 50,
                          ),
                          Text(
                            'Wallet',
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
                                                appState.activeWallet!.alias ??
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
              if (viewData["viewers"].length > 0) ...[
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
                              'Viewer access',
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
                                                for (var i = 0;
                                                    i <
                                                        viewData['viewers']
                                                            .length;
                                                    i++) ...[
                                                  userItem(
                                                      viewData['viewers'][i],
                                                      notifier.getbluecolor)
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
              if (viewData['addApprovers'] == true) ...[
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
                              'Approver access',
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
                                                for (var i = 0;
                                                    i <
                                                        viewData['approvers']
                                                            .length;
                                                    i++) ...[
                                                  userItem(
                                                      viewData['approvers'][i],
                                                      notifier.getbluecolor),
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
                              'Initiator access',
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
                                                for (var i = 0;
                                                    i <
                                                        viewData['initiators']
                                                            .length;
                                                    i++) ...[
                                                  userItem(
                                                      viewData['initiators'][i],
                                                      notifier.getbluecolor)
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
                              'No. of approvals needed',
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
                                child: Text(
                                  viewData['noOfApprovalsNeeded'],
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
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  LanguageEn.password,
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
                  LanguageEn.authorizewithbiometrics,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  LanguageEn.authorize,
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
      sendDataToServer();
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer();
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
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  sendDataToServer() async {
    print('sending to server....');

    try {
      showLoader(context);

      var permissions = [];
      print(viewData);

      for (var i = 0; i < viewData['viewers'].length; i++) {
        print(viewData['viewers'][i]);
        permissions.add(
          {
            "targetUsername": viewData['viewers'][i],
            // "name": "",
            "permission": "VIEW-ONLY",
          },
        );
      }

      for (var i = 0; i < viewData['approvers'].length; i++) {
        print(viewData['approvers'][i]);
        permissions.add(
          {
            "targetUsername": viewData['approvers'][i],
            // "name": "",
            "permission": "APPROVER",
          },
        );
      }

      for (var i = 0; i < viewData['initiators'].length; i++) {
        print(viewData['initiators'][i]);
        permissions.add(
          {
            "targetUsername": viewData['initiators'][i],
            // "name": "",
            "permission": "INITIATOR",
          },
        );
      }

      var postData = {
        "numberOfApprovalsNeeded": int.parse(viewData['noOfApprovalsNeeded']),
        "permissions": permissions,
      };

      String requestBody = jsonEncode(postData);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/users/account',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );
      print(responseData);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;
        print('messagelenth: $messageLength');
        postProcessData(
            context, messageShown, messageLength, responseData['data'],
            callback: () {
          signAndSendToServerAgain(responseData['data']);
        });
        hideLoader(context);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }

  void signAndSendToServerAgain(responseFromServer) async {
    try {
      print('signing and sending....');
      showLoader(context);

      //sign the transaction and the submit again
      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0], // the primary wallet secret key,
        responseFromServer['transaction'],
        responseFromServer['networkPassPhrase'],
      );
      responseFromServer['transactionId'] = "";
      responseFromServer['transactionSignature'] = signature;
      print('second: ${responseFromServer}');

      String requestBody = jsonEncode(responseFromServer);

      print('second: ${requestBody}');

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/users/account',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );

      if (responseData['statusCode'] == 200) {
        print(responseData);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': 'Shared access enabled successfully',
          'message':
              'You have successfully enabled shared access on your wallet (${appState.activeWallet!.alias})!',
        };
        appState.currentAction =
            PageAction(state: PageState.replace, page: SuccessViewPageConfig);
        hideLoader(context);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }
}
