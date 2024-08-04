package controller

import "github.com/daffaromero/retries/services/common/discovery"

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}
