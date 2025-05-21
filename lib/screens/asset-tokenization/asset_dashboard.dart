import 'dart:developer';

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
  int _daysProgress = 0;
  int _totalDays = 0;
  String fiatCurrency = '';
  String regulatorName = '';
  double tokenFee = 0;

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
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = appState.tokenizedAsset!;
    inspect(appState.viewData);
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
    _amountRaised = (tokenizedAsset.quantityOfTokensSold ?? 0) *
        tokenizedAsset.pricePerToken!;
    totalAmountToBeRaised =
        tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!;
    double normalizedProgress =
        _amountRaised / totalAmountToBeRaised; // Convert to 0-1 range
    double normalizedDaysProgress = _daysProgress == _totalDays
        ? 1
        : 1 - (_daysProgress / _totalDays); // Convert to 0-1 range

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(children: [
          CustomAppBar(
            context,
            notifier.getwihitecolor,
            'Asset Details',
            notifier.getbluewhitecolor,
            height: height / 15,
          ).getBar(),
          SizedBox(height: height / 50),
          if (tokenizedAsset.assetLogo != null) ...[
            CircleAvatar(
              radius: 30,
              backgroundColor: tokenizedAsset.assetAlreadyExists == 1
                  ? notifier.getgreencolor
                  : notifier.getbluecolor90,
              child: ClipRRect(
                borderRadius: BorderRadius.circular(100.0),
                child: Image.network(
                  tokenizedAsset.assetLogo!,
                  width: 50,
                  height: 50,
                  fit: BoxFit.fill,
                  errorBuilder: (context, error, stackTrace) {
                    return Image.asset(
                      'assets/images/trovo.png',
                      height: 50,
                      width: 50,
                    );
                  },
                ),
              ),
            ),
          ] else ...[
            Image.asset(
              'assets/images/trovo.png',
              height: 35,
              width: 35,
            ),
          ],
          SizedBox(
            height: height / 70,
          ),
          Text(
            tokenizedAsset.assetCode!.toUpperCase(),
            style: TextStyle(
              fontSize: 16,
              fontFamily: fontsemibold,
              color: notifier.getbluewhitecolor,
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: width / 2.9,
                child: Card(
                  shadowColor: Colors.black,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(15.0),
                  ),
                  color: notifier.isDark
                      ? notifier.getbluecolor90
                      : notifier.getpillbg,
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          tokenizedAsset.assetAlreadyExists == 1
                              ? 'Existing Asset'
                              : 'Upcoming Asset',
                          textAlign: TextAlign.start,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 13,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              Container(
                child: Card(
                  shadowColor: Colors.black,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(15.0),
                  ),
                  color: notifier.isDark
                      ? notifier.getbluecolor90
                      : notifier.getpillbg,
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          getTokenizationStatus(
                              tokenizedAsset.tokenizationStatus!),
                          textAlign: TextAlign.start,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 13,
                            fontFamily: fontbody,
                            color: notifier.getgreencolor,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ],
          ),
          SizedBox(height: height / 70),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Container(
              child: Center(
                child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 30.0),
                          child: Text(
                            'Description',
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
                      padding: const EdgeInsets.all(10.0),
                      child: Text(
                        tokenizedAsset.assetDescription!,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 14,
                          height: 1.4,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                  ],
                ),
              ),
            ),
          ),
          SizedBox(height: height / 50),
          Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  'Statistics',
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ],
          ),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'Total Token Supply',
                  value:
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode}',
                ),
                infoCard(
                  notifier,
                  label: tokenizedAsset.assetAlreadyExists == 1
                      ? 'Asset Value'
                      : 'Total Project Budget',
                  value:
                      '${getFiatValue((tokenizedAsset.numberOfTokenToBeIssued! * tokenizedAsset.pricePerToken!))} ${fiatCurrency}',
                ),
              ],
            ),
          ),
          SizedBox(height: 5),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'Tokens not for Sale',
                  value:
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeIssued! - tokenizedAsset.numberOfTokenToBeSold! - tokenFee)} ${tokenizedAsset.assetCode}',
                ),
                infoCard(
                  notifier,
                  label: tokenizedAsset.assetAlreadyExists == 1
                      ? 'Amount Retained'
                      : 'Equity Contributed',
                  value:
                      '${getFiatValue((tokenizedAsset.assetOwnerRetainedOrContributedValue ?? 0))} ${fiatCurrency}',
                ),
              ],
            ),
          ),
          SizedBox(height: 5),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'Tokens for Sale',
                  value:
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeSold ?? 0)} ${tokenizedAsset.assetCode}',
                ),
                infoCard(
                  notifier,
                  label: 'Amount to be Raised',
                  value:
                      '${getFiatValue(totalAmountToBeRaised)} ${fiatCurrency}',
                ),
              ],
            ),
          ),
          SizedBox(height: 5),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'Funding Currency',
                  value: 'CNGN',
                ),
                infoCard(
                  notifier,
                  label: 'Price Per Token',
                  value:
                      '${truncateToDecimalPlaces(tokenizedAsset.pricePerToken!)} ${fiatCurrency}',
                ),
              ],
            ),
          ),
          SizedBox(height: 5),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'No. of Interest Expressed',
                  value: '${tokenizedAsset.numberOfExpressedInterests} users',
                ),
                infoCard(
                  notifier,
                  label: 'Purchase Commitments',
                  value:
                      '${getFiatValue(tokenizedAsset.expressedInterestAmount ?? 0)} ${tokenizedAsset.assetQuoteCurrency}',
                ),
              ],
            ),
          ),
          SizedBox(height: 5),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'Total Quantity Sold',
                  value:
                      '${truncateToDecimalPlaces(tokenizedAsset.quantityOfTokensSold ?? 0)} ${tokenizedAsset.assetCode}',
                ),
                infoCard(
                  notifier,
                  label: 'Total Amount Raised',
                  value:
                      '${getFiatValue(_amountRaised)} ${tokenizedAsset.assetQuoteCurrency}',
                ),
              ],
            ),
          ),
          SizedBox(height: 5),
          IntrinsicHeight(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  notifier,
                  label: 'No. of Users Bought',
                  value: '${tokenizedAsset.numberOfSubscribers} Users',
                ),
                infoCard(
                  notifier,
                  label: '% of Amount Raised',
                  value:
                      '${getFiatValue((_amountRaised / (totalAmountToBeRaised) * 100))} %',
                ),
              ],
            ),
          ),
          SizedBox(height: height / 50),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12.0, vertical: 10),
            child: Container(
              child: Card(
                shadowColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: notifier.isDark
                    ? notifier.getbluecolor90
                    : notifier.getpillbg,
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 12.0),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      SizedBox(height: 15),
                      Row(
                        children: [
                          Text('Amount Raised',
                              style: TextStyle(
                                fontSize: 14,
                                fontFamily: fontsemibold,
                                color: notifier.getbluewhitecolor,
                              )),
                        ],
                      ),
                      SizedBox(height: 5),
                      // Gesture Detector for Tapping the Progress Bar
                      Container(
                        height: 10,
                        child: LinearProgressIndicator(
                          value: normalizedProgress, // Show progress (0 to 1)
                          minHeight: 10,
                          borderRadius: BorderRadius.circular(10),
                          backgroundColor: Colors.grey[300],
                          color: notifier.getgreencolor,
                        ),
                      ),
                      SizedBox(height: 5),
                      Row(
                        children: [
                          Text(
                              '${formatHistoryNumber(_amountRaised, 6)} ${tokenizedAsset.assetQuoteCurrency} ',
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
                              '${formatHistoryNumber(totalAmountToBeRaised, 6)} ${tokenizedAsset.assetQuoteCurrency}',
                              style: TextStyle(
                                fontSize: 12,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              )),
                        ],
                      ),
                      SizedBox(height: height / 50),
                      if (tokenizedAsset.tokenizationStatus == 5) ...[
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text('Days Remaining',
                                    style: TextStyle(
                                      fontSize: 14,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    )),
                              ],
                            ),
                            SizedBox(height: 5),
                            if (_daysProgress >= 0) ...[
                              // Gesture Detector for Tapping the Progress Bar
                              Container(
                                height: 10,
                                child: LinearProgressIndicator(
                                  value:
                                      normalizedDaysProgress, // Show progress (0 to 1)
                                  minHeight: 10,
                                  borderRadius: BorderRadius.circular(10),
                                  backgroundColor: Colors.grey[300],
                                  color: Colors.blue[400],
                                ),
                              ),
                              SizedBox(height: 5),
                              Row(
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
                            ],
                          ],
                        ),
                        SizedBox(height: height / 50),
                      ],
                    ],
                  ),
                ),
              ),
            ),
          ),
          Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Text(
                  'Details',
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
                    : notifier.getpillbg,
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 10.0),
                  child: Column(
                    children: [
                      categoryTile(
                        notifier,
                        label: 'Asset Information',
                        imageUrl: 'assets/images/asset-info.png',
                        onTap: () {
                          var details = {
                            'Status': tokenizedAsset.assetAlreadyExists == 1
                                ? 'Existing'
                                : 'Upcoming',
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
                                tokenizedAsset.assetAlreadyExists == 0
                                    ? tokenizedAsset.projectExpectedNoOfJobs
                                    : "",
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
                      if (tokenizedAsset.assetAlreadyExists == 0) ...[
                        categoryTile(
                          notifier,
                          label: 'Asset Financial Information',
                          imageUrl: 'assets/images/stakeholders.png',
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
                          imageUrl: 'assets/images/proof.png',
                          onTap: () {
                            var details = {
                              "Legal Risks Identified":
                                  tokenizedAsset.projectIdentifiedLegalRisks ??
                                      '',
                              "Regulatory Risks Identified": tokenizedAsset
                                      .projectIdentifiedRegulatoryRisks ??
                                  '',
                              "Operational/Execution Risks Identified":
                                  tokenizedAsset
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
                      ],
                      categoryTile(
                        notifier,
                        label: 'Stakeholders Information',
                        imageUrl: 'assets/images/stakeholders.png',
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
                            'Legal/Professional Advisor': tokenizedAsset
                                    .legalAdvisor
                                    ?.toLowerCase()
                                    .capitalizeEachWord() ??
                                '',
                            'Rating Agency': '',
                          };
                          displayDetails("Stakeholders Information", details);
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Legal & Compliance Information',
                        imageUrl: 'assets/images/shareholders.png',
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
                        imageUrl: 'assets/images/documents.png',
                        onTap: () {
                          displayDocuments('Verification Documents');
                        },
                      ),
                      SizedBox(height: 10),
                      categoryTile(
                        notifier,
                        label: 'Proof of Payment Documents',
                        imageUrl: 'assets/images/proof.png',
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
      backgroundColor: notifier.getwihitecolor,
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
                    if (item.value.toString().isNotEmpty) ...[
                      infoTile(
                        notifier,
                        item.key,
                        item.value,
                      ),
                    ],
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
      backgroundColor: notifier.getwihitecolor,
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

  String getTokenizationStatus(int status) {
    switch (status) {
      case 0:
        return 'Continue';
      case 1:
        return 'Awaiting Fee';
      case 4:
        return 'Approved';
      case 5:
        return 'Primary Sales';
      case 6:
        return 'Secondary Market';
      case 7:
        return 'Liquidated';
      case 8:
        return 'Refunded';
      default:
        return 'Processing';
    }
  }
}
