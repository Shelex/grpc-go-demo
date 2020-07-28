package logger

import (
	"os"

	"github.com/Shelex/grpc-go-demo/config"
	log "github.com/sirupsen/logrus"
)

func Init(cfg config.Config) {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: cfg.PrettyLogOutput})
	log.SetOutput(os.Stdout)
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Errorf("failed to set log level: %s", err)
	}
	log.SetLevel(logLevel)
}
