import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/ui/responsive/responsive_breakpoints.dart';

void main() {
  group('ResponsiveBreakpoints', () {
    test('mobile breakpoint is 600dp', () {
      expect(ResponsiveBreakpoints.mobile, equals(600.0));
    });

    test('tablet breakpoint is 600dp', () {
      expect(ResponsiveBreakpoints.tablet, equals(600.0));
    });

    test('desktop breakpoint is 1200dp', () {
      expect(ResponsiveBreakpoints.desktop, equals(1200.0));
    });

    test('extraLarge breakpoint is 1920dp', () {
      expect(ResponsiveBreakpoints.extraLarge, equals(1920.0));
    });
  });

  group('ResponsiveHelper', () {
    test('getScreenSize returns mobile for width < 600', () {
      final size = ResponsiveHelper.getScreenSize(500.0);
      expect(size, equals(ResponsiveScreenSize.mobile));
    });

    test('getScreenSize returns tablet for width >= 600 and < 1200', () {
      final size = ResponsiveHelper.getScreenSize(800.0);
      expect(size, equals(ResponsiveScreenSize.tablet));
    });

    test('getScreenSize returns desktop for width >= 1200 and < 1920', () {
      final size = ResponsiveHelper.getScreenSize(1400.0);
      expect(size, equals(ResponsiveScreenSize.desktop));
    });

    test('getScreenSize returns extraLarge for width >= 1920', () {
      final size = ResponsiveHelper.getScreenSize(2000.0);
      expect(size, equals(ResponsiveScreenSize.extraLarge));
    });

    test('isMobile returns true for width < 600', () {
      expect(ResponsiveHelper.isMobile(500.0), isTrue);
    });

    test('isMobile returns false for width >= 600', () {
      expect(ResponsiveHelper.isMobile(600.0), isFalse);
    });

    test('isTablet returns true for width >= 600 and < 1200', () {
      expect(ResponsiveHelper.isTablet(800.0), isTrue);
    });

    test('isTablet returns false for width < 600', () {
      expect(ResponsiveHelper.isTablet(500.0), isFalse);
    });

    test('isTablet returns false for width >= 1200', () {
      expect(ResponsiveHelper.isTablet(1200.0), isFalse);
    });

    test('isDesktop returns true for width >= 1200', () {
      expect(ResponsiveHelper.isDesktop(1400.0), isTrue);
    });

    test('isDesktop returns false for width < 1200', () {
      expect(ResponsiveHelper.isDesktop(1000.0), isFalse);
    });
  });

  group('ResponsiveHelper with BuildContext', () {
    testWidgets('getOrientation returns portrait for portrait orientation',
        (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final orientation = ResponsiveHelper.getOrientation(context);
                expect(orientation, equals(ResponsiveOrientation.portrait));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('isPortrait returns true for portrait orientation',
        (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                expect(ResponsiveHelper.isPortrait(context), isTrue);
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('isLandscape returns false for portrait orientation',
        (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                expect(ResponsiveHelper.isLandscape(context), isFalse);
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('getWidth returns screen width', (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final width = ResponsiveHelper.getWidth(context);
                expect(width, greaterThan(0));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('getHeight returns screen height', (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final height = ResponsiveHelper.getHeight(context);
                expect(height, greaterThan(0));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('getPadding returns correct padding for mobile',
        (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final padding = ResponsiveHelper.getPadding(context);
                expect(padding, equals(const EdgeInsets.all(16.0)));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('getGridColumns returns 2 for mobile', (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final columns = ResponsiveHelper.getGridColumns(context);
                expect(columns, equals(2));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('getFontSize returns mobile size for mobile screen',
        (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final fontSize = ResponsiveHelper.getFontSize(
                  context,
                  mobileSize: 14.0,
                  tabletSize: 16.0,
                  desktopSize: 18.0,
                );
                expect(fontSize, equals(14.0));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });

    testWidgets('getSpacing returns mobile spacing for mobile screen',
        (WidgetTester tester) async {
      tester.binding.window.physicalSizeTestValue = const Size(400, 800);
      addTearDown(tester.binding.window.clearPhysicalSizeTestValue);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Builder(
              builder: (context) {
                final spacing = ResponsiveHelper.getSpacing(
                  context,
                  mobileSpacing: 8.0,
                  tabletSpacing: 12.0,
                  desktopSpacing: 16.0,
                );
                expect(spacing, equals(8.0));
                return const SizedBox();
              },
            ),
          ),
        ),
      );
    });
  });
}
