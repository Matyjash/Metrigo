# Metrigo

Metrigo is a standalone service with gRPC server for getting host metrics and general information such as:

- MemoryUsage
- CPU specs & usage
- Temperature sensors values
- General host info (hostname, os, uptime, etc.)
- Net specs (active interfaces)

## Usage

### CLI

Each metric can be obtained by adequate parameter given as an argument to the binary:

```sh
> ./metrigo host
< Host info:
Hostname: Dev
OS: linux
Platform: ubuntu
Platform version: 24.04
Uptime (s): 6785
```

You can get the full list of possible arguments with:

```sh
> ./metrigo --help
```

### gRPC server

Metrgio can be run as a server when launched with `server` flag on a port `50051` by deafult:

```sh
> ./main --server
< Server is running on port :50051
```

See the available services in the [protobuf file](./pb/metrigo.proto)

## Building

### Prerequisites

- [Go](https://go.dev/doc/install) installed

- Protocol Buffers Compiler

  - Download [protoc](https://github.com/protocolbuffers/protobuf/releases)
  - Install and ensure it's in your path

- Install protoc plugins

    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

### Building project

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
