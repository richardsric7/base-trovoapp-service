import 'dart:convert';
import 'dart:developer';
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
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class CommodityAssetInformationView extends StatefulWidget {
  const CommodityAssetInformationView({Key? key}) : super(key: key);

  @override
  State<CommodityAssetInformationView> createState() =>
      _CommodityAssetInformationView();
}

class _CommodityAssetInformationView
    extends State<CommodityAssetInformationView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  bool formHasError = false;
  late dynamic data = {};

  late String CommodityType;
  late String CommodityDescription;
  late int Quantity;
  late String QualityGrade;
  late String IssuerName;
  late String IssuerType;
  late String IssuerContactInfo;
  late String WarehouseName;
  late String WarehouseOperatorName;
  late String WarehouseLicenseNumber;
  late String WarehouseLocation;
  late String WRNumber;
  late DateTime? WRIssueDate;
  late DateTime? WRExpiryDate;
  late String WRSystemRegistration;
  late String WRRegistrationNumber;
  late String WRVerifier;
  late String StorageCondition;
  late String WarehouseAccreditationBody;
  late double MinimumPurchaseAmount;
  late DateTime? MaturityDate;
  late String AutoRollover;
  late String CurrentBeneficialOwner;
  late String WRCustodianName;
  late String OwnershipRightsRepresented;
  late String TrusteeOrThirdPartyOversight;
  late String LienOrEncumbrances;
  late double AssetValuation;
  late DateTime? ValuationDate;
  late String ValuationMethodology;
  late String TokenizationObjective;
  late int HoldingPeriod;
  late String RedemptionMechanism;
  late String PartiesInvolvedTrustee;
  late String PartiesInvolvedUnderwriter;
  late String PartiesInvolvedAssetManager;
  late String PartiesInvolvedLegalAdvisor;
  late String PartiesInvolvedAuditorVerifier;
  late String PartiesInvolvedRegulator;
  late String RisksMarketRisk;
  late String RisksStorageRisk;
  late String RisksTitleRisk;
  late String RisksFraudRisk;
  late String RisksInsuranceRisk;
  late String RisksOperationalRisk;
  late String RisksRegulatoryRisk;
  late String RisksLiquidityRisk;
  late String RisksForceMajeureRisk;
  late String RisksEarlyRedemptionRisk;
  late String RisksMitigationMeasures;
  late String RisksInsuranceCoverageSummary;
  late String RisksInsuranceProvider;
  late double RisksCoverageValue;

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
                "Commodity Asset Information",
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    SizedBox(
                      width: 350,
                      child: Text(
                        "Commodity, Issuer & Warehouse Information",
                        style: TextStyle(
                          fontSize: 18,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Commodity Type",
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
                      initialValue: CommodityType,
                      onChanged: (value) {
                        setState(() {
                          CommodityType = value;
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
                          CommodityType = value!;
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
                        "Commodity Description",
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
                      initialValue: CommodityDescription,
                      onChanged: (value) {
                        setState(() {
                          CommodityDescription = value;
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
                          CommodityDescription = value!;
                        });
                      },
                    ),
                  ),
                ],
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15),
                child: Container(
                  width: width,
                  child: Text(
                    textAlign: TextAlign.left,
                    "Quantity",
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
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
                      initialValue: Quantity,
                      onChanged: (value) {
                        setState(() {
                          Quantity = value;
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
                          Quantity = value!;
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
                      "Quality Grade / Standard",
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
                      QualityGrade = value.toString();
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
                        "Issuer Name",
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
                      initialValue: IssuerName,
                      onChanged: (value) {
                        setState(() {
                          IssuerName = value;
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
                          IssuerName = value!;
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
                        "Issuer Type",
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
                      IssuerType = value.toString();
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
                        "Issuer Contact Information",
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
                      initialValue: IssuerContactInfo,
                      onChanged: (value) {
                        setState(() {
                          IssuerContactInfo = value;
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
                          IssuerContactInfo = value!;
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
                        "Warehouse Name",
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
                      initialValue: WarehouseName,
                      onChanged: (value) {
                        setState(() {
                          WarehouseName = value;
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
                          WarehouseName = value!;
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
                        "Warehouse Operator Name",
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
                      initialValue: WarehouseOperatorName,
                      onChanged: (value) {
                        setState(() {
                          WarehouseOperatorName = value;
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
                          WarehouseOperatorName = value!;
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
                        "Warehouse License Number",
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
                      initialValue: WarehouseLicenseNumber,
                      onChanged: (value) {
                        setState(() {
                          WarehouseLicenseNumber = value;
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
                          WarehouseLicenseNumber = value!;
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
                        "Warehouse Location",
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
                      initialValue: WarehouseLocation,
                      onChanged: (value) {
                        setState(() {
                          WarehouseLocation = value;
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
                          WarehouseLocation = value!;
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
                        "WR Number",
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
                      initialValue: WRNumber,
                      onChanged: (value) {
                        setState(() {
                          WRNumber = value;
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
                          WRNumber = value!;
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
                        "WR Issue Date",
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
                "Select date",
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
                        WRIssueDate = value;
                      }),
                    },
                  );
                },
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "WR Expiry Date",
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
                "Select date",
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
                        WRExpiryDate = value;
                      }),
                    },
                  );
                },
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "WR System Registration",
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
                      initialValue: WRSystemRegistration,
                      onChanged: (value) {
                        setState(() {
                          WRSystemRegistration = value;
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
                          WRSystemRegistration = value!;
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
                        "WRS Registration Number",
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
                      initialValue: WRRegistrationNumber,
                      onChanged: (value) {
                        setState(() {
                          WRRegistrationNumber = value;
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
                          WRRegistrationNumber = value!;
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
                        "WR Verifier",
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
                      initialValue: WRVerifier,
                      onChanged: (value) {
                        setState(() {
                          WRVerifier = value;
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
                          WRVerifier = value!;
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
                        "Storage Conditions",
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
                      initialValue: StorageCondition,
                      onChanged: (value) {
                        setState(() {
                          StorageCondition = value;
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
                          StorageCondition = value!;
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
                        "Warehouse Accreditation Body",
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
                      initialValue: WarehouseAccreditationBody,
                      onChanged: (value) {
                        setState(() {
                          WarehouseAccreditationBody = value;
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
                          WarehouseAccreditationBody = value!;
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
                        "Minimum Purchase Amount",
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
                      initialValue: MinimumPurchaseAmount,
                      onChanged: (value) {
                        setState(() {
                          MinimumPurchaseAmount = value;
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
                          MinimumPurchaseAmount = value!;
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
                        "Maturity Date",
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
                "Select date",
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
                        MaturityDate = value;
                      }),
                    },
                  );
                },
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Auto-Rollover",
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
                      AutoRollover = value.toString();
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
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Custody, Ownership & Legal",
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
                        "**************",
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
                        "Current Beneficial Owner/Legal Holder of WR",
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
                      initialValue: CurrentBeneficialOwner,
                      onChanged: (value) {
                        setState(() {
                          CurrentBeneficialOwner = value;
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
                          CurrentBeneficialOwner = value!;
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
                        "WR Custodian Name",
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
                      initialValue: WRCustodianName,
                      onChanged: (value) {
                        setState(() {
                          WRCustodianName = value;
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
                          WRCustodianName = value!;
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
                        "Ownership Rights Represented",
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
                      initialValue: OwnershipRightsRepresented,
                      onChanged: (value) {
                        setState(() {
                          OwnershipRightsRepresented = value;
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
                          OwnershipRightsRepresented = value!;
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
                        "Trustee or Third-party Oversight",
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
                      initialValue: TrusteeOrThirdPartyOversight,
                      onChanged: (value) {
                        setState(() {
                          TrusteeOrThirdPartyOversight = value;
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
                          TrusteeOrThirdPartyOversight = value!;
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
                        "Lien or Encumbrances",
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
                      initialValue: LienOrEncumbrances,
                      onChanged: (value) {
                        setState(() {
                          LienOrEncumbrances = value;
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
                          LienOrEncumbrances = value!;
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
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Commodity Valuation & Tokenization Terms",
                        style: TextStyle(
                          fontSize: 18,
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
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "**************",
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
                        "Asset Valuation",
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
                      initialValue: AssetValuation,
                      onChanged: (value) {
                        setState(() {
                          AssetValuation = value;
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
                          AssetValuation = value!;
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
                        "Valuation Date",
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
                "Select date",
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
                        ValuationDate = value;
                      }),
                    },
                  );
                },
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: SizedBox(
                      width: 300,
                      child: Text(
                        "Valuation Methodology",
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
                      initialValue: ValuationMethodology,
                      onChanged: (value) {
                        setState(() {
                          ValuationMethodology = value;
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
                          ValuationMethodology = value!;
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
                      "Tokenization Objective",
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
                      initialValue: TokenizationObjective,
                      onChanged: (value) {
                        setState(() {
                          TokenizationObjective = value;
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
                          TokenizationObjective = value!;
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
                        "Holding Period / Lock-in",
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
                      initialValue: HoldingPeriod,
                      onChanged: (value) {
                        setState(() {
                          HoldingPeriod = value;
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
                          HoldingPeriod = value!;
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
                        "Redemption Mechanism",
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
                      RedemptionMechanism = value.toString();
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
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Parties Involved",
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
                        "**************",
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
                        "Trustee (if applicable)",
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
                      initialValue: PartiesInvolvedTrustee,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedTrustee = value;
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
                          PartiesInvolvedTrustee = value!;
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
                        "Underwriter (if any)",
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
                      initialValue: PartiesInvolvedUnderwriter,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedUnderwriter = value;
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
                          PartiesInvolvedUnderwriter = value!;
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
                        "Asset Manager (optional)",
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
                      initialValue: PartiesInvolvedAssetManager,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedAssetManager = value;
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
                          PartiesInvolvedAssetManager = value!;
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
                        "Legal Advisor",
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
                      initialValue: PartiesInvolvedLegalAdvisor,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedLegalAdvisor = value;
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
                          PartiesInvolvedLegalAdvisor = value!;
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
                        "Auditor / Verifier",
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
                      initialValue: PartiesInvolvedAuditorVerifier,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedAuditorVerifier = value;
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
                          PartiesInvolvedAuditorVerifier = value!;
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
                        "Regulator / Oversight Body",
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
                      initialValue: PartiesInvolvedRegulator,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedRegulator = value;
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
                          PartiesInvolvedRegulator = value!;
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
                      "Risk Disclosures & Insurance",
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
                        "**************",
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
                        "Market Risk",
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
                      initialValue: RisksMarketRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksMarketRisk = value;
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
                          RisksMarketRisk = value!;
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
                        "Storage Risk",
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
                      initialValue: RisksStorageRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksStorageRisk = value;
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
                          RisksStorageRisk = value!;
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
                        "Title Risk",
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
                      initialValue: RisksTitleRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksTitleRisk = value;
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
                          RisksTitleRisk = value!;
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
                        "Fraud Risk",
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
                      initialValue: RisksFraudRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksFraudRisk = value;
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
                          RisksFraudRisk = value!;
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
                        "Insurance Risk",
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
                      initialValue: RisksInsuranceRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksInsuranceRisk = value;
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
                          RisksInsuranceRisk = value!;
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
                        "Operational Risk",
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
                      initialValue: RisksOperationalRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksOperationalRisk = value;
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
                          RisksOperationalRisk = value!;
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
                        "Regulatory Risk",
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
                      initialValue: RisksRegulatoryRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksRegulatoryRisk = value;
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
                          RisksRegulatoryRisk = value!;
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
                        "Liquidity Risk",
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
                      initialValue: RisksLiquidityRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksLiquidityRisk = value;
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
                          RisksLiquidityRisk = value!;
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
                        "Force Majeure Risk",
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
                      initialValue: RisksForceMajeureRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksForceMajeureRisk = value;
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
                          RisksForceMajeureRisk = value!;
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
                        "Early Redemption Risk",
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
                      initialValue: RisksEarlyRedemptionRisk,
                      onChanged: (value) {
                        setState(() {
                          RisksEarlyRedemptionRisk = value;
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
                          RisksEarlyRedemptionRisk = value!;
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
                        "Mitigation Measures",
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
                      initialValue: RisksMitigationMeasures,
                      onChanged: (value) {
                        setState(() {
                          RisksMitigationMeasures = value;
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
                          RisksMitigationMeasures = value!;
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
                        "Insurance Coverage Summary",
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
                      initialValue: RisksInsuranceCoverageSummary,
                      onChanged: (value) {
                        setState(() {
                          RisksInsuranceCoverageSummary = value;
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
                          RisksInsuranceCoverageSummary = value!;
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
                        "Insurance Provider",
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
                      initialValue: RisksInsuranceProvider,
                      onChanged: (value) {
                        setState(() {
                          RisksInsuranceProvider = value;
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
                          RisksInsuranceProvider = value!;
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
                        "Coverage Value",
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
                      initialValue: RisksCoverageValue,
                      onChanged: (value) {
                        setState(() {
                          RisksCoverageValue = value;
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
                          RisksCoverageValue = value!;
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
      // make initial request to the server using the
      // following credential
      //
      var newData = {...data as Map};
      // newData['ownershipType'] = assetOwnership;

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
