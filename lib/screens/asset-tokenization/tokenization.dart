import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
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

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getCurrencies {
    List<DropdownMenuItem<String>> currencies = [];
    appState.fiatRate.forEach((key, value) {
      currencies.add(DropdownMenuItem(
          child: Text(
            key,
            overflow: TextOverflow.ellipsis,
          ),
          value: key));
    });
    return currencies;
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
              txt: 'Asset Tokenization',
              titlecolor: notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            if (true) ...[
              Container(
                height: height / 1.3,
                child: Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Image.asset(
                        "assets/images/wallet.png",
                        height: height / 2.5,
                      ),
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
                                    fontFamily: fontbody,
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
                                  height: height / 50,
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                      SizedBox(
                        height: height / 20,
                      ),
                      Text(
                        'Coming Soon ...',
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                    ],
                  ),
                ),
              )
            ] else ...[
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
                          "welcometoassettokenization".tr(),
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
                          "welcometoassettokenization2".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 13,
                            fontFamily: fontbody,
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
                            SmallButton(
                              'Tokenize asset',
                              notifier.getbluewhitecolor!,
                              notifier.getwihitecolor,
                              onTap: () {
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: WalletPreparationViewPageConfig,
                                );
                              },
                            ),
                            SmallButtonOutlined(
                              'See my assets',
                              notifier.isDark
                                  ? notifier.getbluecolor90
                                  : notifier.getaddsubwalletgrey,
                              notifier.getbluewhitecolor,
                              onTap: () {
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: TokenizedAssetsListViewPageConfig,
                                );
                              },
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
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
                child: TabBar(
                  controller: tabController,
                  labelColor: notifier.getbluewhitecolor,
                  indicatorColor: notifier.getbluewhitecolor,
                  labelStyle: TextStyle(
                    fontSize: 15.sp,
                    fontFamily: fontsemibold,
                  ),
                  tabs: [
                    Tab(
                      height: 20,
                      text: 'Primary Offering',
                    ),
                    Tab(
                      height: 20,
                      text: 'All Assets',
                    ),
                  ],
                ),
              ),
              Container(
                height: height / 2.25,
                child: TabBarView(controller: tabController, children: [
                  SingleChildScrollView(
                    child: Column(
                      children: [
                        for (var i = 10; i >= 0; i--) ...[
                          GestureDetector(
                            onTap: () {
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: TokenizedAssetDetailViewPageConfig,
                              );
                            },
                            child: assetTile(
                                '', 'ASSET $i', 'Property', i % 2 == 0),
                          ),
                        ],
                        SizedBox(height: height / 20),
                      ],
                    ),
                  ),
                  SingleChildScrollView(
                    child: Column(
                      children: [
                        for (var i = 0; i <= 10; i++) ...[
                          GestureDetector(
                              onTap: () {
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: TokenizedAssetDetailViewPageConfig,
                                );
                              },
                              child: assetTile(
                                  '', 'ASSET $i', 'Property', i % 2 == 0)),
                        ],
                        SizedBox(height: height / 20),
                      ],
                    ),
                  ),
                ]),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget assetTile(String imageUrl, String name, String type, isSubscribed) {
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
                  Text(
                    name,
                    style: TextStyle(
                      fontSize: 15,
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
                      onDone: () {},
                    )
                  : showSubscribePopup(
                      context,
                      onDone: () {},
                      dropdownItems: getStandardWallets,
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
                        fontSize: 12,
                        color: notifier.getwihitecolor),
                  ),
                  Icon(
                      isSubscribed
                          ? Icons.check_circle
                          : Icons.add_circle_rounded,
                      size: 20,
                      color: notifier.getwihitecolor),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
