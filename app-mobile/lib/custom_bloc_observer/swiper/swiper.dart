import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/screens/page_view/onboarding_two.dart';
import 'package:trovo_app/screens/page_view/onboarding_three.dart';
import 'package:trovo_app/screens/page_view/onboarding_one.dart';
import 'package:trovo_app/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/loader.dart';

import '../../router/page_actions.dart';
import '../../storage/state.dart';
import '../button/custtom_button.dart';

class Swiper extends StatefulWidget {
  const Swiper({Key? key}) : super(key: key);

  @override
  State<Swiper> createState() => _SwiperState();
}

class _SwiperState extends State<Swiper> {
  late ColorNotifier notifier;

  final PageController _pageController = PageController(initialPage: 0);
  int currentPage = 0;
  int numPages = 3;

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
  }

  List<Widget> _buildPageIndicator() {
    List<Widget> list = [];
    for (int i = 0; i < numPages; i++) {
      list.add(i == currentPage ? _indicator(true) : _indicator(false));
    }
    return list;
  }

  Widget _indicator(bool isActive) {
    return AnimatedContainer(
      duration: const Duration(microseconds: 150),
      margin: const EdgeInsets.symmetric(horizontal: 3.0),
      height: 8.0,
      width: isActive ? 8.0 : 8.0,
      decoration: BoxDecoration(
        color: isActive ? notifier.getbluecolor : Colors.grey,
        borderRadius: const BorderRadius.all(Radius.circular(12)),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: false);
    var appState = Provider.of<DataProvider>(context, listen: false);
    hideLoader(context);
    return Scaffold(
      backgroundColor: notifier.getwihitecolor,
      body: Column(
        children: [
          Container(
            color: Colors.transparent,
            height: height / 1.3,
            child: PageView(
              physics: const ClampingScrollPhysics(),
              controller: _pageController,
              onPageChanged: (int page) {
                setState(() {
                  currentPage = page;
                });
              },
              scrollDirection: Axis.horizontal,
              children: const <Widget>[
                Oneonbonding(),
                Onbondingtwo(),
                Threeonbonding(),
              ],
            ),
          ),
          // SizedBox(height: height / 20.5),
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: _buildPageIndicator(),
              ),
              SizedBox(height: height / 20.5),
              ButtonOutlined(
                currentPage == 2 ? "proceed".tr() : "skip".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: GetStartedViewPageConfig,
                  );
                },
              ),
              SizedBox(height: height / 50),
            ],
          ),
        ],
      ),
    );
  }
}
