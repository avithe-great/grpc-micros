package main

import (
	"context"
	"ex_grpc/hellopb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	hellopb.UnimplementedHelloServiceServer
}

func (*server) Hello(ctx context.Context, request *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	name := request.Name
	response := &hellopb.HelloResponse{
		Greeting: "Hello " + name,
	}
	return response, nil
}

func main() {
	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	fmt.Printf("Server is listening on %v ...", address)

	s := grpc.NewServer()
	hellopb.RegisterHelloServiceServer(s, &server{})

	s.Serve(lis)
}
