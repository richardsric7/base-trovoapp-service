import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/Secret.dart';

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
        appBar: CustomAppBar(
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
                LanguageEn.backupwallets,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  LanguageEn.writeitdown,
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
                    LanguageEn.maynotbedisplayedagain,
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
                  LanguageEn.writeitasfollows,
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
                LanguageEn.continuee,
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
      print(secret);
      Account account = TrovoWalletSDK().parseSecretKey(secret);
      var wlt = user.wallets!
          .firstWhere((wallet) => wallet.publicKey == account.publicKey);
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
