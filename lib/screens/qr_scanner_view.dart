import 'dart:async';
import 'package:easy_localization/easy_localization.dart';
import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:mobile_scanner/mobile_scanner.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

class QrScanner extends StatefulWidget {
  const QrScanner({Key? key}) : super(key: key);

  @override
  _QrScannerState createState() => _QrScannerState();
}

bool expanded = false;

class _QrScannerState extends State<QrScanner> {
  final GlobalKey qrKey = GlobalKey(debugLabel: 'QR');
  // QRViewController? controller;
  DataProvider? appState;
  bool flashLightStatus = false;

  final MobileScannerController controller = MobileScannerController(
    detectionSpeed: DetectionSpeed.normal,
    facing: CameraFacing.back,
  );

  String? scanResult;

  @override
  Future<void> reassemble() async {
    super.reassemble();
    print('resuming camera...');
    controller.start();
    print('camera resumed...');
  }

  @override
  void initState() {
    controller.start();
    print('scanner controller started!');
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    setState(() {
      appState = Provider.of<DataProvider>(context, listen: false);
    });
    return _buildScanner();
  }

  void _handleBarcode(BarcodeCapture scanResult) {
    if (mounted) {
      _handleScanResult(scanResult.barcodes.first.rawValue);
    }
  }

  Widget _buildScanner() {
    return Scaffold(
      body: Stack(
        fit: StackFit.passthrough,
        children: [
          MobileScanner(onDetect: _handleBarcode, controller: controller),
          AppBar(
            iconTheme: IconThemeData(
              color: Colors.grey.shade800, //change your color here
            ),
            backgroundColor: Colors.transparent.withValues(alpha: 0),
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
                      icon: Icon(_getIcon(), color: Colors.white60),
                    ),
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 15.0),
                child: Center(
                  child: Text(
                    "placeqrcode".tr(),
                    style: TextStyle(
                      color: Colors.white60,
                      fontWeight: FontWeight.w500,
                    ),
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
        onPressed: decode,
        child: const Icon(
          Icons.image_outlined,
          size: 35.0,
          color: Colors.white60,
        ),
        backgroundColor: Colors.transparent.withValues(alpha: 0),
        elevation: 0,
      ),
    );
  }

  void _toggleFlashlight() async {
    controller.start();
    await controller.toggleTorch();
    flashLightStatus = !flashLightStatus;
    setState(() {});
  }

  IconData _getIcon() => flashLightStatus
      ? Icons.flashlight_off_rounded
      : Icons.flashlight_on_rounded;

  /// decode from local file
  Future<void> decode() async {
    var image = await ImagePicker().pickImage(source: ImageSource.gallery);
    if (image != null) {
      controller.analyzeImage(image.path).then((result) {
        _handleScanResult(result!.barcodes.first.rawValue);
      });
    }
  }

  void _handleScanResult(String? rawValue) async {
    print('scan result: $scanResult, raw value: $rawValue');
    if (scanResult != null) return;

    if (rawValue != null) {
      scanResult = rawValue;
      runDynamicLinks(Uri.parse(scanResult ?? ''));
    } else {
      Tooltip(
        message: 'Does that look like a QR Code file to you?',
        showDuration: Duration(seconds: 10),
      );
    }
    setState(() {});
  }

  void runDynamicLinks(uri) async {
    try {
      showLoader(context);
      final PendingDynamicLinkData? data = await FirebaseDynamicLinks.instance
          .getDynamicLink(uri);

      if (data != null) {
        final Uri deepLink = data.link;
        print(deepLink.queryParameters);
        appState!.processDeepLink(
          context,
          deepLink,
          rel: 'qrScanner',
          onCancel: () => setState(() {
            scanResult = null;
          }),
        );
        hideLoader(context);
      } else {
        popup(
          context,
          title: 'Error!',
          message:
              'Something went wrong. We could not process the following link [$uri]; Could be caused by bad network or a bad qrcode image.',
          onClose: () {
            setState(() {
              scanResult = null;
            });
          },
        );
        hideLoader(context);
      }
    } catch (e) {
      hideLoader(context);
      print('there was an error $e');
      setState(() {
        scanResult = null;
      });
    }
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }
}
