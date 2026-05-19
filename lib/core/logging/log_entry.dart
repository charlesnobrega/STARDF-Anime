import 'package:stardf_anime_mobile/core/logging/log_level.dart';

/// Represents a single log entry
class LogEntry {
  final DateTime timestamp;
  final LogLevel level;
  final String message;
  final String? tag;
  final String? stackTrace;
  final int threadId;

  LogEntry({
    required this.timestamp,
    required this.level,
    required this.message,
    this.tag,
    this.stackTrace,
    required this.threadId,
  });

  /// Format the log entry as a string
  String format() {
    final timeStr = _formatTime(timestamp);
    final levelStr = level.name.padRight(7);
    final tagStr = tag != null ? '[$tag]' : '';
    final threadStr = '[Thread-$threadId]';
    
    final formatted = '$timeStr [$levelStr] $threadStr $tagStr $message';
    
    if (stackTrace != null && stackTrace!.isNotEmpty) {
      return '$formatted\n$stackTrace';
    }
    
    return formatted;
  }

  /// Format timestamp as ISO 8601 with milliseconds
  String _formatTime(DateTime dt) {
    final year = dt.year.toString().padLeft(4, '0');
    final month = dt.month.toString().padLeft(2, '0');
    final day = dt.day.toString().padLeft(2, '0');
    final hour = dt.hour.toString().padLeft(2, '0');
    final minute = dt.minute.toString().padLeft(2, '0');
    final second = dt.second.toString().padLeft(2, '0');
    final ms = dt.millisecond.toString().padLeft(3, '0');
    
    return '$year-$month-$day $hour:$minute:$second.$ms';
  }

  @override
  String toString() => format();
}
