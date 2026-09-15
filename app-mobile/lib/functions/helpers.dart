import 'dart:io';
import 'package:android_id/android_id.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:flutter/services.dart';
import 'package:package_info_plus/package_info_plus.dart';

Future getDeviceDetails() async {
  String? identifier;
  const _androidIdPlugin = AndroidId();
  final DeviceInfoPlugin deviceInfoPlugin = DeviceInfoPlugin();
  try {
    if (Platform.isAndroid) {
      identifier = await _androidIdPlugin.getId(); //UUID for Android
    } else if (Platform.isIOS) {
      var data = await deviceInfoPlugin.iosInfo;
      identifier = data.identifierForVendor; //UUID for iOS
    }
  } on PlatformException {}

  return identifier;
}

Future<String> getAppVersion() async {
  PackageInfo packageInfo = await PackageInfo.fromPlatform();
  var appVersion = packageInfo.version;
  return appVersion;
}
