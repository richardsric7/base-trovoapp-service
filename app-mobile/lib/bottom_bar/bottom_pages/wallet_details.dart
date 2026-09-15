import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart' hide Trans;
import 'package:trovo_app/bottom_bar/bottom_pages/home.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:trovo_app/widgets/wallet_slides.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  late Asset gas;
  List<Asset> otherTokens = [];
  List<Asset> tokenizedAssets = [];
  late bool localHideBalance;
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;
  String rel = '';

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

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(
        DropdownMenuItem(
          child: Text(wallet.alias!, overflow: TextOverflow.ellipsis),
          value: wallet.publicKey,
        ),
      );
    });
    return wallets;
  }

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    appState = Provider.of<DataProvider>(context, listen: false);
    localHideBalance = appState.hideBalances;

    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    gas = wallet.claimedAssets!.where((asset) => asset.assetCode == '').first;

    rel = appState.viewData!['rel'] != null ? appState.viewData!['rel'] : '';
    tokenizedAssets = wallet.getTokenizedAssets(appState);
    otherTokens = wallet
        .getOtherTokens(appState)
        .where((asset) => asset.assetCode != '' && asset.assetIssuer != '')
        .toList();

    if (tokenizedAssets.isEmpty) {
      listMode = DashboardAssetListMode.OtherAssets;
    }

    reOrderClaimedAssets(wallet.publicKey!);
    reOrderTokenizedAssets(wallet.publicKey!);
  }

  void tabListener() {
    // Tab Changed swiping to a new tab
    activeTabIndex = _tabController.index;
    setState(() {});
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

    if (listMode == DashboardAssetListMode.OtherAssets) {
      // in order to make assets tab length dynamic we have to check
      // for when we have pending asset and then change the tablength
      // to 3 or back to 2 when we do not have pending assets.
      if (wallet.unClaimedAssets != null &&
          wallet.unClaimedAssets!.length > 0) {
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

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          wallet.alias!,
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 80),
              Container(
                constraints: BoxConstraints(maxHeight: height / 5.8),
                child: WalletSlide(
                  backColor: notifier.getbluecolor,
                  foreColor: wihitecolor,
                  alias: wallet.alias!.capitalizeFirst!,
                  isSharedWallet: wallet.isSharedWallet,
                  walletType: wallet.walletType ?? 0,
                  assetCount: wallet.claimedAssets?.length.toString(),
                  totalBalance:
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallet.claimedAssets!)} ${appState.defaultCurrency}',
                  fiatBalance: appState.defaultCurrency == 'USD'
                      ? null
                      : '${totalAccountBalanceInUSD(appState, wallet.claimedAssets!)} USD',
                  initialHiddenState: appState.hideBalances,
                  onHiddenStateChanged: (state) => {
                    setState(() => localHideBalance = state),
                  },
                ),
              ),
              SizedBox(height: height / 80),
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
                                  'walletPublicKey': wallet.publicKey,
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
                                          '${"gas".tr()} ${formatHistoryNumber(gas.amount!, 1000)}',
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
                            "assettokens".tr(),
                            context,
                            null,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 95),
              listMode == DashboardAssetListMode.TokenizedAssets
                  ? showAssets()
                  : cryptoAssets(),
            ],
          ),
        ),
      ),
    );
  }

  void reOrderClaimedAssets(String publicKey) {
    // order asset according to user preference
    if (appState.assetOrderings[publicKey] != null) {
      otherTokens.forEach(
        (asset) => asset.userPreferredIndex =
            appState.assetOrderings[publicKey]![asset.assetCode] ?? 0,
      );
      otherTokens.sort(
        (a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex),
      );
    }
  }

  void reOrderTokenizedAssets(String publicKey) {
    // order asset according to user preference
    if (appState.assetOrderings[publicKey] != null) {
      tokenizedAssets.forEach(
        (asset) => asset.userPreferredIndex =
            appState.assetOrderings[publicKey]![asset.assetCode] ?? 0,
      );
      tokenizedAssets.sort(
        (a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex),
      );
    }
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
      height: height / 1.68,
      child: TabBarView(
        controller: _tabController,
        children: [
          SingleChildScrollView(
            child: Column(
              children: [
                if (tokenizedAssets.length > 0) ...[
                  Container(
                    height: height / 1.76,
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
                          GestureDetector(
                            key: Key(tokenizedAssets[i].assetIssuer!),
                            onTap: () {
                              appState.returnView = PageAction(
                                state: PageState.addAll,
                                pages: [
                                  BottomHomePageConfig,
                                  WalletDetailsViewPageConfig,
                                ],
                              );

                              if (rel == 'sharedWalletView') {
                                appState.returnView = PageAction(
                                  state: PageState.addAll,
                                  pages: [
                                    BottomHomePageConfig,
                                    SharedAccessViewPageConfig,
                                    SharedWalletInfoViewPageConfig,
                                    WalletDetailsViewPageConfig,
                                  ],
                                );
                              }

                              appState.viewData = {
                                'assetCode': tokenizedAssets[i].assetCode,
                                'assetIssuer': tokenizedAssets[i].assetIssuer,
                                'walletPublicKey': wallet.publicKey,
                              };
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: AssetDetailsViewPageConfig,
                              );
                            },
                            child: tiles(tokenizedAssets[i], i),
                          ),
                        ],
                      ],
                    ),
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
                      ),
                    ),
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

  Widget tokenizedAssetTile(
    String imageUrl,
    String name,
    String type,
    isSubscribed,
  ) {
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 5),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(15.0)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
          title: Row(
            children: [
              CircleAvatar(
                radius: 20,
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(100.0),
                  child: Image.network(
                    imageUrl,
                    height: 35,
                    width: 35,
                    fit: BoxFit.fill,
                    errorBuilder: (context, error, stackTrace) {
                      return Image.asset(
                        'assets/images/trovo.png',
                        height: 35,
                        width: 35,
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
                    name,
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getblck,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Container(
                      constraints: BoxConstraints(maxWidth: 100.sp),
                      child: Text(
                        type,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget cryptoAssets() {
    return Container(
      height: height / 1.58,
      child: Column(
        children: [
          Padding(
            padding: EdgeInsets.symmetric(
              horizontal: tabLength > 1 ? 0.0 : 100.0,
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
                  text: "othertokens".tr().toLowerCase().capitalizeEachWord(),
                ),
                if (wallet.unClaimedAssets != null && tabLength == 2) ...[
                  Tab(
                    height: 20,
                    text:
                        '${"pending".tr()} (${wallet.unClaimedAssets == null ? 0 : wallet.unClaimedAssets!.length})',
                  ),
                ],
              ],
            ),
          ),
          Column(
            children: [
              SizedBox(height: height / 40),
              Container(
                height: height / 1.72,
                child: TabBarView(
                  controller: _tabController,
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 10.0, 0, 0),
                      child: Column(
                        children: [
                          if (otherTokens.length > 0) ...[
                            Container(
                              height: height / 1.76,
                              child: ReorderableListView(
                                padding: EdgeInsets.fromLTRB(0, 0, 0, 30),
                                onReorder: (oldIndex, newIndex) {
                                  if (oldIndex < newIndex) {
                                    newIndex -= 1;
                                  }
                                  final Asset item = otherTokens.removeAt(
                                    oldIndex,
                                  );
                                  otherTokens.insert(newIndex, item);
                                  setState(() {});
                                },
                                children: [
                                  for (
                                    var i = 0;
                                    i < otherTokens.length;
                                    i++
                                  ) ...[
                                    GestureDetector(
                                      key: Key(otherTokens[i].assetIssuer!),
                                      onTap: () {
                                        appState.returnView = PageAction(
                                          state: PageState.addAll,
                                          pages: [
                                            BottomHomePageConfig,
                                            WalletDetailsViewPageConfig,
                                          ],
                                        );

                                        if (rel == 'sharedWalletView') {
                                          appState.returnView = PageAction(
                                            state: PageState.addAll,
                                            pages: [
                                              BottomHomePageConfig,
                                              SharedAccessViewPageConfig,
                                              SharedWalletInfoViewPageConfig,
                                              WalletDetailsViewPageConfig,
                                            ],
                                          );
                                        }

                                        appState.viewData = {
                                          'assetCode': otherTokens[i].assetCode,
                                          'assetIssuer':
                                              otherTokens[i].assetIssuer,
                                          'walletPublicKey': wallet.publicKey,
                                        };
                                        appState.currentAction = PageAction(
                                          state: PageState.addPage,
                                          page: AssetDetailsViewPageConfig,
                                        );
                                      },
                                      child: tiles(otherTokens[i], i),
                                    ),
                                  ],
                                ],
                              ),
                            ),
                          ] else ...[
                            Container(
                              height: height / 3,
                              child: Padding(
                                padding: const EdgeInsets.fromLTRB(
                                  10,
                                  28.0,
                                  10,
                                  0,
                                ),
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
                        ],
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
                                          activeTabIndex = _tabController.index;
                                        });

                                        appState.returnView = PageAction(
                                          state: PageState.addAll,
                                          pages: [
                                            BottomHomePageConfig,
                                            WalletDetailsViewPageConfig,
                                          ],
                                        );

                                        if (rel == 'sharedWalletView') {
                                          appState.returnView = PageAction(
                                            state: PageState.addAll,
                                            pages: [
                                              BottomHomePageConfig,
                                              SharedAccessViewPageConfig,
                                              SharedWalletInfoViewPageConfig,
                                              WalletDetailsViewPageConfig,
                                            ],
                                          );
                                        }

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
                                      child: tiles(asset, null),
                                    ),
                                  ],
                                ] else ...[
                                  Container(
                                    height: height / 4,
                                    child: Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                        10,
                                        28.0,
                                        10,
                                        0,
                                      ),
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
                                      ),
                                    ),
                                  ),
                                ],
                                SizedBox(height: height / 22),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget tiles(Asset asset, int? indexOfAsset) {
    if (indexOfAsset != null) {
      if (appState.assetOrderings[wallet.publicKey!] == null) {
        appState.assetOrderings[wallet.publicKey!] = {
          asset.assetCode!: indexOfAsset,
        };
      }
      appState.assetOrderings[wallet.publicKey!]![asset.assetCode!] =
          indexOfAsset;
    }
    return Card(
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      elevation: notifier.isDark ? 0 : 5,
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
                    child: Container(
                      constraints: BoxConstraints(maxWidth: 100.sp),
                      child: Text(
                        "${getFiatRate(asset.usdPrice!.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getblck,
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
                    '${calculateFiatValue(asset.amount.toString(), asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
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
