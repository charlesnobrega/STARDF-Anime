import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:stardf_anime_mobile/core/logging/app_logger.dart';

/// Represents connectivity status
enum ConnectivityStatus {
  online,
  offline,
  unknown,
}

/// Provider for connectivity status
final connectivityProvider = StreamProvider<ConnectivityStatus>((ref) async* {
  final connectivity = Connectivity();
  final logger = AppLogger();

  // Get initial status
  final result = await connectivity.checkConnectivity();
  yield _mapConnectivityResult(result);

  // Listen to connectivity changes
  await for (final result in connectivity.onConnectivityChanged) {
    final status = _mapConnectivityResult(result);
    logger.info('Connectivity changed: $status', tag: 'ConnectivityProvider');
    yield status;
  }
});

/// Provider for online status (boolean)
final isOnlineProvider = Provider<bool>((ref) {
  final connectivity = ref.watch(connectivityProvider);
  return connectivity.when(
    data: (status) => status == ConnectivityStatus.online,
    loading: () => true, // Assume online while loading
    error: (_, __) => true, // Assume online on error
  );
});

/// Provider for offline status (boolean)
final isOfflineProvider = Provider<bool>((ref) {
  return !ref.watch(isOnlineProvider);
});

/// Helper function to map connectivity result
ConnectivityStatus _mapConnectivityResult(ConnectivityResult result) {
  switch (result) {
    case ConnectivityResult.mobile:
    case ConnectivityResult.wifi:
    case ConnectivityResult.ethernet:
      return ConnectivityStatus.online;
    case ConnectivityResult.none:
      return ConnectivityStatus.offline;
    case ConnectivityResult.vpn:
    case ConnectivityResult.bluetooth:
    case ConnectivityResult.other:
      return ConnectivityStatus.online;
  }
}
