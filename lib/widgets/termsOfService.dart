import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:provider/provider.dart';
import '../Custom_BlocObserver/constants.dart';
import '../Custom_BlocObserver/fonts.dart';
import '../Custom_BlocObserver/notifire_clor.dart';
import '../screens/page_view/web_view.dart';
import '../utils/enstring.dart';
import '../utils/medeiaqury/medeiaqury.dart';

class TermsOfService extends StatefulWidget {
  final void Function(bool?)? onChanged;
  bool value;
  bool showError;

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

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: false);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
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
                activeColor: notifier.getbluecolor,
                side: BorderSide(color: notifier.getbluecolor),
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
                      LanguageEn.iagreetothe,
                      style: TextStyle(
                          fontSize: height / 55,
                          color: notifier.getblck,
                          fontFamily: fontbody),
                    ),
                    GestureDetector(
                      onTap: () =>
                          Get.to(() => TrovoWebView(url: termsOfServiceUrl)),
                      child: Text(
                        ' ' + LanguageEn.termsofservices,
                        style: TextStyle(
                            fontFamily: fontbody,
                            fontSize: height / 55,
                            color: notifier.getbluecolor),
                      ),
                    ),
                  ],
                ),
                Row(
                  children: [
                    Text(
                      LanguageEn.and,
                      style: TextStyle(
                          fontFamily: fontbody,
                          fontSize: height / 55,
                          color: notifier.getblck),
                    ),
                    SizedBox(
                      width: 5,
                    ),
                    GestureDetector(
                      onTap: () =>
                          Get.to(() => TrovoWebView(url: privacyPolicyUrl)),
                      child: Text(
                        LanguageEn.privacypolicy,
                        style: TextStyle(
                            fontFamily: fontbody,
                            fontSize: height / 55,
                            color: notifier.getbluecolor),
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
              LanguageEn.termsofserviceerror,
              style: TextStyle(
                  color: Colors.red, fontSize: 12, fontWeight: FontWeight.w400),
            ),
          ),
        ],
      ],
    );
  }
}
