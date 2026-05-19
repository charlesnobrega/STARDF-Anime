import 'package:flutter/material.dart';
import 'responsive_breakpoints.dart';

/// A widget that adapts its layout based on screen size and orientation
class ResponsiveLayout extends StatelessWidget {
  /// Widget to display on mobile screens
  final Widget mobile;

  /// Widget to display on tablet screens (optional)
  final Widget? tablet;

  /// Widget to display on desktop screens (optional)
  final Widget? desktop;

  /// Widget to display on extra large screens (optional)
  final Widget? extraLarge;

  const ResponsiveLayout({
    Key? key,
    required this.mobile,
    this.tablet,
    this.desktop,
    this.extraLarge,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final width = ResponsiveHelper.getWidth(context);
    final screenSize = ResponsiveHelper.getScreenSize(width);

    switch (screenSize) {
      case ResponsiveScreenSize.mobile:
        return mobile;
      case ResponsiveScreenSize.tablet:
        return tablet ?? mobile;
      case ResponsiveScreenSize.desktop:
        return desktop ?? tablet ?? mobile;
      case ResponsiveScreenSize.extraLarge:
        return extraLarge ?? desktop ?? tablet ?? mobile;
    }
  }
}

/// A widget that adapts its layout based on orientation
class OrientationLayout extends StatelessWidget {
  /// Widget to display in portrait orientation
  final Widget portrait;

  /// Widget to display in landscape orientation
  final Widget landscape;

  const OrientationLayout({
    Key? key,
    required this.portrait,
    required this.landscape,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ResponsiveHelper.isPortrait(context) ? portrait : landscape;
  }
}

/// A widget that provides responsive padding based on screen size
class ResponsivePadding extends StatelessWidget {
  /// The child widget
  final Widget child;

  /// Custom padding for mobile (optional)
  final EdgeInsets? mobilePadding;

  /// Custom padding for tablet (optional)
  final EdgeInsets? tabletPadding;

  /// Custom padding for desktop (optional)
  final EdgeInsets? desktopPadding;

  const ResponsivePadding({
    Key? key,
    required this.child,
    this.mobilePadding,
    this.tabletPadding,
    this.desktopPadding,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final width = ResponsiveHelper.getWidth(context);
    final screenSize = ResponsiveHelper.getScreenSize(width);

    EdgeInsets padding;
    switch (screenSize) {
      case ResponsiveScreenSize.mobile:
        padding = mobilePadding ?? const EdgeInsets.all(16.0);
        break;
      case ResponsiveScreenSize.tablet:
        padding = tabletPadding ?? const EdgeInsets.all(24.0);
        break;
      case ResponsiveScreenSize.desktop:
      case ResponsiveScreenSize.extraLarge:
        padding = desktopPadding ?? const EdgeInsets.all(32.0);
        break;
    }

    return Padding(
      padding: padding,
      child: child,
    );
  }
}

/// A widget that provides responsive grid layout
class ResponsiveGrid extends StatelessWidget {
  /// The items to display in the grid
  final List<Widget> children;

  /// Custom column count for mobile (optional)
  final int? mobileColumns;

  /// Custom column count for tablet (optional)
  final int? tabletColumns;

  /// Custom column count for desktop (optional)
  final int? desktopColumns;

  /// The spacing between items
  final double spacing;

  /// The aspect ratio of each grid item
  final double childAspectRatio;

  const ResponsiveGrid({
    Key? key,
    required this.children,
    this.mobileColumns,
    this.tabletColumns,
    this.desktopColumns,
    this.spacing = 16.0,
    this.childAspectRatio = 1.0,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final width = ResponsiveHelper.getWidth(context);
    final screenSize = ResponsiveHelper.getScreenSize(width);

    int columns;
    switch (screenSize) {
      case ResponsiveScreenSize.mobile:
        columns = mobileColumns ?? 2;
        break;
      case ResponsiveScreenSize.tablet:
        columns = tabletColumns ?? 3;
        break;
      case ResponsiveScreenSize.desktop:
      case ResponsiveScreenSize.extraLarge:
        columns = desktopColumns ?? 4;
        break;
    }

    return GridView.count(
      crossAxisCount: columns,
      mainAxisSpacing: spacing,
      crossAxisSpacing: spacing,
      childAspectRatio: childAspectRatio,
      children: children,
    );
  }
}

/// A widget that provides responsive column layout
class ResponsiveColumn extends StatelessWidget {
  /// The children widgets
  final List<Widget> children;

  /// The spacing between children
  final double spacing;

  /// The cross axis alignment
  final CrossAxisAlignment crossAxisAlignment;

  /// The main axis alignment
  final MainAxisAlignment mainAxisAlignment;

  const ResponsiveColumn({
    Key? key,
    required this.children,
    this.spacing = 16.0,
    this.crossAxisAlignment = CrossAxisAlignment.start,
    this.mainAxisAlignment = MainAxisAlignment.start,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: crossAxisAlignment,
      mainAxisAlignment: mainAxisAlignment,
      children: [
        for (int i = 0; i < children.length; i++) ...[
          children[i],
          if (i < children.length - 1) SizedBox(height: spacing),
        ],
      ],
    );
  }
}

/// A widget that provides responsive row layout
class ResponsiveRow extends StatelessWidget {
  /// The children widgets
  final List<Widget> children;

  /// The spacing between children
  final double spacing;

  /// The cross axis alignment
  final CrossAxisAlignment crossAxisAlignment;

  /// The main axis alignment
  final MainAxisAlignment mainAxisAlignment;

  const ResponsiveRow({
    Key? key,
    required this.children,
    this.spacing = 16.0,
    this.crossAxisAlignment = CrossAxisAlignment.center,
    this.mainAxisAlignment = MainAxisAlignment.start,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: crossAxisAlignment,
      mainAxisAlignment: mainAxisAlignment,
      children: [
        for (int i = 0; i < children.length; i++) ...[
          Expanded(child: children[i]),
          if (i < children.length - 1) SizedBox(width: spacing),
        ],
      ],
    );
  }
}
