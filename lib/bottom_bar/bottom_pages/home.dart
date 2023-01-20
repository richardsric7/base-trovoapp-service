import 'dart:async';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/bottom_tab_page.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Home extends StatefulWidget {
  final void Function(int)? onButtonPressed;
  const Home({Key? key, this.onButtonPressed}) : super(key: key);

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
  int tabLength = 1;
  int activeTabIndex = 0;
  int activeWalletIndex = 0;
  var noOfTransactionsToSign;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    _tabController.addListener(tabListener);
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    appState.filterQuery = "&transactionStatus=PENDING";
    appState.getApprovals();
  }

  void tabListener() {
    // Tab Changed swiping to a new tab
    activeTabIndex = _tabController.index;
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
      // if (activeTabIndex == _tabController.length - 1) activeTabIndex = 1;
      tabLength = 2;
    } else {
      tabLength = 1;
      // if (activeTabIndex > tabLength - 1) activeTabIndex = tabLength - 1;
    }

    if (tabLength != _tabController.length) {
      // change the length of tabController too or you will have an error
      _tabController = TabController(length: tabLength, vsync: this);
      _tabController.addListener(tabListener);
    }

    // keep track of the active tab to avoid having it changed
    // on each page rebuild
    // _tabController.animateTo(activeTabIndex);

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
                  height: height / 20,
                ),
                firstRow(),
                SizedBox(
                  height: height / 70,
                ),
                if (userInfo.hasSecurityQuestions == 0) ...[
                  GestureDetector(
                    onTap: () {
                      var primaryWallet = appState.userInfo!.wallets!
                          .firstWhere((wallet) => wallet.primaryWallet == 1);
                      appState.viewData = {
                        SecurityQuestionsViewPageConfig.key: {
                          'signer': primaryWallet.signer,
                          'publicKey': primaryWallet.publicKey,
                          'secretKey': appState.secretKeys[0],
                          'username': appState.userInfo!.username,
                        }
                      };

                      appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: SecurityQuestionsViewPageConfig);
                    },
                    child: Container(
                      color: Colors.red[400],
                      child: Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: Row(
                          children: [
                            SizedBox(
                              width: width / 50,
                            ),
                            Expanded(
                              child: Text(
                                'You have not setup security questions yet. Tap to setup security questions.',
                                style: TextStyle(
                                  fontSize: 13,
                                  fontWeight: FontWeight.bold,
                                  fontFamily: fontsemibold,
                                  color: notifier.getwihitecolor,
                                ),
                              ),
                            ),
                            SizedBox(
                              width: width / 50,
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 70,
                  ),
                ],
                walletSlides(wallets!),
                SizedBox(
                  height: height / 30,
                ),
                // check if the user's xbn balance is 0. This usually is the si-
                // tuation when a new user signs up and has not funded their wallet
                // yet
                if (claimedAssets
                    .where((asset) =>
                        (asset['assetCode'].toString().isEmpty &&
                            asset['assetIssuer'].toString().isEmpty) &&
                        double.parse(asset['amount']) != 0)
                    .isNotEmpty) ...[
                  DefaultTabController(
                    length: tabLength,
                    child: Column(
                      children: [
                        Padding(
                          padding: EdgeInsets.symmetric(
                            horizontal: tabLength == 2 ? 0 : 100,
                          ),
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
                              if (unclaimedAssets != null &&
                                  tabLength == 2) ...[
                                Tab(
                                  height: 20,
                                  text:
                                      '${LanguageEn.pending} (${unclaimedAssets == null ? 0 : unclaimedAssets.length})',
                                ),
                              ],

                              // Tab(
                              //   height: 20,
                              //   text: LanguageEn.nfts,
                              // ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    height: height / 70,
                  ),
                  assetsTabs(),
                ] else ...[
                  showFundWallet(),
                ]
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget assetsTabs() {
    return Container(
      height: height / 1.85,
      child: TabBarView(
        controller: _tabController,
        children: [
          showTokenAssets(),
          if (tabLength == 2) ...[
            Container(
              child: SingleChildScrollView(
                child: Column(
                  children: [
                    if (unclaimedAssets != null &&
                        unclaimedAssets.length > 0) ...[
                      // if assets is greater than 5 then show five assets
                      // and then add a button to view all in the wallet
                      // details view
                      for (var i = 0; i < unclaimedAssets.length; i++) ...[
                        GestureDetector(
                          onTap: () {
                            appState.setActiveWallet = wallets!.firstWhere(
                                (wallet) => wallet.publicKey == activeWallet);
                            appState.viewData = {
                              // since the original asset object
                              // is immutable I create a new assetObj and
                              // copy all the data into it so that
                              // I'll be able to change the data
                              PendingAssetDetailsViewPageConfig.key: {
                                'assetCode': unclaimedAssets[i]['assetCode'],
                                'assetIssuer': unclaimedAssets[i]
                                    ['assetIssuer'],
                                'amount': unclaimedAssets[i]['amount'],
                                'qrCode': unclaimedAssets[i]['qrCode'],
                                'imageUrl': unclaimedAssets[i]['imageUrl'],
                              }
                            };
                            appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: PendingAssetDetailsViewPageConfig,
                            );
                          },
                          child: tiles(unclaimedAssets[i], activeWalletIndex),
                        ),
                      ],
                    ] else ...[
                      Container(
                        height: height / 4,
                        child: Padding(
                            padding: const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
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
          ],
          // SingleChildScrollView(
          //   child: Column(
          //     children: [
          //       // in situations where the blockchain has an issue,
          //       // some values can be returned as null or empty
          //       // so always null check for such situations
          //       if (nfts != null && nfts != {}) ...[
          //         if (nfts[activeWallet] != null &&
          //             nfts[activeWallet].length > 0) ...[
          //           gridView(),
          //           SizedBox(height: 600),
          //         ] else ...[
          //           // showEmptyNFTs(),
          //           gridView(),
          //         ]
          //       ] else ...[
          //         // showEmptyNFTs(),
          //         gridView(),
          //       ],
          //     ],
          //   ),
          // ),
        ],
      ),
    );
  }

  Widget showTokenAssets() {
    return SingleChildScrollView(
      child: Column(
        children: [
          if (claimedAssets.length > 0) ...[
            for (var i = 0; i < claimedAssets.length; i++) ...[
              GestureDetector(
                onTap: () {
                  appState.setActiveWallet = wallets!
                      .firstWhere((wallet) => wallet.publicKey == activeWallet);

                  appState.viewData = {
                    'assetCode': claimedAssets[i]['assetCode'],
                    'assetIssuer': claimedAssets[i]['assetIssuer'],
                    'walletPublicKey': activeWallet,
                  };
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: AssetDetailsViewPageConfig,
                  );
                },
                child: tiles(claimedAssets[i], activeWalletIndex),
              ),
            ],
            SizedBox(
              height: height / 22,
            ),
          ] else ...[
            Container(
              height: height / 3,
              child: Padding(
                  padding: const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
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
                    backgroundColor: notifier.getbluecolor70,
                    child: GestureDetector(
                      onTap: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: ProfileDetailsViewPageConfig);
                      },
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(100.0),
                        child: Image.network(
                          appState.userInfo!.imageThumbnailURL!,
                          width: width / 6.8,
                          // height: width / 10,
                          fit: BoxFit.fill,
                          errorBuilder: (context, error, stackTrace) {
                            return Image.asset(
                              'assets/images/trovo.png',
                              width: width / 9,
                            );
                          },
                        ),
                      ),
                    ),
                  ),
                ),
                SizedBox(
                  width: width / 70,
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      LanguageEn.goodday,
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
                        fontFamily: fontsemibold,
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
                    state: PageState.addPage, page: SharedAccessViewPageConfig);

                if (noOfTransactionsToSign > 0) {
                  // take the user to the pending approvals tab on the shared access view
                  WidgetsBinding.instance.addPostFrameCallback((_) {
                    appState.sharedAccesstabController.animateTo(1,
                        duration: Duration(milliseconds: 500),
                        curve: Curves.easeInOut);
                  });
                }
              },
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 8.0, horizontal: 10.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.people_alt_outlined,
                      size: 25.sp,
                      color: notifier.getbluewhitecolor,
                    ),
                    FutureBuilder<Map>(
                      future: appState.approvals,
                      builder: (context, snapshot) {
                        if (snapshot.connectionState == ConnectionState.done &&
                            snapshot.hasData) {
                          noOfTransactionsToSign =
                              appState.filterQuery.contains('PENDING')
                                  ? snapshot.data!['totalRecords'] ?? 0
                                  : 0;
                          if (noOfTransactionsToSign > 0) {
                            return Text(
                              '($noOfTransactionsToSign)',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.bold,
                                fontFamily: fontsemibold,
                                color: notifier.getbluewhitecolor,
                              ),
                            );
                          }
                        }

                        return Text(
                          '',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.bold,
                            fontFamily: fontsemibold,
                          ),
                        );
                      },
                    ),
                  ],
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
                appState.hasNewAnnouncement = false;
              },
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 8.0, horizontal: 10.0),
                child: SvgPicture.asset(
                  appState.hasNewAnnouncement
                      ? "assets/images/notifications-active.svg"
                      : "assets/images/notifications.svg",
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
                  activeWalletIndex = index == 6 ? index - 1 : index,
                  activeWallet = wallets[activeWalletIndex].publicKey,
                  claimedAssets = assetBalances[activeWallet]['claimed'],
                  unclaimedAssets = assetBalances[activeWallet]['unclaimed'],
                },
              )
            }),
        height: height / 5.4,
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
              return GestureDetector(
                onTap: () {
                  appState.setActiveWallet = wallets[indexOfWallet];
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: WalletDetailsViewPageConfig);
                },
                child: WalletSlide(
                  backColor: colors[wallets.indexOf(wallet)],
                  foreColor: getColor(context, indexOfWallet),
                  alias: wallet.alias!.capitalizeFirst!,
                  totalBalance:
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, claimedAssets)} ${appState.defaultCurrency}',
                  fiatBalance: appState.defaultCurrency == 'USD'
                      ? null
                      : '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, claimedAssets)} USD',
                  initialHiddenState: appState.hideWalletList[indexOfWallet],
                  onHiddenStateChanged: (state) => {
                    setState(
                      () => {
                        appState.hideWalletList[indexOfWallet] = state,
                        StoreData().storeInsertData(
                            'hideWalletList', appState.hideWalletList)
                      },
                    )
                  },
                ),
              );
            }

            return GestureDetector(
              onTap: () {
                changeTabPage(appState, ButtomTabPage.Wallets.index);
                setState(() {});
              },
              child: Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: colors[0],
                    // color: colors[wallets.indexOf(wallet)],
                  ),
                  child: Stack(
                    alignment: AlignmentDirectional.centerEnd,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 35.0, horizontal: 20),
                            child: Image.asset(
                              'assets/images/trovo_white.png',
                              width: 80,
                            ),
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
                    ],
                  ),
                ),
              ),
            );
          },
        );
      }).toList(),
    );
  }

  Widget tiles(asset, indexOfWallet) {
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
                  asset["imageUrl"],
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
                  getBalance(
                      formatHistoryNumber(
                          double.parse(asset["amount"]), 99000000000),
                      indexOfWallet),
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
                        '${calculateFiatValue(asset["amount"], asset["usdPrice"], appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                        indexOfWallet),
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

  String getBalance(String balance, indexOfWallet) {
    String text;
    if (appState.hideBalances) text = hideBalanceText;

    if (indexOfWallet < appState.hideWalletList.length &&
        appState.hideWalletList[indexOfWallet])
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

  Widget showFundWallet() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                children: [
                  Text(
                    'Your wallet is ready!',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  Text(
                    'But you cannot use it for any transaction just yet until it is activated with at least 10 Bantu tokens (XBN)',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 16,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    'You can get Bantu tokens (XBN) for your wallet in 3 easy ways',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                ],
              ),
            ),
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
        Button(
          'Request XBN from Trovo user',
          notifier.getbluecolor,
          wihitecolor,
          onTap: () {
            appState.viewData![RequestSpecificPaymentViewPageConfig.key] = {
              'assetCode': '',
              'assetIssuer': '',
              'publicKey': appState.primaryWallet.publicKey,
              'walletAlias': appState.primaryWallet.alias,
              'isSharedAccess': 0,
            };

            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: RequestSpecificPaymentViewPageConfig);
          },
        ),
        SizedBox(
          height: height / 50,
        ),
        ButtonOutlined(
          'Send XBN to your wallet',
          notifier.getbluecolor80,
          wihitecolor,
          onTap: () {
            Clipboard.setData(
              ClipboardData(
                text: appState.primaryWallet.publicKey,
              ),
            );
            showSnackBar('Public key', context);
          },
        ),
        SizedBox(
          height: height / 50,
        ),
        ButtonOutlined(
          'Buy XBN on TrovoP2P',
          notifier.getwihitecolor,
          notifier.getbluewhitecolor,
          onTap: () => _launchUrl(),
        ),
        SizedBox(
          height: height / 50,
        ),
      ],
    );
  }

  Future<void> _launchUrl() async {
    Uri uri = Uri.https(trovoP2pUrl, '/login');
    if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
      throw 'Could not launch $uri';
    }
  }
}
