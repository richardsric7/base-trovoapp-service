import 'package:flutter/material.dart';
import 'package:flutter_overlay_loader/flutter_overlay_loader.dart';
import 'package:provider/provider.dart';

import '../custom_bloc_observer/notifire_clor.dart';

showLoader(context) {
  ColorNotifier notifier = Provider.of<ColorNotifier>(context, listen: false);
  return Loader.show(
    context,
    progressIndicator: CircularProgressIndicator(
      backgroundColor: notifier.getbluecolor,
      valueColor: new AlwaysStoppedAnimation<Color>(notifier.getgreencolor),
      strokeWidth: 3.0,
    ),
    themeData: Theme.of(context).copyWith(
      colorScheme: ColorScheme.fromSwatch().copyWith(
        secondary: notifier.getbluecolor,
      ),
    ),
  );
  //return navigatorKey.currentContext.loaderOverlay.show();
}

hideLoader(context) {
  return Loader.hide();
}
