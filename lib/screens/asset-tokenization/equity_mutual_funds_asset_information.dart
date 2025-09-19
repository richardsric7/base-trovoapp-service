import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
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
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class EquityMutualFundsAssetInformationView extends StatefulWidget {
  const EquityMutualFundsAssetInformationView({Key? key}) : super(key: key);

  @override
  State<EquityMutualFundsAssetInformationView> createState() =>
      _EquityMutualFundsAssetInformationView();
}

class _EquityMutualFundsAssetInformationView
    extends State<EquityMutualFundsAssetInformationView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  bool formHasError = false;
  late dynamic data = {};

  var formData = {};

  final valueOfAssetController = TextEditingController();
  final miscCostOfAssetController = TextEditingController();
  final percentageFromPromotersController = TextEditingController();
  final assetOwnerRetainedOrContributedValueController =
      TextEditingController();
  final percentValueOfInsuranceController = TextEditingController();

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getFundStructureOptions {
    List<DropdownMenuItem<String>> options = [];
    var data = ['Open-ended', 'Close-ended', 'Interval Fund'];
    for (var i = 0; i < data.length; i++) {
      options.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }

    return options;
  }

  List<DropdownMenuItem<String>> get getEquityStrategyOptions {
    List<DropdownMenuItem<String>> options = [];
    var data = ['Value', 'Growth', 'Income', 'Thematic', 'Index'];
    for (var i = 0; i < data.length; i++) {
      options.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }

    return options;
  }

  List<DropdownMenuItem<String>> get getMarketCapitalizationFocusOptions {
    List<DropdownMenuItem<String>> options = [];
    var data = ['Large-cap', 'Mid-cap', 'Small-cap', 'Mixed'];
    for (var i = 0; i < data.length; i++) {
      options.add(
        DropdownMenuItem(
          child: Text(data[i], overflow: TextOverflow.ellipsis),
          value: data[i],
        ),
      );
    }

    return options;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;
    inspect(data);

    super.initState();
    getdarkmodepreviousstate();

    formData = {
      "assetName": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Fund Name',
        'placeholderText': 'Trovo Funds',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "assetType": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Fund Type',
        'placeholderText': 'Equity funds',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "fundStructure": {
        'type': 'string',
        'widgetType': 'dropdown',
        'label': 'Fund Structure',
        'placeholderText': 'Select fund structure',
        'value': '',
        'options': ['Open-ended', 'Close-ended', 'Interval Fund'],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "assetManagementCompanyName": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Asset Management Company Name',
        'placeholderText': 'e.g. Trovo Fund Management',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "fundManagers": {
        'type': 'string',
        'widgetType': 'list',
        'label': 'Fund Manager(s) (Names and profiles of managing team)',
        'placeholderText': 'Add fund manager',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "regulatoryLicenseNumber": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Regulatory License No. (Issued by SEC or relevant authority)',
        'placeholderText': 'e.g. 112233445566',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "isinOrSecFundCode": {
        'type': 'string',
        'widgetType': 'text',
        'label':
            'ISIN / SEC Fund Code (Unique identifier assigned by regulatory authority)',
        'placeholderText': 'e.g. 112233445566',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "fundLaunchDate": {
        'type': 'string',
        'widgetType': 'datetime',
        'label': 'Fund Launch Date (Official date fund began operations)',
        'placeholderText': "Select launch date",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "totalExpenseRatio": {
        'type': 'double',
        'widgetType': 'text',
        'label':
            'Total Expense Ratio (TER) (Annual % of fund’s operating expenses)',
        'placeholderText': "e.g. 12.5",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "exitLoadRedemptionFee": {
        'type': 'double',
        'widgetType': 'text',
        'label':
            'Exit Load / Redemption Fee (Fee charged on redemption within specific period)',
        'placeholderText': "e.g. 1,000",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "section2": {
        'type': 'string',
        'widgetType': 'section',
        'label': 'Fund Structure and Operations',
        'placeholderText':
            "Provide information about the fund structure of the fund",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "tenure": {
        'type': 'int',
        'widgetType': 'text',
        'label': 'Tenure (Fund duration (e.g. 3 years))',
        'placeholderText': "e.g. 3",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "initialNetAssetValue": {
        'type': 'double',
        'widgetType': 'text',
        'label': 'Initial Net Asset Value (NAV)',
        'placeholderText': "e.g. 1,000,000,000",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "navUpdateFrequency": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'NAV Update Frequency',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "navCalculationMethod": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'NAV Calculation Method',
        'placeholderText': "Enter %",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "minimumInvestmentAmount": {
        'type': 'double',
        'widgetType': 'text',
        'label': 'Minimum Investment Amount (e.g. ₦1,000)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "redemptionRules": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Redemption Rules (Anytime or Specify)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "lockInPeriod": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Lock-In Period (For ELSS or special structures)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "entryLoad": {
        'type': 'double',
        'widgetType': 'text',
        'label': 'Entry Load (Fee charged at time of purchase)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "performanceFee": {
        'type': 'double',
        'widgetType': 'text',
        'label': 'Performance Fee',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "dividendPolicy": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Dividend Policy',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "liquidityProfile": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Liquidity Profile',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "distributionFrequency": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Distribution Frequency',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "distributionMethod": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Distribution Method',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "benchmarkComparisonMethod": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Benchmark Comparison Method',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "feeBreakdownSummary": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Fee Breakdown Summary',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "section3": {
        'type': 'string',
        'widgetType': 'section',
        'label': 'Portfolio Composition & Strategy',
        'placeholderText': "provideassetvalueinfo",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "investmentObjective": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Investment Objective',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "equityStrategy": {
        'type': 'string',
        'widgetType': 'dropdown',
        'label': 'Equity Strategy',
        'placeholderText': "Select type",
        'value': '',
        'options': ['Value', 'Growth', 'Income', 'Thematic', 'Index'],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "marketCapitalizationFocus": {
        'type': 'string',
        'widgetType': 'dropdown',
        'label': 'Market Capitalization Focus',
        'placeholderText': "Select type",
        'value': '',
        'options': ['Large-cap', 'Mid-cap', 'Small-cap', 'Mixed'],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "benchmarkIndex": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Benchmark Index (E.g., NGX ASI, MSCI Frontier Markets)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "sectorExposureLimits": {
        'type': 'string',
        'widgetType': 'text',
        'label':
            'Sector Exposure Limits (Optional regulatory or internal limits)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "topHoldings": {
        'type': 'string',
        'widgetType': 'list',
        'label': 'Top Holdings',
        'placeholderText': 'Add item',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "geographicExposure": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Geographic Exposure (% allocation by country/region)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "riskProfile": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Risk Profile',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "volatilityEstimate": {
        'type': 'double',
        'widgetType': 'text',
        'label': 'Volatility Estimate (Standard deviation or beta)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "dividendYield": {
        'type': 'double',
        'widgetType': 'text',
        'label': 'Dividend Yield (Equity)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "section4": {
        'type': 'string',
        'widgetType': 'section',
        'label': 'Key Entities Involved',
        'placeholderText': "provideassetvalueinfo",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "trusteeName": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Trustee Name',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "custodian": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Custodian',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "auditor": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Auditor',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "fundAdministrator": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Fund Administrator (If different from AMC)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "legalAdvisor": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Legal Advisor (Compliance/legal counsel)',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "ratingAgency": {
        'type': 'string',
        'widgetType': 'text',
        'label': 'Rating Agency',
        'placeholderText': "Enter value",
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
      "investmentCommitteeMembers": {
        'type': 'string',
        'widgetType': 'list',
        'label': 'Investment Committee Members',
        'placeholderText': 'Add committee member',
        'value': '',
        'options': [],
        'autoFormatNumber': false,
        'isFiat': false,
        'controller': null,
      },
    };

    formData['totalExpenseRatio']['value'] = double.parse(
      data['totalExpenseRatio'].toString(),
    );
    formData['exitLoadRedemptionFee']['value'] = double.parse(
      data['exitLoadRedemptionFee'].toString(),
    );
    formData['initialNetAssetValue']['value'] = double.parse(
      data['initialNetAssetValue'].toString(),
    );
    formData['entryLoad']['value'] = double.parse(data['entryLoad'].toString());
    formData['performanceFee']['value'] = double.parse(
      data['performanceFee'].toString(),
    );
    formData['minimumInvestmentAmount']['value'] = double.parse(
      data['minimumInvestmentAmount'].toString(),
    );
    formData['volatilityEstimate']['value'] = double.parse(
      data['volatilityEstimate'].toString(),
    );
    formData['dividendYield']['value'] = double.parse(
      data['dividendYield'].toString(),
    );
    formData['tenure']['value'] = int.parse(data['tenure'].toString());
    formData['topHoldings']['value'] = data['topHoldings'].toString();
    formData['fundManagers']['value'] = data['fundManagers'].toString();
    formData['investmentCommitteeMembers']['value'] =
        data['investmentCommitteeMembers'].toString();
    formData['assetName']['value'] = data['assetName'].toString();
    formData['assetType']['value'] = data['assetType'].toString();
    formData['fundStructure']['value'] = data['fundStructure']
        .toString()
        .nullIfEmpty();
    formData['assetManagementCompanyName']['value'] =
        data['assetManagementCompanyName'].toString();
    formData['regulatoryLicenseNumber']['value'] =
        data['regulatoryLicenseNumber'].toString();
    formData['isinOrSecFundCode']['value'] = data['isinOrSecFundCode']
        .toString();
    formData['navUpdateFrequency']['value'] = data['navUpdateFrequency']
        .toString();
    formData['navCalculationMethod']['value'] = data['navCalculationMethod']
        .toString();
    formData['redemptionRules']['value'] = data['redemptionRules'].toString();
    formData['lockInPeriod']['value'] = data['lockInPeriod'].toString();
    formData['dividendPolicy']['value'] = data['dividendPolicy'].toString();
    formData['liquidityProfile']['value'] = data['liquidityProfile'].toString();
    formData['distributionFrequency']['value'] = data['distributionFrequency']
        .toString();
    formData['distributionMethod']['value'] = data['distributionMethod']
        .toString();
    formData['benchmarkComparisonMethod']['value'] =
        data['benchmarkComparisonMethod'].toString();
    formData['feeBreakdownSummary']['value'] = data['feeBreakdownSummary']
        .toString();
    formData['investmentObjective']['value'] = data['investmentObjective']
        .toString();
    formData['equityStrategy']['value'] = data['equityStrategy']
        .toString()
        .nullIfEmpty();
    formData['marketCapitalizationFocus']['value'] =
        data['marketCapitalizationFocus'].toString().nullIfEmpty();
    formData['benchmarkIndex']['value'] = data['benchmarkIndex'].toString();
    formData['sectorExposureLimits']['value'] = data['sectorExposureLimits']
        .toString();
    formData['geographicExposure']['value'] = data['geographicExposure']
        .toString();
    formData['riskProfile']['value'] = data['riskProfile'].toString();
    formData['trusteeName']['value'] = data['trusteeName'].toString();
    formData['custodian']['value'] = data['custodian'].toString();
    formData['auditor']['value'] = data['auditor'].toString();
    formData['fundAdministrator']['value'] = data['fundAdministrator']
        .toString();
    formData['legalAdvisor']['value'] = data['legalAdvisor'].toString();
    formData['ratingAgency']['value'] = data['ratingAgency'].toString();
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
                "Equity Mutual Funds Information",
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    Text(
                      "basicinformation".tr(),
                      style: TextStyle(
                        fontSize: 18,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 50),
              for (var item in formData.entries) ...[
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      SizedBox(
                        width: 300,
                        child: Text(
                          item.value['label'],
                          style: TextStyle(
                            fontSize: item.value['widgetType'] == 'section'
                                ? 18
                                : 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                      if (item.value['widgetType'] == 'list') ...[
                        TextButton(
                          onPressed: () {
                            addMilestone(
                              label: 'Add ${item.value['label']}',
                              placeholder: item.value['placeholderText'],
                              onDone: (value) {
                                setState(() {
                                  var val =
                                      item.value['value'] == null ||
                                          item.value['value'].toString().isEmpty
                                      ? []
                                      : item.value['value'].toString().split(
                                          ',',
                                        );
                                  val.add(value);
                                  formData[item.key]['value'] = val.join(',');
                                });
                              },
                            );
                          },
                          style: ButtonStyle(
                            padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            minimumSize: WidgetStatePropertyAll(Size.zero),
                          ),
                          child: Icon(Icons.add_circle, size: 20),
                        ),
                      ],
                    ],
                  ),
                ),
                SizedBox(height: height / 70),
                getFormElement(item),
              ],
              Button(
                "saveandcontinuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  var form = _formKey.currentState;
                  if (form!.validate()) {
                    form.save();
                    submitForm();
                  }
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

  Widget getFormElement(MapEntry<dynamic, dynamic> item) {
    switch (item.value['widgetType']) {
      case 'section':
        return Column(
          children: [
            Row(
              children: [
                Container(
                  width: width,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      item.value['placeholderText'],
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(height: height / 50),
          ],
        );
      case 'list':
        var listItems =
            item.value['value'] == null ||
                item.value['value'].toString().isEmpty
            ? []
            : item.value['value'].toString().split(',');
        print('=========> this is slist items ${listItems}');
        return Column(
          children: [
            if (listItems.isEmpty) ...[
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      'This field is required',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: Colors.red,
                      ),
                    ),
                  ),
                ],
              ),
            ] else ...[
              for (var it in listItems) ...[
                SizedBox(height: height / 70),
                listItem(
                  item: it,
                  onDelete: (val) {
                    setState(() {
                      listItems.removeWhere((i) => i == val);
                      item.value['value'] = listItems.join(',');
                    });
                  },
                ),
              ],
            ],
            SizedBox(height: height / 50),
          ],
        );
      case 'dropdown':
        List<DropdownMenuItem<String>> options = [];
        var data = item.value['options'];
        for (var i = 0; i < data.length; i++) {
          options.add(
            DropdownMenuItem(
              child: Text(data[i], overflow: TextOverflow.ellipsis),
              value: data[i],
            ),
          );
        }
        return Column(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    formData[item.key]['value'] = value;
                  });
                },
                options,
                item.value['value'],
                item.value['placeholderText'],
                context,
                null,
                validator: (value) {
                  if (item.value['value'] == null) {
                    return "Please select an item";
                  }
                  return null;
                },
              ),
            ),
            SizedBox(height: height / 50),
          ],
        );
      default:
        if (item.value['type'] == 'double') {
          return Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: CustomTextFormField.textField(
                  item.value['placeholderText'],
                  notifier.getbluecolor,
                  null,
                  notifier.getgrey,
                  null,
                  notifier.getblck,
                  notifier.getgrey,
                  85,
                  300.sp,
                  validator: (value) {
                    if (value.isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                  onSaved: (value) {
                    setState(() {
                      formData[item.key]['value'] = double.parse(value);
                    });
                  },
                  autoFormatNumber: true,
                  controller: TextEditingController(
                    text: item.value['value'] == 0
                        ? ''
                        : formatNumberForInput(item.value['value']),
                  ),
                  keyboardtype: TextInputType.numberWithOptions(decimal: true),
                ),
              ),
            ],
          );
        }

        if (item.value['type'] == 'int') {
          return Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: CustomTextFormField.textField(
                  item.value['placeholderText'],
                  notifier.getbluecolor,
                  null,
                  notifier.getgrey,
                  null,
                  notifier.getblck,
                  notifier.getgrey,
                  85.sp,
                  300.sp,
                  validator: (value) {
                    if (value.isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                  onSaved: (value) {
                    setState(() {
                      formData[item.key]['value'] = int.tryParse(value) ?? 0;
                    });
                  },
                  autoFormatNumber: true,
                  controller: TextEditingController(
                    text: item.value['value'] == 0
                        ? ''
                        : item.value['value'].toString(),
                  ),
                  keyboardtype: TextInputType.numberWithOptions(decimal: true),
                ),
              ),
            ],
          );
        }

        return Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: CustomTextFormField.textField(
                item.value['placeholderText'],
                notifier.getbluecolor,
                null,
                notifier.getgrey,
                null,
                notifier.getblck,
                notifier.getgrey,
                85,
                300.sp,
                initialValue: item.value['value'].toString(),
                onSaved: (value) {
                  setState(() {
                    formData[item.key]['value'] = value;
                  });
                },
                validator: (value) {
                  if (value.isEmpty) {
                    return "fieldcannotbeempty".tr();
                  }
                  return null;
                },
              ),
            ),
          ],
        );
    }
  }

  void submitForm() async {
    try {
      showLoader(context);
      var newData = {...data as Map};

      newData['totalExpenseRatio'] = formData['totalExpenseRatio']['value'];
      newData['exitLoadRedemptionFee'] =
          formData['exitLoadRedemptionFee']['value'];
      newData['initialNetAssetValue'] =
          formData['initialNetAssetValue']['value'];
      newData['entryLoad'] = formData['entryLoad']['value'];
      newData['performanceFee'] = formData['performanceFee']['value'];
      newData['minimumInvestmentAmount'] =
          formData['minimumInvestmentAmount']['value'];
      newData['volatilityEstimate'] = formData['volatilityEstimate']['value'];
      newData['dividendYield'] = formData['dividendYield']['value'];
      newData['tenure'] = formData['tenure']['value'];
      newData['topHoldings'] = formData['topHoldings']['value'];
      newData['investmentCommitteeMembers'] =
          formData['investmentCommitteeMembers']['value'];
      newData['assetName'] = formData['assetName']['value'];
      newData['assetType'] = formData['assetType']['value'];
      newData['fundStructure'] = formData['fundStructure']['value'];
      newData['assetManagementCompanyName'] =
          formData['assetManagementCompanyName']['value'];
      newData['regulatoryLicenseNumber'] =
          formData['regulatoryLicenseNumber']['value'];
      newData['isinOrSecFundCode'] = formData['isinOrSecFundCode']['value'];
      newData['navUpdateFrequency'] = formData['navUpdateFrequency']['value'];
      newData['navCalculationMethod'] =
          formData['navCalculationMethod']['value'];
      newData['redemptionRules'] = formData['redemptionRules']['value'];
      newData['lockInPeriod'] = formData['lockInPeriod']['value'];
      newData['dividendPolicy'] = formData['dividendPolicy']['value'];
      newData['liquidityProfile'] = formData['liquidityProfile']['value'];
      newData['distributionFrequency'] =
          formData['distributionFrequency']['value'];
      newData['distributionMethod'] = formData['distributionMethod']['value'];
      newData['benchmarkComparisonMethod'] =
          formData['benchmarkComparisonMethod']['value'];
      newData['feeBreakdownSummary'] = formData['feeBreakdownSummary']['value'];
      newData['investmentObjective'] = formData['investmentObjective']['value'];
      newData['equityStrategy'] = formData['equityStrategy']['value'];
      newData['marketCapitalizationFocus'] =
          formData['marketCapitalizationFocus']['value'];
      newData['benchmarkIndex'] = formData['benchmarkIndex']['value'];
      newData['sectorExposureLimits'] =
          formData['sectorExposureLimits']['value'];
      newData['geographicExposure'] = formData['geographicExposure']['value'];
      newData['riskProfile'] = formData['riskProfile']['value'];
      newData['trusteeName'] = formData['trusteeName']['value'];
      newData['custodian'] = formData['custodian']['value'];
      newData['auditor'] = formData['auditor']['value'];
      newData['fundAdministrator'] = formData['fundAdministrator']['value'];
      newData['legalAdvisor'] = formData['legalAdvisor']['value'];
      newData['ratingAgency'] = formData['ratingAgency']['value'];

      String requestBody = jsonEncode(newData);
      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
        Navigator.of(context).pop();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
      hideLoader(context);
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
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
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
      child: ElevatedButton(
        onPressed: onClick,
        style: ButtonStyle(
          overlayColor: WidgetStateProperty.all<Color>(notifier.getsplashgrey),
          elevation: WidgetStateProperty.all<double>(0),
          backgroundColor: WidgetStateProperty.all<Color>(backColor),
          foregroundColor: WidgetStateProperty.all<Color>(
            notifier.getwihitecolor,
          ),
          side: WidgetStateProperty.all(
            BorderSide(color: borderColor, width: 1, style: BorderStyle.solid),
          ),
          shape: WidgetStateProperty.all<RoundedRectangleBorder>(
            const RoundedRectangleBorder(
              borderRadius: BorderRadius.all(Radius.circular(10)),
            ),
          ),
        ),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                color: foreColor,
                fontFamily: fontbody,
                fontSize: fontSize,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget CheckboxItem({
    required String label,
    required bool value,
    required void Function(bool?) onChanged,
    String? Function(Object?)? validator,
  }) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Transform.scale(
          scale: 1,
          child: FormField(
            builder: (state) {
              return SizedBox(
                width: 24,
                height: 24,
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
              );
            },
            validator: validator,
          ),
        ),
        SizedBox(width: 10),
        Container(
          width: width / 1.4,
          child: Text(
            label,
            overflow: TextOverflow.visible,
            style: TextStyle(
              fontSize: 15,
              color: formHasError && !value && validator != null
                  ? Colors.red
                  : notifier.getbluewhitecolor,
              fontFamily: fontbody,
            ),
          ),
        ),
      ],
    );
  }

  void addMilestone({
    required String label,
    required String placeholder,
    required void Function(String val) onDone,
  }) {
    String val = '';
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: notifier.getwihitecolor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
          top: Radius.circular(16),
        ), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.75,
          minChildSize: 0.25,
          expand: false,
          builder: (context, scrollController) {
            return SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  SizedBox(height: 10),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 14.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        SizedBox(
                          width: 250.sp,
                          child: Text(
                            label,
                            overflow: TextOverflow.visible,
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(Icons.cancel_outlined),
                          style: ButtonStyle(
                            padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            minimumSize: WidgetStatePropertyAll(Size.zero),
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: 20),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          placeholder,
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onChanged: (value) {
                            val = value;
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 40),
                  Button(
                    "Done",
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: () {
                      Navigator.of(context).pop();
                      onDone(val);
                    },
                  ),
                ],
              ),
            );
          },
        );
      },
    );
  }

  Widget listItem({
    required String item,
    required void Function(String item) onDelete,
  }) {
    return Row(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20),
          child: Container(
            width: width / 1.12,
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10)),
              color: notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        item,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          onDelete(item);
                        },
                        style: ButtonStyle(
                          padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                          minimumSize: WidgetStatePropertyAll(Size.zero),
                        ),
                        child: Icon(
                          CupertinoIcons.trash,
                          size: 15,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}
