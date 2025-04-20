gen:
	protoc --go_out=./internal/pb --go-grpc_out=./internal/pb ./proto/file_transfer.proto

client:
	go build -o client ./cmd/client/main.go