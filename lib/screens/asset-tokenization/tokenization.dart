import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
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
  bool hasInitiatorAccess = false;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
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

  var listOfAssets = <Map<String, String>>[
    {
      "imageUrl": "",
      "assetName": "ATLANTIS 1",
      "assetClass": "Property",
      "status": "Approved"
    },
    {
      "imageUrl": "",
      "assetName": "Cascadia Growth",
      "assetClass": "Property",
      "status": "Approved"
    },
    {
      "imageUrl": "",
      "assetName": "Titan Properties",
      "assetClass": "Property",
      "status": "Pending"
    },
    {
      "imageUrl": "",
      "assetName": "Kings Home",
      "assetClass": "Property",
      "status": "Pending"
    },
    {
      "imageUrl": "",
      "assetName": "Heart Realty",
      "assetClass": "Property",
      "status": "Rejected"
    },
    {
      "imageUrl": "",
      "assetName": "Infinity Productions",
      "assetClass": "Property",
      "status": "Approved",
    },
  ];

  List<String> listMode = [
    'Tokenized Assets',
    'Assets',
  ];

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    tabController = TabController(length: 2, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
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
                              // appState.currentAction = PageAction(
                              //   state: PageState.addPage,
                              //   page: WalletPreparationViewPageConfig,
                              // );
                              setState(() {
                                hasInitiatorAccess = !hasInitiatorAccess;
                              });
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
            SizedBox(height: height / 70),
            ReadyToTokenizeView()
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
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 70,
                  ),
                  ElevatedButton(
                    onPressed: () async {
                      hasInitiatorAccess
                          ? showCreateTokenizationWalletPopup(context)
                          : appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: WalletPreparationViewPageConfig,
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
                  Container(
                    width: width / 3.4,
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
