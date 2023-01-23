import 'package:trovo_wallet/models/curated_asset.dart';

import 'wallet.dart';

class UserInfo {
  String? username;
  String? firstName;
  String? lastName;
  int? hasSecurityQuestions;
  int? accountRecoveryEnabled;
  String? email;
  String? mobile;
  int? mobileVerified;
  String? countryCode;
  String? referrer;
  String? referralLink;
  String? referralQRCode;
  String? publicKey;
  int? corporate;
  String? pushNotificationToken;
  String? imageThumbnailURL;
  int? membershipType;
  DateTime? membershipExpiry;
  int? kycVerified;
  int? walletRecoveryEnabled;
  int? verified;
  int? suspended;
  List<Wallet>? wallets;
  List<Wallet>? sharedWallets;
  List<CuratedAsset>? curatedSwapList;

  UserInfo({
    this.username,
    this.firstName,
    this.lastName,
    this.email,
    this.mobile,
    this.publicKey,
    this.mobileVerified,
    this.countryCode,
    this.referrer,
    this.referralLink,
    this.referralQRCode,
    this.corporate,
    this.pushNotificationToken,
    this.imageThumbnailURL,
    this.membershipType,
    this.membershipExpiry,
    this.kycVerified,
    this.walletRecoveryEnabled,
    this.verified,
    this.suspended,
    this.wallets,
    this.sharedWallets,
    this.hasSecurityQuestions,
    this.accountRecoveryEnabled,
    this.curatedSwapList,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "username": username,
      "firstName": firstName,
      "lastName": lastName,
      "email": email,
      "mobile": mobile,
      "mobileVerified": mobileVerified,
      "hasSecurityQuestions": hasSecurityQuestions,
      "accountRecoveryEnabled": accountRecoveryEnabled,
      "countryCode": countryCode,
      "referrer": referrer,
      "referralLink": referralLink,
      "referralQRCode": referralQRCode,
      "publicKey": publicKey,
      "corporate": corporate,
      "pushNotificationToken": pushNotificationToken,
      "imageThumbnailURL": imageThumbnailURL,
      "membershipType": membershipType,
      "membershipExpiry": membershipExpiry!.toIso8601String(),
      "kycVerified": kycVerified,
      "walletRecoveryEnabled": walletRecoveryEnabled,
      "verified": verified,
      "suspended": suspended,
      "curatedSwapList": curatedSwapList,
    };
  }

  deserializeJson(Map<String, dynamic> m, sharedWallets, assetBalances) {
    return UserInfo(
        username: m['username'],
        firstName: m['firstName'],
        lastName: m['lastName'],
        email: m['email'],
        mobile: m['mobile'],
        mobileVerified: m['mobileVerified'],
        hasSecurityQuestions: m['hasSecurityQuestions'],
        accountRecoveryEnabled: m['accountRecoveryEnabled'],
        countryCode: m['countryCode'],
        referrer: m['referrer'],
        referralLink: m['referralLink'],
        referralQRCode: m['referralQRCode'],
        publicKey: m['publicKey'],
        corporate: m['corporate'],
        pushNotificationToken: m['pushNotificationToken'],
        imageThumbnailURL: m['imageThumbnailURL'],
        membershipType: m['membershipType'],
        membershipExpiry: DateTime.tryParse(m['membershipExpiry']),
        kycVerified: m['kycVerified'],
        walletRecoveryEnabled: m['walletRecoveryEnabled'],
        verified: m['verified'],
        suspended: m['suspended'],
        curatedSwapList: deserializeSwapList(m),
        wallets: deserializeWallets(m, assetBalances),
        sharedWallets: deserializeSharedWallets(sharedWallets));
  }

  List<CuratedAsset> deserializeSwapList(Map<String, dynamic> m) {
    List<CuratedAsset> list = [];
    m['curatedSwapList'].forEach((item) {
      list.add(CuratedAsset().deserializeJson(item));
    });
    return list;
  }

  List<Wallet> deserializeWallets(Map<String, dynamic> m, assetBalances) {
    var userWallets = m['userWallets'];
    var myWallets = <Wallet>[];
    if (userWallets != null) {
      for (var i = 0; i < userWallets.length; i++) {
        var wallet = Wallet().deserializeJson(userWallets[i], assetBalances);
        if (wallet.primaryWallet == 1) {
          // promote the primary wallet to appear first on the list
          myWallets.insert(0, wallet);
          continue;
        }
        myWallets.add(wallet);
      }
    }
    return myWallets;
  }

  List<Wallet> deserializeSharedWallets(m) {
    var myWallets = <Wallet>[];
    for (var i = 0; i < m.length; i++) {
      // get the permission in this current wallet object
      // and make a new list with it
      var accesses = <String>[m[i]['permission']];
      // loop through each of the deserialized wallet object
      myWallets = myWallets.where((item) {
        // if an already deserialized wallet has the same alias as the current
        // wallet object and the accesses list of the deserialized wallet is not
        // empty
        if (item.alias == m[i]['walletAlias'] && item.accesses!.isNotEmpty) {
          // merge the access list on the already deserialized wallet to new
          // we created for the current wallet object
          accesses.addAll(item.accesses!);
          // remove this deserialized wallet from the myWallets list
          return false;
        }
        // return true if no deserialized wallet has the same alias as the current wallet object
        return true;
      }).toList();

      var wallet = Wallet().deserializeSharedJson(m[i], accesses);
      myWallets.add(wallet);
    }
    return myWallets;
  }

  Wallet getWallet(String publicKey) {
    var combinedList = [...wallets!, ...sharedWallets!];
    return combinedList.firstWhere((wallet) => wallet.publicKey == publicKey);
  }

  List<Wallet> get allWallets => [...wallets!, ...sharedWallets!];

  List<Wallet> transactionableWallets() {
    List<Wallet> transWallets = [];
    for (var wallet in wallets!) {
      if (wallet.walletType == 0) {
        if (wallet.walletThreshold == 2 &&
            wallet.permissions!
                .where((perm) =>
                    perm.permission == 'INITIATOR' &&
                    perm.targetUsername == username)
                .isEmpty) {
          continue;
        }

        transWallets.add(wallet);
      }
    }

    sharedWallets!.forEach((wallet) {
      if (wallet.accesses!.contains('INITIATOR')) {
        transWallets.add(wallet);
      }
    });

    return transWallets;
  }
}
