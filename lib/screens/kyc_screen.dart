// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:trovo_app/widgets/popups.dart';

/*
  This screen is for the Api service provider of dojah that enables the KYC
  verificatio of users
*/
class KYCScreen extends StatefulWidget {
  const KYCScreen({Key? key}) : super(key: key);

  @override
  State<KYCScreen> createState() => _KYCScreenState();
}

class _KYCScreenState extends State<KYCScreen> {
  bool granted = false;
  bool isCorporate = false;
  String userID = '';
  DateTime? timestamp = null;
  late DataProvider appState;
  late Future<dynamic> kycConfigFuture;
  late int activeLevel;
  String demoText =
      '''This is a demo process. Your data will not be stored or retained by Dojah and will only be used
        for the purpose of demonstrating this process flow.''';

  @override
  void initState() {
    checkPermission();
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    userID = appState.userInfo!.username!;
    isCorporate = appState.userInfo!.isCorporate;
    activeLevel = appState.userInfo!.kycVerified ?? 0;
    kycConfigFuture = fetchKycConfigs();
  }

  Future<void> fetchKycConfigs() async {
    try {
      var res = await StoreData().storeGetData(
        'lastKyc${activeLevel + 1}Submitted',
      );
      if (res != null) {
        timestamp = DateTime.tryParse(res);
      }

      var uri = '/v1/users/kyc/doja/configs';
      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
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

  // check that the user allows for permission to use their camera
  Future checkPermission() async {
    var status = await Permission.camera.request();
    granted = true;
    if (!status.isGranted) {
      granted = false;
    }

    var locationStatus = await Permission.locationWhenInUse.request();
    if (!locationStatus.isGranted) {
      granted = false;
    }

    var microphoneStatus = await Permission.microphone.request();
    if (!microphoneStatus.isGranted) {
      granted = false;
    }

    setState(() {});
  }

  // the various parameters that are submitted and are fetched from .env variables
  final appID = '67e69361a7d4138770eac9e2';
  final publicKey = 'test_pk_MNldwKATpyKxLjJoEfjkgH8hK';

  InAppWebViewController? _webViewController;

  InAppWebViewSettings settings = InAppWebViewSettings(
    clearCache: true,
    useShouldOverrideUrlLoading: true,
    mediaPlaybackRequiresUserGesture: false,
    useHybridComposition: true,
    allowsInlineMediaPlayback: true,
  );

  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    var height = MediaQuery.of(context).size.height;
    return FutureBuilder<dynamic>(
      future: kycConfigFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return Scaffold(
            body: SafeArea(
              child: SizedBox(
                height: height,
                child: Center(
                  child: CircularProgressIndicator(
                    backgroundColor: notifier.getbluecolor,
                    valueColor: new AlwaysStoppedAnimation<Color>(
                      notifier.getgreencolor,
                    ),
                    strokeWidth: 3.0,
                  ),
                ),
              ),
            ),
          );
        } else if (snapshot.connectionState == ConnectionState.done) {
          if (snapshot.hasError) {
            return Scaffold(
              body: SafeArea(
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: SizedBox(
                    height: height,
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
                              kycConfigFuture = fetchKycConfigs();
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
                  ),
                ),
              ),
            );
          } else if (snapshot.hasData) {
            var data = snapshot.data;
            final kycProgress = data['kycProgress'];
            final widgets = data['widgets'];
            String? widgetId;
            bool isAlreadySubmitted = false;

            for (int level = 1; level <= 4; level++) {
              final levelKey = 'kycLevel${level}Completed';
              activeLevel = level;

              if (kycProgress[levelKey] == 0) {
                // if (kycProgress['kycLevel${level}Submitted'] == 1 ||
                //     (timestamp != null &&
                //         timestamp!
                //             .add(Duration(hours: 1))
                //             .isAfter(DateTime.now()))) {
                //   isAlreadySubmitted = true;
                //   break;
                // }

                if (kycProgress['kycLevel${level}Submitted'] == 1) {
                  isAlreadySubmitted = true;
                  break;
                }

                for (final widget in widgets) {
                  if (widget['level'] == level) {
                    if (!isCorporate || widget['corporate'] == 1) {
                      widgetId = widget['id'];
                      break;
                    }
                  }
                }
                if (widgetId != null) break; // stop once widgetId is found
              }
            }

            if (isAlreadySubmitted) {
              return Scaffold(
                body: SafeArea(
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: SizedBox(
                      height: height,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            'Your level ${activeLevel} KYC is already submitted and is currently processing. Kindly wait for a resolution in the next few minutes.',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 16,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                          ElevatedButton(
                            onPressed: () {
                              Navigator.of(context).pop();
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
                              "close".tr(),
                              style: TextStyle(fontFamily: fontsemibold),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              );
            }

            return Scaffold(
              body: SafeArea(
                child: InAppWebView(
                  // initialOptions: options,
                  initialUrlRequest: URLRequest(
                    url: WebUri("https://widget.dojah.io"),
                  ),
                  initialData: InAppWebViewInitialData(
                    data:
                        """
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta
      name="viewport"
      content="width=device-width, user-scalable=no, initial-scale=1, maximum-scale=1, minimum-scale=1, shrink-to-fit=1"
    />
    <title>Dojah Widget Example</title>
    <style>
      @font-face {
      font-family: Athletics;
      src: url("fonts/Athletics.ttf");
      }
      * {
      font-family: Athletics, sans-serif
      }
      body {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
      }
      .paragraphText {
      font-size: 12px;
      color: #677189;
      padding: 24px 36px;
      margin-top: 0;
      margin-bottom: 0;
      background-color: #3977de05;
      text-align: center;
      line-height: 146.16%;
      border-radius: 8px;
      width: 100%;
      max-width: 268px;
      }
      button {
      margin-top: 32px;
      color: white;
      background-color: #3977de;
      padding: 14px 16px;
      border: none;
      outline: none;
      max-width: 384px;
      width: 100%;
      margin-left: auto;
      margin-right: auto;
      border-radius: 4px;
      cursor: pointer;
      }
      @media (max-width: 768px) {
      .paragraphText {
        padding: 16px;
      }
      }
    </style>
  </head>
  <body>
    <div
      style="
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 100vh";
      "
    >
      <!-- logo goes here -->
      <img src="favicon.ico" alt="" />
      <!-- <p>Logo</p> -->
      <p
        style="
        font-size: 16px;
        color: #1b2a4e;
        font-weight: bold;
        margin-top: 48px;
        margin-bottom: 21px;
        text-align: center;
        "
      >
        Verify your identity with Dojah
      </p>
      <div>
        <p class="paragraphText">
        $demoText
        </p>
      </div>
    
      <button id="custom-btn-connect">Start Level ${activeLevel} KYC</button>
      <p
        style="
        margin-top: 16px;
        font-size: 12px;
        line-height: 17.54px;
        text-align: center;
        "
      >
        By clicking the button above, you agree to <a href="https://www.dojah.io/policy" target="_blank" rel="noopener" style="text-decoration: underline; color: #1b2a4e"> Dojah's Privacy Policy </a>
      </p>
  </div>

  <script src="https://widget.dojah.io/widget.js"></script>
  <script>
    const customOptions = {
      app_id: '$appID',
      p_key: '$publicKey',
      type: 'custom',
      metadata: {
      user_id: '${userID}',
      },
      config: {
      widget_id: "$widgetId"
      },
      onSuccess: function (response) {
      window.flutter_inappwebview.callHandler('onSuccessCallback', response)
      },
      onError: function (err) {
      console.log('Callback Error', err)
      },
      onClose: function (response) {
      window.flutter_inappwebview.callHandler('onCancelCallback', response)
      }
    };
    const connectCustom = new Connect(customOptions);
    const connectBtnCustom = document.querySelector('#custom-btn-connect');
    connectBtnCustom.addEventListener('click', function () {
      connectCustom.setup();
      connectCustom.open();
    });
  </script>
</body>

</html>
                """,
                    historyUrl: WebUri("https://widget.dojah.io"),
                    mimeType: "text/html",
                    baseUrl: WebUri("https://widget.dojah.io"),
                  ),
                  onWebViewCreated: (controller) {
                    _webViewController = controller;

                    _webViewController?.addJavaScriptHandler(
                      handlerName: 'onSuccessCallback',
                      callback: (response) {
                        success(response);
                      },
                    );

                    _webViewController?.addJavaScriptHandler(
                      handlerName: 'onCancelCallback',
                      callback: (response) {
                        cancel(response);
                      },
                    );
                  },
                  onGeolocationPermissionsShowPrompt:
                      (controller, origin) async {
                        return GeolocationPermissionShowPromptResponse(
                          origin: origin,
                          allow: true, // Grant permission
                          retain: true, // Retain permission for future requests
                        );
                      },
                  onPermissionRequest: (controller, request) async {
                    return PermissionResponse(
                      resources: request.resources,
                      action: PermissionResponseAction.GRANT,
                    );
                  },
                ),
              ),
            );
          }
        }
        return Scaffold(
          body: SafeArea(
            child: Text(
              '',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.bold,
                fontFamily: fontsemibold,
              ),
            ),
          ),
        );
      },
    );
  }

  // a function that closes the webview after a successful capture
  void success(arg) {
    String timestamp = DateTime.now().toIso8601String();
    StoreData().storeInsertData('lastKyc${activeLevel}Submitted', timestamp);
    Navigator.of(context).pop();
    showSuccessAlert(
      context,
      text:
          'Your KYC procedure has been recorded successfully. Please wait for a few minutes for your information to be confirmed. Please refresh at intervals by pulling down on the home screen.',
      onTap: () {
        // Navigator.of(context).pop();
      },
    );
  }

  void cancel(arg) {
    Navigator.of(context).pop();
  }
}
