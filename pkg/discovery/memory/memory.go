package memory

import (
	"context"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"sync"
	"time"
)

// ServiceInstance defines a service instance.
type ServiceInstance struct {
	hostPort   string
	lastActive time.Time
}

// Registry defines an in-memory service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*ServiceInstance // map[serviceName]map[instanceID]*ServiceInstance
}

// New creates a new in-memory service registry instance.
func New() *Registry {
	return &Registry{serviceAddrs: make(map[string]map[string]*ServiceInstance)}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = make(map[string]*ServiceInstance)
	}

	r.serviceAddrs[serviceName][instanceID] = &ServiceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(Ctx context.Context, instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return discovery.ErrServiceNotRegistered
	}

	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return discovery.ErrInstanceNotRegistered
	}

	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var addrs []string
	for _, inst := range r.serviceAddrs[serviceName] {
		if inst.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		addrs = append(addrs, inst.hostPort)
	}
	return addrs, nil
}
