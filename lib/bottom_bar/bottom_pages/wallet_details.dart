import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
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

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
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
      tabLength = 3;
    } else {
      tabLength = 2;
    }
    // change the length of tabController too or you will have an error
    _tabController = TabController(length: tabLength, vsync: this);
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
              walletSlides(),
              SizedBox(
                height: height / 30,
              ),
              assetsTabs(),
            ],
          ),
        ),
      ),
    );
  }

  Widget assetsTabs() {
    return Container(
      height: height / 1.5,
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          DefaultTabController(
            length: tabLength,
            child: Container(
              height: height / 10,
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.only(
                  topRight: Radius.circular(40.0),
                  topLeft: Radius.circular(40.0),
                ),
                color: notifier.getbluecolor,
              ),
              // color: Colors.red,
              child: Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 12.0, 0, 10.0),
                    child: TabBar(
                      controller: _tabController,
                      indicatorColor: Colors.transparent,
                      labelStyle: TextStyle(
                        color: notifier.getwihitecolor,
                        fontSize: 12.sp,
                        fontWeight: FontWeight.w600,
                        fontFamily: fontbody,
                      ),
                      tabs: [
                        Tab(
                          height: 20,
                          text: LanguageEn.assets,
                        ),
                        if (tabLength == 3) ...[
                          Tab(
                            height: 20,
                            text:
                                '${LanguageEn.pendingassets} (${unclaimedAssets.length})',
                          ),
                        ],
                        Tab(
                          height: 20,
                          text: LanguageEn.nfts,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          Positioned(
            child: Column(
              children: [
                SizedBox(
                  height: height / 22,
                ),
                Expanded(
                  child: TabBarView(
                    controller: _tabController,
                    children: [
                      Container(
                        decoration: BoxDecoration(
                          borderRadius: const BorderRadius.only(
                            topRight: Radius.circular(40.0),
                            topLeft: Radius.circular(40.0),
                          ),
                          color: notifier.getwihitecolor,
                        ),
                        child: Padding(
                          padding: const EdgeInsets.fromLTRB(0, 28.0, 0, 0),
                          child: Container(
                            child: SingleChildScrollView(
                              child: Column(
                                children: [
                                  if (claimedAssets.length > 0) ...[
                                    for (var asset in claimedAssets) ...[
                                      GestureDetector(
                                          onTap: () {
                                            appState.viewData = {
                                              // since the original asset object
                                              // is immutable I create a new assetObj and
                                              // copy all the data into it so that
                                              // I'll be able to change the data
                                              AssetDetailsViewPageConfig.key: {
                                                'assetCode': asset['assetCode'],
                                                'assetIssuer':
                                                    asset['assetIssuer'],
                                                'amount': asset['amount'],
                                                'qrCode': asset['qrCode'],
                                                'imageUrl': asset['imageUrl'],
                                              }
                                            };
                                            print(appState.viewData);
                                            appState.currentAction = PageAction(
                                              state: PageState.addPage,
                                              page: AssetDetailsViewPageConfig,
                                            );
                                          },
                                          child: tiles(asset)),
                                    ],
                                    SizedBox(
                                      height: height / 20,
                                    ),
                                    Button(
                                      LanguageEn.dashboard,
                                      notifier.getbluecolor,
                                      notifier.getwihitecolor,
                                      onTap: () {
                                        appState.currentAction = PageAction(
                                            state: PageState.replaceAll,
                                            page: BottomHomePageConfig);
                                      },
                                    ),
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
                                ],
                              ),
                            ),
                          ),
                        ),
                      ),
                      if (tabLength == 3) ...[
                        Container(
                          decoration: BoxDecoration(
                            borderRadius: const BorderRadius.only(
                              topRight: Radius.circular(40.0),
                              topLeft: Radius.circular(40.0),
                            ),
                            color: notifier.getwihitecolor,
                          ),
                          child: Padding(
                            padding: const EdgeInsets.fromLTRB(0, 28.0, 0, 0),
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
                                            appState.viewData = {
                                              // since the original asset object
                                              // is immutable I create a new assetObj and
                                              // copy all the data into it so that
                                              // I'll be able to change the data
                                              PendingAssetDetailsViewPageConfig
                                                  .key: {
                                                'assetCode': asset['assetCode'],
                                                'assetIssuer':
                                                    asset['assetIssuer'],
                                                'amount': asset['amount'],
                                                'qrCode': asset['qrCode'],
                                                'imageUrl': asset['imageUrl'],
                                              }
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
                        ),
                      ],
                      Container(
                        height: height / 2,
                        decoration: BoxDecoration(
                          borderRadius: const BorderRadius.only(
                            topRight: Radius.circular(40.0),
                            topLeft: Radius.circular(40.0),
                          ),
                          color: notifier.getwihitecolor,
                        ),
                        child: Padding(
                          padding: const EdgeInsets.fromLTRB(0, 0.0, 0, 0),
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
      height: height / 1.6,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(0, 28.0, 0, 0),
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

  Widget walletSlides() {
    var colors = [notifier.getbluecolor, Colors.red, Colors.green];
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: colors[0],
          // color: colors[i - 1],
        ),
        child: Stack(children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
                child: Image.asset('assets/images/trovo_white.png'),
              ),
            ],
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  activeWallet!.alias!.capitalizeFirst!,
                  style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    Text(
                      LanguageEn.totalbalance,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getwihitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 98.0,
                ),
                Text(
                  appState.hideBalances ? hideBalanceText : '2,082,898 NGN',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: notifier.getwihitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
                SizedBox(height: 2),
                Text(
                  appState.hideBalances ? hideBalanceText : '4,014 USD',
                  style: TextStyle(
                    fontWeight: FontWeight.w300,
                    fontSize: 13,
                    color: notifier.getwihitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ],
            ),
          ),
        ]),
      ),
    );
  }

  Widget tiles(asset) {
    return Card(
      shadowColor: notifier.getblck,
      color: notifier.getwihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      elevation: 5,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
            title: Row(
              children: [
                Image.network(
                  "https://drive.google.com/uc?export=view&id=103fw13pcBoCO2hkTPFX73BUKeWWkVpGZ",
                  height: 35,
                  width: 35,
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
                        '25 NGN',
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
                  appState.hideBalances ? hideBalanceText : asset["amount"],
                  style: TextStyle(
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getblck,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                  child: Text(
                    appState.hideBalances ? hideBalanceText : '146,875 NGN',
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
}
