import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
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
  List<Wallet>? wallets;
  Wallet? activeWallet;
  String activeAsset = '';
  var claimedAssets;

  dynamic selectedWallet = '';
  dynamic selectedAsset = '';

  List<DropdownMenuItem<String>> get assetDropdownItems {
    List<DropdownMenuItem<String>> menuItems = [];
    for (var asset in claimedAssets) {
      var value = asset['assetCode'].toString().isEmpty
          ? "XBN"
          : asset['assetCode'].toString();
      menuItems.add(DropdownMenuItem(
          child: Text(
            value,
            overflow: TextOverflow.ellipsis,
          ),
          value: value));
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
    // TODO: implement initState
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
    activeAsset =
        appState.viewData![AssetDetailsViewPageConfig.key]['assetCode'];
    selectedAsset =
        appState.viewData![AssetDetailsViewPageConfig.key]['assetCode'];
    selectedWallet = activeWallet!.publicKey;
    claimedAssets = assetBalances[activeWallet!.publicKey]['claimed'];

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
                              print('this is new value: $newValue');
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
                              setState(() {
                                print('changing active asset to: $newValue');
                                selectedAsset = newValue!;
                                appState.viewData = {
                                  AssetDetailsViewPageConfig.key: {
                                    'assetCode': newValue.toString(),
                                  }
                                };
                              });
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
                    activeAsset,
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              walletSlides(),
              SizedBox(
                height: height / 30,
              ),
              assetInfo(),
              SizedBox(
                height: height / 20,
              ),
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
        actionButton("assets/images/send.svg", 'Send', () {
          appState.viewData![SendAssetViewPageConfig.key] = {
            'assetCode': activeAsset,
          };

          print(appState.viewData);
          appState.currentAction = PageAction(
            state: PageState.addPage,
            page: SendAssetViewPageConfig,
          );
        }),
        actionButton("assets/images/recieve.svg", 'Recieve', () {
          print('fuck you 2');
        }),
      ],
    );
  }

  Widget actionButton(iconUrl, actionText, action) {
    return ElevatedButton(
      onPressed: action,
      style: ButtonStyle(
        backgroundColor:
            MaterialStateProperty.all<Color>(notifier.getbluecolor!),
        shape: MaterialStateProperty.all<RoundedRectangleBorder>(
          const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(
              Radius.circular(15),
            ),
          ),
        ),
      ),
      child: Container(
        width: width / 3.9,
        height: height / 10,
        child: Padding(
          padding: const EdgeInsets.symmetric(
            vertical: 8.0,
          ),
          child: Column(
            children: [
              SvgPicture.asset(
                iconUrl,
                width: width / 8,
              ),
              Text(
                actionText,
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: notifier.getwihitecolor,
                  fontFamily: fontsemibold,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget walletSlides() {
    var colors = [notifier.getbluecolor, Colors.red, Colors.green];
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: colors[0],
          // color: colors[i - 1],
        ),
        child: Stack(children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
                child: Image.asset('assets/images/trovo_white.png'),
              ),
            ],
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  activeWallet!.alias!.capitalizeFirst!,
                  style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    Text(
                      LanguageEn.totalbalance,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getwihitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 98.0,
                ),
                Text(
                  '2,082,898 NGN',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: notifier.getwihitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
                SizedBox(height: 2),
                Text(
                  '4,014 USD',
                  style: TextStyle(
                    fontWeight: FontWeight.w300,
                    fontSize: 13,
                    color: notifier.getwihitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ],
            ),
          ),
        ]),
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
          color: notifier.getaddsubwalletgrey,
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
                    '$activeAsset Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluecolor,
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
                        color: notifier.getbluecolor,
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
                      color: notifier.getbluecolor,
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
