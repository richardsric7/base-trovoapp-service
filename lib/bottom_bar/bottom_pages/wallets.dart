import 'dart:convert';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/home.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class Wallets extends StatefulWidget {
  const Wallets({Key? key}) : super(key: key);

  @override
  State<Wallets> createState() => _WalletsState();
}

enum WalletAction { import, createNew }

enum WalletView { listWallets, addSubWallet, confirmAddSubWallet }

class _WalletsState extends State<Wallets> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  WalletAction? action = WalletAction.createNew;
  final Authenticator _authenticator = Authenticator();
  bool isTileView = false;
  bool isImport = true;
  String? tag;
  String? description;
  String? secretKey;
  late List<Wallet> carouselWallets;
  late List<Wallet> wallets;
  String? activeWallet;
  late TabController _tabController;
  int activeWalletIndex = 0;
  late Asset gas;
  int isAssetIssuerWallet = 0;
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  String password = '';
  var noXbnBalance = false;
  late Account primaryWalletKeyPair;
  late Account newSubWalletKeyPair;
  late DataProvider appState;
  List<Asset>? unclaimedAssets;
  List<Asset>? claimedAssets;
  late UserInfo userInfo;
  final _formKey = GlobalKey<FormState>();
  final _formKey2 = GlobalKey<FormState>();
  int tabLength = 2;
  int activeTabIndex = 0;
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;
  late RefreshController _refreshController;
  String selectedWalletMode = "My wallets";
  List<String> walletListMode = [
    'My wallets',
    'Shared wallets',
  ];
  late List<WalletTileColor> colors;
  late List<String> walletTypes = [
    'Standard',
    'Minting/Asset Tokenization',
    'Market Making/Trade',
    'Bulk Payment'
  ];

  var listOfAssets = <Map<String, String>>[
    {
      "imageUrl": "",
      "assetName": "Animal Farm",
      "assetClass": "Agriculture",
      "balance": "2049"
    },
    {
      "imageUrl": "",
      "assetName": "Beacon Homes",
      "assetClass": "Property",
      "balance": "3250"
    },
    {
      "imageUrl": "",
      "assetName": "C-Vitals",
      "assetClass": "Health",
      "balance": "100"
    },
    {
      "imageUrl": "",
      "assetName": "Drinkfly",
      "assetClass": "Beverage",
      "balance": "4000"
    }
  ];

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

  Map<String, DashboardAssetListMode> listModes = {
    'Tokenized Assets': DashboardAssetListMode.TokenizedAssets,
    'Other Assets': DashboardAssetListMode.OtherAssets,
  };

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

  int selectedWalletType = 0;

  List<DropdownMenuItem<String>> get walletTypeDropdownItems {
    var dropdownItems = walletTypes
        .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  wallet,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
            value: walletTypes.indexOf(wallet).toString()))
        .toList();

    return dropdownItems;
  }

  List<DropdownMenuItem<String>> get accessModeDropdownItems {
    return walletListMode
        .map<DropdownMenuItem<String>>((item) => DropdownMenuItem(
            child: Text(
              item,
              overflow: TextOverflow.ellipsis,
            ),
            value: item))
        .toList();
  }

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    _tabController.addListener(tabListener);
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    if ((appState.returnView != null && appState.returnView!.pages != null) &&
        appState.returnView!.pages!.contains(WalletPreparationViewPageConfig)) {
      appState.walletView.actionIcon = Icons.cancel_outlined;
      appState.walletView.actionText = "cancel".tr();
      appState.walletView.view = WalletView.addSubWallet;
    }
    gas = appState.primaryWallet.claimedAssets!
        .where((asset) => asset.assetCode == '')
        .first;
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
    userInfo = appState.userInfo!;
    wallets = userInfo.wallets!;
    carouselWallets =
        wallets.length > 6 ? wallets.getRange(0, 6).toList() : wallets;

    // if ((activeWallet == null && wallets.length > 0)) {
    //   activeWallet = wallets[0].publicKey;
    // }

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
        appBar: CustomAppBarWithoutLeading(
          context,
          notifier.getwihitecolor,
          height: height / 15,
          scaffoldKey: key,
          showMenu: true,
          txt: "wallets".tr(),
          titlecolor: notifier.getbluewhitecolor,
        ).getBar(),
        body: SmartRefresher(
          enablePullDown: true,
          controller: _refreshController,
          onRefresh: refreshData,
          child: ListView(
            children: [
              if (appState.walletView.view == WalletView.listWallets) ...[
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: Container(
                      color: notifier.getfavorites,
                      padding: EdgeInsets.all(8),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          Expanded(
                            child: DropdownButtonFormField(
                              isExpanded: true,
                              dropdownColor: notifier.isDark
                                  ? darktilewhitecolor
                                  : notifier.getaddsubwalletgrey,
                              decoration: InputDecoration(
                                contentPadding: EdgeInsets.symmetric(
                                    vertical: 0, horizontal: 20),
                                enabledBorder: OutlineInputBorder(
                                  borderSide: BorderSide.none,
                                  borderRadius: BorderRadius.circular(20),
                                ),
                                border: OutlineInputBorder(
                                  borderSide: BorderSide.none,
                                  borderRadius: BorderRadius.circular(20),
                                ),
                                filled: true,
                                fillColor: notifier.isDark
                                    ? darktilewhitecolor
                                    : notifier.getaddsubwalletgrey,
                              ),
                              value: selectedWalletMode,
                              icon: Icon(
                                Icons.keyboard_arrow_down_rounded,
                                color: notifier.getbluewhitecolor,
                              ),
                              elevation: 0,
                              style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 15,
                                  fontFamily: fontsemibold,
                                  fontWeight: FontWeight.w500),
                              onChanged: (newValue) {
                                setState(() {
                                  selectedWalletMode = newValue.toString();
                                });
                              },
                              items: accessModeDropdownItems,
                            ),
                          ),
                          SizedBox(
                            width: width / 15,
                          ),
                          Expanded(
                            child: Container(
                              child: Card(
                                  shadowColor: Colors.black,
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20.0),
                                  ),
                                  color: notifier.isDark
                                      ? notifier.getbluecolor90
                                      : notifier.getaddsubwalletgrey,
                                  child: TextButton(
                                    onPressed: () {
                                      appState.viewData = {
                                        'assetCode': '',
                                        'assetIssuer': '',
                                        'walletPublicKey':
                                            appState.primaryWallet.publicKey,
                                      };
                                      appState.currentAction = PageAction(
                                          state: PageState.addPage,
                                          page: AssetDetailsViewPageConfig);
                                    },
                                    child: Column(
                                      mainAxisAlignment:
                                          MainAxisAlignment.spaceAround,
                                      children: [
                                        Row(
                                          mainAxisAlignment:
                                              MainAxisAlignment.spaceAround,
                                          children: [
                                            Icon(
                                              Icons.local_gas_station,
                                              color: notifier.getbluecolor,
                                            ),
                                            Container(
                                              width: 100,
                                              child: Text(
                                                '${formatHistoryNumber(gas.amount!, 1000)} ${getAssetCode(gas.assetCode)}',
                                                overflow: TextOverflow.visible,
                                                style: TextStyle(
                                                  fontSize: 13,
                                                  fontWeight: FontWeight.bold,
                                                  fontFamily: fontsemibold,
                                                  color: notifier.getbluecolor,
                                                  overflow:
                                                      TextOverflow.visible,
                                                ),
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
                      )),
                )
              ] else ...[
                SizedBox(
                  height: height / 40,
                )
              ],
              walletSlides(wallets),
              GestureDetector(
                onTap: () {
                  if (appState.walletView.view == WalletView.listWallets) {
                    appState.walletView.actionIcon = Icons.cancel_outlined;
                    appState.walletView.actionText = "cancel".tr();
                    appState.walletView.view = WalletView.addSubWallet;
                  } else if (appState.walletView.view ==
                      WalletView.addSubWallet) {
                    appState.walletView.actionIcon =
                        Icons.add_circle_outline_sharp;
                    appState.walletView.actionText = "addsubwallet".tr();
                    appState.walletView.view = WalletView.listWallets;
                    resetForm();
                  } else if (appState.walletView.view ==
                      WalletView.confirmAddSubWallet) {
                    appState.walletView.actionIcon = Icons.cancel_outlined;
                    appState.walletView.actionText = "cancel".tr();
                    appState.walletView.view = WalletView.addSubWallet;
                  }

                  setState(() {});
                },
                child: Column(
                  children: [
                    Icon(
                      appState.walletView.actionIcon,
                      color: notifier.getbluewhitecolor,
                      size: 30,
                    ),
                    SizedBox(
                      height: 5,
                    ),
                    Text(
                      appState.walletView.actionText,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 90,
              ),
              if (appState.walletView.view == WalletView.addSubWallet) ...[
                addSubwallet()
              ] else if (appState.walletView.view ==
                  WalletView.confirmAddSubWallet) ...[
                confirmAddSubwallet()
              ] else ...[
                // isTileView ? gridView() : walletListView(),
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
                listMode == DashboardAssetListMode.TokenizedAssets
                    ? showAssets()
                    : cryptoAssets(),
              ],
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ));
  }

  Widget cryptoAssets() {
    return Column(
      children: [
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
                    if (unclaimedAssets != null && tabLength == 2) ...[
                      Tab(
                        height: 20,
                        text:
                            '${"pending".tr()} (${unclaimedAssets == null ? 0 : unclaimedAssets!.length})',
                      ),
                    ],
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
    );
  }

  Widget assetTile(String imageUrl, String name, String type, String balance) {
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
                  Container(
                    width: width / 3.4,
                    child: Text(
                      name,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
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
          trailing: Container(
            width: width / 4,
            child: Text(
              balance,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 12,
                fontFamily: fontsemibold,
                color: notifier.getbluecolor,
                overflow: TextOverflow.visible,
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget tokenizedAssetTile(
      String imageUrl, String name, String type, isSubscribed) {
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

  Widget showAssets() {
    return Column(
      children: [
        SizedBox(
          height: height / 70,
        ),
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
        // Row(
        //   children: [
        //     Padding(
        //       padding: const EdgeInsets.symmetric(horizontal: 20),
        //       child: Text(
        //         'Sort by:',
        //         style: TextStyle(
        //           fontSize: 13,
        //           fontFamily: fontbody,
        //           color: notifier.getbluecolor,
        //         ),
        //       ),
        //     ),
        //     Text(
        //       'Asset Tokens',
        //       style: TextStyle(
        //         fontSize: 13,
        //         fontFamily: fontsemibold,
        //         color: notifier.getbluecolor,
        //       ),
        //     ),
        //   ],
        // ),
        // SizedBox(
        //   height: height / 90,
        // ),
      ],
    );
  }

  Widget listingTabs() {
    return Container(
      height: height / 2.4,
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
                    child: tokenizedAssetTile(
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
                    child: tokenizedAssetTile(
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
                assetCount: wallet.claimedAssets?.length.toString(),
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
          },
        );
      }).toList(),
    );
  }

  Widget assetsTabs() {
    return Container(
      height: height / 2.4,
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
        ],
      ),
    );
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

  Widget showTokenAssets() {
    return SingleChildScrollView(
      child: Column(
        children: [
          if (claimedAssets!.length > 0) ...[
            Container(
              height: height / 2.4,
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

  void reOrderClaimedAssets(String publicKey) {
    // order asset according to user preference
    if (appState.assetOrderings[publicKey] != null) {
      claimedAssets!.forEach((asset) => asset.userPreferredIndex =
          appState.assetOrderings[publicKey]![asset.assetCode] ?? 0);
      claimedAssets!
          .sort((a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex));
    }
  }

  Widget gridView() {
    colors = [
      notifier.getstructuredbluecolor,
      notifier.getstructuredbluecolor90,
      notifier.getstructuredbluecolor80,
      notifier.getstructuredbluecolor70,
      notifier.getstructuredbluecolor60,
      notifier.getstructuredbluecolor50,
      notifier.getstructuredgreencolor,
      notifier.getstructuredgreencolor90,
      notifier.getstructuredgreencolor80,
      notifier.getstructuredgreencolor70,
      notifier.getstructuredgreencolor60,
      notifier.getstructuredgreencolor50,
      notifier.getorangecolor,
      notifier.getorangecolor90,
      notifier.getorangecolor80,
      notifier.getorangecolor70,
      notifier.getorangecolor60,
      notifier.getorangecolor50,
      notifier.getpinkcolor,
      notifier.getpinkcolor90,
      notifier.getpinkcolor80,
      notifier.getpinkcolor70,
      notifier.getpinkcolor60,
      notifier.getpinkcolor50,
    ];

    return Container(
      height: height / 2,
      child: GridView.count(
        primary: true,
        padding: const EdgeInsets.fromLTRB(15, 20, 15, 70),
        crossAxisCount: 2,
        mainAxisSpacing: 20,
        crossAxisSpacing: 20,
        childAspectRatio: 1.05,
        children: selectedWalletMode == 'My wallets'
            ? getWallets(userInfo.wallets!, true)
            : getSharedWallets(userInfo.sharedWallets!, true),
      ),
    );
  }

  Widget walletTile(
    walletName,
    usdBal,
    preferredFiatBal,
    WalletTileColor color, {
    isShared = false,
  }) {
    return Card(
      elevation: 5,
      shadowColor: Colors.black,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      color: color.backColor,
      child: Stack(
        children: [
          Center(
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
              child: Image.asset(
                'assets/images/trovo_white.png',
                height: 100,
                width: 100,
                color: Color(0x3CFFFFFF),
              ),
            ),
          ),
          Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Container(
                  width: width / 3,
                  child: Wrap(
                    alignment: WrapAlignment.center,
                    children: [
                      Text(
                        walletName,
                        textAlign: TextAlign.center,
                        overflow: TextOverflow.visible,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: color.foreColor,
                        ),
                      ),
                      if (isShared) ...[
                        SizedBox(width: width / 90),
                        Icon(
                          Icons.people_alt_outlined,
                          color: color.foreColor,
                          size: 20,
                        ),
                      ],
                    ],
                  ),
                ),
                Padding(
                  padding:
                      const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                  child: Text(
                    appState.hideBalances ? hideBalanceText : preferredFiatBal,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontbody,
                      color: color.foreColor,
                    ),
                  ),
                ),
                if (appState.defaultCurrency != 'USD') ...[
                  Padding(
                    padding:
                        const EdgeInsets.symmetric(vertical: 8, horizontal: 10),
                    child: Text(
                      appState.hideBalances ? hideBalanceText : usdBal,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontbody,
                        color: color.foreColor,
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget walletListView() {
    colors = [
      notifier.getstructuredbluecolor,
      notifier.getstructuredbluecolor90,
      notifier.getstructuredbluecolor80,
      notifier.getstructuredbluecolor70,
      notifier.getstructuredbluecolor60,
      notifier.getstructuredbluecolor50,
      notifier.getstructuredgreencolor,
      notifier.getstructuredgreencolor90,
      notifier.getstructuredgreencolor80,
      notifier.getstructuredgreencolor70,
      notifier.getstructuredgreencolor60,
      notifier.getstructuredgreencolor50,
      notifier.getorangecolor,
      notifier.getorangecolor90,
      notifier.getorangecolor80,
      notifier.getorangecolor70,
      notifier.getorangecolor60,
      notifier.getorangecolor50,
      notifier.getpinkcolor,
      notifier.getpinkcolor90,
      notifier.getpinkcolor80,
      notifier.getpinkcolor70,
      notifier.getpinkcolor60,
      notifier.getpinkcolor50,
    ];
    return Column(
      children: selectedWalletMode == 'My wallets'
          ? getWallets(userInfo.wallets!, false)
          : getSharedWallets(userInfo.sharedWallets!, false),
    );
  }

  Widget walletListItem(
    walletName,
    balanceUsd,
    preferredFiatBal,
    WalletTileColor color, {
    isShared = false,
  }) {
    return Container(
      // height: height / 6.6,
      margin: EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        borderRadius: const BorderRadius.all(Radius.circular(20.0)),
        color: color.backColor,
      ),
      child: Stack(
        alignment: AlignmentDirectional.centerEnd,
        children: [
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Image.asset(
                    'assets/images/trovo_white.png',
                    height: 80,
                    width: 80,
                  ),
                  SizedBox(
                    width: width / 20,
                  ),
                ],
              ),
            ],
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 2,
                  child: Wrap(
                    children: [
                      Text(
                        walletName,
                        style: TextStyle(
                          fontSize: 16,
                          color: color.foreColor,
                          fontFamily: fontsemibold,
                        ),
                      ),
                      if (isShared) ...[
                        SizedBox(width: width / 90),
                        Icon(
                          Icons.people_alt_outlined,
                          color: color.foreColor,
                          size: 20,
                        ),
                      ],
                    ],
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    Text(
                      appState.hideBalances
                          ? hideBalanceText
                          : preferredFiatBal,
                      style: TextStyle(
                        fontSize: 20,
                        color: color.foreColor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 80),
                if (appState.defaultCurrency != 'USD') ...[
                  Text(
                    appState.hideBalances ? hideBalanceText : balanceUsd,
                    style: TextStyle(
                      fontWeight: FontWeight.w300,
                      fontSize: 13,
                      color: color.foreColor,
                      fontFamily: fontbody,
                    ),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  // show the add subwallet view as if its a new page
  // the app's back button dispatcher has been overriden to make this page
  // behave as if is a new separate page when you press the back button
  Widget addSubwallet() {
    return Form(
      key: _formKey2,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Card(
              shadowColor: Colors.black,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(15.0),
              ),
              color: notifier.isDark
                  ? notifier.getbluecolor90
                  : notifier.getaddsubwalletgrey,
              child: Center(
                child: Column(
                  children: [
                    SizedBox(
                      height: height / 50,
                    ),
                    Container(
                      width: width / 1.4,
                      child: Text(
                        "abouttocreatesubwallet".tr(),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      "chooseamethod".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      children: [
                        SizedBox(
                          width: width / 10,
                        ),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.import,
                            groupValue: action,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor),
                            onChanged: (value) => {
                              setState(
                                () {
                                  action = value;
                                },
                              )
                            },
                          ),
                        ),
                        Text(
                          "importexistingwallet".tr(),
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    Row(
                      children: [
                        SizedBox(
                          width: width / 10,
                        ),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.createNew,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor),
                            groupValue: action,
                            onChanged: (value) => {
                              setState(
                                () {
                                  action = value;
                                },
                              )
                            },
                          ),
                        ),
                        Text(
                          "createnewwallet".tr(),
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Container(
                              width: width / 1.3,
                              decoration: BoxDecoration(
                                border: Border.all(
                                  color: notifier.getbluecolor,
                                ),
                                borderRadius: const BorderRadius.all(
                                    Radius.circular(15.0)),
                              ),
                              child: dropdown(
                                (newValue) async {
                                  selectedWalletType =
                                      int.parse(newValue.toString());
                                },
                                walletTypeDropdownItems,
                                selectedWalletType.toString(),
                                null,
                                context,
                                null,
                              ),
                            ),
                          ],
                        ),
                        // SizedBox(
                        //   height: height / 50,
                        // ),
                        // Row(
                        //   children: [
                        //     SizedBox(
                        //       width: width / 4.4,
                        //     ),
                        //     GestureDetector(
                        //       onTap: () {
                        //         mintWalletExplainerPopup(context);
                        //       },
                        //       child: Text(
                        //         'What does it mean?',
                        //         style: TextStyle(
                        //           decoration: TextDecoration.underline,
                        //           color: notifier.getbluewhitecolor,
                        //           fontSize: 12,
                        //           fontWeight: FontWeight.w500,
                        //           fontFamily: fontbody,
                        //         ),
                        //       ),
                        //     ),
                        //   ],
                        // ),
                        SizedBox(
                          height: height / 50,
                        ),
                      ],
                    )
                  ],
                ),
              ),
            ),
          ),
          SizedBox(
            height: height / 30,
          ),
          // Tag name
          CustomTextFormField.textField(
            "tag".tr(),
            notifier.getbluecolor,
            Icons.tag,
            notifier.getgrey,
            notifier.getbluewhitecolor,
            notifier.getblck,
            notifier.getgrey,
            70,
            300,
            initialValue: tag,
            onChanged: (value) {
              setState(() {
                tag = value.trim().replaceAll(' ', '');
              });
            },
            onSaved: (value) {
              print('tag: $value');
              tag = value.trim().replaceAll(' ', '');
            },
            keyboardtype: TextInputType.text,
            maxLength: 12,
            validator: validateTag,
            helperText: tag == null || tag!.isEmpty
                ? ''
                : "${appState.userInfo!.username}_$tag",
          ),
          SizedBox(height: height / 50),
          CustomTextFormField.textField(
            "description".tr(),
            notifier.getbluecolor,
            Icons.description,
            notifier.getgrey,
            notifier.getbluewhitecolor,
            notifier.getblck,
            notifier.getgrey,
            70,
            300,
            initialValue: description,
            onSaved: (value) {
              print('description: $value');
              description = value;
            },
            keyboardtype: TextInputType.text,
            maxLength: 100,
            validator: validateDescription,
          ),
          if (action == WalletAction.import) ...[
            SizedBox(height: height / 50),
            // Secret Key
            CustomPasswordFormField(
              "secretkey".tr(),
              notifier.getbluecolor,
              Icons.lock,
              notifier.getgrey,
              notifier.getbluewhitecolor,
              notifier.getblck,
              70,
              300,
              validator: (value) {
                var trimmedVal = value!.trim().replaceAll(' ', '');
                if (trimmedVal.isEmpty) {
                  return "entersecretkeyempty".tr();
                }

                if (trimmedVal.length < 56) {
                  return "secretkeyinvalid".tr();
                }

                try {
                  TrovoWalletSDK().parseSecretKey(value);
                } catch (e) {
                  return 'Secret Key is invalid';
                }

                return null;
              },
              onSaved: (value) {
                print('${"email".tr()}: $value');
                secretKey = value!.trim().replaceAll(' ', '');
              },
              maxLength: 56,
            ),
          ],
          SizedBox(height: height / 30),
          Button(
            "continuee".tr(),
            notifier.getbluecolor,
            wihitecolor,
            onTap: () => submitForm(),
          ),
          SizedBox(height: height / 20),
        ],
      ),
    );
  }

  // show the confirm add subwallet view as if its a new page
  // the app's back button dispatcher has been overriden to make this page
  // behave as if is a separate page when you press the back button
  Widget confirmAddSubwallet() {
    return Column(children: [
      Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20.0),
        child: Card(
          shadowColor: Colors.black,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15.0),
          ),
          color: notifier.isDark
              ? notifier.getbluecolor90
              : notifier.getaddsubwalletgrey,
          child: Center(
            child: Form(
              key: _formKey,
              child: Column(
                children: [
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.4,
                    child: Text(
                      "requesttocreatesubwallet".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    "tag".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Text(
                    "${userInfo.username!}_$tag",
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    "description".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Text(
                    description ?? '',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    "method".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Text(
                    action == WalletAction.import
                        ? "importsubwallet".tr()
                        : "createnewsubwallet".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    "wallettype".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Text(
                    walletTypes[selectedWalletType],
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    "publickey".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 15.0),
                    child: Text(
                      newSubWalletKeyPair.publicKey,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 20,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 15.0),
                    child: Text(
                      "willattractcharges".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 30,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
      SizedBox(height: height / 50),
      // Secret Key
      CustomPasswordFormField(
        "password".tr(),
        notifier.getbluecolor,
        Icons.lock,
        notifier.getgrey,
        notifier.getbluewhitecolor,
        notifier.getblck,
        70,
        300,
        validator: validatePassword,
        textInputAction: TextInputAction.done,
        onChanged: (value) {
          setState(() {
            password = value!.trim().replaceAll(' ', '');
          });
        },
        // onSubmitted: (value) {
        //   print('email: $value');
        //   secretKey = value!.trim().replaceAll(' ', '');
        // },
        onSaved: (value) {
          print('email: $value');
          secretKey = value!.trim().replaceAll(' ', '');
        },
      ),
      SizedBox(height: height / 30),

      if (appState.biometricEnabled && password.isEmpty) ...[
        Button(
          "authorizewithbiometrics".tr(),
          notifier.getbluecolor,
          wihitecolor,
          onTap: toggleSwitch,
        ),
      ] else ...[
        Button(
          "authorize".tr(),
          notifier.getbluecolor,
          wihitecolor,
          onTap: handleAuthorization,
        ),
      ],
      SizedBox(height: height / 20),
    ]);
  }

  List<Widget> getWallets(
    List<Wallet> wallets,
    isTileMode,
  ) {
    List<Wallet> filteredWallets = appState.userInfo!.mySolelyOwnedWallets;
    return [
      for (var i = 0; i < filteredWallets.length; i++) ...[
        GestureDetector(
            onTap: () {
              appState.viewData = {
                'walletPublicKey': filteredWallets[i].publicKey
              };
              appState.currentAction = PageAction(
                  state: PageState.addPage, page: WalletDetailsViewPageConfig);
            },
            child: isTileMode
                ? walletTile(
                    filteredWallets[i].alias!.capitalizeFirst!,
                    '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, filteredWallets[i].claimedAssets!)} USD',
                    '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, filteredWallets[i].claimedAssets!)} ${appState.defaultCurrency}',
                    i % 2 == 0
                        ? colors[((i + 1) % colors.length)]
                        : colors[((i) % colors.length)],
                  )
                : Column(children: [
                    walletListItem(
                      filteredWallets[i].alias!.capitalizeFirst!,
                      '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, filteredWallets[i].claimedAssets!)} USD',
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, filteredWallets[i].claimedAssets!)} ${appState.defaultCurrency}',
                      colors[((i + 1) % colors.length)],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                  ])),
      ]
      // shared wallets
    ];
  }

  List<Widget> getSharedWallets(List<Wallet> sharedWallets, isTileMode) {
    var walletTiles = <Widget>[];
    int i = 0;
    sharedWallets.forEach((wallet) {
      walletTiles.add(
        GestureDetector(
          onTap: () {
            appState.viewData = {'walletPublicKey': wallet.publicKey};
            appState.currentAction = PageAction(
                state: PageState.addPage, page: WalletDetailsViewPageConfig);
          },
          child: isTileMode
              ? walletTile(
                  wallet.alias.toString().capitalizeFirst!,
                  '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallet.claimedAssets!)} USD',
                  '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallet.claimedAssets!)} ${appState.defaultCurrency}',
                  i % 2 == 0
                      ? colors[((i + 1) % colors.length)]
                      : colors[((i) % colors.length)],
                  isShared: true,
                )
              : Column(
                  children: [
                    walletListItem(
                      wallet.alias.toString().capitalizeFirst!,
                      '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallet.claimedAssets!)} USD',
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallet.claimedAssets!)} ${appState.defaultCurrency}',
                      i % 2 == 0
                          ? colors[((i + 1) % colors.length)]
                          : colors[((i) % colors.length)],
                      isShared: true,
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                  ],
                ),
        ),
      );
      i++;
    });

    return walletTiles;
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return "pleaseenteryourpassword".tr();

    if (value.length < 6) return "use6charsormoreforpassword".tr();

    return null;
  }

  submitForm() async {
    var form = _formKey2.currentState;
    if (!form!.validate()) {
      return;
    }
    form.save();
    generateKeyPairs();
    setState(() {
      appState.walletView.view = WalletView.confirmAddSubWallet;
      appState.walletView.actionIcon = Icons.arrow_circle_left_outlined;
      appState.walletView.actionText = "back".tr();
      password = '';
      secretKey = '';
    });
  }

  generateKeyPairs() {
    primaryWalletKeyPair =
        TrovoWalletSDK().parseSecretKey(appState.secretKeys[0]);
    setState(() {
      if (action == WalletAction.import) {
        try {
          // parse supplied secret to get the keypair
          newSubWalletKeyPair = TrovoWalletSDK().parseSecretKey(secretKey);
        } catch (e) {
          popup(context, title: "error".tr(), message: "invalidsecretkey".tr());
        }
      } else {
        // generate keypair for the new subwallet
        newSubWalletKeyPair = TrovoWalletSDK().createAccount();
      }
    });
  }

  String? validateTag(String? value) {
    if (value!.isEmpty) return "enterwallettag".tr();

    String pattern = r'^[a-zA-Z0-9\_]*$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "invalidtagname".tr();
    }

    return null;
  }

  String? validateDescription(String? value) {
    if (value!.isEmpty) return "enterwalletdesc".tr();

    return null;
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer();
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  void handleAuthorization() {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer();
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  void sendDataToServer() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "publickey": newSubWalletKeyPair.publicKey,
        "walletTag": tag,
        "WalletDescription": description,
        // "assetIssuerWallet": isAssetIssuerWallet,
        "walletType": selectedWalletType,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      print('response: $responseData');

      if (responseData['statusCode'] == 200) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        await postProcessData(
            messageShown, messageLength, responseData['data']);
        // print('sending full data to server.........');
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }
    hideLoader(context);
  }

  postProcessData(messageShown, messageLength, data) {
    // we would like to display all messages returned from the initial
    // request to server using a popup. In order to achieve that we
    // employ the use of a little recursion here. Please recursive
    // functions can turn into a nightmare fast so be carefull here.
    if (messageShown <= messageLength - 1) {
      showResponseMessage(
          context,
          data['messages'][messageShown],
          () => {
                print('postProcessData: $messageShown'),
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    sendFullDataToServer(data);
    return;
  }

  void sendFullDataToServer(responseBody) async {
    showLoader(context);

    try {
      // get primary signature
      var primarySignature = TrovoWalletSDK().signBase64Txn(
        primaryWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      // get secondary signature
      var subWalletSignature = TrovoWalletSDK().signBase64Txn(
        newSubWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      responseBody['primarySignature'] = primarySignature;
      responseBody['subWalletSignature'] = subWalletSignature;

      String requestBody = jsonEncode(responseBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      if (responseData['statusCode'] == 200) {
        // add the secret key of this new subwallet to
        // the existing list of secrets
        appState.secretKeys.add(newSubWalletKeyPair.secretKey);
        // store back the list of secret keys but this time it
        // contains the secret key of the newly created subwallet
        await StoreData().storeInsertData('secretKey', appState.secretKeys);
        await updateUserInfo(
          appState.primaryWallet.signer,
          appState.secretKeys[0],
          appState.primaryWallet.publicKey,
          userInfo.username,
          appState,
          forceRefresh: true,
        );
        // add the new subwallet to appState and
        // set the newly created subwallet as the activeWallet
        appState.activeWallet = appState.userInfo!.wallets!.firstWhere(
            (wallet) => wallet.publicKey == newSubWalletKeyPair.publicKey);
        appState.activeWallet!.secretKey = newSubWalletKeyPair.secretKey;
        // move to next page
        appState.currentAction = PageAction(
            state: PageState.addPage, page: CongratulationsPageConfig);
        resetForm();
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['error']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }

  void resetForm() {
    tag = '';
    description = '';
    secretKey = '';
    appState.walletView.actionIcon = Icons.add_circle_outline_sharp;
    appState.walletView.actionText = "addsubwallet".tr();
    appState.walletView.view = WalletView.listWallets;
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
