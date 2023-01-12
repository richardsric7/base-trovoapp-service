import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  var viewData;

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
    viewData = appState.viewData![TransactionStatusViewPageConfig.key];
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
                    viewData['transactionDirection'] as TransactionDirection ==
                            TransactionDirection.Deposit
                        ? 'Your ${viewData['assetCode']} deposit is being processed'
                        : 'Your ${viewData['assetCode']} withdrawal is being processed',
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
                    viewData['transactionDirection'] as TransactionDirection ==
                            TransactionDirection.Deposit
                        ? 'You will receive funds in your wallet once your deposit transaction has been confirmed'
                        : 'Your wallet will be debited  once your Withdrawal transaction has been confirmed',
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
                          keyValuePair(
                              'Wallet', viewData['walletInfo']['alias']),
                          SizedBox(
                            height: height / 90,
                          ),
                          keyValuePair(
                              'Deposit Address',
                              truncate(viewData['depositAddress'], length: 5) +
                                  viewData['depositAddress'].substring(
                                      viewData['depositAddress'].length - 5)),
                          SizedBox(
                            height: height / 90,
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'Asset',
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Image.asset(
                                    viewData['assetImage'],
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
                                    viewData['assetCode'],
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
                        viewData['transactionDirection']
                                    as TransactionDirection ==
                                TransactionDirection.Deposit
                            ? "Make deposit"
                            : "Authorize withdrawal",
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
                        viewData['transactionDirection']
                                    as TransactionDirection ==
                                TransactionDirection.Deposit
                            ? "Processing deposit "
                            : "Processing withdrawal",
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

  Widget keyValuePair(String key, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          key,
          style: TextStyle(
              fontSize: 15,
              color: notifier.getbluewhitecolor,
              fontFamily: fontbody),
        ),
        Text(
          value,
          style: TextStyle(
              fontSize: 15,
              color: notifier.getbluewhitecolor,
              fontFamily: fontsemibold),
        ),
      ],
    );
  }
}
