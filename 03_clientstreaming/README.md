# Client Streaming gRPC 실행 가이드

## 서버 실행 방법

```bash
cd yaeji_gRPC/03_clientstreaming
go build
./03_clientstreaming.exe
```

서버는 포트 50051에서 실행됩니다.

---

## 클라이언트 실행 방법

```bash
cd yaeji_gRPC/03_clientstreaming
go build -tags client -o client.exe
./client.exe
```

---

## 전체 실행 순서

1. **터미널 1에서 서버 실행:**
   ```bash
   cd yaeji_gRPC/03_clientstreaming
   go run .
   ```

2. **터미널 2에서 클라이언트 실행:**
   ```bash
   cd yaeji_gRPC/03_clientstreaming
   go run -tags client .
   ```

---

## 빌드 태그 설명

- `server.go`: `//go:build !client` - 기본 빌드 시 포함 (client 태그가 없을 때)
- `client.go`: `//go:build client` - `-tags client` 플래그 사용 시에만 포함

이를 통해 같은 패키지에서도 독립적으로 실행 가능합니다.

