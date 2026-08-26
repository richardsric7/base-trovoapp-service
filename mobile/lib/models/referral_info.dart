class ReferralInfoObject {
  Map<String, dynamic> downlines;
  Map<String, dynamic> uplines;

  ReferralInfoObject({required this.downlines, required this.uplines});

  ReferralInfoObject.deserializeJson(Map<String, dynamic> m)
      : this(
          downlines: m['downlines'],
          uplines: m['uplines'],
        );
}
