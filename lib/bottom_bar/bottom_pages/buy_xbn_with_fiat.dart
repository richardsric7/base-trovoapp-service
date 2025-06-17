import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';

import 'package:trovo_app/storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

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

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
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
                    'buyxbnwithfiat'.tr(),
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
                Container(
                  width: width,
                  child: Text(
                    'howmuchdoyouwanttobuy'.tr(),
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontWeight: FontWeight.w500,
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ),
                SizedBox(height: height / 90),
                CustomTextFormField.textField(
                  null,
                  notifier.getbluecolor,
                  null,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  notifier.getgrey,
                  70.sp,
                  350.sp,
                  hintText: 'N2000 - N5000',
                  onSaved: (value) {},
                ),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Flexible(
                      child: Text(
                        "250 XBN",
                        textScaler: TextScaler.linear(1.0),
                        style: TextStyle(
                          color: notifier.getdarkgrey,
                          fontWeight: FontWeight.w400,
                          fontSize: 12.0,
                        ),
                      ),
                    ),
                    Flexible(
                      child: Visibility(
                        visible: true,
                        replacement: Container(),
                        child: Text(
                          "0.5 TROV",
                          textScaler: TextScaler.linear(1.0),
                          textAlign: TextAlign.right,
                          style: TextStyle(
                            color: notifier.getdarkgrey,
                            fontSize: 12.0,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 20),
                Button(
                  "continuee".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  width: width - 40,
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: ConfirmBuyXBNWithFiatViewPageConfig,
                    );

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
}
