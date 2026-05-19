/// Log levels for the logging system
enum LogLevel {
  debug,
  info,
  warning,
  error;

  /// Get the string representation of the log level
  String get name {
    switch (this) {
      case LogLevel.debug:
        return 'DEBUG';
      case LogLevel.info:
        return 'INFO';
      case LogLevel.warning:
        return 'WARNING';
      case LogLevel.error:
        return 'ERROR';
    }
  }

  /// Get the numeric value for comparison
  int get value {
    switch (this) {
      case LogLevel.debug:
        return 0;
      case LogLevel.info:
        return 1;
      case LogLevel.warning:
        return 2;
      case LogLevel.error:
        return 3;
    }
  }

  /// Check if this level should be logged based on minimum level
  bool shouldLog(LogLevel minimumLevel) => value >= minimumLevel.value;
}
