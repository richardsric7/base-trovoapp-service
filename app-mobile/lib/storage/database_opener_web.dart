import 'package:sembast_web/sembast_web.dart';

// dart:io (and therefore path_provider/sembast_io) has no filesystem to open
// a database file on for the web target. sembast_web is a real, working
// alternative backed by IndexedDB - not a stub - so this gives Flutter Web
// genuinely functional local storage, matching what native builds get.
Future<Database> openAppDatabase() async {
  return databaseFactoryWeb.openDatabase('trovoWallet.db');
}
