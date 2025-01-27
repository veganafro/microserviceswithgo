package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
	"movieexample.com/rating/pkg/model"
	"net/http"
)

// Gateway defines a rating service gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// GetAggregatedRating returns the aggregated rating for a record or gateway.ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, err
	}
	url := "http://" + addrs[rand.IntN(len(addrs))] + "/rating"
	slog.Info("calling rating service. request = GET", slog.String("url", url))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	req.URL.RawQuery = values.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close() // ignoring error
	if res.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response: %v", res)
	}
	var rating float64
	err = json.NewDecoder(res.Body).Decode(&rating)
	if err != nil {
		return 0, err
	}
	return rating, nil
}

// PutRating writes a rating.
func (g *Gateway) PutRating(
	ctx context.Context,
	recordID model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}
	url := "http://" + addrs[rand.IntN(len(addrs))] + "/rating"
	slog.Info("calling rating service. request = PUT", slog.String("url", url))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	values.Add("userId", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close() // ignoring error
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", res)
	}
	return nil
}
