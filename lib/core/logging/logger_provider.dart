import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:stardf_anime_mobile/core/logging/app_logger.dart';
import 'package:stardf_anime_mobile/core/logging/log_level.dart';

/// Provider for the AppLogger singleton
final appLoggerProvider = Provider<AppLogger>((ref) {
  return AppLogger();
});

/// Provider for initializing the logger
final loggerInitializationProvider = FutureProvider<void>((ref) async {
  final logger = ref.watch(appLoggerProvider);
  await logger.initialize(minimumLevel: LogLevel.debug);
});

/// Provider for log level state
final logLevelProvider = StateProvider<LogLevel>((ref) {
  return LogLevel.debug;
});

/// Provider for memory logs
final memoryLogsProvider = FutureProvider<List<String>>((ref) async {
  final logger = ref.watch(appLoggerProvider);
  final logs = logger.getMemoryLogs();
  return logs.map((e) => e.format()).toList();
});

/// Provider for file logs
final fileLogsProvider = FutureProvider<List<String>>((ref) async {
  final logger = ref.watch(appLoggerProvider);
  return logger.getFileLogs(maxLines: 1000);
});

/// Provider for log statistics
final logStatisticsProvider = FutureProvider<Map<String, dynamic>>((ref) async {
  final logger = ref.watch(appLoggerProvider);
  return logger.getStatistics();
});
