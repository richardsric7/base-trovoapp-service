import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/popups.dart';

import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AuthorizeSubscription extends StatefulWidget {
  const AuthorizeSubscription({Key? key}) : super(key: key);

  @override
  State<AuthorizeSubscription> createState() => _AuthorizeSubscriptionState();
}

class _AuthorizeSubscriptionState extends State<AuthorizeSubscription> {
  late ColorNotifier notifier;
  late DataProvider appState;
  String password = '';
  final formKey = GlobalKey<FormState>();

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
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
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
            title: Container(
              width: width / 1.7,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Padding(
                    padding: EdgeInsets.all(8),
                    child: Column(children: [
                      Text(
                        'Trovo Patron',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 22.sp,
                            fontFamily: fontsemibold),
                      ),
                    ]),
                  ),
                ],
              ),
            ),
          ),
        ),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                height: height / 15,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'You have chosen ',
                    style: TextStyle(
                        fontSize: 22,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody),
                  ),
                ],
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'Platinum',
                    style: TextStyle(
                        fontSize: 22,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: 20,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    border: Border.all(
                        color: notifier.getbluewhitecolor, width: 1.5),
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 10.0, vertical: 15.0),
                    child: Row(
                      children: [
                        Container(
                          width: width / 1.2,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              SizedBox(
                                width: width / 50,
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    'Subscription Plan',
                                    style: TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontsemibold,
                                    ),
                                  )
                                ],
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    'Monthly',
                                    style: TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  )
                                ],
                              ),
                              SizedBox(
                                height: 50,
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    'Expiry Date',
                                    style: TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontsemibold,
                                    ),
                                  )
                                ],
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    '14th Feb 2023',
                                    style: TextStyle(
                                      fontSize: 20,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  )
                                ],
                              ),
                              SizedBox(
                                height: 20,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.all(30.0),
                child: Text(
                  'Once your subscription expires, you will no longer have access to the Trovo Patron perks and will need to resubscribe to keep enjoying the benefits.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontSize: 12,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
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
                  // validator: validatePassword,
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
                  onTap: () {
                    showSuccessAlert(context, onTap: () {});
                  },
                ),
              ] else ...[
                Button(
                  LanguageEn.authorize,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    showSuccessAlert(context, onTap: () {});
                  },
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
}
