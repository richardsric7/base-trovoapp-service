import 'package:flutter/foundation.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:webview_flutter/webview_flutter.dart';
import '../../custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';

class TrovoWebView extends StatefulWidget {
  TrovoWebView({Key? key}) : super(key: key);

  @override
  TrovoWebViewState createState() => TrovoWebViewState();
}

class TrovoWebViewState extends State<TrovoWebView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  bool isLoading = false;
  late WebViewController controller;
  final Set<Factory<OneSequenceGestureRecognizer>> gestureRecognizers = {
    Factory(() => EagerGestureRecognizer())
  };
  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    getdarkmodepreviousstate();

    controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onPageStarted: (String url) {
            isLoading = true;
          },
          onPageFinished: (String url) {
            isLoading = false;
            hideLoader(context);
          },
          onHttpError: (HttpResponseError error) {},
          onWebResourceError: (WebResourceError error) {
            print('error $error');
          },
        ),
      )
      ..loadRequest(Uri.parse(appState.initialUrl));
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
          resizeToAvoidBottomInset: false,
          backgroundColor: notifier.getwihitecolor,
          appBar: CustomAppBarWithoutBanner(
              context, notifier.getwihitecolor, "", notifier.getblck,
              height: height / 15),
          body: Stack(
            children: [
              SingleChildScrollView(
                child: SizedBox(
                  height: height - 90,
                  child: WebViewWidget(
                    gestureRecognizers: gestureRecognizers,
                    controller: controller,
                  ),
                ),
              ),
              if (isLoading) ...[
                Container(
                  child: showLoader(context),
                )
              ],
            ],
          )),
    );
  }
}
