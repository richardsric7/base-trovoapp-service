import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SecurityQuestionsForInactiveAccounts extends StatefulWidget {
  const SecurityQuestionsForInactiveAccounts({Key? key}) : super(key: key);

  @override
  State<SecurityQuestionsForInactiveAccounts> createState() =>
      _SecurityQuestionsForInactiveAccounts();
}

class _SecurityQuestionsForInactiveAccounts
    extends State<SecurityQuestionsForInactiveAccounts> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  late Future<List<Map>> questions;
  String password = '';
  // late var primaryWallet;
  late String? username;
  late String? address;
  late String? signer;
  late String? secretKey;
  var questionsMap = {
    1: {"q": "", "a": "", "e": false},
    2: {"q": "", "a": "", "e": false},
    3: {"q": "", "a": "", "e": false},
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
    address =
        appState.viewData![SecurityQuestionsForInactiveAccountsViewPageConfig
            .key]['address'];
    signer =
        appState.viewData![SecurityQuestionsForInactiveAccountsViewPageConfig
            .key]['signer'];
    secretKey =
        appState.viewData![SecurityQuestionsForInactiveAccountsViewPageConfig
            .key]['secretKey'];
    username =
        appState.viewData![SecurityQuestionsForInactiveAccountsViewPageConfig
            .key]['username'];
    questions = fetchQuestions(signer, secretKey, address, username);
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
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 20,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 50),
                Text(
                  "setup".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30,
                    fontFamily: fontsemibold,
                  ),
                ),
                Text(
                  "securityquestions".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontSize: 30,
                    fontFamily: fontsemibold,
                  ),
                ),
                SizedBox(height: height / 20),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: const BorderRadius.all(
                        Radius.circular(15.0),
                      ),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 20.0,
                            vertical: 15.0,
                          ),
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
                                    fontFamily: fontbody,
                                  ),
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
                                  fontFamily: fontbody,
                                ),
                              ),
                              ElevatedButton(
                                onPressed: () {
                                  setState(() {
                                    questions = fetchQuestions(
                                      address,
                                      secretKey,
                                      address,
                                      username,
                                    );
                                  });
                                },
                                style: ButtonStyle(
                                  backgroundColor:
                                      WidgetStateProperty.all<Color>(
                                        notifier.getbluecolor!,
                                      ),
                                  foregroundColor:
                                      WidgetStateProperty.all<Color>(
                                        notifier.getwihitecolor,
                                      ),
                                ),
                                child: Text(
                                  "retry".tr(),
                                  style: TextStyle(fontFamily: fontsemibold),
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
                  },
                ),
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
                    bottom: MediaQuery.of(context).viewInsets.bottom,
                  ),
                ),
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
              contentPadding: EdgeInsets.symmetric(
                vertical: 20,
                horizontal: 20,
              ),
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
            70,
            300,
            // validator: validateEmail,
            onSaved: (value) {
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
    List<Map<dynamic, dynamic>> questionsList,
    int rel,
  ) {
    var filteredQuestions = questionsList.where((question) {
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
            child: Text(question['Question']),
            value: question['ID'].toString(),
          ),
        )
        .toList();
  }

  void trySubmit() async {
    final form = _formKey.currentState;
    if (!form!.validate()) {
      return;
    }

    form.save();
    // sendToServer();
    var map = {
      "q1": int.parse(questionsMap[1]!['q'].toString()),
      "a1": questionsMap[1]!['a'],
      "q2": int.parse(questionsMap[2]!['q'].toString()),
      "a2": questionsMap[2]!['a'],
      "q3": int.parse(questionsMap[3]!['q'].toString()),
      "a3": questionsMap[3]!['a'],
    };
    appState.setTempSecurityQuestionsAndAnswers = map;
    appState.viewData![EnsurePrivacyPageConfig.key] = {
      'rel': 'restoreUnactivatedAccount',
    };
    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: RequestBackupViewPageConfig,
    );
  }

  Future<List<Map>> fetchQuestions(signer, secretKey, address, username) async {
    Map responseData = await makeGetRequest(
      uri: '/v1/security-questions/$username',
      signer: signer,
      secretKey: secretKey, // the primary wallet secret key
      address: address!,
    );

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
