import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';

class FirebaseDynamicLinkInitializer {
  void initializeDeeplinking() async {
    try {
      // final dynamicLinkParams = DynamicLinkParameters(
      DynamicLinkParameters(
        link: Uri.parse("https://trovowallet.page.link"),
        uriPrefix: "https://trovowallet.page.link",
        androidParameters: const AndroidParameters(
          packageName: "com.trovo.wallet",
          minimumVersion: 1,
        ),
        iosParameters: const IOSParameters(
          bundleId: "com.trovo.wallet",
          appStoreId: "6443621693",
          minimumVersion: "0.0.1",
        ),
        // googleAnalyticsParameters: const GoogleAnalyticsParameters(
        //   source: "twitter",
        //   medium: "social",
        //   campaign: "example-promo",
        // ),
        // socialMetaTagParameters: SocialMetaTagParameters(
        //   title: "Example of a Dynamic Link",
        //   imageUrl: Uri.parse("https://example.com/image.png"),
        // ),
      );
      // final dynamicLink =
      //     await FirebaseDynamicLinks.instance.buildShortLink(dynamicLinkParams);
    } catch (e) {
      print('===============> no network or something...$e');
    }
  }

  Future<PendingDynamicLinkData?> getInitialLink() async {
    // Get any initial links
    return await FirebaseDynamicLinks.instance.getInitialLink();
  }
}
