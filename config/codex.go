package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// CodexEnv defines the OTEL environment for Codex.
func CodexEnv(userID, endpoint string) map[string]string {
	// Codex uses /v1/traces endpoint
	codexEndpoint := endpoint
	if !strings.HasSuffix(endpoint, "/v1/traces") {
		codexEndpoint = endpoint + "/v1/traces"
	}

	env := baseOTELEnv(userID, codexEndpoint, "codex")
	env["CODEX_ENABLE_TELEMETRY"] = "1"
	env["CODEX_OTEL_LOG_USER_PROMPT"] = "false"
	env["CODEX_OTEL_EXPORT_USAGE_METRICS"] = "true"
	env["CODEX_OTEL_INCLUDE_TOKEN_COUNTS"] = "true"
	return env
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

	// Codex uses /v1/traces endpoint
	codexEndpoint := endpoint
	if !strings.HasSuffix(endpoint, "/v1/traces") {
		codexEndpoint = endpoint + "/v1/traces"
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

	// Create [otel] section with otlp-http exporter
	otel := map[string]interface{}{
		"environment":     "prod",
		"log_user_prompt": false,
		"exporter": map[string]interface{}{
			"otlp-http": map[string]interface{}{
				"endpoint": codexEndpoint,
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
