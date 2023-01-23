import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
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

    asset = wallet.unClaimedAssets!.firstWhere(
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
            title: Text(
              LanguageEn.pendingassets,
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
              showNotice(),
              SizedBox(
                height: height / 10,
              ),
              // assetInfo(),
              // SizedBox(
              //   height: height / 20,
              // ),
              if (!wallet.isSharedWallet || wallet.isInitiator) ...[
                Button(
                  LanguageEn.claimasset,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: claimAsset,
                ),
                SizedBox(height: height / 50),
                ButtonOutlined(
                  'Reject asset',
                  notifier.getwihitecolor,
                  notifier.getbluewhitecolor,
                  onTap: () {
                    rejectAsset();
                  },
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
              SizedBox(height: height / 10),
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
                      LanguageEn.pendingassetwarning
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
                  if (!wallet.isSharedWallet || wallet.isInitiator) ...[
                    Container(
                      width: width / 1.3,
                      child: Text(
                        LanguageEn.pendingassetwarning2,
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
                        'You do not have enough permission to claim this asset on [walletAlias].'
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
                    '${getAssetCode(asset!.assetCode)} Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  Text(
                    'www.trovotech.io',
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
                      'TROV token (TROV) is the utility token that powers the Trovotech ecosystem. TROV token is used to access discounts, voting rights, airdrops, NFTs and other community incentives. ',
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
                              truncate(asset!.assetIssuer!, length: 5) +
                                  asset!.assetIssuer.toString().substring(
                                      asset!.assetIssuer.toString().length - 5),
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
                                  text: asset!.assetIssuer!,
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
        uri: wallet.isSharedWallet
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

      if (wallet.isSharedWallet) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makePutRequest(
        uri: wallet.isSharedWallet
            ? '/v1/shared-access/users/actions/claim-asset'
            : '/v1/users/actions/claim-asset',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        await updateUserInfo(
            appState.primaryWallet.signer!,
            appState.secretKeys[0],
            appState.primaryWallet.publicKey!,
            userInfo.username,
            appState);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': LanguageEn.success,
          'message': LanguageEn.trustassetsuccess
              .replaceAll('asset', asset!.assetCode!),
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
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
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
        uri: wallet.isSharedWallet
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
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
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

      if (wallet.isSharedWallet) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makeDeleteRequest(
        uri: wallet.isSharedWallet
            ? '/v1/shared-access/users/actions/reject-asset'
            : '/v1/users/actions/reject-asset',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        await updateUserInfo(
            appState.primaryWallet.signer!,
            appState.secretKeys[0],
            appState.primaryWallet.publicKey!,
            userInfo.username,
            appState);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': wallet.isSharedWallet
              ? 'Request submitted'
              : 'asset successfully rejected'
                  .replaceAll('asset', asset!.assetCode!),
          'message': wallet.isSharedWallet
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
    appState.viewData![PendingAssetDetailsViewPageConfig.key] = null;
  }
}
