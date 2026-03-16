# StarDF-Anime Mobile — Implementation Plan

This document outlines the strategy for developing the mobile version of StarDF-Anime.

## Objective
Provide a unified anime/streaming experience across Desktop (TUI) and Mobile (Android), leveraging the existing scraping core.

## 📱 Technology Stack
- **Framework**: **Flutter** (Targeting Android 8.0+)
- **Core Strategy**: 
    - Export `pkg/stardf` as a shared library.
    - Use `gomobile` to generate bindings for Android (AAR).
    - Flutter FFI or Method Channels to communicate with the Go core.
- **UI Architecture**: BLoC or Provider pattern for state management.
- **Local Storage**: SQLite (drift/sqflite) synced with the Go core's logic.

## 🎨 Design Concept
The mobile app will follow a **Dark Glassmorphism** aesthetic, matching the premium feel of the TUI.

![Mobile App Concept](file:///C:/Users/charles.nobrega/.gemini/antigravity/brain/c740cc9e-bd2a-46bd-93c8-0d8264b46936/mockup_mobile_1773689068952.png)

## 🛤 Road Map
1. **Phase 1: Bridge**
   - Verify `pkg/stardf` compatibility with `gomobile`.
   - Create a basic Flutter scaffold.
2. **Phase 2: UI/UX**
   - Implement the search and discovery screens.
   - Integration with native video players.
3. **Phase 3: Synchronization**
   - AniList real-time sync.
   - Watchlist state management.

> [!TIP]
> Reusing the `pkg/stardf` library ensures that scrapers are maintained in a single place, reducing technical debt.
