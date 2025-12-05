//go:build !client
// +build !client

// Bidirectional streaming gRPC server

package main

import (
	fmt "fmt"
	net "net"

	grpc "google.golang.org/grpc"
)

type BidirectionalService struct {
	UnimplementedBidirectionalServer
}

/*
func MyFunction
*/

func serve() {
	server := grpc.NewServer()
	RegisterBidirectionalServer(server, &BidirectionalService{})
	fmt.Println("Starting server. Listening on port 50051.")
	lis, _ := net.Listen("tcp", ":50051")
	server.Serve(lis)
}

func main() {
	serve()
}
