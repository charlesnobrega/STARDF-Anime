import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/core/logging/log_level.dart';

void main() {
  group('LogLevel', () {
    test('LogLevel.debug has correct name', () {
      expect(LogLevel.debug.name, equals('DEBUG'));
    });

    test('LogLevel.info has correct name', () {
      expect(LogLevel.info.name, equals('INFO'));
    });

    test('LogLevel.warning has correct name', () {
      expect(LogLevel.warning.name, equals('WARNING'));
    });

    test('LogLevel.error has correct name', () {
      expect(LogLevel.error.name, equals('ERROR'));
    });

    test('LogLevel values are in correct order', () {
      expect(LogLevel.debug.value, equals(0));
      expect(LogLevel.info.value, equals(1));
      expect(LogLevel.warning.value, equals(2));
      expect(LogLevel.error.value, equals(3));
    });

    test('shouldLog returns true for equal or higher levels', () {
      expect(LogLevel.debug.shouldLog(LogLevel.debug), isTrue);
      expect(LogLevel.info.shouldLog(LogLevel.debug), isTrue);
      expect(LogLevel.warning.shouldLog(LogLevel.debug), isTrue);
      expect(LogLevel.error.shouldLog(LogLevel.debug), isTrue);
    });

    test('shouldLog returns false for lower levels', () {
      expect(LogLevel.debug.shouldLog(LogLevel.info), isFalse);
      expect(LogLevel.debug.shouldLog(LogLevel.warning), isFalse);
      expect(LogLevel.debug.shouldLog(LogLevel.error), isFalse);
      expect(LogLevel.info.shouldLog(LogLevel.warning), isFalse);
    });

    test('shouldLog with same level returns true', () {
      expect(LogLevel.warning.shouldLog(LogLevel.warning), isTrue);
      expect(LogLevel.error.shouldLog(LogLevel.error), isTrue);
    });
  });
}
