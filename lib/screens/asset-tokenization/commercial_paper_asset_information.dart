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

class CommercialPapersAssetInformationView extends StatefulWidget {
  const CommercialPapersAssetInformationView({Key? key}) : super(key: key);

  @override
  State<CommercialPapersAssetInformationView> createState() =>
      _CommercialPapersAssetInformationView();
}

class _CommercialPapersAssetInformationView
    extends State<CommercialPapersAssetInformationView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  bool formHasError = false;
  late dynamic data = {};

  String TitleOfIssuance = '';
  String TypeOfCommercialPaper = '';
  double TotalIssueSize = 0;
  DateTime? IssueDate;
  DateTime? MaturityDate;
  int Tenor = 0;
  double PricingYield = 0;
  String RedemptionStructure = '';
  String UseOfProceeds = '';
  String Ranking = '';
  String BackingSecurity = '';
  String TaxTreatment = '';
  String IssuerName = '';
  String IssuerType = '';
  String IssuerRegistrationNumber = '';
  DateTime? IncorporationDate;
  String RcNumber = '';
  String TaxIdNumber = '';
  String OfficeAddress = '';
  String IssuerContactPerson = '';
  String ContactEmail = '';
  String ContactPhoneNumber = '';
  String BriefCompanyOverview = '';
  // FileUpload SummaryOfFinancials = '';
  String Rating = '';
  String RatingAgency = '';
  String PartiesInvolvedIssuer = '';
  String PartiesInvolvedArranger = '';
  String PartiesInvolvedLegalAdviser = '';
  String PartiesInvolvedAuditor = '';
  String PartiesInvolvedRatingAgency = '';
  String PartiesInvolvedCustodian = '';
  String PartiesInvolvedTrustee = '';
  String PartiesInvolvedAuditorVerifier = '';
  String SecurityRiskLegalBacking = '';
  String SecurityRiskCollateral = '';
  String SecurityRiskDefaultHistory = '';
  String SecurityRiskCreditRating = '';
  String SecurityRiskRiskFactorsSummary = '';
  String SecurityRiskBusinessRisk = '';
  String SecurityRiskDefaultRisk = '';
  String SecurityRiskLiquidityRisk = '';
  String SecurityRiskRegulatoryRisk = '';
  String SecurityRiskMarketRisk = '';
  String SecurityRiskOperationalRisk = '';
  String SecurityRiskMitigationMeasures = '';

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
                "Commercial Paper Asset Information",
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
                        "Commercial Paper & Profile and Issuer Information",
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
                        "Title of Issuance",
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
                      initialValue: TitleOfIssuance,
                      onChanged: (value) {
                        setState(() {
                          TitleOfIssuance = value;
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
                          TitleOfIssuance = value!;
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
                        "Type of Commercial Paper",
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
                      TypeOfCommercialPaper = value.toString();
                    });
                  },
                  [],
                  null,
                  'Select asset ownership',
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
                        "Total Issue Size",
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
                      initialValue: TotalIssueSize,
                      onChanged: (value) {
                        setState(() {
                          TotalIssueSize = value;
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
                          TotalIssueSize = value!;
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
                      "Issue Date",
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
                        IssueDate = value;
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
                        "Tenor",
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
                      initialValue: Tenor,
                      onChanged: (value) {
                        setState(() {
                          Tenor = value;
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
                          Tenor = value!;
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
                        "Pricing Yield",
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
                      initialValue: PricingYield,
                      onChanged: (value) {
                        setState(() {
                          PricingYield = value;
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
                          PricingYield = value!;
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
                        "Redemption Structure",
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
                      RedemptionStructure = value.toString();
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
                        "Use of Proceeds",
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
                      initialValue: UseOfProceeds,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          UseOfProceeds = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Ranking",
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
                      Ranking = value.toString();
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
                        "Backing / Security",
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
                      BackingSecurity = value.toString();
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
                        "Tax Treatment",
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
                      initialValue: TaxTreatment,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          TaxTreatment = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: multilineInput(
                      'Enter details',
                      notifier.getbluecolor,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      100.sp,
                      width / 1.12,
                      initialValue: IssuerType,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          IssuerType = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Issuer Registration No.",
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
                      initialValue: IssuerRegistrationNumber,
                      onChanged: (value) {
                        setState(() {
                          IssuerRegistrationNumber = value;
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
                          IssuerRegistrationNumber = value!;
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
                        "Incorporation Date",
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
                        IncorporationDate = value;
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
                        "RC Number",
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
                      initialValue: RcNumber,
                      onChanged: (value) {
                        setState(() {
                          RcNumber = value;
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
                          RcNumber = value!;
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
                        "Tax ID Number",
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
                      initialValue: TaxIdNumber,
                      onChanged: (value) {
                        setState(() {
                          TaxIdNumber = value;
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
                          TaxIdNumber = value!;
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
                        "Office Address",
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
                      initialValue: OfficeAddress,
                      onChanged: (value) {
                        setState(() {
                          OfficeAddress = value;
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
                          OfficeAddress = value!;
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
                        "Issuer Contact Person",
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
                      initialValue: IssuerContactPerson,
                      onChanged: (value) {
                        setState(() {
                          IssuerContactPerson = value;
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
                          IssuerContactPerson = value!;
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
                        "Contact Email",
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
                      initialValue: ContactEmail,
                      onChanged: (value) {
                        setState(() {
                          ContactEmail = value;
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
                          ContactEmail = value!;
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
                        "Contact Phone Number",
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
                      initialValue: ContactPhoneNumber,
                      onChanged: (value) {
                        setState(() {
                          ContactPhoneNumber = value;
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
                          ContactPhoneNumber = value!;
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
                        "Brief Company Overview",
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
                      initialValue: BriefCompanyOverview,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          BriefCompanyOverview = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Summary of Financials",
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
                      // SummaryOfFinancials = value.toString();
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
                        "Rating (if any)",
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
                      initialValue: Rating,
                      onChanged: (value) {
                        setState(() {
                          Rating = value;
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
                          Rating = value!;
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
                        "Rating Agency (if any)",
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
                        "Issuer",
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
                      initialValue: PartiesInvolvedIssuer,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedIssuer = value;
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
                          PartiesInvolvedIssuer = value!;
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
                        "Arranger / Dealer",
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
                      initialValue: PartiesInvolvedArranger,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedArranger = value;
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
                          PartiesInvolvedArranger = value!;
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
                        "Legal Adviser",
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
                      initialValue: PartiesInvolvedLegalAdviser,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedLegalAdviser = value;
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
                          PartiesInvolvedLegalAdviser = value!;
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
                        "Auditor / Accountant",
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
                      initialValue: PartiesInvolvedAuditor,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedAuditor = value;
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
                          PartiesInvolvedAuditor = value!;
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
                      initialValue: PartiesInvolvedRatingAgency,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedRatingAgency = value;
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
                          PartiesInvolvedRatingAgency = value!;
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
                      initialValue: PartiesInvolvedCustodian,
                      onChanged: (value) {
                        setState(() {
                          PartiesInvolvedCustodian = value;
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
                          PartiesInvolvedCustodian = value!;
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
              SizedBox(height: height / 30),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "Security & Risk Profile",
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
                        "Credit Rating",
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
                      initialValue: SecurityRiskCreditRating,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          SecurityRiskCreditRating = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Legal Backing",
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
                      initialValue: SecurityRiskLegalBacking,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskLegalBacking = value;
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
                          SecurityRiskLegalBacking = value!;
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
                        "Collateral",
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
                      initialValue: SecurityRiskCollateral,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          SecurityRiskCollateral = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Default History",
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
                      initialValue: SecurityRiskDefaultHistory,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          SecurityRiskDefaultHistory = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Risk Factors Summary",
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
                      initialValue: SecurityRiskRiskFactorsSummary,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          SecurityRiskRiskFactorsSummary = value!;
                        });
                      },
                      minLines: 3,
                      maxLines: null,
                      keyboardtype: TextInputType.multiline,
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
                        "Business Risk",
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
                      initialValue: SecurityRiskBusinessRisk,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskBusinessRisk = value;
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
                          SecurityRiskBusinessRisk = value!;
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
                        "Default Risk",
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
                      initialValue: SecurityRiskDefaultRisk,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskDefaultRisk = value;
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
                          SecurityRiskDefaultRisk = value!;
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
                      initialValue: SecurityRiskLiquidityRisk,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskLiquidityRisk = value;
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
                          SecurityRiskLiquidityRisk = value!;
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
                      initialValue: SecurityRiskRegulatoryRisk,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskRegulatoryRisk = value;
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
                          SecurityRiskRegulatoryRisk = value!;
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
                      initialValue: SecurityRiskMarketRisk,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskMarketRisk = value;
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
                          SecurityRiskMarketRisk = value!;
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
                      initialValue: SecurityRiskOperationalRisk,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskOperationalRisk = value;
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
                          SecurityRiskOperationalRisk = value!;
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
                      initialValue: SecurityRiskMitigationMeasures,
                      onChanged: (value) {
                        setState(() {
                          SecurityRiskMitigationMeasures = value;
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
                          SecurityRiskMitigationMeasures = value!;
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
