import 'package:flutter/services.dart';
import 'dart:convert';

class GoBridge {
  static const MethodChannel _channel = MethodChannel('com.stardf.anime/bridge');

  static Future<List<dynamic>> search(String query) async {
    final String jsonResponse = await _channel.invokeMethod('search', {'query': query});
    final Map<String, dynamic> response = jsonDecode(jsonResponse);
    if (response.containsKey('error')) {
      throw Exception(response['error']);
    }
    return response['data'];
  }

  static Future<List<dynamic>> getEpisodes(String animeUrl, String source) async {
    final String jsonResponse = await _channel.invokeMethod('getEpisodes', {
      'animeURL': animeUrl,
      'source': source,
    });
    final Map<String, dynamic> response = jsonDecode(jsonResponse);
    if (response.containsKey('error')) {
      throw Exception(response['error']);
    }
    return response['data'];
  }

  static Future<Map<String, dynamic>> getStream(Map<String, dynamic> anime, Map<String, dynamic> episode) async {
    final String jsonResponse = await _channel.invokeMethod('getStream', {
      'animeJSON': jsonEncode(anime),
      'episodeJSON': jsonEncode(episode),
    });
    final Map<String, dynamic> response = jsonDecode(jsonResponse);
    if (response.containsKey('error')) {
      throw Exception(response['error']);
    }
    return response['data'];
  }
}
