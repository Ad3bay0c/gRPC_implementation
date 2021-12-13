package main

import (
	"context"
	"fmt"
	"github.com/Ad3bay0c/gRPC/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" // certificate Authority Trust Certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr.Error())
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)
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
	//doClientStreaming(c)

	// Bi-directional Streaming RPC Implementation
	//doBiDirectionalStreaming(c)

	// Unary With Deadline RPC Implementation
	doUnaryWithDeadline(c, 5 * time.Second) // should complete
	//doUnaryWithDeadline(c, 1 * time.Second) // should timeout
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

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName: "Doe",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Printf("Timeout was hit! Deadline was exceeded\n")
			} else {
				fmt.Printf("Unexpected Error: %v", statusErr)
			}
		}else {
			log.Fatalf("Error while calling GreetWithDeadline RPC: %v", err)
		}
		return
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

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
				},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
				LastName:  "Doe",
			},
				},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
				LastName:  "Doe",
			},
				},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jill",
				LastName:  "Doe",
			},
				},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Smith",
			},
				},
	}
	// We create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err.Error())
	}
	waitc := make(chan struct{})
	// We send a bunch of messages to the client (goroutine)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
            stream.Send(req)
            time.Sleep(1000 * time.Millisecond)
        }
		stream.CloseSend()
	}()
	// We receive a bunch of messages from the client (goroutine)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error while receiving: %v", err.Error())
				break
			}
			log.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
}