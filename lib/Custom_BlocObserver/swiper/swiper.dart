import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/create_password.dart';
import 'package:trovo_wallet/screens/page_view/onbonding_two.dart';
import 'package:trovo_wallet/screens/page_view/onbondingthree.dart';
import 'package:trovo_wallet/screens/page_view/one_onbonding.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../screens/Auth/login.dart';
import '../../utils/enstring.dart';
import '../button/custtom_button.dart';
import '../fonts.dart';

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
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return Scaffold(
      backgroundColor: notifier.getwihitecolor,
      body: Column(
        children: [
          Container(
            color: Colors.transparent,
            height: height / 1.5,
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
          SizedBox(height: height / 20.5),
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: _buildPageIndicator(),
              ),
              SizedBox(height: height / 30.5),
              Button(
                LanguageEn.getstarted,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const CreatePassword(),
                    ),
                  );
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                LanguageEn.importwallet,
                notifier.getwihitecolor,
                notifier.getbluecolor,
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => const Login(),
                    ),
                  );
                },
              ),
            ],
          ),
        ],
      ),
    );
  }
}
