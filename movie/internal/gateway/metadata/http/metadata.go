package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
	"net/http"
)

// Gateway defines a movie metadata service gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	url := "http://" + addrs[rand.IntN(len(addrs))] + "/metadata"
	slog.Info("calling metadata service. request = GET", slog.String("url", url))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // ignoring error
	if res.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", res)
	}
	var metadata *model.Metadata
	err = json.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
