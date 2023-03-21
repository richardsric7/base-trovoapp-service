import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
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
            Text(
              'Animal Farm ',
              style: TextStyle(
                fontSize: 20,
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(
              height: height / 70,
            ),
            Text(
              'Animal Farm token',
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
              '3 Tokens',
              style: TextStyle(
                fontSize: 15,
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Column(
                  children: [
                    TextButton(
                      onPressed: () {
                        showSubscribePopup(context,
                            onDone: () {}, dropdownItems: getStandardWallets);
                      },
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
                      onPressed: () {
                        showBuyTokenPopup(context, onDone: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: BuyTokensViewPageConfig,
                          );
                        }, dropdownItems: getStandardWallets);
                      },
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
                      SizedBox(
                        height: height / 70,
                      ),
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
                      SizedBox(
                        height: height / 70,
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
                        'ATLANTIS 1',
                      ),
                      infoTile(
                        notifier,
                        'Category',
                        'Real Estate',
                      ),
                      infoTile(
                        notifier,
                        'Country',
                        'Nigeria',
                      ),
                      infoTile(
                        notifier,
                        'Address',
                        'No. 10 Maitama, Abuja',
                      ),
                      infoTile(
                        notifier,
                        'Issuer',
                        'Atlantis Developers',
                      ),
                      infoTile(
                        notifier,
                        'Issuer Website',
                        'www.atlantis.com',
                      ),
                      infoTile(
                        notifier,
                        'Total Supply',
                        '1000',
                      ),
                      infoTile(
                        notifier,
                        'Quantity Purchased',
                        '400',
                      ),
                      infoTile(
                        notifier,
                        'Total Subscribed',
                        '2,000',
                      ),
                      infoTile(
                        notifier,
                        'Price per Asset',
                        '100 TROV',
                      ),
                      infoTile(
                        notifier,
                        'Funding Method',
                        'TROV',
                      ),
                      infoTile(
                        notifier,
                        'Sales Window',
                        '12/01/2023 - 30/03/2023',
                      ),
                      infoTile(
                        notifier,
                        'Cap Quantity',
                        '5 Tokens',
                      ),
                      infoTile(
                        notifier,
                        'Cap Duration',
                        '12/01/2023 - 20/01/2023',
                      ),
                      infoTile(
                        notifier,
                        'Proceed Payout Cycle',
                        'Monthly',
                      ),
                      infoTile(
                        notifier,
                        'Payout Currency',
                        'TROV',
                      ),
                      infoTile(
                        notifier,
                        'Countries Exempted',
                        'See list',
                      ),
                      Card(
                        elevation: notifier.isDark ? 0 : 3,
                        shadowColor: Colors.black,
                        color: notifier.gettilewihitecolor,
                        margin:
                            EdgeInsets.symmetric(vertical: 10, horizontal: 20),
                        child: Padding(
                          padding: const EdgeInsets.symmetric(vertical: 8.0),
                          child: ListTile(
                            title: Row(
                              children: [
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      'Countries Exempted',
                                      style: TextStyle(
                                        fontSize: 13,
                                        fontFamily: fontsemibold,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          0, 3.0, 0, 0),
                                      child: Text(
                                        'See list',
                                        style: TextStyle(
                                          decoration: TextDecoration.underline,
                                          fontSize: 13,
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
                      ),
                      Card(
                        elevation: notifier.isDark ? 0 : 3,
                        shadowColor: Colors.black,
                        color: notifier.gettilewihitecolor,
                        margin:
                            EdgeInsets.symmetric(vertical: 10, horizontal: 20),
                        child: Padding(
                          padding: const EdgeInsets.symmetric(vertical: 8.0),
                          child: ListTile(
                            title: Row(
                              children: [
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      'Proof of Existence',
                                      style: TextStyle(
                                        fontSize: 13,
                                        fontFamily: fontsemibold,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          0, 3.0, 0, 0),
                                      child: Text(
                                        'C of O',
                                        style: TextStyle(
                                          decoration: TextDecoration.underline,
                                          fontSize: 13,
                                          fontFamily: fontbody,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          0, 3.0, 0, 0),
                                      child: Text(
                                        'Survey Plan',
                                        style: TextStyle(
                                          decoration: TextDecoration.underline,
                                          fontSize: 13,
                                          fontFamily: fontbody,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                    ),
                                    Padding(
                                      padding: const EdgeInsets.fromLTRB(
                                          0, 3.0, 0, 0),
                                      child: Text(
                                        'Governor\'s Consent',
                                        style: TextStyle(
                                          decoration: TextDecoration.underline,
                                          fontSize: 13,
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
                          '35 TROV',
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
                          '35 TROV',
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
                          '35 TROV',
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
                          '35 TROV',
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
                          '35 TROV',
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
}
