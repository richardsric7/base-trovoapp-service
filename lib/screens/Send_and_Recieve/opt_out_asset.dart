import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class OptOutAsset extends StatefulWidget {
  const OptOutAsset({Key? key}) : super(key: key);

  @override
  State<OptOutAsset> createState() => _OptOutAssetState();
}

class _OptOutAssetState extends State<OptOutAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet wallet;
  late Asset asset;
  late bool hasAvailableBalance;

  @override
  void initState() {
    super.initState();

    appState = Provider.of<DataProvider>(context, listen: false);
    userInfo = appState.userInfo!;
    wallet = userInfo.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    asset = wallet.claimedAssets!.firstWhere(
      (claimedAsset) =>
          claimedAsset.assetCode == appState.viewData!['assetCode'] &&
          claimedAsset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    hasAvailableBalance = asset.amount! > 0;
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
          'Remove [${asset.assetCode}]',
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                height: height / 50,
              ),
              if (wallet.canInitiate && !hasAvailableBalance) ...[
                showNotice(),
              ] else ...[
                showBurnNotice(),
              ],
              SizedBox(
                height: height / 20,
              ),
              if (wallet.canInitiate && !hasAvailableBalance) ...[
                Button(
                  LanguageEn.removeasset,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: optOutAsset,
                ),
              ] else ...[
                Button(
                  LanguageEn.back,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    Navigator.of(context).pop();
                  },
                ),
              ],
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }

  Widget showNotice() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
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
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    width: width / 1.3,
                    child: Text(
                      LanguageEn.optoutinfo
                          .replaceAll('assetCode', asset.assetCode!),
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
                  if (wallet.canInitiate) ...[
                    Container(
                      width: width / 1.3,
                      child: Text(
                        LanguageEn.optoutinfo2
                            .replaceAll('assetCode', asset.assetCode!)
                            .replaceAll('walletAlias', wallet.alias!),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w400,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ),
                  ] else ...[
                    Container(
                      width: width / 1.3,
                      child: Text(
                        LanguageEn.notenoughpermission
                            .replaceAll('walletAlias', wallet.alias!),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w400,
                          color: Colors.red,
                          fontFamily: fontbody,
                        ),
                      ),
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

  Widget showBurnNotice() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
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
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    width: width / 1.3,
                    child: Text(
                      LanguageEn.burninfo
                          .replaceAll('assetCode', asset.assetCode!)
                          .replaceAll('walletAlias', wallet.alias!)
                          .replaceAll('amount', asset.amount.toString()),
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
                    width: width / 1.3,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Text(
                            truncate(asset.assetIssuer!, length: 5) +
                                asset.assetIssuer!.toString().substring(
                                    asset.assetIssuer!.toString().length - 5),
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontbody,
                            ),
                          ),
                        ),
                        IconButton(
                          padding: EdgeInsets.zero,
                          onPressed: () => {
                            Clipboard.setData(
                              ClipboardData(
                                text: asset.assetIssuer!,
                              ),
                            ),
                            showSnackBar('Issuer public key', context),
                          },
                          icon: Icon(Icons.copy),
                          color: notifier.getbluewhitecolor,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  optOutAsset() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "assetCode": asset.assetCode!,
        "assetIssuer": asset.assetIssuer!,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makeDeleteRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/asset/opt-out'
            : '/v1/users/asset/opt-out',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');

      if (responseData['statusCode'] == 200 ||
          responseData['statusCode'] == 202) {
        completeClaimAsset(responseData['data']);
        // print('sending full data to server.........');
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }

  void completeClaimAsset(responseBody) async {
    try {
      showLoader(context);

      var signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0],
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      print('this is primary sign: $signature');
      responseBody['transactionSignature'] = signature;

      if (wallet.isSharedWalletAndCanInitiate) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makeDeleteRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/asset/opt-out'
            : '/v1/users/asset/opt-out',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        updateUserInfo(
          appState.primaryWallet.signer!,
          appState.secretKeys[0],
          appState.primaryWallet.publicKey!,
          userInfo.username,
          appState,
          forceRefresh: true,
        );
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': LanguageEn.success,
          'message': wallet.isSharedWalletAndCanInitiate
              ? LanguageEn.optoutassetsuccessshared
                  .replaceAll('asset', asset.assetCode!)
              : LanguageEn.optoutassetsuccess
                  .replaceAll('asset', asset.assetCode!),
          'useOnDone': true,
          'onDone': () {
            appState.currentAction = PageAction(
              state: PageState.addAll,
              pages: [BottomHomePageConfig, OptInOutAssetViewPageConfig],
            );
          },
        };
        appState.currentAction =
            PageAction(state: PageState.addPage, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }

    hideLoader(context);
  }

  @override
  void dispose() {
    super.dispose();
    print('disposing...');
    appState.viewData![OptOutAssetViewPageConfig.key] = null;
  }
}
