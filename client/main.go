package main

import (
	"os"

	client "github.com/Shelex/grpc-go-demo/client/service"
	"github.com/Shelex/grpc-go-demo/config"
	"github.com/Shelex/grpc-go-demo/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	cfg := config.GetEnv()

	logger.Init(cfg)
	log.Info("starting client...")

	conn, err := client.ConnectGRPCService(cfg)
	if err != nil {
		log.Errorf("failed to create connection to server: %s", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err := client.Start(cfg, conn); err != nil {
		log.Errorf("failed to start employee client: %s", err)
		os.Exit(1)
	}
}
