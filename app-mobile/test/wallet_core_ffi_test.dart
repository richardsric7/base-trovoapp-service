// Exercises WalletCoreFFI's real bindings against the compiled wallet-core
// native library (../../wallet-core - build it with `cargo build --release`
// there first). This is a Dart-VM unit test (`flutter test`), so it runs
// wherever a Linux/macOS/Windows build of libwallet_core can be loaded -
// it does not need a device, emulator, or the Android/iOS build of the
// library the app itself ships.
//
// Requires the built library to be discoverable via the dynamic linker's
// standard search path (WalletCoreFFI opens it by bare filename, matching
// how it's loaded in a real desktop build): run with, e.g.
//   LD_LIBRARY_PATH=../wallet-core/target/release flutter test test/wallet_core_ffi_test.dart
// on Linux. Skips (rather than fails) if the library cannot be loaded, so
// this doesn't break `flutter test` runs where wallet-core hasn't been
// built - see /wallet-core/README.md.
import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';
import 'package:trovo_app/functions/wallet_core_ffi.dart';

/// True if WalletCoreFFI can load the native library at all. Computed once
/// so every test's `skip:` sees the same answer; each test body still
/// calls `WalletCoreFFI.instance` itself (a cheap cached getter after the
/// first call) rather than capturing a nullable local, so the analyzer
/// never has to reason about a captured-and-possibly-null closure variable.
bool _available() {
  try {
    WalletCoreFFI.instance;
    return true;
  } catch (_) {
    return false;
  }
}

void main() {
  final skipReason = _available()
      ? false
      : 'libwallet_core not found on LD_LIBRARY_PATH';

  test(
    'generateKeypair returns a well-formed Base address and matching import',
    () {
      final core = WalletCoreFFI.instance;
      final kp = core.generateKeypair();
      expect(kp.address, startsWith('0x'));
      expect(kp.address.length, 42);
      expect(kp.privateKeyHex, startsWith('0x'));
      expect(kp.privateKeyHex.length, 66);

      final imported = core.keypairFromPrivateKey(kp.privateKeyHex);
      expect(imported.address, kp.address);

      final addressOnly = core.addressFromPrivateKey(kp.privateKeyHex);
      expect(addressOnly, kp.address);
    },
    skip: skipReason,
  );

  test('signPersonal / verifyPersonal round-trip, tampering rejected', () {
    final core = WalletCoreFFI.instance;
    final kp = core.generateKeypair();
    final sig = core.signPersonal(kp.privateKeyHex, 'hello base ffi');

    expect(core.verifyPersonal(kp.address, 'hello base ffi', sig), isTrue);
    expect(core.verifyPersonal(kp.address, 'tampered', sig), isFalse);

    final recovered = core.recoverPersonalSigner('hello base ffi', sig);
    expect(recovered.toLowerCase(), kp.address.toLowerCase());
  }, skip: skipReason);

  test(
    'signPersonalBytes signs a raw binary digest (NUL + high bytes) intact',
    () {
      final core = WalletCoreFFI.instance;
      final kp = core.generateKeypair();
      // A digest-shaped payload containing a NUL byte and bytes >127 - the
      // exact case that would corrupt if round-tripped through a Dart
      // String (see the fix this test guards: signBase64Txn moved off
      // signPersonal onto signPersonalBytes for exactly this reason).
      final digest = List<int>.generate(32, (i) => (i * 37 + 5) % 256);
      final raw = Uint8List.fromList(<int>[0, ...digest, 255, 254, 0]);

      final sig1 = core.signPersonalBytes(kp.privateKeyHex, raw);
      final sig2 = core.signPersonalBytes(kp.privateKeyHex, raw);
      // k256 uses RFC6979 deterministic nonces, so signing the same digest
      // twice must produce the same signature.
      expect(sig1, sig2);
    },
    skip: skipReason,
  );

  test(
    'generateMnemonic / keypairFromMnemonic derivation is deterministic',
    () {
      final core = WalletCoreFFI.instance;
      final mnemonic = core.generateMnemonic();
      expect(mnemonic.split(' ').length, 12);

      final kp1 = core.keypairFromMnemonic(mnemonic, 0);
      final kp2 = core.keypairFromMnemonic(mnemonic, 0);
      expect(kp1.address, kp2.address);

      final kpOtherIndex = core.keypairFromMnemonic(mnemonic, 1);
      expect(kpOtherIndex.address, isNot(kp1.address));
    },
    skip: skipReason,
  );

  test('keypairFromPrivateKey tolerates an uppercased "0X..." input (app-mobile call sites uppercase before parsing)', () {
    final core = WalletCoreFFI.instance;
    final kp = core.generateKeypair();
    final uppercased = kp.privateKeyHex.toUpperCase();
    final imported = core.keypairFromPrivateKey(uppercased);
    expect(imported.address, kp.address);
  }, skip: skipReason);
}
