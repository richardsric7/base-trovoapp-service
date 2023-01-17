import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/models/transaction.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/models/wallets_list_view_data.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallets.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/models/user.dart';
import '../router/page_actions.dart';
import 'cache.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  bool isLoggedIn = false;
  List<String> secretKeys = [];
  bool isDark = false;
  bool biometricEnabled = false;
  bool isFirstTime = true;
  String timeout = '5'; // 5 minutes
  String? password;
  String phoneVersion = '';
  var assetBalances;
  var nfts;
  Wallet get primaryWallet =>
      userInfo!.wallets!.firstWhere((wallet) => wallet.primaryWallet == 1);
  Map _transactionableWallets = {};
  Map get transactionableWallets {
    var wallets = userInfo!.wallets;

    for (var i = 0; i < wallets!.length; i++) {
      // get just the standard wallets since they are the only ones we can
      // enable shared access on
      if (wallets[i].walletType == 0) {
        // do not add user wallets where user doesn't have initiator access
        if (wallets[i].walletThreshold == 2 &&
            wallets[i]
                .permissions!
                .where((perm) =>
                    perm.permission == 'INITIATOR' &&
                    perm.targetUsername == userInfo!.username)
                .isEmpty) {
          continue;
        }

        _transactionableWallets[wallets[i].publicKey!] = {
          'publicKey': wallets[i].publicKey,
          'alias': wallets[i].alias,
          'threshold': wallets[i].walletThreshold,
          'sharedAccessEnabled': wallets[i].primaryWallet == 1
              ? 0
              : wallets[i].sharedAccessEnabled,
          'claimedAssets': assetBalances[wallets[i].publicKey!]['claimed'],
        };
      }
    }

    // then get all the shared wallets where I have initiator access on
    for (var i = 0; i < sharedWallets.length; i++) {
      if (sharedWallets[i]['permission'] == 'INITIATOR') {
        _transactionableWallets[sharedWallets[i]['walletPublicKey']] = {
          'publicKey': sharedWallets[i]['walletPublicKey'],
          'alias': '${sharedWallets[i]['walletAlias']}',
          'permission': sharedWallets[i]['permission'],
          'threshold': sharedWallets[i]['walletSettings']['walletThreshold'],
          'sharedAccessEnabled': 1,
          'claimedAssets': sharedWallets[i]['assetBalances']['claimed'],
        };
      }
    }
    return _transactionableWallets;
  }

  Map _allWallets = {}; // both shared and non-shared
  Map get allWallets {
    var wallets = userInfo!.wallets;

    for (var i = 0; i < wallets!.length; i++) {
      if (wallets[i].walletThreshold == 2 &&
          wallets[i]
              .permissions!
              .where((perm) =>
                  perm.permission == 'INITIATOR' &&
                  perm.targetUsername == userInfo!.username)
              .isEmpty) {
        continue;
      }
      _allWallets[wallets[i].publicKey!] = {
        'walletType': wallets[i].walletType,
        'publicKey': wallets[i].publicKey,
        'alias': wallets[i].alias,
        'threshold': wallets[i].walletThreshold,
        'sharedAccessEnabled':
            wallets[i].primaryWallet == 1 ? 0 : wallets[i].sharedAccessEnabled,
        'claimedAssets': assetBalances[wallets[i].publicKey!]['claimed'],
      };
    }

    for (var i = 0; i < sharedWallets.length; i++) {
      _allWallets[sharedWallets[i]['walletPublicKey']] = {
        // since the wallet type is unknown give it a number that can't make transactions
        'walletType': sharedWallets[i]['walletSettings']?['walletType'] ?? 2,
        'publicKey': sharedWallets[i]['walletPublicKey'],
        'alias': '${sharedWallets[i]['walletAlias']}',
        'permission': sharedWallets[i]['permission'],
        'threshold':
            sharedWallets[i]['walletSettings']?['walletThreshold'] ?? 0,
        'sharedAccessEnabled': 1,
        'claimedAssets': sharedWallets[i]['assetBalances']['claimed'],
      };
    }
    return _allWallets;
  }

  bool dialogOpen = false;
  WalletsListViewData walletView = WalletsListViewData(
      view: WalletView.listWallets,
      actionIcon: Icons.add_circle_outline_sharp,
      actionText: LanguageEn.addsubwallet);

  bool hideBalances = false;
  set sethideBalances(bool value) {
    hideBalances = value;
    notifyListeners();
  }

  void updateListeners() => notifyListeners();

  bool introducedSharedAccess = false;
  set setIntroducedSharedAccess(value) {
    introducedSharedAccess = value;
    notifyListeners();
  }

  var sharedWallets;
  set setSharedWallets(wallets) {
    sharedWallets = wallets;
    notifyListeners();
  }

  bool hasNewAnnouncement = false;
  set setHasNewAnnouncement(bool value) {
    hasNewAnnouncement = value;
    notifyListeners();
  }

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

  String tempUsername = '';
  set setTempUsername(value) {
    tempUsername = value;
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

  String tempReferrerUsername = '';
  set setTempReferrerUsername(value) {
    tempReferrerUsername = value;
    notifyListeners();
  }

  bool tempInvalidateOldSigner = true;
  set setTempInvalidateOldSigner(value) {
    tempInvalidateOldSigner = value;
    notifyListeners();
  }

  String tempEmailOtp = '';
  set setTempEmailOtp(value) {
    tempEmailOtp = value;
    notifyListeners();
  }

  Map tempSecurityQuestionsAndAnswers = {};
  set setTempSecurityQuestionsAndAnswers(value) {
    tempSecurityQuestionsAndAnswers = value;
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

  var defaultAssets = [];
  set setDefaultAssets(assets) {
    defaultAssets = assets;
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

  Map fiatRate = {};
  set setFiatRate(value) {
    fiatRate = value;
    notifyListeners();
  }

  String defaultCurrency = '';
  set setDefaultCurrency(value) {
    defaultCurrency = value;
    notifyListeners();
  }

// region transaction history filter
  DateTime? filterStartDate;
  set setFilterStartDate(value) {
    filterStartDate = value;
    notifyListeners();
  }

  DateTime? filterEndDate;
  set setFilterEndDate(value) {
    filterEndDate = value;
    notifyListeners();
  }

  String? filterMinAmount;
  set setFilterMinAmount(value) {
    filterMinAmount = value;
    notifyListeners();
  }

  String? filterMaxAmount;
  set setFilterMaxAmount(value) {
    filterMaxAmount = value;
    notifyListeners();
  }

  String? filterUsername;
  set setFilterUsername(value) {
    filterUsername = value;
    notifyListeners();
  }

  String? filterMemo;
  set setFilterMemo(value) {
    filterMemo = value;
    notifyListeners();
  }

  String? filterFromPublicKey;
  set setFilterFromPublicKey(value) {
    filterFromPublicKey = value;
    notifyListeners();
  }

  String? filterToPublicKey;
  set setFilterToPublicKey(value) {
    filterToPublicKey = value;
    notifyListeners();
  }

  String filterAsset = "*|*";
  set setFilterAsset(value) {
    filterAsset = value;
    notifyListeners();
  }

  String filterQuery = "";
  set setFilterQuery(value) {
    filterQuery = value;
    notifyListeners();
  }

  // used to keep track of the current bottom navigation index
  // this variable is currently used in back_dispatcher to know when
  // to handle the back button
  int currentBottomTabIndex = 0;

  List<TransactionInfo> historyData = <TransactionInfo>[];
  int limit = 20;
  int currentPage = 1;
  int? totalRecords = 0;

  getHistory(context, String forPublicKey, {void Function()? onDone}) async {
    showLoader(context);
    await fetchHistory(
      context,
      forPublicKey,
      limit: limit.toString(),
      query: filterQuery,
    );

    notifyListeners();
    hideLoader(context);
    // scroll to the top of the list if historyData is not null
    if (historyData.length > 0 && onDone != null) onDone();
  }

  Future<void> fetchHistory(
    context,
    String forPublicKey, {
    String? limit,
    String? query,
  }) async {
    try {
      // viewData![PaymentHistoryViewPageConfig.key] will not be null when the
      // the payment history view is opened from shared wallet. So we use the
      // viewData to get the public key of the shared wallet and fetch its transaction
      // history.
      // var publicKey = viewData![PaymentHistoryViewPageConfig.key] != null
      //     ? viewData![PaymentHistoryViewPageConfig.key]['walletPublicKey']
      //     : activeWallet!.publicKey!;
      var uri = '/v1/users/payments/${forPublicKey}?limit=$limit${query}';
      if (!filterAsset.contains("*")) {
        var splitAssetInfo = filterAsset.split("|");
        uri +=
            "&assetIssuer=${splitAssetInfo[0]}&assetCode=${splitAssetInfo[1].isEmpty ? "XBN" : splitAssetInfo[1]}";
      }
      Map responseData = await makeGetRequest(
          uri: uri,
          signer: activeWallet!.signer!,
          publicKey: forPublicKey,
          secretKey: secretKeys[0]);

      if (responseData['statusCode'] == 200) {
        totalRecords = responseData['data']['totalRecords'];
        currentPage = responseData['data']['currentPage'];
        var transactions = <TransactionInfo>[];
        for (var i = 0; i < responseData['data']['records'].length; i++) {
          transactions.add(TransactionInfo()
              .deserializeJson(responseData['data']['records'][i]));
        }

        historyData = transactions;
        notifyListeners();
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {}
  }
// end region transaction history filter

  Future<void> refreshData() async {
    try {
      await updateUserInfo(userInfo!.wallets![0].signer, secretKeys[0],
          userInfo!.wallets![0].publicKey, userInfo!.username, this);
    } catch (e) {
      // print(e);
    }
  }

  String filterTransactionType = "";
  set setFilterTransactionType(value) {
    filterTransactionType = value;
    notifyListeners();
  }

  String filterTransactionId = "";
  set setFilterTransactionId(value) {
    filterTransactionId = value;
    notifyListeners();
  }

  String filterInitiatorUsername = "";
  set setFilterInitiatorUsername(value) {
    filterInitiatorUsername = value;
    notifyListeners();
  }

  String filterDescription = "";
  set setFilterDescription(value) {
    filterDescription = value;
    notifyListeners();
  }

  String filterWalletAlias = "";
  set setFilterWalletAlias(value) {
    filterWalletAlias = value;
    notifyListeners();
  }

  String filterWalletPublicKey = "";
  set setFilterWalletPublicKey(value) {
    filterWalletPublicKey = value;
    notifyListeners();
  }

  String filterTransactionStatus = "";
  set setFilterTransactionStatus(value) {
    filterTransactionStatus = value;
    notifyListeners();
  }

  // excludes the ones that users have signed even if the
  // transaction is still pending because it has not yet gotten
  // number of required approvals
  int excludeUserApproved = 1;
  set setExcludeUserApproved(value) {
    excludeUserApproved = value;
    notifyListeners();
  }

  late Future<Map> approvals;

  Future<Map> fetchApprovals({String? limit, String? query}) async {
    try {
      var uri =
          '/v1/shared-access/approvals?&excludeUserApproved=${excludeUserApproved}&limit=$limit${query}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: activeWallet!.signer!,
        secretKey: secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  getApprovals({void Function()? onDone}) {
    approvals = fetchApprovals(limit: limit.toString(), query: filterQuery);
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

  late TabController sharedAccesstabController;

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

  void initFirebaseListener(BuildContext context) {
    FirebaseDynamicLinks.instance.onLink.listen((dynamicLinkData) async {
      try {
        await StoreData().storeInsertData(
            'initialDynamicLink', dynamicLinkData.link.toString());
        processDeepLink(context, dynamicLinkData.link);
      } catch (e) {
        popup(context,
            title: 'Error', message: 'error processing dynamic link');
      }
    }).onError((error) {
      // Handle errors
    });
  }

  void processDeepLink(BuildContext context, Uri initialDynamicLink,
      {String? rel}) {
    // action login
    if (initialDynamicLink.queryParameters['action'] == 'login') {
      setSplashFinished();
      viewData![AuthorizeLoginViewPageConfig.key] = {
        'action': initialDynamicLink.queryParameters['action'],
        'loginId': initialDynamicLink.queryParameters['loginId'],
        'description': initialDynamicLink.queryParameters['description'],
        'deviceInfo': initialDynamicLink.queryParameters['deviceInfo'],
        'targetUser': initialDynamicLink.queryParameters['targetUser'],
        'ownerUsername': initialDynamicLink.queryParameters['ownerUsername'],
        'serviceShortName':
            initialDynamicLink.queryParameters['serviceShortName'],
      };
      currentAction = PageAction(
          state: PageState.addAll,
          pages: [LoginPageConfig, AuthorizeLoginViewPageConfig]);
    } else if (initialDynamicLink.queryParameters['action'] == 'payment') {
      // action payment
      if (initialDynamicLink.queryParameters['assetCode'] != '' &&
          initialDynamicLink.queryParameters['assetCode'] != null) {
        showChooseWalletPopup(
            context,
            initialDynamicLink.queryParameters['assetCode'] == 'XBN'
                ? ''
                : initialDynamicLink.queryParameters['assetCode'],
            initialDynamicLink.queryParameters['assetIssuer'],
            onDone: (walletPublicKey, isSharedWallet) {
          var deeplinkInfo = {
            "assetCode": initialDynamicLink.queryParameters['assetCode'],
            "assetIssuer": initialDynamicLink.queryParameters['assetIssuer'],
            "source": "qr2",
            "receiver":
                initialDynamicLink.queryParameters['paymentDestination'],
            "amount": initialDynamicLink
                .queryParameters['amount'], // amount we want to send
            "memo": initialDynamicLink.queryParameters['memo'],
            'action': 'payment',
            'sendingWallet': walletPublicKey,
          };

          var claimedAssets =
              transactionableWallets[walletPublicKey]['claimedAssets'];

          var deeplinkAssetCode = deeplinkInfo['assetCode'] == 'XBN'
              ? ''
              : deeplinkInfo['assetCode'];

          for (var asset in claimedAssets) {
            if (asset['assetCode'] == deeplinkAssetCode &&
                asset['assetIssuer'] == deeplinkInfo['assetIssuer']) {
              viewData![SendAssetViewPageConfig.key] = {
                'assetCode': asset['assetCode'],
                'assetIssuer': asset['assetIssuer'],
                'amount': asset['amount'],
                'imageUrl': asset['imageUrl'],
                'usdPrice': asset['usdPrice'],
                'walletInfo': {'isSharedWallet': isSharedWallet},
              };

              // exit the loop immediately we get what we are looking for
              break;
            }
          }

          viewData![SendAssetViewPageConfig.key]['deepLinkInfo'] = deeplinkInfo;
          currentAction = PageAction(
              state: PageState.addAll,
              pages: isLoggedIn
                  ? [BottomHomePageConfig, SendAssetViewPageConfig]
                  : [LoginPageConfig, SendAssetViewPageConfig]);
          setSplashFinished();
        }, onCancel: () {
          if (rel == 'qrScanner') {
            Navigator.of(context).pop();
          } else {
            currentAction = PageAction(
                state: PageState.addAll,
                pages: isLoggedIn ? [BottomHomePageConfig] : [LoginPageConfig]);
          }
          setSplashFinished();
        });
      }
    } else if (initialDynamicLink.queryParameters['action'] == 'authorize') {
      viewData![AuthorizeActionViewPageConfig.key] = {
        'action': initialDynamicLink.queryParameters['action'],
        'authId': initialDynamicLink.queryParameters['authId'],
        'description': initialDynamicLink.queryParameters['description'],
        'deviceInfo': initialDynamicLink.queryParameters['deviceInfo'],
        'targetUser': initialDynamicLink.queryParameters['targetUser'],
        'ownerUsername': initialDynamicLink.queryParameters['ownerUsername'],
        'serviceShortName':
            initialDynamicLink.queryParameters['serviceShortName']
      };
      // action authorize
      setSplashFinished();
      currentAction = PageAction(
          state: PageState.addAll,
          pages: [LoginPageConfig, AuthorizeActionViewPageConfig]);
    } else if (initialDynamicLink.queryParameters['action'] == 'register') {
      setSplashFinished();
      setTempReferrerUsername = initialDynamicLink.queryParameters['referrer'];

      // action register
      currentAction = PageAction(
          state: PageState.addAll,
          pages: [LoginPageConfig, CreatePasswordPageConfig]);
    }
  }
}
