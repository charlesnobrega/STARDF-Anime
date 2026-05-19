import 'dart:async';
import 'dart:convert';
import 'package:flutter/services.dart';
import 'package:logger/logger.dart';

/// Response structure from Go backend
class BridgeResponse {
  final String status;
  final dynamic data;
  final BridgeError? error;
  final DateTime timestamp;

  BridgeResponse({
    required this.status,
    this.data,
    this.error,
    required this.timestamp,
  });

  factory BridgeResponse.fromJson(Map<String, dynamic> json) {
    return BridgeResponse(
      status: json['status'] as String,
      data: json['data'],
      error: json['error'] != null
          ? BridgeError.fromJson(json['error'] as Map<String, dynamic>)
          : null,
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }

  bool get isSuccess => status == 'success';
  bool get isError => status == 'error';
}

/// Error structure from Go backend
class BridgeError {
  final String code;
  final String message;
  final Map<String, dynamic>? details;

  BridgeError({
    required this.code,
    required this.message,
    this.details,
  });

  factory BridgeError.fromJson(Map<String, dynamic> json) {
    return BridgeError(
      code: json['code'] as String,
      message: json['message'] as String,
      details: json['details'] as Map<String, dynamic>?,
    );
  }

  @override
  String toString() => 'BridgeError($code): $message';
}

/// Exception thrown when bridge call fails
class BridgeException implements Exception {
  final String code;
  final String message;
  final Map<String, dynamic>? details;

  BridgeException({
    required this.code,
    required this.message,
    this.details,
  });

  @override
  String toString() => 'BridgeException($code): $message';
}

/// BridgeService handles communication between Flutter and Go backend via GoMobile
class BridgeService {
  static const platform = MethodChannel('com.stardf.anime/bridge');
  static const Duration _timeout = Duration(milliseconds: 500);
  
  final Logger _logger = Logger();

  /// Calls a method on the Go backend with timeout
  Future<BridgeResponse> _callMethod(
    String method, [
    Map<String, dynamic>? arguments,
  ]) async {
    try {
      _logger.i('[Bridge] Calling method: $method with args: $arguments');

      final result = await platform
          .invokeMethod<String>(method, arguments)
          .timeout(_timeout);

      if (result == null) {
        throw BridgeException(
          code: 'BRIDGE_NULL_RESPONSE',
          message: 'Received null response from bridge',
        );
      }

      final json = jsonDecode(result) as Map<String, dynamic>;
      final response = BridgeResponse.fromJson(json);

      if (response.isError) {
        _logger.e(
          '[Bridge] Method $method failed: ${response.error}',
        );
        throw BridgeException(
          code: response.error!.code,
          message: response.error!.message,
          details: response.error!.details,
        );
      }

      _logger.i('[Bridge] Method $method succeeded');
      return response;
    } on TimeoutException {
      _logger.e('[Bridge] Method $method timed out after ${_timeout.inMilliseconds}ms');
      throw BridgeException(
        code: 'BRIDGE_TIMEOUT',
        message: 'Bridge call exceeded ${_timeout.inMilliseconds}ms timeout',
      );
    } on PlatformException catch (e) {
      _logger.e('[Bridge] Platform exception: ${e.code} - ${e.message}');
      throw BridgeException(
        code: 'BRIDGE_PLATFORM_ERROR',
        message: e.message ?? 'Unknown platform error',
        details: {'platformCode': e.code},
      );
    } on BridgeException {
      rethrow;
    } catch (e) {
      _logger.e('[Bridge] Unexpected error: $e');
      throw BridgeException(
        code: 'BRIDGE_UNKNOWN_ERROR',
        message: 'Unexpected error: $e',
      );
    }
  }

  /// Gets list of animes from backend
  Future<List<dynamic>> getAnimes() async {
    final response = await _callMethod('getAnimes');
    return response.data as List<dynamic>? ?? [];
  }

  /// Adds anime to watchlist
  Future<Map<String, dynamic>> addToWatchlist(String animeId) async {
    final response = await _callMethod('addToWatchlist', {'animeId': animeId});
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Marks episode as watched
  Future<Map<String, dynamic>> markAsWatched(
    String animeId,
    String episodeId,
  ) async {
    final response = await _callMethod('markAsWatched', {
      'animeId': animeId,
      'episodeId': episodeId,
    });
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Gets current sync status
  Future<Map<String, dynamic>> getSyncStatus() async {
    final response = await _callMethod('getSyncStatus');
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Synchronizes with desktop backend
  Future<Map<String, dynamic>> syncWithDesktop() async {
    final response = await _callMethod('syncWithDesktop');
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Starts the local web server
  Future<Map<String, dynamic>> startWebServer(int port) async {
    final response = await _callMethod('startWebServer', {'port': port});
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Pauses the web server
  Future<Map<String, dynamic>> pauseWebServer() async {
    final response = await _callMethod('pauseWebServer');
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Resumes the web server
  Future<Map<String, dynamic>> resumeWebServer() async {
    final response = await _callMethod('resumeWebServer');
    return response.data as Map<String, dynamic>? ?? {};
  }

  /// Shuts down the web server
  Future<Map<String, dynamic>> shutdownWebServer() async {
    final response = await _callMethod('shutdownWebServer');
    return response.data as Map<String, dynamic>? ?? {};
  }
}
