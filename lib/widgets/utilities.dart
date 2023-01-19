import 'dart:io';
import 'dart:ui';
import 'package:expandable/expandable.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:share/share.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:path_provider/path_provider.dart' as syspaths;
import 'package:pdf/widgets.dart' as pw;

import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/popups.dart';

import '../utils/medeiaqury/medeiaqury.dart';

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

void showSnackBarForInfo(String message, BuildContext context) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      backgroundColor: notifier.getbluecolor,
      content: Text(
        message,
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

formatHistoryNumber(double number, double trimNum) {
  // if number is greater than 1million return 1m or 1.2m
  if (number >= trimNum) {
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
  WidgetsBinding.instance.addPostFrameCallback((_) {
    if (appState.bottomTabPageController!.hasClients) {
      appState.bottomTabPageController!.animateToPage(index,
          duration: Duration(milliseconds: 500), curve: Curves.easeInOut);

      // set this to the current tab page index
      appState.currentBottomTabIndex = index;
    }
  });
}

void handleDynamicLinkData(Uri parsedUri) {
  print('action: ${parsedUri.queryParameters['action']}');
  print('description: ${parsedUri.queryParameters['description']}');
  print('deviceInfo: ${parsedUri.queryParameters['deviceInfo']}');
  print('targetUser: ${parsedUri.queryParameters['targetUser']}');
  print('ownerUsername: ${parsedUri.queryParameters['ownerUsername']}');
  print('serviceShortName: ${parsedUri.queryParameters['serviceShortName']}');
}

String calculateFiatValue(String assetBalance, String usdPrice, String currency,
        DataProvider appState) =>
    formatNumber(double.parse(getFiatRate(usdPrice, currency, appState)) *
            double.parse(assetBalance))
        .toString();

String getFiatRate(String usdPrice, String currency, DataProvider appState) {
  usdPrice = usdPrice.isEmpty ? '0' : usdPrice;
  return NumberFormat("#,##0.00000", "en_US")
      .format(appState.fiatRate[currency] * double.parse(usdPrice))
      .toString();
}

String getTotalFiatBalanceOfAllAssetsInWallet(
    String currency, DataProvider appState, dynamic assets) {
  double balance = 0;
  if (assets.length > 0) {
    for (var asset in assets) {
      balance += double.parse(calculateFiatValue(asset['amount'].toString(),
              asset['usdPrice'].toString(), currency, appState)
          .replaceAll(',', ''));
    }
  }
  return formatNumber(balance);
}

Widget buildExpandable(context) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;

  return ExpandableNotifier(
      child: ScrollOnExpand(
    child: Container(
      child: Column(
        children: <Widget>[
          ExpandablePanel(
            theme: const ExpandableThemeData(
              headerAlignment: ExpandablePanelHeaderAlignment.center,
              tapBodyToExpand: true,
              tapBodyToCollapse: true,
              hasIcon: false,
            ),
            header: Container(
              child: Padding(
                padding: const EdgeInsets.all(10.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    Text(
                      "Learn more",
                      style: TextStyle(
                        color: notifier.getbluecolor,
                        fontFamily: fontbody,
                        fontSize: 13.sp,
                      ),
                    ),
                    ExpandableIcon(
                      theme: ExpandableThemeData(
                        expandIcon: Icons.keyboard_arrow_right,
                        collapseIcon: Icons.keyboard_arrow_down_outlined,
                        iconColor: notifier.getbluecolor,
                        iconSize: 28.0,
                        iconRotationAngle: 1.9 / 2,
                        iconPadding: EdgeInsets.only(right: 5),
                        hasIcon: false,
                      ),
                    )
                  ],
                ),
              ),
            ),
            collapsed: Container(),
            expanded: Container(
              child: Row(
                children: [
                  SizedBox(
                    width: width / 6,
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        "What is shared access?",
                        style: TextStyle(
                          color: notifier.getbluecolor90,
                          fontFamily: fontsemibold,
                          fontSize: 13.sp,
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      SizedBox(
                        width: width / 1.7,
                        child: Text(
                          "Lorem ipsum dolor emmet what does shared access mean?",
                          style: TextStyle(
                            color: notifier.getbluecolor90,
                            fontFamily: fontbody,
                            fontSize: 13.sp,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      SizedBox(
                        width: width / 1.7,
                        child: Text(
                          "We can also explain more or emphasise very important information here.",
                          style: TextStyle(
                            color: notifier.getbluecolor90,
                            fontFamily: fontbody,
                            fontSize: 13.sp,
                            fontStyle: FontStyle.italic,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    ),
  ));
}

postProcessData(context, messageShown, messageLength, data,
    {required void Function() callback}) {
  // we would like to display all messages returned from the initial
  // request to server using a popup. In order to achieve that we
  // employ the use of a little recursion here. Please recursive
  // functions can turn into a nightmare fast so be carefull here.
  if (messageShown <= messageLength - 1) {
    showResponseMessage(
        context,
        data['messages'][messageShown],
        () => {
              print('postProcessData: $messageShown'),
              postProcessData(context, messageShown, messageLength, data,
                  callback: callback),
            });

    messageShown++;
    return;
  }

  callback();
}

Widget userItem(
  String name,
  void Function()? onClick, {
  required Color backColor,
  required Color foreColor,
  double? fontSize: 12,
  bool restoreMode = false,
}) {
  return Padding(
    padding: const EdgeInsets.all(3.0),
    child: Container(
      decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(10.0)),
          color: backColor),
      child: Padding(
        padding: const EdgeInsets.all(5.0),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                  color: foreColor, fontFamily: fontbody, fontSize: fontSize),
            ),
            SizedBox(
              width: width / 70,
            ),
            if (onClick != null) ...[
              GestureDetector(
                  onTap: () {
                    onClick();
                  },
                  child: Icon(
                    restoreMode ? Icons.replay_sharp : Icons.cancel_outlined,
                    color: wihitecolor,
                  ))
            ]
          ],
        ),
      ),
    ),
  );
}

String trimLeft(String from, String pattern) {
  if ((from).isEmpty || (pattern).isEmpty || pattern.length > from.length)
    return from;

  while (from.startsWith(pattern)) {
    from = from.substring(pattern.length);
  }
  return from;
}

String trimRight(String from, String pattern) {
  if ((from).isEmpty || (pattern).isEmpty || pattern.length > from.length)
    return from;

  while (from.endsWith(pattern)) {
    from = from.substring(0, from.length - pattern.length);
  }
  return from;
}

String trim(String from, String pattern) {
  return trimLeft(trimRight(from, pattern), pattern);
}

Widget dropdown(
  void Function(Object?) onChanged,
  List<DropdownMenuItem<Object>> items,
  Object? value,
  String? hint,
  BuildContext context,
  List<Widget> Function(BuildContext)? selectedItemBuilder,
) {
  var notifier = Provider.of<ColorNotifier>(context, listen: true);
  return Padding(
    padding: const EdgeInsets.symmetric(horizontal: 5.0),
    child: DropdownButtonFormField(
      selectedItemBuilder: selectedItemBuilder,
      isDense: true,
      isExpanded: true,
      hint: Container(
        // width: 150, //and here
        child: hint != null
            ? Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    hint,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              )
            : null,
      ),
      dropdownColor:
          notifier.isDark ? darktilewhitecolor : notifier.getaddsubwalletgrey,
      decoration: InputDecoration(
        contentPadding: EdgeInsets.symmetric(vertical: 0, horizontal: 10),
        enabledBorder: OutlineInputBorder(
          borderSide: BorderSide.none,
          borderRadius: BorderRadius.circular(10),
        ),
        border: OutlineInputBorder(
          borderSide: BorderSide.none,
          borderRadius: BorderRadius.circular(10),
        ),
        filled: true,
        fillColor:
            notifier.isDark ? darktilewhitecolor : notifier.getaddsubwalletgrey,
      ),
      value: value,
      icon: Icon(
        Icons.keyboard_arrow_down_rounded,
        color: notifier.getbluewhitecolor,
      ),
      elevation: 0,
      style: TextStyle(
        color: notifier.getbluewhitecolor,
        fontSize: 15,
        fontFamily: fontsemibold,
        fontWeight: FontWeight.w500,
      ),
      onChanged: onChanged,
      items: items,
    ),
  );
}

Future<void> share(String message, GlobalKey snapshotAreaKey) async {
  final appDir = await syspaths.getExternalStorageDirectory();
  String fileName = '${appDir!.path}/receipt.png';
  RenderRepaintBoundary boundary = snapshotAreaKey.currentContext!
      .findRenderObject()! as RenderRepaintBoundary;

  var image = await boundary.toImage();
  var byteData = await image.toByteData(format: ImageByteFormat.png);
  File file = await File(fileName).create();
  file.writeAsBytesSync(byteData!.buffer.asUint8List());
  print('========================================$fileName');

  Share.shareFiles([fileName], text: message);
}

Future<void> sharePDF(String message, GlobalKey snapshotAreaKey) async {
  final appDir = await syspaths.getExternalStorageDirectory();
  String fileName = '${appDir!.path}/receipt.pdf';
  RenderRepaintBoundary boundary = snapshotAreaKey.currentContext!
      .findRenderObject()! as RenderRepaintBoundary;
  final pdf = pw.Document();

  var image = await boundary.toImage();
  var byteData = await image.toByteData(format: ImageByteFormat.png);

  final pdfImage = pw.MemoryImage(
    byteData!.buffer.asUint8List(),
  );

  pdf.addPage(pw.Page(build: (pw.Context context) {
    return pw.Center(
      child: pw.Image(pdfImage),
    ); // Center
  }));

  File file = await File(fileName).create();
  file.writeAsBytesSync(await pdf.save());

  Share.shareFiles([fileName], text: message);
}
