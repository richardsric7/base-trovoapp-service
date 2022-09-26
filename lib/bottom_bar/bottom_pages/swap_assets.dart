import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:syncfusion_flutter_core/theme.dart';
import 'package:syncfusion_flutter_sliders/sliders.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class SwapAssets extends StatefulWidget {
  const SwapAssets({Key? key}) : super(key: key);

  @override
  State<SwapAssets> createState() => _SwapAssetsState();
}

class _SwapAssetsState extends State<SwapAssets> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  List<Wallet>? wallets;
  Wallet? activeWallet;
  var activeAsset;
  var claimedAssets;
  bool amountError = false;
  double amount = 0;
  var sourceAsset;
  var destinationAsset;
  dynamic selectedWallet = '';
  final formKey = GlobalKey<FormState>();
  bool sourceErr = false;
  bool destErr = false;
  final GlobalKey<FormFieldState> key1 = GlobalKey<FormFieldState>();
  final GlobalKey<FormFieldState> key2 = GlobalKey<FormFieldState>();
  final GlobalKey<FormFieldState> key3 = GlobalKey<FormFieldState>();
  final textController = TextEditingController();
  double? sliderValue = 0;

  List<DropdownMenuItem<String>> get walletDropdownItems {
    return wallets!
        .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
            child: Text(
              wallet.alias!,
              overflow: TextOverflow.ellipsis,
            ),
            value: wallet.publicKey))
        .toList();
  }

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    activeWallet = appState.activeWallet;
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    userInfo = appState.userInfo!;
    wallets = userInfo.wallets!;
    activeWallet = appState.activeWallet;
    selectedWallet = activeWallet!.publicKey;
    assetBalances = appState.assetBalances;
    claimedAssets = assetBalances[activeWallet!.publicKey]['claimed'];

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SafeArea(
          child: SingleChildScrollView(
            child: Column(
              children: [
                SizedBox(
                  height: height / 50,
                ),
                Container(
                  width: width,
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      SizedBox(
                        width: width / 20,
                      ),
                      Text(
                        "Swap",
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold,
                        ),
                      ),
                      SizedBox(
                        width: width / 10,
                      ),
                      Expanded(
                        child: DropdownButtonFormField(
                          isExpanded: true,
                          dropdownColor: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                          decoration: InputDecoration(
                            contentPadding: EdgeInsets.symmetric(
                                vertical: 0, horizontal: 20),
                            enabledBorder: OutlineInputBorder(
                              borderSide: BorderSide.none,
                              borderRadius: BorderRadius.circular(20),
                            ),
                            border: OutlineInputBorder(
                              borderSide: BorderSide.none,
                              borderRadius: BorderRadius.circular(20),
                            ),
                            filled: true,
                            fillColor: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                          ),
                          value: selectedWallet,
                          icon: Icon(
                            Icons.keyboard_arrow_down_rounded,
                            color: notifier.getbluewhitecolor,
                          ),
                          elevation: 0,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              fontWeight: FontWeight.w500),
                          onChanged: (newValue) {
                            setState(() {
                              selectedWallet = newValue!;
                              key1.currentState!.reset();
                              key2.currentState!.reset();
                              key3.currentState!.reset();
                              amount = 0;
                              textController.text = amount.toString();
                              sourceAsset = destinationAsset = null;
                              appState.activeWallet = wallets!.firstWhere(
                                  (wallet) => wallet.publicKey == newValue);
                            });
                          },
                          items: walletDropdownItems,
                        ),
                      ),
                      SizedBox(
                        width: width / 20,
                      ),
                    ],
                  ),
                ),
                Form(
                  key: formKey,
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 50,
                      ),
                      SizedBox(
                        height: height / 30,
                      ),
                      swap(),
                      SizedBox(
                        height: height / 50,
                      ),
                      CustomTextFormField.textField(
                        LanguageEn.amount,
                        notifier.getbluecolor,
                        Icons.currency_exchange,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        notifier.getgrey,
                        // dynamically change the size
                        // of the textbox so it will
                        // consistent when showing an
                        // error message
                        amountError ? 70.sp : 58.sp,
                        300.sp,
                        onChanged: (value) {
                          if (value != null && value.toString().isNotEmpty) {
                            setState(() {
                              amount = double.tryParse(value) ?? 0.0;
                            });
                          }
                        },
                        key: key3,
                        controller: textController,
                        // inputFormatters: [
                        //   doubleTypeFormatter(),
                        // ],
                        keyboardtype:
                            TextInputType.numberWithOptions(decimal: true),
                        validator: validateAmount,
                        onSaved: (value) =>
                            amount = value.trim().replaceAll(' ', ''),
                      ),
                      if (sourceAsset != null && !appState.hideBalances) ...[
                        availableBalance(),
                        SizedBox(
                          height: height / 50,
                        ),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: SfSliderTheme(
                            data: SfSliderThemeData(
                              activeTickColor: notifier.getdarkgrey,
                              inactiveTickColor: notifier.getdarkgrey,
                              activeMinorTickColor: notifier.getdarkgrey,
                              inactiveMinorTickColor: notifier.getdarkgrey,
                              inactiveLabelStyle: TextStyle(
                                  color: notifier.getdarkgrey,
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  fontWeight: FontWeight.w500),
                              activeLabelStyle: TextStyle(
                                  color: notifier.getdarkgrey,
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  fontWeight: FontWeight.w500),
                            ),
                            child: SfSlider(
                              min: 0,
                              max: 100,
                              value: sliderValue,
                              onChanged: _onSliderChanged,
                              interval: 25,
                              stepSize: 1,
                              inactiveColor: notifier.getdarkgrey,
                              showTicks: true,
                              tooltipTextFormatterCallback: _setToolTip,
                              showLabels: true,
                              enableTooltip: true,
                              minorTicksPerInterval: 1,
                            ),
                          ),
                        ),
                        const SizedBox(
                          height: 20.0,
                        ),
                      ],
                      SizedBox(
                        height: height / 20,
                      ),
                      Button(
                        LanguageEn.proceed,
                        notifier.getbluecolor,
                        wihitecolor,
                        onTap: () {
                          handleSubmit();
                        },
                      ),
                      SizedBox(
                        height: height / 20,
                      ),
                      Padding(
                          padding: EdgeInsets.only(
                              bottom:
                                  MediaQuery.of(context).viewInsets.bottom)),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget swap() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        // height: height / 2.5,
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.swapfrom,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                  Container(
                    width: width / 1.6,
                    child: DropdownButtonFormField<String>(
                      key: key1,
                      isExpanded: true,
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      decoration: InputDecoration(
                        contentPadding:
                            EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                          borderRadius: BorderRadius.circular(20),
                        ),
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                          borderRadius: BorderRadius.circular(20),
                        ),
                        filled: true,
                        fillColor: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                        errorStyle: TextStyle(
                          fontFamily: fontbody,
                          fontSize: 12,
                          overflow: TextOverflow.visible,
                        ),
                      ),
                      hint: Container(
                        width: 150, //and here
                        child: Text(
                          LanguageEn.chooseasset,
                          style: TextStyle(
                            color: sourceErr
                                ? Colors.red
                                : notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                          textAlign: TextAlign.end,
                        ),
                      ),
                      icon: Icon(
                        Icons.keyboard_arrow_down_rounded,
                        color:
                            sourceErr ? Colors.red : notifier.getbluewhitecolor,
                      ),
                      elevation: 0,
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          fontWeight: FontWeight.w500),
                      validator: (value) {
                        if (sourceAsset == null) {
                          setState(() {
                            sourceErr = true;
                          });
                          return '';
                        }

                        setState(() {
                          sourceErr = false;
                        });

                        return null;
                      },
                      onChanged: (newValue) {
                        print('this is new value: $newValue');
                        var splitNewValue = newValue!.split('|');
                        setState(() {
                          sourceAsset = claimedAssets.firstWhere(
                            (asset) =>
                                asset['assetIssuer'] == splitNewValue[0] &&
                                asset['assetCode'] == splitNewValue[1],
                          );
                          sourceErr = false;
                        });
                        print('this is source asset: $sourceAsset');
                      },
                      items:
                          dropdownItemBuilder(claimedAssets, destinationAsset),
                    ),
                  ),
                  SizedBox(height: height / 50),
                  SvgPicture.asset(
                    "assets/images/swapicon.svg",
                    // height: height / 40,
                  ),
                  SizedBox(height: height / 25),
                  Text(
                    LanguageEn.swapto,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                  Container(
                    width: width / 1.6,
                    child: DropdownButtonFormField<String>(
                      key: key2,
                      isExpanded: true,
                      dropdownColor: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                      decoration: InputDecoration(
                        contentPadding:
                            EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                        enabledBorder: OutlineInputBorder(
                          borderSide: BorderSide.none,
                          borderRadius: BorderRadius.circular(20),
                        ),
                        border: OutlineInputBorder(
                          borderSide: BorderSide.none,
                          borderRadius: BorderRadius.circular(20),
                        ),
                        filled: true,
                        fillColor: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                        errorStyle: TextStyle(
                          fontFamily: fontbody,
                          fontSize: 12,
                          overflow: TextOverflow.visible,
                        ),
                      ),
                      hint: Container(
                        width: 150, //and here
                        child: Text(
                          LanguageEn.chooseasset,
                          style: TextStyle(
                            color: destErr
                                ? Colors.red
                                : notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                          textAlign: TextAlign.end,
                        ),
                      ),
                      icon: Icon(
                        Icons.keyboard_arrow_down_rounded,
                        color:
                            destErr ? Colors.red : notifier.getbluewhitecolor,
                      ),
                      elevation: 0,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                      validator: (value) {
                        if (destinationAsset == null) {
                          setState(() {
                            destErr = true;
                          });
                          return '';
                        }

                        setState(() {
                          destErr = false;
                        });

                        return null;
                      },
                      onChanged: (newValue) {
                        print('this is new destination value: $newValue');
                        var splitNewValue = newValue!.split('|');
                        setState(() {
                          destinationAsset = claimedAssets.firstWhere(
                            (asset) =>
                                asset['assetIssuer'] == splitNewValue[0] &&
                                asset['assetCode'] == splitNewValue[1],
                          );
                          destErr = false;
                        });
                        print('this is destination asset: $sourceAsset');
                      },
                      items: dropdownItemBuilder(claimedAssets, sourceAsset),
                    ),
                  ),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void handleSubmit() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    submitForm();
  }

  submitForm() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "destinationAssetCode": destinationAsset['assetCode'] == 'XBN'
            ? ''
            : destinationAsset['assetCode'],
        "destinationAssetIssuer": destinationAsset['assetIssuer'],
        "sourceAssetCode":
            sourceAsset['assetCode'] == 'XBN' ? '' : sourceAsset['assetCode'],
        "sourceAssetIssuer": sourceAsset['assetIssuer'],
        "sourceAmount": amount.toStringAsFixed(4),
      };
      String requestBody = jsonEncode(map);
      print('this is request body $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/swap',
        body: requestBody,
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  postProcessData(messageShown, messageLength, data) {
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
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    // go to the definition of appState.viewData
    // to learn more about viewData
    appState.viewData![ConfirmSwapViewPageConfig.key] = data;

    appState.currentAction =
        PageAction(state: PageState.addPage, page: ConfirmSwapViewPageConfig);
  }

  _onSliderChanged(dynamic newValue) {
    setState(() {
      sliderValue = newValue;
      amount = (newValue / 100) * double.parse(sourceAsset['amount']);
      textController.text = amount.toStringAsFixed(4);
      textController.selection = TextSelection.fromPosition(
          TextPosition(offset: textController.text.length));
    });
  }

  String _setToolTip(dynamic actualValue, String formattedText) {
    actualValue = actualValue.round();
    return '$actualValue%';
  }

  List<DropdownMenuItem<String>> dropdownItemBuilder(assets, assetToSkip) {
    var assetList = assets;

    // filter assetList to remove the ones already selected
    if (assetToSkip != null) {
      assetList = assets
          .where((asset) => asset['assetIssuer'] != assetToSkip['assetIssuer'])
          .toList();
    }
    return assetList.map<DropdownMenuItem<String>>((asset) {
      return DropdownMenuItem<String>(
        value: '${asset['assetIssuer']}|${asset["assetCode"]}',
        child: Row(
          children: [
            CircleAvatar(
              maxRadius: 15,
              child: SvgPicture.asset(
                "assets/images/swapicon.svg",
                // height: height / 40,
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(8.0, 0, 0, 0),
              child: Text(
                asset["assetCode"].toString().isEmpty
                    ? 'XBN'
                    : asset["assetCode"],
                style: TextStyle(
                  fontSize: 15,
                  // fontWeight: FontWeight.bold,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                ),
              ),
            ),
          ],
        ),
      );
    }).toList();
  }

  String? validateAmount(String? value) {
    if (value!.isEmpty || double.tryParse(value)! <= 0) {
      return 'Please enter amount to swap';
    }

    if (double.tryParse(value) == null) {
      return 'Please enter a valid amount';
    }

    if (double.tryParse(value)! > (double.parse(sourceAsset['amount']) - 6)) {
      return 'You don\'t have sufficient balance';
    }
    return null;
  }

  Widget availableBalance() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 40.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Flexible(
            child: Text(
              amount.toString().isNotEmpty
                  ? "≈ ${amount.toStringAsFixed(4)} ${getAssetCode(sourceAsset['assetCode'])}"
                  : "≈ 0.0000 ${getAssetCode(sourceAsset['assetCode'])}",
              textScaleFactor: 1.0,
              style: TextStyle(
                  color: notifier.getdarkgrey,
                  fontWeight: FontWeight.w400,
                  fontSize: 12.0.sp),
            ),
          ),
          Flexible(
              child: Visibility(
            visible: true,
            replacement: Container(),
            child: Text(
              "${formatNumber(double.parse(sourceAsset['amount']))} ${getAssetCode(sourceAsset['assetCode'])}",
              textScaleFactor: 1.0,
              textAlign: TextAlign.right,
              style: TextStyle(color: notifier.getdarkgrey, fontSize: 12.0.sp),
            ),
          )),
        ],
      ),
    );
  }
}
