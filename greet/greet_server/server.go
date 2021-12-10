package main

import (
	"context"
	"fmt"
	"github.com/Ad3bay0c/gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type server struct {}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes Invoked\n")
	firstName := req.GetGreeting().GetFirstName()
	for i := 1; i<=10; i++ {
		result := "Hello " + firstName + " Number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// send final message
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error Wjile reading client stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Printf("Greet Everyone invoked\n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
            log.Fatalf("Error while reading client stream: %v", err)
            return err
        }
		result := "Hello " + req.GetGreeting().GetFirstName() + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
		}
	}
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
