// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/storage/state.dart';
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
  String userID = '';

  @override
  void initState() {
    checkPermission();
    super.initState();
    var appState = Provider.of<DataProvider>(context, listen: false);
    userID = appState.userInfo!.username!;
  }

  // check that the user allows for permission to use their camera
  Future checkPermission() async {
    var status = await Permission.camera.request();
    granted = true;
    print('status =============> $status');
    if (!status.isGranted) {
      granted = false;
    }

    var locationStatus = await Permission.locationWhenInUse.request();
    print('location status =============> $locationStatus');
    if (!locationStatus.isGranted) {
      granted = false;
    }

    var microphoneStatus = await Permission.microphone.request();
    print('microphone status =============> $microphoneStatus');
    if (!microphoneStatus.isGranted) {
      granted = false;
    }

    setState(() {});
  }

  // the various parameters that are submitted and are fetched from .env variables
  final appID = '67e69361a7d4138770eac9e2';
  final publicKey = 'test_pk_MNldwKATpyKxLjJoEfjkgH8hK';
  final widgetID = '67e6968cc5f45aec8ae35e0a';

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
    return Scaffold(
      body: SafeArea(
        child: InAppWebView(
          // initialOptions: options,
          initialUrlRequest: URLRequest(url: WebUri("https://widget.dojah.io")),
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
      Verify your identity with this Demo
    </p>
    <div>
      <p class="paragraphText">
      Your data will not be stored or retained by Dojah and will only be used
      for the purpose of demonstrating this process flow.
      </p>
    </div>
  
    <button id="custom-btn-connect">Custom Widget</button>
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
      widget_id: "$widgetID"
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
          onGeolocationPermissionsShowPrompt: (controller, origin) async {
            return GeolocationPermissionShowPromptResponse(
              origin: origin,
              allow: true, // Grant permission
              retain: true, // Retain permission for future requests
            );
          },
          // onPermissionRequest: (
          //   controller,
          //   PermissionResponse(
          //       resources: resources, action: PermissionResponseAction.GRANT)
          // ),
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

  // a function that closes the webview after a successful capture
  void success(arg) {
    showSuccessAlert(
      context,
      onTap:
          'Your KYC procedure has been recorded successfully, you will be notified with further information',
    );
    Navigator.of(context).pop();
  }

  void cancel(arg) {
    Navigator.of(context).pop();
  }
}
