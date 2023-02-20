// import 'dart:ffi';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/announcement.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';

Future<void> updateUserInfo(signer, secretKey, publicKey, username, appState,
    {bool forceRefresh = false}) async {
  String uri = '/v1/users/$username';
  if (forceRefresh) uri += '?type=refresh';

  Map responseData = await makeGetRequest(
    uri: uri,
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
  // print('userInfoMap: ${userInfoMap['userData']}');
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
  state.setUser = UserInfo()
      .deserializeJson(userInfo, walletsSharedWithUser, assetBalances);
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

  // print('response: ${responseData}');

  if (responseData['statusCode'] == 200) {
    appState.setFiatRate = responseData['data'];
    await StoreData().storeInsertData('fiatRate', responseData['data']);
  }
}

Future<void> fetchNotifications(DataProvider appState) async {
  print('fetching announcements...');
  var uri = '/v1/announcements';

  Map responseData = await makeUnSecuredGetRequest(
    Uri.encodeFull(uri),
  );

  // print('response: ${responseData}');

  if (responseData['statusCode'] == 200) {
    //  get the date when the user viewed announcements last
    DateTime? lastNotificationViewDate = DateTime.tryParse(
        await StoreData().storeGetData('lastNotificationViewDate') ??
            DateTime.now().add(Duration(days: -30)).toIso8601String());
    // get all announcements
    var announcementsMap = await StoreData().storeGetData('announcements');
    if (announcementsMap == null) {
      await StoreData().storeInsertData('announcements', responseData['data']);
      appState.setHasNewAnnouncement = true;
      return;
    }

    // deserialize the existing announcements stored in the phone storage
    var announcements = Announcement().deserializeJsonList(announcementsMap);
    // loop through the announcements
    for (var i = 0; i < responseData['data'].length; i++) {
      // deserialize the current item in this iteration
      var announcement =
          Announcement().deserializeJson(responseData['data'][i]);
      // if there's been any new announcements since the user opened announcements
      // last, then setHasNewAnnouncements to true
      if (lastNotificationViewDate!.isBefore(announcement.createdAt!)) {
        announcements.add(announcement);
        appState.setHasNewAnnouncement = true;
      }
    }

    StoreData().storeInsertData(
      'announcements',
      Announcement().toJSONEncodableList(announcements),
    );
  }
}

Future<void> fetchVersionInfo(DataProvider appState) async {
  var versionInfo = await makeUnSecuredGetRequest('/v1/app-version');

  // print('this is response $versionInfo');

  StoreData().storeInsertData(
    'appVersion',
    versionInfo['data'],
  );
}
