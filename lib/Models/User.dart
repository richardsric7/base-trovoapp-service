import 'Wallet.dart';

class UserInfo {
  String? username;
  String? firstName;
  String? lastName;
  int? hasSecretQuestions;
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
    this.hasSecretQuestions,
    this.accountRecoveryEnabled,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "username": username,
      "firstName": firstName,
      "lastName": lastName,
      "email": email,
      "mobile": mobile,
      "mobileVerified": mobileVerified,
      "hasSecretQuestions": hasSecretQuestions,
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
    };
  }

  deserializeJson(Map<String, dynamic> m) {
    var userWallets = m['userWallets'];
    var wallets = <Wallet>[];
    if (userWallets != null) {
      for (var i = 0; i < userWallets.length; i++) {
        var wallet = Wallet().deserializeJson(userWallets[i]);
        if (wallet.primaryWallet == 1) {
          // promote the primary wallet to appear first on the list
          wallets.insert(0, wallet);
          continue;
        }
        wallets.add(wallet);
      }
      wallets.forEach((wallet) {
        print('${wallet.publicKey} ${wallet.primaryWallet}');
      });
    }
    return UserInfo(
      username: m['username'],
      firstName: m['firstName'],
      lastName: m['lastName'],
      email: m['email'],
      mobile: m['mobile'],
      mobileVerified: m['mobileVerified'],
      hasSecretQuestions: m['hasSecretQuestions'],
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
      wallets: wallets,
    );
  }
}
