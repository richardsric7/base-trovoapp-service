import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
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
              '${tokenizedAsset.assetName} token',
              style: TextStyle(
                fontSize: 15,
                fontFamily: fontbody,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Column(
                  children: [
                    TextButton(
                      onPressed: () {
                        showSubscribePopup(context,
                            assetCode: tokenizedAsset.assetCode!,
                            onDone: (publicKey) {
                          setState(() {
                            tokenizedAsset.isSubscribed = true;
                          });
                        }, dropdownItems: getStandardWallets);
                      },
                      child: Column(
                        children: [
                          Icon(
                            tokenizedAsset.isSubscribed ?? false
                                ? Icons.check_circle
                                : Icons.add_circle_rounded,
                            size: 25,
                            color: notifier.getbluewhitecolor,
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                          Text(
                            tokenizedAsset.isSubscribed ?? false
                                ? 'Subscribed'
                                : 'Subscribe',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                Column(
                  children: [
                    TextButton(
                      onPressed: () {
                        showBuyTokenPopup(context,
                            assetCode: tokenizedAsset.assetCode!,
                            onDone: (publicKey) {
                          var wallet = appState.userInfo!.getWallet(publicKey);
                          appState.setActiveWallet = wallet;
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: BuyTokensViewPageConfig,
                          );
                        }, dropdownItems: getStandardWallets);
                      },
                      child: Column(
                        children: [
                          Icon(
                            CupertinoIcons.cart_fill,
                            size: 25,
                            color: notifier.getbluewhitecolor,
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                          Text(
                            'Buy',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                )
              ],
            ),
            SizedBox(height: height / 90),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Container(
                // shadowColor: Colors.black,
                // shape: RoundedRectangleBorder(
                //   borderRadius: BorderRadius.circular(15.0),
                // ),
                // color: notifier.isDark
                //     ? notifier.getbluecolor90
                //     : notifier.getaddsubwalletgrey,
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
                  label: 'Total Supply',
                  value:
                      '${tokenizedAsset.numberOfTokenToBeSold} ${tokenizedAsset.assetCode}',
                  extraValue: '',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  label: 'Total Subscribed',
                  value: '800',
                  extraValue: '',
                ),
              ],
            ),
            SizedBox(height: height / 70),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  label: 'Price Per Asset',
                  value:
                      '${tokenizedAsset.pricePerToken} ${tokenizedAsset.assetQuoteCurrency}',
                  extraValue: '',
                  // extraValue: '\$2,205',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  label: 'Funding Currency',
                  value: '${tokenizedAsset.assetQuoteCurrency}',
                  extraValue: '',
                ),
              ],
            ),
            SizedBox(height: height / 70),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                infoCard(
                  label: 'Subscription Amount',
                  value: '0 ${tokenizedAsset.assetQuoteCurrency}',
                  extraValue: '',
                  // extraValue: '\$2,205',
                ),
                SizedBox(
                  width: width / 50,
                ),
                infoCard(
                  label: 'Amount Bought',
                  value: '0 ${tokenizedAsset.assetQuoteCurrency}',
                  extraValue: '',
                ),
              ],
            ),
            // SizedBox(height: height / 70),
            // Row(
            //   children: [
            //     Padding(
            //       padding: const EdgeInsets.symmetric(horizontal: 30.0),
            //       child: Text(
            //         'Purchase with cNGN',
            //         style: TextStyle(
            //           fontSize: 15,
            //           fontFamily: fontsemibold,
            //           color: notifier.getbluewhitecolor,
            //         ),
            //       ),
            //     ),
            //   ],
            // ),
            // Row(
            //   children: [
            //     Padding(
            //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
            //       child: TextButton(
            //         onPressed: () {},
            //         child: Text(
            //           'Tap to fund wallet now',
            //           style: TextStyle(
            //             decoration: TextDecoration.underline,
            //             fontSize: 15,
            //             fontFamily: fontbody,
            //             color: notifier.getbluewhitecolor,
            //           ),
            //         ),
            //       ),
            //     ),
            //   ],
            // ),
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
              'Asset Category',
              tokenizedAsset.assetSector ?? '',
            ),
            infoTile(notifier, 'Asset Country',
                tokenizedAsset.assetCountryLocation ?? ''),
            infoTile(
              notifier,
              'Asset Location Address',
              tokenizedAsset.assetPhysicalAddress ?? '',
            ),
            infoTile(
              notifier,
              'Asset Issuer',
              tokenizedAsset.assetIssuer ?? '',
            ),
            infoTile(
              notifier,
              'Asset Issuer Website',
              'www.${tokenizedAsset.assetCode!.toLowerCase()}.com',
            ),
            infoTile(notifier, 'Asset Token Total Supply',
                '${tokenizedAsset.numberOfTokenToBeIssued.toString()} ${tokenizedAsset.assetCode}'),
            infoTile(
              notifier,
              'Asset Tokens Quantity Purchased',
              '400',
            ),
            infoTile(
              notifier,
              'Total Subscribed Users',
              '2,000',
            ),
            infoTile(
              notifier,
              'Price Per Asset Token',
              '${tokenizedAsset.pricePerToken} ${tokenizedAsset.assetQuoteCurrency}',
            ),
            infoTile(
              notifier,
              'Asset Token Purchase Method',
              tokenizedAsset.assetQuoteCurrency ?? '',
            ),
            infoTile(
              notifier,
              'Asset Token Sales Window',
              '${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesStart!)} - ${DateFormat('yyyy-MM-dd').format(tokenizedAsset.salesEnd!)}',
            ),
            infoTile(
              notifier,
              'Token Sale Cap',
              '${tokenizedAsset.capQuantity} ${tokenizedAsset.assetCode}',
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
              'Payout Method',
              tokenizedAsset.proceedPayoutCurrency ?? '',
            ),
            infoTile(
              notifier,
              'Exempted Countries',
              '${tokenizedAsset.exemptedCountries!.replaceAll(',', ', ')}',
            ),
            infoTile(
              notifier,
              'Additional Requirements',
              tokenizedAsset.additionalKYCRequirements!.replaceAll(',', ', '),
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
                            'Verified Proof of Existence',
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
                  Text(
                    value,
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
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
}
