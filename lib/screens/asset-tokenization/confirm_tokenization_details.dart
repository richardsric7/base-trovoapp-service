import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ConfirmTokenizationDetails extends StatefulWidget {
  const ConfirmTokenizationDetails({Key? key}) : super(key: key);

  @override
  State<ConfirmTokenizationDetails> createState() =>
      _ConfirmTokenizationDetails();
}

class _ConfirmTokenizationDetails extends State<ConfirmTokenizationDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  late TokenizedAsset tokenizedAsset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = TokenizedAsset().deserializeJson(appState.viewData!);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    var isAlreadySubmitted = tokenizedAsset.tokenizationStatus != null;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
                context,
                notifier.getwihitecolor,
                isAlreadySubmitted
                    ? tokenizedAsset.assetName!
                    : "confirmyourinformation".tr(),
                notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 40,
              ),
              item("assetname".tr(), '${tokenizedAsset.assetName}'),
              SizedBox(
                height: height / 40,
              ),
              item("assetCode".tr(), '${tokenizedAsset.assetCode}'),
              SizedBox(
                height: height / 40,
              ),
              item("totaltokenstobeissued".tr(),
                  '${tokenizedAsset.numberOfTokenToBeIssued} ${tokenizedAsset.assetCode}'),
              SizedBox(
                height: height / 50,
              ),
              item("totaltokenstobesold".tr(),
                  '${tokenizedAsset.numberOfTokenToBeSold} ${tokenizedAsset.assetCode}'),
              SizedBox(
                height: height / 50,
              ),
              item("pricepertoken".tr(),
                  '${(double.parse(tokenizedAsset.assetCurrentValue.toString()) / tokenizedAsset.numberOfTokenToBeIssued!)} ${tokenizedAsset.assetQuoteCurrency}'),
              SizedBox(
                height: height / 50,
              ),
              item("totalamounttoberaised".tr(),
                  '${tokenizedAsset.numberOfTokenToBeIssued} ${tokenizedAsset.assetQuoteCurrency}'),
              SizedBox(
                height: height / 50,
              ),
              item("tokenizationfee".tr(),
                  "N1,000,000.00 + 1,500,000.00 ${tokenizedAsset.assetCode}"),
              SizedBox(
                height: height / 50,
              ),
              item("primaryofferingstartdate".tr(),
                  '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesStart!)}'),
              SizedBox(
                height: height / 50,
              ),
              item("primaryofferingenddate".tr(),
                  '${DateFormat('MMMM dd, yyyy').format(tokenizedAsset.salesEnd!)}'),
              SizedBox(
                height: height / 20,
              ),
              if (isAlreadySubmitted) ...[
                Button(
                  "viewpaymentdetails".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: TokenizationFeePaymentViewPageConfig,
                    );
                  },
                ),
                SizedBox(height: height / 70),
                ButtonOutlined(
                  'back'.tr(),
                  notifier.getwihitecolor,
                  notifier.getbluewhitecolor,
                  borderColor: notifier.getbluewhitecolor,
                  onTap: () {
                    Navigator.of(context).pop();
                  },
                ),
              ] else ...[
                Button(
                  "proceedtopay".tr(),
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: TokenizationFeePaymentViewPageConfig,
                    );
                  },
                ),
              ],
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

  Widget item(String key, String value) {
    return Column(
      children: [
        Text(
          key,
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: 15,
            fontWeight: FontWeight.w400,
            color: notifier.getbluewhitecolor,
            fontFamily: fontbody,
          ),
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(
                      vertical: 10.0, horizontal: 15),
                  child: Container(
                    width: width / 1.3,
                    child: Text(
                      value,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontWeight: FontWeight.w500,
                        color: notifier.getbluewhitecolor,
                        fontSize: 15.sp,
                        fontFamily: fontsemibold,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
