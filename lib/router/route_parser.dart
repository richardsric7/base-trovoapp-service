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
      default:
        return const RouteInformation(location: SplashPath);
    }
  }
}
