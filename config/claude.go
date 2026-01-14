package config

// ClaudeEnv defines the OTEL environment for Claude Code.
func ClaudeEnv(userID, endpoint string) map[string]string {
	env := baseOTELEnv(userID, endpoint, "claude-code")
	env["CLAUDE_CODE_ENABLE_TELEMETRY"] = "1"
	return env
}
