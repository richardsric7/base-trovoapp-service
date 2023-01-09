import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WrappedAsset extends StatefulWidget {
  const WrappedAsset({Key? key}) : super(key: key);

  @override
  State<WrappedAsset> createState() => _WrappedAssetState();
}

class _WrappedAssetState extends State<WrappedAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  Map activeWallet = {};
  Map activeAsset = {};
  var claimedAssets;
  bool isInitiator = false;
  bool isSharedWallet = false;
  dynamic selectedWallet = '';
  dynamic selectedAsset = '';

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
    userInfo = appState.userInfo!;
    assetBalances = appState.assetBalances;

    if (activeWallet.isEmpty) {
      activeWallet = appState.allWallets[appState.activeWallet!.publicKey!];
    }
    selectedWallet = activeWallet['publicKey'];
    claimedAssets = activeWallet['claimedAssets'];
    activeAsset = appState.viewData![WrappedAssetViewPageConfig.key];

    if (appState.viewData![WrappedAssetViewPageConfig.key] != null) {
      selectedAsset = "${getAssetCode(
        appState.viewData![WrappedAssetViewPageConfig.key]['assetCode'],
      )}|${getAssetIssuer(
        appState.viewData![WrappedAssetViewPageConfig.key]['assetIssuer'],
      )}";
    }

    isSharedWallet = activeWallet['sharedAccessEnabled'] == 1;

    // if this is a shared wallet
    if (isSharedWallet) {
      if (activeWallet['permission'] == 'INITIATOR')
        isInitiator = true;
      else
        isInitiator = false;
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(height / 15),
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (activeAsset["imageUrl"].toString().isNotEmpty) ...[
                    Image.network(
                      activeAsset["imageUrl"],
                      height: 30,
                      width: 30,
                      errorBuilder: (context, error, stackTrace) {
                        return Image.asset(
                          'assets/images/trovo.png',
                          height: 80,
                          width: 80,
                        );
                      },
                    ),
                  ],
                  SizedBox(
                    width: width / 50.0,
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
              SizedBox(height: height / 30),
              Image.asset("assets/images/rafiki.png", height: height / 4),
              Padding(
                padding: const EdgeInsets.symmetric(
                    vertical: 15.0, horizontal: 25.0),
                child: RichText(
                  text: TextSpan(
                    text: 'Welcome to the USDC ',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                    children: [
                      TextSpan(
                        text: 'Deposit ',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: 'and ',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: 'Withdraw ',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: 'page. Here, you can do the following:',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  textAlign: TextAlign.justify,
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Image.asset("assets/images/one.png", height: 30),
                    SizedBox(width: width / 90),
                    Container(
                      width: width / 1.2,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Deposit and Withdraw',
                            style: TextStyle(
                                fontSize: 15,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold),
                          ),
                          SizedBox(height: height / 90),
                          Text(
                            'Deposit USDC on Trovo wallet from anywhere outside the Bantu ecosystem. You can also withdraw USDC to any USDC wallet of your choice outside the Bantu ecosystem at any time. Tap on the deposit/withdraw button below to begin ',
                            textAlign: TextAlign.justify,
                            style: TextStyle(
                                fontSize: 15,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody),
                          ),
                        ],
                      ),
                    )
                  ],
                ),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Image.asset("assets/images/two.png", height: 30),
                    SizedBox(width: width / 90),
                    Container(
                      width: width / 1.2,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Deposit/Withdraw History',
                            style: TextStyle(
                                fontSize: 15,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold),
                          ),
                          SizedBox(height: height / 90),
                          Text(
                            'View your deposit and withdrawal history by clicking on the transaction history button here.',
                            textAlign: TextAlign.justify,
                            style: TextStyle(
                                fontSize: 15,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody),
                          ),
                        ],
                      ),
                    )
                  ],
                ),
              ),
              SizedBox(height: height / 20),
              actionButtons(),
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
        actionButton("assets/images/deposit.png", 'Deposit', () {
          appState.viewData![GenerateDepositAddressViewPageConfig.key] =
              appState.viewData![WrappedAssetViewPageConfig.key];
          appState.viewData![GenerateDepositAddressViewPageConfig.key]
              ['walletInfo'] = activeWallet;

          appState.currentAction = PageAction(
            state: PageState.addPage,
            page: GenerateDepositAddressViewPageConfig,
          );
        }),
        actionButton("assets/images/withdraw.png", 'Withdraw', () {}),
        actionButton("assets/images/history-btn.png", 'History', () {}),
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
                  fontSize: 10,
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
}
