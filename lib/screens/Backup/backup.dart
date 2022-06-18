import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Models/User.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/fonts.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
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
              for (var secret in secrets) ...[Secret(user.username!, secret)],

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

  // Widget walletSecret(alias, secret, hiddenText, isHidden) {
  //   return Card(
  //     margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
  //     elevation: 5,
  //     child: ListTile(
  //       title: Padding(
  //         padding: const EdgeInsets.symmetric(vertical: 8.0),
  //         child: Text(
  //           alias,
  //           style: TextStyle(
  //             fontFamily: fontbody,
  //           ),
  //         ),
  //       ),
  //       subtitle: Padding(
  //         padding: const EdgeInsets.symmetric(vertical: 8.0),
  //         child: Text(
  //           isHidden ? hiddenText : secret,
  //           style: TextStyle(
  //             fontFamily: fontbody,
  //           ),
  //         ),
  //       ),
  //       trailing: Row(
  //         mainAxisSize: MainAxisSize.min,
  //         children: [
  //           IconButton(
  //               onPressed: () => {
  //                     Clipboard.setData(
  //                       ClipboardData(text: secret),
  //                     ),
  //                     showSnackBar('Secret', context),
  //                   },
  //               icon: Icon(Icons.copy)),
  //           IconButton(
  //               onPressed: () {},
  //               icon: Icon(
  //                   isHidden ? CupertinoIcons.eye_slash : CupertinoIcons.eye)),
  //         ],
  //       ),
  //     ),
  //   );
  // }
}
