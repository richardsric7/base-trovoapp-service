import 'dart:async';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:terminate_restart/terminate_restart.dart';
import 'package:trovo_app/models/curated_asset.dart';
import 'package:trovo_app/models/deposit_transaction_model.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/models/transaction.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/models/wallets_list_view_data.dart';
import 'package:trovo_app/models/withdrawal_transaction_model.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/models/user.dart';
import '../router/page_actions.dart';
import 'cache.dart';

class DataProvider with ChangeNotifier {
  UserInfo? userInfo;
  bool isLoggedIn = false;
  List<String> secretKeys = [];
  bool isDark = false;
  bool biometricEnabled = false;
  bool restartedAfterSwitch = false;
  List<CuratedAsset> curatedSwapList = [];
  Map<String, CuratedAsset> curatedSwapListMap = {};
  bool isFirstTime = true;
  String timeout = '5'; // 5 minutes
  String? password;
  String appVersion = '';
  String linkId = '';
  // keep track of the view you'd like to return a user to after certain operations
  PageAction? returnView;
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
            wallets[i].permissions!
                .where(
                  (perm) =>
                      perm.permission == 'INITIATOR' &&
                      perm.targetUsername == userInfo!.username,
                )
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

    if (sharedWallets != null) {
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
    }
    return _transactionableWallets;
  }

  List<CuratedAsset> deserializeSwapList(List<dynamic> m) {
    List<CuratedAsset> list = [];
    m.forEach((item) {
      list.add(CuratedAsset().deserializeJson(item));
    });
    return list;
  }

  Map<String, Map<String, int>> assetOrderings = {};
  set setAssetOrderings(immutableMap) {
    if (immutableMap != null) {
      immutableMap.forEach(
        (key, valueMap) => {
          valueMap.forEach((key2, value2) {
            if (assetOrderings[key] == null) {
              assetOrderings[key] = {key2: int.parse(value2.toString())};
            }
            assetOrderings[key]![key2] = int.parse(value2.toString());
          }),
        },
      );
    }
    notifyListeners();
  }

  bool dialogOpen = false;
  WalletsListViewData walletView = WalletsListViewData(
    view: WalletView.listWallets,
    actionIcon: Icons.add_circle_outline_sharp,
    actionText: "addsubwallet".tr(),
  );

  String walletMode = 'Testnet';

  Future<void> changeWalletMode(String value, {bool isReversed = false}) async {
    try {
      StoreData().storeInsertData('walletMode', value);
      walletMode = value;
      currentAction = PageAction(
        state: PageState.addPage,
        page: SplashPageConfig,
      );
      Timer(const Duration(seconds: 4), () async {
        StoreData().storeInsertData('restartedAfterSwitch', !isReversed);
        await TerminateRestart.instance.restartApp(
          options: const TerminateRestartOptions(terminate: true),
        );
      });
      notifyListeners();
    } catch (e) {}
  }

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

  void setPage({
    PageState state = PageState.addPage,
    required PageConfiguration? page,
    Widget? widget,
  }) {
    if (page == null && widget == null) {
      throw Exception('Please supply a page or a widget!');
    }
    _currentAction = PageAction(state: state, page: page, widget: widget);
    notifyListeners();
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

  List<String> backupSecrets = [];

  TokenizedAsset? tokenizedAsset = null;
  List<TokenizedAsset> tempTokenizedAssetList = [];

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

  String? tempSigner = '';
  set setTempSigner(value) {
    tempSigner = value;
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

  Map tokenizationData = {};
  set setTokenizationData(data) {
    tokenizationData = data;
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
    print('initialUrl: =====> $initialUrl');
    currentAction = PageAction(
      state: PageState.addPage,
      page: WebViewPageConfig,
    );
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

  String defaultLanguage = '';
  set setDefaultLanguage(value) {
    defaultLanguage = value;
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
  int totalRecords = 0;

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
        secretKey: secretKeys[0],
      );

      if (responseData['statusCode'] == 200) {
        totalRecords = responseData['data']['totalRecords'];
        currentPage = responseData['data']['currentPage'];
        var transactions = <TransactionInfo>[];
        for (var i = 0; i < responseData['data']['records'].length; i++) {
          transactions.add(
            TransactionInfo().deserializeJson(
              responseData['data']['records'][i],
            ),
          );
        }

        historyData = transactions;
        notifyListeners();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {}
  }
  // end region transaction history filter

  Future<void> refreshData() async {
    try {
      await Future.wait([
        updateUserInfo(
          userInfo!.wallets![0].signer,
          secretKeys[0],
          userInfo!.wallets![0].publicKey,
          userInfo!.username,
          this,
          forceRefresh: true,
        ),
        getFiatRates(this),
      ]);

      notifyListeners();
    } catch (e) {}
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
          '/v1/shared-access/approvals?&excludeUserApproved=${excludeUserApproved}&limit=$limit${query}&page=${currentPage}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: activeWallet!.signer!,
        secretKey: secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );

      if (responseData['statusCode'] == 200) {
        inspect(responseData['data']);
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  getApprovals() {
    approvals = fetchApprovals(limit: limit.toString(), query: filterQuery);
  }

  late List<DepositTransactionModel> depositHistoryData =
      <DepositTransactionModel>[];

  Future<void> fetchDepositHistory(
    context, {
    required String publicKey,
    required String? currency,
  }) async {
    try {
      showLoader(context);
      var uri =
          '/v1/crypto/deposit-history/$currency/$publicKey?limit=$limit${filterQuery}';

      Map responseData = await makeGetRequest(
        uri: uri,
        signer: activeWallet!.signer!,
        secretKey: secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );

      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        totalRecords = responseData['data']['totalRecords'];
        currentPage = responseData['data']['currentPage'];
        var list = <DepositTransactionModel>[];
        for (var i = 0; i < responseData['data']['records'].length; i++) {
          list.add(
            DepositTransactionModel.deserializeJson(
              responseData['data']['records'][i],
            ),
          );
        }
        depositHistoryData = list;
        notifyListeners();
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      hideLoader(context);
      return Future.error('Error! ${e}');
    }
  }

  late List<WithdrawalTransactionModel> withdrawalHistoryData =
      <WithdrawalTransactionModel>[];

  String filterWithdrawalAddress = "";
  set setFilterWithdrawalAddress(value) {
    filterWithdrawalAddress = value;
    notifyListeners();
  }

  String filterWithdrawalStatus = "";
  set setFilterWithdrawalStatus(value) {
    filterWithdrawalStatus = value;
    notifyListeners();
  }

  Future<void> fetchWithdrawalHistory(
    context, {
    required String publicKey,
    required String? currency,
  }) async {
    try {
      showLoader(context);
      var uri =
          '/v1/crypto/withdrawal-history/$currency/$publicKey?limit=$limit${filterQuery}';

      Map responseData = await makeGetRequest(
        uri: uri,
        signer: activeWallet!.signer!,
        secretKey: secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.signer!,
      );

      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        totalRecords = responseData['data']['totalRecords'];
        currentPage = responseData['data']['currentPage'];
        var list = <WithdrawalTransactionModel>[];
        for (var i = 0; i < responseData['data']['records'].length; i++) {
          list.add(
            WithdrawalTransactionModel.deserializeJson(
              responseData['data']['records'][i],
            ),
          );
        }
        withdrawalHistoryData = list;
        notifyListeners();
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      hideLoader(context);
      return Future.error('Error! ${e}');
    }
  }

  String? pdfUrl;

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
  bool clearAccessList = false;

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

  void processDeepLink(
    BuildContext context,
    Uri initialDynamicLink, {
    String? rel,
    void Function()? onCancel,
  }) async {
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
        pages: [LoginPageConfig, AuthorizeLoginViewPageConfig],
      );
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
              "assetCode":
                  initialDynamicLink.queryParameters['assetCode'] == 'XBN'
                  ? ''
                  : initialDynamicLink.queryParameters['assetCode'],
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

            viewData![SendAssetViewPageConfig.key]['deepLinkInfo'] =
                deeplinkInfo;
            currentAction = PageAction(
              state: PageState.addAll,
              pages: isLoggedIn
                  ? [BottomHomePageConfig, SendAssetViewPageConfig]
                  : [LoginPageConfig, SendAssetViewPageConfig],
            );
            setSplashFinished();
          },
          onCancel: () {
            if (rel == 'qrScanner') {
              Navigator.of(context).pop();
              if (onCancel != null) onCancel();
            } else {
              currentAction = PageAction(
                state: PageState.addAll,
                pages: isLoggedIn ? [BottomHomePageConfig] : [LoginPageConfig],
              );
            }
            setSplashFinished();
          },
        );
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
            initialDynamicLink.queryParameters['serviceShortName'],
      };
      // action authorize
      setSplashFinished();
      currentAction = PageAction(
        state: PageState.addAll,
        pages: [LoginPageConfig, AuthorizeActionViewPageConfig],
      );
    } else if (initialDynamicLink.queryParameters['action'] == 'register') {
      setSplashFinished();
      setTempReferrerUsername = initialDynamicLink.queryParameters['referrer'];

      // action register
      currentAction = PageAction(
        state: PageState.addAll,
        pages: [LoginPageConfig, CreatePasswordPageConfig],
      );
    } else if (initialDynamicLink.queryParameters['action'] ==
        'tokenizedAsset') {
      try {
        var assetCode = initialDynamicLink.queryParameters['assetCode'];
        tokenizedAsset = await fetchTokenizedAsset(assetCode: assetCode!);
        setSplashFinished();

        currentAction = PageAction(
          state: PageState.addAll,
          pages: isLoggedIn
              ? [BottomHomePageConfig, TokenizedAssetDetailViewPageConfig]
              : [LoginPageConfig, TokenizedAssetDetailViewPageConfig],
        );
        notifyListeners;
      } catch (e) {
        setSplashFinished();
      }
    }
  }

  Future<void> fetchTokenizationData() async {
    var uri = '/v1/tokenization';

    Map responseData = await makeGetRequest(
      uri: Uri.encodeFull(uri),
      signer: primaryWallet.signer!,
      secretKey: secretKeys[0], // the primary wallet secret key
      publicKey: primaryWallet.signer!,
    );
    if (responseData['statusCode'] == 200) {
      tokenizationData = responseData['data'];
    }
  }

  Future<TokenizedAsset> fetchTokenizedAsset({
    required String assetCode,
  }) async {
    try {
      await Future.wait([
        if (tokenizationData.isEmpty) fetchTokenizationData(),
      ]);
      var uri =
          '/v1/tokenization/list?onlyWithUserPermission=0&assetCode=$assetCode';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: primaryWallet.signer!,
        secretKey: secretKeys[0], // the primary wallet secret key
        publicKey: primaryWallet.signer!,
      );

      inspect(responseData['data']);

      if (responseData['statusCode'] == 200) {
        return TokenizedAsset().deserializeJson(
          responseData['data']['records'][0],
        );
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Map expressedInterests = {};
  Map subscriptions = {};
}
