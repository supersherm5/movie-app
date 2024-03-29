package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/supersherm5/movie-app/gen"
	"github.com/supersherm5/movie-app/metadata/internal/controller/metadata"
	grpcHandler "github.com/supersherm5/movie-app/metadata/internal/handler/grpc"
	"github.com/supersherm5/movie-app/metadata/internal/repository/memory"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"github.com/supersherm5/movie-app/pkg/discovery/consul"
	"github.com/supersherm5/movie-app/test_data"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "Metadata API handler port")
	flag.Parse()
	log.Printf("Starting metadata service on port %d...", port)
	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("localhost:%d", port)
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		log.Fatalf("Failed to register %v service: %v", serviceName, err)
	}
	repo := memory.NewWithMetadata(test_data.FakeMetaData) // Add this line to use fake data
	ctrl := metadata.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	gen.RegisterMetadataServiceServer(server, h)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("%v service failed to serve: %v", serviceName, err)
		}
	}()

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Printf("Error reporting healthy state for the %s at instance, %s: %v", serviceName, instanceID, err)
			}
			<-time.After(1 * time.Second)
		}
	}()

	defer func (registry *consul.Registry, ctx context.Context, instanceID string, serviceName string)  {
		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Printf("Error deregistering service, %s: %v", serviceName, err)
		}
		log.Println("Deregistration of metadata service is complete.")
	}(registry, ctx, instanceID, serviceName)



	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan
	log.Println("Shutting down metadata service...")
	server.GracefulStop()
	log.Println("Metadata service stopped serving requests.")
}

// Main function for non-grpc server
/*
import (
	"context"
	"flag"
	"fmt"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"github.com/supersherm5/movie-app/pkg/discovery/consul"
	"log"
	"net/http"
	"time"

	"github.com/supersherm5/movie-app/metadata/internal/controller/metadata"
	httphandler "github.com/supersherm5/movie-app/metadata/internal/handler/http"
	"github.com/supersherm5/movie-app/metadata/internal/repository/memory"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "Metadata API handler port")
	flag.Parse()
	log.Println("Starting metadata service...")
	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}

	hostPort := fmt.Sprintf("localhost:%d", port)
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Printf("Error reporting healthy state for the %s at instance, %s: %v", serviceName, instanceID, err)
			}
			<-time.After(1 * time.Second)
		}
	}()

	defer func(registry *consul.Registry, ctx context.Context, instanceID string, serviceName string) {
		err := registry.Deregister(ctx, instanceID, serviceName)
		if err != nil {
			log.Printf("Error deregistering service, %s: %v", serviceName, err)
		}
	}(registry, ctx, instanceID, serviceName)

	repo := memory.New()
	ctlr := metadata.New(repo)
	h := httphandler.New(ctlr)
	http.HandleFunc("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
*/
