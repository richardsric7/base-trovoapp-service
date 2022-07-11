import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
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

class Home extends StatefulWidget {
  const Home({Key? key}) : super(key: key);

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late TabController _tabController;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  List<Wallet>? wallets;
  String? activeWallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int touchedIndex = -1;

  @override
  void initState() {
    // TODO: implement initState
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
    wallets = userInfo.wallets!;
    nfts = appState.nfts;
    if (activeWallet == null && wallets!.length > 0) {
      activeWallet = wallets![0].publicKey;
    }
    claimedAssets = assetBalances[activeWallet]['claimed'];
    unclaimedAssets = assetBalances[activeWallet]['unclaimed'];
    if (unclaimedAssets.length > 0) {
      setState(() {
        tabLength = 3;
        _tabController = TabController(length: tabLength, vsync: this);
      });
    }
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 15,
              ),
              firstRow(),
              SizedBox(
                height: height / 50,
              ),
              walletSlides(wallets!),
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
      height: height / 1.9,
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
                                          appState.setActiveWallet = wallets!
                                              .firstWhere((wallet) =>
                                                  wallet.publicKey ==
                                                  activeWallet);
                                          appState.currentAction = PageAction(
                                            state: PageState.addPage,
                                            page: AssetDetailsViewPageConfig,
                                          );
                                        },
                                        child: tiles(asset),
                                      ),
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
                                        tiles(asset),
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
                                if (nfts != null && nfts != {}) ...[
                                  if (nfts[activeWallet] != null &&
                                      nfts[activeWallet].length > 0) ...[
                                    gridView(),
                                    SizedBox(height: 600),
                                  ] else ...[
                                    showEmptyNFTs(),
                                  ]
                                ] else ...[
                                  showEmptyNFTs(),
                                ],
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

  Widget firstRow() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Padding(
                  padding: EdgeInsets.fromLTRB(width / 18, 0, 0, 0),
                  child: SvgPicture.asset(
                    "assets/images/default.svg",
                    width: width / 6,
                  ),
                ),
                SizedBox(
                  width: width / 70,
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      LanguageEn.goodevening,
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getblck,
                        fontSize: 14.sp,
                        fontFamily: fontbody,
                      ),
                    ),
                    SizedBox(
                      height: 5,
                    ),
                    Text(
                      userInfo.firstName!.capitalizeFirst!,
                      style: TextStyle(
                        color: notifier.getbluecolor,
                        fontSize: 17.sp,
                        fontWeight: FontWeight.w600,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                )
              ],
            ),
          ],
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            GestureDetector(
              onTap: () {
                appState.currentAction = PageAction(
                    state: PageState.addPage, page: SearchViewPageConfig);
              },
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 8.0, horizontal: 10.0),
                child: SvgPicture.asset(
                  "assets/images/search.svg",
                  color: notifier.getbluecolor,
                  height: height / 40,
                ),
              ),
            ),
            GestureDetector(
              onTap: () {
                appState.currentAction = PageAction(
                    state: PageState.addPage, page: QrScannerPageConfig);
              },
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 8.0, horizontal: 10.0),
                child: SvgPicture.asset(
                  "assets/images/scan.svg",
                  color: notifier.getbluecolor,
                  height: height / 40,
                ),
              ),
            ),
            GestureDetector(
              onTap: () {
                appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: NotificationsViewPageConfig);
              },
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 8.0, horizontal: 10.0),
                // child: Image.asset("assets/images/notifications.png",
                //     color: notifier.getbluecolor),
                child: SvgPicture.asset(
                  "assets/images/notifications-active.svg",
                  color: notifier.getbluecolor,
                  height: height / 40,
                ),
              ),
            ),
            SizedBox(
              width: height / 50,
            ),
          ],
        )
      ],
    );
  }

  Widget gridView() {
    return Container(
      height: height / 2,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
        child: GridView(
          padding: const EdgeInsets.fromLTRB(0, 0, 0, 70),
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

  Widget walletSlides(List<Wallet> wallets) {
    // var thiswallets = [1, 2, 3];
    var colors = [notifier.getbluecolor, Colors.red, Colors.green];
    return CarouselSlider(
      options: CarouselOptions(
        onPageChanged: ((index, reason) => {
              setState(
                () => {
                  activeWallet = wallets[index].publicKey,
                  claimedAssets = assetBalances[activeWallet]['claimed'],
                  unclaimedAssets = assetBalances[activeWallet]['unclaimed'],
                },
              )
            }),
        height: height / 4.7,
        padEnds: false,
        enableInfiniteScroll: false,
        clipBehavior: Clip.antiAlias,
        viewportFraction: wallets.length > 1 ? 0.9 : 1,
      ),
      items: wallets.map((wallet) {
        return Builder(
          builder: (BuildContext context) {
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
                        padding: const EdgeInsets.symmetric(
                            vertical: 35.0, horizontal: 20),
                        child: Image.asset('assets/images/trovo_white.png'),
                      ),
                    ],
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 20.0, vertical: 35.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          wallet.alias!.capitalizeFirst!,
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
                          '2,082,898 NGN',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                            color: notifier.getwihitecolor,
                            fontFamily: fontsemibold,
                          ),
                        ),
                        SizedBox(height: 2),
                        Text(
                          '4,014 USD',
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
          },
        );
      }).toList(),
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
                  asset["amount"],
                  style: TextStyle(
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getblck,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                  child: Text(
                    '146,875 NGN',
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
