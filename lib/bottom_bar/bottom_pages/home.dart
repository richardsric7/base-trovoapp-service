import 'package:carousel_slider/carousel_slider.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart' hide Trans;
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/bottom_tab_page.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  late List<Wallet> carouselWallets;
  late List<Wallet> wallets;
  List<Asset>? unclaimedAssets;
  List<Asset>? claimedAssets;
  String? activeWallet;
  int tabLength = 2;
  int activeTabIndex = 0;
  int activeWalletIndex = 0;
  var noOfTransactionsToSign;
  var noXbnBalance = false;
  late Asset gas;
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;

  Map<String, DashboardAssetListMode> listModes = {
    'Tokenized Assets': DashboardAssetListMode.TokenizedAssets,
    'Other Assets': DashboardAssetListMode.OtherAssets,
  };

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    _tabController.addListener(tabListener);
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    appState.filterQuery = "&transactionStatus=PENDING";
    appState.getApprovals();
    gas = appState.primaryWallet.claimedAssets!
        .where((asset) => asset.assetCode == '')
        .first;
  }

  void tabListener() {
    // Tab Changed swiping to a new tab
    activeTabIndex = _tabController.index;
    setState(() {});
  }

  List<DropdownMenuItem<DashboardAssetListMode>> get getItems {
    List<DropdownMenuItem<DashboardAssetListMode>> items = [];
    listModes.forEach((key, value) {
      items.add(DropdownMenuItem(
          child: Text(
            key,
            overflow: TextOverflow.ellipsis,
          ),
          value: value));
    });
    return items;
  }

  @override
  void dispose() {
    _tabController.dispose();
    StoreData().storeInsertData('assetOrderings', appState.assetOrderings);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    userInfo = appState.userInfo!;
    wallets = userInfo.wallets!;
    carouselWallets =
        wallets.length > 6 ? wallets.getRange(0, 6).toList() : wallets;

    if ((activeWallet == null && wallets.length > 0) || noXbnBalance) {
      activeWallet = wallets[0].publicKey;
      claimedAssets = wallets[0].claimedAssets;
      unclaimedAssets = wallets[0].unClaimedAssets;
      noXbnBalance = claimedAssets!
              .firstWhere((asset) =>
                  asset.assetCode!.isEmpty && asset.assetIssuer!.isEmpty)
              .amount ==
          0;

      reOrderClaimedAssets(activeWallet!);
    }

    if (listMode == DashboardAssetListMode.OtherAssets) {
      // in order to make assets tab length dynamic we have to check
      // for when we have pending asset and then change the tablength
      // to 3 or back to 2 when we do not have pending assets.
      if (unclaimedAssets != null && unclaimedAssets!.length > 0) {
        tabLength = 2;
      } else {
        tabLength = 1;
      }
    } else if (listMode == DashboardAssetListMode.TokenizedAssets) {
      tabLength = 2;
    }

    if (tabLength != _tabController.length) {
      // change the length of tabController too or you will have an error
      _tabController = TabController(length: tabLength, vsync: this);
      _tabController.addListener(tabListener);
    }

    return Scaffold(
      key: key,
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      drawer: getDrawer(context, appState, notifier),
      body: SmartRefresher(
        enablePullDown: true,
        controller: _refreshController,
        onRefresh: refreshData,
        child: SingleChildScrollView(
          child: Stack(
            children: [
              Column(
                children: [
                  SizedBox(
                    height: 45,
                  ),
                  firstRow(),
                  if (userInfo.hasSecurityQuestions == 0) ...[
                    SizedBox(
                      height: 5,
                    ),
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
                  ],
                  walletSlides(wallets),
                  SizedBox(
                    height: height / 50,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Asset Mode',
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluecolor,
                          ),
                        ),
                        Container(
                          width: width / 2,
                          child: dropdown(
                            (value) {
                              setState(() {
                                listMode = value as DashboardAssetListMode;
                              });
                            },
                            getItems,
                            null,
                            'Tokenized Assets',
                            context,
                            null,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  // check if the user's xbn balance is 0. This usually is the si-
                  // tuation when a new user signs up and has not funded their wallet
                  // yet
                  if (!noXbnBalance) ...[
                    if (listMode == DashboardAssetListMode.TokenizedAssets) ...[
                      // area of new screen
                      // //////////////////
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 12.0),
                        child: Column(
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Expanded(
                                  child: Container(
                                    height: height / 14,
                                    child: Card(
                                        shadowColor: Colors.black,
                                        shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(10.0),
                                        ),
                                        color: notifier.isDark
                                            ? notifier.getbluecolor90
                                            : notifier.getaddsubwalletgrey,
                                        child: TextButton(
                                          onPressed: () {
                                            appState.viewData = {
                                              'assetCode': '',
                                              'assetIssuer': '',
                                              'walletPublicKey': activeWallet,
                                            };

                                            appState.currentAction = PageAction(
                                                state: PageState.addPage,
                                                page:
                                                    AssetDetailsViewPageConfig);
                                          },
                                          child: Column(
                                            mainAxisAlignment:
                                                MainAxisAlignment.spaceAround,
                                            children: [
                                              Row(
                                                children: [
                                                  Icon(
                                                    Icons.local_gas_station,
                                                    color:
                                                        notifier.getbluecolor,
                                                    size: 20,
                                                  ),
                                                  Text(
                                                    "Gas",
                                                    style: TextStyle(
                                                      fontSize: 13,
                                                      fontWeight:
                                                          FontWeight.bold,
                                                      fontFamily: fontsemibold,
                                                      color:
                                                          notifier.getbluecolor,
                                                      overflow:
                                                          TextOverflow.visible,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                              Row(
                                                children: [
                                                  Text(
                                                    '${formatHistoryNumber(gas.amount!, 1000)} ${getAssetCode(gas.assetCode)}',
                                                    style: TextStyle(
                                                      fontSize: 10,
                                                      fontWeight:
                                                          FontWeight.bold,
                                                      fontFamily: fontsemibold,
                                                      color:
                                                          notifier.getbluecolor,
                                                      overflow:
                                                          TextOverflow.visible,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                            ],
                                          ),
                                        )),
                                  ),
                                ),
                                Expanded(
                                  child: Container(
                                    height: height / 14,
                                    child: Card(
                                        shadowColor: Colors.black,
                                        shape: RoundedRectangleBorder(
                                          borderRadius:
                                              BorderRadius.circular(10.0),
                                        ),
                                        color: notifier.isDark
                                            ? notifier.getbluecolor90
                                            : notifier.getaddsubwalletgrey,
                                        child: TextButton(
                                          onPressed: () {
                                            changeTabPage(appState,
                                                ButtomTabPage.Wallets.index);
                                          },
                                          child: Column(
                                            mainAxisAlignment:
                                                MainAxisAlignment.spaceAround,
                                            children: [
                                              Row(
                                                children: [
                                                  Card(
                                                    shadowColor: Colors.black,
                                                    shape:
                                                        RoundedRectangleBorder(
                                                      borderRadius:
                                                          BorderRadius.circular(
                                                              10.0),
                                                    ),
                                                    color:
                                                        notifier.getbluecolor90,
                                                    child: Icon(
                                                      Icons.arrow_outward,
                                                      color: notifier
                                                          .getwihitecolor,
                                                      size: 15,
                                                    ),
                                                  ),
                                                  Text(
                                                    "See Assets",
                                                    style: TextStyle(
                                                      fontSize: 13,
                                                      fontWeight:
                                                          FontWeight.bold,
                                                      fontFamily: fontsemibold,
                                                      color:
                                                          notifier.getbluecolor,
                                                      overflow:
                                                          TextOverflow.visible,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                              Row(
                                                children: [
                                                  Text(
                                                    "24 Assets",
                                                    style: TextStyle(
                                                      fontSize: 10,
                                                      fontWeight:
                                                          FontWeight.bold,
                                                      fontFamily: fontsemibold,
                                                      color:
                                                          notifier.getbluecolor,
                                                      overflow:
                                                          TextOverflow.visible,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                            ],
                                          ),
                                        )),
                                  ),
                                ),
                              ],
                            ),
                            SizedBox(height: height / 50),
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
                                      indicatorColor:
                                          notifier.getbluewhitecolor,
                                      labelStyle: TextStyle(
                                        fontSize: 14,
                                        fontWeight: FontWeight.w600,
                                        fontFamily: fontsemibold,
                                      ),
                                      tabs: [
                                        Tab(
                                          height: 20,
                                          text: "Primary Listing".tr(),
                                        ),
                                        Tab(
                                          height: 20,
                                          text: "Secondary Listing".tr(),
                                        ),
                                      ],
                                    ),
                                  ),
                                  listingTabs(),
                                ],
                              ),
                            ),
                            SizedBox(
                              height: height / 70,
                            ),
                          ],
                        ),
                      ),
                    ] else ...[
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
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  fontFamily: fontsemibold,
                                ),
                                tabs: [
                                  Tab(
                                    height: 20,
                                    text: "assets".tr(),
                                  ),
                                  if (unclaimedAssets != null &&
                                      tabLength == 2) ...[
                                    Tab(
                                      height: 20,
                                      text:
                                          '${"pending".tr()} (${unclaimedAssets == null ? 0 : unclaimedAssets!.length})',
                                    ),
                                  ],
                                  // Tab(
                                  //   height: 20,
                                  //   text: "nfts".tr(),
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
                    ],
                  ] else ...[
                    showFundWallet(),
                  ]
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  var listOfAssets = <Map<String, String>>[
    {"imageUrl": "", "assetName": "ATLANTIS 1", "assetClass": "Property"},
    {"imageUrl": "", "assetName": "Orchard Estate", "assetClass": "Property"},
    {"imageUrl": "", "assetName": "Beacon Homes", "assetClass": "Property"},
    {"imageUrl": "", "assetName": "Animal Farm", "assetClass": "Property"},
  ];

  Widget listingTabs() {
    return Container(
      height: height / 2.21,
      child: TabBarView(
        controller: _tabController,
        children: [
          SingleChildScrollView(
            child: Column(
              children: [
                for (var i = 0; i < listOfAssets.length; i++) ...[
                  GestureDetector(
                    onTap: () {
                      appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: TokenizedAssetDetailViewPageConfig,
                      );
                    },
                    child: assetTile(
                        listOfAssets[i]['imageUrl'] ?? '',
                        listOfAssets[i]['assetName'] ?? '',
                        'Property',
                        i % 2 == 0),
                  ),
                ],
                SizedBox(height: height / 20),
              ],
            ),
          ),
          SingleChildScrollView(
            child: Column(
              children: [
                for (var i = 0; i < listOfAssets.length; i++) ...[
                  GestureDetector(
                    onTap: () {
                      appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: TokenizedAssetDetailViewPageConfig,
                      );
                    },
                    child: assetTile(
                        listOfAssets[i]['imageUrl'] ?? '',
                        listOfAssets[i]['assetName'] ?? '',
                        'Property',
                        i % 2 == 0),
                  ),
                ],
                SizedBox(height: height / 20),
              ],
            ),
          ),
        ],
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
                        unclaimedAssets!.length > 0) ...[
                      // if assets is greater than 5 then show five assets
                      // and then add a button to view all in the wallet
                      // details view
                      for (var i = 0; i < unclaimedAssets!.length; i++) ...[
                        GestureDetector(
                          onTap: () {
                            appState.setActiveWallet = wallets.firstWhere(
                                (wallet) => wallet.publicKey == activeWallet);

                            appState.viewData = {
                              'assetCode': unclaimedAssets![i].assetCode,
                              'assetIssuer': unclaimedAssets![i].assetIssuer,
                              'walletPublicKey': activeWallet,
                            };
                            appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: PendingAssetDetailsViewPageConfig,
                            );
                          },
                          child: tiles(
                              unclaimedAssets![i], null, activeWalletIndex),
                        ),
                      ],
                    ] else ...[
                      Container(
                        height: height / 4,
                        child: Padding(
                            padding: const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
                            child: Center(
                              child: Text(
                                "nopendingassets".tr(),
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
          if (claimedAssets!.length > 0) ...[
            Container(
              height: height / 1.85,
              child: ReorderableListView(
                padding: EdgeInsets.fromLTRB(0, 0, 0, 30),
                onReorder: (oldIndex, newIndex) {
                  if (oldIndex < newIndex) {
                    newIndex -= 1;
                  }
                  final Asset item = claimedAssets!.removeAt(oldIndex);
                  claimedAssets!.insert(newIndex, item);
                  setState(() {});
                },
                children: [
                  for (var i = 0; i < claimedAssets!.length; i++) ...[
                    GestureDetector(
                      key: Key(i.toString()),
                      onTap: () {
                        appState.setActiveWallet = wallets.firstWhere(
                            (wallet) => wallet.publicKey == activeWallet);

                        appState.viewData = {
                          'assetCode': claimedAssets![i].assetCode,
                          'assetIssuer': claimedAssets![i].assetIssuer,
                          'walletPublicKey': activeWallet,
                        };
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: AssetDetailsViewPageConfig,
                        );
                      },
                      child: tiles(claimedAssets![i], i, activeWalletIndex),
                    ),
                  ],
                ],
              ),
            ),
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
                      "noassets".tr(),
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
              "noNFTs".tr(),
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
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    SizedBox(
                      width: width / 40,
                    ),
                    IconButton(
                      onPressed: () {
                        key.currentState!.openDrawer();
                      },
                      icon: Icon(
                        Icons.menu,
                        size: 35,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      width: width / 20,
                    ),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          "goodday".tr(),
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 14,
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
                            fontSize: 17,
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
                        state: PageState.addPage,
                        page: SharedAccessViewPageConfig);

                    if (noOfTransactionsToSign != null &&
                        noOfTransactionsToSign > 0) {
                      // take the user to the pending approvals tab on the shared access view
                      WidgetsBinding.instance.addPostFrameCallback((_) {
                        appState.sharedAccesstabController.animateTo(1,
                            duration: Duration(milliseconds: 500),
                            curve: Curves.easeInOut);
                      });
                    }
                  },
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                        vertical: 8.0, horizontal: 10.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.people_alt_outlined,
                          size: 25,
                          color: notifier.getbluewhitecolor,
                        ),
                        FutureBuilder<Map>(
                          future: appState.approvals,
                          builder: (context, snapshot) {
                            if (snapshot.connectionState ==
                                    ConnectionState.done &&
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
                    padding: const EdgeInsets.symmetric(
                        vertical: 8.0, horizontal: 10.0),
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
                    padding: const EdgeInsets.symmetric(
                        vertical: 8.0, horizontal: 10.0),
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
                if (appState.walletMode == "Testnet") ...[
                  Visibility(
                    visible: true,
                    child: Container(
                      color: Color(0xFFAA453E),
                      width: 18,
                      child: RotatedBox(
                        quarterTurns: 1,
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 7),
                              child: Text(
                                "testnet".tr(),
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                  color: wihitecolor,
                                  fontSize: 11,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ]
              ],
            ),
          ],
        ),
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

  Widget assetTile(String imageUrl, String name, String type, isSubscribed) {
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 5),
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
                      fontSize: 12,
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
          trailing: ElevatedButton(
            onPressed: () async {
              isSubscribed
                  ? showUnSubscribePopup(
                      context,
                      onDone: () {},
                    )
                  : showSubscribePopup(
                      context,
                      onDone: () {},
                      dropdownItems: getStandardWallets,
                    );
            },
            style: ButtonStyle(
              overlayColor:
                  MaterialStateProperty.all<Color>(notifier.getsplashgrey),
              backgroundColor:
                  MaterialStateProperty.all<Color>(notifier.getbluewhitecolor),
              side: MaterialStateProperty.all(
                BorderSide(
                    color: notifier.getbluewhitecolor,
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
              width: width / 4,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    isSubscribed ? 'Subscribed' : 'Subscribe',
                    style: TextStyle(
                        fontFamily: fontsemibold,
                        fontSize: 11,
                        color: notifier.getwihitecolor),
                  ),
                  Icon(
                      isSubscribed
                          ? Icons.check_circle
                          : Icons.add_circle_rounded,
                      size: 18,
                      color: notifier.getwihitecolor),
                ],
              ),
            ),
          ),
        ),
      ),
    );
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
                  activeWalletIndex = index == 5 ? index - 1 : index,
                  activeWallet = wallets[activeWalletIndex].publicKey,
                  claimedAssets = wallets[activeWalletIndex].claimedAssets,
                  unclaimedAssets = wallets[activeWalletIndex].unClaimedAssets,
                },
              ),
              reOrderClaimedAssets(activeWallet!),
            }),
        height: height / 5.9,
        padEnds: false,
        enableInfiniteScroll: false,
        clipBehavior: Clip.antiAlias,
        viewportFraction: wallets.length > 1 ? 0.9 : 1,
      ),
      items: carouselWallets.map((wallet) {
        var indexOfWallet = wallets.indexOf(wallet);
        return Builder(
          builder: (BuildContext context) {
            if (indexOfWallet < 5) {
              return GestureDetector(
                onTap: () {
                  appState.viewData = {
                    'walletPublicKey': wallets[indexOfWallet].publicKey
                  };
                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: WalletDetailsViewPageConfig);
                },
                child: WalletSlide(
                  backColor: colors[wallets.indexOf(wallet)],
                  foreColor: getColor(context, indexOfWallet),
                  alias: wallet.alias!.capitalizeFirst!,
                  totalBalance:
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallets[indexOfWallet].claimedAssets!)} ${appState.defaultCurrency}',
                  fiatBalance: appState.defaultCurrency == 'USD'
                      ? null
                      : '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallets[indexOfWallet].claimedAssets!)} USD',
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
                                  "taptoviewall".tr(),
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

  void reOrderClaimedAssets(String publicKey) {
    // order asset according to user preference
    if (appState.assetOrderings[publicKey] != null) {
      claimedAssets!.forEach((asset) => asset.userPreferredIndex =
          appState.assetOrderings[publicKey]![asset.assetCode] ?? 0);
      claimedAssets!
          .sort((a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex));
    }
  }

  Widget tiles(Asset asset, int? indexOfAsset, int indexOfWallet) {
    if (indexOfAsset != null) {
      if (appState.assetOrderings[wallets[indexOfWallet].publicKey!] == null) {
        appState.assetOrderings[wallets[indexOfWallet].publicKey!] = {
          asset.assetCode!: indexOfAsset
        };
      }
      appState.assetOrderings[wallets[indexOfWallet].publicKey!]![
          asset.assetCode!] = indexOfAsset;
    }
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
                      getAssetCode(asset.assetCode),
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        "${getFiatRate(asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
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
                  getBalance(formatHistoryNumber(asset.amount!, 99000000000),
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
                        '${calculateFiatValue(asset.amount.toString(), asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
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
      await appState.getApprovals();
      _refreshController.refreshCompleted();
      appState.updateListeners();
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
                    "yourwalletisready".tr(),
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
                    "butyoucannotuseityet".tr(),
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
                    "youcangetbantutokens".tr(),
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
          "requestfromuser".tr(),
          notifier.getbluecolor,
          wihitecolor,
          onTap: () {
            appState.viewData = {
              'assetCode': '',
              'assetIssuer': '',
              'walletPublicKey': activeWallet,
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
          "sendxbntoyourwallet".tr(),
          notifier.getbluecolor80,
          wihitecolor,
          onTap: () {
            Clipboard.setData(
              ClipboardData(
                text: appState.primaryWallet.publicKey!,
              ),
            );
            showSnackBar("publickey".tr(), context);
          },
        ),
        SizedBox(
          height: height / 50,
        ),
        ButtonOutlined(
          "buyfromp2p".tr(),
          notifier.getwihitecolor,
          notifier.getbluewhitecolor,
          onTap: () {
            popup(
              context,
              title: "comingsoon".tr(),
              message: "p2pwillbelaunchingsoon".tr(),
              bodyColor: notifier.getbluewhitecolor,
            );
            // _launchUrl();
          },
        ),
        SizedBox(
          height: height / 50,
        ),
      ],
    );
  }
}

enum DashboardAssetListMode { TokenizedAssets, OtherAssets }
