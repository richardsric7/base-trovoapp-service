import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart' hide Trans;
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/models/patronInfo.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SubscriptionPlanOptions extends StatefulWidget {
  const SubscriptionPlanOptions({Key? key}) : super(key: key);

  @override
  State<SubscriptionPlanOptions> createState() =>
      _SubscriptionPlanOptionsState();
}

class _SubscriptionPlanOptionsState extends State<SubscriptionPlanOptions> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late PatronInfo patronInfo;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    patronInfo = appState.viewData!['patronInfo'] as PatronInfo;
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          '${patronInfo.patronPackage.capitalizeFirst} ${"options".tr()}',
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '${patronInfo.patronPackage.capitalizeFirst} Patron options ',
                      style: TextStyle(
                          fontSize: 18,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
              ),
              for (var tier in patronInfo.patronTiers.reversed) ...[
                SizedBox(
                  height: height / 50,
                ),
                planOption(
                  '${tier.tier.capitalizeFirst} ${"subscription".tr().toLowerCase()} ${tier.tier == 'LIFETIME' ? '(recommended)' : ''}',
                  '\$${tier.price}',
                  () {
                    // appState.viewData = {
                    //   'patronInfo': patronInfo,
                    //   'selectedTier': tier,
                    // };
                    // appState.currentAction = PageAction(
                    //     state: PageState.addPage,
                    //     page: AuthorizeSubscriptionViewPageConfig);
                    popup(
                      context,
                      title: "comingsoon".tr(),
                      message: "comingsoondetails".tr(),
                      bodyColor: notifier.getbluewhitecolor,
                    );
                  },
                ),
              ]
            ],
          ),
        ),
      ),
    );
  }

  Widget planOption(
    String name,
    String price,
    void Function() onTap,
  ) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              border: Border.all(color: notifier.getbluewhitecolor, width: 1.5),
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 10.0, vertical: 15.0),
              child: Row(
                children: [
                  Container(
                    width: width / 1.2,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        SizedBox(
                          width: width / 50,
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              name,
                              style: TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            )
                          ],
                        ),
                        SizedBox(
                          height: 20,
                        ),
                        Text(
                          price,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(
                          height: 20,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        SizedBox(height: 5),
        Button(
          'Subscribe',
          notifier.getbluecolor,
          wihitecolor,
          onTap: onTap,
        ),
      ],
    );
  }
}
