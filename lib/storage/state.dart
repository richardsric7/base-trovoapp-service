import 'package:flutter/cupertino.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  bool isDark = false;
  bool biometricEnabled = false;

  set setUser(info) {
    print('setting user...');
    userInfo = info;
    notifyListeners();
  }
}
