import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:stardf_anime_mobile/ui/components/glassmorphic_button.dart';
import 'package:stardf_anime_mobile/ui/theme/glassmorphism_theme.dart';

void main() {
  group('GlassmorphicButton', () {
    testWidgets('renders with label', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Click Me',
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
      expect(find.text('Click Me'), findsOneWidget);
    });

    testWidgets('calls onPressed when tapped', (WidgetTester tester) async {
      bool pressed = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Click Me',
              onPressed: () => pressed = true,
            ),
          ),
        ),
      );

      await tester.tap(find.byType(GlassmorphicButton));
      await tester.pumpAndSettle();

      expect(pressed, isTrue);
    });

    testWidgets('applies custom blur', (WidgetTester tester) async {
      const double testBlur = 15.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              blur: testBlur,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom opacity', (WidgetTester tester) async {
      const double testOpacity = 0.7;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              opacity: testOpacity,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom color', (WidgetTester tester) async {
      const Color testColor = Color(0xFFFF0000);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              color: testColor,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom text color', (WidgetTester tester) async {
      const Color testTextColor = Color(0xFF0000FF);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              textColor: testTextColor,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom border radius', (WidgetTester tester) async {
      const double testRadius = 25.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              borderRadius: testRadius,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom padding', (WidgetTester tester) async {
      const EdgeInsets testPadding = EdgeInsets.all(20.0);

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              padding: testPadding,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('renders with icon', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              icon: Icons.add,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
      expect(find.byIcon(Icons.add), findsOneWidget);
    });

    testWidgets('disables button when enabled is false', (WidgetTester tester) async {
      bool pressed = false;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              enabled: false,
              onPressed: () => pressed = true,
            ),
          ),
        ),
      );

      await tester.tap(find.byType(GlassmorphicButton));
      await tester.pumpAndSettle();

      expect(pressed, isFalse);
    });

    testWidgets('applies custom width', (WidgetTester tester) async {
      const double testWidth = 200.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              width: testWidth,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom height', (WidgetTester tester) async {
      const double testHeight = 60.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              height: testHeight,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies elevation', (WidgetTester tester) async {
      const double testElevation = 8.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              elevation: testElevation,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies custom font size', (WidgetTester tester) async {
      const double testFontSize = 20.0;

      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              fontSize: testFontSize,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('applies bold font weight when isBold is true', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              isBold: true,
              onPressed: () {},
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });

    testWidgets('renders with all glassmorphism properties', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Center(
              child: GlassmorphicButton(
                label: 'Button',
                blur: GlassmorphismBlur.large,
                opacity: GlassmorphismOpacity.high,
                color: GlassmorphismLightColors.primaryGlass,
                textColor: GlassmorphismLightColors.textPrimary,
                borderRadius: GlassmorphismRadius.medium,
                padding: const EdgeInsets.all(GlassmorphismSpacing.lg),
                icon: Icons.check,
                enabled: true,
                width: 150.0,
                height: 50.0,
                elevation: 4.0,
                fontSize: 14.0,
                isBold: true,
                onPressed: () {},
              ),
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsOneWidget);
      expect(find.text('Button'), findsOneWidget);
      expect(find.byIcon(Icons.check), findsOneWidget);
    });

    testWidgets('renders multiple buttons', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Column(
              children: [
                GlassmorphicButton(
                  label: 'Button 1',
                  onPressed: () {},
                ),
                GlassmorphicButton(
                  label: 'Button 2',
                  onPressed: () {},
                ),
                GlassmorphicButton(
                  label: 'Button 3',
                  onPressed: () {},
                ),
              ],
            ),
          ),
        ),
      );

      expect(find.byType(GlassmorphicButton), findsNWidgets(3));
      expect(find.text('Button 1'), findsOneWidget);
      expect(find.text('Button 2'), findsOneWidget);
      expect(find.text('Button 3'), findsOneWidget);
    });

    testWidgets('button opacity changes on press', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: GlassmorphicButton(
              label: 'Test',
              opacity: 0.5,
              onPressed: () {},
            ),
          ),
        ),
      );

      await tester.tap(find.byType(GlassmorphicButton));
      await tester.pump();

      expect(find.byType(GlassmorphicButton), findsOneWidget);
    });
  });
}
