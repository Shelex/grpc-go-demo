package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shelex/grpc-go-demo/config"
	"github.com/Shelex/grpc-go-demo/logger"
	"github.com/Shelex/grpc-go-demo/proto"
	"github.com/Shelex/grpc-go-demo/storage"
	"github.com/Shelex/grpc-go-demo/storage/documents"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewEmployeeService(cfg config.Config, ctx context.Context) (*employeeService, error) {
	creds, err := credentials.NewServerTLSFromFile(cfg.PathToTLSCertFile, cfg.PathToTLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get tls credentials: %w", err)
	}

	var repo storage.Storage
	var dbErr error

	switch cfg.Env {
	case "test":
		repo, dbErr = storage.NewInMemStorage()
	default:
		repo, dbErr = storage.NewMongoStorage(cfg.MongoURL)
	}
	if dbErr != nil {
		return nil, fmt.Errorf("failed to connect to repository: %w", dbErr)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := ConnectGRPCWithContext(ctx, opts...)

	srv := &employeeService{
		repository: repo,
		documents:  documents.NewLocalFS(),
		grpc:       s,
	}

	return srv, nil
}

func Start(ctx context.Context) {
	cfg := config.GetEnv()

	logger.Init(cfg)

	log.Info("starting service...")

	log.Printf("connecting to %s env", cfg.Env)
	srv, err := NewEmployeeService(cfg, ctx)
	if err != nil {
		log.Errorf("failed to create employee service: %s", err)
		os.Exit(1)
	}

	proto.RegisterEmployeeServiceServer(srv.grpc, srv)

	lis, err := net.Listen("tcp", ":"+cfg.DomainServicePort)
	if err != nil {
		log.Errorf("failed to connect tcp on port %s: %s", cfg.DomainServicePort, err)
		os.Exit(1)
	}

	log.Printf("starting server on port %s", cfg.DomainServicePort)
	if err := srv.grpc.Serve(lis); err != nil {
		log.Errorf("failed to serve http: %s", err)
		os.Exit(1)
	}
}

func ConnectGRPCWithContext(ctx context.Context, opts ...grpc.ServerOption) *grpc.Server {
	server := grpc.NewServer(opts...)
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	return server
}

func SetupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("got interrupt signal")
		stop()
	}()
}
