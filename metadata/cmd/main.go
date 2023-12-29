package main

import (
	"log"
	"net/http"

	"github.com/supersherm5/movie-app/metadata/internal/controller/metadata"
	httphandler "github.com/supersherm5/movie-app/metadata/internal/handler/http"
	"github.com/supersherm5/movie-app/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting metadata service...")
	repo := memory.New()
	ctlr := metadata.New(repo)
	h := httphandler.New(ctlr)
	http.HandleFunc("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
