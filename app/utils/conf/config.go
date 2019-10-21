package conf

import (
	"io/ioutil"
	"path/filepath"

	"github.com/growerlab/backend/app/common/errors"
	"github.com/growerlab/backend/app/utils"
	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigPath = "conf/config.yaml"
)

var (
	config *Config
)

type DB struct {
	URL string `yaml:"url,omitempty"`
}

type Config struct {
	Database *DB `yaml:"db,omitempty"`
}

func GetConf() *Config {
	return config
}

func LoadConfig() error {
	path := filepath.Join(utils.BasePath(), DefaultConfigPath)
	confBody, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Trace(err)
	}
	err = yaml.Unmarshal(confBody, config)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
