package main

import (
	"context"
	v1 "github.com/luohaocode/Go-000/Week04/api/demo/v1"
	"google.golang.org/grpc"

	"log"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := v1.NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(),&v1.HelloRequest{
		Name: "hello, golang",
	})
	if err != nil {
		log.Fatalf("client.HelloService err %v", err)
	}

	log.Printf("resp: %s", resp.GetMessage())
}