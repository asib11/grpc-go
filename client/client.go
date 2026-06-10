package main

import (
	"context"
	"encoding/json"
	"fmt"
	proto "grpc/protoc"
	"io"
	"net/http"
	"time"

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

	stream, err := client.ServerReply(context.TODO(), &proto.HelloRequest{SomeString: "Hello Server!"})
	if err != nil {
		http.Error(w, "Error calling ServerReply: "+err.Error(), http.StatusInternalServerError)
		return
	}

	count := 0
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error receiving stream: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Received from server:", message)
		time.Sleep(1 * time.Second)
		count++
	}

	response := map[string]interface{}{
		"message": "Finished receiving messages from server",
		"count":   count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
