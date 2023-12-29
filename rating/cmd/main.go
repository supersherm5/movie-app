package main

import (
	ratingCntrl "github.com/supersherm5/movie-app/rating/internal/controller/rating"
	ratingHandler "github.com/supersherm5/movie-app/rating/internal/handler/http"
	"github.com/supersherm5/movie-app/rating/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting rating service...")
	repo := memory.New()
	ctrl := ratingCntrl.New(repo)
	h := ratingHandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
