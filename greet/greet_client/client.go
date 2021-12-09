package main

import (
	"github.com/Ad3bay0c/gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err.Error())
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	log.Printf("Created Client: %f\n", c)
}
