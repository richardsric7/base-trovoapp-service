import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DisableAccountRecovery extends StatefulWidget {
  const DisableAccountRecovery({Key? key}) : super(key: key);

  @override
  State<DisableAccountRecovery> createState() => _DisableAccountRecovery();
}

class _DisableAccountRecovery extends State<DisableAccountRecovery> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  String password = '';
  late var primaryWallet;
  late Future<Map> securityQuestionsMap;
  var questionsMap = {
    1: {
      "q": "",
      "a": "",
      "e": false,
    },
    2: {
      "q": "",
      "a": "",
      "e": false,
    },
    3: {
      "q": "",
      "a": "",
      "e": false,
    },
  };

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
    securityQuestionsMap = fetchQuestions(
        primaryWallet.signer,
        appState.secretKeys[0],
        primaryWallet.publicKey,
        appState.userInfo!.username);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 20),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 50),
                Text(
                  LanguageEn.disable,
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
                SizedBox(height: height / 20),
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
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20.0, vertical: 15.0),
                          child: Column(
                            children: [
                              Container(
                                width: width / 1.3,
                                child: Text(
                                  LanguageEn
                                      .answersecurityquestionstodisableaccountrecovery,
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
                ),
                SizedBox(height: height / 50),
                FutureBuilder<Map>(
                    future: securityQuestionsMap,
                    builder: (context, snapshot) {
                      if (snapshot.connectionState == ConnectionState.waiting) {
                        return CircularProgressIndicator(
                          backgroundColor: notifier.getbluecolor,
                          valueColor: new AlwaysStoppedAnimation<Color>(
                            notifier.getgreencolor,
                          ),
                          strokeWidth: 3.0,
                        );
                      } else if (snapshot.connectionState ==
                          ConnectionState.done) {
                        if (snapshot.hasError) {
                          return Padding(
                            padding: const EdgeInsets.all(8.0),
                            child: Column(
                              children: [
                                Text(
                                  LanguageEn.somethingwentwrong,
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
                                ),
                                ElevatedButton(
                                  onPressed: () {
                                    setState(() {
                                      securityQuestionsMap = fetchQuestions(
                                          primaryWallet.signer,
                                          appState.secretKeys[0],
                                          primaryWallet.publicKey,
                                          appState.userInfo!.username);
                                    });
                                  },
                                  style: ButtonStyle(
                                    backgroundColor:
                                        MaterialStateProperty.all<Color>(
                                            notifier.getbluecolor!),
                                  ),
                                  child: Text(
                                    LanguageEn.retry,
                                    style: TextStyle(
                                      fontFamily: fontsemibold,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          );
                        } else if (snapshot.hasData) {
                          var securityQuestions =
                              snapshot.data!['securityQuestions'];
                          var userSecurityAnswers =
                              snapshot.data!['userSecurityAnswers'];
                          return Column(
                            children: [
                              for (var i = 0;
                                  i < securityQuestions!.length;
                                  i++) ...[
                                if (securityQuestions[i]['ID'] ==
                                    userSecurityAnswers['q1']) ...[
                                  questionView(securityQuestions[i]['Question'],
                                      securityQuestions[i]['ID'], 1),
                                ],
                                if (securityQuestions[i]['ID'] ==
                                    userSecurityAnswers['q2']) ...[
                                  questionView(securityQuestions[i]['Question'],
                                      securityQuestions[i]['ID'], 2),
                                ],
                                if (securityQuestions[i]['ID'] ==
                                    userSecurityAnswers['q3']) ...[
                                  questionView(securityQuestions[i]['Question'],
                                      securityQuestions[i]['ID'], 3),
                                ]
                              ]
                            ],
                          );
                        } else {
                          return const Text('Empty data');
                        }
                      } else {
                        return Text('State: ${snapshot.connectionState}');
                      }
                    }),
                SizedBox(height: height / 20),
                Button(
                  '${LanguageEn.disable} ${LanguageEn.accountrecovery}',
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    trySubmit();
                  },
                ),
                SizedBox(height: height / 10),
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget questionView(question, int qNumber, int rel) {
    questionsMap[rel]!['q'] = qNumber;
    return Column(
      children: [
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 30.0),
              child: Container(
                width: width / 1.2,
                child: Text(
                  '$question',
                  style: TextStyle(
                      fontSize: 16,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
            ),
          ],
        ),
        SizedBox(height: height / 50),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: CustomTextFormField.textField(
            LanguageEn.enteranswer,
            notifier.getbluecolor,
            Icons.question_answer_outlined,
            notifier.getgrey,
            notifier.getprefixicon,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            // validator: validateEmail,
            onSaved: (value) {
              questionsMap[rel]!['a'] = value.toString().trim();
              print('email: $questionsMap');
            },
            validator: (value) {
              if (value.toString().isEmpty) {
                return 'Please enter answer to the question';
              }
              return null;
            },
            keyboardtype: TextInputType.text,
          ),
        ),
      ],
    );
  }

  void trySubmit() async {
    final form = _formKey.currentState;
    if (!form!.validate()) {
      return;
    }

    form.save();
    sendToServer();
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

  sendToServer() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "transaction": '',
        "transactionSignature": '',
        "transactionId": '',
        "networkPassPhrase": '',
        "securityAnswers": {
          "q1": int.parse(questionsMap[1]!['q'].toString()),
          "a1": questionsMap[1]!['a'],
          "q2": int.parse(questionsMap[2]!['q'].toString()),
          "a2": questionsMap[2]!['a'],
          "q3": int.parse(questionsMap[3]!['q'].toString()),
          "a3": questionsMap[3]!['a']
        },
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      Map responseData = await makeDeleteRequest(
        uri: '/v1/users/account/recovery',
        body: requestBody,
        signer: primaryWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: primaryWallet!.publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        // sendFullDataToServer(responseData['data']);
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        hideLoader(context);
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
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

      Map responseData = await makeDeleteRequest(
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
        hideLoader(context);
        showSuccessAlert(context, onTap: () {
          appState.currentAction = PageAction(
              state: PageState.replaceAll, page: BottomHomePageConfig);
        });
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

  Future<Map> fetchQuestions(signer, secretKey, publicKey, username) async {
    Map responseData = await makeGetRequest(
      uri: '/v1/security-questions/$username',
      signer: signer,
      secretKey: secretKey, // the primary wallet secret key
      publicKey: publicKey!,
    );

    print('response: ${responseData}');

    return responseData['data'];
  }
}
