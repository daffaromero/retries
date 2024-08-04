package memory

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Registry struct {
	sync.RWMutex
	addrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{addrs: map[string]map[string]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, instanceId, serviceName, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.addrs[serviceName]; !ok {
		r.addrs[serviceName] = map[string]*serviceInstance{}
	}

	r.addrs[serviceName][instanceId] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceId, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.addrs[serviceName]; !ok {
		return nil
	}

	delete(r.addrs[serviceName], instanceId)
	return nil
}

func (r *Registry) HealthCheck(instanceId, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.addrs[serviceName]; !ok {
		return errors.New("service not registered")
	}

	if _, ok := r.addrs[serviceName][instanceId]; !ok {
		return errors.New("service instance not registered")
	}

	r.addrs[serviceName][instanceId].lastActive = time.Now()
	return nil
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.addrs[serviceName]) == 0 {
		return nil, errors.New("no service address found")
	}

	var res []string
	for _, i := range r.addrs[serviceName] {
		res = append(res, i.hostPort)
	}

	return res, nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.addrs[serviceName]) == 0 {
		return nil, errors.New("no service address found")
	}

	var res []string
	for _, i := range r.addrs[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}

	return res, nil
}
