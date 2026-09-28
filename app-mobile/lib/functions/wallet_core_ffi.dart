// Conditional export: dart:ffi (used by wallet_core_ffi_io.dart) cannot
// compile for the web target, so native builds get the real FFI bindings
// and everything else (currently just web) gets a stub that reports
// unsupported instead of failing the whole app's compilation. See
// wallet_core_ffi_web.dart for why this isn't WASM-backed yet.
export 'wallet_core_ffi_web.dart' if (dart.library.io) 'wallet_core_ffi_io.dart';
