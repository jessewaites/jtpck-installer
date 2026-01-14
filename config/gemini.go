package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GeminiEnv defines the telemetry environment variables for the Gemini CLI.
func GeminiEnv(userID, endpoint string) map[string]string {
	return map[string]string{
		"GEMINI_TELEMETRY_ENABLED":       "true",
		"GEMINI_TELEMETRY_TARGET":        "local",
		"GEMINI_TELEMETRY_OTLP_ENDPOINT": endpoint,
		"GEMINI_TELEMETRY_OTLP_PROTOCOL": "http",
		"GEMINI_TELEMETRY_USE_COLLECTOR": "true",
		"OTEL_EXPORTER_OTLP_HEADERS":     fmt.Sprintf("Authorization=Bearer %s", userID),
		"OTEL_RESOURCE_ATTRIBUTES":       fmt.Sprintf("user.private_uuid=%s", userID),
	}
}

// GeminiSettingsDir returns the path to the .gemini directory.
func GeminiSettingsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gemini")
}

// GeminiSettingsPath returns the path to the Gemini CLI settings file.
func GeminiSettingsPath() string {
	return filepath.Join(GeminiSettingsDir(), "settings.json")
}

// EnableGeminiTelemetry writes telemetry settings for the Gemini CLI.
func EnableGeminiTelemetry(userID, endpoint string, demoMode bool, logger func(string)) error {
	if logger != nil {
		logger("Enabling Gemini CLI telemetry")
	}

	if demoMode {
		return nil
	}

	if err := os.MkdirAll(GeminiSettingsDir(), 0755); err != nil {
		return fmt.Errorf("creating Gemini settings directory: %w", err)
	}

	settings := map[string]interface{}{}

	data, err := os.ReadFile(GeminiSettingsPath())
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("parsing existing Gemini settings: %w", err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading Gemini settings: %w", err)
	}

	var telemetry map[string]interface{}
	if existing, ok := settings["telemetry"].(map[string]interface{}); ok {
		telemetry = existing
	} else {
		telemetry = map[string]interface{}{}
	}

	telemetry["enabled"] = true
	telemetry["target"] = "local"
	telemetry["otlpEndpoint"] = endpoint
	telemetry["otlpProtocol"] = "http"
	telemetry["useCollector"] = true

	settings["telemetry"] = telemetry

	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding Gemini settings: %w", err)
	}

	if err := os.WriteFile(GeminiSettingsPath(), output, 0644); err != nil {
		return fmt.Errorf("writing Gemini settings: %w", err)
	}

	return nil
}
