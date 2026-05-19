import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:stardf_anime_mobile/core/database/database.dart';

/// Provider for the AppDatabase singleton
final databaseProvider = Provider<AppDatabase>((ref) {
  return AppDatabase();
});

/// Provider for anime list
final animeListProvider = FutureProvider<List<AnimeData>>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getAllAnimes();
});

/// Provider for watchlist
final watchlistProvider = FutureProvider<List<WatchlistData>>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getWatchlist();
});

/// Provider for history
final historyProvider = FutureProvider<List<HistoryData>>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getHistory();
});

/// Provider for pending sync items
final pendingSyncProvider = FutureProvider<List<SyncQueueData>>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getPendingSyncItems();
});

/// Provider for VIP status
final vipStatusProvider = FutureProvider<VipStatusData?>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getVipStatus();
});

/// Provider for unresolved sync conflicts
final syncConflictsProvider = FutureProvider<List<SyncConflictData>>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getUnresolvedConflicts();
});

/// Provider for database statistics
final databaseStatsProvider = FutureProvider<Map<String, int>>((ref) async {
  final db = ref.watch(databaseProvider);
  return db.getStatistics();
});
