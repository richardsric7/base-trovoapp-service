import 'dart:async';
import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart' hide Trans;
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/bottom_tab_page.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  late RefreshController _refreshController;
  late DataProvider appState;
  late UserInfo userInfo;
  late List<Wallet> wallets;
  late List<Wallet> sharedWallets;
  List<Asset>? unclaimedAssets;
  List<Asset>? claimedAssets;
  String? activeWallet;
  var noOfTransactionsToSign;
  var noXbnBalance = false;
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  DashboardAssetListMode listMode = DashboardAssetListMode.TokenizedAssets;
  late Future<List<TokenizedAsset>> primaryOffersListFuture;
  late Future<List<TokenizedAsset>> secondaryListItemsFuture;
  List<TokenizedAsset> tokenizedAssets = [];

  Map<String, DashboardAssetListMode> listModes = {
    'Asset Tokens': DashboardAssetListMode.TokenizedAssets,
    'Other Tokens': DashboardAssetListMode.OtherAssets,
  };
  final Authenticator _authenticator = Authenticator();
  bool showNewUserView = false;

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    appState.filterQuery = "&transactionStatus=PENDING";
    appState.getApprovals();
    userInfo = appState.userInfo!;
    wallets = userInfo.wallets!;
    sharedWallets = userInfo.sharedWallets!;

    if ((activeWallet == null && wallets.length > 0) || noXbnBalance) {
      activeWallet = wallets[0].publicKey;
      claimedAssets = wallets[0].claimedAssets;
      unclaimedAssets = wallets[0].unClaimedAssets;
      noXbnBalance =
          claimedAssets!
              .firstWhere(
                (asset) =>
                    asset.assetCode!.isEmpty && asset.assetIssuer!.isEmpty,
              )
              .amount ==
          0;

      reOrderClaimedAssets(activeWallet!);
    }

    if (!noXbnBalance) {
      primaryOffersListFuture = fetchTokenizationList(status: 0);
      secondaryListItemsFuture = fetchTokenizationList(status: 1);
    }
  }

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

  @override
  void dispose() {
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
                  SizedBox(height: 5),
                  firstRow(),
                  if (userInfo.hasSecurityQuestions == 0 ||
                      showNewUserView) ...[
                    SizedBox(height: 5),
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
                          },
                        };

                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: SecurityQuestionsViewPageConfig,
                        );
                      },
                      child: Container(
                        color: Colors.red[400],
                        child: Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Row(
                            children: [
                              SizedBox(width: width / 50),
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
                              SizedBox(width: width / 50),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ],
                  SizedBox(height: height / 50),
                  GestureDetector(
                    onTap: () {
                      changeTabPage(appState, ButtomTabPage.Wallets.index);
                      setState(() {});
                    },
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius: const BorderRadius.all(
                            Radius.circular(15.0),
                          ),
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
                                    vertical: 35.0,
                                    horizontal: 20,
                                  ),
                                  child: Image.asset(
                                    'assets/images/trovo_white.png',
                                    width: 40,
                                  ),
                                ),
                              ],
                            ),
                            Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 20.0,
                                vertical: 20,
                              ),
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
                                                fontFamily: fontsemibold,
                                              ),
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
                                                    for (
                                                      var i = 0;
                                                      i <
                                                          appState
                                                              .hideWalletList
                                                              .length;
                                                      i++
                                                    ) {
                                                      appState.hideWalletList[i] =
                                                          true;
                                                    }
                                                    StoreData().storeInsertData(
                                                      'hideWalletList',
                                                      appState.hideWalletList,
                                                    );
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
                                  SizedBox(height: height / 50),
                                  Container(
                                    width: width / 1.8,
                                    child: Text(
                                      getBalance(
                                        '${totalAccountBalanceInLocalCurrency} ${appState.defaultCurrency}',
                                      ),
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
                                        '${totalAccountBalanceInUSD} USD',
                                      ),
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
                  SizedBox(height: height / 50),
                  // check if the user's xbn balance is 0. This usually is the si-
                  // tuation when a new user signs up and has not funded their wallet
                  // yet
                  if (!noXbnBalance) ...[
                    SizedBox(height: height / 50),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 12.0),
                      child: Column(
                        children: [
                          Column(
                            children: [
                              Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceBetween,
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  // todo: remove the gesture detector
                                  Flexible(
                                    child: Text(
                                      "Tokenized Assets",
                                      textScaler: TextScaler.linear(1.0),
                                      style: TextStyle(
                                        fontSize: 20,
                                        color: notifier.getbluewhitecolor,
                                        fontWeight: FontWeight.w600,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              SizedBox(height: height / 50),
                              primaryOffers(),
                              SizedBox(height: height / 50),
                              secondaryListing(),
                            ],
                          ),
                          SizedBox(height: height / 70),
                        ],
                      ),
                    ),
                  ] else ...[
                    showFundWallet(),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget primaryOffers() {
    return SingleChildScrollView(
      child: Column(
        children: [
          FutureBuilder<List<TokenizedAsset>>(
            future: primaryOffersListFuture,
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
              } else if (snapshot.connectionState == ConnectionState.done) {
                if (snapshot.hasError) {
                  return Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: SizedBox(
                      height: height / 6,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "somethingwentwrong".tr(),
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 16,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                          ElevatedButton(
                            onPressed: () {
                              setState(() {
                                primaryOffersListFuture = fetchTokenizationList(
                                  status: 0,
                                );
                              });
                            },
                            style: ButtonStyle(
                              backgroundColor: WidgetStateProperty.all<Color>(
                                notifier.getbluecolor!,
                              ),
                              foregroundColor: WidgetStateProperty.all<Color>(
                                notifier.getwihitecolor,
                              ),
                            ),
                            child: Text(
                              "retry".tr(),
                              style: TextStyle(fontFamily: fontsemibold),
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                } else if (snapshot.hasData) {
                  tokenizedAssets = snapshot.data!;
                  return Column(
                    children: [
                      if (tokenizedAssets.isNotEmpty) ...[
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Flexible(
                              child: Text(
                                "Primary Offers",
                                textScaler: TextScaler.linear(1.0),
                                style: TextStyle(
                                  fontSize: 14,
                                  color: notifier.getbluewhitecolor,
                                  fontWeight: FontWeight.w600,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                            TextButton(
                              onPressed: () {
                                appState.viewData = {'rel': 0};
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: SeeAllTokenizedAssetsViewPageConfig,
                                );
                              },
                              style: TextButton.styleFrom(
                                padding:
                                    EdgeInsets.zero, // removes default padding
                                minimumSize: Size(
                                  0,
                                  0,
                                ), // removes minimum size constraints
                                tapTargetSize: MaterialTapTargetSize
                                    .shrinkWrap, // adjusts tap target size
                              ),
                              child: Text(
                                "View all",
                                textScaler: TextScaler.linear(1.0),
                                textAlign: TextAlign.right,
                                style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  decoration: TextDecoration.underline,
                                  fontSize: 12.0,
                                ),
                              ),
                            ),
                          ],
                        ),
                        for (var item in tokenizedAssets) ...[
                          GestureDetector(
                            onTap: () {
                              appState.tokenizedAsset = item;
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: TokenizedAssetDetailViewPageConfig,
                              );
                            },
                            child: tokenizedAssetTile(
                              notifier: notifier,
                              asset: item,
                              onSubscribe: () {
                                if (userInfo.kycVerified == null ||
                                    userInfo.kycVerified == 0) {
                                  kycUnverifiedErrorPop(context);
                                  return;
                                }

                                showSubscribePopup(
                                  context,
                                  asset: item,
                                  onDone: (amount) async {
                                    await subscribeTokenizedAsset(
                                      amount: double.parse(amount),
                                      tokenizedAssetID: item.id!,
                                    );
                                    setState(() {
                                      item.expressedInterest = true;
                                      item.expressedInterestAmount =
                                          double.parse(amount);
                                    });
                                  },
                                );
                              },
                              onBuyToken: () {
                                if (userInfo.kycVerified == null ||
                                    userInfo.kycVerified == 0) {
                                  kycUnverifiedErrorPop(context);
                                  return;
                                }

                                showBuyTokenPopup(
                                  context,
                                  assetCode: item.assetCode!,
                                  onDone: (wallet) {
                                    appState.setActiveWallet = wallet;
                                    appState.tokenizedAsset = item;
                                    appState.currentAction = PageAction(
                                      state: PageState.addPage,
                                      page: BuyTokensViewPageConfig,
                                    );
                                  },
                                  dropdownItems: getStandardWallets,
                                );
                              },
                            ),
                          ),
                        ],
                      ] else ...[
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Image.asset('assets/images/folder.png'),
                              Text(
                                "notokenizedasset".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                              SizedBox(height: 10),
                              Text(
                                "There are no tokenized assets in primary offering yet. Keep an eye out for future listings!"
                                    .tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody,
                                ),
                              ),
                              SizedBox(height: 30),
                            ],
                          ),
                        ),
                      ],
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
    );
  }

  Widget secondaryListing() {
    return SingleChildScrollView(
      child: Column(
        children: [
          FutureBuilder<List<TokenizedAsset>>(
            future: secondaryListItemsFuture,
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
              } else if (snapshot.connectionState == ConnectionState.done) {
                if (snapshot.hasData) {
                  var records = snapshot.data!;
                  return Column(
                    children: [
                      if (records.isNotEmpty) ...[
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Flexible(
                              child: Text(
                                "Secondary Listing",
                                textScaler: TextScaler.linear(1.0),
                                style: TextStyle(
                                  fontSize: 14,
                                  color: notifier.getbluewhitecolor,
                                  fontWeight: FontWeight.w600,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                            TextButton(
                              onPressed: () {
                                appState.viewData = {'rel': 1};
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: SeeAllTokenizedAssetsViewPageConfig,
                                );
                              },
                              style: TextButton.styleFrom(
                                padding:
                                    EdgeInsets.zero, // removes default padding
                                minimumSize: Size(
                                  0,
                                  0,
                                ), // removes minimum size constraints
                                tapTargetSize: MaterialTapTargetSize
                                    .shrinkWrap, // adjusts tap target size
                              ),
                              child: Text(
                                "View all",
                                textScaler: TextScaler.linear(1.0),
                                textAlign: TextAlign.right,
                                style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  decoration: TextDecoration.underline,
                                  fontSize: 12.0,
                                ),
                              ),
                            ),
                          ],
                        ),
                        for (var item in records) ...[
                          GestureDetector(
                            onTap: () {
                              appState.tokenizedAsset = item;
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: TokenizedAssetDetailViewPageConfig,
                              );
                            },
                            child: tokenizedAssetTile(
                              notifier: notifier,
                              asset: item,
                            ),
                          ),
                        ],
                      ],
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
                    SizedBox(width: width / 40),
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
                    SizedBox(width: width / 20),
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
                        SizedBox(height: 5),
                        GestureDetector(
                          onLongPress: () {
                            setState(() {
                              if (appState.walletMode.toLowerCase() ==
                                  'testnet') {
                                showNewUserView = !showNewUserView;
                                if (showNewUserView) {
                                  noXbnBalance = true;
                                } else {
                                  noXbnBalance = false;
                                }
                              }
                            });
                          },
                          child: Text(
                            truncate(
                              userInfo.username!.capitalizeFirst!,
                              length: 8,
                            ),
                            style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 17,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                      ],
                    ),
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
                      page: SharedAccessViewPageConfig,
                    );

                    if (noOfTransactionsToSign != null &&
                        noOfTransactionsToSign > 0) {
                      // take the user to the pending approvals tab on the shared access view
                      WidgetsBinding.instance.addPostFrameCallback((_) {
                        appState.sharedAccesstabController.animateTo(
                          1,
                          duration: Duration(milliseconds: 500),
                          curve: Curves.easeInOut,
                        );
                      });
                    }
                  },
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      vertical: 8.0,
                      horizontal: 10.0,
                    ),
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
                      state: PageState.addPage,
                      page: QrScannerPageConfig,
                    );
                  },
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      vertical: 8.0,
                      horizontal: 10.0,
                    ),
                    child: Icon(
                      Icons.qr_code_scanner_sharp,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
                GestureDetector(
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: NotificationsViewPageConfig,
                    );
                    appState.hasNewAnnouncement = false;
                  },
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      vertical: 8.0,
                      horizontal: 10.0,
                    ),
                    // child: SvgPicture.asset(
                    //   appState.hasNewAnnouncement
                    //       ? "assets/images/notifications-active.svg"
                    //       : "assets/images/notifications.svg",

                    //   height: height / 40,
                    // ),
                    child: Stack(
                      children: [
                        Icon(
                          Icons.notifications_none,
                          size: 27,
                          color: notifier.getbluewhitecolor,
                        ),
                        if (appState.hasNewAnnouncement) ...[
                          Positioned(
                            right: 0,
                            child: Container(
                              width: 10,
                              height: 10,
                              decoration: BoxDecoration(
                                borderRadius: const BorderRadius.all(
                                  Radius.circular(15.0),
                                ),
                                border: Border.all(
                                  color: notifier.getwihitecolor,
                                  width: 1,
                                ),
                                color: notifier.getgreencolor,
                              ),
                            ),
                          ),
                        ],
                      ],
                    ),
                  ),
                ),
                SizedBox(width: height / 50),
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
                              padding: const EdgeInsets.symmetric(
                                horizontal: 7,
                              ),
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
                ],
              ],
            ),
          ],
        ),
      ],
    );
  }

  List<DropdownMenuItem<Wallet>> get getStandardWallets {
    List<DropdownMenuItem<Wallet>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(
        DropdownMenuItem(
          child: Text(wallet.alias!, overflow: TextOverflow.ellipsis),
          value: wallet,
        ),
      );
    });
    return wallets;
  }

  void reOrderClaimedAssets(String publicKey) {
    // order asset according to user preference
    if (appState.assetOrderings[publicKey] != null) {
      claimedAssets!.forEach(
        (asset) => asset.userPreferredIndex =
            appState.assetOrderings[publicKey]![asset.assetCode] ?? 0,
      );
      claimedAssets!.sort(
        (a, b) => a.userPreferredIndex.compareTo(b.userPreferredIndex),
      );
    }
  }

  void refreshData() async {
    try {
      if (!noXbnBalance) {
        primaryOffersListFuture = fetchTokenizationList(status: 0);
        secondaryListItemsFuture = fetchTokenizationList(status: 1);
      }
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
              padding: const EdgeInsets.symmetric(
                horizontal: 20.0,
                vertical: 15.0,
              ),
              child: Column(
                children: [
                  SizedBox(height: height / 50),
                  Image.asset(
                    'assets/images/rafiki-buy-xbn.png',
                    // height: 50,
                    width: 180,
                  ),
                  SizedBox(height: height / 60),
                  Text(
                    "yourwalletisready".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(height: height / 90),
                  Text(
                    "butyoucannotuseityet".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(height: height / 50),
                  Text(
                    "youcangetbantutokens".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ),
          ),
        ),
        SizedBox(height: height / 50),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            HalfButtonWithIcon(
              "buywithfiat".tr(),
              notifier.getbluecolor,
              wihitecolor,
              'assets/images/cash.png',
              width: width / 2.2,
              height: 60,
              onTap: () async {
                await fetchFiatAmountForActivation();
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: BuyXBNWithFiatViewPageConfig,
                );
              },
            ),
            SizedBox(width: 10),
            HalfButtonWithIcon(
              "buyfromp2p".tr(),
              notifier.getbluecolor,
              wihitecolor,
              'assets/images/peers.png',
              width: width / 2.2,
              height: 60,
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
          ],
        ),
        SizedBox(height: 10),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            HalfButtonWithIcon(
              "requestfromuser".tr(),
              notifier.getbluecolor,
              wihitecolor,
              'assets/images/arrow-diagonal-down.png',
              width: width / 2.2,
              height: 60,
              onTap: () {
                appState.viewData = {
                  'assetCode': '',
                  'assetIssuer': '',
                  'walletPublicKey': activeWallet,
                };
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: RequestSpecificPaymentViewPageConfig,
                );
              },
            ),
            SizedBox(width: 10),
            HalfButtonWithIcon(
              "sendxbntoyourwallet".tr(),
              notifier.getbluecolor,
              wihitecolor,
              'assets/images/arrow-diagonal-up.png',
              width: width / 2.2,
              height: 60,
              onTap: () {
                Clipboard.setData(
                  ClipboardData(text: appState.primaryWallet.publicKey!),
                );
                showSnackBar("publickey".tr(), context);
              },
            ),
          ],
        ),
        SizedBox(height: height / 20),
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
          balance += double.parse(
            calculateFiatValue(
              asset.amount.toString(),
              asset.usdPrice.toString(),
              appState.defaultCurrency,
              appState,
            ).replaceAll(',', ''),
          );
        }
      }
    }
    return formatHistoryNumber(
      double.parse(balance.toString()),
      1000000,
      isShort: true,
    );
  }

  String get totalAccountBalanceInUSD {
    double balance = 0;
    for (var wallet in userInfo.allWallets) {
      if (wallet.claimedAssets!.length > 0) {
        for (var asset in wallet.claimedAssets!) {
          balance += double.parse(
            calculateFiatValue(
              asset.amount.toString(),
              asset.usdPrice.toString(),
              'USD',
              appState,
            ).replaceAll(',', ''),
          );
        }
      }
    }
    return formatHistoryNumber(
      double.parse(balance.toString()),
      1000000,
      isShort: true,
    );
  }

  Future<void> fetchTokenizationData() async {
    var uri = '/v1/tokenization';

    Map responseData = await makeGetRequest(
      uri: Uri.encodeFull(uri),
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.signer!,
    );
    if (responseData['statusCode'] == 200) {
      appState.tokenizationData = responseData['data'];
    }
  }

  Future<List<TokenizedAsset>> fetchTokenizationList({
    required int status,
  }) async {
    try {
      Future.wait([
        if (appState.tokenizationData.isEmpty) fetchTokenizationData(),
        if (appState.expressedInterests.isEmpty) fetchExpressedInterests(),
        if (appState.expressedInterests.isEmpty) fetchSubscriptions(),
      ]);
      var uri =
          '/v1/tokenization/list?onlyWithUserPermission=0&salesList=$status';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      // print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        List<TokenizedAsset> tokenizedAssets = [];
        var assets = responseData['data']['records'];
        if (assets != null) {
          for (int i = 0; i < assets.length; i++) {
            var a = TokenizedAsset().deserializeJson(assets[i]);
            if (appState.expressedInterests[a.id] != null) {
              a.expressedInterest = appState.expressedInterests[a.id] != null;
              a.expressedInterestAmount = double.parse(
                appState.expressedInterests[a.id]['amount'].toString(),
              );
            }

            if (appState.subscriptions[a.id] != null) {
              a.isSubscribed = appState.subscriptions[a.id] != null;
              a.subscriptionAmount = double.parse(
                appState.subscriptions[a.id]['amount'].toString(),
              );
            }
            tokenizedAssets.add(a);
          }
        }
        return tokenizedAssets;
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      print('error');
      print(e);
      return Future.error('Error! ${e}');
    }
  }

  Future<void> fetchExpressedInterests() async {
    try {
      var uri = '/v1/tokenization/expressed-interests';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        setState(() {
          appState.expressedInterests = {};
          var records = responseData['data']['records'];
          for (var i = 0; i < records.length; i++) {
            appState.expressedInterests[records[i]['tokenizedAssetId']] =
                records[i];
          }
        });
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      print('error');
      print(e);
      return Future.error('Error! ${e}');
    }
  }

  Future<void> fetchSubscriptions() async {
    try {
      var uri = '/v1/tokenization/subscriptions';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        appState.subscriptions = {};
        setState(() {
          var records = responseData['data']['records'];
          for (var i = 0; i < records.length; i++) {
            appState.subscriptions[records[i]['tokenizedAssetId']] = records[i];
          }
        });
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      print('error');
      print(e);
      return Future.error('Error! ${e}');
    }
  }

  subscribeTokenizedAsset({
    required double amount,
    required String tokenizedAssetID,
  }) async {
    try {
      showLoader(context);

      String requestBody = jsonEncode({'amount': amount});

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization/expressed-interests/${tokenizedAssetID}',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.publicKey!,
      );

      print('==============>response: $responseData');

      if (responseData['statusCode'] == 200) {
        hideLoader(context);
      } else {
        hideLoader(context);
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'].toString().isEmpty
              ? responseData['data']['error']
              : responseData['data']['message'],
        );
      }
    } catch (e) {
      // print(e);
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Future<void> fetchFiatAmountForActivation() async {
    ;
    try {
      showLoader(context);
      var uri = '/v1/users/activate/fiat';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      hideLoader(context);
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        appState.viewData = responseData['data'];
      }
    } catch (e) {
      print('error');
      print(e);
      hideLoader(context);
    }
  }
}

enum DashboardAssetListMode { TokenizedAssets, OtherAssets }
