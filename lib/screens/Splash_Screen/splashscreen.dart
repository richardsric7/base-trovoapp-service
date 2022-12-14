import 'dart:async';
import 'dart:math';
import 'package:firebase_remote_config/firebase_remote_config.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/storage/cache.dart';
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
  String? initialDynamicLink;
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
    appState = Provider.of<DataProvider>(context, listen: false);
    runAsync();
    super.initState();
    getdarkmodepreviousstate();

    controller = AnimationController(
      vsync: this,
      duration: Duration(milliseconds: 11000),
    );

    controller.addListener(() {
      setState(() {});
    });

    controller.repeat();
    Timer(const Duration(seconds: 4), () {
      if (appState.splashFinished) {
        if (initialDynamicLink != null) {
          appState.processDeepLink(context, Uri.parse(initialDynamicLink!));
        } else {
          appState.currentAction = landingPage;
        }
      }
    });
  }

  // one way to await an async functions inside
  // initState is to place and await the async function from inside
  // another function which will not be awaited in initState
  runAsync() async {
    await getVal();
    await initFirebaseTools();
  }

  getVal() async {
    try {
      appState.isFirstTime =
          await StoreData().storeGetData('isFirstTime') ?? true;
      initialDynamicLink = await StoreData().storeGetData('initialDynamicLink');
      appState.timeout = await StoreData().storeGetData('timeOut') ?? '5';
      appState.setDefaultCurrency =
          await StoreData().storeGetData('defaultCurrency') ?? 'USD';
      appState.sethideWalletList = List.filled(6, appState.hideBalances);

      if (!appState.appIsOpen) appState.initFirebaseListener(context);

      print(
          'first time here: ${await StoreData().storeGetData('isFirstTime')}');

      if (appState.isFirstTime) {
        print('first time here indeed: ${appState.isFirstTime}');
        landingPage =
            PageAction(state: PageState.replaceAll, page: OnboardingPageConfig);
        appState.setSplashFinished();
      } else {
        var data = await StoreData().storeGetData('userInfo');
        appState.setUser = UserInfo().deserializeJson(data);
        appState.setSecretKeys = await StoreData().storeGetData('secretKey');
        appState.setPassword = await StoreData().storeGetData('password');
        appState.biometricEnabled =
            await StoreData().storeGetData('biometricsEnabled') ?? false;
        appState.hideBalances =
            await StoreData().storeGetData('hideBalances') ?? false;
        appState.assetBalances =
            await StoreData().storeGetData('assetBalances');
        appState.setDefaultAssets =
            await StoreData().storeGetData('defaultAssets');
        appState.setNFTs = await StoreData().storeGetData('nfts');
        appState.setFiatRate = await StoreData().storeGetData('fiatRate');
        appState.setSharedWallets =
            await StoreData().storeGetData('walletsSharedWithUser');
        print('=====================shared wallet ${appState.sharedWallets}');
        appState.sethideWalletList =
            await StoreData().storeGetData('hideWalletList') ??
                List.filled(6, appState.hideBalances);
        var primaryWallet = appState.userInfo!.wallets!.firstWhere(
            (wallet) => wallet.primaryWallet == 1,
            orElse: () => appState.userInfo!.wallets![0]);
        updateUserInfo(primaryWallet.signer, appState.secretKeys[0],
            primaryWallet.publicKey, appState.userInfo!.username!, appState);
        getFiatRates(primaryWallet.signer, appState.secretKeys[0],
            primaryWallet.publicKey, appState.userInfo!.username!, appState);
        appState.activeWallet = primaryWallet;
        // check if app was not already open
        // if app was not already open then move to the next view
        // else wait for the dynamiclink handler to take over
        print(
            '----------------------------------------appIsOpen = $initialDynamicLink');
        // if (initialDynamicLink == null) {
        //   landingPage =
        //       PageAction(state: PageState.replaceAll, page: LoginPageConfig);
        //   appState.setSplashFinished();
        // } else {
        //   appState.processDeepLink(context, Uri.parse(initialDynamicLink!));
        // }
        appState.setSplashFinished();
        appState.appIsOpen = true;
      }
    } catch (e) {
      print('[getVal]getVal exception:' + e.toString());
    }
  }

  initFirebaseTools() async {
    try {
      // initialize firebase remote config
      final remoteConfig = FirebaseRemoteConfig.instance;
      await remoteConfig.setConfigSettings(RemoteConfigSettings(
        fetchTimeout: const Duration(minutes: 1),
        minimumFetchInterval: const Duration(minutes: 1),
      ));

      await await FirebaseRemoteConfig.instance.fetchAndActivate();
    } catch (e) {
      print('firebase error: $e');
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
          ],
        )),
      ),
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
