package proto

//go:generate protoc --proto_path=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. event.proto
