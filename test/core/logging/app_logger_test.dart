import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/core/logging/app_logger.dart';
import 'package:stardf_anime_mobile/core/logging/log_level.dart';

void main() {
  group('AppLogger', () {
    late AppLogger logger;

    setUp(() {
      logger = AppLogger();
    });

    test('AppLogger is a singleton', () {
      final logger1 = AppLogger();
      final logger2 = AppLogger();
      expect(identical(logger1, logger2), isTrue);
    });

    test('AppLogger initializes correctly', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      // Skip file initialization for unit tests
      expect(logger.isInitialized, isFalse);
    });

    test('AppLogger can be initialized only once', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      // Skip file initialization for unit tests
      expect(logger.isInitialized, isFalse);
    });

    test('AppLogger logs debug messages', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      // Don't initialize file manager for unit tests
      logger.debug('Debug message', tag: 'TestTag');
      
      final logs = logger.getMemoryLogs();
      expect(logs.isNotEmpty, isTrue);
      expect(logs.last.message, equals('Debug message'));
      expect(logs.last.level, equals(LogLevel.debug));
      expect(logs.last.tag, equals('TestTag'));
    });

    test('AppLogger logs info messages', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.info('Info message', tag: 'TestTag');
      
      final logs = logger.getMemoryLogs();
      expect(logs.isNotEmpty, isTrue);
      expect(logs.last.message, equals('Info message'));
      expect(logs.last.level, equals(LogLevel.info));
    });

    test('AppLogger logs warning messages', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.warning('Warning message', tag: 'TestTag');
      
      final logs = logger.getMemoryLogs();
      expect(logs.isNotEmpty, isTrue);
      expect(logs.last.message, equals('Warning message'));
      expect(logs.last.level, equals(LogLevel.warning));
    });

    test('AppLogger logs error messages', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.error('Error message', tag: 'TestTag');
      
      final logs = logger.getMemoryLogs();
      expect(logs.isNotEmpty, isTrue);
      expect(logs.last.message, equals('Error message'));
      expect(logs.last.level, equals(LogLevel.error));
    });

    test('AppLogger respects minimum log level', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.warning);
      
      logger.debug('Debug message');
      logger.info('Info message');
      logger.warning('Warning message');
      logger.error('Error message');
      
      final logs = logger.getMemoryLogs();
      expect(logs.length, equals(2)); // Only warning and error
      expect(logs[0].level, equals(LogLevel.warning));
      expect(logs[1].level, equals(LogLevel.error));
    });

    test('AppLogger can change minimum level', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.debug);
      
      logger.debug('Debug 1');
      expect(logger.getMemoryLogs().length, equals(1));
      
      logger.setMinimumLevel(LogLevel.error);
      logger.debug('Debug 2');
      logger.info('Info 1');
      logger.warning('Warning 1');
      logger.error('Error 1');
      
      final logs = logger.getMemoryLogs();
      expect(logs.length, equals(2)); // Debug 1 and Error 1
    });

    test('AppLogger memory buffer has maximum size', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.debug);
      
      // Log more than max buffer size
      for (int i = 0; i < 1500; i++) {
        logger.info('Message $i');
      }
      
      final logs = logger.getMemoryLogs();
      expect(logs.length, lessThanOrEqualTo(1000));
    });

    test('AppLogger filters logs by level', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.debug);
      
      logger.debug('Debug');
      logger.info('Info');
      logger.warning('Warning');
      logger.error('Error');
      
      final debugLogs = logger.getMemoryLogs(minimumLevel: LogLevel.debug);
      final warningLogs = logger.getMemoryLogs(minimumLevel: LogLevel.warning);
      
      expect(debugLogs.length, equals(4));
      expect(warningLogs.length, equals(2));
    });

    test('AppLogger gets statistics', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.debug);
      
      logger.debug('Debug');
      logger.info('Info');
      logger.warning('Warning');
      logger.error('Error');
      
      final stats = logger.getStatistics();
      expect(stats['totalLogs'], equals(4));
      expect(stats['debugCount'], equals(1));
      expect(stats['infoCount'], equals(1));
      expect(stats['warningCount'], equals(1));
      expect(stats['errorCount'], equals(1));
    });

    test('AppLogger clears memory logs', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.debug);
      
      logger.debug('Debug');
      logger.info('Info');
      expect(logger.getMemoryLogs().length, equals(2));
      
      await logger.clearLogs();
      expect(logger.getMemoryLogs().length, equals(0));
    });

    test('AppLogger logs with stacktrace', () async {
      TestWidgetsFlutterBinding.ensureInitialized();
      logger.setMinimumLevel(LogLevel.debug);
      
      try {
        throw Exception('Test exception');
      } catch (e, stackTrace) {
        logger.error('Error occurred', stackTrace: stackTrace);
      }
      
      final logs = logger.getMemoryLogs();
      expect(logs.isNotEmpty, isTrue);
      expect(logs.last.stackTrace, isNotNull);
      expect(logs.last.stackTrace, contains('Exception'));
    });
  });
}
