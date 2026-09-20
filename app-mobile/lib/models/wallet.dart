import 'package:trovo_app/models/tokenizedAsset.dart';
import 'package:trovo_app/storage/state.dart';

import 'asset.dart';
import 'permission.dart';

class Wallet {
  DateTime? createdAt;
  String? address;
  String? linkedWalletAddress;
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
  List<String>? accesses; // plural, for shared wallet
  String? permission; // singular, for shared wallet
  String? owner; // for shared wallet
  List<Asset>? claimedAssets;
  List<Asset>? unClaimedAssets;

  Wallet({
    this.createdAt,
    this.address,
    this.linkedWalletAddress,
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
      "address": address,
      "linkedWalletAddress": linkedWalletAddress,
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

  Wallet deserializeJson(
    Map<String, dynamic> m,
    assetBalances,
    String username,
  ) {
    return Wallet(
      createdAt: DateTime.parse(m["createdAt"]),
      address: m["address"],
      linkedWalletAddress: m["linkedWalletAddress"],
      secretKey: m["secretKey"],
      tag: m["tag"],
      description: m["description"],
      owner: username,
      alias: m["alias"],
      accesses: getAccesses(m["permissions"], username),
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
      claimedAssets: deserializeAssetList(
        assetBalances[m["address"]]['claimed'],
      ),
      unClaimedAssets: deserializeAssetList(
        assetBalances[m["address"]]['unclaimed'],
      ),
    );
  }

  bool get isInitiator => accesses!.contains('INITIATOR');

  bool get isApprover => accesses!.contains('APPROVER');

  bool get isSharedWallet => sharedAccessEnabled == 1;

  int get totalAssets => claimedAssets?.length ?? 0;

  bool get isSharedWalletAndCanInitiate =>
      sharedAccessEnabled == 1 &&
      walletThreshold == 2 &&
      accesses!.contains('INITIATOR');

  bool get canInitiate =>
      (!isSharedWallet ||
          isInitiator ||
          isPrimaryWallet ||
          walletThreshold == 1) &&
      walletType == 0;

  bool get isPrimaryWallet => primaryWallet == 1;

  Wallet deserializeSharedJson(m, List<String> accesses) {
    return Wallet(
      address: m["walletAddress"],
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
      unClaimedAssets: deserializeAssetList(m["assetBalances"]["unclaimed"]),
    );
  }

  List<String> getAccesses(permissions, String username) {
    var accesses = <String>[];

    if (permissions != null) {
      for (var i = 0; i < permissions.length; i++) {
        if (permissions[i]['targetUsername'] == username &&
            !accesses.contains(permissions[i]['permission'])) {
          accesses.add(permissions[i]['permission']);
        }
      }
    }
    return accesses;
  }

  List<Permission> getPermissionList(permissionArrayString) {
    var permissions = <Permission>[];
    if (permissionArrayString != null) {
      for (var i = 0; i < permissionArrayString.length; i++) {
        permissions.add(
          Permission(
            createdAt: DateTime.parse(permissionArrayString[i]['createdAt']),
            updatedAt: DateTime.parse(permissionArrayString[i]['updatedAt']),
            walletAddress: permissionArrayString[i]['walletAddress'],
            targetUsername: permissionArrayString[i]['targetUsername'],
            fullName: permissionArrayString[i]['fullName'],
            permission: permissionArrayString[i]['permission'],
          ),
        );
      }
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

  List<TokenizedAsset> deserializeTokenizedAssetList(assets) {
    var assetsList = <TokenizedAsset>[];
    if (assets != null) {
      for (var i = 0; i < assets.length; i++) {
        assetsList.add(TokenizedAsset().deserializeJson(assets[i]));
      }
    }
    return assetsList;
  }

  List<Asset> getTokenizedAssets(DataProvider appState) {
    var tokenizedAssets = <Asset>[];
    claimedAssets!.forEach((asset) {
      if (appState
              .curatedSwapListMap['${asset.assetIssuer}|${asset.assetCode}']
              ?.assetClassId ==
          3) {
        tokenizedAssets.add(asset);
      }
    });
    return tokenizedAssets;
  }

  List<Asset> getOtherTokens(DataProvider appState) {
    var tokenizedAssets = <Asset>[];
    claimedAssets!.forEach((asset) {
      if (appState
              .curatedSwapListMap['${asset.assetIssuer}|${asset.assetCode}']
              ?.assetClassId !=
          3) {
        tokenizedAssets.add(asset);
      }
    });
    return tokenizedAssets;
  }
}
