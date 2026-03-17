import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:html/parser.dart' as html_parser;

/// Pure Dart implementation of the anime scrapers.
/// Updated to use animefire.io and robust CSS selectors.
class ScraperService {
  static const _baseUrl = 'https://animefire.io';
  static const _userAgent =
      'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36';

  static final _client = http.Client();

  static Map<String, String> get _headers => {
        'User-Agent': _userAgent,
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        'Accept-Language': 'pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7',
      };

  // ─── AnimeFire Scraper ─────────────────────────────────────────────

  static Future<List<Map<String, dynamic>>> searchAnimeFire(String query) async {
    final searchParam = query.toLowerCase().replaceAll(' ', '-');
    final url = '$_baseUrl/pesquisar/$searchParam';

    try {
      final response = await _client
          .get(Uri.parse(url), headers: _headers)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode != 200) return [];

      final document = html_parser.parse(response.body);
      final results = <Map<String, dynamic>>[];

      // Select columns containing the anime cards
      final cards = document.querySelectorAll('.divCardUltimosEps');
      
      for (final card in cards) {
        final linkElement = card.querySelector('article > a');
        final titleElement = card.querySelector('h3.animeTitle');
        final imgElement = card.querySelector('img.card-img-top');

        final title = titleElement?.text.trim() ?? '';
        final href = linkElement?.attributes['href'] ?? '';
        
        // Handle lazy loading image attributes
        final imageUrl = imgElement?.attributes['data-src'] ?? 
                         imgElement?.attributes['src'] ?? '';

        if (title.isNotEmpty && href.isNotEmpty) {
          results.add({
            'name': title,
            'url': href,
            'imageUrl': imageUrl,
            'source': 'AnimeFire',
          });
        }
      }
      return results;
    } catch (e) {
      debugPrint('AnimeFire search error: $e');
      return [];
    }
  }

  static Future<List<Map<String, dynamic>>> getAnimeFireEpisodes(
      String animeUrl) async {
    try {
      // Ensure the URL is absolute
      final targetUrl = animeUrl.startsWith('http') ? animeUrl : '$_baseUrl$animeUrl';
      
      final response = await _client
          .get(Uri.parse(targetUrl), headers: _headers)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode != 200) return [];

      final document = html_parser.parse(response.body);
      final results = <Map<String, dynamic>>[];

      // Episode links usually have class .lEp.epT
      final episodes = document.querySelectorAll('.div_video_list a.lEp');
      
      for (int i = 0; i < episodes.length; i++) {
        final ep = episodes[i];
        final href = ep.attributes['href'] ?? '';
        final title = ep.text.trim();

        if (href.isNotEmpty) {
          results.add({
            'number': i + 1,
            'title': title.isNotEmpty ? title : 'Episódio ${i + 1}',
            'url': href,
          });
        }
      }
      return results;
    } catch (e) {
      debugPrint('AnimeFire episodes error: $e');
      return [];
    }
  }

  static Future<String?> getAnimeFireStreamUrl(String episodeUrl) async {
    try {
      final targetUrl = episodeUrl.startsWith('http') ? episodeUrl : '$_baseUrl$episodeUrl';
      
      final response = await _client
          .get(Uri.parse(targetUrl), headers: _headers)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode != 200) return null;

      final body = response.body;

      // 1. Look for the iframe/video source page link (data-video-src)
      final videoPageRegex = RegExp(r'data-video-src="([^"]+)"');
      final match = videoPageRegex.firstMatch(body);
      
      if (match != null) {
        final videoPageUrl = match.group(1)!;
        final videoResponse = await _client
            .get(Uri.parse(videoPageUrl), headers: _headers)
            .timeout(const Duration(seconds: 10));

        // 2. Look for .m3u8 or .mp4 in the video source page
        final streamRegex = RegExp(r'(https?://[^\s"]+\.(?:m3u8|mp4)[^\s"]*)');
        final streamMatch = streamRegex.firstMatch(videoResponse.body);
        if (streamMatch != null) return streamMatch.group(1);
      }

      // 3. Fallback: direct video URLs in the main page body
      final directRegex = RegExp(r"""(https?://[^\s"']+\.(?:m3u8|mp4)[^\s"']*)""");
      final directMatch = directRegex.firstMatch(body);
      if (directMatch != null) return directMatch.group(1);

      return null;
    } catch (e) {
      debugPrint('AnimeFire stream error: $e');
      return null;
    }
  }

  // ─── Unified Search ────────────────────────────────────────────────

  static Future<List<Map<String, dynamic>>> searchAll(String query) async {
    final results = await Future.wait([
      searchAnimeFire(query),
    ]);

    return results.expand((list) => list).toList();
  }
}
