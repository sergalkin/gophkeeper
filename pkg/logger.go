package pkg

import (
	"fmt"
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

func NewLogger() *zap.Logger {
	once.Do(func() {
		var l logger
		var zlConfig zap.Config

		err := env.Parse(&l)
		if err != nil {
			fmt.Println(err.Error())
		}

		if l.IsProd == true {
			zlConfig = zap.NewProductionConfig()
		} else {
			zlConfig = zap.NewDevelopmentConfig()
		}

		zlConfig.OutputPaths = []string{"stderr", l.LogPath}
		zl, err = zlConfig.Build()
		if err != nil {
			fmt.Println(err)
		}
	})

	return zl
}
