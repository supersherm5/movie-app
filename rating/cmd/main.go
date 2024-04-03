package main

import (
	"context"
	"errors"
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
	"github.com/supersherm5/movie-app/pkg/discovery"
	"github.com/supersherm5/movie-app/pkg/discovery/consul"
	ratingCntrl "github.com/supersherm5/movie-app/rating/internal/controller/rating"
	ratingHandler "github.com/supersherm5/movie-app/rating/internal/handler/grpc"
	"github.com/supersherm5/movie-app/rating/internal/repository/memory"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "Rating API handler port")
	flag.Parse()
	log.Printf("Starting rating service on port %d...", port)
	registry, err := consul.New("localhost:8500")
	if err != nil {
		log.Fatalf("Failed to load consul for %v service: %v", serviceName, err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := fmt.Sprintf("localhost:%d", port)
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		log.Fatalf("Failed to register %v service: %v", serviceName, err)
	}
	repo := memory.New()
	ctrl := ratingCntrl.New(repo, nil)
	h := ratingHandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("Rating service failed to listen: %v", err)
	}
	server := grpc.NewServer()
	gen.RegisterRatingServiceServer(server, h)

	go func() {
		if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("Rating service failed to serve: %v", err)
		}
	}()

	go func ()  {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Fatalf("Failed to report healthy state for %v service: %v", serviceName, err)
			}
			<- time.After(1 * time.Second)
		}
	}()

	defer func (registry *consul.Registry, ctx context.Context, instanceID string, serviceName string)  {
		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Fatalf("Failed to deregister %v service: %v", serviceName, err)
		}
		log.Println("Rating service deregistered")
	}(registry, ctx, instanceID, serviceName)


	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan
	
	log.Println("Rating service shutting down...")
	server.GracefulStop()
	log.Println("Rating service shut down complete.")

}

/*

	var port int
	flag.IntVar(&port, "port", 8082, "Rating API handler port")
	flag.Parse()
	log.Println("Starting rating service...")

	// Create a new Consul-based service registry.
	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}
	addrPort := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr: addrPort,
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

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not start rating service: %v", err)
		}
		log.Println("Rating service server stopped serving new requests")
	}()

	defer func(registry *consul.Registry, ctx context.Context, instanceID string, serviceName string) {
		err := registry.Deregister(ctx, instanceID, serviceName)
		if err != nil {
			log.Printf("Error unregistering rating service, %s: %v", serviceName, err)
		}
	}(registry, ctx, instanceID, serviceName)

	repo := memory.New()
	ctrl := ratingCtrl.New(repo)
	h := ratingHandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down the server: %v", err)
	}
	log.Println("Rating service server shut down complete.")
*/
