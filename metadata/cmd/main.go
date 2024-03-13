package main

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

const serviceName = "metadata"

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
