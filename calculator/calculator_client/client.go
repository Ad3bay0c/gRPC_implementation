package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err.Error())
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)
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

