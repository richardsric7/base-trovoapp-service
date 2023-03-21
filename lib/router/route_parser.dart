import 'package:flutter/material.dart';
import 'ui_pages.dart';

class TrovoWalletRouteParser extends RouteInformationParser<PageConfiguration> {
  @override
  Future<PageConfiguration> parseRouteInformation(
      RouteInformation routeInformation) async {
    final uri = Uri.parse(routeInformation.location!);
    if (uri.pathSegments.isEmpty) {
      return SplashPageConfig;
    }
    final path = '/' + uri.pathSegments[0];
    switch (path) {
      case SplashPath:
        return SplashPageConfig;
      case LoginPath:
        return LoginPageConfig;
      case OnboardingPath:
        return OnboardingPageConfig;
      case SignupPath:
        return SignupPageConfig;
      case CreatePasswordPath:
        return CreatePasswordPageConfig;
      case ImportWalletPath:
        return ImportWalletPageConfig;
      case EnsurePrivacyPath:
        return EnsurePrivacyPageConfig;
      case CongratulationsPath:
        return CongratulationsPageConfig;
      case FingerprintPath:
        return FingerprintPageConfig;
      case BottomHomePath:
        return BottomHomePageConfig;
      case WebViewPath:
        return WebViewPageConfig;
      case QrScannerPath:
        return QrScannerPageConfig;
      case SearchViewPath:
        return SearchViewPageConfig;
      case NotificationsViewPath:
        return NotificationsViewPageConfig;
      case CreateSubWalletSuccessViewPath:
        return CreateSubWalletSuccessViewPageConfig;
      case WalletDetailsViewPath:
        return WalletDetailsViewPageConfig;
      case AssetDetailsViewPath:
        return AssetDetailsViewPageConfig;
      case SendAssetViewPath:
        return SendAssetViewPageConfig;
      case ConfirmTransactionViewPath:
        return ConfirmTransactionViewPageConfig;
      case TransactionSuccessViewPath:
        return TransactionSuccessViewPageConfig;
      case ReceiveAssetViewPath:
        return ReceiveAssetViewPageConfig;
      case PendingAssetDetailsViewPath:
        return PendingAssetDetailsViewPageConfig;
      case PaymentHistoryViewPath:
        return PaymentHistoryViewPageConfig;
      case PaymentDetailsViewPath:
        return PaymentDetailsViewPageConfig;
      case SwapAssetsViewPath:
        return SwapAssetsViewPageConfig;
      case ConfirmSwapViewPath:
        return ConfirmSwapViewPageConfig;
      case SwapSuccessViewPath:
        return SwapSuccessViewPageConfig;
      case ProfileDetailsViewPath:
        return ProfileDetailsViewPageConfig;
      case ReferralInfoViewPath:
        return ReferralInfoViewPageConfig;
      case PasswordMgtViewPath:
        return PasswordMgtViewPageConfig;
      case BackupAllViewPath:
        return BackupAllViewPageConfig;
      case AuthorizeLoginViewPath:
        return AuthorizeLoginViewPageConfig;
      case AuthorizeActionViewPath:
        return AuthorizeActionViewPageConfig;
      case RequestSpecificPaymentViewPath:
        return RequestSpecificPaymentViewPageConfig;
      case RequestSpecificPaymentDetailsViewPath:
        return RequestSpecificPaymentDetailsViewPageConfig;
      case SecurityQuestionsViewPath:
        return SecurityQuestionsViewPageConfig;
      case AccountRecoverySuccessViewPath:
        return AccountRecoverySuccessViewPageConfig;
      case SharedAccessViewPath:
        return SharedAccessViewPageConfig;
      case SetupAccountRecoveryViewPath:
        return SetupAccountRecoveryViewPageConfig;
      case DisableAccountRecoveryViewPath:
        return DisableAccountRecoveryViewPageConfig;
      case RecoverAccountViewPath:
        return RecoverAccountViewPageConfig;
      case AnswerSecurityQuestionsViewPath:
        return AnswerSecurityQuestionsViewPageConfig;
      case RequestBackupViewPath:
        return RequestBackupViewPageConfig;
      case BackupRecoverySecretViewPath:
        return BackupRecoverySecretViewPageConfig;
      case CompleteAccountRecoveryViewPath:
        return CompleteAccountRecoveryViewPageConfig;
      case DisableAccountRecoveryInfoViewPath:
        return DisableAccountRecoveryInfoViewPageConfig;
      case SuccessViewPath:
        return SuccessViewPageConfig;
      case SecurityQuestionsForInactiveAccountsViewPath:
        return SecurityQuestionsForInactiveAccountsViewPageConfig;
      case AddSharedAccessDetailsViewPath:
        return AddSharedAccessDetailsViewPageConfig;
      case SharedWalletInfoViewPath:
        return SharedWalletInfoViewPageConfig;
      case ApprovalDetailsViewPath:
        return ApprovalDetailsViewPageConfig;
      case UpdateSharedAccessViewPath:
        return UpdateSharedAccessViewPageConfig;
      case UpdateSharedAccessDetailsViewPath:
        return UpdateSharedAccessDetailsViewPageConfig;
      case WelcomeToSharedAccessViewPath:
        return WelcomeToSharedAccessViewPageConfig;
      case GetStartedViewPath:
        return GetStartedViewPageConfig;
      case ShareReceiptViewPath:
        return ShareReceiptViewPageConfig;
      case AnnouncementViewPath:
        return AnnouncementViewPageConfig;
      case WrappedAssetViewPath:
        return WrappedAssetViewPageConfig;
      case GenerateDepositAddressViewPath:
        return GenerateDepositAddressViewPageConfig;
      case SelectDepositAddressViewPath:
        return SelectDepositAddressViewPageConfig;
      case WithdrawAssetViewPath:
        return WithdrawAssetViewPageConfig;
      case ConfirmWithdrawViewPath:
        return ConfirmWithdrawViewPageConfig;
      case TransactionStatusViewPath:
        return TransactionStatusViewPageConfig;
      case DepositWithdrawHistoryViewPath:
        return DepositWithdrawHistoryViewPageConfig;
      case DepositWithdrawDetailsViewPath:
        return DepositWithdrawDetailsViewPageConfig;
      case WelcomeSubscriptionsViewPath:
        return WelcomeSubscriptionsViewPageConfig;
      case SubscriptionPlansViewPath:
        return SubscriptionPlansViewPageConfig;
      case SubscriptionPlanBenefitsViewPath:
        return SubscriptionPlanBenefitsViewPageConfig;
      case SubscriptionPlanOptionsViewPath:
        return SubscriptionPlanOptionsViewPageConfig;
      case AuthorizeSubscriptionViewPath:
        return AuthorizeSubscriptionViewPageConfig;
      case OptInAssetViewPath:
        return OptInAssetViewPageConfig;
      case OptOutAssetViewPath:
        return OptOutAssetViewPageConfig;
      case OptInOutAssetViewPath:
        return OptInOutAssetViewPageConfig;
      case TokenizationWelcomeViewPath:
        return TokenizationWelcomeViewPageConfig;
      case SettingsViewPath:
        return SettingsViewPageConfig;
      case TokenizeAssetViewPath:
        return TokenizeAssetViewPageConfig;
      case AssetVerificationDocumentsViewPath:
        return AssetVerificationDocumentsViewPageConfig;
      case TokenizedAssetDetailViewPath:
        return TokenizedAssetDetailViewPageConfig;
      case BuyTokensViewPath:
        return BuyTokensViewPageConfig;
      case ConfirmBuyViewPath:
        return ConfirmBuyViewPageConfig;
      case TokenizedAssetsListViewPath:
        return TokenizedAssetsListViewPageConfig;
      case AssetDashboardViewPath:
        return AssetDashboardViewPageConfig;
      case AssetSubscribersViewPath:
        return AssetSubscribersViewPageConfig;
      case TotalSalesViewPath:
        return TotalSalesViewPageConfig;
      case ProceedsPayOutViewPath:
        return ProceedsPayOutViewPageConfig;
      case LiquidateAssetViewPath:
        return LiquidateAssetViewPageConfig;
      case LiquidateAssetViewPath:
        return LiquidateAssetViewPageConfig;
      case WalletPreparationViewPath:
        return WalletPreparationViewPageConfig;
      case MyAssetTokenDetailsViewPath:
        return MyAssetTokenDetailsViewPageConfig;
      default:
        return SplashPageConfig;
    }
  }

  @override
  RouteInformation restoreRouteInformation(PageConfiguration configuration) {
    switch (configuration.uiPage) {
      case Pages.Splash:
        return const RouteInformation(location: SplashPath);
      case Pages.Login:
        return const RouteInformation(location: LoginPath);
      case Pages.Onboarding:
        return const RouteInformation(location: OnboardingPath);
      case Pages.Signup:
        return const RouteInformation(location: SignupPath);
      case Pages.CreatePassword:
        return const RouteInformation(location: CreatePasswordPath);
      case Pages.ImportWallet:
        return const RouteInformation(location: ImportWalletPath);
      case Pages.Verification:
        return const RouteInformation(location: VerificationPath);
      case Pages.EnsurePrivacy:
        return const RouteInformation(location: EnsurePrivacyPath);
      case Pages.Congratulations:
        return const RouteInformation(location: CongratulationsPath);
      case Pages.Fingerprint:
        return const RouteInformation(location: FingerprintPath);
      case Pages.BottomHome:
        return const RouteInformation(location: BottomHomePath);
      case Pages.WebView:
        return const RouteInformation(location: WebViewPath);
      case Pages.QrScanner:
        return const RouteInformation(location: QrScannerPath);
      case Pages.SearchView:
        return const RouteInformation(location: SearchViewPath);
      case Pages.NotificationsView:
        return const RouteInformation(location: NotificationsViewPath);
      case Pages.CreateSubWalletSuccessView:
        return const RouteInformation(location: CreateSubWalletSuccessViewPath);
      case Pages.WalletDetailsView:
        return const RouteInformation(location: WalletDetailsViewPath);
      case Pages.AssetDetailsView:
        return const RouteInformation(location: AssetDetailsViewPath);
      case Pages.SendAssetView:
        return const RouteInformation(location: SendAssetViewPath);
      case Pages.ConfirmTransactionView:
        return const RouteInformation(location: ConfirmTransactionViewPath);
      case Pages.TransactionSuccessView:
        return const RouteInformation(location: TransactionSuccessViewPath);
      case Pages.ReceiveAssetView:
        return const RouteInformation(location: ReceiveAssetViewPath);
      case Pages.PendingAssetDetailsView:
        return const RouteInformation(location: PendingAssetDetailsViewPath);
      case Pages.PaymentHistoryView:
        return const RouteInformation(location: PaymentHistoryViewPath);
      case Pages.PaymentDetailsView:
        return const RouteInformation(location: PaymentDetailsViewPath);
      case Pages.SwapAssetsView:
        return const RouteInformation(location: SwapAssetsViewPath);
      case Pages.ConfirmSwapView:
        return const RouteInformation(location: ConfirmSwapViewPath);
      case Pages.SwapSuccessView:
        return const RouteInformation(location: SwapSuccessViewPath);
      case Pages.ProfileDetailsView:
        return const RouteInformation(location: ProfileDetailsViewPath);
      case Pages.ReferralInfoView:
        return const RouteInformation(location: ReferralInfoViewPath);
      case Pages.PasswordMgtView:
        return const RouteInformation(location: PasswordMgtViewPath);
      case Pages.BackupAllView:
        return const RouteInformation(location: BackupAllViewPath);
      case Pages.AuthorizeLoginView:
        return const RouteInformation(location: AuthorizeLoginViewPath);
      case Pages.AuthorizeActionView:
        return const RouteInformation(location: AuthorizeActionViewPath);
      case Pages.RequestSpecificPaymentView:
        return const RouteInformation(location: RequestSpecificPaymentViewPath);
      case Pages.RequestSpecificPaymentDetailsView:
        return const RouteInformation(
            location: RequestSpecificPaymentDetailsViewPath);
      case Pages.SecurityQuestionsView:
        return const RouteInformation(location: SecurityQuestionsViewPath);
      case Pages.RequestOtpView:
        return const RouteInformation(location: RequestOtpViewPath);
      case Pages.AccountRecoverySuccessView:
        return const RouteInformation(location: AccountRecoverySuccessViewPath);
      case Pages.SharedAccessView:
        return const RouteInformation(location: SharedAccessViewPath);
      case Pages.SetupAccountRecoveryView:
        return const RouteInformation(location: SetupAccountRecoveryViewPath);
      case Pages.DisableAccountRecoveryView:
        return const RouteInformation(location: DisableAccountRecoveryViewPath);
      case Pages.RecoverAccountView:
        return const RouteInformation(location: RecoverAccountViewPath);
      case Pages.AnswerSecurityQuestionsView:
        return const RouteInformation(
            location: AnswerSecurityQuestionsViewPath);
      case Pages.RequestBackupView:
        return const RouteInformation(location: RequestBackupViewPath);
      case Pages.BackupRecoverySecretView:
        return const RouteInformation(location: BackupRecoverySecretViewPath);
      case Pages.CompleteAccountRecoveryView:
        return const RouteInformation(
            location: CompleteAccountRecoveryViewPath);
      case Pages.DisableAccountRecoveryInfoView:
        return const RouteInformation(
            location: DisableAccountRecoveryInfoViewPath);
      case Pages.SuccessView:
        return const RouteInformation(location: SuccessViewPath);
      case Pages.SecurityQuestionsForInactiveAccountsView:
        return const RouteInformation(
            location: SecurityQuestionsForInactiveAccountsViewPath);
      case Pages.AddSharedAccessDetailsView:
        return const RouteInformation(location: AddSharedAccessDetailsViewPath);
      case Pages.SharedWalletInfoView:
        return const RouteInformation(location: SharedWalletInfoViewPath);
      case Pages.ApprovalDetailsView:
        return const RouteInformation(location: ApprovalDetailsViewPath);
      case Pages.UpdateSharedAccessView:
        return const RouteInformation(location: UpdateSharedAccessViewPath);
      case Pages.UpdateSharedAccessDetailsView:
        return const RouteInformation(
            location: UpdateSharedAccessDetailsViewPath);
      case Pages.WelcomeToSharedAccessView:
        return const RouteInformation(location: WelcomeToSharedAccessViewPath);
      case Pages.GetStartedView:
        return const RouteInformation(location: GetStartedViewPath);
      case Pages.ShareReceiptView:
        return const RouteInformation(location: ShareReceiptViewPath);
      case Pages.AnnouncementView:
        return const RouteInformation(location: AnnouncementViewPath);
      case Pages.WrappedAssetView:
        return const RouteInformation(location: WrappedAssetViewPath);
      case Pages.GenerateDepositAddressView:
        return const RouteInformation(location: GenerateDepositAddressViewPath);
      case Pages.SelectDepositAddressView:
        return const RouteInformation(location: SelectDepositAddressViewPath);
      case Pages.WithdrawAssetView:
        return const RouteInformation(location: WithdrawAssetViewPath);
      case Pages.ConfirmWithdrawView:
        return const RouteInformation(location: ConfirmWithdrawViewPath);
      case Pages.TransactionStatusView:
        return const RouteInformation(location: TransactionStatusViewPath);
      case Pages.DepositWithdrawHistoryView:
        return const RouteInformation(location: DepositWithdrawHistoryViewPath);
      case Pages.DepositWithdrawDetailsView:
        return const RouteInformation(location: DepositWithdrawDetailsViewPath);
      case Pages.WelcomeSubscriptionsView:
        return const RouteInformation(location: WelcomeSubscriptionsViewPath);
      case Pages.SubscriptionPlansView:
        return const RouteInformation(location: SubscriptionPlansViewPath);
      case Pages.SubscriptionPlanBenefitsView:
        return const RouteInformation(
            location: SubscriptionPlanBenefitsViewPath);
      case Pages.SubscriptionPlanOptionsView:
        return const RouteInformation(
            location: SubscriptionPlanOptionsViewPath);
      case Pages.AuthorizeSubscriptionView:
        return const RouteInformation(location: AuthorizeSubscriptionViewPath);
      case Pages.OptInAssetView:
        return const RouteInformation(location: OptInAssetViewPath);
      case Pages.OptOutAssetView:
        return const RouteInformation(location: OptOutAssetViewPath);
      case Pages.OptInOutAssetView:
        return const RouteInformation(location: OptInOutAssetViewPath);
      case Pages.TokenizationWelcomeView:
        return const RouteInformation(location: TokenizationWelcomeViewPath);
      case Pages.SettingsView:
        return const RouteInformation(location: SettingsViewPath);
      case Pages.TokenizeAssetView:
        return const RouteInformation(location: TokenizeAssetViewPath);
      case Pages.AssetInformationView:
        return const RouteInformation(location: AssetInformationViewPath);
      case Pages.AssetVerificationDocumentsView:
        return const RouteInformation(
            location: AssetVerificationDocumentsViewPath);
      case Pages.TokenizedAssetDetailView:
        return const RouteInformation(location: TokenizedAssetDetailViewPath);
      case Pages.BuyTokensView:
        return const RouteInformation(location: BuyTokensViewPath);
      case Pages.ConfirmBuyView:
        return const RouteInformation(location: ConfirmBuyViewPath);
      case Pages.TokenizedAssetsListView:
        return const RouteInformation(location: TokenizedAssetsListViewPath);
      case Pages.AssetDashboardView:
        return const RouteInformation(location: AssetDashboardViewPath);
      case Pages.AssetSubscribersView:
        return const RouteInformation(location: AssetSubscribersViewPath);
      case Pages.TotalSalesView:
        return const RouteInformation(location: TotalSalesViewPath);
      case Pages.ProceedsPayOutView:
        return const RouteInformation(location: ProceedsPayOutViewPath);
      case Pages.LiquidateAssetView:
        return const RouteInformation(location: LiquidateAssetViewPath);
      case Pages.WalletPreparationView:
        return const RouteInformation(location: WalletPreparationViewPath);
      case Pages.MyAssetTokenDetailsView:
        return const RouteInformation(location: MyAssetTokenDetailsViewPath);
      default:
        return const RouteInformation(location: SplashPath);
    }
  }
}
