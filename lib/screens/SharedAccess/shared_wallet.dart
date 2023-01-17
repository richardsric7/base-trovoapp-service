import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/User.dart';
import 'package:trovo_wallet/models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/WalletSlides.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SharedWallet extends StatefulWidget {
  const SharedWallet({Key? key}) : super(key: key);

  @override
  State<SharedWallet> createState() => _SharedWalletState();
}

class _SharedWalletState extends State<SharedWallet>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late TabController _tabController;
  late RefreshController _refreshController;
  late DataProvider appState;
  late UserInfo userInfo;
  // var nfts;
  Wallet? activeWallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 1;
  int activeTabIndex = 0;
  late bool localHideBalance;
  var viewData;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    localHideBalance = appState.hideBalances;
    viewData = appState.viewData![SharedWalletDetailsViewPageConfig.key];
    claimedAssets = viewData!['claimed'];
    unclaimedAssets = viewData!['unclaimed'];
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
    appState = Provider.of<DataProvider>(context, listen: true);
    // nfts = appState.nfts;

    // in order to make assets tab length dynamic we have to check
    // for when we have pending asset and then change the tablength
    // to 3 or back to 2 when we do not have pending assets.
    if (unclaimedAssets != null && unclaimedAssets.length > 0) {
      tabLength = 2;
    } else {
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
          'Shared Wallet',
          notifier.getbluewhitecolor,
          height: height / 15,
        ),
        body: SmartRefresher(
          enablePullDown: true,
          controller: _refreshController,
          onRefresh: refreshData,
          child: SingleChildScrollView(
            child: Column(
              children: [
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
      height: height / 1.1,
      child: showWallet(),
    );
  }

  Widget showWallet() {
    return Column(
      children: [
        SizedBox(
          height: 20.sp,
        ),
        Padding(
          padding:
              EdgeInsets.symmetric(horizontal: tabLength > 1 ? 0.0 : 100.0),
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
              if (unclaimedAssets != null && tabLength == 2) ...[
                Tab(
                  height: 20,
                  text: '${LanguageEn.pending} (${unclaimedAssets.length})',
                ),
              ],
              // Tab(
              //   height: 20,
              //   text: LanguageEn.nfts,
              // ),
            ],
          ),
        ),
        SizedBox(
          height: 20.sp,
        ),
        Expanded(
          child: TabBarView(
            controller: _tabController,
            children: [
              Container(
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(0, 10.0, 0, 0),
                  child: SingleChildScrollView(
                    child: Column(
                      children: [
                        WalletSlide(
                          backColor: notifier.getbluecolor,
                          foreColor: wihitecolor,
                          alias: viewData['walletAlias']!,
                          totalBalance:
                              '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, claimedAssets)} ${appState.defaultCurrency}',
                          fiatBalance: appState.defaultCurrency == 'USD'
                              ? null
                              : '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, claimedAssets)} USD',
                          initialHiddenState: appState.hideBalances,
                          onHiddenStateChanged: (state) => {
                            setState(
                              () => {
                                localHideBalance = state,
                              },
                            )
                          },
                        ),
                        SizedBox(
                          height: height / 30,
                        ),
                        if (claimedAssets.length > 0) ...[
                          for (var asset in claimedAssets) ...[
                            GestureDetector(
                                onTap: () {
                                  // since the original asset object
                                  // is immutable I create a new assetObj and
                                  // copy all the data into it so that
                                  // I'll be able to change the data
                                  appState.viewData![
                                      SharedWalletAssetDetailsViewPageConfig
                                          .key] = {
                                    'assetCode': asset['assetCode'],
                                    'assetIssuer': asset['assetIssuer'],
                                    'amount': asset['amount'],
                                    'usdPrice': asset['usdPrice'],
                                    'qrCode': asset['qrCode'],
                                    'imageUrl': asset['imageUrl'],
                                    'walletInfo': viewData,
                                    'rel': viewData['rel'],
                                  };
                                  appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page:
                                        SharedWalletAssetDetailsViewPageConfig,
                                  );
                                },
                                child: tiles(asset)),
                          ],
                        ] else ...[
                          Container(
                            height: height / 3,
                            child: Padding(
                                padding:
                                    const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
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
                        Button(
                          LanguageEn.back,
                          notifier.getbluecolor,
                          wihitecolor,
                          onTap: () {
                            Navigator.of(context).pop();
                          },
                        ),
                        SizedBox(
                          height: height / 10,
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              if (tabLength == 2) ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 8, 0, 0),
                  child: Container(
                    child: SingleChildScrollView(
                      child: Column(
                        children: [
                          if (unclaimedAssets.length > 0) ...[
                            for (var asset in unclaimedAssets) ...[
                              GestureDetector(
                                onTap: () {
                                  setState(() {
                                    activeTabIndex = _tabController.index;
                                  });
                                  // since the original asset object
                                  // is immutable I create a new assetObj and
                                  // copy all the data into it so that
                                  // I'll be able to change the data
                                  appState.viewData![
                                      PendingAssetDetailsViewPageConfig.key] = {
                                    'assetCode': asset['assetCode'],
                                    'assetIssuer': asset['assetIssuer'],
                                    'amount': asset['amount'],
                                    'qrCode': asset['qrCode'],
                                    'imageUrl': asset['imageUrl'],
                                    'walletInfo': viewData,
                                    'rel': SharedWalletDetailsViewPageConfig.key
                                  };
                                  appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page: PendingAssetDetailsViewPageConfig,
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

  Widget tiles(asset) {
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
                      asset["assetCode"].toString().isEmpty
                          ? 'XBN'
                          : asset["assetCode"],
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
                    getBalance(
                        '${calculateFiatValue(asset["amount"], asset["usdPrice"], appState.defaultCurrency, appState)} ${appState.defaultCurrency}'),
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

  refreshData() async {
    try {
      var responseData = await fetchWalletBalance(
          signer: appState.activeWallet!.signer!,
          secretKey: appState.secretKeys[0],
          publicKey: viewData['walletPublicKey']);

      claimedAssets = responseData['assetBalances']['claimed'];
      unclaimedAssets = responseData['assetBalances']['unclaimed'];
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  // we need to check that the username entered here is a valid
  // username of an active trovo account
  Future<Map> fetchWalletBalance(
      {required String signer,
      required String secretKey,
      required String publicKey}) async {
    try {
      Map responseData = await makeGetRequest(
        uri: '/v1/shared-access/wallet-balances',
        signer: signer,
        secretKey: secretKey, // the primary wallet secret key
        publicKey: publicKey,
      );

      print('response: ${responseData}');

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }
}
