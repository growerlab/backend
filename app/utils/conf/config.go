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

type Redis struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	CacheDB      int    `yaml:"cache_db"`
	QueueDB      int    `yaml:"queue_db"`
	PermissionDB int    `yaml:"permission_db"`
	MaxIdle      int    `yaml:"max_idle"`
	MaxActive    int    `yaml:"max_active"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}

type Config struct {
	Debug      bool   `yaml:"debug"`
	Database   *DB    `yaml:"db"`
	Redis      *Redis `yaml:"redis"`
	WebsiteURL string `yaml:"website_url"`
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
