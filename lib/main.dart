import 'dart:async';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_storage/get_storage.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/back_dispatcher.dart';
import 'package:trovo_wallet/router/route_parser.dart';
import 'package:trovo_wallet/router/router_delegate.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/notifications/firebase_dynamic_links.dart';
import 'package:trovo_wallet/screens/notifications/firebase_notifications.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'Custom_BlocObserver/notifire_clor.dart';
import 'firebase_options.dart';
import 'storage/store.dart';

void main() async {
  await GetStorage.init();
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
  FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);
  await StoreData().storeDeleteItem('initialDynamicLink');
  FlutterError.onError = FirebaseCrashlytics.instance.recordFlutterFatalError;
  var dynamicLink = await FirebaseDynamicLinkInitializer().getInitialLink();
  if (dynamicLink != null) {
    await StoreData()
        .storeInsertData('initialDynamicLink', dynamicLink.link.toString());
  }

  BlocOverrides.runZoned(
    // () => runApp(
    //   DevicePreview(
    //     enabled: true,
    //     tools: [
    //       ...DevicePreview.defaultTools,
    //     ],
    //     builder: (context) => const App(),
    //   ),
    // ),
    () => runApp(const App()),
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
    return MultiProvider(
      providers: [
        ChangeNotifierProvider<ColorNotifier>(create: (_) => ColorNotifier()),
        ChangeNotifierProvider<DataProvider>(create: (_) => appState)
      ],
      child: GestureDetector(
        onTap: _initializeTimer,
        onPanDown: (_) => _initializeTimer(),
        onScaleStart: (_) => _initializeTimer(),
        child: MaterialApp.router(
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
    } else {
      setState(() {
        appState.currentAction =
            PageAction(state: PageState.replaceAll, page: LoginPageConfig);
        appState.isLoggedIn = false;
      });
    }
  }
}
