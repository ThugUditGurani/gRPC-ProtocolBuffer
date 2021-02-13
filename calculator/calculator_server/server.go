package main

import (
	"Github/gRPC-ProtocolBuffer/calculator/calculatorpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type calculate struct {

}

func (c calculate) SquareRoot(ctx context.Context, request *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := request.GetNumber()
	if number < 0 {
		return nil,status.Errorf(codes.InvalidArgument,"Received a negative number")
	}
	return &calculatorpb.SquareRootResponse{NumberRoot: math.Sqrt(float64(number))}, nil
}

func (c calculate) ComputeAverage(server calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Received Cmpute Average RPC")
	sum := int32(0)
	count := 0

	for{
		req, err := server.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return server.SendAndClose(&calculatorpb.ComputeAverageResponse{Average: average})
		}
		if err != nil {
			log.Fatalf("Error while reading client server: %v",err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (c calculate) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberDecompositionRequest, server calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecompositon RPC: %v",request)
	number := request.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			server.Send(&calculatorpb.PrimeNumberDecompositionResponse{PrimeFactor: divisor,
				})
			number = number/divisor
		}else {
			divisor++
			fmt.Printf("Divisor has increased to %v",divisor)
		}
	}
	return nil
}

func (c calculate) Calculator(ctx context.Context, request *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	firstNumber := request.GetCalculator().FirstNumber
	secondNumber := request.GetCalculator().SecondNumber

	result := firstNumber + secondNumber
	res := &calculatorpb.CalculatorResponse{
	Sum: result,
	}
	return res , nil
}

func main() {

	con , err := net.Listen("tcp","0.0.0.0:50051")
	if err != nil {
		log.Fatalf("can't Connect with sercer %v",con)
	}
	grpcConn := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(grpcConn,&calculate{})
	if err := grpcConn.Serve(con); err != nil {
		log.Fatalf("cant connect with rpc: %v",err)
	}
}
