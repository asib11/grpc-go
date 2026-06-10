package main

import (
	"fmt"
	"time"

	proto "grpc/protoc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedExampleServer
}

func (s *server) ServerReply(req *proto.HelloRequest, stream proto.Example_ServerReplyServer) error {
	fmt.Println("Received request:", req.SomeString)
	time.Sleep(5 * time.Second)

	serverReply1 := []*proto.HelloResponse{
		{Reply: "Response 1"},
		{Reply: "Response 2"},
		{Reply: "Response 3"},
		{Reply: "Response 4"},
		{Reply: "Response 5"},
	}

	for _, r := range serverReply1 {
		if err := stream.Send(r)
		err != nil {
			return err
		}
	}
	return nil
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
