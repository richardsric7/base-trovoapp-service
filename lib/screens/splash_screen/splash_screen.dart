import 'dart:async';
import 'dart:convert';
import 'dart:math';
import 'package:easy_localization/easy_localization.dart';
import 'package:firebase_remote_config/firebase_remote_config.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/services/push_fcm_service.dart';
import 'package:trovo_app/storage/cache.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/user.dart';
import '../../router/page_actions.dart';
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
  bool timerIsDone = false;
  String? initialDynamicLink;

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
      timerIsDone = true;
    });
  }

  // one way to await an async functions inside
  // initState is to place and await the async function from inside
  // another function which will not be awaited in initState
  runAsync() async {
    await initializeAppData();
    await initFirebaseTools();
  }

  initializeAppData() async {
    try {
      fetchVersionInfo(appState);
      if (await StoreData().storeGetData('walletMode') == null) {
        await StoreData().storeInsertData('walletMode', "Testnet");
      }
      appState.walletMode = await StoreData().storeGetData('walletMode');
      appState.restartedAfterSwitch =
          await StoreData().storeGetData('restartedAfterSwitch') ?? false;
      appState.isFirstTime =
          await StoreData().storeGetData('isFirstTime') ?? true;
      initialDynamicLink = await StoreData().storeGetData('initialDynamicLink');
      appState.timeout = await StoreData().storeGetData('timeOut') ?? '5';
      appState.setDefaultCurrency =
          await StoreData().storeGetData('defaultCurrency') ?? 'USD';
      appState.setDefaultLanguage =
          await StoreData().storeGetData('defaultLanguage') ?? 'en';
      appState.sethideWalletList = List.filled(6, appState.hideBalances);

      if (!appState.appIsOpen) appState.initFirebaseListener(context);

      if (appState.isFirstTime) {
        Timer.periodic(Duration(milliseconds: 200), (timer) {
          if (timerIsDone) {
            timer.cancel();
            appState.setSplashFinished();
            appState.currentAction = PageAction(
                state: PageState.replaceAll, page: OnboardingPageConfig);
          }
        });
      } else {
        PackageInfo packageInfo = await PackageInfo.fromPlatform();
        appState.appVersion = packageInfo.version;
        var data = await StoreData().storeGetData('userInfo');
        var sharedWallets =
            await StoreData().storeGetData('walletsSharedWithUser');
        var assetBalances = await StoreData().storeGetData('assetBalances');
        appState.setUser = UserInfo().deserializeJson(
          data,
          sharedWallets,
          assetBalances,
        );
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
        appState.setAssetOrderings =
            await StoreData().storeGetData('assetOrderings');
        appState.setNFTs = await StoreData().storeGetData('nfts');
        appState.setFiatRate = await StoreData().storeGetData('fiatRate') ?? {};
        print('fiatRates ${appState.fiatRate['NGN']}');
        appState.introducedSharedAccess =
            await StoreData().storeGetData('introducedSharedAccess') ?? false;
        appState.sethideWalletList =
            await StoreData().storeGetData('hideWalletList') ??
                List.filled(6, appState.hideBalances);
        var primaryWallet = appState.userInfo!.wallets!.firstWhere(
            (wallet) => wallet.primaryWallet == 1,
            orElse: () => appState.userInfo!.wallets![0]);
        appState.activeWallet = primaryWallet;
        // check if app was not already open
        // if app was not already open then move to the next view
        // else wait for the dynamiclink handler to take over
        Timer.periodic(Duration(milliseconds: 200), (timer) async {
          if (timerIsDone) {
            timer.cancel();
            appState.setSplashFinished();
            appState.appIsOpen = true;

            if (appState.restartedAfterSwitch) {
              print('importing after switch......');
              await importWalletAfterSwitch(appState, context);
            } else {
              String result = await FCM().getPushNotificationToken();

              var token = result.split('|').first;
              DateTime createdAt = DateTime.parse(result.split('|').last);
              var dateDifference = DateTime.now().difference(createdAt);

              updateUserInfo(
                primaryWallet.signer,
                appState.secretKeys[0],
                primaryWallet.publicKey,
                appState.userInfo!.username!,
                appState,
                pnt: dateDifference.inDays > 10 ? token : null,
              );
              fetchNotifications(appState);

              if (initialDynamicLink != null) {
                appState.processDeepLink(
                    context, Uri.parse(initialDynamicLink!));
              } else {
                appState.currentAction = PageAction(
                    state: PageState.replaceAll, page: LoginPageConfig);
              }
            }
          }
        });
      }
    } catch (e) {
      print('initializeAppData exception:' + e.toString());
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
              "Trovo",
              style: TextStyle(
                  color: notifier.getdarkgrey,
                  fontFamily: 'Matahari_Semi_Bold',
                  fontSize: 35.sp),
            ),
            Text(
              "App",
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

  Future<void> importWalletAfterSwitch(
    DataProvider appState,
    BuildContext context,
  ) async {
    var username = appState.userInfo!.username;
    var signer = appState.primaryWallet.signer!;
    var publicKey = appState.primaryWallet.publicKey!;
    var secretKey = appState.secretKeys[0];

    String result = await FCM().getPushNotificationToken();
    var token = result.split('|').first;

    Map responseData = await makeGetRequest(
        uri: '/v1/users/${username}?type=import&pnt=$token',
        signer: signer,
        publicKey: publicKey,
        secretKey: secretKey);
    // print('response==================> $responseData');

    if (responseData['statusCode'] == 200) {
      appState.tempPublicKey = publicKey;
      appState.tempSecretKey = secretKey;
      appState.tempSigner = signer;
      appState.tempPassword = appState.password!;

      fetchNotifications(appState);
      getFiatRates(appState);
      storeUserInfo(responseData['data'], appState);
      await StoreData()
          .storeInsertData('biometricsEnabled', appState.biometricEnabled);

      appState.currentAction =
          PageAction(state: PageState.addPage, page: LoginPageConfig);
    } else if (responseData['statusCode'] == 404) {
      accountNotFoundAfterSwitchPopup(
        context,
        onContinueWithCredentials: () async =>
            await createUserAccountAfterSwitch(),
        onImportNewCredential: () => {
          appState.currentAction =
              PageAction(state: PageState.addPage, page: ImportWalletPageConfig)
        },
        onGoBackToPrevEnvironment: () {
          appState.changeWalletMode(
            appState.walletMode == 'Testnet' ? 'Mainnet' : 'Testnet',
            isReversed: true,
          );
        },
      );
    } else {
      // must be some sort of server error
      // let's throw it
      accountNotFoundAfterSwitchPopup(
        context,
        message: responseData['data']['message'],
        onContinueWithCredentials: () {},
        onImportNewCredential: () => {
          appState.currentAction =
              PageAction(state: PageState.addPage, page: ImportWalletPageConfig)
        },
        onGoBackToPrevEnvironment: () {
          appState.changeWalletMode(
            appState.walletMode == 'Testnet' ? 'Mainnet' : 'Testnet',
            isReversed: true,
          );
        },
      );
    }
  }

  Future<void> createUserAccountAfterSwitch() async {
    try {
      showLoader(context);
      appState.userInfo!.pushNotificationToken =
          await FCM().getPushNotificationToken();
      Map map = {
        'username': appState.userInfo!.username,
        'email': appState.userInfo!.email,
        'firstName': appState.userInfo!.firstName,
        'lastName': appState.userInfo!.lastName,
        'mobile': appState.userInfo!.mobile,
        'mobileCountryCode': appState.userInfo!.countryCode,
        'referrer': appState.userInfo!.referrer,
        'pushNotificationToken': appState.userInfo!.pushNotificationToken,
        'corporate': appState.userInfo!.corporate,
        'verificationCode': '',
      };

      Account? creds = parseKey(context, appState.secretKeys[0])!;

      print('creating user account after switch... ${map}');

      String jsonBody = jsonEncode(map);

      Map responseData = await makePostRequest(
          uri: '/v1/users',
          body: jsonBody,
          signer: creds.publicKey,
          publicKey: creds.publicKey,
          secretKey: creds.secretKey);

      // print('$responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        appState.tempPublicKey = creds.publicKey;
        appState.tempSecretKey = creds.secretKey;
        appState.tempSigner = creds.publicKey;
        appState.tempPassword = appState.password!;

        appState.currentAction =
            PageAction(state: PageState.addPage, page: VerificationPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(context,
          title: "error".tr(),
          message: e.toString().contains('firebase')
              ? 'Network error! Please check your connection and try again.'
              : e.toString());
    }
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
