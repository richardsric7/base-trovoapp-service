import 'dart:async';
import 'dart:io';
import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:provider/provider.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';
import 'package:scan/scan.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';

class QrScanner extends StatefulWidget {
  const QrScanner({Key? key}) : super(key: key);

  @override
  _QrScannerState createState() => _QrScannerState();
}

bool expanded = false;

class _QrScannerState extends State<QrScanner> {
  final GlobalKey qrKey = GlobalKey(debugLabel: 'QR');
  QRViewController? controller;
  DataProvider? appState;
  bool? flashLightStatus;

  Future<String?>? scanResult;
  // In order to get hot reload to work we need to pause the camera if the platform
  // is android, or resume the camera if the platform is iOS.
  @override
  Future<void> reassemble() async {
    super.reassemble();
    if (Platform.isAndroid) {
      await controller?.pauseCamera();
    }
    controller?.resumeCamera();
  }

  @override
  Widget build(BuildContext context) {
    setState(() {
      appState = Provider.of<DataProvider>(context, listen: false);
    });
    return FutureBuilder(
      future: scanResult,
      builder: ((context, AsyncSnapshot<String?> snapshot) {
        return _buildScanner(snapshot);
      }),
    );
  }

  Widget _buildScanner(AsyncSnapshot<String?> snapshot) {
    return Scaffold(
      body: Stack(
        fit: StackFit.passthrough,
        children: [
          QRView(
            cameraFacing: CameraFacing.back,
            key: qrKey,
            onQRViewCreated: _onQRViewCreated,
            overlay: QrScannerOverlayShape(
              overlayColor: const Color.fromARGB(69, 0, 0, 0),
              borderRadius: 10,
              borderWidth: 5,
              borderColor: Colors.white,
              cutOutSize: MediaQuery.of(context).size.width - 45,
              cutOutBottomOffset: 60,
            ),
          ),
          AppBar(
            iconTheme: IconThemeData(
              color: Colors.grey.shade800, //change your color here
            ),
            backgroundColor: Colors.transparent.withOpacity(0.0),
          ),
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                width: MediaQuery.of(context).size.width - 20,
                height: MediaQuery.of(context).size.width - 20,
                child: Container(
                  margin: const EdgeInsets.all(15.0),
                  padding: const EdgeInsets.all(3.0),
                  child: Align(
                    alignment: Alignment.bottomCenter,
                    child: IconButton(
                      onPressed: _toggleFlashlight,
                      icon: Icon(
                        _getIcon(),
                        color: Colors.white60,
                      ),
                    ),
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 15.0),
                child: Center(
                  child: Text(
                    LanguageEn.placeqrcode,
                    style: TextStyle(
                        color: Colors.white60, fontWeight: FontWeight.w500),
                    textAlign: TextAlign.center,
                  ),
                ),
              ),
              SizedBox(
                width: 200,
                height: 74,
                child: Container(
                  margin: const EdgeInsets.all(15.0),
                  padding: const EdgeInsets.all(3.0),
                ),
              ),
            ],
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed:
            snapshot.connectionState == ConnectionState.none ? decode : null,
        child: const Icon(
          Icons.image_outlined,
          size: 35.0,
        ),
        backgroundColor: Colors.transparent.withOpacity(0.0).withOpacity(0.0),
        elevation: 0,
      ),
    );
  }

  void _onQRViewCreated(QRViewController controller) {
    setState(() {
      this.controller = controller;
    });
    controller.scannedDataStream.listen((scanData) {
      _handleScanResult(scanData.code);
    });
    controller.resumeCamera();
  }

  void _toggleFlashlight() async {
    await controller?.toggleFlash();
    flashLightStatus = await controller?.getFlashStatus();
    // run setState to update state after all the awaitables finish running
    setState(() {});
  }

  IconData _getIcon() {
    if (flashLightStatus == null) {
      return Icons.flashlight_on_rounded;
    } else {
      return flashLightStatus!
          ? Icons.flashlight_off_rounded
          : Icons.flashlight_on_rounded;
    }
  }

  /// decode from local file
  Future<void> decode() async {
    var image = await ImagePicker().pickImage(source: ImageSource.gallery);
    if (image != null) {
      Scan.parse(image.path).then((scanResult) {
        _handleScanResult(scanResult);
      });
    }
  }

  void _handleScanResult(String? scanResult) async {
    controller!.pauseCamera();
    print('scan result:');
    print(scanResult);
    if (scanResult != null) {
      runDynamicLinks(Uri.parse(scanResult));
    }

    Tooltip(
        message: 'Does that look like a QR Code file to you?',
        showDuration: Duration(seconds: 10));
    controller!.resumeCamera();
  }

  void runDynamicLinks(uri) async {
    try {
      showLoader(context);
      final PendingDynamicLinkData? data =
          await FirebaseDynamicLinks.instance.getDynamicLink(uri);

      if (data != null) {
        final Uri deepLink = data.link;
        print('deeplink... $deepLink');

        // print('The deepLink data on success is $deepLink');
        print(deepLink.queryParameters);

        if (deepLink.queryParameters['action'] == 'payment') {
          if (deepLink.queryParameters['assetCode'] != '' &&
              deepLink.queryParameters['assetCode'] != null) {
            var deeplinkInfo = {
              "assetCode": deepLink.queryParameters['assetCode'],
              "assetIssuer": deepLink.queryParameters['assetIssuer'],
              "source": "qr2",
              "receiver": deepLink.queryParameters['paymentDestination'],
              "amount":
                  deepLink.queryParameters['amount'], // amount we want to send
              "memo": deepLink.queryParameters['memo'],
              'action': 'payment'
            };
            print('this is deeplinkInfo: $deeplinkInfo');
            var assetInfo = null;

            var assetBalances = appState!.assetBalances;
            var claimedAssets =
                assetBalances[appState!.activeWallet!.publicKey]['claimed'];

            var deeplinkAssetCode = deeplinkInfo['assetCode'] == 'XBN'
                ? ''
                : deeplinkInfo['assetCode'];

            for (var asset in claimedAssets) {
              print('this is asset: $asset');
              if (asset['assetCode'] == deeplinkAssetCode &&
                  asset['assetIssuer'] == deeplinkInfo['assetIssuer']) {
                assetInfo = {
                  'assetCode': asset['assetCode'],
                  'assetIssuer': asset['assetIssuer'],
                  'amount': asset['amount'], // balance amount in the wallet
                  'qrCode': asset['qrCode'],
                  'imageUrl': asset['imageUrl'],
                };

                // exit the loop immediately we get what we are looking for
                break;
              }
            }

            print('this is assetInfo: $assetInfo');

            appState!.viewData = {
              SendAssetViewPageConfig.key: {
                'assetCode': assetInfo['assetCode'],
                'assetIssuer': assetInfo['assetIssuer'],
                'amount': assetInfo['amount'],
                'imageUrl': assetInfo['imageUrl'],
                'deepLinkInfo': deeplinkInfo,
              },
              // to avoid unexpected behaviour in the assetdetails page
              // add the AssetDetailsViewPageConfig view data.
              // The issue occurs when user goes through assetDetailsPage => sendAsset => scanQr
              // apparently the previous page has to be rebuilt when you navigate using
              // appState?.currentAction = PageAction(state: PageState.replace, page: SendAssetViewPageConfig);
              // with PageState.replace.
              AssetDetailsViewPageConfig.key: {
                'assetCode': assetInfo['assetCode'],
                'assetIssuer': assetInfo['assetIssuer'],
                'amount': assetInfo['amount'],
                'qrCode': assetInfo['qrCode'],
                'imageUrl': assetInfo['imageUrl'],
              }
            };

            hideLoader(context);

            appState?.currentAction = PageAction(
                state: PageState.replace, page: SendAssetViewPageConfig);
          }
        } else {
          hideLoader(context);
          popup(
            context,
            title: 'Error!',
            message:
                'The QR code is not meant for ${deepLink.queryParameters['assetCode']} payment. Please scan the correct QR code!',
          );
        }
      } else {
        popup(context,
            title: 'Error!',
            message: 'The QR code is not meant for payment with Trovo Wallet');
        hideLoader(context);
      }
    } catch (e) {
      hideLoader(context);
      print('there was an error $e');
    }
  }

  @override
  void dispose() {
    controller?.dispose();
    super.dispose();
  }
}
