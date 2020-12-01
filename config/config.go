package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	OldRedis          OldRedis `yaml:"old_redis"`
	ConcurrentWorkers int      `yaml:"concurrent_workers"`
	NewRedis          NewRedis `yaml:"new_redis"`
}

type NewRedis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type OldRedis struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Databases []int  `yaml:"databases"`
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
