import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/core/logging/log_entry.dart';
import 'package:stardf_anime_mobile/core/logging/log_level.dart';

void main() {
  group('LogEntry', () {
    test('LogEntry formats correctly without tag or stacktrace', () {
      final now = DateTime(2024, 1, 15, 12, 30, 45, 123);
      final entry = LogEntry(
        timestamp: now,
        level: LogLevel.info,
        message: 'Test message',
        tag: null,
        stackTrace: null,
        threadId: 1,
      );

      final formatted = entry.format();
      expect(formatted, contains('2024-01-15 12:30:45.123'));
      expect(formatted, contains('[INFO'));
      expect(formatted, contains('[Thread-1]'));
      expect(formatted, contains('Test message'));
    });

    test('LogEntry formats correctly with tag', () {
      final now = DateTime(2024, 1, 15, 12, 30, 45, 123);
      final entry = LogEntry(
        timestamp: now,
        level: LogLevel.debug,
        message: 'Debug message',
        tag: 'BridgeService',
        stackTrace: null,
        threadId: 2,
      );

      final formatted = entry.format();
      expect(formatted, contains('[DEBUG'));
      expect(formatted, contains('[BridgeService]'));
      expect(formatted, contains('Debug message'));
    });

    test('LogEntry formats correctly with stacktrace', () {
      final now = DateTime(2024, 1, 15, 12, 30, 45, 123);
      final stackTrace = 'at main() in main.dart:10';
      final entry = LogEntry(
        timestamp: now,
        level: LogLevel.error,
        message: 'Error message',
        tag: 'ErrorHandler',
        stackTrace: stackTrace,
        threadId: 3,
      );

      final formatted = entry.format();
      expect(formatted, contains('[ERROR'));
      expect(formatted, contains('[ErrorHandler]'));
      expect(formatted, contains('Error message'));
      expect(formatted, contains(stackTrace));
    });

    test('LogEntry formats all log levels correctly', () {
      final now = DateTime(2024, 1, 15, 12, 30, 45, 123);

      final debugEntry = LogEntry(
        timestamp: now,
        level: LogLevel.debug,
        message: 'Debug',
        threadId: 1,
      );
      expect(debugEntry.format(), contains('[DEBUG'));

      final infoEntry = LogEntry(
        timestamp: now,
        level: LogLevel.info,
        message: 'Info',
        threadId: 1,
      );
      expect(infoEntry.format(), contains('[INFO'));

      final warningEntry = LogEntry(
        timestamp: now,
        level: LogLevel.warning,
        message: 'Warning',
        threadId: 1,
      );
      expect(warningEntry.format(), contains('[WARNING'));

      final errorEntry = LogEntry(
        timestamp: now,
        level: LogLevel.error,
        message: 'Error',
        threadId: 1,
      );
      expect(errorEntry.format(), contains('[ERROR'));
    });

    test('LogEntry toString returns formatted string', () {
      final now = DateTime(2024, 1, 15, 12, 30, 45, 123);
      final entry = LogEntry(
        timestamp: now,
        level: LogLevel.info,
        message: 'Test',
        threadId: 1,
      );

      expect(entry.toString(), equals(entry.format()));
    });

    test('LogEntry formats timestamp with correct padding', () {
      final now = DateTime(2024, 1, 5, 9, 5, 3, 5);
      final entry = LogEntry(
        timestamp: now,
        level: LogLevel.info,
        message: 'Test',
        threadId: 1,
      );

      final formatted = entry.format();
      expect(formatted, contains('2024-01-05 09:05:03.005'));
    });
  });
}
