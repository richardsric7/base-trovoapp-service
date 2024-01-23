import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/models/wallets_list_view_data.dart';
import 'package:trovo_wallet/storage/state.dart';
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

    // check if any dialog is open when the back button is pressed
    if (appState.dialogOpen) {
      Navigator.of(
        _routerDelegate.navigatorKey.currentContext!,
        rootNavigator: true,
      ).pop(true);
      appState.dialogOpen = false;
      return true;
    }

    if (_routerDelegate.pages.length <= 1) {
      // check if the close app dialog is open
      if (dialogOpen) {
        Navigator.of(
          _routerDelegate.navigatorKey.currentContext!,
          rootNavigator: true,
        ).pop(true);
        dialogOpen = false;
        return true;
      }

      // handle the back button on the wallets view to make sure the create
      // subwallet views/modals are consistently handled
      if (appState.currentBottomTabIndex == 1) {
        if (appState.walletView.view == WalletView.addSubWallet) {
          appState.walletView.actionIcon = Icons.add_circle_outline_sharp;
          appState.walletView.actionText = "addsubwallet".tr();
          appState.walletView.view = WalletView.listWallets;
          appState.updateListeners();
          return true;
        } else if (appState.walletView.view == WalletView.confirmAddSubWallet) {
          appState.walletView.actionIcon = Icons.cancel_outlined;
          appState.walletView.actionText = "cancel".tr();
          appState.walletView.view = WalletView.addSubWallet;
          appState.updateListeners();
          return true;
        }
      }

      dialogOpen = true;
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
                  "yes".tr(),
                  style: TextStyle(fontSize: 16.0, color: trovoblue90),
                ),
                onPressed: () => Navigator.of(
                  context,
                  rootNavigator: true,
                ).pop(false),
              ),
              TextButton(
                  child: Text(
                    "no".tr(),
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
                "areyousure?".tr(),
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
