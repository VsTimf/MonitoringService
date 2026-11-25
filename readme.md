# MonitorService

## Суть проекта

## Зависимости (libs)

1. "github.com/lmittmann/tint"
2. "github.com/shirou/gopsutil/v3/cpu"
3. "github.com/shirou/gopsutil/v3/mem"
4. "modernc.org/sqltie"
5. google.golang.org/grpc/reflection

## Зависимости (tools)

1. Компилятор protoc
sudo dnf install -y protobuf-compiler

2. Плагин для генерации структур данных Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

3. Плагин для генерации gRPC-сервера и клиента
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

4. Утилита для тестирования gRPC
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
