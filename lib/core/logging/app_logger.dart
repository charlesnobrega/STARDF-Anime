import 'dart:async';
import 'dart:io';
import 'package:stardf_anime_mobile/core/logging/log_entry.dart';
import 'package:stardf_anime_mobile/core/logging/log_file_manager.dart';
import 'package:stardf_anime_mobile/core/logging/log_level.dart';

/// Main logger class for the application
class AppLogger {
  static final AppLogger _instance = AppLogger._internal();
  
  late LogFileManager _fileManager;
  LogLevel _minimumLevel = LogLevel.debug;
  bool _initialized = false;
  final List<LogEntry> _memoryBuffer = [];
  static const int _maxMemoryBufferSize = 1000;

  AppLogger._internal();

  /// Get the singleton instance
  factory AppLogger() {
    return _instance;
  }

  /// Initialize the logger
  Future<void> initialize({LogLevel minimumLevel = LogLevel.debug}) async {
    if (_initialized) return;

    _minimumLevel = minimumLevel;
    _fileManager = LogFileManager();
    await _fileManager.initialize();
    _initialized = true;
  }

  /// Check if logger is initialized
  bool get isInitialized => _initialized;

  /// Set the minimum log level
  void setMinimumLevel(LogLevel level) {
    _minimumLevel = level;
  }

  /// Log a debug message
  void debug(String message, {String? tag, StackTrace? stackTrace}) {
    _log(LogLevel.debug, message, tag: tag, stackTrace: stackTrace);
  }

  /// Log an info message
  void info(String message, {String? tag, StackTrace? stackTrace}) {
    _log(LogLevel.info, message, tag: tag, stackTrace: stackTrace);
  }

  /// Log a warning message
  void warning(String message, {String? tag, StackTrace? stackTrace}) {
    _log(LogLevel.warning, message, tag: tag, stackTrace: stackTrace);
  }

  /// Log an error message
  void error(String message, {String? tag, StackTrace? stackTrace}) {
    _log(LogLevel.error, message, tag: tag, stackTrace: stackTrace);
  }

  /// Internal logging method
  void _log(
    LogLevel level,
    String message, {
    String? tag,
    StackTrace? stackTrace,
  }) {
    if (!level.shouldLog(_minimumLevel)) {
      return;
    }

    final entry = LogEntry(
      timestamp: DateTime.now(),
      level: level,
      message: message,
      tag: tag,
      stackTrace: stackTrace?.toString(),
      threadId: _getThreadId(),
    );

    // Add to memory buffer
    _memoryBuffer.add(entry);
    if (_memoryBuffer.length > _maxMemoryBufferSize) {
      _memoryBuffer.removeAt(0);
    }

    // Print to console
    print(entry.format());

    // Write to file asynchronously
    if (_initialized) {
      _fileManager.writeLog(entry).catchError((e) {
        print('Error writing log to file: $e');
      });
    }
  }

  /// Get thread ID (simplified - returns a pseudo thread ID)
  int _getThreadId() {
    // In Dart, we don't have direct access to thread IDs
    // This returns a simplified identifier based on the isolate
    return hashCode % 10000;
  }

  /// Get all logs from memory buffer
  List<LogEntry> getMemoryLogs({LogLevel? minimumLevel}) {
    if (minimumLevel == null) {
      return List.from(_memoryBuffer);
    }
    return _memoryBuffer.where((e) => e.level.shouldLog(minimumLevel)).toList();
  }

  /// Get logs from file
  Future<List<String>> getFileLogs({int? maxLines}) async {
    try {
      final files = await _fileManager.getLogFiles();
      final logs = <String>[];

      for (final file in files) {
        final lines = await file.readAsLines();
        logs.addAll(lines);
      }

      if (maxLines != null && logs.length > maxLines) {
        return logs.sublist(logs.length - maxLines);
      }

      return logs;
    } catch (e) {
      print('Error reading log files: $e');
      return [];
    }
  }

  /// Export logs to a file
  Future<File?> exportLogs({String? fileName}) async {
    try {
      final timestamp = DateTime.now().millisecondsSinceEpoch;
      final name = fileName ?? 'logs-export-$timestamp.txt';
      
      final directory = Directory(_fileManager.logDirectoryPath);
      final exportFile = File('${directory.path}/$name');

      final logs = await getFileLogs();
      await exportFile.writeAsString(logs.join('\n'));

      return exportFile;
    } catch (e) {
      print('Error exporting logs: $e');
      return null;
    }
  }

  /// Clear all logs
  Future<void> clearLogs() async {
    _memoryBuffer.clear();
    await _fileManager.clearAllLogs();
  }

  /// Get log statistics
  Map<String, dynamic> getStatistics() {
    final memoryLogs = _memoryBuffer;
    final debugCount = memoryLogs.where((e) => e.level == LogLevel.debug).length;
    final infoCount = memoryLogs.where((e) => e.level == LogLevel.info).length;
    final warningCount = memoryLogs.where((e) => e.level == LogLevel.warning).length;
    final errorCount = memoryLogs.where((e) => e.level == LogLevel.error).length;

    return {
      'totalLogs': memoryLogs.length,
      'debugCount': debugCount,
      'infoCount': infoCount,
      'warningCount': warningCount,
      'errorCount': errorCount,
      'logDirectory': _fileManager.logDirectoryPath,
      'currentLogFile': _fileManager.currentLogFilePath,
    };
  }
}
