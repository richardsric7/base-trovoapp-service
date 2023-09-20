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
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/patronInfo.dart';
import 'package:trovo_wallet/models/patronMembership.dart';
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
  int currentTab = 0;
  late PatronMembership? patronMembership;

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
    patronMembership = appState.userInfo?.patronMembership;
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
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  MyTab('Monthly', 0, currentTab == 0),
                  MyTab('Yearly', 1, currentTab == 1),
                  MyTab('Lifetime', 2, currentTab == 2),
                ],
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
                      var paymentAssets =
                          snapshot.data!['subscriptionPaymentAssets'];
                      var patronInfoList = <PatronInfo>[];
                      var paymentAssetsList = <Asset>[];
                      var myset = Set<String>();

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
                            tierList.add(
                              PatronTier(
                                id: grade['id'],
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
                            packageList: [],
                            packageListTitle: '',
                            logo: '',
                          );

                          for (var p in patronPackages) {
                            if (p['id'] == grade['patronPackage']) {
                              package.description = p['description'];
                              package.packageListTitle = p['packageListTitle'];
                              package.packageList =
                                  p['packageList'].toString().split('|');
                              package.logo =
                                  p['id'].toString().toLowerCase() == 'gold'
                                      ? "assets/images/gold.png"
                                      : p['id'].toString().toLowerCase() ==
                                              'diamond'
                                          ? "assets/images/diamond.png"
                                          : "assets/images/platinum.png";
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
                          for (var info in patronInfoList) ...[
                            SizedBox(
                              height: height / 50,
                            ),
                            planItem(
                              info,
                              getColor(info.patronPackage),
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
                                  'selectedTier': info.patronTiers[currentTab],
                                  'paymentAssets': paymentAssetsList,
                                };
                                appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page: AuthorizeSubscriptionViewPageConfig);
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

  Widget planItem(PatronInfo info, Color color,
      {required void Function() onTap, required void Function() onReadMore}) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Column(
              children: [
                Container(
                  height: 10,
                  decoration: BoxDecoration(
                      borderRadius: const BorderRadius.only(
                          topLeft: Radius.circular(15.0),
                          topRight: Radius.circular(15.0)),
                      color: color),
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 10.0, vertical: 15.0),
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
                                Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      '${info.patronPackage.capitalizeFirst!} ${"patron".tr()}',
                                      style: TextStyle(
                                        fontSize: 14,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                    Text(
                                      currentTab == 0
                                          ? "amountpermonth".tr(args: [
                                              info.patronTiers[currentTab].price
                                                  .toString(),
                                            ])
                                          : currentTab == 1
                                              ? "amountperyear".tr(args: [
                                                  info.patronTiers[currentTab]
                                                      .price
                                                      .toString()
                                                ])
                                              : "lifetimeplan".tr(),
                                      style: TextStyle(
                                        fontSize: 19,
                                        fontWeight: FontWeight.w400,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold,
                                      ),
                                    ),
                                  ],
                                ),
                                Image.asset(info.logo),
                              ],
                            ),
                            if (info.patronPackage ==
                                    appState.userInfo?.patronMembership
                                        ?.patronPackageId &&
                                info.patronTiers[currentTab].price ==
                                    appState
                                        .userInfo?.patronMembership?.price) ...[
                              SizedBox(height: 10),
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
                            SizedBox(height: 20),
                            Container(
                              width: width / 1.5,
                              child: Text(
                                info.packageListTitle,
                                // '\$${info.patronTiers[0].price} ${"permonth".tr()} / \$${info.patronTiers[1].price} per year / \$${info.patronTiers[2].price} lifetime.',
                                overflow: TextOverflow.visible,
                                style: TextStyle(
                                  fontSize: 15,
                                  color: notifier.getgrey,
                                  fontFamily: fontbody,
                                ),
                              ),
                            ),
                            SizedBox(height: width / 50),
                            for (var item in info.packageList) ...[
                              Row(
                                children: [
                                  Image.asset(
                                    "assets/images/charm_tick.png",
                                    width: 25,
                                  ),
                                  SizedBox(width: 5),
                                  Container(
                                    width: width / 1.5,
                                    child: Text(
                                      item,
                                      overflow: TextOverflow.visible,
                                      style: TextStyle(
                                        fontSize: 15,
                                        fontWeight: FontWeight.w400,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody,
                                      ),
                                    ),
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
                SizedBox(height: 5),
                if (info.patronTiers[currentTab].price !=
                    appState.userInfo?.patronMembership?.price) ...[
                  ButtonOutlined(
                    getActionVerb(info),
                    notifier.getaddsubwalletgrey,
                    notifier.getbluecolor,
                    borderColor: notifier.getbluecolor,
                    onTap: onTap,
                    width: width / 1.4,
                    height: height / 20,
                  ),
                ],
                SizedBox(height: 5),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    TextButton(
                      onPressed: onReadMore,
                      style: TextButton.styleFrom(
                          padding: EdgeInsets.zero,
                          minimumSize: Size(50, 30),
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                          alignment: Alignment.centerLeft),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          Text(
                            "seemore".tr(),
                            style: TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                            ),
                          ),
                          SizedBox(
                            width: 5,
                          ),
                          Icon(
                            Icons.arrow_forward_ios,
                            size: 12,
                          )
                        ],
                      ),
                    ),
                  ],
                ),
                SizedBox(height: height / 50),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget MyTab(String name, int number, bool isActive) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Container(
            child: ElevatedButton(
              onPressed: () => setState(() => currentTab = number),
              style: ButtonStyle(
                backgroundColor: MaterialStateProperty.all<Color>(
                  notifier.isDark
                      ? darktilewhitecolor
                      : isActive
                          ? notifier.getaddsubwalletgrey
                          : wihitecolor,
                ),
                side: MaterialStateProperty.all(
                  BorderSide(
                      color: isActive
                          ? notifier.getbluewhitecolor
                          : notifier.getaddsubwalletgrey,
                      width: 1.5,
                      style: BorderStyle.solid),
                ),
                shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                  const RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(
                      Radius.circular(15),
                    ),
                  ),
                ),
              ),
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 10.0, horizontal: 10),
                child: Text(
                  name,
                  overflow: TextOverflow.visible,
                  style: TextStyle(
                    fontSize: 15,
                    color: isActive
                        ? notifier.getbluewhitecolor
                        : notifier.getgrey,
                    fontFamily: isActive ? fontsemibold : fontbody,
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  String getActionVerb(PatronInfo info) {
    if (info.id > patronMembership!.id!) {
      return '${"downgradeto".tr()}${info.patronPackage.capitalizeFirst!}';
    } else if (info.id < patronMembership!.id!) {
      return '${"upgradeto".tr()}${info.patronPackage.capitalizeFirst!}';
    } else
      return '${"subscribeto".tr()} ${info.patronPackage.capitalizeFirst!}';
  }

  Color getColor(String patronPlan) {
    if (patronPlan.toLowerCase() == 'gold')
      return notifier.getgoldcolor;
    else if (patronPlan.toLowerCase() == 'platinum')
      return notifier.getplatinumcolor;
    else
      return notifier.getdiamondcolor;
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
