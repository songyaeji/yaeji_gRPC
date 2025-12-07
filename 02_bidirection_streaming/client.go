//go:build client
// +build client

// Bidirectional streaming gRPC client

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

// client: gRPC 클라이언트 객체 (서버와 통신하기 위한 인터페이스)
// ctx: context 객체 (요청의 생명주기와 취소를 관리하는 go언어의 표준 패턴)
func send_message(client BidirectionalClient, ctx context.Context) {
	// client.MyFunction(ctx): 서버와의 bidirectional streaming 연결을 생성
	// stream: 서버와 주고받을 수 있는 스트림 객체
	// err: 연결 생성 시 발생할 수 있는 에러
	// go언어는 다중 값 반환 가능!! -> server.go에서 한 거랑 똑같음
	stream, err := client.MyFunction(ctx)
	if err != nil {
		// log.Fatalf: 에러를 출력하고 프로그램을 즉시 종료 (치명적 오류)
		// %v: go언어의 기본 값 포맷터 (에러 메시지를 출력)
		log.Fatalf("MyFunction call failed: %v", err)
	}

	// generate_messages(): 전송할 메시지 리스트를 생성
	// messages: Message 객체들의 배열
	messages := generate_messages()

	// done := make(chan bool): boolean 값을 주고받을 수 있는 채널을 생성
	// 채널(Channel): go언어의 동시성 프로그래밍에서 goroutine 간 통신을 위한 도구
	// bool 타입: true 또는 false 값을 전달할 수 있음
	done := make(chan bool)

	// go func() { ... }(): goroutine을 생성하여 익명 함수를 병렬로 실행
	// goroutine(중요 별 다섯개): go언어의 경량 스레드 -> 다른 코드와 동시에 실행될 수 있음
	// 여기서는 메시지를 보내는 작업을 별도의 goroutine에서 실행 -> 메시지를 받는 작업과 동시에 진행할 수 있게 함 -> 병렬 처리
	go func() {
		// for _, msg := range messages: messages 배열의 각 메시지를 순회
		// _: 인덱스는 사용하지 않으므로 _로 무시 -> 변수 선언 후 사용 안 함
		for _, msg := range messages {
			// stream.Send(msg): 서버로 메시지를 전송
			// if err := ...: Send()의 반환값이 에러인지 확인
			// err != nil이면 전송 실패
			if err := stream.Send(msg); err != nil {
				// 전송 실패 시 에러를 출력하고 프로그램 종료 (치명적 오류)
				log.Fatalf("Failed to send message: %v", err)
			}
		}
		// stream.CloseSend(): 더 이상 보낼 메시지가 없음을 서버에 알림
		// 서버는 이 신호를 받으면 클라이언트가 보내는 메시지 수신을 중단
		stream.CloseSend()
		// done <- true: 채널에 true 값을 보냄 (작업 완료 신호)
		// <-: 채널에 값을 보내는 연산자 -> 송신
		done <- true
	}() // goroutine 함수 종료 -> 메시지 전송 작업 완료

	// 무한 루프: 서버로부터 메시지를 계속 받기 위해 사용
	// for {}: 조건 없음 -> 무한 루프 만들기 가능!!
	// break 문으로만 빠져나올 수 있음!!
	for {
		response, err := stream.Recv()
		if err != nil {
			// 에러가 발생하면 (스트림 종료 또는 네트워크 오류) 루프를 종료
			break
		}
		fmt.Printf("[server to client] %s\n", response.Message)
	}
	// <-done: 채널에서 값을 받을 때까지 대기 -> 수신
	// goroutine이 done 채널에 값을 보낼 때까지 여기서 멈춤 -> 메시지 전송 goroutine이 완료될 때까지 기다림
	<-done
}

func run() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewBidirectionalClient(conn)
	ctx := context.Background()
	send_message(client, ctx)
}

func main() {
	run()
}
