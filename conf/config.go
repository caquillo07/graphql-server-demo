package conf

import (
	"log"
	"strings"

	"github.com/spf13/viper"

	"github.com/caquillo07/graphql-server-demo/database"
)

// Config contains all of the configurable pieces of the application.
type Config struct {
	Server struct {
		// Port is the which the server will listen to
		Port string
	}

	CORS struct {
		Enabled        bool
		Debug          bool
		AllowedOrigins []string
	}

	GraphQL struct {
		Playground bool
		LogQueries bool
	}

	Auth struct {
		// Enabled indicates whether or not the API will require a token
		// for each request
		Enabled bool
	}

	// Config provides database configuration
	Database database.Config
}

// LoadConfig loads configuration from the viper instance.
func LoadConfig(v *viper.Viper) (Config, error) {
	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

// InitViper initializes the global viper instance with the configuration file
func InitViper(configFile string) {
	if configFile == "" {

		// default to one in present directory named 'config'
		configFile = "config.yaml"
	}
	viper.SetConfigFile(configFile)

	// Default settings
	viper.SetDefault("server.allowCORS", true)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	log.Printf("using config file: %s", viper.ConfigFileUsed())
}
