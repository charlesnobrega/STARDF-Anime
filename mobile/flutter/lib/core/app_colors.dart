import 'package:flutter/material.dart';

class AppColors {
  static const Color background = Color(0xFF0F0F13);
  static const Color surface = Color(0xFF1C1C23);
  static const Color primary = Color(0xFF6C5CE7);
  static const Color accent = Color(0xFF00D2D3);
  static const Color textMain = Color(0xFFFFFFFF);
  static const Color textMuted = Color(0xFFA0A0A0);
  static const Color glassBorder = Color(0x33FFFFFF);
  static const Color glassBackground = Color(0x1AFFFFFF);
  
  static const Color error = Color(0xFFFF7675);
  static const Color success = Color(0xFF55EFC4);
  static const Color warning = Color(0xFFFFEAA7);

  static LinearGradient primaryGradient = const LinearGradient(
    colors: [primary, Color(0xFFA29BFE)],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}
