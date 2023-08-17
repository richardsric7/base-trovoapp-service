import 'package:flutter/cupertino.dart';
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

class TokenizedAssetDetail extends StatefulWidget {
  const TokenizedAssetDetail({Key? key}) : super(key: key);

  @override
  State<TokenizedAssetDetail> createState() => _TokenizedAssetDetail();
}

class _TokenizedAssetDetail extends State<TokenizedAssetDetail>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;

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
            SizedBox(
              height: height / 50,
            ),
            Text(
              'ATLANTIS 1',
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
              'Atlantis Estate 1 token',
              style: TextStyle(
                fontSize: 15,
                fontFamily: fontbody,
                color: notifier.getbluewhitecolor,
              ),
            ),
            SizedBox(
              height: height / 30,
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
                            CupertinoIcons.add_circled_solid,
                            size: 25,
                            color: notifier.getbluewhitecolor,
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                          Text(
                            'Subscribe',
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
                          Icon(
                            CupertinoIcons.cart_fill,
                            size: 25,
                            color: notifier.getbluewhitecolor,
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                          Text(
                            'Buy',
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
            SizedBox(height: height / 50),
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
                          'Atlantis Estate 1 tokens are fractional tokens that represent part ownership (via investment) of our real estate development project at Atlantis Estate, Lekki, Lagos, Nigeria. Subscribe to this token to earn rental income monthly pushed to your Trovo Wallet. Also earn benefit from asset appreciation with the ability to sell or buy any portion of your tokens at anytime.',
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
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 30.0),
                  child: Text(
                    'Purchase with cNGN',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: TextButton(
                    onPressed: () {},
                    child: Text(
                      'Tap to fund wallet now',
                      style: TextStyle(
                        decoration: TextDecoration.underline,
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(height: height / 50),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 30.0),
                  child: Text(
                    'Details of Asset',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            infoTile(
              notifier,
              'Asset Code',
              'ATLANTIS 1',
            ),
            infoTile(
              notifier,
              'Asset Category',
              'Real Estate',
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
              'www.atlantis.com',
            ),
            infoTile(
              notifier,
              'Asset Token Total Supply',
              '1000',
            ),
            infoTile(
              notifier,
              'Asset Tokens Quantity Purchased',
              '400',
            ),
            infoTile(
              notifier,
              'Total Subscribed Users',
              '2,000',
            ),
            infoTile(
              notifier,
              'Price Per Asset Token',
              '100 cNGN',
            ),
            infoTile(
              notifier,
              'Asset Token Purchase Method',
              'cNGN',
            ),
            infoTile(
              notifier,
              'Asset Token Sales Window',
              '12/01/2023 - 30/03/2023',
            ),
            infoTile(
              notifier,
              'Token Sale Cap',
              '5 [ATLANTIS 1] Token',
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
              'Payout Method',
              'cNGN',
            ),
            Card(
              elevation: notifier.isDark ? 0 : 3,
              shadowColor: Colors.black,
              color: notifier.gettilewihitecolor,
              margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: ListTile(
                  title: Row(
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Exempted Countries',
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
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
              margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: ListTile(
                  title: Row(
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Additional Requirements',
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
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
              margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: ListTile(
                  title: Row(
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Verified Proof of Existence',
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              'C of O',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              'Survey Plan',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              'Governor\'s Consent',
                              style: TextStyle(
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
    );
  }
}
