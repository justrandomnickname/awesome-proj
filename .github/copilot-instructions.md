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
-   **Hierarchical World Structure**: Cluster → SubCluster → Point architecture
-   **NPC Generation**: Race-based system with temper traits and special traits
-   **Point-based Navigation**: Movement between connected points with interaction history
-   **AI Content Generation**: Dynamic descriptions for locations and NPCs via prompt service
-   **Save/Load System**: Complete JSON-based game state persistence with hierarchy support
-   **Interactive History**: Point-based interaction tracking and state management
-   **Trait System**: Weighted trait generation system with trait name resolution
-   **Map Visualization**: Metro-style interactive map display
-   **Automatic Save Recovery**: Fallback to new game generation when saves are corrupted/missing

### 🏗️ Architecture Structure

```
app/
├── app.go                                    # Wails app context - exposes all frontend methods
├── main.go                                   # Application entry point
├── game/
│   ├── engine.go                            # GameEngine - central orchestrator with all game logic
│   └── game_state.go                       # GameState + PointState management
├── domain/
│   ├── entities/                            # Core game entities
│   │   ├── location_hierarchy.go           # Cluster, SubCluster, Point structures
│   │   ├── location.go                     # Legacy location entities
│   │   ├── npc.go                          # NPC entity with traits and temper system
│   │   ├── interaction.go                  # Player interaction history tracking
│   │   ├── temper_traits.go               # Character temperament system (-10 to +10)
│   │   └── save.go                         # Save file structure and metadata
│   ├── aggregates/
│   │   └── world.go                        # World aggregate containing full hierarchy + NPCs
│   └── services/
│       ├── world_generation_service.go     # World creation orchestrator
│       ├── save_service.go                 # JSON persistence and save management
│       ├── location_interaction_service.go # Player action handling
│       └── prompts/
│           └── subcluster_static_description_service.go  # AI prompt generation for content
└── infrastructure/
    └── builders/
        ├── location_builder.go             # Hierarchical world generation with connections
        ├── npc_builder.go                  # NPC creation with race logic and traits
        ├── trait_system.go                 # Weighted trait generation and name resolution
        └── static/
            ├── npc.go                      # Static NPC data pools (names, descriptions)
            └── taverns.go                  # Static tavern generation data

frontend/
├── src/
│   ├── App.svelte                          # Main application component
│   └── lib/
│       ├── LocationScreen.svelte           # Main game interface with point navigation
│       ├── MapViewer.svelte                # Interactive hierarchical world map
│       └── SaveMenu.svelte                 # Save/load game interface
└── wailsjs/                                # Auto-generated Wails TypeScript bindings
```

### 📊 Current Data Flow

1. **World Generation**: `GameEngine.NewGameEngine()` → `WorldGenerationService.GenerateWorld()` → `LocationBuilder.GenerateLocationHierarchyWithNPCs()`
2. **NPC Creation**: `NPCBuilder.GenerateNPCsForPoint()` → generates race-based NPCs → applies `TemperTraits` → assigns special traits via `TraitSystem`
3. **Save System**: `SaveService` → serializes complete `World` aggregate with hierarchy and point states → handles save recovery
4. **AI Content**: `SubclusterStaticDescriptionService.GenerateSubclusterDescriptionPrompt()` → collects all NPCs in subcluster → generates rich prompts with trait data
5. **Navigation**: `GetCurrentPoint()` → returns Point with NPCs and connections → `MoveToPoint()` → updates game state
6. **Interaction**: Player actions → `LocationInteractionService` → creates `Interaction` records → stored in `PointState`

### 🎲 NPC Generation System

**Race Distribution**: 40% human, 30% dwarf, 30% skaven
**Names**: Race-specific name pools in `builders/static/npc.go`
**Temper Traits**: 6 characteristics (Prudence, Emotionality, Independence, Optimism, Flexibility, Aggressiveness) with values -10 to +10
**Special Traits**: Weighted system with trait categories:

-   Nature Lover, Bookworm, Experienced Warrior, etc.
-   Trait name resolution via `TraitSystem.GetTraitNamesByIDs()`

### 🗺️ World Structure

**Hierarchy**: 3-level structure

-   **Cluster**: Large world areas (forest, city, dungeon)
-   **SubCluster**: Districts within clusters (tavern district, heart of forest)
-   **Point**: Specific locations (tavern hall, forest clearing)

**Connections**: Points connected via `connections` array, enabling navigation
**Entry Points**: Each SubCluster has designated entry points for initial access

### 💾 Save System

**Format**: JSON files in `saves/` directory
**Structure**: Complete world state + game state + point interaction history
**Recovery**: Automatic new game generation when saves are corrupted or missing
**Validation**: Save compatibility checking against current world structure

## Key Architecture Principles

-   **Domain-Driven Design**: Clear separation of concerns across layers
-   **Game Engine**: Central orchestrator pattern - handles all game logic and state
-   **Point-Based Architecture**: Pure hierarchical navigation system (Cluster → SubCluster → Point)
-   **Trait-Driven NPCs**: Rich character system with temper traits and special abilities
-   **JSON Persistence**: File-based saves with complete world state serialization
-   **AI Content Integration**: Dynamic description generation with trait-aware prompts
-   **Robust Save Recovery**: Automatic fallback systems prevent data loss

## Development Guidelines

### Code Style

-   **Go**: Follow standard Go conventions, use interfaces for testability
-   **TypeScript**: Strict typing for all frontend code
-   **Svelte**: Component-based architecture with clear prop interfaces
-   **No Magic Numbers**: Use constants and configuration objects

### Architecture Rules

-   **No Circular Dependencies**: Strict layer separation (entities → services → builders)
-   **Engine as Orchestrator**: All game logic flows through GameEngine
-   **Builder Pattern**: Use builders for complex object creation (locations, NPCs)
-   **Interface Segregation**: Small, focused interfaces for trait systems
-   **Save Compatibility**: Always validate saves against current world structure

### Current Implementation Details

-   **World Generation**: Seed-based deterministic generation with entry point detection
-   **NPC Traits**: ID-based trait system with human-readable name resolution
-   **Point Navigation**: Connection-based movement with interaction history
-   **AI Prompts**: Comprehensive context generation including all subcluster NPCs
-   **Error Handling**: Graceful degradation with automatic new game fallbacks

## Recent Major Changes (September 2025)

-   **Enhanced Trait System**: Fixed trait name display, added comprehensive trait resolution
-   **Save System Overhaul**: Added validation, corruption recovery, and automatic new game generation
-   **NPC Data Completeness**: All NPCs in subcluster now included in AI prompts
-   **Game State Validation**: Entry point detection and compatibility checking
-   **Interaction History**: Full point-based interaction tracking system

## Future Plans (Do Not Implement Yet)
