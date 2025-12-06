//go:build !client
// +build !client

// Server streaming gRPC server

package main

import (
	fmt "fmt"
	net "net"

	grpc "google.golang.org/grpc"
)

func make_message(message string) *Message {
	return &Message{
		Message: message,
	}
}

type ServerStreamingService struct {
	UnimplementedServerStreamingServer
}

func (s *ServerStreamingService) GetServerResponse(request *Number, stream grpc.ServerStreamingServer[Message]) error {
	messages := []*Message{
		make_message("message #1"),
		make_message("message #2"),
		make_message("message #3"),
		make_message("message #4"),
		make_message("message #5"),
	}
	fmt.Printf("Server processing gRPC server-streaming {%d}.\n", request.Value)
	for _, message := range messages {
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func serve() {
	server := grpc.NewServer()
	RegisterServerStreamingServer(server, &ServerStreamingService{})
	fmt.Println("Starting server. Listening on port 50051.")
	lis, _ := net.Listen("tcp", ":50051")
	server.Serve(lis)
}

func main() {
	serve()
}
