import 'dart:async';

import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetDashboard extends StatefulWidget {
  const AssetDashboard({Key? key}) : super(key: key);

  @override
  State<AssetDashboard> createState() => _AssetDashboardState();
}

class _AssetDashboardState extends State<AssetDashboard>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  String assetType = '';
  late TokenizedAsset tokenizedAsset;
  double _amountRaised = 0;
  double totalAmountToBeRaised = 0;
  bool _showValue = false;
  int _daysProgress = 0;
  int _totalDays = 0;
  bool _showDaysValue = false;
  String fiatCurrency = '';
  String regulatorName = '';

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  void _toggleValueDisplay() {
    setState(() {
      _showValue = !_showValue;
    });

    // Hide value after 2 seconds
    Future.delayed(Duration(seconds: 2), () {
      setState(() {
        _showValue = false;
      });
    });
  }

  void _toggleDaysValueDisplay() {
    setState(() {
      _showDaysValue = !_showDaysValue;
    });

    // Hide value after 2 seconds
    Future.delayed(Duration(seconds: 2), () {
      setState(() {
        _showDaysValue = false;
      });
    });
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = appState.tokenizedAsset!;
    var assetTypes = appState.tokenizationData['assetTypes'];
    for (var i = 0; i < assetTypes.length; i++) {
      if (assetTypes[i]['id'].toString() == tokenizedAsset.assetType) {
        assetType = assetTypes[i]['assetType'];
      }
    }

    _daysProgress =
        DateTime.now().difference(tokenizedAsset.salesStart!).inDays;
    _totalDays =
        tokenizedAsset.salesEnd!.difference(tokenizedAsset.salesStart!).inDays;

    var quoteCurrencyCode = '';

    for (var i = 0;
        i < appState.tokenizationData['countryConfigs'].length;
        i++) {
      if (appState.tokenizationData['countryConfigs'][i]['countryCode']
              .toString()
              .toLowerCase() ==
          tokenizedAsset.assetCountryLocation.toString().toLowerCase()) {
        quoteCurrencyCode =
            appState.tokenizationData['countryConfigs'][i]['quoteCurrencyCode'];
        regulatorName =
            appState.tokenizationData['countryConfigs'][i]['regulatorName'];
      }
    }

    for (var i = 0;
        i < appState.tokenizationData['tokenizationCurrencies'].length;
        i++) {
      if (appState.tokenizationData['tokenizationCurrencies'][i]['assetCode']
              .toString()
              .toLowerCase() ==
          quoteCurrencyCode.toString().toLowerCase()) {
        fiatCurrency = appState.tokenizationData['tokenizationCurrencies'][i]
                ['label']
            .toString()
            .toUpperCase();
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    _amountRaised = tokenizedAsset.subscriptionAmount ?? 0;
    totalAmountToBeRaised =
        tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!;
    double normalizedProgress =
        _amountRaised / totalAmountToBeRaised; // Convert to 0-1 range
    double normalizedDaysProgress = _daysProgress == _totalDays
        ? 1
        : _daysProgress / _totalDays; // Convert to 0-1 range

    print('fasdfsd=>>>>>$_daysProgress>>>>>>>>> $normalizedDaysProgress');

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(children: [
          CustomAppBar(
            context,
            notifier.getwihitecolor,
            tokenizedAsset.assetName ?? 'Asset',
            notifier.getbluewhitecolor,
            height: height / 15,
          ).getBar(),
          SizedBox(height: height / 50),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Text(
              // 'Atlantis Estate 1 tokens are fractional tokens that represent part ownership (via investment) of our real estate development project at Atlantis Estate, Lekki, Lagos, Nigeria. ',
              tokenizedAsset.assetDescription!,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 15,
                fontFamily: fontbody,
                color: notifier.getbluewhitecolor,
              ),
            ),
          ),
          SizedBox(height: height / 50),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              infoCard(
                notifier,
                label: 'Asset Value',
                value:
                    '${getFiatValue((tokenizedAsset.numberOfTokenToBeIssued! * tokenizedAsset.pricePerToken!))} ${fiatCurrency}',
                extraValue: '\$4,390.23',
              ),
              SizedBox(
                width: width / 50,
              ),
              infoCard(
                notifier,
                label: 'Total Supply',
                value:
                    '${getFiatValue(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode}',
                extraValue: '',
              ),
            ],
          ),
          SizedBox(height: height / 70),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              infoCard(
                notifier,
                label: 'Amount to be Raised',
                value: '${getFiatValue(totalAmountToBeRaised)} ${fiatCurrency}',
                extraValue: '',
              ),
              SizedBox(
                width: width / 50,
              ),
              infoCard(
                notifier,
                label: 'Tokens for Sale',
                value:
                    '${getFiatValue(tokenizedAsset.numberOfTokenToBeSold ?? 0)} ${tokenizedAsset.assetCode}',
                extraValue: '',
              ),
            ],
          ),
          SizedBox(height: height / 70),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              infoCard(
                notifier,
                label: 'Funding Currency',
                value: 'CNGN',
                extraValue: '',
              ),
              SizedBox(
                width: width / 50,
              ),
              infoCard(
                notifier,
                label: 'Price Per Token',
                value:
                    '${getFiatValue(tokenizedAsset.pricePerToken!)} ${fiatCurrency}',
                extraValue: '',
              ),
            ],
          ),
          SizedBox(height: height / 70),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              infoCard(
                notifier,
                label: 'Interest Expressed Amount',
                value:
                    '${getFiatValue(tokenizedAsset.expressedInterestAmount ?? 0)} ${tokenizedAsset.assetQuoteCurrency}',
                extraValue: '',
              ),
              SizedBox(
                width: width / 50,
              ),
              infoCard(
                notifier,
                label: 'Number of Interest Expressed',
                value: '???? users',
                extraValue: '',
              ),
            ],
          ),
          SizedBox(height: height / 70),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              infoCard(
                notifier,
                label: 'Total Amount Raised',
                value:
                    '${getFiatValue(_amountRaised)} ${tokenizedAsset.assetQuoteCurrency}',
                extraValue: '',
              ),
              SizedBox(
                width: width / 50,
              ),
              infoCard(
                notifier,
                label: 'Total Quantity Sold',
                value: '${'20'} ${tokenizedAsset.assetCode}',
                extraValue: '',
              ),
            ],
          ),
          SizedBox(height: height / 70),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              infoCard(
                notifier,
                label: 'No. of Users Bought',
                value: '10 Users',
                extraValue: '',
              ),
              SizedBox(
                width: width / 50,
              ),
              infoCard(
                notifier,
                label: '% of Amount Raised',
                value: '${(_amountRaised / (totalAmountToBeRaised) * 100)} %',
                extraValue: '',
              ),
            ],
          ),
          SizedBox(height: height / 50),
          Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                        '${getFiatValue(_amountRaised)} ${tokenizedAsset.assetQuoteCurrency} ',
                        style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor)),
                    Text('raised out of ',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        )),
                    Text(
                        '${getFiatValue(totalAmountToBeRaised)} ${tokenizedAsset.assetQuoteCurrency}',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        )),
                  ],
                ),
                SizedBox(height: 5),
                // Gesture Detector for Tapping the Progress Bar
                GestureDetector(
                  onTap: _toggleValueDisplay,
                  onLongPress: _toggleValueDisplay,
                  child: Stack(
                    alignment: Alignment.center,
                    children: [
                      Container(
                        width: 300,
                        height: 20,
                        child: LinearProgressIndicator(
                          value: normalizedProgress, // Show progress (0 to 1)
                          minHeight: 20,
                          borderRadius: BorderRadius.circular(10),
                          backgroundColor: Colors.grey[300],
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      // Show progress value when tapped/long pressed
                      if (_showValue)
                        Container(
                          width: 300,
                          height: 20,
                          decoration: BoxDecoration(
                              color: Colors.black54.withOpacity(
                                  0.7), // Semi-transparent background
                              borderRadius: BorderRadius.circular(10)),
                          alignment: Alignment.center,
                          child: Text(
                            "${getFiatValue(_amountRaised)} / ${getFiatValue(totalAmountToBeRaised)}",
                            style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold),
                          ),
                        ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          SizedBox(height: height / 50),
          if (tokenizedAsset.tokenizationStatus == 5) ...[
            Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (_daysProgress >= 0) ...[
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text((_totalDays - _daysProgress).toString(),
                            style: TextStyle(
                                fontSize: 12,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor)),
                        Text(' days remaining out of ',
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            )),
                        Text('${_totalDays} days',
                            style: TextStyle(
                                fontSize: 12,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor)),
                      ],
                    ),
                    SizedBox(height: 5),
                    // Gesture Detector for Tapping the Progress Bar
                    GestureDetector(
                      onTap: _toggleDaysValueDisplay,
                      onLongPress: _toggleDaysValueDisplay,
                      child: Stack(
                        alignment: Alignment.center,
                        children: [
                          Container(
                            width: 300,
                            height: 20,
                            child: LinearProgressIndicator(
                              value:
                                  normalizedDaysProgress, // Show progress (0 to 1)
                              minHeight: 20,
                              borderRadius: BorderRadius.circular(10),
                              backgroundColor: Colors.grey[300],
                              color: notifier.getbluewhitecolor,
                            ),
                          ),

                          // Show progress value when tapped/long pressed
                          if (_showDaysValue)
                            Container(
                              width: 300,
                              height: 20,
                              decoration: BoxDecoration(
                                  color: Colors.black54.withOpacity(
                                      0.7), // Semi-transparent background
                                  borderRadius: BorderRadius.circular(10)),
                              alignment: Alignment.center,
                              child: Text(
                                "${_daysProgress} / ${_totalDays}",
                                style: TextStyle(
                                    color: Colors.white,
                                    fontWeight: FontWeight.bold),
                              ),
                            ),
                        ],
                      ),
                    ),
                  ],
                ],
              ),
            ),
            SizedBox(height: height / 50),
          ],
          Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  'Details of Asset',
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ],
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10.0),
            child: Container(
              child: Card(
                shadowColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: notifier.isDark
                    ? notifier.getbluecolor90
                    : notifier.getaddsubwalletgrey,
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 10.0),
                  child: Column(
                    children: [
                      categoryTile(
                        notifier,
                        label: 'Asset Information',
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          var details = {
                            'Status': tokenizedAsset.assetAlreadyExists == 1
                                ? 'Existing'
                                : 'Not Existing',
                            'Sector': tokenizedAsset.assetSector ?? '',
                            'Sub-Sector': tokenizedAsset.assetSubSector ?? '',
                            'Type': assetType,
                            'Asset Country': iso2Countries[
                                    tokenizedAsset.assetCountryLocation] ??
                                "",
                            'Address':
                                tokenizedAsset.assetPhysicalAddress ?? '',
                            'Project Strategic Objectives':
                                tokenizedAsset.projectStrategicObjectives ?? "",
                            "Project Development Timeline":
                                tokenizedAsset.projectDevelopmentTimeline ?? "",
                            "Key Milestones & Dates":
                                tokenizedAsset.projectKeyMilestoneAndDates ??
                                    "",
                            "Project Scope": tokenizedAsset.projectScope ?? "",
                            "Project Economic Benefits":
                                tokenizedAsset.projectEconomicBenefits ?? "",
                            "Expected No. of Job to be Created":
                                tokenizedAsset.projectExpectedNoOfJobs ?? "",
                            "Project Intended Social Benefits":
                                tokenizedAsset.projectIntendedSocialBenefits ??
                                    "",
                            "Technical Partners":
                                tokenizedAsset.projectTechnicalPartners ?? "",
                            "Financial Partners":
                                tokenizedAsset.projectFinancialPartners ?? "",
                          };

                          details[tokenizedAsset.assetAlreadyExists == 1
                                  ? 'Original Asset Owner'
                                  : 'Project Sponsor'] =
                              tokenizedAsset.assetOwnerName ?? '';

                          displayDetails('Asset Information', details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Asset Token & Sale Information',
                        imageUrl: 'assets/images/token-info.png',
                        onTap: () {
                          var details = {
                            'Sales Window':
                                '${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesStart!)} - ${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesEnd!)}',
                            'Cap Amount':
                                '${getFiatValue(double.parse(tokenizedAsset.capAmountInFiat!.toString()))} ${fiatCurrency}',
                            'Cap Quantity':
                                '${getFiatValue(tokenizedAsset.capQuantity!)} ${tokenizedAsset.assetCode!.toUpperCase()}',
                            'Cap Duration':
                                '${tokenizedAsset.capDurationInDays} days',
                            'Proceed Payout Cycle': tokenizedAsset.proceedCycle
                                    ?.toLowerCase()
                                    .capitalizeEachWord() ??
                                '',
                            'Exempted Countries':
                                '${tokenizedAsset.exemptedCountries!.replaceAll(',', ', ')}',
                          };
                          displayDetails(
                              'Asset Token & Sale Information', details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Asset Financial Information',
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          var details = {
                            "Estimated Project IRR": formatNumber(
                                tokenizedAsset.estimatedProjectIRR ?? 0),
                            "Estimated Project ROI": formatNumber(
                                tokenizedAsset.estimatedProjectROI ?? 0),
                            "Estimated Project NPV at Launch (Day 1)":
                                formatNumber(
                                    tokenizedAsset.estimatedProjectNPV ?? 0),
                            "Estimated Project Payback Periods (in Months)":
                                tokenizedAsset
                                        .estimatedProjectPaybackPeriodsInMonths ??
                                    '',
                            "All Key Assumptions Including Values Assumed":
                                tokenizedAsset.keyAssumptionsList ?? '',
                          };
                          displayDetails(
                              'Asset Financial Information', details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: "Project Risk Assessment",
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          var details = {
                            "Legal Risks Identified":
                                tokenizedAsset.projectIdentifiedLegalRisks ??
                                    '',
                            "Regulatory Risks Identified": tokenizedAsset
                                    .projectIdentifiedRegulatoryRisks ??
                                '',
                            "Operational/Execution Risks Identified": tokenizedAsset
                                    .projectIdentifiedOperationalOrExecutionRisks ??
                                '',
                            "Market Risks Identified":
                                tokenizedAsset.projectIdentifiedMarketRisks ??
                                    '',
                            "Other Relevant Risks Identified": tokenizedAsset
                                    .projectIdentifiedOtherRelevantRisks ??
                                '',
                          };
                          displayDetails("Project Risk Assessment", details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Stakeholders Information',
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          var details = {
                            'Regulator': regulatorName
                                .toLowerCase()
                                .capitalizeEachWord(),
                            'Asset Custodian': tokenizedAsset
                                    .approvedAssetCustodianInfo
                                    ?.assetCustodianName ??
                                '',
                            'Asset Manager': tokenizedAsset
                                    .assetManagerInfo?.assetManagerName ??
                                '',
                            'Issuing House': tokenizedAsset
                                    .assetIssuingHouseInfo
                                    ?.assetIssuingHouseName
                                    ?.toLowerCase()
                                    .capitalizeEachWord() ??
                                '',
                            'Rating Agency': '',
                          };
                          displayDetails("Project Risk Assessment", details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Legal & Compliance Information',
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          var details = {
                            'I confirm that this asset is free of liens, mortgages, and outstanding loans.':
                                '${tokenizedAsset.isFreeFromLiensAndEncumbrances == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset is not pledged as collateral for any debts and has no use restrictions.':
                                '${tokenizedAsset.undertakingNotCollateral == 1 ? 'Yes' : 'No'}',
                            'I confirm that no third party has any claims, rights, or interests in this asset.':
                                '${tokenizedAsset.undertakingNoClaims == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset is free of any foreclosure, bankruptcy proceedings, legal disputes, judgments, or court-ordered payments.':
                                '${tokenizedAsset.undertakingNoForeclosure == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset complies with all environmental and land use regulations and is free of violations.':
                                '${tokenizedAsset.complianceNoViolation == 1 ? 'Yes' : 'No'}',
                            'I confirm that all necessary permits, licenses, and approvals for the use and ownership of this asset are in place.':
                                '${tokenizedAsset.complianceNoViolation == 1 ? 'Yes' : 'No'}',
                            'I confirm that there are no unpaid taxes, utility bills, fees, or other property-related expenses associated with this asset.':
                                '${tokenizedAsset.outstandingFinancialRespNoDebts == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset does not have any hidden liabilities or obligations that have not been disclosed. ':
                                '${tokenizedAsset.outstandingFinancialRespNoHiddenLiabilities == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset is adequately insured against risks such as fire, theft, and natural disasters.':
                                '${tokenizedAsset.riskManagementFullyInsured == 1 ? 'Yes' : 'No'}',
                            'I confirm that the declared value of this asset reflects its current market value and condition.':
                                '${tokenizedAsset.riskManagementDeclaredValue == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset is not affected by undisclosed easements, rights of way, expropriation, or condemnation.':
                                '${tokenizedAsset.physicalConditionNoUndisclosedEasements == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset is structurally sound and has no unresolved maintenance or safety issues.':
                                '${tokenizedAsset.physicalConditionSound == 1 ? 'Yes' : 'No'}',
                            'I confirm that this asset is not subject to any agreements, such as leases or contracts, that could limit its use or transfer.':
                                '${tokenizedAsset.physicalConditionNolease == 1 ? 'Yes' : 'No'}',
                          };
                          displayDetails(
                              "Legal & Compliance Information", details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Verification Documents',
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          displayDocuments('Verification Documents');
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Proof of Payment Documents',
                        imageUrl: tokenizedAsset.assetCode!.toUpperCase(),
                        onTap: () {
                          displayDocuments('Proof of Payment Documents');
                        },
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
          SizedBox(
            height: height / 30,
          ),
          // Button(
          //   'Payout Proceeds',
          //   notifier.getbluecolor,
          //   wihitecolor,
          //   onTap: () {
          //     appState.currentAction = PageAction(
          //       state: PageState.addPage,
          //       page: ProceedsPayOutViewPageConfig,
          //     );
          //   },
          // ),
          SizedBox(
            height: height / 10,
          ),
        ]),
      ),
    );
  }

  void displayDetails(String label, Map<String, dynamic> items) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)), // Rounded top corners
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
                          label,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(
                            Icons.cancel_outlined,
                          ),
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.zero,
                            minimumSize: Size.zero,
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            visualDensity: VisualDensity.compact,
                            alignment: Alignment.centerLeft,
                            foregroundColor: notifier.getbluewhitecolor,
                            backgroundColor: Colors.transparent,
                            shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.zero),
                          ),
                        )
                      ],
                    ),
                  ),
                  SizedBox(height: 20),
                  for (var item in items.entries) ...[
                    infoTile(
                      notifier,
                      item.key,
                      item.value,
                    ),
                  ],
                  SizedBox(height: 60),
                ],
              ),
            );
          },
        );
      },
    );
  }

  void displayDocuments(String label) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.75,
          minChildSize: 0.25,
          expand: false,
          builder: (context, scrollController) {
            return SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14.0),
                child: Column(
                  crossAxisAlignment:
                      tokenizedAsset.assetTokenizationDocuments!.isNotEmpty
                          ? CrossAxisAlignment.start
                          : CrossAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    SizedBox(height: 10),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          label,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(
                            Icons.cancel_outlined,
                          ),
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.zero,
                            minimumSize: Size.zero,
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            visualDensity: VisualDensity.compact,
                            alignment: Alignment.centerLeft,
                            foregroundColor: notifier.getbluewhitecolor,
                            backgroundColor: Colors.transparent,
                            shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.zero),
                          ),
                        )
                      ],
                    ),
                    SizedBox(height: 20),
                    if (tokenizedAsset
                        .assetTokenizationDocuments!.isNotEmpty) ...[
                      for (var item
                          in tokenizedAsset.assetTokenizationDocuments!) ...[
                        TextButton(
                          style: TextButton.styleFrom(
                              padding: EdgeInsets.zero,
                              minimumSize: Size(50, 30),
                              tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                              alignment: Alignment.centerLeft),
                          onPressed: () {
                            var fileUrl = item.documentUrl;
                            if (fileUrl!.isNotEmpty &&
                                fileUrl.endsWith('.pdf')) {
                              appState.pdfUrl = fileUrl;
                              appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: PdfViewPageConfig);

                              return;
                            }

                            appState.goToWebView(fileUrl);
                          },
                          child: Text(
                            item.documentTitle ?? '',
                            style: TextStyle(
                              decoration: TextDecoration.underline,
                              fontSize: 12,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                      ],
                    ] else ...[
                      Text(
                        'No documents here...',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                    SizedBox(height: 60),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }

  void displayProofOfPaymentDocuments(String label) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.75,
          minChildSize: 0.25,
          expand: false,
          builder: (context, scrollController) {
            return SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 14.0),
                child: Column(
                  crossAxisAlignment:
                      tokenizedAsset.assetTokenizationDocuments!.isNotEmpty
                          ? CrossAxisAlignment.start
                          : CrossAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    SizedBox(height: 10),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          label,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(
                            Icons.cancel_outlined,
                          ),
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.zero,
                            minimumSize: Size.zero,
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            visualDensity: VisualDensity.compact,
                            alignment: Alignment.centerLeft,
                            foregroundColor: notifier.getbluewhitecolor,
                            backgroundColor: Colors.transparent,
                            shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.zero),
                          ),
                        )
                      ],
                    ),
                    SizedBox(height: 20),
                    if (tokenizedAsset
                        .assetTokenizationDocuments!.isNotEmpty) ...[
                      for (var item
                          in tokenizedAsset.proofOfPaymentDocuments!) ...[
                        TextButton(
                          style: TextButton.styleFrom(
                              padding: EdgeInsets.zero,
                              minimumSize: Size(50, 30),
                              tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                              alignment: Alignment.centerLeft),
                          onPressed: () {
                            var fileUrl = item.documentUrl;
                            if (fileUrl!.isNotEmpty &&
                                fileUrl.endsWith('.pdf')) {
                              appState.pdfUrl = fileUrl;
                              appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: PdfViewPageConfig);

                              return;
                            }

                            appState.goToWebView(fileUrl);
                          },
                          child: Text(
                            truncateString(item.documentUrl) ?? '',
                            style: TextStyle(
                              decoration: TextDecoration.underline,
                              fontSize: 12,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                      ],
                    ] else ...[
                      Text(
                        'No documents here...',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                    SizedBox(height: 60),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }
}
