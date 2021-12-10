package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
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

	//streaming
	doServerStreaming(c)

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
