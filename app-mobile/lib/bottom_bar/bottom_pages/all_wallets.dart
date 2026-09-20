import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/constants.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/cache.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import 'package:trovo_app/widgets/wallet_slides.dart';

import '../../custom_bloc_observer/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

import 'package:local_auth/error_codes.dart' as auth_error;

class AllWalletsView extends StatefulWidget {
  const AllWalletsView({Key? key}) : super(key: key);

  @override
  State<AllWalletsView> createState() => _AllWalletsView();
}

enum WalletAction { import, createNew }

enum WalletView { listWallets, addSubWallet, confirmAddSubWallet }

class _AllWalletsView extends State<AllWalletsView>
    with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  WalletAction? action = WalletAction.createNew;
  final Authenticator _authenticator = Authenticator();
  bool isTileView = false;
  bool displaySearch = false;
  String searchText = '';
  bool isImport = true;
  String? tag;
  String? description;
  String? secretKey;
  int isAssetIssuerWallet = 0;
  final GlobalKey<ScaffoldState> key = GlobalKey(); // Create a key
  String password = '';
  late Account primaryWalletKeyPair;
  late Account newSubWalletKeyPair;
  late DataProvider appState;
  late UserInfo userInfo;
  final _formKey = GlobalKey<FormState>();
  final _formKey2 = GlobalKey<FormState>();
  late RefreshController _refreshController;
  String selectedWalletMode = "My wallets";
  List<String> walletListMode = ['My wallets', 'Shared wallets', 'All wallets'];
  late List<WalletTileColor> colors;
  late List<String> walletTypes = [
    'Standard',
    'Minting/Asset Tokenization',
    'Market Making/Trade',
    'Bulk Payment',
  ];
  final List<IconData> icons = [
    Icons.token_outlined,
    Icons.fire_truck_outlined,
    Icons.account_tree_outlined,
  ];
  int selectedWalletType = 0;

  List<DropdownMenuItem<String>> get walletTypeDropdownItems {
    var dropdownItems = walletTypes
        .map<DropdownMenuItem<String>>(
          (wallet) => DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [Text(wallet, overflow: TextOverflow.ellipsis)],
            ),
            value: walletTypes.indexOf(wallet).toString(),
          ),
        )
        .toList();

    return dropdownItems;
  }

  List<DropdownMenuItem<String>> get walletListModeDropdownItems {
    return walletListMode
        .map<DropdownMenuItem<String>>(
          (item) => DropdownMenuItem(
            child: Text(item, overflow: TextOverflow.ellipsis),
            value: item,
          ),
        )
        .toList();
  }

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    if ((appState.returnView != null && appState.returnView!.pages != null) &&
        appState.returnView!.pages!.contains(WalletPreparationViewPageConfig)) {
      appState.walletView.actionIcon = Icons.cancel_outlined;
      appState.walletView.actionText = "cancel".tr();
    }

    selectedWalletMode = appState.viewData?['walletMode'] ?? 'My wallets';
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    userInfo = appState.userInfo!;

    return Scaffold(
      key: key,
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      drawer: getDrawer(context, appState, notifier),
      appBar: CustomAppBar(
        context,
        notifier.getwihitecolor,
        "allwallets".tr(),
        notifier.getbluewhitecolor,
        height: height / 15,
      ).getBar(),
      body: SmartRefresher(
        enablePullDown: true,
        controller: _refreshController,
        onRefresh: refreshData,
        child: ListView(
          children: [
            Padding(
              padding: const EdgeInsets.all(10.0),
              child: Container(
                color: notifier.getfavorites,
                padding: EdgeInsets.all(8.sp),
                child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        Expanded(
                          child: DropdownButtonFormField(
                            isExpanded: true,
                            dropdownColor: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                            decoration: InputDecoration(
                              contentPadding: EdgeInsets.symmetric(
                                vertical: 0,
                                horizontal: 20,
                              ),
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
                            value: selectedWalletMode,
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
                            onChanged: (newValue) {
                              setState(() {
                                selectedWalletMode = newValue.toString();
                              });
                            },
                            items: walletListModeDropdownItems,
                          ),
                        ),
                        Container(
                          child: Card(
                            // shadowColor: Colors.black,
                            elevation: 0,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(10.0),
                            ),
                            color: notifier.isDark
                                ? notifier.getbluecolor90
                                : notifier.getaddsubwalletgrey,
                            child: TextButton(
                              onPressed: () => setState(() {
                                displaySearch = !displaySearch;
                                setState(() {
                                  searchText = '';
                                });
                              }),
                              child: Icon(
                                Icons.search,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: 5),
                    if (displaySearch) ...[
                      CustomTextFormField.textFieldWithoutIcon(
                        'searchwallets'.tr(),
                        notifier.getbluecolor,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        notifier.getgrey,
                        40.sp,
                        width,
                        onChanged: (value) {
                          if (value != null && value.toString().isNotEmpty) {
                            setState(() {
                              searchText = value;
                            });
                          }
                        },
                        keyboardtype: TextInputType.text,
                      ),
                    ],
                  ],
                ),
              ),
            ),
            walletListView(),
            Padding(
              padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget walletTile(
    walletName,
    usdBal,
    preferredFiatBal,
    WalletTileColor color, {
    isShared = false,
  }) {
    return Card(
      elevation: 5,
      shadowColor: Colors.black,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(15.0)),
      color: color.backColor,
      child: Stack(
        children: [
          Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(
                vertical: 35.0,
                horizontal: 20,
              ),
              child: Image.asset(
                'assets/images/trovo_white.png',
                height: 100,
                width: 100,
                color: Color(0x3CFFFFFF),
              ),
            ),
          ),
          Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Container(
                  width: width / 3,
                  child: Wrap(
                    alignment: WrapAlignment.center,
                    children: [
                      Text(
                        walletName,
                        textAlign: TextAlign.center,
                        overflow: TextOverflow.visible,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: color.foreColor,
                        ),
                      ),
                      if (isShared) ...[
                        SizedBox(width: width / 90),
                        Icon(
                          Icons.people_alt_outlined,
                          color: color.foreColor,
                          size: 20,
                        ),
                      ],
                    ],
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(
                    vertical: 10,
                    horizontal: 10,
                  ),
                  child: Text(
                    appState.hideBalances ? hideBalanceText : preferredFiatBal,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontbody,
                      color: color.foreColor,
                    ),
                  ),
                ),
                if (appState.defaultCurrency != 'USD') ...[
                  Padding(
                    padding: const EdgeInsets.symmetric(
                      vertical: 8,
                      horizontal: 10,
                    ),
                    child: Text(
                      appState.hideBalances ? hideBalanceText : usdBal,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontbody,
                        color: color.foreColor,
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget walletListView() {
    colors = [
      notifier.getstructuredbluecolor,
      notifier.getstructuredbluecolor90,
      notifier.getstructuredbluecolor80,
      notifier.getstructuredbluecolor70,
      notifier.getstructuredbluecolor60,
      notifier.getstructuredbluecolor50,
      notifier.getstructuredgreencolor,
      notifier.getstructuredgreencolor90,
      notifier.getstructuredgreencolor80,
      notifier.getstructuredgreencolor70,
      notifier.getstructuredgreencolor60,
      notifier.getstructuredgreencolor50,
      notifier.getorangecolor,
      notifier.getorangecolor90,
      notifier.getorangecolor80,
      notifier.getorangecolor70,
      notifier.getorangecolor60,
      notifier.getorangecolor50,
      notifier.getpinkcolor,
      notifier.getpinkcolor90,
      notifier.getpinkcolor80,
      notifier.getpinkcolor70,
      notifier.getpinkcolor60,
      notifier.getpinkcolor50,
    ];
    return Column(
      children: getWallets(
        displaySearch && searchText.length > 0
            ? userInfo.allWallets
                  .where((wallet) => wallet.alias!.contains(searchText))
                  .toList()
            : switch (selectedWalletMode) {
                    'My wallets' => userInfo.mySolelyOwnedWallets,
                    'Shared wallets' => userInfo.sharedWallets,
                    'All wallets' => userInfo.allWallets,
                    _ => [],
                  } ??
                  [],
        false,
      ),
    );
  }

  Widget walletListItem(
    walletName,
    int walletType,
    balanceUsd,
    preferredFiatBal,
    assetCount,
    WalletTileColor color, {
    isShared = false,
  }) {
    return Container(
      height: height / 5.8,
      margin: EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        borderRadius: const BorderRadius.all(Radius.circular(20.0)),
        color: color.backColor,
      ),
      child: Stack(
        alignment: AlignmentDirectional.topStart,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Container(
                child: Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 20.0,
                    vertical: 20.0,
                  ),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Container(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          crossAxisAlignment: CrossAxisAlignment.end,
                          children: [
                            if (assetCount != null) ...[
                              Container(
                                child: Text(
                                  '${assetCount} ${'assetplural'.tr()}',
                                  style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w600,
                                    color: color.foreColor,
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                              ),
                            ],
                            SizedBox(height: 7),
                            Image.asset(
                              'assets/images/trovo_white.png',
                              height: 40,
                              width: 40,
                            ),
                            SizedBox(height: 7),
                            Container(
                              width: 40,
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  if (isShared) ...[
                                    SizedBox(width: 2),
                                    Icon(
                                      Icons.people_outline,
                                      size: 17,
                                      color: color.foreColor,
                                    ),
                                  ],
                                  if (walletType != 0) ...[
                                    SizedBox(width: 2),
                                    Icon(
                                      icons[walletType - 1],
                                      size: 17,
                                      color: color.foreColor,
                                    ),
                                  ] else ...[
                                    SizedBox(height: 10),
                                  ],
                                  if (walletName.contains('-distribution')) ...[
                                    Icon(
                                      icons[2],
                                      size: 17,
                                      color: color.foreColor,
                                    ),
                                  ],
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
          Container(
            child: Padding(
              padding: const EdgeInsets.symmetric(
                horizontal: 20.0,
                vertical: 20.0,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Container(
                            constraints: BoxConstraints(maxWidth: width / 2.5),
                            child: Text(
                              walletName,
                              style: TextStyle(
                                fontSize: 16,
                                color: color.foreColor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          SizedBox(width: 10),
                          IconButton(
                            padding: EdgeInsets.zero,
                            color: color.foreColor,
                            constraints: BoxConstraints(),
                            onPressed: () => {
                              Clipboard.setData(
                                ClipboardData(text: walletName),
                              ),
                              showSnackBar("walletalias".tr(), context),
                            },
                            icon: Icon(Icons.copy, fill: 1.0, size: 15),
                          ),
                        ],
                      ),
                    ],
                  ),
                  SizedBox(height: height / 25),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        appState.hideBalances
                            ? hideBalanceText
                            : preferredFiatBal,
                        style: TextStyle(
                          fontSize: 20,
                          color: color.foreColor,
                          fontFamily: fontbody,
                        ),
                      ),
                      SizedBox(height: height / 50),
                    ],
                  ),
                  if (appState.defaultCurrency != 'USD') ...[
                    SizedBox(height: height / 80),
                    Text(
                      appState.hideBalances ? hideBalanceText : balanceUsd,
                      style: TextStyle(
                        fontWeight: FontWeight.w300,
                        fontSize: 13,
                        color: color.foreColor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  // show the add subwallet view as if its a new page
  // the app's back button dispatcher has been overriden to make this page
  // behave as if is a new separate page when you press the back button
  Widget addSubwallet() {
    return Form(
      key: _formKey2,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Card(
              shadowColor: Colors.black,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(15.0),
              ),
              color: notifier.isDark
                  ? notifier.getbluecolor90
                  : notifier.getaddsubwalletgrey,
              child: Center(
                child: Column(
                  children: [
                    SizedBox(height: height / 50),
                    Container(
                      width: width / 1.4,
                      child: Text(
                        "abouttocreatesubwallet".tr(),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Text(
                      "chooseamethod".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Row(
                      children: [
                        SizedBox(width: width / 10),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.import,
                            groupValue: action,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            onChanged: (value) => {
                              setState(() {
                                action = value;
                              }),
                            },
                          ),
                        ),
                        Text(
                          "importexistingwallet".tr(),
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    Row(
                      children: [
                        SizedBox(width: width / 10),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.createNew,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateColor.resolveWith(
                              (states) => notifier.getbluewhitecolor,
                            ),
                            groupValue: action,
                            onChanged: (value) => {
                              setState(() {
                                action = value;
                              }),
                            },
                          ),
                        ),
                        Text(
                          "createnewwallet".tr(),
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 50),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Container(
                              width: width / 1.3,
                              decoration: BoxDecoration(
                                border: Border.all(
                                  color: notifier.getbluecolor,
                                ),
                                borderRadius: const BorderRadius.all(
                                  Radius.circular(15.0),
                                ),
                              ),
                              child: dropdown(
                                (newValue) async {
                                  selectedWalletType = int.parse(
                                    newValue.toString(),
                                  );
                                },
                                walletTypeDropdownItems,
                                selectedWalletType.toString(),
                                null,
                                context,
                                null,
                              ),
                            ),
                          ],
                        ),
                        // SizedBox(
                        //   height: height / 50,
                        // ),
                        // Row(
                        //   children: [
                        //     SizedBox(
                        //       width: width / 4.4,
                        //     ),
                        //     GestureDetector(
                        //       onTap: () {
                        //         mintWalletExplainerPopup(context);
                        //       },
                        //       child: Text(
                        //         'What does it mean?',
                        //         style: TextStyle(
                        //           decoration: TextDecoration.underline,
                        //           color: notifier.getbluewhitecolor,
                        //           fontSize: 12.sp,
                        //           fontWeight: FontWeight.w500,
                        //           fontFamily: fontbody,
                        //         ),
                        //       ),
                        //     ),
                        //   ],
                        // ),
                        SizedBox(height: height / 50),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
          SizedBox(height: height / 30),
          // Tag name
          CustomTextFormField.textField(
            "tag".tr(),
            notifier.getbluecolor,
            Icons.tag,
            notifier.getgrey,
            notifier.getbluewhitecolor,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            initialValue: tag,
            onChanged: (value) {
              setState(() {
                tag = value.trim().replaceAll(' ', '');
              });
            },
            onSaved: (value) {
              tag = value.trim().replaceAll(' ', '');
            },
            keyboardtype: TextInputType.text,
            maxLength: 12,
            validator: validateTag,
            helperText: tag == null || tag!.isEmpty
                ? ''
                : "${appState.userInfo!.username}_$tag",
          ),
          SizedBox(height: height / 50),
          CustomTextFormField.textField(
            "description".tr(),
            notifier.getbluecolor,
            Icons.description,
            notifier.getgrey,
            notifier.getbluewhitecolor,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            initialValue: description,
            onSaved: (value) {
              description = value;
            },
            keyboardtype: TextInputType.text,
            maxLength: 100,
            validator: validateDescription,
          ),
          if (action == WalletAction.import) ...[
            SizedBox(height: height / 50),
            // Secret Key
            CustomPasswordFormField(
              "secretkey".tr(),
              notifier.getbluecolor,
              Icons.lock,
              notifier.getgrey,
              notifier.getbluewhitecolor,
              notifier.getblck,
              70.sp,
              300.sp,
              validator: (value) {
                var trimmedVal = value!.trim().replaceAll(' ', '');
                if (trimmedVal.isEmpty) {
                  return "entersecretkeyempty".tr();
                }

                if (trimmedVal.length < 64) {
                  return "secretkeyinvalid".tr();
                }

                try {
                  TrovoWalletSDK().parseSecretKey(value);
                } catch (e) {
                  return 'Secret Key is invalid';
                }

                return null;
              },
              onSaved: (value) {
                secretKey = value!.trim().replaceAll(' ', '');
              },
              maxLength: 66,
            ),
          ],
          SizedBox(height: height / 30),
          Button(
            "continuee".tr(),
            notifier.getbluecolor,
            wihitecolor,
            onTap: () {},
          ),
          SizedBox(height: height / 20),
        ],
      ),
    );
  }

  // show the confirm add subwallet view as if its a new page
  // the app's back button dispatcher has been overriden to make this page
  // behave as if is a separate page when you press the back button
  Widget confirmAddSubwallet() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Card(
            shadowColor: Colors.black,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(15.0),
            ),
            color: notifier.isDark
                ? notifier.getbluecolor90
                : notifier.getaddsubwalletgrey,
            child: Center(
              child: Form(
                key: _formKey,
                child: Column(
                  children: [
                    SizedBox(height: height / 50),
                    Container(
                      width: width / 1.4,
                      child: Text(
                        "requesttocreatesubwallet".tr(),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Text(
                      "tag".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Text(
                      "${userInfo.username!}_$tag",
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Text(
                      "description".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Text(
                      description ?? '',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Text(
                      "method".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Text(
                      action == WalletAction.import
                          ? "importsubwallet".tr()
                          : "createnewsubwallet".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Text(
                      "wallettype".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Text(
                      walletTypes[selectedWalletType],
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(height: height / 50),
                    Text(
                      "publickey".tr(),
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 15.0),
                      child: Text(
                        newSubWalletKeyPair.address,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(height: height / 20),
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 15.0),
                      child: Text(
                        "willattractcharges".tr(),
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(height: height / 30),
                  ],
                ),
              ),
            ),
          ),
        ),
        SizedBox(height: height / 50),
        // Secret Key
        CustomPasswordFormField(
          "password".tr(),
          notifier.getbluecolor,
          Icons.lock,
          notifier.getgrey,
          notifier.getbluewhitecolor,
          notifier.getblck,
          70.sp,
          300.sp,
          validator: validatePassword,
          textInputAction: TextInputAction.done,
          onChanged: (value) {
            setState(() {
              password = value!.trim().replaceAll(' ', '');
            });
          },
          onSaved: (value) {
            secretKey = value!.trim().replaceAll(' ', '');
          },
        ),
        SizedBox(height: height / 30),

        if (appState.biometricEnabled && password.isEmpty) ...[
          Button(
            "authorizewithbiometrics".tr(),
            notifier.getbluecolor,
            wihitecolor,
            onTap: toggleSwitch,
          ),
        ] else ...[
          Button(
            "authorize".tr(),
            notifier.getbluecolor,
            wihitecolor,
            onTap: handleAuthorization,
          ),
        ],
        SizedBox(height: height / 20),
      ],
    );
  }

  List<Widget> getWallets(List<Wallet> filteredWallets, isTileMode) {
    return [
      if (filteredWallets.isNotEmpty) ...[
        for (var i = 0; i < filteredWallets.length; i++) ...[
          GestureDetector(
            onTap: () {
              appState.viewData = {'walletAddress': filteredWallets[i].address};
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: WalletDetailsViewPageConfig,
              );
            },
            child: isTileMode
                ? walletTile(
                    filteredWallets[i].alias!.capitalizeFirst!,
                    '${totalAccountBalanceInUSD(appState, filteredWallets[i].claimedAssets!)} USD',
                    '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, filteredWallets[i].claimedAssets!)} ${appState.defaultCurrency}',
                    i % 2 == 0
                        ? colors[((i + 1) % colors.length)]
                        : colors[((i) % colors.length)],
                  )
                : Column(
                    children: [
                      Container(
                        constraints: BoxConstraints(maxHeight: height / 5.8),
                        child: WalletSlide(
                          alias: filteredWallets[i].alias!.capitalizeFirst!,
                          walletType: filteredWallets[i].walletType ?? 0,
                          fiatBalance:
                              '${totalAccountBalanceInUSD(appState, filteredWallets[i].claimedAssets!)} USD',
                          totalBalance:
                              '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, filteredWallets[i].claimedAssets!)} ${appState.defaultCurrency}',
                          assetCount: filteredWallets[i].claimedAssets!.length
                              .toString(),
                          backColor:
                              colors[((i + 1) % colors.length)].backColor,
                          foreColor:
                              colors[((i + 1) % colors.length)].foreColor,
                          isSharedWallet: filteredWallets[i].isSharedWallet,
                          initialHiddenState: appState.hideBalances,
                        ),
                      ),
                      SizedBox(height: height / 50),
                    ],
                  ),
          ),
        ],
      ] else ...[
        Container(
          height: height / 2,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(
                "nowalletshere".tr(),
                style: TextStyle(
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ],
          ),
        ),
      ],

      // shared wallets
    ];
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return "pleaseenteryourpassword".tr();

    if (value.length < 6) return "use6charsormoreforpassword".tr();

    return null;
  }

  generateKeyPairs() {
    primaryWalletKeyPair = TrovoWalletSDK().parseSecretKey(
      appState.secretKeys[0],
    );
    setState(() {
      if (action == WalletAction.import) {
        try {
          // parse supplied secret to get the keypair
          newSubWalletKeyPair = TrovoWalletSDK().parseSecretKey(secretKey);
        } catch (e) {
          popup(context, title: "error".tr(), message: "invalidsecretkey".tr());
        }
      } else {
        // generate keypair for the new subwallet
        newSubWalletKeyPair = TrovoWalletSDK().createAccount();
      }
    });
  }

  String? validateTag(String? value) {
    if (value!.isEmpty) return "enterwallettag".tr();

    String pattern = r'^[a-zA-Z0-9\_]*$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return "invalidtagname".tr();
    }

    return null;
  }

  String? validateDescription(String? value) {
    if (value!.isEmpty) return "enterwalletdesc".tr();

    return null;
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer();
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  void handleAuthorization() {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer();
    } else {
      popup(context, title: "oops".tr(), message: "invalidpassword".tr());
    }
  }

  void sendDataToServer() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "publickey": newSubWalletKeyPair.address,
        "walletTag": tag,
        "WalletDescription": description,
        "walletType": selectedWalletType,
      };
      String requestBody = jsonEncode(map);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.address,
        secretKey: primaryWalletKeyPair.secretKey,
        address: primaryWalletKeyPair.address,
      );

      if (responseData['statusCode'] == 200) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        await postProcessData(
          messageShown,
          messageLength,
          responseData['data'],
        );
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
    }
    hideLoader(context);
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
        () => {postProcessData(messageShown, messageLength, data)},
      );

      messageShown++;
      return;
    }

    sendFullDataToServer(data);
    return;
  }

  void sendFullDataToServer(responseBody) async {
    showLoader(context);

    try {
      // get primary signature
      var primarySignature = TrovoWalletSDK().signBase64Txn(
        primaryWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      // get secondary signature
      var subWalletSignature = TrovoWalletSDK().signBase64Txn(
        newSubWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      responseBody['primarySignature'] = primarySignature;
      responseBody['subWalletSignature'] = subWalletSignature;

      String requestBody = jsonEncode(responseBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.address,
        secretKey: primaryWalletKeyPair.secretKey,
        address: primaryWalletKeyPair.address,
      );

      if (responseData['statusCode'] == 200) {
        // add the secret key of this new subwallet to
        // the existing list of secrets
        appState.secretKeys.add(newSubWalletKeyPair.secretKey);
        // store back the list of secret keys but this time it
        // contains the secret key of the newly created subwallet
        await StoreData().storeInsertData('secretKey', appState.secretKeys);
        await updateUserInfo(
          appState.primaryWallet.address,
          appState.secretKeys[0],
          appState.primaryWallet.address,
          userInfo.username,
          appState,
          forceRefresh: true,
        );
        // add the new subwallet to appState and
        // set the newly created subwallet as the activeWallet
        appState.activeWallet = appState.userInfo!.wallets!.firstWhere(
          (wallet) => wallet.address == newSubWalletKeyPair.address,
        );
        appState.activeWallet!.secretKey = newSubWalletKeyPair.secretKey;
        // move to next page
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: CongratulationsPageConfig,
        );
        resetForm();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['error'],
        );
      }
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
    }

    hideLoader(context);
  }

  void resetForm() {
    tag = '';
    description = '';
    secretKey = '';
    appState.walletView.actionIcon = Icons.add_circle_outline_sharp;
    appState.walletView.actionText = "addsubwallet".tr();
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
