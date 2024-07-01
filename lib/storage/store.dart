import 'dart:async';
import 'package:path/path.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sembast/sembast.dart';
import 'package:sembast/sembast_io.dart';

class AppDatabase {
// Singleton instance
  static final AppDatabase _singleton = AppDatabase._();

// Singleton accessor
  static AppDatabase get instance => _singleton;

  // Completer is used for transforming synchronous code into asynchronous code.
  Completer<Database>? _dbOpenCompleter;

  // A private constructor. Allows us to create instances of AppDatabase
  // only from within the AppDatabase class itself.
  AppDatabase._();

  // Database object accessor
  Future<Database> get database async {
    // If completer is null, AppDatabaseClass is newly instantiated, so database is not yet opened
    if (_dbOpenCompleter == null) {
      _dbOpenCompleter = Completer();
      // Calling _openDatabase will also complete the completer with database instance
      _openDatabase();
    }
    // If the database is already opened, awaiting the future will happen instantly.
    // Otherwise, awaiting the returned future will take some time - until complete() is called
    // on the Completer in _openDatabase() below.
    return _dbOpenCompleter!.future;
  }

  Future _openDatabase() async {
    // Get a platform-specific directory where persistent app data can be stored
    final appDocumentDir = await getApplicationDocumentsDirectory();
    // Path with the form: /platform-specific-directory/demo.db
    final dbPath = join(appDocumentDir.path, 'trovoWallet.db');

    final database = await databaseFactoryIo.openDatabase(dbPath);

    // Any code awaiting the Completer's future will now start executing
    _dbOpenCompleter!.complete(database);
  }
}

class StoreData {
  var store = StoreRef.main();

  Future<Database> get _db async => await AppDatabase.instance.database;

  storeInsertData(key, value) async {
    await store.record(key).put(await _db, value, merge: true);

    print('$key data Inserted successfully !!');
  }

  storeGetData(key) async {
    // var count = await store.count(await _db);
    // print("Store count is $count");
    try {
      var data = await store.record(key).get(await _db);
      //print('$key data retrived successfully !!');

      return data;
    } catch (e) {
      throw e;
    }
  }

  storeDeleteItem(String item) async {
    await store.record(item).delete(await _db);

    print('$item deleted successfully !!');
  }

  storeDeleteData() async {
    await store.delete(await _db);

    print('data deleted successfully !!');
  }
}
