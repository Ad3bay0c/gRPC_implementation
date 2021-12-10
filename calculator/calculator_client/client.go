package main

import (
	"context"
	"fmt"
	"github.com/Ad3bay0c/gRPC/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err.Error())
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	// unary
	//doUnary(c)

	//server streaming
	//doServerStreaming(c)

	//client streaming
	//doClientStreaming(c)

	//BiDi streaming
	doBiDiStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
        FirstNumber:  5,
        SecondNumber: 10,
    }
    res, err := c.Sum(context.Background(), req)
    if err != nil {
        log.Fatalf("Error while calling Sum RPC: %v", err.Error())
    }
    log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	res, err := c.PrimeNumber(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("An Error Occurred: %v", err)
		}
		log.Printf("Factor is: %v", msg.PrimeNumber)
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
        log.Fatalf("Error while calling ComputeAverage: %v", err.Error())
    }
	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
            Number: 1,
        },
		&calculatorpb.ComputeAverageRequest{
            Number: 2,
        },
		&calculatorpb.ComputeAverageRequest{
            Number: 3,
        },
		&calculatorpb.ComputeAverageRequest{
            Number: 4,
        },
	}

	for _, req := range requests {
		fmt.Printf("Sending Request %v\n", req)
		stream.Send(req)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err.Error())
	}
	fmt.Printf("The total Average is: %v\n", res.Result)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling FindMaximum: %v", err.Error())
	}

	numbers := []int32{1,5,3,6,2,20}
	wait := make(chan struct{})
	// send a bunch of messages to the server
	go func() {
		for _, n := range numbers {
			fmt.Printf("Sending number: %v\n", n)
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: n,
			})
			if err != nil {
				log.Fatalf("Error while sending message: %v", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of messages from the server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error while receiving: %v", err)
				break
			}
			log.Printf("Received: %v", res.GetMaximum())
		}
		close(wait)
	}()

	<- wait
}