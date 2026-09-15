class PatronTier {
  int id;
  String tier;
  int canExpire;
  int inactive;
  int price;

  PatronTier({
    required this.id,
    required this.tier,
    required this.canExpire,
    required this.inactive,
    required this.price,
  });

  bool get isExpirable => this.canExpire == 1;
}
