# Links wallet-core (keys and signing) into the app once it is built with
# scripts/build_wallet_core.sh ios. The Podfile only adds this pod when
# WalletCore.xcframework exists.
Pod::Spec.new do |s|
  s.name             = 'WalletCore'
  s.version          = '0.1.0'
  s.summary          = "Trovo's wallet-core: keys, recovery phrases, signing and Safe addresses."
  s.homepage         = 'https://github.com/richardsric7/base-trovoapp-service'
  s.license          = { :type => 'Proprietary' }
  s.author           = 'Trovo'
  s.source           = { :path => '.' }
  s.platform         = :ios, '12.0'
  s.static_framework = true
  s.vendored_frameworks = 'WalletCore.xcframework'
  # Dart calls these functions through FFI (DynamicLibrary.process()), so no
  # Swift code references them: -force_load keeps the linker from dropping
  # them from the app.
  s.user_target_xcconfig = {
    'OTHER_LDFLAGS' => '$(inherited) -force_load "$(PODS_XCFRAMEWORKS_BUILD_DIR)/WalletCore/libwallet_core.a"'
  }
end
