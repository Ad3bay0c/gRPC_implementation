package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type server struct {}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}
func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "0.0.0.0:50051"
	}

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err.Error())
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	log.Printf("Server Started on port: %v\n", port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err.Error())
	}
}
