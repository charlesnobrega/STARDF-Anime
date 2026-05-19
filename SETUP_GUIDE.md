# StarDF-Anime Mobile - Setup Guide

## Prerequisites

Before setting up the project, ensure you have the following installed:

### Required Tools
- **Flutter SDK** (3.0 or later)
  - Download from: https://flutter.dev/docs/get-started/install
  - Verify: `flutter --version`

- **Go SDK** (1.21 or later)
  - Download from: https://golang.org/dl
  - Verify: `go version`

- **Android SDK** (API level 26+)
  - Install via Android Studio or command line
  - Set `ANDROID_SDK_ROOT` environment variable

- **GoMobile** (for Go-Android integration)
  - Install: `go install golang.org/x/mobile/cmd/gomobile@latest`
  - Install NDK: `gomobile init`

### Optional Tools
- **Android Studio** (for emulator and debugging)
- **VS Code** with Flutter extension
- **Git** for version control

## Step 1: Environment Setup

### 1.1 Set Environment Variables

**Windows (PowerShell)**:
```powershell
$env:ANDROID_SDK_ROOT = "C:\Android\sdk"
$env:ANDROID_HOME = "C:\Android\sdk"
$env:JAVA_HOME = "C:\Program Files\Android\Android Studio\jbr"
```

**macOS/Linux (bash)**:
```bash
export ANDROID_SDK_ROOT=$HOME/Android/sdk
export ANDROID_HOME=$HOME/Android/sdk
export JAVA_HOME=$ANDROID_SDK_ROOT/../jbr
```

### 1.2 Verify Flutter Setup

```bash
flutter doctor
```

Ensure all required components are installed (Flutter SDK, Android SDK, etc.)

## Step 2: Project Setup

### 2.1 Install Flutter Dependencies

```bash
flutter pub get
```

### 2.2 Generate Code

Generate Drift database code:
```bash
flutter pub run build_runner build
```

Generate JSON serialization code:
```bash
flutter pub run build_runner build
```

### 2.3 Configure Android

Copy and edit local.properties:
```bash
cp android/local.properties.template android/local.properties
```

Edit `android/local.properties`:
```properties
sdk.dir=/path/to/android/sdk
flutter.sdk=/path/to/flutter/sdk
flutter.versionCode=1
flutter.versionName=1.0.0
```

## Step 3: GoMobile Setup

### 3.1 Initialize GoMobile

```bash
gomobile init
```

This installs the Android NDK and sets up GoMobile.

### 3.2 Build Go Bindings

Navigate to the go directory and build the Android library:

```bash
cd go
gomobile bind -target=android -o=../android/app/src/main/jniLibs/libgo.aar .
cd ..
```

This creates the GoMobile bindings that Flutter can call.

## Step 4: Android Emulator Setup

### 4.1 Create Emulator

Using Android Studio:
1. Open Android Studio
2. Tools → Device Manager
3. Create Virtual Device
4. Select device (e.g., Pixel 4)
5. Select API level 26+ (Android 8.0+)
6. Finish

Or using command line:
```bash
sdkmanager "system-images;android-26;google_apis;arm64-v8a"
avdmanager create avd -n android26 -k "system-images;android-26;google_apis;arm64-v8a"
```

### 4.2 Start Emulator

```bash
emulator -avd android26
```

Or from Android Studio: Device Manager → Play button

## Step 5: Run the Application

### 5.1 List Available Devices

```bash
flutter devices
```

### 5.2 Run on Emulator

```bash
flutter run
```

Or specify device:
```bash
flutter run -d emulator-5554
```

### 5.3 Run on Physical Device

1. Enable USB Debugging on device
2. Connect device via USB
3. Run: `flutter run`

## Step 6: Verify Setup

### 6.1 Check Bridge Communication

The app should:
1. Start successfully
2. Initialize the Go backend
3. Start the local web server on port 8080
4. Display the home screen

### 6.2 Test Bridge Calls

In the app, test:
- Fetching animes list
- Adding to watchlist
- Marking episodes as watched
- Checking sync status

### 6.3 Check Logs

View logs:
```bash
flutter logs
```

Look for bridge communication logs starting with `[Bridge]` or `[Server]`

## Step 7: Development Workflow

### 7.1 Hot Reload

During development, use hot reload:
```bash
flutter run
# Press 'r' to hot reload
# Press 'R' to hot restart
```

### 7.2 Run Tests

Unit tests:
```bash
flutter test
```

Integration tests:
```bash
flutter test integration_test/app_test.dart
```

### 7.3 Build APK

Debug APK:
```bash
flutter build apk --debug
```

Release APK:
```bash
flutter build apk --release
```

## Troubleshooting

### Issue: Flutter SDK not found

**Solution**: Set `FLUTTER_HOME` environment variable:
```bash
export FLUTTER_HOME=/path/to/flutter
export PATH=$PATH:$FLUTTER_HOME/bin
```

### Issue: Android SDK not found

**Solution**: Set `ANDROID_SDK_ROOT`:
```bash
export ANDROID_SDK_ROOT=/path/to/android/sdk
```

### Issue: GoMobile not found

**Solution**: Install GoMobile:
```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

### Issue: Bridge calls timeout

**Possible causes**:
- Go backend not started
- Web server not running
- Network connectivity issue

**Solution**:
1. Check logs: `flutter logs`
2. Verify Go backend is running
3. Check network connectivity
4. Increase timeout if needed

### Issue: Emulator too slow

**Solution**:
- Use hardware acceleration (KVM on Linux, HAXM on Windows)
- Allocate more RAM to emulator
- Use physical device for testing

### Issue: Build fails with NDK error

**Solution**:
1. Install NDK: `sdkmanager "ndk;25.1.8937393"`
2. Set NDK path in `android/local.properties`:
   ```properties
   ndk.dir=/path/to/ndk
   ```

## Next Steps

1. Review the [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for architecture overview
2. Start implementing UI screens in `lib/screens/`
3. Implement state management in `lib/providers/`
4. Add database models in `lib/models/`
5. Write tests in `test/` and `integration_test/`
6. Configure CI/CD for automated testing

## Additional Resources

- [Flutter Documentation](https://flutter.dev/docs)
- [GoMobile Guide](https://pkg.go.dev/golang.org/x/mobile)
- [Android Development](https://developer.android.com/)
- [Dart Language](https://dart.dev/guides)
- [Go Language](https://golang.org/doc/)

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review Flutter documentation
3. Check GoMobile documentation
4. Search GitHub issues
5. Ask in Flutter community forums
