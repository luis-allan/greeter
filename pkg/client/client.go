package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/thingful/greeter/pkg/greeter"
)

const (
	address = "localhost:5555"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "foobar"})
	if err != nil {
		log.Fatalf("Could not get thing: %v", err)
	}

	log.Println(r)
}
