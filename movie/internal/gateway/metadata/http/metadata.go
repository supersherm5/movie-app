package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/supersherm5/movie-app/metadata/pkg/model"
	"github.com/supersherm5/movie-app/movie/internal/gateway"
	"net/http"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	addr string
}

// New creates a new HTTP gateway for a movie metadata service.
func New(addr string) *Gateway {
	return &Gateway{addr}
}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	MetadataPath := fmt.Sprintf("%s/metadata", g.addr)
	req, reqErr := http.NewRequest(http.MethodGet, MetadataPath, nil)
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
