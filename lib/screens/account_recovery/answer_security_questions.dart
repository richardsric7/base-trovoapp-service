import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AnswerSecurityQuestions extends StatefulWidget {
  const AnswerSecurityQuestions({Key? key}) : super(key: key);

  @override
  State<AnswerSecurityQuestions> createState() => _AnswerSecurityQuestions();
}

class _AnswerSecurityQuestions extends State<AnswerSecurityQuestions> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  String password = '';
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
    securityQuestionsMap = fetchQuestions(appState.tempPublicKey,
        appState.tempSecretKey, appState.tempPublicKey, appState.tempUsername);
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
                height: height / 18)
            .getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 50),
                Text(
                  "answer".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 30.sp,
                      fontFamily: fontsemibold),
                ),
                Text(
                  "securityquestions".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
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
                                  "answersecurityquestionsdescription".tr(),
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
                                  "somethingwentwrong".tr(),
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
                                          appState.tempPublicKey,
                                          appState.tempSecretKey,
                                          appState.tempPublicKey,
                                          appState.tempUsername);
                                    });
                                  },
                                  style: ButtonStyle(
                                    backgroundColor:
                                        MaterialStateProperty.all<Color>(
                                            notifier.getbluecolor!),
                                  ),
                                  child: Text(
                                    "retry".tr(),
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

                          if (userSecurityAnswers['q1'] == 0 &&
                              userSecurityAnswers['q2'] == 0 &&
                              userSecurityAnswers['q1'] == 0) {
                            return Column(
                              children: [
                                Text(
                                    'You have not setup security questions yet.'),
                                SizedBox(height: height / 20),
                                Button(
                                  "back".tr(),
                                  notifier.getbluecolor,
                                  wihitecolor,
                                  onTap: () {
                                    Navigator.of(context).pop();
                                  },
                                ),
                              ],
                            );
                          }

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
                              ],
                              SizedBox(height: height / 20),
                              Button(
                                "continuee".tr(),
                                notifier.getbluecolor,
                                wihitecolor,
                                onTap: () {
                                  trySubmit();
                                },
                              ),
                              SizedBox(height: height / 10),
                              Padding(
                                  padding: EdgeInsets.only(
                                      bottom: MediaQuery.of(context)
                                          .viewInsets
                                          .bottom)),
                            ],
                          );
                        } else {
                          return const Text('Empty data');
                        }
                      } else {
                        return Text('State: ${snapshot.connectionState}');
                      }
                    }),
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
            "enteranswer".tr(),
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
                return "pleaseenteranswer".tr();
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

  sendToServer() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "q1": int.parse(questionsMap[1]!['q'].toString()),
        "a1": questionsMap[1]!['a'],
        "q2": int.parse(questionsMap[2]!['q'].toString()),
        "a2": questionsMap[2]!['a'],
        "q3": int.parse(questionsMap[3]!['q'].toString()),
        "a3": questionsMap[3]!['a']
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/verify-answers/${appState.tempUsername}',
        body: requestBody,
        signer: appState.tempPublicKey,
        secretKey: appState.tempSecretKey, // the primary wallet secret key
        publicKey: appState.tempPublicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        appState.setTempSecurityQuestionsAndAnswers = map;
        appState.viewData![EnsurePrivacyPageConfig.key] = {
          'rel': 'accountRecovery',
        };
        appState.currentAction = PageAction(
            state: PageState.addPage, page: RequestBackupViewPageConfig);
      } else {
        hideLoader(context);
        popup(context,
            title: "error".tr(), message: responseData['data']['error']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
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
