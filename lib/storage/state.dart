import 'package:flutter/cupertino.dart';
import 'package:trovo_wallet/Models/Transaction.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../Models/User.dart';
import '../router/PageActions.dart';
import 'cache.dart';

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
    print('going to: $url');
    initialUrl = url;
    currentAction =
        PageAction(state: PageState.addPage, page: WebViewPageConfig);
  }

  Future<void> refreshData() async {
    try {
      await updateUserInfo(userInfo!.wallets![0].signer, secretKeys[0],
          userInfo!.wallets![0].publicKey, userInfo!.username, this);
    } catch (e) {
      print(e);
    }
  }

  List<TransactionInfo> historyData = <TransactionInfo>[];
  int limit = 20;
  int currentPage = 1;
  int? totalRecords = 0;

  getHistory() async {
    // var data =
    //     await StoreData().storeGetData('historyData${activeWallet!.alias}');

    // if (data == null || data.length <= 0) {
    await fetchHistory(limit);
    //   return;
    // }

    // for (var i = 0; i < data.length; i++) {
    //   historyData.clear();
    //   historyData.add(TransactionInfo().deserializeJson(data[i]));
    // }

    // totalRecords =
    //     await StoreData().storeGetData('totalRecords${activeWallet!.alias}');
    // currentPage =
    //     await StoreData().storeGetData('currentPage${activeWallet!.alias}');
    notifyListeners();
  }

  Future<void> fetchHistory(limit) async {
    print('fetching history for: ${activeWallet!.publicKey!}');
    Map responseData = await makeGetRequest(
        uri: '/v1/users/payments/${activeWallet!.publicKey}?limit=$limit',
        signer: activeWallet!.signer!,
        publicKey: activeWallet!.publicKey!,
        secretKey: secretKeys[0]);

    print('response: ${responseData['data']}');
    if (responseData['statusCode'] == 200) {
      totalRecords = responseData['data']['totalRecords'];
      currentPage = responseData['data']['currentPage'];
      var transactions = <TransactionInfo>[];
      for (var i = 0; i < responseData['data']['records'].length; i++) {
        transactions.add(TransactionInfo()
            .deserializeJson(responseData['data']['records'][i]));
      }

      // if (limit <= 20) {
      //   await StoreData().storeInsertData('historyData${activeWallet!.alias}',
      //       responseData['data']['records']);
      //   await StoreData().storeInsertData('totalRecords${activeWallet!.alias}',
      //       responseData['data']['totalRecords']);
      //   await StoreData().storeInsertData('currentPage${activeWallet!.alias}',
      //       responseData['data']['currentPage']);
      // }

      print('transactions: $transactions');

      historyData = transactions;
      notifyListeners();
    }
  }

  // view data is where all the data that a particular view needs
  // to do its work is. So when you want to pass any data from one view to
  // another, assign it to viewData and then get it back when you get
  // to the view. viewData is of type Map<String, dynamic>? where the string key
  // is the ViewPageConfig.key and the value is the data you want to pass to the
  // view. The value is of dynamic type so you can pass any data type you want.
  Map<String, dynamic>? viewData;
}
