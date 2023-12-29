package main

import (
	"github.com/supersherm5/movie-app/movie/internal/controller/movie"
	MetadataGateway "github.com/supersherm5/movie-app/movie/internal/gateway/metadata/http"
	RatingGateway "github.com/supersherm5/movie-app/movie/internal/gateway/rating/http"
	HttpHandler "github.com/supersherm5/movie-app/movie/internal/handler/http"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting the movie service....")

	// Instantiate the metadata and rating gateways.
	metadataGateway := MetadataGateway.New("http://localhost:8081")
	ratingGateway := RatingGateway.New("http://localhost:8082")

	// Assign the gateways to the controller.
	ctrl := movie.New(ratingGateway, metadataGateway)

	// Instantiate the http handler and pass the controller to it.
	h := HttpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}

}
