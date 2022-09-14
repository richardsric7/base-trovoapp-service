import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/Models/Transaction.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallets.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/Auth/AuthorizeActionView.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import '../Models/User.dart';
import '../Models/WalletsListViewData.dart';
import '../router/PageActions.dart';
import 'cache.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  List<String> secretKeys = [];
  bool isDark = false;
  bool biometricEnabled = false;
  String timeout = '5'; // 5 minutes
  String? password;
  var assetBalances;
  var nfts;
  bool dialogOpen = false;
  WalletsListViewData walletView = WalletsListViewData(
      view: WalletView.listWallets,
      actionIcon: Icons.add_circle_outline_sharp,
      actionText: LanguageEn.addsubwallet);

  bool hideBalances = false;
  set sethideBalances(bool value) {
    hideBalances = value;
    print('notifying listeners...');
    notifyListeners();
  }

  void updateListeners() => notifyListeners();

  // bool hideActiveWalletBalance = false;
  // set toggleActiveBalances(bool value) {
  //   hideActiveWalletBalance = value;
  //   notifyListeners();
  // }

  // void resetActiveWalletBalances() {
  //   hideActiveWalletBalance = hideBalances;
  // }

  // used to check if dynamic link was used while the app is open
  // for some reason the splashscreen finishes before the firebase dynamiclink
  // handler is processed so we will use this flag to check on the splashscreen
  // whether the app is open so the splashscreen will wait till the dynamiclink
  // handler finishes and then move to the appropriate next screen.
  bool appIsOpen = false;
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

  // used to keep track of the current bottom navigation index
  // this variable is currently used in back_dispatcher to know when
  // to handle the back button
  int currentBottomTabIndex = 0;

  List<TransactionInfo> historyData = <TransactionInfo>[];
  int limit = 20;
  int currentPage = 1;
  int? totalRecords = 0;

  getHistory() async {
    await fetchHistory(limit);
    notifyListeners();
  }

  Future<void> fetchHistory(limit) async {
    try {
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

        print('transactions: $transactions');

        historyData = transactions;
        notifyListeners();
      }
    } catch (e) {
      print('................................in transaction history: $e');
    }
  }

  // view data is where all the data that a particular view needs
  // to do its work is. So when you want to pass any data from one view to
  // another, assign it to viewData and then get it back when you get
  // to the view. viewData is of type Map<String, dynamic>? where the string key
  // is the ViewPageConfig.key and the value is the data you want to pass to the
  // view. The value is of dynamic type so you can pass any data type you want.
  Map<String, dynamic>? viewData = {};

  // in order to make it possible for bottom navigation tabs to be changed from
  // anywhere in the app we will bring this function here where everybody can
  // reach it from anywhere.
  PageController? bottomTabPageController;

  // use this to keep track of individual wallets' hidden state used
  // especially on the dashboard screen to track which wallet is set to hidden
  // by the user
  List<bool> hideWalletList = [];
  set sethideWalletList(list) {
    hideWalletList.clear();
    if (list != null) {
      for (var i = 0; i < list.length; i++) {
        hideWalletList.add(list[i]);
      }
    }
    notifyListeners();
  }

  void initFirebaseListener() {
    print('initing firebaselistener..............................');
    FirebaseDynamicLinks.instance.onLink.listen((dynamicLinkData) async {
      try {
        print('one 1');
        await StoreData().storeDeleteItem('initialDynamicLink');
        // Navigator.pushNamed(context, dynamicLinkData.link.path);
        print('this is dynamicLinkData: $dynamicLinkData');
        print(
            'current action is login ${dynamicLinkData.link.queryParameters['action']}');
        await StoreData().storeInsertData(
            'initialDynamicLink', dynamicLinkData.link.toString());
        var deepLinkView = getDeepLinkView(dynamicLinkData.link);

        currentAction = deepLinkView;
      } catch (e) {
        print('error processing dynamic link');
      }
    }).onError((error) {
      // Handle errors
      print('this is dynamicLink error: $error');
    });
  }

  PageAction getDeepLinkView(Uri initialDynamicLink) {
    PageAction pageAction;
    if (initialDynamicLink.queryParameters['action'] == 'login')
      pageAction = PageAction(
          state: PageState.addAll,
          pages: [LoginPageConfig, AuthorizeLoginViewPageConfig]);
    else
      pageAction = PageAction(
          state: PageState.addAll,
          pages: [LoginPageConfig, AuthorizeActionViewPageConfig]);

    viewData![pageAction.pages![1].key] = {
      'action': initialDynamicLink.queryParameters['action'],
      'loginId': initialDynamicLink.queryParameters['loginId'],
      'authId': initialDynamicLink.queryParameters['authId'],
      'description': initialDynamicLink.queryParameters['description'],
      'deviceInfo': initialDynamicLink.queryParameters['deviceInfo'],
      'targetUser': initialDynamicLink.queryParameters['targetUser'],
      'ownerUsername': initialDynamicLink.queryParameters['ownerUsername'],
      'serviceShortName': initialDynamicLink.queryParameters['serviceShortName']
    };
    return pageAction;
  }
}
