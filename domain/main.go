package main

import (
	"context"

	"github.com/Shelex/grpc-go-demo/domain/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	service.SetupGracefulShutdown(cancel)
	service.Start(ctx)
}
