import 'patronTier.dart';

class PatronInfo {
  int id;
  String patronPackage;
  List<PatronTier> patronTiers;
  String description;

  PatronInfo({
    required this.id,
    required this.patronPackage,
    required this.patronTiers,
    required this.description,
  });

  PatronInfo.deserializeJson(m, int canExpire, int isInactive)
      : this(
          id: m['id'],
          patronPackage: m['patronPackage'],
          patronTiers: m['patronTier'],
          description: '',
        );
}
