package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/supersherm5/movie-app/gen"
	"github.com/supersherm5/movie-app/movie/internal/controller/movie"
	metadatagateway "github.com/supersherm5/movie-app/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/supersherm5/movie-app/movie/internal/gateway/rating/grpc"
	grpchandler "github.com/supersherm5/movie-app/movie/internal/handler/grpc"
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
		log.Fatalf("Could not start movie service: %v", err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	hostPort := "localhost:8083"
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		log.Fatalf("Could not register movie service: %v", err)
	}

	
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("Could not start movie service: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Could not start movie service: %v", err)
		}
	}()

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Fatalf("Could not report healthy state for movie service: %v", err)
			}
			<- time.After(1 * time.Second)
		}
	}()


defer func(registry *consul.Registry, ctx context.Context, instanceID string, serviceName string) {
	if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
		log.Fatalf("Could not deregister movie service: %v", err)
	}
	log.Println("Movie service deregistered")
}(registry, ctx, instanceID, serviceName)

signalChan := make(chan os.Signal, 1)
signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
<-signalChan

log.Println("Shutting down movie service...")
srv.GracefulStop()
log.Println("Movie service stopped serving requests.")

}

//func main() {
//	var port int
//	flag.IntVar(&port, "port", 8083, "Movie API handler port")
//	flag.Parse()
//
//	log.Println("Starting the movie service....")
//	registry, err := consul.New("localhost:8500")
//	if err != nil {
//		panic(err)
//	}
//	addrPort := fmt.Sprintf(":%d", port)
//	server := &http.Server{Addr: addrPort}
//	hostPort := fmt.Sprintf("localhost:%d", port)
//	ctx := context.Background()
//	instanceID := discovery.GenerateInstanceID(serviceName)
//
//	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
//		panic(err)
//	}
//	go func() {
//		for {
//			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
//				log.Printf("Error reporting healthy state for the %s at instance, %s: %v", serviceName, instanceID, err)
//			}
//			<-time.After(1 * time.Second)
//		}
//	}()
//
//	go func() {
//		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
//			log.Fatalf("Could not start movie service: %v", err)
//		}
//		log.Println("Movie service server stopped serving new requests")
//	}()
//
//	defer func(registry *consul.Registry, ctx context.Context, instanceID string, serviceName string) {
//		err := registry.Deregister(ctx, instanceID, serviceName)
//		if err != nil {
//			log.Printf("Error deregistering movie service, %s: %v", serviceName, err)
//		}
//	}(registry, ctx, instanceID, serviceName)
//
//	// Instantiate the metadata and rating gateways.
//	metadataGateway := MetadataGateway.New(registry)
//	ratingGateway := RatingGateway.New(registry)
//
//	// Assign the gateways to the controller.
//	ctrl := movie.New(ratingGateway, metadataGateway)
//
//	// Instantiate the http handler and pass the controller to it.
//	h := HttpHandler.New(ctrl)
//	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
//
//	signalChan := make(chan os.Signal, 1)
//	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
//	<-signalChan
//
//	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
//	defer shutdownRelease()
//
//	if err := server.Shutdown(shutdownCtx); err != nil {
//		log.Printf("Error shutting down the server: %v", err)
//	}
//}
