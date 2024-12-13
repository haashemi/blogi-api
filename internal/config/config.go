package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseConfig
	API APIConfig
}

type BaseConfig struct {
	DBConn string `env:"BLOGI_DB_CONN"`
}

type APIConfig struct {
	APIAddr    string `env:"BLOGI_API_ADDR" envDefault:":8080"`
	HMAC       string `env:"BLOGI_API_HMAC"`
	IsDevBuild bool   `env:"BLOGI_IS_DEV_BUILD" envDefault:"false"`
}

func Load() (config Config, err error) {
	if err = godotenv.Load(".env"); err != nil {
		return
	}

	if err = env.Parse(&config.BaseConfig); err != nil {
		return
	} else if err = env.Parse(&config.API); err != nil {
		return
	}
	return
}
