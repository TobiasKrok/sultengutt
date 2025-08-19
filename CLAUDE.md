# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sultengutt is a cross-platform desktop reminder application for ordering surprise dinners. It uses Go with the Cobra CLI framework for commands and platform-specific GUI popups. The application schedules periodic reminders based on user-configured days and times, with support for dynamic mantras loaded from external config files.

## Architecture

The application follows a modular internal package structure:
- **cmd/main.go**: CLI entry point with Cobra commands (install, execute, pause, resume, status)
- **internal/config**: Configuration management with JSON persistence in ~/.sultengutt/
- **internal/installer**: Interactive TUI installer using Charm's Huh framework
- **internal/popup**: Cross-platform GUI popup reminders
  - macOS/Linux: Modern Fyne-based GUI with custom theme and typography
  - Windows: PowerShell-generated WPF/Windows Forms popups with modern dark styling
- **internal/mantras**: Dynamic mantra loading from config/mantras.json with fallback locations
- **internal/scheduler**: OS-specific task scheduling (Windows Task Scheduler, macOS launchd)
- **internal/utils**: Enhanced utilities for duration parsing (supports minutes, hours, days, weeks, months) and executable resolution

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

1. **Platform Support**: Full cross-platform support with OS-specific implementations:
   - **Windows**: Task Scheduler via schtasks command, PowerShell-based modern popups
   - **macOS**: launchd via plist files, Fyne-based GUI with custom themes
   - **Linux**: Planned (currently uses macOS implementation)

2. **Configuration**: 
   - Main config: `~/.sultengutt/sultengutt.json` with pause state tracking
   - Mantras: `~/.sultengutt/mantras.json` (auto-copied from config/ during install)
   - Package manager friendly: supports Homebrew (/opt/homebrew/etc/), Scoop (persist), system-wide locations

3. **Popup Design Philosophy**:
   - **macOS/Linux**: Modern Fyne-based with custom themes, larger fonts, white text on dark backgrounds
   - **Windows**: WPF-first with Windows Forms fallback, modern dark styling, proper typography
   - **Mantras**: Displays "Your mantra for today" header with random mantra from config file

4. **State Management**: 
   - PausedUntil: -1 (not paused), 0 (paused indefinitely), >0 (unix timestamp)
   - Fresh install detection via isFreshInstall flag
   - Dynamic mantra loading with graceful fallbacks

5. **Duration Parsing**: Simplified `utils.ParseDuration(interface{})` function:
   - Single regex handles all formats: `30m`, `2h`, `1d`, `30 minutes`, `2 hours`, etc.
   - Accepts both `string` and `[]string` inputs
   - Clean error messages and streamlined logic

6. **Dependencies**: 
   - **CLI/TUI**: Cobra, Charm libraries (Huh, Lipgloss)
   - **GUI**: Fyne (macOS/Linux), PowerShell WPF/Windows Forms (Windows)
   - **Build constraints**: Proper platform separation using Go build tags