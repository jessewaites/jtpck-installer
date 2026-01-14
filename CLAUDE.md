# JTPCK Installer

CLI tool for configuring OpenTelemetry in Claude Code, Codex, and Gemini CLI.

## Architecture

- **Go + Cobra CLI** - Command framework
- **Bubble Tea TUI** - Interactive UI (animation, input, success screens)
- **Telemetry setup** - Creates wrapper scripts, modifies config files

## Core Flow

1. Validates UUID format (if provided as arg)
2. Checks for existing config
3. Validates tool installations (claude/codex/gemini)
4. Shows animation
5. Collects user UUID (if not provided)
6. Configures telemetry:
   - Creates wrapper scripts in `~/.jtpck/wrappers/`
   - Updates `~/.codex/config.toml` [otel] section
   - Updates `~/.gemini/settings.json` telemetry section
   - Installs shell aliases to `~/.zshrc`
7. Saves config to `~/.jtpck/config.json`

## Commands

- `jtpck [UUID]` - Setup (interactive if no UUID)
- `jtpck configure` - Reconfigure with new UUID
- `jtpck uninstall` - Remove all configs/wrappers/aliases
- `jtpck --demo` - Preview UI without writing files
- `jtpck --version` - Show version

## Key Files

- `JTPCK.go` - Main entry
- `cmd/root.go` - Setup command (line 54)
- `cmd/configure.go` - Reconfigure command
- `cmd/uninstall.go` - Cleanup command
- `config/*.go` - Config management per tool
- `wrapper/*.go` - Wrapper script creation
- `shell/*.go` - Shell detection and alias installation
- `ui/*.go` - Bubble Tea TUI components
- `validator/*.go` - Tool installation validation

## Telemetry Endpoint

`https://JTPCK.com/api/v1/telemetry`

## UUID Format

Lowercase, hyphenated UUID: `^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`

## Homebrew Tap

Preparing for Homebrew release. Binary name: `jtpck`

## Dev Notes

- Version in `cmd/root.go:20` (currently 0.1.0)
- Wrappers inject env vars per tool (see `config/apps.go`)
- Shell auto-detection supports zsh/bash
- Demo mode skips all file writes for testing UI
