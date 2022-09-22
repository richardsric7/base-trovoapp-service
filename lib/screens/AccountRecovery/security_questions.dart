import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
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
  String password = '';
  String question1 = '';
  String question2 = '';
  String question3 = '';
  String question4 = '';
  String question5 = '';

  List<DropdownMenuItem<String>> get listOfQuestions {
    return [
      DropdownMenuItem<String>(
          child: Text(
            'What was the name of your street?',
            overflow: TextOverflow.ellipsis,
          ),
          value: 'Question 1'),
      DropdownMenuItem<String>(
          child: Text(
            'What is your favorite destination spot?',
            overflow: TextOverflow.ellipsis,
          ),
          value: 'Question 2'),
      DropdownMenuItem<String>(
          child: Text(
            'What was the name of your favorite teacher?',
            overflow: TextOverflow.ellipsis,
          ),
          value: 'Question 3'),
      DropdownMenuItem<String>(
          child: Text(
            'What is your favorite color?',
            overflow: TextOverflow.ellipsis,
          ),
          value: 'Question 4'),
      DropdownMenuItem<String>(
          child: Text(
            'What is your grandma\'s first name?',
            overflow: TextOverflow.ellipsis,
          ),
          value: 'Question 5'),
    ];
  }

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
                  LanguageEn.setup,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluecolor,
                      fontSize: 30.sp,
                      fontFamily: fontsemibold),
                ),
                Text(
                  LanguageEn.securityquestions,
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
                                  LanguageEn.setupsecurityquestionsdescription,
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
                          EdgeInsets.symmetric(vertical: 0, horizontal: 20),
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
                    ),
                    hint: Container(
                      child: Text(
                        '${LanguageEn.choosequestion} 1',
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500,
                    ),
                    onChanged: (newValue) async {
                      setState(() {
                        question1 = newValue!.toString();
                      });
                    },
                    items: listOfQuestions,
                  ),
                ),
                SizedBox(height: height / 50),
                if (question1.isNotEmpty) ...[
                  CustomTextFormField.textField(
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
                      print('email: $value');
                      // email = value.trim().replaceAll(' ', '');
                    },
                    keyboardtype: TextInputType.text,
                  ),
                ],
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
                          EdgeInsets.symmetric(vertical: 0, horizontal: 20),
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
                    ),
                    hint: Container(
                      child: Text(
                        '${LanguageEn.choosequestion} 2',
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500,
                    ),
                    onChanged: (newValue) async {
                      setState(() {
                        question2 = newValue!.toString();
                      });
                    },
                    items: listOfQuestions,
                  ),
                ),
                SizedBox(height: height / 50),
                if (question2.isNotEmpty) ...[
                  CustomTextFormField.textField(
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
                      print('email: $value');
                      // email = value.trim().replaceAll(' ', '');
                    },
                    keyboardtype: TextInputType.text,
                  ),
                ],
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
                          EdgeInsets.symmetric(vertical: 0, horizontal: 20),
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
                    ),
                    hint: Container(
                      child: Text(
                        '${LanguageEn.choosequestion} 3',
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500,
                    ),
                    onChanged: (newValue) async {
                      setState(() {
                        question3 = newValue!.toString();
                      });
                    },
                    items: listOfQuestions,
                  ),
                ),
                SizedBox(height: height / 50),
                if (question3.isNotEmpty) ...[
                  CustomTextFormField.textField(
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
                      print('email: $value');
                      // email = value.trim().replaceAll(' ', '');
                    },
                    keyboardtype: TextInputType.text,
                  ),
                ],
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
                          EdgeInsets.symmetric(vertical: 0, horizontal: 20),
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
                    ),
                    hint: Container(
                      child: Text(
                        '${LanguageEn.choosequestion} 4',
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500,
                    ),
                    onChanged: (newValue) async {
                      setState(() {
                        question4 = newValue!.toString();
                      });
                    },
                    items: listOfQuestions,
                  ),
                ),
                SizedBox(height: height / 50),
                if (question4.isNotEmpty) ...[
                  CustomTextFormField.textField(
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
                      print('email: $value');
                      // email = value.trim().replaceAll(' ', '');
                    },
                    keyboardtype: TextInputType.text,
                  ),
                ],
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
                          EdgeInsets.symmetric(vertical: 0, horizontal: 20),
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
                    ),
                    hint: Container(
                      child: Text(
                        '${LanguageEn.choosequestion} 5',
                        style: TextStyle(
                          color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500,
                    ),
                    onChanged: (newValue) async {
                      setState(() {
                        question5 = newValue!.toString();
                      });
                    },
                    items: listOfQuestions,
                  ),
                ),
                SizedBox(height: height / 50),
                if (question5.isNotEmpty) ...[
                  CustomTextFormField.textField(
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
                      print('email: $value');
                      // email = value.trim().replaceAll(' ', '');
                    },
                    keyboardtype: TextInputType.text,
                  ),
                ],
                SizedBox(height: height / 20),
                Button(
                  LanguageEn.continuee,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    print('object');
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: RequestOtpViewPageConfig);
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

  String? validatePassword(value) {
    print('password: $value');
    if (value.isEmpty) {
      //return "Enter a password";
      return LanguageEn.passwordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    return null;
  }

  String? validateConfirmPassword(value) {
    print('confirm password: ${value.trim().replaceAll(' ', '')} & $password');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.confirmpasswordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    if (password != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return LanguageEn.passwordmismatcherror;
    }

    return null;
  }

  bool validate() {
    final form = _formKey.currentState;
    if (form!.validate()) {
      form.save();
      return true;
    }
    return false;
  }

  void saveAndProceed() async {}
}
