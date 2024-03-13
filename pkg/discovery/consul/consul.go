package consul

import (
	"context"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"strconv"
	"strings"
)

// Registry defines a Consul-based service registry.
type Registry struct {
	client *consul.Client
}

// New creates a new Consul-based service registry instance.
func New(addr string) (*Registry, error) {
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
		return discovery.ErrHostPortFormattedIncorrectly
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TTL:     "5s",
		},
	})
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	return r.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	var addrs []string
	for _, entry := range entries {
		addr := fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "", "pass")
}
