#!/bin/bash
# Cleanup script to remove JTPCK installer artifacts

set -e

echo "🧹 Cleaning up JTPCK installer artifacts..."

# 1. Remove .jtpck directory
if [ -d "$HOME/.jtpck" ]; then
    echo "  Removing ~/.jtpck/"
    rm -rf "$HOME/.jtpck"
else
    echo "  ~/.jtpck/ not found (skip)"
fi

# 2. Restore .zshrc from backup or remove JTPCK section
if [ -f "$HOME/.zshrc.jtpck-backup" ]; then
    echo "  Restoring ~/.zshrc from backup"
    cp "$HOME/.zshrc.jtpck-backup" "$HOME/.zshrc"
    rm "$HOME/.zshrc.jtpck-backup"
elif [ -f "$HOME/.zshrc" ]; then
    echo "  Removing JTPCK aliases from ~/.zshrc"
    # Remove JTPCK section between START and END markers
    sed -i.bak '/# JTPCK Telemetry Aliases - START/,/# JTPCK Telemetry Aliases - END/d' "$HOME/.zshrc"
    rm -f "$HOME/.zshrc.bak"
else
    echo "  ~/.zshrc not found (skip)"
fi

# 3. Remove [otel] section from ~/.codex/config.toml
if [ -f "$HOME/.codex/config.toml" ]; then
    echo "  Removing [otel] from ~/.codex/config.toml"
    # Remove [otel] section and everything until next section or EOF
    sed -i.bak '/^\[otel\]/,/^\[/{ /^\[otel\]/d; /^\[/!d; }' "$HOME/.codex/config.toml"
    rm -f "$HOME/.codex/config.toml.bak"
else
    echo "  ~/.codex/config.toml not found (skip)"
fi

# 4. Remove telemetry section from ~/.gemini/settings.json
if [ -f "$HOME/.gemini/settings.json" ]; then
    echo "  Removing telemetry from ~/.gemini/settings.json"
    # Check if jq is available
    if command -v jq &> /dev/null; then
        # Use jq to remove telemetry section
        jq 'del(.telemetry)' "$HOME/.gemini/settings.json" > "$HOME/.gemini/settings.json.tmp"
        mv "$HOME/.gemini/settings.json.tmp" "$HOME/.gemini/settings.json"
    else
        echo "  ⚠️  jq not found - manually remove telemetry section from ~/.gemini/settings.json"
    fi
else
    echo "  ~/.gemini/settings.json not found (skip)"
fi

echo "✓ Cleanup complete!"
echo ""
echo "You can now run the installer again for testing."
