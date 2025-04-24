import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  bool hasEnoughTrov = false;
  double trovUsdPrice = 0;
  String filterValue = '';
  DateTime? filterStartDate;
  DateTime? filterEndDate;
  double filterMinAmount = 0;
  double filterMaxAmount = 0;
  double tokenizationApplicationFee = 0;
  String tokenizationApplicationFeeAsset = '';
  List<TokenizedAsset> records = [];
  bool showFilter = false;
  late List<Wallet> wallets;
  String selectedWallet = '';
  TokenizationFilterMode filterMode =
      TokenizationFilterMode.AssetTokenizationStatus;
  TokenizationStatus activeStatus = TokenizationStatus.All;
  ScrollController scrollController = new ScrollController();
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

  List<DropdownMenuItem<TokenizationStatus>> get getTokenizationStatusList {
    List<DropdownMenuItem<TokenizationStatus>> items = [];
    tokenizationStatuses.forEach((key, value) {
      items.add(DropdownMenuItem(
          child: Text(
            key,
            overflow: TextOverflow.ellipsis,
          ),
          value: value));
    });
    return items;
  }

  List<DropdownMenuItem<TokenizationFilterMode>> get getTokenizationFilterList {
    List<DropdownMenuItem<TokenizationFilterMode>> items = [];
    tokenizationFilters.forEach((key, value) {
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

  Map<String, TokenizationStatus> tokenizationStatuses = {
    'All': TokenizationStatus.All,
    'Awaiting Fee': TokenizationStatus.AwaitingFeePayment,
    'Awaiting Approval': TokenizationStatus.AwaitingMinting,
    'Primary Sales': TokenizationStatus.PrimarySales,
    'Secondary Market': TokenizationStatus.SecondaryMarket,
    'Liquidated': TokenizationStatus.Liquidated,
    'Refunded': TokenizationStatus.Refunded,
  };

  Map<String, TokenizationFilterMode> tokenizationFilters = {
    'Tokenization Status': TokenizationFilterMode.AssetTokenizationStatus,
    'Amount Range': TokenizationFilterMode.AmountRange,
    'Asset Code': TokenizationFilterMode.AssetCode,
    'Asset Name': TokenizationFilterMode.AssetName,
    'Asset Type': TokenizationFilterMode.AssetType,
    'Initiator Username': TokenizationFilterMode.InitiatorUsername,
    'Offering Type': TokenizationFilterMode.OfferingType,
    'Has Sec Approval': TokenizationFilterMode.HasSecApproval,
    'Asset Description': TokenizationFilterMode.AssetDescription,
    'Created Between': TokenizationFilterMode.DateRange,
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
    wallets = appState.userInfo!.allWallets;
    for (var asset in appState.primaryWallet.claimedAssets!) {
      if (asset.assetCode!.toUpperCase() == tokenizationApplicationFeeAsset) {
        trovUsdPrice = asset.usdPrice!;
        if (asset.amount! >= (tokenizationApplicationFee / asset.usdPrice!)) {
          hasEnoughTrov = true;
          break;
        }
      }
    }

    for (var i = 0;
        i < appState.tokenizationData['countryConfigs'].length;
        i++) {
      if (appState.tokenizationData['countryConfigs'][i]['countryCode']
              .toString()
              .toLowerCase() ==
          appState.userInfo?.countryCode?.toLowerCase()) {
        tokenizationApplicationFee = appState.tokenizationData['countryConfigs']
            [i]['tokenizationApplicationFee'];
        tokenizationApplicationFeeAsset = appState
            .tokenizationData['countryConfigs'][i]
                ['tokenizationApplicationFeeAsset']
            .toString()
            .split(':')[0];
      }
    }
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

  void startTokenizationPressed() {
    if (appState.tokenizationData.isEmpty) {
      popup(context,
          title: "error".tr(),
          message:
              "Cannot initiate this process at the moment. Please check your network, refresh this view and try again.");
      return;
    }

    if (!canCreateNewTokenization) {
      popup(context,
          title: "error".tr(),
          message:
              "You must complete the active tokenization process before starting a new one.");
      return;
    }

    if (!hasEnoughTrov) {
      popup(context,
          title: "error".tr(),
          message:
              "You must have at least ${formatNumber(tokenizationApplicationFee)} $tokenizationApplicationFeeAsset tokens (\$500 worth) in your wallet ${appState.primaryWallet.alias!.toUpperCase()} to begin a new tokenization process.");
      return;
    }

    appState.viewData = {};

    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: SetupAndComplianceViewPageConfig,
    );
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
                        if (records.length > 0) ...[
                          SizedBox(
                            height: height / 70,
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceAround,
                            children: [
                              ElevatedButton(
                                onPressed: startTokenizationPressed,
                                style: ButtonStyle(
                                  overlayColor:
                                      MaterialStateProperty.all<Color>(
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
                                        "Tokenize Another Asset",
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
                          )
                        ],
                        SizedBox(
                          height: height / 50,
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              if (showFilter) ...[
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
                              filterMode = newValue as TokenizationFilterMode;
                              showPopup(newValue);
                            });
                          },
                          getTokenizationFilterList,
                          null,
                          'Tokenization Status',
                          context,
                          null,
                        ),
                      ),
                      Expanded(
                        flex: 2,
                        child: getContent(filterMode),
                      ),
                      SizedBox(
                        width: width / 50,
                      ),
                    ],
                  ),
                ),
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
                      // records =
                      //     snapshot.data!['records'] as List<TokenizedAsset>;
                      canCreateNewTokenization = true;
                      if (records.length > 0) {
                        var assets = <Widget>[];
                        for (var i = 0; i < records.length; i++) {
                          var asset = records[i];
                          if (asset.tokenizationStatus! == 0) {
                            canCreateNewTokenization = false;
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

                              if (records[i].tokenizationStatus! <= 3) {
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
                                '${records[i].assetName?.length == 0 ? 'No name yet' : records[i].assetName} ${records[i].assetCode?.length == 0 ? '' : '(${records[i].assetCode})'}',
                                '${records[i].assetSubSector}',
                                records[i].tokenizationStatus == 1 &&
                                        records[i].vettingStatus == 0
                                    ? 'Pending Vetting'
                                    : getTokenizationStatus(
                                        records[i].tokenizationStatus!)),
                            // getTokenizationStatus(
                            //     records[i].tokenizationStatus)),
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
                                    'My Tokenized Assets',
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
                      onPressed: startTokenizationPressed,
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
      if (responseData['statusCode'] == 200) {
        await fetchTokenizationData();
        List<TokenizedAsset> assets = [];
        tokenizedAssets = responseData['data']['records'];
        if (tokenizedAssets != null) {
          for (int i = 0; i < tokenizedAssets.length; i++) {
            var a = TokenizedAsset().deserializeJson(tokenizedAssets[i]);
            a.usdPrice = 1.47;
            a.assetIssuer = a.walletToHoldAssetsNotForSale ?? '';
            a.pricePerToken = (double.parse(a.assetCurrentValue.toString()) /
                a.numberOfTokenToBeIssued!);
            assets.add(a);
          }
        }
        setState(() {
          records = assets;
        });
        return {"records": assets};
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
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
        return 'Secondary Market';
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
                tokenizationFilterTextFieldPopup(context,
                    label: 'Asset Description',
                    placeholder: 'Enter text here', onDone: (value) async {
                  setState(() {
                    filterValue = value!;
                  });
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty
                          ? "enterdescription".tr()
                          : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
                tokenizationFilterTextFieldPopup(context,
                    label: 'assetname'.tr(),
                    placeholder: 'Enter text here', onDone: (value) async {
                  setState(() {
                    filterValue = value!;
                  });
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty ? "enterassetname".tr() : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
                tokenizationFilterTextFieldPopup(context,
                    label: 'assetcode'.tr(),
                    placeholder: 'Enter text here', onDone: (value) async {
                  setState(() {
                    filterValue = value!;
                  });
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty ? "enterassetcode".tr() : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
      case TokenizationFilterMode.AssetSubSector:
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
                tokenizationFilterTextFieldPopup(context,
                    label: 'assetsubsector'.tr(),
                    placeholder: 'Enter text here', onDone: (value) async {
                  setState(() {
                    filterValue = value!;
                  });
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty
                          ? "enterassetsubsector".tr()
                          : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
                tokenizationFilterTextFieldPopup(context,
                    label: 'assettype'.tr(),
                    placeholder: 'Enter text here', onDone: (value) async {
                  setState(() {
                    filterValue = value!;
                  });
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty ? "enterassettype".tr() : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
      case TokenizationFilterMode.InitiatorUsername:
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
                tokenizationFilterTextFieldPopup(context,
                    label: 'initiatorusername'.tr(),
                    placeholder: 'Enter text here', onDone: (value) async {
                  setState(() {
                    filterValue = value!;
                  });
                  if (value != null && value.isNotEmpty) {
                    appState.setFilterQuery = "&name=${value}";
                  }
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty ? "enterusername".tr() : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
      case TokenizationFilterMode.OfferingType:
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
                tokenizationOfferingTypePopup(
                  context,
                  onSelected: (value) {
                    setState(() {
                      filterValue = value;
                    });
                    Navigator.of(context).pop(); // dismiss dialog,
                  },
                );
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      filterValue.isEmpty
                          ? "chooseofferingtype".tr()
                          : filterValue,
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: filterValue.isEmpty ? 15 : 15,
                        fontFamily: fontsemibold,
                      ),
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
                tokenizationAmountRangePopup(context,
                    minAmount: filterMinAmount, maxAmount: filterMaxAmount,
                    onDone: (minAmount, maxAmount) async {
                  setState(() {
                    filterMinAmount = minAmount;
                    filterMaxAmount = maxAmount;
                    filterValue =
                        "&amount=${filterMinAmount}%7C${filterMaxAmount}";
                  });
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    constraints: BoxConstraints(
                      maxWidth: width / 2.9,
                    ),
                    child: Text(
                      appState.filterMinAmount != null &&
                              appState.filterMaxAmount != null
                          ? "${appState.filterMinAmount} - ${appState.filterMaxAmount} "
                          : "enterrange".tr(),
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
                tokenizationCustomDateRangePopup(context,
                    initialStartDate: filterStartDate,
                    initialEndDate: filterEndDate,
                    onDone: (startDate, endDate) async {
                  setState(() {
                    filterStartDate = startDate;
                    filterEndDate = endDate;
                    filterValue =
                        "&dateBetween=${DateFormat('yyyy-MM-dd').format(startDate)}%7C${DateFormat('yyyy-MM-dd').format(endDate)}";
                  });
                });
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    width: 130,
                    child: Text(
                      filterStartDate != null && filterEndDate != null
                          ? '${DateFormat('yyyy-MM-dd').format(filterStartDate!)} : ${DateFormat('yyyy-MM-dd').format(filterEndDate!)}'
                          : "enterrange".tr(),
                      textAlign: TextAlign.start,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: appState.filterStartDate != null &&
                                  appState.filterEndDate != null
                              ? 13
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
                tokenizationStatusPopup(
                  context,
                  onSelected: (value) {
                    setState(() {
                      filterValue = value;
                    });
                    Navigator.of(context).pop(); // dismiss dialog,
                  },
                );
              },
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    getTokenizationValue(),
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

  getAmountRangeValue() {
    if (appState.filterMinAmount != null && appState.filterMaxAmount != null) {
      return "${appState.filterMinAmount} - ${appState.filterMaxAmount} ";
    }

    return "enterrange".tr();
  }

  getTokenizationValue() {
    if (filterValue.isNotEmpty) {
      switch (filterValue) {
        case '1':
          return 'Awaiting Fee';
        case '2':
          return 'Processing';
        case '5':
          return 'Primary Sales';
        case '6':
          return 'Secondary Market';
        case '7':
          return 'Liquidated';
        case '8':
          return 'Refunded';
        default:
          return 'All';
      }
    }

    return 'All';
  }

  void showPopup(TokenizationFilterMode type) {
    switch (type) {
      case TokenizationFilterMode.AssetDescription:
        tokenizationFilterTextFieldPopup(context,
            label: 'Asset Description',
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      case TokenizationFilterMode.AmountRange:
        tokenizationAmountRangePopup(context,
            minAmount: filterMinAmount,
            maxAmount: filterMaxAmount, onDone: (minAmount, maxAmount) async {
          setState(() {
            filterMinAmount = minAmount;
            filterMaxAmount = maxAmount;
            filterValue = "&amount=${filterMinAmount}%7C${filterMaxAmount}";
          });
        });
        break;
      case TokenizationFilterMode.AssetCode:
        tokenizationFilterTextFieldPopup(context,
            label: 'assetcode'.tr(),
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      case TokenizationFilterMode.DateRange:
        tokenizationCustomDateRangePopup(context,
            initialStartDate: filterStartDate,
            initialEndDate: filterEndDate, onDone: (startDate, endDate) async {
          setState(() {
            filterStartDate = startDate;
            filterEndDate = endDate;
            filterValue =
                "&dateBetween=${DateFormat('yyyy-MM-dd').format(startDate)}%7C${DateFormat('yyyy-MM-dd').format(endDate)}";
          });
        });
        break;
      case TokenizationFilterMode.AssetSector:
        tokenizationFilterTextFieldPopup(context,
            label: 'assetsector'.tr(),
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      case TokenizationFilterMode.AssetSubSector:
        tokenizationFilterTextFieldPopup(context,
            label: 'assetsubsector'.tr(),
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      case TokenizationFilterMode.AssetType:
        tokenizationFilterTextFieldPopup(context,
            label: 'assettype'.tr(),
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      case TokenizationFilterMode.InitiatorUsername:
        tokenizationFilterTextFieldPopup(context,
            label: 'initiatorusername'.tr(),
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      case TokenizationFilterMode.OfferingType:
        tokenizationOfferingTypePopup(
          context,
          onSelected: (value) {
            setState(() {
              filterValue = value;
            });
            Navigator.of(context).pop(); // dismiss dialog,
          },
        );
        break;
      case TokenizationFilterMode.AssetName:
        tokenizationFilterTextFieldPopup(context,
            label: 'assetname'.tr(),
            placeholder: 'Enter text here', onDone: (value) async {
          setState(() {
            filterValue = value!;
          });
          if (value != null && value.isNotEmpty) {
            appState.setFilterQuery = "&name=${value}";
          }
        });
        break;
      default:
        tokenizationStatusPopup(
          context,
          onSelected: (value) {
            setState(() {
              filterValue = value;
            });
            Navigator.of(context).pop(); // dismiss dialog,
          },
        );
        break;
    }
  }
}

enum TokenizationStatus {
  Open,
  AwaitingFeePayment,
  AwaitingFeePaymentConfirmation,
  AwaitingDD,
  AwaitingMinting,
  PrimarySales,
  SecondaryMarket,
  Liquidated,
  Refunded,
  DateRange,
  AmountRange,
  All,
}

enum TokenizationFilterMode {
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
