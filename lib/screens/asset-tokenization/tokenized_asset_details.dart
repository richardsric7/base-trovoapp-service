import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/countdown.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TokenizedAssetDetail extends StatefulWidget {
  const TokenizedAssetDetail({Key? key}) : super(key: key);

  @override
  State<TokenizedAssetDetail> createState() => _TokenizedAssetDetail();
}

class _TokenizedAssetDetail extends State<TokenizedAssetDetail>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  String assetType = '';
  String fiatCurrency = '';
  late TokenizedAsset tokenizedAsset;
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

  List<DropdownMenuItem<Wallet>> get getStandardWallets {
    List<DropdownMenuItem<Wallet>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet));
    });
    return wallets;
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = appState.tokenizedAsset!;
    var assetTypes = appState.tokenizationData['assetTypes'];
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

    for (var i = 0; i < assetTypes.length; i++) {
      if (assetTypes[i]['id'].toString() == tokenizedAsset.assetType) {
        assetType = assetTypes[i]['assetType'];
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    inspect(appState.tokenizedAsset);

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              '',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(
              height: height / 50,
            ),
            // if (tokenizedAsset.assetLogo != null) ...[
            //   Image.memory(
            //     base64Decode(tokenizedAsset.assetLogo!),
            //     height: 50,
            //     width: 50,
            //     errorBuilder: (context, error, stackTrace) {
            //       return Image.asset(
            //         'assets/images/trovo.png',
            //         height: 50,
            //         width: 50,
            //       );
            //     },
            //   ),
            // ] else ...[
            //   Image.asset(
            //     'assets/images/trovo.png',
            //     height: 45,
            //     width: 45,
            //   ),
            // ],
            // SizedBox(
            //   height: height / 70,
            // ),
            Text(
              tokenizedAsset.assetName!.capitalizeEachWord(),
              style: TextStyle(
                fontSize: 20,
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(
              height: height / 70,
            ),
            Text(
              tokenizedAsset.assetCode!.toUpperCase(),
              style: TextStyle(
                fontSize: 15,
                fontFamily: fontbody,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(height: height / 70),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                if (tokenizedAsset.tokenizationStatus == 4) ...[
                  Countdown(
                    startDate: tokenizedAsset.salesStart!,
                    isColumn: true,
                  ),
                ] else if (tokenizedAsset.tokenizationStatus == 5) ...[
                  Card(
                    margin: EdgeInsets.zero,
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(5.0),
                        side: BorderSide(
                          color: notifier.getbluewhitecolor,
                          width: 1,
                        )),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: Padding(
                      padding: const EdgeInsets.all(5.0),
                      child: Text(
                        'Available',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                          overflow: TextOverflow.visible,
                        ),
                      ),
                    ),
                  ),
                ]
              ],
            ),
            if (tokenizedAsset.tokenizationStatus == 5) ...[
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  ElevatedButton(
                    onPressed: () {
                      if (appState.userInfo?.kycVerified == null ||
                          appState.userInfo?.kycVerified == 0) {
                        kycUnverifiedErrorPop(context);
                        return;
                      }
                      showBuyTokenPopup(context,
                          assetCode: tokenizedAsset.assetCode!.toUpperCase(),
                          onDone: (wallet) {
                        appState.setActiveWallet = wallet;
                        appState.tokenizedAsset = tokenizedAsset;
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: BuyTokensViewPageConfig,
                        );
                      }, dropdownItems: getStandardWallets);
                    },
                    style: ButtonStyle(
                      padding: MaterialStateProperty.all(
                        EdgeInsets.symmetric(vertical: 10, horizontal: 85),
                      ),
                      overlayColor: MaterialStateProperty.all<Color>(
                          notifier.getbluecolor90),
                      backgroundColor: MaterialStateProperty.all<Color>(
                          notifier.getbluewhitecolor),
                      side: MaterialStateProperty.all(
                        BorderSide(
                            color: notifier.getbluewhitecolor,
                            width: 1,
                            style: BorderStyle.solid),
                      ),
                      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                        const RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(
                            Radius.circular(10),
                          ),
                        ),
                      ),
                    ),
                    child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Image.asset('assets/images/money.png'),
                          SizedBox(width: 10),
                          Text(
                            'Buy',
                            style: TextStyle(
                              fontFamily: fontsemibold,
                              fontSize: 12,
                              color: notifier.getwihitecolor,
                            ),
                          ),
                        ]),
                  ),
                ],
              ),
            ] else if (tokenizedAsset.tokenizationStatus == 4) ...[
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  ElevatedButton(
                    onPressed: () {
                      if (appState.userInfo?.kycVerified == null ||
                          appState.userInfo?.kycVerified == 0) {
                        kycUnverifiedErrorPop(context);
                        return;
                      }

                      showSubscribePopup(
                        context,
                        asset: tokenizedAsset,
                        onDone: (amount) async {
                          await subscribeTokenizedAsset(
                            amount: double.parse(amount),
                            tokenizedAssetID: tokenizedAsset.id!,
                          );
                          setState(() {
                            tokenizedAsset.expressedInterest = true;
                            tokenizedAsset.expressedInterestAmount =
                                double.parse(amount);
                          });
                        },
                      );
                    },
                    style: ButtonStyle(
                      padding: MaterialStateProperty.all(
                        EdgeInsets.symmetric(vertical: 10, horizontal: 50),
                      ),
                      overlayColor: MaterialStateProperty.all<Color>(
                          notifier.getsplashgrey),
                      backgroundColor: MaterialStateProperty.all<Color>(
                        tokenizedAsset.expressedInterest ?? false
                            ? notifier.getaddsubwalletgrey
                            : notifier.getwihitecolor,
                      ),
                      side: MaterialStateProperty.all(
                        BorderSide(
                            color: tokenizedAsset.expressedInterest ?? false
                                ? notifier.getaddsubwalletgrey
                                : notifier.getbluewhitecolor,
                            width: 1,
                            style: BorderStyle.solid),
                      ),
                      shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                        const RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(
                            Radius.circular(10),
                          ),
                        ),
                      ),
                    ),
                    child: Container(
                      child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            if (tokenizedAsset.expressedInterest ?? false) ...[
                              Icon(
                                Icons.check_circle_rounded,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                              SizedBox(width: 10),
                              Text(
                                'Interest Expressed',
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                  fontSize: 12,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ] else ...[
                              Icon(
                                Icons.add_circle_rounded,
                                size: 20,
                                color: tokenizedAsset.expressedInterest ?? false
                                    ? notifier.getwihitecolor
                                    : notifier.getbluewhitecolor,
                              ),
                              SizedBox(width: 10),
                              Text(
                                'Express Interest',
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                  fontSize: 12,
                                  color:
                                      tokenizedAsset.expressedInterest ?? false
                                          ? notifier.getwihitecolor
                                          : notifier.getbluewhitecolor,
                                ),
                              ),
                            ]
                          ]),
                    ),
                  ),
                ],
              ),
            ] else if (tokenizedAsset.tokenizationStatus == 6) ...[
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  ElevatedButton(
                    onPressed: () {
                      popup(
                        context,
                        title: "comingsoon".tr(),
                        message: "p2pwillbelaunchingsoon".tr(),
                        bodyColor: notifier.getbluewhitecolor,
                      );
                      // _launchUrl();
                    },
                    style: ButtonStyle(
                      padding: MaterialStateProperty.all(
                        EdgeInsets.symmetric(vertical: 10, horizontal: 50),
                      ),
                      overlayColor: MaterialStateProperty.all<Color>(
                          notifier.getbluecolor90),
                      backgroundColor: MaterialStateProperty.all<Color>(
                          notifier.getbluewhitecolor),
                      side: MaterialStateProperty.all(
                        BorderSide(
                            color: notifier.getbluewhitecolor,
                            width: 1,
                            style: BorderStyle.solid),
                      ),
                      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                        const RoundedRectangleBorder(
                          borderRadius: BorderRadius.all(
                            Radius.circular(10),
                          ),
                        ),
                      ),
                    ),
                    child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Image.asset('assets/images/money.png'),
                          SizedBox(width: 10),
                          Text(
                            'Buy on TrovoP2P',
                            style: TextStyle(
                              fontFamily: fontsemibold,
                              fontSize: 12,
                              color: notifier.getwihitecolor,
                            ),
                          ),
                        ]),
                  ),
                ],
              ),
            ],
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Container(
                child: Center(
                  child: Column(
                    children: [
                      SizedBox(height: height / 70),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Padding(
                            padding:
                                const EdgeInsets.symmetric(horizontal: 30.0),
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
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode!.toUpperCase()}',
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
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeSold! * tokenizedAsset.pricePerToken!)} ${fiatCurrency}',
                  extraValue: '',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  notifier,
                  label: 'Tokens for Sale',
                  value:
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeSold ?? 0)} ${tokenizedAsset.assetCode!.toUpperCase()}',
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
                  label: 'Total Quantity Held',
                  value:
                      '${getFiatValue(tokenizedAsset.subscriptionAmount ?? 0)} ${tokenizedAsset.assetCode!.toUpperCase()}',
                  extraValue: '',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  notifier,
                  label: 'Value of Quantity Held',
                  value:
                      '${getFiatValue(tokenizedAsset.subscriptionAmount == null ? 0 : tokenizedAsset.subscriptionAmount! * tokenizedAsset.pricePerToken!)} ${fiatCurrency}',
                  extraValue: '',
                ),
              ],
            ),
            SizedBox(height: height / 50),
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
              tokenizedAsset.assetCode!.toUpperCase(),
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
                iso2Countries[tokenizedAsset.assetCountryLocation] ?? ""),
            infoTile(
              notifier,
              'Address',
              tokenizedAsset.assetPhysicalAddress ?? '',
            ),
            if (tokenizedAsset.assetAlreadyExists == 0) ...[
              infoTile(
                notifier,
                "Project Strategic Objectives",
                tokenizedAsset.projectStrategicObjectives ?? "",
              ),
              infoTile(
                notifier,
                "Project Development Timeline",
                tokenizedAsset.projectDevelopmentTimeline ?? "",
              ),
              infoTile(
                notifier,
                "Key Milestones & Dates",
                tokenizedAsset.projectKeyMilestoneAndDates ?? "",
              ),
              infoTile(
                notifier,
                "Project Scope",
                tokenizedAsset.projectScope ?? "",
              ),
              infoTile(
                notifier,
                "Project Economic Benefits",
                tokenizedAsset.projectEconomicBenefits ?? "",
              ),
              infoTile(
                notifier,
                "Expected No. of Job to be Created",
                tokenizedAsset.projectExpectedNoOfJobs ?? "",
              ),
              infoTile(
                notifier,
                "Project Intended Social Benefits",
                tokenizedAsset.projectIntendedSocialBenefits ?? "",
              ),
              infoTile(
                notifier,
                "Technical Partners",
                tokenizedAsset.projectTechnicalPartners ?? "",
              ),
              infoTile(
                notifier,
                "Financial Partners",
                tokenizedAsset.projectFinancialPartners ?? "",
              ),
              infoTile(
                notifier,
                "Project Intended Social Benefits",
                tokenizedAsset.projectIntendedSocialBenefits ?? "",
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        "Asset Financial Performance",
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
                "Estimated Project IRR",
                formatNumber(tokenizedAsset.estimatedProjectIRR ?? 0),
              ),
              infoTile(
                notifier,
                "Estimated Project ROI",
                formatNumber(tokenizedAsset.estimatedProjectROI ?? 0),
              ),
              infoTile(
                notifier,
                "Estimated Project NPV at Launch (Day 1)",
                formatNumber(tokenizedAsset.estimatedProjectNPV ?? 0),
              ),
              infoTile(
                notifier,
                "Estimated Project Payback Periods (in Months)",
                tokenizedAsset.estimatedProjectPaybackPeriodsInMonths ?? '',
              ),
              infoTile(
                notifier,
                "All Key Assumptions Including Values Assumed",
                tokenizedAsset.keyAssumptionsList ?? '',
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 30.0),
                    child: SizedBox(
                      width: width / 1.2,
                      child: Text(
                        "Project Risk Assessment",
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
                "Legal Risks Identified",
                tokenizedAsset.projectIdentifiedLegalRisks ?? '',
              ),
              infoTile(
                notifier,
                "Regulatory Risks Identified",
                tokenizedAsset.projectIdentifiedRegulatoryRisks ?? '',
              ),
              infoTile(
                notifier,
                "Operational/Execution Risks Identified",
                tokenizedAsset.projectIdentifiedOperationalOrExecutionRisks ??
                    '',
              ),
              infoTile(
                notifier,
                "Market Risks Identified",
                tokenizedAsset.projectIdentifiedMarketRisks ?? '',
              ),
              infoTile(
                notifier,
                "Other Relevant Risks Identified",
                tokenizedAsset.projectIdentifiedOtherRelevantRisks ?? '',
              ),
            ],
            infoTile(
              notifier,
              'Regulator',
              regulatorName.toLowerCase().capitalizeEachWord(),
            ),
            infoTile(
              notifier,
              'Asset Custodian',
              tokenizedAsset.approvedAssetCustodianInfo?.assetCustodianName ??
                  '',
            ),
            infoTile(
              notifier,
              'Asset Manager',
              tokenizedAsset.assetManagerInfo?.assetManagerName ?? '',
            ),
            infoTile(
              notifier,
              'Issuing House',
              tokenizedAsset.assetIssuingHouseInfo?.assetIssuingHouseName
                      ?.toLowerCase()
                      .capitalizeEachWord() ??
                  '',
            ),
            infoTile(
              notifier,
              'Rating Agency',
              tokenizedAsset.assetOwnerName ?? '',
            ),
            infoTile(
              notifier,
              'Sales Window',
              '${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesStart!)} - ${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesEnd!)}',
            ),
            if ((tokenizedAsset.capAmountInFiat ?? 0) > 0) ...[
              infoTile(
                notifier,
                'Cap Amount',
                '${getFiatValue(double.parse(tokenizedAsset.capAmountInFiat!.toString()))} ${fiatCurrency}',
              ),
              infoTile(
                notifier,
                'Cap Quantity',
                '${getFiatValue(tokenizedAsset.capQuantity!)} ${tokenizedAsset.assetCode!.toUpperCase()}',
              ),
              infoTile(
                notifier,
                'Cap Duration',
                '${tokenizedAsset.capDurationInDays} days',
              ),
            ],
            infoTile(
              notifier,
              'Proceed Payout Cycle',
              tokenizedAsset.proceedCycle?.toLowerCase().capitalizeEachWord() ??
                  '',
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
              'Free of liens, mortgages, and outstanding loans.',
              '${tokenizedAsset.isFreeFromLiensAndEncumbrances == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'Not pledged as collateral for any debts and has no use restrictions.',
              '${tokenizedAsset.undertakingNotCollateral == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'No third party has any claims, rights, or interests in this asset.',
              '${tokenizedAsset.undertakingNoClaims == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'Free of any foreclosure, bankruptcy proceedings, legal disputes, judgments, or court-ordered payments.',
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
              'Complies with all environmental and land use regulations and is free of violations.',
              '${tokenizedAsset.complianceNoViolation == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'All necessary permits, licenses, and approvals for the use and ownership of this asset are in place.',
              '${tokenizedAsset.complianceNoViolation == 1 ? 'Yes' : 'No'}',
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
              'No unpaid taxes, utility bills, fees, or other property-related expenses associated with this asset.',
              '${tokenizedAsset.outstandingFinancialRespNoDebts == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'Asset does not have any hidden liabilities or obligations that have not been disclosed. ',
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
              'Asset is adequately insured against risks such as fire, theft, and natural disasters.',
              '${tokenizedAsset.riskManagementFullyInsured == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'The declared value of this asset reflects its current market value and condition.',
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
              'Asset is not affected by undisclosed easements, rights of way, expropriation, or condemnation.',
              '${tokenizedAsset.physicalConditionNoUndisclosedEasements == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'Asset is structurally sound and has no unresolved maintenance or safety issues.',
              '${tokenizedAsset.physicalConditionSound == 1 ? 'Yes' : 'No'}',
            ),
            infoTile(
              notifier,
              'Asset is not subject to any agreements, such as leases or contracts, that could limit its use or transfer.',
              '${tokenizedAsset.physicalConditionNolease == 1 ? 'Yes' : 'No'}',
            ),
            if (tokenizedAsset.assetTokenizationDocuments != null &&
                tokenizedAsset.assetTokenizationDocuments!.isNotEmpty) ...[
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
                              'Asset Verification Documents',
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
            ],
            // Card(
            //   elevation: notifier.isDark ? 0 : 3,
            //   shadowColor: Colors.black,
            //   color: notifier.gettilewihitecolor,
            //   margin: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            //   child: Padding(
            //     padding: const EdgeInsets.symmetric(vertical: 8.0),
            //     child: ListTile(
            //       title: Row(
            //         children: [
            //           Column(
            //             crossAxisAlignment: CrossAxisAlignment.start,
            //             children: [
            //               Text(
            //                 'Other Documents',
            //                 style: TextStyle(
            //                   fontSize: 13,
            //                   fontFamily: fontsemibold,
            //                   color: notifier.getbluewhitecolor,
            //                 ),
            //               ),
            //               for (var item in tokenizedAsset
            //                   .assetTokenizationDocuments!) ...[
            //                 TextButton(
            //                   style: TextButton.styleFrom(
            //                       padding: EdgeInsets.zero,
            //                       minimumSize: Size(50, 30),
            //                       tapTargetSize:
            //                           MaterialTapTargetSize.shrinkWrap,
            //                       alignment: Alignment.centerLeft),
            //                   onPressed: () {
            //                     var fileUrl = item.documentUrl;
            //                     if (fileUrl!.isNotEmpty &&
            //                         fileUrl.endsWith('.pdf')) {
            //                       appState.pdfUrl = fileUrl;
            //                       appState.currentAction = PageAction(
            //                           state: PageState.addPage,
            //                           page: PdfViewPageConfig);

            //                       return;
            //                     }

            //                     appState.goToWebView(fileUrl);
            //                   },
            //                   child: Text(
            //                     item.documentTitle ?? '',
            //                     style: TextStyle(
            //                       decoration: TextDecoration.underline,
            //                       fontSize: 12,
            //                       fontFamily: fontbody,
            //                       color: notifier.getbluewhitecolor,
            //                     ),
            //                   ),
            //                 ),
            //               ],
            //             ],
            //           ),
            //         ],
            //       ),
            //     ),
            //   ),
            // ),
            SizedBox(height: height / 20),
          ],
        ),
      ),
    );
  }

  subscribeTokenizedAsset(
      {required double amount, required String tokenizedAssetID}) async {
    try {
      showLoader(context);

      String requestBody = jsonEncode({
        'amount': amount,
      });

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization/expressed-interests/${tokenizedAssetID}',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.publicKey!,
      );

      if (responseData['statusCode'] == 200) {
        hideLoader(context);
      } else {
        hideLoader(context);
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'].toString().isEmpty
              ? responseData['data']['error']
              : responseData['data']['message'],
        );
      }
    } catch (e) {
      // print(e);
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}
