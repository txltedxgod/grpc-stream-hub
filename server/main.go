package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type MessageRequest struct {
	RoomId    string
	Sender    string
	Content   string
	Timestamp int64
}

type MessageResponse struct {
	MessageId string
	RoomId    string
	Sender    string
	Content   string
	Timestamp int64
}

func main() {
	port := flag.Int("port", 50051, "gRPC server listen port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", *port, err)
	}

	hub := NewHub()
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(AuthStreamInterceptor),
	)

	log.Printf("[gRPC Hub] Listening for streaming connections on :%d\n", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
