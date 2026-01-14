package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config represents the JTPCK configuration
type Config struct {
	UserID   string    `json:"user_id"`
	Endpoint string    `json:"endpoint"`
	Created  time.Time `json:"created_at"`
	Updated  time.Time `json:"updated_at"`
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".jtpck", "config.json")
}

// ConfigDir returns the path to the .jtpck directory
func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".jtpck")
}

// Exists checks if config file exists
func Exists() bool {
	_, err := os.Stat(ConfigPath())
	return err == nil
}

// Load reads config from disk
func Load() (*Config, error) {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes config to disk
func (c *Config) Save() error {
	// Ensure directory exists
	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return err
	}

	c.Updated = time.Now()
	if c.Created.IsZero() {
		c.Created = time.Now()
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), data, 0644)
}

// New creates a new config instance
func New(userID, endpoint string) *Config {
	return &Config{
		UserID:   userID,
		Endpoint: endpoint,
		Created:  time.Now(),
		Updated:  time.Now(),
	}
}
