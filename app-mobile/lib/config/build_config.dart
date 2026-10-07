// Values chosen when the app is built, with --dart-define, e.g.
//
//   flutter build appbundle \
//     --dart-define=TROVO_API_URL_MAINNET=https://api.trovo.example.com \
//     --dart-define=FLUTTERWAVE_PUBLIC_KEY_MAINNET=FLWPUBK-xxxx-X
//
// Every value has a default, so a plain `flutter run` keeps working. See
// CONFIGURATION.md.
class BuildConfig {
  /// app-backend's address used while the app is in Testnet mode.
  static const String apiUrlTestnet = String.fromEnvironment(
    'TROVO_API_URL_TESTNET',
    defaultValue: 'https://api.dev.trovo.app',
  );

  /// app-backend's address used while the app is in Mainnet mode.
  static const String apiUrlMainnet = String.fromEnvironment(
    'TROVO_API_URL_MAINNET',
    defaultValue: 'https://api.trovotechnologies.com',
  );

  /// Flutterwave's public test key, for card payments in Testnet mode.
  static const String flutterwavePublicKeyTestnet = String.fromEnvironment(
    'FLUTTERWAVE_PUBLIC_KEY_TESTNET',
    defaultValue: 'FLWPUBK_TEST-45bd332ee4bdefdcacd6d2513944cd16-X',
  );

  /// Flutterwave's live public key, for card payments in Mainnet mode. No
  /// default: without it card payment is not offered on Mainnet.
  static const String flutterwavePublicKeyMainnet = String.fromEnvironment(
    'FLUTTERWAVE_PUBLIC_KEY_MAINNET',
  );

  static String apiUrl(String? walletMode) =>
      _trimSlash(walletMode == 'Testnet' ? apiUrlTestnet : apiUrlMainnet);

  static String flutterwavePublicKey(String? walletMode) =>
      walletMode == 'Testnet'
          ? flutterwavePublicKeyTestnet
          : flutterwavePublicKeyMainnet;

  static String _trimSlash(String url) =>
      url.endsWith('/') ? url.substring(0, url.length - 1) : url;
}
