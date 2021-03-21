package main

import (
	"context"
	"ex_grpc/hellopb"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello client ...")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := hellopb.NewHelloServiceClient(cc)
	request := &hellopb.HelloRequest{Name: "Avinash"}

	resp, _ := client.Hello(context.Background(), request)
	fmt.Printf("Receive response => [%v]", resp.Greeting)
}
