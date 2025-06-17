import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutterwave_standard/core/flutterwave.dart';
import 'package:flutterwave_standard/models/requests/customer.dart';
import 'package:flutterwave_standard/models/requests/customizations.dart';
import 'package:flutterwave_standard/models/responses/charge_response.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ConfirmBuyXBNWithFiat extends StatefulWidget {
  const ConfirmBuyXBNWithFiat({Key? key}) : super(key: key);

  @override
  State<ConfirmBuyXBNWithFiat> createState() => _ConfirmBuyXBNWithFiat();
}

class _ConfirmBuyXBNWithFiat extends State<ConfirmBuyXBNWithFiat>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  var viewData;

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
    appState = Provider.of<DataProvider>(context, listen: true);
    inspect(appState.viewData);

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
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Text(
                      "confirmyourtransaction".tr(),
                      style: TextStyle(
                        fontSize: 20.sp,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 20),
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
                              'N2,500',
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
                            Text(
                              '250 XBN & 0.5 TROV',
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
                SizedBox(height: height / 20),
                Button(
                  "makepayment".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  width: width - 40,
                  onTap: () {
                    handlePaymentInitialization();
                    // appState.currentAction = PageAction(
                    //   state: PageState.addPage,
                    //   page: ConfirmBuyXBNWithFiatViewPageConfig,
                    // );

                    // appState.viewData![ShareReceiptViewPageConfig.key] = viewData;
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

  handlePaymentInitialization() async {
    final Customer customer = Customer(
      name: "Flutterwave Developer",
      phoneNumber: "1234566677777",
      email: "customer@customer.com",
    );
    final Flutterwave flutterwave = Flutterwave(
      publicKey: "FLWPUBK_TEST-45bd332ee4bdefdcacd6d2513944cd16-X",
      currency: "ngn",
      redirectUrl: "trovo.app.link",
      txRef: "xdvdsw3422d",
      amount: '2500',
      customer: customer,
      paymentOptions: "ussd, card, bank transfer",
      customization: Customization(title: "Buy XBN and TROV"),
      isTestMode: true,
    );

    final ChargeResponse response = await flutterwave.charge(context);
    print(response);
    // Handle the response
    if (response.success == true) {
      // Payment was successful
      print('charge successful');
    } else {
      print('charge unsuccessful');
      // Payment failed or was cancelled
    }
  }
}
