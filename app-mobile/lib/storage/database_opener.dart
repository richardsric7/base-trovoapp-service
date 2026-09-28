// Conditional export: sembast_io (used by database_opener_io.dart) depends
// on dart:io for a real filesystem, which doesn't exist on the web target.
// database_opener_web.dart backs the same API with sembast_web (IndexedDB)
// instead - a real implementation, not a stub.
export 'database_opener_web.dart' if (dart.library.io) 'database_opener_io.dart';
