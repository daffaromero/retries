package controller

import (
	"io"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"google.golang.org/grpc"
)

type RestCategoryServer struct {
	grpc.ServerStream
	results chan *pb.GetCategoryResponse
}

func (x *RestCategoryServer) Send(m *pb.GetCategoryResponse) error {
	x.results <- m
	return nil
}

func NewRestCategoryServer() *RestCategoryServer {
	return &RestCategoryServer{
		results: make(chan *pb.GetCategoryResponse),
	}
}

func (x *RestCategoryServer) Recv() (*pb.GetCategoryResponse, error) {
	resp, ok := <-x.results
	if !ok {
		return nil, io.EOF
	}
	return resp, nil
}
