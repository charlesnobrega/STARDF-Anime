import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'glassmorphism_theme.dart';

/// Enum for app theme modes
enum AppThemeMode {
  light,
  dark,
  system,
}

/// State class for theme management
class ThemeState {
  final AppThemeMode mode;
  final bool isDark;

  const ThemeState({
    required this.mode,
    required this.isDark,
  });

  ThemeState copyWith({
    AppThemeMode? mode,
    bool? isDark,
  }) {
    return ThemeState(
      mode: mode ?? this.mode,
      isDark: isDark ?? this.isDark,
    );
  }
}

/// Notifier for theme state management
class ThemeNotifier extends StateNotifier<ThemeState> {
  ThemeNotifier()
      : super(
          const ThemeState(
            mode: AppThemeMode.system,
            isDark: false,
          ),
        );

  void setThemeMode(AppThemeMode mode) {
    state = state.copyWith(mode: mode);
  }

  void setDarkMode(bool isDark) {
    state = state.copyWith(isDark: isDark);
  }

  void toggleTheme() {
    state = state.copyWith(isDark: !state.isDark);
  }
}

/// Provider for theme state
final themeProvider = StateNotifierProvider<ThemeNotifier, ThemeState>((ref) {
  return ThemeNotifier();
});

/// Provider for light theme data
final lightThemeProvider = Provider<ThemeData>((ref) {
  return ThemeData(
    useMaterial3: true,
    brightness: Brightness.light,
    primaryColor: GlassmorphismLightColors.primaryGlass,
    scaffoldBackgroundColor: const Color(0xFFF8F8F8),
    appBarTheme: const AppBarTheme(
      backgroundColor: Color(0xFFF8F8F8),
      elevation: 0,
      centerTitle: true,
      titleTextStyle: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 20,
        fontWeight: FontWeight.bold,
      ),
    ),
    textTheme: const TextTheme(
      displayLarge: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 32,
        fontWeight: FontWeight.bold,
      ),
      displayMedium: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 28,
        fontWeight: FontWeight.bold,
      ),
      displaySmall: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 24,
        fontWeight: FontWeight.bold,
      ),
      headlineMedium: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 20,
        fontWeight: FontWeight.bold,
      ),
      headlineSmall: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 18,
        fontWeight: FontWeight.w600,
      ),
      titleLarge: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 16,
        fontWeight: FontWeight.w600,
      ),
      titleMedium: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 14,
        fontWeight: FontWeight.w500,
      ),
      titleSmall: TextStyle(
        color: GlassmorphismLightColors.textSecondary,
        fontSize: 12,
        fontWeight: FontWeight.w500,
      ),
      bodyLarge: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 16,
      ),
      bodyMedium: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 14,
      ),
      bodySmall: TextStyle(
        color: GlassmorphismLightColors.textSecondary,
        fontSize: 12,
      ),
      labelLarge: TextStyle(
        color: GlassmorphismLightColors.textPrimary,
        fontSize: 14,
        fontWeight: FontWeight.w500,
      ),
    ),
    colorScheme: const ColorScheme.light(
      primary: Color(0xFF6366F1),
      secondary: Color(0xFF8B5CF6),
      tertiary: Color(0xFFEC4899),
      surface: GlassmorphismLightColors.primaryGlass,
      background: Color(0xFFF8F8F8),
      error: Color(0xFFEF4444),
      onPrimary: Colors.white,
      onSecondary: Colors.white,
      onTertiary: Colors.white,
      onSurface: GlassmorphismLightColors.textPrimary,
      onBackground: GlassmorphismLightColors.textPrimary,
      onError: Colors.white,
    ),
  );
});

/// Provider for dark theme data
final darkThemeProvider = Provider<ThemeData>((ref) {
  return ThemeData(
    useMaterial3: true,
    brightness: Brightness.dark,
    primaryColor: GlassmorphismDarkColors.primaryGlass,
    scaffoldBackgroundColor: const Color(0xFF121212),
    appBarTheme: const AppBarTheme(
      backgroundColor: Color(0xFF1E1E1E),
      elevation: 0,
      centerTitle: true,
      titleTextStyle: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 20,
        fontWeight: FontWeight.bold,
      ),
    ),
    textTheme: const TextTheme(
      displayLarge: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 32,
        fontWeight: FontWeight.bold,
      ),
      displayMedium: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 28,
        fontWeight: FontWeight.bold,
      ),
      displaySmall: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 24,
        fontWeight: FontWeight.bold,
      ),
      headlineMedium: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 20,
        fontWeight: FontWeight.bold,
      ),
      headlineSmall: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 18,
        fontWeight: FontWeight.w600,
      ),
      titleLarge: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 16,
        fontWeight: FontWeight.w600,
      ),
      titleMedium: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 14,
        fontWeight: FontWeight.w500,
      ),
      titleSmall: TextStyle(
        color: GlassmorphismDarkColors.textSecondary,
        fontSize: 12,
        fontWeight: FontWeight.w500,
      ),
      bodyLarge: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 16,
      ),
      bodyMedium: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 14,
      ),
      bodySmall: TextStyle(
        color: GlassmorphismDarkColors.textSecondary,
        fontSize: 12,
      ),
      labelLarge: TextStyle(
        color: GlassmorphismDarkColors.textPrimary,
        fontSize: 14,
        fontWeight: FontWeight.w500,
      ),
    ),
    colorScheme: const ColorScheme.dark(
      primary: Color(0xFF818CF8),
      secondary: Color(0xFFA78BFA),
      tertiary: Color(0xFFF472B6),
      surface: GlassmorphismDarkColors.primaryGlass,
      background: Color(0xFF121212),
      error: Color(0xFFFCA5A5),
      onPrimary: Color(0xFF1E1E1E),
      onSecondary: Color(0xFF1E1E1E),
      onTertiary: Color(0xFF1E1E1E),
      onSurface: GlassmorphismDarkColors.textPrimary,
      onBackground: GlassmorphismDarkColors.textPrimary,
      onError: Color(0xFF1E1E1E),
    ),
  );
});

/// Provider for current theme data based on theme state
final currentThemeProvider = Provider<ThemeData>((ref) {
  final themeState = ref.watch(themeProvider);
  final lightTheme = ref.watch(lightThemeProvider);
  final darkTheme = ref.watch(darkThemeProvider);

  return themeState.isDark ? darkTheme : lightTheme;
});

/// Provider for current glassmorphism colors based on theme state
final currentGlassmorphismColorsProvider = Provider<({
  Color primaryGlass,
  Color secondaryGlass,
  Color accentGlass,
  Color textPrimary,
  Color textSecondary,
  Color borderColor,
  Color shadowColor,
})>((ref) {
  final themeState = ref.watch(themeProvider);

  if (themeState.isDark) {
    return (
      primaryGlass: GlassmorphismDarkColors.primaryGlass,
      secondaryGlass: GlassmorphismDarkColors.secondaryGlass,
      accentGlass: GlassmorphismDarkColors.accentGlass,
      textPrimary: GlassmorphismDarkColors.textPrimary,
      textSecondary: GlassmorphismDarkColors.textSecondary,
      borderColor: GlassmorphismDarkColors.borderColor,
      shadowColor: GlassmorphismDarkColors.shadowColor,
    );
  } else {
    return (
      primaryGlass: GlassmorphismLightColors.primaryGlass,
      secondaryGlass: GlassmorphismLightColors.secondaryGlass,
      accentGlass: GlassmorphismLightColors.accentGlass,
      textPrimary: GlassmorphismLightColors.textPrimary,
      textSecondary: GlassmorphismLightColors.textSecondary,
      borderColor: GlassmorphismLightColors.borderColor,
      shadowColor: GlassmorphismLightColors.shadowColor,
    );
  }
});
