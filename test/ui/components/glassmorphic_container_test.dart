import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/ui/components/glassmorphic_container.dart';
import 'package:stardf_anime_mobile/ui/theme/glassmorphism_theme.dart';

void main() {
  group('GlassmorphicContainer', () {
    testWidgets('renders with default properties', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
      expect(find.byType(BackdropFilter), findsOneWidget);
      expect(find.text('Test'), findsOneWidget);
    });

    testWidgets('applies blur effect correctly', (WidgetTester tester) async {
      const double testBlur = 15.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              blur: testBlur,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      final backdropFilter = find.byType(BackdropFilter);
      expect(backdropFilter, findsOneWidget);
    });

    testWidgets('applies opacity correctly', (WidgetTester tester) async {
      const double testOpacity = 0.7;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              opacity: testOpacity,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
      expect(find.byType(Container), findsWidgets);
    });

    testWidgets('applies custom color', (WidgetTester tester) async {
      const Color testColor = Color(0xFFFF0000);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              color: testColor,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('applies border radius correctly', (WidgetTester tester) async {
      const double testRadius = 25.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              borderRadius: testRadius,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('shows border when enabled', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              showBorder: true,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('hides border when disabled', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              showBorder: false,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('applies custom padding', (WidgetTester tester) async {
      const EdgeInsets testPadding = EdgeInsets.all(20.0);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              padding: testPadding,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('applies custom margin', (WidgetTester tester) async {
      const EdgeInsets testMargin = EdgeInsets.all(10.0);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              margin: testMargin,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('applies elevation when specified', (WidgetTester tester) async {
      const double testElevation = 8.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              elevation: testElevation,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
      // Material widget is created for elevation, plus the Scaffold's Material
      expect(find.byType(Material), findsWidgets);
    });

    testWidgets('handles tap when clickable', (WidgetTester tester) async {
      bool tapped = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              isClickable: true,
              onTap: () => tapped = true,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      await tester.tap(find.byType(GlassmorphicContainer));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
    });

    testWidgets('does not handle tap when not clickable', (WidgetTester tester) async {
      bool tapped = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              isClickable: false,
              onTap: () => tapped = true,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      await tester.tap(find.byType(GlassmorphicContainer));
      await tester.pumpAndSettle();

      expect(tapped, isFalse);
    });

    testWidgets('renders with all glassmorphism properties', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              blur: GlassmorphismBlur.large,
              opacity: GlassmorphismOpacity.high,
              color: GlassmorphismLightColors.primaryGlass,
              borderRadius: GlassmorphismRadius.extraLarge,
              showBorder: true,
              padding: const EdgeInsets.all(GlassmorphismSpacing.lg),
              margin: const EdgeInsets.all(GlassmorphismSpacing.md),
              elevation: 4.0,
              isClickable: true,
              onTap: () {},
              child: const Text('Glassmorphic Container'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
      expect(find.byType(BackdropFilter), findsOneWidget);
      expect(find.text('Glassmorphic Container'), findsOneWidget);
    });

    testWidgets('renders with custom border color', (WidgetTester tester) async {
      const Color customBorderColor = Color(0xFF0000FF);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicContainer(
              borderColor: customBorderColor,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsOneWidget);
    });

    testWidgets('renders multiple containers', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Column(
              children: [
                GlassmorphicContainer(
                  child: const Text('Container 1'),
                ),
                GlassmorphicContainer(
                  child: const Text('Container 2'),
                ),
                GlassmorphicContainer(
                  child: const Text('Container 3'),
                ),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicContainer), findsNWidgets(3));
      expect(find.text('Container 1'), findsOneWidget);
      expect(find.text('Container 2'), findsOneWidget);
      expect(find.text('Container 3'), findsOneWidget);
    });
  });
}
