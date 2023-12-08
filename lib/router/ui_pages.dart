import 'page_actions.dart';

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
const String SecurityQuestionsForInactiveAccountsViewPath =
    '/SecurityQuestionsForInactiveAccounts';
const String SharedAccessViewPath = '/SharedAccessView';
const String AddSharedAccessDetailsViewPath = '/AddSharedAccessDetailsView';
const String SharedWalletInfoViewPath = '/SharedWalletInfoView';
const String ApprovalDetailsViewPath = '/ApprovalDetailsView';
const String UpdateSharedAccessViewPath = '/UpdateSharedAccessView';
const String UpdateSharedAccessDetailsViewPath =
    '/UpdateSharedAccessDetailsView';
const String WelcomeToSharedAccessViewPath = '/WelcomeToSharedAccessView';
const String GetStartedViewPath = '/GetStartedView';
const String ShareReceiptViewPath = '/ShareReceiptView';
const String AnnouncementViewPath = '/AnnouncementView';
const String WrappedAssetViewPath = '/WrappedAssetView';
const String GenerateDepositAddressViewPath = '/GenerateDepositAddressView';
const String SelectDepositAddressViewPath = '/SelectDepositAddressView';
const String WithdrawAssetViewPath = '/WithdrawAssetView';
const String ConfirmWithdrawViewPath = '/ConfirmWithdrawView';
const String TransactionStatusViewPath = '/TransactionStatusView';
const String DepositWithdrawHistoryViewPath = '/DepositWithdrawHistoryView';
const String DepositWithdrawDetailsViewPath = '/DepositWithdrawDetailsView';
const String WelcomeSubscriptionsViewPath = '/WelcomeSubscriptionsView';
const String SubscriptionPlansViewPath = '/SubscriptionPlansView';
const String SubscriptionPlanBenefitsViewPath = '/SubscriptionPlanBenefitsView';
const String AuthorizeSubscriptionViewPath = '/AuthorizeSubscriptionView';
const String OptInAssetViewPath = '/OptInAssetView';
const String OptOutAssetViewPath = '/OptOutAssetView';
const String OptInOutAssetViewPath = '/OptInOutAssetView';
const String TokenizationWelcomeViewPath = '/TokenizationWelcomeView';
const String SettingsViewPath = '/SettingsView';
const String TokenizeAssetViewPath = '/TokenizeAssetView';
const String AssetInformationViewPath = '/AssetInformationView';
const String AssetTokenInformationViewPath = '/AssetTokenInformationView';
const String AssetVerificationDocumentsViewPath =
    '/AssetVerificationDocumentsView';
const String TokenizedAssetDetailViewPath = '/TokenizedAssetDetailView';
const String BuyTokensViewPath = '/BuyTokensView';
const String ConfirmBuyViewPath = '/ConfirmBuyView';
const String TokenizedAssetsListViewPath = '/TokenizedAssetsListView';
const String AssetDashboardViewPath = '/AssetDashboardView';
const String AssetSubscribersViewPath = '/AssetSubscribersView';
const String TotalSalesViewPath = '/TotalSalesView';
const String ProceedsPayOutViewPath = '/ProceedsPayOutView';
const String LiquidateAssetViewPath = '/LiquidateAssetView';
const String WalletPreparationViewPath = '/WalletPreparationView';
const String MyAssetTokenDetailsViewPath = '/MyAssetTokenDetailsView';
const String MarketTradeViewPath = '/MarketTradeView';
const String MarketTradeInfoViewPath = '/MarketTradeInfoView';
const String MarketPairsViewPath = '/MarketPairsView';
const String SetupAndComplianceViewPath = '/SetupAndComplianceView';

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
  SecurityQuestionsForInactiveAccountsView,
  AddSharedAccessDetailsView,
  SharedWalletInfoView,
  SharedWalletDetailsView,
  SharedWalletAssetDetailsView,
  SendAssetSharedWalletView,
  RecieveAssetSharedWalletView,
  ConfirmInitiatePaymentView,
  ApprovalDetailsView,
  UpdateSharedAccessView,
  UpdateSharedAccessDetailsView,
  WelcomeToSharedAccessView,
  GetStartedView,
  ShareReceiptView,
  AnnouncementView,
  WrappedAssetView,
  GenerateDepositAddressView,
  SelectDepositAddressView,
  WithdrawAssetView,
  ConfirmWithdrawView,
  TransactionStatusView,
  DepositWithdrawHistoryView,
  DepositWithdrawDetailsView,
  WelcomeSubscriptionsView,
  SubscriptionPlansView,
  SubscriptionPlanBenefitsView,
  AuthorizeSubscriptionView,
  OptInAssetView,
  OptOutAssetView,
  OptInOutAssetView,
  TokenizationWelcomeView,
  SettingsView,
  TokenizeAssetView,
  AssetInformationView,
  AssetTokenInformationView,
  AssetVerificationDocumentsView,
  TokenizedAssetDetailView,
  BuyTokensView,
  ConfirmBuyView,
  TokenizedAssetsListView,
  AssetDashboardView,
  AssetSubscribersView,
  TotalSalesView,
  ProceedsPayOutView,
  LiquidateAssetView,
  WalletPreparationView,
  MyAssetTokenDetailsView,
  MarketTradeView,
  MarketTradeInfoView,
  MarketPairsView,
  SetupAndComplianceView,
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
PageConfiguration SecurityQuestionsForInactiveAccountsViewPageConfig =
    PageConfiguration(
        key: 'SecurityQuestionsForInactiveAccountsView',
        path: SecurityQuestionsForInactiveAccountsViewPath,
        uiPage: Pages.SecurityQuestionsForInactiveAccountsView,
        currentPageAction: null);
PageConfiguration AddSharedAccessDetailsViewPageConfig = PageConfiguration(
    key: 'AddSharedAccessDetailsView',
    path: AddSharedAccessDetailsViewPath,
    uiPage: Pages.AddSharedAccessDetailsView,
    currentPageAction: null);
PageConfiguration SharedWalletInfoViewPageConfig = PageConfiguration(
    key: 'SharedWalletInfoView',
    path: SharedWalletInfoViewPath,
    uiPage: Pages.SharedWalletInfoView,
    currentPageAction: null);
PageConfiguration ApprovalDetailsViewPageConfig = PageConfiguration(
    key: 'ApprovalDetailsView',
    path: ApprovalDetailsViewPath,
    uiPage: Pages.ApprovalDetailsView,
    currentPageAction: null);
PageConfiguration UpdateSharedAccessViewPageConfig = PageConfiguration(
    key: 'UpdateSharedAccessView',
    path: UpdateSharedAccessViewPath,
    uiPage: Pages.UpdateSharedAccessView,
    currentPageAction: null);
PageConfiguration UpdateSharedAccessDetailsViewPageConfig = PageConfiguration(
    key: 'UpdateSharedAccessDetailsView',
    path: UpdateSharedAccessDetailsViewPath,
    uiPage: Pages.UpdateSharedAccessDetailsView,
    currentPageAction: null);
PageConfiguration WelcomeToSharedAccessViewPageConfig = PageConfiguration(
    key: 'WelcomeToSharedAccessView',
    path: WelcomeToSharedAccessViewPath,
    uiPage: Pages.WelcomeToSharedAccessView,
    currentPageAction: null);
PageConfiguration GetStartedViewPageConfig = PageConfiguration(
    key: 'GetStartedView',
    path: GetStartedViewPath,
    uiPage: Pages.GetStartedView,
    currentPageAction: null);
PageConfiguration ShareReceiptViewPageConfig = PageConfiguration(
    key: 'ShareReceiptView',
    path: ShareReceiptViewPath,
    uiPage: Pages.ShareReceiptView,
    currentPageAction: null);
PageConfiguration AnnouncementViewPageConfig = PageConfiguration(
    key: 'AnnouncementView',
    path: AnnouncementViewPath,
    uiPage: Pages.AnnouncementView,
    currentPageAction: null);
PageConfiguration WrappedAssetViewPageConfig = PageConfiguration(
    key: 'WrappedAssetView',
    path: WrappedAssetViewPath,
    uiPage: Pages.WrappedAssetView,
    currentPageAction: null);
PageConfiguration GenerateDepositAddressViewPageConfig = PageConfiguration(
    key: 'GenerateDepositAddressView',
    path: GenerateDepositAddressViewPath,
    uiPage: Pages.GenerateDepositAddressView,
    currentPageAction: null);
PageConfiguration SelectDepositAddressViewPageConfig = PageConfiguration(
    key: 'SelectDepositAddressView',
    path: SelectDepositAddressViewPath,
    uiPage: Pages.SelectDepositAddressView,
    currentPageAction: null);
PageConfiguration WithdrawAssetViewPageConfig = PageConfiguration(
    key: 'WithdrawAssetView',
    path: WithdrawAssetViewPath,
    uiPage: Pages.WithdrawAssetView,
    currentPageAction: null);
PageConfiguration ConfirmWithdrawViewPageConfig = PageConfiguration(
    key: 'ConfirmWithdrawView',
    path: ConfirmWithdrawViewPath,
    uiPage: Pages.ConfirmWithdrawView,
    currentPageAction: null);
PageConfiguration TransactionStatusViewPageConfig = PageConfiguration(
    key: 'TransactionStatusView',
    path: TransactionStatusViewPath,
    uiPage: Pages.TransactionStatusView,
    currentPageAction: null);
PageConfiguration DepositWithdrawHistoryViewPageConfig = PageConfiguration(
    key: 'DepositWithdrawHistoryView',
    path: DepositWithdrawHistoryViewPath,
    uiPage: Pages.DepositWithdrawHistoryView,
    currentPageAction: null);
PageConfiguration DepositWithdrawDetailsViewPageConfig = PageConfiguration(
    key: 'DepositWithdrawDetailsView',
    path: DepositWithdrawDetailsViewPath,
    uiPage: Pages.DepositWithdrawDetailsView,
    currentPageAction: null);
PageConfiguration WelcomeSubscriptionsViewPageConfig = PageConfiguration(
    key: 'WelcomeSubscriptionsView',
    path: WelcomeSubscriptionsViewPath,
    uiPage: Pages.WelcomeSubscriptionsView,
    currentPageAction: null);
PageConfiguration SubscriptionPlansViewPageConfig = PageConfiguration(
    key: 'SubscriptionPlansView',
    path: SubscriptionPlansViewPath,
    uiPage: Pages.SubscriptionPlansView,
    currentPageAction: null);
PageConfiguration SubscriptionPlanBenefitsViewPageConfig = PageConfiguration(
    key: 'SubscriptionPlanBenefitsView',
    path: SubscriptionPlanBenefitsViewPath,
    uiPage: Pages.SubscriptionPlanBenefitsView,
    currentPageAction: null);
PageConfiguration AuthorizeSubscriptionViewPageConfig = PageConfiguration(
    key: 'AuthorizeSubscriptionView',
    path: AuthorizeSubscriptionViewPath,
    uiPage: Pages.AuthorizeSubscriptionView,
    currentPageAction: null);
PageConfiguration OptInAssetViewPageConfig = PageConfiguration(
    key: 'OptInAssetView',
    path: OptInAssetViewPath,
    uiPage: Pages.OptInAssetView,
    currentPageAction: null);
PageConfiguration OptOutAssetViewPageConfig = PageConfiguration(
    key: 'OptOutAssetView',
    path: OptOutAssetViewPath,
    uiPage: Pages.OptOutAssetView,
    currentPageAction: null);
PageConfiguration OptInOutAssetViewPageConfig = PageConfiguration(
    key: 'OptInOutAssetView',
    path: OptInOutAssetViewPath,
    uiPage: Pages.OptInOutAssetView,
    currentPageAction: null);
PageConfiguration TokenizationWelcomeViewPageConfig = PageConfiguration(
    key: 'TokenizationWelcomeView',
    path: TokenizationWelcomeViewPath,
    uiPage: Pages.TokenizationWelcomeView,
    currentPageAction: null);
PageConfiguration SettingsViewPageConfig = PageConfiguration(
    key: 'SettingsView',
    path: SettingsViewPath,
    uiPage: Pages.SettingsView,
    currentPageAction: null);
PageConfiguration TokenizeAssetViewPageConfig = PageConfiguration(
    key: 'TokenizeAssetView',
    path: TokenizeAssetViewPath,
    uiPage: Pages.TokenizeAssetView,
    currentPageAction: null);
PageConfiguration AssetInformationViewPageConfig = PageConfiguration(
    key: 'AssetInformationView',
    path: AssetInformationViewPath,
    uiPage: Pages.AssetInformationView,
    currentPageAction: null);
PageConfiguration AssetTokenInformationViewPageConfig = PageConfiguration(
    key: 'AssetTokenInformationView',
    path: AssetTokenInformationViewPath,
    uiPage: Pages.AssetTokenInformationView,
    currentPageAction: null);
PageConfiguration AssetVerificationDocumentsViewPageConfig = PageConfiguration(
    key: 'AssetVerificationDocumentsView',
    path: AssetVerificationDocumentsViewPath,
    uiPage: Pages.AssetVerificationDocumentsView,
    currentPageAction: null);
PageConfiguration TokenizedAssetDetailViewPageConfig = PageConfiguration(
    key: 'TokenizedAssetDetailView',
    path: TokenizedAssetDetailViewPath,
    uiPage: Pages.TokenizedAssetDetailView,
    currentPageAction: null);
PageConfiguration BuyTokensViewPageConfig = PageConfiguration(
    key: 'BuyTokensView',
    path: BuyTokensViewPath,
    uiPage: Pages.BuyTokensView,
    currentPageAction: null);
PageConfiguration ConfirmBuyViewPageConfig = PageConfiguration(
    key: 'ConfirmBuyView',
    path: ConfirmBuyViewPath,
    uiPage: Pages.ConfirmBuyView,
    currentPageAction: null);
PageConfiguration TokenizedAssetsListViewPageConfig = PageConfiguration(
    key: 'TokenizedAssetsListView',
    path: TokenizedAssetsListViewPath,
    uiPage: Pages.TokenizedAssetsListView,
    currentPageAction: null);
PageConfiguration AssetDashboardViewPageConfig = PageConfiguration(
    key: 'AssetDashboardView',
    path: AssetDashboardViewPath,
    uiPage: Pages.AssetDashboardView,
    currentPageAction: null);
PageConfiguration AssetSubscribersViewPageConfig = PageConfiguration(
    key: 'AssetSubscribersView',
    path: AssetSubscribersViewPath,
    uiPage: Pages.AssetSubscribersView,
    currentPageAction: null);
PageConfiguration TotalSalesViewPageConfig = PageConfiguration(
    key: 'TotalSalesView',
    path: TotalSalesViewPath,
    uiPage: Pages.TotalSalesView,
    currentPageAction: null);
PageConfiguration ProceedsPayOutViewPageConfig = PageConfiguration(
    key: 'ProceedsPayOutView',
    path: ProceedsPayOutViewPath,
    uiPage: Pages.ProceedsPayOutView,
    currentPageAction: null);
PageConfiguration LiquidateAssetViewPageConfig = PageConfiguration(
    key: 'LiquidateAssetView',
    path: LiquidateAssetViewPath,
    uiPage: Pages.LiquidateAssetView,
    currentPageAction: null);
PageConfiguration WalletPreparationViewPageConfig = PageConfiguration(
    key: 'WalletPreparationView',
    path: WalletPreparationViewPath,
    uiPage: Pages.WalletPreparationView,
    currentPageAction: null);
PageConfiguration MyAssetTokenDetailsViewPageConfig = PageConfiguration(
    key: 'MyAssetTokenDetailsView',
    path: MyAssetTokenDetailsViewPath,
    uiPage: Pages.MyAssetTokenDetailsView,
    currentPageAction: null);
PageConfiguration MarketTradeViewPageConfig = PageConfiguration(
    key: 'MarketTradeView',
    path: MarketTradeViewPath,
    uiPage: Pages.MarketTradeView,
    currentPageAction: null);
PageConfiguration MarketTradeInfoViewPageConfig = PageConfiguration(
    key: 'MarketTradeInfoView',
    path: MarketTradeInfoViewPath,
    uiPage: Pages.MarketTradeInfoView,
    currentPageAction: null);
PageConfiguration MarketPairsViewPageConfig = PageConfiguration(
    key: 'MarketPairsView',
    path: MarketPairsViewPath,
    uiPage: Pages.MarketPairsView,
    currentPageAction: null);
PageConfiguration SetupAndComplianceViewPageConfig = PageConfiguration(
    key: 'SetupAndComplianceView',
    path: SetupAndComplianceViewPath,
    uiPage: Pages.SetupAndComplianceView,
    currentPageAction: null);
