import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sembast/sembast_io.dart';

Future<Database> openAppDatabase() async {
  // Get a platform-specific directory where persistent app data can be stored
  final appDocumentDir = await getApplicationDocumentsDirectory();
  // Path with the form: /platform-specific-directory/demo.db
  final dbPath = join(appDocumentDir.path, 'trovoWallet.db');

  return databaseFactoryIo.openDatabase(dbPath);
}
