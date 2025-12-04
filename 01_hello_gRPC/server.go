//go:build !client
// +build !client

package main

import (
	context "context"
	fmt "fmt"
	net "net"

	grpc "google.golang.org/grpc"
)

// (4) protoc가 생성한 Servicer 클래스를 base class로 해서 원격 호출될 함수들을 멤버로 갖는 서버 구조체를 생성함
type server struct {
	UnimplementedMyServiceServer
}

// (5) 서버 구조체에 원격 호출될 함수에 대한 rpc 함수를 작성함
// (5.1) proto 파일내 정의한 rpc 함수 이름에 대응하는 메서드를 작성함
func (s *server) MyFunction(ctx context.Context, request *MyNumber) (*MyNumber, error) {
	// (5.2) proto 파일내 message 이름과 동일한 message 구조체를 생성하여 응답 전달 용도로 사용함
	response := &MyNumber{}
	// (5.3) proto 파일내 message 이름과 동일한 message 구조체의 변수에 원격 함수의 수행 결과를 저장함
	// 앞서 (3)의 원격 호출할 함수에게 client로 부터 받은 입력 파라메타를 전달하고 결과를 가져옴
	response.Value = int32(my_func(int(request.Value)))
	// (5.4) 원격 함수 호출 결과를 client에게 돌려줌
	return response, nil
}

func main() {
	// (6) grpc.server를 생성함
	// Go에서는 grpc.NewServer()를 사용하여 서버를 생성함
	grpcServer := grpc.NewServer()

	// (7) RegisterMyServiceServer()를 사용해서, grpc.server에 (4)의 Servicer를 추가함
	RegisterMyServiceServer(grpcServer, &server{})

	// (8) grpc.server의 통신 포트를 열고, start()로 서버를 실행함
	lis, _ := net.Listen("tcp", ":50051")
	fmt.Println("Starting server. Listening on port 50051.")
	grpcServer.Serve(lis)

	// (9) grpc.server가 유지되도록 프로그램 실행을 유지함
	// server.Serve()가 블로킹되어 자동으로 유지됨
}
