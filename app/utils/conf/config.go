package conf

import (
	"io/ioutil"

	"github.com/growerlab/backend/app/common/errors"
	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigPath = "conf/config.yaml"
)

var (
	config *Config
)

type DB struct {
	URL string `yaml:"url"`
}

type Config struct {
	Database *DB `yaml:"db"`
}

func GetConf() *Config {
	return config
}

func LoadConfig() error {
	confBody, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		return errors.Trace(err)
	}
	err = yaml.Unmarshal(confBody, &config)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
