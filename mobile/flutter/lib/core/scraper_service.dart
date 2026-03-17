import 'package:http/http.dart' as http;
import 'package:html/parser.dart' as html_parser;

/// Pure Dart implementation of the anime scrapers.
/// This replaces the Go bridge approach which fails due to
/// incompatible TUI dependencies in the Go dependency tree.
class ScraperService {
  static const _userAgent =
      'Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36';

  static final _client = http.Client();

  static Map<String, String> get _headers => {
        'User-Agent': _userAgent,
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        'Accept-Language': 'pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7',
      };

  // ─── AnimeFire Scraper ─────────────────────────────────────────────

  static Future<List<Map<String, dynamic>>> searchAnimeFire(String query) async {
    final encoded = Uri.encodeComponent(query.toLowerCase().replaceAll(' ', '-'));
    final url = 'https://animefire.plus/pesquisar/$encoded';
    
    try {
      final response = await _client.get(Uri.parse(url), headers: _headers)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode != 200) return [];

      final document = html_parser.parse(response.body);
      final results = <Map<String, dynamic>>[];

      final cards = document.querySelectorAll('.row.ml-1.mr-1 a');
      for (final card in cards) {
        final title = card.querySelector('.animeTitle')?.text.trim() ?? '';
        final href = card.attributes['href'] ?? '';
        final img = card.querySelector('img');
        final imageUrl = img?.attributes['data-src'] ?? img?.attributes['src'] ?? '';

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
      print('AnimeFire search error: $e');
      return [];
    }
  }

  static Future<List<Map<String, dynamic>>> getAnimeFireEpisodes(String animeUrl) async {
    try {
      final response = await _client.get(Uri.parse(animeUrl), headers: _headers)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode != 200) return [];

      final document = html_parser.parse(response.body);
      final results = <Map<String, dynamic>>[];

      final episodes = document.querySelectorAll('.div_video_list a');
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
      print('AnimeFire episodes error: $e');
      return [];
    }
  }

  static Future<String?> getAnimeFireStreamUrl(String episodeUrl) async {
    try {
      final response = await _client.get(Uri.parse(episodeUrl), headers: _headers)
          .timeout(const Duration(seconds: 10));

      if (response.statusCode != 200) return null;

      final body = response.body;
      
      // Look for video data-video-src or video source
      final videoRegex = RegExp(r'data-video-src="([^"]+)"');
      final match = videoRegex.firstMatch(body);
      if (match != null) {
        final videoPageUrl = match.group(1)!;
        // Fetch the video page to get the actual stream URL
        final videoResponse = await _client.get(Uri.parse(videoPageUrl), headers: _headers)
            .timeout(const Duration(seconds: 10));
        
        // Look for .m3u8 or .mp4 URLs
        final m3u8Regex = RegExp(r'(https?://[^\s"]+\.m3u8[^\s"]*)');
        final m3u8Match = m3u8Regex.firstMatch(videoResponse.body);
        if (m3u8Match != null) {
          return m3u8Match.group(1);
        }
        
        final mp4Regex = RegExp(r'(https?://[^\s"]+\.mp4[^\s"]*)');
        final mp4Match = mp4Regex.firstMatch(videoResponse.body);
        if (mp4Match != null) {
          return mp4Match.group(1);
        }
      }
      
      // Alternative: look for direct video URLs in the page
      final directM3u8 = RegExp(r'(https?://[^\s"\']+\.m3u8[^\s"\']*)');
      final directMatch = directM3u8.firstMatch(body);
      if (directMatch != null) {
        return directMatch.group(1);
      }

      return null;
    } catch (e) {
      print('AnimeFire stream error: $e');
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
