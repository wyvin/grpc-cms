package server

import (
	"context"
	pb "grpc-content/proto"
)

type helloService struct{}

func NewHelloService() *helloService {
	return &helloService{}
}

func (h helloService) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.ApiResponse, error) {
	return &pb.ApiResponse{
		Code: 0,
		Msg:  "success",
	}, nil
}
