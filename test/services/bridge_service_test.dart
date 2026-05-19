import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/services/bridge_service.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  group('BridgeService', () {
    late BridgeService bridgeService;

    setUp(() {
      bridgeService = BridgeService();
    });

    group('JSON Serialization', () {
      test('BridgeResponse.fromJson deserializes success response', () {
        final json = {
          'status': 'success',
          'data': {'id': '1', 'title': 'Test Anime'},
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);

        expect(response.status, 'success');
        expect(response.isSuccess, true);
        expect(response.isError, false);
        expect(response.data, {'id': '1', 'title': 'Test Anime'});
      });

      test('BridgeResponse.fromJson deserializes error response', () {
        final json = {
          'status': 'error',
          'error': {
            'code': 'TEST_ERROR',
            'message': 'Test error message',
            'details': {'key': 'value'},
          },
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);

        expect(response.status, 'error');
        expect(response.isError, true);
        expect(response.isSuccess, false);
        expect(response.error?.code, 'TEST_ERROR');
        expect(response.error?.message, 'Test error message');
      });

      test('BridgeError.fromJson deserializes error details', () {
        final json = {
          'code': 'BRIDGE_TIMEOUT',
          'message': 'Operation timed out',
          'details': {'timeout_ms': 500},
        };

        final error = BridgeError.fromJson(json);

        expect(error.code, 'BRIDGE_TIMEOUT');
        expect(error.message, 'Operation timed out');
        expect(error.details?['timeout_ms'], 500);
      });

      test('BridgeError.toString formats error correctly', () {
        final error = BridgeError(
          code: 'TEST_CODE',
          message: 'Test message',
        );

        expect(error.toString(), 'BridgeError(TEST_CODE): Test message');
      });

      test('BridgeException.toString formats exception correctly', () {
        final exception = BridgeException(
          code: 'TEST_CODE',
          message: 'Test message',
        );

        expect(exception.toString(), 'BridgeException(TEST_CODE): Test message');
      });
    });

    group('Error Handling', () {
      test('BridgeResponse identifies success status correctly', () {
        final successResponse = BridgeResponse(
          status: 'success',
          data: {},
          timestamp: DateTime.now(),
        );

        expect(successResponse.isSuccess, true);
        expect(successResponse.isError, false);
      });

      test('BridgeResponse identifies error status correctly', () {
        final errorResponse = BridgeResponse(
          status: 'error',
          error: BridgeError(
            code: 'TEST_ERROR',
            message: 'Test',
          ),
          timestamp: DateTime.now(),
        );

        expect(errorResponse.isError, true);
        expect(errorResponse.isSuccess, false);
      });

      test('BridgeException can be thrown and caught', () {
        expect(
          () => throw BridgeException(
            code: 'TEST',
            message: 'Test error',
          ),
          throwsA(isA<BridgeException>()),
        );
      });
    });

    group('Data Type Support', () {
      test('BridgeResponse supports string data', () {
        final json = {
          'status': 'success',
          'data': 'test string',
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.data, 'test string');
      });

      test('BridgeResponse supports int data', () {
        final json = {
          'status': 'success',
          'data': 42,
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.data, 42);
      });

      test('BridgeResponse supports float data', () {
        final json = {
          'status': 'success',
          'data': 3.14,
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.data, 3.14);
      });

      test('BridgeResponse supports bool data', () {
        final json = {
          'status': 'success',
          'data': true,
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.data, true);
      });

      test('BridgeResponse supports array data', () {
        final json = {
          'status': 'success',
          'data': [1, 2, 3, 'test'],
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.data, [1, 2, 3, 'test']);
      });

      test('BridgeResponse supports object data', () {
        final json = {
          'status': 'success',
          'data': {
            'nested': {
              'key': 'value',
              'number': 123,
            },
          },
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.data['nested']['key'], 'value');
        expect(response.data['nested']['number'], 123);
      });
    });

    group('Timestamp Handling', () {
      test('BridgeResponse parses ISO 8601 timestamp', () {
        final json = {
          'status': 'success',
          'data': {},
          'timestamp': '2024-01-15T10:30:45Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.timestamp.year, 2024);
        expect(response.timestamp.month, 1);
        expect(response.timestamp.day, 15);
      });

      test('BridgeResponse handles null error', () {
        final json = {
          'status': 'success',
          'data': {},
          'error': null,
          'timestamp': '2024-01-01T00:00:00Z',
        };

        final response = BridgeResponse.fromJson(json);
        expect(response.error, null);
      });
    });
  });
}
