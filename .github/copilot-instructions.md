# Copilot Instructions

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Overview

This is a text-based RPG game built with:

-   **Backend**: Go + Wails v2.10.1 for the desktop application framework
-   **Frontend**: Svelte + TypeScript for the user interface
-   **Architecture**: Clean DDD architecture with Game Engine as central orchestrator

## Current Implementation Status

### ✅ Implemented Features

-   **Game Engine**: Central orchestrator in `app/game/engine.go`
-   **NPC Generation**: Race-based system (human, dwarf, skaven) with weighted probabilities
-   **Location Generation**: Single starting location with random attributes
-   **Frontend Display**: LocationScreen.svelte shows location info and NPC cards
-   **Simplified Architecture**: Removed movement system, levels, emoji for minimalism

### 🏗️ Architecture Structure

```
app/
├── app.go                          # Wails app context (exposes GetCurrentLocation)
├── game/
│   ├── engine.go                   # GameEngine - central orchestrator
│   └── game_state.go              # GameState + LocationInfo/NPCInfo structs
├── domain/
│   ├── entities/                   # Location, NPC entities (no levels/emoji)
│   ├── aggregates/                 # World aggregate
│   └── services/                   # WorldGenerationService
└── infrastructure/
    └── builders/                   # LocationBuilder, NPCBuilder

frontend/
├── src/
│   ├── App.svelte                 # Main component
│   └── lib/
│       └── LocationScreen.svelte  # Shows location + NPC cards (simplified)
└── wailsjs/                       # Auto-generated Wails bindings
```

### 📊 Current Data Flow

1. `GameEngine.NewGameEngine()` → creates world via `WorldGenerationService`
2. `LocationBuilder.GenerateRandomLocations()` → creates single "start" location
3. `NPCBuilder.GenerateNPCsForLocation()` → creates 2-5 NPCs with race logic
4. Frontend calls `GetCurrentLocation()` → returns LocationInfo with NPCs
5. `LocationScreen.svelte` displays location name, description, NPC cards

### 🎲 NPC Generation Logic

-   **Race Distribution**: 40% human, 30% dwarf, 30% skaven
-   **Names**: Race-specific name pools (Артур/Торин/Скритч etc.)
-   **Descriptions**: Race-based templates (мудрый старик/крепкий воин/подозрительный крыс)
-   **Fields**: Only ID, Name, Race, Description (no levels/emoji removed)

## Key Architecture Principles

-   **Domain-Driven Design**: Clear separation of concerns across layers
-   **Game Engine**: Central orchestrator pattern - handles all game state
-   **No Premature Optimization**: Only implement features when actually needed
-   **JSON Persistence**: File-based saves (no database for core game data)
-   **Minimal UI**: Clean, simple interface without unnecessary visual elements

## Development Guidelines

### Code Style

-   **Go**: Follow standard Go conventions, use interfaces for testability
-   **TypeScript**: Strict typing for all frontend code
-   **Svelte**: Component-based architecture with clear prop interfaces
-   **No Magic Numbers**: Use constants and configuration objects

### Architecture Rules

-   **No Circular Dependencies**: Strict layer separation
-   **Engine as Orchestrator**: All game logic flows through GameEngine
-   **Builder Pattern**: Use builders for complex object creation
-   **Interface Segregation**: Small, focused interfaces

### Current Simplifications

-   **No Movement System**: Single location only (removed exits/navigation)
-   **No Character Levels**: NPCs have no gameplay stats
-   **No Visual Emoji**: Clean text-only interface
-   **No Future-Proofing**: Add features only when needed

## Future Plans (Do Not Implement Yet)

-   **location_state**: AI-generated location descriptions based on player actions
-   **Player Interactions**: Dialogue system with NPCs
-   **Save/Load System**: JSON-based game state persistence
-   **Location States**: Dynamic location descriptions
