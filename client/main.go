package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	serverAddr := flag.String("server", "localhost:50051", "gRPC server target address")
	room := flag.String("room", "general", "Chat room / topic name")
	username := flag.String("user", "dev", "User identifier")
	token := flag.String("token", "secret-token", "Authentication token")
	flag.Parse()

	log.Printf("[Client] Connecting to gRPC Stream Hub on %s (Room: %s, User: %s)...\n", *serverAddr, *room, *username)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-api-token", *token)
	conn, err := grpc.DialContext(ctx, *serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	fmt.Printf("[+] Connected! Type messages and press ENTER to broadcast:\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if text == "/exit" {
			break
		}
		fmt.Printf("[%s] Sent: %s\n", time.Now().Format("15:04:05"), text)
	}
}
