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
import 'package:trovo_wallet/models/patronTier.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SubscriptionPlans extends StatefulWidget {
  const SubscriptionPlans({Key? key}) : super(key: key);

  @override
  State<SubscriptionPlans> createState() => _SubscriptionPlansState();
}

class _SubscriptionPlansState extends State<SubscriptionPlans> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Future<Map> fetchPlansFuture;

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
    fetchPlansFuture = fetchPatronPlans();
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
          "trovopatronplans".tr(),
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
                      "subscriptiontypes".tr(),
                      style: TextStyle(
                          fontSize: 20,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                  ],
                ),
              ),
              FutureBuilder<Map>(
                future: fetchPlansFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return Container(
                      width: width / 1.2,
                      height: height / 1.7,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          CircularProgressIndicator(
                            backgroundColor: notifier.getbluecolor,
                            valueColor: new AlwaysStoppedAnimation<Color>(
                              notifier.getgreencolor,
                            ),
                            strokeWidth: 3.0,
                          ),
                        ],
                      ),
                    );
                  } else if (snapshot.connectionState == ConnectionState.done) {
                    if (snapshot.hasError) {
                      return Container(
                        width: width / 1.2,
                        height: height / 1.7,
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              "somethingwentwrong".tr(),
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            ElevatedButton(
                              onPressed: () {
                                setState(() {
                                  fetchPlansFuture = fetchPatronPlans();
                                });
                              },
                              style: ButtonStyle(
                                backgroundColor:
                                    MaterialStateProperty.all<Color>(
                                        notifier.getbluecolor!),
                              ),
                              child: Text(
                                "retry".tr(),
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                          ],
                        ),
                      );
                    } else if (snapshot.hasData) {
                      var membershipGrades = snapshot.data!['membershipGrades'];
                      var tiers = snapshot.data!['patronTiers'];
                      var patronPackages = snapshot.data!['patronPackages'];
                      var list = <PatronInfo>[];
                      var myset = Set<String>();

                      for (var grade in membershipGrades) {
                        var tierList = <PatronTier>[];
                        for (var tier in tiers) {
                          if (tier['tier'] == grade['patronTier']) {
                            tierList.add(
                              PatronTier(
                                id: tier['id'],
                                tier: tier['tier'],
                                canExpire: tier['canExpire'],
                                inactive: tier['inactive'],
                                price: grade['price'],
                              ),
                            );
                          }
                        }

                        if (myset.add(grade['patronPackage'])) {
                          var package = PatronInfo(
                            id: grade['id'],
                            patronPackage: grade['patronPackage'],
                            patronTiers: tierList,
                            description: '',
                          );

                          for (var p in patronPackages) {
                            if (p['id'] == grade['patronPackage']) {
                              package.description = p['description'];
                            }
                          }

                          list.add(package);
                        } else {
                          list
                              .firstWhere((item) =>
                                  item.patronPackage == grade['patronPackage'])
                              .patronTiers
                              .addAll(
                                tierList,
                              );
                        }
                      }

                      return Column(
                        children: [
                          for (var info in list) ...[
                            SizedBox(
                              height: height / 50,
                            ),
                            planItem(
                              '${info.patronPackage.capitalizeFirst!} ${"patron".tr()}',
                              '\$${info.patronTiers[0].price} ${"permonth".tr()} / \$${info.patronTiers[1].price} per year / \$${info.patronTiers[2].price} lifetime.',
                              '${"subscribeto".tr()} ${info.patronPackage.capitalizeFirst!}',
                              onReadMore: () {
                                appState.viewData = {
                                  'patronInfo': info,
                                };
                                appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page:
                                        SubscriptionPlanBenefitsViewPageConfig);
                              },
                              onTap: () {
                                appState.viewData = {
                                  'patronInfo': info,
                                };
                                appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page:
                                        SubscriptionPlanOptionsViewPageConfig);
                              },
                            ),
                          ],
                          SizedBox(
                            height: height / 10,
                          ),
                        ],
                      );
                    }
                  }

                  return Text(
                    '',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontsemibold,
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

  Widget planItem(String name, String description, String buttonText,
      {required void Function() onTap, required void Function() onReadMore}) {
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
                                fontSize: 19,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                            TextButton(
                              onPressed: onReadMore,
                              style: TextButton.styleFrom(
                                  padding: EdgeInsets.zero,
                                  minimumSize: Size(50, 30),
                                  tapTargetSize:
                                      MaterialTapTargetSize.shrinkWrap,
                                  alignment: Alignment.centerLeft),
                              child: Text(
                                "readmore".tr(),
                                style: TextStyle(
                                  fontStyle: FontStyle.italic,
                                  fontSize: 15,
                                  fontWeight: FontWeight.w400,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                          ],
                        ),
                        SizedBox(
                          height: 20,
                        ),
                        Text(
                          description,
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
          buttonText,
          notifier.getbluecolor,
          wihitecolor,
          onTap: onTap,
        ),
      ],
    );
  }

  Future<Map> fetchPatronPlans() async {
    try {
      Map responseData = await makeGetRequest(
        uri: '/v1/patron',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.publicKey!,
      );

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error("somethingwentwrong".tr());
      }
    } catch (e) {
      return Future.error('${"error".tr()} ${e}');
    }
  }
}
