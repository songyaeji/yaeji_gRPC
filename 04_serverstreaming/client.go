//go:build client
// +build client

// Server streaming gRPC client

package main

import (
	context "context"
	fmt "fmt"
	log "log"

	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials/insecure"
)

func recv_message(client ServerStreamingClient, ctx context.Context) {
	request := &Number{Value: 5}
	stream, err := client.GetServerResponse(ctx, request)
	if err != nil {
		log.Fatalf("GetServerResponse call failed: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("[server to client] %s\n", response.Message)
	}
}

func run() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewServerStreamingClient(conn)
	ctx := context.Background()
	recv_message(client, ctx)
}

func main() {
	run()
}
