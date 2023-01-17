import 'dart:math';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:trovo_wallet/models/bottom_tab_page.dart';
import 'package:trovo_wallet/firebase_options.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

final FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin =
    FlutterLocalNotificationsPlugin();

late BuildContext _context;
late DataProvider _appState;

String? selectedNotificationPayload;

Future<void> initAppNotification(context, appState) async {
  _context = context;
  _appState = appState;
  print('................App Notification initialized...................');
  // needed if you intend to initialize in the `main` function
  WidgetsFlutterBinding.ensureInitialized();

  final NotificationAppLaunchDetails? notificationAppLaunchDetails =
      await flutterLocalNotificationsPlugin.getNotificationAppLaunchDetails();
  //String initialRoute = HomePage.routeName;
  if (notificationAppLaunchDetails!.didNotificationLaunchApp) {
    // selectedNotificationPayload = notificationAppLaunchDetails.payload!;
    // initialRoute = SecondPage.routeName;
  }

  const AndroidInitializationSettings initializationSettingsAndroid =
      AndroidInitializationSettings('@mipmap/ic_launcher');

  /// Note: permissions aren't requested here just to demonstrate that can be
  /// done later
  final IOSInitializationSettings initializationSettingsIOS =
      IOSInitializationSettings(
          requestAlertPermission: false,
          requestBadgePermission: false,
          requestSoundPermission: false,
          onDidReceiveLocalNotification: onDidReceiveLocalNotification);
  const MacOSInitializationSettings initializationSettingsMacOS =
      MacOSInitializationSettings(
          requestAlertPermission: false,
          requestBadgePermission: false,
          requestSoundPermission: false);
  final InitializationSettings initializationSettings = InitializationSettings(
      android: initializationSettingsAndroid,
      iOS: initializationSettingsIOS,
      macOS: initializationSettingsMacOS);
  await flutterLocalNotificationsPlugin.initialize(initializationSettings,
      onSelectNotification: selectNotification);
  initMyNotification(context);
}

initMyNotification(BuildContext context) {
  FirebaseMessaging.instance.setForegroundNotificationPresentationOptions(
      alert: true, badge: true, sound: true);

  FirebaseMessaging.onMessage.listen((RemoteMessage message) {
    print('Got a message whilst in the foreground! ${message.toMap()}');

    showNotification(message);

    if (message.notification != null) {
      print('Message also contained a notification: ${message.notification}');
    }
  });

  //When the app is in the background, but not terminated.
  FirebaseMessaging.onMessageOpenedApp.listen(
    (message) {
      print('==============message opened app:${message.notification!.title}');
      print(message.toMap());
      goToPageRoute(message.data['route'] ?? '');
      return;
    },
    cancelOnError: false,
    onDone: () {},
  );

  FirebaseMessaging.onBackgroundMessage((message) {
    print('Got a message whilst in the background!');
    print('Message data: ${message.notification!.body}');
    print('Message data: ${message.notification!.title}');
    // lets just return something that makes the compiler
    // happy.
    return Future.sync(() {});
  });
}

void selectNotification(String? route) async {
  if (route != null) {
    print('================notification payload: $route');
    goToPageRoute(route);
    return;
  }
}

Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  print("Handling a background message");
  await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
}

Future<void> showNotification(RemoteMessage payload) async {
  const AndroidNotificationDetails androidPlatformChannelSpecifics =
      AndroidNotificationDetails('Trovo Wallet', 'Trovo Technologies',
          channelDescription:
              'Trovo Wallet is an app to manage all your crypto assets',
          importance: Importance.max,
          priority: Priority.high,
          ticker: 'ticker');
  const NotificationDetails platformChannelSpecifics =
      NotificationDetails(android: androidPlatformChannelSpecifics);
  var rand = Random().nextInt(999999);
  await flutterLocalNotificationsPlugin.show(rand, payload.notification!.title!,
      payload.notification!.body!, platformChannelSpecifics,
      payload: '${payload.data['route']}');
}

void onDidReceiveLocalNotification(
    int id, String? title, String? body, String? payload) async {
  // display a dialog with the notification details, tap ok to go to another page
  showDialog(
    context: _context,
    builder: (BuildContext context) => CupertinoAlertDialog(
      title: Text(title!),
      content: Text(body!),
      actions: [
        CupertinoDialogAction(
          isDefaultAction: true,
          child: Text('Ok'),
          onPressed: () async {
            Navigator.of(context, rootNavigator: true).pop();
          },
        )
      ],
    ),
  );
}

void requestPermissions() {
  flutterLocalNotificationsPlugin
      .resolvePlatformSpecificImplementation<
          IOSFlutterLocalNotificationsPlugin>()
      ?.requestPermissions(
        alert: true,
        badge: true,
        sound: true,
      );
  flutterLocalNotificationsPlugin
      .resolvePlatformSpecificImplementation<
          MacOSFlutterLocalNotificationsPlugin>()
      ?.requestPermissions(
        alert: true,
        badge: true,
        sound: true,
      );
}

void goToPageRoute(String route) {
  if (_appState.isLoggedIn) {
    switch (route) {
      case 'basicTransactionHistory':
        _appState.currentAction =
            PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
        changeTabPage(_appState, ButtomTabPage.TransactionHistory.index);
        break;
      case 'pendingApproval':
        _appState.currentAction = PageAction(
          state: PageState.addAll,
          pages: [BottomHomePageConfig, SharedAccessViewPageConfig],
        );
        // take the user to the pending approvals tab on the shared access view
        WidgetsBinding.instance.addPostFrameCallback((_) {
          _appState.sharedAccesstabController.animateTo(1,
              duration: Duration(milliseconds: 500), curve: Curves.easeInOut);
        });
        break;
      default:
        _appState.currentAction =
            PageAction(state: PageState.replaceAll, page: BottomHomePageConfig);
        changeTabPage(_appState, ButtomTabPage.Dashboard.index);
    }
  } else {
    _appState.currentAction =
        PageAction(state: PageState.replaceAll, page: LoginPageConfig);
  }
}
