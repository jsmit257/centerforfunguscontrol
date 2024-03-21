package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HuautlaHost string `envconfig:"HUAUTLA_HOST" default:"localhost"`
	HuautlaPort int    `envconfig:"HUAUTLA_PORT" default:"5432"`
	HuautlaUser string `envconfig:"HUAUTLA_USER" default:"postgres"`
	HuautlaPass string `envconfig:"HUAUTLA_PASS" default:"root"`
	HuautlaSSL  string `envconfig:"HUAUTLA_SSL" default:"disable"`

	HTTPHost string `envconfig:"HTTP_HOST" default:"127.0.0.1"`
	HTTPPort int    `envconfig:"HTTP_PORT" default:"8080"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`
}

func NewConfig() *Config {
	result := &Config{}
	if err := envconfig.Process("CFFC", result); err != nil {
		panic(err)
	}
	return result
}
