import 'dart:ui';
import 'package:flutter/material.dart';
import '../theme/glassmorphism_theme.dart';

/// A widget that implements the glassmorphism design pattern with blur and transparency effects.
/// 
/// The GlassmorphicContainer provides a frosted glass effect using BackdropFilter
/// and customizable blur, opacity, and color properties.
class GlassmorphicContainer extends StatelessWidget {
  /// The child widget to display inside the container
  final Widget child;

  /// The blur amount for the glassmorphism effect (default: 10.0)
  final double blur;

  /// The opacity of the glass effect (default: 0.5)
  final double opacity;

  /// The background color of the glass (default: white)
  final Color color;

  /// The border radius of the container (default: 20.0)
  final double borderRadius;

  /// Whether to show a border (default: true)
  final bool showBorder;

  /// The border color (default: white with 0.2 opacity)
  final Color? borderColor;

  /// The padding inside the container (default: 16.0)
  final EdgeInsets padding;

  /// The margin around the container (default: 0.0)
  final EdgeInsets margin;

  /// Optional shadow elevation (default: 0.0)
  final double elevation;

  /// Optional callback when the container is tapped
  final VoidCallback? onTap;

  /// Whether the container should be clickable (default: false)
  final bool isClickable;

  const GlassmorphicContainer({
    Key? key,
    required this.child,
    this.blur = GlassmorphismBlur.medium,
    this.opacity = GlassmorphismOpacity.medium,
    this.color = GlassmorphismLightColors.primaryGlass,
    this.borderRadius = GlassmorphismRadius.large,
    this.showBorder = true,
    this.borderColor,
    this.padding = const EdgeInsets.all(GlassmorphismSpacing.lg),
    this.margin = EdgeInsets.zero,
    this.elevation = 0.0,
    this.onTap,
    this.isClickable = false,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final effectiveBorderColor = borderColor ?? Colors.white.withOpacity(0.2);

    Widget container = BackdropFilter(
      filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
      child: Container(
        decoration: BoxDecoration(
          color: color.withOpacity(opacity),
          borderRadius: BorderRadius.circular(borderRadius),
          border: showBorder
              ? Border.all(
                  color: effectiveBorderColor,
                  width: 1.0,
                )
              : null,
        ),
        padding: padding,
        child: child,
      ),
    );

    if (elevation > 0) {
      container = Material(
        elevation: elevation,
        borderRadius: BorderRadius.circular(borderRadius),
        child: container,
      );
    }

    if (isClickable && onTap != null) {
      container = GestureDetector(
        onTap: onTap,
        child: container,
      );
    }

    if (margin != EdgeInsets.zero) {
      container = Padding(
        padding: margin,
        child: container,
      );
    }

    return container;
  }
}
