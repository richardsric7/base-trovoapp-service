import 'dart:convert';
import 'dart:math';
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
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
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
  TextEditingController _utf8TextController = TextEditingController();
  TextEditingController toController = TextEditingController();
  TextEditingController sendingWalletController = TextEditingController();
  final amountController = TextEditingController();
  bool isSharedWallet = false;
  dynamic selectedNetwork = '';
  late Future<Map> fetchNetworksFuture;
  int index = 0;

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
              value: '${item['network']}|$index'),
        )
        .toList();
  }

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    viewData = appState.viewData![WithdrawAssetViewPageConfig.key];
    fetchNetworksFuture = fetchNetworks();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    viewData = appState.viewData![WithdrawAssetViewPageConfig.key];
    isSharedWallet = viewData['walletInfo']['sharedAccessEnabled'] == 1;
    if (selectedNetwork.toString().isNotEmpty) {
      index = int.parse(selectedNetwork.split('|')[1]);
    }

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
                              LanguageEn.somethingwentwrong,
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            ElevatedButton(
                              onPressed: () {
                                setState(() {
                                  fetchNetworksFuture = fetchNetworks();
                                });
                              },
                              style: ButtonStyle(
                                backgroundColor:
                                    MaterialStateProperty.all<Color>(
                                        notifier.getbluecolor!),
                              ),
                              child: Text(
                                LanguageEn.retry,
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                ),
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
                          SizedBox(
                            height: height / 40,
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
                            validator: validateDropdown,
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
                    if (selectedNetwork.toString().isNotEmpty) ...[
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
                                myKeyValueRow('Network fee',
                                    '${networks[index]['withdrawFee']} ${viewData['assetCode']}'),
                                SizedBox(
                                  height: height / 90,
                                ),
                                myKeyValueRow('Min',
                                    '${networks[index]['withdrawMin']} ${viewData['assetCode']}'),
                                SizedBox(
                                  height: height / 90,
                                ),
                                myKeyValueRow('Max',
                                    '${networks[index]['withdrawMax']} ${viewData['assetCode']}'),
                                SizedBox(
                                  height: height / 90,
                                ),
                                myKeyValueRow('Estimated arrival time',
                                    '${networks[index]['estimatedArrivalTime'].toString()} min(s)'),
                              ],
                            ),
                          ),
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                    ],
                    GestureDetector(
                      child: CustomTextFormField.textField(
                        "${LanguageEn.withdraw} ${LanguageEn.to}",
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
                    // SizedBox(
                    //   height: height / 50,
                    // ),
                    // Padding(
                    //   padding: const EdgeInsets.symmetric(
                    //     horizontal: 5,
                    //   ),
                    //   child: Container(
                    //     decoration: BoxDecoration(
                    //       borderRadius:
                    //           const BorderRadius.all(Radius.circular(15.0)),
                    //       color: notifier.isDark
                    //           ? darktilewhitecolor
                    //           : notifier.getaddsubwalletgrey,
                    //     ),
                    //     child: Padding(
                    //       padding: const EdgeInsets.all(20.0),
                    //       child: Column(
                    //         children: [
                    //           Row(
                    //             mainAxisAlignment:
                    //                 MainAxisAlignment.spaceBetween,
                    //             children: [
                    //               Text(
                    //                 'Total',
                    //                 style: TextStyle(
                    //                     fontSize: 13,
                    //                     color: notifier.getbluewhitecolor,
                    //                     fontFamily: fontsemibold),
                    //               ),
                    //               Text(
                    //                 '0.0000 ${viewData['assetCode']}',
                    //                 style: TextStyle(
                    //                     fontSize: 13,
                    //                     color: notifier.getbluewhitecolor,
                    //                     fontFamily: fontbody),
                    //               ),
                    //             ],
                    //           ),
                    //           SizedBox(
                    //             height: height / 90,
                    //           ),
                    //         ],
                    //       ),
                    //     ),
                    //   ),
                    // ),
                    SizedBox(height: height / 20),
                  ],
                ),
              )),
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
              fontFamily: fontsemibold),
        ),
        Text(
          value,
          style: TextStyle(
              fontSize: 13,
              color: notifier.getbluewhitecolor,
              fontFamily: fontbody),
        ),
      ],
    );
  }

  Future<Map> fetchNetworks() async {
    try {
      Map responseData = await makeGetRequest(
        uri: '/v1/crypto/withdrawal-networks/${viewData['assetCode']}',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: viewData['walletInfo']['publicKey'],
      );

      print(responseData['data']);

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  String? validateMemo(String? value) {
    if (value!.isEmpty) return null;

    RegExp regex = new RegExp(networks[index]['memoRegex']);
    if (!regex.hasMatch(value)) return 'Text is too long or has invalid chars';
  }

  String? validateDropdown(String? _) {
    return selectedNetwork.toString().isNotEmpty
        ? null
        : 'Please select network';
  }

  String? validateTo(String? value) {
    RegExp regex = new RegExp(networks[index]['addressRegex']);

    // reciever cannot be empty
    if (value!.isEmpty) return 'Please enter destination address';

    if (!regex.hasMatch(value.trim().replaceAll(' ', '')))
      return 'Invalid destination address';

    return null;
  }

  String? validateAmount(String? value) {
    var minValue = networks[index]['withdrawMin'];
    var maxValue = networks[index]['withdrawMax'];
    if (value!.isEmpty) {
      setState(() {
        amountError = true;
      });
      return 'Please enter amount to send';
    }

    if (double.tryParse(value) == null) {
      setState(() {
        amountError = true;
      });
      return 'Please enter a valid amount';
    }

    if (double.tryParse(value)! < double.parse(minValue)) {
      setState(() {
        amountError = true;
      });
      return 'Value less than min withdrawable';
    }

    if (double.tryParse(value)! > double.parse(maxValue)) {
      setState(() {
        amountError = true;
      });
      return 'Value greater than max withdrawable';
    }

    if (double.tryParse(value)! > (double.parse(viewData['amount']))) {
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
        "currency": viewData['assetCode'].toString(),
        "amountSubmitted": double.parse(amount),
        "withdrawalAddress": to,
        "withdrawalNetwork": selectedNetwork.toString().split('|')[0],
        // "withdrawalMemo": memo,
        "withdrawalServiceFee": serviceFee,
        "withdrawalNetworkFee": double.parse(networks[index]['withdrawFee']),
      };

      String requestBody = jsonEncode(map);
      print('===============> map: $map');

      Map responseData = await makePostRequest(
        uri: isSharedWallet
            ? '/v1/shared-access/crypto/withdrawals'
            : '/v1/crypto/withdrawals',
        body: requestBody,
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: viewData['walletInfo']['publicKey'],
      );

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
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    appState.viewData![ConfirmWithdrawViewPageConfig.key] =
        appState.viewData![WithdrawAssetViewPageConfig.key];
    appState.viewData![ConfirmWithdrawViewPageConfig.key]['data'] = data;

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

class _Utf8LengthLimitingTextInputFormatter extends TextInputFormatter {
  _Utf8LengthLimitingTextInputFormatter(this.maxLength)
      : assert(maxLength == null || maxLength == -1 || maxLength > 0);

  final int maxLength;

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    if (maxLength != null &&
        maxLength > 0 &&
        bytesLength(newValue.text) > maxLength) {
      // If already at the maximum and tried to enter even more, keep the old value.
      if (bytesLength(oldValue.text) == maxLength) {
        return oldValue;
      }
      return truncate(newValue, maxLength);
    }
    return newValue;
  }

  static TextEditingValue truncate(TextEditingValue value, int maxLength) {
    var newValue = '';
    if (bytesLength(value.text) > maxLength) {
      var length = 0;

      value.text.characters.takeWhile((char) {
        var nbBytes = bytesLength(char);
        if (length + nbBytes <= maxLength) {
          newValue += char;
          length += nbBytes;
          return true;
        }
        return false;
      });
    }
    return TextEditingValue(
      text: newValue,
      selection: value.selection.copyWith(
        baseOffset: min(value.selection.start, newValue.length),
        extentOffset: min(value.selection.end, newValue.length),
      ),
      composing: TextRange.empty,
    );
  }

  static int bytesLength(String value) {
    return utf8.encode(value).length;
  }
}
