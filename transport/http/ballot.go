package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/welaw/welaw/endpoints"
)

func decodeHTTPCreateVoteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateVoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateVotesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateVotesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCreateVoteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.CreateVoteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPDeleteVoteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteVoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPDeleteVoteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.DeleteVoteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPGetVoteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetVoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetVoteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.GetVoteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPUpdateVoteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateVoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUpdateVoteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.UpdateVoteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPListVotesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListVotesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPListVotesResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp endpoints.ListVotesResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
