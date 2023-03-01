import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/curated_asset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import 'package:trovo_wallet/widgets/top_drop_downs.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetDetails extends StatefulWidget {
  const AssetDetails({Key? key}) : super(key: key);

  @override
  State<AssetDetails> createState() => _AssetDetailsState();
}

class _AssetDetailsState extends State<AssetDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet wallet;
  late Asset? asset;
  late CuratedAsset? curatedAsset;
  String selectedWallet = '';
  String selectedAsset = '';

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    userInfo = appState.userInfo!;
    wallet = userInfo.getWallet(
      appState.viewData!['walletPublicKey'],
    );
    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    curatedAsset = userInfo.curatedSwapList!.firstWhereOrNull(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    // free the memory..... lol
    appState.viewData = {};
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    if (selectedWallet.isEmpty) {
      selectedWallet = wallet.publicKey!;
    }

    if (asset == null) {
      asset = wallet.claimedAssets!.firstWhere(
        (asset) => asset.assetCode == '' && asset.assetIssuer == '',
      );
      selectedAsset = '';
    }

    this.curatedAsset = userInfo.curatedSwapList!.firstWhereOrNull(
      (curatedAsset) =>
          curatedAsset.assetCode == asset!.assetCode &&
          curatedAsset.assetIssuer == asset!.assetIssuer,
    );

    if (selectedAsset.isEmpty) {
      selectedAsset =
          "${getAssetCode(asset!.assetCode)}|${getAssetIssuer(asset!.assetIssuer)}";
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
              actions: [
                Container(
                  width: width / 1.2,
                  child: Row(
                    children: [
                      TopDropdowns(
                        onWalletChanged: (newValue) {
                          selectedWallet = newValue;
                          this.wallet = userInfo.getWallet(selectedWallet);

                          this.asset = wallet.claimedAssets!.firstWhereOrNull((x) =>
                              "${getAssetCode(x.assetCode)}|${getAssetIssuer(x.assetIssuer)}" ==
                              selectedAsset);

                          setState(() {});
                        },
                        onAssetChanged: (newValue) {
                          setState(() {
                            selectedAsset = newValue;
                            newValue = newValue.toString().contains('XBN')
                                ? '|'
                                : newValue;
                            for (var asset in wallet.claimedAssets!) {
                              var splitNewValue =
                                  newValue.toString().split('|');
                              if (asset.assetCode == splitNewValue[0] &&
                                  asset.assetIssuer == splitNewValue[1]) {
                                this.asset = asset;
                              }
                            }
                          });
                        },
                        claimedAssets: wallet.claimedAssets!,
                        selectedAsset: selectedAsset,
                        selectedWallet: selectedWallet,
                      ),
                    ],
                  ),
                ),
                Visibility(
                  visible: true,
                  child: Padding(
                    padding: EdgeInsets.only(top: 5),
                    child: Banner(
                      location: BannerLocation.topEnd,
                      message: "Testnet",
                    ),
                  ),
                ),
              ]),
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
                    getAssetCode(asset!.assetCode),
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
                alias: wallet.alias.toString().capitalizeFirst!,
                totalBalance:
                    '${formatNumber(asset!.amount!)} ${getAssetCode(asset!.assetCode)}',
                fiatBalance:
                    '${calculateFiatValue(asset!.amount!.toString(), asset!.usdPrice!.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                initialHiddenState: appState.hideBalances,
              ),
              SizedBox(
                height: height / 30,
              ),
              if (curatedAsset != null) curatedAssetInfo() else assetInfo(),
              SizedBox(
                height: height / 20,
              ),
              if ((!wallet.isSharedWallet ||
                      wallet.isInitiator ||
                      wallet.isPrimaryWallet) &&
                  wallet.walletType == 0) ...[
                actionButtons(),
              ] else ...[
                Button(
                  'Receive',
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.viewData = {
                      'walletPublicKey': wallet.publicKey,
                      'assetCode': asset!.assetCode,
                      'assetIssuer': asset!.assetIssuer,
                    };

                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: ReceiveAssetViewPageConfig,
                    );
                  },
                ),
              ],
              SizedBox(
                height: height / 20,
              ),
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
          actionButton("assets/images/send.png", 'Send', () {
            appState.viewData = {
              'walletPublicKey': wallet.publicKey,
              'assetCode': asset!.assetCode,
              'assetIssuer': asset!.assetIssuer,
            };

            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: SendAssetViewPageConfig,
            );
          }),
          if (curatedAsset != null &&
              (curatedAsset!.isWithdrawable ||
                  curatedAsset!.canGenerateDepositAddresses == 1)) ...[
            actionButton(
                "assets/images/dep-with-button.png", 'Deposit/Withdraw', () {
              appState.viewData = {
                'walletPublicKey': wallet.publicKey,
                'assetCode': asset!.assetCode,
                'assetIssuer': asset!.assetIssuer,
              };

              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: WrappedAssetViewPageConfig,
              );
            }),
          ],
          actionButton("assets/images/receive.png", 'Receive', () {
            appState.viewData = {
              'walletPublicKey': wallet.publicKey,
              'assetCode': asset!.assetCode,
              'assetIssuer': asset!.assetIssuer,
            };

            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: ReceiveAssetViewPageConfig,
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
        height: height / 9,
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
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 12,
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

  Widget curatedAssetInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        constraints: BoxConstraints(minHeight: height / 2.5),
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
                    '${getAssetCode(curatedAsset!.assetCode)} Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  if (asset!.imageUrl != null) ...[
                    SizedBox(
                      height: height / 50.0,
                    ),
                    Container(
                      width: width / 1.3,
                      child: Image.network(
                        asset!.imageUrl!,
                        height: 50,
                        width: 50,
                        errorBuilder: (context, error, stackTrace) {
                          return Image.asset(
                            'assets/images/trovo.png',
                            height: 50,
                            width: 50,
                          );
                        },
                      ),
                    ),
                  ],
                  SizedBox(
                    height: height / 50.0,
                  ),
                  Text(
                    curatedAsset!.website!,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.3,
                    child: Text(
                      curatedAsset!.description!,
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
                  if (asset!.assetIssuer.toString().isNotEmpty) ...[
                    Text(
                      'Issuer Public Key',
                      style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(
                      width: width / 1.3,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: [
                          SizedBox(
                            width: width / 20,
                          ),
                          Expanded(
                            flex: 3,
                            child: Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
                              child: Text(
                                truncate(asset!.assetIssuer!, length: 5) +
                                    asset!.assetIssuer!.toString().substring(
                                        asset!.assetIssuer!.toString().length -
                                            5),
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 15.sp,
                                  fontFamily: fontbody,
                                ),
                              ),
                            ),
                          ),
                          Expanded(
                            flex: 1,
                            child: IconButton(
                              padding: EdgeInsets.zero,
                              onPressed: () => {
                                Clipboard.setData(
                                  ClipboardData(
                                    text: asset!.assetIssuer!,
                                  ),
                                ),
                                showSnackBar('Issuer public key', context),
                              },
                              icon: Icon(Icons.copy),
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                        ],
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    if (curatedAsset!.contactEmail!.toString().isNotEmpty) ...[
                      Text(
                        'Contact Email',
                        style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.w600,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold),
                      ),
                      SizedBox(
                        width: width / 1.3,
                        child: Row(
                          children: [
                            Expanded(
                              flex: 3,
                              child: Padding(
                                padding: const EdgeInsets.symmetric(
                                    horizontal: 20.0),
                                child: Text(
                                  curatedAsset!.contactEmail!,
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontWeight: FontWeight.w500,
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 15.sp,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              ),
                            )
                          ],
                        ),
                      ),
                    ],
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget assetInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        constraints: BoxConstraints(minHeight: height / 2.5),
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
                    '${getAssetCode(asset!.assetCode!)} Token',
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
                    child: Image.network(
                      asset!.imageUrl!,
                      height: 50,
                      width: 50,
                      errorBuilder: (context, error, stackTrace) {
                        return Image.asset(
                          'assets/images/trovo.png',
                          height: 50,
                          width: 50,
                        );
                      },
                    ),
                  ),
                  SizedBox(
                    height: height / 50.0,
                  ),
                  if (asset!.assetIssuer!.toString().isNotEmpty) ...[
                    Text(
                      'Issuer Public Key',
                      style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(
                      width: width / 1.7,
                      child: Row(
                        children: [
                          Expanded(
                            flex: 3,
                            child: Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
                              child: Text(
                                truncate(asset!.assetIssuer!, length: 5) +
                                    asset!.assetIssuer!.toString().substring(
                                        asset!.assetIssuer!.toString().length -
                                            5),
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 15.sp,
                                  fontFamily: fontbody,
                                ),
                              ),
                            ),
                          ),
                          Expanded(
                            flex: 1,
                            child: IconButton(
                              padding: EdgeInsets.zero,
                              onPressed: () => {
                                Clipboard.setData(
                                  ClipboardData(
                                    text: asset!.assetIssuer!,
                                  ),
                                ),
                                showSnackBar('Issuer public key', context),
                              },
                              icon: Icon(Icons.copy),
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
