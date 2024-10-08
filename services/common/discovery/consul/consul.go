package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{client: client}, nil
}

func (r *Registry) Register(ctx context.Context, instanceId, serviceName, hostPort string, tags []string) error {
	host, portString, found := strings.Cut(hostPort, ":")
	if !found {
		return errors.New("invalid host:port format")
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return err
	}

	registration := &consul.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Port:    port,
		Address: host,
		Tags:    tags,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceId,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	err = r.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceId string, serviceName string) error {
	log.Printf("deregistering service %s", instanceId)
	return r.client.Agent().CheckDeregister(instanceId)
}

func (r *Registry) HealthCheck(instanceId string, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceId, "online", consul.HealthPassing)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]*consul.ServiceEntry, error) {
	instances, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	return instances, nil
}
