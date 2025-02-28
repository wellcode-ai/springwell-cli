package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the SpringWell CLI configuration
type Config struct {
	Project struct {
		Package           string `mapstructure:"package"`
		DefaultsDirectory string `mapstructure:"defaultsDirectory"`
	} `mapstructure:"project"`

	Code struct {
		Style struct {
			Indentation int `mapstructure:"indentation"`
			LineWidth   int `mapstructure:"lineWidth"`
		} `mapstructure:"style"`
		Lombok            bool `mapstructure:"lombok"`
		StandardizeFields bool `mapstructure:"standardizeFields"`
	} `mapstructure:"code"`

	Templates struct {
		Directory string `mapstructure:"directory"`
	} `mapstructure:"templates"`

	AWS struct {
		Region          string   `mapstructure:"region"`
		DefaultServices []string `mapstructure:"defaultServices"`
	} `mapstructure:"aws"`

	Plugins []string `mapstructure:"plugins"`
}

// LoadConfig loads the configuration from the .springwell.yml file
func LoadConfig(projectDir string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(".springwell")
	v.SetConfigType("yml")
	v.AddConfigPath(projectDir)

	// Set default values
	v.SetDefault("project.defaultsDirectory", ".springwell/templates")
	v.SetDefault("code.style.indentation", 4)
	v.SetDefault("code.style.lineWidth", 120)
	v.SetDefault("code.lombok", true)
	v.SetDefault("code.standardizeFields", true)
	v.SetDefault("templates.directory", ".springwell/templates")
	v.SetDefault("aws.region", "us-east-1")
	v.SetDefault("aws.defaultServices", []string{"s3", "secretsManager"})

	// Try to read the config file
	if err := v.ReadInConfig(); err != nil {
		// If config file doesn't exist, that's fine. We use defaults.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	config := &Config{}
	config.Project.Package = "com.example.service"
	config.Project.DefaultsDirectory = ".springwell/templates"

	config.Code.Style.Indentation = 4
	config.Code.Style.LineWidth = 120
	config.Code.Lombok = true
	config.Code.StandardizeFields = true

	config.Templates.Directory = ".springwell/templates"

	config.AWS.Region = "us-east-1"
	config.AWS.DefaultServices = []string{"s3", "secretsManager"}

	return config
}

// SaveConfig saves the configuration to the .springwell.yml file
func SaveConfig(config *Config, projectDir string) error {
	v := viper.New()
	v.SetConfigFile(filepath.Join(projectDir, ".springwell.yml"))

	// Convert config struct to map
	v.Set("project", config.Project)
	v.Set("code", config.Code)
	v.Set("templates", config.Templates)
	v.Set("aws", config.AWS)
	v.Set("plugins", config.Plugins)

	// Create the directory if it doesn't exist
	configDir := filepath.Join(projectDir, ".springwell")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return v.WriteConfig()
}
