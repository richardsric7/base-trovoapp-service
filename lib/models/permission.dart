class Permission {
  DateTime? createdAt;
  DateTime? updatedAt;
  String? walletPublicKey;
  String targetUsername;
  String fullName;
  String permission;
  PermissionState? permissionState;
  Permission(
      {required this.targetUsername,
      required this.fullName,
      required this.permission,
      this.createdAt,
      this.updatedAt,
      this.walletPublicKey,
      this.permissionState});
  Permission.clone(Permission original)
      : this(
          targetUsername: original.targetUsername,
          fullName: original.fullName,
          permission: original.permission,
          createdAt: original.createdAt,
          updatedAt: original.updatedAt,
          walletPublicKey: original.walletPublicKey,
          permissionState: original.permissionState,
        );
}

enum PermissionState { Revoked, Modified, Added }
