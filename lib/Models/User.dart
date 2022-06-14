class UserInfo {
  String? username;
  String? firstName;
  String? lastName;
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
  String? imageThumbNail;
  int? membershipType;
  DateTime? membershipExpiry;
  int? kycVerified;
  int? walletRecoveryEnabled;
  int? verified;
  int? suspended;

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
    this.imageThumbNail,
    this.membershipType,
    this.membershipExpiry,
    this.kycVerified,
    this.walletRecoveryEnabled,
    this.verified,
    this.suspended,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "username": username,
      "firstName": firstName,
      "lastName": lastName,
      "email": email,
      "mobile": mobile,
      "mobileVerified": mobileVerified,
      "countryCode": countryCode,
      "referrer": referrer,
      "referralLink": referralLink,
      "referralQRCode": referralQRCode,
      "publicKey": publicKey,
      "corporate": corporate,
      "pushNotificationToken": pushNotificationToken,
      "imageThumbNail": imageThumbNail,
      "membershipType": membershipType,
      "membershipExpiry": membershipExpiry,
      "kycVerified": kycVerified,
      "walletRecoveryEnabled": walletRecoveryEnabled,
      "verified": verified,
      "suspended": suspended,
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
      countryCode: m['countryCode'],
      referrer: m['referrer'],
      referralLink: m['referralLink'],
      referralQRCode: m['referralQRCode'],
      publicKey: m['publicKey'],
      corporate: m['corporate'],
      pushNotificationToken: m['pushNotificationToken'],
      imageThumbNail: m['imageThumbNail'],
      membershipType: m['membershipType'],
      membershipExpiry: DateTime.tryParse(m['membershipExpiry']),
      kycVerified: m['kycVerified'],
      walletRecoveryEnabled: m['walletRecoveryEnabled'],
      verified: m['verified'],
      suspended: m['suspended'],
    );
  }
}
