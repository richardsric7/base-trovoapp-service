import 'dart:convert';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
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

class RequestEarlyExitView extends StatefulWidget {
  const RequestEarlyExitView({Key? key}) : super(key: key);

  @override
  State<RequestEarlyExitView> createState() => _RequestEarlyExitView();
}

class _RequestEarlyExitView extends State<RequestEarlyExitView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late Asset? asset;
  final formKey = GlobalKey<FormState>();
  String to = ''; // the reciever
  String amount = '';
  String accountNumber = '';
  bool amountError = false;
  String preferredPaymentMethod = 'Bank';
  String? memo;
  var deeplinkInfo;
  TextEditingController _utf8TextController = TextEditingController();
  TextEditingController toController = TextEditingController();
  TextEditingController sendingWalletController = TextEditingController();
  final amountController = TextEditingController();
  final FocusNode focusNode = FocusNode();

  List<DropdownMenuItem<String>> get getPaymentMethods {
    List<DropdownMenuItem<String>> options = [];
    ['Bank', 'Wallet'].forEach((item) {
      options.add(
        DropdownMenuItem(
          child: Text(
            item.replaceAll(' ', ''),
            overflow: TextOverflow.ellipsis,
          ),
          value: item,
        ),
      );
    });
    return options;
  }

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (appState.viewData![SendAssetViewPageConfig.key]?['deepLinkInfo'] ==
            null &&
        appState.viewData!['walletPublicKey'] != null) {
      wallet = appState.userInfo!.getWallet(
        appState.viewData!['walletPublicKey'],
      );
      asset = wallet.claimedAssets!.firstWhere(
        (asset) =>
            asset.assetCode == appState.viewData!['assetCode'] &&
            asset.assetIssuer == appState.viewData!['assetIssuer'],
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    if (appState.viewData![SendAssetViewPageConfig.key] != null &&
        appState.viewData![SendAssetViewPageConfig.key]['deepLinkInfo'] !=
            null) {
      deeplinkInfo =
          appState.viewData![SendAssetViewPageConfig.key]['deepLinkInfo'];
      toController.text = deeplinkInfo['receiver'];
      amountController.text = deeplinkInfo['amount'];
      _utf8TextController.text = deeplinkInfo['memo'];
      wallet = appState.userInfo!.getWallet(deeplinkInfo['sendingWallet']);
      asset = wallet.claimedAssets!.firstWhere(
        (asset) =>
            asset.assetCode == deeplinkInfo['assetCode'] &&
            asset.assetIssuer == deeplinkInfo['assetIssuer'],
      );
      appState.viewData![SendAssetViewPageConfig.key]['deepLinkInfo'] = null;
    }

    sendingWalletController.text = wallet.alias!;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          '',
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Request Early Exit (${getAssetCode(asset!.assetCode)})',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: height / 50),
              formFields(),
              SizedBox(height: height / 10),
              Button(
                "proceed".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.setPage(page: EarlyExitSummaryViewPageConfig);
                },
              ),
              SizedBox(height: height / 20),
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
            amount.isNotEmpty
                ? "≈ ${formatNumber(double.tryParse(amount) ?? 0.0)} ${getAssetCode(asset!.assetCode)}"
                : "≈ 0.0000 ${getAssetCode(asset!.assetCode)}",
            textScaler: TextScaler.linear(1.0),
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
              textScaler: TextScaler.linear(1.0),
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
                children: [
                  Row(
                    children: [
                      Text(
                        "Number of tokens to exit with",
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 70),
                  CustomTextFormField.textField(
                    "",
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    notifier.getprefixicon,
                    notifier.getblck,
                    notifier.getgrey,
                    60.sp,
                    300.sp,
                    onChanged: (value) {
                      setState(() {
                        var splitText = value.split('.');
                        if (splitText.length > 2) {
                          // remove all dots except the first one.
                          value = '${splitText[0]}.${splitText[1]}';
                        }

                        amount = trim(value.toString(), '.');
                      });
                    },
                    controller: amountController,
                    readOnly:
                        deeplinkInfo != null &&
                        deeplinkInfo['amount'].toString().isNotEmpty,
                    autoFormatNumber: true,
                    keyboardtype: TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    validator: validateAmount,
                    onSaved: (value) =>
                        amount = value.trim().replaceAll(' ', ''),
                  ),
                  if (!appState.hideBalances) ...[availableBalance()],
                  SizedBox(height: height / 50),
                  dropdown(
                    (value) {
                      setState(() {
                        preferredPaymentMethod = value.toString();
                      });
                    },
                    getPaymentMethods,
                    preferredPaymentMethod.isEmpty
                        ? null
                        : preferredPaymentMethod,
                    'Select payment method',
                    context,
                    null,
                  ),
                  SizedBox(height: height / 50),
                  if (preferredPaymentMethod.toLowerCase() == 'bank') ...[
                    Row(
                      children: [
                        Text(
                          "accountnumber".tr(),
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 70),
                    CustomTextFormField.textField(
                      "",
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      60.sp,
                      300.sp,
                      initialValue: accountNumber,
                      onChanged: (value) {
                        setState(() {
                          accountNumber = value;
                        });
                      },
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          accountNumber = value!;
                        });
                      },
                    ),
                    Row(
                      children: [
                        Text(
                          "JOHN DOE KENECHUKWU",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 50),
                  ] else ...[
                    Row(
                      children: [
                        Text(
                          "Wallet username/address",
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(height: height / 70),
                    DestinationTextInput(
                      controller: toController,
                      focusNode: focusNode,
                      isReadOnly: false,
                      validator: validateTo,
                      onSaved: (value) =>
                          to = value.toLowerCase().trim().replaceAll(' ', ''),
                    ),
                    SizedBox(height: height / 50),
                  ],
                  SizedBox(height: height / 20),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  String? validateTo(String? value) {
    // reciever cannot be empty
    if (value!.isEmpty) return "enterreceiverusername".tr();

    if (value.length < 3) return "invalidusernameorpublickey".tr();

    return null;
  }

  String? validateAmount(String? value) {
    if (value!.isEmpty) {
      setState(() {
        amountError = true;
      });
      "pleaseenteramounttosend".tr();
    }

    if (double.tryParse(value) == null) {
      setState(() {
        amountError = true;
      });
      return "pleaseentervalidamount".tr();
    }

    if (double.tryParse(value)! <= 0) {
      setState(() {
        amountError = true;
      });
      return "valuemustbegreaterthan".tr(args: ['0']);
    }

    if (double.tryParse(value)! > asset!.amount!) {
      setState(() {
        amountError = true;
      });
      return "youdonthavesufficientbalance".tr();
    }

    if ((getAssetCode(asset!.assetCode) == 'XBN') &&
        double.tryParse(value)! > (asset!.amount! - 7)) {
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
      // following credential
      Map map = {
        "destination": to,
        "memo": memo,
        "amount": amount.toString(),
        "assetCode": asset!.assetCode == 'XBN' ? '' : asset!.assetCode,
        "assetIssuer": asset!.assetIssuer,
      };
      String requestBody = jsonEncode(map);

      Map responseData = await makePostRequest(
        uri: wallet.isSharedWalletAndCanInitiate
            ? '/v1/shared-access/payment'
            : '/v1/users/payment',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: wallet.publicKey!,
      );

      hideLoader(context);

      if (responseData['statusCode'] == 202) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        postProcessData(messageShown, messageLength, responseData['data']);
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {
      hideLoader(context);
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
        () => {postProcessData(messageShown, messageLength, data)},
      );

      messageShown++;
      return;
    }

    appState.viewData = {
      'walletPublicKey': wallet.publicKey,
      'assetCode': asset!.assetCode,
      'assetIssuer': asset!.assetIssuer,
      'rel': 'dashboard',
      'transactionData': data,
    };
    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: ConfirmTransactionViewPageConfig,
    );
  }
}

class DestinationTextInput extends StatefulWidget {
  final TextEditingController controller;
  final bool isReadOnly;
  final String? Function(String? value) validator;
  final FocusNode focusNode;
  final void Function(String value) onSaved;
  const DestinationTextInput({
    super.key,
    required this.controller,
    required this.focusNode,
    required this.isReadOnly,
    required this.validator,
    required this.onSaved,
  });

  @override
  State<DestinationTextInput> createState() => _DestinationTextInputState();
}

class _DestinationTextInputState extends State<DestinationTextInput> {
  late ColorNotifier notifier;
  late DataProvider appState;
  OverlayEntry? _overlayEntry = null;
  final LayerLink _layerLink = LayerLink();
  bool showAutocomplete = false;

  @override
  void initState() {
    super.initState();
    widget.focusNode.addListener(() {
      if (!widget.focusNode.hasFocus && _overlayEntry != null) {
        _overlayEntry!.remove();
        _overlayEntry = null;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    appState = Provider.of<DataProvider>(context, listen: true);
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return CompositedTransformTarget(
      link: this._layerLink,
      child: GestureDetector(
        child: CustomTextFormField.textField(
          "",
          notifier.getbluecolor,
          null,
          notifier.getgrey,
          notifier.getprefixicon,
          notifier.getblck,
          notifier.getgrey,
          75.sp,
          300.sp,
          onChanged: (value) {
            if (value.toString().startsWith(
              '${appState.primaryWallet.alias}_',
            )) {
              if (_overlayEntry != null) {
                _overlayEntry!.remove();
              }
              _overlayEntry = null;
              _overlayEntry = _createOverlayEntry();
              Overlay.of(context).insert(_overlayEntry!);
            } else {
              if (_overlayEntry != null) {
                _overlayEntry!.remove();
                _overlayEntry = null;
              }
            }
          },
          controller: widget.controller,
          readOnly: widget.isReadOnly,
          validator: widget.validator,
          onSaved: widget.onSaved,
        ),
      ),
    );
  }

  OverlayEntry _createOverlayEntry() {
    RenderBox renderBox = context.findRenderObject() as RenderBox;
    var size = renderBox.size;
    var wallets = appState.userInfo!.allWallets
        .where((wallet) => wallet.alias!.startsWith(widget.controller.text))
        .toList();
    return OverlayEntry(
      builder: (context) => Positioned(
        width: size.width,
        child: CompositedTransformFollower(
          link: this._layerLink,
          showWhenUnlinked: false,
          offset: Offset(0.0, size.height + 5.0),
          child: Material(
            elevation: 4.0,
            color: notifier.getaddsubwalletgrey,
            child: Container(
              constraints: BoxConstraints(maxHeight: height / 4),
              child: ListView.separated(
                padding: EdgeInsets.zero,
                shrinkWrap: true,
                itemCount: wallets.length,
                itemBuilder: (context, index) {
                  return TextButton(
                    onPressed: () {
                      widget.controller.text = wallets[index].alias!;
                      _overlayEntry!.remove();
                      _overlayEntry = null;
                    },
                    child: Row(
                      children: [
                        Text(
                          wallets[index].alias!,
                          textAlign: TextAlign.start,
                          style: TextStyle(color: notifier.getbluewhitecolor),
                        ),
                      ],
                    ),
                  );
                },
                separatorBuilder: (context, index) {
                  return Divider();
                },
              ),
            ),
          ),
        ),
      ),
    );
  }
}
