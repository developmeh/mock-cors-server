package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config holds all configuration for the server
type Config struct {
	Port   int      `mapstructure:"port"`
	Routes []Route  `mapstructure:"routes"`
	CORS   CORSConfig `mapstructure:"cors"`
	Version string  `mapstructure:"version"`
}

// Route represents a single route configuration
type Route struct {
	Path        string     `mapstructure:"path"`
	Type        string     `mapstructure:"type"`        // "static", "json", or "dummy"
	FilePath    string     `mapstructure:"file_path"`   // For static files
	JSONContent string     `mapstructure:"json_content"` // For JSON blob responses
	ContentType string     `mapstructure:"content_type"`
	CORS        *CORSConfig `mapstructure:"cors"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Port: 8081,
		Routes: []Route{
			{
				Path:        "/v1/json/begin",
				Type:        "dummy",
				ContentType: "application/json",
			},
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "site-token", "client-id", "placement-id", "integrator-id", "oauth-type"},
			AllowCredentials: true,
			MaxAge:           86400,
		},
		Version: "1.0.0",
	}
}

// LoadConfig loads the configuration from file and environment variables
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dummy_http_passkeys")
	viper.AddConfigPath("/etc/dummy_http_passkeys")

	// Environment variables
	viper.SetEnvPrefix("MOCK_CORS")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Override with environment variables if they exist
	if viper.IsSet("port") {
		config.Port = viper.GetInt("port")
	}
	if viper.IsSet("version") {
		config.Version = viper.GetString("version")
	}

	// Unmarshal the rest of the config (routes, CORS, etc.)
	var tempConfig Config
	if err := viper.Unmarshal(&tempConfig); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Merge non-zero values from tempConfig
	if len(tempConfig.Routes) > 0 {
		config.Routes = tempConfig.Routes
	}
	if tempConfig.CORS.AllowOrigins != nil {
		config.CORS = tempConfig.CORS
	}

	return config, nil
}
