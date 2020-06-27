package config

import (
	"github.com/caarlos0/env"
)

type config struct {
	Addr            string `env:"SERVICE_HOST" envDefault:"0.0.0.0"`
	Port            int    `env:"SERVICE_PORT" envDefault:"8080"`
	ReadTimeout     int    `env:"TIMEOUT_SERVER_READ" envDefault:"60"`
	WriteTimeout    int    `env:"TIMEOUT_SERVER_WRITE" envDefault:"60"`
}

var Conf *config

func init() {
	var err error
	if Conf, err = ParseConfig(); err != nil {
		panic(err.Error())
	}
}

func ParseConfig() (*config, error) {
	var cfg config
	err := env.Parse(&cfg)
	return &cfg, err
}
