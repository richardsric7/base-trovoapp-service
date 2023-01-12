import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/WalletSlides.dart';
import 'package:trovo_wallet/widgets/topDropdowns.dart';
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
  var assetBalances;
  Map activeWallet = {};
  Map activeAsset = {};
  var claimedAssets;
  late Map curatedAsset;
  bool isInitiator = false;
  bool isSharedWallet = false;
  dynamic selectedWallet = '';
  dynamic selectedAsset = '';

  List<DropdownMenuItem<String>> assetDropdownItems(bool isSelected) {
    List<DropdownMenuItem<String>> menuItems = [];
    for (var asset in claimedAssets) {
      menuItems.add(DropdownMenuItem(
          child: Text(
            isSelected
                ? truncate(
                    getAssetCode(asset['assetCode']),
                    length: 3,
                  )
                : getAssetCode(asset['assetCode']),
            overflow: TextOverflow.visible,
          ),
          value:
              '${getAssetCode(asset['assetCode'])}|${getAssetIssuer(asset['assetIssuer'])}'));
    }
    return menuItems;
  }

  List<DropdownMenuItem<String>> walletDropdownItems(bool isSelected) {
    var walletsList = <DropdownMenuItem<String>>[];
    appState.allWallets.forEach((key, value) {
      walletsList.add(
        DropdownMenuItem(
          child: Row(
            children: [
              Container(
                constraints: isSelected
                    ? BoxConstraints(maxWidth: width / 4)
                    : BoxConstraints(maxWidth: width / 2.5),
                child: Text(
                  value['alias'],
                  overflow:
                      isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
                ),
              ),
              if (value['sharedAccessEnabled'] == 1) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.people_outline,
                  size: 17,
                  color: notifier.getbluecolor,
                )
              ],
              if (!isSelected && key == selectedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.check,
                  size: 18,
                  color: notifier.getbluecolor,
                )
              ],
            ],
          ),
          value: key,
        ),
      );
    });

    return walletsList;
  }

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

    curatedAsset = userInfo.curatedSwapList!.firstWhere(
        (asset) =>
            asset['assetCode'] ==
                appState.viewData![AssetDetailsViewPageConfig.key]
                    ['assetCode'] &&
            asset['assetIssuer'] ==
                appState.viewData![AssetDetailsViewPageConfig.key]
                    ['assetIssuer'],
        orElse: () => {});

    if (activeWallet.isEmpty) {
      activeWallet = appState.allWallets[appState.activeWallet!.publicKey!];
    }
    selectedWallet = activeWallet['publicKey'];
    claimedAssets = activeWallet['claimedAssets'];
    activeAsset = appState.viewData![AssetDetailsViewPageConfig.key];

    if (appState.viewData![AssetDetailsViewPageConfig.key] != null) {
      selectedAsset = "${getAssetCode(
        appState.viewData![AssetDetailsViewPageConfig.key]['assetCode'],
      )}|${getAssetIssuer(
        appState.viewData![AssetDetailsViewPageConfig.key]['assetIssuer'],
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
              actions: [
                Container(
                  width: width / 1.2,
                  child: Row(
                    children: [
                      TopDropdowns(
                        onWalletChanged: (newValue) {
                          selectedWallet = newValue;
                          activeWallet = appState.allWallets[newValue];
                          claimedAssets = activeWallet['claimedAssets'];

                          for (var asset in claimedAssets) {
                            // we need to somehow take care of the selected asset
                            // when switching wallets because of scenarios
                            // where one wallet has an asset that is not listed
                            // on the other. Here we are checking whether the
                            // newly selected wallet contains the currently
                            // selected asset and if it doesn't we switch
                            // back to the default asset which is XBN
                            if (asset['assetIssuer'] == selectedAsset ||
                                asset['assetIssuer'] == '') {
                              appState.viewData![
                                  AssetDetailsViewPageConfig.key] = asset;
                              break;
                            }
                          }
                          setState(() {});
                        },
                        onAssetChanged: (newValue) {
                          setState(() {
                            newValue = newValue.toString().contains('XBN')
                                ? '|'
                                : newValue;
                            for (var asset in claimedAssets) {
                              var splitNewValue =
                                  newValue.toString().split('|');
                              if (asset['assetCode'] == splitNewValue[0] &&
                                  asset['assetIssuer'] == splitNewValue[1]) {
                                appState.viewData![
                                    AssetDetailsViewPageConfig.key] = asset;
                              }
                            }
                          });
                        },
                        claimedAssets: claimedAssets,
                        selectedAsset: selectedAsset,
                        selectedWallet: selectedWallet,
                      ),
                    ],
                  ),
                )
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
                alias: activeWallet['alias'].toString().capitalizeFirst!,
                totalBalance:
                    '${formatNumber(double.parse(activeAsset['amount']))} ${getAssetCode(activeAsset['assetCode'])}',
                fiatBalance:
                    '${calculateFiatValue(activeAsset['amount'], activeAsset['usdPrice'], appState.defaultCurrency, appState)} ${appState.defaultCurrency}',
                initialHiddenState: appState.hideBalances,
              ),
              SizedBox(
                height: height / 30,
              ),
              if (curatedAsset.isNotEmpty) curatedAssetInfo() else assetInfo(),
              SizedBox(
                height: height / 20,
              ),
              if ((!isSharedWallet || isInitiator) &&
                  activeWallet['walletType'] == 0) ...[
                actionButtons(),
              ] else ...[
                Button(
                  'Receive',
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.viewData![ReceiveAssetViewPageConfig.key] =
                        appState.viewData![AssetDetailsViewPageConfig.key];
                    appState.viewData![ReceiveAssetViewPageConfig.key]
                        ['walletInfo'] = activeWallet;

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
            appState.viewData![SendAssetViewPageConfig.key] =
                appState.viewData![AssetDetailsViewPageConfig.key];
            appState.viewData![SendAssetViewPageConfig.key]['walletInfo'] =
                activeWallet;

            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: SendAssetViewPageConfig,
            );
          }),
          actionButton("assets/images/receive.png", 'Receive', () {
            appState.viewData![ReceiveAssetViewPageConfig.key] =
                appState.viewData![AssetDetailsViewPageConfig.key];
            appState.viewData![ReceiveAssetViewPageConfig.key]['walletInfo'] =
                activeWallet;

            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: ReceiveAssetViewPageConfig,
            );
          }),
          if (curatedAsset.isNotEmpty &&
              (curatedAsset['withdrawable'] == 1 ||
                  curatedAsset['generateDepositAddress'] == 1)) ...[
            actionButton(
                "assets/images/dep-with-button.png", 'Deposit/Withdraw', () {
              appState.viewData![WrappedAssetViewPageConfig.key] = curatedAsset;
              appState.viewData![WrappedAssetViewPageConfig.key]['usdPrice'] =
                  activeAsset['usdPrice'];
              appState.viewData![WrappedAssetViewPageConfig.key]['amount'] =
                  activeAsset['amount'];
              appState.viewData![WrappedAssetViewPageConfig.key]
                      ['cryptoWalletDepositAddresses'] =
                  activeAsset['cryptoWalletDepositAddresses'];
              appState.viewData![WrappedAssetViewPageConfig.key]['walletInfo'] =
                  activeWallet;

              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: WrappedAssetViewPageConfig,
              );
            }),
          ]
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
                    '${getAssetCode(curatedAsset['assetCode'])} Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  if (activeAsset["imageUrl"].toString().isNotEmpty) ...[
                    SizedBox(
                      height: height / 50.0,
                    ),
                    Container(
                      width: width / 1.3,
                      child: Image.network(
                        activeAsset["imageUrl"],
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
                    curatedAsset['website'],
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
                      curatedAsset['description'],
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
                  if (activeAsset['assetIssuer'].toString().isNotEmpty) ...[
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
                                truncate(activeAsset['assetIssuer'],
                                        length: 5) +
                                    activeAsset['assetIssuer']
                                        .toString()
                                        .substring(activeAsset['assetIssuer']
                                                .toString()
                                                .length -
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
                                    text: activeAsset['assetIssuer'],
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
                    if (curatedAsset['contactEmail'].toString().isNotEmpty) ...[
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
                                  curatedAsset['contactEmail'],
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
                  SizedBox(height: 2),
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
                    child: Image.network(
                      activeAsset["imageUrl"],
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
                  if (activeAsset['assetIssuer'].toString().isNotEmpty) ...[
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
                                truncate(activeAsset['assetIssuer'],
                                        length: 5) +
                                    activeAsset['assetIssuer']
                                        .toString()
                                        .substring(activeAsset['assetIssuer']
                                                .toString()
                                                .length -
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
                                    text: activeAsset['assetIssuer'],
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
