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
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class QuickBuyView extends StatefulWidget {
  const QuickBuyView({Key? key}) : super(key: key);

  @override
  State<QuickBuyView> createState() => _QuickBuyView();
}

class _QuickBuyView extends State<QuickBuyView> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  late Map viewData;
  String amount = '';
  double assetValue = 0.0;
  late Asset asset;
  final amountController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    viewData = appState.viewData!;
    asset = viewData['asset'] as Asset;
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
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(height: height / 40),
                  Container(
                    width: width,
                    child: Text(
                      'Buy ${getAssetCode(asset.assetCode)} With Fiat',
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
                      borderRadius: const BorderRadius.all(
                        Radius.circular(15.0),
                      ),
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
                            "An easy option to quickly buy virtual assets on Trovo App. You will get ${asset.assetCode} when you send fiat.",
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
                  Row(
                    children: [
                      Text(
                        "How much do you want to buy?",
                        textAlign: TextAlign.start,
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w400,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 50),
                  CustomTextFormField.textField(
                    "Amount (NGN)",
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    notifier.getprefixicon,
                    notifier.getblck,
                    notifier.getgrey,
                    75.sp,
                    double.infinity,
                    onChanged: (value) {
                      setState(() {
                        var splitText = value.split('.');
                        if (splitText.length > 2) {
                          // remove all dots except the first one.
                          value = '${splitText[0]}.${splitText[1]}';
                        }

                        amount = trim(value.toString(), '.');
                        assetValue =
                            double.tryParse(
                              calculateFiatValue(
                                amount,
                                asset.usdPrice.toString(),
                                'NGN',
                                appState,
                              ).replaceAll(',', ''),
                            ) ??
                            0.0;
                      });
                    },
                    controller: amountController,
                    autoFormatNumber: true,
                    keyboardtype: TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    validator: (String? value) {
                      if (value!.isEmpty) {
                        "Please enter how much you want to buy";
                      }

                      if (double.tryParse(value) == null) {
                        return "pleaseentervalidamount".tr();
                      }
                    },
                    onSaved: (value) =>
                        amount = value.trim().replaceAll(' ', ''),
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Flexible(
                        child: Text(
                          amount.isNotEmpty
                              ? "You will get ${formatNumber(assetValue)} ${getAssetCode(asset.assetCode)}"
                              : "≈ 0.0000 ${getAssetCode(asset.assetCode)}",

                          textScaler: TextScaler.linear(1.0),
                          style: TextStyle(
                            color: notifier.getdarkgrey,
                            fontWeight: FontWeight.w400,
                            fontSize: 12.0.sp,
                          ),
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 20),
                  Button(
                    "makepayment".tr(),
                    notifier.getbluecolor,
                    wihitecolor,
                    width: width - 40,
                    onTap: () {
                      final form = _formKey.currentState;
                      if (!form!.validate()) return;

                      appState.viewData!['amount'] = amount;
                      appState.viewData!['assetValue'] = assetValue;
                      appState.setPage(page: ConfirmQuickBuyViewPageConfig);
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
      ),
    );
  }
}
