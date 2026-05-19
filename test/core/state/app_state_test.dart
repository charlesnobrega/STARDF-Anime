import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/core/state/app_state.dart';

void main() {
  group('AppState', () {
    test('creates instance with default values', () {
      const state = AppState();

      expect(state.isInitialized, false);
      expect(state.isOnline, true);
      expect(state.isSyncing, false);
      expect(state.currentError, null);
      expect(state.lastSyncTime, null);
    });

    test('copyWith updates specified fields', () {
      const state = AppState(
        isInitialized: true,
        isOnline: false,
        isSyncing: true,
      );

      final updated = state.copyWith(
        isOnline: true,
        currentError: 'Test error',
      );

      expect(updated.isInitialized, true);
      expect(updated.isOnline, true);
      expect(updated.isSyncing, true);
      expect(updated.currentError, 'Test error');
    });

    test('copyWith preserves unspecified fields', () {
      final now = DateTime.now();
      final state = AppState(
        isInitialized: true,
        lastSyncTime: now,
      );

      final updated = state.copyWith(isOnline: false);

      expect(updated.isInitialized, true);
      expect(updated.lastSyncTime, now);
      expect(updated.isOnline, false);
    });

    test('toString returns formatted string', () {
      const state = AppState(
        isInitialized: true,
        isOnline: false,
      );

      final str = state.toString();

      expect(str, contains('isInitialized: true'));
      expect(str, contains('isOnline: false'));
    });
  });

  group('AppStateNotifier', () {
    test('initializes with default AppState', () {
      final notifier = AppStateNotifier();

      expect(notifier.state.isInitialized, false);
      expect(notifier.state.isOnline, true);
    });

    test('initialize sets isInitialized to true', () {
      final notifier = AppStateNotifier();

      notifier.initialize();

      expect(notifier.state.isInitialized, true);
    });

    test('setOnlineStatus updates isOnline', () {
      final notifier = AppStateNotifier();

      notifier.setOnlineStatus(false);

      expect(notifier.state.isOnline, false);
    });

    test('setSyncingStatus updates isSyncing', () {
      final notifier = AppStateNotifier();

      notifier.setSyncingStatus(true);

      expect(notifier.state.isSyncing, true);
    });

    test('setError updates currentError', () {
      final notifier = AppStateNotifier();

      notifier.setError('Test error');

      expect(notifier.state.currentError, 'Test error');
    });

    test('clearError sets currentError to null', () {
      final notifier = AppStateNotifier();
      notifier.setError('Test error');

      notifier.clearError();

      expect(notifier.state.currentError, null);
    });

    test('updateLastSyncTime sets lastSyncTime to now', () {
      final notifier = AppStateNotifier();
      final before = DateTime.now().subtract(const Duration(seconds: 1));

      notifier.updateLastSyncTime();

      final after = DateTime.now().add(const Duration(seconds: 1));
      expect(notifier.state.lastSyncTime, isNotNull);
      expect(notifier.state.lastSyncTime!.isAfter(before), true);
      expect(notifier.state.lastSyncTime!.isBefore(after), true);
    });

    test('reset returns to initial state', () {
      final notifier = AppStateNotifier();
      notifier.initialize();
      notifier.setOnlineStatus(false);
      notifier.setSyncingStatus(true);
      notifier.setError('Test error');

      notifier.reset();

      expect(notifier.state.isInitialized, false);
      expect(notifier.state.isOnline, true);
      expect(notifier.state.isSyncing, false);
      expect(notifier.state.currentError, null);
    });

    test('multiple operations maintain state consistency', () {
      final notifier = AppStateNotifier();

      notifier.initialize();
      notifier.setOnlineStatus(false);
      notifier.setSyncingStatus(true);
      notifier.setError('Error 1');
      notifier.updateLastSyncTime();

      expect(notifier.state.isInitialized, true);
      expect(notifier.state.isOnline, false);
      expect(notifier.state.isSyncing, true);
      expect(notifier.state.currentError, 'Error 1');
      expect(notifier.state.lastSyncTime, isNotNull);

      notifier.clearError();
      notifier.setSyncingStatus(false);

      expect(notifier.state.currentError, null);
      expect(notifier.state.isSyncing, false);
      expect(notifier.state.isInitialized, true);
    });
  });
}
