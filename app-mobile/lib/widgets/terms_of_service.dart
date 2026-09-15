import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import '../custom_bloc_observer/constants.dart';
import '../custom_bloc_observer/fonts.dart';
import '../custom_bloc_observer/notifire_clor.dart';
import '../storage/state.dart';
import '../utils/medeiaqury/medeiaqury.dart';

class TermsOfService extends StatefulWidget {
  final void Function(bool?)? onChanged;
  final bool value;
  final bool showError;

  TermsOfService({
    Key? key,
    required this.value,
    required this.showError,
    this.onChanged,
  }) : super(key: key);

  @override
  State<TermsOfService> createState() => _TermsOfServiceState();
}

class _TermsOfServiceState extends State<TermsOfService> {
  late ColorNotifier notifier;
  late DataProvider appState;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: false);
    return termsOfService();
  }

  Widget termsOfService() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Transform.scale(
              scale: 1.sp,
              child: Checkbox(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.all(
                    Radius.circular(5.sp),
                  ),
                ),
                activeColor: notifier.isDark
                    ? notifier.getbluecolor50
                    : notifier.getbluecolor90,
                side: BorderSide(
                  color: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor90,
                ),
                value: widget.value,
                onChanged: widget.onChanged,
              ),
            ),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(
                      "iagreetothe".tr(),
                      style: TextStyle(
                          fontSize: height / 55,
                          color: notifier.getblck,
                          fontFamily: fontbody),
                    ),
                    GestureDetector(
                      onTap: () {
                        appState.goToWebView(termsOfServiceUrl);
                      },
                      child: Text(
                        ' ' + "termsofservices".tr(),
                        style: TextStyle(
                          fontFamily: fontbody,
                          fontSize: height / 55,
                          color: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
                        ),
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Text(
                      "and".tr(),
                      style: TextStyle(
                          fontFamily: fontbody,
                          fontSize: height / 55,
                          color: notifier.getblck),
                    ),
                    SizedBox(
                      width: 5,
                    ),
                    GestureDetector(
                      onTap: () {
                        appState.goToWebView(privacyPolicyUrl);
                      },
                      child: Text(
                        "privacypolicy".tr(),
                        style: TextStyle(
                          fontFamily: fontbody,
                          fontSize: height / 55,
                          color: notifier.isDark
                              ? notifier.getbluecolor50
                              : notifier.getbluecolor90,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            )
          ],
        ),
        // You need to accept terms
        if (widget.showError) ...[
          Padding(
            padding: EdgeInsets.fromLTRB(10.0, 0, 0, 0),
            child: Text(
              "termsofserviceerror".tr(),
              style: TextStyle(
                  color: Colors.red, fontSize: 12, fontWeight: FontWeight.w400),
            ),
          ),
        ],
      ],
    );
  }
}
