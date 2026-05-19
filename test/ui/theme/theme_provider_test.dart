import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/ui/theme/theme_provider.dart';

void main() {
  group('ThemeState', () {
    test('creates ThemeState with correct values', () {
      const state = ThemeState(
        mode: AppThemeMode.light,
        isDark: false,
      );

      expect(state.mode, equals(AppThemeMode.light));
      expect(state.isDark, isFalse);
    });

    test('copyWith creates new instance with updated values', () {
      const state = ThemeState(
        mode: AppThemeMode.light,
        isDark: false,
      );

      final newState = state.copyWith(isDark: true);

      expect(newState.mode, equals(AppThemeMode.light));
      expect(newState.isDark, isTrue);
      expect(identical(state, newState), isFalse);
    });

    test('copyWith preserves unchanged values', () {
      const state = ThemeState(
        mode: AppThemeMode.dark,
        isDark: true,
      );

      final newState = state.copyWith(mode: AppThemeMode.system);

      expect(newState.mode, equals(AppThemeMode.system));
      expect(newState.isDark, isTrue);
    });
  });

  group('ThemeNotifier', () {
    test('initializes with system theme mode and light theme', () {
      final notifier = ThemeNotifier();

      expect(notifier.state.mode, equals(AppThemeMode.system));
      expect(notifier.state.isDark, isFalse);
    });

    test('setThemeMode updates theme mode', () {
      final notifier = ThemeNotifier();

      notifier.setThemeMode(AppThemeMode.dark);

      expect(notifier.state.mode, equals(AppThemeMode.dark));
    });

    test('setDarkMode updates isDark flag', () {
      final notifier = ThemeNotifier();

      notifier.setDarkMode(true);

      expect(notifier.state.isDark, isTrue);
    });

    test('toggleTheme toggles isDark flag', () {
      final notifier = ThemeNotifier();

      expect(notifier.state.isDark, isFalse);

      notifier.toggleTheme();
      expect(notifier.state.isDark, isTrue);

      notifier.toggleTheme();
      expect(notifier.state.isDark, isFalse);
    });
  });

  group('Theme Providers', () {
    testWidgets('lightThemeProvider provides light theme data',
        (WidgetTester tester) async {
      late ThemeData lightTheme;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                lightTheme = ref.watch(lightThemeProvider);
                return const Scaffold(body: SizedBox());
              },
            ),
          ),
        ),
      );

      expect(lightTheme.brightness, equals(Brightness.light));
    });

    testWidgets('darkThemeProvider provides dark theme data',
        (WidgetTester tester) async {
      late ThemeData darkTheme;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                darkTheme = ref.watch(darkThemeProvider);
                return const Scaffold(body: SizedBox());
              },
            ),
          ),
        ),
      );

      expect(darkTheme.brightness, equals(Brightness.dark));
    });

    testWidgets('currentThemeProvider returns light theme when isDark is false',
        (WidgetTester tester) async {
      late ThemeData currentTheme;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                currentTheme = ref.watch(currentThemeProvider);
                return const Scaffold(body: SizedBox());
              },
            ),
          ),
        ),
      );

      expect(currentTheme.brightness, equals(Brightness.light));
    });

    testWidgets('currentThemeProvider returns dark theme when isDark is true',
        (WidgetTester tester) async {
      late ThemeData currentTheme;
      late ThemeNotifier themeNotifier;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                themeNotifier = ref.read(themeProvider.notifier);
                currentTheme = ref.watch(currentThemeProvider);
                return const Scaffold(body: SizedBox());
              },
            ),
          ),
        ),
      );

      // Modify provider outside of build phase
      themeNotifier.setDarkMode(true);
      await tester.pumpAndSettle();

      expect(currentTheme.brightness, equals(Brightness.dark));
    });

    testWidgets('currentGlassmorphismColorsProvider returns light colors when isDark is false',
        (WidgetTester tester) async {
      late var colors;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                colors = ref.watch(currentGlassmorphismColorsProvider);
                return const Scaffold(body: SizedBox());
              },
            ),
          ),
        ),
      );

      expect(colors.textPrimary, equals(const Color(0xFF1A1A1A)));
    });

    testWidgets('currentGlassmorphismColorsProvider returns dark colors when isDark is true',
        (WidgetTester tester) async {
      late var colors;
      late ThemeNotifier themeNotifier;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                themeNotifier = ref.read(themeProvider.notifier);
                colors = ref.watch(currentGlassmorphismColorsProvider);
                return const Scaffold(body: SizedBox());
              },
            ),
          ),
        ),
      );

      // Modify provider outside of build phase
      themeNotifier.setDarkMode(true);
      await tester.pumpAndSettle();

      expect(colors.textPrimary, equals(const Color(0xFFFFFFFF)));
    });
  });

  group('Theme Switching', () {
    testWidgets('switching theme updates currentThemeProvider',
        (WidgetTester tester) async {
      late ThemeData currentTheme;

      await tester.pumpWidget(
        ProviderScope(
          child: MaterialApp(
            home: Consumer(
              builder: (context, ref, child) {
                currentTheme = ref.watch(currentThemeProvider);

                return Scaffold(
                  body: Center(
                    child: ElevatedButton(
                      onPressed: () {
                        ref.read(themeProvider.notifier).toggleTheme();
                      },
                      child: const Text('Toggle Theme'),
                    ),
                  ),
                );
              },
            ),
          ),
        ),
      );

      expect(currentTheme.brightness, equals(Brightness.light));

      await tester.tap(find.text('Toggle Theme'));
      await tester.pumpAndSettle();

      expect(currentTheme.brightness, equals(Brightness.dark));
    });
  });
}
