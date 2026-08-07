import 'patronTier.dart';

class PatronInfo {
  int id;
  String patronPackage;
  List<PatronTier> patronTiers;
  String description;
  String packageListTitle;
  String logo;
  List<String> packageList;

  PatronInfo({
    required this.id,
    required this.patronPackage,
    required this.patronTiers,
    required this.packageListTitle,
    required this.packageList,
    required this.logo,
    required this.description,
  });

  PatronInfo.deserializeJson(m, int canExpire, int isInactive)
      : this(
          id: m['id'],
          patronPackage: m['patronPackage'],
          patronTiers: m['patronTier'],
          packageListTitle: m['packageListTitle'],
          packageList: m['packageList'].toString().split('|'),
          description: '',
          logo: '',
        );
}
