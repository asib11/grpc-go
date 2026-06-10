package main

import (
	"errors"
	"fmt"
	"io"

	proto "grpc/protoc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedExampleServer
}

func (s *server) ServerReply(stream proto.Example_ServerReplyServer) error {
	for i := 0; i < 5; i++ {
		err := stream.Send(&proto.HelloResponse{Reply: fmt.Sprintf("Hello from server %d", i+1)})

		if err != nil {
			return errors.New("Error sending message: " + err.Error())
		}

	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("Error receiving message: " + err.Error())
		}
		fmt.Println(req.SomeString)
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
