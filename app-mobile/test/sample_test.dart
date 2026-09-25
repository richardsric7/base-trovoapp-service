import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
// import 'package:mocktail/mocktail.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/bottom_bar/bottom_pages/swap_assets.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/storage/state.dart';

import 'bootstrap_test.dart';
import 'mock_data_provider.dart';

void main() {
  test('simple addition test', () {
    final result = 2 + 2;
    expect(result, 4);
  });

  testWidgets('SwapAssets renders correctly', (WidgetTester tester) async {
    TestWidgetsFlutterBinding.ensureInitialized();

    // final mock = MockDataProvider();

    // when(() => mock.walletMode).thenReturn("Testnet");
    // when(() => mock.primaryWallet).thenReturn(
    //   Wallet(
    //     publicKey: "GD73I6JM5RX22XDZCT6RT4SSOYYRYDZI3A6NO4TITDFLZRM6NKMYMNSY",
    //     claimedAssets: [
    //       Asset().deserializeJson({
    //         "contractAddress": "",
    //         "assetCode": "",
    //         "amount": "821672.236766",
    //         "inTrade": {
    //           "sellingLiabilities": "0.0000000",
    //           "buyingLiabilities": "0.0000000",
    //         },
    //         "qrCode":
    //             "https://storage.googleapis.com/trovo-wallet-testnet/trovo-wallet-testnet/9b3cc40b-c439-4b63-a929-0fb1745c4b9c1903206415.png",
    //         "imageUrl":
    //             "https://drive.google.com/uc?export=view&id=103fw13pcBoCO2hkTPFX73BUKeWWkVpGZ",
    //         "usdPrice": "0.001",
    //         "nativePrice": "1",
    //         "cryptoWalletDepositAddresses": [],
    //         "closedGroup": "",
    //         "quoteCurrency": "",
    //         "tokenizedAsset": 0,
    //         "fundingStructure": 0,
    //         "exitWithFiat": 0,
    //       }),
    //       Asset().deserializeJson({
    //         "contractAddress":
    //             "GCS7P6422J2MHBCOTPZMMP6D65UP4TX2JX3INNK7RKPGFLOIC2PO37RU",
    //         "assetCode": "GLEN",
    //         "amount": "85.0701829",
    //         "inTrade": {
    //           "sellingLiabilities": "0.0000000",
    //           "buyingLiabilities": "0.0000000",
    //         },
    //         "qrCode":
    //             "https://storage.googleapis.com/trovo-wallet-testnet/trovo-wallet-testnet/454fedd3-9b6d-4f48-86c7-af19e66de7cd604482872.png",
    //         "imageUrl":
    //             "https://storage.googleapis.com/trovo-wallet-testnet/trovo-wallet-testnet/391a4475-24e9-447d-a856-e69628699fe9thundeyy-logo-0613f0ed-2a97-4d39-859f-a24c4f39aba4.jpg",
    //         "usdPrice": "23.51",
    //         "nativePrice": "23.51",
    //         "cryptoWalletDepositAddresses": [],
    //         "closedGroup": "",
    //         "quoteCurrency": "CNGN",
    //         "tokenizedAsset": 1,
    //         "fundingStructure": 0,
    //         "exitWithFiat": 0,
    //       }),
    //     ],
    //   ),
    // );

    // await bootstrapTestApp(
    //   tester: tester,
    //   child: SwapAssets(),
    //   // dataProvider: mock,
    //   colorNotifier: ColorNotifier(),
    // );

    expect(find.text("Swap From"), findsOneWidget);

    await tester.tap(
      find.widgetWithText(DropdownButtonFormField<String>, "Choose asset"),
    );
    await tester.pump();

    expect(find.text('GAS'), findsOneWidget);
  });
}
