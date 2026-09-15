import 'dart:async';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:trovo_app/storage/store.dart';

class FCM {
  final FirebaseMessaging _firebaseMessaging = FirebaseMessaging.instance;
  final streamCtlr = StreamController<String>.broadcast();

  Future<String> getPushNotificationToken() async {
    String? walletMode = await StoreData().storeGetData('walletMode');
    String token = await StoreData().storeGetData('${walletMode}-token') ?? '';
    if (token.isEmpty) {
      try {
        var result = await _firebaseMessaging.getToken();
        var createdAt = DateTime.now();
        token = '${result!}|$createdAt';
        await StoreData().storeInsertData('${walletMode}-token', token);
      } catch (e) {}
    }
    return token;
  }
}
