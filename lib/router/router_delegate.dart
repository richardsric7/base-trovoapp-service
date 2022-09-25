import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/Custom_BlocObserver/swiper/swiper.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/confirm_swap.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/change_password.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_detail.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/profile_details.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/referral_info.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/subwalletCreateSuccess.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/notificationsView.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/searchview.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/swap_assets.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/swap_success.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallet_details.dart';
import 'package:trovo_wallet/bottom_bar/bottombar.dart';
import 'package:trovo_wallet/screens/AccountRecovery/account_recovery_success.dart';
import 'package:trovo_wallet/screens/Auth/AuthorizeActionView.dart';
import 'package:trovo_wallet/screens/Auth/AuthorizeLoginView.dart';
import 'package:trovo_wallet/screens/Auth/create_password.dart';
import 'package:trovo_wallet/screens/Auth/signup.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';
import 'package:trovo_wallet/screens/Backup/backup_all.dart';
import 'package:trovo_wallet/screens/Backup/congratulation.dart';
import 'package:trovo_wallet/screens/Backup/ensure_privacy.dart';
import 'package:trovo_wallet/screens/ImportWallet/importwallet.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/asset_details.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/confirm_transaction.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/recieve_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/request_specific_payment_details.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/request_specific_payment.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/send_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/transaction_success.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/trust_asset.dart';
import 'package:trovo_wallet/screens/AccountRecovery/request_otp.dart';
import 'package:trovo_wallet/screens/AccountRecovery/security_questions.dart';
import 'package:trovo_wallet/screens/SharedAccess/shared_access.dart';
import 'package:trovo_wallet/screens/page_view/web_view.dart';
import 'package:trovo_wallet/storage/state.dart';
import '../screens/Auth/fingerprint.dart';
import '../screens/Auth/login.dart';
import '../screens/Backup/backup.dart';
import '../screens/Splash_Screen/splashscreen.dart';
import '../screens/qr_scanner_view.dart';
import 'PageActions.dart';
import 'back_dispatcher.dart';
import 'ui_pages.dart';

class TrovoWalletRouterDelegate extends RouterDelegate<PageConfiguration>
    with ChangeNotifier, PopNavigatorRouterDelegateMixin<PageConfiguration> {
  final List<MaterialPage> _pages = [];
  TrovoWalletBackButtonDispatcher? backButtonDispatcher;

  @override
  final GlobalKey<NavigatorState> navigatorKey;
  final DataProvider appState;

  TrovoWalletRouterDelegate(this.appState) : navigatorKey = GlobalKey() {
    appState.addListener(() {
      notifyListeners();
    });
  }

  /// Getter for a list that cannot be changed
  List<MaterialPage> get pages => List.unmodifiable(_pages);

  /// Number of pages function
  int numPages() => _pages.length;

  @override
  PageConfiguration get currentConfiguration =>
      _pages.last.arguments as PageConfiguration;

  @override
  Widget build(BuildContext context) {
    return Navigator(
      key: navigatorKey,
      onPopPage: _onPopPage,
      pages: buildPages(),
    );
  }

  bool _onPopPage(Route<dynamic> route, result) {
    final didPop = route.didPop(result);
    if (!didPop) {
      return false;
    }
    if (canPop()) {
      pop();
      return true;
    } else {
      return false;
    }
  }

  void _removePage(MaterialPage page) {
    if (page != null) {
      _pages.remove(page);
    }
  }

  void pop() {
    if (canPop()) {
      _removePage(_pages.last);
    }
  }

  bool canPop() {
    return _pages.length > 1;
  }

  @override
  Future<bool> popRoute() {
    if (canPop()) {
      _removePage(_pages.last);
      notifyListeners();
      return Future.value(true);
    }
    return Future.value(false);
  }

  MaterialPage _createPage(Widget child, PageConfiguration pageConfig) {
    return MaterialPage(
        child: child,
        key: ValueKey(pageConfig.key),
        name: pageConfig.path,
        arguments: pageConfig);
  }

  void _addPageData(Widget child, PageConfiguration pageConfig) {
    _pages.add(
      _createPage(child, pageConfig),
    );
  }

  void addPage(PageConfiguration pageConfig) {
    final shouldAddPage = _pages.isEmpty ||
        (_pages.last.arguments as PageConfiguration).uiPage !=
            pageConfig.uiPage;
    if (shouldAddPage) {
      switch (pageConfig.uiPage) {
        case Pages.Splash:
          _addPageData(const SplashScreen(), SplashPageConfig);
          break;
        case Pages.Login:
          _addPageData(const Login(), LoginPageConfig);
          break;
        case Pages.Onboarding:
          _addPageData(const Swiper(), OnboardingPageConfig);
          break;
        case Pages.Signup:
          _addPageData(const SignUp(), SignupPageConfig);
          break;
        case Pages.CreatePassword:
          _addPageData(const CreatePassword(), CreatePasswordPageConfig);
          break;
        case Pages.ImportWallet:
          _addPageData(const ImportWallet(), ImportWalletPageConfig);
          break;
        case Pages.Verification:
          _addPageData(const Veryfication(), VerificationPageConfig);
          break;
        case Pages.Congratulations:
          _addPageData(Congratulations(), CongratulationsPageConfig);
          break;
        case Pages.EnsurePrivacy:
          _addPageData(EnsurePrivacy(), EnsurePrivacyPageConfig);
          break;
        case Pages.Backup:
          _addPageData(Backup(), BackupPageConfig);
          break;
        case Pages.Fingerprint:
          _addPageData(FingerPrint(), FingerprintPageConfig);
          break;
        case Pages.BottomHome:
          _addPageData(BottomHome(), BottomHomePageConfig);
          break;
        case Pages.WebView:
          _addPageData(TrovoWebView(), WebViewPageConfig);
          break;
        case Pages.QrScanner:
          _addPageData(QrScanner(), QrScannerPageConfig);
          break;
        case Pages.SearchView:
          _addPageData(SearchView(), SearchViewPageConfig);
          break;
        case Pages.NotificationsView:
          _addPageData(NotificationsView(), NotificationsViewPageConfig);
          break;
        case Pages.CreateSubWalletSuccessView:
          _addPageData(CreateSubWalletSuccessView(),
              CreateSubWalletSuccessViewPageConfig);
          break;
        case Pages.WalletDetailsView:
          _addPageData(WalletDetails(), WalletDetailsViewPageConfig);
          break;
        case Pages.AssetDetailsView:
          _addPageData(AssetDetails(), AssetDetailsViewPageConfig);
          break;
        case Pages.SendAssetView:
          _addPageData(SendAsset(), SendAssetViewPageConfig);
          break;
        case Pages.ConfirmTransactionView:
          _addPageData(ConfirmTransaction(), ConfirmTransactionViewPageConfig);
          break;
        case Pages.TransactionSuccessView:
          _addPageData(TransactionSuccess(), TransactionSuccessViewPageConfig);
          break;
        case Pages.ReceiveAssetView:
          _addPageData(ReceiveAsset(), ReceiveAssetViewPageConfig);
          break;
        case Pages.PendingAssetDetailsView:
          _addPageData(
              PendingAssetDetails(), PendingAssetDetailsViewPageConfig);
          break;
        case Pages.PaymentHistoryView:
          _addPageData(PaymentHistory(), PaymentHistoryViewPageConfig);
          break;
        case Pages.PaymentDetailsView:
          _addPageData(PaymentDetails(), PaymentDetailsViewPageConfig);
          break;
        case Pages.SwapAssetsView:
          _addPageData(SwapAssets(), SwapAssetsViewPageConfig);
          break;
        case Pages.ConfirmSwapView:
          _addPageData(ConfirmSwap(), ConfirmSwapViewPageConfig);
          break;
        case Pages.SwapSuccessView:
          _addPageData(SwapSuccess(), SwapSuccessViewPageConfig);
          break;
        case Pages.ProfileDetailsView:
          _addPageData(ProfileDetails(), ProfileDetailsViewPageConfig);
          break;
        case Pages.ReferralInfoView:
          _addPageData(ReferralInfo(), ReferralInfoViewPageConfig);
          break;
        case Pages.PasswordMgtView:
          _addPageData(PasswordMgtView(), PasswordMgtViewPageConfig);
          break;
        case Pages.BackupAllView:
          _addPageData(BackupAll(), BackupAllViewPageConfig);
          break;
        case Pages.AuthorizeLoginView:
          _addPageData(AuthorizeLoginView(), AuthorizeLoginViewPageConfig);
          break;
        case Pages.AuthorizeActionView:
          _addPageData(AuthorizeActionView(), AuthorizeActionViewPageConfig);
          break;
        case Pages.RequestSpecificPaymentView:
          _addPageData(
              RequestSpecificPayment(), RequestSpecificPaymentViewPageConfig);
          break;
        case Pages.RequestSpecificPaymentDetailsView:
          _addPageData(RequestSpecificPaymentDetails(),
              RequestSpecificPaymentDetailsViewPageConfig);
          break;
        case Pages.SecurityQuestionsView:
          _addPageData(SecurityQuestions(), SecurityQuestionsViewPageConfig);
          break;
        case Pages.RequestOtpView:
          _addPageData(RequestOtp(), RequestOtpViewPageConfig);
          break;
        case Pages.AccountRecoverySuccessView:
          _addPageData(
              AccountRecoverySuccess(), AccountRecoverySuccessViewPageConfig);
          break;
        case Pages.SharedAccessView:
          _addPageData(SharedAccess(), SharedAccessViewPageConfig);
          break;
        default:
          break;
      }
    }
  }

  void replace(PageConfiguration newRoute) {
    if (_pages.isNotEmpty) {
      _pages.removeLast();
    }
    addPage(newRoute);
  }

  void setPath(List<MaterialPage> path) {
    _pages.clear();
    _pages.addAll(path);
  }

  void replaceAll(PageConfiguration newRoute) {
    setNewRoutePath(newRoute);
  }

  void push(PageConfiguration newRoute) {
    addPage(newRoute);
  }

  void pushWidget(Widget child, PageConfiguration newRoute) {
    _addPageData(child, newRoute);
  }

  void addAll(List<PageConfiguration> routes) {
    _pages.clear();
    routes.forEach((route) {
      addPage(route);
    });
  }

  @override
  Future<void> setNewRoutePath(PageConfiguration configuration) {
    final shouldAddPage = _pages.isEmpty ||
        (_pages.last.arguments as PageConfiguration).uiPage !=
            configuration.uiPage;
    if (shouldAddPage) {
      _pages.clear();
      addPage(configuration);
    }
    return SynchronousFuture(null);
  }

  void _setPageAction(PageAction action) {
    switch (action.page?.uiPage) {
      case Pages.Splash:
        SplashPageConfig.currentPageAction = action;
        break;
      case Pages.Login:
        LoginPageConfig.currentPageAction = action;
        break;
      case Pages.Signup:
        SignupPageConfig.currentPageAction = action;
        break;
      case Pages.Onboarding:
        OnboardingPageConfig.currentPageAction = action;
        break;
      case Pages.CreatePassword:
        CreatePasswordPageConfig.currentPageAction = action;
        break;
      case Pages.ImportWallet:
        ImportWalletPageConfig.currentPageAction = action;
        break;
      case Pages.Verification:
        VerificationPageConfig.currentPageAction = action;
        break;
      case Pages.Congratulations:
        CongratulationsPageConfig.currentPageAction = action;
        break;
      case Pages.EnsurePrivacy:
        EnsurePrivacyPageConfig.currentPageAction = action;
        break;
      case Pages.Backup:
        BackupPageConfig.currentPageAction = action;
        break;
      case Pages.Fingerprint:
        FingerprintPageConfig.currentPageAction = action;
        break;
      case Pages.BottomHome:
        BottomHomePageConfig.currentPageAction = action;
        break;
      case Pages.WebView:
        WebViewPageConfig.currentPageAction = action;
        break;
      case Pages.QrScanner:
        QrScannerPageConfig.currentPageAction = action;
        break;
      case Pages.SearchView:
        SearchViewPageConfig.currentPageAction = action;
        break;
      case Pages.NotificationsView:
        NotificationsViewPageConfig.currentPageAction = action;
        break;
      case Pages.CreateSubWalletSuccessView:
        CreateSubWalletSuccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.WalletDetailsView:
        WalletDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.AssetDetailsView:
        AssetDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.SendAssetView:
        SendAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.ConfirmTransactionView:
        ConfirmTransactionViewPageConfig.currentPageAction = action;
        break;
      case Pages.TransactionSuccessView:
        TransactionSuccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.ReceiveAssetView:
        ReceiveAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.PendingAssetDetailsView:
        PendingAssetDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.PaymentHistoryView:
        PaymentHistoryViewPageConfig.currentPageAction = action;
        break;
      case Pages.PaymentDetailsView:
        PaymentDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.SwapAssetsView:
        SwapAssetsViewPageConfig.currentPageAction = action;
        break;
      case Pages.ConfirmSwapView:
        ConfirmSwapViewPageConfig.currentPageAction = action;
        break;
      case Pages.SwapSuccessView:
        SwapSuccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.ProfileDetailsView:
        ProfileDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.ReferralInfoView:
        ReferralInfoViewPageConfig.currentPageAction = action;
        break;
      case Pages.PasswordMgtView:
        PasswordMgtViewPageConfig.currentPageAction = action;
        break;
      case Pages.BackupAllView:
        BackupAllViewPageConfig.currentPageAction = action;
        break;
      case Pages.AuthorizeLoginView:
        AuthorizeLoginViewPageConfig.currentPageAction = action;
        break;
      case Pages.AuthorizeActionView:
        AuthorizeActionViewPageConfig.currentPageAction = action;
        break;
      case Pages.RequestSpecificPaymentView:
        RequestSpecificPaymentViewPageConfig.currentPageAction = action;
        break;
      case Pages.RequestSpecificPaymentDetailsView:
        RequestSpecificPaymentDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.SecurityQuestionsView:
        SecurityQuestionsViewPageConfig.currentPageAction = action;
        break;
      case Pages.RequestOtpView:
        RequestOtpViewPageConfig.currentPageAction = action;
        break;
      case Pages.AccountRecoverySuccessView:
        AccountRecoverySuccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.SharedAccessView:
        SharedAccessViewPageConfig.currentPageAction = action;
        break;
      default:
        break;
    }
  }

  List<Page> buildPages() {
    if (!appState.splashFinished) {
      replaceAll(SplashPageConfig);
    } else {
      switch (appState.currentAction.state) {
        case PageState.none:
          break;
        case PageState.addPage:
          _setPageAction(appState.currentAction);
          addPage(appState.currentAction.page!);
          break;
        case PageState.pop:
          pop();
          break;
        case PageState.replace:
          _setPageAction(appState.currentAction);
          replace(appState.currentAction.page!);
          break;
        case PageState.replaceAll:
          _setPageAction(appState.currentAction);
          replaceAll(appState.currentAction.page!);
          break;
        case PageState.addWidget:
          _setPageAction(appState.currentAction);
          pushWidget(
              appState.currentAction.widget!, appState.currentAction.page!);
          break;
        case PageState.addAll:
          addAll(appState.currentAction.pages!);
          break;
      }
    }
    appState.resetCurrentAction();
    return List.of(_pages);
  }

  void parseRoute(Uri uri) {
    if (uri.pathSegments.isEmpty) {
      setNewRoutePath(SplashPageConfig);
      return;
    }

    // Handle navapp://deeplinks/details/#
    if (uri.pathSegments.length == 2) {
      // if (uri.pathSegments[0] == 'details') {
      //   pushWidget(Details(int.parse(uri.pathSegments[1])), DetailsPageConfig);
      // }
    } else if (uri.pathSegments.length == 1) {
      final path = uri.pathSegments[0];
      switch (path) {
        case 'splash':
          replaceAll(SplashPageConfig);
          break;
        case 'login':
          replaceAll(LoginPageConfig);
          break;
        case 'onboarding':
          setPath([
            _createPage(const Swiper(), OnboardingPageConfig),
          ]);
          break;
        case 'signup':
          setPath([
            _createPage(const Login(), LoginPageConfig),
            _createPage(const CreatePassword(), CreatePasswordPageConfig)
          ]);
          break;
        case 'createPassword':
          setPath([
            _createPage(const Login(), LoginPageConfig),
            _createPage(const CreatePassword(), CreatePasswordPageConfig)
          ]);
          break;
        case 'importWallet':
          setPath([
            _createPage(const Login(), LoginPageConfig),
            _createPage(const ImportWallet(), ImportWalletPageConfig)
          ]);
          break;
        case 'verification':
          setPath([
            _createPage(const Login(), LoginPageConfig),
            _createPage(const CreatePassword(), CreatePasswordPageConfig)
          ]);
          break;
        case 'congratulations':
          setPath([
            _createPage(const Login(), LoginPageConfig),
            _createPage(const CreatePassword(), CreatePasswordPageConfig)
          ]);
          break;
        case 'ensurePrivacy':
          setPath([
            _createPage(const Login(), LoginPageConfig),
          ]);
          break;
        case 'backup':
          setPath([
            _createPage(const Login(), LoginPageConfig),
          ]);
          break;
        case 'fingerPrint':
          setPath([
            _createPage(const Login(), LoginPageConfig),
            _createPage(const FingerPrint(), FingerprintPageConfig),
          ]);
          break;
        case 'webview':
          setPath([
            _createPage(TrovoWebView(), WebViewPageConfig),
          ]);
          break;
        case 'assetDetailsView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const AssetDetails(), AssetDetailsViewPageConfig),
          ]);
          break;
        case 'sendAssetView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const AssetDetails(), AssetDetailsViewPageConfig),
            _createPage(const SendAsset(), SendAssetViewPageConfig),
          ]);
          break;
        case 'confirmTransactionView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const AssetDetails(), AssetDetailsViewPageConfig),
            _createPage(const SendAsset(), SendAssetViewPageConfig),
            _createPage(
                const ConfirmTransaction(), ConfirmTransactionViewPageConfig),
          ]);
          break;
        case 'ProfileDetailsView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const ProfileDetails(), ProfileDetailsViewPageConfig),
          ]);
          break;
        case 'ReferralInfoView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const ReferralInfo(), ReferralInfoViewPageConfig),
          ]);
          break;
        case 'PasswordMgtView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const PasswordMgtView(), PasswordMgtViewPageConfig),
          ]);
          break;
        case 'BackupAllView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
            _createPage(const BackupAll(), BackupAllViewPageConfig),
          ]);
          break;
        case 'AuthorizeLoginView':
          setPath([
            _createPage(AuthorizeLoginView(), AuthorizeLoginViewPageConfig),
          ]);
          break;
        case 'AuthorizeActionView':
          setPath([
            _createPage(AuthorizeActionView(), AuthorizeActionViewPageConfig),
          ]);
          break;
        case 'home':
        case 'qrscanner':
        case 'searchview':
        case 'notificationsview':
        case 'walletDetailsView':
        case 'transactionSuccessView':
        case 'createSubWalletSuccessView':
        case 'recieveAssetView':
        case 'requestSpecificPaymentView':
        case 'requestSpecificPaymentDetatilsView':
        case 'pendingAssetDetailsView':
        case 'PaymentHistoryView':
        case 'PaymentDetailsView':
        case 'SwapAssetsView':
        case 'ConfirmSwapView':
        case 'SwapSuccessView':
        case 'SecurityQuestionsView':
        case 'RequestOtpView':
        case 'AccountRecoverySuccessView':
        case 'SharedAccessView':
          setPath([
            _createPage(const BottomHome(), BottomHomePageConfig),
          ]);
          break;
        default:
          setPath([
            _createPage(const SplashScreen(), SplashPageConfig),
          ]);
      }
    }
  }
}
