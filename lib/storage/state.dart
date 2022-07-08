import 'package:flutter/cupertino.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../Models/User.dart';
import '../router/PageActions.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  List<String> secretKeys = [];
  bool isDark = false;
  bool biometricEnabled = false;
  String? password;
  var assetBalances;
  var nfts;

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

  Wallet? activeWallet;
  set setActiveWallet(value) {
    activeWallet = value;
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

  set setassetBalances(balances) {
    assetBalances = balances;
    notifyListeners();
  }

  set setNFTs(newNfts) {
    nfts = newNfts;
    notifyListeners();
  }

  String secretToBackup = '';
  set setSecretToBackup(secret) {
    secretToBackup = secret;
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
