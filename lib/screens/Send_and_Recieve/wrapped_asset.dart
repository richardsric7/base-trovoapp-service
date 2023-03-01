import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/curated_asset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
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
  late Wallet wallet;
  late Asset? asset;
  late CuratedAsset? curatedAsset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    curatedAsset = appState.userInfo!.curatedSwapList!.firstWhereOrNull(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (curatedAsset!.realAssetImageUrl
                      .toString()
                      .isNotEmpty) ...[
                    Image.network(
                      curatedAsset!.realAssetImageUrl!,
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
                    getAssetCode(curatedAsset!.assetCode),
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
                    text: 'Welcome to the assetCode '
                        .replaceAll('assetCode', asset!.assetCode!),
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
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Column(
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Image.asset("assets/images/one.png", height: 30),
                        SizedBox(width: width / 90),
                        Container(
                          width: width / 1.3,
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
                                curatedAsset!.assetRedemptionInstructions!,
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
                    SizedBox(height: height / 50),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Image.asset("assets/images/two.png", height: 30),
                        SizedBox(width: width / 90),
                        Container(
                          width: width / 1.3,
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
    return Container(
      constraints: BoxConstraints(maxWidth: width / 1.3),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceAround,
        children: [
          actionButton("assets/images/deposit.png", 'Deposit', () {
            appState.viewData = {
              'walletPublicKey': wallet.publicKey,
              'assetCode': asset!.assetCode,
              'assetIssuer': asset!.assetIssuer,
            };

            if (asset!.cryptoWalletDepositAddresses!.length > 0) {
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: SelectDepositAddressViewPageConfig,
              );
            } else {
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: GenerateDepositAddressViewPageConfig,
              );
            }
          }),
          actionButton("assets/images/withdraw.png", 'Withdraw', () {
            appState.viewData = {
              'walletPublicKey': wallet.publicKey,
              'assetCode': asset!.assetCode,
              'assetIssuer': asset!.assetIssuer,
            };

            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: WithdrawAssetViewPageConfig,
            );
          }),
          actionButton("assets/images/history-btn.png", 'History', () {
            appState.fetchDepositHistory(
              context,
              publicKey: wallet.publicKey!,
              currency: asset!.assetCode,
            );
            appState.viewData = {
              'walletPublicKey': wallet.publicKey,
              'assetCode': asset!.assetCode,
              'assetIssuer': asset!.assetIssuer,
            };

            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: DepositWithdrawHistoryViewPageConfig,
            );
          }),
        ],
      ),
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
