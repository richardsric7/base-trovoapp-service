import 'package:flutter/material.dart';
import 'package:collection/collection.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WithdrawAsset extends StatefulWidget {
  const WithdrawAsset({Key? key}) : super(key: key);

  @override
  State<WithdrawAsset> createState() => _WithdrawAsset();
}

class _WithdrawAsset extends State<WithdrawAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  String to = ''; // the reciever
  String amount = '';
  bool amountError = false;
  String? memo;
  var viewData;
  TextEditingController toController = TextEditingController();
  TextEditingController sendingWalletController = TextEditingController();
  final amountController = TextEditingController();
  bool isSharedWallet = false;
  dynamic selectedNetwork = '';

  List<dynamic> networks = [];

  List<DropdownMenuItem<String>> get networksDropdownItems {
    return networks
        .mapIndexed<DropdownMenuItem<String>>(
          (index, item) => DropdownMenuItem(
              child: Text(
                item['network'].toString(),
                overflow: TextOverflow.ellipsis,
              ),
              value: '${item['depositAddress']}|$index'),
        )
        .toList();
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    print(appState.viewData![WithdrawAssetViewPageConfig.key]['data']);
    networks =
        appState.viewData![WithdrawAssetViewPageConfig.key]['data'].toList();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    viewData = appState.viewData![WithdrawAssetViewPageConfig.key];
    isSharedWallet = viewData['walletInfo']['sharedAccessEnabled'] == 1;

    print('viewData ${viewData}');

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(height / 15),
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            leading: Navigator.canPop(context)
                ? GestureDetector(
                    onTap: () {
                      Navigator.of(context).pop();
                    },
                    child: Image.asset("assets/images/back.png", scale: 5),
                  )
                : null,
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '${LanguageEn.withdraw} ${getAssetCode(viewData['assetCode'])}',
                      style: TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    )
                  ],
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              formFields(),
              SizedBox(
                height: height / 10,
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
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Widget availableBalance() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Flexible(
          child: Text(
            amount.isNotEmpty
                ? "≈ ${formatNumber(double.parse(amount))} ${getAssetCode(viewData['assetCode'])}"
                : "≈ 0.0000 ${getAssetCode(viewData['assetCode'])}",
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
            "${formatNumber(double.parse(viewData['amount']))} ${getAssetCode(viewData['assetCode'])}",
            textScaleFactor: 1.0,
            textAlign: TextAlign.right,
            style: TextStyle(color: notifier.getdarkgrey, fontSize: 12.0.sp),
          ),
        )),
      ],
    );
  }

  Widget formFields() {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
              width: 300.sp,
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              ),
              child: Form(
                key: formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      'Network',
                      style: TextStyle(
                          fontSize: 15,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
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
                            hint: Text(
                              'Select network',
                              style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody,
                              ),
                              textAlign: TextAlign.end,
                            ),
                            icon: Icon(
                              Icons.keyboard_arrow_down_rounded,
                              color: notifier.getbluewhitecolor,
                            ),
                            elevation: 0,
                            style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontSize: 15,
                                fontFamily: fontbody,
                                fontWeight: FontWeight.w500),
                            onChanged: (newValue) {
                              setState(() {
                                selectedNetwork = newValue!;
                              });
                            },
                            items: networksDropdownItems,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    GestureDetector(
                      child: CustomTextFormField.textField(
                        "${LanguageEn.withdraw} ${LanguageEn.to}",
                        notifier.getbluecolor,
                        Icons.send,
                        notifier.getgrey,
                        notifier.getprefixicon,
                        notifier.getblck,
                        notifier.getgrey,
                        70.sp,
                        300.sp,
                        controller: toController,
                        validator: validateTo,
                        onSaved: (value) =>
                            to = value.trim().replaceAll(' ', ''),
                      ),
                    ),
                    SizedBox(height: height / 50),
                    CustomTextFormField.textField(
                      LanguageEn.amount,
                      notifier.getbluecolor,
                      Icons.currency_exchange,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      onChanged: (value) {
                        setState(() {
                          amount = trim(value.toString(), '.');
                        });
                      },
                      controller: amountController,
                      keyboardtype:
                          TextInputType.numberWithOptions(decimal: true),
                      validator: validateAmount,
                      onSaved: (value) =>
                          amount = value.trim().replaceAll(' ', ''),
                      inputFormatters: [
                        FilteringTextInputFormatter.allow(RegExp(r'[0-9 \.]'))
                      ],
                    ),
                    if (!appState.hideBalances) ...[availableBalance()],
                    SizedBox(height: height / 30),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Flexible(
                          child: Text(
                            "+1% fee (0.000 USDC)",
                            textScaleFactor: 1.0,
                            style: TextStyle(
                                color: notifier.getdarkgrey,
                                fontWeight: FontWeight.w400,
                                fontSize: 12.0.sp),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Padding(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 5,
                      ),
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius:
                              const BorderRadius.all(Radius.circular(15.0)),
                          color: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                        ),
                        child: Padding(
                          padding: const EdgeInsets.all(20.0),
                          child: Column(
                            children: [
                              Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    'Total',
                                    style: TextStyle(
                                        fontSize: 15,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold),
                                  ),
                                  Text(
                                    '0.0000 USDC',
                                    style: TextStyle(
                                        fontSize: 15,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold),
                                  ),
                                ],
                              ),
                              SizedBox(
                                height: height / 90,
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.end,
                                children: [
                                  Text(
                                    '\$0.0000',
                                    style: TextStyle(
                                        fontSize: 13,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    SizedBox(height: height / 20),
                  ],
                ),
              )),
        ),
      ],
    );
  }

  String? validateTo(String? value) {
    // reciever cannot be empty
    if (value!.isEmpty) return 'Please enter reciever username or public key';

    if (value.length < 3) return 'Invalid username or public key';

    return null;
  }

  String? validateAmount(String? value) {
    if (value!.isEmpty) {
      setState(() {
        amountError = true;
      });
      'Please enter amount to send';
    }

    if (double.tryParse(value) == null) {
      setState(() {
        amountError = true;
      });
      return 'Please enter a valid amount';
    }

    if (double.tryParse(value)! <= 0) {
      setState(() {
        amountError = true;
      });
      return 'Value must be greater than 0';
    }

    if (double.tryParse(value)! > (double.parse(viewData['amount']) - 6)) {
      setState(() {
        amountError = true;
      });
      return 'You don\'t have sufficient balance';
    }

    setState(() {
      amountError = false;
    });
    return null;
  }

  void handleSubmit() {
    appState.viewData![ConfirmWithdrawViewPageConfig.key] =
        appState.viewData![WithdrawAssetViewPageConfig.key];

    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: ConfirmWithdrawViewPageConfig,
    );
  }

  @override
  void dispose() {
    super.dispose();
    viewData?['deepLinkInfo'] = null;
  }
}
