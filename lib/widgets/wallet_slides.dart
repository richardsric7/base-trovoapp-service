import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../custom_bloc_observer/constants.dart';
import '../custom_bloc_observer/fonts.dart';
import '../custom_bloc_observer/notifire_clor.dart';
import '../storage/state.dart';
import '../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

import 'popups.dart';

class WalletSlide extends StatefulWidget {
  final String alias;
  final String totalBalance;
  final String? fiatBalance;
  final String? assetCount;
  final Color backColor;
  final Color foreColor;
  final bool isSharedWallet;
  final int walletType;
  final bool initialHiddenState;
  final void Function(bool)? onHiddenStateChanged;

  WalletSlide({
    Key? key,
    required this.alias,
    required this.totalBalance,
    this.fiatBalance,
    this.assetCount,
    required this.isSharedWallet,
    required this.walletType,
    required this.backColor,
    required this.foreColor,
    required this.initialHiddenState,
    this.onHiddenStateChanged,
  }) : super(key: key);

  @override
  State<WalletSlide> createState() => _WalletSlideState();
}

class _WalletSlideState extends State<WalletSlide> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late bool localHideBalance;
  final Authenticator _authenticator = Authenticator();
  final List<IconData> icons = [
    Icons.token_outlined,
    Icons.fire_truck_outlined,
    Icons.account_tree_outlined,
  ];

  @override
  void initState() {
    super.initState();
    localHideBalance = widget.initialHiddenState;
    appState = Provider.of<DataProvider>(context, listen: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return walletSlide();
  }

  Widget walletSlide() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 8, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: widget.backColor,
        ),
        child: Stack(
          alignment: AlignmentDirectional.centerEnd,
          children: [
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 5,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      // Image.asset(
                      //   'assets/images/trovo_white.png',
                      //   width: 40,
                      // ),
                    ],
                  ),
                ),
              ],
            ),
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Container(
                        width: width / 2.3,
                        child: Row(
                          children: [
                            ConstrainedBox(
                              constraints:
                                  BoxConstraints(maxWidth: width / 3.0),
                              child: Container(
                                child: Text(
                                  widget.alias,
                                  overflow: TextOverflow.ellipsis,
                                  style: TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.w600,
                                      color: widget.foreColor,
                                      fontFamily: fontsemibold),
                                ),
                              ),
                            ),
                            // Spacer(),
                            IconButton(
                              padding: EdgeInsets.zero,
                              color: widget.foreColor,
                              constraints: BoxConstraints(),
                              onPressed: () => {
                                Clipboard.setData(
                                  ClipboardData(
                                    text: widget.alias,
                                  ),
                                ),
                                showSnackBar("walletalias".tr(), context),
                              },
                              icon: Icon(
                                Icons.copy,
                                fill: 1.0,
                                size: 15,
                              ),
                            ),
                          ],
                        ),
                      ),
                      if (widget.assetCount != null) ...[
                        Container(
                          child: Text(
                            '${widget.assetCount} Assets',
                            style: TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w600,
                                color: widget.foreColor,
                                fontFamily: fontsemibold),
                          ),
                        ),
                      ],
                    ],
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  Row(
                    children: [
                      Text(
                        "totalbalance".tr(),
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w400,
                          color: widget.foreColor,
                          fontFamily: fontbody,
                        ),
                      ),
                      SizedBox(
                        width: 15,
                      ),
                      GestureDetector(
                        onTap: () {
                          if (localHideBalance) {
                            authenticateAndToggle();
                          } else
                            toggleHideBalance();
                        },
                        child: Icon(
                          getIcon(),
                          size: 18,
                          color: widget.foreColor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 2),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Container(
                            width: width / 1.8,
                            child: Text(
                              getBalance(widget.totalBalance),
                              style: TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                                color: widget.foreColor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          if (widget.fiatBalance != null) ...[
                            SizedBox(height: 2),
                            Text(
                              getBalance(widget.fiatBalance!),
                              style: TextStyle(
                                fontWeight: FontWeight.w300,
                                fontSize: 13,
                                color: widget.foreColor,
                                fontFamily: fontbody,
                              ),
                            ),
                          ],
                        ],
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      Container(
                        width: 40,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            if (widget.isSharedWallet) ...[
                              Icon(
                                Icons.people_outline,
                                size: 17,
                                color: widget.foreColor,
                              )
                            ],
                            if (widget.walletType != 0) ...[
                              Icon(
                                icons[widget.walletType - 1],
                                size: 17,
                                color: widget.foreColor,
                              )
                            ],
                            if (widget.alias.contains('-distribution')) ...[
                              Icon(
                                icons[2],
                                size: 17,
                                color: widget.foreColor,
                              )
                            ],
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  IconData getIcon() {
    IconData icon;
    if (appState.hideBalances) icon = CupertinoIcons.eye;

    if (localHideBalance)
      icon = CupertinoIcons.eye;
    else
      icon = CupertinoIcons.eye_slash;

    return icon;
  }

  void authenticateAndToggle() {
    if (appState.biometricEnabled) {
      toggleBiometrics();
      return;
    }

    showPasswordDialog(context, () {
      toggleHideBalance();
    });
  }

  void toggleBiometrics() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        toggleHideBalance();
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context, callback: toggleHideBalance);
      }
    }
  }

  toggleHideBalance() {
    setState(() {
      localHideBalance = !localHideBalance;
      if (widget.onHiddenStateChanged != null) {
        widget.onHiddenStateChanged!(localHideBalance);
      }
    });
  }

  String getBalance(String balance) {
    String text;
    if (appState.hideBalances) text = hideBalanceText;

    if (localHideBalance)
      text = hideBalanceText;
    else
      text = balance;

    return text;
  }
}
