import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ReceiveAsset extends StatefulWidget {
  const ReceiveAsset({Key? key}) : super(key: key);

  @override
  State<ReceiveAsset> createState() => _ReceiveAssetState();
}

class _ReceiveAssetState extends State<ReceiveAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  List<Wallet>? wallets;
  Wallet? activeWallet;
  var activeAsset;
  var claimedAssets;

  dynamic selectedWallet = '';
  dynamic selectedAsset = '';

  List<DropdownMenuItem<String>> get assetDropdownItems {
    List<DropdownMenuItem<String>> menuItems = [];
    for (var asset in claimedAssets) {
      menuItems.add(DropdownMenuItem(
          child: Text(
            getAssetCode(asset['assetCode']),
            overflow: TextOverflow.ellipsis,
          ),
          value: getAssetIssuer(asset['assetIssuer'])));
    }
    return menuItems;
  }

  List<DropdownMenuItem<String>> get walletDropdownItems {
    return wallets!
        .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
            child: Text(
              wallet.alias!,
              overflow: TextOverflow.ellipsis,
            ),
            value: wallet.publicKey))
        .toList();
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
    wallets = userInfo.wallets!;
    activeWallet = appState.activeWallet;
    selectedWallet = activeWallet!.publicKey;
    claimedAssets = assetBalances[activeWallet!.publicKey]['claimed'];
    activeAsset = appState.viewData![ReceiveAssetViewPageConfig.key];
    selectedAsset = getAssetIssuer(
      appState.viewData![ReceiveAssetViewPageConfig.key]['assetIssuer'],
    );

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
                      Expanded(
                        child: DropdownButtonFormField(
                          isExpanded: true,
                          decoration: InputDecoration(
                            contentPadding: EdgeInsets.symmetric(
                                vertical: 0, horizontal: 20),
                            enabledBorder: OutlineInputBorder(
                              borderSide: BorderSide.none,
                              borderRadius: BorderRadius.circular(20),
                            ),
                            border: OutlineInputBorder(
                              borderSide: BorderSide.none,
                              borderRadius: BorderRadius.circular(20),
                            ),
                            filled: true,
                            fillColor: notifier.getaddsubwalletgrey,
                          ),
                          value: selectedWallet,
                          icon: Icon(
                            Icons.keyboard_arrow_down_rounded,
                            color: notifier.getbluecolor,
                          ),
                          elevation: 0,
                          style: TextStyle(
                              color: notifier.getbluecolor,
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              fontWeight: FontWeight.w500),
                          onChanged: (newValue) {
                            setState(() {
                              selectedWallet = newValue!;
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
                                      ReceiveAssetViewPageConfig.key] = asset;
                                  break;
                                }
                              }
                              appState.activeWallet = wallets!.firstWhere(
                                  (wallet) => wallet.publicKey == newValue);
                            });
                          },
                          items: walletDropdownItems,
                        ),
                      ),
                      SizedBox(
                        width: width / 20,
                      ),
                      Expanded(
                        child: DropdownButtonFormField(
                            decoration: InputDecoration(
                              contentPadding: EdgeInsets.symmetric(
                                  vertical: 0, horizontal: 20),
                              enabledBorder: OutlineInputBorder(
                                borderSide: BorderSide.none,
                                borderRadius: BorderRadius.circular(20),
                              ),
                              border: OutlineInputBorder(
                                borderSide: BorderSide.none,
                                borderRadius: BorderRadius.circular(20),
                              ),
                              filled: true,
                              fillColor: notifier.getaddsubwalletgrey,
                            ),
                            value: selectedAsset,
                            icon: Icon(
                              Icons.keyboard_arrow_down_rounded,
                              color: notifier.getbluecolor,
                            ),
                            style: TextStyle(
                              color: notifier.getbluecolor,
                              fontSize: 15,
                              fontFamily: fontsemibold,
                            ),
                            onChanged: (newValue) {
                              print('this is newValue $newValue');
                              setState(() {
                                print('changing active asset to: $newValue');
                                newValue = newValue == nativeAssetIssuer
                                    ? ''
                                    : newValue;
                                for (var asset in claimedAssets) {
                                  print('this is newValue $newValue');
                                  if (asset['assetIssuer'] == newValue) {
                                    appState.viewData![
                                        ReceiveAssetViewPageConfig.key] = asset;
                                  }
                                }
                              });
                              print(
                                  'this is new viewdata: ${appState.viewData}');
                            },
                            borderRadius: BorderRadius.all(
                              Radius.circular(15),
                            ),
                            items: assetDropdownItems),
                      ),
                      SizedBox(
                        width: width / 20,
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
                    "Receive " + getAssetCode(activeAsset['assetCode']),
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
              showReceivingWallet(),
              SizedBox(
                height: height / 50,
              ),
              showUsername(),
              SizedBox(
                height: height / 50,
              ),
              showPublicKey(),
              SizedBox(
                height: height / 50,
              ),
              showQrCode(),
              SizedBox(
                height: height / 20,
              ),
              Button(
                LanguageEn.dashboard,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: BottomHomePageConfig);
                },
              ),
              SizedBox(height: height / 20),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Padding showReceivingWallet() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20, vertical: 10.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    'Receiving Wallet',
                    style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(height: height / 90),
                  Text(
                    activeWallet!.alias!,
                    style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Padding showPublicKey() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20, vertical: 10.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    'Receive from non Trovo wallet',
                    style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(height: height / 90),
                  Row(
                    children: [
                      Container(
                        width: 250,
                        child: Text(
                          activeWallet!.publicKey!,
                          style: TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w600,
                              color: notifier.getbluecolor,
                              fontFamily: fontsemibold),
                        ),
                      ),
                      IconButton(
                        onPressed: () {
                          Clipboard.setData(
                            ClipboardData(
                              text: activeWallet!.publicKey!,
                            ),
                          );
                          showSnackBar('Public key', context);
                        },
                        icon: Icon(Icons.copy, size: 20),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Padding showUsername() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20, vertical: 10.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    'Receive with Trovo username',
                    style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                  Row(
                    children: [
                      Text(
                        userInfo.username!,
                        style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.w600,
                            color: notifier.getbluecolor,
                            fontFamily: fontsemibold),
                      ),
                      IconButton(
                        onPressed: () {
                          Clipboard.setData(
                            ClipboardData(
                              text: userInfo.username,
                            ),
                          );
                          showSnackBar('Username', context);
                        },
                        icon: Icon(Icons.copy, size: 20),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget showQrCode() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.getaddsubwalletgrey,
          ),
          child: Image.memory(
              base64.decode(activeAsset['qrCode'].split(',').last))),
    );
  }
}
