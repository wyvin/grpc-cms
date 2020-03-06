package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc-content/proto"
	"log"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("../certs/cert.pem", "localhost")
	if err != nil {
		log.Printf("Failed to create TLS credentials %v", err)
	}
	conn, err := grpc.Dial(":50052", grpc.WithTransportCredentials(creds))
	defer conn.Close()

	if err != nil {
		log.Println(err)
	}

	//c := pb.NewHelloWorldClient(conn)
	//ctx := context.Background()
	//body := &pb.HelloWorldRequest{
	//}

	//r, err := c.SayHelloWorld(ctx, body)
	//if err != nil {
	//	log.Println(err)
	//}

	c := pb.NewAdClient(conn)
	ctx := context.Background()

	//body := &pb.AddAdRequest{
	//	Name:        "testName",
	//	Title:       "testTitle",
	//	Description: "desc",
	//	Remark:      "remark",
	//	Cover:       "cover",
	//	Url:         "url",
	//	Priority:    1,
	//	Display:     1,
	//}
	//
	//rep, err := c.AddAd(ctx, body)
	body := &pb.GetAdPlacementListRequest{
		Page:    0,
		PerPage: 0,
	}
	rep, err := c.GetAdPlacementList(ctx, body)
	log.Println(rep)
	log.Println(err)
}
