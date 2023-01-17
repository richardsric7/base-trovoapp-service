import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/models/User.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/Secret.dart';

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
                LanguageEn.backupwallet,
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
              SizedBox(height: height / 20),
              Secret(state.activeWallet!.alias!, state.activeWallet!.secretKey!,
                  state.activeWallet!.publicKey!),
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

  gotoNext() async {
    if (state.isFirstTime) {
      state.currentAction =
          PageAction(state: PageState.addPage, page: FingerprintPageConfig);
    } else {
      state.currentAction =
          PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
    }
  }
}
