package client

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Shelex/grpc-go-demo/client/graph"
	"github.com/Shelex/grpc-go-demo/client/graph/generated"
	"github.com/Shelex/grpc-go-demo/config"
	"github.com/Shelex/grpc-go-demo/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Start(cfg config.Config, conn *grpc.ClientConn) error {
	log.Println("starting graphql server...")

	client := proto.NewEmployeeServiceClient(conn)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(client)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://%s:%s/ for GraphQL playground", cfg.Host, cfg.ClientAPIPort)
	return http.ListenAndServe(":"+cfg.ClientAPIPort, nil)
}

func ConnectGRPCService(cfg config.Config) (*grpc.ClientConn, error) {
	log.Println("connecting to grpc...")
	creds, err := credentials.NewClientTLSFromFile(cfg.PathToTLSCertFile, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse tls credentials: %w", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial(cfg.Host+":"+cfg.DomainServicePort, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}
	return conn, nil
}
