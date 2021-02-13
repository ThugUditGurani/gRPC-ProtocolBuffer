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
)

func main() {
	conn , err := grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v",err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(c)

	//doServerStreaming(c)
	//doClientStreaming(c)

	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Square root Unary RPC")
	n := 10
	res , err := c.SquareRoot(context.Background(),&calculatorpb.SquareRootRequest{Number: int32(n)})
	if err != nil {
		status , ok := status.FromError(err)
		if ok {
			fmt.Println(status.Message())
			fmt.Println(status.Code())
			if status.Code() == codes.InvalidArgument {
				fmt.Println("we probalbly sent a negative number")
			}
		}else {
			log.Fatalf("Big Error calling squareroot: %v",err)
		}
	}
	fmt.Printf("Result of square root of %v : %v",n,res.GetNumberRoot())
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a computeAverage Client Streaming RPC ..")
	stream , error := c.ComputeAverage(context.Background())
	if error != nil{
		log.Fatalf("Error while opening stream %v",error)
	}
	numbers := []int32{3,5,9,54,23}
	for _,number  := range numbers {
		stream.Send(&calculatorpb.ComputeAverageRequest{Number: number})
	}
	res,err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving number %v",error)
	}
	fmt.Printf("The Average is : %v",res.GetAverage())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{Number: 1200}
	stream , err := c.PrimeNumberDecomposition(context.Background(),req)
	if err != nil {
		log.Fatalf("Error while calling Prime Decompositon %v",err)
	}
	for {
		res,err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened %v",err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalculatorRequest{Calculator: &calculatorpb.Calculator{
		FirstNumber:  1,
		SecondNumber: 2,
	}}

	res , err := c.Calculator(context.Background(),req)
	if err != nil {
		log.Fatalf("Error in res: %v",err)
	}
	log.Printf("Response : %v",res)
}
