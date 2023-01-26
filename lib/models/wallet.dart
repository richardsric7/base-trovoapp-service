import 'asset.dart';
import 'permission.dart';

class Wallet {
  DateTime? createdAt;
  String? publicKey;
  String? secretKey;
  String? tag;
  String? description;
  String? alias;
  String? signer;
  String? userId;
  int? sharedAccessEnabled;
  int? walletType;
  int? walletThreshold;
  int? numberOfApprovalsNeeded;
  int? primaryWallet;
  DateTime? sharedAccessCreatedAt;
  DateTime? sharedAccessUpdatedAt;
  List<Permission>? permissions;
  // not part of the user data but needed for
  // when a user has more than one access on
  // a shared wallet i.e. [APPROVER, INITITOR]
  List<String>? accesses; // singular, for shared wallet
  String? permission; // singular, for shared wallet
  String? owner; // for shared wallet
  List<Asset>? claimedAssets;
  List<Asset>? unClaimedAssets;

  Wallet({
    this.createdAt,
    this.publicKey,
    this.secretKey,
    this.tag,
    this.description,
    this.alias,
    this.signer,
    this.userId,
    this.sharedAccessEnabled,
    this.walletType,
    this.walletThreshold,
    this.numberOfApprovalsNeeded,
    this.primaryWallet,
    this.permissions,
    this.permission,
    this.accesses,
    this.sharedAccessCreatedAt,
    this.sharedAccessUpdatedAt,
    this.owner,
    this.claimedAssets,
    this.unClaimedAssets,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "createdAt": createdAt!.toIso8601String(),
      "publicKey": publicKey,
      "secretKey": secretKey,
      "tag": tag,
      "description": description,
      "alias": alias,
      "signer": signer,
      "userId": userId,
      "sharedAccessEnabled": sharedAccessEnabled,
      "walletType": walletType,
      "walletThreshold": walletThreshold,
      "numberOfApprovalsNeeded": numberOfApprovalsNeeded,
      "primaryWallet": primaryWallet,
      "sharedAccessCreatedAt": sharedAccessCreatedAt!.toIso8601String(),
      "sharedAccessUpdatedAt": sharedAccessUpdatedAt!.toIso8601String(),
      "permissions": permissions,
      "permission": permission,
      "owner": owner,
    };
  }

  Wallet deserializeJson(Map<String, dynamic> m, assetBalances) {
    return Wallet(
      createdAt: DateTime.parse(m["createdAt"]),
      publicKey: m["publicKey"],
      secretKey: m["secretKey"],
      tag: m["tag"],
      description: m["description"],
      alias: m["alias"],
      signer: m["signer"],
      userId: m["userId"],
      sharedAccessEnabled: m["sharedAccessEnabled"],
      walletType: m["walletType"],
      walletThreshold: m["walletThreshold"],
      numberOfApprovalsNeeded: m["numberOfApprovalsNeeded"],
      primaryWallet: m["primaryWallet"],
      sharedAccessCreatedAt: DateTime.parse(m["sharedAccessCreatedAt"]),
      sharedAccessUpdatedAt: DateTime.parse(m["sharedAccessUpdatedAt"]),
      permissions: getPermissionList(m["permissions"]),
      claimedAssets:
          deserializeAssetList(assetBalances[m["publicKey"]]['claimed']),
      unClaimedAssets:
          deserializeAssetList(assetBalances[m["publicKey"]]['unclaimed']),
    );
  }

  bool get isInitiator => accesses!.contains('INITIATOR');

  bool get isSharedWallet => sharedAccessEnabled == 1;

  bool get isPrimaryWallet => primaryWallet == 1;

  Wallet deserializeSharedJson(m, List<String> accesses) {
    return Wallet(
        publicKey: m["walletPublicKey"],
        alias: m["walletAlias"],
        permission: m["permission"],
        accesses: accesses,
        description: m["walletDescription"],
        owner: m["owner"],
        sharedAccessEnabled: 1,
        numberOfApprovalsNeeded: m["walletSettings"] != null
            ? m["walletSettings"]["numberOfApprovalsNeeded"]
            : null,
        walletType: m["walletSettings"] != null
            ? m["walletSettings"]["walletType"]
            : null,
        walletThreshold: m["walletSettings"] != null
            ? m["walletSettings"]["walletThreshold"]
            : null,
        permissions: m["walletSettings"] != null
            ? getPermissionList(m["walletSettings"]["permissions"])
            : null,
        claimedAssets: deserializeAssetList(m["assetBalances"]["claimed"]),
        unClaimedAssets: deserializeAssetList(m["assetBalances"]["unclaimed"]));
  }

  List<Permission> getPermissionList(permissionArrayString) {
    var permissions = <Permission>[];
    for (var i = 0; i < permissionArrayString.length; i++) {
      permissions.add(Permission(
          createdAt: DateTime.parse(permissionArrayString[i]['createdAt']),
          updatedAt: DateTime.parse(permissionArrayString[i]['updatedAt']),
          walletPublicKey: permissionArrayString[i]['walletPublicKey'],
          targetUsername: permissionArrayString[i]['targetUsername'],
          fullName: permissionArrayString[i]['fullName'],
          permission: permissionArrayString[i]['permission']));
    }
    return permissions;
  }

  List<Asset> deserializeAssetList(assets) {
    var assetsList = <Asset>[];
    if (assets != null) {
      for (var i = 0; i < assets.length; i++) {
        assetsList.add(Asset().deserializeJson(assets[i]));
      }
    }
    return assetsList;
  }
}
