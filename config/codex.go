package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// CodexEnv defines environment variables for the Codex wrapper.
// Note: Codex reads OTEL config from ~/.codex/config.toml, not env vars.
// We only set Codex-specific env vars here; OTEL settings are in config.toml.
func CodexEnv(userID, endpoint string) map[string]string {
	return map[string]string{
		"CODEX_ENABLE_TELEMETRY":          "1",
		"CODEX_OTEL_LOG_USER_PROMPT":      "false",
		"CODEX_OTEL_EXPORT_USAGE_METRICS": "true",
		"CODEX_OTEL_INCLUDE_TOKEN_COUNTS": "true",
	}
}

// CodexConfigPath returns the path to Codex config.toml
func CodexConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".codex", "config.toml")
}

// EnableCodexTelemetry writes [otel] section to Codex config.toml
func EnableCodexTelemetry(userID, endpoint string, demoMode bool, logger func(string)) error {
	if logger != nil {
		logger("Enabling Codex telemetry")
	}

	if demoMode {
		return nil
	}

	configPath := CodexConfigPath()

	// Read existing config
	var config map[string]interface{}
	data, err := os.ReadFile(configPath)
	if err == nil && len(data) > 0 {
		if err := toml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("parsing existing Codex config: %w", err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading Codex config: %w", err)
	} else {
		config = make(map[string]interface{})
	}

	// Codex needs both logs and traces endpoints
	logsEndpoint := endpoint + "/v1/logs"
	tracesEndpoint := endpoint + "/v1/traces"

	// Create [otel] section with both log and trace exporters
	otel := map[string]interface{}{
		"environment":     "prod",
		"log_user_prompt": false,
		"exporter": map[string]interface{}{
			"otlp-http": map[string]interface{}{
				"endpoint": logsEndpoint,
				"protocol": "binary",
				"headers": map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", userID),
				},
			},
		},
		"trace_exporter": map[string]interface{}{
			"otlp-http": map[string]interface{}{
				"endpoint": tracesEndpoint,
				"protocol": "binary",
				"headers": map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", userID),
				},
			},
		},
	}

	config["otel"] = otel

	// Write updated config
	output, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("encoding Codex config: %w", err)
	}

	if err := os.WriteFile(configPath, output, 0644); err != nil {
		return fmt.Errorf("writing Codex config: %w", err)
	}

	return nil
}
