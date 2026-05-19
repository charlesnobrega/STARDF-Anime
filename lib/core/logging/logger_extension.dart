import 'package:stardf_anime_mobile/core/logging/app_logger.dart';

/// Extension methods for easier logging
extension LoggerExtension on Object {
  /// Get the logger instance
  AppLogger get logger => AppLogger();

  /// Log a debug message with this object's type as tag
  void logDebug(String message, {StackTrace? stackTrace}) {
    logger.debug(message, tag: runtimeType.toString(), stackTrace: stackTrace);
  }

  /// Log an info message with this object's type as tag
  void logInfo(String message, {StackTrace? stackTrace}) {
    logger.info(message, tag: runtimeType.toString(), stackTrace: stackTrace);
  }

  /// Log a warning message with this object's type as tag
  void logWarning(String message, {StackTrace? stackTrace}) {
    logger.warning(message, tag: runtimeType.toString(), stackTrace: stackTrace);
  }

  /// Log an error message with this object's type as tag
  void logError(String message, {StackTrace? stackTrace}) {
    logger.error(message, tag: runtimeType.toString(), stackTrace: stackTrace);
  }
}
