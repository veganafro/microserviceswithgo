package main

import (
	"log/slog"
	"movieexample.com/metadata/internal/controller/metadata"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"net/http"
)

func main() {
	slog.Info("starting metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	handler := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(handler.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
