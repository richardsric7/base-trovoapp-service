import 'package:flutter/material.dart';
import 'ui_pages.dart';

class FakeBantuPayRouteParser
    extends RouteInformationParser<PageConfiguration> {
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
      case CreateAccountPath:
        return CreateAccountPageConfig;
      case AgreementsPath:
        return AgreementsPageConfig;
      case ImportWalletPath:
        return ImportWalletPageConfig;
      case MyQRViewPath:
        return MyQRViewPageConfig;
      case ScanResultViewPath:
        return ScanResultPageConfig;
      case ConfirmTransactionPath:
        return ConfirmTransactionPageConfig;
      case ConfirmSwapPath:
        return ConfirmSwapPageConfig;
      case TransactionSuccessPath:
        return TransactionSuccessPageConfig;
      case SwapSuccessPath:
        return SwapSuccessPageConfig;
      case DashboardPath:
        return DashboardPageConfig;
      case TransactionDetailPath:
        return TransactionDetailPageConfig;
      case RecieveAssetPath:
        return RecieveAssetPageConfig;
      case ReceiveSpecificAmountPath:
        return ReceiveSpecificAmountPageConfig;
      case RequestSpecificAmountConfirmPath:
        return RequestSpecificAmountConfirmPageConfig;
      case AccountInfoPath:
        return AccountInfoPageConfig;
      case OnboardingPath:
        return OnboardingPageConfig;
      case UserProfilePath:
        return UserProfilePageConfig;
      case EditProfilePath:
        return EditProfilePageConfig;
      case EnableBiometricsPath:
        return EnableBiometricsPageConfig;
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
      case Pages.CreateAccount:
        return const RouteInformation(location: CreateAccountPath);
      case Pages.Agreements:
        return const RouteInformation(location: AgreementsPath);
      case Pages.ImportWallet:
        return const RouteInformation(location: ImportWalletPath);
      case Pages.MyQRView:
        return const RouteInformation(location: MyQRViewPath);
      case Pages.ScanResultPage:
        return const RouteInformation(location: ScanResultViewPath);
      case Pages.ConfirmTransaction:
        return const RouteInformation(location: ConfirmTransactionPath);
      case Pages.ConfirmSwap:
        return const RouteInformation(location: ConfirmSwapPath);
      case Pages.TransactionSuccess:
        return const RouteInformation(location: TransactionSuccessPath);
      case Pages.SwapSuccess:
        return const RouteInformation(location: SwapSuccessPath);
      case Pages.Dashboard:
        return const RouteInformation(location: DashboardPath);
      case Pages.TransactionDetail:
        return const RouteInformation(location: TransactionDetailPath);
      case Pages.RecieveAsset:
        return const RouteInformation(location: RecieveAssetPath);
      case Pages.ReceiveSpecificAmount:
        return const RouteInformation(location: ReceiveSpecificAmountPath);
      case Pages.AccountInfo:
        return const RouteInformation(location: AccountInfoPath);
      case Pages.Onboarding:
        return const RouteInformation(location: OnboardingPath);
      case Pages.UserProfile:
        return const RouteInformation(location: UserProfilePath);
      case Pages.EditProfile:
        return const RouteInformation(location: EditProfilePath);
      case Pages.EnableBiometrics:
        return const RouteInformation(location: EnableBiometricsPath);
      default:
        return const RouteInformation(location: SplashPath);
    }
  }
}
