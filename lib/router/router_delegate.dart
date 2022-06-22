// import 'package:fake_bantu_pay/ui/agreements.dart';
// import 'package:fake_bantu_pay/ui/scan_result_view.dart';
// import 'package:flutter/foundation.dart';
// import 'package:flutter/material.dart';
// import 'package:trovo_wallet/storage/state.dart';
// import '../app_state.dart';
// import '../ui/AccountInfo.dart';
// import '../ui/Dashboard.dart';
// import '../ui/EditProfile.dart';
// import '../ui/Onboarding.dart';
// import '../ui/UserProfile.dart';
// import '../ui/asset_details.dart';
// import '../ui/confirm_swap.dart';
// import '../ui/confirm_transaction.dart';
// import '../ui/create_account.dart';
// import '../ui/enableBiometrics.dart';
// import '../ui/import_wallet.dart';
// import '../ui/login.dart';
// import '../ui/qr_scanner_view.dart';
// import '../ui/receive_specific_amount.dart';
// import '../ui/recieve_asset_view.dart';
// import '../ui/request_specific_amount_confirm.dart';
// import '../ui/splash.dart';
// import '../ui/swap_success.dart';
// import '../ui/transaction_details.dart';
// import '../ui/transaction_success.dart';
// import 'PageActions.dart';
// import 'back_dispatcher.dart';
// import 'ui_pages.dart';

// class FakeBantuPayRouterDelegate extends RouterDelegate<PageConfiguration>
//     with ChangeNotifier, PopNavigatorRouterDelegateMixin<PageConfiguration> {
//   final List<MaterialPage> _pages = [];
//   FakeBantuPayBackButtonDispatcher? backButtonDispatcher;

//   @override
//   final GlobalKey<NavigatorState> navigatorKey;
//   final DataProvider appState;

//   FakeBantuPayRouterDelegate(this.appState) : navigatorKey = GlobalKey() {
//     appState.addListener(() {
//       notifyListeners();
//     });
//   }

//   /// Getter for a list that cannot be changed
//   List<MaterialPage> get pages => List.unmodifiable(_pages);

//   /// Number of pages function
//   int numPages() => _pages.length;

//   @override
//   PageConfiguration get currentConfiguration =>
//       _pages.last.arguments as PageConfiguration;

//   @override
//   Widget build(BuildContext context) {
//     return Navigator(
//       key: navigatorKey,
//       onPopPage: _onPopPage,
//       pages: buildPages(),
//     );
//   }

//   bool _onPopPage(Route<dynamic> route, result) {
//     final didPop = route.didPop(result);
//     if (!didPop) {
//       return false;
//     }
//     if (canPop()) {
//       pop();
//       return true;
//     } else {
//       return false;
//     }
//   }

//   void _removePage(MaterialPage page) {
//     if (page != null) {
//       _pages.remove(page);
//     }
//   }

//   void pop() {
//     if (canPop()) {
//       _removePage(_pages.last);
//     }
//   }

//   bool canPop() {
//     return _pages.length > 1;
//   }

//   @override
//   Future<bool> popRoute() {
//     if (canPop()) {
//       _removePage(_pages.last);
//       notifyListeners();
//       return Future.value(true);
//     }
//     return Future.value(false);
//   }

//   MaterialPage _createPage(Widget child, PageConfiguration pageConfig) {
//     return MaterialPage(
//         child: child,
//         key: ValueKey(pageConfig.key),
//         name: pageConfig.path,
//         arguments: pageConfig);
//   }

//   void _addPageData(Widget child, PageConfiguration pageConfig) {
//     _pages.add(
//       _createPage(child, pageConfig),
//     );
//   }

//   void addPage(PageConfiguration pageConfig) {
//     final shouldAddPage = _pages.isEmpty ||
//         (_pages.last.arguments as PageConfiguration).uiPage !=
//             pageConfig.uiPage;
//     if (shouldAddPage) {
//       switch (pageConfig.uiPage) {
//         case Pages.Login:
//           _addPageData(const Login(), LoginPageConfig);
//           break;
//         case Pages.CreateAccount:
//           _addPageData(const CreateAccount(), CreateAccountPageConfig);
//           break;
//         case Pages.Agreements:
//           _addPageData(Agreements(), AgreementsPageConfig);
//           break;
//         case Pages.ImportWallet:
//           _addPageData(const ImportWallet(), ImportWalletPageConfig);
//           break;
//         case Pages.MyQRView:
//           _addPageData(const MyQRView(), MyQRViewPageConfig);
//           break;
//         case Pages.ScanResultPage:
//           _addPageData(const ScanResultPage(), ScanResultPageConfig);
//           break;
//         case Pages.ConfirmTransaction:
//           _addPageData(
//               const ConfirmTransactionPage(), ConfirmTransactionPageConfig);
//           break;
//         case Pages.ConfirmSwap:
//           _addPageData(const ConfirmSwapPage(), ConfirmSwapPageConfig);
//           break;
//         case Pages.TransactionSuccess:
//           _addPageData(
//               const TransactionSuccessPage(), TransactionSuccessPageConfig);
//           break;
//         case Pages.SwapSuccess:
//           _addPageData(const SwapSuccessPage(), SwapSuccessPageConfig);
//           break;
//         case Pages.Dashboard:
//           _addPageData(const Dashboard(), DashboardPageConfig);
//           break;
//         case Pages.TransactionDetail:
//           _addPageData(TransactionDetail(), TransactionDetailPageConfig);
//           break;
//         case Pages.AssetDetails:
//           _addPageData(AssetDetails(), AssetDetailsPageConfig);
//           break;
//         case Pages.RecieveAsset:
//           _addPageData(const RecieveAsset(), RecieveAssetPageConfig);
//           break;
//         case Pages.ReceiveSpecificAmount:
//           _addPageData(
//               const ReceiveSpecificAmount(), ReceiveSpecificAmountPageConfig);
//           break;
//         case Pages.RequestSpecificAmountConfirm:
//           _addPageData(const RequestSpecificAmountConfirm(),
//               RequestSpecificAmountConfirmPageConfig);
//           break;
//         case Pages.AccountInfo:
//           _addPageData(const AccountInfo(), AccountInfoPageConfig);
//           break;
//         case Pages.Onboarding:
//           _addPageData(const Onboarding(), OnboardingPageConfig);
//           break;
//         case Pages.UserProfile:
//           _addPageData(const UserProfile(), UserProfilePageConfig);
//           break;
//         case Pages.EditProfile:
//           _addPageData(const EditProfile(), EditProfilePageConfig);
//           break;
//         case Pages.EnableBiometrics:
//           _addPageData(const EnableBiometrics(), EnableBiometricsPageConfig);
//           break;
//         default:
//           break;
//       }
//     }
//   }

//   void replace(PageConfiguration newRoute) {
//     if (_pages.isNotEmpty) {
//       _pages.removeLast();
//     }
//     addPage(newRoute);
//   }

//   void setPath(List<MaterialPage> path) {
//     _pages.clear();
//     _pages.addAll(path);
//   }

//   void replaceAll(PageConfiguration newRoute) {
//     setNewRoutePath(newRoute);
//   }

//   void push(PageConfiguration newRoute) {
//     addPage(newRoute);
//   }

//   void pushWidget(Widget child, PageConfiguration newRoute) {
//     _addPageData(child, newRoute);
//   }

//   void addAll(List<PageConfiguration> routes) {
//     _pages.clear();
//     routes.forEach((route) {
//       addPage(route);
//     });
//   }

//   @override
//   Future<void> setNewRoutePath(PageConfiguration configuration) {
//     final shouldAddPage = _pages.isEmpty ||
//         (_pages.last.arguments as PageConfiguration).uiPage !=
//             configuration.uiPage;
//     if (shouldAddPage) {
//       _pages.clear();
//       addPage(configuration);
//     }
//     return SynchronousFuture(null);
//   }

//   void _setPageAction(PageAction action) {
//     switch (action.page?.uiPage) {
//       case Pages.Splash:
//         SplashPageConfig.currentPageAction = action;
//         break;
//       case Pages.Login:
//         LoginPageConfig.currentPageAction = action;
//         break;
//       case Pages.CreateAccount:
//         CreateAccountPageConfig.currentPageAction = action;
//         break;
//       case Pages.Agreements:
//         AgreementsPageConfig.currentPageAction = action;
//         break;
//       case Pages.ImportWallet:
//         ImportWalletPageConfig.currentPageAction = action;
//         break;
//       case Pages.MyQRView:
//         MyQRViewPageConfig.currentPageAction = action;
//         break;
//       case Pages.ScanResultPage:
//         ScanResultPageConfig.currentPageAction = action;
//         break;
//       case Pages.ConfirmTransaction:
//         ConfirmTransactionPageConfig.currentPageAction = action;
//         break;
//       case Pages.ConfirmSwap:
//         ConfirmSwapPageConfig.currentPageAction = action;
//         break;
//       case Pages.TransactionSuccess:
//         TransactionSuccessPageConfig.currentPageAction = action;
//         break;
//       case Pages.SwapSuccess:
//         SwapSuccessPageConfig.currentPageAction = action;
//         break;
//       case Pages.Dashboard:
//         DashboardPageConfig.currentPageAction = action;
//         break;
//       case Pages.TransactionDetail:
//         TransactionDetailPageConfig.currentPageAction = action;
//         break;
//       case Pages.AssetDetails:
//         AssetDetailsPageConfig.currentPageAction = action;
//         break;
//       case Pages.RecieveAsset:
//         RecieveAssetPageConfig.currentPageAction = action;
//         break;
//       case Pages.ReceiveSpecificAmount:
//         ReceiveSpecificAmountPageConfig.currentPageAction = action;
//         break;
//       case Pages.RequestSpecificAmountConfirm:
//         RequestSpecificAmountConfirmPageConfig.currentPageAction = action;
//         break;
//       case Pages.AccountInfo:
//         AccountInfoPageConfig.currentPageAction = action;
//         break;
//       case Pages.Onboarding:
//         OnboardingPageConfig.currentPageAction = action;
//         break;
//       case Pages.UserProfile:
//         UserProfilePageConfig.currentPageAction = action;
//         break;
//       case Pages.EditProfile:
//         EditProfilePageConfig.currentPageAction = action;
//         break;
//       case Pages.EnableBiometrics:
//         EnableBiometricsPageConfig.currentPageAction = action;
//         break;
//       default:
//         break;
//     }
//   }

//   List<Page> buildPages() {
//     // if (!appState.splashFinished) {
//       replaceAll(SplashPageConfig);
//     } else {
//       switch (appState.currentAction.state) {
//         case PageState.none:
//           break;
//         case PageState.addPage:
//           _setPageAction(appState.currentAction);
//           addPage(appState.currentAction.page!);
//           break;
//         case PageState.pop:
//           pop();
//           break;
//         case PageState.replace:
//           _setPageAction(appState.currentAction);
//           replace(appState.currentAction.page!);
//           break;
//         case PageState.replaceAll:
//           _setPageAction(appState.currentAction);
//           replaceAll(appState.currentAction.page!);
//           break;
//         case PageState.addWidget:
//           _setPageAction(appState.currentAction);
//           pushWidget(
//               appState.currentAction.widget!, appState.currentAction.page!);
//           break;
//         case PageState.addAll:
//           addAll(appState.currentAction.pages!);
//           break;
//       }
//     }
//     appState.resetCurrentAction();
//     return List.of(_pages);
//   }

//   void parseRoute(Uri uri) {
//     if (uri.pathSegments.isEmpty) {
//       setNewRoutePath(SplashPageConfig);
//       return;
//     }

//     // Handle navapp://deeplinks/details/#
//     if (uri.pathSegments.length == 2) {
//       // if (uri.pathSegments[0] == 'details') {
//       //   pushWidget(Details(int.parse(uri.pathSegments[1])), DetailsPageConfig);
//       // }
//     } else if (uri.pathSegments.length == 1) {
//       final path = uri.pathSegments[0];
//       switch (path) {
//         case 'splash':
//           replaceAll(SplashPageConfig);
//           break;
//         case 'login':
//           replaceAll(LoginPageConfig);
//           break;
//         case 'createAccount':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const CreateAccount(), CreateAccountPageConfig)
//           ]);
//           break;
//         case 'agreements':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const CreateAccount(), CreateAccountPageConfig),
//             _createPage(Agreements(), AgreementsPageConfig)
//           ]);
//           break;
//         case 'importWallet':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const ImportWallet(), ImportWalletPageConfig),
//           ]);
//           break;
//         case 'myQrView':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const MyQRView(), MyQRViewPageConfig),
//           ]);
//           break;
//         case 'scanResultPage':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const ScanResultPage(), ScanResultPageConfig),
//           ]);
//           break;
//         case 'confirmTransaction':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const ScanResultPage(), ScanResultPageConfig),
//             _createPage(
//                 const ConfirmTransactionPage(), ConfirmTransactionPageConfig),
//           ]);
//           break;
//         case 'accountInfo':
//           setPath([
//             _createPage(const Dashboard(), DashboardPageConfig),
//             _createPage(const AccountInfo(), AccountInfoPageConfig),
//           ]);
//           break;
//         case 'onboarding':
//           setPath([
//             _createPage(const Onboarding(), OnboardingPageConfig),
//           ]);
//           break;
//         case 'userProfile':
//           setPath([
//             _createPage(const Dashboard(), DashboardPageConfig),
//             _createPage(const UserProfile(), UserProfilePageConfig),
//           ]);
//           break;
//         case 'editProfile':
//           setPath([
//             _createPage(const Dashboard(), DashboardPageConfig),
//             _createPage(const UserProfile(), UserProfilePageConfig),
//             _createPage(const EditProfile(), EditProfilePageConfig),
//           ]);
//           break;
//         case 'enableBiometrics':
//           setPath([
//             _createPage(const Login(), LoginPageConfig),
//             _createPage(const CreateAccount(), CreateAccountPageConfig),
//             _createPage(const AccountInfo(), AccountInfoPageConfig),
//             _createPage(const EnableBiometrics(), EnableBiometricsPageConfig),
//           ]);
//           break;
//         case 'confirmSwap':
//         case 'transactionSuccess':
//         case 'swapSuccess':
//         case 'dashboard':
//         case 'transactionDetail':
//         case 'assetDetails':
//         case 'recieveAsset':
//         case 'receiveSpecificAmount':
//         case 'requestSpecificAmountConfirm':
//           setPath([
//             _createPage(const Dashboard(), DashboardPageConfig),
//           ]);
//           break;
//       }
//     }
//   }
// }
