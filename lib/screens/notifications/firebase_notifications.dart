import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/material.dart';
import 'package:flutter/src/foundation/key.dart';
import 'package:flutter/src/widgets/framework.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';

// create an instance
FirebaseMessaging messaging = FirebaseMessaging.instance;
late FlutterLocalNotificationsPlugin fltNotification;

void initMessaging() {
  var androiInit =
      AndroidInitializationSettings('@mipmap/ic_launcher'); //for logo
  var iosInit = IOSInitializationSettings();
  var initSetting = InitializationSettings(android: androiInit, iOS: iosInit);
  fltNotification = FlutterLocalNotificationsPlugin();
  fltNotification.initialize(initSetting);
  var androidDetails = AndroidNotificationDetails('1', 'fcm_default_channel',
      channelDescription: 'channel Description');
  var iosDetails = IOSNotificationDetails();
  var generalNotificationDetails =
      NotificationDetails(android: androidDetails, iOS: iosDetails);
  FirebaseMessaging.onMessage.listen((RemoteMessage message) {
    RemoteNotification notification = message.notification!;
    // AndroidNotification android = message.notification!.android!;

    if (notification.android != null) {
      fltNotification.show(notification.hashCode, notification.title,
          notification.body, generalNotificationDetails);
    }

    // if (notification != null && android != null) {
    //   fltNotification.show(notification.hashCode, notification.title,
    //       notification.body, generalNotificationDetails);
    // }
  });

  FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);
}

Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  print("Handling a background message");
}
