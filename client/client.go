package main

import (
	"context"
	proto "grpc/protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	
	client := proto.NewExampleClient(conn)
	req := &proto.HelloRequest{SomeString: "Hello from the client!"}
	client.ServerReply(context.TODO(), req)

}