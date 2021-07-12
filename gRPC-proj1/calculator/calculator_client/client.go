package main

import (
	"context"
	"fmt"
	"gRPC-proj1/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	// "google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // connecting with server
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewResultServiceClient(cc) // creating request client

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	doErrorUnary(c)

}

func doErrorUnary(c calculatorpb.ResultServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC....")

	//correct call
	doErrorCall(c, 10)

	//error call
	doErrorCall(c, -8)

}

// custom grpc compatible error handling
func doErrorCall(c calculatorpb.ResultServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of SquareRoot of %v: %v\n", n, res.GetNumberRoot())

}

func doBiDiStreaming(c calculatorpb.ResultServiceClient) {
	fmt.Println("Starting to do a FindMaximum BiDi stream RPC......")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	//send go routine
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number.....: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum of.....: %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc

}

func doClientStreaming(c calculatorpb.ResultServiceClient) {
	fmt.Println("Starting to do a ComputerAverage Client Streaming RPC...")

	stream, err := c.ComputerAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opeing stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputerAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())

}

func doServerStreaming(c calculatorpb.ResultServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{ // creating request
		Number: 15,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF { // end loop when END OF FILE is reached
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doUnary(c calculatorpb.ResultServiceClient) {
	fmt.Println("Starting to do Unary RPC....")
	req := &calculatorpb.ResultRequest{ //creating request
		FirstNumber:  10,
		SecondNumber: 5,
	}

	res, err := c.Result(context.Background(), req) // sending and receiving response from server
	if err != nil {
		log.Fatalf("Error while calling Result RPC: %v", err)
	}

	log.Printf("Response from Server: %v", res)
}
