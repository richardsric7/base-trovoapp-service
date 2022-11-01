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
