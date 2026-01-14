package config

// AppEnvs returns per-application environment variables keyed by tool name.
func AppEnvs(userID, endpoint string) map[string]map[string]string {
	return map[string]map[string]string{
		"claude": ClaudeEnv(userID, endpoint),
		"codex":  CodexEnv(userID, endpoint),
		"gemini": GeminiEnv(userID, endpoint),
	}
}
