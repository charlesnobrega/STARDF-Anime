import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Represents the overall application state
class AppState {
  final bool isInitialized;
  final bool isOnline;
  final bool isSyncing;
  final String? currentError;
  final DateTime? lastSyncTime;

  const AppState({
    this.isInitialized = false,
    this.isOnline = true,
    this.isSyncing = false,
    this.currentError,
    this.lastSyncTime,
  });

  /// Create a copy of AppState with modified fields
  AppState copyWith({
    bool? isInitialized,
    bool? isOnline,
    bool? isSyncing,
    String? currentError,
    DateTime? lastSyncTime,
    bool clearError = false,
    bool clearLastSyncTime = false,
  }) {
    return AppState(
      isInitialized: isInitialized ?? this.isInitialized,
      isOnline: isOnline ?? this.isOnline,
      isSyncing: isSyncing ?? this.isSyncing,
      currentError: clearError ? null : (currentError ?? this.currentError),
      lastSyncTime: clearLastSyncTime ? null : (lastSyncTime ?? this.lastSyncTime),
    );
  }

  @override
  String toString() {
    return 'AppState(isInitialized: $isInitialized, isOnline: $isOnline, '
        'isSyncing: $isSyncing, currentError: $currentError, lastSyncTime: $lastSyncTime)';
  }
}

/// StateNotifier for managing AppState
class AppStateNotifier extends StateNotifier<AppState> {
  AppStateNotifier() : super(const AppState());

  /// Initialize the application
  void initialize() {
    state = state.copyWith(isInitialized: true);
  }

  /// Update online status
  void setOnlineStatus(bool isOnline) {
    state = state.copyWith(isOnline: isOnline);
  }

  /// Set syncing status
  void setSyncingStatus(bool isSyncing) {
    state = state.copyWith(isSyncing: isSyncing);
  }

  /// Set current error
  void setError(String? error) {
    state = state.copyWith(currentError: error);
  }

  /// Clear current error
  void clearError() {
    state = state.copyWith(clearError: true);
  }

  /// Update last sync time
  void updateLastSyncTime() {
    state = state.copyWith(lastSyncTime: DateTime.now());
  }

  /// Reset to initial state
  void reset() {
    state = const AppState();
  }
}

/// Provider for AppState
final appStateProvider = StateNotifierProvider<AppStateNotifier, AppState>((ref) {
  return AppStateNotifier();
});
