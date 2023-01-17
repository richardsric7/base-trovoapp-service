import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:intl/intl.dart';
import 'package:trovo_wallet/custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
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
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class ApprovalDetails extends StatefulWidget {
  const ApprovalDetails({Key? key}) : super(key: key);

  @override
  State<ApprovalDetails> createState() => _ApprovalDetails();
}

class _ApprovalDetails extends State<ApprovalDetails>
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
  var rejectReason;

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
    viewData = appState.viewData![ApprovalDetailsViewPageConfig.key];

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
                    'Approve transaction request',
                    style: TextStyle(
                        fontSize: 20.sp,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
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
                    children: [
                      SizedBox(
                        width: width / 30,
                      ),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(
                            height: height / 50,
                          ),
                          displayInfo(key: 'Wallet', value: viewData['alias']),
                          displayInfo(
                              key: 'Transaction type',
                              value: viewData['transactionType']
                                  .toString()
                                  .capitalizeFirst!),
                          displayInfo(
                              key: 'Initiator', value: viewData['initiator']),
                          displayInfo(
                              key: 'Initiated',
                              value:
                                  '${DateFormat('yyyy-MM-dd hh:mm a').format(DateTime.parse(viewData['createdAt']))}'),
                          // displayInfo(
                          //     key: 'Description',
                          //     value: viewData['description']),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Description: ',
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: 5,
                              ),
                              Container(
                                width: width / 1.2,
                                child: Wrap(
                                  crossAxisAlignment: WrapCrossAlignment.center,
                                  children: [
                                    Text(
                                      viewData['description'],
                                      overflow: TextOverflow.visible,
                                      style: TextStyle(
                                        fontSize: 13,
                                        fontWeight: FontWeight.w400,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              SizedBox(
                                height: height / 50.0,
                              ),
                            ],
                          ),
                          displayInfo(
                              key: 'Approval status',
                              value:
                                  '${viewData['approvalsGotten'].toString()} out of ${viewData['approvalsNeeded'].toString()} approvals recieved'),
                          if (viewData['approvedBy'].toString().isNotEmpty) ...[
                            displayInfo(
                                key: 'Approved by',
                                value: '${viewData['approvedBy']}'),
                          ],
                          if (viewData['rejectedBy'].toString().isNotEmpty) ...[
                            displayInfo(
                                key: 'Rejected by',
                                value: '${viewData['rejectedBy']}'),
                            displayInfo(
                                key: 'Reason for rejection',
                                value: viewData['reasonForRejection']),
                          ],
                          displayInfo(
                              key: 'Transaction status',
                              value: viewData['transactionStatus']
                                  .toString()
                                  .capitalizeFirst!),
                          SizedBox(
                            height: height / 50,
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              // if transaction is still pending or transaction is not already
              // signed by current user (whether approved or rejected) then show
              // the approve or reject buttons
              if (viewData['transactionStatus'] == 'PENDING' &&
                  !(viewData['approvedBy']
                          .toString()
                          .contains(appState.userInfo!.username!) ||
                      viewData['rejectedBy']
                          .toString()
                          .contains(appState.userInfo!.username!))) ...[
                SizedBox(
                  height: height / 20,
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
                    'Approve with biometrics',
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: toggleSwitch,
                  ),
                ] else ...[
                  Button(
                    'Approve',
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: handleAuthorization,
                  ),
                ],
                SizedBox(
                  height: height / 50,
                ),
                if (appState.biometricEnabled && password.isEmpty) ...[
                  ButtonOutlined(
                    'Reject with biometrics',
                    notifier.getwihitecolor,
                    notifier.getbluewhitecolor,
                    onTap: () {
                      rejectionReasonPopup(context, (reason) {
                        setState(() {
                          rejectReason = reason;
                        });
                        toggleSwitch(approve: false);
                      });
                    },
                  ),
                ] else ...[
                  ButtonOutlined(
                    'Reject',
                    notifier.getwihitecolor,
                    notifier.getbluewhitecolor,
                    onTap: () {
                      rejectionReasonPopup(context, (reason) {
                        setState(() {
                          rejectReason = reason;
                        });
                        handleAuthorization(approve: false);
                      });
                    },
                  ),
                ],
              ],
              SizedBox(
                height: height / 15,
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

  Widget displayInfo({required String key, required String value}) {
    return Column(
      children: [
        Container(
          width: width / 1.2,
          child: Wrap(
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              Text(
                '$key: ',
                style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w700,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold),
              ),
              Text(
                value,
                overflow: TextOverflow.visible,
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w400,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
            ],
          ),
        ),
        SizedBox(
          height: height / 50.0,
        ),
      ],
    );
  }

  String getHeadlineLabel(transactionStatus) {
    switch (transactionStatus) {
      case 'PENDING':
        return (viewData['approvedBy']
                    .toString()
                    .contains(appState.userInfo!.username!) ||
                viewData['rejectedBy']
                    .toString()
                    .contains(appState.userInfo!.username!))
            ? 'You have already signed this transaction '
            : 'Your approval has been requested for this transaction';
      case 'REJECTED':
        return 'This transaction has already been rejected';
      case 'COMPLETED':
        return 'This transaction has been completed';
      default:
        return '';
    }
  }

  void handleAuthorization({bool approve = true}) {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      if (!approve) {
        sendRejectToServer();
        return;
      }
      sendDataToServer();
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }

  void toggleSwitch({bool approve = true}) async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        if (approve) {
          sendDataToServer();
        } else {
          sendRejectToServer();
        }
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

      String requestBody = jsonEncode({});

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/approval/${viewData['id']}',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );
      print('responseData: ${responseData}');

      if (responseData['statusCode'] == 202) {
        sendDataToServerAgain(responseData['data']);
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

  sendDataToServerAgain(Map data) async {
    print('sending to server again....');

    try {
      showLoader(context);
      // sign transaction
      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0], // the primary wallet secret key,
        data['transaction'],
        data['networkPassPhrase'],
      );

      data['transactionSignature'] = signature;

      String requestBody = jsonEncode(data);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/shared-access/approval/${viewData['id']}',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );

      print('responseData: ${responseData}');
      if (responseData['statusCode'] == 200) {
        appState.getApprovals();
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': 'Transaction approval submitted',
          'message':
              'You have successfully submitted your own approval for this transaction.',
          'useOnDone': true,
          'onDone': () {
            appState.currentAction = PageAction(
                state: PageState.replace, page: SharedAccessViewPageConfig);
          },
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

  sendRejectToServer() async {
    print('sending to server....');

    try {
      showLoader(context);

      String requestBody = jsonEncode({
        'rejectionReason': rejectReason,
      });

      print(requestBody);

      Map responseData = await makeDeleteRequest(
        uri: '/v1/shared-access/approval/${viewData['id']}',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );
      print('responseData: ${responseData}');

      if (responseData['statusCode'] == 200) {
        appState.getApprovals();
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': 'Rejection submitted',
          'message': 'You have successfully rejected this transaction.',
          'useOnDone': true,
          'onDone': () {
            appState.currentAction = PageAction(
                state: PageState.replace, page: SharedAccessViewPageConfig);
          },
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

  Future<void> updateUserInfo() async {
    Map responseData = await makeGetRequest(
      uri:
          '/v1/users/${appState.userInfo!.username!.trim().replaceAll(' ', '')}',
      signer: activeWallet!.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: activeWallet!.publicKey!,
    );

    print('secretkey: ${appState.secretKeys[0]}');

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200) {
      await storeUserInfo(responseData['data'], appState);
    }
  }
}
