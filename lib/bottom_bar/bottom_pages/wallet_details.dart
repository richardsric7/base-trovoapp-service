import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/WalletSlides.dart';
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
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  Wallet? activeWallet;
  List<Wallet>? wallets;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int activeTabIndex = 0;
  late bool localHideBalance;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    appState = Provider.of<DataProvider>(context, listen: false);
    localHideBalance = appState.hideBalances;
  }

  void tabListener() {
    print('adding event listeners...');
    print("${_tabController.index}");
    // Tab Changed swiping to a new tab
    activeTabIndex = _tabController.index;
    print('index changed.');
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    userInfo = appState.userInfo!;
    assetBalances = appState.assetBalances;
    activeWallet = appState.activeWallet;
    nfts = appState.nfts;
    wallets = userInfo.wallets!;

    claimedAssets = assetBalances[activeWallet!.publicKey]['claimed'];
    unclaimedAssets = assetBalances[activeWallet!.publicKey]['unclaimed'];
    // in order to make assets tab length dynamic we have to check
    // for when we have pending asset and then change the tablength
    // to 3 or back to 2 when we do not have pending assets.
    if (unclaimedAssets != null && unclaimedAssets.length > 0) {
      if (activeTabIndex == _tabController.length - 1) activeTabIndex = 2;
      tabLength = 3;
    } else {
      tabLength = 2;
      if (activeTabIndex > tabLength - 1) activeTabIndex = tabLength - 1;
    }

    if (tabLength != _tabController.length) {
      // change the length of tabController too or you will have an error
      _tabController = TabController(length: tabLength, vsync: this);
      _tabController.addListener(tabListener);
    }
    // keep track of the active tab to avoid having it changed
    // on each page rebuild
    _tabController.animateTo(activeTabIndex);
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
          TabBar(
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
              if (unclaimedAssets != null && tabLength == 3) ...[
                Tab(
                  height: 20,
                  text: '${LanguageEn.pending} (${unclaimedAssets.length})',
                ),
              ],
              Tab(
                height: 20,
                text: LanguageEn.nfts,
              ),
            ],
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
                                  alias: activeWallet!.alias!.capitalizeFirst!,
                                  totalBalance:
                                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, claimedAssets)} ${appState.defaultCurrency}',
                                  fiatBalance: appState.defaultCurrency == 'USD'
                                      ? null
                                      : '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, claimedAssets)} USD',
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
                                if (claimedAssets.length > 0) ...[
                                  for (var asset in claimedAssets) ...[
                                    GestureDetector(
                                        onTap: () {
                                          // since the original asset object
                                          // is immutable I create a new assetObj and
                                          // copy all the data into it so that
                                          // I'll be able to change the data
                                          appState.viewData![
                                              AssetDetailsViewPageConfig
                                                  .key] = {
                                            'assetCode': asset['assetCode'],
                                            'assetIssuer': asset['assetIssuer'],
                                            'amount': asset['amount'],
                                            'usdPrice': asset['usdPrice'],
                                            'qrCode': asset['qrCode'],
                                            'imageUrl': asset['imageUrl'],
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
                      if (tabLength == 3) ...[
                        Padding(
                          padding: const EdgeInsets.fromLTRB(0, 8, 0, 0),
                          child: Container(
                            child: SingleChildScrollView(
                              child: Column(
                                children: [
                                  if (unclaimedAssets.length > 0) ...[
                                    for (var asset in unclaimedAssets) ...[
                                      GestureDetector(
                                        onTap: () {
                                          setState(() {
                                            activeTabIndex =
                                                _tabController.index;
                                          });
                                          // since the original asset object
                                          // is immutable I create a new assetObj and
                                          // copy all the data into it so that
                                          // I'll be able to change the data
                                          appState.viewData![
                                              PendingAssetDetailsViewPageConfig
                                                  .key] = {
                                            'assetCode': asset['assetCode'],
                                            'assetIssuer': asset['assetIssuer'],
                                            'amount': asset['amount'],
                                            'qrCode': asset['qrCode'],
                                            'imageUrl': asset['imageUrl'],
                                          };
                                          print(appState.viewData);
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
                      Container(
                        height: height / 2,
                        child: SingleChildScrollView(
                          child: Column(
                            children: [
                              // in situations where the blockchain has an issue,
                              // some values can be returned as null or empty
                              // so always null check for such situations
                              // if (nfts != null && nfts != {}) ...[
                              //   if (nfts[activeWallet!.publicKey] != null &&
                              //       nfts[activeWallet!.publicKey].length >
                              //           0) ...[
                              gridView(),
                              //     SizedBox(height: 600),
                              //   ] else ...[
                              //     showEmptyNFTs(),
                              //   ]
                              // ] else ...[
                              //   showEmptyNFTs(),
                              // ],
                            ],
                          ),
                        ),
                      ),
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

  Widget tiles(asset) {
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
                  "https://drive.google.com/uc?export=view&id=103fw13pcBoCO2hkTPFX73BUKeWWkVpGZ",
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
                      asset["assetCode"].toString().isEmpty
                          ? 'XBN'
                          : asset["assetCode"],
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        "${getFiatRate(asset["usdPrice"], appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
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
                  getBalance(formatNumber(double.parse(asset["amount"]))),
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
                        '${calculateFiatValue(asset["amount"], asset["usdPrice"], appState.defaultCurrency, appState)} ${appState.defaultCurrency}'),
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
