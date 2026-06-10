package main

import (
	"context"
	"fmt"
	proto "grpc/protoc"
	"net/http"
	"encoding/json"
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
	mux.HandleFunc("GET /sent", clientConnectionServer)
	fmt.Println("HTTP server is running on port 8081...")
	http.ListenAndServe(":8081", mux)
}

func clientConnectionServer(w http.ResponseWriter, r *http.Request) {
	req := []*proto.HelloRequest{
		{SomeString: "Request 1"},
		{SomeString: "Request 2"},
		{SomeString: "Request 3"},
		{SomeString: "Request 4"},
		{SomeString: "Request 5"},
		{SomeString: "Request 6"},
	}

	stream, err := client.ServerReply(context.TODO())
	if err != nil {
		http.Error(w, "Error calling ServerReply: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, r := range req {
		if err := stream.Send(r); err != nil {
			http.Error(w, "Error sending request: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		http.Error(w, "Error receiving response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"message count": resp})
}