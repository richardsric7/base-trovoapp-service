import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class CreateSubWalletSuccessView extends StatefulWidget {
  const CreateSubWalletSuccessView({Key? key}) : super(key: key);

  @override
  State<CreateSubWalletSuccessView> createState() =>
      _CreateSubWalletSuccessViewState();
}

class _CreateSubWalletSuccessViewState
    extends State<CreateSubWalletSuccessView> {
  late ColorNotifier notifier;
  late DataProvider appState;

  @override
  Widget build(BuildContext context) {
    appState = Provider.of<DataProvider>(context, listen: true);
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 6),
              Center(
                child: Image.asset("assets/images/success.gif",
                    height: height / 3.3),
              ),
              SizedBox(height: height / 11),
              Text(
                "createsuccess".tr(),
                style: TextStyle(
                    color: notifier.getblck,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Text(
                "youhavecreatedsuccessfully".tr(),
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: fontbody),
              ),
              SizedBox(height: height / 4.3),
              Button(
                "continuee".tr(),
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: BottomHomePageConfig);
                },
              )
            ],
          ),
        ),
      ),
    );
  }
}
