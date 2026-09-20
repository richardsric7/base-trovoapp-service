import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';

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
  Map expressedInterests = {};
  Map subscriptions = {};
  late Future<List<TokenizedAsset>> tokenizedAssetListFuture;
  List<TokenizedAsset> tokenizedAssets = [];

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);

    tokenizedAssetListFuture = fetchTokenizationList(
      status: appState.viewData!['rel'],
    );
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
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          appState.viewData!['rel'] == 0
              ? "Primary Offers"
              : "Secondary Listing",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10),
            child: Column(
              children: [
                Container(
                  constraints: BoxConstraints(
                    minWidth: double.infinity,
                    minHeight: 500,
                  ),
                  child: primaryOffers(),
                ),
                SizedBox(height: height / 20),
                Padding(
                  padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom,
                  ),
                ),
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
            future: tokenizedAssetListFuture,
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
                      height: 400,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "somethingwentwrong".tr(),
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 16,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                          ElevatedButton(
                            onPressed: () {
                              setState(() {
                                tokenizedAssetListFuture =
                                    fetchTokenizationList(status: 0);
                              });
                            },
                            style: ButtonStyle(
                              backgroundColor: WidgetStateProperty.all<Color>(
                                notifier.getbluecolor80!,
                              ),
                              foregroundColor: WidgetStateProperty.all<Color>(
                                notifier.getwihitecolor,
                              ),
                            ),
                            child: Text(
                              "retry".tr(),
                              style: TextStyle(fontFamily: fontsemibold),
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                } else if (snapshot.hasData) {
                  tokenizedAssets = snapshot.data!;
                  return Column(
                    children: [
                      if (tokenizedAssets.isNotEmpty) ...[
                        for (var item in tokenizedAssets) ...[
                          GestureDetector(
                            onTap: () {
                              appState.tokenizedAsset = item;
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: TokenizedAssetDetailViewPageConfig,
                              );
                            },
                            child: tokenizedAssetTile(
                              notifier: notifier,
                              asset: item,
                              onSubscribe: () {
                                showSubscribePopup(
                                  context,
                                  asset: item,
                                  onDone: (amount) async {
                                    await subscribeTokenizedAsset(
                                      amount: double.parse(amount),
                                      tokenizedAssetID: item.id!,
                                    );
                                    setState(() {
                                      item.expressedInterest = true;
                                      item.expressedInterestAmount =
                                          double.parse(amount);
                                    });
                                  },
                                );
                              },
                            ),
                          ),
                        ],
                      ] else ...[
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: SizedBox(
                            height: 400,
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                Text(
                                  "nothingtoshowhere2".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontSize: 16,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ],
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

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(
        DropdownMenuItem(
          child: Text(wallet.alias!, overflow: TextOverflow.ellipsis),
          value: wallet.address,
        ),
      );
    });
    return wallets;
  }

  Future<List<TokenizedAsset>> fetchTokenizationList({
    required int status,
  }) async {
    try {
      await fetchExpressedInterests();
      await fetchSubscriptions();
      var uri =
          '/v1/tokenization/list?onlyWithUserPermission=0&salesList=$status';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        address: appState.primaryWallet.signer!,
      );
      if (responseData['statusCode'] == 200) {
        List<TokenizedAsset> tokenizedAssets = [];
        var assets = responseData['data']['records'];
        if (assets != null) {
          for (int i = 0; i < assets.length; i++) {
            var a = TokenizedAsset().deserializeJson(assets[i]);
            if (expressedInterests[a.id] != null) {
              a.expressedInterest = expressedInterests[a.id] != null;
              a.expressedInterestAmount = double.parse(
                expressedInterests[a.id]['amount'].toString(),
              );
            }

            if (subscriptions[a.id] != null) {
              a.isSubscribed = subscriptions[a.id] != null;
              a.subscriptionAmount = double.parse(
                subscriptions[a.id]['amount'].toString(),
              );
            }
            tokenizedAssets.add(a);
          }
        }
        return tokenizedAssets;
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Future<void> fetchExpressedInterests() async {
    try {
      var uri = '/v1/tokenization/expressed-interests';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        address: appState.primaryWallet.signer!,
      );
      if (responseData['statusCode'] == 200) {
        setState(() {
          expressedInterests = {};
          var records = responseData['data']['records'];
          for (var i = 0; i < records.length; i++) {
            expressedInterests[records[i]['tokenizedAssetId']] = records[i];
          }
        });
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Future<void> fetchSubscriptions() async {
    try {
      var uri = '/v1/tokenization/subscriptions';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        address: appState.primaryWallet.signer!,
      );
      if (responseData['statusCode'] == 200) {
        subscriptions = {};
        setState(() {
          var records = responseData['data']['records'];
          for (var i = 0; i < records.length; i++) {
            subscriptions[records[i]['tokenizedAssetId']] = records[i];
          }
        });
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  subscribeTokenizedAsset({
    required double amount,
    required String tokenizedAssetID,
  }) async {
    try {
      showLoader(context);

      String requestBody = jsonEncode({'amount': amount});

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization/expressed-interests/${tokenizedAssetID}',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        address: appState.primaryWallet.address!,
      );

      if (responseData['statusCode'] == 200) {
        hideLoader(context);
      } else {
        hideLoader(context);
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'].toString().isEmpty
              ? responseData['data']['error']
              : responseData['data']['message'],
        );
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}
