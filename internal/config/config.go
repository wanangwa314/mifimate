package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Device DeviceConfig `mapstructure:"device"`
	App    AppConfig    `mapstructure:"app"`
}

// DeviceConfig holds device-specific configuration
type DeviceConfig struct {
	DefaultIP         string `mapstructure:"default_ip"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	ConnectionTimeout int    `mapstructure:"connection_timeout"`
	PollInterval      int    `mapstructure:"poll_interval"` // seconds
	AutoReconnect     bool   `mapstructure:"auto_reconnect"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Theme             string `mapstructure:"theme"`    // light, dark, system
	Language          string `mapstructure:"language"` // en, es, fr, etc.
	AutoStart         bool   `mapstructure:"auto_start"`
	MinimizeToTray    bool   `mapstructure:"minimize_to_tray"`
	ShowNotifications bool   `mapstructure:"show_notifications"`
	LogLevel          string `mapstructure:"log_level"` // debug, info, warn, error
}

func DefaultConfig() *Config {
	return &Config{
		Device: DeviceConfig{
			DefaultIP:         "192.168.1.1",
			Username:          "admin",
			Password:          "admin",
			ConnectionTimeout: 10,
			PollInterval:      3,
			AutoReconnect:     true,
		},
		App: AppConfig{
			Theme:             "system",
			Language:          "en",
			AutoStart:         false,
			MinimizeToTray:    true,
			ShowNotifications: true,
			LogLevel:          "info",
		},
	}
}

func Load() (*Config, error) {
	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set defaults
	cfg := DefaultConfig()
	viper.SetDefault("device", cfg.Device)
	viper.SetDefault("app", cfg.App)

	// Try to read existing config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create with defaults
			configPath := filepath.Join(configDir, "config.yaml")
			if err := viper.WriteConfigAs(configPath); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Set values in viper
	viper.Set("device", c.Device)
	viper.Set("app", c.App)

	// Write to file
	configPath := filepath.Join(configDir, "config.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Use platform-specific config directory
	configDir := filepath.Join(homeDir, ".config", "mifi-manager")

	return configDir, nil
}
