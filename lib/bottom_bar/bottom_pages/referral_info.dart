import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../storage/state.dart';
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
    var appState = Provider.of<DataProvider>(context, listen: true);
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
                child: Image.asset("assets/images/avatar.png",
                    height: height / 10),
              ),
              Text(
                'Obi Enechi',
                style: TextStyle(
                    color: notifier.getbluecolor,
                    fontFamily: 'Gilroy_Bold',
                    fontSize: 16.sp),
              ),
              SizedBox(height: height / 50),
              GestureDetector(
                onTap: () {
                  share();
                },
                child: invitefriend(notifier.getbluecolor,
                    LanguageEn.invitefriends, notifier.getwihitecolor),
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
                          labelColor: notifier.getbluecolor,
                          labelStyle: TextStyle(
                            fontSize: 13.sp,
                            fontWeight: FontWeight.w600,
                            fontFamily: fontbody,
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
                child: Expanded(
                  child: TabBarView(controller: _tabController, children: [
                    socials(),
                    socials(),
                  ]),
                ),
              ),
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }

  Widget socials() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
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
                  // Twitter
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Image.asset(
                                "assets/images/twitter.png",
                                height: height / 30,
                                color: notifier.getbluecolor,
                              ),
                            ],
                          ),
                        ),
                        Container(
                          child: Column(
                            children: [
                              Text(
                                LanguageEn.unverified,
                                style: TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.w600,
                                    color: notifier.getbluecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: 5,
                              ),
                              GestureDetector(
                                onTap: () {},
                                child: Text(
                                  LanguageEn.taptoverify,
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.w400,
                                    color: notifier.getbluecolor,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              )
                            ],
                          ),
                        )
                      ],
                    ),
                  ),
                  // Instagram
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Image.asset(
                            "assets/images/instagram.png",
                            height: height / 30,
                            color: notifier.getbluecolor,
                          ),
                        ),
                        Container(
                          child: Column(
                            children: [
                              Image.asset(
                                "assets/images/tick.png",
                                height: height / 30,
                                color: notifier.getbluecolor,
                              ),
                            ],
                          ),
                        )
                      ],
                    ),
                  ),
                  // Instagram
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Image.asset(
                            "assets/images/gemcave.png",
                            height: height / 30,
                            color: notifier.getbluecolor,
                          ),
                        ),
                        GestureDetector(
                          onTap: () {},
                          child: Text(
                            LanguageEn.taptoconnect,
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        )
                      ],
                    ),
                  ),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> share() async {
    await FlutterShare.share(
        title: 'Example share',
        text: 'Example share text',
        linkUrl: 'https://flutter.dev/',
        chooserTitle: 'Example Chooser Title');
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
                      color: notifier.getwihitecolor,
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
