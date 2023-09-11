import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_html/flutter_html.dart';
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

import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SubscriptionPlanBenefits extends StatefulWidget {
  const SubscriptionPlanBenefits({Key? key}) : super(key: key);

  @override
  State<SubscriptionPlanBenefits> createState() =>
      _SubscriptionPlanBenefitsState();
}

class _SubscriptionPlanBenefitsState extends State<SubscriptionPlanBenefits> {
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
          '',
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    '${patronInfo.patronPackage.capitalizeFirst!} ${"patron".tr()}',
                    style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 22,
                      fontFamily: fontsemibold,
                    ),
                  ),
                  Image.asset(
                    patronInfo.logo,
                    height: 30,
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Html(data: patronInfo.description, style: {
                  "*": Style(
                      color: notifier.getbluewhitecolor,
                      fontSize: FontSize.large,
                      lineHeight: LineHeight.number(1.2),
                      wordSpacing: 1.2,
                      textAlign: TextAlign.justify),
                  "h1, h2, h3, h4": Style(
                    fontFamily: fontbody,
                    fontSize: FontSize.large,
                    fontWeight: FontWeight.w800,
                  ),
                }),
              ),
              SizedBox(height: height / 30),
              tile(
                  '\$${patronInfo.patronTiers[0].price} per month',
                  'Billed per month',
                  appState.userInfo?.patronMembership?.patronTierId!
                          .toLowerCase() ==
                      patronInfo.patronTiers[0].tier.toLowerCase()),
              tile(
                  '\$${patronInfo.patronTiers[1].price} per year',
                  'Billed annually',
                  appState.userInfo?.patronMembership?.patronTierId!
                          .toLowerCase() ==
                      patronInfo.patronTiers[1].tier.toLowerCase()),
              tile(
                  '\$${patronInfo.patronTiers[2].price} per lifetime',
                  'Billed once',
                  appState.userInfo?.patronMembership?.patronTierId!
                          .toLowerCase() ==
                      patronInfo.patronTiers[2].tier.toLowerCase()),
              SizedBox(height: height / 10),
            ],
          ),
        ),
      ),
    );
  }

  Widget tile(String main, String sub, bool showButton) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 10.0, vertical: 12.0),
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
                            Container(
                              width: width / 1.5,
                              child: Text(
                                main,
                                style: TextStyle(
                                  fontSize: 13,
                                  // fontWeight: FontWeight.w400,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                          ],
                        ),
                        SizedBox(height: 5),
                        Text(
                          sub,
                          style: TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.w400,
                            color: notifier.getgrey,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(height: 5),
                        if (showButton) ...[
                          Row(
                            children: [
                              Button(
                                'currentplan'.tr(),
                                notifier.getbluewhitecolor,
                                notifier.getwihitecolor,
                                onTap: () {},
                                width: 130,
                                height: height / 25,
                              ),
                            ],
                          ),
                        ],
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}
