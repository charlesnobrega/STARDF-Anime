import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:matcher/matcher.dart' as matcher;
import 'package:stardf_anime_mobile/core/database/database.dart';

void main() {
  late AppDatabase database;

  setUp(() {
    database = AppDatabase.forTesting();
  });

  tearDown(() async {
    await database.close();
  });

  group('Animes Table', () {
    test('insert and retrieve anime', () async {
      final anime = AnimesCompanion(
        title: const Value('Test Anime'),
        description: const Value('Test Description'),
        episodes: const Value(12),
        imageUrl: const Value('https://example.com/image.jpg'),
        genre: const Value('Action'),
        status: const Value('ongoing'),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.insertAnime(anime);
      final animes = await database.getAllAnimes();

      expect(animes.length, 1);
      expect(animes[0].title, 'Test Anime');
      expect(animes[0].episodes, 12);
    });

    test('get anime by id', () async {
      final anime = AnimesCompanion(
        title: const Value('Test Anime'),
        episodes: const Value(12),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.insertAnime(anime);
      final retrieved = await database.getAnimeById(1);

      expect(retrieved, matcher.isNotNull);
      expect(retrieved!.title, 'Test Anime');
    });

    test('update anime', () async {
      final anime = AnimesCompanion(
        title: const Value('Original Title'),
        episodes: const Value(12),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.insertAnime(anime);
      final retrieved = await database.getAnimeById(1);

      final updated = retrieved!.copyWith(
        title: 'Updated Title',
        updatedAt: DateTime.now(),
      );
      await database.updateAnime(updated);

      final result = await database.getAnimeById(1);
      expect(result!.title, 'Updated Title');
    });

    test('delete anime', () async {
      final anime = AnimesCompanion(
        title: const Value('Test Anime'),
        episodes: const Value(12),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.insertAnime(anime);
      await database.deleteAnime(1);

      final animes = await database.getAllAnimes();
      expect(animes.length, 0);
    });
  });

  group('Watchlist Table', () {
    test('add to watchlist', () async {
      final item = WatchlistCompanion(
        animeId: const Value(1),
        status: const Value('watching'),
        currentEpisode: const Value(5),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addToWatchlist(item);
      final watchlist = await database.getWatchlist();

      expect(watchlist.length, 1);
      expect(watchlist[0].animeId, 1);
      expect(watchlist[0].status, 'watching');
    });

    test('get watchlist item by anime id', () async {
      final item = WatchlistCompanion(
        animeId: const Value(1),
        status: const Value('watching'),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addToWatchlist(item);
      final retrieved = await database.getWatchlistItem(1);

      expect(retrieved, matcher.isNotNull);
      expect(retrieved!.animeId, 1);
    });

    test('update watchlist item', () async {
      final item = WatchlistCompanion(
        animeId: const Value(1),
        status: const Value('watching'),
        currentEpisode: const Value(5),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addToWatchlist(item);
      final retrieved = await database.getWatchlistItem(1);

      final updated = retrieved!.copyWith(
        currentEpisode: 10,
        updatedAt: DateTime.now(),
      );
      await database.updateWatchlistItem(updated);

      final result = await database.getWatchlistItem(1);
      expect(result!.currentEpisode, 10);
    });

    test('remove from watchlist', () async {
      final item = WatchlistCompanion(
        animeId: const Value(1),
        status: const Value('watching'),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addToWatchlist(item);
      await database.removeFromWatchlist(1);

      final watchlist = await database.getWatchlist();
      expect(watchlist.length, 0);
    });

    test('unique constraint on animeId', () async {
      final item1 = WatchlistCompanion(
        animeId: const Value(1),
        status: const Value('watching'),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      final item2 = WatchlistCompanion(
        animeId: const Value(1),
        status: const Value('completed'),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addToWatchlist(item1);
      
      // Second insert with same animeId should fail or replace
      expect(
        () => database.addToWatchlist(item2),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('History Table', () {
    test('add history entry', () async {
      final entry = HistoryCompanion(
        animeId: const Value(1),
        episodeNumber: const Value(5),
        episodeTitle: const Value('Episode 5'),
        watchedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addHistoryEntry(entry);
      final history = await database.getHistory();

      expect(history.length, 1);
      expect(history[0].animeId, 1);
      expect(history[0].episodeNumber, 5);
    });

    test('get history for anime', () async {
      final entry1 = HistoryCompanion(
        animeId: const Value(1),
        episodeNumber: const Value(1),
        watchedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      final entry2 = HistoryCompanion(
        animeId: const Value(1),
        episodeNumber: const Value(2),
        watchedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      final entry3 = HistoryCompanion(
        animeId: const Value(2),
        episodeNumber: const Value(1),
        watchedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addHistoryEntry(entry1);
      await database.addHistoryEntry(entry2);
      await database.addHistoryEntry(entry3);

      final animeHistory = await database.getHistoryForAnime(1);

      expect(animeHistory.length, 2);
      expect(animeHistory.every((e) => e.animeId == 1), true);
    });

    test('delete history entry', () async {
      final entry = HistoryCompanion(
        animeId: const Value(1),
        episodeNumber: const Value(5),
        watchedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );

      await database.addHistoryEntry(entry);
      await database.deleteHistoryEntry(1);

      final history = await database.getHistory();
      expect(history.length, 0);
    });
  });

  group('SyncQueue Table', () {
    test('add sync queue item', () async {
      final item = SyncQueueCompanion(
        operation: const Value('add'),
        entityType: const Value('watchlist'),
        entityId: const Value('1'),
        data: const Value('{"animeId": 1}'),
        createdAt: Value(DateTime.now()),
      );

      await database.addSyncQueueItem(item);
      final pending = await database.getPendingSyncItems();

      expect(pending.length, 1);
      expect(pending[0].operation, 'add');
      expect(pending[0].synced, false);
    });

    test('mark sync item as synced', () async {
      final item = SyncQueueCompanion(
        operation: const Value('add'),
        entityType: const Value('watchlist'),
        entityId: const Value('1'),
        data: const Value('{"animeId": 1}'),
        createdAt: Value(DateTime.now()),
      );

      await database.addSyncQueueItem(item);
      await database.markSyncItemAsSynced(1);

      final pending = await database.getPendingSyncItems();
      expect(pending.length, 0);
    });

    test('increment sync retry', () async {
      final item = SyncQueueCompanion(
        operation: const Value('add'),
        entityType: const Value('watchlist'),
        entityId: const Value('1'),
        data: const Value('{"animeId": 1}'),
        createdAt: Value(DateTime.now()),
      );

      await database.addSyncQueueItem(item);
      await database.incrementSyncRetry(1, 'Connection timeout');

      final pending = await database.getPendingSyncItems();
      expect(pending[0].retryCount, 1);
      expect(pending[0].lastError, 'Connection timeout');
    });

    test('delete sync queue item', () async {
      final item = SyncQueueCompanion(
        operation: const Value('add'),
        entityType: const Value('watchlist'),
        entityId: const Value('1'),
        data: const Value('{"animeId": 1}'),
        createdAt: Value(DateTime.now()),
      );

      await database.addSyncQueueItem(item);
      await database.deleteSyncQueueItem(1);

      final pending = await database.getPendingSyncItems();
      expect(pending.length, 0);
    });

    test('clear sync queue', () async {
      for (int i = 0; i < 5; i++) {
        final item = SyncQueueCompanion(
          operation: const Value('add'),
          entityType: const Value('watchlist'),
          entityId: Value('$i'),
          data: const Value('{}'),
          createdAt: Value(DateTime.now()),
        );
        await database.addSyncQueueItem(item);
      }

      await database.clearSyncQueue();

      final pending = await database.getPendingSyncItems();
      expect(pending.length, 0);
    });
  });

  group('VIP Status Table', () {
    test('update vip status', () async {
      final status = VipStatusCompanion(
        isVip: const Value(true),
        vipTier: const Value('premium'),
        vipExpiresAt: Value(DateTime.now().add(const Duration(days: 30))),
      );

      await database.updateVipStatus(status);
      final retrieved = await database.getVipStatus();

      expect(retrieved, matcher.isNotNull);
      expect(retrieved!.isVip, true);
      expect(retrieved.vipTier, 'premium');
    });

    test('get vip status returns null when not set', () async {
      final status = await database.getVipStatus();
      expect(status, matcher.isNull);
    });
  });

  group('Database Statistics', () {
    test('get statistics', () async {
      // Add some data
      final anime = AnimesCompanion(
        title: const Value('Test'),
        episodes: const Value(12),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );
      await database.insertAnime(anime);

      final watchlistItem = WatchlistCompanion(
        animeId: const Value(1),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );
      await database.addToWatchlist(watchlistItem);

      final stats = await database.getStatistics();

      expect(stats['animes'], 1);
      expect(stats['watchlist'], 1);
      expect(stats['history'], 0);
      expect(stats['pendingSync'], 0);
      expect(stats['conflicts'], 0);
    });
  });

  group('Database Cleanup', () {
    test('clear all data', () async {
      // Add some data
      final anime = AnimesCompanion(
        title: const Value('Test'),
        episodes: const Value(12),
        addedAt: Value(DateTime.now()),
        updatedAt: Value(DateTime.now()),
      );
      await database.insertAnime(anime);

      await database.clearAllData();

      final animes = await database.getAllAnimes();
      expect(animes.length, 0);
    });
  });
}
