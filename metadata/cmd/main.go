package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"movieexample.com/metadata/internal/controller/metadata"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"net/http"
	"time"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	slog.Info("starting metadata service")
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
	repo := memory.New()
	ctrl := metadata.New(repo)
	handler := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(handler.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
