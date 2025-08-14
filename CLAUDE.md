# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sultengutt is a cross-platform desktop reminder application for ordering surprise dinners. It uses Go with the Cobra CLI framework for commands and Fyne for GUI popups. The application schedules periodic reminders based on user-configured days and times.

## Architecture

The application follows a modular internal package structure:
- **cmd/main.go**: CLI entry point with Cobra commands (install, execute, pause, resume, status)
- **internal/config**: Configuration management with JSON persistence in ~/.sultengutt/
- **internal/installer**: Interactive TUI installer using Charm's Huh framework
- **internal/popup**: Cross-platform GUI popup reminders using Fyne
- **internal/scheduler**: OS-specific task scheduling (Windows Task Scheduler, macOS launchd planned)
- **internal/utils**: Shared utilities for duration parsing and executable resolution

## Development Commands

```bash
# Build the application
go build -o sultengutt cmd/main.go

# Run the application
go run cmd/main.go [command]

# Install dependencies
go mod download
go mod tidy

# Test (no tests currently exist)
go test ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Key Implementation Notes

1. **Platform Support**: Currently Windows-only scheduler. The Mac scheduler (internal/scheduler/mac.go) is incomplete and has compilation errors.

2. **Configuration**: Stored in `~/.sultengutt/sultengutt.json` with pause state tracking.

3. **Scheduling**: Uses native OS task schedulers - Windows Task Scheduler via schtasks command.

4. **State Management**: 
   - PausedUntil: -1 (not paused), 0 (paused indefinitely), >0 (unix timestamp)
   - Fresh install detection via isFreshInstall flag

5. **Dependencies**: Heavy use of Charm libraries (Huh, Lipgloss) for TUI and Fyne for GUI popups.