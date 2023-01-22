import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WalletDetails extends StatefulWidget {
  const WalletDetails({Key? key}) : super(key: key);

  @override
  State<WalletDetails> createState() => _WalletDetailsState();
}

class _WalletDetailsState extends State<WalletDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late TabController _tabController;
  late DataProvider appState;
  late Wallet wallet;
  int tabLength = 1;
  int activeTabIndex = 0;
  late bool localHideBalance;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    appState = Provider.of<DataProvider>(context, listen: false);
    localHideBalance = appState.hideBalances;

    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );
  }

  void tabListener() {
    // Tab Changed swiping to a new tab
    activeTabIndex = _tabController.index;
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    if (wallet.unClaimedAssets != null && wallet.unClaimedAssets!.length > 0) {
      tabLength = 2;
    } else {
      tabLength = 1;
    }

    if (tabLength != _tabController.length) {
      // change the length of tabController too or you will have an error
      _tabController = TabController(length: tabLength, vsync: this);
      _tabController.addListener(tabListener);
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              assetsTabs(),
            ],
          ),
        ),
      ),
    );
  }

  Widget assetsTabs() {
    return Container(
      height: height / 1.1,
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          Padding(
            padding:
                EdgeInsets.symmetric(horizontal: tabLength > 1 ? 0.0 : 100.0),
            child: TabBar(
              controller: _tabController,
              labelColor: notifier.getbluewhitecolor,
              indicatorColor: notifier.getbluewhitecolor,
              labelStyle: TextStyle(
                fontSize: 14.sp,
                fontWeight: FontWeight.w600,
                fontFamily: fontsemibold,
              ),
              tabs: [
                Tab(
                  height: 20,
                  text: LanguageEn.assets,
                ),
                if (wallet.unClaimedAssets != null && tabLength == 2) ...[
                  Tab(
                    height: 20,
                    text:
                        '${LanguageEn.pending} (${wallet.unClaimedAssets == null ? 0 : wallet.unClaimedAssets!.length})',
                  ),
                ],
                // Tab(
                //   height: 20,
                //   text: LanguageEn.nfts,
                // ),
              ],
            ),
          ),
          Positioned(
            child: Column(
              children: [
                SizedBox(
                  height: height / 40,
                ),
                Expanded(
                  child: TabBarView(
                    controller: _tabController,
                    children: [
                      Container(
                        child: Padding(
                          padding: const EdgeInsets.fromLTRB(0, 10.0, 0, 0),
                          child: SingleChildScrollView(
                            child: Column(
                              children: [
                                WalletSlide(
                                  backColor: notifier.getbluecolor,
                                  foreColor: wihitecolor,
                                  alias: wallet.alias!.capitalizeFirst!,
                                  totalBalance:
                                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallet.claimedAssets!)} ${appState.defaultCurrency}',
                                  fiatBalance: appState.defaultCurrency == 'USD'
                                      ? null
                                      : '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallet.claimedAssets!)} USD',
                                  initialHiddenState: appState.hideBalances,
                                  onHiddenStateChanged: (state) => {
                                    setState(
                                      () => {
                                        localHideBalance = state,
                                      },
                                    )
                                  },
                                ),
                                SizedBox(
                                  height: height / 30,
                                ),
                                if (wallet.claimedAssets!.length > 0) ...[
                                  for (var asset in wallet.claimedAssets!) ...[
                                    GestureDetector(
                                        onTap: () {
                                          appState.viewData = {
                                            'assetCode': asset.assetCode,
                                            'assetIssuer': asset.assetIssuer,
                                            'walletPublicKey': wallet.publicKey,
                                          };
                                          appState.currentAction = PageAction(
                                            state: PageState.addPage,
                                            page: AssetDetailsViewPageConfig,
                                          );
                                        },
                                        child: tiles(asset)),
                                  ],
                                ] else ...[
                                  Container(
                                    height: height / 3,
                                    child: Padding(
                                        padding: const EdgeInsets.fromLTRB(
                                            10, 28.0, 10, 0),
                                        child: Center(
                                          child: Text(
                                            LanguageEn.noassets,
                                            style: TextStyle(
                                              fontSize: 13,
                                              fontWeight: FontWeight.bold,
                                              fontFamily: fontsemibold,
                                              color: notifier.getblck,
                                            ),
                                          ),
                                        )),
                                  ),
                                ],
                                SizedBox(
                                  height: height / 22,
                                ),
                                Button(
                                  LanguageEn.back,
                                  notifier.getbluecolor,
                                  wihitecolor,
                                  onTap: () {
                                    appState.currentAction = PageAction(
                                        state: PageState.replaceAll,
                                        page: BottomHomePageConfig);
                                  },
                                ),
                                SizedBox(
                                  height: height / 10,
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                      if (tabLength == 2) ...[
                        Padding(
                          padding: const EdgeInsets.fromLTRB(0, 8, 0, 0),
                          child: Container(
                            child: SingleChildScrollView(
                              child: Column(
                                children: [
                                  if (wallet.unClaimedAssets != null &&
                                      wallet.unClaimedAssets!.length > 0) ...[
                                    for (var asset
                                        in wallet.unClaimedAssets!) ...[
                                      GestureDetector(
                                        onTap: () {
                                          setState(() {
                                            activeTabIndex =
                                                _tabController.index;
                                          });

                                          appState.viewData = {
                                            'assetCode': asset.assetCode,
                                            'assetIssuer': asset.assetIssuer,
                                            'walletPublicKey': wallet.publicKey,
                                          };
                                          appState.currentAction = PageAction(
                                            state: PageState.addPage,
                                            page:
                                                PendingAssetDetailsViewPageConfig,
                                          );
                                        },
                                        child: tiles(asset),
                                      ),
                                    ],
                                  ] else ...[
                                    Container(
                                      height: height / 4,
                                      child: Padding(
                                          padding: const EdgeInsets.fromLTRB(
                                              10, 28.0, 10, 0),
                                          child: Center(
                                            child: Text(
                                              LanguageEn.nopendingassets,
                                              style: TextStyle(
                                                fontSize: 13,
                                                fontWeight: FontWeight.bold,
                                                fontFamily: fontsemibold,
                                                color: notifier.getblck,
                                              ),
                                            ),
                                          )),
                                    ),
                                  ],
                                  SizedBox(
                                    height: height / 22,
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      ],
                      // Container(
                      //   height: height / 2,
                      //   child: SingleChildScrollView(
                      //     child: Column(
                      //       children: [
                      //         // in situations where the blockchain has an issue,
                      //         // some values can be returned as null or empty
                      //         // so always null check for such situations
                      //         // if (nfts != null && nfts != {}) ...[
                      //         //   if (nfts[activeWallet!.publicKey] != null &&
                      //         //       nfts[activeWallet!.publicKey].length >
                      //         //           0) ...[
                      //         gridView(),
                      //         //     SizedBox(height: 600),
                      //         //   ] else ...[
                      //         //     showEmptyNFTs(),
                      //         //   ]
                      //         // ] else ...[
                      //         //   showEmptyNFTs(),
                      //         // ],
                      //       ],
                      //     ),
                      //   ),
                      // ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget showEmptyNFTs() {
    return Container(
      height: height / 3,
      child: Padding(
          padding: const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
          child: Center(
            child: Text(
              LanguageEn.noNFTs,
              style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.bold,
                fontFamily: fontsemibold,
                color: notifier.getblck,
              ),
            ),
          )),
    );
  }

  Widget gridView() {
    return Container(
      height: height / 1.13,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(0, 10.0, 0, 0),
        child: GridView(
          padding: const EdgeInsets.fromLTRB(15, 0, 15, 70),
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              mainAxisSpacing: 20,
              crossAxisSpacing: 20,
              childAspectRatio: 1.05),
          children: [
            nftCard(
              "assets/images/awka-paws.svg",
              'AWKA PAWS',
              'GBCVE....UJKKGA',
              Colors.blue,
            ),
            nftCard(
              "assets/images/warri-wolves.svg",
              'WARRI WOLVES',
              'GBCVE....UJKKGA',
              Colors.green,
            ),
            nftCard(
              "assets/images/accra-goats.svg",
              'ACCRA GOATS',
              'GBCVE....UJKKGA',
              Colors.red,
            ),
            nftCard(
              "assets/images/awka-paws.svg",
              'AWKA PAWS',
              'GBCVE....UJKKGA',
              Colors.blue,
            ),
            nftCard(
              "assets/images/warri-wolves.svg",
              'WARRI WOLVES',
              'GBCVE....UJKKGA',
              Colors.green,
            ),
            nftCard(
              "assets/images/accra-goats.svg",
              'ACCRA GOATS',
              'GBCVE....UJKKGA',
              Colors.red,
            ),
            nftCard(
              "assets/images/awka-paws.svg",
              'AWKA PAWS',
              'GBCVE....UJKKGA',
              Colors.blue,
            ),
            nftCard(
              "assets/images/warri-wolves.svg",
              'WARRI WOLVES',
              'GBCVE....UJKKGA',
              Colors.green,
            ),
            nftCard(
              "assets/images/accra-goats.svg",
              'ACCRA GOATS',
              'GBCVE....UJKKGA',
              Colors.red,
            ),
          ],
        ),
      ),
    );
  }

  Widget nftCard(image, title, subtitle, color) {
    return Card(
      elevation: 5,
      shadowColor: Colors.black,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      color: color,
      child: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SvgPicture.asset(
              image,
              width: 75,
              height: 75,
            ),
            SizedBox(
              height: 10,
            ),
            Text(
              title,
              style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.bold,
                fontFamily: fontsemibold,
                color: notifier.getblck,
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
              child: Text(
                subtitle,
                style: TextStyle(
                  fontSize: 10,
                  fontFamily: fontbody,
                  color: notifier.getblck,
                ),
              ),
            ),
          ],
        ),
      ), //SizedBox
    );
  }

  Widget tiles(Asset asset) {
    return Card(
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      elevation: notifier.isDark ? 0 : 5,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
            title: Row(
              children: [
                Image.network(
                  asset.imageUrl!,
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
                      asset.assetCode.toString().isEmpty
                          ? 'XBN'
                          : asset.assetCode!,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        "${getFiatRate(asset.usdPrice!.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                  ],
                )
              ],
            ),
            trailing: Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  getBalance(formatNumber(asset.amount!)),
                  style: TextStyle(
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getblck,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                  child: Text(
                    getBalance(
                        '${calculateFiatValue(asset.amount.toString(), asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}'),
                    style: TextStyle(
                      fontSize: 9,
                      fontFamily: fontbody,
                      color: notifier.getblck,
                    ),
                  ),
                ),
              ],
            )),
      ),
    );
  }

  String getBalance(String balance) {
    String text;
    if (appState.hideBalances) text = hideBalanceText;

    if (localHideBalance)
      text = hideBalanceText;
    else
      text = balance;

    return text;
  }
}
