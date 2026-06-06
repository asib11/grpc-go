package main

import (
	"context"
	"errors"
	"fmt"

	proto "grpc/protoc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedExampleServer
}

func (s *server) ServerReply(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {

	fmt.Printf("Received request with SomeString: %s\n", req.SomeString)
	fmt.Printf("Hello from the server!\n")
	return &proto.HelloResponse{}, errors.New("")
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterExampleServer(s, &server{})
	reflection.Register(s)

	fmt.Println("Server is running on port 9000...")

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
