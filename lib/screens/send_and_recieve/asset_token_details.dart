import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart' hide Trans;
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/wallet_slides.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetTokenDetails extends StatefulWidget {
  const AssetTokenDetails({Key? key}) : super(key: key);

  @override
  State<AssetTokenDetails> createState() => _AssetTokenDetailsState();
}

class _AssetTokenDetailsState extends State<AssetTokenDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet wallet;
  late TokenizedAsset? asset;
  // String selectedWallet = '';
  // String selectedAsset = '';

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    userInfo = appState.userInfo!;
    wallet = userInfo.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    // asset = wallet.tokenizedAssets!.firstWhere(
    //   (asset) =>
    //       asset.assetCode == appState.viewData!['assetCode'] &&
    //       asset.assetIssuer == appState.viewData!['assetIssuer'],
    // );
    asset = appState.viewData!["tokenizedAsset"];

    // free the memory..... lol
    appState.viewData = {};
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    // if (selectedWallet.isEmpty) {
    //   selectedWallet = wallet.publicKey!;
    // }

    // if (asset == null) {
    //   asset = wallet.tokenizedAssets!.firstWhere(
    //     (asset) => asset.assetCode == '' && asset.assetIssuer == '',
    //   );
    //   selectedAsset = '';
    // }

    // if (selectedAsset.isEmpty) {
    //   selectedAsset =
    //       "${getAssetCode(asset!.assetCode)}|${getAssetIssuer(asset!.assetIssuer)}";
    // }

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
            // actions: [
            //   Container(
            //     width: width / 1.2,
            //     child: Row(
            //       children: [
            //         TopDropdowns(
            //           onWalletChanged: (newValue) {
            //             selectedWallet = newValue;
            //             this.wallet = userInfo.getWallet(selectedWallet);

            //             this.asset = wallet.tokenizedWallets.firstWhereOrNull((x) =>
            //                 "${getAssetCode(x.assetCode)}|${getAssetIssuer(x.assetIssuer)}" ==
            //                 selectedAsset);

            //             setState(() {});
            //           },
            //           onAssetChanged: (newValue) {
            //             setState(() {
            //               selectedAsset = newValue;
            //               newValue = newValue.toString().contains('XBN')
            //                   ? '|'
            //                   : newValue;
            //               for (var asset in wallet.tokenizedAssets!) {
            //                 var splitNewValue = newValue.toString().split('|');
            //                 if (asset.assetCode == splitNewValue[0] &&
            //                     asset.assetIssuer == splitNewValue[1]) {
            //                   this.asset = asset;
            //                 }
            //               }
            //             });
            //           },
            //           claimedAssets: wallet.tokenizedAssets!,
            //           selectedAsset: selectedAsset,
            //           selectedWallet: selectedWallet,
            //         ),
            //       ],
            //     ),
            //   ),
            //   if (appState.walletMode == "Testnet") ...[
            //     Visibility(
            //       visible: true,
            //       child: Container(
            //         color: Color(0xFFAA453E),
            //         width: 18,
            //         child: RotatedBox(
            //           quarterTurns: 1,
            //           child: Column(
            //             mainAxisAlignment: MainAxisAlignment.center,
            //             children: [
            //               Padding(
            //                 padding: const EdgeInsets.symmetric(horizontal: 2),
            //                 child: Text(
            //                   "testnet".tr(),
            //                   style: TextStyle(
            //                     fontFamily: fontsemibold,
            //                     color: wihitecolor,
            //                     fontSize: 11,
            //                   ),
            //                 ),
            //               ),
            //             ],
            //           ),
            //         ),
            //       ),
            //     ),
            //   ]
            // ],
          ),
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
                isSharedWallet: wallet.isSharedWallet,
                walletType: wallet.walletType!,
                totalBalance:
                    '${formatNumber(asset!.subscriptionAmount!)} ${getAssetCode(asset!.assetCode)}',
                fiatBalance:
                    '${calculateFiatValue(asset!.subscriptionAmount!.toString(), asset!.usdPrice!.toString(), appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                initialHiddenState: appState.hideBalances,
              ),
              SizedBox(
                height: height / 30,
              ),
              assetInfo(),
              SizedBox(
                height: height / 20,
              ),
              // if ((!wallet.isSharedWallet ||
              //         wallet.isInitiator ||
              //         wallet.isPrimaryWallet) &&
              //     wallet.walletType == 0) ...[
              //   actionButtons(),
              // ] else ...[
              //   Button(
              //     "receive".tr(),
              //     notifier.getbluecolor,
              //     wihitecolor,
              //     onTap: () {
              //       appState.viewData = {
              //         'walletPublicKey': wallet.publicKey,
              //         'assetCode': asset!.assetCode,
              //         'assetIssuer': asset!.assetIssuer,
              //       };

              //       appState.currentAction = PageAction(
              //         state: PageState.addPage,
              //         page: ReceiveAssetViewPageConfig,
              //       );
              //     },
              //   ),
              // ],
              Button(
                "back".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  Navigator.of(context).pop();
                },
              ),
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
          // if (curatedAsset != null &&
          //     (curatedAsset!.isWithdrawable ||
          //         curatedAsset!.canGenerateDepositAddresses == 1)) ...[
          //   actionButton(
          //       "assets/images/dep-with-button.png", 'Deposit/Withdraw', () {
          //     appState.viewData = {
          //       'walletPublicKey': wallet.publicKey,
          //       'assetCode': asset!.assetCode,
          //       'assetIssuer': asset!.assetIssuer,
          //     };

          //     appState.currentAction = PageAction(
          //       state: PageState.addPage,
          //       page: WrappedAssetViewPageConfig,
          //     );
          //   }),
          // ],
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
                    '${getAssetCode(asset!.assetCode!)} ${"token".tr()}',
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
                    child: Image.memory(
                      base64Decode(asset!.assetLogo!),
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
                  Text(
                    'Price',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(
                    height: height / 90.0,
                  ),
                  SizedBox(
                    width: width / 1.3,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        '1 ${asset?.assetCode?.toUpperCase()} = ${formatNumberShort(double.parse(getFiatRate(asset!.usdPrice.toString(), appState.defaultCurrency, appState, getUnFormatted: true)))} ${appState.defaultCurrency}',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontWeight: FontWeight.w500,
                          color: notifier.getbluewhitecolor,
                          fontSize: 15.sp,
                          fontFamily: fontbody,
                        ),
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  if (asset!.assetIssuer!.toString().isNotEmpty) ...[
                    Text(
                      "issuerpubkey".tr(),
                      style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(
                      width: width / 1.5,
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
                                showSnackBar("issuerpubkey".tr(), context),
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
