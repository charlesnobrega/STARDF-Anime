class Anime {
  final String name;
  final String url;
  final String imageUrl;
  final String source;
  final List<Episode> episodes;

  Anime({
    required this.name,
    required this.url,
    required this.imageUrl,
    required this.source,
    required this.episodes,
  });

  factory Anime.fromJson(Map<String, dynamic> json) {
    return Anime(
      name: json['name'] ?? json['Name'] ?? '',
      url: json['url'] ?? json['URL'] ?? '',
      imageUrl: json['imageUrl'] ?? json['ImageURL'] ?? '',
      source: json['source'] ?? json['Source'] ?? '',
      episodes: (json['episodes'] as List?)
              ?.map((e) => Episode.fromJson(e))
              .toList() ??
          (json['Episodes'] as List?)
              ?.map((e) => Episode.fromJson(e))
              .toList() ??
          [],
    );
  }

  Map<String, dynamic> toJson() => {
    'name': name,
    'url': url,
    'imageUrl': imageUrl,
    'source': source,
  };
}

class Episode {
  final int number;
  final String title;
  final String url;

  Episode({
    required this.number,
    required this.title,
    required this.url,
  });

  factory Episode.fromJson(Map<String, dynamic> json) {
    return Episode(
      number: json['number'] ?? json['Num'] ?? 0,
      title: json['title'] ?? json['Number'] ?? '',
      url: json['url'] ?? json['URL'] ?? '',
    );
  }

  Map<String, dynamic> toJson() => {
    'number': number,
    'title': title,
    'url': url,
  };
}
