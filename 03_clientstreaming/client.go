//go:build client
// +build client

// Client streaming gRPC client

package main

import (
	context "context"
	fmt "fmt"
	log "log"

	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials/insecure"
)

func make_message(message string) *Message {
	return &Message{
		Message: message,
	}
}

func generate_messages() []*Message {
	messages := []*Message{
		make_message("message #1"),
		make_message("message #2"),
		make_message("message #3"),
		make_message("message #4"),
		make_message("message #5"),
	}
	for _, msg := range messages {
		fmt.Printf("[client to server] %s\n", msg.Message)
	}
	return messages
}

func send_message(client ClientStreamingClient, ctx context.Context) {
	stream, err := client.MyFunction(ctx)
	if err != nil {
		log.Fatalf("MyFunction call failed: %v", err)
	}

	messages := generate_messages()
	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}
	fmt.Printf("[server to client] %d\n", response.Value)
}

func run() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewClientStreamingClient(conn)
	ctx := context.Background()
	send_message(client, ctx)
}

func main() {
	run()
}
