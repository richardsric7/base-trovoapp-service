import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';

class FirebaseDynamicLinkInitializer {
  void initializeDeeplinking(String walletMode) async {
    try {
      walletMode == "Testnet"
          ? DynamicLinkParameters(
              link: Uri.parse("https://trovowalletdev.page.link"),
              uriPrefix: "https://trovowalletdev.page.link",
              androidParameters: const AndroidParameters(
                packageName: "com.trovo.wallet",
                minimumVersion: 1,
              ),
              iosParameters: const IOSParameters(
                bundleId: "com.trovo.wallet",
                appStoreId: "6443621693",
                minimumVersion: "0.0.1",
              ),
            )
          :
          // final dynamicLinkParams = DynamicLinkParameters(
          DynamicLinkParameters(
              link: Uri.parse("https://trovoapp.page.link"),
              uriPrefix: "https://trovoapp.page.link",
              androidParameters: const AndroidParameters(
                packageName: "com.trovo.wallet",
                minimumVersion: 1,
              ),
              iosParameters: const IOSParameters(
                bundleId: "com.trovo.wallet",
                appStoreId: "6443621693",
                minimumVersion: "0.0.1",
              ),
            );
    } catch (e) {
      print('===============> no network or something...$e');
    }
  }

  Future<PendingDynamicLinkData?> getInitialLink() async {
    // Get any initial links
    return await FirebaseDynamicLinks.instance.getInitialLink();
  }
}
