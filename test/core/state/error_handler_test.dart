import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/core/logging/app_logger.dart';
import 'package:stardf_anime_mobile/core/state/error_handler_provider.dart';

void main() {
  group('AppError', () {
    test('creates instance with required fields', () {
      final error = AppError(
        code: 'TEST_ERROR',
        message: 'Test error message',
        category: ErrorCategory.network,
      );

      expect(error.code, 'TEST_ERROR');
      expect(error.message, 'Test error message');
      expect(error.category, ErrorCategory.network);
      expect(error.timestamp, isNotNull);
    });

    test('creates instance with optional fields', () {
      final stackTrace = StackTrace.current;
      final error = AppError(
        code: 'TEST_ERROR',
        message: 'Test error message',
        details: 'Additional details',
        category: ErrorCategory.bridge,
        stackTrace: stackTrace,
      );

      expect(error.details, 'Additional details');
      expect(error.stackTrace, stackTrace);
    });

    test('toString returns formatted string', () {
      final error = AppError(
        code: 'TEST_ERROR',
        message: 'Test message',
        category: ErrorCategory.sync,
      );

      final str = error.toString();

      expect(str, contains('TEST_ERROR'));
      expect(str, contains('Test message'));
      expect(str, contains('sync'));
    });
  });

  group('ErrorHandler', () {
    late ErrorHandler errorHandler;

    setUp(() {
      final logger = AppLogger();
      errorHandler = ErrorHandler(logger);
    });

    test('handleError creates and stores error', () {
      final error = errorHandler.handleError(
        code: 'TEST_ERROR',
        message: 'Test message',
        category: ErrorCategory.network,
      );

      expect(error.code, 'TEST_ERROR');
      expect(error.message, 'Test message');
      expect(errorHandler.getErrorHistory().length, 1);
    });

    test('handleError with details stores details', () {
      final error = errorHandler.handleError(
        code: 'TEST_ERROR',
        message: 'Test message',
        details: 'Error details',
        category: ErrorCategory.bridge,
      );

      expect(error.details, 'Error details');
    });

    test('handleException creates error from exception', () {
      final exception = NetworkException('Connection failed');

      final error = errorHandler.handleException(
        exception,
        category: ErrorCategory.network,
      );

      expect(error.code, 'NETWORK_ERROR');
      expect(error.message, contains('Connection failed'));
    });

    test('getUserMessage returns appropriate message for category', () {
      final networkError = AppError(
        code: 'NET_ERROR',
        message: 'Network failed',
        category: ErrorCategory.network,
      );

      final message = errorHandler.getUserMessage(networkError);

      expect(message, contains('conexão'));
    });

    test('getErrorHistory returns all errors', () {
      errorHandler.handleError(
        code: 'ERROR1',
        message: 'First error',
        category: ErrorCategory.network,
      );
      errorHandler.handleError(
        code: 'ERROR2',
        message: 'Second error',
        category: ErrorCategory.bridge,
      );

      final history = errorHandler.getErrorHistory();

      expect(history.length, 2);
      expect(history[0].code, 'ERROR1');
      expect(history[1].code, 'ERROR2');
    });

    test('clearErrorHistory removes all errors', () {
      errorHandler.handleError(
        code: 'ERROR1',
        message: 'Error',
        category: ErrorCategory.network,
      );

      errorHandler.clearErrorHistory();

      expect(errorHandler.getErrorHistory().length, 0);
    });

    test('getErrorsByCategory filters errors', () {
      errorHandler.handleError(
        code: 'NET_ERROR',
        message: 'Network error',
        category: ErrorCategory.network,
      );
      errorHandler.handleError(
        code: 'BRIDGE_ERROR',
        message: 'Bridge error',
        category: ErrorCategory.bridge,
      );
      errorHandler.handleError(
        code: 'NET_ERROR2',
        message: 'Another network error',
        category: ErrorCategory.network,
      );

      final networkErrors = errorHandler.getErrorsByCategory(ErrorCategory.network);

      expect(networkErrors.length, 2);
      expect(networkErrors.every((e) => e.category == ErrorCategory.network), true);
    });

    test('getRecentErrors returns limited errors', () {
      for (int i = 0; i < 15; i++) {
        errorHandler.handleError(
          code: 'ERROR$i',
          message: 'Error $i',
          category: ErrorCategory.network,
        );
      }

      final recent = errorHandler.getRecentErrors(limit: 5);

      expect(recent.length, 5);
      expect(recent.last.code, 'ERROR14');
    });

    test('error history respects max size limit', () {
      // Create more than max errors
      for (int i = 0; i < 150; i++) {
        errorHandler.handleError(
          code: 'ERROR$i',
          message: 'Error $i',
          category: ErrorCategory.network,
        );
      }

      final history = errorHandler.getErrorHistory();

      expect(history.length, 100); // Max history size
    });
  });

  group('Custom Exceptions', () {
    test('NetworkException creates with message', () {
      final exception = NetworkException('Connection timeout');

      expect(exception.toString(), 'Connection timeout');
    });

    test('BridgeException creates with message and code', () {
      final exception = BridgeException('Bridge call failed', code: 'BRIDGE_TIMEOUT');

      expect(exception.toString(), 'Bridge call failed');
      expect(exception.code, 'BRIDGE_TIMEOUT');
    });

    test('TimeoutException creates with message', () {
      final exception = TimeoutException('Operation timed out');

      expect(exception.toString(), 'Operation timed out');
    });

    test('SyncException creates with message', () {
      final exception = SyncException('Sync failed');

      expect(exception.toString(), 'Sync failed');
    });

    test('ResourceException creates with message', () {
      final exception = ResourceException('Resource not found');

      expect(exception.toString(), 'Resource not found');
    });
  });

  group('ErrorCategory', () {
    test('all error categories are defined', () {
      expect(ErrorCategory.network, isNotNull);
      expect(ErrorCategory.bridge, isNotNull);
      expect(ErrorCategory.sync, isNotNull);
      expect(ErrorCategory.resource, isNotNull);
      expect(ErrorCategory.unknown, isNotNull);
    });
  });
}
