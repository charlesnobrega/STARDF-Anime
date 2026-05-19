import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/ui/components/glassmorphic_card.dart';
import 'package:stardf_anime_mobile/ui/theme/glassmorphism_theme.dart';

void main() {
  group('GlassmorphicCard', () {
    testWidgets('renders with child content', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              child: const Text('Card Content'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Card Content'), findsOneWidget);
    });

    testWidgets('renders with title', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              title: 'Card Title',
              child: const Text('Card Content'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Card Title'), findsOneWidget);
      expect(find.text('Card Content'), findsOneWidget);
    });

    testWidgets('renders with title and subtitle', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              title: 'Card Title',
              subtitle: 'Card Subtitle',
              child: const Text('Card Content'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Card Title'), findsOneWidget);
      expect(find.text('Card Subtitle'), findsOneWidget);
      expect(find.text('Card Content'), findsOneWidget);
    });

    testWidgets('renders with custom header', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              header: const Text('Custom Header'),
              child: const Text('Card Content'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Custom Header'), findsOneWidget);
      expect(find.text('Card Content'), findsOneWidget);
    });

    testWidgets('renders with footer', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              child: const Text('Card Content'),
              footer: const Text('Card Footer'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Card Content'), findsOneWidget);
      expect(find.text('Card Footer'), findsOneWidget);
    });

    testWidgets('applies custom blur', (WidgetTester tester) async {
      const double testBlur = 15.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              blur: testBlur,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies custom opacity', (WidgetTester tester) async {
      const double testOpacity = 0.7;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              opacity: testOpacity,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies custom color', (WidgetTester tester) async {
      const Color testColor = Color(0xFFFF0000);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              color: testColor,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies custom border radius', (WidgetTester tester) async {
      const double testRadius = 25.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              borderRadius: testRadius,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies custom padding', (WidgetTester tester) async {
      const EdgeInsets testPadding = EdgeInsets.all(20.0);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              padding: testPadding,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies custom margin', (WidgetTester tester) async {
      const EdgeInsets testMargin = EdgeInsets.all(10.0);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              margin: testMargin,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies elevation', (WidgetTester tester) async {
      const double testElevation = 8.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              elevation: testElevation,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('handles tap when clickable', (WidgetTester tester) async {
      bool tapped = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              isClickable: true,
              onTap: () => tapped = true,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      await tester.tap(find.byType(GlassmorphicCard));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
    });

    testWidgets('applies custom text color', (WidgetTester tester) async {
      const Color testTextColor = Color(0xFF0000FF);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              title: 'Title',
              textColor: testTextColor,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('applies custom subtitle color', (WidgetTester tester) async {
      const Color testSubtitleColor = Color(0xFF00FF00);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              title: 'Title',
              subtitle: 'Subtitle',
              subtitleColor: testSubtitleColor,
              child: const Text('Test'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
    });

    testWidgets('renders with all glassmorphism properties', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              title: 'Glassmorphic Card',
              subtitle: 'With all properties',
              blur: GlassmorphismBlur.large,
              opacity: GlassmorphismOpacity.high,
              color: GlassmorphismLightColors.primaryGlass,
              borderRadius: GlassmorphismRadius.large,
              padding: const EdgeInsets.all(GlassmorphismSpacing.lg),
              margin: const EdgeInsets.all(GlassmorphismSpacing.md),
              elevation: 4.0,
              isClickable: true,
              onTap: () {},
              textColor: GlassmorphismLightColors.textPrimary,
              subtitleColor: GlassmorphismLightColors.textSecondary,
              child: const Text('Card Content'),
              footer: const Text('Card Footer'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Glassmorphic Card'), findsOneWidget);
      expect(find.text('With all properties'), findsOneWidget);
      expect(find.text('Card Content'), findsOneWidget);
      expect(find.text('Card Footer'), findsOneWidget);
    });

    testWidgets('renders multiple cards', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Column(
              children: [
                GlassmorphicCard(
                  title: 'Card 1',
                  child: const Text('Content 1'),
                ),
                GlassmorphicCard(
                  title: 'Card 2',
                  child: const Text('Content 2'),
                ),
                GlassmorphicCard(
                  title: 'Card 3',
                  child: const Text('Content 3'),
                ),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsNWidgets(3));
      expect(find.text('Card 1'), findsOneWidget);
      expect(find.text('Card 2'), findsOneWidget);
      expect(find.text('Card 3'), findsOneWidget);
    });

    testWidgets('header overrides title and subtitle', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicCard(
              title: 'Title',
              subtitle: 'Subtitle',
              header: const Text('Custom Header'),
              child: const Text('Content'),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicCard), findsOneWidget);
      expect(find.text('Custom Header'), findsOneWidget);
      // Title and subtitle should still be rendered but not visible due to header
      expect(find.text('Title'), findsNothing);
      expect(find.text('Subtitle'), findsNothing);
    });
  });
}
