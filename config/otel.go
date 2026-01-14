package config

import "fmt"

// baseOTELEnv generates a baseline OTEL environment variable map for the given service.
func baseOTELEnv(userID, endpoint, serviceName string) map[string]string {
	return map[string]string{
		"OTEL_SERVICE_NAME":           serviceName,
		"OTEL_RESOURCE_ATTRIBUTES":    fmt.Sprintf("user.private_uuid=%s", userID),
		"OTEL_EXPORTER_OTLP_PROTOCOL": "http/json",
		"OTEL_EXPORTER_OTLP_ENDPOINT": endpoint,
		"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": endpoint + "/v1/traces",
		"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": endpoint + "/v1/metrics",
		"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": endpoint + "/v1/logs",
		"OTEL_EXPORTER_OTLP_HEADERS":  fmt.Sprintf("Authorization=Bearer %s", userID),
		"OTEL_TRACES_EXPORTER":        "otlp",
		"OTEL_METRICS_EXPORTER":       "otlp",
		"OTEL_LOGS_EXPORTER":          "otlp",
		"OTEL_TRACES_SAMPLER":         "parentbased_traceidratio",
		"OTEL_TRACES_SAMPLER_ARG":     "1.0",
	}
}

// OTELEnv is kept for backward compatibility and returns the Claude Code defaults.
func OTELEnv(userID, endpoint string) map[string]string {
	return baseOTELEnv(userID, endpoint, "claude-code")
}
