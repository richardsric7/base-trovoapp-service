import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TransactionStatus extends StatefulWidget {
  const TransactionStatus({Key? key}) : super(key: key);

  @override
  State<TransactionStatus> createState() => _TransactionStatus();
}

class _TransactionStatus extends State<TransactionStatus> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late Asset? asset;
  var transactionInfo = {};

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

    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    transactionInfo = appState.viewData!['transactionData'];
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
        body: SingleChildScrollView(
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
              Image.asset('assets/images/trovo.png', height: height / 8.5),
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 50),
                child: Text(
                  'Your [${transactionInfo['currency']}] withdrawal request is being processed',
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
                  'Your wallet will be debited  once your Withdrawal transaction has been confirmed',
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
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Column(
                      children: [
                        keyValuePair('Wallet', wallet.alias!),
                        SizedBox(
                          height: height / 90,
                        ),
                        keyValuePair(
                            'Withdraw Address',
                            truncate(transactionInfo['withdrawalAddress'],
                                    length: 5) +
                                transactionInfo['withdrawalAddress'].substring(
                                    transactionInfo['withdrawalAddress']
                                            .length -
                                        5)),
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
                                Text(
                                  transactionInfo['currency'],
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
              SizedBox(
                height: height / 50,
              ),
              ButtonOutlined(
                LanguageEn.dashboard,
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.replaceAll,
                    page: BottomHomePageConfig,
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
