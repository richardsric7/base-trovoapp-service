// import 'package:firebase_dynamic_links/firebase_dynamic_links.dart';

// class DynamicLinkService {
//   Future handleDynamicLink() async {
//     // Get initial dynamic link if the app is started using the link
//     final PendingDynamicLinkData? data =
//         await FirebaseDynamicLinks.instance.getInitialLink();

//     _handleDeepLink(data);

//     // INTO FOREGROUND FROM DYNAMIC LINK LOGIC
//     FirebaseDynamicLinks.instance.onLink.listen(
//         (PendingDynamicLinkData dynamicLinkData) {
//       _handleDeepLink(dynamicLinkData);
//     }, onError: (e) async {
//       print('Dynamic Link Failed: ${e.message}');
//     });
//   }

//   void _handleDeepLink(PendingDynamicLinkData? data) {
//     final Uri deeplink = data!.link;
//     print('_handleDeepLink | deeplink: $deeplink');
//   }
// }
