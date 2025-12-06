//go:build !client
// +build !client

// Server streaming gRPC server

package main

func make_message(message string) *Message {
	return &Message{
		Message: message,
	}
}

type ServerStreamingService struct {
	UnimplementedServerStreamingServer
}
