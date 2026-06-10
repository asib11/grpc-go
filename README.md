# gRPC in Go

A hands-on demonstration of **gRPC** communication patterns implemented in Go, covering Unary RPC, Server-Side Streaming, Client-Side Streaming, and Bidirectional Streaming.

---

## Table of Contents

- [What is gRPC?](#what-is-grpc)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Protocol Buffers (`.proto`)](#protocol-buffers-proto)
- [Communication Patterns](#communication-patterns)
  - [1. Unary RPC](#1-unary-rpc)
  - [2. Server-Side Streaming](#2-server-side-streaming)
  - [3. Client-Side Streaming](#3-client-side-streaming)
  - [4. Bidirectional Streaming](#4-bidirectional-streaming)
- [Running the Project](#running-the-project)
- [Dependencies](#dependencies)

---

## What is gRPC?

**gRPC** (Google Remote Procedure Call) is a high-performance, open-source RPC framework that runs over **HTTP/2**. It uses **Protocol Buffers** (protobuf) as its Interface Definition Language (IDL) for defining service contracts and serializing structured data.

Key advantages over traditional REST:

- Strongly typed contracts via `.proto` files
- Efficient binary serialization (smaller payload than JSON)
- Built-in support for streaming (server, client, and bidirectional)
- Auto-generated client and server stubs in multiple languages

---

## Project Structure

```
grpc-go/
├── client/          # gRPC client implementations
│   └── main.go
├── server/          # gRPC server implementations
│   └── main.go
├── protoc/          # Protocol Buffer definitions and generated code
│   ├── *.proto      # Service definitions
│   ├── *.pb.go      # Generated protobuf message code
│   └── *_grpc.pb.go # Generated gRPC service code
├── go.mod
└── go.sum
```

---

## Prerequisites

Make sure the following tools are installed:

- [Go 1.20+](https://go.dev/dl/)
- [Protocol Buffer Compiler (`protoc`)](https://grpc.io/docs/protoc-installation/)
- Go plugins for `protoc`:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Ensure `$GOPATH/bin` is in your `PATH`.

---

## Installation & Setup

```bash
# Clone the repository
git clone https://github.com/asib11/grpc-go.git
cd grpc-go

# Install dependencies
go mod tidy
```

### Regenerating Proto Files

If you modify any `.proto` file, regenerate the Go code by running:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       protoc/*.proto
```

---

## Protocol Buffers (`.proto`)

All service definitions live in the `protoc/` directory. A `.proto` file defines:

- **Messages** — the request and response data structures
- **Services** — the RPC methods exposed by the server

Example structure of a `.proto` file used in this project:

```proto
syntax = "proto3";

package grpc;
option go_package = "./protoc";

// Unary
rpc SayHello (HelloRequest) returns (HelloResponse);

// Server streaming
rpc ServerStream (HelloRequest) returns (stream HelloResponse);

// Client streaming
rpc ClientStream (stream HelloRequest) returns (HelloResponse);

// Bidirectional streaming
rpc BiDiStream (stream HelloRequest) returns (stream HelloResponse);
```

---

## Communication Patterns

### 1. Unary RPC

> **One request → One response**

The simplest and most familiar pattern — similar to a regular HTTP request/response cycle. The client sends a single request and waits for a single response from the server.

**Use cases:** Authentication, fetching a single record, simple computations.

**How it works:**

```
Client ──── Request ────► Server
Client ◄─── Response ─── Server
```

**Proto definition:**

```proto
rpc SayHello (HelloRequest) returns (HelloResponse);
```

**Server-side (Go):**

```go
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    return &pb.HelloResponse{Message: "Hello, " + req.Name}, nil
}
```

**Client-side (Go):**

```go
resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "World"})
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Println(resp.Message)
```

---

### 2. Server-Side Streaming

> **One request → Multiple responses (stream)**

The client sends a single request and the server responds with a stream of messages. The client reads from the stream until the server signals it is done.

**Use cases:** Live data feeds, sending large files in chunks, real-time notifications, paginated results.

**How it works:**

```
Client ──── Request ────► Server
Client ◄─── Response 1 ── Server
Client ◄─── Response 2 ── Server
Client ◄─── Response 3 ── Server
           ...
Client ◄─── (stream end) ─ Server
```

**Proto definition:**

```proto
rpc ServerStream (HelloRequest) returns (stream HelloResponse);
```

**Server-side (Go):**

```go
func (s *Server) ServerStream(req *pb.HelloRequest, stream pb.Service_ServerStreamServer) error {
    for i := 0; i < 5; i++ {
        err := stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Message %d", i)})
        if err != nil {
            return err
        }
        time.Sleep(time.Second)
    }
    return nil
}
```

**Client-side (Go):**

```go
stream, err := client.ServerStream(ctx, &pb.HelloRequest{Name: "World"})
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Println(resp.Message)
}
```

---

### 3. Client-Side Streaming

> **Multiple requests (stream) → One response**

The client sends a stream of messages to the server. Once the client finishes sending, the server processes them all and returns a single response.

**Use cases:** Uploading files in chunks, aggregating sensor readings, batch data ingestion.

**How it works:**

```
Client ──── Request 1 ───► Server
Client ──── Request 2 ───► Server
Client ──── Request 3 ───► Server
           ...
Client ──── (stream end) ► Server
Client ◄─── Response ───── Server
```

**Proto definition:**

```proto
rpc ClientStream (stream HelloRequest) returns (HelloResponse);
```

**Server-side (Go):**

```go
func (s *Server) ClientStream(stream pb.Service_ClientStreamServer) error {
    var names []string
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.HelloResponse{
                Message: "Received: " + strings.Join(names, ", "),
            })
        }
        names = append(names, req.Name)
    }
}
```

**Client-side (Go):**

```go
stream, _ := client.ClientStream(ctx)
for _, name := range []string{"Alice", "Bob", "Charlie"} {
    stream.Send(&pb.HelloRequest{Name: name})
}
resp, err := stream.CloseAndRecv()
fmt.Println(resp.Message)
```

---

### 4. Bidirectional Streaming

> **Multiple requests (stream) ↔ Multiple responses (stream)**

Both the client and the server send independent streams of messages to each other simultaneously. The two streams operate independently — the server doesn't have to wait for all client messages before responding, and vice versa.

**Use cases:** Chat applications, multiplayer games, collaborative editing, real-time telemetry with acknowledgment.

**How it works:**

```
Client ──── Request 1 ───► Server
Client ◄─── Response A ─── Server
Client ──── Request 2 ───► Server
Client ◄─── Response B ─── Server
           ...  (concurrent, full-duplex)
```

**Proto definition:**

```proto
rpc BiDiStream (stream HelloRequest) returns (stream HelloResponse);
```

**Server-side (Go):**

```go
func (s *Server) BiDiStream(stream pb.Service_BiDiStreamServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        stream.Send(&pb.HelloResponse{
            Message: "Echo: " + req.Name,
        })
    }
}
```

**Client-side (Go):**

```go
stream, _ := client.BiDiStream(ctx)

// Send messages concurrently
go func() {
    for _, name := range []string{"Alice", "Bob", "Charlie"} {
        stream.Send(&pb.HelloRequest{Name: name})
    }
    stream.CloseSend()
}()

// Receive responses
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Println(resp.Message)
}
```

---

## Running the Project

**Start the server:**

```bash
go run server/main.go
```

**Run the client (in a separate terminal):**

```bash
go run client/main.go
```

The server listens on `:50051` by default.

---

## Dependencies

| Package | Version | Purpose |
|---|---|---|
| `google.golang.org/grpc` | latest | gRPC runtime for Go |
| `google.golang.org/protobuf` | latest | Protocol Buffers runtime |
| `google.golang.org/grpc/cmd/protoc-gen-go-grpc` | latest | Code generator plugin |

Install with:

```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

---

## Summary of gRPC Patterns

| Pattern | Request | Response | Use Case |
|---|---|---|---|
| **Unary** | Single | Single | CRUD, auth, simple queries |
| **Server Streaming** | Single | Stream | Live feeds, large downloads |
| **Client Streaming** | Stream | Single | File uploads, batch ingestion |
| **Bidirectional** | Stream | Stream | Chat, real-time collaboration |

---

> Built with ❤️ using [Go](https://go.dev) and [gRPC](https://grpc.io)
