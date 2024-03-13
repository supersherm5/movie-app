package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/supersherm5/movie-app/metadata/pkg/model"
	"github.com/supersherm5/movie-app/movie/internal/gateway"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"log"
	"math/rand"
	"net/http"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata/"
	log.Printf("Calling metadata service. Request GET %s", url)

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, gateway.ErrServiceNotReachable
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return nil, respErr
	}

	defer func(res *http.Response) {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp)

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, gateway.ErrNotFound

	case resp.StatusCode/100 != 2:
		return nil, fmt.Errorf("non-200 status code: %v", resp)
	}

	var metadata *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}
