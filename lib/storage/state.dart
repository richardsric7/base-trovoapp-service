import 'package:flutter/cupertino.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import '../Models/User.dart';
import '../router/PageActions.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  List<String> secretKeys = [];
  bool isDark = false;
  bool biometricEnabled = false;
  String? password;

  bool _splashFinished = false;
  bool get splashFinished => _splashFinished;
  void setSplashFinished() {
    _splashFinished = true;
    print(_splashFinished);
    notifyListeners();
  }

  PageAction _currentAction = PageAction();
  PageAction get currentAction => _currentAction;
  set currentAction(PageAction action) {
    _currentAction = action;
    notifyListeners();
  }

  void resetCurrentAction() {
    _currentAction = PageAction();
  }

  set setPassword(pswd) {
    password = pswd;
    notifyListeners();
  }

  String tempPassword = '';
  set setTempPassword(value) {
    tempPassword = value;
    notifyListeners();
  }

  String tempPublicKey = '';
  set setTempPublicKey(value) {
    tempPublicKey = value;
    notifyListeners();
  }

  String tempSecretKey = '';
  set setTempSecretKey(value) {
    tempSecretKey = value;
    notifyListeners();
  }

  set setUser(info) {
    userInfo = info;
    notifyListeners();
  }

  set setSecretKeys(secrets) {
    secretKeys.clear();
    if (secrets != null) {
      for (var i = 0; i < secrets.length; i++) {
        secretKeys.add(secrets[i]);
      }
    }
    notifyListeners();
  }

  void addSecrets(List<String> secrets) {
    for (var i = 0; i < secrets.length; i++) {
      secretKeys.add(secrets[i]);
    }
    notifyListeners();
  }

  String initialUrl = "";
  goToWebView(url) {
    initialUrl = url;
    currentAction =
        PageAction(state: PageState.addPage, page: WebViewPageConfig);
  }
}
