import 'package:local_auth/local_auth.dart';

class Authenticator {
  final LocalAuthentication _localAuthentication = LocalAuthentication();

  Future<bool> checkingForBioMetrics() async {
    bool canCheckBiometrics = await _localAuthentication.canCheckBiometrics;
    return canCheckBiometrics;
  }

  //this method opens a dialog for fingerprint authentication.
  //we do not need to create a dialog nut it popsup from device natively.
  Future<bool> authenticateMe() async {
    try {
      return await _localAuthentication.authenticate(
        localizedReason: 'Please authenticate to complete this action',
        biometricOnly: true,
        persistAcrossBackgrounding: true, // native process (was stickyAuth)
        // useErrorDialogs is gone in local_auth 3.x - implementations now
        // always behave as if it were false.
      );
    } catch (e) {
      return false;
    }
  }

  Future<bool> canCheckBiometrics() async =>
      _localAuthentication.canCheckBiometrics;
}
