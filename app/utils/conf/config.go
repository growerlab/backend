package conf

import (
	"io/ioutil"
	"net/url"

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
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Namespace   string `yaml:"namespace"`
	MaxIdle     int    `yaml:"max_idle"`
	MaxActive   int    `yaml:"max_active"`
	IdleTimeout int    `yaml:"idle_timeout"`
}

type Mensa struct {
	SshUser  string `yaml:"ssh_user"`
	SshHost  string `yaml:"ssh_host"`
	SshPort  int    `yaml:"ssh_port"`
	HttpHost string `yaml:"http_host"`
	HttpPort int    `yaml:"http_port"`
}

type Config struct {
	Debug      bool   `yaml:"debug"`
	WebsiteURL string `yaml:"website_url"`
	websiteURL *url.URL

	Port     int    `yaml:"port"`
	Database *DB    `yaml:"db"`
	Redis    *Redis `yaml:"redis"`
	Mensa    *Mensa `yaml:"mensa"`
}

func (c *Config) EnableHTTPS() bool {
	if c.websiteURL == nil {
		var err error
		c.websiteURL, err = url.Parse(c.WebsiteURL)
		if err != nil {
			panic(err)
		}
	}
	return c.websiteURL.Scheme == "https"
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
