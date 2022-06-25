import 'dart:async';
import 'dart:math';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/screens/Backup/congratulation.dart';
import 'package:trovo_wallet/storage/state.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../Models/User.dart';
import '../../router/PageActions.dart';
import '../../router/ui_pages.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({Key? key}) : super(key: key);

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen>
    with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late AnimationController controller;
  PageAction landingPage =
      PageAction(state: PageState.replaceAll, page: LoginPageConfig);

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

    controller = AnimationController(
      vsync: this,
      duration: Duration(milliseconds: 11000),
    );

    controller.addListener(() {
      setState(() {});
    });

    controller.repeat();
    Timer(const Duration(seconds: 4), () {
      // appState.currentAction = PageAction(
      //     state: PageState.replaceAll, page: CongratulationsPageConfig);
      appState.currentAction = landingPage;
      appState.setSplashFinished();
    });
  }

  runAsync() async {
    await getVal();
  }

  getVal() async {
    bool isFirstTime;

    try {
      isFirstTime = await StoreData().storeGetData('isFirstTime') ?? true;
      // isFirstTime = await StoreData().storeGetData('isFirstTime') ?? true;

      print('first time here: $isFirstTime');

      if (isFirstTime) {
        print('first time here indeed: $isFirstTime');
        landingPage =
            PageAction(state: PageState.replaceAll, page: OnboardingPageConfig);
      } else {
        var data = await StoreData().storeGetData('userInfo');
        appState.setUser = UserInfo().deserializeJson(data);
        appState.setSecretKeys = await StoreData().storeGetData('secretKey');
        appState.setPassword = await StoreData().storeGetData('password');
        appState.biometricEnabled =
            await StoreData().storeGetData('biometricsEnabled') ?? false;
        landingPage =
            PageAction(state: PageState.replaceAll, page: LoginPageConfig);
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
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: Center(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            RotationTransition(
              turns: Tween(
                begin: 0.0,
                end: 2 * pi,
              ).animate(controller),
              child:
                  Image.asset("assets/images/trovo.png", height: height / 13),
            ),
            SizedBox(height: height / 45),
            Text(
              "Trovo Wallet",
              style: TextStyle(
                  color: notifier.getdarkgrey,
                  fontFamily: 'Matahari_Semi_Bold',
                  fontSize: 35.sp),
            ),
            // ElevatedButton(
            //   onPressed: () => {
            //     if (controller.isCompleted) {controller.reset()},
            //     controller.forward(),
            //   },
            //   child: Text('again'),
            // ),
            // Stack(
            //   children: [
            //     Column(
            //       children: [
            //         Center(
            //             child: Image.asset("assets/images/trovo.png",
            //                 height: height / 13)),
            //         SizedBox(height: height / 45),
            //         Text(
            //           "Trovo Wallet",
            //           style: TextStyle(
            //               color: notifier.getdarkgrey,
            //               fontFamily: 'Matahari_Semi_Bold',
            //               fontSize: 35.sp),
            //         ),
            //       ],
            //     ),
            //     Center(
            //         child: AnimatedBuilder(
            //             animation: controller, builder: _blurAnimationBuilder))
            //   ],
            // ),
          ],
        )),
      ),
    );
  }

  // Widget _blurAnimationBuilder(context, child) {
  //   double startValue = 10.0;
  //   return BackdropFilter(
  //     filter: ImageFilter.blur(
  //       sigmaX: startValue - controller.value * 10,
  //       sigmaY: startValue - controller.value * 10,
  //     ),
  //     child: Container(color: Colors.transparent),
  //   );
  // }

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
