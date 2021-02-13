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
	"net"
	"time"
)

type server struct {

}

func (s *server) GreetWithDeadline(ctx context.Context, request *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline Funstion is invoked with %v",request)
	firstName := request.GetGreeting().GetFirstName()
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The Client canceled the request")
			return nil, status.Error(codes.Canceled,"The Client canceled the request")
		}
		time.Sleep(1 & time.Second)
	}
	result := "Hello" + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res , nil
}

func (s *server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryOne function was invoked with a Streaming request")

	for  {
		req,err := everyoneServer.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v",err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "hello" + firstName + "!"
		errSend := everyoneServer.Send(&greetpb.GreetEveryoneResponse{Result: result})
		if errSend != nil {
			log.Fatalf("Error while sending data to client : %v",err)
			return err
		}
	}
}

func (s *server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet Server %v",greetServer)
	result := ""
	for {
		req,err := greetServer.Recv()
		if err == io.EOF {
			//We have finished reading the client stream
			greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			break
		}
		if err != nil {
			log.Fatalf("Error while reading %v",err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello" + firstName + "!"
	}
	return nil
}

func (s *server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	firstName := request.GetGreeting().GetFirstName()
	for i := 0 ; i < 10; i++ {
		result := "hello" + firstName + "number" + string(i)
		res := &greetpb.GreetManyTimesResponse{Result: result}
		timesServer.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) Greet(ctx context.Context,req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	fmt.Printf("Greet Funstion is invoked with %v",req)
   firstName := req.GetGreeting().GetFirstName()

   result := "Hello" + firstName
   res := &greetpb.GreetResponse{
   	     Result: result,
   }
   return res , nil
}


func main() {
	fmt.Println("Hello World")

	lis , err := net.Listen("tcp","0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v",err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s,&server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v",err)
	}
}
