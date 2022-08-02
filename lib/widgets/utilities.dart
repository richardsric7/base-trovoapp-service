import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/widgets/popups.dart';

void showSnackBar(String rel, BuildContext context) {
  notifier = Provider.of<ColorNotifier>(context, listen: false);
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      backgroundColor: notifier.getbluecolor,
      content: Text(
        '$rel copied successfully',
        style: TextStyle(
          color: notifier.getwihitecolor,
          fontSize: 12.sp,
          fontWeight: FontWeight.w500,
          fontFamily: fontbody,
        ),
      ),
      action: SnackBarAction(
        label: 'DISMISS',
        textColor: Colors.white,
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
    NumberFormat("#,##0.000", "en_US").format(number);

String truncate(String text, {length: 7, omission: '...'}) {
  if (length >= text.length) {
    return text;
  }
  return text.replaceRange(length, text.length, omission);
}
