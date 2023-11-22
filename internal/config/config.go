package config

import (
	"log"
	"os"

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
func ReadConfigFile() (Conf, error) {
	var configData Conf
	conf, err := os.ReadFile("config.yaml")

	if err != nil {
		return Conf{}, err
	}
	err = yaml.Unmarshal([]byte(conf), &configData)

	if err != nil {
		return Conf{}, err
	}
	log.Println("INFO: Successfully read the cofig.yaml file")

	return configData, nil
}
