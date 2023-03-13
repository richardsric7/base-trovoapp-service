import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/searchview.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/settings.dart';
import 'package:trovo_wallet/custom_bloc_observer/swiper/swiper.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/confirm_swap.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/change_password.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_detail.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/profile_details.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/referral_info.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/share_receipt.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/subwallet_create_success.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/swap_assets.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/swap_success.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallet_details.dart';
import 'package:trovo_wallet/bottom_bar/bottom_bar.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/deposit_withdraw_details.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/opt_in_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/opt_in_out_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/opt_out_asset.dart';
import 'package:trovo_wallet/screens/account_recovery/account_recovery_success.dart';
import 'package:trovo_wallet/screens/account_recovery/answer_security_questions.dart';
import 'package:trovo_wallet/screens/account_recovery/backup_recovery_secret.dart';
import 'package:trovo_wallet/screens/account_recovery/complete_account_recovery.dart';
import 'package:trovo_wallet/screens/account_recovery/disable_account_recovery.dart';
import 'package:trovo_wallet/screens/account_recovery/disable_account_recovery_info.dart';
import 'package:trovo_wallet/screens/account_recovery/recover_account.dart';
import 'package:trovo_wallet/screens/account_recovery/request_backup.dart';
import 'package:trovo_wallet/screens/account_recovery/security_questions_for_inactive_accounts.dart';
import 'package:trovo_wallet/screens/account_recovery/setup_account_recovery.dart';
import 'package:trovo_wallet/screens/Auth/AuthorizeActionView.dart';
import 'package:trovo_wallet/screens/Auth/AuthorizeLoginView.dart';
import 'package:trovo_wallet/screens/Auth/create_password.dart';
import 'package:trovo_wallet/screens/Auth/get_started.dart';
import 'package:trovo_wallet/screens/Auth/signup.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';
import 'package:trovo_wallet/screens/Backup/backup_all.dart';
import 'package:trovo_wallet/screens/Backup/congratulation.dart';
import 'package:trovo_wallet/screens/Backup/ensure_privacy.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/asset_details.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/confirm_transaction.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/confirm_withdraw.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/deposit_withdrawal_history.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/select_deposit_address.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/generate_deposit_address.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/recieve_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/request_specific_payment_details.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/request_specific_payment.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/send_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/transaction_status.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/transaction_success.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/trust_asset.dart';
import 'package:trovo_wallet/screens/account_recovery/security_questions.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/withdraw_asset.dart';
import 'package:trovo_wallet/screens/Send_and_Recieve/wrapped_asset.dart';
import 'package:trovo_wallet/screens/announcements/announcementsView.dart';
import 'package:trovo_wallet/screens/asset-tokenization/asset_information.dart';
import 'package:trovo_wallet/screens/asset-tokenization/asset_token_information.dart';
import 'package:trovo_wallet/screens/asset-tokenization/asset_verification_documents.dart';
import 'package:trovo_wallet/screens/asset-tokenization/tokenization.dart';
import 'package:trovo_wallet/screens/asset-tokenization/tokenize_asset_view.dart';
import 'package:trovo_wallet/screens/import_wallet/import_wallet.dart';
import 'package:trovo_wallet/screens/page_view/success_view.dart';
import 'package:trovo_wallet/screens/page_view/web_view.dart';
import 'package:trovo_wallet/screens/shared_access/add_shared_access_details.dart';
import 'package:trovo_wallet/screens/shared_access/approval_details.dart';
import 'package:trovo_wallet/screens/shared_access/shared_access.dart';
import 'package:trovo_wallet/screens/shared_access/shared_wallet_info.dart';
import 'package:trovo_wallet/screens/shared_access/update_shared_access.dart';
import 'package:trovo_wallet/screens/shared_access/update_shared_access_details.dart';
import 'package:trovo_wallet/screens/shared_access/welcome_to_shared_access.dart';
import 'package:trovo_wallet/screens/subscriptions/authorize_subscription.dart';
import 'package:trovo_wallet/screens/subscriptions/subscription_benefits.dart';
import 'package:trovo_wallet/screens/subscriptions/subscription_plan_options.dart';
import 'package:trovo_wallet/screens/subscriptions/subscription_plans.dart';
import 'package:trovo_wallet/screens/subscriptions/welcome.dart';
import 'package:trovo_wallet/storage/state.dart';
import '../screens/Auth/fingerprint.dart';
import '../screens/Auth/login.dart';
import '../screens/Backup/backup.dart';
import '../screens/Splash_Screen/splash_screen.dart';
import '../screens/announcements/announcementView.dart';
import '../screens/qr_scanner_view.dart';
import 'page_actions.dart';
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
          _addPageData(AnnouncementsView(), NotificationsViewPageConfig);
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
        case Pages.AccountRecoverySuccessView:
          _addPageData(
              AccountRecoverySuccess(), AccountRecoverySuccessViewPageConfig);
          break;
        case Pages.SharedAccessView:
          _addPageData(SharedAccess(), SharedAccessViewPageConfig);
          break;
        case Pages.SetupAccountRecoveryView:
          _addPageData(
              SetupAccountRecovery(), SetupAccountRecoveryViewPageConfig);
          break;
        case Pages.DisableAccountRecoveryView:
          _addPageData(
              DisableAccountRecovery(), DisableAccountRecoveryViewPageConfig);
          break;
        case Pages.RecoverAccountView:
          _addPageData(RecoverAccount(), RecoverAccountViewPageConfig);
          break;
        case Pages.AnswerSecurityQuestionsView:
          _addPageData(
              AnswerSecurityQuestions(), AnswerSecurityQuestionsViewPageConfig);
          break;
        case Pages.RequestBackupView:
          _addPageData(RequestBackup(), RequestBackupViewPageConfig);
          break;
        case Pages.BackupRecoverySecretView:
          _addPageData(
              BackupRecoverySecret(), BackupRecoverySecretViewPageConfig);
          break;
        case Pages.CompleteAccountRecoveryView:
          _addPageData(
              CompleteAccountRecovery(), CompleteAccountRecoveryViewPageConfig);
          break;
        case Pages.DisableAccountRecoveryInfoView:
          _addPageData(DisableAccountRecoveryInfo(),
              DisableAccountRecoveryInfoViewPageConfig);
          break;
        case Pages.SuccessView:
          _addPageData(SuccessView(), SuccessViewPageConfig);
          break;
        case Pages.SecurityQuestionsForInactiveAccountsView:
          _addPageData(SecurityQuestionsForInactiveAccounts(),
              SecurityQuestionsForInactiveAccountsViewPageConfig);
          break;
        case Pages.AddSharedAccessDetailsView:
          _addPageData(
              AddSharedAccessDetails(), AddSharedAccessDetailsViewPageConfig);
          break;
        case Pages.SharedWalletInfoView:
          _addPageData(SharedWalletInfo(), SharedWalletInfoViewPageConfig);
          break;
        case Pages.ApprovalDetailsView:
          _addPageData(ApprovalDetails(), ApprovalDetailsViewPageConfig);
          break;
        case Pages.UpdateSharedAccessView:
          _addPageData(UpdateSharedAccess(), UpdateSharedAccessViewPageConfig);
          break;
        case Pages.UpdateSharedAccessDetailsView:
          _addPageData(UpdateSharedAccessDetails(),
              UpdateSharedAccessDetailsViewPageConfig);
          break;
        case Pages.WelcomeToSharedAccessView:
          _addPageData(
              WelcomeToSharedAccess(), WelcomeToSharedAccessViewPageConfig);
          break;
        case Pages.GetStartedView:
          _addPageData(GetStarted(), GetStartedViewPageConfig);
          break;
        case Pages.ShareReceiptView:
          _addPageData(ShareReceipt(), ShareReceiptViewPageConfig);
          break;
        case Pages.AnnouncementView:
          _addPageData(AnnouncementView(), AnnouncementViewPageConfig);
          break;
        case Pages.WrappedAssetView:
          _addPageData(WrappedAsset(), WrappedAssetViewPageConfig);
          break;
        case Pages.GenerateDepositAddressView:
          _addPageData(
              GenerateDepositAddress(), GenerateDepositAddressViewPageConfig);
          break;
        case Pages.SelectDepositAddressView:
          _addPageData(
              SelectDepositAddress(), SelectDepositAddressViewPageConfig);
          break;
        case Pages.WithdrawAssetView:
          _addPageData(WithdrawAsset(), WithdrawAssetViewPageConfig);
          break;
        case Pages.ConfirmWithdrawView:
          _addPageData(ConfirmWithdrawal(), ConfirmWithdrawViewPageConfig);
          break;
        case Pages.TransactionStatusView:
          _addPageData(TransactionStatus(), TransactionStatusViewPageConfig);
          break;
        case Pages.DepositWithdrawHistoryView:
          _addPageData(
              DepositWithdrawHistory(), DepositWithdrawHistoryViewPageConfig);
          break;
        case Pages.DepositWithdrawDetailsView:
          _addPageData(
              DepositWithdrawDetails(), DepositWithdrawDetailsViewPageConfig);
          break;
        case Pages.WelcomeSubscriptionsView:
          _addPageData(
              WelcomeSubscriptions(), WelcomeSubscriptionsViewPageConfig);
          break;
        case Pages.SubscriptionPlansView:
          _addPageData(SubscriptionPlans(), SubscriptionPlansViewPageConfig);
          break;
        case Pages.SubscriptionPlanBenefitsView:
          _addPageData(SubscriptionPlanBenefits(),
              SubscriptionPlanBenefitsViewPageConfig);
          break;
        case Pages.SubscriptionPlanOptionsView:
          _addPageData(
              SubscriptionPlanOptions(), SubscriptionPlanOptionsViewPageConfig);
          break;
        case Pages.AuthorizeSubscriptionView:
          _addPageData(
              AuthorizeSubscription(), AuthorizeSubscriptionViewPageConfig);
          break;
        case Pages.OptInAssetView:
          _addPageData(OptInAsset(), OptInAssetViewPageConfig);
          break;
        case Pages.OptOutAssetView:
          _addPageData(OptOutAsset(), OptOutAssetViewPageConfig);
          break;
        case Pages.OptInOutAssetView:
          _addPageData(OptInOutAsset(), OptInOutAssetViewPageConfig);
          break;
        case Pages.TokenizationWelcomeView:
          _addPageData(
              TokenizationWelcome(), TokenizationWelcomeViewPageConfig);
          break;
        case Pages.SettingsView:
          _addPageData(Settings(), SettingsViewPageConfig);
          break;
        case Pages.TokenizeAssetView:
          _addPageData(TokenizeAsset(), TokenizeAssetViewPageConfig);
          break;
        case Pages.AssetInformationView:
          _addPageData(AssetInformation(), AssetInformationViewPageConfig);
          break;
        case Pages.AssetTokenInformationView:
          _addPageData(
              AssetTokenInformation(), AssetTokenInformationViewPageConfig);
          break;
        case Pages.AssetVerificationDocumentsView:
          _addPageData(AssetVerificationDocuments(),
              AssetVerificationDocumentsViewPageConfig);
          break;
        default:
          break;
      }
    }
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
      case Pages.AccountRecoverySuccessView:
        AccountRecoverySuccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.SharedAccessView:
        SharedAccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.SetupAccountRecoveryView:
        SetupAccountRecoveryViewPageConfig.currentPageAction = action;
        break;
      case Pages.DisableAccountRecoveryView:
        DisableAccountRecoveryViewPageConfig.currentPageAction = action;
        break;
      case Pages.RecoverAccountView:
        RecoverAccountViewPageConfig.currentPageAction = action;
        break;
      case Pages.AnswerSecurityQuestionsView:
        AnswerSecurityQuestionsViewPageConfig.currentPageAction = action;
        break;
      case Pages.RequestBackupView:
        RequestBackupViewPageConfig.currentPageAction = action;
        break;
      case Pages.BackupRecoverySecretView:
        BackupRecoverySecretViewPageConfig.currentPageAction = action;
        break;
      case Pages.CompleteAccountRecoveryView:
        CompleteAccountRecoveryViewPageConfig.currentPageAction = action;
        break;
      case Pages.DisableAccountRecoveryInfoView:
        DisableAccountRecoveryInfoViewPageConfig.currentPageAction = action;
        break;
      case Pages.SuccessView:
        SuccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.AddSharedAccessDetailsView:
        AddSharedAccessDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.SecurityQuestionsForInactiveAccountsView:
        SecurityQuestionsForInactiveAccountsViewPageConfig.currentPageAction =
            action;
        break;
      case Pages.SharedWalletInfoView:
        SharedWalletInfoViewPageConfig.currentPageAction = action;
        break;
      case Pages.ApprovalDetailsView:
        ApprovalDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.UpdateSharedAccessView:
        UpdateSharedAccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.UpdateSharedAccessDetailsView:
        UpdateSharedAccessDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.WelcomeToSharedAccessView:
        WelcomeToSharedAccessViewPageConfig.currentPageAction = action;
        break;
      case Pages.GetStartedView:
        GetStartedViewPageConfig.currentPageAction = action;
        break;
      case Pages.ShareReceiptView:
        ShareReceiptViewPageConfig.currentPageAction = action;
        break;
      case Pages.AnnouncementView:
        AnnouncementViewPageConfig.currentPageAction = action;
        break;
      case Pages.WrappedAssetView:
        WrappedAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.GenerateDepositAddressView:
        GenerateDepositAddressViewPageConfig.currentPageAction = action;
        break;
      case Pages.SelectDepositAddressView:
        SelectDepositAddressViewPageConfig.currentPageAction = action;
        break;
      case Pages.WithdrawAssetView:
        WithdrawAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.ConfirmWithdrawView:
        ConfirmWithdrawViewPageConfig.currentPageAction = action;
        break;
      case Pages.TransactionStatusView:
        TransactionStatusViewPageConfig.currentPageAction = action;
        break;
      case Pages.DepositWithdrawHistoryView:
        DepositWithdrawHistoryViewPageConfig.currentPageAction = action;
        break;
      case Pages.DepositWithdrawDetailsView:
        DepositWithdrawDetailsViewPageConfig.currentPageAction = action;
        break;
      case Pages.WelcomeSubscriptionsView:
        WelcomeSubscriptionsViewPageConfig.currentPageAction = action;
        break;
      case Pages.SubscriptionPlansView:
        SubscriptionPlansViewPageConfig.currentPageAction = action;
        break;
      case Pages.SubscriptionPlanBenefitsView:
        SubscriptionPlanBenefitsViewPageConfig.currentPageAction = action;
        break;
      case Pages.SubscriptionPlanOptionsView:
        SubscriptionPlanOptionsViewPageConfig.currentPageAction = action;
        break;
      case Pages.AuthorizeSubscriptionView:
        AuthorizeSubscriptionViewPageConfig.currentPageAction = action;
        break;
      case Pages.OptInAssetView:
        OptInAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.OptOutAssetView:
        OptOutAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.OptInOutAssetView:
        OptInOutAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.TokenizationWelcomeView:
        TokenizationWelcomeViewPageConfig.currentPageAction = action;
        break;
      case Pages.SettingsView:
        SettingsViewPageConfig.currentPageAction = action;
        break;
      case Pages.TokenizeAssetView:
        TokenizeAssetViewPageConfig.currentPageAction = action;
        break;
      case Pages.AssetInformationView:
        AssetInformationViewPageConfig.currentPageAction = action;
        break;
      case Pages.AssetTokenInformationView:
        AssetTokenInformationViewPageConfig.currentPageAction = action;
        break;
      case Pages.AssetVerificationDocumentsView:
        AssetVerificationDocumentsViewPageConfig.currentPageAction = action;
        break;
      default:
        break;
    }
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
}
