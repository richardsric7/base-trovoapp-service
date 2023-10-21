import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart' hide Trans;
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/secret.dart';

class BackupAll extends StatefulWidget {
  const BackupAll({Key? key}) : super(key: key);

  @override
  State<BackupAll> createState() => _BackupAllState();
}

class _BackupAllState extends State<BackupAll> {
  late DataProvider state;
  late UserInfo user;
  late List<String> secrets;

  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    state = Provider.of<DataProvider>(context, listen: false);
    user = state.userInfo!;
    secrets = state.secretKeys;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBarWithoutBanner(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Text(
                "backupwallets".tr(),
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  "writeitdown".tr(),
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
                ),
              ),
              // display only for subwallets
              if (state.activeWallet!.primaryWallet == 0) ...[
                SizedBox(height: height / 50),
                Container(
                  width: width / 1.2,
                  child: Text(
                    "maynotbedisplayedagain".tr(),
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 15.sp,
                        fontFamily: fontbody),
                  ),
                ),
              ],
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  "writeitasfollows".tr(),
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(height: height / 20),
              for (var wallet in getUserWallets()) ...[
                Secret(wallet.alias!, wallet.secretKey!, wallet.publicKey!)
              ],
              SizedBox(height: height / 20),
              Button(
                "continuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  gotoNext();
                },
              ),
              SizedBox(height: height / 20),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  List<Wallet> getUserWallets() {
    var wallets = <Wallet>[];
    secrets.forEach((secret) {
      Account account = TrovoWalletSDK().parseSecretKey(secret);
      var wlt = user.wallets!
          .firstWhereOrNull((wallet) => wallet.publicKey == account.publicKey);
      if (wlt == null) {
        wlt = state.primaryWallet;
      }
      wlt.secretKey = secret;
      wallets.add(wlt);
    });
    return wallets;
  }

  gotoNext() async {
    var isFirstTime = await StoreData().storeGetData('isFirstTime') ?? true;
    if (isFirstTime) {
      state.currentAction =
          PageAction(state: PageState.addPage, page: FingerprintPageConfig);
    } else {
      state.currentAction =
          PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
    }
    state.viewData![EnsurePrivacyPageConfig.key] = null;
  }
}
