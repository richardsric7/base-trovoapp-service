import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart' hide Trans;
import 'package:trovo_wallet/bottom_bar/bottom_pages/home.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
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
  List<Asset> claimedAssets = [];
  late bool localHideBalance;
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;
  String rel = '';

  var listOfAssets = <Map<String, String>>[
    // {
    //   "imageUrl": "",
    //   "assetName": "Animal Farm",
    //   "assetClass": "Agriculture",
    //   "balance": "2049"
    // },
    // {
    //   "imageUrl": "",
    //   "assetName": "Beacon Homes",
    //   "assetClass": "Property",
    //   "balance": "3250"
    // },
    // {
    //   "imageUrl": "",
    //   "assetName": "C-Vitals",
    //   "assetClass": "Health",
    //   "balance": "100"
    // },
    // {
    //   "imageUrl": "",
    //   "assetName": "Drinkfly",
    //   "assetClass": "Beverage",
    //   "balance": "4000"
    // }
  ];

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
    claimedAssets = wallet.claimedAssets!
        .where((asset) => asset.assetCode != '' && asset.assetIssuer != '')
        .toList();
    reOrderClaimedAssets(wallet.publicKey!);
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
              SizedBox(
                height: height / 80,
              ),
              Container(
                constraints: BoxConstraints(
                  maxHeight: height / 5.8,
                ),
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
                      : '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallet.claimedAssets!)} USD',
                  initialHiddenState: appState.hideBalances,
                  onHiddenStateChanged: (state) => {
                    setState(
                      () => localHideBalance = state,
                    )
                  },
                ),
              ),
              SizedBox(
                height: height / 80,
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
                                      'walletPublicKey': wallet.publicKey,
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
                                              '${"gas".tr()} ${formatHistoryNumber(gas.amount!, 1000)}',
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
                              "assettokens".tr(),
                              context,
                              null,
                            ),
                          ),
                        ),
                      ],
                    )),
              ),
              SizedBox(
                height: height / 95,
              ),
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
      claimedAssets.forEach((asset) => asset.userPreferredIndex =
          appState.assetOrderings[publicKey]![asset.assetCode] ?? 0);
      claimedAssets
          .sort((a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex));
      print('testing the microphone $claimedAssets');
    }
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
                      text: "assettokens".tr().toUpperCase(),
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
                if (listOfAssets.isNotEmpty) ...[
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
                ] else ...[
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: SizedBox(
                      height: height / 2,
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
                SizedBox(height: height / 20),
              ],
            ),
          ),
        ],
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
                      assetCode: 'asset.assetCode!',
                      onDone: (walletPublicKey) {},
                      dropdownItems: getStandardWallets,
                    )
                  : showSubscribePopup(
                      context,
                      assetCode: 'tokenizedAsset.assetCode!',
                      onDone: (walletPublicKey) {},
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

  Widget cryptoAssets() {
    return Container(
      height: height / 1.58,
      child: Column(
        children: [
          Padding(
            padding:
                EdgeInsets.symmetric(horizontal: tabLength > 1 ? 0.0 : 100.0),
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
              SizedBox(
                height: height / 40,
              ),
              Container(
                height: height / 1.72,
                child: TabBarView(
                  controller: _tabController,
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 10.0, 0, 0),
                      child: Column(
                        children: [
                          if (claimedAssets.length > 0) ...[
                            Container(
                              height: height / 1.76,
                              child: ReorderableListView(
                                padding: EdgeInsets.fromLTRB(0, 0, 0, 30),
                                onReorder: (oldIndex, newIndex) {
                                  if (oldIndex < newIndex) {
                                    newIndex -= 1;
                                  }
                                  final Asset item =
                                      claimedAssets.removeAt(oldIndex);
                                  claimedAssets.insert(newIndex, item);
                                  setState(() {});
                                },
                                children: [
                                  for (var i = 0;
                                      i < claimedAssets.length;
                                      i++) ...[
                                    GestureDetector(
                                      key: Key(claimedAssets[i].assetIssuer!),
                                      onTap: () {
                                        appState.returnView = PageAction(
                                          state: PageState.addAll,
                                          pages: [
                                            BottomHomePageConfig,
                                            WalletDetailsViewPageConfig
                                          ],
                                        );

                                        if (rel == 'sharedWalletView') {
                                          appState.returnView = PageAction(
                                              state: PageState.addAll,
                                              pages: [
                                                BottomHomePageConfig,
                                                SharedAccessViewPageConfig,
                                                SharedWalletInfoViewPageConfig,
                                                WalletDetailsViewPageConfig
                                              ]);
                                        }

                                        appState.viewData = {
                                          'assetCode':
                                              claimedAssets[i].assetCode,
                                          'assetIssuer':
                                              claimedAssets[i].assetIssuer,
                                          'walletPublicKey': wallet.publicKey,
                                        };
                                        appState.currentAction = PageAction(
                                          state: PageState.addPage,
                                          page: AssetDetailsViewPageConfig,
                                        );
                                      },
                                      child: tiles(claimedAssets[i], i),
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
                                      10, 28.0, 10, 0),
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
                                            WalletDetailsViewPageConfig
                                          ],
                                        );

                                        if (rel == 'sharedWalletView') {
                                          appState.returnView = PageAction(
                                              state: PageState.addAll,
                                              pages: [
                                                BottomHomePageConfig,
                                                SharedAccessViewPageConfig,
                                                SharedWalletInfoViewPageConfig,
                                                WalletDetailsViewPageConfig
                                              ]);
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
                                            10, 28.0, 10, 0),
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
                      ),
                    ],
                    // Container(
                    //   height: height / 2,
                    //   child: SingleChildScrollView(
                    //     child: Column(
                    //       children: [
                    //         // in situations where the blockchain has an issue,
                    //         // some values can be returned as null or empty
                    //         // so always null check for such situations
                    //         // if (nfts != null && nfts != {}) ...[
                    //         //   if (nfts[activeWallet!.publicKey] != null &&
                    //         //       nfts[activeWallet!.publicKey].length >
                    //         //           0) ...[
                    //         gridView(),
                    //         //     SizedBox(height: 600),
                    //         //   ] else ...[
                    //         //     showEmptyNFTs(),
                    //         //   ]
                    //         // ] else ...[
                    //         //   showEmptyNFTs(),
                    //         // ],
                    //       ],
                    //     ),
                    //   ),
                    // ),
                  ],
                ),
              ),
            ],
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

  Widget tiles(Asset asset, int? indexOfAsset) {
    if (indexOfAsset != null) {
      if (appState.assetOrderings[wallet.publicKey!] == null) {
        appState.assetOrderings[wallet.publicKey!] = {
          asset.assetCode!: indexOfAsset
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
                      child: Text(
                        "${getFiatRate(asset.usdPrice!.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}",
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
                        '${calculateFiatValue(asset.amount.toString(), asset.usdPrice.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}'),
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
