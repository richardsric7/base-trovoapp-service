import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SharedWalletAssetDetails extends StatefulWidget {
  const SharedWalletAssetDetails({Key? key}) : super(key: key);

  @override
  State<SharedWalletAssetDetails> createState() =>
      _SharedWalletAssetDetailsState();
}

class _SharedWalletAssetDetailsState extends State<SharedWalletAssetDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  var assetBalances;
  var activeAsset;
  var claimedAssets;
  var walletDetails; // details of the current shared wallet
  bool isInitiator = false;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    assetBalances = appState.assetBalances;
    activeAsset =
        appState.viewData![SharedWalletAssetDetailsViewPageConfig.key];
    walletDetails = activeAsset['walletInfo'];
    for (var i = 0; i < walletDetails['permissions'].length; i++) {
      if (walletDetails['permissions'][i] == 'INITIATOR') isInitiator = true;
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          'Shared Wallet',
          notifier.getbluewhitecolor,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  SizedBox(
                    width: 20,
                  ),
                  Text(
                    getAssetCode(activeAsset['assetCode']),
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              WalletSlide(
                backColor: notifier.getbluecolor,
                foreColor: wihitecolor,
                alias: walletDetails['walletAlias'],
                totalBalance:
                    '${formatNumber(double.parse(activeAsset['amount']))} ${getAssetCode(activeAsset['assetCode'])}',
                fiatBalance:
                    '${calculateFiatValue(activeAsset['amount'], activeAsset['usdPrice'], appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                initialHiddenState: appState.hideBalances,
              ),
              SizedBox(
                height: height / 30,
              ),
              assetInfo(),
              SizedBox(
                height: height / 20,
              ),
              if (isInitiator) ...[
                actionButtons(),
              ] else ...[
                Button(
                  'Receive',
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.viewData![
                        RecieveAssetSharedWalletViewPageConfig
                            .key] = appState
                        .viewData![SharedWalletAssetDetailsViewPageConfig.key];

                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: RecieveAssetSharedWalletViewPageConfig,
                    );
                  },
                ),
              ]
            ],
          ),
        ),
      ),
    );
  }

  Widget actionButtons() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        actionButton("assets/images/send.png", 'Send', () {
          appState.viewData![SendAssetSharedWalletViewPageConfig.key] =
              appState.viewData![SharedWalletAssetDetailsViewPageConfig.key];

          appState.currentAction = PageAction(
            state: PageState.addPage,
            page: SendAssetSharedWalletViewPageConfig,
          );
        }),
        actionButton("assets/images/receive.png", 'Receive', () {
          appState.viewData![RecieveAssetSharedWalletViewPageConfig.key] =
              appState.viewData![SharedWalletAssetDetailsViewPageConfig.key];

          appState.currentAction = PageAction(
            state: PageState.addPage,
            page: RecieveAssetSharedWalletViewPageConfig,
          );
        }),
      ],
    );
  }

  Widget actionButton(iconUrl, actionText, action) {
    return GestureDetector(
      onTap: action,
      child: Container(
        width: width / 3.9,
        height: height / 10,
        child: Padding(
          padding: const EdgeInsets.symmetric(
            vertical: 8.0,
          ),
          child: Column(
            children: [
              Image.asset(
                iconUrl,
                width: width / 8,
                color: notifier.getbluewhitecolor,
              ),
              Text(
                actionText,
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget assetInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        height: height / 2.5,
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    '${getAssetCode(activeAsset['assetCode'])} Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.3,
                    child: Text(
                      'TROV token (TROV) is the utility token that powers the Trovotech ecosystem. TROV token is used to access discounts, voting rights, airdrops, NFTs and other community incentives. ',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50.0,
                  ),
                  Text(
                    'www.trovotech.io',
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
