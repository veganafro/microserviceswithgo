package consul

import (
	"context"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"movieexample.com/pkg/discovery"
	"strconv"
	"strings"
)

// Registry defines a Consul-based service registry.
type Registry struct {
	client *consul.Client
}

// NewRegistry creates a new Consul-based service registry instance.
func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client}, nil
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8081")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	err = r.client.Agent().ServiceRegister(
		&consul.AgentServiceRegistration{
			Address: parts[0],
			Port:    port,
			Name:    serviceName,
			ID:      instanceID,
			Check: &consul.AgentServiceCheck{
				CheckID: instanceID,
				TTL:     "5s",
			},
		},
	)
	return err
}

// Deregister removes a service instance record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	err := r.client.Agent().ServiceDeregister(instanceID)
	return err
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, entry := range entries {
		res = append(res, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return res, nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(ctx context.Context, instanceID string, serviceName string) error {
	err := r.client.Agent().PassTTL(instanceID, "")
	return err
}
