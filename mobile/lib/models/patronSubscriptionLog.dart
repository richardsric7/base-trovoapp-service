class PatronSubscriptionLog {
  String id;
  DateTime createdAt;
  DateTime updatedAt;
  String username;
  String patronPackageId;
  String patronTierId;
  String activePatronPackageId;
  String activePatronTierId;
  DateTime effectiveDate;
  DateTime validTill;

  PatronSubscriptionLog({
    required this.id,
    required this.createdAt,
    required this.updatedAt,
    required this.username,
    required this.patronPackageId,
    required this.patronTierId,
    required this.activePatronPackageId,
    required this.activePatronTierId,
    required this.effectiveDate,
    required this.validTill,
  });

  PatronSubscriptionLog.deserializeJson(m)
      : this(
          id: m['id'],
          createdAt: DateTime.parse(m['createdAt']),
          updatedAt: DateTime.parse(m['updatedAt']),
          username: m['username'],
          patronPackageId: m['patronPackageId'],
          patronTierId: m['patronTierId'],
          activePatronPackageId: m['activePatronPackageId'],
          activePatronTierId: m['activePatronTierId'],
          effectiveDate: DateTime.parse(m['effectiveDate']),
          validTill: DateTime.parse(m['validTill']),
        );
}
