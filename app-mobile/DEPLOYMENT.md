# Deployment Guide

This guide assumes you have **never built or released a Flutter app
before**. It covers: installing the tools, running the app in debug mode,
running static analysis and tests, and producing signed release builds for
Android and iOS.

All commands below are run from the `app-mobile/` directory unless said
otherwise.

## 1. Install the Flutter SDK

`pubspec.yaml` requires Dart `>=3.8.0`, which ships with Flutter 3.32+.
**Install a current Flutter release on the `stable` channel** (3.32 or
newer); developers pin the exact version with [FVM](https://fvm.app/).

1. Install Flutter following the official instructions for your OS:
   https://docs.flutter.dev/get-started/install
   - Pick the **stable** channel.
   - If the installer gives you the latest stable release and it's newer
     than 3.32.x, that's usually fine for day-to-day work, but if you hit
     analyzer/build differences from other developers, pin the team's version
     (e.g. via [FVM](https://fvm.app/), which lets you pin an exact
     Flutter version per project: `fvm install 3.32.0 && fvm use 3.32.0`).
2. Add Flutter to your `PATH` (the installer's instructions cover this
   per-OS).
3. Verify the install and let Flutter tell you what else is missing:

   ```bash
   flutter --version
   flutter doctor
   ```

   Run `flutter doctor` again after each step below — it's the
   authoritative "what's still missing" checklist.

## 2. Install platform tooling

This project has both an `android/` and an `ios/` folder, so it targets
both platforms.

### Android (any OS)

1. Install **Android Studio**: https://developer.android.com/studio
2. On first launch, use its **SDK Manager** to install an Android SDK
   platform and the Android SDK command-line tools (`flutter doctor` will
   tell you if anything's missing, including license acceptance —
   `flutter doctor --android-licenses`).
3. Install a JDK if you don't already have one — **Temurin/OpenJDK 17**
   is the version the Android build is set up for.
4. Create an emulator via Android Studio's **Device Manager**, or plug in
   a physical Android device with USB debugging enabled.

### iOS (macOS only — Xcode does not run on Windows/Linux)

1. Install **Xcode** from the Mac App Store.
2. Open it once to accept the license and let it install additional
   components.
3. Install CocoaPods, which Flutter uses to manage this project's iOS
   dependencies (`ios/Podfile`):

   ```bash
   sudo gem install cocoapods
   ```

4. To run on the iOS **Simulator**, no Apple account is needed. To run on
   a **physical iPhone** or produce a release build, you need an Apple
   Developer account and code signing set up in Xcode — see
   [§6 iOS release build](#6-ios-release-build-ipa) below.

If you're on Windows or Linux, you can still do all Android work; skip
the iOS sections.

## 3. Get project dependencies

```bash
cd app-mobile
flutter pub get
```

This reads `pubspec.yaml`/`pubspec.lock` and downloads all Dart/Flutter
packages. Re-run it any time you pull changes that touch `pubspec.yaml`.

## 4. Run in debug mode

1. Start an Android emulator or iOS Simulator, or connect a physical
   device.
2. Confirm Flutter sees it:

   ```bash
   flutter devices
   ```

3. Run the app with hot reload:

   ```bash
   flutter run
   ```

   Or, in VS Code: open the Run and Debug panel and press F5 (this repo's
   original README already documented this shortcut). In Android Studio:
   use the Run ▶ button with your device selected.

## 5. Static analysis and tests

There is no CI pipeline in this repository (the inherited workflows were
removed - see the root `ARCHITECTURE.md`); run these yourself after
`flutter pub get`:

```bash
# Analyze — fail only on errors, not warnings/infos
# (there's a known backlog of MaterialState*->WidgetState* and
# Share->SharePlus deprecation infos that don't block PRs today)
flutter analyze --no-fatal-warnings --no-fatal-infos

# Run the test suite (see test/)
flutter test
```

Run both before opening a PR. Plain `flutter analyze` (no flags) is
stricter and worth running too.

## 6. Build a release Android APK/AppBundle

### 6.1 Generate a signing keystore (one-time, per release identity)

There is no `key.properties.example` or checked-in keystore in this repo
— `android/key.properties` and `android/**/local.properties` are
git-ignored (see `.gitignore`). For release builds, generate your own
keystore once with the JDK's `keytool`:

```bash
keytool -genkey -v -keystore ~/trovo-release-key.jks \
  -keyalg RSA -keysize 2048 -validity 10000 \
  -alias trovo-release
```

You'll be prompted for a keystore password, your name/org details, and a
key password. **Keep the resulting `.jks` file and passwords private —
never commit them.**

### 6.2 Point the build at your keystore

Create `android/key.properties` (this exact path — `android/app/build.gradle.kts`
reads it from `rootProject.file("key.properties")`, i.e. `android/key.properties`):

```properties
storePassword=<your keystore password>
keyPassword=<your key password>
keyAlias=trovo-release
storeFile=<path to your .jks file, e.g. /Users/you/trovo-release-key.jks>
```

This file is already covered by `.gitignore` — do not remove that entry,
and never commit real passwords.

Without a `key.properties` file, `flutter build apk`/`appbundle` will
still produce a build, but Gradle's `signingConfig` block
(`android/app/build.gradle.kts`) will fail to resolve a keystore for the
`release` build type — create the file first.

### 6.3 Build

```bash
# Android App Bundle — what Google Play expects
flutter build appbundle

# or, a plain installable APK
flutter build apk
```

Output lands under `build/app/outputs/bundle/release/` (AAB) or
`build/app/outputs/flutter-apk/` (APK).

## 7. iOS release build (`.ipa`)

This requires a **Mac with Xcode**, an **Apple Developer Program
account**, and code signing configured in Xcode. Full Apple signing setup
(certificates, provisioning profiles, App Store Connect app records) is
its own large topic — Apple's own documentation is the source of truth,
so this section describes the shape of the process rather than exact,
unverified menu paths:

1. Enroll in the **Apple Developer Program** if you haven't already
   (apple.com), and have access to the relevant Apple Developer account
   / App Store Connect app record for `com.trovo.wallet` (the bundle ID
   set in `ios/Runner.xcodeproj`).
2. Open `ios/Runner.xcworkspace` in Xcode (not the `.xcodeproj` — Flutter
   projects use CocoaPods, so the workspace is the file to open).
3. In the Runner target's **Signing & Capabilities** tab, set your team
   and let Xcode manage signing (simplest path for a first build), or
   configure a manual provisioning profile if your team uses one.
4. Build the IPA:

   ```bash
   flutter build ipa
   ```

   Output lands under `build/ios/ipa/`.
5. Upload to TestFlight / App Store Connect via Xcode's Organizer, or
   `xcrun altool` / `xcrun notarytool`, per Apple's current upload
   instructions.

There is no automated iOS release pipeline in this repository (the
inherited GitHub Actions workflow and fastlane setup were removed).

## 8. Switching Flutter versions / stale build artifacts

This project has a `build/` directory at its root (git-ignored — it's
Flutter's build output, not source). If you switch Flutter SDK versions,
or hit strange build errors after a Flutter upgrade/downgrade, clear
cached build state first:

```bash
flutter clean
flutter pub get
```

`flutter clean` removes `build/`, `.dart_tool/`, and other generated
directories so the next build starts fresh.
