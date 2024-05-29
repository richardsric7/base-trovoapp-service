import 'dart:async';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get_storage/get_storage.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/back_dispatcher.dart';
import 'package:trovo_wallet/router/route_parser.dart';
import 'package:trovo_wallet/router/router_delegate.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/asset-tokenization/state.dart';
import 'package:trovo_wallet/screens/notifications/firebase_dynamic_links.dart';
import 'package:trovo_wallet/screens/notifications/firebase_notifications.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'custom_bloc_observer/notifire_clor.dart';
import 'firebase_options.dart';
import 'storage/store.dart';
import 'package:easy_localization/easy_localization.dart';

void main() async {
  await GetStorage.init();
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(
      name: await StoreData().storeGetData('walletMode') ?? "Mainnet",
      options: DefaultFirebaseOptions.currentPlatform(
          await StoreData().storeGetData('walletMode') ?? "Mainnet"));
  print(
      '-----------------------------------------------this is the initialized app from main method:  ${DefaultFirebaseOptions.currentPlatform(await StoreData().storeGetData('walletMode') ?? "Testnet")}');
  FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);
  await StoreData().storeDeleteItem('initialDynamicLink');
  FlutterError.onError = FirebaseCrashlytics.instance.recordFlutterFatalError;
  var dynamicLink = await FirebaseDynamicLinkInitializer().getInitialLink();
  if (dynamicLink != null) {
    await StoreData()
        .storeInsertData('initialDynamicLink', dynamicLink.link.toString());
  }

  await EasyLocalization.ensureInitialized();

  runApp(
    EasyLocalization(
        supportedLocales: [Locale('en', 'US')],
        path:
            'assets/translations', // <-- change the path of the translation files
        fallbackLocale: Locale('en-US'),
        child: App()),
  );
}

class App extends StatefulWidget {
  const App({Key? key}) : super(key: key);

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> {
  TrovoWalletBackButtonDispatcher? backButtonDispatcher;
  final appState = DataProvider();
  Timer? _timer;
  late FirebaseMessaging messaging;
  TrovoWalletRouterDelegate? delegate;
  final parser = TrovoWalletRouteParser();

  _AppState() {
    delegate = TrovoWalletRouterDelegate(appState);
    delegate?.setNewRoutePath(SplashPageConfig);
    backButtonDispatcher = TrovoWalletBackButtonDispatcher(delegate!);
  }

  @override
  void initState() {
    super.initState();
    initAppNotification(context, appState);
  }

  @override
  Widget build(BuildContext context) {
    SystemChrome.setPreferredOrientations([
      DeviceOrientation.portraitUp,
    ]);
    return MultiProvider(
      providers: [
        ChangeNotifierProvider<ColorNotifier>(create: (_) => ColorNotifier()),
        ChangeNotifierProvider<DataProvider>(create: (_) => appState),
        ChangeNotifierProvider<AssetTokenizationViewsState>(
            create: (_) => AssetTokenizationViewsState())
      ],
      child: GestureDetector(
        onTap: _initializeTimer,
        onPanDown: (_) => _initializeTimer(),
        onScaleStart: (_) => _initializeTimer(),
        child: MaterialApp.router(
          localizationsDelegates: context.localizationDelegates,
          supportedLocales: context.supportedLocales,
          locale: context.locale,
          routerDelegate: delegate!,
          routeInformationParser: parser,
          backButtonDispatcher: backButtonDispatcher,
          debugShowCheckedModeBanner: false,
          theme: ThemeData(
            colorScheme: ThemeData().colorScheme.copyWith(primary: trovoblue),
            fontFamily: fontbody,
          ),
        ),
      ),
    );
  }

  void _initializeTimer() async {
    print('------------------timer initialized--------------------');
    String? time = await StoreData().storeGetData('timeOut');
    print('this is timeout: $time');
    int? timeOut = int.tryParse(time == null ? '5' : time);

    // print('Timer init And Timeout is $timeOut........... ');

    if (_timer != null) {
      _timer!.cancel();
    }
    // setup action after 5 minutes
    _timer = Timer(Duration(minutes: timeOut == null ? 5 : timeOut),
        () => _handleInactivity());
  }

  void _handleInactivity() async {
    _timer?.cancel();
    _timer = null;

    if (appState.isFirstTime) {
      setState(() {
        appState.currentAction =
            PageAction(state: PageState.replaceAll, page: OnboardingPageConfig);
      });
    } else if (appState.restartedAfterSwitch) {
      // do nothing
    } else {
      setState(() {
        appState.currentAction =
            PageAction(state: PageState.replaceAll, page: LoginPageConfig);
        appState.isLoggedIn = false;
      });
    }
  }

  // void registerNotification() async {
  //   // 1. Initialize the Firebase app
  //   // await Firebase.initializeApp();
  //   var app = await Firebase.initializeApp(
  //       name: 'Mainnet',
  //       options: DefaultFirebaseOptions.currentPlatform("Mainnet"));
  //   print(
  //       '-----------------------------------------------this is the initialized app from registerNotification:  ${app.name}');

  //   // 2. Instantiate Firebase Messaging
  //   messaging = FirebaseMessaging.instance;
  //   // 3. On iOS, this helps to take the user permissions
  //   NotificationSettings settings = await messaging.requestPermission(
  //     alert: true,
  //     badge: true,
  //     provisional: false,
  //     sound: true,
  //   );

  //   if (settings.authorizationStatus == AuthorizationStatus.authorized) {
  //     print('User granted permission');
  //     // TODO: handle the received notifications
  //   } else {
  //     print('User declined or has not accepted permission');
  //   }
  // }
}
