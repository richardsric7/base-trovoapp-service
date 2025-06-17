import 'dart:async';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:trovo_app/storage/store.dart';

class FCM {
  final FirebaseMessaging _firebaseMessaging = FirebaseMessaging.instance;
  final streamCtlr = StreamController<String>.broadcast();

  // static Future<dynamic> onBackgroundMessage(
  //     Map<String, dynamic> message) async {
  //   print("onMessage: $message");

  //   if (message.containsKey('data')) {
  //     // Handle data message
  //     final dynamic data = message['data'];
  //   }

  //   if (message.containsKey('notification')) {
  //     // Handle notification message
  //     final dynamic notification = message['notification'];
  //   }

  //   // Or do other work.
  // }

  //   setNotifications() {
  //     _firebaseMessaging.configure(
  //       onMessage: (message) async {
  //         print("onMessage: $message");
  //         //streamCtlr.sink.add(message['data']['msg']);
  //         print("onMessage: " + message['data'].toString());
  //         String notificationTitle = message['aps'] == null
  //             ? message['notification']['body'] == null ||
  //                     message['notification']['body'] == ''
  //                 ? 'BantuPay Notification'
  //                 : message['notification']['title'] == '' ||
  //                         message['notification']['title'] == null
  //                     ? 'BantuPay Notification'
  //                     : message['notification']['title']
  //             : message['aps']['alert']['title'];
  //         String notificationBody = message['aps'] == null
  //             ? message['notification']['body'] == null ||
  //                     message['notification']['body'] == ''
  //                 ? message['notification']['title']
  //                 : message['notification']['body']
  //             : message['aps']['alert']['body'];
  //         await showNotification(notificationTitle, notificationBody);
  //         updateAllCache();
  //       },
  //       onBackgroundMessage: Platform.isIOS ? null : onBackgroundMessage,
  //       onLaunch: (message) async {
  //         print("onLaunch: $message");
  //         updateAllCache();
  //       },
  //       onResume: (message) async {
  //         print("onResume: $message");
  //         updateAllCache();
  //       },
  //     );

  //     final token =
  //         _firebaseMessaging.getToken().then((value) => saveToken(value));
  //     _firebaseMessaging.getToken().then((value) => print(value));
  //     _firebaseMessaging.requestNotificationPermissions(
  //         const IosNotificationSettings(
  //             sound: true, badge: true, alert: true, provisional: true));
  //     _firebaseMessaging.onIosSettingsRegistered
  //         .listen((IosNotificationSettings settings) {
  //       print("Settings registered: $settings");
  //     });
  //   }

  //   dispose() {
  //     streamCtlr?.close();
  //   }
  // }

  Future<String> getPushNotificationToken() async {
    String? walletMode = await StoreData().storeGetData('walletMode');
    String token = await StoreData().storeGetData('${walletMode}-token') ?? '';
    if (token.isEmpty) {
      try {
        var result = await _firebaseMessaging.getToken();
        var createdAt = DateTime.now();
        token = '${result!}|$createdAt';
        await StoreData().storeInsertData('${walletMode}-token', token);
      } catch (e) {
        print(e);
      }
    }
    print('FCM Token ${token}');
    return token;
  }
}
