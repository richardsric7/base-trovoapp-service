import 'package:flutter/material.dart';
import 'package:trovo_wallet/widgets/popups.dart';

void showSnackBar(String rel, BuildContext context) {
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      backgroundColor: notifier.getbluecolor,
      content: Text('$rel copied successfully'),
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
