package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HuautlaHost string `envconfig:"HUAUTLA_HOST" default:"localhost" json:"huautla_host,omitempty"`
	HuautlaPort uint   `envconfig:"HUAUTLA_PORT" default:"5432" json:"huautla_port,omitempty"`
	HuautlaUser string `envconfig:"HUAUTLA_USER" default:"postgres" json:"huautla_user,omitempty"`
	HuautlaPass string `envconfig:"HUAUTLA_PASS" default:"root" json:"-"`
	HuautlaSSL  string `envconfig:"HUAUTLA_SSL" default:"disable" json:"huautla_ssl,omitempty"`

	AuthnHost string `envconfig:"AUTHN_HOST" json:"authn_host,omitempty"`
	AuthnPort uint16 `envconfig:"AUTHN_PORT" json:"authn_port,omitempty"`
	AuthnPath string `envconfig:"AUTHN_PATH" required:"true" json:"authn_path,omitempty"`

	HTTPHost string `envconfig:"HTTP_HOST" default:"127.0.0.1" json:"http_host,omitempty"`
	HTTPPort int    `envconfig:"HTTP_PORT" default:"8080" json:"http_port,omitempty"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO" json:"log_level,omitempty"`

	PhotoDir string `envconfig:"PHOTO_ROOT" default:"album/" json:"photo_location,omitempty"` // XXX: is this too much info to make public?
}

func NewConfig() (*Config, error) {
	result := &Config{}
	if err := envconfig.Process("CFFC", result); err != nil {
		return nil, err
	}

	if result.AuthnHost == "" {
		// so much for required; if there's no host to contact, the web shouldn't try
		// to check valid; get this value from /settings
		result.AuthnPath = ""
	}

	// XXX: replacing other os's path-separator with our separator is out
	//  of scope for this endeavor; just make sure it ends with ours
	if result.PhotoDir[len(result.PhotoDir)-1] != os.PathSeparator {
		result.PhotoDir = string(append(([]byte)(result.PhotoDir), os.PathSeparator))
	}

	return result, nil
}
