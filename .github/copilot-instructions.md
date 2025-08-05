# Copilot Instructions

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Overview

This is a text-based RPG game built with:

- **Backend**: Go + Wails for the desktop application framework
- **Frontend**: Svelte + TypeScript for the user interface
- **Architecture**: Game Engine with layered architecture to avoid circular dependencies

## Key Architecture Principles

- Use domain-driven design with clear separation of concerns
- Game Engine acts as the central orchestrator
- JSON files for game state persistence (no database for core game data)
- AI integration for dynamic content generation (NPCs, locations, descriptions)

## Code Style Guidelines

- Use TypeScript for all frontend code
- Follow Go best practices for backend code
- Maintain clean architecture with proper layer separation
- Use interfaces for dependency injection and testing
