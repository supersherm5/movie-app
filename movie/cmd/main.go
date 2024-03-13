package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supersherm5/movie-app/movie/internal/controller/movie"
	MetadataGateway "github.com/supersherm5/movie-app/movie/internal/gateway/metadata/http"
	RatingGateway "github.com/supersherm5/movie-app/movie/internal/gateway/rating/http"
	HttpHandler "github.com/supersherm5/movie-app/movie/internal/handler/http"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"github.com/supersherm5/movie-app/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "Movie API handler port")
	flag.Parse()

	log.Println("Starting the movie service....")
	registry, err := consul.New("localhost:8500")
	if err != nil {
		panic(err)
	}
	addrPort := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: addrPort}
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
			log.Fatalf("Could not start movie service: %v", err)
		}
		log.Println("Movie service server stopped serving new requests")
	}()

	defer func(registry *consul.Registry, ctx context.Context, instanceID string, serviceName string) {
		err := registry.Deregister(ctx, instanceID, serviceName)
		if err != nil {
			log.Printf("Error deregistering movie service, %s: %v", serviceName, err)
		}
	}(registry, ctx, instanceID, serviceName)

	// Instantiate the metadata and rating gateways.
	metadataGateway := MetadataGateway.New(registry)
	ratingGateway := RatingGateway.New(registry)

	// Assign the gateways to the controller.
	ctrl := movie.New(ratingGateway, metadataGateway)

	// Instantiate the http handler and pass the controller to it.
	h := HttpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down the server: %v", err)
	}
}
