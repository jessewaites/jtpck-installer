package wrapper

import (
	"fmt"
	"strings"
	"time"
)

// shellEscape escapes a string for safe use in shell double quotes
func shellEscape(s string) string {
	// Escape backslashes, double quotes, dollar signs, and backticks
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, `$`, `\$`)
	s = strings.ReplaceAll(s, "`", "\\`")
	return s
}

// GenerateScript creates a wrapper script for the given tool
func GenerateScript(toolName, toolPath string, env map[string]string) string {
	var sb strings.Builder

	sb.WriteString("#!/bin/bash\n")
	sb.WriteString(fmt.Sprintf("# JTPCK Telemetry Wrapper for %s\n", toolName))
	sb.WriteString(fmt.Sprintf("# Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	// Export environment variables with proper escaping
	for key, value := range env {
		escapedValue := shellEscape(value)
		sb.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, escapedValue))
	}

	sb.WriteString("\n")
	// toolPath comes from exec.LookPath() so it's trusted, but escape anyway
	sb.WriteString(fmt.Sprintf("exec \"%s\" \"$@\"\n", shellEscape(toolPath)))

	return sb.String()
}
