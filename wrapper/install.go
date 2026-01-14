package wrapper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// WrapperDir returns the path to .jtpck directory
func WrapperDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".jtpck")
}

// WrapperPath returns the path for a specific tool wrapper
func WrapperPath(toolName string) string {
	return filepath.Join(WrapperDir(), fmt.Sprintf("%s-wrapper", toolName))
}

// CreateWrappers creates wrapper scripts for all tools with per-tool env vars.
func CreateWrappers(envs map[string]map[string]string, tools []string) error {
	// Ensure directory exists
	if err := os.MkdirAll(WrapperDir(), 0755); err != nil {
		return fmt.Errorf("failed to create wrapper directory: %w", err)
	}

	for _, tool := range tools {
		// Find tool path
		toolPath, err := exec.LookPath(tool)
		if err != nil {
			// Tool not found, skip but continue
			continue
		}

		env := envs[tool]
		// Generate script
		script := GenerateScript(tool, toolPath, env)

		// Write to file
		wrapperPath := WrapperPath(tool)
		if err := os.WriteFile(wrapperPath, []byte(script), 0755); err != nil {
			return fmt.Errorf("failed to write %s wrapper: %w", tool, err)
		}
	}

	return nil
}

// GetInstalledTools returns which tools have wrappers created
func GetInstalledTools(tools []string) []string {
	var installed []string
	for _, tool := range tools {
		if _, err := os.Stat(WrapperPath(tool)); err == nil {
			installed = append(installed, tool)
		}
	}
	return installed
}
