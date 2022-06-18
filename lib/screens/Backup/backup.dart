import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../Models/Wallet.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/Secret.dart';
import '../Auth/fingerprint.dart';

class Backup extends StatefulWidget {
  const Backup({Key? key}) : super(key: key);

  @override
  State<Backup> createState() => _BackupState();
}

class _BackupState extends State<Backup> {
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
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 6),
              Text(
                LanguageEn.secretkey,
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.2,
                child: Text(
                  LanguageEn.youysecrethasbeengenerated,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 15.sp,
                      fontFamily: fontbody),
                ),
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
              // SizedBox(height: height / 50),
              // Container(
              //   width: width / 1.2,
              //   child: Text(
              //     LanguageEn.writeitasfollows,
              //     style: TextStyle(
              //         color: notifier.getgrey,
              //         fontSize: 15.sp,
              //         fontFamily: fontbody),
              //   ),
              // ),
              SizedBox(height: height / 20),
              // Padding(
              //   padding: const EdgeInsets.symmetric(horizontal: 30),
              //   child: TextButton(
              //     onPressed: () {
              //       setState(() {
              //         showAll = !showAll;
              //       });
              //     },
              //     child: Text(
              //       showAll
              //           ? LanguageEn.taptorevealsecretkeys
              //           : LanguageEn.taptohidesecretkeys,
              //       style: TextStyle(
              //           color: notifier.getbluecolor,
              //           fontSize: 15.sp,
              //           fontFamily: fontbody),
              //     ),
              //   ),
              // ),
              for (var wallet in getUserWallets()) ...[
                Secret(wallet.alias!, wallet.secretKey!)
              ],

              SizedBox(height: height / 20),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  Navigator.pushReplacement(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const FingerPrint(),
                    ),
                  );
                },
              ),
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
}
