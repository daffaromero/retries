package controller

import (
	"io"

	pb "github.com/daffaromero/retries/services/common/genproto/orders"
	"google.golang.org/grpc"
)

type RestOrderServer struct {
	grpc.ServerStream
	results chan *pb.GetOrderResponse
}

func (x *RestOrderServer) Send(m *pb.GetOrderResponse) error {
	x.results <- m
	return nil
}

func NewRestOrderServer() *RestOrderServer {
	return &RestOrderServer{
		results: make(chan *pb.GetOrderResponse),
	}
}

func (x *RestOrderServer) Recv() (*pb.GetOrderResponse, error) {
	resp, ok := <-x.results
	if !ok {
		return nil, io.EOF
	}
	return resp, nil
}
