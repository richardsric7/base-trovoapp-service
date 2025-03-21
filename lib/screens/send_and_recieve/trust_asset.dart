import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:get/get_utils/get_utils.dart' hide Trans;
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
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
import '../../utils/medeiaqury/medeiaqury.dart';

class PendingAssetDetails extends StatefulWidget {
  const PendingAssetDetails({Key? key}) : super(key: key);

  @override
  State<PendingAssetDetails> createState() => _PendingAssetDetailsState();
}

class _PendingAssetDetailsState extends State<PendingAssetDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet wallet;
  late Asset? asset;

  @override
  void initState() {
    super.initState();

    appState = Provider.of<DataProvider>(context, listen: false);
    userInfo = appState.userInfo!;
    wallet = userInfo.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    asset = wallet.unClaimedAssets!.firstWhereOrNull(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      appBar: PreferredSize(
        preferredSize: Size.fromHeight(height / 15),
        child: AppBar(
          centerTitle: true,
          elevation: 0,
          backgroundColor: notifier.getwihitecolor,
          title: Text(
            "pendingassets".tr(),
            style: TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold),
          ),
          leading: GestureDetector(
            onTap: () {
              Navigator.of(context).pop();
            },
            child: Image.asset("assets/images/back.png", scale: 5),
          ),
        ),
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            SizedBox(
              height: height / 20,
            ),
            if (asset != null) ...[
              showNotice(),
              SizedBox(
                height: height / 10,
              ),
              Button(
                "claimasset".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: claimAsset,
              ),
              SizedBox(height: height / 50),
              ButtonOutlined(
                "rejectasset".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  rejectAsset();
                },
              ),
              SizedBox(height: height / 10),
            ] else ...[
              Center(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    "nothingtoshowhere".tr(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ),
              SizedBox(
                height: height / 10,
              ),
              Button(
                "back".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  Navigator.of(context).pop();
                },
              ),
            ],
          ],
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
                      "pendingassetwarning"
                          .tr()
                          .replaceAll('assetCode', asset!.assetCode!)
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
                  SizedBox(
                    height: height / 50.0,
                  ),
                  Container(
                    width: width / 1.3,
                    child: Text(
                      "pendingassetwarning2".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontbody,
                      ),
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

  claimAsset() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "assetCode": asset!.assetCode!,
        "assetIssuer": asset!.assetIssuer!,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePutRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/actions/claim-asset'
            : '/v1/users/actions/claim-asset',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');

      if (responseData['statusCode'] == 202) {
        completeClaimAsset(responseData['data']);
        // print('sending full data to server.........');
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
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

      print('this is primary sign: $signature');
      responseBody['transactionSignature'] = signature;

      if (wallet.isSharedWalletAndCanInitiate) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makePutRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/actions/claim-asset'
            : '/v1/users/actions/claim-asset',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        updateUserInfo(appState.primaryWallet.signer!, appState.secretKeys[0],
            appState.primaryWallet.publicKey!, userInfo.username, appState);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': "success".tr(),
          'message':
              "trustassetsuccess".tr().replaceAll('asset', asset!.assetCode!),
          'useOnDone': true,
          'onDone': () {
            appState.currentAction = appState.returnView ??
                PageAction(
                  state: PageState.addAll,
                  pages: [BottomHomePageConfig],
                );
          },
        };
        appState.currentAction = PageAction(
            state: PageState.replaceAll, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }

  rejectAsset() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "assetCode": asset!.assetCode!,
        "assetIssuer": asset!.assetIssuer!,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makeDeleteRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/actions/reject-asset'
            : '/v1/users/actions/reject-asset',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');

      if (responseData['statusCode'] == 202) {
        completeRejectAsset(responseData['data']);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
      hideLoader(context);
    }
  }

  void completeRejectAsset(responseBody) async {
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
            ? '/v1/shared-access/users/actions/reject-asset'
            : '/v1/users/actions/reject-asset',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        updateUserInfo(appState.primaryWallet.signer!, appState.secretKeys[0],
            appState.primaryWallet.publicKey!, userInfo.username, appState);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': wallet.isSharedWalletAndCanInitiate
              ? 'Request submitted'
              : 'asset successfully rejected'
                  .replaceAll('asset', asset!.assetCode!),
          'message': wallet.isSharedWalletAndCanInitiate
              ? 'Your request to reject asset has been successfully submitted. This transaction will be completed when it gets the required number of approvals by those who have approver access on this wallet.'
                  .replaceAll('asset', asset!.assetCode!)
              : 'You have successfully rejected this asset. Your wallet will not hold this asset.',
          'useOnDone': true,
          'onDone': () {
            appState.currentAction = appState.returnView ??
                PageAction(
                  state: PageState.addAll,
                  pages: [BottomHomePageConfig],
                );
          },
        };
        appState.currentAction = PageAction(
            state: PageState.replaceAll, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }

  @override
  void dispose() {
    super.dispose();
    print('disposing...');
    appState.viewData![PendingAssetDetailsViewPageConfig.key] = null;
  }
}
