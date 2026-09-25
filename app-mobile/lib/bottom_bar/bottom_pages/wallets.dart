import 'package:carousel_slider/carousel_slider.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/bottom_bar/bottom_pages/home.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/wallets_list_view_data.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import 'package:trovo_app/widgets/wallet_slides.dart';

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
  int isContractAddressWallet = 0;
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  String password = '';
  var noXbnBalance = false;
  late Account primaryWalletKeyPair;
  late Account newSubWalletKeyPair;
  late DataProvider appState;
  List<Asset> otherTokens = [];
  late UserInfo userInfo;
  final carouselController = CarouselSliderController();
  int tabLength = 1;
  int activeTabIndex = 0;
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;
  late RefreshController _refreshController;
  String selectedWalletMode = "My wallets";
  List<String> walletListMode = ['My wallets', 'Shared wallets', 'All wallets'];
  late List<WalletTileColor> colors;
  List<Asset> tokenizedAssets = [];

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(
        DropdownMenuItem(
          child: Text(wallet.alias!, overflow: TextOverflow.ellipsis),
          value: wallet.address,
        ),
      );
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
      items.add(
        DropdownMenuItem(
          child: Text(key, overflow: TextOverflow.ellipsis),
          value: value,
        ),
      );
    });
    return items;
  }

  List<DropdownMenuItem<String>> get walletListModeDropdownItems {
    return walletListMode
        .map<DropdownMenuItem<String>>(
          (item) => DropdownMenuItem(
            child: Text(item, overflow: TextOverflow.ellipsis),
            value: item,
          ),
        )
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

    WidgetsBinding.instance.addPostFrameCallback((_) {
      carouselController.animateToPage(
        0,
        duration: Duration(milliseconds: 500),
        curve: Curves.easeInOut,
      );
    });
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
    wallets =
        switch (selectedWalletMode) {
          'My wallets' => userInfo.mySolelyOwnedWallets,
          'Shared wallets' => userInfo.sharedWallets,
          'All wallets' => userInfo.allWallets,
          _ => [],
        } ??
        [];
    carouselWallets = wallets.length > 6
        ? wallets.getRange(0, 6).toList()
        : wallets;

    // if ((activeWallet == null && wallets.length > 0)) {
    //   activeWallet = wallets[0].address;
    // }

    if ((activeWallet == null && wallets.length > 0) || noXbnBalance) {
      activeWallet = wallets[0].address;
      otherTokens = wallets[0]
          .getOtherTokens(appState)
          .where((asset) => asset.assetCode != '' && asset.contractAddress != '')
          .toList();
      tokenizedAssets = wallets[0].getTokenizedAssets(appState);
      noXbnBalance =
          wallets[0].claimedAssets!
              .firstWhere(
                (asset) =>
                    asset.assetCode!.isEmpty && asset.contractAddress!.isEmpty,
              )
              .amount ==
          0;

      if (tokenizedAssets.isEmpty) {
        listMode = DashboardAssetListMode.OtherAssets;
      }

      reOrderClaimedAssets(activeWallet!);
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
                            vertical: 0,
                            horizontal: 20,
                          ),
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
                          fontWeight: FontWeight.w500,
                        ),
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
                    SizedBox(width: width / 15),
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
                              mainAxisAlignment: MainAxisAlignment.spaceAround,
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
                                          color: notifier.getbluewhitecolor,
                                          overflow: TextOverflow.visible,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
            if (wallets.isNotEmpty) ...[
              walletSlides(wallets),
              SizedBox(height: height / 70),
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
                                  'contractAddress': '',
                                  'walletAddress': activeWallet,
                                };
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: AssetDetailsViewPageConfig,
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
                                            color: notifier.getbluewhitecolor,
                                            overflow: TextOverflow.visible,
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      ),
                      SizedBox(width: width / 20),
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
                            listMode,
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
            SizedBox(height: height / 80),
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
                bottom: MediaQuery.of(context).viewInsets.bottom,
              ),
            ),
          ],
        ),
      ),
    );
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
                      text: "othertokens"
                          .tr()
                          .toLowerCase()
                          .capitalizeEachWord(),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
        SizedBox(height: height / 70),
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
                      text: "assettokens"
                          .tr()
                          .toLowerCase()
                          .capitalizeEachWord(),
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
                if (tokenizedAssets.length > 0) ...[
                  Container(
                    height: height / 2.2,
                    child: ReorderableListView(
                      padding: EdgeInsets.fromLTRB(0, 0, 0, 30),
                      onReorder: (oldIndex, newIndex) {
                        if (oldIndex < newIndex) {
                          newIndex -= 1;
                        }
                        final Asset item = tokenizedAssets.removeAt(oldIndex);
                        tokenizedAssets.insert(newIndex, item);
                        setState(() {});
                      },
                      children: [
                        for (var i = 0; i < tokenizedAssets.length; i++) ...[
                          if (tokenizedAssets[i].assetCode != '') ...[
                            GestureDetector(
                              key: Key(i.toString()),
                              onTap: () {
                                appState.setActiveWallet = wallets.firstWhere(
                                  (wallet) => wallet.address == activeWallet,
                                );

                                appState.viewData = {
                                  'assetCode': tokenizedAssets[i].assetCode,
                                  'contractAddress': tokenizedAssets[i].contractAddress,
                                  'walletAddress': activeWallet,
                                };
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: AssetDetailsViewPageConfig,
                                );
                              },
                              child: tiles(
                                tokenizedAssets[i],
                                i,
                                activeWalletIndex,
                              ),
                            ),
                          ],
                        ],
                      ],
                    ),
                  ),
                  SizedBox(height: height / 22),
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
                      ),
                    ),
                  ),
                ],
                SizedBox(height: height / 22),
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
      carouselController: carouselController,
      options: CarouselOptions(
        onPageChanged: ((index, reason) => {
          setState(() {
            activeWalletIndex = index == 5 ? index - 1 : index;
            activeWallet = wallets[activeWalletIndex].address;
            otherTokens = wallets[activeWalletIndex]
                .getOtherTokens(appState)
                .where(
                  (asset) => asset.assetCode != '' && asset.contractAddress != '',
                )
                .toList();
            tokenizedAssets = wallets[activeWalletIndex].getTokenizedAssets(
              appState,
            );
          }),
          // reOrderClaimedAssets(activeWallet!),
        }),
        height: 153,
        initialPage: wallets.length > 1 ? 1 : 0,
        padEnds: false,
        enableInfiniteScroll: false,
        clipBehavior: Clip.antiAlias,
        viewportFraction: 1,
      ),
      items: carouselWallets.map((wallet) {
        var indexOfWallet = wallets.indexOf(wallet);
        return Builder(
          builder: (BuildContext context) {
            if (indexOfWallet < 5) {
              return GestureDetector(
                onTap: () {
                  appState.viewData = {
                    'walletAddress': wallets[indexOfWallet].address,
                  };
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: WalletDetailsViewPageConfig,
                  );
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
                      : '${totalAccountBalanceInUSD(appState, wallets[indexOfWallet].claimedAssets!)} USD',
                  initialHiddenState: appState.hideWalletList[indexOfWallet],
                  onHiddenStateChanged: (state) => {
                    setState(() {
                      appState.hideWalletList[indexOfWallet] = state;
                      StoreData().storeInsertData(
                        'hideWalletList',
                        appState.hideWalletList,
                      );
                    }),
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
                              vertical: 35.0,
                              horizontal: 20,
                            ),
                            child: Image.asset(
                              'assets/images/trovo_white.png',
                              width: 80,
                            ),
                          ),
                        ],
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 20.0,
                          vertical: 25.0,
                        ),
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
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                                SizedBox(width: width / 50),
                                Icon(Icons.arrow_forward, color: wihitecolor),
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
        children: [showTokenAssets()],
      ),
    );
  }

  Widget tiles(Asset asset, int? indexOfAsset, int indexOfWallet) {
    if (indexOfAsset != null) {
      if (appState.assetOrderings[wallets[indexOfWallet].address!] == null) {
        appState.assetOrderings[wallets[indexOfWallet].address!] = {
          asset.assetCode!: indexOfAsset,
        };
      }
      appState.assetOrderings[wallets[indexOfWallet].address!]![asset
              .assetCode!] =
          indexOfAsset;
    }
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(15.0)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
          title: Row(
            children: [
              CircleAvatar(
                radius: 23,
                backgroundColor: Colors.white,
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(100.0),
                  child: Image.network(
                    asset.imageUrl!,
                    height: 40,
                    width: 40,
                    fit: BoxFit.fill,
                    errorBuilder: (context, error, stackTrace) {
                      return Image.asset(
                        'assets/images/trovo.png',
                        height: 40,
                        width: 40,
                      );
                    },
                  ),
                ),
              ),
              SizedBox(width: 5),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    getAssetCode(asset.assetCode),
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getblck,
                      overflow: TextOverflow.visible,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Container(
                      constraints: BoxConstraints(maxWidth: 100.sp),
                      child: Text(
                        asset.tokenizedAsset
                            ? asset.usdPrice.toString()
                            : "${getFiatRate(asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                          overflow: TextOverflow.visible,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
          trailing: Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                getBalance(
                  formatHistoryNumber(asset.amount!, 99000000000),
                  indexOfWallet,
                ),
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
                    asset.tokenizedAsset
                        ? '${calculateFiatValue(asset.amount.toString(), asset.usdPrice.toString(), 'USD', appState)} NGN'
                        : '${calculateFiatValue(asset.amount.toString(), asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                    indexOfWallet,
                  ),
                  style: TextStyle(
                    fontSize: 9,
                    fontFamily: fontbody,
                    color: notifier.getblck,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget showTokenAssets() {
    return SingleChildScrollView(
      child: Column(
        children: [
          if (otherTokens.length > 0) ...[
            Container(
              height: height / 2.2,
              child: ReorderableListView(
                padding: EdgeInsets.fromLTRB(0, 0, 0, 30),
                onReorder: (oldIndex, newIndex) {
                  if (oldIndex < newIndex) {
                    newIndex -= 1;
                  }
                  final Asset item = otherTokens.removeAt(oldIndex);
                  otherTokens.insert(newIndex, item);
                  setState(() {});
                },
                children: [
                  for (var i = 0; i < otherTokens.length; i++) ...[
                    if (otherTokens[i].assetCode != '') ...[
                      GestureDetector(
                        key: Key(i.toString()),
                        onTap: () {
                          appState.setActiveWallet = wallets.firstWhere(
                            (wallet) => wallet.address == activeWallet,
                          );

                          appState.viewData = {
                            'assetCode': otherTokens[i].assetCode,
                            'contractAddress': otherTokens[i].contractAddress,
                            'walletAddress': activeWallet,
                          };
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDetailsViewPageConfig,
                          );
                        },
                        child: tiles(otherTokens[i], i, activeWalletIndex),
                      ),
                    ],
                  ],
                ],
              ),
            ),
            SizedBox(height: height / 22),
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
                ),
              ),
            ),
          ],
          SizedBox(height: height / 22),
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

  void reOrderClaimedAssets(String address) {
    // order asset according to user preference
    if (appState.assetOrderings[address] != null) {
      otherTokens.forEach(
        (asset) => asset.userPreferredIndex =
            appState.assetOrderings[address]![asset.assetCode] ?? 0,
      );
      otherTokens.sort(
        (a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex),
      );
    }
  }

  void refreshData() async {
    try {
      await appState.refreshData();
      _refreshController.refreshCompleted();
      setState(() {});
    } catch (e) {
      _refreshController.refreshFailed();
      setState(() {});
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    StoreData().storeInsertData('assetOrderings', appState.assetOrderings);
    super.dispose();
  }
}
