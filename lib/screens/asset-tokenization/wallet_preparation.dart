import 'package:country_picker/country_picker.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/bottom_tab_page.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WalletPreparation extends StatefulWidget {
  const WalletPreparation({Key? key}) : super(key: key);

  @override
  State<WalletPreparation> createState() => _WalletPreparationState();
}

class _WalletPreparationState extends State<WalletPreparation>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late TabController tabController;
  String selectedCountry = 'Nigeria';
  List<String> assetCategories = [
    'Agriculture and Farming',
    'Art and Collectibles',
    'Automotive and Transportation',
    'Commodities',
    'Eduction and Learning',
    'Environmental and Renewable Energy',
    'Financing and Banking',
    'Gaming and Virtual Reality',
    'Healthcare and Medical',
    'Real Estate',
  ];

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getMintingWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getMintingWallets.forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet.publicKey));
    });
    return wallets;
  }

  List<DropdownMenuItem<String>> get getMarketMakingWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getMarketMakingWallets.forEach((wallet) {
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
              'Wallet Preparation',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Please select wallets for your asset token in order to create an asset token',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Text(
                  'Select Minting Wallet',
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                Text(
                  'What does this mean?',
                  style: TextStyle(
                    decoration: TextDecoration.underline,
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getMintingWallets,
                null,
                appState.userInfo!.getMintingWallets.length > 0
                    ? appState.userInfo!.getMintingWallets.first.alias
                    : '',
                context,
                null,
              ),
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: TextButton(
                    onPressed: () async {
                      appState.returnView = PageAction(
                          state: PageState.addAll,
                          pages: [
                            BottomHomePageConfig,
                            WalletPreparationViewPageConfig
                          ]);
                      appState.currentAction =
                          PageAction(state: PageState.addAll, pages: [
                        BottomHomePageConfig,
                      ]);
                      changeTabPage(appState, ButtomTabPage.Wallets.index);
                      setState(() {});
                    },
                    child: Text(
                      'Or create a new wallet for this',
                      style: TextStyle(
                        decoration: TextDecoration.underline,
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Text(
                  'Select Minting Wallet',
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                Text(
                  'What does this mean?',
                  style: TextStyle(
                    decoration: TextDecoration.underline,
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getMarketMakingWallets,
                null,
                appState.userInfo!.getMarketMakingWallets.length > 0
                    ? appState.userInfo!.getMarketMakingWallets.first.alias
                    : '',
                context,
                null,
              ),
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 10.0),
                  child: TextButton(
                    onPressed: () {
                      appState.returnView = PageAction(
                          state: PageState.addAll,
                          pages: [
                            BottomHomePageConfig,
                            WalletPreparationViewPageConfig
                          ]);
                      appState.currentAction =
                          PageAction(state: PageState.addAll, pages: [
                        BottomHomePageConfig,
                      ]);
                      changeTabPage(appState, ButtomTabPage.Wallets.index);
                      setState(() {});
                    },
                    child: Text(
                      'Or create a new wallet for this',
                      style: TextStyle(
                        decoration: TextDecoration.underline,
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 20,
            ),
            Button(
              'Continue',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: TokenizeAssetViewPageConfig);
              },
            ),
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
  }
}
