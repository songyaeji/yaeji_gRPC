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

//send_message

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
