import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SecurityQuestions extends StatefulWidget {
  const SecurityQuestions({Key? key}) : super(key: key);

  @override
  State<SecurityQuestions> createState() => _SecurityQuestions();
}

class _SecurityQuestions extends State<SecurityQuestions> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  late Future<List<Map>> questions;
  String password = '';
  // late var primaryWallet;
  late String? username;
  late String? publicKey;
  late String? signer;
  late String? secretKey;
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
    // primaryWallet = appState.userInfo!.wallets!
    //     .firstWhere((wallet) => wallet.primaryWallet == 1);
    publicKey =
        appState.viewData![SecurityQuestionsViewPageConfig.key]['publicKey'];
    signer = appState.viewData![SecurityQuestionsViewPageConfig.key]['signer'];
    secretKey =
        appState.viewData![SecurityQuestionsViewPageConfig.key]['secretKey'];
    username =
        appState.viewData![SecurityQuestionsViewPageConfig.key]['username'];
    questions = fetchQuestions(signer, secretKey, publicKey, username);
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
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(height / 15),
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                Text(
                  "setup".tr(),
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
                                  "setupsecurityquestionsdescription".tr(),
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
                FutureBuilder<List<Map>>(
                    future: questions,
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
                                      questions = fetchQuestions(publicKey,
                                          secretKey, publicKey, username);
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
                          return Column(
                            children: [
                              questionView(snapshot.data!, 1),
                              questionView(snapshot.data!, 2),
                              questionView(snapshot.data!, 3),
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
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget questionView(List<Map> questions, int rel) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: DropdownButtonFormField(
            isDense: true,
            isExpanded: true,
            dropdownColor: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
            decoration: InputDecoration(
              contentPadding:
                  EdgeInsets.symmetric(vertical: 20, horizontal: 20),
              enabledBorder: OutlineInputBorder(
                borderSide: BorderSide.none,
                borderRadius: BorderRadius.circular(10),
              ),
              border: OutlineInputBorder(
                borderSide: BorderSide.none,
                borderRadius: BorderRadius.circular(10),
              ),
              filled: true,
              fillColor: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
              errorStyle: TextStyle(
                fontFamily: fontbody,
                fontSize: 12,
                overflow: TextOverflow.visible,
              ),
            ),
            hint: Container(
              child: Text(
                // Choose question n
                '${"choosequestion".tr()} $rel',
                style: TextStyle(
                  color: questionsMap[rel]!['e'] == true
                      ? Colors.red
                      : notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
                textAlign: TextAlign.end,
              ),
            ),
            icon: Icon(
              Icons.keyboard_arrow_down_rounded,
              color: notifier.getbluewhitecolor,
            ),
            elevation: 0,
            style: TextStyle(
              overflow: TextOverflow.visible,
              color: notifier.getbluewhitecolor,
              fontSize: 15,
              fontFamily: fontsemibold,
              fontWeight: FontWeight.w500,
            ),
            onChanged: (newValue) async {
              setState(() {
                questionsMap[rel]!['q'] = newValue!.toString();
              });
            },
            validator: (value) {
              if (questionsMap[rel]!['q'].toString().isEmpty) {
                setState(() {
                  questionsMap[rel]!['e'] = true;
                });
                return '';
              }

              setState(() {
                questionsMap[rel]!['e'] = false;
              });

              return null;
            },
            items: getQuestions(questions, rel),
          ),
        ),
        SizedBox(height: height / 50),
        if (questionsMap[rel]!['q'].toString().isNotEmpty) ...[
          CustomTextFormField.textField(
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
              print('email: $value');
              questionsMap[rel]!['a'] = value.toString().trim();
            },
            validator: (value) {
              if (value.toString().isEmpty) {
                return 'Please enter anwser to the question';
              }
              return null;
            },
            keyboardtype: TextInputType.text,
          ),
        ],
      ],
    );
  }

  List<DropdownMenuItem<String>> getQuestions(
      List<Map<dynamic, dynamic>> questionsList, int rel) {
    var filteredQuestions = questionsList.where((question) {
      print('Q: ${questionsMap[1]!['q']} A: ${question['ID']}');
      if (rel != 1 && questionsMap[1]!['q'] == question['ID'].toString())
        return false;
      if (rel != 2 && questionsMap[2]!['q'] == question['ID'].toString())
        return false;
      if (rel != 3 && questionsMap[3]!['q'] == question['ID'].toString())
        return false;
      return true;
    }).toList();
    return filteredQuestions
        .map(
          (question) => DropdownMenuItem<String>(
              child: Text(
                question['Question'],
              ),
              value: question['ID'].toString()),
        )
        .toList();
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
        uri: '/v1/security-questions',
        body: requestBody,
        signer: signer!,
        secretKey: secretKey!, // the primary wallet secret key
        publicKey: publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        await updateUserInfo(signer, secretKey, publicKey, username, appState);
        appState.viewData = {
          SuccessViewPageConfig.key: {
            'title': "success".tr(),
            'message': "securityquestionssuccessmessage".tr(),
            'useOnDone': true,
            'onDone': () {
              appState.currentAction = appState.returnView ??
                  PageAction(
                    state: PageState.addAll,
                    pages: [BottomHomePageConfig],
                  );
            },
          }
        };
        appState.currentAction = PageAction(
            state: PageState.replaceAll, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Future<List<Map>> fetchQuestions(
      signer, secretKey, publicKey, username) async {
    Map responseData = await makeGetRequest(
      uri: '/v1/security-questions/$username',
      signer: signer,
      secretKey: secretKey, // the primary wallet secret key
      publicKey: publicKey!,
    );

    print('response: ${responseData}');
    var questionsList = <Map>[];

    if (responseData['statusCode'] == 200) {
      var questions = responseData['data']['securityQuestions'];
      for (var i = 0; i < questions.length; i++) {
        questionsList.add(responseData['data']['securityQuestions'][i]);
      }
    }
    return questionsList;
  }
}
