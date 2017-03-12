package sisito

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Database DatabaseConfig
	User     []User
}

type DatabaseConfig struct {
	Host     string
	Port     int64
	Database string
	User     string
	Password string
}

type User struct {
	Userid   string
	Password string
}

func LoadConfig(path string) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(path, config)

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

	if database.User == "" {
		database.User = "root"
	}

	return
}
