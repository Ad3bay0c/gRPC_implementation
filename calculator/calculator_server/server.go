package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)
type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	return &calculatorpb.SumResponse{
		SumResult: req.GetFirstNumber() + req.GetSecondNumber(),
	}, nil
}

func (*server) PrimeNumber(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberServer) error {
	n := req.Number
	var k int32 = 2
	for n > 1 {
		if n % k == 0 {
			result := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumber: k,
			}
			err := stream.Send(result)
			n = n / k
			if err != nil {
                return err
            }
		} else {
			k = k +	1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	var average int64
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float64(average) / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("Error while collecting value from client")
		}
		average += req.Number
		count++
	}
	return nil
}
func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err.Error())
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	log.Printf("Server Started\n")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Start Server: %v", err.Error())
	}
}
