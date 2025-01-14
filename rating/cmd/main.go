package main

import (
	"log/slog"
	"movieexample.com/rating/internal/controller/rating"
	httphandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
	"net/http"
)

func main() {
	slog.Info("starting rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	handler := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(handler.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
