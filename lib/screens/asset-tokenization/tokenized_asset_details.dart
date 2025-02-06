import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/countdown.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
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

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    var tokenizedAsset = appState.tokenizedAsset!;
    var countdown =
        tokenizedAsset.salesStart!.difference(DateTime.now()).inDays;
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
              tokenizedAsset.assetName!,
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
              '${tokenizedAsset.assetCode} token',
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
                Countdown(
                  startDate: tokenizedAsset.salesStart!,
                  isColumn: true,
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            if (countdown < 0) ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  ElevatedButton(
                    onPressed: () {
                      showBuyTokenPopup(context,
                          assetCode: tokenizedAsset.assetCode!,
                          onDone: (publicKey) {
                        var wallet = appState.userInfo!.getWallet(publicKey);
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
                        EdgeInsets.symmetric(vertical: 10, horizontal: 95),
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
            ] else ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  ElevatedButton(
                    onPressed: tokenizedAsset.isSubscribed ?? false
                        ? null
                        : () {
                            showSubscribePopup(
                              context,
                              assetCode: tokenizedAsset.assetCode!,
                              onDone: (amount) async {
                                setState(() {
                                  print(amount);
                                  tokenizedAsset.isSubscribed = true;
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
                        tokenizedAsset.isSubscribed ?? false
                            ? notifier.getaddsubwalletgrey
                            : notifier.getwihitecolor,
                      ),
                      side: MaterialStateProperty.all(
                        BorderSide(
                            color: tokenizedAsset.isSubscribed ?? false
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
                            if (tokenizedAsset.isSubscribed ?? false) ...[
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
                                color: tokenizedAsset.isSubscribed ?? false
                                    ? notifier.getwihitecolor
                                    : notifier.getbluewhitecolor,
                              ),
                              SizedBox(width: 10),
                              Text(
                                'Express Interest',
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                  fontSize: 12,
                                  color: tokenizedAsset.isSubscribed ?? false
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
            ],
            SizedBox(height: height / 90),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Container(
                child: Center(
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 70,
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
                  label: 'Asset Value',
                  value:
                      '${tokenizedAsset.assetQuoteCurrency}${getFiatValue(tokenizedAsset.assetCurrentValue!)}',
                  extraValue: '\$4,390.23',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  label: 'Total Supply',
                  value:
                      '${getFiatValue(tokenizedAsset.numberOfTokenToBeSold!)} ${tokenizedAsset.assetCode}',
                  extraValue: '',
                ),
              ],
            ),
            SizedBox(height: height / 70),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  label: 'Funding Currency',
                  value: 'CNGN',
                  extraValue: '',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  label: 'Price Per Token',
                  value:
                      '${getFiatValue(tokenizedAsset.pricePerToken!)} ${tokenizedAsset.assetQuoteCurrency}',
                  extraValue: '\$0.12',
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
              tokenizedAsset.assetType ?? '',
            ),
            infoTile(notifier, 'Asset Country',
                tokenizedAsset.assetCountryLocation ?? ''),
            infoTile(
              notifier,
              'Address',
              tokenizedAsset.assetPhysicalAddress ?? '',
            ),
            infoTile(
              notifier,
              'Issuer',
              tokenizedAsset.assetIssuer ?? '',
            ),
            infoTile(
              notifier,
              'Issuer Website',
              'www.${tokenizedAsset.assetCode!.toLowerCase()}.com',
            ),
            infoTile(notifier, 'Asset Token Total Supply',
                '${getFiatValue(tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetCode}'),
            infoTile(
              notifier,
              'Asset Tokens Quantity Purchased',
              '${getFiatValue(tokenizedAsset.amount ?? 0)}',
            ),
            infoTile(
              notifier,
              'Total Subscribed Users',
              '2,000',
            ),
            infoTile(
              notifier,
              'Sales Window',
              '${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesStart!)} - ${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesEnd!)}',
            ),
            infoTile(
              notifier,
              'Cap Quantity',
              '${getFiatValue(tokenizedAsset.capQuantity!)} ${tokenizedAsset.assetCode}',
            ),
            infoTile(
              notifier,
              'Cap Duration',
              '${tokenizedAsset.capDurationInDays} days',
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

  Widget infoCard(
      {required String label, required String value, String? extraValue}) {
    return Container(
      width: width / 2.3,
      // height: height / 5.5,
      child: Card(
        shadowColor: Colors.black,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(15.0),
        ),
        color: notifier.isDark
            ? notifier.getbluecolor90
            : notifier.getaddsubwalletgrey,
        child: TextButton(
          onPressed: () {
            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: TotalSalesViewPageConfig,
            );
          },
          child: Row(
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(
                    height: height / 70,
                  ),
                  Text(
                    label,
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 70,
                  ),
                  SizedBox(
                    width: 130,
                    child: Text(
                      value,
                      textAlign: TextAlign.start,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  if (extraValue != null) ...[
                    SizedBox(
                      height: height / 70,
                    ),
                    Text(
                      extraValue,
                      textAlign: TextAlign.start,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  String getFiatValue(double amount) {
    if (amount < 99000000000) {
      return formatNumberShort(amount);
    }
    return formatHistoryNumber(amount, 99000000000)
        .toString()
        .replaceAll(',', '');
  }
}
