package main

import (
	"context"
	"fmt"
	"github.com/Ad3bay0c/gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err.Error())
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	// Unary RPC Implementation
	//doUnary(c)

	// Server Streaming RPC Implementation
	//doServerStreaming(c)

	// Client Streaming RPC Implementation
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName: "Doe",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err.Error())
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName: "Doe",
		},
	}
	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("An error occurred: %v", err.Error())
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("An Error occurred: %v", err.Error())
		}
		log.Printf("Response : %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err.Error())
	}
	requests := []*greetpb.LongGreetRequest{
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "John",
                LastName: "Doe",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "Jane",
                LastName: "Doe",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "Jack",
                LastName: "Doe",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "Jill",
                LastName: "Doe",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "John",
                LastName: "Smith",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "Jane",
                LastName: "Smith",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "Jack",
                LastName: "Smith",
            },
        },
        &greetpb.LongGreetRequest{
            Greeting: &greetpb.Greeting{
                FirstName: "Jill",
                LastName: "Smith",
            },
        },
    }

	// We iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending Request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
        log.Fatalf("Error while receiving response from LongGreet: %v", err.Error())
    }
	fmt.Printf("LongGreet Response: %v\n", res)
}