import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/constants.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtom_password.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class Wallets extends StatefulWidget {
  const Wallets({Key? key}) : super(key: key);

  @override
  State<Wallets> createState() => _WalletsState();
}

enum WalletAction { import, createNew }

enum WalletView { listWallets, addSubWallet, confirmAddSubWallet }

class _WalletsState extends State<Wallets> with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  WalletAction? action = WalletAction.createNew;
  final Authenticator _authenticator = Authenticator();
  bool isTileView = false;
  bool isImport = true;
  String? tag;
  String? description;
  String? secretKey;
  int isAssetIssuerWallet = 0;
  String password = '';
  late Account primaryWalletKeyPair;
  late Account newSubWalletKeyPair;
  late DataProvider appState;
  late UserInfo userInfo;
  final _formKey = GlobalKey<FormState>();
  final _formKey2 = GlobalKey<FormState>();
  late RefreshController _refreshController;
  String selectedWalletMode = "My wallets";
  List<String> walletListMode = [
    'My wallets',
    'Shared wallets',
  ];
  late List<WalletTileColor> colors;
  late List<String> walletTypes = [
    'Standard',
    'Minting/Asset Tokenization',
    'Market Making',
    'Bulk Payment'
  ];
  int selectedWalletType = 0;

  List<DropdownMenuItem<String>> get walletTypeDropdownItems {
    var dropdownItems = walletTypes
        .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  wallet,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
            value: walletTypes.indexOf(wallet).toString()))
        .toList();

    return dropdownItems;
  }

  List<DropdownMenuItem<String>> get accessModeDropdownItems {
    return walletListMode
        .map<DropdownMenuItem<String>>((item) => DropdownMenuItem(
            child: Text(
              item,
              overflow: TextOverflow.ellipsis,
            ),
            value: item))
        .toList();
  }

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    userInfo = appState.userInfo!;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
          resizeToAvoidBottomInset: false,
          backgroundColor: notifier.getwihitecolor,
          appBar: AppBar(
            centerTitle: true,
            title: Text(
              LanguageEn.wallets,
              style: TextStyle(
                  color: notifier.getblck,
                  fontWeight: FontWeight.bold,
                  fontFamily: fontsemibold),
            ),
            backgroundColor: notifier.getfavorites,
            elevation: 0,
          ),
          body: SmartRefresher(
            enablePullDown: true,
            controller: _refreshController,
            onRefresh: refreshData,
            child: ListView(
              children: [
                if (appState.walletView.view == WalletView.listWallets) ...[
                  Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: Container(
                        color: notifier.getfavorites,
                        padding: EdgeInsets.all(8.sp),
                        child: Row(
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
                                value: selectedWalletMode,
                                icon: Icon(
                                  Icons.keyboard_arrow_down_rounded,
                                  color: notifier.getbluewhitecolor,
                                ),
                                elevation: 0,
                                style: TextStyle(
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 15.sp,
                                    fontFamily: fontsemibold,
                                    fontWeight: FontWeight.w500),
                                onChanged: (newValue) {
                                  setState(() {
                                    selectedWalletMode = newValue.toString();
                                  });
                                },
                                items: accessModeDropdownItems,
                              ),
                            ),
                            SizedBox(
                              width: width / 15,
                            ),
                            GestureDetector(
                              onTap: () => setState(() {
                                isTileView = !isTileView;
                              }),
                              child: Padding(
                                padding: const EdgeInsets.all(8.0),
                                child: SvgPicture.asset(
                                  isTileView
                                      ? "assets/images/listview.svg"
                                      : "assets/images/tileview.svg",
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ),
                          ],
                        )),
                  )
                ] else ...[
                  SizedBox(
                    height: height / 40,
                  )
                ],
                GestureDetector(
                  onTap: () {
                    appState.viewData = {
                      'walletPublicKey': appState.primaryWallet.publicKey
                    };
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: WalletDetailsViewPageConfig);
                  },
                  child: walletListItem(
                    appState.primaryWallet.alias!.capitalizeFirst,
                    '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, appState.primaryWallet.claimedAssets!)} USD',
                    '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, appState.primaryWallet.claimedAssets!)} ${appState.defaultCurrency}',
                    notifier.getstructuredbluecolor,
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                GestureDetector(
                  onTap: () {
                    if (appState.walletView.view == WalletView.listWallets) {
                      appState.walletView.actionIcon = Icons.cancel_outlined;
                      appState.walletView.actionText = LanguageEn.cancel;
                      appState.walletView.view = WalletView.addSubWallet;
                    } else if (appState.walletView.view ==
                        WalletView.addSubWallet) {
                      appState.walletView.actionIcon =
                          Icons.add_circle_outline_sharp;
                      appState.walletView.actionText = LanguageEn.addsubwallet;
                      appState.walletView.view = WalletView.listWallets;
                      resetForm();
                    } else if (appState.walletView.view ==
                        WalletView.confirmAddSubWallet) {
                      appState.walletView.actionIcon = Icons.cancel_outlined;
                      appState.walletView.actionText = LanguageEn.cancel;
                      appState.walletView.view = WalletView.addSubWallet;
                    }

                    setState(() {});
                  },
                  child: Column(
                    children: [
                      Icon(
                        appState.walletView.actionIcon,
                        color: notifier.getbluewhitecolor,
                        size: 35,
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Text(
                        appState.walletView.actionText,
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w500,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                if (appState.walletView.view == WalletView.addSubWallet) ...[
                  addSubwallet()
                ] else if (appState.walletView.view ==
                    WalletView.confirmAddSubWallet) ...[
                  confirmAddSubwallet()
                ] else ...[
                  isTileView ? gridView() : walletListView(),
                ],
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          )),
    );
  }

  Widget gridView() {
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

    return Container(
      height: height / 2,
      child: GridView.count(
        primary: true,
        padding: const EdgeInsets.fromLTRB(15, 20, 15, 70),
        crossAxisCount: 2,
        mainAxisSpacing: 20,
        crossAxisSpacing: 20,
        childAspectRatio: 1.05,
        children: selectedWalletMode == 'My wallets'
            ? getWallets(userInfo.wallets!, true)
            : getSharedWallets(userInfo.sharedWallets!, true),
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
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      color: color.backColor,
      child: Stack(
        children: [
          Center(
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
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
                  padding:
                      const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
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
                    padding:
                        const EdgeInsets.symmetric(vertical: 8, horizontal: 10),
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
      children: selectedWalletMode == 'My wallets'
          ? getWallets(userInfo.wallets!, false)
          : getSharedWallets(userInfo.sharedWallets!, false),
    );
  }

  Widget walletListItem(
    walletName,
    balanceUsd,
    preferredFiatBal,
    WalletTileColor color, {
    isShared = false,
  }) {
    return Container(
      // height: height / 6.6,
      margin: EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        borderRadius: const BorderRadius.all(Radius.circular(20.0)),
        color: color.backColor,
      ),
      child: Stack(
        alignment: AlignmentDirectional.centerEnd,
        children: [
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Image.asset(
                    'assets/images/trovo_white.png',
                    height: 80,
                    width: 80,
                  ),
                  SizedBox(
                    width: width / 20,
                  ),
                ],
              ),
            ],
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 20.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 2,
                  child: Wrap(
                    children: [
                      Text(
                        walletName,
                        style: TextStyle(
                          fontSize: 16,
                          color: color.foreColor,
                          fontFamily: fontsemibold,
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
                SizedBox(
                  height: height / 50,
                ),
                Row(
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
                  ],
                ),
                SizedBox(height: height / 80),
                if (appState.defaultCurrency != 'USD') ...[
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
                    SizedBox(
                      height: height / 50,
                    ),
                    Container(
                      width: width / 1.4,
                      child: Text(
                        LanguageEn.abouttocreatesubwallet,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      LanguageEn.chooseamethod,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      children: [
                        SizedBox(
                          width: width / 10,
                        ),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.import,
                            groupValue: action,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor),
                            onChanged: (value) => {
                              setState(
                                () {
                                  action = value;
                                },
                              )
                            },
                          ),
                        ),
                        Text(
                          LanguageEn.importexistingwallet,
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
                        SizedBox(
                          width: width / 10,
                        ),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.createNew,
                            activeColor: notifier.getbluewhitecolor,
                            fillColor: MaterialStateColor.resolveWith(
                                (states) => notifier.getbluewhitecolor),
                            groupValue: action,
                            onChanged: (value) => {
                              setState(
                                () {
                                  action = value;
                                },
                              )
                            },
                          ),
                        ),
                        Text(
                          LanguageEn.createnewwallet,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Container(
                              width: width / 1.7,
                              decoration: BoxDecoration(
                                border: Border.all(
                                  color: notifier.getbluecolor,
                                ),
                                borderRadius: const BorderRadius.all(
                                    Radius.circular(15.0)),
                              ),
                              child: dropdown(
                                (newValue) async {
                                  selectedWalletType =
                                      int.parse(newValue.toString());
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
                        SizedBox(
                          height: height / 50,
                        ),
                      ],
                    )
                  ],
                ),
              ),
            ),
          ),
          SizedBox(
            height: height / 30,
          ),
          // Tag name
          CustomTextFormField.textField(
            LanguageEn.tag,
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
              print('tag: $value');
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
            LanguageEn.description,
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
              print('description: $value');
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
              LanguageEn.secretkey,
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
                  return LanguageEn.entersecretkeyempty;
                }

                if (trimmedVal.length < 56) {
                  return LanguageEn.secretkeyinvalid;
                }

                try {
                  TrovoWalletSDK().parseSecretKey(value);
                } catch (e) {
                  return 'Secret Key is invalid';
                }

                return null;
              },
              onSaved: (value) {
                print('email: $value');
                secretKey = value!.trim().replaceAll(' ', '');
              },
              maxLength: 56,
            ),
          ],
          SizedBox(height: height / 30),
          Button(
            LanguageEn.continuee,
            notifier.getbluecolor,
            wihitecolor,
            onTap: () => submitForm(),
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
    return Column(children: [
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
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.4,
                    child: Text(
                      LanguageEn.requesttocreatesubwallet,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.tag,
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
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.description,
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
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.method,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Text(
                    action == WalletAction.import
                        ? LanguageEn.importsubwallet
                        : LanguageEn.createnewsubwallet,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    'Wallet type',
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
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.publickey,
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
                      newSubWalletKeyPair.publicKey,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 20,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 15.0),
                    child: Text(
                      'Please note that completing this process will attract some charges.',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 13,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 30,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
      SizedBox(height: height / 50),
      // Secret Key
      CustomPasswordFormField(
        LanguageEn.password,
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
        // onSubmitted: (value) {
        //   print('email: $value');
        //   secretKey = value!.trim().replaceAll(' ', '');
        // },
        onSaved: (value) {
          print('email: $value');
          secretKey = value!.trim().replaceAll(' ', '');
        },
      ),
      SizedBox(height: height / 30),

      if (appState.biometricEnabled && password.isEmpty) ...[
        Button(
          LanguageEn.authorizewithbiometrics,
          notifier.getbluecolor,
          wihitecolor,
          onTap: toggleSwitch,
        ),
      ] else ...[
        Button(
          LanguageEn.authorize,
          notifier.getbluecolor,
          wihitecolor,
          onTap: handleAuthorization,
        ),
      ],
      SizedBox(height: height / 20),
    ]);
  }

  List<Widget> getWallets(
    List<Wallet> wallets,
    isTileMode,
  ) {
    List<Wallet> filteredWallets = wallets
        .where((wallet) => !wallet.isSharedWallet && !wallet.isPrimaryWallet)
        .toList();
    return [
      for (var i = 0; i < filteredWallets.length; i++) ...[
        GestureDetector(
            onTap: () {
              appState.viewData = {
                'walletPublicKey': filteredWallets[i].publicKey
              };
              appState.currentAction = PageAction(
                  state: PageState.addPage, page: WalletDetailsViewPageConfig);
            },
            child: isTileMode
                ? walletTile(
                    filteredWallets[i].alias!.capitalizeFirst!,
                    '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, filteredWallets[i].claimedAssets!)} USD',
                    '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, filteredWallets[i].claimedAssets!)} ${appState.defaultCurrency}',
                    i % 2 == 0
                        ? colors[((i + 1) % colors.length)]
                        : colors[((i) % colors.length)],
                  )
                : Column(children: [
                    walletListItem(
                      filteredWallets[i].alias!.capitalizeFirst!,
                      '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, filteredWallets[i].claimedAssets!)} USD',
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, filteredWallets[i].claimedAssets!)} ${appState.defaultCurrency}',
                      colors[((i + 1) % colors.length)],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                  ])),
      ]
      // shared wallets
    ];
  }

  List<Widget> getSharedWallets(List<Wallet> sharedWallets, isTileMode) {
    var walletTiles = <Widget>[];
    int i = 0;
    sharedWallets.forEach((wallet) {
      walletTiles.add(
        GestureDetector(
          onTap: () {
            appState.viewData = {'walletPublicKey': wallet.publicKey};
            appState.currentAction = PageAction(
                state: PageState.addPage, page: WalletDetailsViewPageConfig);
          },
          child: isTileMode
              ? walletTile(
                  wallet.alias.toString().capitalizeFirst!,
                  '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallet.claimedAssets!)} USD',
                  '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallet.claimedAssets!)} ${appState.defaultCurrency}',
                  i % 2 == 0
                      ? colors[((i + 1) % colors.length)]
                      : colors[((i) % colors.length)],
                  isShared: true,
                )
              : Column(
                  children: [
                    walletListItem(
                      wallet.alias.toString().capitalizeFirst!,
                      '${getTotalFiatBalanceOfAllAssetsInWallet('USD', appState, wallet.claimedAssets!)} USD',
                      '${getTotalFiatBalanceOfAllAssetsInWallet(appState.defaultCurrency, appState, wallet.claimedAssets!)} ${appState.defaultCurrency}',
                      i % 2 == 0
                          ? colors[((i + 1) % colors.length)]
                          : colors[((i) % colors.length)],
                      isShared: true,
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                  ],
                ),
        ),
      );
      i++;
    });

    return walletTiles;
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  submitForm() async {
    print('submitting...');
    var form = _formKey2.currentState;
    if (!form!.validate()) {
      return;
    }
    form.save();
    generateKeyPairs();
    setState(() {
      appState.walletView.view = WalletView.confirmAddSubWallet;
      appState.walletView.actionIcon = Icons.arrow_circle_left_outlined;
      appState.walletView.actionText = LanguageEn.back;
      password = '';
      secretKey = '';
    });
  }

  generateKeyPairs() {
    primaryWalletKeyPair =
        TrovoWalletSDK().parseSecretKey(appState.secretKeys[0]);
    setState(() {
      if (action == WalletAction.import) {
        try {
          // parse supplied secret to get the keypair
          newSubWalletKeyPair = TrovoWalletSDK().parseSecretKey(secretKey);
        } catch (e) {
          popup(context, title: 'Error!', message: 'Secret Key is invalid');
        }
      } else {
        // generate keypair for the new subwallet
        newSubWalletKeyPair = TrovoWalletSDK().createAccount();
      }
    });
  }

  String? validateTag(String? value) {
    print('validating tag...');
    if (value!.isEmpty) return 'Enter wallet tag';

    String pattern = r'^[a-zA-Z0-9\_]*$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return 'Invalid tag name';
    }

    return null;
  }

  String? validateDescription(String? value) {
    print('validating description...');
    if (value!.isEmpty) return 'Enter wallet description';

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
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }

  void sendDataToServer() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "publickey": newSubWalletKeyPair.publicKey,
        "walletTag": tag,
        "WalletDescription": description,
        // "assetIssuerWallet": isAssetIssuerWallet,
        "walletType": selectedWalletType,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        print('new dialog $messageLength');
        print('messagecount $messageShown');
        await postProcessData(
            messageShown, messageLength, responseData['data']);
        // print('sending full data to server.........');
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
    print('messageShown: $messageShown messageLength $messageLength');
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

    sendFullDataToServer(data);
    return;
  }

  void sendFullDataToServer(responseBody) async {
    try {
      showLoader(context);
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

      print('this is primary sign: $primarySignature');
      print('this is subwallet sign: $subWalletSignature');

      responseBody['primarySignature'] = primarySignature;
      responseBody['subWalletSignature'] = subWalletSignature;

      print(responseBody);

      String requestBody = jsonEncode(responseBody);

      print('this is request body: $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        // add the secret key of this new subwallet to
        // the existing list of secrets
        appState.secretKeys.add(newSubWalletKeyPair.secretKey);
        // store back the list of secret keys but this time it
        // contains the secret key of the newly created subwallet
        await StoreData().storeInsertData('secretKey', appState.secretKeys);
        await updateUserInfo();
        // add the new subwallet to appState and
        // set the newly created subwallet as the activeWallet
        appState.activeWallet = appState.userInfo!.wallets!.firstWhere(
            (wallet) => wallet.publicKey == newSubWalletKeyPair.publicKey);
        appState.activeWallet!.secretKey = newSubWalletKeyPair.secretKey;
        // move to next page
        appState.currentAction = PageAction(
            state: PageState.addPage, page: CongratulationsPageConfig);
        resetForm();
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['error']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }

    hideLoader(context);
  }

  void resetForm() {
    tag = '';
    description = '';
    secretKey = '';
    appState.walletView.actionIcon = Icons.add_circle_outline_sharp;
    appState.walletView.actionText = LanguageEn.addsubwallet;
    appState.walletView.view = WalletView.listWallets;
  }

  Future<void> updateUserInfo() async {
    var keyPair =
        TrovoWalletSDK().parseSecretKey(primaryWalletKeyPair.secretKey);
    Map responseData = await makeGetRequest(
        uri: '/v1/users/${userInfo.username!.trim().replaceAll(' ', '')}',
        signer: keyPair.publicKey,
        publicKey: keyPair.publicKey,
        secretKey: keyPair.secretKey);

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200) {
      await storeUserInfo(responseData['data'], appState);
    }
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
