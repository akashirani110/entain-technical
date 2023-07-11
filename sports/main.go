package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9001", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9001")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
