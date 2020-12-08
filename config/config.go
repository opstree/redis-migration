package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	OldRedis          Redis `yaml:"old_redis"`
	ConcurrentWorkers int   `yaml:"concurrent_workers"`
	NewRedis          Redis `yaml:"new_redis"`
	Databases         []int `yaml:"migration_databases"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

// ParseConfig is to parse YAML configuration file
func ParseConfig(configFile string) Configuration {
	var configContent Configuration
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		logrus.Errorf("Error while reading file %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &configContent)
	if err != nil {
		logrus.Errorf("Error in parsing file to yaml content %v", err)
	}
	return configContent
}
