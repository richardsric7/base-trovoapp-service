import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/widgets/popups.dart';

void showSnackBar(String rel, BuildContext context) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      backgroundColor: notifier.getbluecolor,
      content: Text(
        '$rel copied successfully',
        style: TextStyle(
          color: wihitecolor,
          fontSize: 12.sp,
          fontWeight: FontWeight.w500,
          fontFamily: fontbody,
        ),
      ),
      action: SnackBarAction(
        label: 'DISMISS',
        textColor: wihitecolor,
        onPressed: () => {
          ScaffoldMessenger.of(context).clearSnackBars(),
        },
      ),
    ),
  );
}

getAssetCode(assetCode) {
  // assign XBN to the asset which has an
  // empty assetCode value.
  // native token of the bantu blockchain
  // has empty values as assetCode and
  // assetIssuer
  return assetCode.toString().isEmpty ? nativeAssetCode : assetCode.toString();
}

Color getColor(context, indexOfWallet) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);

  // return only white color if app is in dark mode
  if (notifier.isDark) {
    return wihitecolor;
  }

  return indexOfWallet > 2 ? notifier.getbluecolor : notifier.getwihitecolor;
}

getAssetIssuer(assetIssuer) {
  // assign 'Native Token' to the asset which has an
  // empty assetIssuer value.
  // Native token of the bantu blockchain
  // has empty values as assetCode and
  // assetIssuer
  return assetIssuer.toString().isEmpty
      ? nativeAssetIssuer
      : assetIssuer.toString();
}

formatNumber(double number) =>
    NumberFormat("#,##0.0000", "en_US").format(number);

formatHistoryNumber(double number) {
  // if number is greater than 1million return 1m or 1.2m
  if (number >= 1000000) {
    return NumberFormat.compact().format(number);
  }

  return NumberFormat("#,##0", "en_US").format(number);
}

String truncate(String text, {length: 7, omission: '...'}) {
  if (length >= text.length) {
    return text;
  }
  return text.replaceRange(length, text.length, omission);
}

class doubleTypeFormatter extends TextInputFormatter {
  doubleTypeFormatter();
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    print('this is old value ${oldValue.text}');
    print('this is new value ${newValue.text}');
    return TextEditingValue(
        text: formatNumber(double.parse(newValue.text.replaceAll(',', ''))),
        selection: TextSelection.collapsed(offset: newValue.selection.end + 1));
  }
}

void changeTabPage(appState, index) {
  // moves user to the wallets list tab
  appState.bottomTabPageController!.animateToPage(index,
      duration: const Duration(milliseconds: 500), curve: Curves.ease);
  // set this to the wallets list tab index
  appState.currentBottomTabIndex = index;
}

void handleDynamicLinkData(Uri parsedUri) {
  print('action: ${parsedUri.queryParameters['action']}');
  print('description: ${parsedUri.queryParameters['description']}');
  print('deviceInfo: ${parsedUri.queryParameters['deviceInfo']}');
  print('targetUser: ${parsedUri.queryParameters['targetUser']}');
  print('ownerUsername: ${parsedUri.queryParameters['ownerUsername']}');
  print('serviceShortName: ${parsedUri.queryParameters['serviceShortName']}');
}
