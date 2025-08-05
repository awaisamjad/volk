package config

import (
	"fmt"
	// "log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const configFileName = "volk_config.toml"

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port        int `toml:"port"`
	ReadTimeout int `toml:"read_timeout"` // seconds
}

// FileServerConfig holds file serving configuration
type FileServerConfig struct {
	DocumentRoot string `toml:"document_root"`
	DefaultFile  string `toml:"default_file"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Format     string `toml:"format"`      // plain, verbose
	FilePath   string `toml:"file_path"`   // Path to log file, empty for stdout
	AccessLogs bool   `toml:"access_logs"` // Enable HTTP access logging
}

// Config is the root configuration structure
type Config struct {
	Server     ServerConfig     `toml:"server"`
	FileServer FileServerConfig `toml:"file_server"`
	Logging    LogConfig        `toml:"logging"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:        6543,
			ReadTimeout: 30,
		},
		FileServer: FileServerConfig{
			DocumentRoot: ".",
			DefaultFile:  "index.html",
		},
		Logging: LogConfig{
			Format:     "plain",
			FilePath:   "",
			AccessLogs: true,
		},
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`
[server]
port = %d
read_timeout = %d

[file_server]
document_root = "%s"
default_file = "%s"

[logging]
format = "%s"
file_path = "%s"
access_logs = %t`,
		c.Server.Port, c.Server.ReadTimeout,
		c.FileServer.DocumentRoot, c.FileServer.DefaultFile,
		c.Logging.Format, c.Logging.FilePath, c.Logging.AccessLogs)
}

// LoadConfig loads configuration from a TOML file.
// It returns the configuration and an error, if any.
func LoadConfig() (Config, error) {
	config := DefaultConfig()

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	_, err := toml.DecodeFile(configFileName, &config)
	if err != nil {
		return config, fmt.Errorf("error decoding config file: %w", err)
	}

	if !filepath.IsAbs(config.FileServer.DocumentRoot) {
		absPath, err := filepath.Abs(config.FileServer.DocumentRoot)
		if err == nil {
			config.FileServer.DocumentRoot = absPath
		} else {
			return config, fmt.Errorf(`could not determine absolute path for document_root, using default config.
			error: %w`, err)
		}
	}

	return config, nil
}

// SaveConfig saves the configuration to a TOML file
func SaveConfig(config Config) error {
	f, err := os.Create(configFileName)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}

	return nil
}

// CreateDefaultConfigFile creates a default configuration file if it doesn't exist
func CreateDefaultConfigFile() error {
	if _, err := os.Stat(configFileName); err == nil {
		return nil
	}

	config := DefaultConfig()
	return SaveConfig(config)
}
