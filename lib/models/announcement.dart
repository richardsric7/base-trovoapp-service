
class Announcement {
  int? id;
  String? message;
  String? title;
  String? broadcastLevel;
  String? level;
  bool isViewed;
  DateTime? expiry;
  DateTime? createdAt;

  Announcement({
    this.id,
    this.message,
    this.title,
    this.broadcastLevel,
    this.level,
    this.expiry,
    this.createdAt,
    this.isViewed = false,
  });

  Map<String, dynamic> toJSONEncodable() {
    return <String, dynamic>{
      "id": id,
      "message": message,
      "title": title,
      "broadcastLevel": broadcastLevel,
      "level": level,
      "isViewed": isViewed,
      "expiry": expiry!.toIso8601String(),
      "createdAt": createdAt!.toIso8601String(),
    };
  }

  Announcement deserializeJson(Map<String, dynamic> m) {
    return Announcement(
      id: m['id'],
      message: m['message'],
      title: m['title'],
      broadcastLevel: m['broadcastLevel'],
      level: m['level'],
      isViewed: m['isViewed'] ?? false,
      expiry: DateTime.tryParse(m['expiry']),
      createdAt: DateTime.tryParse(m['createdAt']),
    );
  }

  List<Announcement> deserializeJsonList(m) {
    List<Announcement> announcements = [];
    for (var i = 0; i < m.length; i++) {
      announcements.add(
        Announcement(
          id: m[i]['id'],
          message: m[i]['message'],
          title: m[i]['title'],
          broadcastLevel: m[i]['broadcastLevel'],
          level: m[i]['level'],
          isViewed: m[i]['isViewed'] ?? false,
          expiry: DateTime.tryParse(m[i]['expiry']),
          createdAt: DateTime.tryParse(m[i]['createdAt']),
        ),
      );
    }

    return announcements;
  }

  toJSONEncodableList(List<Announcement> announcements) {
    var myList = [];
    announcements.forEach((announcement) {
      myList.add(announcement.toJSONEncodable());
    });
    return myList;
  }
}
