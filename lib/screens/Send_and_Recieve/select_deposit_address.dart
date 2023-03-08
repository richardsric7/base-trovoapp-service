import 'dart:convert';
import 'package:collection/collection.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/crypto_wallet_address.dart';
import 'package:trovo_wallet/models/curated_asset.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
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
  dynamic selectedNetwork = '';
  late Wallet wallet;
  late Asset? asset;
  late CuratedAsset? curatedAsset;

  List<CryptoWalletDepositAddress> networks = [];

  List<DropdownMenuItem<String>> get networksDropdownItems {
    return networks
        .mapIndexed<DropdownMenuItem<String>>(
          (index, item) => DropdownMenuItem(
              child: Text(
                item.network.toString(),
                overflow: TextOverflow.ellipsis,
              ),
              value: '${item.depositAddress}|$index'),
        )
        .toList();
  }

  @override
  void initState() {
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

    networks = asset!.cryptoWalletDepositAddresses!;

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

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
                    'Deposit ${getAssetCode(asset!.assetCode)}',
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
                      child: Image.network(networks[int.parse(
                        selectedNetwork.toString().split('|')[1],
                      )]
                          .qrCode!)),
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
