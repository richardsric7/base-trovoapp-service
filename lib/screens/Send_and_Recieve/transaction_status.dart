import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
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

class TransactionStatus extends StatefulWidget {
  const TransactionStatus({Key? key}) : super(key: key);

  @override
  State<TransactionStatus> createState() => _TransactionStatus();
}

class _TransactionStatus extends State<TransactionStatus> {
  late ColorNotifier notifier;
  bool isChecked = false;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;

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
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Form(
            key: _formKey,
            child: Column(
              children: [
                SizedBox(height: height / 10),
                Text(
                  'Transaction Status',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 25.sp,
                      fontFamily: fontsemibold),
                ),
                SizedBox(
                  height: height / 20,
                ),
                Image.asset('assets/images/trovo.png', height: height / 7.5),
                SizedBox(
                  height: height / 30,
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 50),
                  child: Text(
                    'Your USDC Withdrawal is being processed',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        fontSize: 20,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ),
                SizedBox(
                  height: height / 30,
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 30),
                  child: Text(
                    'Your wallet will be debited once your Withdrawal transaction has been confirmed',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        fontSize: 16,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody),
                  ),
                ),
                SizedBox(height: height / 30),
                Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 30,
                  ),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Column(
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'Amount',
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Text(
                                '0.0000 USDC',
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                            ],
                          ),
                          SizedBox(
                            height: height / 90,
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'To',
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Text(
                                'GAVEB........KHR2Y',
                                style: TextStyle(
                                    fontSize: 13,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                            ],
                          ),
                          SizedBox(
                            height: height / 90,
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'From',
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Image.asset(
                                    'assets/images/trovo.png',
                                    height: 25,
                                    width: 25,
                                    errorBuilder: (context, error, stackTrace) {
                                      return Image.asset(
                                        'assets/images/trovo.png',
                                        height: 25,
                                        width: 25,
                                      );
                                    },
                                  ),
                                  SizedBox(
                                    width: width / 50.0,
                                  ),
                                  Text(
                                    'TROV',
                                    style: TextStyle(
                                        fontSize: 15,
                                        fontWeight: FontWeight.bold,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                SizedBox(height: height / 30),
                Stepper(
                  currentStep: 2,
                  controlsBuilder: (context, _) {
                    return Column(
                      children: [],
                    );
                  },
                  steps: <Step>[
                    Step(
                      state: StepState.complete,
                      isActive: true,
                      title: Text(
                        "Authorize withdrawal",
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 15.sp),
                      ),
                      content: Text(""),
                    ),
                    Step(
                      state: StepState.complete,
                      isActive: true,
                      title: Text(
                        "Processing withdrawal",
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 15.sp),
                      ),
                      content: Text(""),
                    ),
                    Step(
                      state: StepState.complete,
                      isActive: false,
                      title: Text(
                        "Transaction completed",
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 15.sp),
                      ),
                      content: Text(""),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 20,
                ),
                Button(
                  'View in History',
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: DepositWithdrawHistoryViewPageConfig,
                    );
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
}
