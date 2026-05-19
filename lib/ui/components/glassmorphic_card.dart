import 'dart:ui';
import 'package:flutter/material.dart';
import '../theme/glassmorphism_theme.dart';
import 'glassmorphic_container.dart';

/// A card widget that implements the glassmorphism design pattern.
/// 
/// The GlassmorphicCard provides a frosted glass effect card with
/// optional header, title, and footer sections.
class GlassmorphicCard extends StatelessWidget {
  /// The main content of the card
  final Widget child;

  /// Optional title to display at the top of the card
  final String? title;

  /// Optional subtitle to display below the title
  final String? subtitle;

  /// Optional header widget (overrides title/subtitle if provided)
  final Widget? header;

  /// Optional footer widget
  final Widget? footer;

  /// The blur amount for the glassmorphism effect (default: 10.0)
  final double blur;

  /// The opacity of the glass effect (default: 0.5)
  final double opacity;

  /// The background color of the glass (default: white)
  final Color color;

  /// The border radius of the card (default: 20.0)
  final double borderRadius;

  /// The padding inside the card (default: 16.0)
  final EdgeInsets padding;

  /// The margin around the card (default: 0.0)
  final EdgeInsets margin;

  /// Optional shadow elevation (default: 0.0)
  final double elevation;

  /// Optional callback when the card is tapped
  final VoidCallback? onTap;

  /// Whether the card should be clickable (default: false)
  final bool isClickable;

  /// The text color for title and subtitle (default: black)
  final Color textColor;

  /// The text color for subtitle (default: gray)
  final Color subtitleColor;

  const GlassmorphicCard({
    Key? key,
    required this.child,
    this.title,
    this.subtitle,
    this.header,
    this.footer,
    this.blur = GlassmorphismBlur.medium,
    this.opacity = GlassmorphismOpacity.medium,
    this.color = GlassmorphismLightColors.primaryGlass,
    this.borderRadius = GlassmorphismRadius.large,
    this.padding = const EdgeInsets.all(GlassmorphismSpacing.lg),
    this.margin = EdgeInsets.zero,
    this.elevation = 0.0,
    this.onTap,
    this.isClickable = false,
    this.textColor = GlassmorphismLightColors.textPrimary,
    this.subtitleColor = GlassmorphismLightColors.textSecondary,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GlassmorphicContainer(
      blur: blur,
      opacity: opacity,
      color: color,
      borderRadius: borderRadius,
      padding: EdgeInsets.zero,
      margin: margin,
      elevation: elevation,
      onTap: onTap,
      isClickable: isClickable,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header section
          if (header != null)
            Padding(
              padding: padding,
              child: header!,
            )
          else if (title != null)
            Padding(
              padding: padding,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title!,
                    style: TextStyle(
                      color: textColor,
                      fontSize: 18.0,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  if (subtitle != null) ...[
                    const SizedBox(height: GlassmorphismSpacing.sm),
                    Text(
                      subtitle!,
                      style: TextStyle(
                        color: subtitleColor,
                        fontSize: 14.0,
                      ),
                    ),
                  ],
                ],
              ),
            ),
          // Content section
          Padding(
            padding: header != null || title != null
                ? padding.copyWith(top: 0)
                : padding,
            child: child,
          ),
          // Footer section
          if (footer != null)
            Padding(
              padding: padding.copyWith(top: 0),
              child: footer!,
            ),
        ],
      ),
    );
  }
}
