// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'database.dart';

// ignore_for_file: type=lint
class $AnimesTable extends Animes with TableInfo<$AnimesTable, AnimeData> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $AnimesTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
      'id', aliasedName, false,
      hasAutoIncrement: true,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('PRIMARY KEY AUTOINCREMENT'));
  static const VerificationMeta _titleMeta = const VerificationMeta('title');
  @override
  late final GeneratedColumn<String> title = GeneratedColumn<String>(
      'title', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _descriptionMeta =
      const VerificationMeta('description');
  @override
  late final GeneratedColumn<String> description = GeneratedColumn<String>(
      'description', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _episodesMeta =
      const VerificationMeta('episodes');
  @override
  late final GeneratedColumn<int> episodes = GeneratedColumn<int>(
      'episodes', aliasedName, false,
      type: DriftSqlType.int, requiredDuringInsert: true);
  static const VerificationMeta _imageUrlMeta =
      const VerificationMeta('imageUrl');
  @override
  late final GeneratedColumn<String> imageUrl = GeneratedColumn<String>(
      'image_url', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _genreMeta = const VerificationMeta('genre');
  @override
  late final GeneratedColumn<String> genre = GeneratedColumn<String>(
      'genre', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _statusMeta = const VerificationMeta('status');
  @override
  late final GeneratedColumn<String> status = GeneratedColumn<String>(
      'status', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _addedAtMeta =
      const VerificationMeta('addedAt');
  @override
  late final GeneratedColumn<DateTime> addedAt = GeneratedColumn<DateTime>(
      'added_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _syncedAtMeta =
      const VerificationMeta('syncedAt');
  @override
  late final GeneratedColumn<DateTime> syncedAt = GeneratedColumn<DateTime>(
      'synced_at', aliasedName, true,
      type: DriftSqlType.dateTime, requiredDuringInsert: false);
  static const VerificationMeta _updatedAtMeta =
      const VerificationMeta('updatedAt');
  @override
  late final GeneratedColumn<DateTime> updatedAt = GeneratedColumn<DateTime>(
      'updated_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  @override
  List<GeneratedColumn> get $columns => [
        id,
        title,
        description,
        episodes,
        imageUrl,
        genre,
        status,
        addedAt,
        syncedAt,
        updatedAt
      ];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'animes';
  @override
  VerificationContext validateIntegrity(Insertable<AnimeData> instance,
      {bool isInserting = false}) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('title')) {
      context.handle(
          _titleMeta, title.isAcceptableOrUnknown(data['title']!, _titleMeta));
    } else if (isInserting) {
      context.missing(_titleMeta);
    }
    if (data.containsKey('description')) {
      context.handle(
          _descriptionMeta,
          description.isAcceptableOrUnknown(
              data['description']!, _descriptionMeta));
    }
    if (data.containsKey('episodes')) {
      context.handle(_episodesMeta,
          episodes.isAcceptableOrUnknown(data['episodes']!, _episodesMeta));
    } else if (isInserting) {
      context.missing(_episodesMeta);
    }
    if (data.containsKey('image_url')) {
      context.handle(_imageUrlMeta,
          imageUrl.isAcceptableOrUnknown(data['image_url']!, _imageUrlMeta));
    }
    if (data.containsKey('genre')) {
      context.handle(
          _genreMeta, genre.isAcceptableOrUnknown(data['genre']!, _genreMeta));
    }
    if (data.containsKey('status')) {
      context.handle(_statusMeta,
          status.isAcceptableOrUnknown(data['status']!, _statusMeta));
    }
    if (data.containsKey('added_at')) {
      context.handle(_addedAtMeta,
          addedAt.isAcceptableOrUnknown(data['added_at']!, _addedAtMeta));
    } else if (isInserting) {
      context.missing(_addedAtMeta);
    }
    if (data.containsKey('synced_at')) {
      context.handle(_syncedAtMeta,
          syncedAt.isAcceptableOrUnknown(data['synced_at']!, _syncedAtMeta));
    }
    if (data.containsKey('updated_at')) {
      context.handle(_updatedAtMeta,
          updatedAt.isAcceptableOrUnknown(data['updated_at']!, _updatedAtMeta));
    } else if (isInserting) {
      context.missing(_updatedAtMeta);
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  AnimeData map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return AnimeData(
      id: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}id'])!,
      title: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}title'])!,
      description: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}description']),
      episodes: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}episodes'])!,
      imageUrl: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}image_url']),
      genre: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}genre']),
      status: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}status']),
      addedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}added_at'])!,
      syncedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}synced_at']),
      updatedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}updated_at'])!,
    );
  }

  @override
  $AnimesTable createAlias(String alias) {
    return $AnimesTable(attachedDatabase, alias);
  }
}

class AnimeData extends DataClass implements Insertable<AnimeData> {
  final int id;
  final String title;
  final String? description;
  final int episodes;
  final String? imageUrl;
  final String? genre;
  final String? status;
  final DateTime addedAt;
  final DateTime? syncedAt;
  final DateTime updatedAt;
  const AnimeData(
      {required this.id,
      required this.title,
      this.description,
      required this.episodes,
      this.imageUrl,
      this.genre,
      this.status,
      required this.addedAt,
      this.syncedAt,
      required this.updatedAt});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['title'] = Variable<String>(title);
    if (!nullToAbsent || description != null) {
      map['description'] = Variable<String>(description);
    }
    map['episodes'] = Variable<int>(episodes);
    if (!nullToAbsent || imageUrl != null) {
      map['image_url'] = Variable<String>(imageUrl);
    }
    if (!nullToAbsent || genre != null) {
      map['genre'] = Variable<String>(genre);
    }
    if (!nullToAbsent || status != null) {
      map['status'] = Variable<String>(status);
    }
    map['added_at'] = Variable<DateTime>(addedAt);
    if (!nullToAbsent || syncedAt != null) {
      map['synced_at'] = Variable<DateTime>(syncedAt);
    }
    map['updated_at'] = Variable<DateTime>(updatedAt);
    return map;
  }

  AnimesCompanion toCompanion(bool nullToAbsent) {
    return AnimesCompanion(
      id: Value(id),
      title: Value(title),
      description: description == null && nullToAbsent
          ? const Value.absent()
          : Value(description),
      episodes: Value(episodes),
      imageUrl: imageUrl == null && nullToAbsent
          ? const Value.absent()
          : Value(imageUrl),
      genre:
          genre == null && nullToAbsent ? const Value.absent() : Value(genre),
      status:
          status == null && nullToAbsent ? const Value.absent() : Value(status),
      addedAt: Value(addedAt),
      syncedAt: syncedAt == null && nullToAbsent
          ? const Value.absent()
          : Value(syncedAt),
      updatedAt: Value(updatedAt),
    );
  }

  factory AnimeData.fromJson(Map<String, dynamic> json,
      {ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return AnimeData(
      id: serializer.fromJson<int>(json['id']),
      title: serializer.fromJson<String>(json['title']),
      description: serializer.fromJson<String?>(json['description']),
      episodes: serializer.fromJson<int>(json['episodes']),
      imageUrl: serializer.fromJson<String?>(json['imageUrl']),
      genre: serializer.fromJson<String?>(json['genre']),
      status: serializer.fromJson<String?>(json['status']),
      addedAt: serializer.fromJson<DateTime>(json['addedAt']),
      syncedAt: serializer.fromJson<DateTime?>(json['syncedAt']),
      updatedAt: serializer.fromJson<DateTime>(json['updatedAt']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'title': serializer.toJson<String>(title),
      'description': serializer.toJson<String?>(description),
      'episodes': serializer.toJson<int>(episodes),
      'imageUrl': serializer.toJson<String?>(imageUrl),
      'genre': serializer.toJson<String?>(genre),
      'status': serializer.toJson<String?>(status),
      'addedAt': serializer.toJson<DateTime>(addedAt),
      'syncedAt': serializer.toJson<DateTime?>(syncedAt),
      'updatedAt': serializer.toJson<DateTime>(updatedAt),
    };
  }

  AnimeData copyWith(
          {int? id,
          String? title,
          Value<String?> description = const Value.absent(),
          int? episodes,
          Value<String?> imageUrl = const Value.absent(),
          Value<String?> genre = const Value.absent(),
          Value<String?> status = const Value.absent(),
          DateTime? addedAt,
          Value<DateTime?> syncedAt = const Value.absent(),
          DateTime? updatedAt}) =>
      AnimeData(
        id: id ?? this.id,
        title: title ?? this.title,
        description: description.present ? description.value : this.description,
        episodes: episodes ?? this.episodes,
        imageUrl: imageUrl.present ? imageUrl.value : this.imageUrl,
        genre: genre.present ? genre.value : this.genre,
        status: status.present ? status.value : this.status,
        addedAt: addedAt ?? this.addedAt,
        syncedAt: syncedAt.present ? syncedAt.value : this.syncedAt,
        updatedAt: updatedAt ?? this.updatedAt,
      );
  AnimeData copyWithCompanion(AnimesCompanion data) {
    return AnimeData(
      id: data.id.present ? data.id.value : this.id,
      title: data.title.present ? data.title.value : this.title,
      description:
          data.description.present ? data.description.value : this.description,
      episodes: data.episodes.present ? data.episodes.value : this.episodes,
      imageUrl: data.imageUrl.present ? data.imageUrl.value : this.imageUrl,
      genre: data.genre.present ? data.genre.value : this.genre,
      status: data.status.present ? data.status.value : this.status,
      addedAt: data.addedAt.present ? data.addedAt.value : this.addedAt,
      syncedAt: data.syncedAt.present ? data.syncedAt.value : this.syncedAt,
      updatedAt: data.updatedAt.present ? data.updatedAt.value : this.updatedAt,
    );
  }

  @override
  String toString() {
    return (StringBuffer('AnimeData(')
          ..write('id: $id, ')
          ..write('title: $title, ')
          ..write('description: $description, ')
          ..write('episodes: $episodes, ')
          ..write('imageUrl: $imageUrl, ')
          ..write('genre: $genre, ')
          ..write('status: $status, ')
          ..write('addedAt: $addedAt, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(id, title, description, episodes, imageUrl,
      genre, status, addedAt, syncedAt, updatedAt);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is AnimeData &&
          other.id == this.id &&
          other.title == this.title &&
          other.description == this.description &&
          other.episodes == this.episodes &&
          other.imageUrl == this.imageUrl &&
          other.genre == this.genre &&
          other.status == this.status &&
          other.addedAt == this.addedAt &&
          other.syncedAt == this.syncedAt &&
          other.updatedAt == this.updatedAt);
}

class AnimesCompanion extends UpdateCompanion<AnimeData> {
  final Value<int> id;
  final Value<String> title;
  final Value<String?> description;
  final Value<int> episodes;
  final Value<String?> imageUrl;
  final Value<String?> genre;
  final Value<String?> status;
  final Value<DateTime> addedAt;
  final Value<DateTime?> syncedAt;
  final Value<DateTime> updatedAt;
  const AnimesCompanion({
    this.id = const Value.absent(),
    this.title = const Value.absent(),
    this.description = const Value.absent(),
    this.episodes = const Value.absent(),
    this.imageUrl = const Value.absent(),
    this.genre = const Value.absent(),
    this.status = const Value.absent(),
    this.addedAt = const Value.absent(),
    this.syncedAt = const Value.absent(),
    this.updatedAt = const Value.absent(),
  });
  AnimesCompanion.insert({
    this.id = const Value.absent(),
    required String title,
    this.description = const Value.absent(),
    required int episodes,
    this.imageUrl = const Value.absent(),
    this.genre = const Value.absent(),
    this.status = const Value.absent(),
    required DateTime addedAt,
    this.syncedAt = const Value.absent(),
    required DateTime updatedAt,
  })  : title = Value(title),
        episodes = Value(episodes),
        addedAt = Value(addedAt),
        updatedAt = Value(updatedAt);
  static Insertable<AnimeData> custom({
    Expression<int>? id,
    Expression<String>? title,
    Expression<String>? description,
    Expression<int>? episodes,
    Expression<String>? imageUrl,
    Expression<String>? genre,
    Expression<String>? status,
    Expression<DateTime>? addedAt,
    Expression<DateTime>? syncedAt,
    Expression<DateTime>? updatedAt,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (title != null) 'title': title,
      if (description != null) 'description': description,
      if (episodes != null) 'episodes': episodes,
      if (imageUrl != null) 'image_url': imageUrl,
      if (genre != null) 'genre': genre,
      if (status != null) 'status': status,
      if (addedAt != null) 'added_at': addedAt,
      if (syncedAt != null) 'synced_at': syncedAt,
      if (updatedAt != null) 'updated_at': updatedAt,
    });
  }

  AnimesCompanion copyWith(
      {Value<int>? id,
      Value<String>? title,
      Value<String?>? description,
      Value<int>? episodes,
      Value<String?>? imageUrl,
      Value<String?>? genre,
      Value<String?>? status,
      Value<DateTime>? addedAt,
      Value<DateTime?>? syncedAt,
      Value<DateTime>? updatedAt}) {
    return AnimesCompanion(
      id: id ?? this.id,
      title: title ?? this.title,
      description: description ?? this.description,
      episodes: episodes ?? this.episodes,
      imageUrl: imageUrl ?? this.imageUrl,
      genre: genre ?? this.genre,
      status: status ?? this.status,
      addedAt: addedAt ?? this.addedAt,
      syncedAt: syncedAt ?? this.syncedAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (title.present) {
      map['title'] = Variable<String>(title.value);
    }
    if (description.present) {
      map['description'] = Variable<String>(description.value);
    }
    if (episodes.present) {
      map['episodes'] = Variable<int>(episodes.value);
    }
    if (imageUrl.present) {
      map['image_url'] = Variable<String>(imageUrl.value);
    }
    if (genre.present) {
      map['genre'] = Variable<String>(genre.value);
    }
    if (status.present) {
      map['status'] = Variable<String>(status.value);
    }
    if (addedAt.present) {
      map['added_at'] = Variable<DateTime>(addedAt.value);
    }
    if (syncedAt.present) {
      map['synced_at'] = Variable<DateTime>(syncedAt.value);
    }
    if (updatedAt.present) {
      map['updated_at'] = Variable<DateTime>(updatedAt.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('AnimesCompanion(')
          ..write('id: $id, ')
          ..write('title: $title, ')
          ..write('description: $description, ')
          ..write('episodes: $episodes, ')
          ..write('imageUrl: $imageUrl, ')
          ..write('genre: $genre, ')
          ..write('status: $status, ')
          ..write('addedAt: $addedAt, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }
}

class $WatchlistTable extends Watchlist
    with TableInfo<$WatchlistTable, WatchlistData> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $WatchlistTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
      'id', aliasedName, false,
      hasAutoIncrement: true,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('PRIMARY KEY AUTOINCREMENT'));
  static const VerificationMeta _animeIdMeta =
      const VerificationMeta('animeId');
  @override
  late final GeneratedColumn<int> animeId = GeneratedColumn<int>(
      'anime_id', aliasedName, false,
      type: DriftSqlType.int, requiredDuringInsert: true);
  static const VerificationMeta _statusMeta = const VerificationMeta('status');
  @override
  late final GeneratedColumn<String> status = GeneratedColumn<String>(
      'status', aliasedName, false,
      type: DriftSqlType.string,
      requiredDuringInsert: false,
      defaultValue: const Constant('watching'));
  static const VerificationMeta _currentEpisodeMeta =
      const VerificationMeta('currentEpisode');
  @override
  late final GeneratedColumn<int> currentEpisode = GeneratedColumn<int>(
      'current_episode', aliasedName, false,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultValue: const Constant(0));
  static const VerificationMeta _addedAtMeta =
      const VerificationMeta('addedAt');
  @override
  late final GeneratedColumn<DateTime> addedAt = GeneratedColumn<DateTime>(
      'added_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _syncedAtMeta =
      const VerificationMeta('syncedAt');
  @override
  late final GeneratedColumn<DateTime> syncedAt = GeneratedColumn<DateTime>(
      'synced_at', aliasedName, true,
      type: DriftSqlType.dateTime, requiredDuringInsert: false);
  static const VerificationMeta _updatedAtMeta =
      const VerificationMeta('updatedAt');
  @override
  late final GeneratedColumn<DateTime> updatedAt = GeneratedColumn<DateTime>(
      'updated_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  @override
  List<GeneratedColumn> get $columns =>
      [id, animeId, status, currentEpisode, addedAt, syncedAt, updatedAt];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'watchlist';
  @override
  VerificationContext validateIntegrity(Insertable<WatchlistData> instance,
      {bool isInserting = false}) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('anime_id')) {
      context.handle(_animeIdMeta,
          animeId.isAcceptableOrUnknown(data['anime_id']!, _animeIdMeta));
    } else if (isInserting) {
      context.missing(_animeIdMeta);
    }
    if (data.containsKey('status')) {
      context.handle(_statusMeta,
          status.isAcceptableOrUnknown(data['status']!, _statusMeta));
    }
    if (data.containsKey('current_episode')) {
      context.handle(
          _currentEpisodeMeta,
          currentEpisode.isAcceptableOrUnknown(
              data['current_episode']!, _currentEpisodeMeta));
    }
    if (data.containsKey('added_at')) {
      context.handle(_addedAtMeta,
          addedAt.isAcceptableOrUnknown(data['added_at']!, _addedAtMeta));
    } else if (isInserting) {
      context.missing(_addedAtMeta);
    }
    if (data.containsKey('synced_at')) {
      context.handle(_syncedAtMeta,
          syncedAt.isAcceptableOrUnknown(data['synced_at']!, _syncedAtMeta));
    }
    if (data.containsKey('updated_at')) {
      context.handle(_updatedAtMeta,
          updatedAt.isAcceptableOrUnknown(data['updated_at']!, _updatedAtMeta));
    } else if (isInserting) {
      context.missing(_updatedAtMeta);
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  List<Set<GeneratedColumn>> get uniqueKeys => [
        {animeId},
      ];
  @override
  WatchlistData map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return WatchlistData(
      id: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}id'])!,
      animeId: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}anime_id'])!,
      status: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}status'])!,
      currentEpisode: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}current_episode'])!,
      addedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}added_at'])!,
      syncedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}synced_at']),
      updatedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}updated_at'])!,
    );
  }

  @override
  $WatchlistTable createAlias(String alias) {
    return $WatchlistTable(attachedDatabase, alias);
  }
}

class WatchlistData extends DataClass implements Insertable<WatchlistData> {
  final int id;
  final int animeId;
  final String status;
  final int currentEpisode;
  final DateTime addedAt;
  final DateTime? syncedAt;
  final DateTime updatedAt;
  const WatchlistData(
      {required this.id,
      required this.animeId,
      required this.status,
      required this.currentEpisode,
      required this.addedAt,
      this.syncedAt,
      required this.updatedAt});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['anime_id'] = Variable<int>(animeId);
    map['status'] = Variable<String>(status);
    map['current_episode'] = Variable<int>(currentEpisode);
    map['added_at'] = Variable<DateTime>(addedAt);
    if (!nullToAbsent || syncedAt != null) {
      map['synced_at'] = Variable<DateTime>(syncedAt);
    }
    map['updated_at'] = Variable<DateTime>(updatedAt);
    return map;
  }

  WatchlistCompanion toCompanion(bool nullToAbsent) {
    return WatchlistCompanion(
      id: Value(id),
      animeId: Value(animeId),
      status: Value(status),
      currentEpisode: Value(currentEpisode),
      addedAt: Value(addedAt),
      syncedAt: syncedAt == null && nullToAbsent
          ? const Value.absent()
          : Value(syncedAt),
      updatedAt: Value(updatedAt),
    );
  }

  factory WatchlistData.fromJson(Map<String, dynamic> json,
      {ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return WatchlistData(
      id: serializer.fromJson<int>(json['id']),
      animeId: serializer.fromJson<int>(json['animeId']),
      status: serializer.fromJson<String>(json['status']),
      currentEpisode: serializer.fromJson<int>(json['currentEpisode']),
      addedAt: serializer.fromJson<DateTime>(json['addedAt']),
      syncedAt: serializer.fromJson<DateTime?>(json['syncedAt']),
      updatedAt: serializer.fromJson<DateTime>(json['updatedAt']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'animeId': serializer.toJson<int>(animeId),
      'status': serializer.toJson<String>(status),
      'currentEpisode': serializer.toJson<int>(currentEpisode),
      'addedAt': serializer.toJson<DateTime>(addedAt),
      'syncedAt': serializer.toJson<DateTime?>(syncedAt),
      'updatedAt': serializer.toJson<DateTime>(updatedAt),
    };
  }

  WatchlistData copyWith(
          {int? id,
          int? animeId,
          String? status,
          int? currentEpisode,
          DateTime? addedAt,
          Value<DateTime?> syncedAt = const Value.absent(),
          DateTime? updatedAt}) =>
      WatchlistData(
        id: id ?? this.id,
        animeId: animeId ?? this.animeId,
        status: status ?? this.status,
        currentEpisode: currentEpisode ?? this.currentEpisode,
        addedAt: addedAt ?? this.addedAt,
        syncedAt: syncedAt.present ? syncedAt.value : this.syncedAt,
        updatedAt: updatedAt ?? this.updatedAt,
      );
  WatchlistData copyWithCompanion(WatchlistCompanion data) {
    return WatchlistData(
      id: data.id.present ? data.id.value : this.id,
      animeId: data.animeId.present ? data.animeId.value : this.animeId,
      status: data.status.present ? data.status.value : this.status,
      currentEpisode: data.currentEpisode.present
          ? data.currentEpisode.value
          : this.currentEpisode,
      addedAt: data.addedAt.present ? data.addedAt.value : this.addedAt,
      syncedAt: data.syncedAt.present ? data.syncedAt.value : this.syncedAt,
      updatedAt: data.updatedAt.present ? data.updatedAt.value : this.updatedAt,
    );
  }

  @override
  String toString() {
    return (StringBuffer('WatchlistData(')
          ..write('id: $id, ')
          ..write('animeId: $animeId, ')
          ..write('status: $status, ')
          ..write('currentEpisode: $currentEpisode, ')
          ..write('addedAt: $addedAt, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(
      id, animeId, status, currentEpisode, addedAt, syncedAt, updatedAt);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is WatchlistData &&
          other.id == this.id &&
          other.animeId == this.animeId &&
          other.status == this.status &&
          other.currentEpisode == this.currentEpisode &&
          other.addedAt == this.addedAt &&
          other.syncedAt == this.syncedAt &&
          other.updatedAt == this.updatedAt);
}

class WatchlistCompanion extends UpdateCompanion<WatchlistData> {
  final Value<int> id;
  final Value<int> animeId;
  final Value<String> status;
  final Value<int> currentEpisode;
  final Value<DateTime> addedAt;
  final Value<DateTime?> syncedAt;
  final Value<DateTime> updatedAt;
  const WatchlistCompanion({
    this.id = const Value.absent(),
    this.animeId = const Value.absent(),
    this.status = const Value.absent(),
    this.currentEpisode = const Value.absent(),
    this.addedAt = const Value.absent(),
    this.syncedAt = const Value.absent(),
    this.updatedAt = const Value.absent(),
  });
  WatchlistCompanion.insert({
    this.id = const Value.absent(),
    required int animeId,
    this.status = const Value.absent(),
    this.currentEpisode = const Value.absent(),
    required DateTime addedAt,
    this.syncedAt = const Value.absent(),
    required DateTime updatedAt,
  })  : animeId = Value(animeId),
        addedAt = Value(addedAt),
        updatedAt = Value(updatedAt);
  static Insertable<WatchlistData> custom({
    Expression<int>? id,
    Expression<int>? animeId,
    Expression<String>? status,
    Expression<int>? currentEpisode,
    Expression<DateTime>? addedAt,
    Expression<DateTime>? syncedAt,
    Expression<DateTime>? updatedAt,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (animeId != null) 'anime_id': animeId,
      if (status != null) 'status': status,
      if (currentEpisode != null) 'current_episode': currentEpisode,
      if (addedAt != null) 'added_at': addedAt,
      if (syncedAt != null) 'synced_at': syncedAt,
      if (updatedAt != null) 'updated_at': updatedAt,
    });
  }

  WatchlistCompanion copyWith(
      {Value<int>? id,
      Value<int>? animeId,
      Value<String>? status,
      Value<int>? currentEpisode,
      Value<DateTime>? addedAt,
      Value<DateTime?>? syncedAt,
      Value<DateTime>? updatedAt}) {
    return WatchlistCompanion(
      id: id ?? this.id,
      animeId: animeId ?? this.animeId,
      status: status ?? this.status,
      currentEpisode: currentEpisode ?? this.currentEpisode,
      addedAt: addedAt ?? this.addedAt,
      syncedAt: syncedAt ?? this.syncedAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (animeId.present) {
      map['anime_id'] = Variable<int>(animeId.value);
    }
    if (status.present) {
      map['status'] = Variable<String>(status.value);
    }
    if (currentEpisode.present) {
      map['current_episode'] = Variable<int>(currentEpisode.value);
    }
    if (addedAt.present) {
      map['added_at'] = Variable<DateTime>(addedAt.value);
    }
    if (syncedAt.present) {
      map['synced_at'] = Variable<DateTime>(syncedAt.value);
    }
    if (updatedAt.present) {
      map['updated_at'] = Variable<DateTime>(updatedAt.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('WatchlistCompanion(')
          ..write('id: $id, ')
          ..write('animeId: $animeId, ')
          ..write('status: $status, ')
          ..write('currentEpisode: $currentEpisode, ')
          ..write('addedAt: $addedAt, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }
}

class $HistoryTable extends History with TableInfo<$HistoryTable, HistoryData> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $HistoryTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
      'id', aliasedName, false,
      hasAutoIncrement: true,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('PRIMARY KEY AUTOINCREMENT'));
  static const VerificationMeta _animeIdMeta =
      const VerificationMeta('animeId');
  @override
  late final GeneratedColumn<int> animeId = GeneratedColumn<int>(
      'anime_id', aliasedName, false,
      type: DriftSqlType.int, requiredDuringInsert: true);
  static const VerificationMeta _episodeNumberMeta =
      const VerificationMeta('episodeNumber');
  @override
  late final GeneratedColumn<int> episodeNumber = GeneratedColumn<int>(
      'episode_number', aliasedName, false,
      type: DriftSqlType.int, requiredDuringInsert: true);
  static const VerificationMeta _episodeTitleMeta =
      const VerificationMeta('episodeTitle');
  @override
  late final GeneratedColumn<String> episodeTitle = GeneratedColumn<String>(
      'episode_title', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _watchedAtMeta =
      const VerificationMeta('watchedAt');
  @override
  late final GeneratedColumn<DateTime> watchedAt = GeneratedColumn<DateTime>(
      'watched_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _syncedAtMeta =
      const VerificationMeta('syncedAt');
  @override
  late final GeneratedColumn<DateTime> syncedAt = GeneratedColumn<DateTime>(
      'synced_at', aliasedName, true,
      type: DriftSqlType.dateTime, requiredDuringInsert: false);
  static const VerificationMeta _updatedAtMeta =
      const VerificationMeta('updatedAt');
  @override
  late final GeneratedColumn<DateTime> updatedAt = GeneratedColumn<DateTime>(
      'updated_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  @override
  List<GeneratedColumn> get $columns => [
        id,
        animeId,
        episodeNumber,
        episodeTitle,
        watchedAt,
        syncedAt,
        updatedAt
      ];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'history';
  @override
  VerificationContext validateIntegrity(Insertable<HistoryData> instance,
      {bool isInserting = false}) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('anime_id')) {
      context.handle(_animeIdMeta,
          animeId.isAcceptableOrUnknown(data['anime_id']!, _animeIdMeta));
    } else if (isInserting) {
      context.missing(_animeIdMeta);
    }
    if (data.containsKey('episode_number')) {
      context.handle(
          _episodeNumberMeta,
          episodeNumber.isAcceptableOrUnknown(
              data['episode_number']!, _episodeNumberMeta));
    } else if (isInserting) {
      context.missing(_episodeNumberMeta);
    }
    if (data.containsKey('episode_title')) {
      context.handle(
          _episodeTitleMeta,
          episodeTitle.isAcceptableOrUnknown(
              data['episode_title']!, _episodeTitleMeta));
    }
    if (data.containsKey('watched_at')) {
      context.handle(_watchedAtMeta,
          watchedAt.isAcceptableOrUnknown(data['watched_at']!, _watchedAtMeta));
    } else if (isInserting) {
      context.missing(_watchedAtMeta);
    }
    if (data.containsKey('synced_at')) {
      context.handle(_syncedAtMeta,
          syncedAt.isAcceptableOrUnknown(data['synced_at']!, _syncedAtMeta));
    }
    if (data.containsKey('updated_at')) {
      context.handle(_updatedAtMeta,
          updatedAt.isAcceptableOrUnknown(data['updated_at']!, _updatedAtMeta));
    } else if (isInserting) {
      context.missing(_updatedAtMeta);
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  HistoryData map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return HistoryData(
      id: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}id'])!,
      animeId: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}anime_id'])!,
      episodeNumber: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}episode_number'])!,
      episodeTitle: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}episode_title']),
      watchedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}watched_at'])!,
      syncedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}synced_at']),
      updatedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}updated_at'])!,
    );
  }

  @override
  $HistoryTable createAlias(String alias) {
    return $HistoryTable(attachedDatabase, alias);
  }
}

class HistoryData extends DataClass implements Insertable<HistoryData> {
  final int id;
  final int animeId;
  final int episodeNumber;
  final String? episodeTitle;
  final DateTime watchedAt;
  final DateTime? syncedAt;
  final DateTime updatedAt;
  const HistoryData(
      {required this.id,
      required this.animeId,
      required this.episodeNumber,
      this.episodeTitle,
      required this.watchedAt,
      this.syncedAt,
      required this.updatedAt});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['anime_id'] = Variable<int>(animeId);
    map['episode_number'] = Variable<int>(episodeNumber);
    if (!nullToAbsent || episodeTitle != null) {
      map['episode_title'] = Variable<String>(episodeTitle);
    }
    map['watched_at'] = Variable<DateTime>(watchedAt);
    if (!nullToAbsent || syncedAt != null) {
      map['synced_at'] = Variable<DateTime>(syncedAt);
    }
    map['updated_at'] = Variable<DateTime>(updatedAt);
    return map;
  }

  HistoryCompanion toCompanion(bool nullToAbsent) {
    return HistoryCompanion(
      id: Value(id),
      animeId: Value(animeId),
      episodeNumber: Value(episodeNumber),
      episodeTitle: episodeTitle == null && nullToAbsent
          ? const Value.absent()
          : Value(episodeTitle),
      watchedAt: Value(watchedAt),
      syncedAt: syncedAt == null && nullToAbsent
          ? const Value.absent()
          : Value(syncedAt),
      updatedAt: Value(updatedAt),
    );
  }

  factory HistoryData.fromJson(Map<String, dynamic> json,
      {ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return HistoryData(
      id: serializer.fromJson<int>(json['id']),
      animeId: serializer.fromJson<int>(json['animeId']),
      episodeNumber: serializer.fromJson<int>(json['episodeNumber']),
      episodeTitle: serializer.fromJson<String?>(json['episodeTitle']),
      watchedAt: serializer.fromJson<DateTime>(json['watchedAt']),
      syncedAt: serializer.fromJson<DateTime?>(json['syncedAt']),
      updatedAt: serializer.fromJson<DateTime>(json['updatedAt']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'animeId': serializer.toJson<int>(animeId),
      'episodeNumber': serializer.toJson<int>(episodeNumber),
      'episodeTitle': serializer.toJson<String?>(episodeTitle),
      'watchedAt': serializer.toJson<DateTime>(watchedAt),
      'syncedAt': serializer.toJson<DateTime?>(syncedAt),
      'updatedAt': serializer.toJson<DateTime>(updatedAt),
    };
  }

  HistoryData copyWith(
          {int? id,
          int? animeId,
          int? episodeNumber,
          Value<String?> episodeTitle = const Value.absent(),
          DateTime? watchedAt,
          Value<DateTime?> syncedAt = const Value.absent(),
          DateTime? updatedAt}) =>
      HistoryData(
        id: id ?? this.id,
        animeId: animeId ?? this.animeId,
        episodeNumber: episodeNumber ?? this.episodeNumber,
        episodeTitle:
            episodeTitle.present ? episodeTitle.value : this.episodeTitle,
        watchedAt: watchedAt ?? this.watchedAt,
        syncedAt: syncedAt.present ? syncedAt.value : this.syncedAt,
        updatedAt: updatedAt ?? this.updatedAt,
      );
  HistoryData copyWithCompanion(HistoryCompanion data) {
    return HistoryData(
      id: data.id.present ? data.id.value : this.id,
      animeId: data.animeId.present ? data.animeId.value : this.animeId,
      episodeNumber: data.episodeNumber.present
          ? data.episodeNumber.value
          : this.episodeNumber,
      episodeTitle: data.episodeTitle.present
          ? data.episodeTitle.value
          : this.episodeTitle,
      watchedAt: data.watchedAt.present ? data.watchedAt.value : this.watchedAt,
      syncedAt: data.syncedAt.present ? data.syncedAt.value : this.syncedAt,
      updatedAt: data.updatedAt.present ? data.updatedAt.value : this.updatedAt,
    );
  }

  @override
  String toString() {
    return (StringBuffer('HistoryData(')
          ..write('id: $id, ')
          ..write('animeId: $animeId, ')
          ..write('episodeNumber: $episodeNumber, ')
          ..write('episodeTitle: $episodeTitle, ')
          ..write('watchedAt: $watchedAt, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(
      id, animeId, episodeNumber, episodeTitle, watchedAt, syncedAt, updatedAt);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is HistoryData &&
          other.id == this.id &&
          other.animeId == this.animeId &&
          other.episodeNumber == this.episodeNumber &&
          other.episodeTitle == this.episodeTitle &&
          other.watchedAt == this.watchedAt &&
          other.syncedAt == this.syncedAt &&
          other.updatedAt == this.updatedAt);
}

class HistoryCompanion extends UpdateCompanion<HistoryData> {
  final Value<int> id;
  final Value<int> animeId;
  final Value<int> episodeNumber;
  final Value<String?> episodeTitle;
  final Value<DateTime> watchedAt;
  final Value<DateTime?> syncedAt;
  final Value<DateTime> updatedAt;
  const HistoryCompanion({
    this.id = const Value.absent(),
    this.animeId = const Value.absent(),
    this.episodeNumber = const Value.absent(),
    this.episodeTitle = const Value.absent(),
    this.watchedAt = const Value.absent(),
    this.syncedAt = const Value.absent(),
    this.updatedAt = const Value.absent(),
  });
  HistoryCompanion.insert({
    this.id = const Value.absent(),
    required int animeId,
    required int episodeNumber,
    this.episodeTitle = const Value.absent(),
    required DateTime watchedAt,
    this.syncedAt = const Value.absent(),
    required DateTime updatedAt,
  })  : animeId = Value(animeId),
        episodeNumber = Value(episodeNumber),
        watchedAt = Value(watchedAt),
        updatedAt = Value(updatedAt);
  static Insertable<HistoryData> custom({
    Expression<int>? id,
    Expression<int>? animeId,
    Expression<int>? episodeNumber,
    Expression<String>? episodeTitle,
    Expression<DateTime>? watchedAt,
    Expression<DateTime>? syncedAt,
    Expression<DateTime>? updatedAt,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (animeId != null) 'anime_id': animeId,
      if (episodeNumber != null) 'episode_number': episodeNumber,
      if (episodeTitle != null) 'episode_title': episodeTitle,
      if (watchedAt != null) 'watched_at': watchedAt,
      if (syncedAt != null) 'synced_at': syncedAt,
      if (updatedAt != null) 'updated_at': updatedAt,
    });
  }

  HistoryCompanion copyWith(
      {Value<int>? id,
      Value<int>? animeId,
      Value<int>? episodeNumber,
      Value<String?>? episodeTitle,
      Value<DateTime>? watchedAt,
      Value<DateTime?>? syncedAt,
      Value<DateTime>? updatedAt}) {
    return HistoryCompanion(
      id: id ?? this.id,
      animeId: animeId ?? this.animeId,
      episodeNumber: episodeNumber ?? this.episodeNumber,
      episodeTitle: episodeTitle ?? this.episodeTitle,
      watchedAt: watchedAt ?? this.watchedAt,
      syncedAt: syncedAt ?? this.syncedAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (animeId.present) {
      map['anime_id'] = Variable<int>(animeId.value);
    }
    if (episodeNumber.present) {
      map['episode_number'] = Variable<int>(episodeNumber.value);
    }
    if (episodeTitle.present) {
      map['episode_title'] = Variable<String>(episodeTitle.value);
    }
    if (watchedAt.present) {
      map['watched_at'] = Variable<DateTime>(watchedAt.value);
    }
    if (syncedAt.present) {
      map['synced_at'] = Variable<DateTime>(syncedAt.value);
    }
    if (updatedAt.present) {
      map['updated_at'] = Variable<DateTime>(updatedAt.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('HistoryCompanion(')
          ..write('id: $id, ')
          ..write('animeId: $animeId, ')
          ..write('episodeNumber: $episodeNumber, ')
          ..write('episodeTitle: $episodeTitle, ')
          ..write('watchedAt: $watchedAt, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }
}

class $SyncQueueTable extends SyncQueue
    with TableInfo<$SyncQueueTable, SyncQueueData> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $SyncQueueTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
      'id', aliasedName, false,
      hasAutoIncrement: true,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('PRIMARY KEY AUTOINCREMENT'));
  static const VerificationMeta _operationMeta =
      const VerificationMeta('operation');
  @override
  late final GeneratedColumn<String> operation = GeneratedColumn<String>(
      'operation', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _entityTypeMeta =
      const VerificationMeta('entityType');
  @override
  late final GeneratedColumn<String> entityType = GeneratedColumn<String>(
      'entity_type', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _entityIdMeta =
      const VerificationMeta('entityId');
  @override
  late final GeneratedColumn<String> entityId = GeneratedColumn<String>(
      'entity_id', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _dataMeta = const VerificationMeta('data');
  @override
  late final GeneratedColumn<String> data = GeneratedColumn<String>(
      'data', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _createdAtMeta =
      const VerificationMeta('createdAt');
  @override
  late final GeneratedColumn<DateTime> createdAt = GeneratedColumn<DateTime>(
      'created_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _syncedMeta = const VerificationMeta('synced');
  @override
  late final GeneratedColumn<bool> synced = GeneratedColumn<bool>(
      'synced', aliasedName, false,
      type: DriftSqlType.bool,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('CHECK ("synced" IN (0, 1))'),
      defaultValue: const Constant(false));
  static const VerificationMeta _retryCountMeta =
      const VerificationMeta('retryCount');
  @override
  late final GeneratedColumn<int> retryCount = GeneratedColumn<int>(
      'retry_count', aliasedName, false,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultValue: const Constant(0));
  static const VerificationMeta _lastErrorMeta =
      const VerificationMeta('lastError');
  @override
  late final GeneratedColumn<String> lastError = GeneratedColumn<String>(
      'last_error', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  @override
  List<GeneratedColumn> get $columns => [
        id,
        operation,
        entityType,
        entityId,
        data,
        createdAt,
        synced,
        retryCount,
        lastError
      ];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'sync_queue';
  @override
  VerificationContext validateIntegrity(Insertable<SyncQueueData> instance,
      {bool isInserting = false}) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('operation')) {
      context.handle(_operationMeta,
          operation.isAcceptableOrUnknown(data['operation']!, _operationMeta));
    } else if (isInserting) {
      context.missing(_operationMeta);
    }
    if (data.containsKey('entity_type')) {
      context.handle(
          _entityTypeMeta,
          entityType.isAcceptableOrUnknown(
              data['entity_type']!, _entityTypeMeta));
    } else if (isInserting) {
      context.missing(_entityTypeMeta);
    }
    if (data.containsKey('entity_id')) {
      context.handle(_entityIdMeta,
          entityId.isAcceptableOrUnknown(data['entity_id']!, _entityIdMeta));
    } else if (isInserting) {
      context.missing(_entityIdMeta);
    }
    if (data.containsKey('data')) {
      context.handle(
          _dataMeta, this.data.isAcceptableOrUnknown(data['data']!, _dataMeta));
    } else if (isInserting) {
      context.missing(_dataMeta);
    }
    if (data.containsKey('created_at')) {
      context.handle(_createdAtMeta,
          createdAt.isAcceptableOrUnknown(data['created_at']!, _createdAtMeta));
    } else if (isInserting) {
      context.missing(_createdAtMeta);
    }
    if (data.containsKey('synced')) {
      context.handle(_syncedMeta,
          synced.isAcceptableOrUnknown(data['synced']!, _syncedMeta));
    }
    if (data.containsKey('retry_count')) {
      context.handle(
          _retryCountMeta,
          retryCount.isAcceptableOrUnknown(
              data['retry_count']!, _retryCountMeta));
    }
    if (data.containsKey('last_error')) {
      context.handle(_lastErrorMeta,
          lastError.isAcceptableOrUnknown(data['last_error']!, _lastErrorMeta));
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  SyncQueueData map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return SyncQueueData(
      id: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}id'])!,
      operation: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}operation'])!,
      entityType: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}entity_type'])!,
      entityId: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}entity_id'])!,
      data: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}data'])!,
      createdAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}created_at'])!,
      synced: attachedDatabase.typeMapping
          .read(DriftSqlType.bool, data['${effectivePrefix}synced'])!,
      retryCount: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}retry_count'])!,
      lastError: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}last_error']),
    );
  }

  @override
  $SyncQueueTable createAlias(String alias) {
    return $SyncQueueTable(attachedDatabase, alias);
  }
}

class SyncQueueData extends DataClass implements Insertable<SyncQueueData> {
  final int id;
  final String operation;
  final String entityType;
  final String entityId;
  final String data;
  final DateTime createdAt;
  final bool synced;
  final int retryCount;
  final String? lastError;
  const SyncQueueData(
      {required this.id,
      required this.operation,
      required this.entityType,
      required this.entityId,
      required this.data,
      required this.createdAt,
      required this.synced,
      required this.retryCount,
      this.lastError});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['operation'] = Variable<String>(operation);
    map['entity_type'] = Variable<String>(entityType);
    map['entity_id'] = Variable<String>(entityId);
    map['data'] = Variable<String>(data);
    map['created_at'] = Variable<DateTime>(createdAt);
    map['synced'] = Variable<bool>(synced);
    map['retry_count'] = Variable<int>(retryCount);
    if (!nullToAbsent || lastError != null) {
      map['last_error'] = Variable<String>(lastError);
    }
    return map;
  }

  SyncQueueCompanion toCompanion(bool nullToAbsent) {
    return SyncQueueCompanion(
      id: Value(id),
      operation: Value(operation),
      entityType: Value(entityType),
      entityId: Value(entityId),
      data: Value(data),
      createdAt: Value(createdAt),
      synced: Value(synced),
      retryCount: Value(retryCount),
      lastError: lastError == null && nullToAbsent
          ? const Value.absent()
          : Value(lastError),
    );
  }

  factory SyncQueueData.fromJson(Map<String, dynamic> json,
      {ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return SyncQueueData(
      id: serializer.fromJson<int>(json['id']),
      operation: serializer.fromJson<String>(json['operation']),
      entityType: serializer.fromJson<String>(json['entityType']),
      entityId: serializer.fromJson<String>(json['entityId']),
      data: serializer.fromJson<String>(json['data']),
      createdAt: serializer.fromJson<DateTime>(json['createdAt']),
      synced: serializer.fromJson<bool>(json['synced']),
      retryCount: serializer.fromJson<int>(json['retryCount']),
      lastError: serializer.fromJson<String?>(json['lastError']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'operation': serializer.toJson<String>(operation),
      'entityType': serializer.toJson<String>(entityType),
      'entityId': serializer.toJson<String>(entityId),
      'data': serializer.toJson<String>(data),
      'createdAt': serializer.toJson<DateTime>(createdAt),
      'synced': serializer.toJson<bool>(synced),
      'retryCount': serializer.toJson<int>(retryCount),
      'lastError': serializer.toJson<String?>(lastError),
    };
  }

  SyncQueueData copyWith(
          {int? id,
          String? operation,
          String? entityType,
          String? entityId,
          String? data,
          DateTime? createdAt,
          bool? synced,
          int? retryCount,
          Value<String?> lastError = const Value.absent()}) =>
      SyncQueueData(
        id: id ?? this.id,
        operation: operation ?? this.operation,
        entityType: entityType ?? this.entityType,
        entityId: entityId ?? this.entityId,
        data: data ?? this.data,
        createdAt: createdAt ?? this.createdAt,
        synced: synced ?? this.synced,
        retryCount: retryCount ?? this.retryCount,
        lastError: lastError.present ? lastError.value : this.lastError,
      );
  SyncQueueData copyWithCompanion(SyncQueueCompanion data) {
    return SyncQueueData(
      id: data.id.present ? data.id.value : this.id,
      operation: data.operation.present ? data.operation.value : this.operation,
      entityType:
          data.entityType.present ? data.entityType.value : this.entityType,
      entityId: data.entityId.present ? data.entityId.value : this.entityId,
      data: data.data.present ? data.data.value : this.data,
      createdAt: data.createdAt.present ? data.createdAt.value : this.createdAt,
      synced: data.synced.present ? data.synced.value : this.synced,
      retryCount:
          data.retryCount.present ? data.retryCount.value : this.retryCount,
      lastError: data.lastError.present ? data.lastError.value : this.lastError,
    );
  }

  @override
  String toString() {
    return (StringBuffer('SyncQueueData(')
          ..write('id: $id, ')
          ..write('operation: $operation, ')
          ..write('entityType: $entityType, ')
          ..write('entityId: $entityId, ')
          ..write('data: $data, ')
          ..write('createdAt: $createdAt, ')
          ..write('synced: $synced, ')
          ..write('retryCount: $retryCount, ')
          ..write('lastError: $lastError')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(id, operation, entityType, entityId, data,
      createdAt, synced, retryCount, lastError);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is SyncQueueData &&
          other.id == this.id &&
          other.operation == this.operation &&
          other.entityType == this.entityType &&
          other.entityId == this.entityId &&
          other.data == this.data &&
          other.createdAt == this.createdAt &&
          other.synced == this.synced &&
          other.retryCount == this.retryCount &&
          other.lastError == this.lastError);
}

class SyncQueueCompanion extends UpdateCompanion<SyncQueueData> {
  final Value<int> id;
  final Value<String> operation;
  final Value<String> entityType;
  final Value<String> entityId;
  final Value<String> data;
  final Value<DateTime> createdAt;
  final Value<bool> synced;
  final Value<int> retryCount;
  final Value<String?> lastError;
  const SyncQueueCompanion({
    this.id = const Value.absent(),
    this.operation = const Value.absent(),
    this.entityType = const Value.absent(),
    this.entityId = const Value.absent(),
    this.data = const Value.absent(),
    this.createdAt = const Value.absent(),
    this.synced = const Value.absent(),
    this.retryCount = const Value.absent(),
    this.lastError = const Value.absent(),
  });
  SyncQueueCompanion.insert({
    this.id = const Value.absent(),
    required String operation,
    required String entityType,
    required String entityId,
    required String data,
    required DateTime createdAt,
    this.synced = const Value.absent(),
    this.retryCount = const Value.absent(),
    this.lastError = const Value.absent(),
  })  : operation = Value(operation),
        entityType = Value(entityType),
        entityId = Value(entityId),
        data = Value(data),
        createdAt = Value(createdAt);
  static Insertable<SyncQueueData> custom({
    Expression<int>? id,
    Expression<String>? operation,
    Expression<String>? entityType,
    Expression<String>? entityId,
    Expression<String>? data,
    Expression<DateTime>? createdAt,
    Expression<bool>? synced,
    Expression<int>? retryCount,
    Expression<String>? lastError,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (operation != null) 'operation': operation,
      if (entityType != null) 'entity_type': entityType,
      if (entityId != null) 'entity_id': entityId,
      if (data != null) 'data': data,
      if (createdAt != null) 'created_at': createdAt,
      if (synced != null) 'synced': synced,
      if (retryCount != null) 'retry_count': retryCount,
      if (lastError != null) 'last_error': lastError,
    });
  }

  SyncQueueCompanion copyWith(
      {Value<int>? id,
      Value<String>? operation,
      Value<String>? entityType,
      Value<String>? entityId,
      Value<String>? data,
      Value<DateTime>? createdAt,
      Value<bool>? synced,
      Value<int>? retryCount,
      Value<String?>? lastError}) {
    return SyncQueueCompanion(
      id: id ?? this.id,
      operation: operation ?? this.operation,
      entityType: entityType ?? this.entityType,
      entityId: entityId ?? this.entityId,
      data: data ?? this.data,
      createdAt: createdAt ?? this.createdAt,
      synced: synced ?? this.synced,
      retryCount: retryCount ?? this.retryCount,
      lastError: lastError ?? this.lastError,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (operation.present) {
      map['operation'] = Variable<String>(operation.value);
    }
    if (entityType.present) {
      map['entity_type'] = Variable<String>(entityType.value);
    }
    if (entityId.present) {
      map['entity_id'] = Variable<String>(entityId.value);
    }
    if (data.present) {
      map['data'] = Variable<String>(data.value);
    }
    if (createdAt.present) {
      map['created_at'] = Variable<DateTime>(createdAt.value);
    }
    if (synced.present) {
      map['synced'] = Variable<bool>(synced.value);
    }
    if (retryCount.present) {
      map['retry_count'] = Variable<int>(retryCount.value);
    }
    if (lastError.present) {
      map['last_error'] = Variable<String>(lastError.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('SyncQueueCompanion(')
          ..write('id: $id, ')
          ..write('operation: $operation, ')
          ..write('entityType: $entityType, ')
          ..write('entityId: $entityId, ')
          ..write('data: $data, ')
          ..write('createdAt: $createdAt, ')
          ..write('synced: $synced, ')
          ..write('retryCount: $retryCount, ')
          ..write('lastError: $lastError')
          ..write(')'))
        .toString();
  }
}

class $VipStatusTable extends VipStatus
    with TableInfo<$VipStatusTable, VipStatusData> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $VipStatusTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
      'id', aliasedName, false,
      hasAutoIncrement: true,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('PRIMARY KEY AUTOINCREMENT'));
  static const VerificationMeta _isVipMeta = const VerificationMeta('isVip');
  @override
  late final GeneratedColumn<bool> isVip = GeneratedColumn<bool>(
      'is_vip', aliasedName, false,
      type: DriftSqlType.bool,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('CHECK ("is_vip" IN (0, 1))'),
      defaultValue: const Constant(false));
  static const VerificationMeta _vipExpiresAtMeta =
      const VerificationMeta('vipExpiresAt');
  @override
  late final GeneratedColumn<DateTime> vipExpiresAt = GeneratedColumn<DateTime>(
      'vip_expires_at', aliasedName, true,
      type: DriftSqlType.dateTime, requiredDuringInsert: false);
  static const VerificationMeta _vipTierMeta =
      const VerificationMeta('vipTier');
  @override
  late final GeneratedColumn<String> vipTier = GeneratedColumn<String>(
      'vip_tier', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _syncedAtMeta =
      const VerificationMeta('syncedAt');
  @override
  late final GeneratedColumn<DateTime> syncedAt = GeneratedColumn<DateTime>(
      'synced_at', aliasedName, true,
      type: DriftSqlType.dateTime, requiredDuringInsert: false);
  static const VerificationMeta _updatedAtMeta =
      const VerificationMeta('updatedAt');
  @override
  late final GeneratedColumn<DateTime> updatedAt = GeneratedColumn<DateTime>(
      'updated_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  @override
  List<GeneratedColumn> get $columns =>
      [id, isVip, vipExpiresAt, vipTier, syncedAt, updatedAt];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'vip_status';
  @override
  VerificationContext validateIntegrity(Insertable<VipStatusData> instance,
      {bool isInserting = false}) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('is_vip')) {
      context.handle(
          _isVipMeta, isVip.isAcceptableOrUnknown(data['is_vip']!, _isVipMeta));
    }
    if (data.containsKey('vip_expires_at')) {
      context.handle(
          _vipExpiresAtMeta,
          vipExpiresAt.isAcceptableOrUnknown(
              data['vip_expires_at']!, _vipExpiresAtMeta));
    }
    if (data.containsKey('vip_tier')) {
      context.handle(_vipTierMeta,
          vipTier.isAcceptableOrUnknown(data['vip_tier']!, _vipTierMeta));
    }
    if (data.containsKey('synced_at')) {
      context.handle(_syncedAtMeta,
          syncedAt.isAcceptableOrUnknown(data['synced_at']!, _syncedAtMeta));
    }
    if (data.containsKey('updated_at')) {
      context.handle(_updatedAtMeta,
          updatedAt.isAcceptableOrUnknown(data['updated_at']!, _updatedAtMeta));
    } else if (isInserting) {
      context.missing(_updatedAtMeta);
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  VipStatusData map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return VipStatusData(
      id: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}id'])!,
      isVip: attachedDatabase.typeMapping
          .read(DriftSqlType.bool, data['${effectivePrefix}is_vip'])!,
      vipExpiresAt: attachedDatabase.typeMapping.read(
          DriftSqlType.dateTime, data['${effectivePrefix}vip_expires_at']),
      vipTier: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}vip_tier']),
      syncedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}synced_at']),
      updatedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}updated_at'])!,
    );
  }

  @override
  $VipStatusTable createAlias(String alias) {
    return $VipStatusTable(attachedDatabase, alias);
  }
}

class VipStatusData extends DataClass implements Insertable<VipStatusData> {
  final int id;
  final bool isVip;
  final DateTime? vipExpiresAt;
  final String? vipTier;
  final DateTime? syncedAt;
  final DateTime updatedAt;
  const VipStatusData(
      {required this.id,
      required this.isVip,
      this.vipExpiresAt,
      this.vipTier,
      this.syncedAt,
      required this.updatedAt});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['is_vip'] = Variable<bool>(isVip);
    if (!nullToAbsent || vipExpiresAt != null) {
      map['vip_expires_at'] = Variable<DateTime>(vipExpiresAt);
    }
    if (!nullToAbsent || vipTier != null) {
      map['vip_tier'] = Variable<String>(vipTier);
    }
    if (!nullToAbsent || syncedAt != null) {
      map['synced_at'] = Variable<DateTime>(syncedAt);
    }
    map['updated_at'] = Variable<DateTime>(updatedAt);
    return map;
  }

  VipStatusCompanion toCompanion(bool nullToAbsent) {
    return VipStatusCompanion(
      id: Value(id),
      isVip: Value(isVip),
      vipExpiresAt: vipExpiresAt == null && nullToAbsent
          ? const Value.absent()
          : Value(vipExpiresAt),
      vipTier: vipTier == null && nullToAbsent
          ? const Value.absent()
          : Value(vipTier),
      syncedAt: syncedAt == null && nullToAbsent
          ? const Value.absent()
          : Value(syncedAt),
      updatedAt: Value(updatedAt),
    );
  }

  factory VipStatusData.fromJson(Map<String, dynamic> json,
      {ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return VipStatusData(
      id: serializer.fromJson<int>(json['id']),
      isVip: serializer.fromJson<bool>(json['isVip']),
      vipExpiresAt: serializer.fromJson<DateTime?>(json['vipExpiresAt']),
      vipTier: serializer.fromJson<String?>(json['vipTier']),
      syncedAt: serializer.fromJson<DateTime?>(json['syncedAt']),
      updatedAt: serializer.fromJson<DateTime>(json['updatedAt']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'isVip': serializer.toJson<bool>(isVip),
      'vipExpiresAt': serializer.toJson<DateTime?>(vipExpiresAt),
      'vipTier': serializer.toJson<String?>(vipTier),
      'syncedAt': serializer.toJson<DateTime?>(syncedAt),
      'updatedAt': serializer.toJson<DateTime>(updatedAt),
    };
  }

  VipStatusData copyWith(
          {int? id,
          bool? isVip,
          Value<DateTime?> vipExpiresAt = const Value.absent(),
          Value<String?> vipTier = const Value.absent(),
          Value<DateTime?> syncedAt = const Value.absent(),
          DateTime? updatedAt}) =>
      VipStatusData(
        id: id ?? this.id,
        isVip: isVip ?? this.isVip,
        vipExpiresAt:
            vipExpiresAt.present ? vipExpiresAt.value : this.vipExpiresAt,
        vipTier: vipTier.present ? vipTier.value : this.vipTier,
        syncedAt: syncedAt.present ? syncedAt.value : this.syncedAt,
        updatedAt: updatedAt ?? this.updatedAt,
      );
  VipStatusData copyWithCompanion(VipStatusCompanion data) {
    return VipStatusData(
      id: data.id.present ? data.id.value : this.id,
      isVip: data.isVip.present ? data.isVip.value : this.isVip,
      vipExpiresAt: data.vipExpiresAt.present
          ? data.vipExpiresAt.value
          : this.vipExpiresAt,
      vipTier: data.vipTier.present ? data.vipTier.value : this.vipTier,
      syncedAt: data.syncedAt.present ? data.syncedAt.value : this.syncedAt,
      updatedAt: data.updatedAt.present ? data.updatedAt.value : this.updatedAt,
    );
  }

  @override
  String toString() {
    return (StringBuffer('VipStatusData(')
          ..write('id: $id, ')
          ..write('isVip: $isVip, ')
          ..write('vipExpiresAt: $vipExpiresAt, ')
          ..write('vipTier: $vipTier, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode =>
      Object.hash(id, isVip, vipExpiresAt, vipTier, syncedAt, updatedAt);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is VipStatusData &&
          other.id == this.id &&
          other.isVip == this.isVip &&
          other.vipExpiresAt == this.vipExpiresAt &&
          other.vipTier == this.vipTier &&
          other.syncedAt == this.syncedAt &&
          other.updatedAt == this.updatedAt);
}

class VipStatusCompanion extends UpdateCompanion<VipStatusData> {
  final Value<int> id;
  final Value<bool> isVip;
  final Value<DateTime?> vipExpiresAt;
  final Value<String?> vipTier;
  final Value<DateTime?> syncedAt;
  final Value<DateTime> updatedAt;
  const VipStatusCompanion({
    this.id = const Value.absent(),
    this.isVip = const Value.absent(),
    this.vipExpiresAt = const Value.absent(),
    this.vipTier = const Value.absent(),
    this.syncedAt = const Value.absent(),
    this.updatedAt = const Value.absent(),
  });
  VipStatusCompanion.insert({
    this.id = const Value.absent(),
    this.isVip = const Value.absent(),
    this.vipExpiresAt = const Value.absent(),
    this.vipTier = const Value.absent(),
    this.syncedAt = const Value.absent(),
    required DateTime updatedAt,
  }) : updatedAt = Value(updatedAt);
  static Insertable<VipStatusData> custom({
    Expression<int>? id,
    Expression<bool>? isVip,
    Expression<DateTime>? vipExpiresAt,
    Expression<String>? vipTier,
    Expression<DateTime>? syncedAt,
    Expression<DateTime>? updatedAt,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (isVip != null) 'is_vip': isVip,
      if (vipExpiresAt != null) 'vip_expires_at': vipExpiresAt,
      if (vipTier != null) 'vip_tier': vipTier,
      if (syncedAt != null) 'synced_at': syncedAt,
      if (updatedAt != null) 'updated_at': updatedAt,
    });
  }

  VipStatusCompanion copyWith(
      {Value<int>? id,
      Value<bool>? isVip,
      Value<DateTime?>? vipExpiresAt,
      Value<String?>? vipTier,
      Value<DateTime?>? syncedAt,
      Value<DateTime>? updatedAt}) {
    return VipStatusCompanion(
      id: id ?? this.id,
      isVip: isVip ?? this.isVip,
      vipExpiresAt: vipExpiresAt ?? this.vipExpiresAt,
      vipTier: vipTier ?? this.vipTier,
      syncedAt: syncedAt ?? this.syncedAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (isVip.present) {
      map['is_vip'] = Variable<bool>(isVip.value);
    }
    if (vipExpiresAt.present) {
      map['vip_expires_at'] = Variable<DateTime>(vipExpiresAt.value);
    }
    if (vipTier.present) {
      map['vip_tier'] = Variable<String>(vipTier.value);
    }
    if (syncedAt.present) {
      map['synced_at'] = Variable<DateTime>(syncedAt.value);
    }
    if (updatedAt.present) {
      map['updated_at'] = Variable<DateTime>(updatedAt.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('VipStatusCompanion(')
          ..write('id: $id, ')
          ..write('isVip: $isVip, ')
          ..write('vipExpiresAt: $vipExpiresAt, ')
          ..write('vipTier: $vipTier, ')
          ..write('syncedAt: $syncedAt, ')
          ..write('updatedAt: $updatedAt')
          ..write(')'))
        .toString();
  }
}

class $SyncConflictsTable extends SyncConflicts
    with TableInfo<$SyncConflictsTable, SyncConflictData> {
  @override
  final GeneratedDatabase attachedDatabase;
  final String? _alias;
  $SyncConflictsTable(this.attachedDatabase, [this._alias]);
  static const VerificationMeta _idMeta = const VerificationMeta('id');
  @override
  late final GeneratedColumn<int> id = GeneratedColumn<int>(
      'id', aliasedName, false,
      hasAutoIncrement: true,
      type: DriftSqlType.int,
      requiredDuringInsert: false,
      defaultConstraints:
          GeneratedColumn.constraintIsAlways('PRIMARY KEY AUTOINCREMENT'));
  static const VerificationMeta _entityTypeMeta =
      const VerificationMeta('entityType');
  @override
  late final GeneratedColumn<String> entityType = GeneratedColumn<String>(
      'entity_type', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _entityIdMeta =
      const VerificationMeta('entityId');
  @override
  late final GeneratedColumn<String> entityId = GeneratedColumn<String>(
      'entity_id', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _mobileDataMeta =
      const VerificationMeta('mobileData');
  @override
  late final GeneratedColumn<String> mobileData = GeneratedColumn<String>(
      'mobile_data', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _desktopDataMeta =
      const VerificationMeta('desktopData');
  @override
  late final GeneratedColumn<String> desktopData = GeneratedColumn<String>(
      'desktop_data', aliasedName, false,
      type: DriftSqlType.string, requiredDuringInsert: true);
  static const VerificationMeta _mobileTimestampMeta =
      const VerificationMeta('mobileTimestamp');
  @override
  late final GeneratedColumn<DateTime> mobileTimestamp =
      GeneratedColumn<DateTime>('mobile_timestamp', aliasedName, false,
          type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _desktopTimestampMeta =
      const VerificationMeta('desktopTimestamp');
  @override
  late final GeneratedColumn<DateTime> desktopTimestamp =
      GeneratedColumn<DateTime>('desktop_timestamp', aliasedName, false,
          type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _resolutionMeta =
      const VerificationMeta('resolution');
  @override
  late final GeneratedColumn<String> resolution = GeneratedColumn<String>(
      'resolution', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _resolvedDataMeta =
      const VerificationMeta('resolvedData');
  @override
  late final GeneratedColumn<String> resolvedData = GeneratedColumn<String>(
      'resolved_data', aliasedName, true,
      type: DriftSqlType.string, requiredDuringInsert: false);
  static const VerificationMeta _createdAtMeta =
      const VerificationMeta('createdAt');
  @override
  late final GeneratedColumn<DateTime> createdAt = GeneratedColumn<DateTime>(
      'created_at', aliasedName, false,
      type: DriftSqlType.dateTime, requiredDuringInsert: true);
  static const VerificationMeta _resolvedAtMeta =
      const VerificationMeta('resolvedAt');
  @override
  late final GeneratedColumn<DateTime> resolvedAt = GeneratedColumn<DateTime>(
      'resolved_at', aliasedName, true,
      type: DriftSqlType.dateTime, requiredDuringInsert: false);
  @override
  List<GeneratedColumn> get $columns => [
        id,
        entityType,
        entityId,
        mobileData,
        desktopData,
        mobileTimestamp,
        desktopTimestamp,
        resolution,
        resolvedData,
        createdAt,
        resolvedAt
      ];
  @override
  String get aliasedName => _alias ?? actualTableName;
  @override
  String get actualTableName => $name;
  static const String $name = 'sync_conflicts';
  @override
  VerificationContext validateIntegrity(Insertable<SyncConflictData> instance,
      {bool isInserting = false}) {
    final context = VerificationContext();
    final data = instance.toColumns(true);
    if (data.containsKey('id')) {
      context.handle(_idMeta, id.isAcceptableOrUnknown(data['id']!, _idMeta));
    }
    if (data.containsKey('entity_type')) {
      context.handle(
          _entityTypeMeta,
          entityType.isAcceptableOrUnknown(
              data['entity_type']!, _entityTypeMeta));
    } else if (isInserting) {
      context.missing(_entityTypeMeta);
    }
    if (data.containsKey('entity_id')) {
      context.handle(_entityIdMeta,
          entityId.isAcceptableOrUnknown(data['entity_id']!, _entityIdMeta));
    } else if (isInserting) {
      context.missing(_entityIdMeta);
    }
    if (data.containsKey('mobile_data')) {
      context.handle(
          _mobileDataMeta,
          mobileData.isAcceptableOrUnknown(
              data['mobile_data']!, _mobileDataMeta));
    } else if (isInserting) {
      context.missing(_mobileDataMeta);
    }
    if (data.containsKey('desktop_data')) {
      context.handle(
          _desktopDataMeta,
          desktopData.isAcceptableOrUnknown(
              data['desktop_data']!, _desktopDataMeta));
    } else if (isInserting) {
      context.missing(_desktopDataMeta);
    }
    if (data.containsKey('mobile_timestamp')) {
      context.handle(
          _mobileTimestampMeta,
          mobileTimestamp.isAcceptableOrUnknown(
              data['mobile_timestamp']!, _mobileTimestampMeta));
    } else if (isInserting) {
      context.missing(_mobileTimestampMeta);
    }
    if (data.containsKey('desktop_timestamp')) {
      context.handle(
          _desktopTimestampMeta,
          desktopTimestamp.isAcceptableOrUnknown(
              data['desktop_timestamp']!, _desktopTimestampMeta));
    } else if (isInserting) {
      context.missing(_desktopTimestampMeta);
    }
    if (data.containsKey('resolution')) {
      context.handle(
          _resolutionMeta,
          resolution.isAcceptableOrUnknown(
              data['resolution']!, _resolutionMeta));
    }
    if (data.containsKey('resolved_data')) {
      context.handle(
          _resolvedDataMeta,
          resolvedData.isAcceptableOrUnknown(
              data['resolved_data']!, _resolvedDataMeta));
    }
    if (data.containsKey('created_at')) {
      context.handle(_createdAtMeta,
          createdAt.isAcceptableOrUnknown(data['created_at']!, _createdAtMeta));
    } else if (isInserting) {
      context.missing(_createdAtMeta);
    }
    if (data.containsKey('resolved_at')) {
      context.handle(
          _resolvedAtMeta,
          resolvedAt.isAcceptableOrUnknown(
              data['resolved_at']!, _resolvedAtMeta));
    }
    return context;
  }

  @override
  Set<GeneratedColumn> get $primaryKey => {id};
  @override
  SyncConflictData map(Map<String, dynamic> data, {String? tablePrefix}) {
    final effectivePrefix = tablePrefix != null ? '$tablePrefix.' : '';
    return SyncConflictData(
      id: attachedDatabase.typeMapping
          .read(DriftSqlType.int, data['${effectivePrefix}id'])!,
      entityType: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}entity_type'])!,
      entityId: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}entity_id'])!,
      mobileData: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}mobile_data'])!,
      desktopData: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}desktop_data'])!,
      mobileTimestamp: attachedDatabase.typeMapping.read(
          DriftSqlType.dateTime, data['${effectivePrefix}mobile_timestamp'])!,
      desktopTimestamp: attachedDatabase.typeMapping.read(
          DriftSqlType.dateTime, data['${effectivePrefix}desktop_timestamp'])!,
      resolution: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}resolution']),
      resolvedData: attachedDatabase.typeMapping
          .read(DriftSqlType.string, data['${effectivePrefix}resolved_data']),
      createdAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}created_at'])!,
      resolvedAt: attachedDatabase.typeMapping
          .read(DriftSqlType.dateTime, data['${effectivePrefix}resolved_at']),
    );
  }

  @override
  $SyncConflictsTable createAlias(String alias) {
    return $SyncConflictsTable(attachedDatabase, alias);
  }
}

class SyncConflictData extends DataClass
    implements Insertable<SyncConflictData> {
  final int id;
  final String entityType;
  final String entityId;
  final String mobileData;
  final String desktopData;
  final DateTime mobileTimestamp;
  final DateTime desktopTimestamp;
  final String? resolution;
  final String? resolvedData;
  final DateTime createdAt;
  final DateTime? resolvedAt;
  const SyncConflictData(
      {required this.id,
      required this.entityType,
      required this.entityId,
      required this.mobileData,
      required this.desktopData,
      required this.mobileTimestamp,
      required this.desktopTimestamp,
      this.resolution,
      this.resolvedData,
      required this.createdAt,
      this.resolvedAt});
  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    map['id'] = Variable<int>(id);
    map['entity_type'] = Variable<String>(entityType);
    map['entity_id'] = Variable<String>(entityId);
    map['mobile_data'] = Variable<String>(mobileData);
    map['desktop_data'] = Variable<String>(desktopData);
    map['mobile_timestamp'] = Variable<DateTime>(mobileTimestamp);
    map['desktop_timestamp'] = Variable<DateTime>(desktopTimestamp);
    if (!nullToAbsent || resolution != null) {
      map['resolution'] = Variable<String>(resolution);
    }
    if (!nullToAbsent || resolvedData != null) {
      map['resolved_data'] = Variable<String>(resolvedData);
    }
    map['created_at'] = Variable<DateTime>(createdAt);
    if (!nullToAbsent || resolvedAt != null) {
      map['resolved_at'] = Variable<DateTime>(resolvedAt);
    }
    return map;
  }

  SyncConflictsCompanion toCompanion(bool nullToAbsent) {
    return SyncConflictsCompanion(
      id: Value(id),
      entityType: Value(entityType),
      entityId: Value(entityId),
      mobileData: Value(mobileData),
      desktopData: Value(desktopData),
      mobileTimestamp: Value(mobileTimestamp),
      desktopTimestamp: Value(desktopTimestamp),
      resolution: resolution == null && nullToAbsent
          ? const Value.absent()
          : Value(resolution),
      resolvedData: resolvedData == null && nullToAbsent
          ? const Value.absent()
          : Value(resolvedData),
      createdAt: Value(createdAt),
      resolvedAt: resolvedAt == null && nullToAbsent
          ? const Value.absent()
          : Value(resolvedAt),
    );
  }

  factory SyncConflictData.fromJson(Map<String, dynamic> json,
      {ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return SyncConflictData(
      id: serializer.fromJson<int>(json['id']),
      entityType: serializer.fromJson<String>(json['entityType']),
      entityId: serializer.fromJson<String>(json['entityId']),
      mobileData: serializer.fromJson<String>(json['mobileData']),
      desktopData: serializer.fromJson<String>(json['desktopData']),
      mobileTimestamp: serializer.fromJson<DateTime>(json['mobileTimestamp']),
      desktopTimestamp: serializer.fromJson<DateTime>(json['desktopTimestamp']),
      resolution: serializer.fromJson<String?>(json['resolution']),
      resolvedData: serializer.fromJson<String?>(json['resolvedData']),
      createdAt: serializer.fromJson<DateTime>(json['createdAt']),
      resolvedAt: serializer.fromJson<DateTime?>(json['resolvedAt']),
    );
  }
  @override
  Map<String, dynamic> toJson({ValueSerializer? serializer}) {
    serializer ??= driftRuntimeOptions.defaultSerializer;
    return <String, dynamic>{
      'id': serializer.toJson<int>(id),
      'entityType': serializer.toJson<String>(entityType),
      'entityId': serializer.toJson<String>(entityId),
      'mobileData': serializer.toJson<String>(mobileData),
      'desktopData': serializer.toJson<String>(desktopData),
      'mobileTimestamp': serializer.toJson<DateTime>(mobileTimestamp),
      'desktopTimestamp': serializer.toJson<DateTime>(desktopTimestamp),
      'resolution': serializer.toJson<String?>(resolution),
      'resolvedData': serializer.toJson<String?>(resolvedData),
      'createdAt': serializer.toJson<DateTime>(createdAt),
      'resolvedAt': serializer.toJson<DateTime?>(resolvedAt),
    };
  }

  SyncConflictData copyWith(
          {int? id,
          String? entityType,
          String? entityId,
          String? mobileData,
          String? desktopData,
          DateTime? mobileTimestamp,
          DateTime? desktopTimestamp,
          Value<String?> resolution = const Value.absent(),
          Value<String?> resolvedData = const Value.absent(),
          DateTime? createdAt,
          Value<DateTime?> resolvedAt = const Value.absent()}) =>
      SyncConflictData(
        id: id ?? this.id,
        entityType: entityType ?? this.entityType,
        entityId: entityId ?? this.entityId,
        mobileData: mobileData ?? this.mobileData,
        desktopData: desktopData ?? this.desktopData,
        mobileTimestamp: mobileTimestamp ?? this.mobileTimestamp,
        desktopTimestamp: desktopTimestamp ?? this.desktopTimestamp,
        resolution: resolution.present ? resolution.value : this.resolution,
        resolvedData:
            resolvedData.present ? resolvedData.value : this.resolvedData,
        createdAt: createdAt ?? this.createdAt,
        resolvedAt: resolvedAt.present ? resolvedAt.value : this.resolvedAt,
      );
  SyncConflictData copyWithCompanion(SyncConflictsCompanion data) {
    return SyncConflictData(
      id: data.id.present ? data.id.value : this.id,
      entityType:
          data.entityType.present ? data.entityType.value : this.entityType,
      entityId: data.entityId.present ? data.entityId.value : this.entityId,
      mobileData:
          data.mobileData.present ? data.mobileData.value : this.mobileData,
      desktopData:
          data.desktopData.present ? data.desktopData.value : this.desktopData,
      mobileTimestamp: data.mobileTimestamp.present
          ? data.mobileTimestamp.value
          : this.mobileTimestamp,
      desktopTimestamp: data.desktopTimestamp.present
          ? data.desktopTimestamp.value
          : this.desktopTimestamp,
      resolution:
          data.resolution.present ? data.resolution.value : this.resolution,
      resolvedData: data.resolvedData.present
          ? data.resolvedData.value
          : this.resolvedData,
      createdAt: data.createdAt.present ? data.createdAt.value : this.createdAt,
      resolvedAt:
          data.resolvedAt.present ? data.resolvedAt.value : this.resolvedAt,
    );
  }

  @override
  String toString() {
    return (StringBuffer('SyncConflictData(')
          ..write('id: $id, ')
          ..write('entityType: $entityType, ')
          ..write('entityId: $entityId, ')
          ..write('mobileData: $mobileData, ')
          ..write('desktopData: $desktopData, ')
          ..write('mobileTimestamp: $mobileTimestamp, ')
          ..write('desktopTimestamp: $desktopTimestamp, ')
          ..write('resolution: $resolution, ')
          ..write('resolvedData: $resolvedData, ')
          ..write('createdAt: $createdAt, ')
          ..write('resolvedAt: $resolvedAt')
          ..write(')'))
        .toString();
  }

  @override
  int get hashCode => Object.hash(
      id,
      entityType,
      entityId,
      mobileData,
      desktopData,
      mobileTimestamp,
      desktopTimestamp,
      resolution,
      resolvedData,
      createdAt,
      resolvedAt);
  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      (other is SyncConflictData &&
          other.id == this.id &&
          other.entityType == this.entityType &&
          other.entityId == this.entityId &&
          other.mobileData == this.mobileData &&
          other.desktopData == this.desktopData &&
          other.mobileTimestamp == this.mobileTimestamp &&
          other.desktopTimestamp == this.desktopTimestamp &&
          other.resolution == this.resolution &&
          other.resolvedData == this.resolvedData &&
          other.createdAt == this.createdAt &&
          other.resolvedAt == this.resolvedAt);
}

class SyncConflictsCompanion extends UpdateCompanion<SyncConflictData> {
  final Value<int> id;
  final Value<String> entityType;
  final Value<String> entityId;
  final Value<String> mobileData;
  final Value<String> desktopData;
  final Value<DateTime> mobileTimestamp;
  final Value<DateTime> desktopTimestamp;
  final Value<String?> resolution;
  final Value<String?> resolvedData;
  final Value<DateTime> createdAt;
  final Value<DateTime?> resolvedAt;
  const SyncConflictsCompanion({
    this.id = const Value.absent(),
    this.entityType = const Value.absent(),
    this.entityId = const Value.absent(),
    this.mobileData = const Value.absent(),
    this.desktopData = const Value.absent(),
    this.mobileTimestamp = const Value.absent(),
    this.desktopTimestamp = const Value.absent(),
    this.resolution = const Value.absent(),
    this.resolvedData = const Value.absent(),
    this.createdAt = const Value.absent(),
    this.resolvedAt = const Value.absent(),
  });
  SyncConflictsCompanion.insert({
    this.id = const Value.absent(),
    required String entityType,
    required String entityId,
    required String mobileData,
    required String desktopData,
    required DateTime mobileTimestamp,
    required DateTime desktopTimestamp,
    this.resolution = const Value.absent(),
    this.resolvedData = const Value.absent(),
    required DateTime createdAt,
    this.resolvedAt = const Value.absent(),
  })  : entityType = Value(entityType),
        entityId = Value(entityId),
        mobileData = Value(mobileData),
        desktopData = Value(desktopData),
        mobileTimestamp = Value(mobileTimestamp),
        desktopTimestamp = Value(desktopTimestamp),
        createdAt = Value(createdAt);
  static Insertable<SyncConflictData> custom({
    Expression<int>? id,
    Expression<String>? entityType,
    Expression<String>? entityId,
    Expression<String>? mobileData,
    Expression<String>? desktopData,
    Expression<DateTime>? mobileTimestamp,
    Expression<DateTime>? desktopTimestamp,
    Expression<String>? resolution,
    Expression<String>? resolvedData,
    Expression<DateTime>? createdAt,
    Expression<DateTime>? resolvedAt,
  }) {
    return RawValuesInsertable({
      if (id != null) 'id': id,
      if (entityType != null) 'entity_type': entityType,
      if (entityId != null) 'entity_id': entityId,
      if (mobileData != null) 'mobile_data': mobileData,
      if (desktopData != null) 'desktop_data': desktopData,
      if (mobileTimestamp != null) 'mobile_timestamp': mobileTimestamp,
      if (desktopTimestamp != null) 'desktop_timestamp': desktopTimestamp,
      if (resolution != null) 'resolution': resolution,
      if (resolvedData != null) 'resolved_data': resolvedData,
      if (createdAt != null) 'created_at': createdAt,
      if (resolvedAt != null) 'resolved_at': resolvedAt,
    });
  }

  SyncConflictsCompanion copyWith(
      {Value<int>? id,
      Value<String>? entityType,
      Value<String>? entityId,
      Value<String>? mobileData,
      Value<String>? desktopData,
      Value<DateTime>? mobileTimestamp,
      Value<DateTime>? desktopTimestamp,
      Value<String?>? resolution,
      Value<String?>? resolvedData,
      Value<DateTime>? createdAt,
      Value<DateTime?>? resolvedAt}) {
    return SyncConflictsCompanion(
      id: id ?? this.id,
      entityType: entityType ?? this.entityType,
      entityId: entityId ?? this.entityId,
      mobileData: mobileData ?? this.mobileData,
      desktopData: desktopData ?? this.desktopData,
      mobileTimestamp: mobileTimestamp ?? this.mobileTimestamp,
      desktopTimestamp: desktopTimestamp ?? this.desktopTimestamp,
      resolution: resolution ?? this.resolution,
      resolvedData: resolvedData ?? this.resolvedData,
      createdAt: createdAt ?? this.createdAt,
      resolvedAt: resolvedAt ?? this.resolvedAt,
    );
  }

  @override
  Map<String, Expression> toColumns(bool nullToAbsent) {
    final map = <String, Expression>{};
    if (id.present) {
      map['id'] = Variable<int>(id.value);
    }
    if (entityType.present) {
      map['entity_type'] = Variable<String>(entityType.value);
    }
    if (entityId.present) {
      map['entity_id'] = Variable<String>(entityId.value);
    }
    if (mobileData.present) {
      map['mobile_data'] = Variable<String>(mobileData.value);
    }
    if (desktopData.present) {
      map['desktop_data'] = Variable<String>(desktopData.value);
    }
    if (mobileTimestamp.present) {
      map['mobile_timestamp'] = Variable<DateTime>(mobileTimestamp.value);
    }
    if (desktopTimestamp.present) {
      map['desktop_timestamp'] = Variable<DateTime>(desktopTimestamp.value);
    }
    if (resolution.present) {
      map['resolution'] = Variable<String>(resolution.value);
    }
    if (resolvedData.present) {
      map['resolved_data'] = Variable<String>(resolvedData.value);
    }
    if (createdAt.present) {
      map['created_at'] = Variable<DateTime>(createdAt.value);
    }
    if (resolvedAt.present) {
      map['resolved_at'] = Variable<DateTime>(resolvedAt.value);
    }
    return map;
  }

  @override
  String toString() {
    return (StringBuffer('SyncConflictsCompanion(')
          ..write('id: $id, ')
          ..write('entityType: $entityType, ')
          ..write('entityId: $entityId, ')
          ..write('mobileData: $mobileData, ')
          ..write('desktopData: $desktopData, ')
          ..write('mobileTimestamp: $mobileTimestamp, ')
          ..write('desktopTimestamp: $desktopTimestamp, ')
          ..write('resolution: $resolution, ')
          ..write('resolvedData: $resolvedData, ')
          ..write('createdAt: $createdAt, ')
          ..write('resolvedAt: $resolvedAt')
          ..write(')'))
        .toString();
  }
}

abstract class _$AppDatabase extends GeneratedDatabase {
  _$AppDatabase(QueryExecutor e) : super(e);
  $AppDatabaseManager get managers => $AppDatabaseManager(this);
  late final $AnimesTable animes = $AnimesTable(this);
  late final $WatchlistTable watchlist = $WatchlistTable(this);
  late final $HistoryTable history = $HistoryTable(this);
  late final $SyncQueueTable syncQueue = $SyncQueueTable(this);
  late final $VipStatusTable vipStatus = $VipStatusTable(this);
  late final $SyncConflictsTable syncConflicts = $SyncConflictsTable(this);
  @override
  Iterable<TableInfo<Table, Object?>> get allTables =>
      allSchemaEntities.whereType<TableInfo<Table, Object?>>();
  @override
  List<DatabaseSchemaEntity> get allSchemaEntities =>
      [animes, watchlist, history, syncQueue, vipStatus, syncConflicts];
}

typedef $$AnimesTableCreateCompanionBuilder = AnimesCompanion Function({
  Value<int> id,
  required String title,
  Value<String?> description,
  required int episodes,
  Value<String?> imageUrl,
  Value<String?> genre,
  Value<String?> status,
  required DateTime addedAt,
  Value<DateTime?> syncedAt,
  required DateTime updatedAt,
});
typedef $$AnimesTableUpdateCompanionBuilder = AnimesCompanion Function({
  Value<int> id,
  Value<String> title,
  Value<String?> description,
  Value<int> episodes,
  Value<String?> imageUrl,
  Value<String?> genre,
  Value<String?> status,
  Value<DateTime> addedAt,
  Value<DateTime?> syncedAt,
  Value<DateTime> updatedAt,
});

class $$AnimesTableFilterComposer
    extends Composer<_$AppDatabase, $AnimesTable> {
  $$AnimesTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get title => $composableBuilder(
      column: $table.title, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get description => $composableBuilder(
      column: $table.description, builder: (column) => ColumnFilters(column));

  ColumnFilters<int> get episodes => $composableBuilder(
      column: $table.episodes, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get imageUrl => $composableBuilder(
      column: $table.imageUrl, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get genre => $composableBuilder(
      column: $table.genre, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get status => $composableBuilder(
      column: $table.status, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get addedAt => $composableBuilder(
      column: $table.addedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnFilters(column));
}

class $$AnimesTableOrderingComposer
    extends Composer<_$AppDatabase, $AnimesTable> {
  $$AnimesTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get title => $composableBuilder(
      column: $table.title, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get description => $composableBuilder(
      column: $table.description, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<int> get episodes => $composableBuilder(
      column: $table.episodes, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get imageUrl => $composableBuilder(
      column: $table.imageUrl, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get genre => $composableBuilder(
      column: $table.genre, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get status => $composableBuilder(
      column: $table.status, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get addedAt => $composableBuilder(
      column: $table.addedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnOrderings(column));
}

class $$AnimesTableAnnotationComposer
    extends Composer<_$AppDatabase, $AnimesTable> {
  $$AnimesTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get title =>
      $composableBuilder(column: $table.title, builder: (column) => column);

  GeneratedColumn<String> get description => $composableBuilder(
      column: $table.description, builder: (column) => column);

  GeneratedColumn<int> get episodes =>
      $composableBuilder(column: $table.episodes, builder: (column) => column);

  GeneratedColumn<String> get imageUrl =>
      $composableBuilder(column: $table.imageUrl, builder: (column) => column);

  GeneratedColumn<String> get genre =>
      $composableBuilder(column: $table.genre, builder: (column) => column);

  GeneratedColumn<String> get status =>
      $composableBuilder(column: $table.status, builder: (column) => column);

  GeneratedColumn<DateTime> get addedAt =>
      $composableBuilder(column: $table.addedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get syncedAt =>
      $composableBuilder(column: $table.syncedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get updatedAt =>
      $composableBuilder(column: $table.updatedAt, builder: (column) => column);
}

class $$AnimesTableTableManager extends RootTableManager<
    _$AppDatabase,
    $AnimesTable,
    AnimeData,
    $$AnimesTableFilterComposer,
    $$AnimesTableOrderingComposer,
    $$AnimesTableAnnotationComposer,
    $$AnimesTableCreateCompanionBuilder,
    $$AnimesTableUpdateCompanionBuilder,
    (AnimeData, BaseReferences<_$AppDatabase, $AnimesTable, AnimeData>),
    AnimeData,
    PrefetchHooks Function()> {
  $$AnimesTableTableManager(_$AppDatabase db, $AnimesTable table)
      : super(TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$AnimesTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$AnimesTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$AnimesTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<String> title = const Value.absent(),
            Value<String?> description = const Value.absent(),
            Value<int> episodes = const Value.absent(),
            Value<String?> imageUrl = const Value.absent(),
            Value<String?> genre = const Value.absent(),
            Value<String?> status = const Value.absent(),
            Value<DateTime> addedAt = const Value.absent(),
            Value<DateTime?> syncedAt = const Value.absent(),
            Value<DateTime> updatedAt = const Value.absent(),
          }) =>
              AnimesCompanion(
            id: id,
            title: title,
            description: description,
            episodes: episodes,
            imageUrl: imageUrl,
            genre: genre,
            status: status,
            addedAt: addedAt,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          createCompanionCallback: ({
            Value<int> id = const Value.absent(),
            required String title,
            Value<String?> description = const Value.absent(),
            required int episodes,
            Value<String?> imageUrl = const Value.absent(),
            Value<String?> genre = const Value.absent(),
            Value<String?> status = const Value.absent(),
            required DateTime addedAt,
            Value<DateTime?> syncedAt = const Value.absent(),
            required DateTime updatedAt,
          }) =>
              AnimesCompanion.insert(
            id: id,
            title: title,
            description: description,
            episodes: episodes,
            imageUrl: imageUrl,
            genre: genre,
            status: status,
            addedAt: addedAt,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          withReferenceMapper: (p0) => p0
              .map((e) => (e.readTable(table), BaseReferences(db, table, e)))
              .toList(),
          prefetchHooksCallback: null,
        ));
}

typedef $$AnimesTableProcessedTableManager = ProcessedTableManager<
    _$AppDatabase,
    $AnimesTable,
    AnimeData,
    $$AnimesTableFilterComposer,
    $$AnimesTableOrderingComposer,
    $$AnimesTableAnnotationComposer,
    $$AnimesTableCreateCompanionBuilder,
    $$AnimesTableUpdateCompanionBuilder,
    (AnimeData, BaseReferences<_$AppDatabase, $AnimesTable, AnimeData>),
    AnimeData,
    PrefetchHooks Function()>;
typedef $$WatchlistTableCreateCompanionBuilder = WatchlistCompanion Function({
  Value<int> id,
  required int animeId,
  Value<String> status,
  Value<int> currentEpisode,
  required DateTime addedAt,
  Value<DateTime?> syncedAt,
  required DateTime updatedAt,
});
typedef $$WatchlistTableUpdateCompanionBuilder = WatchlistCompanion Function({
  Value<int> id,
  Value<int> animeId,
  Value<String> status,
  Value<int> currentEpisode,
  Value<DateTime> addedAt,
  Value<DateTime?> syncedAt,
  Value<DateTime> updatedAt,
});

class $$WatchlistTableFilterComposer
    extends Composer<_$AppDatabase, $WatchlistTable> {
  $$WatchlistTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnFilters(column));

  ColumnFilters<int> get animeId => $composableBuilder(
      column: $table.animeId, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get status => $composableBuilder(
      column: $table.status, builder: (column) => ColumnFilters(column));

  ColumnFilters<int> get currentEpisode => $composableBuilder(
      column: $table.currentEpisode,
      builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get addedAt => $composableBuilder(
      column: $table.addedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnFilters(column));
}

class $$WatchlistTableOrderingComposer
    extends Composer<_$AppDatabase, $WatchlistTable> {
  $$WatchlistTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<int> get animeId => $composableBuilder(
      column: $table.animeId, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get status => $composableBuilder(
      column: $table.status, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<int> get currentEpisode => $composableBuilder(
      column: $table.currentEpisode,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get addedAt => $composableBuilder(
      column: $table.addedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnOrderings(column));
}

class $$WatchlistTableAnnotationComposer
    extends Composer<_$AppDatabase, $WatchlistTable> {
  $$WatchlistTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<int> get animeId =>
      $composableBuilder(column: $table.animeId, builder: (column) => column);

  GeneratedColumn<String> get status =>
      $composableBuilder(column: $table.status, builder: (column) => column);

  GeneratedColumn<int> get currentEpisode => $composableBuilder(
      column: $table.currentEpisode, builder: (column) => column);

  GeneratedColumn<DateTime> get addedAt =>
      $composableBuilder(column: $table.addedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get syncedAt =>
      $composableBuilder(column: $table.syncedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get updatedAt =>
      $composableBuilder(column: $table.updatedAt, builder: (column) => column);
}

class $$WatchlistTableTableManager extends RootTableManager<
    _$AppDatabase,
    $WatchlistTable,
    WatchlistData,
    $$WatchlistTableFilterComposer,
    $$WatchlistTableOrderingComposer,
    $$WatchlistTableAnnotationComposer,
    $$WatchlistTableCreateCompanionBuilder,
    $$WatchlistTableUpdateCompanionBuilder,
    (
      WatchlistData,
      BaseReferences<_$AppDatabase, $WatchlistTable, WatchlistData>
    ),
    WatchlistData,
    PrefetchHooks Function()> {
  $$WatchlistTableTableManager(_$AppDatabase db, $WatchlistTable table)
      : super(TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$WatchlistTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$WatchlistTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$WatchlistTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<int> animeId = const Value.absent(),
            Value<String> status = const Value.absent(),
            Value<int> currentEpisode = const Value.absent(),
            Value<DateTime> addedAt = const Value.absent(),
            Value<DateTime?> syncedAt = const Value.absent(),
            Value<DateTime> updatedAt = const Value.absent(),
          }) =>
              WatchlistCompanion(
            id: id,
            animeId: animeId,
            status: status,
            currentEpisode: currentEpisode,
            addedAt: addedAt,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          createCompanionCallback: ({
            Value<int> id = const Value.absent(),
            required int animeId,
            Value<String> status = const Value.absent(),
            Value<int> currentEpisode = const Value.absent(),
            required DateTime addedAt,
            Value<DateTime?> syncedAt = const Value.absent(),
            required DateTime updatedAt,
          }) =>
              WatchlistCompanion.insert(
            id: id,
            animeId: animeId,
            status: status,
            currentEpisode: currentEpisode,
            addedAt: addedAt,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          withReferenceMapper: (p0) => p0
              .map((e) => (e.readTable(table), BaseReferences(db, table, e)))
              .toList(),
          prefetchHooksCallback: null,
        ));
}

typedef $$WatchlistTableProcessedTableManager = ProcessedTableManager<
    _$AppDatabase,
    $WatchlistTable,
    WatchlistData,
    $$WatchlistTableFilterComposer,
    $$WatchlistTableOrderingComposer,
    $$WatchlistTableAnnotationComposer,
    $$WatchlistTableCreateCompanionBuilder,
    $$WatchlistTableUpdateCompanionBuilder,
    (
      WatchlistData,
      BaseReferences<_$AppDatabase, $WatchlistTable, WatchlistData>
    ),
    WatchlistData,
    PrefetchHooks Function()>;
typedef $$HistoryTableCreateCompanionBuilder = HistoryCompanion Function({
  Value<int> id,
  required int animeId,
  required int episodeNumber,
  Value<String?> episodeTitle,
  required DateTime watchedAt,
  Value<DateTime?> syncedAt,
  required DateTime updatedAt,
});
typedef $$HistoryTableUpdateCompanionBuilder = HistoryCompanion Function({
  Value<int> id,
  Value<int> animeId,
  Value<int> episodeNumber,
  Value<String?> episodeTitle,
  Value<DateTime> watchedAt,
  Value<DateTime?> syncedAt,
  Value<DateTime> updatedAt,
});

class $$HistoryTableFilterComposer
    extends Composer<_$AppDatabase, $HistoryTable> {
  $$HistoryTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnFilters(column));

  ColumnFilters<int> get animeId => $composableBuilder(
      column: $table.animeId, builder: (column) => ColumnFilters(column));

  ColumnFilters<int> get episodeNumber => $composableBuilder(
      column: $table.episodeNumber, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get episodeTitle => $composableBuilder(
      column: $table.episodeTitle, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get watchedAt => $composableBuilder(
      column: $table.watchedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnFilters(column));
}

class $$HistoryTableOrderingComposer
    extends Composer<_$AppDatabase, $HistoryTable> {
  $$HistoryTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<int> get animeId => $composableBuilder(
      column: $table.animeId, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<int> get episodeNumber => $composableBuilder(
      column: $table.episodeNumber,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get episodeTitle => $composableBuilder(
      column: $table.episodeTitle,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get watchedAt => $composableBuilder(
      column: $table.watchedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnOrderings(column));
}

class $$HistoryTableAnnotationComposer
    extends Composer<_$AppDatabase, $HistoryTable> {
  $$HistoryTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<int> get animeId =>
      $composableBuilder(column: $table.animeId, builder: (column) => column);

  GeneratedColumn<int> get episodeNumber => $composableBuilder(
      column: $table.episodeNumber, builder: (column) => column);

  GeneratedColumn<String> get episodeTitle => $composableBuilder(
      column: $table.episodeTitle, builder: (column) => column);

  GeneratedColumn<DateTime> get watchedAt =>
      $composableBuilder(column: $table.watchedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get syncedAt =>
      $composableBuilder(column: $table.syncedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get updatedAt =>
      $composableBuilder(column: $table.updatedAt, builder: (column) => column);
}

class $$HistoryTableTableManager extends RootTableManager<
    _$AppDatabase,
    $HistoryTable,
    HistoryData,
    $$HistoryTableFilterComposer,
    $$HistoryTableOrderingComposer,
    $$HistoryTableAnnotationComposer,
    $$HistoryTableCreateCompanionBuilder,
    $$HistoryTableUpdateCompanionBuilder,
    (HistoryData, BaseReferences<_$AppDatabase, $HistoryTable, HistoryData>),
    HistoryData,
    PrefetchHooks Function()> {
  $$HistoryTableTableManager(_$AppDatabase db, $HistoryTable table)
      : super(TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$HistoryTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$HistoryTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$HistoryTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<int> animeId = const Value.absent(),
            Value<int> episodeNumber = const Value.absent(),
            Value<String?> episodeTitle = const Value.absent(),
            Value<DateTime> watchedAt = const Value.absent(),
            Value<DateTime?> syncedAt = const Value.absent(),
            Value<DateTime> updatedAt = const Value.absent(),
          }) =>
              HistoryCompanion(
            id: id,
            animeId: animeId,
            episodeNumber: episodeNumber,
            episodeTitle: episodeTitle,
            watchedAt: watchedAt,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          createCompanionCallback: ({
            Value<int> id = const Value.absent(),
            required int animeId,
            required int episodeNumber,
            Value<String?> episodeTitle = const Value.absent(),
            required DateTime watchedAt,
            Value<DateTime?> syncedAt = const Value.absent(),
            required DateTime updatedAt,
          }) =>
              HistoryCompanion.insert(
            id: id,
            animeId: animeId,
            episodeNumber: episodeNumber,
            episodeTitle: episodeTitle,
            watchedAt: watchedAt,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          withReferenceMapper: (p0) => p0
              .map((e) => (e.readTable(table), BaseReferences(db, table, e)))
              .toList(),
          prefetchHooksCallback: null,
        ));
}

typedef $$HistoryTableProcessedTableManager = ProcessedTableManager<
    _$AppDatabase,
    $HistoryTable,
    HistoryData,
    $$HistoryTableFilterComposer,
    $$HistoryTableOrderingComposer,
    $$HistoryTableAnnotationComposer,
    $$HistoryTableCreateCompanionBuilder,
    $$HistoryTableUpdateCompanionBuilder,
    (HistoryData, BaseReferences<_$AppDatabase, $HistoryTable, HistoryData>),
    HistoryData,
    PrefetchHooks Function()>;
typedef $$SyncQueueTableCreateCompanionBuilder = SyncQueueCompanion Function({
  Value<int> id,
  required String operation,
  required String entityType,
  required String entityId,
  required String data,
  required DateTime createdAt,
  Value<bool> synced,
  Value<int> retryCount,
  Value<String?> lastError,
});
typedef $$SyncQueueTableUpdateCompanionBuilder = SyncQueueCompanion Function({
  Value<int> id,
  Value<String> operation,
  Value<String> entityType,
  Value<String> entityId,
  Value<String> data,
  Value<DateTime> createdAt,
  Value<bool> synced,
  Value<int> retryCount,
  Value<String?> lastError,
});

class $$SyncQueueTableFilterComposer
    extends Composer<_$AppDatabase, $SyncQueueTable> {
  $$SyncQueueTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get operation => $composableBuilder(
      column: $table.operation, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get entityType => $composableBuilder(
      column: $table.entityType, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get entityId => $composableBuilder(
      column: $table.entityId, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get data => $composableBuilder(
      column: $table.data, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get createdAt => $composableBuilder(
      column: $table.createdAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<bool> get synced => $composableBuilder(
      column: $table.synced, builder: (column) => ColumnFilters(column));

  ColumnFilters<int> get retryCount => $composableBuilder(
      column: $table.retryCount, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get lastError => $composableBuilder(
      column: $table.lastError, builder: (column) => ColumnFilters(column));
}

class $$SyncQueueTableOrderingComposer
    extends Composer<_$AppDatabase, $SyncQueueTable> {
  $$SyncQueueTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get operation => $composableBuilder(
      column: $table.operation, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get entityType => $composableBuilder(
      column: $table.entityType, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get entityId => $composableBuilder(
      column: $table.entityId, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get data => $composableBuilder(
      column: $table.data, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get createdAt => $composableBuilder(
      column: $table.createdAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<bool> get synced => $composableBuilder(
      column: $table.synced, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<int> get retryCount => $composableBuilder(
      column: $table.retryCount, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get lastError => $composableBuilder(
      column: $table.lastError, builder: (column) => ColumnOrderings(column));
}

class $$SyncQueueTableAnnotationComposer
    extends Composer<_$AppDatabase, $SyncQueueTable> {
  $$SyncQueueTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get operation =>
      $composableBuilder(column: $table.operation, builder: (column) => column);

  GeneratedColumn<String> get entityType => $composableBuilder(
      column: $table.entityType, builder: (column) => column);

  GeneratedColumn<String> get entityId =>
      $composableBuilder(column: $table.entityId, builder: (column) => column);

  GeneratedColumn<String> get data =>
      $composableBuilder(column: $table.data, builder: (column) => column);

  GeneratedColumn<DateTime> get createdAt =>
      $composableBuilder(column: $table.createdAt, builder: (column) => column);

  GeneratedColumn<bool> get synced =>
      $composableBuilder(column: $table.synced, builder: (column) => column);

  GeneratedColumn<int> get retryCount => $composableBuilder(
      column: $table.retryCount, builder: (column) => column);

  GeneratedColumn<String> get lastError =>
      $composableBuilder(column: $table.lastError, builder: (column) => column);
}

class $$SyncQueueTableTableManager extends RootTableManager<
    _$AppDatabase,
    $SyncQueueTable,
    SyncQueueData,
    $$SyncQueueTableFilterComposer,
    $$SyncQueueTableOrderingComposer,
    $$SyncQueueTableAnnotationComposer,
    $$SyncQueueTableCreateCompanionBuilder,
    $$SyncQueueTableUpdateCompanionBuilder,
    (
      SyncQueueData,
      BaseReferences<_$AppDatabase, $SyncQueueTable, SyncQueueData>
    ),
    SyncQueueData,
    PrefetchHooks Function()> {
  $$SyncQueueTableTableManager(_$AppDatabase db, $SyncQueueTable table)
      : super(TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$SyncQueueTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$SyncQueueTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$SyncQueueTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<String> operation = const Value.absent(),
            Value<String> entityType = const Value.absent(),
            Value<String> entityId = const Value.absent(),
            Value<String> data = const Value.absent(),
            Value<DateTime> createdAt = const Value.absent(),
            Value<bool> synced = const Value.absent(),
            Value<int> retryCount = const Value.absent(),
            Value<String?> lastError = const Value.absent(),
          }) =>
              SyncQueueCompanion(
            id: id,
            operation: operation,
            entityType: entityType,
            entityId: entityId,
            data: data,
            createdAt: createdAt,
            synced: synced,
            retryCount: retryCount,
            lastError: lastError,
          ),
          createCompanionCallback: ({
            Value<int> id = const Value.absent(),
            required String operation,
            required String entityType,
            required String entityId,
            required String data,
            required DateTime createdAt,
            Value<bool> synced = const Value.absent(),
            Value<int> retryCount = const Value.absent(),
            Value<String?> lastError = const Value.absent(),
          }) =>
              SyncQueueCompanion.insert(
            id: id,
            operation: operation,
            entityType: entityType,
            entityId: entityId,
            data: data,
            createdAt: createdAt,
            synced: synced,
            retryCount: retryCount,
            lastError: lastError,
          ),
          withReferenceMapper: (p0) => p0
              .map((e) => (e.readTable(table), BaseReferences(db, table, e)))
              .toList(),
          prefetchHooksCallback: null,
        ));
}

typedef $$SyncQueueTableProcessedTableManager = ProcessedTableManager<
    _$AppDatabase,
    $SyncQueueTable,
    SyncQueueData,
    $$SyncQueueTableFilterComposer,
    $$SyncQueueTableOrderingComposer,
    $$SyncQueueTableAnnotationComposer,
    $$SyncQueueTableCreateCompanionBuilder,
    $$SyncQueueTableUpdateCompanionBuilder,
    (
      SyncQueueData,
      BaseReferences<_$AppDatabase, $SyncQueueTable, SyncQueueData>
    ),
    SyncQueueData,
    PrefetchHooks Function()>;
typedef $$VipStatusTableCreateCompanionBuilder = VipStatusCompanion Function({
  Value<int> id,
  Value<bool> isVip,
  Value<DateTime?> vipExpiresAt,
  Value<String?> vipTier,
  Value<DateTime?> syncedAt,
  required DateTime updatedAt,
});
typedef $$VipStatusTableUpdateCompanionBuilder = VipStatusCompanion Function({
  Value<int> id,
  Value<bool> isVip,
  Value<DateTime?> vipExpiresAt,
  Value<String?> vipTier,
  Value<DateTime?> syncedAt,
  Value<DateTime> updatedAt,
});

class $$VipStatusTableFilterComposer
    extends Composer<_$AppDatabase, $VipStatusTable> {
  $$VipStatusTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnFilters(column));

  ColumnFilters<bool> get isVip => $composableBuilder(
      column: $table.isVip, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get vipExpiresAt => $composableBuilder(
      column: $table.vipExpiresAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get vipTier => $composableBuilder(
      column: $table.vipTier, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnFilters(column));
}

class $$VipStatusTableOrderingComposer
    extends Composer<_$AppDatabase, $VipStatusTable> {
  $$VipStatusTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<bool> get isVip => $composableBuilder(
      column: $table.isVip, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get vipExpiresAt => $composableBuilder(
      column: $table.vipExpiresAt,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get vipTier => $composableBuilder(
      column: $table.vipTier, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get syncedAt => $composableBuilder(
      column: $table.syncedAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get updatedAt => $composableBuilder(
      column: $table.updatedAt, builder: (column) => ColumnOrderings(column));
}

class $$VipStatusTableAnnotationComposer
    extends Composer<_$AppDatabase, $VipStatusTable> {
  $$VipStatusTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<bool> get isVip =>
      $composableBuilder(column: $table.isVip, builder: (column) => column);

  GeneratedColumn<DateTime> get vipExpiresAt => $composableBuilder(
      column: $table.vipExpiresAt, builder: (column) => column);

  GeneratedColumn<String> get vipTier =>
      $composableBuilder(column: $table.vipTier, builder: (column) => column);

  GeneratedColumn<DateTime> get syncedAt =>
      $composableBuilder(column: $table.syncedAt, builder: (column) => column);

  GeneratedColumn<DateTime> get updatedAt =>
      $composableBuilder(column: $table.updatedAt, builder: (column) => column);
}

class $$VipStatusTableTableManager extends RootTableManager<
    _$AppDatabase,
    $VipStatusTable,
    VipStatusData,
    $$VipStatusTableFilterComposer,
    $$VipStatusTableOrderingComposer,
    $$VipStatusTableAnnotationComposer,
    $$VipStatusTableCreateCompanionBuilder,
    $$VipStatusTableUpdateCompanionBuilder,
    (
      VipStatusData,
      BaseReferences<_$AppDatabase, $VipStatusTable, VipStatusData>
    ),
    VipStatusData,
    PrefetchHooks Function()> {
  $$VipStatusTableTableManager(_$AppDatabase db, $VipStatusTable table)
      : super(TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$VipStatusTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$VipStatusTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$VipStatusTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<bool> isVip = const Value.absent(),
            Value<DateTime?> vipExpiresAt = const Value.absent(),
            Value<String?> vipTier = const Value.absent(),
            Value<DateTime?> syncedAt = const Value.absent(),
            Value<DateTime> updatedAt = const Value.absent(),
          }) =>
              VipStatusCompanion(
            id: id,
            isVip: isVip,
            vipExpiresAt: vipExpiresAt,
            vipTier: vipTier,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          createCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<bool> isVip = const Value.absent(),
            Value<DateTime?> vipExpiresAt = const Value.absent(),
            Value<String?> vipTier = const Value.absent(),
            Value<DateTime?> syncedAt = const Value.absent(),
            required DateTime updatedAt,
          }) =>
              VipStatusCompanion.insert(
            id: id,
            isVip: isVip,
            vipExpiresAt: vipExpiresAt,
            vipTier: vipTier,
            syncedAt: syncedAt,
            updatedAt: updatedAt,
          ),
          withReferenceMapper: (p0) => p0
              .map((e) => (e.readTable(table), BaseReferences(db, table, e)))
              .toList(),
          prefetchHooksCallback: null,
        ));
}

typedef $$VipStatusTableProcessedTableManager = ProcessedTableManager<
    _$AppDatabase,
    $VipStatusTable,
    VipStatusData,
    $$VipStatusTableFilterComposer,
    $$VipStatusTableOrderingComposer,
    $$VipStatusTableAnnotationComposer,
    $$VipStatusTableCreateCompanionBuilder,
    $$VipStatusTableUpdateCompanionBuilder,
    (
      VipStatusData,
      BaseReferences<_$AppDatabase, $VipStatusTable, VipStatusData>
    ),
    VipStatusData,
    PrefetchHooks Function()>;
typedef $$SyncConflictsTableCreateCompanionBuilder = SyncConflictsCompanion
    Function({
  Value<int> id,
  required String entityType,
  required String entityId,
  required String mobileData,
  required String desktopData,
  required DateTime mobileTimestamp,
  required DateTime desktopTimestamp,
  Value<String?> resolution,
  Value<String?> resolvedData,
  required DateTime createdAt,
  Value<DateTime?> resolvedAt,
});
typedef $$SyncConflictsTableUpdateCompanionBuilder = SyncConflictsCompanion
    Function({
  Value<int> id,
  Value<String> entityType,
  Value<String> entityId,
  Value<String> mobileData,
  Value<String> desktopData,
  Value<DateTime> mobileTimestamp,
  Value<DateTime> desktopTimestamp,
  Value<String?> resolution,
  Value<String?> resolvedData,
  Value<DateTime> createdAt,
  Value<DateTime?> resolvedAt,
});

class $$SyncConflictsTableFilterComposer
    extends Composer<_$AppDatabase, $SyncConflictsTable> {
  $$SyncConflictsTableFilterComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnFilters<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get entityType => $composableBuilder(
      column: $table.entityType, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get entityId => $composableBuilder(
      column: $table.entityId, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get mobileData => $composableBuilder(
      column: $table.mobileData, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get desktopData => $composableBuilder(
      column: $table.desktopData, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get mobileTimestamp => $composableBuilder(
      column: $table.mobileTimestamp,
      builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get desktopTimestamp => $composableBuilder(
      column: $table.desktopTimestamp,
      builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get resolution => $composableBuilder(
      column: $table.resolution, builder: (column) => ColumnFilters(column));

  ColumnFilters<String> get resolvedData => $composableBuilder(
      column: $table.resolvedData, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get createdAt => $composableBuilder(
      column: $table.createdAt, builder: (column) => ColumnFilters(column));

  ColumnFilters<DateTime> get resolvedAt => $composableBuilder(
      column: $table.resolvedAt, builder: (column) => ColumnFilters(column));
}

class $$SyncConflictsTableOrderingComposer
    extends Composer<_$AppDatabase, $SyncConflictsTable> {
  $$SyncConflictsTableOrderingComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  ColumnOrderings<int> get id => $composableBuilder(
      column: $table.id, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get entityType => $composableBuilder(
      column: $table.entityType, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get entityId => $composableBuilder(
      column: $table.entityId, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get mobileData => $composableBuilder(
      column: $table.mobileData, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get desktopData => $composableBuilder(
      column: $table.desktopData, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get mobileTimestamp => $composableBuilder(
      column: $table.mobileTimestamp,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get desktopTimestamp => $composableBuilder(
      column: $table.desktopTimestamp,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get resolution => $composableBuilder(
      column: $table.resolution, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<String> get resolvedData => $composableBuilder(
      column: $table.resolvedData,
      builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get createdAt => $composableBuilder(
      column: $table.createdAt, builder: (column) => ColumnOrderings(column));

  ColumnOrderings<DateTime> get resolvedAt => $composableBuilder(
      column: $table.resolvedAt, builder: (column) => ColumnOrderings(column));
}

class $$SyncConflictsTableAnnotationComposer
    extends Composer<_$AppDatabase, $SyncConflictsTable> {
  $$SyncConflictsTableAnnotationComposer({
    required super.$db,
    required super.$table,
    super.joinBuilder,
    super.$addJoinBuilderToRootComposer,
    super.$removeJoinBuilderFromRootComposer,
  });
  GeneratedColumn<int> get id =>
      $composableBuilder(column: $table.id, builder: (column) => column);

  GeneratedColumn<String> get entityType => $composableBuilder(
      column: $table.entityType, builder: (column) => column);

  GeneratedColumn<String> get entityId =>
      $composableBuilder(column: $table.entityId, builder: (column) => column);

  GeneratedColumn<String> get mobileData => $composableBuilder(
      column: $table.mobileData, builder: (column) => column);

  GeneratedColumn<String> get desktopData => $composableBuilder(
      column: $table.desktopData, builder: (column) => column);

  GeneratedColumn<DateTime> get mobileTimestamp => $composableBuilder(
      column: $table.mobileTimestamp, builder: (column) => column);

  GeneratedColumn<DateTime> get desktopTimestamp => $composableBuilder(
      column: $table.desktopTimestamp, builder: (column) => column);

  GeneratedColumn<String> get resolution => $composableBuilder(
      column: $table.resolution, builder: (column) => column);

  GeneratedColumn<String> get resolvedData => $composableBuilder(
      column: $table.resolvedData, builder: (column) => column);

  GeneratedColumn<DateTime> get createdAt =>
      $composableBuilder(column: $table.createdAt, builder: (column) => column);

  GeneratedColumn<DateTime> get resolvedAt => $composableBuilder(
      column: $table.resolvedAt, builder: (column) => column);
}

class $$SyncConflictsTableTableManager extends RootTableManager<
    _$AppDatabase,
    $SyncConflictsTable,
    SyncConflictData,
    $$SyncConflictsTableFilterComposer,
    $$SyncConflictsTableOrderingComposer,
    $$SyncConflictsTableAnnotationComposer,
    $$SyncConflictsTableCreateCompanionBuilder,
    $$SyncConflictsTableUpdateCompanionBuilder,
    (
      SyncConflictData,
      BaseReferences<_$AppDatabase, $SyncConflictsTable, SyncConflictData>
    ),
    SyncConflictData,
    PrefetchHooks Function()> {
  $$SyncConflictsTableTableManager(_$AppDatabase db, $SyncConflictsTable table)
      : super(TableManagerState(
          db: db,
          table: table,
          createFilteringComposer: () =>
              $$SyncConflictsTableFilterComposer($db: db, $table: table),
          createOrderingComposer: () =>
              $$SyncConflictsTableOrderingComposer($db: db, $table: table),
          createComputedFieldComposer: () =>
              $$SyncConflictsTableAnnotationComposer($db: db, $table: table),
          updateCompanionCallback: ({
            Value<int> id = const Value.absent(),
            Value<String> entityType = const Value.absent(),
            Value<String> entityId = const Value.absent(),
            Value<String> mobileData = const Value.absent(),
            Value<String> desktopData = const Value.absent(),
            Value<DateTime> mobileTimestamp = const Value.absent(),
            Value<DateTime> desktopTimestamp = const Value.absent(),
            Value<String?> resolution = const Value.absent(),
            Value<String?> resolvedData = const Value.absent(),
            Value<DateTime> createdAt = const Value.absent(),
            Value<DateTime?> resolvedAt = const Value.absent(),
          }) =>
              SyncConflictsCompanion(
            id: id,
            entityType: entityType,
            entityId: entityId,
            mobileData: mobileData,
            desktopData: desktopData,
            mobileTimestamp: mobileTimestamp,
            desktopTimestamp: desktopTimestamp,
            resolution: resolution,
            resolvedData: resolvedData,
            createdAt: createdAt,
            resolvedAt: resolvedAt,
          ),
          createCompanionCallback: ({
            Value<int> id = const Value.absent(),
            required String entityType,
            required String entityId,
            required String mobileData,
            required String desktopData,
            required DateTime mobileTimestamp,
            required DateTime desktopTimestamp,
            Value<String?> resolution = const Value.absent(),
            Value<String?> resolvedData = const Value.absent(),
            required DateTime createdAt,
            Value<DateTime?> resolvedAt = const Value.absent(),
          }) =>
              SyncConflictsCompanion.insert(
            id: id,
            entityType: entityType,
            entityId: entityId,
            mobileData: mobileData,
            desktopData: desktopData,
            mobileTimestamp: mobileTimestamp,
            desktopTimestamp: desktopTimestamp,
            resolution: resolution,
            resolvedData: resolvedData,
            createdAt: createdAt,
            resolvedAt: resolvedAt,
          ),
          withReferenceMapper: (p0) => p0
              .map((e) => (e.readTable(table), BaseReferences(db, table, e)))
              .toList(),
          prefetchHooksCallback: null,
        ));
}

typedef $$SyncConflictsTableProcessedTableManager = ProcessedTableManager<
    _$AppDatabase,
    $SyncConflictsTable,
    SyncConflictData,
    $$SyncConflictsTableFilterComposer,
    $$SyncConflictsTableOrderingComposer,
    $$SyncConflictsTableAnnotationComposer,
    $$SyncConflictsTableCreateCompanionBuilder,
    $$SyncConflictsTableUpdateCompanionBuilder,
    (
      SyncConflictData,
      BaseReferences<_$AppDatabase, $SyncConflictsTable, SyncConflictData>
    ),
    SyncConflictData,
    PrefetchHooks Function()>;

class $AppDatabaseManager {
  final _$AppDatabase _db;
  $AppDatabaseManager(this._db);
  $$AnimesTableTableManager get animes =>
      $$AnimesTableTableManager(_db, _db.animes);
  $$WatchlistTableTableManager get watchlist =>
      $$WatchlistTableTableManager(_db, _db.watchlist);
  $$HistoryTableTableManager get history =>
      $$HistoryTableTableManager(_db, _db.history);
  $$SyncQueueTableTableManager get syncQueue =>
      $$SyncQueueTableTableManager(_db, _db.syncQueue);
  $$VipStatusTableTableManager get vipStatus =>
      $$VipStatusTableTableManager(_db, _db.vipStatus);
  $$SyncConflictsTableTableManager get syncConflicts =>
      $$SyncConflictsTableTableManager(_db, _db.syncConflicts);
}
