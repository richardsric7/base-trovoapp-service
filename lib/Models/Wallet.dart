import 'Permission.dart';

class Wallet {
  DateTime? createdAt;
  String? publicKey;
  String? secretKey;
  String? tag;
  String? description;
  String? alias;
  String? signer;
  String? userId;
  int? managedAccessEnabled;
  int? primaryWallet;
  List<Permission>? permissions;

  Wallet(
      {this.createdAt,
      this.publicKey,
      this.secretKey,
      this.tag,
      this.description,
      this.alias,
      this.signer,
      this.userId,
      this.managedAccessEnabled,
      this.primaryWallet,
      this.permissions});

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
      "managedAccessEnabled": managedAccessEnabled,
      "primaryWallet": primaryWallet,
      "permissions": permissions,
    };
  }

  Wallet deserializeJson(Map<String, dynamic> m) {
    return Wallet(
      createdAt: DateTime.parse(m["createdAt"]),
      publicKey: m["publicKey"],
      secretKey: m["secretKey"],
      tag: m["tag"],
      description: m["description"],
      alias: m["alias"],
      signer: m["signer"],
      userId: m["userId"],
      managedAccessEnabled: m["managedAccessEnabled"],
      primaryWallet: m["primaryWallet"],
      permissions: getPermissionList(m["permissions"]),
    );
  }

  List<Permission> getPermissionList(permissionArrayString) {
    var permissions = <Permission>[];
    for (var i = 0; i < permissionArrayString.length; i++) {
      permissions.add(Permission(
          createdAt: DateTime.parse(permissionArrayString[i]['createdAt']),
          updatedAt: DateTime.parse(permissionArrayString[i]['updatedAt']),
          walletPublicKey: permissionArrayString[i]['walletPublicKey'],
          targetUsername: permissionArrayString[i]['targetUsername'],
          permission: permissionArrayString[i]['permission']));
    }
    return permissions;
  }
}
