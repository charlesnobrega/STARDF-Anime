import 'dart:io';
import 'package:path_provider/path_provider.dart';
import 'package:stardf_anime_mobile/core/logging/log_entry.dart';

/// Manages log file storage and rotation
class LogFileManager {
  static const String _logDirName = 'logs';
  static const String _logFilePrefix = 'app';
  static const String _logFileExtension = '.log';
  static const int _maxFileSizeBytes = 10 * 1024 * 1024; // 10MB
  static const int _maxLogFiles = 7; // Keep 7 days of logs

  late Directory _logDirectory;
  late File _currentLogFile;
  int _currentFileSizeBytes = 0;

  /// Initialize the log file manager
  Future<void> initialize() async {
    try {
      final appDocDir = await getApplicationDocumentsDirectory();
      _logDirectory = Directory('${appDocDir.path}/$_logDirName');
      
      // Create logs directory if it doesn't exist
      if (!await _logDirectory.exists()) {
        await _logDirectory.create(recursive: true);
      }

      // Get or create current log file
      _currentLogFile = await _getOrCreateLogFile();
      _currentFileSizeBytes = await _currentLogFile.length();
    } catch (e) {
      print('Error initializing log file manager: $e');
      rethrow;
    }
  }

  /// Get or create the current log file
  Future<File> _getOrCreateLogFile() async {
    final now = DateTime.now();
    final dateStr = _formatDate(now);
    final fileName = '$_logFilePrefix-$dateStr$_logFileExtension';
    final file = File('${_logDirectory.path}/$fileName');

    if (!await file.exists()) {
      await file.create(recursive: true);
    }

    return file;
  }

  /// Format date as YYYY-MM-DD
  String _formatDate(DateTime dt) {
    final year = dt.year.toString().padLeft(4, '0');
    final month = dt.month.toString().padLeft(2, '0');
    final day = dt.day.toString().padLeft(2, '0');
    return '$year-$month-$day';
  }

  /// Write a log entry to file
  Future<void> writeLog(LogEntry entry) async {
    try {
      // Check if we need to rotate the file
      final entrySize = entry.format().length + 1; // +1 for newline
      if (_currentFileSizeBytes + entrySize > _maxFileSizeBytes) {
        await _rotateLogFile();
      }

      // Write the log entry
      await _currentLogFile.writeAsString(
        '${entry.format()}\n',
        mode: FileMode.append,
      );

      _currentFileSizeBytes += entrySize;
    } catch (e) {
      print('Error writing log: $e');
    }
  }

  /// Rotate log files when size limit is reached
  Future<void> _rotateLogFile() async {
    try {
      // Get current timestamp for rotation
      final now = DateTime.now();
      final timestamp = now.millisecondsSinceEpoch;
      
      // Rename current file with timestamp
      final rotatedFileName = '$_logFilePrefix-${_formatDate(now)}-$timestamp$_logFileExtension';
      final rotatedFile = File('${_logDirectory.path}/$rotatedFileName');
      
      await _currentLogFile.rename(rotatedFile.path);

      // Create new log file
      _currentLogFile = await _getOrCreateLogFile();
      _currentFileSizeBytes = 0;

      // Clean up old log files
      await _cleanupOldLogs();
    } catch (e) {
      print('Error rotating log file: $e');
    }
  }

  /// Clean up log files older than retention period
  Future<void> _cleanupOldLogs() async {
    try {
      final files = _logDirectory.listSync();
      final logFiles = files
          .whereType<File>()
          .where((f) => f.path.contains(_logFilePrefix) && f.path.endsWith(_logFileExtension))
          .toList();

      // Sort by modification time (oldest first)
      logFiles.sort((a, b) => a.statSync().modified.compareTo(b.statSync().modified));

      // Delete files older than retention period
      if (logFiles.length > _maxLogFiles) {
        for (int i = 0; i < logFiles.length - _maxLogFiles; i++) {
          try {
            await logFiles[i].delete();
          } catch (e) {
            print('Error deleting old log file: $e');
          }
        }
      }
    } catch (e) {
      print('Error cleaning up old logs: $e');
    }
  }

  /// Get all log files
  Future<List<File>> getLogFiles() async {
    try {
      final files = _logDirectory.listSync();
      return files
          .whereType<File>()
          .where((f) => f.path.contains(_logFilePrefix) && f.path.endsWith(_logFileExtension))
          .toList();
    } catch (e) {
      print('Error getting log files: $e');
      return [];
    }
  }

  /// Get the current log file path
  String get currentLogFilePath => _currentLogFile.path;

  /// Get the logs directory path
  String get logDirectoryPath => _logDirectory.path;

  /// Clear all log files
  Future<void> clearAllLogs() async {
    try {
      final files = await getLogFiles();
      for (final file in files) {
        await file.delete();
      }
      _currentFileSizeBytes = 0;
      _currentLogFile = await _getOrCreateLogFile();
    } catch (e) {
      print('Error clearing logs: $e');
    }
  }
}
