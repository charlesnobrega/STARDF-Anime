import 'scraper_service.dart';

/// Bridge abstraction layer.
/// Previously used MethodChannel to communicate with Go core.
/// Now uses pure Dart scrapers (ScraperService) for Android compatibility.
class GoBridge {
  static Future<List<dynamic>> search(String query) async {
    return await ScraperService.searchAll(query);
  }

  static Future<List<dynamic>> getEpisodes(String animeUrl, String source) async {
    // Route to the correct scraper based on source
    switch (source) {
      case 'AnimeFire':
        return await ScraperService.getAnimeFireEpisodes(animeUrl);
      default:
        return await ScraperService.getAnimeFireEpisodes(animeUrl);
    }
  }

  static Future<Map<String, dynamic>> getStream(
      Map<String, dynamic> anime, Map<String, dynamic> episode) async {
    final source = anime['source'] ?? 'AnimeFire';
    String? url;

    switch (source) {
      case 'AnimeFire':
        url = await ScraperService.getAnimeFireStreamUrl(episode['url'] ?? '');
        break;
      default:
        url = await ScraperService.getAnimeFireStreamUrl(episode['url'] ?? '');
    }

    if (url == null) {
      throw Exception('Não foi possível obter a URL do stream');
    }

    return {
      'url': url,
      'metadata': {'source': source},
    };
  }
}
