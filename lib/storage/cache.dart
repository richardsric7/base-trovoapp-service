import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/store.dart';

Future<void> updateUserInfo(
    signer, secretKey, publicKey, username, appState) async {
  Map responseData = await makeGetRequest(
    uri: '/v1/users/$username',
    signer: signer,
    secretKey: secretKey, // the primary wallet secret key
    publicKey: publicKey!,
  );

  print('response: ${responseData}');

  if (responseData['statusCode'] == 200) {
    await storeUserInfo(responseData['data'], appState);
  }
}

storeUserInfo(userInfoMap, state) async {
  print('userInfoMap: ${userInfoMap['userData']}');
  var userInfo = userInfoMap['userData'] ?? {};
  var assetBalances = userInfoMap['assetBalances'] ?? {};
  var nfts = userInfoMap['nfts'] ?? {};
  var walletsSharedWithUser = userInfoMap['walletsSharedWithUser'] ?? [];
  var defaultAssets = userInfoMap['defaultAssets'] ?? [];

  await StoreData().storeInsertData('userInfo', userInfo);
  await StoreData().storeInsertData('assetBalances', assetBalances);
  await StoreData().storeInsertData('nfts', nfts);
  await StoreData()
      .storeInsertData('walletsSharedWithUser', walletsSharedWithUser);
  await StoreData().storeInsertData('isFirstTime', false);
  await StoreData().storeInsertData('defaultAssets', defaultAssets);

  // save useInfo to appstate
  state.setUser = UserInfo().deserializeJson(userInfo);
  state.setNFTs = nfts;
  state.setSharedWallets = walletsSharedWithUser;
  state.setassetBalances = assetBalances;
}

Future<void> getFiatRates(
    signer, secretKey, publicKey, username, appState) async {
  Map responseData = await makeGetRequest(
    uri: '/v1/rates',
    signer: signer,
    secretKey: secretKey, // the primary wallet secret key
    publicKey: publicKey!,
  );

  print('response: ${responseData}');

  if (responseData['statusCode'] == 200) {
    appState.setFiatRate = responseData['data'];
    await StoreData().storeInsertData('fiatRate', responseData['data']);
  }
}
