import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:stardf_anime_mobile/core/logging/app_logger.dart';
import 'package:stardf_anime_mobile/core/logging/logger_extension.dart';

/// Error categories for different types of errors
enum ErrorCategory {
  network,
  bridge,
  sync,
  resource,
  unknown,
}

/// Represents an application error
class AppError {
  final String code;
  final String message;
  final String? details;
  final ErrorCategory category;
  final StackTrace? stackTrace;
  final DateTime timestamp;

  AppError({
    required this.code,
    required this.message,
    this.details,
    required this.category,
    this.stackTrace,
    DateTime? timestamp,
  }) : timestamp = timestamp ?? DateTime.now();

  @override
  String toString() {
    return 'AppError(code: $code, message: $message, category: $category, timestamp: $timestamp)';
  }
}

/// Manages error handling and recovery
class ErrorHandler {
  final AppLogger _logger;
  final List<AppError> _errorHistory = [];
  static const int _maxErrorHistory = 100;

  ErrorHandler(this._logger);

  /// Handle an error
  AppError handleError({
    required String code,
    required String message,
    String? details,
    required ErrorCategory category,
    StackTrace? stackTrace,
  }) {
    final error = AppError(
      code: code,
      message: message,
      details: details,
      category: category,
      stackTrace: stackTrace,
    );

    _addToHistory(error);
    _logError(error);

    return error;
  }

  /// Handle an exception
  AppError handleException(
    Object exception, {
    String? code,
    String? message,
    String? details,
    required ErrorCategory category,
    StackTrace? stackTrace,
  }) {
    final errorCode = code ?? _getErrorCode(exception);
    final errorMessage = message ?? _getErrorMessage(exception);

    return handleError(
      code: errorCode,
      message: errorMessage,
      details: details,
      category: category,
      stackTrace: stackTrace,
    );
  }

  /// Get user-friendly error message
  String getUserMessage(AppError error) {
    switch (error.category) {
      case ErrorCategory.network:
        return 'Erro de conexão. Verifique sua internet e tente novamente.';
      case ErrorCategory.bridge:
        return 'Erro ao comunicar com o backend. Tente novamente.';
      case ErrorCategory.sync:
        return 'Erro ao sincronizar dados. Tente novamente.';
      case ErrorCategory.resource:
        return 'Recurso não disponível. Tente novamente mais tarde.';
      case ErrorCategory.unknown:
        return 'Erro desconhecido. Tente novamente.';
    }
  }

  /// Get error history
  List<AppError> getErrorHistory() {
    return List.unmodifiable(_errorHistory);
  }

  /// Clear error history
  void clearErrorHistory() {
    _errorHistory.clear();
  }

  /// Get errors by category
  List<AppError> getErrorsByCategory(ErrorCategory category) {
    return _errorHistory.where((e) => e.category == category).toList();
  }

  /// Get recent errors
  List<AppError> getRecentErrors({int limit = 10}) {
    return _errorHistory.skip((_errorHistory.length - limit).clamp(0, _errorHistory.length)).toList();
  }

  /// Private helper methods

  void _addToHistory(AppError error) {
    _errorHistory.add(error);
    if (_errorHistory.length > _maxErrorHistory) {
      _errorHistory.removeAt(0);
    }
  }

  void _logError(AppError error) {
    final logMessage = '${error.code}: ${error.message}';
    final details = error.details != null ? '\nDetails: ${error.details}' : '';

    _logger.error(
      '$logMessage$details',
      tag: 'ErrorHandler',
      stackTrace: error.stackTrace,
    );
  }

  String _getErrorCode(Object exception) {
    if (exception is TimeoutException) {
      return 'TIMEOUT_ERROR';
    } else if (exception is NetworkException) {
      return 'NETWORK_ERROR';
    } else if (exception is BridgeException) {
      return 'BRIDGE_ERROR';
    }
    return 'UNKNOWN_ERROR';
  }

  String _getErrorMessage(Object exception) {
    return exception.toString();
  }
}

/// Custom exception classes

class NetworkException implements Exception {
  final String message;
  NetworkException(this.message);

  @override
  String toString() => message;
}

class BridgeException implements Exception {
  final String message;
  final String? code;
  BridgeException(this.message, {this.code});

  @override
  String toString() => message;
}

class TimeoutException implements Exception {
  final String message;
  TimeoutException(this.message);

  @override
  String toString() => message;
}

class SyncException implements Exception {
  final String message;
  SyncException(this.message);

  @override
  String toString() => message;
}

class ResourceException implements Exception {
  final String message;
  ResourceException(this.message);

  @override
  String toString() => message;
}

/// Provider for ErrorHandler
final errorHandlerProvider = Provider<ErrorHandler>((ref) {
  final logger = AppLogger();
  return ErrorHandler(logger);
});

/// Provider for error history
final errorHistoryProvider = StateProvider<List<AppError>>((ref) {
  return [];
});

/// Provider for current error
final currentErrorProvider = StateProvider<AppError?>((ref) {
  return null;
});

/// Provider for error statistics
final errorStatisticsProvider = Provider<Map<String, dynamic>>((ref) {
  final errorHandler = ref.watch(errorHandlerProvider);
  final history = errorHandler.getErrorHistory();

  final categoryCount = <ErrorCategory, int>{};
  for (final error in history) {
    categoryCount[error.category] = (categoryCount[error.category] ?? 0) + 1;
  }

  return {
    'totalErrors': history.length,
    'byCategory': categoryCount,
    'recentErrors': errorHandler.getRecentErrors(limit: 5),
  };
});
