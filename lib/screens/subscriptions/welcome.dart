import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:timeago/timeago.dart' as timeago;
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/patronInfo.dart';
import 'package:trovo_wallet/models/patronMembership.dart';
import 'package:trovo_wallet/models/patronSubscriptionLog.dart';
import 'package:trovo_wallet/models/patronTier.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';

import '../../custom_bloc_observer/fonts.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WelcomeSubscriptions extends StatefulWidget {
  const WelcomeSubscriptions({Key? key}) : super(key: key);

  @override
  State<WelcomeSubscriptions> createState() => _WelcomeSubscriptionsState();
}

class _WelcomeSubscriptionsState extends State<WelcomeSubscriptions> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Future<Map> fetchPlansFuture;
  late PatronMembership? patronMembership;
  late PatronTier currentPatronTier;

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
    patronMembership = appState.userInfo?.patronMembership;
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
          "trovopatron".tr(),
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: patronMembership == null ? Welcome() : CurrentPlan(),
        ),
      ),
    );
  }

  /* this part needs total refactoring if not overhauling */
  Widget CurrentPlan() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        FutureBuilder<Map>(
          future: fetchPlansFuture,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return Container(
                width: width,
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
                          backgroundColor: MaterialStateProperty.all<Color>(
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
                var paymentAssets = snapshot.data!['subscriptionPaymentAssets'];
                var logs = snapshot.data!['patronSubscriptionLogs'];
                var patronInfoList = <PatronInfo>[];
                var patronSubscriptionLogs = <PatronSubscriptionLog>[];
                var paymentAssetsList = <Asset>[];
                var myset = Set<String>();

                for (var log in logs) {
                  patronSubscriptionLogs
                      .add(PatronSubscriptionLog.deserializeJson(log));
                }

                for (var asset in paymentAssets) {
                  paymentAssetsList.add(
                    Asset(
                      assetCode: asset['assetCode'],
                      assetIssuer: asset['assetIssuer'],
                    ),
                  );
                }

                for (var grade in membershipGrades) {
                  var tierList = <PatronTier>[];
                  for (var tier in tiers) {
                    if (tier['tier'] == grade['patronTier']) {
                      var t = PatronTier(
                        id: tier['id'],
                        tier: tier['tier'],
                        canExpire: tier['canExpire'],
                        inactive: tier['inactive'],
                        price: grade['price'],
                      );
                      tierList.add(t);
                    }
                  }

                  if (grade['patronPackage'] ==
                          patronMembership?.patronPackageId &&
                      grade['patronTier'] == patronMembership?.patronTierId)
                    patronMembership?.price = grade['price'];

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

                    patronInfoList.add(package);
                  } else {
                    patronInfoList
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
                    SizedBox(height: height / 50),
                    Padding(
                      padding: EdgeInsets.symmetric(horizontal: 15),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            "currentplan".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                                color: notifier.getgrey,
                                fontSize: 13.sp,
                                fontFamily: fontbody),
                          ),
                          Row(
                            children: [
                              Image.asset("assets/images/gold.png",
                                  height: height / 30),
                              SizedBox(
                                width: 5,
                              ),
                              Text(
                                "currentpatronplan".tr(args: [
                                  patronMembership?.patronPackageId
                                          .toString()
                                          .capitalizeFirst ??
                                      'null',
                                  patronMembership?.patronTierId
                                          .toString()
                                          .capitalizeFirst ??
                                      'null',
                                ]),
                                textAlign: TextAlign.start,
                                style: TextStyle(
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 17.sp,
                                    fontFamily: fontsemibold),
                              ),
                            ],
                          ),
                          SizedBox(height: height / 70),
                          Text(
                            "billingamount".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                                color: notifier.getgrey,
                                fontSize: 13.sp,
                                fontFamily: fontbody),
                          ),
                          Text(
                            patronMembership?.patronTierId!.toLowerCase() ==
                                    'monthly'
                                ? "amountpermonth".tr(args: [
                                    patronMembership?.price.toString() ?? '',
                                  ])
                                : patronMembership?.patronTierId!
                                            .toLowerCase() ==
                                        'annual'
                                    ? "amountperyear".tr(args: [
                                        patronMembership?.price.toString() ??
                                            '',
                                      ])
                                    : "lifetimenobilling".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontSize: 17.sp,
                                fontFamily: fontsemibold),
                          ),
                          SizedBox(height: height / 70),
                          Text(
                            "nextbilldate".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                                color: notifier.getgrey,
                                fontSize: 13.sp,
                                fontFamily: fontbody),
                          ),
                          Text(
                            DateFormat('dd MMM, y')
                                .format(patronMembership!.validTill!),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontSize: 17.sp,
                                fontFamily: fontsemibold),
                          ),
                        ],
                      ),
                    ),
                    SizedBox(height: height / 10),
                    ButtonOutlined(
                      "cancelsubscription".tr(),
                      notifier.getwihitecolor,
                      notifier.getbluewhitecolor,
                      onTap: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: SubscriptionPlansViewPageConfig);
                      },
                    ),
                    SizedBox(height: height / 50),
                    Button(
                      "viewpatronplans".tr(),
                      notifier.getbluecolor,
                      wihitecolor,
                      onTap: () {
                        appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: SubscriptionPlansViewPageConfig);
                      },
                    ),
                    SizedBox(height: height / 15),
                    Padding(
                      padding: EdgeInsets.symmetric(horizontal: 18),
                      child: Row(
                        children: [
                          Text(
                            "activity".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontSize: 17.sp,
                                fontFamily: fontsemibold),
                          ),
                        ],
                      ),
                    ),
                    for (var log in patronSubscriptionLogs) ...[
                      logItem(log),
                    ]
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
    );
  }

  Widget Welcome() {
    return Column(
      children: [
        SizedBox(height: height / 8),
        Image.asset("assets/images/unlock.png", height: height / 3.7),
        SizedBox(height: height / 10),
        Padding(
          padding: EdgeInsets.symmetric(horizontal: width / 15),
          child: Column(children: [
            Text(
              "unlockfulltrovotechpotential".tr(),
              textAlign: TextAlign.center,
              style: TextStyle(
                  color: notifier.getblck,
                  fontSize: 22.sp,
                  fontFamily: fontsemibold),
            ),
          ]),
        ),
        SizedBox(height: height / 15),
        Padding(
          padding: EdgeInsets.symmetric(horizontal: width / 15),
          child: Column(children: [
            Text(
              "youarenotatrovopatron".tr(),
              textAlign: TextAlign.center,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontSize: 17.sp,
                  fontFamily: fontbody),
            ),
          ]),
        ),
        SizedBox(height: height / 10),
        Button(
          "viewpatronplans".tr(),
          notifier.getbluecolor,
          wihitecolor,
          onTap: () {
            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: SubscriptionPlansViewPageConfig);
          },
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
        print(
            '=======================> patron response: ${responseData['data']}');
        return responseData['data'];
      } else {
        return Future.error("somethingwentwrong".tr());
      }
    } catch (e) {
      return Future.error('${"error".tr()} ${e}');
    }
  }

  Widget logItem(PatronSubscriptionLog log) {
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
                                "subscribedto".tr(args: [
                                  log.patronPackageId,
                                  log.patronTierId.capitalizeFirst!
                                ]),
                                style: TextStyle(
                                  fontSize: 13,
                                  // fontWeight: FontWeight.w400,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                            Text(
                              '\$${patronMembership?.price}',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ],
                        ),
                        SizedBox(
                          height: 20,
                        ),
                        Text(
                          timeago.format(log.createdAt),
                          style: TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                        ),
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
