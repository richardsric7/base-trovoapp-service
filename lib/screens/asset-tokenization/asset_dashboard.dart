import 'dart:async';

import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
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
  double _progress = 15000000;
  double _maxValue = 0;
  bool _showValue = false;
  int _daysProgress = 0;
  int _totalDays = 0;
  bool _showDaysValue = false;

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

    _daysProgress = tokenizedAsset.salesEnd!.difference(DateTime.now()).inDays;
    _totalDays =
        tokenizedAsset.salesEnd!.difference(tokenizedAsset.salesStart!).inDays;
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    _maxValue =
        tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!;
    double normalizedProgress = _progress / _maxValue; // Convert to 0-1 range
    double normalizedDaysProgress =
        _daysProgress / _totalDays; // Convert to 0-1 range

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
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
            if (tokenizedAsset.tokenizationStatus! >= 4) ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  infoCard(
                    notifier,
                    label: 'Asset Value',
                    value:
                        '${getFiatValue((tokenizedAsset.numberOfTokenToBeIssued! * tokenizedAsset.pricePerToken!))} ${appState.defaultCurrency}',
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
                    value:
                        '${getFiatValue(tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!)} ${appState.defaultCurrency}',
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
                    extraValue: '\$0.12',
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
                        '${getFiatValue(tokenizedAsset.pricePerToken!)} ${appState.defaultCurrency}',
                    extraValue: '\$0.12',
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
                    value: '10 Users',
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
                        '${getFiatValue(tokenizedAsset.subscriptionAmount ?? 0)} ${tokenizedAsset.assetQuoteCurrency}',
                    extraValue: '\$100.12',
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
                    value: '%10',
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
                            '${getFiatValue(_progress)} ${tokenizedAsset.assetQuoteCurrency}',
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
                            '${getFiatValue(_maxValue)} ${tokenizedAsset.assetQuoteCurrency}',
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
                              value:
                                  normalizedProgress, // Show progress (0 to 1)
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
                                "${getFiatValue(_progress)} / ${getFiatValue(_maxValue)}",
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
              infoTile(
                notifier,
                'Asset Code',
                tokenizedAsset.assetCode ?? '',
              ),
              infoTile(
                notifier,
                'Sector',
                tokenizedAsset.assetSector ?? '',
              ),
              infoTile(
                notifier,
                'Sub-Sector',
                tokenizedAsset.assetSubSector ?? '',
              ),
              infoTile(
                notifier,
                'Type',
                assetType,
              ),
              infoTile(notifier, 'Asset Country',
                  tokenizedAsset.assetCountryLocation ?? ''),
              infoTile(
                notifier,
                'Address',
                tokenizedAsset.assetPhysicalAddress ?? '',
              ),
              if (tokenizedAsset.assetAlreadyExists == 1) ...[
                infoTile(
                  notifier,
                  'Original Asset Owner',
                  tokenizedAsset.assetOwnerName ?? '',
                ),
              ] else ...[
                infoTile(
                  notifier,
                  'Project Sponsor',
                  tokenizedAsset.assetOwnerName ?? '',
                ),
              ],
              infoTile(
                notifier,
                'Sales Window',
                '${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesStart!)} - ${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesEnd!)}',
              ),
              infoTile(
                notifier,
                'Cap Amount',
                '${getFiatValue(tokenizedAsset.capQuantity!)} ${tokenizedAsset.assetCode}',
              ),
              infoTile(
                notifier,
                'Cap Quantity',
                '${getFiatValue(tokenizedAsset.capQuantity!)} ${tokenizedAsset.assetCode}',
              ),
              infoTile(
                notifier,
                'Cap Duration',
                '${tokenizedAsset.assetLogo} days',
              ),
              infoTile(
                notifier,
                'Proceed Payout Cycle',
                tokenizedAsset.proceedCycle ?? '',
              ),
              infoTile(
                notifier,
                'Exempted Countries',
                '${tokenizedAsset.exemptedCountries!.replaceAll(',', ', ')}',
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        'Ownership and Legal Independence',
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              infoTile(
                notifier,
                'I confirm that this asset is free of all liens, mortgages, and outstanding loans.',
                '${tokenizedAsset.isFreeFromLiensAndEncumbrances == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that this asset is not pledged as collateral for any debts and has no use restrictions.',
                '${tokenizedAsset.undertakingNotCollateral == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that no third party has any claims, rights, or interests in this asset.',
                '${tokenizedAsset.undertakingNoClaims == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that this asset is free of any foreclosure, bankruptcy proceedings, legal disputes, judgments, or court-ordered payments.',
                '${tokenizedAsset.undertakingNoForeclosure == 1 ? 'Yes' : 'No'}',
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        'Regulatory Compliance and Approvals',
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              infoTile(
                notifier,
                'I confirm that this asset complies with all environmental and land use regulations and is free of violations.',
                '${tokenizedAsset.complianceNoViolation == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that all necessary permits, licenses, and approvals for the use and ownership of this asset are in place.',
                '${tokenizedAsset.complianceAllPermits == 1 ? 'Yes' : 'No'}',
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        'Outstanding Financial Responsibilities',
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              infoTile(
                notifier,
                'I confirm that there are no unpaid taxes, utility bills, fees, or other property-related expenses associated with this asset.',
                '${tokenizedAsset.outstandingFinancialRespNoDebts == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that this asset does not have any hidden liabilities or obligations that have not been disclosed.',
                '${tokenizedAsset.outstandingFinancialRespNoHiddenLiabilities == 1 ? 'Yes' : 'No'}',
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        'Risk Management and Insurance Coverage',
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              infoTile(
                notifier,
                'I confirm that this asset is adequately insured against risks such as fire, theft, and natural disasters.',
                '${tokenizedAsset.riskManagementFullyInsured == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that the declared value of this asset reflects its current market value and condition.',
                '${tokenizedAsset.riskManagementDeclaredValue == 1 ? 'Yes' : 'No'}',
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        'Physical Condition and Legal Status',
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              infoTile(
                notifier,
                'I confirm that this asset is not affected by undisclosed easements, rights of way, expropriation, or condemnation.',
                '?????',
              ),
              infoTile(
                notifier,
                'I confirm that this asset is structurally sound and has no unresolved maintenance or safety issues.',
                '${tokenizedAsset.physicalConditionSound == 1 ? 'Yes' : 'No'}',
              ),
              infoTile(
                notifier,
                'I confirm that this asset is not subject to any agreements, such as leases or contracts, that could limit its use or transfer.',
                '${tokenizedAsset.physicalConditionNolease == 1 ? 'Yes' : 'No'}',
              ),
              Card(
                elevation: notifier.isDark ? 0 : 3,
                shadowColor: Colors.black,
                color: notifier.gettilewihitecolor,
                margin: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 8.0),
                  child: ListTile(
                    title: Row(
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Proof of Existence',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontsemibold,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                            for (var item in tokenizedAsset
                                .assetTokenizationDocuments!) ...[
                              TextButton(
                                style: TextButton.styleFrom(
                                    padding: EdgeInsets.zero,
                                    minimumSize: Size(50, 30),
                                    tapTargetSize:
                                        MaterialTapTargetSize.shrinkWrap,
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
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              Card(
                elevation: notifier.isDark ? 0 : 3,
                shadowColor: Colors.black,
                color: notifier.gettilewihitecolor,
                margin: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 8.0),
                  child: ListTile(
                    title: Row(
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Proof of Payment Documents',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontsemibold,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                            for (var item
                                in tokenizedAsset.proofOfPaymentDocuments!) ...[
                              TextButton(
                                style: TextButton.styleFrom(
                                    padding: EdgeInsets.zero,
                                    minimumSize: Size(50, 30),
                                    tapTargetSize:
                                        MaterialTapTargetSize.shrinkWrap,
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
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Button(
                'Payout Proceeds',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: ProceedsPayOutViewPageConfig,
                  );
                },
              ),
              // SizedBox(height: height / 70),
              // ButtonOutlined(
              //   'Liquidate Asset',
              //   notifier.getwihitecolor,
              //   Colors.red,
              //   borderColor: Colors.red,
              //   onTap: () {
              //     appState.currentAction = PageAction(
              //       state: PageState.addPage,
              //       page: LiquidateAssetViewPageConfig,
              //     );
              //   },
              // ),
              SizedBox(
                height: height / 10,
              ),
            ]
          ],
        ),
      ),
    );
  }

  Widget assetTile(String imageUrl, String name, String type, isSubscribed) {
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
          title: Row(
            children: [
              Image.network(
                imageUrl,
                height: 35,
                width: 35,
                errorBuilder: (context, error, stackTrace) {
                  return Image.asset(
                    'assets/images/trovo.png',
                    height: 35,
                    width: 35,
                  );
                },
              ),
              SizedBox(width: 20),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    name,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getblck,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      type,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
          trailing: TextButton(
            onPressed: () async {},
            child: Container(
              width: width / 4,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    isSubscribed ? 'Subscribed' : 'Subscribe',
                    style: TextStyle(
                        fontFamily: fontsemibold,
                        fontSize: 12,
                        color: notifier.getblck),
                  ),
                  Icon(
                      isSubscribed
                          ? Icons.check_circle
                          : Icons.add_circle_rounded,
                      size: 20,
                      color: notifier.getblck),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
