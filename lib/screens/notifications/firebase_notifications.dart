// import 'package:firebase_messaging/firebase_messaging.dart';
// import 'package:flutter/material.dart';
// import 'package:flutter/src/foundation/key.dart';
// import 'package:flutter/src/widgets/framework.dart';
// import 'package:flutter_local_notifications/flutter_local_notifications.dart';

// // create an instance
// FirebaseMessaging messaging = FirebaseMessaging.instance;
// late FlutterLocalNotificationsPlugin fltNotification;

// void initMessaging() {
//   var androiInit =
//       AndroidInitializationSettings('@mipmap/ic_launcher'); //for logo
//   var iosInit = IOSInitializationSettings();
//   var initSetting = InitializationSettings(android: androiInit, iOS: iosInit);
//   fltNotification = FlutterLocalNotificationsPlugin();
//   fltNotification.initialize(initSetting);
//   var androidDetails = AndroidNotificationDetails('1', 'fcm_default_channel',
//       channelDescription: 'channel Description');
//   var iosDetails = IOSNotificationDetails();
//   var generalNotificationDetails =
//       NotificationDetails(android: androidDetails, iOS: iosDetails);
//   FirebaseMessaging.onMessage.listen((RemoteMessage message) {
//     print('notification recieved!');
//     RemoteNotification notification = message.notification!;
//     // AndroidNotification android = message.notification!.android!;

//     if (notification.android != null) {
//       fltNotification.show(notification.hashCode, notification.title,
//           notification.body, generalNotificationDetails);
//     }

//     // if (notification != null && android != null) {
//     //   fltNotification.show(notification.hashCode, notification.title,
//     //       notification.body, generalNotificationDetails);
//     // }
//   });

//   FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);
// }

// Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
//   print("Handling a background message");
// }

import 'dart:math';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/screens/Auth/signup.dart';

final FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin =
    FlutterLocalNotificationsPlugin();

late BuildContext _context;

String? selectedNotificationPayload;

Future<void> initAppNotification(context) async {
  _context = context;
  print('................App Notification initialized...................');
  // needed if you intend to initialize in the `main` function
  WidgetsFlutterBinding.ensureInitialized();

  final NotificationAppLaunchDetails? notificationAppLaunchDetails =
      await flutterLocalNotificationsPlugin.getNotificationAppLaunchDetails();
  //String initialRoute = HomePage.routeName;
  if (notificationAppLaunchDetails!.didNotificationLaunchApp) {
    selectedNotificationPayload = notificationAppLaunchDetails.payload!;
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
  //     onSelectNotification: (String payload) async {
  //   if (payload != null) {
  //     debugPrint('notification payload: $payload');
  //     await updateAllCache();
  //   }
  //   selectedNotificationPayload = payload;
  //   // selectNotificationSubject.add(payload);
  // });
  initMyNotification(context);
}

initMyNotification(BuildContext context) {
  FirebaseMessaging.instance.setForegroundNotificationPresentationOptions(
      alert: true, badge: true, sound: true);

  FirebaseMessaging.onMessage.listen((RemoteMessage message) {
    print('Got a message whilst in the foreground!');
    print('Message data: ${message.notification!.body}');
    print('Message data: ${message.notification!.title}');

    String notificationTitle = message.notification!.title!;
    String notificationBody = message.notification!.body!;
    showNotification(notificationTitle, notificationBody);

    if (message.notification != null) {
      print('Message also contained a notification: ${message.notification}');
    }
  });
}

void selectNotification(String? payload) async {
  if (payload != null) {
    print('notification payload: $payload');
  } else {
    print("Notification Done");
  }
  Get.to(() => SignUp(), arguments: payload);
}

Future<void> showNotification(title, message) async {
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
  await flutterLocalNotificationsPlugin.show(
      rand, title, message, platformChannelSpecifics,
      payload: '$title&$message');
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
