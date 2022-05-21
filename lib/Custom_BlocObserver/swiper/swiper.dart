import 'package:flutter/material.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/page_view/onbonding_two.dart';
import 'package:gocrypto/screens/page_view/onbondingthree.dart';
import 'package:gocrypto/screens/page_view/one_onbonding.dart';
import 'package:gocrypto/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class Swiper extends StatefulWidget {
  const Swiper({Key? key}) : super(key: key);

  @override
  State<Swiper> createState() => _SwiperState();
}

class _SwiperState extends State<Swiper> {
  late ColorNotifier notifier;

  final PageController _pageController = PageController(initialPage: 0);
  int currentPage = 0;

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

  // List<Widget> _buildPageIndicator() {
  //   List<Widget> list = [];
  //   for (int i = 0; i < _numPages; i++) {
  //     list.add(i == _currentPage ? _indicator(true) : _indicator(false));
  //   }
  //   return list;
  // }
  //
  // Widget _indicator(bool isActive) {
  //   return AnimatedContainer(
  //     duration: const Duration(microseconds: 150),
  //     margin: const EdgeInsets.symmetric(horizontal: 3.0),
  //     height: 8.0,
  //     width: isActive ? 8.0 : 8.0,
  //     decoration: BoxDecoration(
  //       color: isActive ? Colors.red : Colors.grey,
  //       borderRadius: const BorderRadius.all(Radius.circular(12)),
  //     ),
  //   );
  // }
  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return Scaffold(
      body: Column(
        children: [
          Container(
            color: Colors.transparent,
            height: height,
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
                Threeonbonding(),
                Onbondingtwo(),
              ],
            ),
            // ),Row(
            //   mainAxisAlignment: MainAxisAlignment.center,
            //   children: _buildPageIndicator(),
            // ),
            // Container(
            //   height: 50,
            //   child: PageView(
            //     physics: const ClampingScrollPhysics(),
            //     controller: _pageController,
            //     onPageChanged: (int page) {
            //       setState(() {
            //         _currentPage = page;
            //       });
            //     },
            //     scrollDirection: Axis.horizontal,
            //     children:   <Widget>[
            //
            //       oneonbonding(),
            //       onbondingtwo(),
            //     ],
            //
            //   ),
          ),
        ],
      ),
    );
  }
}
