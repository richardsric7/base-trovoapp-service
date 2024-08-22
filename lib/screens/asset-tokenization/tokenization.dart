import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/storage/store.dart';
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
  bool hasInitiatorAccess = false;
  late var savedAssets;
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

  List<DropdownMenuItem<String>> get getItems {
    List<DropdownMenuItem<String>> items = [];
    listMode.forEach((key) {
      items.add(DropdownMenuItem(
          child: Text(
            key,
            overflow: TextOverflow.ellipsis,
          ),
          value: key));
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

  List<String> listMode = [
    'Tokenized Assets',
    'Assets',
  ];

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    tabController = TabController(length: 2, vsync: this);
    appState = Provider.of<DataProvider>(context, listen: false);
    listOfTokenizations = fetchTokenizationList();
    hasInitiatorAccess = checkHasInitiatorAccess();
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
      body: SingleChildScrollView(
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
            // SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Card(
                shadowColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: notifier.isDark
                    ? notifier.getaddsubwalletgrey
                    : notifier.getbluecolor90,
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
                          color: notifier.getwihitecolor,
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
                          color: notifier.getwihitecolor,
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

                              hasInitiatorAccess
                                  ? appState.currentAction = PageAction(
                                      state: PageState.addPage,
                                      page: WalletPreparationViewPageConfig,
                                    )
                                  : showCreateTokenizationWalletPopup(context);
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
                                    color: notifier.getwihitecolor,
                                  ),
                                  SizedBox(
                                    width: 4,
                                  ),
                                  Text(
                                    "tokenizeasset".tr(),
                                    style: TextStyle(
                                        fontFamily: fontsemibold,
                                        fontSize: 12,
                                        color: notifier.getwihitecolor),
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
                                  listOfTokenizations = fetchTokenizationList();
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
                      return Column(
                        children: [
                          SizedBox(
                            height: height / 50,
                          ),
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.spaceBetween,
                              children: [
                                Text(
                                  'Tokenized Assets',
                                  style: TextStyle(
                                    fontSize: 17,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluecolor,
                                  ),
                                ),
                                Container(
                                  width: width / 2.5,
                                  child: dropdown(
                                    (value) {},
                                    getItems,
                                    null,
                                    getItems.first.value,
                                    context,
                                    null,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Row(
                            children: [
                              Padding(
                                padding:
                                    const EdgeInsets.symmetric(horizontal: 20),
                                child: Text(
                                  'Sort by:',
                                  style: TextStyle(
                                    fontSize: 13,
                                    fontFamily: fontbody,
                                    color: notifier.getbluecolor,
                                  ),
                                ),
                              ),
                              Text(
                                'All',
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontsemibold,
                                  color: notifier.getbluecolor,
                                ),
                              ),
                            ],
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Container(
                            height: height / 2.1,
                            child: SingleChildScrollView(
                              child: Column(
                                children: [
                                  for (var i = 0; i < records.length; i++) ...[
                                    GestureDetector(
                                      onTap: () async {
                                        appState.viewData = savedAssets[i];
                                        appState.setActiveTokenizationWalletPublicKey =
                                            appState.viewData![
                                                'issuingWalletPublicKey'];
                                        appState.setActiveDistributionWalletPublicKey =
                                            appState.viewData![
                                                'marketMakingWallet'];

                                        if (records[i].tokenizationStatus ==
                                            null) {
                                          appState.currentAction = PageAction(
                                            state: PageState.addPage,
                                            page:
                                                SetupAndComplianceViewPageConfig,
                                          );

                                          return;
                                        }

                                        if (records[i].tokenizationStatus ==
                                            1) {
                                          appState.tokenizedAsset = records[i];
                                          appState.currentAction = PageAction(
                                            state: PageState.addPage,
                                            page: AssetDashboardViewPageConfig,
                                          );
                                          return;
                                        }

                                        var list = [];

                                        for (var j = 0;
                                            j < savedAssets.length;
                                            j++) {
                                          var data = Map.from(savedAssets[j]);
                                          if (data['tokenizationStatus'] !=
                                              null) {
                                            if (data['id'] == records[i].id) {
                                              data['tokenizationStatus'] = 1;
                                            }
                                            list.add(data);
                                          }
                                        }

                                        await StoreData()
                                            .storeDeleteItem('tokenizedAsset');
                                        await StoreData().storeInsertData(
                                            'tokenizedAsset', list);

                                        appState.currentAction = PageAction(
                                          state: PageState.addPage,
                                          page:
                                              ConfirmTokenizationDetailsViewPageConfig,
                                        );
                                      },
                                      child: assetTile(
                                        records[i].assetLogo ?? '',
                                        '${records[i].assetName.length == 0 ? 'No name' : records[i].assetName} (${records[i].assetCode.length == 0 ? 'Nill' : records[i].assetCode})',
                                        '${records[i].assetSubSector}',
                                        records[i].tokenizationStatus == null
                                            ? 'Continue'
                                            : records[i].tokenizationStatus == 0
                                                ? 'Pending'
                                                : records[i].tokenizationStatus ==
                                                        1
                                                    ? 'Approved'
                                                    : 'Rejected',
                                      ),
                                    ),
                                  ],
                                  SizedBox(height: height / 20),
                                ],
                              ),
                            ),
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
    );
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
                ? notifier.getaddsubwalletgrey
                : notifier.getbluecolor50,
            child: Center(
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
                      hasInitiatorAccess
                          ? appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: WalletPreparationViewPageConfig,
                            )
                          : showCreateTokenizationWalletPopup(context);
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
                      shape: MaterialStateProperty.all<RoundedRectangleBorder>(
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
                Image.memory(
                  base64Decode(imageUrl),
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
            child: Card(
              shadowColor: Colors.black,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(5.0),
              ),
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
                    color: getStatusColor(status),
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
    // await fetchTokenizationData();
    // savedAssets = await StoreData().storeGetData('tokenizedAsset');
    // List<TokenizedAsset> tokenizedAssets = [];
    // if (savedAssets != null) {
    //   for (int i = 0; i < savedAssets.length; i++) {
    //     print(savedAssets[i]);
    //     var a = TokenizedAsset().deserializeJson(savedAssets[i]);
    //     a.usdPrice = 1.47;
    //     a.assetIssuer = a.walletToHoldAssetsNotForSale ?? '';
    //     a.pricePerToken = (double.parse(a.assetCurrentValue.toString()) /
    //         a.numberOfTokenToBeIssued!);
    //     tokenizedAssets.add(a);
    //   }
    // }
    // return {"records": tokenizedAssets};
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
        List<TokenizedAsset> tokenizedAssets = [];
        savedAssets = responseData['data']['records'];
        await inspect(savedAssets);
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
        return {"records": tokenizedAssets};
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      print('error');
      print(e);
      return Future.error('Error! ${e}');
    }
  }

  Color getStatusColor(String status) {
    switch (status.toLowerCase()) {
      case 'rejected':
        return Colors.red;
      case 'approved':
        return notifier.getgreencolor;
      default: // pending
        return notifier.getbluecolor;
    }
  }
}
