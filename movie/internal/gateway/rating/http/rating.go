package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/supersherm5/movie-app/movie/internal/gateway"
	"github.com/supersherm5/movie-app/pkg/discovery"
	"github.com/supersherm5/movie-app/rating/pkg/model"
	"math/rand"
	"net/http"
)

// Gateway defines an HTTP gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry,
	}
}

// GetAggregatedRating gets the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating/"
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return 0, gateway.ErrServiceNotReachable
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	req.URL.RawQuery = values.Encode()
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return 0, respErr
	}
	defer func(res *http.Response) {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp)

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return 0, gateway.ErrNotFound

	case resp.StatusCode/100 != 2:
		return 0, fmt.Errorf("non-200 status code: %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// PutRating writes a rating
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating/"
	req, reqErr := http.NewRequest(http.MethodPut, url, nil)
	if reqErr != nil {
		return gateway.ErrServiceNotReachable
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))
	values.Add("userID", string(rating.UserID))
	values.Add("value", fmt.Sprintf("%v", rating.Value))
	req.URL.RawQuery = values.Encode()

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return respErr
	}
	defer func(res *http.Response) {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp)

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return gateway.ErrNotFound

	case resp.StatusCode/100 != 2:
		return fmt.Errorf("non-200 status code: %v", resp)
	}

	return nil
}
