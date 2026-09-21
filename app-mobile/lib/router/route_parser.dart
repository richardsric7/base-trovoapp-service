import 'package:flutter/material.dart';

import 'ui_pages.dart';

class TrovoWalletRouteParser extends RouteInformationParser<PageConfiguration> {
  @override
  Future<PageConfiguration> parseRouteInformation(
    RouteInformation routeInformation,
  ) async {
    final uri = routeInformation.uri;
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
      case PdfViewPath:
        return PdfViewPageConfig;
      case QrScannerPath:
        return QrScannerPageConfig;
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
      case AuthorizeSubscriptionViewPath:
        return AuthorizeSubscriptionViewPageConfig;
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
      case WalletPreparationViewPath:
        return WalletPreparationViewPageConfig;
      case MarketTradeViewPath:
        return MarketTradeViewPageConfig;
      case MarketTradeInfoViewPath:
        return MarketTradeInfoViewPageConfig;
      case MarketPairsViewPath:
        return MarketPairsViewPageConfig;
      case SetupAndComplianceViewPath:
        return SetupAndComplianceViewPageConfig;
      case AllWalletsViewPath:
        return AllWalletsViewPageConfig;
      case ConfirmTokenizationDetailsViewPath:
        return ConfirmTokenizationDetailsViewPageConfig;
      case TokenizationFeePaymentViewPath:
        return TokenizationFeePaymentViewPageConfig;
      case DeleteAccountViewPath:
        return DeleteAccountViewPageConfig;
      default:
        return SplashPageConfig;
    }
  }

  @override
  RouteInformation restoreRouteInformation(PageConfiguration configuration) {
    switch (configuration.uiPage) {
      case Pages.Splash:
        return RouteInformation(uri: Uri.parse(SplashPath));
      case Pages.Login:
        return RouteInformation(uri: Uri.parse(LoginPath));
      case Pages.Onboarding:
        return RouteInformation(uri: Uri.parse(OnboardingPath));
      case Pages.Signup:
        return RouteInformation(uri: Uri.parse(SignupPath));
      case Pages.CreatePassword:
        return RouteInformation(uri: Uri.parse(CreatePasswordPath));
      case Pages.ImportWallet:
        return RouteInformation(uri: Uri.parse(ImportWalletPath));
      case Pages.Verification:
        return RouteInformation(uri: Uri.parse(VerificationPath));
      case Pages.EnsurePrivacy:
        return RouteInformation(uri: Uri.parse(EnsurePrivacyPath));
      case Pages.Congratulations:
        return RouteInformation(uri: Uri.parse(CongratulationsPath));
      case Pages.Fingerprint:
        return RouteInformation(uri: Uri.parse(FingerprintPath));
      case Pages.BottomHome:
        return RouteInformation(uri: Uri.parse(BottomHomePath));
      case Pages.WebView:
        return RouteInformation(uri: Uri.parse(WebViewPath));
      case Pages.PdfView:
        return RouteInformation(uri: Uri.parse(PdfViewPath));
      case Pages.QrScanner:
        return RouteInformation(uri: Uri.parse(QrScannerPath));
      case Pages.NotificationsView:
        return RouteInformation(uri: Uri.parse(NotificationsViewPath));
      case Pages.CreateSubWalletSuccessView:
        return RouteInformation(uri: Uri.parse(CreateSubWalletSuccessViewPath));
      case Pages.WalletDetailsView:
        return RouteInformation(uri: Uri.parse(WalletDetailsViewPath));
      case Pages.AssetDetailsView:
        return RouteInformation(uri: Uri.parse(AssetDetailsViewPath));
      case Pages.SendAssetView:
        return RouteInformation(uri: Uri.parse(SendAssetViewPath));
      case Pages.ConfirmTransactionView:
        return RouteInformation(uri: Uri.parse(ConfirmTransactionViewPath));
      case Pages.TransactionSuccessView:
        return RouteInformation(uri: Uri.parse(TransactionSuccessViewPath));
      case Pages.ReceiveAssetView:
        return RouteInformation(uri: Uri.parse(ReceiveAssetViewPath));
      case Pages.PaymentHistoryView:
        return RouteInformation(uri: Uri.parse(PaymentHistoryViewPath));
      case Pages.PaymentDetailsView:
        return RouteInformation(uri: Uri.parse(PaymentDetailsViewPath));
      case Pages.SwapAssetsView:
        return RouteInformation(uri: Uri.parse(SwapAssetsViewPath));
      case Pages.ConfirmSwapView:
        return RouteInformation(uri: Uri.parse(ConfirmSwapViewPath));
      case Pages.SwapSuccessView:
        return RouteInformation(uri: Uri.parse(SwapSuccessViewPath));
      case Pages.ProfileDetailsView:
        return RouteInformation(uri: Uri.parse(ProfileDetailsViewPath));
      case Pages.ReferralInfoView:
        return RouteInformation(uri: Uri.parse(ReferralInfoViewPath));
      case Pages.PasswordMgtView:
        return RouteInformation(uri: Uri.parse(PasswordMgtViewPath));
      case Pages.BackupAllView:
        return RouteInformation(uri: Uri.parse(BackupAllViewPath));
      case Pages.AuthorizeLoginView:
        return RouteInformation(uri: Uri.parse(AuthorizeLoginViewPath));
      case Pages.AuthorizeActionView:
        return RouteInformation(uri: Uri.parse(AuthorizeActionViewPath));
      case Pages.RequestSpecificPaymentView:
        return RouteInformation(uri: Uri.parse(RequestSpecificPaymentViewPath));
      case Pages.RequestSpecificPaymentDetailsView:
        return RouteInformation(
          uri: Uri.parse(RequestSpecificPaymentDetailsViewPath),
        );
      case Pages.SecurityQuestionsView:
        return RouteInformation(uri: Uri.parse(SecurityQuestionsViewPath));
      case Pages.RequestOtpView:
        return RouteInformation(uri: Uri.parse(RequestOtpViewPath));
      case Pages.AccountRecoverySuccessView:
        return RouteInformation(uri: Uri.parse(AccountRecoverySuccessViewPath));
      case Pages.SharedAccessView:
        return RouteInformation(uri: Uri.parse(SharedAccessViewPath));
      case Pages.SetupAccountRecoveryView:
        return RouteInformation(uri: Uri.parse(SetupAccountRecoveryViewPath));
      case Pages.DisableAccountRecoveryView:
        return RouteInformation(uri: Uri.parse(DisableAccountRecoveryViewPath));
      case Pages.RecoverAccountView:
        return RouteInformation(uri: Uri.parse(RecoverAccountViewPath));
      case Pages.AnswerSecurityQuestionsView:
        return RouteInformation(
          uri: Uri.parse(AnswerSecurityQuestionsViewPath),
        );
      case Pages.RequestBackupView:
        return RouteInformation(uri: Uri.parse(RequestBackupViewPath));
      case Pages.BackupRecoverySecretView:
        return RouteInformation(uri: Uri.parse(BackupRecoverySecretViewPath));
      case Pages.CompleteAccountRecoveryView:
        return RouteInformation(
          uri: Uri.parse(CompleteAccountRecoveryViewPath),
        );
      case Pages.DisableAccountRecoveryInfoView:
        return RouteInformation(
          uri: Uri.parse(DisableAccountRecoveryInfoViewPath),
        );
      case Pages.SuccessView:
        return RouteInformation(uri: Uri.parse(SuccessViewPath));
      case Pages.SecurityQuestionsForInactiveAccountsView:
        return RouteInformation(
          uri: Uri.parse(SecurityQuestionsForInactiveAccountsViewPath),
        );
      case Pages.AddSharedAccessDetailsView:
        return RouteInformation(uri: Uri.parse(AddSharedAccessDetailsViewPath));
      case Pages.SharedWalletInfoView:
        return RouteInformation(uri: Uri.parse(SharedWalletInfoViewPath));
      case Pages.ApprovalDetailsView:
        return RouteInformation(uri: Uri.parse(ApprovalDetailsViewPath));
      case Pages.UpdateSharedAccessView:
        return RouteInformation(uri: Uri.parse(UpdateSharedAccessViewPath));
      case Pages.UpdateSharedAccessDetailsView:
        return RouteInformation(
          uri: Uri.parse(UpdateSharedAccessDetailsViewPath),
        );
      case Pages.WelcomeToSharedAccessView:
        return RouteInformation(uri: Uri.parse(WelcomeToSharedAccessViewPath));
      case Pages.GetStartedView:
        return RouteInformation(uri: Uri.parse(GetStartedViewPath));
      case Pages.ShareReceiptView:
        return RouteInformation(uri: Uri.parse(ShareReceiptViewPath));
      case Pages.AnnouncementView:
        return RouteInformation(uri: Uri.parse(AnnouncementViewPath));
      case Pages.WrappedAssetView:
        return RouteInformation(uri: Uri.parse(WrappedAssetViewPath));
      case Pages.GenerateDepositAddressView:
        return RouteInformation(uri: Uri.parse(GenerateDepositAddressViewPath));
      case Pages.SelectDepositAddressView:
        return RouteInformation(uri: Uri.parse(SelectDepositAddressViewPath));
      case Pages.WithdrawAssetView:
        return RouteInformation(uri: Uri.parse(WithdrawAssetViewPath));
      case Pages.ConfirmWithdrawView:
        return RouteInformation(uri: Uri.parse(ConfirmWithdrawViewPath));
      case Pages.TransactionStatusView:
        return RouteInformation(uri: Uri.parse(TransactionStatusViewPath));
      case Pages.DepositWithdrawHistoryView:
        return RouteInformation(uri: Uri.parse(DepositWithdrawHistoryViewPath));
      case Pages.DepositWithdrawDetailsView:
        return RouteInformation(uri: Uri.parse(DepositWithdrawDetailsViewPath));
      case Pages.WelcomeSubscriptionsView:
        return RouteInformation(uri: Uri.parse(WelcomeSubscriptionsViewPath));
      case Pages.SubscriptionPlansView:
        return RouteInformation(uri: Uri.parse(SubscriptionPlansViewPath));
      case Pages.SubscriptionPlanBenefitsView:
        return RouteInformation(
          uri: Uri.parse(SubscriptionPlanBenefitsViewPath),
        );
      case Pages.AuthorizeSubscriptionView:
        return RouteInformation(uri: Uri.parse(AuthorizeSubscriptionViewPath));
      case Pages.TokenizationWelcomeView:
        return RouteInformation(uri: Uri.parse(TokenizationWelcomeViewPath));
      case Pages.SettingsView:
        return RouteInformation(uri: Uri.parse(SettingsViewPath));
      case Pages.TokenizeAssetView:
        return RouteInformation(uri: Uri.parse(TokenizeAssetViewPath));
      case Pages.AssetInformationView:
        return RouteInformation(uri: Uri.parse(AssetInformationViewPath));
      case Pages.AssetVerificationDocumentsView:
        return RouteInformation(
          uri: Uri.parse(AssetVerificationDocumentsViewPath),
        );
      case Pages.TokenizedAssetDetailView:
        return RouteInformation(uri: Uri.parse(TokenizedAssetDetailViewPath));
      case Pages.BuyTokensView:
        return RouteInformation(uri: Uri.parse(BuyTokensViewPath));
      case Pages.ConfirmBuyView:
        return RouteInformation(uri: Uri.parse(ConfirmBuyViewPath));
      case Pages.AssetDashboardView:
        return RouteInformation(uri: Uri.parse(AssetDashboardViewPath));
      case Pages.AssetSubscribersView:
        return RouteInformation(uri: Uri.parse(AssetSubscribersViewPath));
      case Pages.TotalSalesView:
        return RouteInformation(uri: Uri.parse(TotalSalesViewPath));
      case Pages.ProceedsPayOutView:
        return RouteInformation(uri: Uri.parse(ProceedsPayOutViewPath));
      case Pages.LiquidateAssetView:
        return RouteInformation(uri: Uri.parse(LiquidateAssetViewPath));
      case Pages.WalletPreparationView:
        return RouteInformation(uri: Uri.parse(WalletPreparationViewPath));
      case Pages.MyAssetTokenDetailsView:
        return RouteInformation(uri: Uri.parse(MyAssetTokenDetailsViewPath));
      case Pages.MarketTradeView:
        return RouteInformation(uri: Uri.parse(MarketTradeViewPath));
      case Pages.MarketTradeInfoView:
        return RouteInformation(uri: Uri.parse(MarketTradeInfoViewPath));
      case Pages.MarketPairsView:
        return RouteInformation(uri: Uri.parse(MarketPairsViewPath));
      case Pages.SetupAndComplianceView:
        return RouteInformation(uri: Uri.parse(SetupAndComplianceViewPath));
      case Pages.AllWalletsView:
        return RouteInformation(uri: Uri.parse(AllWalletsViewPath));
      case Pages.ConfirmTokenizationDetailsView:
        return RouteInformation(
          uri: Uri.parse(ConfirmTokenizationDetailsViewPath),
        );
      case Pages.TokenizationFeePaymentView:
        return RouteInformation(uri: Uri.parse(TokenizationFeePaymentViewPath));
      case Pages.DeleteAccountPrerequisitesView:
        return RouteInformation(
          uri: Uri.parse(DeleteAccountPrerequisitesViewPath),
        );
      case Pages.DeleteAccountView:
        return RouteInformation(uri: Uri.parse(DeleteAccountViewPath));
      default:
        return RouteInformation(uri: Uri.parse(SplashPath));
    }
  }
}
