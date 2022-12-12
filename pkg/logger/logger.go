package logger

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type logger struct {
	IsProd  bool   `env:"IS_PROD" envDefault:"false"`
	LogPath string `env:"LOG_PATH" envDefault:"./.tmp/log.txt"`
}

var (
	zl   *zap.Logger
	once sync.Once
)

// NewLogger - is wrap up for creating Logger.
//
// Only first call of this function is actually creating an instance of  *zap.Logger, all other calls
// returns already created instance of *zap.Logger.
//
// Realisation of Singleton pattern.
func NewLogger() *zap.Logger {
	once.Do(func() {
		var zlogger logger
		var zlConfig zap.Config

		err := env.Parse(&zlogger)
		if err != nil {
			log.Println(err.Error())
		}

		if zlogger.IsProd {
			zlConfig = zap.NewProductionConfig()
		} else {
			zlConfig = zap.NewDevelopmentConfig()
		}

		zlConfig.OutputPaths = []string{"stderr", zlogger.LogPath}
		zl, err = zlConfig.Build()
		if err != nil {
			log.Println(err.Error())
		}
	})

	return zl
}
