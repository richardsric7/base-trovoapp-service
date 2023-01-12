import 'dart:convert';
import 'package:collection/collection.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/topDropdowns.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SelectDepositAddress extends StatefulWidget {
  const SelectDepositAddress({Key? key}) : super(key: key);

  @override
  State<SelectDepositAddress> createState() => _SelectDepositAddressState();
}

class _SelectDepositAddressState extends State<SelectDepositAddress>
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
  dynamic selectedNetwork = '';

  List<dynamic> networks = [];

  List<DropdownMenuItem<String>> get networksDropdownItems {
    return networks
        .mapIndexed<DropdownMenuItem<String>>(
          (index, item) => DropdownMenuItem(
              child: Text(
                item['network'].toString(),
                overflow: TextOverflow.ellipsis,
              ),
              value: '${item['depositAddress']}|$index'),
        )
        .toList();
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    print(appState.viewData![SelectDepositAddressViewPageConfig.key]['data']);
    networks = appState.viewData![SelectDepositAddressViewPageConfig.key]
            ['data']
        .toList();
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
    activeAsset = appState.viewData![SelectDepositAddressViewPageConfig.key];

    if (appState.viewData![SelectDepositAddressViewPageConfig.key] != null) {
      selectedAsset = "${getAssetCode(
        appState.viewData![SelectDepositAddressViewPageConfig.key]['assetCode'],
      )}|${getAssetIssuer(
        appState.viewData![SelectDepositAddressViewPageConfig.key]
            ['assetIssuer'],
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
                                      SelectDepositAddressViewPageConfig.key] =
                                  asset;
                              break;
                            }
                          }
                          setState(() {});
                        },
                        selectedWallet: selectedWallet,
                      ),
                    ],
                  ),
                ),
                SizedBox(width: width / 50),
              ]),
        ),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
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
                    'Deposit ${getAssetCode(activeAsset['assetCode'])}',
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 0, 0, 0),
                child: Text(
                  'Please make your deposit to the address displayed below on the selected network.',
                  style: TextStyle(
                      fontSize: 15,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(
                height: 20,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 0, 0, 0),
                child: Text(
                  'Network',
                  style: TextStyle(
                      fontSize: 15,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold),
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Row(
                  children: [
                    Expanded(
                      child: DropdownButtonFormField(
                        isExpanded: true,
                        dropdownColor: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                        decoration: InputDecoration(
                          contentPadding:
                              EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                          enabledBorder: OutlineInputBorder(
                            borderSide: BorderSide.none,
                            borderRadius: BorderRadius.circular(10),
                          ),
                          border: OutlineInputBorder(
                            borderSide: BorderSide.none,
                            borderRadius: BorderRadius.circular(10),
                          ),
                          filled: true,
                          fillColor: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                        ),
                        // value: selectedNetwork,
                        hint: Text(
                          'Select network',
                          style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                          textAlign: TextAlign.end,
                        ),
                        icon: Icon(
                          Icons.keyboard_arrow_down_rounded,
                          color: notifier.getbluewhitecolor,
                        ),
                        elevation: 0,
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontSize: 15,
                            fontFamily: fontbody,
                            fontWeight: FontWeight.w500),
                        onChanged: (newValue) {
                          setState(() {
                            selectedNetwork = newValue!;
                          });
                        },
                        items: networksDropdownItems,
                      ),
                    ),
                  ],
                ),
              ),
              if (selectedNetwork.toString().isNotEmpty) ...[
                SizedBox(
                  height: height / 50,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 0, 0, 0),
                  child: Text(
                    'Deposit Address',
                    style: TextStyle(
                        fontSize: 15,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20, vertical: 10.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Container(
                                    width: 250,
                                    child: Text(
                                      selectedNetwork.toString().split('|')[0],
                                      overflow: TextOverflow.visible,
                                      style: TextStyle(
                                          fontSize: 15,
                                          color: notifier.getbluewhitecolor,
                                          fontFamily: fontbody),
                                    ),
                                  ),
                                  IconButton(
                                    onPressed: () {
                                      Clipboard.setData(
                                        ClipboardData(
                                          text: selectedNetwork,
                                        ),
                                      );
                                      showSnackBar('Deposit address', context);
                                    },
                                    icon: Icon(Icons.copy,
                                        size: 20,
                                        color: notifier.getbluewhitecolor),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      'OR',
                      style: TextStyle(
                          fontSize: 15,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 50,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 0, 0, 0),
                  child: Text(
                    'Scan QR Code',
                    style: TextStyle(
                        fontSize: 15,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(15.0)),
                        color: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      child: Image.memory(base64.decode(networks[int.parse(
                        selectedNetwork.toString().split('|')[1],
                      )]['qrCode']
                          .split(',')
                          .last))),
                ),
                SizedBox(
                  height: height / 10,
                ),
                Button(
                  'Done',
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.replaceAll,
                      page: BottomHomePageConfig,
                    );
                  },
                ),
              ] else ...[
                SizedBox(
                  height: height / 50,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20, vertical: 10.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Container(
                                    width: width / 1.28,
                                    child: Text(
                                      'To get the deposit address, please select the network where you would like to make the deposit.',
                                      textAlign: TextAlign.justify,
                                      style: TextStyle(
                                          fontSize: 15,
                                          color: notifier.getbluewhitecolor,
                                          fontFamily: fontbody),
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
