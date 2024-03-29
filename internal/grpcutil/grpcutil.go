package grpcutil

import (
	"context"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
)

// ServiceConnection attempts to select a random service
// instance and returns a gRPC connections to it.
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	instances, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	instance := instances[rand.Intn(len(instances))]
	conn, err := grpc.Dial(instance, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
