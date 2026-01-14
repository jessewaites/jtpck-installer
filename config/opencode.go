package config

// OpenCodeEnv defines the OTEL environment for OpenCode.
func OpenCodeEnv(userID, endpoint string) map[string]string {
	return baseOTELEnv(userID, endpoint, "opencode")
}
