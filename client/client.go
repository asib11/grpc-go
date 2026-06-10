package main

import (
	"context"
	"encoding/json"
	"fmt"
	proto "grpc/protoc"
	"io"
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
	mux.HandleFunc("GET /sent", clientConnectionServer)
	fmt.Println("HTTP server is running on port 8081...")
	http.ListenAndServe(":8081", mux)
}

func clientConnectionServer(w http.ResponseWriter, r *http.Request) {

	stream, err := client.ServerReply(context.TODO())
	if err != nil {
		http.Error(w, "Error calling ServerReply: "+err.Error(), http.StatusInternalServerError)
		return
	}

	send, receive := 0, 0
	for i := 0; i < 10; i++ {
		err := stream.Send(&proto.HelloRequest{SomeString: fmt.Sprintf("Send %d from client", i+1)})
		if err != nil {
			http.Error(w, "Error sending message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		send++
	}

	if err := stream.CloseSend(); err != nil {
		http.Error(w, "Error closing send stream: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		fmt.Println("Received from server:", resp.Reply)
		receive++
	}

	response := map[string]interface{}{
		"message":http.StatusText(http.StatusOK),
		"messages_sent":     send,
		"messages_received": receive,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
