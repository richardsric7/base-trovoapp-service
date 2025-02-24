import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TokenizationWelcome extends StatefulWidget {
  const TokenizationWelcome({Key? key}) : super(key: key);

  @override
  State<TokenizationWelcome> createState() => _TokenizationWelcomeState();
}

class _TokenizationWelcomeState extends State<TokenizationWelcome>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  late TabController tabController;
  late Future<Map> listOfTokenizations;
  bool canCreateNewTokenization = true;
  bool showFilter = false;
  late List<Wallet> wallets;
  String selectedWallet = '';
  TokenizationFilterMode filterType = TokenizationFilterMode.All;
  TokenizedAssetListMode listMode = TokenizedAssetListMode.All;
  ScrollController scrollController = new ScrollController();
  // bool hasInitiatorAccess = false;
  late var tokenizedAssets;
  late RefreshController _refreshController;
  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  bool checkHasInitiatorAccess() {
    for (var wallet in appState.userInfo!.getMintingWallets) {
      if (wallet.isSharedWalletAndCanInitiate) {
        return true;
      }
    }
    return false;
  }

  List<DropdownMenuItem<TokenizedAssetListMode>> get getItems {
    List<DropdownMenuItem<TokenizedAssetListMode>> items = [];
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

  Map<String, TokenizedAssetListMode> listModes = {
    'All': TokenizedAssetListMode.All,
    'Awaiting Fee Payment': TokenizedAssetListMode.AwaitingFeePayment,
    'Awaiting Fee Payment Confirmation':
        TokenizedAssetListMode.AwaitingFeePaymentConfirmation,
    'Awaiting DD': TokenizedAssetListMode.AwaitingDD,
    'Awaiting Approval/Minting': TokenizedAssetListMode.AwaitingMinting,
    'Primary Sales': TokenizedAssetListMode.PrimarySales,
    'Secondary Sales': TokenizedAssetListMode.SecondarySales,
    'Liquidated': TokenizedAssetListMode.Liquidated,
    'Refunded': TokenizedAssetListMode.Refunded,
  };

  List<DropdownMenuItem<String>> walletDropdownItems(bool isSelected) {
    var walletsList = <DropdownMenuItem<String>>[];

    wallets.forEach((wallet) {
      walletsList.add(
        DropdownMenuItem(
          child: Row(
            children: [
              Container(
                constraints: isSelected
                    ? BoxConstraints(maxWidth: width / 4)
                    : BoxConstraints(maxWidth: width / 2.5),
                child: Text(
                  wallet.alias!,
                  overflow:
                      isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
                ),
              ),
              if (wallet.isSharedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.people_outline,
                  size: 17,
                  color: notifier.getbluewhitecolor,
                )
              ],
              if (!isSelected && wallet.publicKey == selectedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.check,
                  size: 18,
                  color: notifier.getbluecolor,
                )
              ],
            ],
          ),
          value: wallet.publicKey,
        ),
      );
    });

    return walletsList;
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    tabController = TabController(length: 2, vsync: this);
    appState = Provider.of<DataProvider>(context, listen: false);
    _refreshController = RefreshController(initialRefresh: false);
    listOfTokenizations = fetchTokenizationList();
    // hasInitiatorAccess = checkHasInitiatorAccess();
    wallets = appState.userInfo!.allWallets;
  }

  void refreshData() async {
    try {
      listOfTokenizations = fetchTokenizationList();
      _refreshController.refreshCompleted();
      appState.updateListeners();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
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
          child: Column(
            children: [
              CustomAppBarWithoutLeading(
                context,
                notifier.getwihitecolor,
                scaffoldKey: key,
                showMenu: true,
                txt: '',
                titlecolor: notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
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
                          height: height / 70,
                        ),
                        Text(
                          "welcometoassettokenization2".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        SizedBox(
                          height: height / 70,
                        ),
                        Text(
                          "welcometoassettokenization3".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 13,
                            height: 1.4,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        SizedBox(
                          height: height / 70,
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            ElevatedButton(
                              onPressed: () async {
                                if (appState.tokenizationData.isEmpty) {
                                  popup(context,
                                      title: "error".tr(),
                                      message:
                                          "Cannot initiate this process at the moment. Please check your network, refresh this view and try again.");
                                  return;
                                }

                                // if (!canCreateNewTokenization) {
                                //   popup(context,
                                //       title: "error".tr(),
                                //       message:
                                //           "You must complete the active tokenization process before starting a new one.");
                                //   return;
                                // }

                                if (!canCreateNewTokenization) {
                                  popup(context,
                                      title: "error".tr(),
                                      message:
                                          "You must have at least \$500 worth of TROV on any of your wallets to begin a new tokenization process.");
                                  return;
                                }

                                appState.viewData = {};

                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: SetupAndComplianceViewPageConfig,
                                );
                              },
                              style: ButtonStyle(
                                overlayColor: MaterialStateProperty.all<Color>(
                                    notifier.getsplashgrey),
                                backgroundColor:
                                    MaterialStateProperty.all<Color>(
                                  notifier.isDark
                                      ? notifier.getbluecolor90
                                      : notifier.getaddsubwalletgrey,
                                ),
                                side: MaterialStateProperty.all(
                                  BorderSide(
                                      color: notifier.getbluewhitecolor,
                                      width: 1,
                                      style: BorderStyle.solid),
                                ),
                                shape: MaterialStateProperty.all<
                                    RoundedRectangleBorder>(
                                  const RoundedRectangleBorder(
                                    borderRadius: BorderRadius.all(
                                      Radius.circular(10),
                                    ),
                                  ),
                                ),
                              ),
                              child: Container(
                                width: width / 1.5,
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Icon(
                                      Icons.add_circle_rounded,
                                      size: 20,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                    SizedBox(
                                      width: 4,
                                    ),
                                    Text(
                                      "tokenizeasset".tr(),
                                      style: TextStyle(
                                          fontFamily: fontsemibold,
                                          fontSize: 12,
                                          color: notifier.getbluewhitecolor),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ],
                        ),
                        SizedBox(
                          height: height / 50,
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              if (showFilter) ...[
                Container(
                  width: width,
                  child: Row(
                    children: [
                      SizedBox(
                        width: width / 50,
                      ),
                      // hide the dropdown when we view this page from shared
                      // wallet
                      Expanded(
                        flex: 2,
                        child: dropdown(
                          (newValue) async {
                            appState.filterAsset = "*|*";
                            showLoader(context);
                            appState.limit = 20;
                            appState.totalRecords = 0;
                            appState.currentPage = 1;
                            hideLoader(context);

                            if (mounted) {
                              setState(() {});
                            }
                          },
                          walletDropdownItems(false),
                          null,
                          null,
                          context,
                          (context) {
                            return walletDropdownItems(true);
                          },
                        ),
                      ),
                      Expanded(
                        flex: 2,
                        child: dropdown(
                          (newValue) async {
                            appState.filterAsset = "*|*";
                            showLoader(context);
                            appState.limit = 20;
                            appState.totalRecords = 0;
                            appState.currentPage = 1;
                            hideLoader(context);

                            if (mounted) {
                              setState(() {});
                            }
                          },
                          walletDropdownItems(false),
                          null,
                          null,
                          context,
                          (context) {
                            return walletDropdownItems(true);
                          },
                        ),
                      ),
                      SizedBox(
                        width: width / 50,
                      ),
                    ],
                  ),
                ),
                SizedBox(height: height / 50),
                Container(
                  width: width,
                  child: Row(
                    children: [
                      SizedBox(
                        width: width / 50,
                      ),
                      Expanded(
                        flex: 2,
                        child: dropdown(
                          (newValue) async {
                            setState(() {
                              listMode = newValue as TokenizedAssetListMode;
                            });
                          },
                          getItems,
                          null,
                          'All',
                          context,
                          null,
                        ),
                      ),
                      Expanded(
                        flex: 2,
                        child: getContent(filterType),
                      ),
                      SizedBox(
                        width: width / 50,
                      ),
                    ],
                  ),
                ),
                SizedBox(height: height / 35),
              ],
              FutureBuilder<Map>(
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
                  } else if (snapshot.connectionState == ConnectionState.done) {
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
                      var records = snapshot.data!['records'];
                      if (records.length > 0) {
                        var assets = <Widget>[];

                        for (var i = 0; i < records.length; i++) {
                          var asset = records[i] as TokenizedAsset;
                          print(
                              'Can create new tokenization: ${asset.assetCode} ${asset.tokenizationStatus} ');
                          if (asset.tokenizationStatus! <= 2) {
                            canCreateNewTokenization = false;
                            print(
                                'Can create new tokenization: $canCreateNewTokenization');
                          }
                          assets.add(GestureDetector(
                            onTap: () async {
                              appState.viewData = tokenizedAssets[i];

                              if (records[i].tokenizationStatus == 0) {
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: SetupAndComplianceViewPageConfig,
                                );

                                return;
                              }

                              if (records[i].tokenizationStatus == 1) {
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page:
                                      ConfirmTokenizationDetailsViewPageConfig,
                                );
                                return;
                              }

                              appState.tokenizedAsset = records[i];
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: AssetDashboardViewPageConfig,
                              );
                            },
                            child: assetTile(
                              records[i].assetLogo ?? '',
                              '${records[i].assetName.length == 0 ? 'No name yet' : records[i].assetName} ${records[i].assetCode.length == 0 ? '' : '(${records[i].assetCode})'}',
                              '${records[i].assetSubSector}',
                              getTokenizationStatus(
                                  records[i].tokenizationStatus),
                            ),
                          ));
                        }

                        return Column(
                          children: [
                            SizedBox(
                              height: height / 50,
                            ),
                            Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20),
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    'Tokenized Assets',
                                    style: TextStyle(
                                      fontSize: 17,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            SizedBox(height: height / 90),
                            Column(
                              children: [
                                ...assets,
                                SizedBox(height: height / 20),
                              ],
                            ),
                          ],
                        );
                      } else {
                        return ReadyToTokenizeView();
                      }
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
              SizedBox(height: height / 70),
            ],
          ),
        ),
      ),
    );
  }

  adjustScrollPosition() {
    if (scrollController.hasClients)
      scrollController.jumpTo(scrollController.position.minScrollExtent);
  }

  Widget ReadyToTokenizeView() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Card(
            shadowColor: Colors.black,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(15.0),
            ),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
            child: Center(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: Column(
                  children: [
                    SizedBox(
                      height: height / 50,
                    ),
                    Image.asset(
                      'assets/images/tokenize.png',
                      // height: 50,
                      width: 250,
                    ),
                    SizedBox(
                      height: height / 60,
                    ),
                    Text(
                      "welcometoassettokenization4".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 13,
                        // fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                    Text(
                      "welcometoassettokenization5".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 13,
                        height: 1.4,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                    Text(
                      "welcometoassettokenization6".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 13,
                        height: 1.4,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                    Text(
                      "welcometoassettokenization7".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 13,
                        height: 1.4,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 70,
                    ),
                    ElevatedButton(
                      onPressed: () async {
                        appState.viewData = {};

                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: SetupAndComplianceViewPageConfig,
                        );
                      },
                      style: ButtonStyle(
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluewhitecolor),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getbluewhitecolor,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Container(
                        width: width / 1.5,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.add_circle_rounded,
                              size: 20,
                              color: notifier.getwihitecolor,
                            ),
                            SizedBox(
                              width: 4,
                            ),
                            Text(
                              "proceedtokenizeasset".tr(),
                              style: TextStyle(
                                  fontFamily: fontsemibold,
                                  fontSize: 12,
                                  color: notifier.getwihitecolor),
                            ),
                          ],
                        ),
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
        ),
      ],
    );
  }

  Widget assetTile(String imageUrl, String name, String type, String status) {
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
              if (imageUrl.length == 0) ...[
                Image.asset(
                  'assets/images/trovo.png',
                  height: 35,
                  width: 35,
                ),
              ] else ...[
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
              ],
              SizedBox(width: 20),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    width: width / 2.7,
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
                    child: Container(
                      width: width / 2.7,
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
          trailing: Container(
            width: width / 4,
            child: Card(
              shadowColor: Colors.black,
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(5.0),
                  side: BorderSide(
                    color: notifier.getbluewhitecolor,
                    width: 1,
                  )),
              color: notifier.isDark
                  ? notifier.getbluecolor90
                  : notifier.getaddsubwalletgrey,
              child: Padding(
                padding: const EdgeInsets.all(5.0),
                child: Text(
                  status,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                    overflow: TextOverflow.visible,
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
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
      inspect(responseData['data']);
      appState.tokenizationData = responseData['data'];
    }
  }

  Future<Map> fetchTokenizationList() async {
    try {
      var uri = '/v1/tokenization/list';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        await fetchTokenizationData();
        List<TokenizedAsset> assets = [];
        tokenizedAssets = responseData['data']['records'];
        await inspect(tokenizedAssets);
        if (tokenizedAssets != null) {
          for (int i = 0; i < tokenizedAssets.length; i++) {
            inspect(tokenizedAssets[i]);
            var a = TokenizedAsset().deserializeJson(tokenizedAssets[i]);
            a.usdPrice = 1.47;
            a.assetIssuer = a.walletToHoldAssetsNotForSale ?? '';
            a.pricePerToken = (double.parse(a.assetCurrentValue.toString()) /
                a.numberOfTokenToBeIssued!);
            assets.add(a);
          }
        }
        return {"records": assets};
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      print('error');
      print(e);
      return Future.error('Error! ${e}');
    }
  }

  String getTokenizationStatus(int status) {
    switch (status) {
      case 0:
        return 'Continue';
      case 1:
        return 'Awaiting Fee';
      case 4:
        return 'Approved';
      case 5:
        return 'Primary Sales';
      case 6:
        return 'Secondary Sales';
      case 7:
        return 'Liquidated';
      case 8:
        return 'Refunded';
      default:
        return 'Processing';
    }
  }

  Widget getContent(TokenizationFilterMode type) {
    switch (type) {
      case TokenizationFilterMode.AssetDescription:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // textFieldPopup(context, rel: HistoryFilterType.Username,
                //     onDone: (value) async {
                //   appState.setFilterUsername = value;
                //   if (value != null && value.isNotEmpty) {
                //     appState.setFilterQuery = "&name=${value}";
                //     await appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //   }
                // });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      appState.filterUsername == null
                          ? "enterusername2".tr()
                          : appState.filterUsername!,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterUsername != null ? 12 : 15,
                          fontFamily: fontsemibold),
                    ),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
      case TokenizationFilterMode.AssetName:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // textFieldPopup(context, rel: HistoryFilterType.Memo,
                //     onDone: (value) async {
                //   appState.setFilterMemo = value;
                //   if (value != null && value.isNotEmpty) {
                //     appState.setFilterQuery = "&memo=${value}";
                //     await appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //   }
                // });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      appState.filterMemo == null
                          ? "entermemo".tr()
                          : truncate(appState.filterMemo!, length: 30),
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterMemo != null ? 12 : 15,
                          fontFamily: fontsemibold),
                    ),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
      case TokenizationFilterMode.AssetCode:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // textFieldPopup(context, rel: TokenizationFilterMode.AssetType,
                //     onDone: (value) async {
                //   if (value != null && value.toString().isNotEmpty) {
                //     appState.setFilterFromPublicKey = value;
                //     appState.setFilterQuery = "&fromPublicKey=$value";
                //     await appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //   }
                // });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      'getTruncatedPublicKey(appState.filterFromPublicKey)',
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize:
                              appState.filterFromPublicKey != null ? 12 : 15,
                          fontFamily: fontsemibold),
                    ),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
      case TokenizationFilterMode.AssetType:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // textFieldPopup(context, rel: HistoryFilterType.ToPublicKey,
                //     onDone: (value) async {
                //   if (value != null && value.toString().isNotEmpty) {
                //     appState.setFilterToPublicKey = value;
                //     appState.setFilterQuery = "&toPublicKey=$value";
                //     await appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //   }
                // });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      'getTruncatedPublicKey(appState.filterToPublicKey)',
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize:
                              appState.filterToPublicKey != null ? 12 : 15,
                          fontFamily: fontsemibold),
                    ),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
      case TokenizationFilterMode.AmountRange:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // amountRangePopup(context, onDone: () async {
                //   if (appState.filterMinAmount != null &&
                //       appState.filterMaxAmount != null) {
                //     appState.setFilterQuery =
                //         "&amount=${appState.filterMinAmount}%7C${appState.filterMaxAmount}";
                //     await appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //   }
                // });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      'truncate(getAmountRangeValue(), length: 30)',
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterMinAmount != null &&
                                  appState.filterMaxAmount != null
                              ? 12
                              : 15,
                          fontFamily: fontsemibold),
                    ),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
      case TokenizationFilterMode.DateRange:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // customDateRangePopup(context, onDone: () async {
                //   appState.setFilterQuery =
                //       "&dateBetween=${DateFormat('yyyy-MM-dd').format(appState.filterStartDate!)}%7C${DateFormat('yyyy-MM-dd').format(appState.filterEndDate!)}";
                //   await appState.getHistory(
                //     context,
                //     selectedWallet,
                //     onDone: () => adjustScrollPosition(),
                //   );
                // });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'getDateRangeValue()',
                    textAlign: TextAlign.start,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: appState.filterStartDate != null &&
                                appState.filterEndDate != null
                            ? 13
                            : 15,
                        fontFamily: fontsemibold),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
      // HistoryFilterType.TransactionDirection
      default:
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: TextButton(
              onPressed: () {
                // transactionTypePopup(
                //   context,
                //   onAllSelected: () {
                //     appState.setFilterQuery = "";
                //     appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //     Navigator.of(context).pop(); // dismiss dialog,
                //   },
                //   onPaymentSelected: () {
                //     appState.setFilterQuery = "&transactionType=payment";
                //     appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //     Navigator.of(context).pop(); // dismiss dialog,
                //   },
                //   onSwapSelected: () {
                //     appState.setFilterQuery = "&transactionType=swap";
                //     appState.getHistory(
                //       context,
                //       selectedWallet,
                //       onDone: () => adjustScrollPosition(),
                //     );
                //     Navigator.of(context).pop(); // dismiss dialog,
                //   },
                // );
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'getTransactionDirectionValue()',
                    textAlign: TextAlign.start,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: appState.filterStartDate != null &&
                                appState.filterEndDate != null
                            ? 13
                            : 15,
                        fontFamily: fontsemibold),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                ],
              ),
            ),
          ),
        );
    }
  }
}

enum TokenizedAssetListMode {
  Open,
  AwaitingFeePayment,
  AwaitingFeePaymentConfirmation,
  AwaitingDD,
  AwaitingMinting,
  PrimarySales,
  SecondarySales,
  Liquidated,
  Refunded,
  DateRange,
  AmountRange,
  All,
}

enum TokenizationFilterMode {
  All,
  DateRange,
  AmountRange,
  AssetDescription,
  AssetType,
  AssetName,
  AssetCode,
  AssetSector,
  AssetSubSector,
  InitiatorUsername,
  OfferingType,
  HasSecApproval,
  CreatedBetween,
  AssetTokenizationStatus,
}
