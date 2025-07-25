package config

import (
	"fmt"
	// "log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port           int    `toml:"port"`
	Host           string `toml:"host"`
	ReadTimeout    int    `toml:"read_timeout"`  // seconds
	WriteTimeout   int    `toml:"write_timeout"` // seconds
	MaxConnections int    `toml:"max_connections"`
}

// FileServerConfig holds file serving configuration
type FileServerConfig struct {
	DocumentRoot      string            `toml:"document_root"`
	DefaultFile       string            `toml:"default_file"`
	AllowListing      bool              `toml:"allow_directory_listing"`
	MimeTypeOverrides map[string]string `toml:"mime_type_overrides"`
}

// SecurityConfig holds security-related settings
type SecurityConfig struct {
	AllowDirectoryTraversal bool     `toml:"allow_directory_traversal"` // Should be false in production
	MaxRequestSize          int      `toml:"max_request_size"`          // in bytes
	RateLimit               int      `toml:"rate_limit"`                // requests per minute
	AllowedOrigins          []string `toml:"allowed_origins"`           // CORS origins
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
	Security   SecurityConfig   `toml:"security"`
	Logging    LogConfig        `toml:"logging"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:           6543,
			Host:           "0.0.0.0",
			ReadTimeout:    30,
			WriteTimeout:   30,
			MaxConnections: 100,
		},
		FileServer: FileServerConfig{
			DocumentRoot:      ".",
			DefaultFile:       "index.html",
			AllowListing:      false,
			MimeTypeOverrides: map[string]string{},
		},
		Security: SecurityConfig{
			AllowDirectoryTraversal: false,
			MaxRequestSize:          1048576, // 1MB
			RateLimit:               60,      // 60 requests per minute
			AllowedOrigins:          []string{"*"},
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
host = %s
read_timeout = %d
write_timeout = %d
max_connections = %d

[file_server]
document_root = %s
default_file = %s
allow_directory_listing = %t
mime_type_overrides = %v

[security]
allow_directory_traversal = %t
max_request_size = %d
rate_limit = %d
allowed_origins = %v

[logging]
format = %s
file_path = %s
access_logs = %t`,
		c.Server.Port, c.Server.Host, c.Server.ReadTimeout, c.Server.WriteTimeout, c.Server.MaxConnections,
		c.FileServer.DocumentRoot, c.FileServer.DefaultFile, c.FileServer.AllowListing, c.FileServer.MimeTypeOverrides,
		c.Security.AllowDirectoryTraversal, c.Security.MaxRequestSize, c.Security.RateLimit, c.Security.AllowedOrigins,
		c.Logging.Format, c.Logging.FilePath, c.Logging.AccessLogs)
}

// LoadConfig loads configuration from a TOML file
func LoadConfig(filePath string) (Config, error) {
	config := DefaultConfig()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	_, err := toml.DecodeFile(filePath, &config)
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
func SaveConfig(config Config, filePath string) error {
	f, err := os.Create(filePath)
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
func CreateDefaultConfigFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		// File exists, don't overwrite
		return nil
	}

	config := DefaultConfig()
	return SaveConfig(config, filePath)
}
