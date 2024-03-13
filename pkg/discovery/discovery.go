package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry defines a service registry.
type Registry interface {
	// Register creates a service instance record in the registry.
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error

	// Deregister removes a service instance record from the registry.
	Deregister(ctx context.Context, instanceID, serviceName string) error

	// ServiceAddresses returns the list of addresses of active instances of the given service.
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)

	// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
	ReportHealthyState(instanceID, serviceName string) error
}

// ErrNotFound is returned when no service addresses are found.
var ErrNotFound = errors.New("no service addresses found")

// ErrServiceNotRegistered is returned when a service is not registered.
var ErrServiceNotRegistered = errors.New("service not registered")

// ErrInstanceNotRegistered is returned when an instance is not registered.
var ErrInstanceNotRegistered = errors.New("service instance not registered")

// ErrHostPortFormattedIncorrectly is returned when a host:port is not formatted correctly.
var ErrHostPortFormattedIncorrectly = errors.New("hostPort must formatted as <host>:<port>, for example: localhost:8080")

// GenerateInstanceID generates a pseudo-random service instance identifier,
// using a service name suffixed by dash and a random number.
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
