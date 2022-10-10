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
      case RequestOtpViewPath:
        return RequestOtpViewPageConfig;
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
      default:
        return const RouteInformation(location: SplashPath);
    }
  }
}
