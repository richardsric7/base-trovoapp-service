import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_svg/svg.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
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
  // late List<Wallet> transactionableWallets;
  late Wallet wallet;
  late Asset asset;
  late List<Asset> claimedAssets;
  late RefreshController _refreshController;
  bool amountError = false;
  double amount = 0;
  Asset? sourceAsset = null;
  Asset? destinationAsset = null;
  String? sourceAssetRawDropdownValue;
  String? destinationAssetRawDropdownValue;
  String selectedWallet = '';
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  final formKey = GlobalKey<FormState>();
  bool sourceErr = false;
  bool destErr = false;
  final GlobalKey<FormFieldState> key1 = GlobalKey<FormFieldState>();
  final GlobalKey<FormFieldState> key2 = GlobalKey<FormFieldState>();
  final GlobalKey<FormFieldState> key3 = GlobalKey<FormFieldState>();
  final textController = TextEditingController();

  List<DropdownMenuItem<String>> walletDropdownItems(bool isSelected) {
    var walletsList = <DropdownMenuItem<String>>[];
    appState.userInfo!.transactionableWallets().forEach((wallet) {
      walletsList.add(
        DropdownMenuItem(
          child: Row(
            children: [
              Container(
                constraints: isSelected
                    ? BoxConstraints(maxWidth: width / 4)
                    : BoxConstraints(maxWidth: width / 2.5),
                child: Text(
                  wallet.alias!,
                  overflow:
                      isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
                ),
              ),
              if (wallet.isSharedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.people_outline,
                  size: 17,
                  color: notifier.getbluecolor,
                )
              ],
              if (!isSelected && wallet.publicKey == selectedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.check,
                  size: 18,
                  color: notifier.getbluecolor,
                )
              ],
            ],
          ),
          value: wallet.publicKey,
        ),
      );
    });

    return walletsList;
  }

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet = appState.primaryWallet;
    selectedWallet = appState.primaryWallet.publicKey!;
    claimedAssets = wallet.claimedAssets!;
    _refreshController = RefreshController(initialRefresh: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    return Scaffold(
      key: key,
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      drawer: getDrawer(context, appState, notifier),
      appBar: CustomAppBarWithoutLeading(
        context,
        notifier.getwihitecolor,
        height: height / 15,
        scaffoldKey: key,
        showMenu: true,
        txt: "swap".tr(),
        titlecolor: notifier.getbluewhitecolor,
      ).getBar(),
      body: SmartRefresher(
        enablePullDown: true,
        controller: _refreshController,
        onRefresh: refreshData,
        child: SafeArea(
          child: SingleChildScrollView(
            child: Column(
              children: [
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    SizedBox(
                      width: width / 15,
                    ),
                    Text(
                      "selectwallet".tr(),
                      style: TextStyle(
                          color: notifier.getbluewhitecolor,
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          fontWeight: FontWeight.w500),
                    ),
                    SizedBox(
                      width: width / 15,
                    ),
                    Expanded(
                        child: dropdown(
                      (newValue) {
                        selectedWallet = newValue.toString();
                        wallet = appState.userInfo!.getWallet(selectedWallet);
                        key1.currentState!.reset();
                        key2.currentState!.reset();
                        key3.currentState!.reset();
                        amount = 0;
                        textController.text = amount.toString();
                        sourceAsset = destinationAsset =
                            sourceAssetRawDropdownValue =
                                destinationAssetRawDropdownValue = null;

                        claimedAssets = wallet.claimedAssets!;
                        setState(() {});
                      },
                      walletDropdownItems(false),
                      selectedWallet.toString().isEmpty ? null : selectedWallet,
                      null,
                      context,
                      (context) {
                        return walletDropdownItems(true);
                      },
                    )),
                    SizedBox(
                      width: width / 15,
                    ),
                  ],
                ),
                Form(
                  key: formKey,
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 30,
                      ),
                      swap(),
                      SizedBox(
                        height: height / 50,
                      ),
                      CustomTextFormField.textField(
                        "amount".tr(),
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
                        70,
                        300,
                        onChanged: (value) {
                          if (value != null && value.toString().isNotEmpty) {
                            setState(() {
                              amount = double.tryParse(value) ?? 0.0;
                            });
                          }
                        },
                        key: key3,
                        controller: textController,
                        inputFormatters: [
                          FilteringTextInputFormatter.allow(RegExp(r'[0-9 \.]'))
                        ],
                        keyboardtype:
                            TextInputType.numberWithOptions(decimal: true),
                        autoFormatNumber: true,
                        validator: validateAmount,
                        onSaved: (value) => amount = value
                            .trim()
                            .replaceAll(' ', '')
                            .replaceAll(',', ''),
                      ),
                      if (sourceAsset != null && !appState.hideBalances) ...[
                        availableBalance(),
                        const SizedBox(
                          height: 20.0,
                        ),
                      ],
                      SizedBox(
                        height: height / 20,
                      ),
                      Button(
                        "proceed".tr(),
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
                    "swapfrom".tr(),
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
                      value: sourceAssetRawDropdownValue,
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
                          "chooseasset".tr(),
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
                        sourceAssetRawDropdownValue = newValue;
                        var splitNewValue = newValue!.split('|');
                        sourceAsset = wallet.claimedAssets!.firstWhere(
                          (asset) =>
                              asset.assetIssuer == splitNewValue[0] &&
                              asset.assetCode == splitNewValue[1],
                        );
                        sourceErr = false;
                        setState(() {});
                      },
                      items: dropdownItemBuilder(
                          claimedAssets, destinationAsset, false),
                    ),
                  ),
                  SizedBox(height: height / 50),
                  GestureDetector(
                    onTap: () {
                      setState(() {
                        if (destinationAssetRawDropdownValue != null) {
                          var splitNewValue =
                              destinationAssetRawDropdownValue!.split('|');
                          if (claimedAssets
                              .where((asset) =>
                                  asset.assetIssuer == splitNewValue[0] &&
                                  asset.assetCode == splitNewValue[1])
                              .isNotEmpty) {
                            var assetHolder = sourceAsset;
                            var rawValueHolder = sourceAssetRawDropdownValue;

                            sourceAsset = destinationAsset;
                            sourceAssetRawDropdownValue =
                                destinationAssetRawDropdownValue;
                            destinationAsset = assetHolder;
                            destinationAssetRawDropdownValue = rawValueHolder;
                          } else {
                            sourceAsset = destinationAsset =
                                sourceAssetRawDropdownValue =
                                    destinationAssetRawDropdownValue = null;
                          }

                          amount = 0;
                        }
                      });
                    },
                    child: Image.asset(
                      "assets/images/swap-rotated.png",
                      height: height / 20,
                    ),
                  ),
                  SizedBox(height: height / 25),
                  Text(
                    "swapto".tr(),
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
                      value: destinationAssetRawDropdownValue,
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
                          "chooseasset".tr(),
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
                        destinationAssetRawDropdownValue = newValue;
                        var splitNewValue = newValue!.split('|');
                        setState(() {
                          destinationAsset = claimedAssets.firstWhere(
                              (asset) =>
                                  asset.assetCode == splitNewValue[1] &&
                                  asset.assetIssuer == splitNewValue[0],
                              orElse: () {
                            // must be a curated swap item
                            return Asset(
                                assetCode: splitNewValue[1],
                                assetIssuer: splitNewValue[0]);
                          });
                          destErr = false;
                        });
                      },
                      items:
                          dropdownItemBuilder(claimedAssets, sourceAsset, true),
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
        "destinationAssetCode": destinationAsset!.assetCode,
        "destinationAssetIssuer": destinationAsset!.assetIssuer,
        "sourceAssetCode": sourceAsset!.assetCode,
        "sourceAssetIssuer": sourceAsset!.assetIssuer,
        "sourceAmount": amount.toStringAsFixed(4),
      };

      String requestBody = jsonEncode(map);

      Map responseData = await makePostRequest(
        uri: getEndpoint(),
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      // print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      // print(e);
      popup(context, title: "error".tr(), message: e.toString());
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
                // print('postProcessData: $messageShown'),
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    // go to the definition of appState.viewData
    // to learn more about viewData
    appState.viewData = {
      'transactionData': data,
      'walletPublicKey': wallet.publicKey,
      'sourceUsdPrice': sourceAsset!.usdPrice,
      'destinationUsdPrice': destinationAsset!.usdPrice
    };

    appState.currentAction =
        PageAction(state: PageState.addPage, page: ConfirmSwapViewPageConfig);
  }

  List<DropdownMenuItem<String>> dropdownItemBuilder(
      List<Asset> assets, Asset? assetToSkip, bool isDestination) {
    var assetsMap = {};
    List<DropdownMenuItem<String>> dropDownItems = [];

    if (isDestination) {
      // add the default assets to the list of destination assets
      appState.userInfo!.curatedSwapList!.forEach((asset) {
        assetsMap['${asset.assetIssuer}|${asset.assetCode}'] = asset.assetCode;
      });
    } else {
      assets.forEach((asset) {
        assetsMap['${asset.assetIssuer}|${asset.assetCode}'] = asset.assetCode;
      });
    }

    // remove the ones already selected as source or destination asset
    if (assetToSkip != null) {
      assetsMap.remove('${assetToSkip.assetIssuer}|${assetToSkip.assetCode}');
    }

    assetsMap.forEach((key, value) {
      dropDownItems.add(
        DropdownMenuItem<String>(
          // to make each asset in the list unique we combine both the assetCode
          // and the assetIssuer using '|' as the separator so we get something like
          // "asset|assetIssuer" as the value of each dropdown item
          value: key,
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
                  value.toString().isEmpty ? 'XBN' : value,
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
        ),
      );
    });

    return dropDownItems;
  }

  String getEndpoint() {
    // // if the wallet is share enabled and the user has initiator access
    // if (appState.transactionableWallets[selectedWallet]
    //             ['sharedAccessEnabled'] ==
    //         1 &&
    //     appState.transactionableWallets[selectedWallet]['threshold'] == 2) {
    //   return '/v1/shared-access/swap';
    // }
    return wallet.isSharedWalletAndCanInitiate
        ? '/v1/shared-access/swap'
        : '/v1/users/swap';
  }

  String? validateAmount(String? value) {
    if (value!.isEmpty || double.tryParse(value)! <= 0) {
      return "enteramounttoswap".tr();
    }

    if (double.tryParse(value) == null) {
      return "pleaseentervalidamount".tr();
    }

    if (sourceAsset == null) {
      return "chooseassettoswap".tr();
    }

    if (double.tryParse(value)! > (sourceAsset!.amount!)) {
      return "youdonthavesufficientbalance".tr();
    }

    if (getAssetCode(sourceAsset!.assetCode) == 'XBN' &&
        double.tryParse(value)! > (sourceAsset!.amount! - 6)) {
      return "youdonthavesufficientbalance".tr();
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
                  ? "≈ ${formatNumber(amount)} ${getAssetCode(sourceAsset!.assetCode)}"
                  : "≈ 0.0000 ${getAssetCode(sourceAsset!.assetCode)}",
              textScaleFactor: 1.0,
              style: TextStyle(
                  color: notifier.getdarkgrey,
                  fontWeight: FontWeight.w400,
                  fontSize: 12.0),
            ),
          ),
          Flexible(
              child: Visibility(
            visible: true,
            replacement: Container(),
            child: Text(
              "${formatNumber(sourceAsset!.amount!)} ${getAssetCode(sourceAsset!.assetCode)}",
              textScaleFactor: 1.0,
              textAlign: TextAlign.right,
              style: TextStyle(color: notifier.getdarkgrey, fontSize: 12.0),
            ),
          )),
        ],
      ),
    );
  }

  void refreshData() async {
    try {
      await appState.refreshData();
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }
}
