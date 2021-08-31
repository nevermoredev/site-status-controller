package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	User string
	Host string
	Port string
	Password string
}

var configPath ="./build/rmqConfig.yaml"


func GetConfig() *Config {


	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil
	}

	return config

}




