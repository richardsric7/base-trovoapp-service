import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallets.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'router_delegate.dart';

class TrovoWalletBackButtonDispatcher extends RootBackButtonDispatcher {
  final TrovoWalletRouterDelegate _routerDelegate;

  TrovoWalletBackButtonDispatcher(this._routerDelegate) : super();
  bool dialogOpen = false;

  @override
  Future<bool> didPopRoute() async {
    var appState = Provider.of<DataProvider>(
        _routerDelegate.navigatorKey.currentContext!,
        listen: false);
    print(
        'currentBottomIndex: ${appState.currentBottomTabIndex}; pagesLength: ${_routerDelegate.pages.length}');
    if (_routerDelegate.pages.length <= 1) {
      if (dialogOpen) {
        Navigator.of(
          _routerDelegate.navigatorKey.currentContext!,
          rootNavigator: true,
        ).pop(true);
        dialogOpen = false;
        return true;
      }
      dialogOpen = true;

      // if (appState.currentBottomTabIndex == 1) {
      //   if (appState.walletView.view == WalletView.listWallets &&
      //       dialogOpen == false) {
      //     await _confirmAppExit() ?? true;
      //   } else if (appState.walletView.view == WalletView.addSubWallet) {
      //     appState.walletView.actionIcon = Icons.add_circle_outline_sharp;
      //     appState.walletView.actionText = LanguageEn.addsubwallet;
      //     appState.walletView.view = WalletView.listWallets;
      //     appState.notifyListeners();
      //     var dialogContext;
      //     showDialog<bool>(
      //         context: _routerDelegate.navigatorKey.currentContext!,
      //         builder: (context) {
      //           dialogContext = context;
      //           return Container(
      //             decoration: const BoxDecoration(
      //               color: Colors.white,
      //               borderRadius: BorderRadius.all(
      //                 Radius.circular(30),
      //               ),
      //             ),
      //             child: Text(
      //               'Are you sure you want to close this application?',
      //               style: TextStyle(
      //                 fontSize: 18.0,
      //                 fontWeight: FontWeight.w400,
      //               ),
      //             ),
      //           );
      //         }).then(
      //       (value) => Navigator.of(
      //         dialogContext,
      //         rootNavigator: true,
      //       ).pop(false),
      //     );
      //     return true;
      //   } else if (appState.walletView.view == WalletView.confirmAddSubWallet) {
      //     appState.walletView.actionIcon = Icons.cancel_outlined;
      //     appState.walletView.actionText = LanguageEn.cancel;
      //     appState.walletView.view = WalletView.addSubWallet;
      //     appState.notifyListeners();
      //     return true;
      //   }
      // }
      return await _confirmAppExit() ?? true;
    } else {
      return _routerDelegate.popRoute();
    }
  }

  Future<bool?> _confirmAppExit() {
    return showDialog<bool>(
        context: _routerDelegate.navigatorKey.currentContext!,
        builder: (context) {
          return AlertDialog(
            shape: const RoundedRectangleBorder(
                borderRadius: BorderRadius.all(Radius.circular(25.0))),
            actions: [
              TextButton(
                child: Text(
                  'Yes',
                  style: TextStyle(fontSize: 16.0, color: trovoblue90),
                ),
                onPressed: () => Navigator.of(
                  context,
                  rootNavigator: true,
                ).pop(false),
              ),
              TextButton(
                  child: Text(
                    'No',
                    style: TextStyle(fontSize: 16.0, color: trovoblue90),
                  ),
                  onPressed: () {
                    Navigator.of(
                      context,
                      rootNavigator: true,
                    ).pop(true);
                    dialogOpen = false;
                  }),
            ],
            content: Container(
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.all(
                  Radius.circular(30),
                ),
              ),
              child: Text(
                'Are you sure you want to close this application?',
                style: TextStyle(
                  fontSize: 18.0,
                  fontWeight: FontWeight.w400,
                ),
              ),
            ),
          );
        });
  }
}
