package sisito

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Database DatabaseConfig
	User     []UserConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int64
	Database string
	Username string
	Password string
}

type UserConfig struct {
	Userid   string
	Password string
}

func LoadConfig(flags *Flags) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(flags.Config, config)

	if err != nil {
		return
	}

	database := config.Database

	if database.Host == "" {
		database.Host = "localhost"
	}

	if database.Port == 0 {
		database.Port = 3306
	}

	if database.Username == "" {
		database.Username = "root"
	}

	return
}
