# Metrigo

## Prerequisites

- Go `1.24.5`

- Protocol Buffers Compiler

  - Download from <https://github.com/protocolbuffers/protobuf/releases>
  - Install and ensure it's in your path

- Protoc plugins

    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

## Building

- Download module dependencies

  ```bash
  go mod download
  ```

- Generate protobuf files

  ```bash
  protoc \
  --go_out=paths=source_relative:. \ 
  --go-grpc_out=paths=source_relative:. \ 
  pb/metrigo.proto
  ```

- Build the binary

  ```bash
  go build ./cmd/main.go
  ```
