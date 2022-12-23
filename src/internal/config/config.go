package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Conf struct holds all the repositories to configure
type Conf struct {
	Repos []map[string]string
}

// ReadConfigFile reads the repositories to scrape from the configmap attached to the pod as volume
func ReadConfigFile() (Conf, error) {
	var configData Conf
	conf, err := os.ReadFile("config.yaml")
	if err != nil {
		return Conf{}, err
	}

	err = yaml.Unmarshal(conf, &configData)

	if err != nil {
		return Conf{}, err
	}

	return configData, nil
}
