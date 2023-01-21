import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/curated_asset.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/wallet.dart';
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

class GenerateDepositAddress extends StatefulWidget {
  const GenerateDepositAddress({Key? key}) : super(key: key);

  @override
  State<GenerateDepositAddress> createState() => _GenerateDepositAddressState();
}

class _GenerateDepositAddressState extends State<GenerateDepositAddress>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late Asset asset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );
    asset = wallet.claimedAssets!.firstWhere(
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
    appState = Provider.of<DataProvider>(context, listen: true);

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
                    'Deposit ${getAssetCode(asset.assetCode!)}',
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: 20,
              ),
              depositInfo(),
              SizedBox(
                height: height / 10,
              ),
              Button(
                'Generate Deposit Address',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  generateDepositAddress();
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget depositInfo() {
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
              child: Container(
                width: width / 1.3,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Image.asset(
                      "assets/images/deposit-pin.png",
                      height: height / 7,
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      'You don’t have any deposit address yet. Please tap on the button below to generate addresses in order to continue with your deposit transaction.',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                          fontSize: 15,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void generateDepositAddress() async {
    try {
      showLoader(context);

      Map responseData = await makePostRequest(
        uri: '/v1/crypto/generate-addresses/${asset.assetCode}',
        body: "",
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      if (responseData['statusCode'] == 200) {
        await updateUserInfo(
          wallet.signer,
          appState.secretKeys[0],
          wallet.publicKey,
          appState.userInfo!.username,
          appState,
        );
        hideLoader(context);
        appState.viewData = {
          'walletPublicKey': wallet.publicKey,
          'assetCode': asset.assetCode,
          'assetIssuer': asset.assetIssuer,
        };

        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: SelectDepositAddressViewPageConfig,
        );
      } else {
        hideLoader(context);
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }
}
