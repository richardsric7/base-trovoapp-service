import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart' hide Trans;
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/curated_asset.dart';
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

class OptInAsset extends StatefulWidget {
  const OptInAsset({Key? key}) : super(key: key);

  @override
  State<OptInAsset> createState() => _OptInAssetState();
}

class _OptInAssetState extends State<OptInAsset> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet wallet;
  late CuratedAsset asset;
  bool hasInfo = true;

  @override
  void initState() {
    super.initState();

    appState = Provider.of<DataProvider>(context, listen: false);
    userInfo = appState.userInfo!;
    wallet = userInfo.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    var result = appState.curatedSwapList.firstWhereOrNull(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    if (result == null) {
      asset = CuratedAsset(
        assetIssuer: appState.viewData!['assetIssuer'],
        assetCode: appState.viewData!['assetCode'],
      );
      hasInfo = false;
    } else {
      asset = result;
    }
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
          '${"add".tr()} [${asset.assetCode}]',
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              showNotice(),
              SizedBox(
                height: height / 50,
              ),
              hasInfo ? assetInfo() : customAssetInfo(),
              SizedBox(
                height: height / 20,
              ),
              if (wallet.canInitiate) ...[
                Button(
                  "addasset".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: optInAsset,
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
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    width: width / 1.3,
                    child: Text(
                      "optininfo"
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
                  SizedBox(
                    height: height / 50.0,
                  ),
                  Container(
                    width: width / 1.3,
                    child: Text(
                      wallet.canInitiate
                          ? "optininfo2"
                              .tr()
                              .replaceAll('walletAlias', wallet.alias!)
                          : "notenoughpermission"
                              .tr()
                              .replaceAll('walletAlias', wallet.alias!),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: wallet.canInitiate
                            ? notifier.getbluewhitecolor
                            : Colors.red[400]!,
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
        constraints: BoxConstraints(minHeight: height / 2.5),
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
                    '${getAssetCode(asset.assetCode)} ${"token".tr()}',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  if (asset.imageUrl != null) ...[
                    SizedBox(
                      height: height / 50.0,
                    ),
                    Container(
                      width: width / 1.3,
                      child: Image.network(
                        asset.imageUrl!,
                        height: 50,
                        width: 50,
                        errorBuilder: (context, error, stackTrace) {
                          return Image.asset(
                            'assets/images/trovo.png',
                            height: 50,
                            width: 50,
                          );
                        },
                      ),
                    ),
                  ],
                  SizedBox(
                    height: height / 50.0,
                  ),
                  Text(
                    asset.website!,
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
                      asset.description!,
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
                  if (asset.assetIssuer.toString().isNotEmpty) ...[
                    Text(
                      "issuerpubkey".tr(),
                      style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(
                      width: width / 1.3,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: [
                          SizedBox(
                            width: width / 20,
                          ),
                          Expanded(
                            flex: 3,
                            child: Padding(
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 20.0),
                              child: Text(
                                truncate(asset.assetIssuer!, length: 5) +
                                    asset.assetIssuer!.toString().substring(
                                        asset.assetIssuer!.toString().length -
                                            5),
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
                                    text: asset.assetIssuer!,
                                  ),
                                ),
                                showSnackBar("issuerpubkey".tr(), context),
                              },
                              icon: Icon(Icons.copy),
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            width: width / 20,
                          ),
                        ],
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    if (asset.contactEmail!.toString().isNotEmpty) ...[
                      Text(
                        "contactemail".tr(),
                        style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.w600,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold),
                      ),
                      SizedBox(
                        width: width / 1.3,
                        child: Row(
                          children: [
                            Expanded(
                              flex: 3,
                              child: Padding(
                                padding: const EdgeInsets.symmetric(
                                    horizontal: 20.0),
                                child: Text(
                                  asset.contactEmail!,
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontWeight: FontWeight.w500,
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 15.sp,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              ),
                            )
                          ],
                        ),
                      ),
                    ],
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

  Widget customAssetInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        constraints: BoxConstraints(minHeight: height / 2.5),
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
                    '${getAssetCode(asset.assetCode!)} ${"token".tr()}',
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
                    child: Image.asset(
                      'assets/images/trovo.png',
                      height: 50,
                      width: 50,
                    ),
                  ),
                  SizedBox(
                    height: height / 50.0,
                  ),
                  if (asset.assetIssuer!.toString().isNotEmpty) ...[
                    Text(
                      "issuerpubkey".tr(),
                      style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(
                      width: width / 1.3,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Padding(
                            padding:
                                const EdgeInsets.symmetric(horizontal: 20.0),
                            child: Text(
                              truncatePublicKey(asset.assetIssuer!),
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
                              showSnackBar("issuerpubkey".tr(), context),
                            },
                            icon: Icon(Icons.copy),
                            color: notifier.getbluewhitecolor,
                          ),
                        ],
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
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

  optInAsset() async {
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

      Map responseData = await makePostRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/asset/opt-in'
            : '/v1/users/asset/opt-in',
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

      Map responseData = await makePostRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/users/asset/opt-in'
            : '/v1/users/asset/opt-in',
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
          'title': "success".tr(),
          'message': wallet.isSharedWalletAndCanInitiate
              ? "optinassetsuccessshared"
                  .tr()
                  .replaceAll('asset', asset.assetCode!)
              : "optinassetsuccess".tr().replaceAll('asset', asset.assetCode!),
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
    appState.viewData![OptInAssetViewPageConfig.key] = null;
  }
}
