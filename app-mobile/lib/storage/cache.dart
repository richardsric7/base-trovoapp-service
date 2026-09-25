import 'dart:developer';

import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/announcement.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';

Future<void> updateUserInfo(
  signer,
  secretKey,
  address,
  username,
  appState, {
  bool forceRefresh = false,
  String? pnt,
}) async {
  String uri = '/v1/users/$username';
  if (forceRefresh) uri += '?type=refresh';

  if (pnt != null) uri += '?type=import&pnt=$pnt';

  Map responseData = await makeGetRequest(
    uri: uri,
    signer: signer,
    secretKey: secretKey, // the primary wallet secret key
    address: address!,
  );
  inspect(responseData['data']);
  if (responseData['statusCode'] == 200) {
    await storeUserInfo(responseData['data'], appState);
    await fetchCuratedSwapList(appState);
  }
}

Future<void> fetchCuratedSwapList(DataProvider appState) async {
  String uri = '/v1/curated-assets/users';

  Map responseData = await makeGetRequest(
    uri: uri,
    signer: appState.primaryWallet.signer ?? "",
    secretKey: appState.secretKeys[0], // the primary wallet secret key
    address: appState.primaryWallet.address ?? "",
  );
  if (responseData['statusCode'] == 200) {
    if (responseData['data'].length > 0) {
      await StoreData().storeInsertData(
        'curatedSwapList',
        responseData['data'],
      );
      appState.curatedSwapList = appState.deserializeSwapList(
        responseData['data'],
      );
      appState.curatedSwapList.forEach((ca) {
        appState.curatedSwapListMap['${ca.contractAddress}|${ca.assetCode}'] = ca;
      });
    }
  }
}

storeUserInfo(userInfoMap, DataProvider state) async {
  var userInfo = userInfoMap['userData'] as Map<String, dynamic>;
  var assetBalances = userInfoMap['assetBalances'] as Map<String, dynamic>;
  var nfts = userInfoMap['nfts'] ?? {};
  var walletsSharedWithUser = userInfoMap['walletsSharedWithUser'] ?? [];
  var defaultAssets = userInfoMap['defaultAssets'] ?? [];

  if (userInfo.isNotEmpty) {
    await StoreData().storeInsertData('userInfo', userInfo);
  }

  if (assetBalances.isNotEmpty) {
    await StoreData().storeInsertData('assetBalances', assetBalances);
  }
  await StoreData().storeInsertData('nfts', nfts);
  await StoreData().storeInsertData(
    'walletsSharedWithUser',
    walletsSharedWithUser,
  );
  await StoreData().storeInsertData('isFirstTime', false);
  await StoreData().storeInsertData('defaultAssets', defaultAssets);
  await StoreData().storeInsertData('restartedAfterSwitch', false);
  // save useInfo to appstate
  state.setUser = UserInfo().deserializeJson(
    userInfo,
    walletsSharedWithUser,
    assetBalances,
  );
  state.activeWallet = state.primaryWallet;
  state.setNFTs = nfts;
  state.setSharedWallets = walletsSharedWithUser;
  state.setassetBalances = assetBalances;
}

Future<void> getFiatRates(appState) async {
  Map responseData = await makeUnSecuredGetRequest('/v1/rates');

  if (responseData['statusCode'] == 200) {
    if (responseData['data'] != null) {
      appState.setFiatRate = responseData['data'];
      await StoreData().storeInsertData('fiatRate', responseData['data']);
    }
  }
}

Future<void> fetchNotifications(DataProvider appState) async {
  var uri = '/v1/announcements';

  Map responseData = await makeUnSecuredGetRequest(Uri.encodeFull(uri));

  if (responseData['statusCode'] == 200) {
    //  get the date when the user viewed announcements last
    DateTime? lastNotificationViewDate = DateTime.tryParse(
      await StoreData().storeGetData('lastNotificationViewDate') ??
          DateTime.now().add(Duration(days: -30)).toIso8601String(),
    );
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
      var announcement = Announcement().deserializeJson(
        responseData['data'][i],
      );
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

  StoreData().storeInsertData('appVersion', versionInfo['data']);
}
