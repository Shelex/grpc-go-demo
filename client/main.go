package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Shelex/grpc-go-demo/client/graph"
	"github.com/Shelex/grpc-go-demo/client/graph/generated"
	"github.com/Shelex/grpc-go-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultPort = "8080"
const serverPort = "9000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost:"+serverPort, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewEmployeeServiceClient(conn)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(client)}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
