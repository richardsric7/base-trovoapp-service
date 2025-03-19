import 'dart:convert';

import 'package:carousel_slider/carousel_slider.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/home.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/models/wallets_list_view_data.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Wallets extends StatefulWidget {
  const Wallets({Key? key}) : super(key: key);

  @override
  State<Wallets> createState() => _WalletsState();
}

class _WalletsState extends State<Wallets> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  bool isTileView = false;
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
  List<Asset> claimedAssets = [];
  late UserInfo userInfo;
  final carouselController = CarouselController();
  int tabLength = 2;
  int activeTabIndex = 0;
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;
  late RefreshController _refreshController;
  String selectedWalletMode = "My wallets";
  List<String> walletListMode = ['My wallets', 'Shared wallets', 'All wallets'];
  late List<WalletTileColor> colors;
  List<TokenizedAsset> tokenizedAssets = [];

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
    'Asset Tokens': DashboardAssetListMode.TokenizedAssets,
    'Other Tokens': DashboardAssetListMode.OtherAssets,
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

  List<DropdownMenuItem<String>> get walletListModeDropdownItems {
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
    wallets = switch (selectedWalletMode) {
          'My wallets' => userInfo.mySolelyOwnedWallets,
          'Shared wallets' => userInfo.sharedWallets,
          'All wallets' => userInfo.allWallets,
          _ => [],
        } ??
        [];
    carouselWallets =
        wallets.length > 6 ? wallets.getRange(0, 6).toList() : wallets;

    // if ((activeWallet == null && wallets.length > 0)) {
    //   activeWallet = wallets[0].publicKey;
    // }

    if ((activeWallet == null && wallets.length > 0) || noXbnBalance) {
      activeWallet = wallets[0].publicKey;
      claimedAssets = wallets[0]
          .claimedAssets!
          .where((asset) => asset.assetCode != '' && asset.assetIssuer != '')
          .toList();
      unclaimedAssets = wallets[0].unClaimedAssets;
      tokenizedAssets = wallets[0].tokenizedAssets ?? [];
      noXbnBalance = wallets[0]
              .claimedAssets!
              .firstWhere((asset) =>
                  asset.assetCode!.isEmpty && asset.assetIssuer!.isEmpty)
              .amount ==
          0;

      if (tokenizedAssets.isEmpty) {
        listMode = DashboardAssetListMode.OtherAssets;
      }

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
      tabLength = 1;
    }

    if (tabLength != _tabController.length) {
      // change the length of tabController too or you will have an error
      _tabController = TabController(length: tabLength, vsync: this);
      _tabController.addListener(tabListener);
    }

    return Scaffold(
        key: key,
        resizeToAvoidBottomInset: false,
        floatingActionButton: FloatingActionButton(
          onPressed: () {
            addSubWalletPopup(context);
            setState(() {});
          },
          backgroundColor: notifier.getbluecolor,
          child: Icon(
            appState.walletView.actionIcon,
            color: wihitecolor,
            size: 30,
          ),
        ),
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
                                borderRadius: BorderRadius.circular(10),
                              ),
                              border: OutlineInputBorder(
                                borderSide: BorderSide.none,
                                borderRadius: BorderRadius.circular(10),
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
                              if (wallets.length > 1)
                                carouselController.animateToPage(0);
                            },
                            items: walletListModeDropdownItems,
                          ),
                        ),
                        SizedBox(
                          width: width / 15,
                        ),
                        Expanded(
                          child: Container(
                            child: Card(
                                // shadowColor: Colors.black,
                                elevation: 0,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(10.0),
                                ),
                                color: notifier.isDark
                                    ? notifier.getbluecolor90
                                    : notifier.getaddsubwalletgrey,
                                child: TextButton(
                                  onPressed: () {
                                    appState.currentAction = PageAction(
                                      state: PageState.addPage,
                                      page: AllWalletsViewPageConfig,
                                    );
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
                                            Icons.list,
                                            color: notifier.getbluewhitecolor,
                                          ),
                                          Container(
                                            width: 100,
                                            child: Text(
                                              'viewall'.tr(),
                                              overflow: TextOverflow.visible,
                                              style: TextStyle(
                                                fontSize: 13,
                                                fontWeight: FontWeight.bold,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                                overflow: TextOverflow.visible,
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
              ),
              if (wallets.isNotEmpty) ...[
                walletSlides(wallets),
                SizedBox(
                  height: height / 70,
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: Container(
                    color: notifier.getfavorites,
                    // padding: EdgeInsets.all(5),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        Expanded(
                          child: Container(
                            child: Card(
                                // shadowColor: Colors.black,
                                elevation: 0,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(10.0),
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
                                            color: notifier.getbluewhitecolor,
                                          ),
                                          Container(
                                            width: 100,
                                            child: Text(
                                              '${"gas".tr()} ${formatHistoryNumber(wallets[activeWalletIndex].claimedAssets!.where((w) => w.assetCode == '').first.amount!, 1000, isShort: true)}',
                                              overflow: TextOverflow.visible,
                                              style: TextStyle(
                                                fontSize: 15,
                                                fontWeight: FontWeight.bold,
                                                fontFamily: fontsemibold,
                                                color:
                                                    notifier.getbluewhitecolor,
                                                overflow: TextOverflow.visible,
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
                        SizedBox(
                          width: width / 20,
                        ),
                        Expanded(
                          child: Container(
                            width: width / 2,
                            child: dropdown(
                              (value) {
                                setState(() {
                                  listMode = value as DashboardAssetListMode;
                                });
                              },
                              getItems,
                              null,
                              "othertokens".tr(),
                              context,
                              null,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
              SizedBox(
                height: height / 80,
              ),
              if (wallets.isNotEmpty) ...[
                listMode == DashboardAssetListMode.TokenizedAssets
                    ? showAssets()
                    : cryptoAssets(),
              ] else ...[
                Container(
                  height: height / 2,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        "nowalletshere".tr(),
                        style: TextStyle(
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                ),
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
                      text:
                          "othertokens".tr().toLowerCase().capitalizeEachWord(),
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

  Widget showAssets() {
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
                      text:
                          "assettokens".tr().toLowerCase().capitalizeEachWord(),
                    ),
                  ],
                ),
              ),
              listingTabs(),
            ],
          ),
        ),
      ],
    );
  }

  Widget listingTabs() {
    return Container(
      height: height / 2.19,
      child: TabBarView(
        controller: _tabController,
        children: [
          SingleChildScrollView(
            child: Column(
              children: [
                SizedBox(
                  height: height / 90,
                ),
                if (tokenizedAssets.isNotEmpty) ...[
                  for (var i = 0; i < tokenizedAssets.length; i++) ...[
                    GestureDetector(
                      onTap: () {
                        appState.viewData = {
                          'assetCode': tokenizedAssets[i].assetCode,
                          'assetIssuer': tokenizedAssets[i].assetIssuer,
                          'tokenizedAsset': tokenizedAssets[i],
                          'walletPublicKey': activeWallet,
                        };
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: AssetTokenDetailsViewPageConfig,
                        );
                      },
                      child: tokenizedAssetTile(
                          tokenizedAssets[i], i, activeWalletIndex),
                    ),
                  ],
                  SizedBox(height: height / 20),
                ] else ...[
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: SizedBox(
                      height: height / 4,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "nothingtoshowhere2".tr(),
                            textAlign: TextAlign.center,
                            style: TextStyle(
                                fontSize: 16,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody),
                          ),
                        ],
                      ),
                    ),
                  )
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget tokenizedAssetTile(
      TokenizedAsset asset, int? indexOfAsset, int indexOfWallet) {
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
                if (asset.assetLogo != null) ...[
                  Image.memory(
                    base64Decode(asset.assetLogo!),
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
                ] else ...[
                  Image.asset(
                    'assets/images/trovo.png',
                    height: 35,
                    width: 35,
                  ),
                ],
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
                        "${getFiatRate(asset.usdPrice == null ? '1.45' : asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
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
                          asset.subscriptionAmount ?? 0, 99000000000),
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
                        '${calculateFiatValue(asset.subscriptionAmount == null ? '0' : asset.subscriptionAmount.toString(), asset.usdPrice == null ? '1.45' : asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
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
      carouselController: carouselController,
      options: CarouselOptions(
        onPageChanged: ((index, reason) => {
              setState(
                () => {
                  activeWalletIndex = index == 5 ? index - 1 : index,
                  activeWallet = wallets[activeWalletIndex].publicKey,
                  claimedAssets = wallets[activeWalletIndex]
                      .claimedAssets!
                      .where((asset) =>
                          asset.assetCode != '' && asset.assetIssuer != '')
                      .toList(),
                  unclaimedAssets = wallets[activeWalletIndex].unClaimedAssets,
                  tokenizedAssets =
                      wallets[activeWalletIndex].tokenizedAssets ?? [],
                },
              ),
              // reOrderClaimedAssets(activeWallet!),
            }),
        height: height / 5.6,
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
                  isSharedWallet: wallet.isSharedWallet,
                  walletType: wallet.walletType ?? 0,
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
            }

            return GestureDetector(
              onTap: () {
                appState.viewData = {'walletMode': selectedWalletMode};
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: AllWalletsViewPageConfig,
                );
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
                                      fontSize: 15,
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

  Widget assetsTabs() {
    return Container(
      height: height / 2.2,
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
                        fontSize: 13,
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
          if (claimedAssets.length > 0) ...[
            Container(
              height: height / 2.2,
              child: ReorderableListView(
                padding: EdgeInsets.fromLTRB(0, 0, 0, 30),
                onReorder: (oldIndex, newIndex) {
                  print('reodered $oldIndex $newIndex');
                  if (oldIndex < newIndex) {
                    newIndex -= 1;
                  }
                  final Asset item = claimedAssets.removeAt(oldIndex);
                  claimedAssets.insert(newIndex, item);
                  setState(() {});
                },
                children: [
                  for (var i = 0; i < claimedAssets.length; i++) ...[
                    if (claimedAssets[i].assetCode != '') ...[
                      GestureDetector(
                        key: Key(i.toString()),
                        onTap: () {
                          appState.setActiveWallet = wallets.firstWhere(
                              (wallet) => wallet.publicKey == activeWallet);

                          appState.viewData = {
                            'assetCode': claimedAssets[i].assetCode,
                            'assetIssuer': claimedAssets[i].assetIssuer,
                            'walletPublicKey': activeWallet,
                          };
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDetailsViewPageConfig,
                          );
                        },
                        child: tiles(claimedAssets[i], i, activeWalletIndex),
                      ),
                    ]
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
    print('reodering assets... ${appState.assetOrderings}');
    if (appState.assetOrderings[publicKey] != null) {
      claimedAssets.forEach((asset) => asset.userPreferredIndex =
          appState.assetOrderings[publicKey]![asset.assetCode] ?? 0);
      claimedAssets
          .sort((a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex));
    }
  }

  void refreshData() async {
    try {
      await appState.refreshData();
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    StoreData().storeInsertData('assetOrderings', appState.assetOrderings);
    super.dispose();
  }
}
