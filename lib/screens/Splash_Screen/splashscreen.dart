import 'dart:async';
import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/swiper/swiper.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/screens/Auth/login.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SpashScreen extends StatefulWidget {
  const SpashScreen({Key? key}) : super(key: key);

  @override
  State<SpashScreen> createState() => _SpashScreenState();
}

class _SpashScreenState extends State<SpashScreen>
    with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  late AnimationController controller;
  Widget landingPage = Login();

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
    runAsync();
    // FCM firebaseMessaging = FCM();
    // firebaseMessaging.setNotifications();
    // firebaseMessaging.streamCtlr.stream.listen((msgData) {
    //   print('The Firebase data notification message is: $msgData');
    //   updateAllCache();
    // });

    controller = AnimationController(
      vsync: this,
      duration: Duration(milliseconds: 1000),
    );

    controller.addListener(() {
      setState(() {});
    });

    controller.forward();
    Timer(
      const Duration(seconds: 4),
      // () => Navigator.push(
      () => Navigator.pushReplacement(
        context,
        LandingPageRoute(landingPage),
      ),
    );
  }

  runAsync() async {
    await getVal();
  }

  getVal() async {
    String isFirstTime;
    String activeSecret;

    try {
      isFirstTime = await StoreData().storeGetData('password') ?? '';
      activeSecret = await StoreData().storeGetData('secretKey') ?? '';

      print('first time here: ' + isFirstTime);
      print('first time here: ' + activeSecret);

      if (isFirstTime.isEmpty || activeSecret.isEmpty) {
        setState(() {
          print('first time here indeed: ' + isFirstTime);
          print('first time here deedin: ' + activeSecret);
          landingPage = Swiper();
        });
      } else {
        setState(() {
          landingPage = Login();
        });
      }
    } catch (e) {
      print('[getVal]getVal exception:' + e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: Center(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Stack(
              children: [
                // RotationTransition(
                //     turns: Tween(
                //       begin: 0.0,
                //       end: 2 * pi,
                //     ).animate(controller),
                //     child: Image.asset("assets/images/trovo.png",
                //         height: height / 13)),
                Column(
                  children: [
                    Center(
                        child: Image.asset("assets/images/trovo.png",
                            height: height / 13)),
                    SizedBox(height: height / 45),
                    Text(
                      "Trovo Wallet",
                      style: TextStyle(
                          color: notifier.getdarkgrey,
                          fontFamily: 'Matahari_Semi_Bold',
                          fontSize: 35.sp),
                    ),
                  ],
                ),
                Center(
                    child: AnimatedBuilder(
                        animation: controller, builder: _blurAnimationBuilder))
              ],
            ),
          ],
        )),
      ),
    );
  }

  Widget _blurAnimationBuilder(context, child) {
    double startValue = 10.0;
    return BackdropFilter(
      filter: ImageFilter.blur(
        sigmaX: startValue - controller.value * 10,
        sigmaY: startValue - controller.value * 10,
      ),
      child: Container(color: Colors.transparent),
    );
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }
}

class LandingPageRoute extends MaterialPageRoute {
  Widget child;
  LandingPageRoute(this.child)
      : super(builder: (BuildContext context) => child);

  // OPTIONAL IF YOU WISH TO HAVE SOME EXTRA ANIMATION WHILE ROUTING
  @override
  Widget buildPage(BuildContext context, Animation<double> animation,
      Animation<double> secondaryAnimation) {
    return FadeTransition(
      opacity: animation,
      child: child,
    );
  }

  @override
  Duration get transitionDuration => Duration(milliseconds: 1000);
}
