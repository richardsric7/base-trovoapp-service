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

class HybridMutualFundsAssetInformationView extends StatefulWidget {
  const HybridMutualFundsAssetInformationView({Key? key}) : super(key: key);

  @override
  State<HybridMutualFundsAssetInformationView> createState() =>
      _HybridMutualFundsAssetInformationView();
}

class _HybridMutualFundsAssetInformationView
    extends State<HybridMutualFundsAssetInformationView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;

  late String assetName;
  bool assetAlreadyExists = false;
  bool formHasError = false;
  late dynamic data = {};

  late String FundName;
  late String FundType;
  late String FundStructure;
  late String AssetManagementCompanyName;
  late String
  FundManagers; // Names & profiles as a single string or serialized JSON
  late String RegulatoryLicenseNumber;
  late String IsinOrSecFundCode;
  DateTime? FundLaunchDate;
  double? TotalExpenseRatio; // % value
  double? ExitLoadOrRedemptionFee; // % value
  int? Tenure; // Duration in years
  double? InitialNetAssetValue;
  late String NavUpdateFrequency;
  late String NavCalculationMethod;
  double? MinimumInvestmentAmount;
  late String RedemptionRules;
  late String LockInPeriod;
  double? EntryLoad;
  double? PerformanceFee;
  late String DividendPolicy;
  late String FundRiskRating;
  late String LiquidityProfile;
  late String DistributionFrequency;
  late String DistributionMethod;
  late String BenchmarkComparisonMethod;
  late String FeeBreakdownSummary;
  late String InvestmentObjective;
  late String AssetAllocation;
  late String TopHoldings;
  double? AverageMaturity;
  double? YieldToMaturity;
  late String CreditRatingProfile;
  late String
  LockInPeriodPortfolio; // Another Lock-In Period field in portfolio section
  double? PerformanceFeeIfAny;
  late String EquityStrategy;
  late String MarketCapitalizationFocus;
  late String BenchmarkIndex;
  late String SectorExposureLimits;
  late String TopEquityHoldings;
  late String GeographicExposure;
  late String RiskProfile;
  double? VolatilityEstimate;
  late String TargetAllocation;
  late String AllowedAllocationRange;
  late String AssetClassesIncluded;
  late String RebalancingFrequency;
  late String BenchmarkIndexComposite;
  late String TopEquityHoldingsList;
  late String TopDebtHoldings;
  double? DividendYield;
  late String CreditRatingDistribution;
  double? AverageMaturityHybrid;
  double? YieldToMaturityHybrid;
  late String TrusteeName;
  late String Custodian;
  late String Auditor;
  late String FundAdministrator;
  late String LegalAdvisor;
  late String RatingAgency;
  late String InvestmentCommitteeMembers;

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

  List<DropdownMenuItem<String>> get getAssetProtectionOptions {
    List<DropdownMenuItem<String>> assetProtectionOptions = [];
    var data = appState.tokenizationData['assetProtectionOptions'];
    for (var i = 0; i < data.length; i++) {
      assetProtectionOptions.add(
        DropdownMenuItem(
          child: Text(data![i]['id'], overflow: TextOverflow.ellipsis),
          value: data![i]['id'],
        ),
      );
    }
    return assetProtectionOptions;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;
    inspect(data);
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
                "Hybrid Mutual Funds Information",
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Fund Name",
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
                      "Enter fund name",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: FundName,
                      onChanged: (value) {
                        setState(() {
                          FundName = value;
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
                          FundName = value!;
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
                      "Fund Type",
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
                      "Enter fund type",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: FundType,
                      onChanged: (value) {
                        setState(() {
                          FundType = value;
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
                          FundType = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15),
                child: Container(
                  width: width,
                  child: Text(
                    textAlign: TextAlign.left,
                    "Fund Structure",
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ),
              SizedBox(height: height / 70),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      FundStructure = value.toString();
                    });
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
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
                    child: Text(
                      "Asset Management Company Name",
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
                      "Enter fund type",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: AssetManagementCompanyName,
                      onChanged: (value) {
                        setState(() {
                          AssetManagementCompanyName = value;
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
                          AssetManagementCompanyName = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    SizedBox(
                      width: 300,
                      child: Text(
                        "Fund Manager(s) (Names and profiles of managing team)",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    TextButton(
                      onPressed: addMilestone,
                      style: ButtonStyle(
                        padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                        minimumSize: WidgetStatePropertyAll(Size.zero),
                      ),
                      child: Icon(Icons.add_circle, size: 20),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 70),
              if (true) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'Add name and profiles',
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
              SizedBox(height: height / 70),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Container(
                      width: width / 1.12,
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(10),
                        ),
                        color: notifier.getaddsubwalletgrey,
                      ),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 10,
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  "Land acquisition",
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                                TextButton(
                                  onPressed: () {},
                                  style: ButtonStyle(
                                    padding: WidgetStatePropertyAll(
                                      EdgeInsets.all(7),
                                    ),
                                    tapTargetSize:
                                        MaterialTapTargetSize.shrinkWrap,
                                    minimumSize: WidgetStatePropertyAll(
                                      Size.zero,
                                    ),
                                  ),
                                  child: Icon(
                                    CupertinoIcons.trash,
                                    size: 15,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: 10),
                            Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  "January, 2025",
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontsemibold,
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
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Regulatory License No. (Issued by SEC or relevant authority)",
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
                    child: CustomTextFormField.textField(
                      "Enter license no.",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: RegulatoryLicenseNumber,
                      onChanged: (value) {
                        setState(() {
                          RegulatoryLicenseNumber = value;
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
                          RegulatoryLicenseNumber = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "ISIN / SEC Fund Code (Unique identifier assigned by regulatory authority)",
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
                    child: CustomTextFormField.textField(
                      "Enter code",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: IsinOrSecFundCode,
                      onChanged: (value) {
                        setState(() {
                          IsinOrSecFundCode = value;
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
                          IsinOrSecFundCode = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Fund Launch Date (Official date fund began operations)",
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
              ButtonOutlined(
                // fundLaunchDate != null
                //     ? DateFormat('MMMM dd, yyyy').format(fundLaunchDate!)
                //     : "ends".tr(),
                "Select launch date",
                notifier.getwihitecolor,
                notifier.getgrey,
                borderColor: notifier.getgrey,
                width: 320,
                height: 50.sp,
                onTap: () {
                  showDatePicker(
                    context: context,
                    initialDate: DateTime.now(),
                    firstDate: DateTime.fromMicrosecondsSinceEpoch(1000),
                    lastDate: DateTime.now().add(Duration(days: 730)),
                  ).then(
                    (value) => {
                      setState(() {
                        FundLaunchDate = value;
                      }),
                    },
                  );
                },
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Total Expense Ratio (TER) (Annual % of fund’s operating expenses)",
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
                    child: CustomTextFormField.textField(
                      "Enter %",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: TotalExpenseRatio,
                      onChanged: (value) {
                        setState(() {
                          TotalExpenseRatio = value;
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
                          TotalExpenseRatio = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Exit Load / Redemption Fee (Fee charged on redemption within specific period)",
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
                    child: CustomTextFormField.textField(
                      "Enter %",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      width / 1.12,
                      onChanged: (value) {
                        setState(() {
                          ExitLoadOrRedemptionFee = value;
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        ExitLoadOrRedemptionFee = double.parse(
                          value!.toString(),
                        );
                      },
                      autoFormatNumber: true,
                      isFiat: true,
                      controller: valueOfAssetController,
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
                      "Fund Structure and Operations",
                      style: TextStyle(
                        fontSize: 18,
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
                  Container(
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "Provide information about the fund structure of the fund",
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Tenure (Fund duration (e.g. 3 years))",
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
                    child: CustomTextFormField.textField(
                      "Enter %",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: Tenure,
                      onChanged: (value) {
                        setState(() {
                          Tenure = value;
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
                          Tenure = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Initial Net Asset Value (NAV)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: InitialNetAssetValue,
                      onChanged: (value) {
                        setState(() {
                          InitialNetAssetValue = value;
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
                          InitialNetAssetValue = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "NAV Update Frequency",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: NavUpdateFrequency,
                      onChanged: (value) {
                        setState(() {
                          NavUpdateFrequency = value;
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
                          NavUpdateFrequency = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "NAV Calculation Method",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: NavCalculationMethod,
                      onChanged: (value) {
                        setState(() {
                          NavCalculationMethod = value;
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
                          NavCalculationMethod = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Minimum Investment Amount (e.g. ₦1,000)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: MinimumInvestmentAmount,
                      onChanged: (value) {
                        setState(() {
                          MinimumInvestmentAmount = value;
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
                          MinimumInvestmentAmount = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Redemption Rules (Anytime or Specify)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: RedemptionRules,
                      onChanged: (value) {
                        setState(() {
                          RedemptionRules = value;
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
                          RedemptionRules = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Lock-In Period (For ELSS or special structures)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: LockInPeriod,
                      onChanged: (value) {
                        setState(() {
                          LockInPeriod = value;
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
                          LockInPeriod = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Entry Load (Fee charged at time of purchase)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: EntryLoad,
                      onChanged: (value) {
                        setState(() {
                          EntryLoad = value;
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
                          EntryLoad = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Performance Fee",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: PerformanceFee,
                      onChanged: (value) {
                        setState(() {
                          PerformanceFee = value;
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
                          PerformanceFee = value!;
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
                      "Dividend Policy",
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
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: DividendPolicy,
                      onChanged: (value) {
                        setState(() {
                          DividendPolicy = value;
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
                          DividendPolicy = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Fund Risk Rating (External)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: FundRiskRating,
                      onChanged: (value) {
                        setState(() {
                          FundRiskRating = value;
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
                          FundRiskRating = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Liquidity Profile",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: LiquidityProfile,
                      onChanged: (value) {
                        setState(() {
                          LiquidityProfile = value;
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
                          LiquidityProfile = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Distribution Frequency",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: DistributionFrequency,
                      onChanged: (value) {
                        setState(() {
                          DistributionFrequency = value;
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
                          DistributionFrequency = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Distribution Method",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: DistributionMethod,
                      onChanged: (value) {
                        setState(() {
                          DistributionMethod = value;
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
                          DistributionMethod = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Benchmark Comparison Method",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: BenchmarkComparisonMethod,
                      onChanged: (value) {
                        setState(() {
                          BenchmarkComparisonMethod = value;
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
                          BenchmarkComparisonMethod = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Fee Breakdown Summary",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: FeeBreakdownSummary,
                      onChanged: (value) {
                        setState(() {
                          FeeBreakdownSummary = value;
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
                          FeeBreakdownSummary = value!;
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
                      "Portfolio Composition & Strategy",
                      style: TextStyle(
                        fontSize: 18,
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
                  Container(
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "provideassetvalueinfo",
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Investment Objective",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: InvestmentObjective,
                      onChanged: (value) {
                        setState(() {
                          InvestmentObjective = value;
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
                          InvestmentObjective = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Asset Allocation (% in T-bills, Bonds, CPs, etc.)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: AssetAllocation,
                      onChanged: (value) {
                        setState(() {
                          AssetAllocation = value;
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
                          AssetAllocation = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Top Holdings",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: TopHoldings,
                      onChanged: (value) {
                        setState(() {
                          TopHoldings = value;
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
                          TopHoldings = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Average Maturity",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: AverageMaturity,
                      onChanged: (value) {
                        setState(() {
                          AverageMaturity = value;
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
                          AverageMaturity = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Yield to Maturity (YTM)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: YieldToMaturity,
                      onChanged: (value) {
                        setState(() {
                          YieldToMaturity = value;
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
                          YieldToMaturity = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Credit Rating Profile",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: CreditRatingProfile,
                      onChanged: (value) {
                        setState(() {
                          CreditRatingProfile = value;
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
                          CreditRatingProfile = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Lock-in Period",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: LockInPeriodPortfolio,
                      onChanged: (value) {
                        setState(() {
                          LockInPeriodPortfolio = value;
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
                          LockInPeriodPortfolio = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Performance Fee (if any)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: PerformanceFeeIfAny,
                      onChanged: (value) {
                        setState(() {
                          PerformanceFeeIfAny = value;
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
                          PerformanceFeeIfAny = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Equity Strategy",
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
                      EquityStrategy = value.toString();
                    });
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Market Capitalization Focus",
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
                      MarketCapitalizationFocus = value.toString();
                    });
                  },
                  [],
                  null,
                  'Select type',
                  context,
                  null,
                  validator: (value) {
                    // if (selectedAssetSectorId.isEmpty) {
                    //   return "pleaseselectassetsector".tr();
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Benchmark Index (E.g., NGX ASI, MSCI Frontier Markets)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: BenchmarkIndex,
                      onChanged: (value) {
                        setState(() {
                          BenchmarkIndex = value;
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
                          BenchmarkIndex = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Sector Exposure Limits (Optional regulatory or internal limits)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: SectorExposureLimits,
                      onChanged: (value) {
                        setState(() {
                          SectorExposureLimits = value;
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
                          SectorExposureLimits = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Top Holdings",
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
                      'Enter value',
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: TopEquityHoldings,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          TopEquityHoldings = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Geographic Exposure (% allocation by country/region)",
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
                      'Enter value',
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: GeographicExposure,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          GeographicExposure = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Risk Profile",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: RiskProfile,
                      onChanged: (value) {
                        setState(() {
                          RiskProfile = value;
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
                          RiskProfile = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Volatility Estimate (Standard deviation or beta)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: VolatilityEstimate,
                      onChanged: (value) {
                        setState(() {
                          VolatilityEstimate = value;
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
                          VolatilityEstimate = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Target Allocation (e.g 60% Equity / 40% Debt)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: TargetAllocation,
                      onChanged: (value) {
                        setState(() {
                          TargetAllocation = value;
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
                          TargetAllocation = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Allowed Allocation Range (E.g. 50–70% Equity, 30–50% Debt)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: AllowedAllocationRange,
                      onChanged: (value) {
                        setState(() {
                          AllowedAllocationRange = value;
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
                          AllowedAllocationRange = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Asset Classes Included Instruments (Equity, Bonds, T-bills, REITs, etc.)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: AssetClassesIncluded,
                      onChanged: (value) {
                        setState(() {
                          AssetClassesIncluded = value;
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
                          AssetClassesIncluded = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Rebalancing Frequency (Portfolio reallocation frequency)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: RebalancingFrequency,
                      onChanged: (value) {
                        setState(() {
                          RebalancingFrequency = value;
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
                          RebalancingFrequency = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Benchmark Index (Composite or blended index)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: BenchmarkIndexComposite,
                      onChanged: (value) {
                        setState(() {
                          BenchmarkIndexComposite = value;
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
                          BenchmarkIndexComposite = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Top Equity Holdings (Top 5–10 stocks)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: TopEquityHoldingsList,
                      onChanged: (value) {
                        setState(() {
                          TopEquityHoldingsList = value;
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
                          TopEquityHoldingsList = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Top Debt Holdings (Top 5–10 debt instruments)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: TopDebtHoldings,
                      onChanged: (value) {
                        setState(() {
                          TopDebtHoldings = value;
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
                          TopDebtHoldings = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Dividend Yield (Equity)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: DividendYield,
                      onChanged: (value) {
                        setState(() {
                          DividendYield = value;
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
                          DividendYield = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Credit Rating Distribution (Debt) (AAA, AA, etc.)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: CreditRatingDistribution,
                      onChanged: (value) {
                        setState(() {
                          CreditRatingDistribution = value;
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
                          CreditRatingDistribution = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Average Maturity (Hybrid)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: AverageMaturityHybrid,
                      onChanged: (value) {
                        setState(() {
                          AverageMaturityHybrid = value;
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
                          AverageMaturityHybrid = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Yield to Maturity (Hybrid)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: YieldToMaturityHybrid,
                      onChanged: (value) {
                        setState(() {
                          YieldToMaturityHybrid = value;
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
                          YieldToMaturityHybrid = value!;
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
                      "Key Entities Involved",
                      style: TextStyle(
                        fontSize: 18,
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
                  Container(
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "provideassetvalueinfo",
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Trustee Name",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: TrusteeName,
                      onChanged: (value) {
                        setState(() {
                          TrusteeName = value;
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
                          TrusteeName = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Custodian",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: Custodian,
                      onChanged: (value) {
                        setState(() {
                          Custodian = value;
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
                          Custodian = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Auditor",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: Auditor,
                      onChanged: (value) {
                        setState(() {
                          Auditor = value;
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
                          Auditor = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Fund Administrator  (If different from AMC)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: FundAdministrator,
                      onChanged: (value) {
                        setState(() {
                          FundAdministrator = value;
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
                          FundAdministrator = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Legal Advisor  (Compliance/legal counsel)",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: LegalAdvisor,
                      onChanged: (value) {
                        setState(() {
                          LegalAdvisor = value;
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
                          LegalAdvisor = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Rating Agency",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: RatingAgency,
                      onChanged: (value) {
                        setState(() {
                          RatingAgency = value;
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
                          RatingAgency = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Investment Committee Members",
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
                    child: CustomTextFormField.textField(
                      "Enter value",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      85,
                      300.sp,
                      initialValue: InvestmentCommitteeMembers,
                      onChanged: (value) {
                        setState(() {
                          InvestmentCommitteeMembers = value;
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
                          InvestmentCommitteeMembers = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 30),
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

  void submitForm() async {
    try {
      showLoader(context);
      var newData = {...data as Map};

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

  void addMilestone() {
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
                        Text(
                          "Add Milestone",
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
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
                        child: Text(
                          "Milestone",
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
                          "Enter milestone",
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          // initialValue: assetPhysicalAddress,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {});
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
                          "Select Date",
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
                          70.sp,
                          width / 1.12,
                          // initialValue: assetPhysicalAddress,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onSaved: (value) {
                            setState(() {});
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 60),
                ],
              ),
            );
          },
        );
      },
    );
  }
}
