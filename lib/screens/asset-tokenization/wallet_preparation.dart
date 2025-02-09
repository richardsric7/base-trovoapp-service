import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/popups.dart';
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
  List<String> excludedWallets = [];
  List<Wallet> issuingWallets = [];

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<Wallet>> get getIssuingWallets {
    List<DropdownMenuItem<Wallet>> wallets = [];
    issuingWallets.forEach((wallet) {
      if (wallet.isSharedWalletAndCanInitiate) {
        wallets.add(DropdownMenuItem(
            child: Text(
              wallet.alias!,
              overflow: TextOverflow.ellipsis,
            ),
            value: wallet));
      }
    });
    return wallets;
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    excludedWallets = appState.viewData!['excludedWallets'];
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    issuingWallets = appState.userInfo!.getMintingWallets
        .where((wallet) => !excludedWallets.contains(wallet.publicKey))
        .toList();

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              'walletpreparation'.tr(),
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            walletPreparation(),
            SizedBox(
              height: height / 20,
            ),
            Button(
              "continuee".tr(),
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                // if (appState.activeTokenizationWalletPublicKey == null) {
                //   popup(context,
                //       title: 'Error',
                //       message: 'Please select your asset tokenization wallet');
                //   return;
                // }

                // if (appState.activeDistributionWalletPublicKey == null) {
                //   popup(context,
                //       title: 'Error',
                //       message: 'Please select your asset distribution wallet');
                //   return;
                // }

                appState.viewData = {};
                appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: SetupAndComplianceViewPageConfig);
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

  Widget walletPreparation() {
    return Column(
      children: [
        SizedBox(
          height: height / 50,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            "selecttokenizationwallet".tr(),
            textAlign: TextAlign.left,
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
        if (issuingWallets.isNotEmpty) ...[
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 15),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  "selectissuingwallet".tr(),
                  style: TextStyle(
                    fontSize: 13,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                Text(
                  "whatdoesthismean".tr(),
                  style: TextStyle(
                    decoration: TextDecoration.underline,
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ],
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          // Padding(
          //   padding: const EdgeInsets.symmetric(horizontal: 10.0),
          //   child: dropdown(
          //     (value) {
          //       var wallet = value as Wallet;
          //       // appState.setActiveTokenizationWalletPublicKey =
          //       //     wallet.publicKey;
          //       // appState.setActiveDistributionWalletPublicKey =
          //       //     wallet.linkedWalletPublicKey;
          //     },
          //     getIssuingWallets,
          //     null,
          //     getHintTextForMintingWallet(),
          //     context,
          //     null,
          //   ),
          // ),
          SizedBox(
            height: height / 70,
          ),
          Text(
            "oR".tr(),
            style: TextStyle(
              decoration: TextDecoration.underline,
              fontSize: 12,
              fontFamily: fontsemibold,
              color: notifier.getbluewhitecolor,
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
        ],
        TextButton(
          onPressed: () {
            appState.returnView = PageAction(
                state: PageState.addAll,
                pages: [BottomHomePageConfig, WalletPreparationViewPageConfig]);
            addSubWalletPopup(context);
          },
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.add_circle_outline_outlined,
                color: notifier.getbluewhitecolor,
              ),
              Text(
                "Create New Tokenization Wallet".tr(),
                style: TextStyle(
                  decoration: TextDecoration.underline,
                  fontSize: 12,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ],
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
      ],
    );
  }

  // String getHintTextForMintingWallet() {
  //   var wallet = appState.userInfo!.getMintingWallets.where(
  //       (w) => w.publicKey == appState.activeTokenizationWalletPublicKey);

  //   return wallet.length > 0
  //       ? wallet.first.alias!
  //       : 'Select tokenization wallet';
  // }

  Widget CheckItem(
    String name,
    void Function()? onClick, {
    required Color backColor,
    required Color foreColor,
    required Color borderColor,
    required Icon icon,
    double? fontSize = 15,
  }) {
    return Padding(
      padding: const EdgeInsets.all(3.0),
      child: ElevatedButton(
        onPressed: onClick,
        style: ButtonStyle(
          overlayColor:
              MaterialStateProperty.all<Color>(notifier.getsplashgrey),
          elevation: MaterialStateProperty.all<double>(0),
          backgroundColor: MaterialStateProperty.all<Color>(backColor),
          side: MaterialStateProperty.all(
            BorderSide(color: borderColor, width: 1, style: BorderStyle.solid),
          ),
          shape: MaterialStateProperty.all<RoundedRectangleBorder>(
            const RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(10),
              ),
            ),
          ),
        ),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            icon,
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                  color: foreColor, fontFamily: fontbody, fontSize: fontSize),
            ),
          ],
        ),
      ),
    );
  }
}
