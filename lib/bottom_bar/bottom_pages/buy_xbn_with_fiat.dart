import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:uuid/uuid.dart';

class BuyXBNWithFiat extends StatefulWidget {
  const BuyXBNWithFiat({Key? key}) : super(key: key);

  @override
  State<BuyXBNWithFiat> createState() => _BuyXBNWithFiat();
}

class _BuyXBNWithFiat extends State<BuyXBNWithFiat>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  late Map viewData;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    viewData = appState.viewData!;
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
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Column(
              children: [
                SizedBox(height: height / 40),
                Container(
                  width: width,
                  child: Text(
                    'Activate Account',
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontWeight: FontWeight.w500,
                      color: notifier.getbluewhitecolor,
                      fontSize: 20.sp,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ),
                SizedBox(height: height / 50),
                Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 20.0,
                      vertical: 15.0,
                    ),
                    child: Column(
                      children: [
                        SizedBox(height: height / 50),
                        Image.asset(
                          'assets/images/rafiki-buy-xbn.png',
                          // height: 50,
                          width: 180,
                        ),
                        SizedBox(height: height / 60),
                        Text(
                          "aneasyoptiontoactivateaccount".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        SizedBox(height: height / 50),
                      ],
                    ),
                  ),
                ),
                SizedBox(height: height / 50),
                Text(
                  "youpay".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 10, 0, 3),
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
                        Column(
                          children: [
                            SizedBox(height: height / 50),
                            Text(
                              'NGN ${getFiatValue(double.parse(viewData['activationAmount'].toString()))}',
                              style: TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w700,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                            SizedBox(height: height / 50.0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(height: height / 50),
                Text(
                  "youget".tr(),
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: height / 50),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 10, 0, 3),
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
                        Column(
                          children: [
                            SizedBox(height: height / 50),
                            SizedBox(
                              width: 300,
                              child: Text(
                                'NGN ${getFiatValue(getPercentageValue(double.parse(viewData['gasPercent'].toString()), double.parse(viewData['activationAmount'].toString())))} worth of Gas and NGN ${getFiatValue(getPercentageValue(double.parse(viewData['trovTokenPercent'].toString()), double.parse(viewData['activationAmount'].toString())))} worth of TROV',
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 15,
                                  fontWeight: FontWeight.w700,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                            SizedBox(height: height / 50.0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(height: height / 20),
                Button(
                  "makepayment".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  width: width - 40,
                  onTap: () {
                    savePaymentInvoiceAndContinue();
                  },
                ),
                SizedBox(height: height / 20),
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

  Future<void> savePaymentInvoiceAndContinue() async {
    try {
      var uuid = Uuid();
      String uniqueId = uuid.v4();
      showLoader(context);
      String requestBody = jsonEncode({
        'id': uniqueId,
        'amount': viewData['activationAmount'],
        'paymentType': 'ACTIVATION',
      });

      print(requestBody);

      var uri = '/v1/users/fiat/flutterwave';
      Map responseData = await makePostRequest(
        body: requestBody,
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      hideLoader(context);
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        appState.viewData!['id'] = uniqueId;
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: FlutterwaveWebViewPageConfig,
        );
      }
    } catch (e) {
      print('error');
      print(e);
      hideLoader(context);
    }
  }

  double getPercentageValue(double percentage, double amount) {
    return percentage * amount / 100;
  }
}
