import 'dart:ui';
import 'package:flutter/material.dart';
import '../theme/glassmorphism_theme.dart';
import 'glassmorphic_container.dart';

/// A button widget that implements the glassmorphism design pattern.
/// 
/// The GlassmorphicButton provides a frosted glass effect button with
/// customizable blur, opacity, and color properties.
class GlassmorphicButton extends StatefulWidget {
  /// The text to display on the button
  final String label;

  /// The callback when the button is pressed
  final VoidCallback onPressed;

  /// The blur amount for the glassmorphism effect (default: 10.0)
  final double blur;

  /// The opacity of the glass effect (default: 0.5)
  final double opacity;

  /// The background color of the glass (default: white)
  final Color color;

  /// The text color (default: black)
  final Color textColor;

  /// The border radius of the button (default: 12.0)
  final double borderRadius;

  /// The padding inside the button (default: 12.0 horizontal, 8.0 vertical)
  final EdgeInsets padding;

  /// Optional icon to display before the text
  final IconData? icon;

  /// Whether the button is enabled (default: true)
  final bool enabled;

  /// The width of the button (default: null - auto)
  final double? width;

  /// The height of the button (default: 48.0)
  final double height;

  /// Optional shadow elevation (default: 0.0)
  final double elevation;

  /// The font size of the label (default: 16.0)
  final double fontSize;

  /// Whether the label should be bold (default: false)
  final bool isBold;

  const GlassmorphicButton({
    Key? key,
    required this.label,
    required this.onPressed,
    this.blur = GlassmorphismBlur.medium,
    this.opacity = GlassmorphismOpacity.medium,
    this.color = GlassmorphismLightColors.primaryGlass,
    this.textColor = GlassmorphismLightColors.textPrimary,
    this.borderRadius = GlassmorphismRadius.medium,
    this.padding = const EdgeInsets.symmetric(
      horizontal: GlassmorphismSpacing.lg,
      vertical: GlassmorphismSpacing.md,
    ),
    this.icon,
    this.enabled = true,
    this.width,
    this.height = 48.0,
    this.elevation = 0.0,
    this.fontSize = 16.0,
    this.isBold = false,
  }) : super(key: key);

  @override
  State<GlassmorphicButton> createState() => _GlassmorphicButtonState();
}

class _GlassmorphicButtonState extends State<GlassmorphicButton> {
  bool _isPressed = false;

  @override
  Widget build(BuildContext context) {
    final effectiveOpacity = _isPressed ? widget.opacity * 0.8 : widget.opacity;

    Widget buttonContent = Row(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        if (widget.icon != null) ...[
          Icon(
            widget.icon,
            color: widget.textColor,
            size: widget.fontSize,
          ),
          const SizedBox(width: GlassmorphismSpacing.md),
        ],
        Text(
          widget.label,
          style: TextStyle(
            color: widget.textColor,
            fontSize: widget.fontSize,
            fontWeight: widget.isBold ? FontWeight.bold : FontWeight.normal,
          ),
        ),
      ],
    );

    return GestureDetector(
      onTapDown: widget.enabled ? (_) => setState(() => _isPressed = true) : null,
      onTapUp: widget.enabled
          ? (_) {
              setState(() => _isPressed = false);
              widget.onPressed();
            }
          : null,
      onTapCancel: widget.enabled ? () => setState(() => _isPressed = false) : null,
      child: SizedBox(
        width: widget.width,
        height: widget.height,
        child: GlassmorphicContainer(
          blur: widget.blur,
          opacity: effectiveOpacity,
          color: widget.color,
          borderRadius: widget.borderRadius,
          padding: widget.padding,
          elevation: widget.elevation,
          child: buttonContent,
        ),
      ),
    );
  }
}
