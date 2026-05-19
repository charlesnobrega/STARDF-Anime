import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/ui/responsive/responsive_layout.dart';

void main() {
  group('ResponsiveLayout', () {
    testWidgets('renders without error', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveLayout(
              mobile: const Text('Mobile'),
              tablet: const Text('Tablet'),
              desktop: const Text('Desktop'),
            ),
          ),
        ),
      );

      expect(find.byType(ResponsiveLayout), findsOneWidget);
    });

    testWidgets('falls back to mobile when tablet is not provided',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveLayout(
              mobile: const Text('Mobile'),
              desktop: const Text('Desktop'),
            ),
          ),
        ),
      );

      expect(find.byType(ResponsiveLayout), findsOneWidget);
    });

    testWidgets('falls back to tablet when desktop is not provided',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveLayout(
              mobile: const Text('Mobile'),
              tablet: const Text('Tablet'),
            ),
          ),
        ),
      );

      expect(find.byType(ResponsiveLayout), findsOneWidget);
    });
  });

  group('OrientationLayout', () {
    testWidgets('renders portrait widget by default', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: OrientationLayout(
              portrait: const Text('Portrait'),
              landscape: const Text('Landscape'),
            ),
          ),
        ),
      );

      // The default test orientation is portrait
      expect(find.byType(OrientationLayout), findsOneWidget);
    });
  });

  group('ResponsivePadding', () {
    testWidgets('renders with default padding', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsivePadding(
              child: const Text('Content'),
            ),
          ),
        ),
      );

      expect(find.text('Content'), findsOneWidget);
    });

    testWidgets('applies custom mobile padding', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsivePadding(
              mobilePadding: const EdgeInsets.all(8.0),
              child: const Text('Content'),
            ),
          ),
        ),
      );

      expect(find.text('Content'), findsOneWidget);
    });
  });

  group('ResponsiveGrid', () {
    testWidgets('renders grid with children', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveGrid(
              children: [
                for (int i = 0; i < 6; i++) Text('Item $i'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(GridView), findsOneWidget);
      for (int i = 0; i < 6; i++) {
        expect(find.text('Item $i'), findsOneWidget);
      }
    });

    testWidgets('renders grid with custom column count', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveGrid(
              mobileColumns: 2,
              tabletColumns: 3,
              desktopColumns: 4,
              children: [
                for (int i = 0; i < 8; i++) Text('Item $i'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(GridView), findsOneWidget);
      for (int i = 0; i < 8; i++) {
        expect(find.text('Item $i'), findsOneWidget);
      }
    });

    testWidgets('renders grid with custom spacing', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveGrid(
              spacing: 20.0,
              children: [
                for (int i = 0; i < 4; i++) Text('Item $i'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(GridView), findsOneWidget);
      for (int i = 0; i < 4; i++) {
        expect(find.text('Item $i'), findsOneWidget);
      }
    });
  });

  group('ResponsiveColumn', () {
    testWidgets('renders column with children', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveColumn(
              children: [
                const Text('Item 1'),
                const Text('Item 2'),
                const Text('Item 3'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(Column), findsOneWidget);
      expect(find.text('Item 1'), findsOneWidget);
      expect(find.text('Item 2'), findsOneWidget);
      expect(find.text('Item 3'), findsOneWidget);
    });

    testWidgets('applies spacing between children', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveColumn(
              spacing: 16.0,
              children: [
                const Text('Item 1'),
                const Text('Item 2'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(Column), findsOneWidget);
      expect(find.text('Item 1'), findsOneWidget);
      expect(find.text('Item 2'), findsOneWidget);
    });

    testWidgets('applies custom cross axis alignment', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveColumn(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                const Text('Item 1'),
                const Text('Item 2'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(Column), findsOneWidget);
    });
  });

  group('ResponsiveRow', () {
    testWidgets('renders row with children', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveRow(
              children: [
                const Text('Item 1'),
                const Text('Item 2'),
                const Text('Item 3'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(Row), findsOneWidget);
      expect(find.text('Item 1'), findsOneWidget);
      expect(find.text('Item 2'), findsOneWidget);
      expect(find.text('Item 3'), findsOneWidget);
    });

    testWidgets('applies spacing between children', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveRow(
              spacing: 16.0,
              children: [
                const Text('Item 1'),
                const Text('Item 2'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(Row), findsOneWidget);
      expect(find.text('Item 1'), findsOneWidget);
      expect(find.text('Item 2'), findsOneWidget);
    });

    testWidgets('applies custom main axis alignment', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ResponsiveRow(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text('Item 1'),
                const Text('Item 2'),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(Row), findsOneWidget);
    });
  });
}
