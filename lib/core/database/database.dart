import 'dart:io';
import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:path_provider/path_provider.dart';
import 'package:path/path.dart' as p;

part 'database.g.dart';

// Table definitions

/// Animes table - stores anime information
@DataClassName('AnimeData')
class Animes extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get title => text()();
  TextColumn get description => text().nullable()();
  IntColumn get episodes => integer()();
  TextColumn get imageUrl => text().nullable()();
  TextColumn get genre => text().nullable()();
  TextColumn get status => text().nullable()(); // 'ongoing', 'completed', 'upcoming'
  DateTimeColumn get addedAt => dateTime()();
  DateTimeColumn get syncedAt => dateTime().nullable()();
  DateTimeColumn get updatedAt => dateTime()();
}

/// Watchlist table - stores user's watchlist
@DataClassName('WatchlistData')
class Watchlist extends Table {
  IntColumn get id => integer().autoIncrement()();
  IntColumn get animeId => integer()();
  TextColumn get status => text().withDefault(const Constant('watching'))(); // 'watching', 'completed', 'dropped', 'planned'
  IntColumn get currentEpisode => integer().withDefault(const Constant(0))();
  DateTimeColumn get addedAt => dateTime()();
  DateTimeColumn get syncedAt => dateTime().nullable()();
  DateTimeColumn get updatedAt => dateTime()();

  @override
  List<Set<Column>> get uniqueKeys => [
    {animeId}
  ];
}

/// History table - stores watch history
@DataClassName('HistoryData')
class History extends Table {
  IntColumn get id => integer().autoIncrement()();
  IntColumn get animeId => integer()();
  IntColumn get episodeNumber => integer()();
  TextColumn get episodeTitle => text().nullable()();
  DateTimeColumn get watchedAt => dateTime()();
  DateTimeColumn get syncedAt => dateTime().nullable()();
  DateTimeColumn get updatedAt => dateTime()();
}

/// SyncQueue table - stores pending sync operations
@DataClassName('SyncQueueData')
class SyncQueue extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get operation => text()(); // 'add', 'remove', 'update'
  TextColumn get entityType => text()(); // 'watchlist', 'history', 'vip'
  TextColumn get entityId => text()();
  TextColumn get data => text()(); // JSON serialized data
  DateTimeColumn get createdAt => dateTime()();
  BoolColumn get synced => boolean().withDefault(const Constant(false))();
  IntColumn get retryCount => integer().withDefault(const Constant(0))();
  TextColumn get lastError => text().nullable()();
}

/// VIP Status table - stores VIP information
@DataClassName('VipStatusData')
class VipStatus extends Table {
  IntColumn get id => integer().autoIncrement()();
  BoolColumn get isVip => boolean().withDefault(const Constant(false))();
  DateTimeColumn get vipExpiresAt => dateTime().nullable()();
  TextColumn get vipTier => text().nullable()(); // 'basic', 'premium', 'ultimate'
  DateTimeColumn get syncedAt => dateTime().nullable()();
  DateTimeColumn get updatedAt => dateTime()();
}

/// SyncConflict table - stores sync conflicts for resolution
@DataClassName('SyncConflictData')
class SyncConflicts extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get entityType => text()();
  TextColumn get entityId => text()();
  TextColumn get mobileData => text()(); // JSON
  TextColumn get desktopData => text()(); // JSON
  DateTimeColumn get mobileTimestamp => dateTime()();
  DateTimeColumn get desktopTimestamp => dateTime()();
  TextColumn get resolution => text().nullable()(); // 'mobile', 'desktop', 'merged'
  TextColumn get resolvedData => text().nullable()(); // JSON
  DateTimeColumn get createdAt => dateTime()();
  DateTimeColumn get resolvedAt => dateTime().nullable()();
}

/// Main database class
@DriftDatabase(tables: [Animes, Watchlist, History, SyncQueue, VipStatus, SyncConflicts])
class AppDatabase extends _$AppDatabase {
  AppDatabase() : super(_openConnection());

  // For testing
  AppDatabase.forTesting() : super(_testingConnection());

  @override
  int get schemaVersion => 1;

  @override
  MigrationStrategy get migration {
    return MigrationStrategy(
      onCreate: (Migrator m) async {
        await m.createAll();
      },
      onUpgrade: (Migrator m, int from, int to) async {
        // Handle migrations here as schema evolves
      },
    );
  }

  // Anime queries
  Future<List<AnimeData>> getAllAnimes() => select(animes).get();

  Future<AnimeData?> getAnimeById(int id) =>
      (select(animes)..where((t) => t.id.equals(id))).getSingleOrNull();

  Future<void> insertAnime(AnimesCompanion anime) => into(animes).insert(anime);

  Future<bool> updateAnime(AnimeData anime) => update(animes).replace(anime);

  Future<int> deleteAnime(int id) =>
      (delete(animes)..where((t) => t.id.equals(id))).go();

  // Watchlist queries
  Future<List<WatchlistData>> getWatchlist() => select(watchlist).get();

  Future<WatchlistData?> getWatchlistItem(int animeId) =>
      (select(watchlist)..where((t) => t.animeId.equals(animeId))).getSingleOrNull();

  Future<void> addToWatchlist(WatchlistCompanion item) =>
      into(watchlist).insert(item);

  Future<bool> updateWatchlistItem(WatchlistData item) =>
      update(watchlist).replace(item);

  Future<int> removeFromWatchlist(int animeId) =>
      (delete(watchlist)..where((t) => t.animeId.equals(animeId))).go();

  // History queries
  Future<List<HistoryData>> getHistory() => select(history).get();

  Future<List<HistoryData>> getHistoryForAnime(int animeId) =>
      (select(history)..where((t) => t.animeId.equals(animeId))).get();

  Future<void> addHistoryEntry(HistoryCompanion entry) =>
      into(history).insert(entry);

  Future<int> deleteHistoryEntry(int id) =>
      (delete(history)..where((t) => t.id.equals(id))).go();

  // Sync Queue queries
  Future<List<SyncQueueData>> getPendingSyncItems() =>
      (select(syncQueue)..where((t) => t.synced.equals(false))).get();

  Future<void> addSyncQueueItem(SyncQueueCompanion item) =>
      into(syncQueue).insert(item);

  Future<bool> markSyncItemAsSynced(int id) async {
    final item = await (select(syncQueue)..where((t) => t.id.equals(id)))
        .getSingleOrNull();
    if (item == null) return false;

    return update(syncQueue).replace(
      item.copyWith(synced: true, retryCount: 0, lastError: const Value(null)),
    );
  }

  Future<bool> incrementSyncRetry(int id, String? error) async {
    final item = await (select(syncQueue)..where((t) => t.id.equals(id)))
        .getSingleOrNull();
    if (item == null) return false;

    return update(syncQueue).replace(
      item.copyWith(
        retryCount: item.retryCount + 1,
        lastError: Value(error),
      ),
    );
  }

  Future<int> deleteSyncQueueItem(int id) =>
      (delete(syncQueue)..where((t) => t.id.equals(id))).go();

  Future<void> clearSyncQueue() => delete(syncQueue).go();

  // VIP Status queries
  Future<VipStatusData?> getVipStatus() => select(vipStatus).getSingleOrNull();

  Future<void> updateVipStatus(VipStatusCompanion status) async {
    final existing = await getVipStatus();
    if (existing != null) {
      // Create a new VipStatusData with updated values
      final updated = VipStatusData(
        id: existing.id,
        isVip: status.isVip.present ? status.isVip.value : existing.isVip,
        vipExpiresAt: status.vipExpiresAt.present
            ? status.vipExpiresAt.value
            : existing.vipExpiresAt,
        vipTier: status.vipTier.present ? status.vipTier.value : existing.vipTier,
        syncedAt: null,
        updatedAt: DateTime.now(),
      );
      await update(vipStatus).replace(updated);
    } else {
      await into(vipStatus).insert(
        status.copyWith(
          syncedAt: const Value(null),
          updatedAt: Value(DateTime.now()),
        ),
      );
    }
  }

  // Sync Conflict queries
  Future<List<SyncConflictData>> getUnresolvedConflicts() =>
      (select(syncConflicts)..where((t) => t.resolution.isNull())).get();

  Future<void> addSyncConflict(SyncConflictsCompanion conflict) =>
      into(syncConflicts).insert(conflict);

  Future<bool> resolveSyncConflict(
    int id,
    String resolution,
    String? resolvedData,
  ) async {
    final conflict = await (select(syncConflicts)..where((t) => t.id.equals(id)))
        .getSingleOrNull();
    if (conflict == null) return false;

    return update(syncConflicts).replace(
      conflict.copyWith(
        resolution: Value(resolution),
        resolvedData: Value(resolvedData),
        resolvedAt: Value(DateTime.now()),
      ),
    );
  }

  // Utility methods
  Future<void> clearAllData() async {
    await delete(animes).go();
    await delete(watchlist).go();
    await delete(history).go();
    await delete(syncQueue).go();
    await delete(vipStatus).go();
    await delete(syncConflicts).go();
  }

  Future<Map<String, int>> getStatistics() async {
    final animeCount = await select(animes).get().then((list) => list.length);
    final watchlistCount =
        await select(watchlist).get().then((list) => list.length);
    final historyCount = await select(history).get().then((list) => list.length);
    final pendingSyncCount = await getPendingSyncItems()
        .then((list) => list.length);
    final conflictCount = await getUnresolvedConflicts()
        .then((list) => list.length);

    return {
      'animes': animeCount,
      'watchlist': watchlistCount,
      'history': historyCount,
      'pendingSync': pendingSyncCount,
      'conflicts': conflictCount,
    };
  }
}

LazyDatabase _openConnection() {
  return LazyDatabase(() async {
    final dbFolder = await getApplicationDocumentsDirectory();
    final file = File(p.join(dbFolder.path, 'app_database.db'));
    return NativeDatabase(file);
  });
}

LazyDatabase _testingConnection() {
  return LazyDatabase(() async {
    return NativeDatabase.memory();
  });
}
