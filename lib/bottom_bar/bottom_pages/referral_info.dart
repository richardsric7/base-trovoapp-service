import 'package:firebase_remote_config/firebase_remote_config.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ReferralInfo extends StatefulWidget {
  const ReferralInfo({Key? key}) : super(key: key);

  @override
  State<ReferralInfo> createState() => _ReferralInfoState();
}

class _ReferralInfoState extends State<ReferralInfo>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late TabController _tabController;

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
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          LanguageEn.myreferrals,
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 20,
              ),
              Center(
                child: Image.asset(
                  "assets/images/obi.png",
                  height: height / 10,
                  fit: BoxFit.fill,
                ),
              ),
              SizedBox(height: height / 90),
              Text(
                'Obi Enechi',
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 16.sp),
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  share();
                },
                child: invitefriend(notifier.getbluecolor,
                    LanguageEn.invitefriends, wihitecolor),
              ),
              SizedBox(height: height / 50),
              DefaultTabController(
                length: 2,
                child: Container(
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
                        child: TabBar(
                          controller: _tabController,
                          labelColor: notifier.getbluewhitecolor,
                          indicatorColor: notifier.getbluewhitecolor,
                          labelStyle: TextStyle(
                            fontSize: 13.sp,
                            fontWeight: FontWeight.w600,
                            fontFamily: fontsemibold,
                          ),
                          tabs: [
                            Tab(
                              height: 50,
                              text: LanguageEn.referrals,
                            ),
                            Tab(
                              height: 50,
                              text: LanguageEn.rewards,
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              Container(
                height: height / 1.9,
                child: TabBarView(controller: _tabController, children: [
                  Column(
                    children: [
                      SizedBox(height: height / 30),
                      Text(
                        "4 ${LanguageEn.referrals}",
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 16.sp),
                      ),
                      referralList(),
                    ],
                  ),
                  Column(
                    children: [
                      SizedBox(height: height / 30),
                      Text(
                        "25 TROV earned",
                        style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 16.sp),
                      ),
                      commissionList(),
                    ],
                  ),
                ]),
              ),
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }

  Widget referralList() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  referredUser('Thundeyy', 'joined 1 day ago'),
                  // Instagram
                  SizedBox(
                    height: height / 50,
                  ),
                  referredUser('muche', '3 days ago'),
                  // Instagram
                  SizedBox(
                    height: height / 50,
                  ),
                  referredUser('ric', '3 days ago'),
                  SizedBox(
                    height: height / 50,
                  ),
                  referredUser('kennis', 'joined 1 week ago'),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget referredUser(name, timeAgo) {
    return Container(
      width: width / 1.29,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Text(
            name,
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
          ),
          Text(
            timeAgo,
            style: TextStyle(
                fontSize: 11,
                fontWeight: FontWeight.w600,
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold),
          ),
        ],
      ),
    );
  }

  Widget commissionList() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  commissionItem('4 TROV', 'Swap commission', '3 days ago'),
                  // Instagram
                  SizedBox(
                    height: height / 40,
                  ),
                  commissionItem('12 TROV', 'Swap commission', '3 days ago'),
                  SizedBox(
                    height: height / 40,
                  ),
                  commissionItem(
                      '9 TROV', 'Subscription commission', 'joined 1 week ago'),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget commissionItem(amount, rel, timeAgo) {
    return Container(
      width: width / 1.29,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Container(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  amount,
                  style: TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w600,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold),
                ),
                SizedBox(
                  height: 5,
                ),
                GestureDetector(
                  onTap: () {},
                  child: Text(
                    rel,
                    style: TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                )
              ],
            ),
          ),
          Text(
            timeAgo,
            style: TextStyle(
                fontSize: 11,
                fontWeight: FontWeight.w600,
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold),
          ),
        ],
      ),
    );
  }

  Future<void> share() async {
    var label = await FirebaseRemoteConfig.instance
        .getString('wallet_referral_share_label');
    await FlutterShare.share(
      title: 'Trovo Wallet',
      text: label,
    );
  }

  Widget invitefriend(colorbutton, buttontext, buttontextcolor) {
    return Center(
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(15),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            LayoutBuilder(builder: (context, constraints) {
              return Container(
                height: height / 10,
                width: width / 1.1,
                decoration: BoxDecoration(
                  color: colorbutton!,
                  borderRadius: BorderRadius.circular(15),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    Image.asset("assets/images/referrals.png",
                        height: height / 30),
                    Container(
                      width: width / 1.7,
                      child: Text(
                        buttontext!,
                        textAlign: TextAlign.start,
                        style: TextStyle(
                            fontFamily: fontbody,
                            fontSize: 13.sp,
                            color: buttontextcolor),
                      ),
                    ),
                    Icon(
                      Icons.arrow_forward_ios,
                      size: 12.sp,
                      color: wihitecolor,
                    )
                  ],
                ),
              );
            }),
          ],
        ),
      ),
    );
  }
}
