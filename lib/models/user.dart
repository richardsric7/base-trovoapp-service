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
  List<Map<String, dynamic>>? curatedSwapList;

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

  deserializeJson(Map<String, dynamic> m) {
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
      wallets: deserializeWallets(m),
    );
  }

  List<Map<String, dynamic>> deserializeSwapList(Map<String, dynamic> m) {
    List<Map<String, dynamic>> list = [];
    m['curatedSwapList'].forEach((item) {
      list.add({
        'assetIssuer': item['assetIssuer'],
        'assetCode': item['assetCode'],
        'assetName': item['assetName'],
        'description': item['description'],
        'imageUrl': item['imageUrl'],
        'website': item['website'],
        'assetConditions': item['assetConditions'],
        'assetLimit': item['assetLimit'],
        'assetRedemptionInstructions': item['assetRedemptionInstructions'],
        'contactEmail': item['contactEmail'],
        'assetClassId': item['assetClassId'],
        'assetClass': item['assetClass'],
        'organization': item['organization'],
        'withdrawable': item['withdrawable'],
        'decimalPlaces': item['decimalPlaces'],
        'realAssetImageUrl': item['realAssetImageUrl'],
        'generateDepositAddress': item['generateDepositAddress'],
      });
    });
    return list;
  }

  List<Wallet> deserializeWallets(Map<String, dynamic> m) {
    var userWallets = m['userWallets'];
    var myWallets = <Wallet>[];
    if (userWallets != null) {
      for (var i = 0; i < userWallets.length; i++) {
        var wallet = Wallet().deserializeJson(userWallets[i]);
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
}
