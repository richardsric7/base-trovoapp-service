import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:collection/collection.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  TextEditingController toController = TextEditingController();
  TextEditingController sendingWalletController = TextEditingController();
  final amountController = TextEditingController();
  dynamic selectedNetwork = '';
  late Future<Map> fetchNetworksFuture;
  int index = 0;
  late Wallet wallet;
  late Asset? asset;

  List<dynamic> networks = [];
  double serviceFee = 0;

  List<DropdownMenuItem<String>> get networksDropdownItems {
    return networks
        .mapIndexed<DropdownMenuItem<String>>(
          (index, item) => DropdownMenuItem(
            child: Text(
              item['name'].toString(),
              overflow: TextOverflow.ellipsis,
            ),
            value: '${item['network']}|$index',
          ),
        )
        .toList();
  }

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);

    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );

    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    fetchNetworksFuture = fetchNetworks();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    if (selectedNetwork.toString().isNotEmpty) {
      index = int.parse(selectedNetwork.split('|')[1]);
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '${"withdraw".tr()} ${getAssetCode(asset!.assetCode)}',
                      style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 50),
              FutureBuilder<Map>(
                future: fetchNetworksFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return Container(
                      width: width / 1.2,
                      height: height / 1.7,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          CircularProgressIndicator(
                            backgroundColor: notifier.getbluecolor,
                            valueColor: new AlwaysStoppedAnimation<Color>(
                              notifier.getgreencolor,
                            ),
                            strokeWidth: 3.0,
                          ),
                        ],
                      ),
                    );
                  } else if (snapshot.connectionState == ConnectionState.done) {
                    if (snapshot.hasError) {
                      return Container(
                        width: width / 1.2,
                        height: height / 1.7,
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              "somethingwentwrong".tr(),
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 16,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody,
                              ),
                            ),
                            ElevatedButton(
                              onPressed: () {
                                setState(() {
                                  fetchNetworksFuture = fetchNetworks();
                                });
                              },
                              style: ButtonStyle(
                                backgroundColor: WidgetStateProperty.all<Color>(
                                  notifier.getbluecolor!,
                                ),
                                foregroundColor: WidgetStateProperty.all<Color>(
                                  notifier.getwihitecolor,
                                ),
                              ),
                              child: Text(
                                "retry".tr(),
                                style: TextStyle(fontFamily: fontsemibold),
                              ),
                            ),
                          ],
                        ),
                      );
                    } else if (snapshot.hasData) {
                      networks = snapshot.data!['networks'];
                      serviceFee = double.parse(snapshot.data!['serviceFee']);
                      return Column(
                        children: [
                          formFields(),
                          SizedBox(height: height / 40),
                          Button(
                            "proceed".tr(),
                            notifier.getbluecolor,
                            wihitecolor,
                            onTap: () {
                              handleSubmit();
                            },
                          ),
                          SizedBox(height: height / 20),
                        ],
                      );
                    }
                  }

                  return Text(
                    '',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontsemibold,
                    ),
                  );
                },
              ),
              Padding(
                padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom,
                ),
              ),
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
            double.tryParse(amount) != null
                ? "≈ ${formatNumber(double.parse(amount))} ${getAssetCode(asset!.assetCode)}"
                : "≈ 0.0000 ${getAssetCode(asset!.assetCode)}",
            textScaleFactor: 1.0,
            style: TextStyle(
              color: notifier.getdarkgrey,
              fontWeight: FontWeight.w400,
              fontSize: 12.0.sp,
            ),
          ),
        ),
        Flexible(
          child: Visibility(
            visible: true,
            replacement: Container(),
            child: Text(
              "${formatNumber(asset!.amount!)} ${getAssetCode(asset!.assetCode)}",
              textScaleFactor: 1.0,
              textAlign: TextAlign.right,
              style: TextStyle(color: notifier.getdarkgrey, fontSize: 12.0.sp),
            ),
          ),
        ),
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
                  SizedBox(height: height / 50),
                  Text(
                    "network".tr(),
                    style: TextStyle(
                      fontSize: 15,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                  SizedBox(height: height / 50),
                  Row(
                    children: [
                      Expanded(
                        child: DropdownButtonFormField(
                          isExpanded: true,
                          validator: validateDropdown,
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
                            "selectnetwork".tr(),
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
                            fontWeight: FontWeight.w500,
                          ),
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
                  SizedBox(height: height / 50),
                  if (selectedNetwork.toString().isNotEmpty) ...[
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 5),
                      child: Container(
                        decoration: BoxDecoration(
                          borderRadius: const BorderRadius.all(
                            Radius.circular(15.0),
                          ),
                          color: notifier.isDark
                              ? darktilewhitecolor
                              : notifier.getaddsubwalletgrey,
                        ),
                        child: Padding(
                          padding: const EdgeInsets.all(20.0),
                          child: Column(
                            children: [
                              myKeyValueRow(
                                "networkfee".tr(),
                                '${networks[index]['withdrawFee']} ${asset!.assetCode}',
                              ),
                              SizedBox(height: height / 90),
                              myKeyValueRow(
                                "min".tr(),
                                '${networks[index]['withdrawMin']} ${asset!.assetCode}',
                              ),
                              SizedBox(height: height / 90),
                              myKeyValueRow(
                                "max".tr(),
                                '${networks[index]['withdrawMax']} ${asset!.assetCode}',
                              ),
                              SizedBox(height: height / 90),
                              myKeyValueRow(
                                "eta".tr(),
                                '${networks[index]['estimatedArrivalTime'].toString()} min(s)',
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                    SizedBox(height: height / 50),
                  ],
                  GestureDetector(
                    child: CustomTextFormField.textField(
                      "${"withdraw".tr()} ${"to".tr()}",
                      notifier.getbluecolor,
                      Icons.send,
                      notifier.getgrey,
                      notifier.getprefixicon,
                      notifier.getblck,
                      notifier.getgrey,
                      80.sp,
                      300.sp,
                      controller: toController,
                      validator: validateTo,
                      onSaved: (value) => to = value.trim().replaceAll(' ', ''),
                    ),
                  ),
                  SizedBox(height: height / 50),
                  CustomTextFormField.textField(
                    "amount".tr(),
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
                    autoFormatNumber: true,
                    keyboardtype: TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    validator: validateAmount,
                    onSaved: (value) =>
                        amount = value.trim().replaceAll(' ', ''),
                    inputFormatters: [
                      FilteringTextInputFormatter.allow(RegExp(r'[0-9 \.]')),
                    ],
                  ),
                  if (!appState.hideBalances) ...[availableBalance()],
                  SizedBox(height: height / 20),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget myKeyValueRow(String key, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          key,
          style: TextStyle(
            fontSize: 13,
            color: notifier.getbluewhitecolor,
            fontFamily: fontsemibold,
          ),
        ),
        Text(
          value,
          style: TextStyle(
            fontSize: 13,
            color: notifier.getbluewhitecolor,
            fontFamily: fontbody,
          ),
        ),
      ],
    );
  }

  Future<Map> fetchNetworks() async {
    try {
      Map responseData = await makeGetRequest(
        uri: '/v1/crypto/withdrawal-networks/${asset!.assetCode}',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  String? validateDropdown(String? _) {
    return selectedNetwork.toString().isNotEmpty
        ? null
        : 'Please select network';
  }

  String? validateTo(String? value) {
    RegExp regex = new RegExp(networks[index]['addressRegex']);

    // reciever cannot be empty
    if (value!.isEmpty) return "enterdestinationaddress".tr();

    if (!regex.hasMatch(value.trim().replaceAll(' ', '')))
      return "invaliddestinationaddress".tr();

    return null;
  }

  String? validateAmount(String? value) {
    var minValue = networks[index]['withdrawMin'];
    var maxValue = networks[index]['withdrawMax'];
    if (value!.isEmpty) {
      setState(() {
        amountError = true;
      });
      return "pleaseenteramounttosend".tr();
    }

    if (double.tryParse(value) == null) {
      setState(() {
        amountError = true;
      });
      return "pleaseentervalidamount".tr();
    }

    if (double.tryParse(value)! < double.parse(minValue)) {
      setState(() {
        amountError = true;
      });
      return "valuelessthanwithdrawable".tr();
    }

    if (double.tryParse(value)! > double.parse(maxValue)) {
      setState(() {
        amountError = true;
      });
      return "valuegreaterthanwithdrawable".tr();
    }

    if (double.tryParse(value)! > (asset!.amount!)) {
      setState(() {
        amountError = true;
      });
      return "youdonthavesufficientbalance".tr();
    }

    setState(() {
      amountError = false;
    });
    return null;
  }

  void handleSubmit() {
    final form = formKey.currentState;
    if (!form!.validate()) {
      return;
    }

    form.save();
    submit();
  }

  submit() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      Map map = {
        "currency": asset!.assetCode.toString(),
        "amountSubmitted": double.parse(amount),
        "withdrawalAddress": to,
        "withdrawalNetwork": selectedNetwork.toString().split('|')[0],
        "withdrawalServiceFee": serviceFee,
        "withdrawalNetworkFee": double.parse(networks[index]['withdrawFee']),
      };

      String requestBody = jsonEncode(map);

      Map responseData = await makePostRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/crypto/withdrawals'
            : '/v1/crypto/withdrawals',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        appState.viewData = {
          'transactionData': responseData['data'],
          'withdrawalNetworkName': networks[index]['name'],
          'walletPublicKey': wallet.publicKey,
          'assetCode': asset!.assetCode,
          'assetIssuer': asset!.assetIssuer,
        };

        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: ConfirmWithdrawViewPageConfig,
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
  }
}
