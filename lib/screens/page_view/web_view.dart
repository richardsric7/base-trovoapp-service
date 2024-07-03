import 'dart:io';
import 'package:flutter/foundation.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/storage/state.dart';
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
    getdarkmodepreviousstate();
    // Enable virtual display.
    if (Platform.isAndroid) WebView.platform = AndroidWebView();
  }

  // Reference to webview controller
  WebViewController? _controller;

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
              WebView(
                javascriptMode: JavascriptMode.unrestricted,
                gestureRecognizers: gestureRecognizers,
                initialUrl: appState.initialUrl,
                onPageStarted: (value) => {
                  print('loading... $value'),
                  setState(() {
                    isLoading = true;
                  })
                },
                onWebViewCreated: (WebViewController webViewController) {
                  // Get reference to WebView controller to access it globally
                  _controller = webViewController;
                },
                onPageFinished: (value) {
                  print('finished loading.$value');
                  setState(() {
                    print('setting state...');
                    isLoading = false;
                    hideLoader(context);
                  });

                  // In the final result page we check the url to make sure  it is the last page.if (url.contains('/finalresponse.html')) {
                  _controller?.runJavascript('''
                        const pdfjs = require('pdfs-dist');
                        pdfjs.getPdfInfo('(link unavailable)', (info) => {
                          console.log(info);
                        });

                      ''');
                },
                onWebResourceError: (error) => {print('error $error')},
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
