package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// DbConfig represents the configuration structure for the database.
type DbConfig struct {
	URI    string `yaml:"uri"`
	DbName string `yaml:"dbname"`
}

// LoadConfig loads the database configuration from a YAML file using Viper.
func LoadConfig() (DbConfig, error) {
	// Initialize an instance of DbConfig to hold the configuration values
	var config DbConfig

	// Set the configuration file path, name, and type for Viper
	viper.AddConfigPath("./config")
	viper.SetConfigName("database-config")
	viper.SetConfigType("yaml")

	// Enable VIPER to read environment variables for configuration
	viper.AutomaticEnv()

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal the configuration into the DbConfig struct
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// Return the populated DbConfig and any error encountered
	return config, err
}
