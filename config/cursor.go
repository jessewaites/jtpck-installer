package config

// CursorEnv defines the OTEL environment for Cursor.
func CursorEnv(userID, endpoint string) map[string]string {
	return baseOTELEnv(userID, endpoint, "cursor")
}
