import 'package:flutter/material.dart';

/// Responsive breakpoints for different screen sizes
class ResponsiveBreakpoints {
  /// Mobile breakpoint: < 600dp
  static const double mobile = 600.0;

  /// Tablet breakpoint: >= 600dp
  static const double tablet = 600.0;

  /// Desktop breakpoint: >= 1200dp
  static const double desktop = 1200.0;

  /// Extra large breakpoint: >= 1920dp
  static const double extraLarge = 1920.0;
}

/// Responsive orientation helper
enum ResponsiveOrientation {
  portrait,
  landscape,
}

/// Responsive screen size helper
enum ResponsiveScreenSize {
  mobile,
  tablet,
  desktop,
  extraLarge,
}

/// Helper class for responsive design
class ResponsiveHelper {
  /// Get the current screen size based on width
  static ResponsiveScreenSize getScreenSize(double width) {
    if (width < ResponsiveBreakpoints.mobile) {
      return ResponsiveScreenSize.mobile;
    } else if (width < ResponsiveBreakpoints.desktop) {
      return ResponsiveScreenSize.tablet;
    } else if (width < ResponsiveBreakpoints.extraLarge) {
      return ResponsiveScreenSize.desktop;
    } else {
      return ResponsiveScreenSize.extraLarge;
    }
  }

  /// Check if the screen is mobile
  static bool isMobile(double width) {
    return width < ResponsiveBreakpoints.mobile;
  }

  /// Check if the screen is tablet
  static bool isTablet(double width) {
    return width >= ResponsiveBreakpoints.tablet &&
        width < ResponsiveBreakpoints.desktop;
  }

  /// Check if the screen is desktop
  static bool isDesktop(double width) {
    return width >= ResponsiveBreakpoints.desktop;
  }

  /// Get the current orientation
  static ResponsiveOrientation getOrientation(BuildContext context) {
    return MediaQuery.of(context).orientation == Orientation.portrait
        ? ResponsiveOrientation.portrait
        : ResponsiveOrientation.landscape;
  }

  /// Check if the screen is in portrait orientation
  static bool isPortrait(BuildContext context) {
    return MediaQuery.of(context).orientation == Orientation.portrait;
  }

  /// Check if the screen is in landscape orientation
  static bool isLandscape(BuildContext context) {
    return MediaQuery.of(context).orientation == Orientation.landscape;
  }

  /// Get the screen width
  static double getWidth(BuildContext context) {
    return MediaQuery.of(context).size.width;
  }

  /// Get the screen height
  static double getHeight(BuildContext context) {
    return MediaQuery.of(context).size.height;
  }

  /// Get the device pixel ratio
  static double getDevicePixelRatio(BuildContext context) {
    return MediaQuery.of(context).devicePixelRatio;
  }

  /// Get the padding for the current screen size
  static EdgeInsets getPadding(BuildContext context) {
    final width = getWidth(context);
    if (isMobile(width)) {
      return const EdgeInsets.all(16.0);
    } else if (isTablet(width)) {
      return const EdgeInsets.all(24.0);
    } else {
      return const EdgeInsets.all(32.0);
    }
  }

  /// Get the horizontal padding for the current screen size
  static double getHorizontalPadding(BuildContext context) {
    final width = getWidth(context);
    if (isMobile(width)) {
      return 16.0;
    } else if (isTablet(width)) {
      return 24.0;
    } else {
      return 32.0;
    }
  }

  /// Get the vertical padding for the current screen size
  static double getVerticalPadding(BuildContext context) {
    final width = getWidth(context);
    if (isMobile(width)) {
      return 12.0;
    } else if (isTablet(width)) {
      return 16.0;
    } else {
      return 20.0;
    }
  }

  /// Get the grid column count for the current screen size
  static int getGridColumns(BuildContext context) {
    final width = getWidth(context);
    if (isMobile(width)) {
      return 2;
    } else if (isTablet(width)) {
      return 3;
    } else {
      return 4;
    }
  }

  /// Get the font size for the current screen size
  static double getFontSize(
    BuildContext context, {
    required double mobileSize,
    required double tabletSize,
    required double desktopSize,
  }) {
    final width = getWidth(context);
    if (isMobile(width)) {
      return mobileSize;
    } else if (isTablet(width)) {
      return tabletSize;
    } else {
      return desktopSize;
    }
  }

  /// Get the spacing for the current screen size
  static double getSpacing(
    BuildContext context, {
    required double mobileSpacing,
    required double tabletSpacing,
    required double desktopSpacing,
  }) {
    final width = getWidth(context);
    if (isMobile(width)) {
      return mobileSpacing;
    } else if (isTablet(width)) {
      return tabletSpacing;
    } else {
      return desktopSpacing;
    }
  }
}
