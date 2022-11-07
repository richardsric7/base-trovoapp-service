class Permission {
  DateTime createdAt;
  DateTime updatedAt;
  String walletPublicKey;
  String targetUsername;
  String permission;
  Permission({
    required this.createdAt,
    required this.updatedAt,
    required this.walletPublicKey,
    required this.targetUsername,
    required this.permission,
  });
}
