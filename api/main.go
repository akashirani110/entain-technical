package main

import (
	"context"
	"flag"
	"git.neds.sh/matty/entain/api/proto/sports"
	"log"
	"net/http"

	"git.neds.sh/matty/entain/api/proto/racing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var (
	apiEndpoint        = flag.String("api-endpoint", "localhost:8000", "API endpoint")
	grpcEndpointRacing = flag.String("grpc-endpoint-racing", "localhost:9000", "gRPC server endpoint")
	grpcEndpointSports = flag.String("grpc-endpoint-sports", "localhost:9001", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("failed running api server: %s\n", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	// register the racing service api endpoint
	if err := racing.RegisterRacingHandlerFromEndpoint(
		ctx,
		mux,
		*grpcEndpointRacing,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	// register the new sports service api endpoint
	if err := sports.RegisterSportsHandlerFromEndpoint(
		ctx,
		mux,
		*grpcEndpointSports,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	log.Printf("API server listening on: %s\n", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, mux)
}
