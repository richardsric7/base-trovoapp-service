import 'PageActions.dart';

const String SplashPath = '/splash';
const String LoginPath = '/login';
const String OnboardingPath = '/onboarding';
const String CreatePasswordPath = '/createPassword';
const String ImportWalletPath = '/importWallet';
const String SignupPath = '/signup';
const String VerificationPath = '/verification';
const String CongratulationsPath = '/congratulations';
const String EnsurePrivacyPath = '/ensurePrivacy';
const String BackupPath = '/backup';
const String FingerprintPath = '/fingerprint';
const String BottomHomePath = '/home';
const String WebViewPath = '/webview';
const String QrScannerPath = '/qrScanner';
const String SearchViewPath = '/searchview';
const String NotificationsViewPath = '/notificationsview';
const String CreateSubWalletSuccessViewPath = '/createSubWalletSuccessView';
const String WalletDetailsViewPath = '/WalletDetailsView';
const String AssetDetailsViewPath = '/AssetDetailsView';
const String SendAssetViewPath = '/sendAssetView';
const String ConfirmTransactionViewPath = '/confirmTransactionAssetView';
const String TransactionSuccessViewPath = '/transactionSuccessView';
const String ReceiveAssetViewPath = '/recieveAssetView';
const String PendingAssetDetailsViewPath = '/pendingAssetView';
const String PaymentHistoryViewPath = '/paymentHistoryView';
const String PaymentDetailsViewPath = '/PaymentDetailsView';

enum Pages {
  Splash,
  Login,
  Onboarding,
  Signup,
  CreatePassword,
  ImportWallet,
  Verification,
  Congratulations,
  EnsurePrivacy,
  Backup,
  Fingerprint,
  BottomHome,
  WebView,
  QrScanner,
  SearchView,
  NotificationsView,
  CreateSubWalletSuccessView,
  WalletDetailsView,
  AssetDetailsView,
  SendAssetView,
  ConfirmTransactionView,
  TransactionSuccessView,
  ReceiveAssetView,
  PendingAssetDetailsView,
  PaymentHistoryView,
  PaymentDetailsView,
}

class PageConfiguration {
  final String key;
  final String path;
  final Pages uiPage;
  PageAction? currentPageAction;

  PageConfiguration(
      {required this.key,
      required this.path,
      required this.uiPage,
      this.currentPageAction});
}

PageConfiguration SplashPageConfig = PageConfiguration(
    key: 'Splash',
    path: SplashPath,
    uiPage: Pages.Splash,
    currentPageAction: null);
PageConfiguration LoginPageConfig = PageConfiguration(
    key: 'Login',
    path: LoginPath,
    uiPage: Pages.Login,
    currentPageAction: null);
PageConfiguration SignupPageConfig = PageConfiguration(
    key: 'Signup',
    path: SignupPath,
    uiPage: Pages.Signup,
    currentPageAction: null);
PageConfiguration OnboardingPageConfig = PageConfiguration(
    key: 'Onboarding',
    path: OnboardingPath,
    uiPage: Pages.Onboarding,
    currentPageAction: null);
PageConfiguration CreatePasswordPageConfig = PageConfiguration(
    key: 'CreatePassword',
    path: CreatePasswordPath,
    uiPage: Pages.CreatePassword,
    currentPageAction: null);
PageConfiguration ImportWalletPageConfig = PageConfiguration(
    key: 'ImportWallet',
    path: ImportWalletPath,
    uiPage: Pages.ImportWallet,
    currentPageAction: null);
PageConfiguration VerificationPageConfig = PageConfiguration(
    key: 'Verification',
    path: VerificationPath,
    uiPage: Pages.Verification,
    currentPageAction: null);
PageConfiguration CongratulationsPageConfig = PageConfiguration(
    key: 'Congratulations',
    path: CongratulationsPath,
    uiPage: Pages.Congratulations,
    currentPageAction: null);
PageConfiguration EnsurePrivacyPageConfig = PageConfiguration(
    key: 'EnsurePrivacy',
    path: EnsurePrivacyPath,
    uiPage: Pages.EnsurePrivacy,
    currentPageAction: null);
PageConfiguration BackupPageConfig = PageConfiguration(
    key: 'Backup',
    path: BackupPath,
    uiPage: Pages.Backup,
    currentPageAction: null);
PageConfiguration FingerprintPageConfig = PageConfiguration(
    key: 'Fingerprint',
    path: FingerprintPath,
    uiPage: Pages.Fingerprint,
    currentPageAction: null);
PageConfiguration BottomHomePageConfig = PageConfiguration(
    key: 'BottomHome',
    path: BottomHomePath,
    uiPage: Pages.BottomHome,
    currentPageAction: null);
PageConfiguration WebViewPageConfig = PageConfiguration(
    key: 'WebView',
    path: WebViewPath,
    uiPage: Pages.WebView,
    currentPageAction: null);
PageConfiguration QrScannerPageConfig = PageConfiguration(
    key: 'QrScanner',
    path: QrScannerPath,
    uiPage: Pages.QrScanner,
    currentPageAction: null);
PageConfiguration SearchViewPageConfig = PageConfiguration(
    key: 'SearchView',
    path: SearchViewPath,
    uiPage: Pages.SearchView,
    currentPageAction: null);
PageConfiguration NotificationsViewPageConfig = PageConfiguration(
    key: 'NotificationsView',
    path: NotificationsViewPath,
    uiPage: Pages.NotificationsView,
    currentPageAction: null);
PageConfiguration CreateSubWalletSuccessViewPageConfig = PageConfiguration(
    key: 'CreateSubWalletSuccessView',
    path: CreateSubWalletSuccessViewPath,
    uiPage: Pages.CreateSubWalletSuccessView,
    currentPageAction: null);
PageConfiguration WalletDetailsViewPageConfig = PageConfiguration(
    key: 'WalletDetailsView',
    path: WalletDetailsViewPath,
    uiPage: Pages.WalletDetailsView,
    currentPageAction: null);
PageConfiguration AssetDetailsViewPageConfig = PageConfiguration(
    key: 'AssetDetailsView',
    path: AssetDetailsViewPath,
    uiPage: Pages.AssetDetailsView,
    currentPageAction: null);
PageConfiguration SendAssetViewPageConfig = PageConfiguration(
    key: 'SendAssetView',
    path: SendAssetViewPath,
    uiPage: Pages.SendAssetView,
    currentPageAction: null);
PageConfiguration ConfirmTransactionViewPageConfig = PageConfiguration(
    key: 'ConfirmTransactionView',
    path: ConfirmTransactionViewPath,
    uiPage: Pages.ConfirmTransactionView,
    currentPageAction: null);
PageConfiguration TransactionSuccessViewPageConfig = PageConfiguration(
    key: 'TransactionSuccessView',
    path: TransactionSuccessViewPath,
    uiPage: Pages.TransactionSuccessView,
    currentPageAction: null);
PageConfiguration ReceiveAssetViewPageConfig = PageConfiguration(
    key: 'ReceiveAssetView',
    path: ReceiveAssetViewPath,
    uiPage: Pages.ReceiveAssetView,
    currentPageAction: null);
PageConfiguration PendingAssetDetailsViewPageConfig = PageConfiguration(
    key: 'PendingAssetDetailsView',
    path: PendingAssetDetailsViewPath,
    uiPage: Pages.PendingAssetDetailsView,
    currentPageAction: null);
PageConfiguration PaymentHistoryViewPageConfig = PageConfiguration(
    key: 'PaymentHistoryView',
    path: PaymentHistoryViewPath,
    uiPage: Pages.PaymentHistoryView,
    currentPageAction: null);
PageConfiguration PaymentDetailsViewPageConfig = PageConfiguration(
    key: 'PaymentDetailsView',
    path: PaymentDetailsViewPath,
    uiPage: Pages.PaymentDetailsView,
    currentPageAction: null);
