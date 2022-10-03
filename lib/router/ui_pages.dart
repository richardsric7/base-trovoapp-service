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
const String SwapAssetsViewPath = '/SwapAssetsView';
const String ConfirmSwapViewPath = '/ConfirmSwapView';
const String SwapSuccessViewPath = '/SwapSuccessView';
const String ProfileDetailsViewPath = '/ProfileDetailsView';
const String ReferralInfoViewPath = '/ReferralInfoView';
const String PasswordMgtViewPath = '/PasswordMgtView';
const String BackupAllViewPath = '/BackupAllView';
const String AuthorizeLoginViewPath = '/AuthorizeLoginView';
const String AuthorizeActionViewPath = '/AuthorizeActionView';
const String RequestSpecificPaymentViewPath = '/RequestSpecificPaymentView';
const String RequestSpecificPaymentDetailsViewPath =
    '/RequestSpecificPaymentDetailsView';
const String SecurityQuestionsViewPath = '/SecurityQuestionsView';
const String RequestOtpViewPath = '/RequestOtpView';
const String AccountRecoverySuccessViewPath = '/AccountRecoverySuccessView';
const String SharedAccessViewPath = '/SharedAccessView';
const String SetupAccountRecoveryViewPath = '/SetupAccountRecoveryView';
const String DisableAccountRecoveryViewPath = '/DisableAccountRecoveryView';
const String RecoverAccountViewPath = '/RecoverAccountView';
const String AnswerSecurityQuestionsViewPath = '/AnswerSecurityQuestionsView';
const String RequestBackupViewPath = '/RequestBackupView';
const String BackupRecoverySecretViewPath = '/BackupRecoverySecretView';
const String CompleteAccountRecoveryViewPath = '/CompleteAccountRecoveryView';
const String DisableAccountRecoveryInfoViewPath =
    '/DisableAccountRecoveryInfoView';
const String SuccessViewPath = '/SuccessView';

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
  SwapAssetsView,
  ConfirmSwapView,
  SwapSuccessView,
  ProfileDetailsView,
  ReferralInfoView,
  PasswordMgtView,
  BackupAllView,
  AuthorizeLoginView,
  AuthorizeActionView,
  RequestSpecificPaymentView,
  RequestSpecificPaymentDetailsView,
  SecurityQuestionsView,
  RequestOtpView,
  AccountRecoverySuccessView,
  SharedAccessView,
  SetupAccountRecoveryView,
  DisableAccountRecoveryView,
  RecoverAccountView,
  AnswerSecurityQuestionsView,
  RequestBackupView,
  BackupRecoverySecretView,
  CompleteAccountRecoveryView,
  DisableAccountRecoveryInfoView,
  SuccessView,
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
PageConfiguration SwapAssetsViewPageConfig = PageConfiguration(
    key: 'SwapAssetsView',
    path: SwapAssetsViewPath,
    uiPage: Pages.SwapAssetsView,
    currentPageAction: null);
PageConfiguration ConfirmSwapViewPageConfig = PageConfiguration(
    key: 'ConfirmSwapView',
    path: ConfirmSwapViewPath,
    uiPage: Pages.ConfirmSwapView,
    currentPageAction: null);
PageConfiguration SwapSuccessViewPageConfig = PageConfiguration(
    key: 'SwapSuccessView',
    path: SwapSuccessViewPath,
    uiPage: Pages.SwapSuccessView,
    currentPageAction: null);
PageConfiguration ProfileDetailsViewPageConfig = PageConfiguration(
    key: 'ProfileDetailsView',
    path: ProfileDetailsViewPath,
    uiPage: Pages.ProfileDetailsView,
    currentPageAction: null);
PageConfiguration ReferralInfoViewPageConfig = PageConfiguration(
    key: 'ReferralInfoView',
    path: ReferralInfoViewPath,
    uiPage: Pages.ReferralInfoView,
    currentPageAction: null);
PageConfiguration PasswordMgtViewPageConfig = PageConfiguration(
    key: 'PasswordMgtView',
    path: PasswordMgtViewPath,
    uiPage: Pages.PasswordMgtView,
    currentPageAction: null);
PageConfiguration BackupAllViewPageConfig = PageConfiguration(
    key: 'BackupAllView',
    path: BackupAllViewPath,
    uiPage: Pages.BackupAllView,
    currentPageAction: null);
PageConfiguration AuthorizeLoginViewPageConfig = PageConfiguration(
    key: 'AuthorizeLoginView',
    path: AuthorizeLoginViewPath,
    uiPage: Pages.AuthorizeLoginView,
    currentPageAction: null);
PageConfiguration AuthorizeActionViewPageConfig = PageConfiguration(
    key: 'AuthorizeActionView',
    path: AuthorizeActionViewPath,
    uiPage: Pages.AuthorizeActionView,
    currentPageAction: null);
PageConfiguration RequestSpecificPaymentViewPageConfig = PageConfiguration(
    key: 'RequestSpecificPaymentView',
    path: RequestSpecificPaymentViewPath,
    uiPage: Pages.RequestSpecificPaymentView,
    currentPageAction: null);
PageConfiguration RequestSpecificPaymentDetailsViewPageConfig =
    PageConfiguration(
        key: 'RequestSpecificPaymentDetailsView',
        path: RequestSpecificPaymentDetailsViewPath,
        uiPage: Pages.RequestSpecificPaymentDetailsView,
        currentPageAction: null);
PageConfiguration SecurityQuestionsViewPageConfig = PageConfiguration(
    key: 'SecurityQuestionsView',
    path: SecurityQuestionsViewPath,
    uiPage: Pages.SecurityQuestionsView,
    currentPageAction: null);
PageConfiguration RequestOtpViewPageConfig = PageConfiguration(
    key: 'RequestOtpView',
    path: RequestOtpViewPath,
    uiPage: Pages.RequestOtpView,
    currentPageAction: null);
PageConfiguration AccountRecoverySuccessViewPageConfig = PageConfiguration(
    key: 'AccountRecoverySuccessView',
    path: AccountRecoverySuccessViewPath,
    uiPage: Pages.AccountRecoverySuccessView,
    currentPageAction: null);
PageConfiguration SharedAccessViewPageConfig = PageConfiguration(
    key: 'SharedAccessView',
    path: SharedAccessViewPath,
    uiPage: Pages.SharedAccessView,
    currentPageAction: null);
PageConfiguration SetupAccountRecoveryViewPageConfig = PageConfiguration(
    key: 'SetupAccountRecoveryView',
    path: SetupAccountRecoveryViewPath,
    uiPage: Pages.SetupAccountRecoveryView,
    currentPageAction: null);
PageConfiguration DisableAccountRecoveryViewPageConfig = PageConfiguration(
    key: 'DisableAccountRecoveryView',
    path: DisableAccountRecoveryViewPath,
    uiPage: Pages.DisableAccountRecoveryView,
    currentPageAction: null);
PageConfiguration RecoverAccountViewPageConfig = PageConfiguration(
    key: 'RecoverAccountView',
    path: RecoverAccountViewPath,
    uiPage: Pages.RecoverAccountView,
    currentPageAction: null);
PageConfiguration AnswerSecurityQuestionsViewPageConfig = PageConfiguration(
    key: 'AnswerSecurityQuestionsView',
    path: AnswerSecurityQuestionsViewPath,
    uiPage: Pages.AnswerSecurityQuestionsView,
    currentPageAction: null);
PageConfiguration RequestBackupViewPageConfig = PageConfiguration(
    key: 'RequestBackupView',
    path: RequestBackupViewPath,
    uiPage: Pages.RequestBackupView,
    currentPageAction: null);
PageConfiguration BackupRecoverySecretViewPageConfig = PageConfiguration(
    key: 'BackupRecoverySecretView',
    path: BackupRecoverySecretViewPath,
    uiPage: Pages.BackupRecoverySecretView,
    currentPageAction: null);
PageConfiguration CompleteAccountRecoveryViewPageConfig = PageConfiguration(
    key: 'CompleteAccountRecoveryView',
    path: CompleteAccountRecoveryViewPath,
    uiPage: Pages.CompleteAccountRecoveryView,
    currentPageAction: null);
PageConfiguration DisableAccountRecoveryInfoViewPageConfig = PageConfiguration(
    key: 'DisableAccountRecoveryInfoView',
    path: DisableAccountRecoveryInfoViewPath,
    uiPage: Pages.DisableAccountRecoveryInfoView,
    currentPageAction: null);
PageConfiguration SuccessViewPageConfig = PageConfiguration(
    key: 'SuccessView',
    path: SuccessViewPath,
    uiPage: Pages.SuccessView,
    currentPageAction: null);
