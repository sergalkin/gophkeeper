package server

import (
	"github.com/caarlos0/env/v6"

	"github.com/sergalkin/gophkeeper/pkg"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS" envDefault:"localhost"`
	Port    string `env:"SERVER_PORT" envDefault:"8080"`

	SSLCertPath string `env:"SSL_CERT_PATH" envDefault:"cert/localhost.crt"`
	SSLKeyPath  string `env:"SSL_KEY_PATH" envDefault:"cert/localhost.key"`

	DSN string `env:"DSN" envDefault:"postgresql://root:root@localhost:5432/gophkeeper?sslmode=disable"`
}

var cfg Config

func NewConfig() Config {
	cfg.parse()

	return cfg
}

func (c *Config) parse() {
	err := env.Parse(&cfg)
	if err != nil {
		pkg.NewLogger().Error(err.Error())
	}
}