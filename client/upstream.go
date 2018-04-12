package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/welaw/welaw/endpoints"
)

func decodeHTTPCreateUpstreamResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.CreateUpstreamResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func encodeHTTPGetUpstreamRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(endpoints.GetUpstreamRequest)
	values := req.URL.Query()
	values.Add("ident", r.Ident)
	req.URL.RawQuery = values.Encode()
	return nil
}

func decodeHTTPGetUpstreamResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.GetUpstreamResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPListUpstreamsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.ListUpstreamsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPUpdateUpstreamResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.UpdateUpstreamResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
