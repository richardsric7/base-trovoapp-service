import 'dart:convert';
import 'dart:developer';
import 'package:country_picker/country_picker.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:image_picker/image_picker.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  late int proceedPayoutType;
  late int numberOfTokenToBeIssued;
  late String assetCode;
  late String assetName;
  late DateTime? salesStart;
  late DateTime? salesEnd;
  late double? capQuantity;
  late double? capAmountInFiat;
  late String assetQuoteCurrency;
  late int capDurationInDays;
  late String proceedCycle;
  late List<String> exemptedCountries;
  late String additionalKYCRequirements;
  late bool investorAccreditationRequired;
  late bool capOnPurchase;
  late bool attestInformationAccurateAndVerifiable;
  late bool acknowledgedSuitabilityCriteria;
  late bool withholdingTaxDisclosure;
  late String walletToHoldAssetsNotForSale;
  late String authorizedRepresentativeName;
  late String authorizedRepresentativeTitleOrPosition;
  late String authorizedRepresentativeEmail;
  late String minimumKycTier;
  late String investorCategory;
  late int bankId;
  late String accountNumber;
  late String beneficiaryName;
  late int tokenizationFeeId;
  late dynamic data = {};
  final numberOfTokenToBeIssuedController = TextEditingController();
  final capQuantityController = TextEditingController();
  final capAmountInFiatController = TextEditingController();
  final capDurationInDaysController = TextEditingController();
  bool formIsValid = true;

  List<Map> assetQuoteCurrencies = [];
  List<Map> assetPayoutTypes = [
    {'text': 'Crypto', 'value': 0},
    {'text': 'Fiat', 'value': 1},
  ];
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

  List<DropdownMenuItem<int>> get getBanksList {
    List<DropdownMenuItem<int>> options = [];
    appState.tokenizationData['banks']['bankList'].forEach((item) {
      options.add(
        DropdownMenuItem(
          child: Text(item["bankName"], overflow: TextOverflow.ellipsis),
          value: item["id"],
        ),
      );
    });
    return options;
  }

  List<DropdownMenuItem<String>> get getAssetQuoteCurrencies {
    List<DropdownMenuItem<String>> options = [];
    assetQuoteCurrencies.forEach((item) {
      options.add(
        DropdownMenuItem(
          child: Text(
            item["assetCode"].toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: item["assetCode"].toString(),
        ),
      );
    });
    return options;
  }

  List<DropdownMenuItem<String>> get getPayoutCurrencies {
    List<DropdownMenuItem<String>> options = [];
    assetQuoteCurrencies.forEach((item) {
      options.add(
        DropdownMenuItem(
          child: Text(
            item["assetCode"].toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: item["assetCode"].toString(),
        ),
      );
    });
    return options;
  }

  List<DropdownMenuItem<int>> get getAssetPayoutType {
    List<DropdownMenuItem<int>> options = [];
    assetPayoutTypes.forEach((item) {
      options.add(
        DropdownMenuItem(
          child: Text(item['text'], overflow: TextOverflow.ellipsis),
          value: item['value'],
        ),
      );
    });
    return options;
  }

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(
        DropdownMenuItem(
          child: Text(wallet.alias!, overflow: TextOverflow.ellipsis),
          value: wallet.publicKey,
        ),
      );
    });
    return wallets;
  }

  List<DropdownMenuItem<int>> getFeeItems(bool isSelected) {
    List<DropdownMenuItem<int>> items = [];
    for (
      var i = 0;
      i < appState.tokenizationData["tokenizationFees"].length;
      i++
    ) {
      var item = appState.tokenizationData["tokenizationFees"][i];
      var fiatPercentage = item['feeFiatPercentage'];
      var assetPercentage = item['feeAssetPercentage'];
      var fiatFeeCap = double.parse(item['feeFiatCap'].toString());
      var tokenFee = (numberOfTokenToBeIssued * assetPercentage) / 100;
      var fiatFee = (data['assetCurrentValue'] * fiatPercentage) / 100;

      if (i == 0) {
        items.add(
          DropdownMenuItem(
            child: Text(
              "Select fee",
              overflow: isSelected
                  ? TextOverflow.ellipsis
                  : TextOverflow.visible,
            ),
            value: item['id'],
          ),
        );
        continue;
      }

      items.add(
        DropdownMenuItem(
          child: Text(
            "Option ${i} - ${assetQuoteCurrency}${formatNumber(fiatFeeCap > fiatFee ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}",
            overflow: isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
          ),
          value: item['id'],
        ),
      );
    }
    return items;
  }

  List<DropdownMenuItem<String>> get getPayoutCycles {
    List<DropdownMenuItem<String>> cycles = [];
    payoutCycleOptions.forEach((item) {
      cycles.add(
        DropdownMenuItem(
          child: Text(item, overflow: TextOverflow.ellipsis),
          value: item,
        ),
      );
    });
    return cycles;
  }

  List<DropdownMenuItem<String>> get getInvestorCategory {
    List<DropdownMenuItem<String>> categories = [];
    ["Retail", "Qualified", "Institutional"].forEach((item) {
      categories.add(
        DropdownMenuItem(
          child: Text(item, overflow: TextOverflow.ellipsis),
          value: item,
        ),
      );
    });
    return categories;
  }

  List<DropdownMenuItem<String>> get getKYCTiers {
    List<DropdownMenuItem<String>> tiers = [];
    ["Tier 1", "Tier 2", "Tier 3"].forEach((item) {
      tiers.add(
        DropdownMenuItem(
          child: Text(item, overflow: TextOverflow.ellipsis),
          value: item,
        ),
      );
    });
    return tiers;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;

    for (
      var i = 0;
      i < appState.tokenizationData["tokenizationCurrencies"].length;
      i++
    ) {
      assetQuoteCurrencies.add(
        appState.tokenizationData["tokenizationCurrencies"][i],
      );
    }

    for (
      var i = 0;
      i < appState.tokenizationData["assetProceedCycle"].length;
      i++
    ) {
      payoutCycleOptions.add(
        appState.tokenizationData["assetProceedCycle"][i]["id"],
      );
    }

    tokenizationFeeId = data["tokenizationFeeId"] == 0
        ? 1
        : data["tokenizationFeeId"];
    numberOfTokenToBeIssued = data['numberOfTokenToBeIssued'];
    walletToHoldAssetsNotForSale =
        data['walletToHoldAssetsNotForSale'].toString().isEmpty
        ? ''
        : data['walletToHoldAssetsNotForSale'].toString();

    authorizedRepresentativeName = data['authorizedRepresentativeName'] ?? '';

    authorizedRepresentativeTitleOrPosition =
        data['authorizedRepresentativeTitleOrPosition'] ?? '';
    authorizedRepresentativeEmail = data['authorizedRepresentativeEmail'] ?? '';

    minimumKycTier = data['minimumKycTier'] ?? '';
    investorCategory = data['investorCategory'] ?? '';

    assetCode = data['assetCode'];
    assetName = data['assetName'];
    accountNumber = data['accountNumber'];
    beneficiaryName = data['beneficiaryName'];
    bankId = data['bankId'];
    var parsedSalesStart = DateTime.parse(data['salesStart']);
    salesStart = parsedSalesStart.year == DateTime(0001).year
        ? null
        : parsedSalesStart;
    var parsedSalesEnd = DateTime.parse(data['salesEnd']);
    salesEnd = parsedSalesEnd.year == DateTime(0001).year
        ? null
        : parsedSalesEnd;
    capOnPurchase = data['capOnPurchase'] == 1;
    attestInformationAccurateAndVerifiable =
        data['acceptTokenizationTermsAndAgreement'] == 1;
    acknowledgedSuitabilityCriteria =
        data['acceptTokenizationTermsAndAgreement'] == 1;
    capQuantity = double.parse(data['capQuantity'].toString());
    capAmountInFiat = double.parse(data['capAmountInFiat'].toString());
    capDurationInDays = data['capDurationInDays'];
    proceedCycle = data['proceedCycle'];
    assetLogo = data['assetLogo'];
    exemptedCountries = data['exemptedCountries'].toString().isEmpty
        ? []
        : data['exemptedCountries'].toString().split(',');
    hasAdditionalKYCRequirements = data['hasAdditionalKYCRequirements'] == 1;
    withholdingTaxDisclosure = data['withholdingTaxDisclosure'] == 1;
    proceedPayoutCurrency = data['proceedPayoutCurrency'];
    proceedPayoutType = data['proceedPayoutType'];
    assetQuoteCurrency = data['assetQuoteCurrency'];
    additionalKYCRequirements = data['additionalKYCRequirements'];
    investorAccreditationRequired = data['investorAccreditationRequired'] == 1;

    numberOfTokenToBeIssuedController.text = numberOfTokenToBeIssued == 0
        ? ''
        : formatNumberForInput(
            double.parse(numberOfTokenToBeIssued.toString()),
          );
    // numberOfTokenToBeSoldController.text = numberOfTokenToBeSold == 0
    //     ? ''
    //     : formatNumberForInput(double.parse(numberOfTokenToBeSold.toString()));
    capQuantityController.text = capQuantity == 0
        ? ''
        : formatNumberForInput(double.parse(capQuantity.toString()));
    capAmountInFiatController.text = capAmountInFiat == 0
        ? ''
        : formatNumberForInput(double.parse(capAmountInFiat.toString()));
    capDurationInDaysController.text = capDurationInDays == 0
        ? ''
        : capDurationInDays.toString();
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
              SizedBox(height: height / 30),
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
              SizedBox(height: height / 70),
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
                      85,
                      300.sp,
                      initialValue: assetName,
                      readOnly: true,
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
              SizedBox(height: height / 70),
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
                      85,
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
                      getFile();
                    },
                    child: Column(
                      children: [
                        SizedBox(height: height / 50),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                          child: Container(
                            decoration: BoxDecoration(
                              border: Border.all(
                                color: notifier.getbluewhitecolor,
                                width: 1,
                              ),
                              borderRadius: const BorderRadius.all(
                                Radius.circular(15.0),
                              ),
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
                                  ),
                                ),
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
                    getFile();
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Image.network(
                      assetLogo!,
                      width: width / 1.3,
                      height: height / 6,
                    ),
                  ),
                ),
              ],
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 50),
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
                      85,
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
                          numberOfTokenToBeIssued = value.isNotEmpty
                              ? int.parse(value)
                              : 0;
                        });
                      },
                      onSaved: (value) {
                        setState(() {
                          numberOfTokenToBeIssued = int.parse(value!);
                        });
                      },
                      autoFormatNumber: true,
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 70),
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
                          .where(
                            (wallet) =>
                                wallet.value == walletToHoldAssetsNotForSale,
                          )
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
                    appState.returnView = PageAction(
                      state: PageState.addAll,
                      pages: [
                        BottomHomePageConfig,
                        SetupAndComplianceViewPageConfig,
                        TokenizeAssetViewPageConfig,
                        AssetTokenInformationViewPageConfig,
                      ],
                    );
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
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      tokenizationFeeId = int.parse(value.toString());
                    });
                  },
                  getFeeItems(false),
                  getFeeItems(
                        false,
                      ).where((item) => item.value == tokenizationFeeId).isEmpty
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
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "assettokensale".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 70),
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
              SizedBox(height: height / 50),
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
                      SizedBox(height: height / 50),
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
                            firstDate: DateTime.fromMicrosecondsSinceEpoch(
                              1000,
                            ),
                            lastDate: DateTime.now().add(Duration(days: 730)),
                          ).then(
                            (value) => setState(() {
                              salesStart = value;
                            }),
                          );
                        },
                      ),
                      if (!formIsValid &&
                          (salesStart == null ||
                              salesStart == DateTime(0))) ...[
                        Text(
                          "pleaseselectdate".tr(),
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
                      SizedBox(height: height / 50),
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
                            firstDate: DateTime.fromMicrosecondsSinceEpoch(
                              1000,
                            ),
                            lastDate: DateTime.now().add(Duration(days: 730)),
                          ).then(
                            (value) => {
                              setState(() {
                                salesEnd = value;
                              }),
                            },
                          );
                        },
                      ),
                      if (!formIsValid &&
                          (salesEnd == null || salesEnd == DateTime(0))) ...[
                        Text(
                          "pleaseselectdate".tr(),
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
              SizedBox(height: height / 80),
              if (!formIsValid &&
                  salesEnd != null &&
                  salesEnd!.isBefore(salesStart!)) ...[
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        "salesendmustbeaftersalesstart".tr(),
                        textAlign: TextAlign.start,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: Colors.red,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "caponpurchase".tr(),
                      style: TextStyle(
                        fontSize: 12,
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
              SizedBox(height: height / 50),
              if (capOnPurchase) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "capamount".tr(),
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: CustomTextFormField.textField(
                        "amount".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        85,
                        300.sp,
                        controller: capAmountInFiatController,
                        onChanged: (val) {
                          setState(() {
                            capAmountInFiat = val.isNotEmpty
                                ? double.parse(val)
                                : 0;
                          });
                        },
                        autoFormatNumber: true,
                        isFiat: true,
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
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
                SizedBox(height: height / 50),
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
                        85,
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
                        keyboardtype: TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "proceedpayout".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      proceedPayoutCurrency = value.toString();
                    });
                  },
                  getPayoutCurrencies,
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Container(
                    width: width / 1.09,
                    child: checkBoxItem(
                      text: "Withholding Tax Disclosure",
                      value: withholdingTaxDisclosure,
                      onChanged: (bool? value) {
                        setState(() {
                          withholdingTaxDisclosure = !withholdingTaxDisclosure;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 15.0),
                    child: Text(
                      "assetsaleproceeds".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Container(
                  width: width,
                  child: Text(
                    "provideaccountdetails".tr(),
                    textAlign: TextAlign.left,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "selectbank".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      bankId = int.parse(value.toString());
                    });
                  },
                  getBanksList,
                  bankId == 0 ? null : bankId,
                  // null,
                  'pleaseselectbank'.tr(),
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
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "accountnumber".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "accountnumber".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: accountNumber,
                      onChanged: (value) {
                        setState(() {
                          accountNumber = value;
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
                          accountNumber = value!;
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
                      "beneficiaryname".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "beneficiaryname".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: beneficiaryName,
                      onChanged: (value) {
                        setState(() {
                          beneficiaryName = value;
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
                          beneficiaryName = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "primarybuyerrequirement".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Exempted Countries".tr(),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
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
              SizedBox(height: height / 50),
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
                        elevation: MaterialStateProperty.all<double>(0),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            exemptedCountries.length > 0
                                ? iso2Countries[exemptedCountries.last] ?? ""
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
              SizedBox(height: height / 50),
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
                          horizontal: 20.0,
                          vertical: 15.0,
                        ),
                        child: Column(
                          children: [
                            Container(
                              width: width / 1.8,
                              child: Wrap(
                                alignment: WrapAlignment.center,
                                children: [
                                  if (exemptedCountries.isNotEmpty) ...[
                                    for (var country in exemptedCountries) ...[
                                      userItem(
                                        iso2Countries[country] ?? "",
                                        () {
                                          setState(() {
                                            exemptedCountries.remove(country);
                                          });
                                        },
                                        foreColor: wihitecolor,
                                        backColor: notifier.getbluebackcolor,
                                      ),
                                    ],
                                  ] else ...[
                                    Text(
                                      "nameofexemptedcountriesappearhere".tr(),
                                      textAlign: TextAlign.center,
                                      style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody,
                                        fontSize: 15.sp,
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                            ),
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
              SizedBox(height: height / 50),
              Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: SizedBox(
                      height: 25,
                      child: Row(
                        children: [
                          Radio(
                            value: hasAdditionalKYCRequirements,
                            groupValue: true,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateProperty.all(
                              notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) {
                              setState(() {
                                hasAdditionalKYCRequirements = true;
                              });
                            },
                          ),
                          Text(
                            'yes'.tr(),
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  SizedBox(height: 10),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: SizedBox(
                      height: 25,
                      child: Row(
                        children: [
                          Radio(
                            value: hasAdditionalKYCRequirements,
                            groupValue: false,
                            fillColor: MaterialStateProperty.all(
                              notifier.getbluewhitecolor,
                            ),
                            activeColor: notifier.getbluewhitecolor,
                            onChanged: (value) {
                              setState(() {
                                hasAdditionalKYCRequirements = false;
                              });
                            },
                          ),
                          Text(
                            'no'.tr(),
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                  if (hasAdditionalKYCRequirements) ...[
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Container(
                            width: 340,
                            child: Text(
                              "Describe additional KYC requirements (Describe specific documents or identity checks required)",
                              style: TextStyle(
                                fontSize: 12,
                                fontFamily: fontsemibold,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 70),
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Container(
                      width: 340,
                      child: Text(
                        "Investor Category Eligibilty (Define categories of eligible investors)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      investorCategory = value.toString();
                    });
                  },
                  getInvestorCategory,
                  investorCategory.isEmpty ? null : investorCategory,
                  'Select option',
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
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Container(
                      width: 340,
                      child: Text(
                        "Minimum KYC Tier Required (Select platform verification tier requred to invest)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      minimumKycTier = value.toString();
                    });
                  },
                  getKYCTiers,
                  minimumKycTier.isEmpty ? null : minimumKycTier,
                  'Select option',
                  context,
                  null,
                  validator: (value) {
                    // if (value == null || value.toString().isEmpty) {
                    //   return "fieldcannotbeempty".tr();
                    // }
                    return null;
                  },
                ),
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Container(
                      width: 340,
                      child: Text(
                        "Investment Suitability Disclaimer",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Container(
                    width: width / 1.09,
                    child: checkBoxItem(
                      text:
                          "I acknowledge that I have reviewed the suitability criteria",
                      value: acknowledgedSuitabilityCriteria,
                      onChanged: (bool? value) {
                        setState(() {
                          acknowledgedSuitabilityCriteria =
                              !acknowledgedSuitabilityCriteria;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Final Declaration and Attestation",
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Authorized representative name",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: authorizedRepresentativeName,
                      onChanged: (value) {
                        setState(() {
                          authorizedRepresentativeName = value;
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
                          authorizedRepresentativeName = value!;
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
                      "Position/Title of authorized Representative",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: authorizedRepresentativeTitleOrPosition,
                      onChanged: (value) {
                        setState(() {
                          authorizedRepresentativeTitleOrPosition = value;
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
                          authorizedRepresentativeTitleOrPosition = value!;
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
                      "Contact email of authorized representative",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: CustomTextFormField.textField(
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: authorizedRepresentativeEmail,
                      onChanged: (value) {
                        setState(() {
                          authorizedRepresentativeEmail = value;
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
                          authorizedRepresentativeEmail = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 70),
              FormField(
                builder: (state) {
                  return Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      Transform.scale(
                        scale: 1,
                        child: Checkbox(
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(Radius.circular(5)),
                          ),
                          activeColor: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
                          side: BorderSide(
                            color: notifier.isDark
                                ? notifier.getbluecolor50
                                : notifier.getbluecolor90,
                          ),
                          value: attestInformationAccurateAndVerifiable,
                          onChanged: (bool? value) {
                            setState(() {
                              attestInformationAccurateAndVerifiable = value!;
                            });
                          },
                        ),
                      ),
                      Container(
                        width: width / 1.2,
                        child: RichText(
                          text: TextSpan(
                            text:
                                "I attest that the information provided is accurate and verifiable",
                            style: TextStyle(
                              fontSize: 13.sp,
                              fontFamily: fontbody,
                              color:
                                  state.hasError &&
                                      !attestInformationAccurateAndVerifiable
                                  ? Colors.red
                                  : notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                      ),
                    ],
                  );
                },
                validator: (value) {
                  if (!attestInformationAccurateAndVerifiable) {
                    return '';
                  }

                  return null;
                },
              ),
              SizedBox(height: height / 30),
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

                    if (salesEnd == null || salesEnd!.isBefore(salesStart!)) {
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
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom,
                ),
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
            exemptedCountries.add(country.countryCode);
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
      var newData = {...data as Map};
      newData['numberOfTokenToBeIssued'] = numberOfTokenToBeIssued;
      newData['assetCode'] = assetCode;
      newData['assetName'] = assetName;
      newData['accountNumber'] = accountNumber;
      newData['beneficiaryName'] = beneficiaryName;
      newData['bankId'] = bankId;
      newData['salesStart'] = DateFormat(
        "yyyy-MM-ddTHH:mm:ss.SSSSSS'Z'",
      ).format(salesStart!.toUtc());
      newData['salesEnd'] = DateFormat(
        "yyyy-MM-ddTHH:mm:ss.SSSSSS'Z'",
      ).format(salesEnd!);
      newData['capOnPurchase'] = capOnPurchase ? 1 : 0;
      newData['capQuantity'] = capQuantity;
      newData['capAmountInFiat'] = capAmountInFiat;
      newData['capDurationInDays'] = capDurationInDays;
      newData['proceedCycle'] = proceedCycle;
      newData['walletToHoldAssetsNotForSale'] = walletToHoldAssetsNotForSale;
      newData['assetLogo'] = assetLogo;
      newData['minimumKycTier'] = minimumKycTier;
      newData['investorCategory'] = investorCategory;
      newData['exemptedCountries'] = exemptedCountries.join(',');
      newData['hasAdditionalKYCRequirements'] = hasAdditionalKYCRequirements
          ? 1
          : 0;
      newData['assetQuoteCurrency'] = assetQuoteCurrency;
      newData['proceedPayoutCurrency'] = proceedPayoutCurrency;
      newData['proceedPayoutType'] = proceedPayoutType;
      newData['additionalKYCRequirements'] = additionalKYCRequirements;
      newData['investorAccreditationRequired'] = investorAccreditationRequired
          ? 1
          : 0;
      newData['attestInformationAccurateAndVerifiable'] =
          attestInformationAccurateAndVerifiable ? 0 : 1;
      newData['acknowledgedSuitabilityCriteria'] =
          acknowledgedSuitabilityCriteria ? 0 : 1;
      newData['withholdingTaxDisclosure'] = withholdingTaxDisclosure ? 1 : 0;
      newData['tokenizationFeeId'] = tokenizationFeeId;
      newData['authorizedRepresentativeName'] = authorizedRepresentativeName;
      newData['authorizedRepresentativeTitleOrPosition'] =
          authorizedRepresentativeTitleOrPosition;
      newData['authorizedRepresentativeEmail'] = authorizedRepresentativeEmail;

      inspect(newData);

      String requestBody = jsonEncode(newData);

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        Navigator.of(context).pop();
        await refreshCurrentTokenizationInfo();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
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

      if (responseData['statusCode'] == 200) {
        appState.viewData = responseData['data'];
        await inspect(appState.viewData);
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Future<void> uploadAssetLogo(PlatformFile file) async {
    try {
      showLoader(context);
      Map responseData = await makePutRequestForMultipartDocumentUpload(
        uri: '/v1/tokenization/logo',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: appState.primaryWallet.signer!,
        file: file,
        tokenizedAssetId: appState.viewData!['id'],
        documentTitle: "",
        documentType: "",
      );

      if (responseData['statusCode'] == 200) {
        setState(() {
          assetLogo = responseData['data'].toString().replaceAll('\"', '');
        });
        hideLoader(context);
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
        hideLoader(context);
      }
    } catch (e) {
      hideLoader(context);
      popup(
        context,
        title: "error".tr(),
        message: "Sorry, something went wrong. Please try again.",
      );
    }
  }

  Widget checkBoxItem({
    required String text,
    required bool? value,
    required void Function(bool? value) onChanged,
  }) {
    return Container(
      width: width / 1.08,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          Transform.scale(
            scale: 1.sp,
            child: Checkbox(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.all(Radius.circular(5.sp)),
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
                fontFamily: fontbody,
              ),
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
          border: OutlineInputBorder(borderRadius: BorderRadius.circular(15)),
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
    var fiatPercentage = appState
        .tokenizationData["tokenizationFees"][index]['feeFiatPercentage'];
    var assetPercentage = appState
        .tokenizationData["tokenizationFees"][index]['feeAssetPercentage'];
    var fiatFeeCap = double.parse(
      appState.tokenizationData["tokenizationFees"][index]['feeFiatCap']
          .toString(),
    );
    var tokenFee = (numberOfTokenToBeIssued * assetPercentage) / 100;
    var fiatFee = (data['assetCurrentValue'] * fiatPercentage) / 100;
    return "${appState.tokenizationData["tokenizationFees"][index]['feeDescription']} (\$${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(double.parse(tokenFee.toString()))} ${assetCode}).";
  }

  // Future<void> getImage() async {
  //   var image = await ImagePicker().pickImage(source: ImageSource.gallery);
  //   if (image != null) {
  //     var imageBase64Uncompressed = await getBase64Image(image);
  //     await StoreData().storeInsertData('image', imageBase64Uncompressed);
  //     setState(() {
  //       assetLogo = imageBase64Uncompressed;
  //     });
  //   }
  // }

  Future<void> getFile() async {
    FilePickerResult? result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['jpg', 'jpeg', 'gif', 'png', 'pdf'],
      withData: true,
    );

    if (result == null) {
      return null;
    }

    PlatformFile file = result.files.single;

    uploadAssetLogo(file);
  }

  Future<dynamic> getBase64Image(XFile image) async {
    //
    List<int> imageBytes = await image.readAsBytes();
    String imageB64 = base64Encode(imageBytes);
    return imageB64;
  }
}
