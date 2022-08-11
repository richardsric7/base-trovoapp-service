import 'dart:async';

import 'package:firebase_core/firebase_core.dart';
// import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
// import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
// import 'package:fluttertoast/fluttertoast.dart';
import 'package:get_storage/get_storage.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/back_dispatcher.dart';
import 'package:trovo_wallet/router/route_parser.dart';
import 'package:trovo_wallet/router/router_delegate.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/notifications/firebase_notifications.dart';
import 'package:trovo_wallet/storage/state.dart';
// import 'package:uni_links/uni_links.dart';
import 'Custom_BlocObserver/notifire_clor.dart';

void main() async {
  await GetStorage.init();
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp();
  BlocOverrides.runZoned(
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
  // bool _initialURILinkHandled = false;
  // Uri? _initialURI;
  // Uri? _currentURI;
  // Object? _err;

  StreamSubscription? _streamSubscription;
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
    initAppNotification(context);
    // _initURIHandler();
    // _incomingLinkHandler();
    // DynamicLinkService().handleDynamicLink();
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider<ColorNotifier>(create: (_) => ColorNotifier()),
        ChangeNotifierProvider<DataProvider>(create: (_) => appState)
      ],
      child: MaterialApp.router(
        routerDelegate: delegate!,
        routeInformationParser: parser,
        backButtonDispatcher: backButtonDispatcher,
        debugShowCheckedModeBanner: false,
      ),
    );
  }

  // Widget build(BuildContext context) {
  //   return MaterialApp(
  //     home: Scaffold(
  //         appBar: AppBar(
  //           title: Text('Deeplink Tutorial'),
  //         ),
  //         body: Center(
  //             child: Padding(
  //           padding: const EdgeInsets.symmetric(horizontal: 20),
  //           child: Column(
  //             mainAxisAlignment: MainAxisAlignment.center,
  //             children: <Widget>[
  //               // 1
  //               ListTile(
  //                 title: const Text("Initial Link"),
  //                 subtitle: Text(_initialURI.toString()),
  //               ),
  //               // 2
  //               if (!kIsWeb) ...[
  //                 // 3
  //                 ListTile(
  //                   title: const Text("Current Link Host"),
  //                   subtitle: Text('${_currentURI?.host}'),
  //                 ),
  //                 // 4
  //                 ListTile(
  //                   title: const Text("Current Link Scheme"),
  //                   subtitle: Text('${_currentURI?.scheme}'),
  //                 ),
  //                 // 5
  //                 ListTile(
  //                   title: const Text("Current Link"),
  //                   subtitle: Text(_currentURI.toString()),
  //                 ),
  //                 // 6
  //                 ListTile(
  //                   title: const Text("Current Link Path"),
  //                   subtitle: Text('${_currentURI?.path}'),
  //                 )
  //               ],
  //               // 7
  //               if (_err != null)
  //                 ListTile(
  //                   title: const Text('Error',
  //                       style: TextStyle(color: Colors.red)),
  //                   subtitle: Text(_err.toString()),
  //                 ),
  //               const SizedBox(
  //                 height: 20,
  //               ),
  //               const Text("Check the blog for testing instructions")
  //             ],
  //           ),
  //         ))),
  //   );
  // }

  // Future<void> _initURIHandler() async {
  //   // 1
  //   if (!_initialURILinkHandled) {
  //     _initialURILinkHandled = true;
  //     // 2
  //     Fluttertoast.showToast(
  //         msg: "Invoked _initURIHandler",
  //         toastLength: Toast.LENGTH_SHORT,
  //         gravity: ToastGravity.BOTTOM,
  //         timeInSecForIosWeb: 1,
  //         backgroundColor: Colors.green,
  //         textColor: Colors.white);
  //     try {
  //       // 3
  //       final initialURI = await getInitialUri();
  //       // 4
  //       if (initialURI != null) {
  //         debugPrint("Initial URI received $initialURI");
  //         if (!mounted) {
  //           return;
  //         }
  //         setState(() {
  //           _initialURI = initialURI;
  //         });
  //       } else {
  //         debugPrint("Null Initial URI received");
  //       }
  //     } on PlatformException {
  //       // 5
  //       debugPrint("Failed to receive initial uri");
  //     } on FormatException catch (err) {
  //       // 6
  //       if (!mounted) {
  //         return;
  //       }
  //       debugPrint('Malformed Initial URI received');
  //       setState(() => _err = err);
  //     }
  //   }
  // }

  // void _incomingLinkHandler() {
  //   // 1
  //   if (!kIsWeb) {
  //     // 2
  //     _streamSubscription = uriLinkStream.listen((Uri? uri) {
  //       if (!mounted) {
  //         return;
  //       }
  //       debugPrint('Received URI: $uri');
  //       setState(() {
  //         _currentURI = uri;
  //         _err = null;
  //       });
  //       // 3
  //     }, onError: (Object err) {
  //       if (!mounted) {
  //         return;
  //       }
  //       debugPrint('Error occurred: $err');
  //       setState(() {
  //         _currentURI = null;
  //         if (err is FormatException) {
  //           _err = err;
  //         } else {
  //           _err = null;
  //         }
  //       });
  //     });
  //   }
  // }

  // @override
  // void dispose() {
  //   _streamSubscription?.cancel();
  //   super.dispose();
  // }
}
