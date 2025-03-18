import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class BuyTokens extends StatefulWidget {
  const BuyTokens({Key? key}) : super(key: key);

  @override
  State<BuyTokens> createState() => _BuyTokens();
}

class _BuyTokens extends State<BuyTokens> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final _formKey = GlobalKey<FormState>();
  final amountController = TextEditingController();
  final quantityController = TextEditingController();

  double amount = 0;
  double quantity = 0;

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
    var tokenizedAsset = appState.tokenizedAsset!;
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                'Buy ${tokenizedAsset.assetName}',
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SizedBox(
                height: height / 30,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      'Currency',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: 300.sp,
                      height: 55.sp,
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 10),
                            child: Text(
                              tokenizedAsset.assetQuoteCurrency!,
                              style: TextStyle(fontSize: 15),
                            ),
                          ),
                          const SizedBox(height: 2),
                        ],
                      ),
                    ),
                  )
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      'Quantity of ${tokenizedAsset.assetCode}',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      'Quantity',
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      autoFormatNumber: true,
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onChanged: (value) {
                        setState(() {
                          var a = int.tryParse(value);
                          if (a != null) {
                            var pricePerToken = tokenizedAsset.pricePerToken!;
                            amountController.text =
                                formatNumberShort(a * pricePerToken).toString();
                          }
                        });
                      },
                      controller: quantityController,
                      onSaved: (value) {
                        quantity = double.parse(value);
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      'Amount',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      'Amount',
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      onChanged: (value) {
                        setState(() {
                          var a = int.tryParse(value);
                          if (a != null) {
                            var pricePerToken = tokenizedAsset.pricePerToken!;
                            quantityController.text =
                                formatNumberShort(a / pricePerToken).toString();
                          }
                        });
                      },
                      autoFormatNumber: true,
                      isFiat: true,
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      controller: amountController,
                      onSaved: (value) {
                        amount = double.parse(value);
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
              Button(
                'Pay',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  var form = _formKey.currentState;
                  if (form!.validate()) {
                    form.save();
                    appState.viewData = {
                      'amount': amount,
                      'quantity': quantity,
                    };
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: ConfirmBuyViewPageConfig,
                    );
                  }
                },
              ),
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget confirmLiensAndEncumbrance() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(5.sp),
              ),
            ),
            activeColor: notifier.isDark
                ? notifier.getbluecolor50
                : notifier.getbluecolor90,
            side: BorderSide(
              color: notifier.isDark
                  ? notifier.getbluecolor50
                  : notifier.getbluecolor90,
            ),
            value: true,
            onChanged: (bool? value) {
              setState(() {});
            },
          ),
        ),
        Container(
          width: width / 1.2,
          child: Text(
            'I confirm that this asset is completely free of all liens and encumbrance',
            overflow: TextOverflow.visible,
            style: TextStyle(
                fontSize: 15,
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody),
          ),
        ),
      ],
    );
  }
}

Widget CheckItem(
  String name,
  void Function()? onClick, {
  required Color backColor,
  required Color foreColor,
  required Color borderColor,
  double? fontSize = 15,
}) {
  return Padding(
    padding: const EdgeInsets.all(3.0),
    child: Container(
      decoration: BoxDecoration(
          border: Border.all(color: borderColor, width: 1),
          borderRadius: const BorderRadius.all(Radius.circular(10.0)),
          color: backColor),
      child: Padding(
        padding: const EdgeInsets.all(15.0),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                  color: foreColor, fontFamily: fontbody, fontSize: fontSize),
            ),
            SizedBox(
              width: width / 70,
            ),
          ],
        ),
      ),
    ),
  );
}
