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

// 반환값: 오류 발생 -> error 반환, 정상 종료 -> nil 반환 (이거 담당 함수임~~)
func (s *BidirectionalService) MyFunction(stream grpc.BidiStreamingServer[Message, Message]) error {
	// 서버가 bidirectional streaming을 처리하기 시작했음을 출력
	fmt.Println("Server processing gRPC bidirectional streaming.")

	// 무한 루프: 클라이언트로부터 메시지를 계속 받기 위해 사용
	// go언어에서 for {} 는 조건 없음 -> 무한 루프 만들기 가능!!
	// break 문이나 return 문으로만 빠져나올 수 있음!!
	for {
		// stream.Recv(): 클라이언트로부터 메시지를 받음
		// message(받은 메시지)와 err(에러) 두 개의 값을 반환 -> go언어는 다중 값 반환 가능!!
		message, err := stream.Recv()

		// 에러 처리: err가 nil이 아니면 (=에러 발생 시)
		// err가 nil이 아닌 경우: 클라이언트가 스트림을 닫았거나, 네트워크 오류가 발생한 경우
		if err != nil {
			// break: for 루프 빠져나오기
			// 정상적인 경우(스트림 종료)와 오류 경우 모두 여기서 루프 종료
			break
		}

		// stream.Send(): 받은 메시지를 그대로 클라이언트한테 다시 보내기
		// 에러가 발생 -> error 반환 -> 함수 종료
		if err := stream.Send(message); err != nil {
			// return err: 에러 반환 -> 함수 즉시 종료
			return err
		}
	// 정상적으로 루프가 종료되면 nil(에러 없음)을 반환 -> 함수를 종료
	return nil
}

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
