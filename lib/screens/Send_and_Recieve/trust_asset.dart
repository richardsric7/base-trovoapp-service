import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
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
  var assetBalances;
  Wallet? activeWallet;
  var viewData;
  var walletDetails; // will contain information if we are viewing this page from shared wallet
  late bool isSharedWallet;

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
    activeWallet = appState.activeWallet;
    viewData = appState.viewData![PendingAssetDetailsViewPageConfig.key];

    isSharedWallet = appState.viewData![PendingAssetDetailsViewPageConfig.key]
            ['rel'] ==
        SharedWalletDetailsViewPageConfig.key;
    walletDetails = viewData['walletInfo'];

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
          ),
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
                    LanguageEn.pendingassets,
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
              showNotice(),
              SizedBox(
                height: height / 50,
              ),
              assetInfo(),
              SizedBox(
                height: height / 20,
              ),
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
                          .replaceAll('assetCode', viewData['assetCode'])
                          .replaceAll(
                              'walletAlias',
                              isSharedWallet
                                  ? walletDetails['walletAlias']
                                  : activeWallet!.alias!),
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
                    '${getAssetCode(viewData['assetCode'])} Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
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
                        color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
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

  claimAsset() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "assetCode": viewData['assetCode'],
        "assetIssuer": viewData['assetIssuer'],
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePutRequest(
        uri: isSharedWallet
            ? '/v1/shared-access/users/actions/claim-asset'
            : '/v1/users/actions/claim-asset',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: isSharedWallet
            ? walletDetails['walletPublicKey']
            : activeWallet!.publicKey!,
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

      if (isSharedWallet) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makePutRequest(
        uri: isSharedWallet
            ? '/v1/shared-access/users/actions/claim-asset'
            : '/v1/users/actions/claim-asset',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: isSharedWallet
            ? walletDetails['walletPublicKey']
            : activeWallet!.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        await updateUserInfo(activeWallet!.signer!, appState.secretKeys[0],
            activeWallet!.publicKey!, userInfo.username, appState);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': LanguageEn.success,
          'message': LanguageEn.trustassetsuccess
              .replaceAll('asset', viewData['assetCode']),
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
        "assetCode": viewData['assetCode'],
        "assetIssuer": viewData['assetIssuer'],
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makeDeleteRequest(
        uri: isSharedWallet
            ? '/v1/shared-access/users/actions/reject-asset'
            : '/v1/users/actions/reject-asset',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: isSharedWallet
            ? walletDetails['walletPublicKey']
            : activeWallet!.publicKey!,
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

      if (isSharedWallet) {
        responseBody['commit'] = 1;
      }

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makeDeleteRequest(
        uri: isSharedWallet
            ? '/v1/shared-access/users/actions/reject-asset'
            : '/v1/users/actions/reject-asset',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: isSharedWallet
            ? walletDetails['walletPublicKey']
            : activeWallet!.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        await updateUserInfo(activeWallet!.signer!, appState.secretKeys[0],
            activeWallet!.publicKey!, userInfo.username, appState);
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': isSharedWallet
              ? 'Request submitted'
              : 'asset successfully rejected'
                  .replaceAll('asset', viewData['assetCode']),
          'message': isSharedWallet
              ? 'Your request to reject asset has been successfully submitted. This transaction will be completed when it gets the required number of approvals by those who have approver access on this wallet.'
                  .replaceAll('asset', viewData['assetCode'])
              : 'You have successfully rejected this asset. Your wallet will not hold this asset.',
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
