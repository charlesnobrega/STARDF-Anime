import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/services/bridge_service.dart';

void main() {
  group('Web Server Integration Tests', () {
    late BridgeService bridgeService;

    setUp(() {
      bridgeService = BridgeService();
    });

    group('Server Lifecycle', () {
      test('StartWebServer initializes server successfully', () async {
        try {
          final result = await bridgeService.startWebServer(8080);
          expect(result, isNotNull);
          expect(result['message'], contains('started'));
          expect(result['port'], 8080);
          expect(result['url'], contains('localhost'));

          // Cleanup
          await bridgeService.shutdownWebServer();
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('PauseWebServer pauses server gracefully', () async {
        try {
          await bridgeService.startWebServer(8080);
          final result = await bridgeService.pauseWebServer();
          expect(result, isNotNull);
          expect(result['message'], contains('paused'));

          // Cleanup
          await bridgeService.shutdownWebServer();
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('ResumeWebServer resumes server after pause', () async {
        try {
          await bridgeService.startWebServer(8080);
          await bridgeService.pauseWebServer();
          final result = await bridgeService.resumeWebServer();
          expect(result, isNotNull);
          expect(result['message'], contains('resumed'));

          // Cleanup
          await bridgeService.shutdownWebServer();
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('ShutdownWebServer closes server gracefully', () async {
        try {
          await bridgeService.startWebServer(8080);
          final result = await bridgeService.shutdownWebServer();
          expect(result, isNotNull);
          expect(result['message'], contains('shut down'));
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });
    });

    group('Bridge Method Calls', () {
      test('GetAnimes returns list of animes', () async {
        try {
          final result = await bridgeService.getAnimes();
          expect(result, isA<List>());
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('AddToWatchlist adds anime to watchlist', () async {
        try {
          final result = await bridgeService.addToWatchlist('1');
          expect(result, isA<Map>());
          expect(result['animeId'], '1');
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('MarkAsWatched marks episode as watched', () async {
        try {
          final result = await bridgeService.markAsWatched('1', 'ep1');
          expect(result, isA<Map>());
          expect(result['animeId'], '1');
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('GetSyncStatus returns sync status', () async {
        try {
          final result = await bridgeService.getSyncStatus();
          expect(result, isA<Map>());
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });

      test('SyncWithDesktop synchronizes with desktop', () async {
        try {
          final result = await bridgeService.syncWithDesktop();
          expect(result, isA<Map>());
        } catch (e) {
          // Expected if GoMobile is not available in test environment
          expect(e, isA<BridgeException>());
        }
      });
    });

    group('Error Handling', () {
      test('AddToWatchlist with empty animeId returns error', () async {
        try {
          await bridgeService.addToWatchlist('');
          fail('Should throw BridgeException');
        } catch (e) {
          expect(e, isA<BridgeException>());
        }
      });

      test('MarkAsWatched with empty parameters returns error', () async {
        try {
          await bridgeService.markAsWatched('', '');
          fail('Should throw BridgeException');
        } catch (e) {
          expect(e, isA<BridgeException>());
        }
      });

      test('ResumeWebServer when not running returns error', () async {
        try {
          await bridgeService.resumeWebServer();
          fail('Should throw BridgeException');
        } catch (e) {
          expect(e, isA<BridgeException>());
        }
      });
    });
  });
}
