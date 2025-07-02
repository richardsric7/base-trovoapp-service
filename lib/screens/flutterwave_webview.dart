import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:uuid/uuid.dart';
import 'package:webview_flutter/webview_flutter.dart';
import 'dart:io';

import 'package:flutter/foundation.dart';

/// the is paystack webview
class FlutterwaveWebView extends StatefulWidget {
  const FlutterwaveWebView({Key? key}) : super(key: key);

  @override
  State<FlutterwaveWebView> createState() => _FlutterwaveWebViewState();
}

class _FlutterwaveWebViewState extends State<FlutterwaveWebView> {
  late DataProvider appState;
  var uuid = Uuid();
  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    final WebViewController controller = WebViewController();

    controller
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onProgress: (int progress) {
            debugPrint('WebView is loading (progress : $progress%)');
          },
          onPageStarted: (String url) {
            debugPrint('Page started loading: $url');
          },
          onPageFinished: (String url) {
            debugPrint('Page finished loading: $url');
          },
          onNavigationRequest: (NavigationRequest request) {
            if (request.url.startsWith('https://www.youtube.com/')) {
              debugPrint('blocking navigation to ${request.url}');
              return NavigationDecision.prevent;
            }
            debugPrint('allowing navigation to ${request.url}');
            return NavigationDecision.navigate;
          },
          onHttpError: (HttpResponseError error) {
            debugPrint('Error occurred on page: ${error.response?.statusCode}');
          },
          onUrlChange: (UrlChange change) {
            debugPrint('url change to ${change.url}');
          },
          onHttpAuthRequest: (HttpAuthRequest request) {
            // openDialog(request);
          },
        ),
      )
      ..addJavaScriptChannel(
        'Toaster',
        onMessageReceived: (JavaScriptMessage message) {
          ScaffoldMessenger.of(
            context,
          ).showSnackBar(SnackBar(content: Text(message.message)));
        },
      )
      // ..loadRequest(Uri.parse('https://flutter.dev'))
      ..runJavaScript("""
          var link = document.createElement('link');
          link.rel = 'stylesheet';
          link.type = 'text/css';
          link.href = 'https://checkout.paystack.com/assets/index-CT0IlyWf.css';
          document.head.appendChild(link);
        """);

    // setBackgroundColor is not currently supported on macOS.
    if (kIsWeb || !Platform.isMacOS) {
      controller.setBackgroundColor(const Color(0x80000000));
    }
  }

  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    var viewData = appState.viewData!;
    var userInfo = appState.userInfo!;
    String uniqueId = uuid.v4();
    String html =
        '''
          <!DOCTYPE html>
          <html>
          <meta name="viewport" content="width=device-width, initial-scale=1.0">
          <script src="https://checkout.flutterwave.com/v3.js"></script>
          <body onload="makePayment()">
              <script>
                  function makePayment() {
                      FlutterwaveCheckout({
                        public_key: 'FLWPUBK_TEST-45bd332ee4bdefdcacd6d2513944cd16-X',
                        tx_ref: '$uniqueId',
                        amount: ${viewData['activationAmount']},
                        currency: 'NGN',
                        payment_options: 'ussd, card, bank transfer',                        
                        meta: {                          
                          user_id: '${userInfo.username}',
                          transaction_type: 'ACTIVATION',                      
                          product: 'ACTIVATION',
                          destination_wallet_id: '${appState.primaryWallet.publicKey}',
                        },
                        customer: {
                          email: '${userInfo.email}',
                          phone_number: '${userInfo.mobile}',
                          name: '${userInfo.fullName}',                              
                        },
                        customizations: {
                          title: 'Account Activation',
                          description: 'Pay ${viewData['activationAmount']} to activate account.',
                          logo: 'https://trovotech.io/img/Transperant-Logo-1.png',                          
                        },
                        callback: function(response){
                          if(response.status === 'successful'){
                            window.location.href = 'https://trovo.app/success'; 
                          }
                        },
                        onclose: function (incomplete) {
                          window.location.href = 'https://trovo.app/cancel';                          
                        },
                      });
                    }
              </script>
          </body>
          </html>
    ''';
    return Scaffold(
      backgroundColor: notifier.getwihitecolor,
      appBar: AppBar(backgroundColor: notifier.getwihitecolor, elevation: 0),
      body: WebViewWidget(
        controller: WebViewController()
          ..setJavaScriptMode(JavaScriptMode.unrestricted)
          ..setBackgroundColor(notifier.getwihitecolor)
          ..setNavigationDelegate(
            NavigationDelegate(
              onProgress: (int progress) {
                // Update loading bar.
              },
              onPageStarted: (String url) {
                showLoader(context);
              },
              onPageFinished: (url) {
                hideLoader(context);
              },
              onWebResourceError: (WebResourceError error) {},
              onNavigationRequest: (NavigationRequest request) {
                if (request.url.contains('https://trovo.app/success')) {
                  Navigator.pop(context, false);
                  showSuccessAlert(
                    context,
                    text:
                        'Your activation payment of ${viewData['activationAmount']} was successful. Value will be transferred to your primary wallet as soon as the payment is confirmed.',
                    onTap: () {
                      appState.currentAction = PageAction(
                        state: PageState.replaceAll,
                        page: BottomHomePageConfig,
                      );
                    },
                    buttonText: 'Done',
                  );
                  return NavigationDecision.prevent;
                }

                if (request.url.contains('https://trovo.app/cancel')) {
                  Navigator.pop(context, false);
                  // popup(
                  //   context,
                  //   title: 'Cancelled',
                  //   message: 'Operation was cancelled',
                  //   buttonText: 'Done',
                  // );
                  return NavigationDecision.prevent;
                }

                return NavigationDecision.navigate;
              },
            ),
          )
          ..loadHtmlString(html),
      ),
    );
  }
}
