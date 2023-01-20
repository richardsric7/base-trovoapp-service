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
  bool? _isInitiator;
  int? sharedAccessEnabled;
  int? walletType;
  int? walletThreshold;
  int? numberOfApprovalsNeeded;
  int? primaryWallet;
  DateTime? sharedAccessCreatedAt;
  DateTime? sharedAccessUpdatedAt;
  List<Permission>? permissions;
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
    this.sharedAccessCreatedAt,
    this.sharedAccessUpdatedAt,
    this.permission,
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

  bool get isInitiator {
    return permission == "INITIATOR";
  }

  bool get isSharedWallet {
    return sharedAccessEnabled == 1;
  }

  Wallet deserializeSharedJson(m) {
    return Wallet(
        publicKey: m["walletPublicKey"],
        alias: m["walletAlias"],
        permission: m["permission"],
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
