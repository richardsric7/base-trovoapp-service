import 'package:flutter/cupertino.dart';
import '../Models/User.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  List<String> secretKeys = [];
  bool isDark = false;
  bool biometricEnabled = false;
  String? password;

  set setPassword(pswd) {
    password = pswd;
    notifyListeners();
  }

  set setUser(info) {
    print('setting user...');
    userInfo = info;
    notifyListeners();
  }

  set setSecretKeys(secrets) {
    if (secrets != null) {
      for (var i = 0; i < secrets.length; i++) {
        secretKeys.add(secrets[i]);
        notifyListeners();
      }
    }
  }
}
