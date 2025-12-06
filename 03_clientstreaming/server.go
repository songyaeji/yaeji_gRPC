//go:build !client
// +build !client

// Client streaming gRPC server

package main

import (
	fmt "fmt"
	net "net"

	grpc "google.golang.org/grpc"
)

type ClientStreamingServicer struct {
	UnimplementedClientStreamingServer
}

func (s *ClientStreamingServicer) MyFunction(stream grpc.ClientStreamingServer[Message, Number]) error {
	fmt.Println("Server processing gRPC client-streaming.")
	count := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		count++
	}
	return stream.SendAndClose(Number{Value: int32(count)})
}

func serve() {
	server := grpc.NewServer()
	RegisterClientStreamingServer(server, &ClientStreamingServicer{})
	fmt.Println("Starting server. Listening on port 50051.")
	lis, _ := net.Listen("tcp", ":50051")
	server.Serve(lis)
}

func main() {
	serve()
}
