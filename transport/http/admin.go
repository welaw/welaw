package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/welaw/welaw/endpoints"
)

func decodeHTTPGetServerStatsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.GetServerStatsRequest{}, nil
}

func decodeHTTPLoadReposRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.LoadReposRequest{}, nil
}

func decodeHTTPSaveReposRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.SaveReposRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
