import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/utilities.dart';
import 'package:uuid/uuid.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class ConfirmQuickBuyView extends StatefulWidget {
  const ConfirmQuickBuyView({Key? key}) : super(key: key);

  @override
  State<ConfirmQuickBuyView> createState() => _ConfirmQuickBuyView();
}

class _ConfirmQuickBuyView extends State<ConfirmQuickBuyView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  var viewData;
  String amount = '';
  late Asset asset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    viewData = appState.viewData!;
    amount = viewData['amount'];
    asset = viewData['asset'];
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
                              'N${formatNumber(double.tryParse(amount) ?? 0.0)}',
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
                              '${formatNumber(appState.viewData!['assetValue'] ?? 0.0)} ${getAssetCode(asset.assetCode)}',
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
                    showPaymentMethodsPopup(context);
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

  showPaymentMethodsPopup(context) async {
    var notifier = Provider.of<ColorNotifier>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    int selectedGateway = 0;

    Widget buyOption({
      required String iconUrl,
      required String text,
      required bool isActive,
      required void Function() onTap,
    }) {
      return GestureDetector(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(0, 10, 0, 3),
          child: Container(
            width: width,
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    spacing: 10,
                    children: [
                      Image.asset(iconUrl, height: 40, width: 40),
                      Text(
                        text,
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w700,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: 20,
                    child: Transform.scale(
                      scale: 1,
                      child: Radio<bool>(
                        value: isActive,
                        activeColor: notifier.getbluewhitecolor,
                        fillColor: WidgetStateColor.resolveWith(
                          (states) => notifier.getbluewhitecolor,
                        ),
                        groupValue: true,
                        onChanged: (value) => onTap(),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      );
    }

    return showDialog(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (context, setStateForDialog) {
            return AlertDialog(
              backgroundColor: Colors.transparent,
              insetPadding: const EdgeInsets.all(1),
              content: Align(
                alignment: Alignment.bottomCenter,
                child: Container(
                  decoration: BoxDecoration(
                    color: notifier.getwihitecolor,
                    borderRadius: BorderRadius.all(Radius.circular(23)),
                  ),
                  width: double.infinity,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        SizedBox(height: 20),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              "Payment Methods",
                              textAlign: TextAlign.justify,
                              style: TextStyle(
                                fontSize: 16,
                                height: 1.4,
                                fontFamily: fontsemibold,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                            TextButton(
                              onPressed: () => Navigator.of(context).pop(),
                              child: Icon(Icons.cancel_outlined, size: 20),
                              style: TextButton.styleFrom(
                                padding: EdgeInsets.zero,
                                minimumSize: Size.zero,
                                tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                                visualDensity: VisualDensity.compact,
                                alignment: Alignment.centerLeft,
                                foregroundColor: notifier.getbluewhitecolor,
                                backgroundColor: Colors.transparent,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.zero,
                                ),
                              ),
                            ),
                          ],
                        ),
                        SizedBox(height: 30),
                        Text(
                          "Choose your preferred payment gateway to proceed to payment",
                          textAlign: TextAlign.start,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(height: 10),
                        buyOption(
                          iconUrl: 'assets/images/cngn-logo.png',
                          text: 'Paystack',
                          isActive: selectedGateway == 0,
                          onTap: () {
                            setStateForDialog(() {
                              selectedGateway = 0;
                            });
                          },
                        ),
                        SizedBox(height: 5),
                        buyOption(
                          iconUrl: 'assets/images/xbn-logo.png',
                          text: 'Flutterwave',
                          isActive: selectedGateway == 1,
                          onTap: () {
                            setStateForDialog(() {
                              selectedGateway = 1;
                            });
                          },
                        ),
                        SizedBox(height: 20),
                        ElevatedButton(
                          onPressed: () {
                            Navigator.of(context).pop();
                            var uuid = Uuid();
                            String uniqueId = uuid.v4();
                            appState.viewData!['id'] = uniqueId;
                            appState.viewData!['activationAmount'] = amount;
                            appState.setPage(
                              page: FlutterwaveWebViewPageConfig,
                            );
                          },
                          style: ButtonStyle(
                            padding: WidgetStateProperty.all(EdgeInsets.zero),
                            overlayColor: WidgetStateProperty.all<Color>(
                              notifier.getbluecolor90,
                            ),
                            backgroundColor: WidgetStateProperty.all<Color>(
                              notifier.getbluewhitecolor,
                            ),
                            foregroundColor: WidgetStateProperty.all<Color>(
                              notifier.getwihitecolor,
                            ),
                            side: WidgetStateProperty.all(
                              BorderSide(
                                color: notifier.getbluewhitecolor,
                                width: 1,
                                style: BorderStyle.solid,
                              ),
                            ),
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            shape:
                                WidgetStateProperty.all<RoundedRectangleBorder>(
                                  const RoundedRectangleBorder(
                                    borderRadius: BorderRadius.all(
                                      Radius.circular(10),
                                    ),
                                  ),
                                ),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Proceed',
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                  fontSize: 12,
                                  color: notifier.getwihitecolor,
                                ),
                              ),
                            ],
                          ),
                        ),
                        SizedBox(height: height / 50),
                      ],
                    ),
                  ),
                ),
              ),
            );
          },
        );
      },
    );
  }
}
