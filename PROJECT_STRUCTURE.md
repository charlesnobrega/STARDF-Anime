# StarDF-Anime Mobile - Project Structure

## Overview

This is a Flutter mobile application for Android with GoMobile integration for Go backend communication. The project follows a layered architecture with clear separation between UI, business logic, bridge, and backend.

## Directory Structure

```
.
├── lib/                          # Flutter/Dart source code
│   ├── main.dart                # App entry point
│   ├── services/
│   │   └── bridge_service.dart  # GoMobile bridge communication
│   ├── providers/               # State management (Provider/Riverpod)
│   ├── models/                  # Data models
│   ├── screens/                 # UI screens
│   ├── widgets/                 # Reusable widgets
│   └── utils/                   # Utility functions
│
├── go/                          # Go backend source code
│   ├── main.go                 # Go entry point (GoMobile)
│   ├── bridge.go               # Bridge interface for Flutter
│   ├── server.go               # Local HTTP server
│   ├── go.mod                  # Go module definition
│   └── go.sum                  # Go dependencies
│
├── test/                        # Unit tests
│   └── widget_test.dart        # Widget tests
│
├── integration_test/            # Integration tests
│   └── app_test.dart           # Integration test suite
│
├── android/                     # Android native code
│   ├── app/
│   │   ├── build.gradle        # App-level Gradle configuration
│   │   └── src/
│   │       └── main/
│   │           ├── kotlin/
│   │           │   └── com/stardf/anime/
│   │           │       └── MainActivity.kt
│   │           └── AndroidManifest.xml
│   ├── build.gradle            # Project-level Gradle configuration
│   ├── settings.gradle         # Gradle settings
│   ├── gradle.properties       # Gradle properties
│   └── local.properties.template
│
├── pubspec.yaml                # Flutter dependencies
├── analysis_options.yaml       # Dart linting rules
├── .gitignore                  # Git ignore rules
└── PROJECT_STRUCTURE.md        # This file
```

## Key Components

### Flutter Layer (lib/)

- **main.dart**: Application entry point, initializes bridge and UI
- **services/bridge_service.dart**: Handles communication with Go backend via GoMobile
- **providers/**: State management using Provider/Riverpod
- **models/**: Data classes for Anime, Watchlist, History, etc.
- **screens/**: UI screens (Home, Watchlist, History, Settings)
- **widgets/**: Reusable UI components with glassmorphism design

### Go Backend (go/)

- **bridge.go**: Exposes functions to Flutter via GoMobile
  - GetAnimes()
  - AddToWatchlist(animeID)
  - MarkAsWatched(animeID, episodeID)
  - GetSyncStatus()
  - SyncWithDesktop()

- **server.go**: Local HTTP server for REST API
  - StartWebServer(port)
  - PauseWebServer()
  - ResumeWebServer()
  - ShutdownWebServer()
  - Routes: /health, /api/animes, /api/watchlist, /api/history, /api/sync

### Android Configuration (android/)

- **MainActivity.kt**: Android activity with method channel setup
- **AndroidManifest.xml**: App permissions and configuration
- **build.gradle**: Gradle build configuration with GoMobile support
- **settings.gradle**: Gradle settings and plugin configuration

## Dependencies

### Flutter/Dart
- **provider**: State management
- **riverpod**: Alternative state management
- **drift**: SQLite ORM for local database
- **http**: HTTP client for REST API
- **connectivity_plus**: Network connectivity monitoring
- **logger**: Logging framework
- **go_router**: Navigation and routing

### Go
- **golang.org/x/mobile**: GoMobile framework for mobile integration

## Build Configuration

### Android Minimum SDK
- **minSdkVersion**: 26 (Android 8.0)
- **targetSdkVersion**: 34 (Android 14)

### GoMobile Configuration
- **ABI Filters**: arm64-v8a, armeabi-v7a, x86_64
- **Method Channel**: com.stardf.anime/bridge

## Communication Flow

1. **Flutter UI** → User interaction
2. **State Management** → Updates app state
3. **Bridge Service** → Calls Go backend via method channel
4. **Go Backend** → Processes request, returns JSON response
5. **Bridge Service** → Parses response, handles errors
6. **State Management** → Updates UI with new data

## Setup Instructions

### Prerequisites
- Flutter SDK (3.0+)
- Go SDK (1.21+)
- Android SDK (API 26+)
- GoMobile tools

### Initial Setup

1. **Install Flutter dependencies**:
   ```bash
   flutter pub get
   ```

2. **Generate Drift code**:
   ```bash
   flutter pub run build_runner build
   ```

3. **Setup Android local.properties**:
   ```bash
   cp android/local.properties.template android/local.properties
   # Edit android/local.properties with your SDK paths
   ```

4. **Build GoMobile bindings** (requires gomobile):
   ```bash
   cd go
   gomobile bind -target=android -o=../android/app/src/main/jniLibs/libgo.aar .
   ```

### Running the App

```bash
flutter run
```

### Running Tests

```bash
# Unit tests
flutter test

# Integration tests
flutter test integration_test/app_test.dart
```

## Architecture Principles

1. **Separation of Concerns**: UI, business logic, bridge, and backend are clearly separated
2. **Error Handling**: Comprehensive error handling with structured error responses
3. **Logging**: Detailed logging at all layers for debugging
4. **Timeout Management**: 500ms timeout for bridge calls to prevent hanging
5. **Graceful Lifecycle**: Proper handling of app pause/resume/shutdown
6. **Offline Support**: Local queue for operations when offline
7. **Sync Management**: Bidirectional sync with desktop backend

## Next Steps

1. Implement database schema with Drift
2. Create UI screens with glassmorphism design
3. Implement state management providers
4. Add comprehensive error handling
5. Implement sync engine for desktop synchronization
6. Add logging and debugging tools
7. Write unit and integration tests
8. Optimize performance for mobile
9. Test on multiple Android versions
10. Build and release APK/AAB

## References

- [Flutter Documentation](https://flutter.dev/docs)
- [GoMobile Documentation](https://pkg.go.dev/golang.org/x/mobile)
- [Drift Documentation](https://drift.simonbinder.eu/)
- [Provider Documentation](https://pub.dev/packages/provider)
