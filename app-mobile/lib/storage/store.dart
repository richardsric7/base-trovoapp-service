import 'dart:async';
import 'package:sembast/sembast.dart';
import 'database_opener.dart';

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
    final database = await openAppDatabase();

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
    try {
      var data = await store.record(key).get(await _db);

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
