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
import 'package:trovo_wallet/models/permission.dart';
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

class UpdateSharedAccessDetails extends StatefulWidget {
  const UpdateSharedAccessDetails({Key? key}) : super(key: key);

  @override
  State<UpdateSharedAccessDetails> createState() =>
      _UpdateSharedAccessDetails();
}

class _UpdateSharedAccessDetails extends State<UpdateSharedAccessDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  var viewData;
  String sharedAccessModifySuccess =
      'Shared access modification on wallet [alias] was successful.';
  String sharedAccessModifyRequestSuccess =
      'Your request to modify shared access on wallet [alias] has been successfully submitted! This transaction will be completed when it gets the required number of approvals.';

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet =
        appState.userInfo!.getWallet(appState.viewData!['walletPublicKey']);
    viewData = appState.viewData;
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

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
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                child: Text(
                  "youareabouttomodify".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
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
                            "wallet".tr(),
                            style: TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w700,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold),
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
                                                wallet.alias!,
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
                              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                              child: Container(
                                width: width / 1.27,
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    for (var i = 0;
                                        i < viewData['viewers'].length;
                                        i++) ...[
                                      getPermissionInfo(viewData['viewers'][i]),
                                    ]
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
                              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                              child: Container(
                                width: width / 1.27,
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    for (var i = 0;
                                        i < viewData['approvers'].length;
                                        i++) ...[
                                      getPermissionInfo(
                                          viewData['approvers'][i]),
                                    ]
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
                      // mainAxisAlignment: MainAxisAlignment.center,
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
                              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                              child: Container(
                                width: width / 1.27,
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    for (var i = 0;
                                        i < viewData['initiators'].length;
                                        i++) ...[
                                      getPermissionInfo(
                                          viewData['initiators'][i]),
                                    ]
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
                              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                              child: Container(
                                child: Text(
                                  '${viewData['noOfApprovalsNeeded']}/${viewData['noOfApprovers']}',
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
      sendDataToServer();
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
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
    if (value!.isEmpty) return "pleaseenteryourpassword".tr();

    if (value.length < 6) return "use6charsormoreforpassword".tr();

    return null;
  }

  Widget getPermissionInfo(permission) {
    var text = '';
    Color color;
    switch (permission.permissionState) {
      case PermissionState.Added:
        text = "added".tr();
        color = notifier.getgreencolor;
        break;
      case PermissionState.Revoked:
        text = "revoked".tr();
        color = Colors.red;
        break;
      case PermissionState.Modified:
        text = "modified".tr();
        color = Colors.yellow;
        break;
      default:
        text = "active".tr();
        color = notifier.getbluewhitecolor;
    }
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            width: width / 1.75,
            child: Text(
              '${permission.targetUsername} [${permission.fullName}]',
              style: TextStyle(
                  fontSize: 13,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold),
            ),
          ),
          Text(
            text,
            style: TextStyle(fontSize: 13, color: color, fontFamily: fontbody),
          ),
        ],
      ),
    );
  }

  sendDataToServer() async {
    print('sending to server....');

    try {
      showLoader(context);

      var addedPermissions = [];
      var modifiedPermissions = [];
      var revokedPermissions = [];
      print(viewData);

      for (var i = 0; i < viewData['viewers'].length; i++) {
        if (viewData['viewers'][i].permissionState == PermissionState.Added) {
          addedPermissions.add(
            {
              "targetUsername": viewData['viewers'][i].targetUsername,
              // "name": "",
              "permission": "VIEW-ONLY",
            },
          );
        }

        if (viewData['viewers'][i].permissionState ==
            PermissionState.Modified) {
          modifiedPermissions.add(
            {
              "targetUsername": viewData['viewers'][i].targetUsername,
              // "name": "",
              "permission": "VIEW-ONLY",
            },
          );
        }

        if (viewData['viewers'][i].permissionState == PermissionState.Revoked) {
          revokedPermissions.add(
            {
              "targetUsername": viewData['viewers'][i].targetUsername,
              // "name": "",
              "permission": "VIEW-ONLY",
            },
          );
        }
      }

      for (var i = 0; i < viewData['approvers'].length; i++) {
        if (viewData['approvers'][i].permissionState == PermissionState.Added) {
          addedPermissions.add(
            {
              "targetUsername": viewData['approvers'][i].targetUsername,
              // "name": "",
              "permission": "APPROVER",
            },
          );
        }

        if (viewData['approvers'][i].permissionState ==
            PermissionState.Modified) {
          modifiedPermissions.add(
            {
              "targetUsername": viewData['approvers'][i].targetUsername,
              // "name": "",
              "permission": "APPROVER",
            },
          );
        }

        if (viewData['approvers'][i].permissionState ==
            PermissionState.Revoked) {
          revokedPermissions.add(
            {
              "targetUsername": viewData['approvers'][i].targetUsername,
              // "name": "",
              "permission": "APPROVER",
            },
          );
        }
      }

      for (var i = 0; i < viewData['initiators'].length; i++) {
        if (viewData['initiators'][i].permissionState ==
            PermissionState.Added) {
          addedPermissions.add(
            {
              "targetUsername": viewData['initiators'][i].targetUsername,
              // "name": "",
              "permission": "INITIATOR",
            },
          );
        }

        if (viewData['initiators'][i].permissionState ==
            PermissionState.Modified) {
          modifiedPermissions.add(
            {
              "targetUsername": viewData['initiators'][i].targetUsername,
              // "name": "",
              "permission": "INITIATOR",
            },
          );
        }

        if (viewData['initiators'][i].permissionState ==
            PermissionState.Revoked) {
          revokedPermissions.add(
            {
              "targetUsername": viewData['initiators'][i].targetUsername,
              // "name": "",
              "permission": "INITIATOR",
            },
          );
        }
      }

      var postData = {
        "numberOfApprovalsNeeded": viewData['noOfApprovalsNeeded'],
        "revokedPermissions": revokedPermissions,
        "addedPermissions": addedPermissions,
      };

      String requestBody = jsonEncode(postData);

      print(requestBody);

      Map responseData = await makePutRequest(
        uri: '/v1/shared-access/users/account',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );
      print(responseData);

      if (responseData['statusCode'] == 200 ||
          responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;
        print('messagelenth: $messageLength');
        postProcessData(
            context, messageShown, messageLength, responseData['data'],
            callback: () {
          signAndCommitTransaction(responseData['data']);
        });
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
      hideLoader(context);
    }
  }

  void signAndCommitTransaction(responseFromServer) async {
    try {
      print('signing and sending....');
      showLoader(context);

      // wallets with only shared view-only access are still not fully shared-wallets
      // the owner can still solely initiate and complete transactions. So in
      // the case that the wallet owner wants to now modify shared access on thier
      // wallet the can still sign the transaction. If however, the wallet has
      // shared approver and initiator access, then they cannot sign the transaction
      // rather they set commit=1
      if (responseFromServer['signatureRequired'] == 1) {
        responseFromServer['commit'] = 0;
        //sign the transaction and the submit again
        var signature = TrovoWalletSDK().signBase64Txn(
          appState.secretKeys[0], // the primary wallet secret key,
          responseFromServer['transaction'],
          responseFromServer['networkPassPhrase'],
        );
        responseFromServer['transactionId'] = "";
        responseFromServer['transactionSignature'] = signature;
      } else {
        responseFromServer['commit'] = 1;
      }
      print('second: ${responseFromServer}');

      String requestBody = jsonEncode(responseFromServer);

      print('second: ${requestBody}');

      Map responseData = await makePutRequest(
        uri: '/v1/shared-access/users/account',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      if (responseData['statusCode'] == 200 ||
          responseData['statusCode'] == 202) {
        updateUserInfo(
          appState.primaryWallet.signer!,
          appState.secretKeys[0],
          appState.primaryWallet.publicKey,
          appState.userInfo!.username,
          appState,
          forceRefresh: true,
        );
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': 'Request successfully submitted',
          'message': wallet.isPrimaryWallet || wallet.walletThreshold! < 2
              ? sharedAccessModifySuccess.replaceAll('alias', wallet.alias!)
              : sharedAccessModifyRequestSuccess.replaceAll(
                  'alias', wallet.alias!),
          'useOnDone': true,
          'onDone': () {
            appState.currentAction =
                PageAction(state: PageState.addAll, pages: [
              BottomHomePageConfig,
              SharedAccessViewPageConfig,
            ]);
          },
        };
        appState.currentAction =
            PageAction(state: PageState.replace, page: SuccessViewPageConfig);
        hideLoader(context);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
      hideLoader(context);
    }
  }
}
