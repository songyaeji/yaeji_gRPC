package main

import (
	context "context"
	fmt "fmt"
	log "log"

	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials/insecure"
)

func main() {
	// (1) grpc 모듈을 import 함
	// (Go에서는 import 문으로 이미 위에 선언되어 있음)

	// (2) protoc가 생성한 클래스를 import 함
	// (Go에서는 같은 package 내에 생성되어 있으므로 별도 import 불필요)

	// (3) gRPC 통신 채널을 생성함
	// Go에서는 grpc.Dial()을 사용하여 연결을 생성함
	// insecure 채널은 credentials.NewCredentials()를 사용함
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(credentials.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// (4) protoc가 생성한 _grpc.pb.go 화일의 stub 함수를, (3)의 채널을, 사용하여 실행하여 stub를 생성함
	// Go에서는 NewMyServiceClient() 함수를 사용하여 클라이언트(stub)를 생성함
	client := NewMyServiceClient(conn)

	// (5) protoc가 생성한 .pb.go 화일의 메세지 타입에 맞춰서, 원격 함수에 전달할 메시지를 만들고, 전달할 값을 저장함
	request := &MyNumber{
		Value: 4,
	}

	// (6) 원격 함수를 stub을 사용하여 호출함
	ctx := context.Background()
	response, err := client.MyFunction(ctx, request)
	if err != nil {
		log.Fatalf("MyFunction call failed: %v", err)
	}

	// (7) 결과를 활용하는 작업을 수행함 [optional]
	fmt.Printf("gRPC result: %d\n", response.Value)
}
