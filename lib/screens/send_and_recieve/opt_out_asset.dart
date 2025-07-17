import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/cache.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
    wallet = userInfo.getWallet(appState.viewData!['walletPublicKey']);

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
          '${"remove".tr()} [${asset.assetCode}]',
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(height: height / 50),
              if (wallet.canInitiate && !hasAvailableBalance) ...[
                showNotice(),
              ] else ...[
                showBurnNotice(),
              ],
              SizedBox(height: height / 20),
              if (wallet.canInitiate && !hasAvailableBalance) ...[
                Button(
                  "removeasset".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: optOutAsset,
                ),
              ] else ...[
                Button(
                  "back".tr(),
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
              padding: const EdgeInsets.symmetric(
                horizontal: 20.0,
                vertical: 15.0,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    width: width / 1.3,
                    child: Text(
                      "optoutinfo".tr().replaceAll(
                        'assetCode',
                        asset.assetCode!,
                      ),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50.0),
                  if (wallet.canInitiate) ...[
                    Container(
                      width: width / 1.3,
                      child: Text(
                        "optoutinfo2"
                            .tr()
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
                        "notenoughpermission".tr().replaceAll(
                          'walletAlias',
                          wallet.alias!,
                        ),
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
              padding: const EdgeInsets.symmetric(
                horizontal: 20.0,
                vertical: 15.0,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    width: width / 1.3,
                    child: Text(
                      "burninfo"
                          .tr()
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
                                  asset.assetIssuer!.toString().length - 5,
                                ),
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
                              ClipboardData(text: asset.assetIssuer!),
                            ),
                            showSnackBar("issuerpubkey".tr(), context),
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

      Map responseData = await makeDeleteRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/asset/opt-out'
            : '/v1/users/asset/opt-out',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      if (responseData['statusCode'] == 200 ||
          responseData['statusCode'] == 202) {
        completeClaimAsset(responseData['data']);
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
        hideLoader(context);
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
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

      responseBody['transactionSignature'] = signature;

      if (wallet.isSharedWalletAndCanInitiate) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      Map responseData = await makeDeleteRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/asset/opt-out'
            : '/v1/users/asset/opt-out',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

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
          'title': "success".tr(),
          'message': wallet.isSharedWalletAndCanInitiate
              ? "optoutassetsuccessshared".tr().replaceAll(
                  'asset',
                  asset.assetCode!,
                )
              : "optoutassetsuccess".tr().replaceAll('asset', asset.assetCode!),
          'useOnDone': true,
          'onDone': () {
            appState.currentAction = PageAction(
              state: PageState.addAll,
              pages: [BottomHomePageConfig, OptInOutAssetViewPageConfig],
            );
          },
        };
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: SuccessViewPageConfig,
        );
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }

  @override
  void dispose() {
    super.dispose();
    appState.viewData![OptOutAssetViewPageConfig.key] = null;
  }
}
