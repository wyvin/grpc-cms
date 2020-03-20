package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "grpc-cms/proto"
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


	// 广告模块测试
	c := pb.NewAdClient(conn)
	ctx := context.Background()

	// 添加广告
	body := &pb.AddAdRequest{
		AppId:       "jx12345678",
		GroupId:     999,
		Name:        "testName",
		Title:       "testTitle",
		Description: "testDesc",
		Remark:      "testRemark",
		Cover:       "testCover",
		Url:         "testUrl",
		Priority:    99,
		Display:     1,
		State:       2,
	}
	rep, err := c.AddAd(ctx, body)
	log.Println(rep)

	//body := &pb.GetAdPlacementListRequest{
	//	Page:    0,
	//	PerPage: 0,
	//}
	//rep, err := c.GetAdPlacementList(ctx, body)
	//log.Println(rep)
	//log.Println(err)
}
