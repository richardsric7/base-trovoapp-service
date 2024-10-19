import 'dart:convert';
import 'dart:developer';
import 'package:country_picker/country_picker.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:image_picker/image_picker.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetTokenInformation extends StatefulWidget {
  const AssetTokenInformation({Key? key}) : super(key: key);

  @override
  State<AssetTokenInformation> createState() => _AssetTokenInformation();
}

class _AssetTokenInformation extends State<AssetTokenInformation>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final _formKey = GlobalKey<FormState>();
  String? assetLogo;
  bool hasAdditionalKYCRequirements = false;
  late String proceedPayoutCurrency;
  late int numberOfTokenToBeSold;
  late int numberOfTokenToBeIssued;
  late int totalTokenHeldByManager;
  late double pricePerToken;
  late String assetCode;
  late String assetName;
  late DateTime? salesStart;
  late DateTime? salesEnd;
  late int capQuantity;
  late String assetQuoteCurrency;
  late int capDurationInDays;
  late String proceedCycle;
  late List<String> exemptedCountries;
  late String additionalKYCRequirements;
  late bool investorAccreditationRequired;
  late bool capOnPurchase;
  late String walletToHoldAssetsNotForSale;
  late int tokenizationFeeId;
  late dynamic data = {};
  final numberOfTokenToBeIssuedController = TextEditingController();
  final numberOfTokenToBeSoldController = TextEditingController();
  final capQuantityController = TextEditingController();
  final capDurationInDaysController = TextEditingController();
  bool formIsValid = true;

  List<String> assetQuoteCurrencies = [];
  List<String> payoutCycleOptions = [];

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getAssetQuoteCurrencies {
    List<DropdownMenuItem<String>> options = [];
    assetQuoteCurrencies.forEach((item) {
      options.add(DropdownMenuItem(
          child: Text(
            item,
            overflow: TextOverflow.ellipsis,
          ),
          value: item));
    });
    return options;
  }

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet.publicKey));
    });
    return wallets;
  }

  List<DropdownMenuItem<int>> getFeeItems(bool isSelected) {
    List<DropdownMenuItem<int>> items = [];
    for (var i = 0;
        i < appState.tokenizationData["tokenizationFees"].length;
        i++) {
      var item = appState.tokenizationData["tokenizationFees"][i];
      print(item['feeDescription']);
      var fiatPercentage = item['feeFiatPercentage'];
      var assetPercentage = item['feeAssetPercentage'];
      var fiatFeeCap = double.parse(item['feeFiatCap'].toString());
      var tokenFee = (numberOfTokenToBeIssued * assetPercentage) / 100;
      var fiatFee = (data['assetCurrentValue'] * fiatPercentage) / 100;

      if (i == 4) {
        items.add(DropdownMenuItem(
            child: Text(
              "Option ${i + 1} - ${assetQuoteCurrency}${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}",
              // "${item['feeDescription']} (\$${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}).",
              overflow:
                  isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
            ),
            value: item['id']));
        continue;
      }

      items.add(DropdownMenuItem(
          child: Text(
            "Option ${i + 1} - ${assetQuoteCurrency}${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}",
            // "${item['feeDescription']} (\$${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}).",
            overflow: isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
          ),
          value: item['id']));
    }
    return items;
  }

  List<DropdownMenuItem<String>> get getPayoutCycles {
    List<DropdownMenuItem<String>> cycles = [];
    payoutCycleOptions.forEach((item) {
      cycles.add(DropdownMenuItem(
          child: Text(
            item,
            overflow: TextOverflow.ellipsis,
          ),
          value: item));
    });
    return cycles;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    inspect(appState.viewData);
    inspect(appState.tokenizationData);
    data = appState.viewData;

    for (var i = 0;
        i < appState.tokenizationData["tokenizationCurrencies"].length;
        i++) {
      assetQuoteCurrencies
          .add(appState.tokenizationData["tokenizationCurrencies"][i]["label"]);
    }

    for (var i = 0;
        i < appState.tokenizationData["assetProceedCycle"].length;
        i++) {
      payoutCycleOptions
          .add(appState.tokenizationData["assetProceedCycle"][i]["id"]);
    }

    tokenizationFeeId =
        data["tokenizationFeeId"] == 0 ? 1 : data["tokenizationFeeId"];
    numberOfTokenToBeSold = data['numberOfTokenToBeSold'];
    numberOfTokenToBeIssued = data['numberOfTokenToBeIssued'];
    walletToHoldAssetsNotForSale =
        data['walletToHoldAssetsNotForSale'].toString().isEmpty
            ? ''
            : data['walletToHoldAssetsNotForSale'].toString();
    totalTokenHeldByManager = data['totalTokenHeldByManager'];
    pricePerToken = double.parse(data['pricePerToken'].toString());
    assetCode = data['assetCode'];
    assetName = data['assetName'];
    var parsedSalesStart = DateTime.parse(data['salesStart']);
    salesStart =
        parsedSalesStart.year == DateTime(0001).year ? null : parsedSalesStart;
    var parsedSalesEnd = DateTime.parse(data['salesEnd']);
    salesEnd =
        parsedSalesEnd.year == DateTime(0001).year ? null : parsedSalesEnd;
    capOnPurchase = data['capOnPurchase'] == 1;
    capQuantity = data['capQuantity'];
    capDurationInDays = data['capDurationInDays'];
    proceedCycle = data['proceedCycle'];
    assetLogo = data['assetLogo'];
    exemptedCountries = data['exemptedCountries'].toString().isEmpty
        ? []
        : data['exemptedCountries'].toString().split(',');
    hasAdditionalKYCRequirements = data['hasAdditionalKYCRequirements'] == 1;
    proceedPayoutCurrency = data['proceedPayoutCurrency'];
    assetQuoteCurrency = data['assetQuoteCurrency'];
    additionalKYCRequirements = data['additionalKYCRequirements'];
    investorAccreditationRequired = data['investorAccreditationRequired'] == 1;

    numberOfTokenToBeIssuedController.text = numberOfTokenToBeIssued == 0
        ? ''
        : formatNumberForInput(
            double.parse(numberOfTokenToBeIssued.toString()));
    numberOfTokenToBeSoldController.text = numberOfTokenToBeSold == 0
        ? ''
        : formatNumberForInput(double.parse(numberOfTokenToBeSold.toString()));
    capQuantityController.text = capQuantity == 0 ? '' : capQuantity.toString();
    capDurationInDaysController.text =
        capDurationInDays == 0 ? '' : capDurationInDays.toString();
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                "assettokeninfo".tr(),
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
                      "enterassetname".tr(),
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
                height: height / 70,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "enterassetname".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      initialValue: assetName,
                      onChanged: (value) {
                        setState(() {
                          assetName = value;
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          assetName = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "enterassetcode".tr(),
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
                height: height / 70,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "enterassetcode".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      initialValue: assetCode,
                      onChanged: (value) {
                        setState(() {
                          assetCode = value;
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          assetCode = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "uploadassetlogo".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  GestureDetector(
                    onTap: () {
                      getImage();
                    },
                    child: Column(
                      children: [
                        SizedBox(
                          height: height / 50,
                        ),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                          child: Container(
                            decoration: BoxDecoration(
                              border: Border.all(
                                  color: notifier.getbluewhitecolor, width: 1),
                              borderRadius:
                                  const BorderRadius.all(Radius.circular(15.0)),
                              color: notifier.isDark
                                  ? darktilewhitecolor
                                  : notifier.getaddsubwalletgrey,
                            ),
                            child: Column(
                              children: [
                                SizedBox(
                                    width: width / 1.2,
                                    height: height / 6,
                                    child: Center(
                                      child: Wrap(
                                        alignment: WrapAlignment.center,
                                        children: [
                                          Text(
                                            "browsefiles".tr(),
                                            textAlign: TextAlign.center,
                                            style: TextStyle(
                                              color: notifier.getbluewhitecolor,
                                              fontFamily: fontsemibold,
                                              fontSize: 12.sp,
                                            ),
                                          ),
                                        ],
                                      ),
                                    )),
                                const SizedBox(height: 2),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              if (!formIsValid &&
                  (assetLogo == null || assetLogo!.isEmpty)) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "pleaseuploadassetlogo".tr(),
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: Colors.red,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              if (assetLogo != null && assetLogo!.isNotEmpty) ...[
                GestureDetector(
                  onTap: () {
                    getImage();
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Image.memory(
                      base64Decode(assetLogo!),
                      width: width / 1.3,
                      height: height / 6,
                    ),
                  ),
                ),
              ],
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "nooftokentoissue".tr(),
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
                      "notobeissued".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      controller: numberOfTokenToBeIssuedController,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onChanged: (value) {
                        setState(() {
                          var val = value.toString().replaceAll('.', '');
                          numberOfTokenToBeIssued =
                              val.isNotEmpty ? int.parse(val) : 0;
                          totalTokenHeldByManager =
                              numberOfTokenToBeIssued - numberOfTokenToBeSold;

                          pricePerToken = (double.parse(
                                  data['assetCurrentValue'].toString()) /
                              numberOfTokenToBeIssued);
                          numberOfTokenToBeIssuedController.text =
                              val.isNotEmpty
                                  ? formatNumberForInput(double.parse(val))
                                  : val;
                        });
                      },
                      onSaved: (value) {
                        setState(() {
                          numberOfTokenToBeIssued = int.parse(value!);
                        });
                      },
                      autoFormatNumber: true,
                      inputFormatters: [
                        FilteringTextInputFormatter.allow(RegExp(r'[0-9 \.]'))
                      ],
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
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
                      "notobesold".tr(),
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
                      "notobesold".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      controller: numberOfTokenToBeSoldController,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onChanged: (value) {
                        print('this is value $value');
                        setState(() {
                          var val = value.toString().replaceAll('.', '');
                          numberOfTokenToBeSold =
                              val.isNotEmpty ? int.parse(val) : 0;
                          totalTokenHeldByManager =
                              numberOfTokenToBeIssued - numberOfTokenToBeSold;
                          numberOfTokenToBeSoldController.text = val.isNotEmpty
                              ? formatNumberForInput(double.parse(val))
                              : val;
                        });
                      },
                      onSaved: (value) {
                        setState(() {
                          numberOfTokenToBeSold = int.parse(value!);
                        });
                      },
                      autoFormatNumber: true,
                      inputFormatters: [
                        FilteringTextInputFormatter.allow(RegExp(r'[0-9 \.]'))
                      ],
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
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
                      "totaltobeheldbymanager".tr(),
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
                              formatNumberForInput(double.parse(
                                  totalTokenHeldByManager.toString())),
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
                      "selectwallettoholdassetnotforsale".tr(),
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
                height: height / 70,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      walletToHoldAssetsNotForSale = value.toString();
                    });
                  },
                  getStandardWallets,
                  getStandardWallets
                          .where((wallet) =>
                              wallet.value == walletToHoldAssetsNotForSale)
                          .isEmpty
                      ? null
                      : walletToHoldAssetsNotForSale,
                  'selectwallet'.tr(),
                  context,
                  null,
                  validator: (value) {
                    if (value == null || value.toString().isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: TextButton(
                  onPressed: () {
                    appState.returnView =
                        PageAction(state: PageState.addAll, pages: [
                      BottomHomePageConfig,
                      SetupAndComplianceViewPageConfig,
                      TokenizeAssetViewPageConfig,
                      AssetTokenInformationViewPageConfig
                    ]);
                    addSubWalletPopup(context);
                  },
                  style: TextButton.styleFrom(padding: EdgeInsets.zero),
                  child: Row(
                    children: [
                      Icon(
                        Icons.add_circle_outline_outlined,
                        color: notifier.getbluewhitecolor,
                        size: 18,
                      ),
                      SizedBox(width: 3),
                      Text(
                        "orcreateanewwallet".tr(),
                        style: TextStyle(
                          decoration: TextDecoration.underline,
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "selecttokenizationfee".tr(),
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
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      tokenizationFeeId = int.parse(value.toString());
                    });
                  },
                  getFeeItems(false),
                  getFeeItems(false)
                          .where((item) => item.value == tokenizationFeeId)
                          .isEmpty
                      ? null
                      : tokenizationFeeId,
                  'selectfee'.tr(),
                  context,
                  (context) {
                    return getFeeItems(true);
                  },
                  validator: (value) {
                    if (value == null || value.toString().isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                  // itemHeight: 70,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assetsalesandpricing".tr(),
                      style: TextStyle(
                        fontSize: 15,
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
                    child: Text(
                      "pricepertoken".tr(),
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
                              (pricePerToken == 0 || pricePerToken.isNaN)
                                  ? ''
                                  : formatNumberForInput(pricePerToken),
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "calculatedassetvalue".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
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
                    child: Text(
                      "assetquotecurrency".tr(),
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
                height: height / 70,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      assetQuoteCurrency = value.toString();
                    });
                  },
                  getAssetQuoteCurrencies,
                  assetQuoteCurrency.isEmpty ? null : assetQuoteCurrency,
                  'Select currency',
                  context,
                  null,
                  validator: (value) {
                    if (value == null || value == value.toString().isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            "salesstartsfrom".tr(),
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      ButtonOutlined(
                        salesStart != null && salesStart != DateTime(0)
                            ? DateFormat('MMMM dd, yyyy').format(salesStart!)
                            : "starts".tr(),
                        notifier.getwihitecolor,
                        notifier.getgrey,
                        borderColor: notifier.getgrey,
                        width: width / 2.7,
                        height: 50.sp,
                        onTap: () {
                          showDatePicker(
                            context: context,
                            initialDate: DateTime.now(),
                            firstDate:
                                DateTime.fromMicrosecondsSinceEpoch(1000),
                            lastDate: DateTime.now().add(Duration(days: 730)),
                          ).then((value) => setState(() {
                                salesStart = value;
                              }));
                        },
                      ),
                      if (!formIsValid &&
                          (salesStart == null ||
                              salesStart == DateTime(0))) ...[
                        Text(
                          "pleaseuploadassetlogo".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontbody,
                            color: Colors.red,
                          ),
                        ),
                      ],
                    ],
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            "salesendon".tr(),
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      ButtonOutlined(
                        salesEnd != null
                            ? DateFormat('MMMM dd, yyyy').format(salesEnd!)
                            : "ends".tr(),
                        notifier.getwihitecolor,
                        notifier.getgrey,
                        borderColor: notifier.getgrey,
                        width: width / 2.7,
                        height: 50.sp,
                        onTap: () {
                          showDatePicker(
                            context: context,
                            initialDate: DateTime.now(),
                            firstDate:
                                DateTime.fromMicrosecondsSinceEpoch(1000),
                            lastDate: DateTime.now().add(Duration(days: 730)),
                          ).then((value) => {
                                setState(() {
                                  salesEnd = value;
                                })
                              });
                        },
                      ),
                      if (!formIsValid &&
                          (salesEnd == null || salesEnd == DateTime(0))) ...[
                        Text(
                          "pleaseuploadassetlogo".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontbody,
                            color: Colors.red,
                          ),
                        ),
                      ],
                    ],
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
                      "caponpurchase".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  Spacer(),
                  Transform.scale(
                    scale: 0.7,
                    child: CupertinoSwitch(
                      activeColor: notifier.getgreencolor,
                      value: capOnPurchase,
                      onChanged: (val) async {
                        setState(() {
                          capOnPurchase = !capOnPurchase;
                        });
                      },
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Container(
                      width: width / 1.2,
                      child: Text(
                        "capexplained".tr(),
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              if (capOnPurchase) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "capquantity".tr(),
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
                        "quantity".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        70.sp,
                        300.sp,
                        controller: capQuantityController,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            capQuantity = int.parse(value!);
                          });
                        },
                        autoFormatNumber: true,
                        keyboardtype:
                            TextInputType.numberWithOptions(decimal: true),
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
                        "capduration".tr(),
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
                        "noofdays".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        70.sp,
                        300.sp,
                        controller: capDurationInDaysController,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            capDurationInDays = int.parse(value!);
                          });
                        },
                        autoFormatNumber: true,
                        keyboardtype:
                            TextInputType.numberWithOptions(decimal: true),
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 30,
                ),
              ],
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "proceedpayout".tr(),
                      style: TextStyle(
                        fontSize: 15,
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
                    child: Text(
                      "proceedpayoutcycle".tr(),
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
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      proceedCycle = value.toString();
                    });
                  },
                  getPayoutCycles,
                  proceedCycle.isEmpty ? null : proceedCycle,
                  'Select payout cycle',
                  context,
                  null,
                  validator: (value) {
                    if (value == null || value == value.toString().isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "proceedpayoutcurrency".tr(),
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
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      proceedPayoutCurrency = value.toString();
                    });
                  },
                  getAssetQuoteCurrencies,
                  proceedPayoutCurrency.isEmpty ? null : proceedPayoutCurrency,
                  'Select payout currency',
                  context,
                  null,
                  validator: (value) {
                    if (value == null || value == value.toString().isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "primarybuyerrequirement".tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Container(
                      width: width / 1.2,
                      child: Text(
                        "selectexemptedcountries".tr(),
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Container(
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10.0),
                    ),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: TextButton(
                      onPressed: showCountryListPopup,
                      style: ButtonStyle(
                          elevation: MaterialStateProperty.all<double>(0)),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            exemptedCountries.length > 0
                                ? exemptedCountries.last
                                : 'Select countries',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          Icon(Icons.keyboard_arrow_down_rounded),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 15.0),
                        child: Column(
                          children: [
                            Container(
                                width: width / 1.8,
                                child: Wrap(
                                  alignment: WrapAlignment.center,
                                  children: [
                                    if (exemptedCountries.isNotEmpty) ...[
                                      for (var country
                                          in exemptedCountries) ...[
                                        userItem(
                                          country,
                                          () {
                                            setState(() {
                                              exemptedCountries.remove(country);
                                            });
                                          },
                                          foreColor: wihitecolor,
                                          backColor: notifier.getbluebackcolor,
                                        ),
                                      ]
                                    ] else ...[
                                      Text(
                                        "nameofexemptedcountriesappearhere"
                                            .tr(),
                                        textAlign: TextAlign.center,
                                        style: TextStyle(
                                            color: notifier.getbluewhitecolor,
                                            fontFamily: fontbody,
                                            fontSize: 15.sp),
                                      ),
                                    ]
                                  ],
                                )),
                            SizedBox(height: 2),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "arethereadditionkycreq".tr(),
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
              Column(
                children: [
                  ListTile(
                    title: Text(
                      'yes'.tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    leading: Radio(
                      value: hasAdditionalKYCRequirements,
                      groupValue: true,
                      activeColor: notifier.getbluewhitecolor,
                      fillColor:
                          MaterialStateProperty.all(notifier.getbluewhitecolor),
                      onChanged: (value) {
                        setState(() {
                          hasAdditionalKYCRequirements = true;
                        });
                      },
                    ),
                  ),
                  ListTile(
                    title: Text(
                      'no'.tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    leading: Radio(
                      value: hasAdditionalKYCRequirements,
                      groupValue: false,
                      fillColor:
                          MaterialStateProperty.all(notifier.getbluewhitecolor),
                      activeColor: notifier.getbluewhitecolor,
                      onChanged: (value) {
                        setState(() {
                          hasAdditionalKYCRequirements = false;
                          print('addAdditionalKyc: $value');
                        });
                      },
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  if (hasAdditionalKYCRequirements) ...[
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Text(
                            "listrequirements".tr(),
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: multilineInput(
                            additionalKYCRequirements,
                            "Requirements",
                            notifier.getbluecolor,
                            notifier.getgrey,
                            notifier.getblck,
                            notifier.getgrey,
                            100.sp,
                            300.sp,
                            validator: (value) {
                              if (value.isEmpty) {
                                return "fieldcannotbeempty".tr();
                              }
                              return null;
                            },
                            onSaved: (value) {
                              setState(() {
                                additionalKYCRequirements = value!;
                              });
                            },
                            keyboardtype: TextInputType.multiline,
                            minLines: 3,
                            maxLines: null,
                          ),
                        ),
                      ],
                    ),
                  ],
                ],
              ),
              Row(
                children: [
                  Container(
                    width: width / 1.09,
                    child: checkBoxItem(
                      text: "investormustbeaccredited".tr(),
                      value: investorAccreditationRequired,
                      onChanged: (bool? value) {
                        setState(() {
                          investorAccreditationRequired =
                              !investorAccreditationRequired;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
              Button(
                "saveandcontinuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  setState(() {
                    formIsValid = true;
                    if (salesStart == null) {
                      formIsValid = false;
                    }

                    if (salesEnd == null) {
                      formIsValid = false;
                    }

                    if (assetLogo == null) {
                      formIsValid = false;
                    }
                    var form = _formKey.currentState;
                    if (form!.validate() && formIsValid) {
                      form.save();
                      submitForm();
                    }
                  });
                },
              ),
              SizedBox(
                height: height / 10,
              ),
              Padding(
                padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom),
              ),
            ],
          ),
        ),
      ),
    );
  }

  showCountryListPopup() {
    appState.dialogOpen = true;
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) {
      showCountryPicker(
        context: context,
        onSelect: (Country country) {
          setState(() {
            exemptedCountries.add(country.name);
          });
        },
      );
    });
    return SizedBox();
  }

  void submitForm() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credential
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;
      var newData = {...data as Map};

      newData['numberOfTokenToBeSold'] = numberOfTokenToBeSold;
      newData['numberOfTokenToBeIssued'] = numberOfTokenToBeIssued;
      newData['totalTokenHeldByManager'] =
          numberOfTokenToBeIssued - numberOfTokenToBeSold;
      newData['pricePerToken'] = pricePerToken;
      newData['assetCode'] = assetCode;
      newData['assetName'] = assetName;
      newData['salesStart'] = DateFormat("yyyy-MM-ddTHH:mm:ss.SSSSSS'Z'")
          .format(salesStart!.toUtc());
      newData['salesEnd'] =
          DateFormat("yyyy-MM-ddTHH:mm:ss.SSSSSS'Z'").format(salesEnd!);
      newData['capOnPurchase'] = capOnPurchase ? 1 : 0;
      newData['capQuantity'] = capQuantity;
      newData['capDurationInDays'] = capDurationInDays;
      newData['proceedCycle'] = proceedCycle;
      newData['walletToHoldAssetsNotForSale'] = walletToHoldAssetsNotForSale;
      newData['assetLogo'] = assetLogo;
      newData['exemptedCountries'] = exemptedCountries.join(',');
      newData['hasAdditionalKYCRequirements'] =
          hasAdditionalKYCRequirements ? 1 : 0;
      newData['assetQuoteCurrency'] = assetQuoteCurrency;
      newData['proceedPayoutCurrency'] = proceedPayoutCurrency;
      newData['additionalKYCRequirements'] = additionalKYCRequirements;
      newData['investorAccreditationRequired'] =
          investorAccreditationRequired ? 1 : 0;
      newData['tokenizationFeeId'] = tokenizationFeeId;

      String requestBody = jsonEncode(newData);
      print('requestBody  =======> $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: mintingWalletPublicKey,
      );

      hideLoader(context);

      print('responseData token information  ${responseData['data']}');
      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        Navigator.of(context).pop();
        await refreshCurrentTokenizationInfo();
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Future<void> refreshCurrentTokenizationInfo() async {
    try {
      var uri = '/v1/tokenization/detail/${appState.viewData!['id']}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      print('===============> token informationresponse ${responseData}');
      if (responseData['statusCode'] == 200) {
        print('success');
        appState.viewData = responseData['data'];
        await inspect(appState.viewData);
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Widget checkBoxItem(
      {required String text,
      required bool? value,
      required void Function(bool? value) onChanged}) {
    return Container(
      width: width / 1.08,
      child: Row(
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
              value: value,
              onChanged: onChanged,
            ),
          ),
          Container(
            width: width / 1.3,
            child: Text(
              text,
              overflow: TextOverflow.visible,
              style: TextStyle(
                  fontSize: 15,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody),
            ),
          ),
        ],
      ),
    );
  }

  Widget multilineInput(
    initialValue,
    labletext,
    focuscolor,
    lablecolor,
    textcolor,
    bordercolor,
    h,
    w, {
    onChanged,
    maxLength,
    minLines,
    maxLines,
    validator,
    onSaved,
    keyboardtype,
    focusNode,
  }) {
    return Container(
      height: h,
      width: w,
      child: TextFormField(
        focusNode: focusNode,
        maxLength: maxLength,
        initialValue: initialValue,
        minLines: minLines,
        maxLines: maxLines,
        style: TextStyle(color: textcolor, fontFamily: fontbody),
        cursorColor: lablecolor,
        onChanged: onChanged,
        decoration: InputDecoration(
          hintText: labletext,
          hintStyle: TextStyle(color: lablecolor),
          disabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          labelStyle: TextStyle(color: lablecolor),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: bordercolor, width: 1),
            borderRadius: BorderRadius.circular(15),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: focuscolor, width: 1),
            borderRadius: BorderRadius.circular(15),
          ),
        ),
        keyboardType: keyboardtype,
        validator: validator,
        onSaved: onSaved,
      ),
    );
  }

  String getFeeInfo(int index) {
    var fiatPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeFiatPercentage'];
    var assetPercentage = appState.tokenizationData["tokenizationFees"][index]
        ['feeAssetPercentage'];
    var fiatFeeCap = double.parse(appState.tokenizationData["tokenizationFees"]
            [index]['feeFiatCap']
        .toString());
    var tokenFee = (numberOfTokenToBeIssued * assetPercentage) / 100;
    var fiatFee = (data['assetCurrentValue'] * fiatPercentage) / 100;
    return "${appState.tokenizationData["tokenizationFees"][index]['feeDescription']} (\$${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}).";
  }

  Future<void> getImage() async {
    var image = await ImagePicker().pickImage(source: ImageSource.gallery);
    if (image != null) {
      var imageBase64Uncompressed = await getBase64Image(image);
      await StoreData().storeInsertData('image', imageBase64Uncompressed);
      setState(() {
        assetLogo = imageBase64Uncompressed;
      });
    }
  }

  Future<dynamic> getBase64Image(XFile image) async {
    //
    List<int> imageBytes = await image.readAsBytes();
    String imageB64 = base64Encode(imageBytes);
    return imageB64;
  }
}
