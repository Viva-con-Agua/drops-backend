package utils

import (
	"github.com/jinzhu/configor"
)

var (
	Config = struct {
		DB struct {
			Host     string `required:"true"`
			User     string `default:"admin"`
			Password string `default:"admin"`
			Port     int    `default:"3306"`
			Name     string `required:"true"`
		}
		REDIS struct {
			Host string `required:"true"`
			Port int    `default:"6379"`
		}
	}{}
)

func LoadConfig() {
	configor.Load(&Config, "config.yml")
}
