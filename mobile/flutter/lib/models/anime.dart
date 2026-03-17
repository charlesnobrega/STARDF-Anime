class Anime {
  final String name;
  final String url;
  final String imageUrl;
  final String source;
  final List<Episode> episodes;
  final int? anilistId;
  final AnimeDetails? details;

  Anime({
    required this.name,
    required this.url,
    required this.imageUrl,
    required this.source,
    required this.episodes,
    this.anilistId,
    this.details,
  });

  factory Anime.fromJson(Map<String, dynamic> json) {
    return Anime(
      name: json['Name'] ?? '',
      url: json['URL'] ?? '',
      imageUrl: json['ImageURL'] ?? '',
      source: json['Source'] ?? '',
      episodes: (json['Episodes'] as List?)
              ?.map((e) => Episode.fromJson(e))
              .toList() ??
          [],
      anilistId: json['AnilistID'],
      details: json['Details'] != null ? AnimeDetails.fromJson(json['Details']) : null,
    );
  }

  Map<String, dynamic> toJson() => {
    'Name': name,
    'URL': url,
    'ImageURL': imageUrl,
    'Source': source,
    'AnilistID': anilistId,
    'Details': details?.toJson(),
  };
}

class Episode {
  final String number;
  final int num;
  final String url;
  final String? title;
  final bool isFiller;

  Episode({
    required this.number,
    required this.num,
    required this.url,
    this.title,
    this.isFiller = false,
  });

  factory Episode.fromJson(Map<String, dynamic> json) {
    return Episode(
      number: json['Number'] ?? '',
      num: json['Num'] ?? 0,
      url: json['URL'] ?? '',
      title: json['Title']?['English'] ?? json['Title']?['Romaji'],
      isFiller: json['IsFiller'] ?? false,
    );
  }

  Map<String, dynamic> toJson() => {
    'Number': number,
    'Num': num,
    'URL': url,
    'IsFiller': isFiller,
  };
}

class AnimeDetails {
  final String description;
  final List<String> genres;
  final int averageScore;
  final String status;

  AnimeDetails({
    required this.description,
    required this.genres,
    required this.averageScore,
    required this.status,
  });

  factory AnimeDetails.fromJson(Map<String, dynamic> json) {
    return AnimeDetails(
      description: json['Description'] ?? '',
      genres: (json['Genres'] as List?)?.map((e) => e.toString()).toList() ?? [],
      averageScore: json['AverageScore'] ?? 0,
      status: json['Status'] ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
    'Description': description,
    'Genres': genres,
    'AverageScore': averageScore,
    'Status': status,
  };
}
