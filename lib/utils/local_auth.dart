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
    return await _localAuthentication.authenticate(
        localizedReason: 'Please authenticate to complete this action',
        options: AuthenticationOptions(
          biometricOnly: true,
          useErrorDialogs: true, // show error in dialog
          stickyAuth: true, // native process
        ));
  }

  Future<bool> canCheckBiometrics() async =>
      _localAuthentication.canCheckBiometrics;
}
