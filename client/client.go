package main

import (
	"context"
	"fmt"
	proto "grpc/protoc"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.ExampleClient

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	client = proto.NewExampleClient(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{message}", clientConnectionServer)
	fmt.Println("HTTP server is running on port 8081...")
	http.ListenAndServe(":8081", mux)
}

func clientConnectionServer(w http.ResponseWriter, r *http.Request) {
	message := r.PathValue("message")

	req := &proto.HelloRequest{SomeString: message}
	resp, err := client.ServerReply(context.Background(), req)
	if err != nil {
		http.Error(w, "Error calling ServerReply: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp.Reply))
}
