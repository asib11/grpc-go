package main

import (
	"fmt"
	"io"
	"strconv"

	proto "grpc/protoc"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedExampleServer
}

func (s *server) ServerReply(stream proto.Example_ServerReplyServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.HelloResponse{
				Reply: strconv.Itoa(count),
			})
		}
		if err != nil {
			return err
		}
		count++
		fmt.Println(req)
	}
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
