import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import '../custom_bloc_observer/constants.dart';
import '../custom_bloc_observer/fonts.dart';
import '../custom_bloc_observer/notifire_clor.dart';
import '../storage/state.dart';
import '../utils/enstring.dart';
import '../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

import 'popups.dart';

class WalletSlide extends StatefulWidget {
  String alias;
  String totalBalance;
  String? fiatBalance;
  Color backColor;
  Color foreColor;
  bool initialHiddenState;
  void Function(bool)? onHiddenStateChanged;

  WalletSlide({
    Key? key,
    required this.alias,
    required this.totalBalance,
    this.fiatBalance,
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
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: widget.backColor,
        ),
        child: Stack(
          alignment: AlignmentDirectional.centerEnd,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(
                      vertical: 35.0, horizontal: 20),
                  child: Image.asset(
                    'assets/images/trovo_white.png',
                    width: 80,
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
                  Container(
                    width: width / 2,
                    child: Text(
                      widget.alias,
                      style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w600,
                          color: widget.foreColor,
                          fontFamily: fontsemibold),
                    ),
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
                          size: 20,
                          color: widget.foreColor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 98.0,
                  ),
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
                  SizedBox(height: 2),
                  if (widget.fiatBalance != null) ...[
                    Text(
                      getBalance(widget.fiatBalance!),
                      style: TextStyle(
                        fontWeight: FontWeight.w300,
                        fontSize: 13,
                        color: widget.foreColor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ]
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
        biometricsErrorAlert(context);
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
