package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/welaw/welaw/endpoints"
)

func encodeHTTPGetUpstreamRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(endpoints.GetUpstreamRequest)
	values := req.URL.Query()
	values.Add("ident", r.Ident)
	req.URL.RawQuery = values.Encode()
	return nil
}

func decodeHTTPGetUpstreamRequest(_ context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	ident, found := vars["ident"]
	if !found {
		return nil, fmt.Errorf("upstream not found")
	}
	return endpoints.GetUpstreamRequest{Ident: ident}, nil
}

func decodeHTTPGetUpstreamResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.GetUpstreamResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPListUpstreamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.ListUpstreamsRequest{}, nil
}

func decodeHTTPListUpstreamsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.ListUpstreamsResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPUpdateUpstreamRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateUpstreamRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUpdateUpstreamResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.UpdateUpstreamResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
