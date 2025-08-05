# Awesome Proj - Text-based RPG Game

## About

A text-based RPG game with AI-powered content generation built using:

- **Backend**: Go + Wails for desktop application framework
- **Frontend**: Svelte + TypeScript for the user interface
- **AI Integration**: Dynamic content generation for NPCs, locations, and descriptions

## Architecture

The project follows a clean layered architecture:

- **Game Engine**: Central orchestrator for game logic
- **Domain Layer**: Core game entities (Location, NPC, Player)
- **Service Layer**: Business logic and AI integration
- **Storage Layer**: JSON-based persistence (no database for core game data)

## Development

### Live Development

To run in live development mode:

```bash
wails dev
```

This will start:

- Vite development server with hot reload on `http://localhost:5173/`
- Wails DevServer on `http://localhost:34115` for browser development with Go method access

### Building

To build a redistributable, production mode package:

```bash
wails build
```

## Project Structure

```
awesome-proj/
├── main.go              # Application entry point
├── app.go               # Main application context
├── frontend/            # Svelte + TypeScript frontend
│   ├── src/
│   ├── package.json
│   └── vite.config.ts
├── internal/            # Go backend (future game logic)
└── data/               # JSON game data (future saves/templates)
```
