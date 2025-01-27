package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/http"
	httphandler "movieexample.com/movie/internal/handler/http"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"net/http"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	slog.Info("starting movie service")
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(context.Background(), instanceID, serviceName); err != nil {
				slog.Error("failed to report healthy state", slog.String("err", err.Error()))
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName) // ignoring error

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(metadataGateway, ratingGateway)
	handler := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(handler.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
