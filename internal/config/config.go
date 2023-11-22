package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

// Conf struct holds all the repositories to configure
type Conf struct {
	Repos []map[string]struct {
		Name        string `yaml:"name"`
		CurrentFlag string `yaml:"current_flag"`
	} `yaml:"repos"`
}

// ReadConfigFile reads the repositories to scrape from
func ReadConfigFile(yamlData []byte) (Conf, error) {
	var configData Conf

	err := yaml.Unmarshal([]byte(yamlData), &configData)

	if err != nil {
		return Conf{}, err
	}
	log.Println("INFO: Successfully read the cofig.yaml file")

	return configData, nil
}
