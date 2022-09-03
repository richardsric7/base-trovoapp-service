import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
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

class Home extends StatefulWidget {
  final void Function(int) onButtonPressed;
  const Home({Key? key, required this.onButtonPressed}) : super(key: key);

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late TabController _tabController;
  late RefreshController _refreshController;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  List<Wallet>? wallets;
  List<Wallet>? carouselWallets;
  String? activeWallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int activeTabIndex = 0;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    _tabController.addListener(tabListener);
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    appState.resetActiveWalletBalances();
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
  void dispose() {
    _tabController.dispose();
    super.dispose();
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
    carouselWallets =
        wallets!.length > 7 ? wallets!.getRange(0, 7).toList() : wallets;
    nfts = appState.nfts;
    if (activeWallet == null && wallets!.length > 0) {
      activeWallet = wallets![0].publicKey;
    }
    claimedAssets = assetBalances[activeWallet]['claimed'];
    unclaimedAssets = assetBalances[activeWallet]['unclaimed'];

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

    print(
        'tablength: $tabLength, tabcontroller.length: ${_tabController.length}');
    print('activeTabIndex: $activeTabIndex');

    // keep track of the active tab to avoid having it changed
    // on each page rebuild
    _tabController.animateTo(activeTabIndex);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SmartRefresher(
          enablePullDown: true,
          controller: _refreshController,
          onRefresh: refreshData,
          child: SingleChildScrollView(
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
                            text:
                                '${LanguageEn.pending} (${unclaimedAssets.length})',
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
                                    if (unclaimedAssets != null &&
                                        unclaimedAssets.length > 0) ...[
                                      for (var asset in unclaimedAssets) ...[
                                        GestureDetector(
                                          onTap: () {
                                            appState.setActiveWallet = wallets!
                                                .firstWhere((wallet) =>
                                                    wallet.publicKey ==
                                                    activeWallet);
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
                                if (nfts != null && nfts != {}) ...[
                                  if (nfts[activeWallet] != null &&
                                      nfts[activeWallet].length > 0) ...[
                                    gridView(),
                                    SizedBox(height: 600),
                                  ] else ...[
                                    // showEmptyNFTs(),
                                    gridView(),
                                  ]
                                ] else ...[
                                  // showEmptyNFTs(),
                                  gridView(),
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
                    child: CircleAvatar(
                      radius: 30,
                      backgroundColor: notifier.getwihitecolor,
                      foregroundImage: AssetImage("assets/images/obi.png"),
                    )),
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
                        color: notifier.getbluewhitecolor,
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
                        color: notifier.getbluewhitecolor,
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
                  color: notifier.getbluewhitecolor,
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
                  color: notifier.getbluewhitecolor,
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
                  color: notifier.getbluewhitecolor,
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
    var colors = [
      notifier.getbluecolor,
      notifier.getbluecolor90,
      notifier.getbluecolor80,
      notifier.getbluecolor70,
      notifier.getbluecolor60,
      notifier.getbluecolor50,
    ];
    return CarouselSlider(
      options: CarouselOptions(
        onPageChanged: ((index, reason) => {
              setState(
                () => {
                  appState.resetActiveWalletBalances(),
                  activeWallet = wallets[index].publicKey,
                  claimedAssets = assetBalances[activeWallet]['claimed'],
                  unclaimedAssets = assetBalances[activeWallet]['unclaimed'],
                  print('activeWallet: $activeWallet'),
                },
              )
            }),
        height: height / 4.7,
        padEnds: false,
        enableInfiniteScroll: false,
        clipBehavior: Clip.antiAlias,
        viewportFraction: wallets.length > 1 ? 0.9 : 1,
      ),
      items: carouselWallets!.map((wallet) {
        var indexOfWallet = wallets.indexOf(wallet);
        return Builder(
          builder: (BuildContext context) {
            if (indexOfWallet < 6) {
              return WalletSlide(
                backColor: colors[wallets.indexOf(wallet)],
                foreColor: getColor(context, indexOfWallet),
                alias: wallet.alias!.capitalizeFirst!,
                totalBalance: '2,082,898 NGN',
                fiatBalance: '4,014 USD',
              );
            }

            return GestureDetector(
              onTap: () {
                // moves user to the wallets list tab
                widget.onButtonPressed(1);
              },
              child: Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: colors[0],
                    // color: colors[wallets.indexOf(wallet)],
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
                          horizontal: 20.0, vertical: 25.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.center,
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Row(
                            children: [
                              Text(
                                LanguageEn.taptoviewall,
                                style: TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.w600,
                                    color: wihitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                width: width / 50,
                              ),
                              Icon(
                                Icons.arrow_forward,
                                color: wihitecolor,
                              )
                            ],
                          ),
                        ],
                      ),
                    ),
                  ]),
                ),
              ),
            );
          },
        );
      }).toList(),
    );
  }

  Widget tiles(asset) {
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
                      getAssetCode(asset["assetCode"]),
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
                    getBalance('146,875 NGN'),
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
    print(
        'getting bal for active wallet... ${appState.hideActiveWalletBalance}');
    String text;
    if (appState.hideBalances) text = hideBalanceText;

    if (appState.hideActiveWalletBalance)
      text = hideBalanceText;
    else
      text = balance;

    return text;
  }

  void refreshData() async {
    try {
      await appState.refreshData();
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }
}
