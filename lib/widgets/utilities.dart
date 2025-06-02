import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:ui';
import 'package:easy_localization/easy_localization.dart';
import 'package:expandable/expandable.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:share_plus/share_plus.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:path_provider/path_provider.dart' as syspaths;
import 'package:pdf/widgets.dart' as pw;
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/cache.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/countdown.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:flutter_idensic_mobile_sdk_plugin/flutter_idensic_mobile_sdk_plugin.dart';
import '../utils/medeiaqury/medeiaqury.dart';

void showSnackBar(String rel, BuildContext context) {
  var notifier = Provider.of<ColorNotifier>(context, listen: false);
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      backgroundColor: notifier.getbluecolor,
      content: Text(
        '$rel ${"copiedsuccessfully".tr()}',
        style: TextStyle(
          color: wihitecolor,
          fontSize: 12.sp,
          fontWeight: FontWeight.w500,
          fontFamily: fontbody,
        ),
      ),
      action: SnackBarAction(
        label: "dismiss".tr(),
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
        label: "dismiss".tr(),
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

String truncateToDecimalPlaces(double number, {int decimalPlaces = 7}) {
  // Convert number to string with high precision to avoid initial rounding
  String numStr = number.toString();

  // Split into integer and fractional parts
  List<String> parts = numStr.split('.');
  if (parts.length < 2 || decimalPlaces <= 0) {
    return NumberFormat(
            decimalPlaces == 2 ? "#,##0.##" : "#,##0.#######", "en_US")
        .format(number
            .truncateToDouble()); // No decimal part or no decimals requested
  }

  String integerPart = parts[0];
  String fractionalPart = parts[1];

  // Truncate the fractional part to the desired length
  if (fractionalPart.length > decimalPlaces) {
    fractionalPart = fractionalPart.substring(0, decimalPlaces);
  }
  // Combine and parse back to double
  String truncatedStr = '$integerPart.$fractionalPart';
  return NumberFormat(
          decimalPlaces == 2 ? "#,##0.##" : "#,##0.#######", "en_US")
      .format(double.parse(truncatedStr));
}

formatNumber(double number) {
  var formattedString = NumberFormat("#,##0.#######", "en_US").format(number);
  // var splitFormattedString = formattedString.split('.');
  // // if (int.parse(splitFormattedString[1]) == 0) {
  // //   return splitFormattedString[0];
  // // }

  return formattedString;
}

formatNumberShort(double number) {
  var formattedString = NumberFormat("#,##0.##", "en_US").format(number);
  // var splitFormattedString = formattedString.split('.');
  // if (int.parse(splitFormattedString[1]) == 0) {
  //   return splitFormattedString[0];
  // }

  return formattedString;
}

formatNumberForInput(double number) {
  var splitNumber = number.toString().split('.');
  return '${NumberFormat("#,##0", "en_US").format(double.parse(splitNumber[0]))}${splitNumber[1] != '0' ? '.' + splitNumber[1] : ''}';
}

formatHistoryNumber(double number, double trimNum, {bool isShort = false}) {
  // if number is greater than 1million return 1m or 1.2m
  if (number >= trimNum) {
    return NumberFormat.compact().format(number);
  }

  return isShort ? formatNumberShort(number) : formatNumber(number);
}

extension on double {
  // Like [toStringAsFixed] but truncates (toward zero) to the specified
  // number of fractional digits instead of rounding.
  // String toStringAsTruncated(int fractionDigits) {
  //   // Require same limits as [toStringAsFixed].
  //   assert(fractionDigits >= 0);
  //   assert(fractionDigits <= 20);

  //   if (fractionDigits == 0) {
  //     return truncateToDouble().toString();
  //   }

  //   // [toString] will represent very small numbers in exponential form.
  //   // Instead use [toStringAsFixed] with the maximum number of fractional
  //   // digits.
  //   var s = toStringAsFixed(20);

  //   // [toStringAsFixed] will still represent very large numbers in
  //   // exponential form.
  //   if (s.contains('e')) {
  //     // Ignore values in exponential form.
  //     return s;
  //   }

  //   // Ignore unrecognized values (e.g. NaN, +infinity, -infinity).
  //   var i = s.indexOf('.');
  //   if (i == -1) {
  //     return s;
  //   }

  //   return s.substring(0, i + fractionDigits + 1);
  // }
}

truncatePublicKey(String? publicKey) {
  if (publicKey == null) return "enterpublickey".tr();
  if (publicKey.length <= 7) return publicKey;
  return truncate(publicKey, length: 7) +
      publicKey.substring(publicKey.length - 7);
}

truncateString(String? text) {
  if (text == null) return "entertext".tr();
  if (text.length <= 15) return text;
  return truncate(text, length: 15) + text.substring(text.length - 15);
}

String truncate(String text, {length = 7, omission = '...'}) {
  if (length >= text.length) {
    return text;
  }
  return text.replaceRange(length, text.length, omission);
}

Account? parseKey(BuildContext context, String secretKey) {
  try {
    Account account = TrovoWalletSDK().parseSecretKey(secretKey.toUpperCase());
    return account;
  } catch (e) {
    print(e);
    // must be some sort of server error
    // let's throw it
    popup(context, title: "error".tr(), message: "invalidcredentials".tr());
    return null;
  }
}

class doubleTypeFormatter extends TextInputFormatter {
  doubleTypeFormatter();
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    return TextEditingValue(
        text:
            formatNumberShort(double.parse(newValue.text.replaceAll(',', ''))),
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

String calculateFiatValue(String assetBalance, String usdPrice, String currency,
        DataProvider appState) =>
    formatNumber(double.parse(getFiatRate(usdPrice, currency, appState,
                getUnFormatted: true)) *
            double.parse(assetBalance))
        .toString();

String getFiatRate(String usdPrice, String currency, DataProvider appState,
    {bool getUnFormatted = false}) {
  usdPrice = usdPrice.isEmpty ? '0' : usdPrice;
  if (getUnFormatted)
    return (appState.fiatRate[currency] * double.parse(usdPrice)).toString();

  return NumberFormat("#,##0.00", "en_US")
      .format(appState.fiatRate[currency] * double.parse(usdPrice))
      .toString();
}

String getTotalFiatBalanceOfAllAssetsInWallet(
    String currency, DataProvider appState, List<Asset> assets) {
  double balance = 0;
  if (assets.length > 0) {
    for (var asset in assets) {
      balance += double.parse(calculateFiatValue(asset.amount.toString(),
              asset.usdPrice.toString(), currency, appState)
          .replaceAll(',', ''));
    }
  }
  return formatHistoryNumber(
    double.parse(balance.toString()),
    1000000,
    isShort: true,
  );
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
                      "learnmore".tr(),
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
                        "whatissharedaccess".tr(),
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
  double? fontSize = 12,
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

String getExplorerBaseUrl(String walletMode) {
  return walletMode == 'Mainnet'
      ? bantuBlockchainExplorerBaseUrl
      : bantuBlockchainExplorerTestnetBaseUrl;
}

Widget iconDropdown(
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
      hint: Icon(
        Icons.filter_list_outlined,
        color: notifier.getbluewhitecolor,
        size: 20,
      ),
      icon: Container(),
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

Widget dropdown(
  void Function(Object?)? onChanged,
  List<DropdownMenuItem<Object>> items,
  Object? value,
  String? hint,
  BuildContext context,
  List<Widget> Function(BuildContext)? selectedItemBuilder, {
  String? Function(Object?)? validator,
  void Function()? onTap,
  double? itemHeight,
}) {
  var notifier = Provider.of<ColorNotifier>(context, listen: true);
  return Padding(
    padding: const EdgeInsets.symmetric(horizontal: 5.0),
    child: Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        DropdownButtonFormField(
          selectedItemBuilder: selectedItemBuilder,
          isDense: true,
          isExpanded: true,
          itemHeight: itemHeight,
          validator: validator,
          onTap: onTap,
          hint: Container(
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
          dropdownColor: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
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
            fillColor: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
          ),
          value: value,
          icon: onChanged == null
              ? null
              : Icon(
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
      ],
    ),
  );
}

Future<void> share(String message, GlobalKey snapshotAreaKey) async {
  final appDir = await syspaths.getTemporaryDirectory();
  String fileName = '${appDir.path}/receipt.png';

  RenderRepaintBoundary boundary = snapshotAreaKey.currentContext!
      .findRenderObject()! as RenderRepaintBoundary;

  var image = await boundary.toImage();
  var byteData = await image.toByteData(format: ImageByteFormat.png);
  File file = await File(fileName).create();
  file.writeAsBytesSync(byteData!.buffer.asUint8List());
  await Share.shareXFiles([XFile(fileName)],
      text: message, sharePositionOrigin: boundary.paintBounds);
}

Future<void> sharePDF(String message, GlobalKey snapshotAreaKey) async {
  final appDir = await syspaths.getTemporaryDirectory();
  String fileName = '${appDir.path}/receipt.pdf';
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
  await Share.shareXFiles([XFile(fileName)],
      text: message, sharePositionOrigin: boundary.paintBounds);
}

void disableSharedAccess(
    BuildContext context, DataProvider appState, Wallet wallet,
    {bool viewOnly = false}) async {
  try {
    showLoader(context);

    Map responseData = await makeDeleteRequest(
      uri: '/v1/shared-access/users/account',
      body: '{}',
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: wallet.publicKey!,
    );

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200 ||
        responseData['statusCode'] == 202) {
      var messageLength = responseData['data']['messages'].length;
      var messageShown = 0;
      print('messagelenth: $messageLength');
      postProcessData(
          context, messageShown, messageLength, responseData['data'],
          callback: () {
        signAndCommitTransaction(
            responseData['data'], context, appState, wallet, viewOnly);
      });
      hideLoader(context);
    } else {
      popup(context,
          title: "error".tr(), message: responseData['data']['message']);
      hideLoader(context);
    }
  } catch (e) {
    popup(context, title: "error".tr(), message: e.toString());
    hideLoader(context);
  }
}

void signAndCommitTransaction(responseFromServer, BuildContext context,
    DataProvider appState, Wallet wallet, bool viewOnly) async {
  try {
    String viewOnlySuccess =
        "sharedaccessdisabledsuccessfully".tr(args: [wallet.alias!]);
    String sharedAccessSuccess =
        "sharedaccessdisablerequestsuccessful".tr(args: [wallet.alias!]);
    showLoader(context);

    //sign the transaction and the submit again
    var signature = TrovoWalletSDK().signBase64Txn(
      appState.secretKeys[0], // the primary wallet secret key,
      responseFromServer['transaction'],
      responseFromServer['networkPassPhrase'],
    );
    responseFromServer['transactionId'] = "";
    responseFromServer['transactionSignature'] = signature;
    responseFromServer['commit'] = 1;

    String requestBody = jsonEncode(responseFromServer);

    Map responseData = await makeDeleteRequest(
      uri: '/v1/shared-access/users/account',
      body: requestBody,
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: wallet.publicKey!,
    );

    if (responseData['statusCode'] == 200 ||
        responseData['statusCode'] == 202) {
      updateUserInfo(
        appState.primaryWallet.signer!,
        appState.secretKeys[0],
        appState.primaryWallet.publicKey,
        appState.userInfo!.username,
        appState,
        forceRefresh: true,
      );
      appState.viewData![SuccessViewPageConfig.key] = {
        'title': "requestsubmitted".tr(),
        'message': viewOnly ? viewOnlySuccess : sharedAccessSuccess,
        'useOnDone': true,
        'onDone': () {
          appState.currentAction = PageAction(state: PageState.addAll, pages: [
            BottomHomePageConfig,
            SharedAccessViewPageConfig,
          ]);
        },
      };
      appState.currentAction =
          PageAction(state: PageState.replace, page: SuccessViewPageConfig);
      hideLoader(context);
    } else {
      popup(context,
          title: "error".tr(), message: responseData['data']['message']);
      hideLoader(context);
    }
  } catch (e) {
    popup(context, title: "error".tr(), message: e.toString());
    hideLoader(context);
  }
}

Widget tokenizedAssetTile({
  required ColorNotifier notifier,
  required TokenizedAsset asset,
  void Function()? onSubscribe,
  void Function()? onBuyToken,
}) {
  return Card(
    elevation: notifier.isDark ? 0 : 5,
    shadowColor: Colors.black,
    color: notifier.gettilewihitecolor,
    margin: EdgeInsets.symmetric(vertical: 10, horizontal: 5),
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(15.0),
    ),
    child: Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: ListTile(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (asset.assetLogo != null) ...[
                  Stack(
                    alignment: Alignment.topRight,
                    children: [
                      CircleAvatar(
                        backgroundColor: Colors.transparent,
                        child: ClipRRect(
                          borderRadius: BorderRadius.circular(100.0),
                          child: Image.network(
                            asset.assetLogo!,
                            height: 38,
                            width: 38,
                            fit: BoxFit.fill,
                          ),
                        ),
                      ),
                      Container(
                        width: 10,
                        height: 10,
                        decoration: BoxDecoration(
                          borderRadius:
                              const BorderRadius.all(Radius.circular(15.0)),
                          border: Border.all(
                            color: notifier.getwihitecolor,
                            width: 1,
                          ),
                          color: asset.assetAlreadyExists == 1
                              ? notifier.getgreencolor
                              : notifier.getupcomingassetyellow,
                        ),
                      ),
                    ],
                  ),
                ] else ...[
                  Image.asset(
                    'assets/images/trovo.png',
                    height: 35,
                    width: 35,
                  ),
                ],
                SizedBox(width: 10),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Container(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                '${asset.assetName!.capitalizeEachWord()} (${asset.assetCode!.toUpperCase()})',
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontsemibold,
                                  color: notifier.getblck,
                                ),
                              ),
                              Padding(
                                padding:
                                    const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                                child: Text(
                                  asset.assetSector!,
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontbody,
                                    color: notifier.getblck,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ],
            ),
            if (onSubscribe != null) ...[
              if (asset.tokenizationStatus == 4) ...[
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    SizedBox(
                      width: 160,
                      child: Countdown(startDate: asset.salesStart!),
                    ),
                    ElevatedButton(
                      onPressed: () async {
                        onSubscribe();
                      },
                      style: ButtonStyle(
                        padding: MaterialStateProperty.all(
                          EdgeInsets.symmetric(vertical: 0, horizontal: 6),
                        ),
                        overlayColor: MaterialStateProperty.all<Color>(
                            notifier.getsplashgrey),
                        backgroundColor: MaterialStateProperty.all<Color>(
                          asset.expressedInterest ?? false
                              ? notifier.getbluewhitecolor
                              : notifier.getwihitecolor,
                        ),
                        side: MaterialStateProperty.all(
                          BorderSide(
                              color: notifier.getbluewhitecolor,
                              width: 1,
                              style: BorderStyle.solid),
                        ),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            if (asset.expressedInterest ?? false) ...[
                              Container(
                                width: 70,
                                child: Text(
                                  'Interest Expressed',
                                  style: TextStyle(
                                    fontFamily: fontsemibold,
                                    fontSize: 12,
                                    overflow: TextOverflow.visible,
                                    color: asset.expressedInterest ?? false
                                        ? notifier.getwihitecolor
                                        : notifier.getbluewhitecolor,
                                  ),
                                ),
                              ),
                              Icon(
                                Icons.check_circle_rounded,
                                size: 20,
                                color: asset.expressedInterest ?? false
                                    ? notifier.getwihitecolor
                                    : notifier.getbluewhitecolor,
                              ),
                            ] else ...[
                              Container(
                                width: 53,
                                child: Text(
                                  'Express Interest',
                                  style: TextStyle(
                                    fontFamily: fontsemibold,
                                    fontSize: 12,
                                    color: asset.expressedInterest ?? false
                                        ? notifier.getwihitecolor
                                        : notifier.getbluewhitecolor,
                                  ),
                                ),
                              ),
                              Icon(
                                Icons.add_circle_rounded,
                                size: 20,
                                color: asset.expressedInterest ?? false
                                    ? notifier.getwihitecolor
                                    : notifier.getbluewhitecolor,
                              ),
                            ]
                          ]),
                    ),
                  ],
                ),
              ] else if (asset.tokenizationStatus == 5) ...[
                Container(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Card(
                        margin: EdgeInsets.zero,
                        shadowColor: Colors.black,
                        shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(5.0),
                            side: BorderSide(
                              color: notifier.getbluewhitecolor,
                              width: 1,
                            )),
                        color: notifier.isDark
                            ? notifier.getbluecolor90
                            : notifier.getaddsubwalletgrey,
                        child: Padding(
                          padding: const EdgeInsets.all(5.0),
                          child: Text(
                            'Available',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                              overflow: TextOverflow.visible,
                            ),
                          ),
                        ),
                      ),
                      ElevatedButton(
                        onPressed: onBuyToken,
                        style: ButtonStyle(
                          padding: MaterialStateProperty.all(EdgeInsets.zero),
                          overlayColor: MaterialStateProperty.all<Color>(
                              notifier.getbluecolor90),
                          backgroundColor: MaterialStateProperty.all<Color>(
                              notifier.getbluewhitecolor),
                          side: MaterialStateProperty.all(
                            BorderSide(
                                color: notifier.getbluewhitecolor,
                                width: 1,
                                style: BorderStyle.solid),
                          ),
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                          shape:
                              MaterialStateProperty.all<RoundedRectangleBorder>(
                            const RoundedRectangleBorder(
                              borderRadius: BorderRadius.all(
                                Radius.circular(10),
                              ),
                            ),
                          ),
                        ),
                        child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Buy',
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                  fontSize: 12,
                                  color: notifier.getwihitecolor,
                                ),
                              ),
                              Image.asset(
                                'assets/images/money.png',
                                color: notifier.getwihitecolor,
                              ),
                            ]),
                      ),
                    ],
                  ),
                ),
              ],
            ],
          ],
        ),
      ),
    ),
  );
}

String getFiatValue(double amount) {
  if (amount < 99000000000) {
    return formatNumberShort(amount);
  }
  return formatHistoryNumber(amount, 99000000000)
      .toString()
      .replaceAll(',', '');
}

Widget infoCard(ColorNotifier notifier,
    {required String label, required String value, String? extraValue}) {
  return Container(
    width: width / 2.1,
    child: Card(
      shadowColor: Colors.black,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      color: notifier.isDark ? notifier.getbluecolor90 : notifier.getpillbg,
      child: TextButton(
        onPressed: () {},
        child: Row(
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(
                  width: width / 2.5,
                  child: Text(
                    label,
                    textAlign: TextAlign.start,
                    overflow: TextOverflow.visible,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
                SizedBox(
                  height: 6,
                ),
                SizedBox(
                  width: width / 2.5,
                  child: Text(
                    value,
                    textAlign: TextAlign.start,
                    style: TextStyle(
                      fontSize: 14,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    ),
  );
}

Widget categoryTile(
  ColorNotifier notifier, {
  required String label,
  required String imageUrl,
  required void Function() onTap,
}) {
  return GestureDetector(
    onTap: onTap,
    child: Card(
      elevation: notifier.isDark ? 0 : 3,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(horizontal: 10),
      child: ListTile(
        title: Row(
          children: [
            Image.asset(imageUrl, width: 30),
            SizedBox(width: 10),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      label,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
        trailing: Icon(
          Icons.arrow_forward_ios,
          color: notifier.getbluewhitecolor,
          size: 20,
        ),
      ),
    ),
  );
}

Widget infoTile(ColorNotifier notifier, String key, String value) {
  return Card(
    // elevation: notifier.isDark ? 0 : 3,
    shadowColor: Colors.black,
    color: notifier.isDark
        ? notifier.getbluecolor90
        : notifier.getaddsubwalletgrey,
    margin: EdgeInsets.symmetric(vertical: 5, horizontal: 10),
    child: Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: ListTile(
        title: Row(
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(
                  width: width / 1.18,
                  child: Text(
                    key,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
                Container(
                  width: width / 1.2,
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      value.isEmpty ? 'Nil' : value,
                      overflow: TextOverflow.visible,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    ),
  );
}

Widget pill(
  String name, {
  required Color backColor,
  required Color foreColor,
  double? fontSize = 12,
  bool hideDirectionUp = false,
  bool showDirectionDown = false,
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
            Row(
              children: [
                Text(
                  name,
                  textAlign: TextAlign.center,
                  softWrap: true,
                  style: TextStyle(
                      color: foreColor,
                      fontFamily: fontbody,
                      fontSize: fontSize),
                ),
                if (!hideDirectionUp) ...[
                  Icon(
                    showDirectionDown
                        ? Icons.arrow_downward
                        : Icons.arrow_upward,
                    color: foreColor,
                    size: 18,
                  ),
                ],
              ],
            ),
            SizedBox(
              width: width / 70,
            ),
          ],
        ),
      ),
    ),
  );
}

extension StringCasing on String {
  String capitalizeFirstLetter() {
    if (isEmpty) return this;
    return this[0].toUpperCase() + substring(1);
  }

  String capitalizeEachWord() {
    return split(' ').map((word) => word.capitalizeFirstLetter()).join(' ');
  }
}

fetchKycConfig(DataProvider appState) async {
  try {
    var uri = '/v1/users/kyc/sumsub/configs';
    Map responseData = await makeGetRequest(
      uri: Uri.encodeFull(uri),
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.signer!,
    );
    print('===============> response ${responseData}');
    if (responseData['statusCode'] == 200) {
      return responseData['data'];
    } else {
      return {};
    }
  } catch (e) {
    print('error');
    print(e);
    throw e;
  }
}

initiateKyc(context, DataProvider appState, String kycLevel) async {
  try {
    showLoader(context);

    Map responseData = await makePostRequest(
      uri: '/v1/users/kyc/sumsub/initiate/${kycLevel}',
      body: '{}',
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.publicKey!,
    );

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200 ||
        responseData['statusCode'] == 202) {
      hideLoader(context);
      return responseData['data'];
    } else {
      popup(context,
          title: "error".tr(), message: responseData['data']['message']);
      hideLoader(context);
    }
  } catch (e) {
    popup(context, title: "error".tr(), message: e.toString());
    hideLoader(context);
  }
}

completeKyc(context, DataProvider appState, String kycLevel) async {
  try {
    showLoader(context);

    Map responseData = await makePostRequest(
      uri: '/v1/users/kyc/sumsub/complete/${kycLevel}',
      body: '{}',
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.publicKey!,
    );

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200 ||
        responseData['statusCode'] == 202) {
      hideLoader(context);
      return responseData['data'];
    } else {
      popup(context,
          title: "error".tr(), message: responseData['data']['message']);
      hideLoader(context);
    }
  } catch (e) {
    popup(context, title: "error".tr(), message: e.toString());
    hideLoader(context);
  }
}

void launchSDK(context, DataProvider appState) async {
  var response = await fetchKycConfig(appState);
  String kycLevel = '';
  var kycVerified = appState.userInfo!.kycVerified;

  if (kycVerified == 0 && response['kycProgress']['kycLevel1Done'] == 0) {
    kycLevel = 'individual-level-1';
  } else if (kycVerified == 1 &&
      response['kycProgress']['kycLevel2Done'] == 0) {
    kycLevel = 'individual-level-2';
  } else if (kycVerified == 2 &&
      response['kycProgress']['kycLevel3Done'] == 0) {
    kycLevel = 'individual-level-3';
  }

  if (kycLevel.isEmpty) {
    popup(context,
        title: 'Error',
        message:
            'You cannot initiate another KYC at this moment. Please wait for your initiated KYC to complete.');
    return;
  }

  var res = await initiateKyc(context, appState, kycLevel);

  print('res ======> $res');
  String accessToken = res['applicantToken'];
  print('accessToken $accessToken');

  // From your backend get an access token for the applicant to be verified.
  // The token must be generated with `levelName` and `userId`,
  // where `levelName` is the name of a level configured in your dashboard.
  //
  // The sdk will work in the production or in the sandbox environment
  // depend on which one the `accessToken` has been generated on.
  //

  // // The access token has a limited lifespan and when it's expired, you must provide another one.
  // So be prepared to get a new token from your backend.
  final onTokenExpiration = () async {
    // call your backend to fetch a new access token (this is just an example)
    return Future<String>.delayed(Duration(seconds: 2), () async {
      response = await initiateKyc(context, appState, kycLevel);
      return response['applicantToken'];
    });
  };

  final SNSStatusChangedHandler onStatusChanged =
      (SNSMobileSDKStatus newStatus, SNSMobileSDKStatus prevStatus) {
    print("The SDK status was changed: $prevStatus -> $newStatus");
  };

  final snsMobileSDK = SNSMobileSDK.init(accessToken, onTokenExpiration)
      .withHandlers(
          // optional handlers
          onStatusChanged: onStatusChanged)
      .withDebug(true) // set debug mode if required
      .withLocale(Locale(
          "en")) // optional, for cases when you need to override the system locale
      .build();

  final SNSMobileSDKResult result = await snsMobileSDK.launch();

  print(
      "=============================================>>>>>>>>>>>>>>>>>>>>>>Completed with result: $result");
  await completeKyc(context, appState, kycLevel);
  showSuccessAlert(context, onTap: () {});
}

Widget getDrawer(
    BuildContext context, DataProvider appState, ColorNotifier notifier) {
  return Drawer(
    backgroundColor: notifier.getwihitecolor,
    child: ListView(
      // Important: Remove any padding from the ListView.
      padding: EdgeInsets.zero,
      children: [
        GestureDetector(
          onTap: () {
            appState.currentAction = PageAction(
                state: PageState.addPage, page: ProfileDetailsViewPageConfig);
          },
          child: UserAccountsDrawerHeader(
            decoration: BoxDecoration(
              color: notifier.getwihitecolor,
            ),
            margin: const EdgeInsets.only(bottom: 8.0),
            accountName: Text(
              '${appState.userInfo!.fullName.toLowerCase().capitalizeEachWord()} (${appState.userInfo!.username!.toLowerCase().capitalizeEachWord()})',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
            accountEmail: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  appState.userInfo!.email!,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                SizedBox(height: 4),
                Text(
                  appState.userInfo!.kycVerified != null &&
                          appState.userInfo!.kycVerified! > 0
                      ? 'Verified (Level ${appState.userInfo!.kycVerified})'
                      : 'Unverified',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontFamily: fontsemibold,
                    color: appState.userInfo!.kycVerified != null &&
                            appState.userInfo!.kycVerified! > 0
                        ? notifier.getgreencolor
                        : Colors.red,
                  ),
                ),
              ],
            ),
            currentAccountPicture: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  child: Stack(
                    children: [
                      CircleAvatar(
                        radius: 30,
                        backgroundColor: notifier.getbluecolor70,
                        child: GestureDetector(
                          onTap: () {
                            appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: ProfileDetailsViewPageConfig);
                          },
                          child: ClipRRect(
                            borderRadius: BorderRadius.circular(100.0),
                            child: Image.network(
                              appState.userInfo!.imageThumbnailURL!,
                              width: width / 6.8,
                              // height: width / 10,
                              fit: BoxFit.fill,
                              errorBuilder: (context, error, stackTrace) {
                                return Image.asset(
                                  'assets/images/trovo.png',
                                  width: width / 9,
                                );
                              },
                            ),
                          ),
                        ),
                      ),
                      if (appState.userInfo != null &&
                          appState.userInfo!.patronMembership != null) ...[
                        Container(
                          width: width / 6.0,
                          height: height / 12.5,
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.end,
                            crossAxisAlignment: CrossAxisAlignment.end,
                            children: [
                              Image.asset(
                                appState.userInfo!.patronMembership
                                        ?.getLogo() ??
                                    'assets/images/trovo.png',
                                width: 30,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/access.png",
            color: notifier.getgrey.withOpacity(.80),
            scale: 5,
            height: height / 35,
          ),
          title: Text(
            "sharedaccess".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            appState.backupSecrets.clear();
            Navigator.pop(context);
            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: appState.introducedSharedAccess
                  ? SharedAccessViewPageConfig
                  : WelcomeToSharedAccessViewPageConfig,
            );
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/help.png",
            color: notifier.getgrey.withOpacity(.80),
            height: height / 40,
          ),
          title: Text(
            "Verify Account",
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            // launchSDK(context, appState);
            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: KycScreenViewPageConfig,
            );
            Navigator.pop(context);
          },
        ),
        // ListTile(
        //   leading: Icon(
        //     CupertinoIcons.doc_chart,
        //     color: notifier.getgrey.withOpacity(.80),
        //   ),
        //   title: Text(
        //     'DEX Trade',
        //     style: TextStyle(
        //       fontFamily: fontbody,
        //       color: notifier.getbluewhitecolor,
        //     ),
        //   ),
        //   onTap: () {
        //     appState.currentAction = PageAction(
        //       state: PageState.addPage,
        //       page: MarketTradeViewPageConfig,
        //     );
        //     Navigator.pop(context);
        //   },
        // ),
        ListTile(
          leading: Image.asset(
            "assets/images/trovo.png",
            color: notifier.getgrey.withOpacity(.80),
            scale: 5,
            height: height / 40,
          ),
          title: Text(
            "trovopatron".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: WelcomeSubscriptionsViewPageConfig);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/asset.png",
            color: notifier.getgrey.withOpacity(.80),
            scale: 1,
            height: height / 40,
          ),
          title: Text(
            "addremoveasset".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: OptInOutAssetViewPageConfig,
            );
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/import.png",
            color: notifier.getgrey.withOpacity(.80),
            scale: 1,
            height: height / 40,
          ),
          title: Text(
            "importwallet".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.currentAction = PageAction(
                state: PageState.addPage, page: ImportWalletPageConfig);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/backup-wallets.png",
            color: notifier.getgrey.withOpacity(.80),
            scale: 1,
            height: height / 40,
          ),
          title: Text(
            "backupwallet".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            // go to the definition of appState.viewData
            // to learn more about viewData
            if (appState.viewData?[EnsurePrivacyPageConfig.key] == null) {
              appState.viewData = {
                EnsurePrivacyPageConfig.key: {},
              };
            }
            appState.viewData?[EnsurePrivacyPageConfig.key] = {
              'rel': 'backupAll',
            };
            appState.currentAction = PageAction(
                state: PageState.addPage, page: EnsurePrivacyPageConfig);
            Navigator.pop(context);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/history.png",
            color: notifier.getgrey.withOpacity(.80),
            scale: 5,
            height: height / 40,
          ),
          title: Text(
            "accountrecovery".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            // appState.userInfo!.accountRecoveryEnabled == 1
            //           ? "Enabled"
            //           : ""
            if (appState.userInfo!.hasSecurityQuestions == 0) {
              var primaryWallet = appState.userInfo!.wallets!
                  .firstWhere((wallet) => wallet.primaryWallet == 1);

              appState.returnView = PageAction(
                state: PageState.addAll,
                pages: [
                  BottomHomePageConfig,
                  SetupAccountRecoveryViewPageConfig,
                ],
              );

              appState.viewData = {
                SecurityQuestionsViewPageConfig.key: {
                  'signer': primaryWallet.signer,
                  'publicKey': primaryWallet.publicKey,
                  'secretKey': appState.secretKeys[0],
                  'username': appState.userInfo!.username,
                }
              };
              appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: SecurityQuestionsViewPageConfig);
            } else if (appState.userInfo!.accountRecoveryEnabled == 0) {
              appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: SetupAccountRecoveryViewPageConfig);
            } else {
              appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: DisableAccountRecoveryInfoViewPageConfig);
            }
            Navigator.pop(context);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/settings.png",
            color: notifier.getgrey.withOpacity(.80),
            height: height / 40,
          ),
          title: Text(
            "settings".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: SettingsViewPageConfig,
            );
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/help.png",
            color: notifier.getgrey.withOpacity(.80),
            height: height / 40,
          ),
          title: Text(
            "helpandsupport".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.goToWebView(trovoSupportUrl);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/terms.png",
            color: notifier.getgrey.withOpacity(.80),
            height: height / 40,
          ),
          title: Text(
            "termsofuse".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.goToWebView(termsOfServiceUrl);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/copyright.png",
            color: notifier.getgrey.withOpacity(.80),
            height: height / 40,
          ),
          title: Text(
            "abouttrovowallet".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.goToWebView(termsOfServiceUrl);
          },
        ),
        ListTile(
          leading: Image.asset(
            "assets/images/logout.png",
            color: notifier.getgrey.withOpacity(.80),
            height: height / 40,
          ),
          title: Text(
            "logout".tr(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          ),
          onTap: () {
            Navigator.pop(context);
            appState.currentAction =
                PageAction(state: PageState.replaceAll, page: LoginPageConfig);
            appState.isLoggedIn = false;
          },
        ),
        SizedBox(height: height / 20),
      ],
    ),
  );
}
