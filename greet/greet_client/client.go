package main

import (
	"Github/gRPC-ProtocolBuffer/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I am client")
	conn, err := grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v",err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBidirectionalStreaming(c)
	doUnaryWithDeadLine(c)
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "stephane"}},
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "stephaneOne"}},
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "stephaneTwo"}},
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{FirstName: "stephaneThree"}},
	}
	//we create a stream by invoking the client
	stream,error := c.GreetEveryone(context.Background())
	if error != nil {
		log.Fatalf("Error while streaming %v",error)
		return
	}
	waitc := make(chan struct{})
	//we send a bunch of message to the client (go routine)
	go func() {
		//function to send a bunch of message
		for _,req := range requests{
			fmt.Printf("Sending message: %v",req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//we receive a bunch of message from the client (go routine)
	go func() {
		//function to receive a bunch of message
		for {
			res,err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving %v",err)
				break
			}
			fmt.Printf("Receiving : %v",res.GetResult())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do  aclient Streaming RPC")
	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "stephane"}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "stephaneOne"}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "stephaneTwo"}},
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{FirstName: "stephaneThree"}},
	}
	stream , err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while connect %v",err)
	}
	for _,req := range request{
		log.Printf("Sending request : %v",req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res,err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receving %v",err)
	}
	fmt.Printf("Long greet Responses ",res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")
	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{LastName: "Mic", FirstName: "Stephane"}}

	resStream, err := c.GreetManyTimes(context.Background(),req)
	if err != nil {
		log.Fatalf("err %v",err)
	}

	for  {
		msg , err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v",err)
		}

		log.Printf("Response from GreetManyTImes: %v",msg.GetResult())

	}
}

func doUnary(c greetpb.GreetServiceClient){
	//fmt.Printf("Starting to do Unary Rpc client : %f",c)
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Mic",
		},
	}
	res,error := c.Greet(context.Background(),req)
	if error != nil{
		log.Fatalf("error while Calling Greet RPC %v",error)
	}
	log.Printf("Response from Greet: %v",res.Result)
}

func doUnaryWithDeadLine(c greetpb.GreetServiceClient){
	//fmt.Printf("Starting to do Unary Rpc client : %f",c)
	var seconds time.Duration = 5
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Mic",
		},
	}
	ctx,cancel := context.WithTimeout(context.Background(),seconds * time.Second)
	defer cancel()

	res,error := c.GreetWithDeadline(ctx,req)
	if error != nil{
		statusErr , ok := status.FromError(error)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			}else {
				fmt.Printf("Unexpected error: %v",statusErr)
			}
		}
		log.Fatalf("error while Calling Greet RPC %v",error)
	}
	log.Printf("Response from Greet: %v",res.Result)
}
