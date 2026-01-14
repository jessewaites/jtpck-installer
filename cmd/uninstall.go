package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove JTPCK telemetry configuration",
	Long:  `Removes all JTPCK configuration files, wrappers, and shell aliases.`,
	Run:   runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	if demoMode {
		fmt.Println("🧹 Cleaning up JTPCK installer artifacts... (DEMO MODE)")
	} else {
		fmt.Println("🧹 Cleaning up JTPCK installer artifacts...")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// 1. Remove .jtpck directory
	jtpckDir := filepath.Join(home, ".jtpck")
	if _, err := os.Stat(jtpckDir); err == nil {
		fmt.Println("  Removing ~/.jtpck/")
		if !demoMode {
			if err := os.RemoveAll(jtpckDir); err != nil {
				fmt.Printf("  ⚠️  Failed to remove ~/.jtpck/: %v\n", err)
			}
		}
	} else {
		fmt.Println("  ~/.jtpck/ not found (skip)")
	}

	// 2. Restore .zshrc from backup or remove JTPCK section
	zshrcPath := filepath.Join(home, ".zshrc")
	backupPath := filepath.Join(home, ".zshrc.jtpck-backup")

	if _, err := os.Stat(backupPath); err == nil {
		fmt.Println("  Restoring ~/.zshrc from backup")
		if !demoMode {
			input, err := os.ReadFile(backupPath)
			if err != nil {
				fmt.Printf("  ⚠️  Failed to read backup: %v\n", err)
			} else {
				if err := os.WriteFile(zshrcPath, input, 0644); err != nil {
					fmt.Printf("  ⚠️  Failed to restore .zshrc: %v\n", err)
				} else {
					os.Remove(backupPath)
				}
			}
		}
	} else if _, err := os.Stat(zshrcPath); err == nil {
		fmt.Println("  Removing JTPCK aliases from ~/.zshrc")
		if !demoMode {
			if err := removeJTPCKSection(zshrcPath); err != nil {
				fmt.Printf("  ⚠️  Failed to update .zshrc: %v\n", err)
			}
		}
	} else {
		fmt.Println("  ~/.zshrc not found (skip)")
	}

	// 3. Remove [otel] section from ~/.codex/config.toml
	codexConfigPath := filepath.Join(home, ".codex", "config.toml")
	if _, err := os.Stat(codexConfigPath); err == nil {
		fmt.Println("  Removing [otel] from ~/.codex/config.toml")
		if !demoMode {
			if err := removeCodexOtelSection(codexConfigPath); err != nil {
				fmt.Printf("  ⚠️  Failed to update Codex config: %v\n", err)
			}
		}
	} else {
		fmt.Println("  ~/.codex/config.toml not found (skip)")
	}

	// 4. Remove telemetry section from ~/.gemini/settings.json
	geminiSettingsPath := filepath.Join(home, ".gemini", "settings.json")
	if _, err := os.Stat(geminiSettingsPath); err == nil {
		fmt.Println("  Removing telemetry from ~/.gemini/settings.json")
		if !demoMode {
			if err := removeGeminiTelemetry(geminiSettingsPath); err != nil {
				fmt.Printf("  ⚠️  Failed to update Gemini settings: %v\n", err)
			}
		}
	} else {
		fmt.Println("  ~/.gemini/settings.json not found (skip)")
	}

	fmt.Println("✓ Cleanup complete!")
	fmt.Println("")
	if demoMode {
		fmt.Println("(DEMO MODE - No files were modified)")
	} else {
		fmt.Println("JTPCK telemetry has been removed from your system.")
	}
}

// removeJTPCKSection removes lines between JTPCK markers from a file
func removeJTPCKSection(filePath string) error {
	input, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string
	skipUntilEnd := false

	for _, line := range lines {
		if strings.Contains(line, "# JTPCK Telemetry Aliases - START") {
			skipUntilEnd = true
			continue
		}
		if strings.Contains(line, "# JTPCK Telemetry Aliases - END") {
			skipUntilEnd = false
			continue
		}
		if !skipUntilEnd {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// removeCodexOtelSection removes [otel] section from Codex config.toml
func removeCodexOtelSection(filePath string) error {
	input, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string
	inOtelSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if entering [otel] section
		if trimmed == "[otel]" {
			inOtelSection = true
			continue
		}

		// Check if entering a different section
		if strings.HasPrefix(trimmed, "[") && trimmed != "[otel]" {
			inOtelSection = false
		}

		// Skip lines in [otel] section
		if inOtelSection {
			continue
		}

		newLines = append(newLines, line)
	}

	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

// removeGeminiTelemetry removes telemetry section from Gemini settings.json
func removeGeminiTelemetry(filePath string) error {
	input, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Simple approach: read line by line and skip telemetry section
	lines := strings.Split(string(input), "\n")
	var newLines []string
	inTelemetry := false
	braceCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if starting telemetry section
		if strings.Contains(trimmed, `"telemetry"`) {
			inTelemetry = true
			braceCount = 0
			continue
		}

		if inTelemetry {
			// Count braces to find end of telemetry object
			braceCount += strings.Count(line, "{")
			braceCount -= strings.Count(line, "}")

			if braceCount <= 0 {
				inTelemetry = false
			}
			continue
		}

		newLines = append(newLines, line)
	}

	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}
