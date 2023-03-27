import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class MyAssetTokenDetails extends StatefulWidget {
  const MyAssetTokenDetails({Key? key}) : super(key: key);

  @override
  State<MyAssetTokenDetails> createState() => _MyAssetTokenDetails();
}

class _MyAssetTokenDetails extends State<MyAssetTokenDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late TabController tabController;
  bool localHideBalance = false;
  final Authenticator _authenticator = Authenticator();

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
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
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              '',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            assetInfo(
              'ANMF',
              '12.4304324',
              notifier.getbluewhitecolor,
              wihitecolor,
            ),
            SizedBox(
              height: height / 50,
            ),
            // SizedBox(height: height / 50),
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
                      Padding(
                        padding: const EdgeInsets.all(10.0),
                        child: Text(
                          'Animal Farm  tokens are fractional tokens that represent part ownership (via investment) of our Animal farm at Tudun wada, Kano State.',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 14,
                            height: 1.4,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
            SizedBox(height: height / 70),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Column(
                  children: [
                    TextButton(
                      onPressed: () {},
                      child: Column(
                        children: [
                          Icon(
                            CupertinoIcons.cart_fill,
                            size: 25,
                            color: notifier.getbluewhitecolor,
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                          Text(
                            'Sell Asset',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                Column(
                  children: [
                    TextButton(
                      onPressed: () {},
                      child: Column(
                        children: [
                          Image.asset(
                            "assets/images/swap.png",
                            height: height / 40,
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                          Text(
                            'Transfer Asset',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                )
              ],
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
              child: TabBar(
                controller: tabController,
                labelColor: notifier.getbluewhitecolor,
                indicatorColor: notifier.getbluewhitecolor,
                labelStyle: TextStyle(
                  fontSize: 13.sp,
                  fontFamily: fontsemibold,
                ),
                tabs: [
                  Tab(
                    height: 20,
                    text: 'Details of Asset',
                  ),
                  Tab(
                    height: 20,
                    text: 'Proceeds History',
                  ),
                ],
              ),
            ),
            Container(
              height: height / 2.0,
              child: TabBarView(controller: tabController, children: [
                SingleChildScrollView(
                  child: Column(
                    children: [
                      infoTile(
                        notifier,
                        'Asset Code',
                        'ANMF',
                      ),
                      infoTile(
                        notifier,
                        'Asset Category',
                        'Agriculture',
                      ),
                      infoTile(
                        notifier,
                        'Asset Country',
                        'Nigeria',
                      ),
                      infoTile(
                        notifier,
                        'Asset Location Address',
                        'No. 10 Maitama, Abuja',
                      ),
                      infoTile(
                        notifier,
                        'Asset Issuer',
                        'Atlantis Developers',
                      ),
                      infoTile(
                        notifier,
                        'Asset Issuer Website',
                        'www.anmf.com',
                      ),
                      infoTile(
                        notifier,
                        'Asset Token Total Supply',
                        '1000',
                      ),
                      infoTile(
                        notifier,
                        'Proceed Payout Cycle',
                        'Monthly',
                      ),
                      infoTile(
                        notifier,
                        'Payout Method',
                        'cNGN',
                      ),
                      SizedBox(height: height / 20),
                    ],
                  ),
                ),
                SingleChildScrollView(
                  child: Column(
                    children: [
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '35 cNGN',
                          'Paid on 02/03/2022',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '35 cNGN',
                          'Paid on 02/02/2022',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '35 cNGN',
                          'Paid on 02/01/2022',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '35 cNGN',
                          'Paid on 02/12/2021',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '35 cNGN',
                          'Paid on 02/11/2021',
                        ),
                      ),
                      SizedBox(height: height / 20),
                    ],
                  ),
                ),
              ]),
            ),
          ],
        ),
      ),
    );
  }

  Widget assetTile(String name, String type) {
    return Card(
      elevation: notifier.isDark ? 0 : 3,
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
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    truncate(name, length: 14),
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      type,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
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

  Widget assetInfo(
    String assetCode,
    String balance,
    Color backColor,
    Color foreColor,
  ) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: backColor,
        ),
        child: Stack(
          alignment: AlignmentDirectional.centerEnd,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(
                      vertical: 25.0, horizontal: 15),
                  child: Image.asset(
                    'assets/images/trovo_white.png',
                    height: height / 20,
                  ),
                ),
              ],
            ),
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  Container(
                    width: width / 2,
                    child: Text(
                      assetCode,
                      style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w600,
                          color: foreColor,
                          fontFamily: fontsemibold),
                    ),
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  Row(
                    children: [
                      Text(
                        'Animal Farm Token ',
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w400,
                          color: foreColor,
                          fontFamily: fontbody,
                        ),
                      ),
                      SizedBox(
                        width: 15,
                      ),
                      GestureDetector(
                        onTap: () {
                          if (localHideBalance) {
                            authenticateAndToggle();
                          } else
                            toggleHideBalance();
                        },
                        child: Icon(
                          getIcon(),
                          size: 20,
                          color: foreColor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 98.0,
                  ),
                  Container(
                    width: width / 1.8,
                    child: Text(
                      getBalance(balance),
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: foreColor,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  IconData getIcon() {
    IconData icon;
    if (appState.hideBalances) icon = CupertinoIcons.eye;

    if (localHideBalance)
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
      toggleHideBalance();
    });
  }

  void toggleBiometrics() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        toggleHideBalance();
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  toggleHideBalance() {
    setState(() {
      localHideBalance = !localHideBalance;
    });
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
