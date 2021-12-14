package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	log.Printf("FindMaximum Server Invoked\n")

	var max int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while collecting value from client: %v", err)
		}
		if req.GetNumber() > max {
			max = req.GetNumber()
			err := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})
			if err != nil {
				log.Fatalf("Error while sending value to client: %v", err)
			}
		}
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Received a negative number: %v", number)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err.Error())
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)
	log.Printf("Server Started\n")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to Start Server: %v", err.Error())
	}
}
