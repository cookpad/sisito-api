package sisito

import (
	"github.com/BurntSushi/toml"
)

const (
	ConfigFile = "config.tml"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int64
	Database string
	User     string
	Password string
}

func LoadConfig() (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(ConfigFile, config)
	return
}
