import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
              SizedBox(height: height / 30),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Column(
                  children: [
                    Row(
                      children: [
                        Text(
                          'Amount',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 50),
                    Row(
                      children: [
                        CustomTextFormField.textField(
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
                              var a = double.tryParse(value);
                              if (a != null) {
                                amount = a;
                                var pricePerToken =
                                    tokenizedAsset.pricePerToken!;
                                quantity = a / pricePerToken;
                              } else {
                                quantity = 0;
                              }
                            });
                          },
                          keyboardtype: TextInputType.numberWithOptions(
                            decimal: true,
                          ),
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
                      ],
                    ),
                    if (!appState.hideBalances) ...[availableBalance()],
                    SizedBox(height: height / 70),
                    Row(
                      children: [
                        Text(
                          'Currency',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 50),
                    Row(
                      children: [
                        Container(
                          width: 300.sp,
                          height: 55.sp,
                          decoration: BoxDecoration(
                            borderRadius: const BorderRadius.all(
                              Radius.circular(15.0),
                            ),
                            color: notifier.getaddsubwalletgrey,
                          ),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Padding(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 10,
                                ),
                                child: Text(
                                  tokenizedAsset.assetQuoteCurrency!,
                                  style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getblck,
                                  ),
                                ),
                              ),
                              const SizedBox(height: 2),
                            ],
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 50),
                    Row(
                      children: [
                        Text(
                          'Quantity of ${tokenizedAsset.assetCode}',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 50),
                    Row(
                      children: [
                        Container(
                          width: 300.sp,
                          height: 55.sp,
                          decoration: BoxDecoration(
                            borderRadius: const BorderRadius.all(
                              Radius.circular(15.0),
                            ),
                            color: notifier.getaddsubwalletgrey,
                          ),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Padding(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 10,
                                ),
                                child: Text(
                                  formatNumber(quantity),
                                  style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getblck,
                                  ),
                                ),
                              ),
                              const SizedBox(height: 2),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 30),
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
              SizedBox(height: height / 10),
            ],
          ),
        ),
      ),
    );
  }

  Widget availableBalance() {
    var asset = appState.activeWallet!.claimedAssets!.firstWhere(
      (asset) => asset.assetCode!.toLowerCase() == 'cngn',
    );
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Flexible(
          child: Text(
            amount.toString().isNotEmpty
                ? "≈ ${formatNumber(amount)} ${getAssetCode(asset.assetCode)}"
                : "≈ 0.0000 ${getAssetCode(asset.assetCode)}",
            textScaleFactor: 1.0,
            style: TextStyle(
              color: notifier.getdarkgrey,
              fontWeight: FontWeight.w400,
              fontSize: 12.0.sp,
            ),
          ),
        ),
        Flexible(
          child: Visibility(
            visible: true,
            replacement: Container(),
            child: Text(
              "${formatNumber(asset.amount!)} ${getAssetCode(asset.assetCode)}",
              textScaleFactor: 1.0,
              textAlign: TextAlign.right,
              style: TextStyle(color: notifier.getdarkgrey, fontSize: 12.0.sp),
            ),
          ),
        ),
      ],
    );
  }
}
