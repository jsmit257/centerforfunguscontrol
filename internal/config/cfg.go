package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HuautlaHost string `envconfig:"HUAUTLA_HOST" default:"localhost"`
	HuautlaPort uint   `envconfig:"HUAUTLA_PORT" default:"5432"`
	HuautlaUser string `envconfig:"HUAUTLA_USER" default:"postgres"`
	HuautlaPass string `envconfig:"HUAUTLA_PASS" default:"root"`
	HuautlaSSL  string `envconfig:"HUAUTLA_SSL" default:"disable"`

	AuthnHost string `envconfig:"AUTHN_HOST"`
	AuthnPort uint16 `envconfig:"AUTHN_PORT"`

	HTTPHost string `envconfig:"HTTP_HOST" default:"127.0.0.1"`
	HTTPPort int    `envconfig:"HTTP_PORT" default:"8080"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`

	PhotoDir string `envconfig:"PHOTO_ROOT" default:"album/"`
}

func NewConfig() *Config {
	result := &Config{}
	if err := envconfig.Process("CFFC", result); err != nil {
		panic(err)
	}

	if result.PhotoDir[len(result.PhotoDir)-1] != os.PathSeparator {
		result.PhotoDir = string(append(([]byte)(result.PhotoDir), os.PathSeparator))
	}

	return result
}
