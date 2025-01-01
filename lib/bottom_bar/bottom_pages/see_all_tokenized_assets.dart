import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SeeAllTokenizedAssets extends StatefulWidget {
  const SeeAllTokenizedAssets({Key? key}) : super(key: key);

  @override
  State<SeeAllTokenizedAssets> createState() => _SeeAllTokenizedAssets();
}

class _SeeAllTokenizedAssets extends State<SeeAllTokenizedAssets>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late List<Wallet> wallets;
  late DataProvider appState;
  late Future<List<TokenizedAsset>> listOfTokenizations;
  var viewData;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    viewData = appState.viewData!;
    listOfTokenizations = fetchTokenizationList();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    wallets = appState.userInfo!.wallets!;
    inspect(appState.viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(context, notifier.getwihitecolor, "Primary Offers",
                notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10),
            child: Column(
              children: [
                primaryOffers(),
                SizedBox(
                  height: height / 20,
                ),
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget primaryOffers() {
    return SingleChildScrollView(
      child: Column(
        children: [
          FutureBuilder<List<TokenizedAsset>>(
            future: listOfTokenizations,
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return SizedBox(
                  height: height / 2,
                  child: Center(
                    child: CircularProgressIndicator(
                      backgroundColor: notifier.getbluecolor,
                      valueColor: new AlwaysStoppedAnimation<Color>(
                        notifier.getgreencolor,
                      ),
                      strokeWidth: 3.0,
                    ),
                  ),
                );
              } else if (snapshot.connectionState == ConnectionState.done) {
                if (snapshot.hasError) {
                  return Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: SizedBox(
                      height: height / 6,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "somethingwentwrong".tr(),
                            textAlign: TextAlign.center,
                            style: TextStyle(
                                fontSize: 16,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody),
                          ),
                          ElevatedButton(
                            onPressed: () {
                              setState(() {
                                listOfTokenizations = fetchTokenizationList();
                              });
                            },
                            style: ButtonStyle(
                              backgroundColor: MaterialStateProperty.all<Color>(
                                  notifier.getbluecolor!),
                            ),
                            child: Text(
                              "retry".tr(),
                              style: TextStyle(
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                } else if (snapshot.hasData) {
                  var records = snapshot.data!;
                  return Column(
                    children: [
                      if (records.isNotEmpty) ...[
                        for (var item in records) ...[
                          GestureDetector(
                            onTap: () {
                              appState.tokenizedAsset = item;
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: TokenizedAssetDetailViewPageConfig,
                              );
                            },
                            child: assetTile(item),
                          ),
                        ],
                      ] else ...[
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: SizedBox(
                            height: height / 6,
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  "nothingtoshowhere2".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
                                ),
                              ],
                            ),
                          ),
                        )
                      ]
                    ],
                  );
                }
              }
              return Text(
                '',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.bold,
                  fontFamily: fontsemibold,
                ),
              );
            },
          ),
        ],
      ),
    );
  }

  Future<List<TokenizedAsset>> fetchTokenizationList() async {
    List<TokenizedAsset> tokenizedAssets = [];
    if (appState.tempTokenizedAssetList.isEmpty) {
      var savedAssets = await StoreData().storeGetData('tokenizedAsset');
      if (savedAssets != null) {
        for (int i = 0; i < savedAssets.length; i++) {
          print(savedAssets[i]);
          var a = TokenizedAsset().deserializeJson(savedAssets[i]);
          a.usdPrice = 1.47;
          a.assetIssuer = a.walletToHoldAssetsNotForSale ?? '';
          a.pricePerToken = (double.parse(a.assetCurrentValue.toString()) /
              a.numberOfTokenToBeIssued!);
          tokenizedAssets.add(a);
        }
      }
    }
    appState.tempTokenizedAssetList = tokenizedAssets;
    inspect(tokenizedAssets);
    return tokenizedAssets;
  }

  Widget assetTile(TokenizedAsset asset) {
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 5),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
          title: Row(
            children: [
              if (asset.assetLogo != null) ...[
                Image.memory(
                  base64Decode(asset.assetLogo!),
                  height: 35,
                  width: 35,
                  errorBuilder: (context, error, stackTrace) {
                    return Image.asset(
                      'assets/images/trovo.png',
                      height: 35,
                      width: 35,
                    );
                  },
                ),
              ] else ...[
                Image.asset(
                  'assets/images/trovo.png',
                  height: 35,
                  width: 35,
                ),
              ],
              SizedBox(width: 20),
              Container(
                width: width / 3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${asset.assetName!} (${asset.assetCode})',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        asset.assetSector!,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          trailing: ElevatedButton(
            onPressed: () async {
              asset.isSubscribed ?? false
                  ? showUnSubscribePopup(
                      context,
                      assetCode: asset.assetCode!,
                      onDone: (walletPublicKey) {
                        setState(() {
                          var wallet = wallets
                              .where((wallet) =>
                                  wallet.publicKey == walletPublicKey)
                              .first;
                          wallet.tokenizedAssets!.removeWhere((a) =>
                              a.assetCode == asset.assetCode &&
                              a.assetIssuer == asset.assetIssuer);
                          asset.isSubscribed = false;
                        });
                      },
                      dropdownItems: getUnsubscribableWallets(asset),
                    )
                  : showSubscribePopup(
                      context,
                      assetCode: asset.assetCode!,
                      onDone: (walletPublicKey) async {
                        setState(() {
                          var wallet = wallets
                              .where((wallet) =>
                                  wallet.publicKey == walletPublicKey)
                              .first;
                          wallet.tokenizedAssets != null
                              ? wallet.tokenizedAssets!.add(asset)
                              : wallet.tokenizedAssets = [asset];

                          asset.isSubscribed = true;
                        });
                      },
                      dropdownItems: getStandardWallets,
                    );
            },
            style: ButtonStyle(
              padding: MaterialStateProperty.all(
                EdgeInsets.symmetric(vertical: 0, horizontal: 6),
              ),
              overlayColor:
                  MaterialStateProperty.all<Color>(notifier.getsplashgrey),
              backgroundColor: MaterialStateProperty.all<Color>(
                asset.isSubscribed ?? false
                    ? notifier.getbluewhitecolor
                    : notifier.getwihitecolor,
              ),
              side: MaterialStateProperty.all(
                BorderSide(
                    color: notifier.getbluewhitecolor,
                    width: 1,
                    style: BorderStyle.solid),
              ),
              shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                const RoundedRectangleBorder(
                  borderRadius: BorderRadius.all(
                    Radius.circular(10),
                  ),
                ),
              ),
            ),
            child: Container(
              width: width / 3.7,
              child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    if (asset.isSubscribed ?? false) ...[
                      Text(
                        'Reserved',
                        style: TextStyle(
                          fontFamily: fontsemibold,
                          fontSize: 12,
                          color: asset.isSubscribed ?? false
                              ? notifier.getwihitecolor
                              : notifier.getbluewhitecolor,
                        ),
                      ),
                      Icon(
                        Icons.check_circle_rounded,
                        size: 20,
                        color: asset.isSubscribed ?? false
                            ? notifier.getwihitecolor
                            : notifier.getbluewhitecolor,
                      ),
                    ] else ...[
                      Text(
                        'Reserve',
                        style: TextStyle(
                          fontFamily: fontsemibold,
                          fontSize: 12,
                          color: asset.isSubscribed ?? false
                              ? notifier.getwihitecolor
                              : notifier.getbluewhitecolor,
                        ),
                      ),
                      Icon(
                        Icons.add_circle_rounded,
                        size: 20,
                        color: asset.isSubscribed ?? false
                            ? notifier.getwihitecolor
                            : notifier.getbluewhitecolor,
                      ),
                    ]
                  ]),
            ),
          ),
        ),
      ),
    );
  }

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet.publicKey));
    });
    return wallets;
  }

  List<DropdownMenuItem<String>> getUnsubscribableWallets(
      TokenizedAsset asset) {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.where((wallet) {
      return wallet.tokenizedAssets != null &&
          wallet.tokenizedAssets!
              .where((a) =>
                  a.assetName == asset.assetName &&
                  a.assetCode == asset.assetCode)
              .isNotEmpty;
    }).forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet.publicKey));
    });
    return wallets;
  }
}
