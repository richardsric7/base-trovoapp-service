import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
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
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

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
  late List<Wallet> wallets;
  late List<Wallet> sharedWallets;
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
  late Future<List<TokenizedAsset>> listOfTokenizations;

  Map<String, DashboardAssetListMode> listModes = {
    'Asset Tokens': DashboardAssetListMode.TokenizedAssets,
    'Other Tokens': DashboardAssetListMode.OtherAssets,
  };
  final Authenticator _authenticator = Authenticator();

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
    listOfTokenizations = fetchTokenizationList();
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
    listOfTokenizations = fetchTokenizationList();
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
    sharedWallets = userInfo.sharedWallets!;

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
                    height: 5,
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
                  SizedBox(
                    height: height / 50,
                  ),
                  GestureDetector(
                    onTap: () {
                      changeTabPage(appState, ButtomTabPage.Wallets.index);
                      setState(() {});
                    },
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius:
                              const BorderRadius.all(Radius.circular(15.0)),
                          color: notifier.getbluecolor,
                          // color: colors[wallets.indexOf(wallet)],
                        ),
                        child: Stack(
                          alignment: AlignmentDirectional.centerEnd,
                          children: [
                            Row(
                              mainAxisAlignment: MainAxisAlignment.end,
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                Padding(
                                  padding: const EdgeInsets.symmetric(
                                      vertical: 35.0, horizontal: 20),
                                  child: Image.asset(
                                    'assets/images/trovo_white.png',
                                    width: 40,
                                  ),
                                ),
                              ],
                            ),
                            Padding(
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 20.0, vertical: 20),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceEvenly,
                                children: [
                                  Row(
                                    mainAxisAlignment:
                                        MainAxisAlignment.spaceBetween,
                                    children: [
                                      Container(
                                        width: width / 1.8,
                                        child: Row(
                                          mainAxisAlignment:
                                              MainAxisAlignment.spaceBetween,
                                          children: [
                                            Text(
                                              "totalaccountbalance".tr(),
                                              style: TextStyle(
                                                  fontSize: 15,
                                                  fontWeight: FontWeight.w600,
                                                  color: wihitecolor,
                                                  fontFamily: fontsemibold),
                                            ),
                                            SizedBox(),
                                            GestureDetector(
                                              onTap: () {
                                                if (appState.hideBalances) {
                                                  authenticateAndToggle();
                                                } else
                                                  setState(() {
                                                    appState.hideBalances =
                                                        !appState.hideBalances;
                                                    for (var i = 0;
                                                        i <
                                                            appState
                                                                .hideWalletList
                                                                .length;
                                                        i++) {
                                                      appState.hideWalletList[
                                                          i] = true;
                                                    }
                                                    StoreData().storeInsertData(
                                                        'hideWalletList',
                                                        appState
                                                            .hideWalletList);
                                                  });
                                              },
                                              child: Icon(
                                                getIcon(),
                                                size: 20,
                                                color: wihitecolor,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      SizedBox(),
                                    ],
                                  ),
                                  SizedBox(
                                    height: height / 50,
                                  ),
                                  Container(
                                    width: width / 1.8,
                                    child: Text(
                                      getBalance(
                                          '${totalAccountBalanceInLocalCurrency} ${appState.defaultCurrency}'),
                                      style: TextStyle(
                                        fontSize: 16,
                                        fontWeight: FontWeight.bold,
                                        color: wihitecolor,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                  if (appState.defaultCurrency != 'USD') ...[
                                    SizedBox(height: 10),
                                    Text(
                                      getBalance(
                                          '${totalAccountBalanceInUSD} USD'),
                                      style: TextStyle(
                                        fontWeight: FontWeight.w300,
                                        fontSize: 13,
                                        color: wihitecolor,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  // check if the user's xbn balance is 0. This usually is the si-
                  // tuation when a new user signs up and has not funded their wallet
                  // yet
                  if (!noXbnBalance) ...[
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 12.0),
                      child: Column(
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
                                        text: "primaryoffers".tr(),
                                      ),
                                      Tab(
                                        height: 20,
                                        text: "secondarylisting".tr(),
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

  Widget listingTabs() {
    return Container(
      height: height / 1.77,
      child: TabBarView(
        controller: _tabController,
        children: [
          SingleChildScrollView(
            child: Column(
              children: [
                FutureBuilder<List<TokenizedAsset>>(
                  future: listOfTokenizations,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.waiting) {
                      return SizedBox(
                        height: height / 2,
                        child: Center(
                          child: CircularProgressIndicator(
                            backgroundColor: notifier.getbluecolor,
                            valueColor: new AlwaysStoppedAnimation<Color>(
                              notifier.getgreencolor,
                            ),
                            strokeWidth: 3.0,
                          ),
                        ),
                      );
                    } else if (snapshot.connectionState ==
                        ConnectionState.done) {
                      if (snapshot.hasError) {
                        return Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: SizedBox(
                            height: height / 2,
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  "somethingwentwrong".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
                                ),
                                ElevatedButton(
                                  onPressed: () {
                                    setState(() {
                                      listOfTokenizations =
                                          fetchTokenizationList();
                                    });
                                  },
                                  style: ButtonStyle(
                                    backgroundColor:
                                        MaterialStateProperty.all<Color>(
                                            notifier.getbluecolor!),
                                  ),
                                  child: Text(
                                    "retry".tr(),
                                    style: TextStyle(
                                      fontFamily: fontsemibold,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        );
                      } else if (snapshot.hasData) {
                        var records = snapshot.data!;
                        return Column(
                          children: [
                            if (records.isNotEmpty) ...[
                              for (var item in records) ...[
                                GestureDetector(
                                  onTap: () {
                                    appState.tokenizedAsset = item;
                                    appState.currentAction = PageAction(
                                      state: PageState.addPage,
                                      page: TokenizedAssetDetailViewPageConfig,
                                    );
                                  },
                                  child: assetTile(item),
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
                            ]
                          ],
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
          SingleChildScrollView(
            child: Column(
              children: [
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
            ),
          ),
        ],
      ),
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

  Widget assetTile(TokenizedAsset asset) {
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
              Container(
                width: width / 3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${asset.assetName!} (${asset.assetCode})',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        asset.assetSector!,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          trailing: ElevatedButton(
            onPressed: () async {
              asset.isSubscribed ?? false
                  ? showUnSubscribePopup(
                      context,
                      assetCode: asset.assetCode!,
                      onDone: (walletPublicKey) {
                        setState(() {
                          asset.isSubscribed = false;
                        });
                      },
                      dropdownItems: getUnsubscribableWallets(asset),
                    )
                  : showSubscribePopup(
                      context,
                      assetCode: asset.assetCode!,
                      onDone: (walletPublicKey) async {
                        setState(() {
                          asset.isSubscribed = true;
                        });
                      },
                      dropdownItems: getStandardWallets,
                    );
            },
            style: ButtonStyle(
              padding: MaterialStateProperty.all(
                EdgeInsets.symmetric(vertical: 0, horizontal: 6),
              ),
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
              width: width / 5,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    asset.isSubscribed ?? false ? 'Subscribed' : 'Subscribe',
                    style: TextStyle(
                        fontFamily: fontsemibold,
                        fontSize: 10,
                        color: notifier.getwihitecolor),
                  ),
                  Icon(
                      asset.isSubscribed ?? false
                          ? Icons.check_circle
                          : Icons.add_circle_rounded,
                      size: 15,
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

  List<DropdownMenuItem<String>> getUnsubscribableWallets(
      TokenizedAsset asset) {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.where((wallet) {
      return wallet.tokenizedAssets != null &&
          wallet.tokenizedAssets!
              .where((a) =>
                  a.assetName == asset.assetName &&
                  a.assetCode == asset.assetCode)
              .isNotEmpty;
    }).forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet.publicKey));
    });
    return wallets;
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

  int getTotalNumberOfAssets() {
    int total = 0;
    wallets.forEach((wallet) => total += wallet.totalAssets);
    sharedWallets.forEach((wallet) => total += wallet.totalAssets);
    return total;
  }

  IconData getIcon() {
    IconData icon;
    if (appState.hideBalances) icon = CupertinoIcons.eye;

    if (appState.hideBalances)
      icon = CupertinoIcons.eye;
    else
      icon = CupertinoIcons.eye_slash;

    return icon;
  }

  void authenticateAndToggle() {
    if (appState.biometricEnabled) {
      toggleBiometrics();
      return;
    }

    showPasswordDialog(context, () {
      setState(() {
        appState.hideBalances = !appState.hideBalances;
      });
    });
  }

  void toggleBiometrics() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        setState(() {
          appState.hideBalances = !appState.hideBalances;
        });
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  String getBalance(String balance) {
    String text;
    if (appState.hideBalances) text = hideBalanceText;

    if (appState.hideBalances)
      text = hideBalanceText;
    else
      text = balance;

    return text;
  }

  String get totalAccountBalanceInLocalCurrency {
    double balance = 0;
    for (var wallet in userInfo.allWallets) {
      if (wallet.claimedAssets!.length > 0) {
        for (var asset in wallet.claimedAssets!) {
          balance += double.parse(calculateFiatValue(asset.amount.toString(),
                  asset.usdPrice.toString(), appState.defaultCurrency, appState)
              .replaceAll(',', ''));
        }
      }
    }
    return formatHistoryNumber(double.parse(balance.toString()), 1000000);
  }

  String get totalAccountBalanceInUSD {
    double balance = 0;
    for (var wallet in userInfo.allWallets) {
      if (wallet.claimedAssets!.length > 0) {
        for (var asset in wallet.claimedAssets!) {
          balance += double.parse(calculateFiatValue(asset.amount.toString(),
                  asset.usdPrice.toString(), 'USD', appState)
              .replaceAll(',', ''));
        }
      }
    }
    return formatHistoryNumber(double.parse(balance.toString()), 1000000);
  }

  Future<List<TokenizedAsset>> fetchTokenizationList() async {
    List<TokenizedAsset> tokenizedAssets = [];
    if (appState.tempTokenizedAssetList.isEmpty) {
      var savedAssets = await StoreData().storeGetData('tokenizedAsset');
      if (savedAssets != null) {
        for (int i = 0; i < savedAssets.length; i++) {
          print(savedAssets[i]);
          var a = TokenizedAsset().deserializeJson(savedAssets[i]);
          a.usdPrice = 1.47;
          a.assetIssuer = a.walletToHoldAssetsNotForSale ?? '';
          a.pricePerToken = (double.parse(a.assetCurrentValue.toString()) /
              a.numberOfTokenToBeIssued!);
          tokenizedAssets.add(a);
        }
      }
    }
    appState.tempTokenizedAssetList = tokenizedAssets;
    inspect(tokenizedAssets);
    return tokenizedAssets;
  }
}

enum DashboardAssetListMode { TokenizedAssets, OtherAssets }
